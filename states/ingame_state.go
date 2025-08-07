package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"sword/entities"
	"sword/systems"
	"sword/world"
)

/*
InGameState represents the main gameplay state using modular systems.
This implementation uses a clean system architecture for better separation
of concerns and easier testing/maintenance.

Key features:
  - Modular system architecture with clear responsibilities
  - Proper room transition management
  - Reduced coupling between components
  - Better error handling and state management
*/
type InGameState struct {
	stateManager *engine.StateManager

	// Core game entities
	player  *entities.Player
	enemies []entities.Enemy

	// Modular systems
	systemManager     *systems.GameSystemManager
	roomTransitionMgr *world.RoomTransitionManager
	worldMap          *world.WorldMap

	// Rendering systems
	camera           *engine.Camera
	viewportRenderer *engine.ViewportRenderer
	hudManager       *engine.HUDManager

	// Configuration and state
	parallaxConfigIndex int
}

/*
NewInGameState creates a new in-game state.
Initializes all systems and sets up the modular architecture.
*/
func NewInGameState(sm *engine.StateManager) *InGameState {
	// Get window size for viewport setup
	windowWidth, windowHeight := ebiten.WindowSize()
	physicsUnit := engine.GetPhysicsUnit()

	// Create the initial room
	room := world.NewSimpleRoom("main")

	// Calculate spawn position
	tileMap := room.GetTileMap()
	playerSpawnX := (tileMap.Width / 2) * physicsUnit
	playerSpawnY := (tileMap.Height - 2) * physicsUnit

	// For larger rooms, use floor detection
	if tileMap.Width > 10 || tileMap.Height > 10 {
		groundY := room.FindFloorAtX(playerSpawnX)
		if groundY > 0 {
			playerSpawnY = groundY
		}
	}

	// Create core entities
	player := entities.NewPlayer(playerSpawnX, playerSpawnY)

	// Initialize world map
	worldMap := world.NewWorldMap()

	// Initialize room transition system
	roomTransitionMgr := world.NewRoomTransitionManager(worldMap)
	roomTransitionMgr.RegisterRoom(room)
	roomTransitionMgr.SetCurrentRoom(room.GetZoneID())

	// Add spawn points to the room
	roomTransitionMgr.AddSpawnPoint(room.GetZoneID(), world.SpawnPoint{
		ID: "main_spawn",
		X:  playerSpawnX,
		Y:  playerSpawnY,
	})

	worldMap.SetCurrentRoom(room.GetZoneID())

	// Create camera and viewport systems
	camera := engine.NewCamera(windowWidth, windowHeight)
	if tileMap != nil {
		worldWidth := tileMap.Width * physicsUnit
		worldHeight := tileMap.Height * physicsUnit
		camera.SetWorldBounds(worldWidth, worldHeight)
	}

	viewportRenderer := engine.NewViewportRenderer(windowWidth, windowHeight)
	if tileMap != nil {
		viewportRenderer.SetWorldBounds(tileMap.Width*physicsUnit, tileMap.Height*physicsUnit)
	}

	// Initialize HUD system
	hudManager := engine.NewHUDManager()

	// Set up debug HUD
	debugHUD := engine.NewDebugHUD()
	hudManager.AddComponent(debugHUD)

	// Set up minimap - using proper constructor with required parameters
	minimapRenderer := world.NewMiniMapRenderer(worldMap, 200, windowWidth-220, 20)
	hudManager.AddComponent(minimapRenderer)

	// Create the refactored state
	state := &InGameState{
		stateManager:        sm,
		player:              player,
		enemies:             make([]entities.Enemy, 0),
		roomTransitionMgr:   roomTransitionMgr,
		worldMap:            worldMap,
		camera:              camera,
		viewportRenderer:    viewportRenderer,
		hudManager:          hudManager,
		parallaxConfigIndex: 0,
	}

	// Initialize modular systems
	state.initializeSystems()

	return state
}

/*
initializeSystems sets up the modular game systems.
*/
func (ris *InGameState) initializeSystems() {
	ris.systemManager = systems.NewGameSystemManager()

	// Create and register systems
	inputSystem := systems.NewInputSystem(ris.player, ris.roomTransitionMgr)
	physicsSystem := systems.NewPhysicsSystem(ris.player)
	cameraSystem := systems.NewCameraSystem(ris.camera, ris.player)
	roomSystem := systems.NewRoomSystem(ris.roomTransitionMgr, ris.worldMap, ris.player)

	// Set initial room for systems
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom != nil {
		physicsSystem.SetCurrentRoom(currentRoom)
		cameraSystem.SetCurrentRoom(currentRoom)
	}

	// Register systems with manager
	ris.systemManager.AddSystem("input", inputSystem)
	ris.systemManager.AddSystem("physics", physicsSystem)
	ris.systemManager.AddSystem("camera", cameraSystem)
	ris.systemManager.AddSystem("room", roomSystem)

	// Set update order: Input -> Room -> Physics -> Camera
	ris.systemManager.SetUpdateOrder([]string{"Input", "Room", "Physics", "Camera"})
}

