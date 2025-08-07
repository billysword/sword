package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

/*
GameSystem defines the interface for all game systems.
Each system manages a specific aspect of the game (input, physics, camera, etc.)
*/
type GameSystem interface {
	GetName() string
	Update() error
}

/*
InputSystem handles all player input and translates it to game actions.
Manages keyboard input for movement, jumping, and other player actions.
*/
type InputSystem struct {
	player           *entities.Player
	keysPressed      map[ebiten.Key]bool
	roomTransitionMgr *world.RoomTransitionManager
	pauseRequested   bool
}

/*
NewInputSystem creates a new input system instance.
Parameters:
  - player: The player entity to control
  - roomTransitionMgr: Manager for handling room transitions
*/
func NewInputSystem(player *entities.Player, roomTransitionMgr *world.RoomTransitionManager) *InputSystem {
	return &InputSystem{
		player:           player,
		keysPressed:      make(map[ebiten.Key]bool),
		roomTransitionMgr: roomTransitionMgr,
	}
}

func (is *InputSystem) GetName() string {
	return "Input"
}

func (is *InputSystem) Update() error {
	// Handle room transitions first
	if is.roomTransitionMgr != nil {
		is.roomTransitionMgr.Update()
	}
	
	// Handle movement inputs
	config := engine.GameConfig.PlayerPhysics
	physicsUnit := engine.GetPhysicsUnit()
	
	// Left/Right movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		is.player.MoveLeft(config.MoveSpeed * physicsUnit)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		is.player.MoveRight(config.MoveSpeed * physicsUnit)
	}
	
	// Jump input
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		is.player.Jump()
	} else {
		is.player.ReleaseJump()
	}
	
	// Debug toggle inputs
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		if !is.keysPressed[ebiten.KeyF3] {
			engine.GameConfig.ShowDebugInfo = !engine.GameConfig.ShowDebugInfo
			is.keysPressed[ebiten.KeyF3] = true
		}
	} else {
		is.keysPressed[ebiten.KeyF3] = false
	}
	
	if ebiten.IsKeyPressed(ebiten.KeyF4) {
		if !is.keysPressed[ebiten.KeyF4] {
			engine.GameConfig.ShowDebugOverlay = !engine.GameConfig.ShowDebugOverlay
			is.keysPressed[ebiten.KeyF4] = true
		}
	} else {
		is.keysPressed[ebiten.KeyF4] = false
	}
	
	// Grid toggle with G key
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		if !is.keysPressed[ebiten.KeyG] {
			engine.ToggleGrid()
			is.keysPressed[ebiten.KeyG] = true
		}
	} else {
		is.keysPressed[ebiten.KeyG] = false
	}
	
	// Room transition debug keys
	if is.roomTransitionMgr != nil {
		// Quick room transitions for testing
		if ebiten.IsKeyPressed(ebiten.Key1) {
			if !is.keysPressed[ebiten.Key1] {
				is.roomTransitionMgr.StartTransition(world.DirectionNorth, 500)
				is.keysPressed[ebiten.Key1] = true
			}
		} else {
			is.keysPressed[ebiten.Key1] = false
		}
		
		if ebiten.IsKeyPressed(ebiten.Key2) {
			if !is.keysPressed[ebiten.Key2] {
				is.roomTransitionMgr.StartTransition(world.DirectionEast, 500)
				is.keysPressed[ebiten.Key2] = true
			}
		} else {
			is.keysPressed[ebiten.Key2] = false
		}
		
		if ebiten.IsKeyPressed(ebiten.Key3) {
			if !is.keysPressed[ebiten.Key3] {
				is.roomTransitionMgr.StartTransition(world.DirectionSouth, 500)
				is.keysPressed[ebiten.Key3] = true
			}
		} else {
			is.keysPressed[ebiten.Key3] = false
		}
		
		if ebiten.IsKeyPressed(ebiten.Key4) {
			if !is.keysPressed[ebiten.Key4] {
				is.roomTransitionMgr.StartTransition(world.DirectionWest, 500)
				is.keysPressed[ebiten.Key4] = true
			}
		} else {
			is.keysPressed[ebiten.Key4] = false
		}
		
		// R key to reload current room
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			if !is.keysPressed[ebiten.KeyR] {
				if currentRoom := is.roomTransitionMgr.GetCurrentRoom(); currentRoom != nil {
					is.roomTransitionMgr.ReloadCurrentRoom()
				}
				is.keysPressed[ebiten.KeyR] = true
			}
		} else {
			is.keysPressed[ebiten.KeyR] = false
		}
	}
	
	// Pause request handling
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if !is.keysPressed[ebiten.KeyEscape] {
			is.pauseRequested = true
			is.keysPressed[ebiten.KeyEscape] = true
		}
	} else {
		is.keysPressed[ebiten.KeyEscape] = false
	}
	
	return nil
}

func (is *InputSystem) HasPauseRequest() bool {
	return is.pauseRequested
}

func (is *InputSystem) ClearRequests() {
	is.pauseRequested = false
}

/*
PhysicsSystem handles physics simulation for all entities.
Manages gravity, collisions, and movement for the player and enemies.
*/
type PhysicsSystem struct {
	player   *entities.Player
	enemies  []entities.Enemy
	room     world.Room
}

/*
NewPhysicsSystem creates a new physics system instance.
*/
func NewPhysicsSystem(player *entities.Player) *PhysicsSystem {
	return &PhysicsSystem{
		player:     player,
		enemies:    make([]entities.Enemy, 0),
	}
}

func (ps *PhysicsSystem) GetName() string {
	return "Physics"
}

