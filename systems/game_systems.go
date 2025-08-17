package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	player            *entities.Player
	roomTransitionMgr *world.RoomTransitionManager
	pauseRequested    bool
	settingsRequested bool // Add settings request flag

	// UI toggle callback to avoid coupling with HUD
	OnToggleMinimap func()
}

/*
NewInputSystem creates a new input system instance.
Parameters:
  - player: The player entity to control
  - roomTransitionMgr: Manager for handling room transitions
*/
func NewInputSystem(player *entities.Player, roomTransitionMgr *world.RoomTransitionManager) *InputSystem {
	return &InputSystem{
		player:            player,
		roomTransitionMgr: roomTransitionMgr,
		pauseRequested:    false,
		settingsRequested: false, // Initialize settings request flag
	}
}

func (is *InputSystem) GetName() string {
	return "Input"
}

func (is *InputSystem) Update() error {
	// Handle room transitions first
	if is.roomTransitionMgr != nil {
		// Check for room transitions
		is.roomTransitionMgr.CheckTransitions(is.player, ebiten.IsKeyPressed(ebiten.KeyE))
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			is.logKeyPress("E (Interact)")
		}
	}

	// Handle movement inputs - Player handles its own input
	is.player.ProcessInput()

	// Pause request handling
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		is.pauseRequested = true
		is.logKeyPress("Escape (Pause)")
	}

	// Minimap toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		is.logKeyPress("M (Toggle Minimap)")
		if is.OnToggleMinimap != nil {
			is.OnToggleMinimap()
		}
	}

	// Log movement and action keys
	keys := []ebiten.Key{
			ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp, ebiten.KeyDown,
			ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeySpace,
	}
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			is.logKeyPress(k.String())
		}
	}

	return nil
}

func (is *InputSystem) logKeyPress(desc string) {
	playerX, playerY := is.player.GetPosition()
	roomName := ""
	if is.roomTransitionMgr != nil {
		if currentRoom := is.roomTransitionMgr.GetCurrentRoom(); currentRoom != nil {
			roomName = currentRoom.GetZoneID()
		}
	}
	engine.LogPlayerInput(desc, playerX, playerY, roomName)
}

func (is *InputSystem) HasPauseRequest() bool {
	return is.pauseRequested
}

func (is *InputSystem) HasSettingsRequest() bool {
	return is.settingsRequested
}

func (is *InputSystem) ClearRequests() {
	is.pauseRequested = false
	is.settingsRequested = false
}

/*
PhysicsSystem handles physics simulation for all entities.
Manages gravity, collisions, and movement for the player and enemies.
*/
type PhysicsSystem struct {
	player  *entities.Player
	enemies []entities.Enemy
	room    world.Room
}

/*
NewPhysicsSystem creates a new physics system instance.
*/
func NewPhysicsSystem(player *entities.Player) *PhysicsSystem {
	return &PhysicsSystem{
		player:  player,
		enemies: make([]entities.Enemy, 0),
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

// Update updates physics for all entities
func (ps *PhysicsSystem) Update() error {
	// Log player state before physics update
	px, py := ps.player.GetPosition()
	vx, vy := ps.player.GetVelocity()
	onGround := ps.player.IsOnGround()
	
	roomName := "NoRoom"
	if ps.room != nil {
		roomName = ps.room.GetZoneID()
	}
	
	engine.LogDebug(fmt.Sprintf("PHYSICS_BEFORE: Room=%s Pos=(%d,%d) Vel=(%d,%d) OnGround=%v", 
		roomName, px, py, vx, vy, onGround))
	
	// Update player physics with tiles when room is present
	if ps.room != nil {
		if tileProvider, ok := ps.room.(entities.TileProvider); ok {
			engine.LogDebug("PHYSICS: Using UpdateWithTileCollision")
			ps.player.UpdateWithTileCollision(tileProvider)
		} else {
			engine.LogDebug("PHYSICS: Room doesn't implement TileProvider, using basic Update")
			ps.player.Update()
		}
	} else {
		engine.LogDebug("PHYSICS: No room, using basic Update")
		ps.player.Update()
	}
	
	// Log player state after physics update
	px2, py2 := ps.player.GetPosition()
	vx2, vy2 := ps.player.GetVelocity()
	onGround2 := ps.player.IsOnGround()
	
	engine.LogDebug(fmt.Sprintf("PHYSICS_AFTER: Room=%s Pos=(%d,%d) Vel=(%d,%d) OnGround=%v", 
		roomName, px2, py2, vx2, vy2, onGround2))
	
	// Log movement delta if any
	if px != px2 || py != py2 {
		engine.LogDebug(fmt.Sprintf("PHYSICS_DELTA: ΔPos=(%d,%d) ΔVel=(%d,%d)", 
			px2-px, py2-py, vx2-vx, vy2-vy))
	}
	
	// Update enemies
	for _, enemy := range ps.enemies {
		enemy.Update()
	}
	
	// Handle collision detection for enemies if needed (player handled above)
	return nil
}

/*
CameraSystem manages the game camera and viewport.
Handles camera following, boundaries, and smooth transitions.
*/
type CameraSystem struct {
	camera  *engine.Camera
	player  *entities.Player
	room    world.Room
	enabled bool
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
			u := engine.GetPhysicsUnit()
			scale := engine.GameConfig.TileScaleFactor
			cs.camera.SetWorldBounds(int(float64(tileMap.Width*u)*scale), int(float64(tileMap.Height*u)*scale))
		}
	}
}