/*
Update implements the game update loop using modular systems.
*/
func (ris *InGameState) Update() error {
	// Handle camera viewport changes
	ris.updateCameraViewport()

	// Update all systems
	err := ris.systemManager.UpdateAll()
	if err != nil {
		return err
	}

	// Handle state transition requests from input system
	if inputSystem := ris.systemManager.GetSystem("Input"); inputSystem != nil {
		if is, ok := inputSystem.(*systems.InputSystem); ok {
			if is.HasPauseRequest() {
				is.ClearRequests()
				pauseState := NewPauseState(ris.stateManager, ris)
				ris.stateManager.ChangeState(pauseState)
				return nil
			}
			if is.HasSettingsRequest() {
				is.ClearRequests()
				currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
				settingsState := NewSettingsState(ris.stateManager, currentRoom)
				ris.stateManager.ChangeState(settingsState)
				return nil
			}
		}
	}

	// Handle room changes - update other systems when room changes
	if roomSystem := ris.systemManager.GetSystem("Room"); roomSystem != nil {
		if rs, ok := roomSystem.(*systems.RoomSystem); ok {
			currentRoom := rs.GetCurrentRoom()
			if currentRoom != nil {
				// Update physics system with new room
				if physicsSystem := ris.systemManager.GetSystem("Physics"); physicsSystem != nil {
					if ps, ok := physicsSystem.(*systems.PhysicsSystem); ok {
						ps.SetCurrentRoom(currentRoom)
					}
				}

				// Update camera system with new room
				if cameraSystem := ris.systemManager.GetSystem("Camera"); cameraSystem != nil {
					if cs, ok := cameraSystem.(*systems.CameraSystem); ok {
						cs.SetCurrentRoom(currentRoom)
					}
				}
			}
		}
	}

	// Handle additional input that's not part of the input system
	// Enhanced parallax controls
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		playerX, playerY := ris.player.GetPosition()
		roomName := ""
		if currentRoom := ris.roomTransitionMgr.GetCurrentRoom(); currentRoom != nil {
			roomName = currentRoom.GetZoneID()
		}
		engine.LogPlayerInput("D (Toggle Depth of Field)", playerX, playerY, roomName)
		ris.toggleDepthOfField()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		playerX, playerY := ris.player.GetPosition()
		roomName := ""
		if currentRoom := ris.roomTransitionMgr.GetCurrentRoom(); currentRoom != nil {
			roomName = currentRoom.GetZoneID()
		}
		engine.LogPlayerInput("L (Cycle Parallax Layers)", playerX, playerY, roomName)
		ris.cycleParallaxLayers()
	}

	// Update debug HUD with current state
	ris.updateDebugHUD()

	// Update HUD systems
	if err := ris.hudManager.Update(); err != nil {
		return err
	}

	return nil
}

