package states

import (
	"fmt"
	"strings"

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


}

/*
NewInGameState creates a new in-game state.
Initializes all systems and sets up the modular architecture.
*/
func NewInGameState(sm *engine.StateManager) *InGameState {
	// Get window size for viewport setup
	windowWidth, windowHeight := ebiten.WindowSize()

	// Initialize world map and room transition system
	worldMap := world.NewWorldMap()
	roomTransitionMgr := world.NewRoomTransitionManager(worldMap)

	// Load Tiled zone rooms once into the active manager
	mainRoom := world.Room(nil)
	tileMap := (*world.TileMap)(nil)

	playerSpawnX, playerSpawnY := 0, 0

	zoneName := "cradle"
	if err := world.LoadZoneRoomsFromData(roomTransitionMgr, zoneName, "assets"); err != nil {
		panic(fmt.Errorf("failed to load Tiled zone '%s': %w", zoneName, err))
	}
	// Choose a deterministic first room id if available
	mainRoomID := ""
	for _, id := range roomTransitionMgr.ListRoomIDs() {
		if strings.HasSuffix(id, "/r01") {
			mainRoomID = id
			break
		}
	}
	if mainRoomID == "" {
		ids := roomTransitionMgr.ListRoomIDs()
		if len(ids) == 0 {
			panic("no rooms loaded from Tiled zone")
		}
		mainRoomID = ids[0]
	}

	mainRoom = roomTransitionMgr.GetRoom(mainRoomID)
	if mainRoom == nil {
		panic(fmt.Errorf("main room '%s' not found after Tiled load", mainRoomID))
	}
	if err := roomTransitionMgr.SetCurrentRoom(mainRoomID); err != nil {
		panic(err)
	}

	// Compute initial scale framing using the resolved mainRoom
	if tm := mainRoom.GetTileMap(); tm != nil {
		roomPxW := tm.Width * engine.GameConfig.TileSize
		roomPxH := tm.Height * engine.GameConfig.TileSize
		fitScaleW := float64(windowWidth) / float64(roomPxW)
		fitScaleH := float64(windowHeight) / float64(roomPxH)
		fitScale := fitScaleW
		if fitScaleH < fitScale {
			fitScale = fitScaleH
		}
		if fitScale > 1.0 {
			maxScale := 4.0
			if fitScale > maxScale {
				engine.GameConfig.TileScaleFactor = maxScale
			} else {
				engine.GameConfig.TileScaleFactor = fitScale
			}
		}
	}

	// Determine spawn position - find a safe open area  
	tileMap = mainRoom.GetTileMap()
	playerSpawnX = 80  // Try middle of room (tile X=5, 5*16=80)
	playerSpawnY = 32  // Fallback Y position
	
	// Try to find the floor at this X position
	if tileMap != nil {
		if groundY := mainRoom.FindFloorAtX(playerSpawnX); groundY >= 0 {
			playerSpawnY = groundY
			engine.LogInfo(fmt.Sprintf("Found floor at Y=%d for X=%d, spawning player at (%d,%d)", groundY, playerSpawnX, playerSpawnX, playerSpawnY))
		} else {
			engine.LogInfo(fmt.Sprintf("No valid floor found at X=%d, using fallback spawn position (%d,%d)", playerSpawnX, playerSpawnX, playerSpawnY))
		}
	}

	// Create core entities
	player := entities.NewPlayer(playerSpawnX, playerSpawnY)

	// Set world map room and initial trail position
	worldMap.SetCurrentRoom(mainRoom.GetZoneID())
	px, py := player.GetPosition()
	worldMap.AddPlayerPosition(px, py)

	// Create camera and viewport systems
	camera := engine.NewCamera(windowWidth, windowHeight)
	if tileMap != nil {
		u := engine.GetPhysicsUnit()
		scaledWorldWidth := int(float64(tileMap.Width*u) * engine.GameConfig.TileScaleFactor)
		scaledWorldHeight := int(float64(tileMap.Height*u) * engine.GameConfig.TileScaleFactor)
		camera.SetWorldBounds(scaledWorldWidth, scaledWorldHeight)
	}

	viewportRenderer := engine.NewViewportRenderer(windowWidth, windowHeight)
	if tileMap != nil {
		u := engine.GetPhysicsUnit()
		scaledWorldWidth := int(float64(tileMap.Width*u) * engine.GameConfig.TileScaleFactor)
		scaledWorldHeight := int(float64(tileMap.Height*u) * engine.GameConfig.TileScaleFactor)
		viewportRenderer.SetWorldBounds(scaledWorldWidth, scaledWorldHeight)
	}

	// Initialize HUD system
	hudManager := engine.NewHUDManager()
	debugHUD := engine.NewDebugHUD()
	hudManager.AddComponent(debugHUD)
	minimapRenderer := world.NewMiniMapRenderer(worldMap, player, 200, windowWidth-220, 20)
	hudManager.AddComponent(minimapRenderer)
	zoneMapOverlay := world.NewZoneMapOverlay(worldMap, player)
	hudManager.AddComponent(zoneMapOverlay)

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

	// Wire HUD toggles
	inputSystem.OnToggleMinimap = func() {
		if ris.hudManager != nil {
			ris.hudManager.ToggleComponent("minimap")
		}
	}

	// Set initial room for systems
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom != nil {
		physicsSystem.SetCurrentRoom(currentRoom)
		cameraSystem.SetCurrentRoom(currentRoom)
	}

	// Wire systems so room changes update physics and camera
	if roomSystem != nil {
		roomSystem.SetPhysicsSystem(physicsSystem)
		roomSystem.SetCameraSystem(cameraSystem)
	}

	// Register systems with manager
	ris.systemManager.AddSystem("Input", inputSystem)
	ris.systemManager.AddSystem("Physics", physicsSystem)
	ris.systemManager.AddSystem("Camera", cameraSystem)
	ris.systemManager.AddSystem("Room", roomSystem)

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

	// Feed world map player position for trail/minimap
	if ris.worldMap != nil && ris.player != nil {
		px, py := ris.player.GetPosition()
		// Convert to room-local pixel coords if needed; player coordinates are already in pixels
		// but trail is drawn relative to current room in minimap renderer
		ris.worldMap.AddPlayerPosition(px, py)
	}

	// Toggle zone map overlay
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if ris.hudManager != nil {
			ris.hudManager.ToggleComponent("zone_map")
		}
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
					u := engine.GetPhysicsUnit()
		playerTileX := playerX / u
		playerTileY := playerY / u
		playerPos := fmt.Sprintf("World(px): (%d, %d) | Tiles: (%d, %d)",
			playerX, playerY, playerTileX, playerTileY)
		dh.UpdatePlayerPos(playerPos)

			// Update player velocity
			vx, vy := ris.player.GetVelocity()
			playerVelocity := fmt.Sprintf("Velocity (px/frame): (%d, %d)", vx, vy)
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

	// Background is handled via parallax in rooms that support it; only draw the fallback background when no room is active.
	currentRoom := ris.roomTransitionMgr.GetCurrentRoom()
	if currentRoom == nil {
		if engine.GetBackgroundVisible() {
			if backgroundImage := engine.GetBackgroundImage(); backgroundImage != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(engine.GameConfig.TileScaleFactor, engine.GameConfig.TileScaleFactor)
				screen.DrawImage(backgroundImage, op)
			}
		}
	}

	engine.LogDebug("DRAW_LAYER: Room")
	if currentRoom != nil {
		offsetX, offsetY := ris.camera.GetOffset()
		currentRoom.DrawWithCamera(screen, offsetX, offsetY)
	}

	engine.LogDebug("DRAW_LAYER: Player")
	if ris.player != nil {
		offsetX, offsetY := ris.camera.GetOffset()
		ris.player.DrawWithCamera(screen, offsetX, offsetY)
		if engine.GameConfig.ShowDebugOverlay {
			ris.player.DrawDebug(screen, offsetX, offsetY)
		}
	}

	engine.LogDebug(fmt.Sprintf("DRAW_LAYER: Enemies (%d)", len(ris.enemies)))
	// Draw enemies with camera offset
	offsetX, offsetY := ris.camera.GetOffset()
	for _, enemy := range ris.enemies {
		enemy.DrawWithCamera(screen, offsetX, offsetY)
		if engine.GameConfig.ShowDebugOverlay {
			enemy.DrawDebug(screen, offsetX, offsetY)
		}
	}

	// Grid overlay is rendered by the room implementation (camera-relative).

	// Draw viewport frame/borders for small rooms
	if ris.viewportRenderer != nil {
		// Ensure viewport renderer has up-to-date world bounds based on current room
		if currentRoom != nil {
			if tm := currentRoom.GetTileMap(); tm != nil {
				u := engine.GetPhysicsUnit()
				// Keep viewport frame in unscaled physics world since it draws in screen pixels but uses offset in scaled space
				ris.viewportRenderer.SetWorldBounds(tm.Width*u, tm.Height*u)
			}
		}
		offsetX, offsetY := ris.camera.GetOffset()
		// Camera offset is already in screen pixels; pass through
		ris.viewportRenderer.SetOffset(offsetX, offsetY)
		ris.viewportRenderer.DrawFrame(screen)
	}

	engine.LogDebug("DRAW_LAYER: HUD")
	// Draw HUD elements on top
	ris.hudManager.Draw(screen)

	engine.LogDebug("DRAW_LAYER: DebugText")
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
				u := engine.GetPhysicsUnit()
				scaledWorldWidth := int(float64(tileMap.Width*u) * engine.GameConfig.TileScaleFactor)
				scaledWorldHeight := int(float64(tileMap.Height*u) * engine.GameConfig.TileScaleFactor)
				ris.camera.SetWorldBounds(scaledWorldWidth, scaledWorldHeight)
				ris.viewportRenderer.SetWorldBounds(scaledWorldWidth, scaledWorldHeight)

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