func (cs *CameraSystem) SetCurrentRoom(room world.Room) {
	cs.SetRoom(room)
}

// Update updates the camera to follow the player
func (cs *CameraSystem) Update() error {
	// Get player position for camera tracking
	playerX, playerY := cs.player.GetPosition()

	// Update camera position in scaled screen pixels
	s := engine.GameConfig.TileScaleFactor
	cs.camera.Update(int(float64(playerX)*s), int(float64(playerY)*s))

	return nil
}

func (cs *CameraSystem) SetEnabled(enabled bool) {
	cs.enabled = enabled
}

/*
RoomSystem manages room transitions and world state changes.
Handles checking for and processing room transitions and updates
other systems when the room changes.
*/
type RoomSystem struct {
	transitionManager *world.RoomTransitionManager
	worldMap          *world.WorldMap
	player            *entities.Player
	physicsSystem     *PhysicsSystem
	cameraSystem      *CameraSystem
}

func NewRoomSystem(transitionManager *world.RoomTransitionManager, worldMap *world.WorldMap, player *entities.Player) *RoomSystem {
	return &RoomSystem{
		transitionManager: transitionManager,
		worldMap:          worldMap,
		player:            player,
		physicsSystem:     nil, // Will be initialized later
		cameraSystem:      nil, // Will be initialized later
	}
}

func (rs *RoomSystem) GetName() string {
	return "Room"
}

// SetPhysicsSystem sets the physics system reference
func (rs *RoomSystem) SetPhysicsSystem(ps *PhysicsSystem) {
	rs.physicsSystem = ps
}

// SetCameraSystem sets the camera system reference
func (rs *RoomSystem) SetCameraSystem(cs *CameraSystem) {
	rs.cameraSystem = cs
}

// Update checks for room transitions
func (rs *RoomSystem) Update() error {
	if rs.transitionManager != nil && rs.player != nil {
		// Safety: if player fell below current room bounds, portal to safety room
		if current := rs.transitionManager.GetCurrentRoom(); current != nil {
			tm := current.GetTileMap()
			u := engine.GetPhysicsUnit()
			_, py := rs.player.GetPosition()
			if tm != nil {
				maxY := tm.Height * u
				if py > maxY+u { // allow small margin
					// Queue a transition if safety room exists
					if len(rs.transitionManager.GetSpawnPoints("safety")) > 0 {
						// Create a pending transition to safety room's default spawn
						// We use CheckTransitions/ProcessPendingTransition pathway by directly setting pending
						// Not exposed: fallback to direct spawn after SetCurrentRoom
						rs.transitionManager.SetCurrentRoom("safety")
						// Try spawn id "entry" then first spawn
						if err := rs.transitionManager.SpawnPlayerInRoom(rs.player, "safety", "entry"); err != nil {
							// ignore error, player may be placed later by fallback
						}
						// Update camera/physics to new room immediately
						if rs.physicsSystem != nil {
							rs.physicsSystem.SetRoom(rs.transitionManager.GetCurrentRoom())
						}
						if rs.cameraSystem != nil {
							rs.cameraSystem.SetRoom(rs.transitionManager.GetCurrentRoom())
						}
					}
				}
			}
		}

		// Process any pending transitions
		if rs.transitionManager.HasPendingTransition() {
			newRoom, err := rs.transitionManager.ProcessPendingTransition(rs.player)
			if err != nil {
				return fmt.Errorf("failed to process room transition: %w", err)
			}

			if newRoom != nil {
				// Recompute tile scale to fit the new room to the current window (zoom in small rooms)
				if tm := newRoom.GetTileMap(); tm != nil {
					u := engine.GetPhysicsUnit()
					winW, winH := ebiten.WindowSize()
					roomPxW := tm.Width * u
					roomPxH := tm.Height * u
					fitScaleW := float64(winW) / float64(roomPxW)
					fitScaleH := float64(winH) / float64(roomPxH)
					fitScale := math.Min(fitScaleW, fitScaleH)
					// Clamp between 1x and 4x
					if fitScale < 1.0 {
						fitScale = 1.0
					}
					if fitScale > 4.0 {
						fitScale = 4.0
					}
					engine.GameConfig.TileScaleFactor = fitScale
				}

				// Notify other systems about the room change
				if rs.physicsSystem != nil {
					rs.physicsSystem.SetRoom(newRoom)
				}
				if rs.cameraSystem != nil {
					// After updating scale, update camera bounds for new room
					rs.cameraSystem.SetRoom(newRoom)
				}
			}
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
