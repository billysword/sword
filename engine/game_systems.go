package engine

import (
	"fmt"
	"sword/entities"
	"sword/world"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GameSystem represents a modular game system that can be updated
type GameSystem interface {
	Update() error
	GetName() string
	IsEnabled() bool
	SetEnabled(bool)
}

// BaseSystem provides common functionality for all systems
type BaseSystem struct {
	name    string
	enabled bool
}

// GetName returns the system name
func (bs *BaseSystem) GetName() string {
	return bs.name
}

// IsEnabled returns whether the system is enabled
func (bs *BaseSystem) IsEnabled() bool {
	return bs.enabled
}

// SetEnabled sets the system enabled state
func (bs *BaseSystem) SetEnabled(enabled bool) {
	bs.enabled = enabled
}

// InputSystem handles input processing
type InputSystem struct {
	BaseSystem
	player           *entities.Player
	stateManager     *StateManager
	roomTransitionMgr *world.RoomTransitionManager
	lastActionPressed bool
}

// NewInputSystem creates a new input system
func NewInputSystem(player *entities.Player, stateManager *StateManager, roomTransitionMgr *world.RoomTransitionManager) *InputSystem {
	return &InputSystem{
		BaseSystem:       BaseSystem{name: "Input", enabled: true},
		player:           player,
		stateManager:     stateManager,
		roomTransitionMgr: roomTransitionMgr,
	}
}

// Update processes input for the current frame
func (is *InputSystem) Update() error {
	if !is.enabled {
		return nil
	}

	roomName := "Unknown"
	if is.roomTransitionMgr != nil {
		roomName = is.roomTransitionMgr.GetCurrentRoomID()
	}

	playerX, playerY := is.player.GetPosition()

	// State transition inputs
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		LogPlayerInput("Escape (Pause)", playerX, playerY, roomName)
		// Note: This needs access to the current state, will be handled by InGameState
		return fmt.Errorf("pause_requested")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		LogPlayerInput("Tab (Settings)", playerX, playerY, roomName)
		// Note: This needs access to the current room, will be handled by InGameState
		return fmt.Errorf("settings_requested")
	}

	// Debug toggle inputs
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)
	
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		LogPlayerInput("B (Toggle Background)", playerX, playerY, roomName)
		ToggleBackground()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		LogPlayerInput("G (Toggle Grid)", playerX, playerY, roomName)
		ToggleGrid()
	}

	// Character scale adjustments
	if inpututil.IsKeyJustPressed(ebiten.KeyComma) {
		if shiftPressed {
			GameConfig.CharScaleFactor -= 0.05
		} else {
			GameConfig.CharScaleFactor -= 0.1
		}
		if GameConfig.CharScaleFactor < 0.5 {
			GameConfig.CharScaleFactor = 0.5
		}
		LogInfo(fmt.Sprintf("Character scale decreased to: %.2f", GameConfig.CharScaleFactor))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		if shiftPressed {
			GameConfig.CharScaleFactor += 0.05
		} else {
			GameConfig.CharScaleFactor += 0.1
		}
		if GameConfig.CharScaleFactor > 3.0 {
			GameConfig.CharScaleFactor = 3.0
		}
		LogInfo(fmt.Sprintf("Character scale increased to: %.2f", GameConfig.CharScaleFactor))
	}

	// Tile scale adjustments
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		if shiftPressed {
			GameConfig.TileScaleFactor -= 0.1
		} else {
			GameConfig.TileScaleFactor -= 0.5
		}
		if GameConfig.TileScaleFactor < 0.5 {
			GameConfig.TileScaleFactor = 0.5
		}
		LogInfo(fmt.Sprintf("Tile scale decreased to: %.1f", GameConfig.TileScaleFactor))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		if shiftPressed {
			GameConfig.TileScaleFactor += 0.1
		} else {
			GameConfig.TileScaleFactor += 0.5
		}
		if GameConfig.TileScaleFactor > 4.0 {
			GameConfig.TileScaleFactor = 4.0
		}
		LogInfo(fmt.Sprintf("Tile scale increased to: %.1f", GameConfig.TileScaleFactor))
	}

	// Room transition input
	actionPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if is.roomTransitionMgr != nil {
		is.roomTransitionMgr.CheckTransitions(is.player, actionPressed)
	}

	// Handle player input
	is.player.HandleInputWithLogging(roomName)
	
	return nil
}