/*
updateDebugHUD updates the debug HUD with current game state information.
*/
func (ris *InGameState) updateDebugHUD() {
	if debugHUD := ris.hudManager.GetComponent("debug_hud"); debugHUD != nil {
		if dh, ok := debugHUD.(*engine.DebugHUD); ok {
			// Update room info
			roomInfo := "No Room"
			currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
			if currentRoom != nil {
				roomInfo = currentRoom.GetZoneID()
				if tileMap := currentRoom.GetTileMap(); tileMap != nil {
					roomInfo = fmt.Sprintf("%s (%dx%d tiles)", roomInfo, tileMap.Width, tileMap.Height)
				}
			}
			dh.UpdateRoomInfo(roomInfo)

			// Update player position
			playerX, playerY := ris.player.GetPosition()
			physicsUnit := engine.GetPhysicsUnit()
			playerPixelX := float64(playerX) / float64(physicsUnit)
			playerPixelY := float64(playerY) / float64(physicsUnit)
			playerTileX := int(playerPixelX / float64(engine.GameConfig.TileSize) / engine.GameConfig.TileScaleFactor)
			playerTileY := int(playerPixelY / float64(engine.GameConfig.TileSize) / engine.GameConfig.TileScaleFactor)
			playerPos := fmt.Sprintf("Physics: (%d, %d) | Pixels: (%.1f, %.1f) | Tiles: (%d, %d)",
				playerX, playerY, playerPixelX, playerPixelY, playerTileX, playerTileY)
			dh.UpdatePlayerPos(playerPos)

			// Update player velocity
			vx, vy := ris.player.GetVelocity()
			velocityPixelX := float64(vx) / float64(physicsUnit)
			velocityPixelY := float64(vy) / float64(physicsUnit)
			playerVelocity := fmt.Sprintf("Velocity: Physics: (%d, %d) | Pixels/frame: (%.1f, %.1f)",
				vx, vy, velocityPixelX, velocityPixelY)
			dh.UpdatePlayerVelocity(playerVelocity)

			// Update player status
			onGround := "In Air"
			if ris.player.IsOnGround() {
				onGround = "On Ground"
			}
			facing := "Left"
			if ris.player.IsFacingRight() {
				facing = "Right"
			}
			playerStatus := fmt.Sprintf("Status: %s | Facing: %s", onGround, facing)
			dh.UpdatePlayerStatus(playerStatus)

			// Update camera position
			if ris.camera != nil {
				camX, camY := ris.camera.GetPosition()
				cameraPos := fmt.Sprintf("Camera: (%.1f, %.1f)", camX, camY)
				dh.UpdateCameraPos(cameraPos)
			}
		}
	}
}

/*
Draw renders the game using the modular rendering systems.
*/
func (ris *InGameState) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Clear()

	// Draw background if enabled
	if engine.GetBackgroundVisible() {
		if backgroundImage := engine.GetBackgroundImage(); backgroundImage != nil {
			// Scale and draw background
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(engine.GameConfig.TileScaleFactor, engine.GameConfig.TileScaleFactor)
			screen.DrawImage(backgroundImage, op)
		}
	}

	// Draw current room with camera offset
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom != nil {
		cameraX, cameraY := ris.camera.GetPosition()
		currentRoom.DrawWithCamera(screen, float64(cameraX), float64(cameraY))
	}

	// Draw player - Player uses its own Draw method without camera offset
	// The camera offset is handled internally by the player
	if ris.player != nil {
		ris.player.Draw(screen)
	}

	// Draw enemies with camera offset
	cameraX, cameraY := ris.camera.GetPosition()
	for _, enemy := range ris.enemies {
		enemy.DrawWithCamera(screen, float64(cameraX), float64(cameraY))
	}

	// Draw debug grid if enabled
	if engine.GetGridVisible() {
		cameraX, cameraY := ris.camera.GetPosition()
		engine.DrawGridWithCamera(screen, float64(cameraX), float64(cameraY))
	}

	// Draw HUD elements on top
	ris.hudManager.Draw(screen)

	// Draw simple debug info
	playerX, playerY := ris.player.GetPosition()
	debugText := fmt.Sprintf("Player: (%d, %d)", playerX, playerY)
	ebitenutil.DebugPrint(screen, debugText)
}

/*
OnEnter is called when entering the game state.
*/
func (ris *InGameState) OnEnter() {
	engine.LogInfo("Entered refactored in-game state")

	// Initialize camera to follow player
	if ris.player != nil && ris.camera != nil {
		px, py := ris.player.GetPosition()
		ris.camera.Update(px, py)
	}
}

/*
OnExit is called when leaving the game state.
*/
func (ris *InGameState) OnExit() {
	engine.LogInfo("Exited refactored in-game state")

	// Let the current room know we're leaving
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom != nil {
		currentRoom.OnExit(ris.player)
	}
}

/*
GetCurrentRoom returns the current room being played.
*/
func (ris *InGameState) GetCurrentRoom() world.Room {
	return ris.roomTransitionMgr.GetCurrentRoom()
}

/*
AddEnemy adds an enemy to the game state.
*/
func (ris *InGameState) AddEnemy(enemy entities.Enemy) {
	ris.enemies = append(ris.enemies, enemy)

	// Add to physics system
	if physicsSystem := ris.systemManager.GetSystem("Physics"); physicsSystem != nil {
		if ps, ok := physicsSystem.(*systems.PhysicsSystem); ok {
			ps.AddEnemy(enemy)
		}
	}
}

/*
ClearEnemies removes all enemies from the game state.
*/
func (ris *InGameState) ClearEnemies() {
	ris.enemies = nil

	// Clear from physics system
	if physicsSystem := ris.systemManager.GetSystem("Physics"); physicsSystem != nil {
		if ps, ok := physicsSystem.(*systems.PhysicsSystem); ok {
			ps.ClearEnemies()
		}
	}
}