func (ps *PhysicsSystem) AddEnemy(enemy entities.Enemy) {
	ps.enemies = append(ps.enemies, enemy)
}

func (ps *PhysicsSystem) SetRoom(room world.Room) {
	ps.room = room
}

func (ps *PhysicsSystem) SetCurrentRoom(room world.Room) {
	ps.SetRoom(room)
}

func (ps *PhysicsSystem) ClearEnemies() {
	ps.enemies = ps.enemies[:0]
}

func (ps *PhysicsSystem) Update() error {
	// Update player physics
	ps.player.UpdatePhysics()
	
	// Update enemy physics
	for _, enemy := range ps.enemies {
		if err := enemy.Update(); err != nil {
			return fmt.Errorf("enemy update failed: %w", err)
		}
	}
	
	// Handle collisions if room is set
	if ps.room != nil {
		// Player collision with room
		ps.room.CheckCollision(ps.player)
		
		// Enemy collision with room
		for _, enemy := range ps.enemies {
			if collisionEntity, ok := enemy.(entities.CollisionHandler); ok {
				ps.room.CheckCollision(collisionEntity)
			}
		}
	}
	
	return nil
}

/*
CameraSystem manages the game camera and viewport.
Handles camera following, boundaries, and smooth transitions.
*/
type CameraSystem struct {
	camera   *engine.Camera
	player   *entities.Player
	room     world.Room
	enabled  bool
}

func NewCameraSystem(camera *engine.Camera, player *entities.Player) *CameraSystem {
	return &CameraSystem{
		camera:  camera,
		player:  player,
		enabled: true,
	}
}

func (cs *CameraSystem) GetName() string {
	return "Camera"
}

func (cs *CameraSystem) SetRoom(room world.Room) {
	cs.room = room
	// Update camera bounds when room changes
	if cs.camera != nil && room != nil {
		if tileMap := room.GetTileMap(); tileMap != nil {
			worldWidth := tileMap.Width * engine.GetPhysicsUnit()
			worldHeight := tileMap.Height * engine.GetPhysicsUnit()
			cs.camera.SetWorldBounds(worldWidth, worldHeight)
		}
	}
}

func (cs *CameraSystem) SetCurrentRoom(room world.Room) {
	cs.SetRoom(room)
}

func (cs *CameraSystem) Update() error {
	if !cs.enabled || cs.camera == nil || cs.player == nil {
		return nil
	}
	
	// Update camera to follow player
	playerX, playerY := cs.player.GetPosition()
	cs.camera.FollowTarget(float64(playerX), float64(playerY))
	cs.camera.Update()
	
	return nil
}

func (cs *CameraSystem) SetEnabled(enabled bool) {
	cs.enabled = enabled
}

/*
RoomSystem manages room transitions and current room state.
*/
type RoomSystem struct {
	transitionManager *world.RoomTransitionManager
	worldMap          *world.WorldMap
	player            *entities.Player
}

func NewRoomSystem(transitionManager *world.RoomTransitionManager, worldMap *world.WorldMap, player *entities.Player) *RoomSystem {
	return &RoomSystem{
		transitionManager: transitionManager,
		worldMap:          worldMap,
		player:            player,
	}
}

func (rs *RoomSystem) GetName() string {
	return "Room"
}

func (rs *RoomSystem) Update() error {
	// Update current room
	currentRoom := rs.transitionManager.GetCurrentRoom()
	if currentRoom == nil {
		return nil
	}
	
	// Update room logic
	if err := currentRoom.Update(rs.player); err != nil {
		return fmt.Errorf("room update failed: %w", err)
	}
	
	// Check for room transitions
	if transition := currentRoom.CheckTransition(rs.player); transition != nil {
		if err := rs.transitionManager.StartTransitionToRoom(transition.Direction, transition.TargetRoomID, transition.Duration); err != nil {
			return fmt.Errorf("room transition failed: %w", err)
		}
		
		// Position player in new room
		if newRoom := rs.transitionManager.GetCurrentRoom(); newRoom != nil {
			newRoom.PositionPlayerAtEntrance(rs.player, transition.Direction.Opposite())
		}
	}
	
	return nil
}

func (rs *RoomSystem) GetCurrentRoom() world.Room {
	return rs.transitionManager.GetCurrentRoom()
}

/*
GameSystemManager manages all game systems and their update order.
*/
type GameSystemManager struct {
	systems     map[string]GameSystem
	updateOrder []string
}

func NewGameSystemManager() *GameSystemManager {
	return &GameSystemManager{
		systems:     make(map[string]GameSystem),
		updateOrder: make([]string, 0),
	}
}

// AddSystem adds a new system to the manager
func (gsm *GameSystemManager) AddSystem(name string, system GameSystem) {
	gsm.systems[name] = system
	// Add to update order if not already present
	found := false
	for _, n := range gsm.updateOrder {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		gsm.updateOrder = append(gsm.updateOrder, name)
	}
}

// GetSystem returns a system by name
func (gsm *GameSystemManager) GetSystem(name string) GameSystem {
	return gsm.systems[name]
}

// UpdateAll updates all systems in order
func (gsm *GameSystemManager) UpdateAll() error {
	// Update in specified order
	for _, systemName := range gsm.updateOrder {
		if system, exists := gsm.systems[systemName]; exists {
			if err := system.Update(); err != nil {
				return fmt.Errorf("system %s update failed: %w", systemName, err)
			}
		}
	}
	
	// Update any systems not in the order list
	for name, system := range gsm.systems {
		found := false
		for _, orderedName := range gsm.updateOrder {
			if name == orderedName {
				found = true
				break
			}
		}
		if !found {
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