// PhysicsSystem handles physics updates
type PhysicsSystem struct {
	BaseSystem
	player   *entities.Player
	enemies  []entities.Enemy
	room     world.Room
}

// NewPhysicsSystem creates a new physics system
func NewPhysicsSystem(player *entities.Player) *PhysicsSystem {
	return &PhysicsSystem{
		BaseSystem: BaseSystem{name: "Physics", enabled: true},
		player:     player,
		enemies:    make([]entities.Enemy, 0),
	}
}

// SetCurrentRoom sets the room for collision detection
func (ps *PhysicsSystem) SetCurrentRoom(room world.Room) {
	ps.room = room
}

// AddEnemy adds an enemy to the physics simulation
func (ps *PhysicsSystem) AddEnemy(enemy entities.Enemy) {
	ps.enemies = append(ps.enemies, enemy)
}

// ClearEnemies removes all enemies
func (ps *PhysicsSystem) ClearEnemies() {
	ps.enemies = ps.enemies[:0]
}

// Update processes physics for the current frame
func (ps *PhysicsSystem) Update() error {
	if !ps.enabled {
		return nil
	}

	// Update player physics
	ps.player.Update()

	// Update enemy physics
	for _, enemy := range ps.enemies {
		if err := enemy.Update(); err != nil {
			return fmt.Errorf("enemy update failed: %w", err)
		}
	}

	// Handle collisions with room
	if ps.room != nil {
		ps.room.HandleCollisions(ps.player)
		
		// Handle enemy collisions too
		for _, enemy := range ps.enemies {
			if collisionEntity, ok := enemy.(entities.CollisionHandler); ok {
				collisionEntity.HandleCollisions()
			}
		}
	}

	return nil
}

// CameraSystem handles camera updates
type CameraSystem struct {
	BaseSystem
	camera   *Camera
	player   *entities.Player
	room     world.Room
}

// NewCameraSystem creates a new camera system
func NewCameraSystem(camera *Camera, player *entities.Player) *CameraSystem {
	return &CameraSystem{
		BaseSystem: BaseSystem{name: "Camera", enabled: true},
		camera:     camera,
		player:     player,
	}
}

// SetCurrentRoom sets the room for camera bounds
func (cs *CameraSystem) SetCurrentRoom(room world.Room) {
	cs.room = room
	if cs.camera != nil && room != nil {
		if tileMap := room.GetTileMap(); tileMap != nil {
			physicsUnit := GetPhysicsUnit()
			worldWidth := tileMap.Width * physicsUnit
			worldHeight := tileMap.Height * physicsUnit
			cs.camera.SetWorldBounds(worldWidth, worldHeight)
		}
	}
}

// Update processes camera movement for the current frame
func (cs *CameraSystem) Update() error {
	if !cs.enabled || cs.camera == nil || cs.player == nil {
		return nil
	}

	// Update camera to follow player
	playerX, playerY := cs.player.GetPosition()
	cs.camera.Update(playerX, playerY)

	return nil
}

// GetCamera returns the camera instance
func (cs *CameraSystem) GetCamera() *Camera {
	return cs.camera
}

// RoomSystem handles room-specific updates and transitions
type RoomSystem struct {
	BaseSystem
	transitionManager *world.RoomTransitionManager
	worldMap          *world.WorldMap
	player            *entities.Player
	lastRoomID        string
}