/*
RemoveEnemy removes a specific enemy from the game state.
Returns true if the enemy was found and removed, false otherwise.
*/
func (ris *InGameState) RemoveEnemy(enemy entities.Enemy) bool {
	for i, e := range ris.enemies {
		if e == enemy {
			// Remove by swapping with last element and truncating
			ris.enemies[i] = ris.enemies[len(ris.enemies)-1]
			ris.enemies = ris.enemies[:len(ris.enemies)-1]

			// Remove from physics system
			if physicsSystem := ris.systemManager.GetSystem("Physics"); physicsSystem != nil {
				if ps, ok := physicsSystem.(*systems.PhysicsSystem); ok {
					// Need to rebuild the physics system enemy list
					ps.ClearEnemies()
					for _, remaining := range ris.enemies {
						ps.AddEnemy(remaining)
					}
				}
			}

			return true
		}
	}
	return false
}

/*
GetEnemies returns a copy of the current enemies slice.
*/
func (ris *InGameState) GetEnemies() []entities.Enemy {
	result := make([]entities.Enemy, len(ris.enemies))
	copy(result, ris.enemies)
	return result
}

/*
toggleDepthOfField toggles the depth of field effect.
*/
func (ris *InGameState) toggleDepthOfField() {
	engine.GameConfig.EnableDepthOfField = !engine.GameConfig.EnableDepthOfField

	// Update debug HUD to reflect the change
	ris.updateDebugHUD()
}

/*
cycleParallaxLayers cycles through different parallax layer configurations.
*/
func (ris *InGameState) cycleParallaxLayers() {
	// Define layer configurations
	layerConfigs := [][]engine.ParallaxLayer{
		// Config 0: Full layers (default)
		{
			{Speed: 0.1, OffsetY: 0},   // Far background
			{Speed: 0.3, OffsetY: 50},  // Mid background
			{Speed: 0.5, OffsetY: 100}, // Near background
			{Speed: 0.7, OffsetY: 150}, // Foreground
		},
		// Config 1: Simple two-layer
		{
			{Speed: 0.2, OffsetY: 0},   // Background
			{Speed: 0.6, OffsetY: 100}, // Foreground
		},
		// Config 2: Single static background
		{
			{Speed: 0.0, OffsetY: 0}, // Static
		},
		// Config 3: Extreme depth
		{
			{Speed: 0.05, OffsetY: 0},  // Very far
			{Speed: 0.2, OffsetY: 30},  // Far
			{Speed: 0.4, OffsetY: 60},  // Mid
			{Speed: 0.6, OffsetY: 90},  // Near
			{Speed: 0.8, OffsetY: 120}, // Very near
		},
	}

	// Cycle to next configuration
	ris.parallaxConfigIndex = (ris.parallaxConfigIndex + 1) % len(layerConfigs)

	// Apply new configuration
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom != nil {
		if simpleRoom, ok := currentRoom.(*world.SimpleRoom); ok {
			renderer := simpleRoom.GetParallaxRenderer()
			if renderer != nil {
				renderer.SetLayers(layerConfigs[ris.parallaxConfigIndex])
			}
		}
	}
}

/*
updateCameraViewport handles window resize events.
*/
func (ris *InGameState) updateCameraViewport() {
	if ris.camera == nil {
		return
	}

	currentWidth, currentHeight := ebiten.WindowSize()
	cameraWidth, cameraHeight := ris.camera.GetViewportSize()

	// Check if window size has changed
	if currentWidth != cameraWidth || currentHeight != cameraHeight {
		// Create new camera with updated viewport
		ris.camera = engine.NewCamera(currentWidth, currentHeight)

		// Update viewport renderer
		ris.viewportRenderer = engine.NewViewportRenderer(currentWidth, currentHeight)

		// Restore world bounds and camera position
		currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
		if currentRoom != nil {
			tileMap := currentRoom.GetTileMap()
			if tileMap != nil {
				physicsUnit := engine.GetPhysicsUnit()
				worldWidth := tileMap.Width * physicsUnit
				worldHeight := tileMap.Height * physicsUnit
				ris.camera.SetWorldBounds(worldWidth, worldHeight)
				ris.viewportRenderer.SetWorldBounds(worldWidth, worldHeight)

				// Update camera system
				if cameraSystem := ris.systemManager.GetSystem("Camera"); cameraSystem != nil {
					if cs, ok := cameraSystem.(*systems.CameraSystem); ok {
						cs.SetCurrentRoom(currentRoom)
					}
				}
			}
		}
	}
}