// NewRoomSystem creates a new room system
func NewRoomSystem(transitionManager *world.RoomTransitionManager, worldMap *world.WorldMap, player *entities.Player) *RoomSystem {
	return &RoomSystem{
		BaseSystem:        BaseSystem{name: "Room", enabled: true},
		transitionManager: transitionManager,
		worldMap:          worldMap,
		player:            player,
	}
}

// Update processes room updates and transitions
func (rs *RoomSystem) Update() error {
	if !rs.enabled {
		return nil
	}

	currentRoom := rs.transitionManager.GetCurrentRoom()
	if currentRoom == nil {
		return nil
	}

	// Update room-specific logic
	if err := currentRoom.Update(rs.player); err != nil {
		return fmt.Errorf("room update failed: %w", err)
	}

	// Handle room transitions
	if rs.transitionManager.HasPendingTransition() {
		newRoom, err := rs.transitionManager.ProcessPendingTransition(rs.player)
		if err != nil {
			return fmt.Errorf("room transition failed: %w", err)
		}
		if newRoom != nil {
			// Room changed successfully
			LogInfo(fmt.Sprintf("Room transitioned to: %s", newRoom.GetZoneID()))
		}
	}

	// Update world map with player position and room changes
	playerX, playerY := rs.player.GetPosition()
	rs.worldMap.AddPlayerPosition(playerX, playerY)
	
	currentRoomID := rs.transitionManager.GetCurrentRoomID()
	if currentRoomID != rs.lastRoomID {
		rs.worldMap.DiscoverRoom(currentRoom)
		rs.worldMap.SetCurrentRoom(currentRoomID)
		rs.lastRoomID = currentRoomID
	}

	return nil
}

// GetCurrentRoom returns the current room
func (rs *RoomSystem) GetCurrentRoom() world.Room {
	return rs.transitionManager.GetCurrentRoom()
}

// GameSystemManager manages all game systems
type GameSystemManager struct {
	systems        []GameSystem
	updateOrder    []string // Order in which to update systems
	systemRegistry map[string]GameSystem
}

// NewGameSystemManager creates a new system manager
func NewGameSystemManager() *GameSystemManager {
	return &GameSystemManager{
		systems:        make([]GameSystem, 0),
		updateOrder:    []string{"Input", "Physics", "Camera", "Room"},
		systemRegistry: make(map[string]GameSystem),
	}
}

// RegisterSystem adds a system to the manager
func (gsm *GameSystemManager) RegisterSystem(system GameSystem) {
	gsm.systems = append(gsm.systems, system)
	gsm.systemRegistry[system.GetName()] = system
}

// GetSystem returns a system by name
func (gsm *GameSystemManager) GetSystem(name string) GameSystem {
	return gsm.systemRegistry[name]
}

// SetSystemEnabled enables or disables a system
func (gsm *GameSystemManager) SetSystemEnabled(name string, enabled bool) {
	if system, exists := gsm.systemRegistry[name]; exists {
		system.SetEnabled(enabled)
	}
}

// UpdateAll updates all systems in the specified order
func (gsm *GameSystemManager) UpdateAll() error {
	// Update systems in the specified order
	for _, systemName := range gsm.updateOrder {
		if system, exists := gsm.systemRegistry[systemName]; exists {
			if system.IsEnabled() {
				if err := system.Update(); err != nil {
					return fmt.Errorf("system %s update failed: %w", systemName, err)
				}
			}
		}
	}
	
	// Update any remaining systems not in the update order
	for _, system := range gsm.systems {
		found := false
		for _, orderedName := range gsm.updateOrder {
			if system.GetName() == orderedName {
				found = true
				break
			}
		}
		if !found && system.IsEnabled() {
			if err := system.Update(); err != nil {
				return fmt.Errorf("system %s update failed: %w", system.GetName(), err)
			}
		}
	}
	
	return nil
}

// SetUpdateOrder sets the order in which systems are updated
func (gsm *GameSystemManager) SetUpdateOrder(order []string) {
	gsm.updateOrder = order
}