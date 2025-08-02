package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

/*
InGameState represents the actual gameplay state.
Manages the core game loop including player input, physics updates,
camera movement, and rendering. Coordinates between the player,
current room, and camera systems.

Key responsibilities:
  - Player input processing and movement
  - Camera following and world rendering
  - Room management and collision handling
  - Debug feature toggles (background, grid)
  - Pause state transitions
*/
type InGameState struct {
	stateManager *engine.StateManager // Reference to state manager for transitions
	player       *entities.Player     // The player character instance
	enemies      []entities.Enemy     // All enemies in the current room (interface slice)
	currentRoom  world.Room           // Current room/level being played
	camera       *engine.Camera       // Camera for world scrolling
}

/*
NewInGameState creates a new in-game state.
Initializes all core game systems including the player, camera, and room.
Sets up the initial game world with proper camera bounds and player positioning.

Parameters:
  - sm: StateManager instance for handling state transitions

Returns a pointer to the new InGameState instance.
*/
func NewInGameState(sm *engine.StateManager) *InGameState {
	// Get the actual window size for camera viewport
	windowWidth, windowHeight := ebiten.WindowSize()

	physicsUnit := engine.GetPhysicsUnit()
	groundY := engine.GameConfig.GroundLevel * physicsUnit

	return &InGameState{
		stateManager: sm,
		player:       entities.NewPlayer(50*physicsUnit, groundY),
		enemies:      make([]entities.Enemy, 0), // Initialize empty enemies slice
		currentRoom:  world.NewSimpleRoom("main"),
		camera:       engine.NewCamera(windowWidth, windowHeight),
	}
}

/*
Update handles game logic and input.
Processes all game systems in the correct order: input handling,
player physics, camera updates, and room logic. Also handles
debug toggles and pause state transitions.

Input handling:
  - P/ESC: Pause game
  - B: Toggle background rendering
  - G: Toggle debug grid overlay

Returns any error from game systems, or ebiten.Termination to quit.
*/
func (ig *InGameState) Update() error {
	// Get room name for logging
	roomName := ""
	if ig.currentRoom != nil {
		roomName = ig.currentRoom.GetZoneID()
	}

	// Get player position for logging
	playerX, playerY := ig.player.GetPosition()

	// Check for pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		engine.LogPlayerInput("P/ESC (Pause)", playerX, playerY, roomName)
		ig.stateManager.ChangeState(NewPauseState(ig.stateManager, ig))
		return nil
	}

	// Check for quit (Alt+F4 style quit)
	if (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyF4)) {
		engine.LogPlayerInput("Alt+F4 (Quit)", playerX, playerY, roomName)
		return ebiten.Termination
	}

	// Debug toggle keys
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		engine.LogPlayerInput("B (Toggle Background)", playerX, playerY, roomName)
		engine.ToggleBackground()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		engine.LogPlayerInput("G (Toggle Grid)", playerX, playerY, roomName)
		engine.ToggleGrid()
	}

	ig.player.HandleInputWithLogging(roomName)
	ig.player.Update()

	// Update all enemies
	for _, enemy := range ig.enemies {
		enemy.Update()
	}

	// Update camera to follow player
	if ig.camera != nil && ig.player != nil {
		px, py := ig.player.GetPosition()
		// Convert physics units to pixels for camera
		ig.camera.Update(px/engine.GetPhysicsUnit(), py/engine.GetPhysicsUnit())
	}

	// Let the current room handle its own logic
	if ig.currentRoom != nil {
		if err := ig.currentRoom.Update(ig.player); err != nil {
			return err
		}
		// Let the room handle collisions
		ig.currentRoom.HandleCollisions(ig.player)
	}

	return nil
}

/*
Draw renders the game world.
Handles all rendering with proper camera transformation including
background, room tiles, player character, and debug information.
Uses camera offset to create scrolling world effect.

Rendering order:
 1. Room background and tiles (with camera offset)
 2. Player character (with camera offset)
 3. Debug grid overlay (if enabled)
 4. UI and debug information (no camera offset)

Parameters:
  - screen: The target screen/image to render to
*/
func (ig *InGameState) Draw(screen *ebiten.Image) {
	// Apply camera transformation
	cameraOffsetX, cameraOffsetY := float64(0), float64(0)
	if ig.camera != nil {
		cameraOffsetX, cameraOffsetY = ig.camera.GetOffset()
	}

	// Let the current room draw itself with camera offset
	if ig.currentRoom != nil {
		ig.currentRoom.DrawWithCamera(screen, cameraOffsetX, cameraOffsetY)
	} else {
		// Fallback: draw background if no room
		if engine.GetBackgroundImage() != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.5, 0.5)
			op.GeoM.Translate(cameraOffsetX, cameraOffsetY)
			screen.DrawImage(engine.GetBackgroundImage(), op)
		}
	}

	// Draw all enemies with camera offset
	for _, enemy := range ig.enemies {
		enemy.DrawWithCamera(screen, cameraOffsetX, cameraOffsetY)
	}

	// Draw player on top of room and enemies with camera offset
	if ig.player != nil {
		ig.player.DrawWithCamera(screen, cameraOffsetX, cameraOffsetY)
	}

	// Show debug info (HUD elements don't move with camera)
	roomInfo := "No Room"
	if ig.currentRoom != nil {
		roomInfo = ig.currentRoom.GetZoneID()
	}

	backgroundStatus := "ON"
	if !engine.GetBackgroundVisible() {
		backgroundStatus = "OFF"
	}

	gridStatus := "OFF"
	if engine.GetGridVisible() {
		gridStatus = "ON"
	}

	// Add camera position to debug info
	camX, camY := float64(0), float64(0)
	if ig.camera != nil {
		camX, camY = ig.camera.GetPosition()
	}

	msg := fmt.Sprintf("TPS: %0.2f\nRoom: %s\nCamera: (%.0f, %.0f)\nEnemies: %d\nPress SPACE to jump\nR - Switch Room\nP/ESC - Pause\nB - Background: %s\nG - Grid: %s",
		ebiten.ActualTPS(), roomInfo, camX, camY, len(ig.enemies), backgroundStatus, gridStatus)
	ebitenutil.DebugPrint(screen, msg)

	// Add placeholder text for dead zone and HUD areas
	ig.drawPlaceholderText(screen)
}

/*
drawPlaceholderText renders placeholder text for dead zone and HUD areas.
This method visualizes the camera dead zone boundaries and placeholder HUD areas
to help with UI design and debugging the new larger window layout.
*/
func (ig *InGameState) drawPlaceholderText(screen *ebiten.Image) {
	if ig.camera == nil {
		return
	}

	// Get screen dimensions
	screenWidth, screenHeight := ebiten.WindowSize()
	
	// Get camera dead zone dimensions
	deadZoneX, deadZoneY := ig.camera.GetDeadZone()
	
	// Calculate dead zone boundaries on screen
	centerX := screenWidth / 2
	centerY := screenHeight / 2
	
	// Dead zone bounds
	deadZoneLeft := centerX - deadZoneX
	deadZoneRight := centerX + deadZoneX
	deadZoneTop := centerY - deadZoneY
	deadZoneBottom := centerY + deadZoneY
	
	// Draw dead zone corner markers and labels
	ebitenutil.DebugPrintAt(screen, "DEAD ZONE", deadZoneLeft+5, deadZoneTop+5)
	ebitenutil.DebugPrintAt(screen, "↑", centerX, deadZoneTop-20)
	ebitenutil.DebugPrintAt(screen, "↓", centerX, deadZoneBottom+10)
	ebitenutil.DebugPrintAt(screen, "←", deadZoneLeft-20, centerY)
	ebitenutil.DebugPrintAt(screen, "→", deadZoneRight+10, centerY)
	
	// HUD area placeholders
	// Top HUD area
	ebitenutil.DebugPrintAt(screen, "HUD: Health, Items, Minimap", 20, 20)
	
	// Left HUD area  
	ebitenutil.DebugPrintAt(screen, "HUD:", 20, screenHeight/2-60)
	ebitenutil.DebugPrintAt(screen, "Inventory", 20, screenHeight/2-40)
	ebitenutil.DebugPrintAt(screen, "Skills", 20, screenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "Hotkeys", 20, screenHeight/2)
	
	// Right HUD area
	ebitenutil.DebugPrintAt(screen, "HUD:", screenWidth-120, screenHeight/2-60)
	ebitenutil.DebugPrintAt(screen, "Map", screenWidth-120, screenHeight/2-40)
	ebitenutil.DebugPrintAt(screen, "Objectives", screenWidth-120, screenHeight/2-20)
	ebitenutil.DebugPrintAt(screen, "Status", screenWidth-120, screenHeight/2)
	
	// Bottom HUD area
	ebitenutil.DebugPrintAt(screen, "HUD: Chat, Action Bar, System Messages", 20, screenHeight-40)
	
	// Show screen dimensions info
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Screen: %dx%d", screenWidth, screenHeight), screenWidth-150, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Dead Zone: %dx%d", deadZoneX*2, deadZoneY*2), screenWidth-150, 40)
}

/*
OnEnter is called when entering the game state.
Initializes or resets all game systems when starting gameplay.
Sets up camera bounds, player positioning, and room state.
Called when transitioning from menu or resume from pause.
*/
func (ig *InGameState) OnEnter() {
	// Reset player position or load level data
	if ig.player == nil {
		physicsUnit := engine.GetPhysicsUnit()
		groundY := engine.GameConfig.GroundLevel * physicsUnit
		ig.player = entities.NewPlayer(50*physicsUnit, groundY)
	}

	// Initialize room if needed
	if ig.currentRoom == nil {
		ig.currentRoom = world.NewSimpleRoom("main")
	}

	// Spawn some test enemies if the enemies slice is empty
	if len(ig.enemies) == 0 {
		physicsUnit := engine.GetPhysicsUnit()
		groundY := engine.GameConfig.GroundLevel * physicsUnit

		// Spawn different types of enemies to demonstrate the interface system
		ig.enemies = append(ig.enemies, entities.NewSlimeEnemy(300*physicsUnit, groundY))     // Patrol behavior
		ig.enemies = append(ig.enemies, entities.NewWandererEnemy(600*physicsUnit, groundY))  // Random behavior
		ig.enemies = append(ig.enemies, entities.NewSlimeEnemy(900*physicsUnit, groundY))     // Patrol behavior
		ig.enemies = append(ig.enemies, entities.NewWandererEnemy(1200*physicsUnit, groundY)) // Random behavior
	}

	// Set up camera bounds based on room size
	if ig.camera != nil && ig.currentRoom != nil {
		tileMap := ig.currentRoom.GetTileMap()
		if tileMap != nil {
			// Convert tile dimensions to pixel dimensions
			physicsUnit := engine.GetPhysicsUnit()
			worldWidth := tileMap.Width * physicsUnit
			worldHeight := tileMap.Height * physicsUnit
			ig.camera.SetWorldBounds(worldWidth, worldHeight)

			// Center camera on player initially
			px, py := ig.player.GetPosition()
			ig.camera.CenterOn(px/physicsUnit, py/physicsUnit)
		}
	}
}

/*
GetCurrentRoom returns the current room being played.
Provides access to the current room for other states that need tile data.

Returns the current room instance.
*/
func (ig *InGameState) GetCurrentRoom() world.Room {
	return ig.currentRoom
}

/*
OnExit is called when leaving the game state.
Handles cleanup and state preservation when transitioning away
from gameplay (to pause, menu, or quit). Notifies the current
room of the exit and can save game state.
*/
func (ig *InGameState) OnExit() {
	// Let the room know we're leaving
	if ig.currentRoom != nil {
		ig.currentRoom.OnExit(ig.player)
	}
	// Save game state or cleanup resources
}

/*
AddEnemy adds a new enemy to the current room.
Creates an enemy at the specified position and adds it to the enemies slice.

Parameters:
  - enemy: The enemy instance to add (must implement Enemy interface)
*/
func (ig *InGameState) AddEnemy(enemy entities.Enemy) {
	ig.enemies = append(ig.enemies, enemy)
}

/*
AddSlimeEnemy adds a new slime enemy to the current room.
Convenience method for creating and adding slime enemies.

Parameters:
  - x: Horizontal spawn position in physics units
  - y: Vertical spawn position in physics units

Returns a pointer to the newly created slime enemy.
*/
func (ig *InGameState) AddSlimeEnemy(x, y int) *entities.SlimeEnemy {
	slime := entities.NewSlimeEnemy(x, y)
	ig.enemies = append(ig.enemies, slime)
	return slime
}

/*
AddWandererEnemy adds a new wanderer enemy to the current room.
Convenience method for creating and adding wanderer enemies.

Parameters:
  - x: Horizontal spawn position in physics units
  - y: Vertical spawn position in physics units

Returns a pointer to the newly created wanderer enemy.
*/
func (ig *InGameState) AddWandererEnemy(x, y int) *entities.WandererEnemy {
	wanderer := entities.NewWandererEnemy(x, y)
	ig.enemies = append(ig.enemies, wanderer)
	return wanderer
}

/*
RemoveEnemy removes an enemy from the current room.
Finds and removes the specified enemy from the enemies slice.

Parameters:
  - enemy: The enemy to remove (Enemy interface)

Returns true if the enemy was found and removed, false otherwise.
*/
func (ig *InGameState) RemoveEnemy(enemy entities.Enemy) bool {
	for i, e := range ig.enemies {
		if e == enemy {
			// Remove enemy by swapping with last element and truncating
			ig.enemies[i] = ig.enemies[len(ig.enemies)-1]
			ig.enemies = ig.enemies[:len(ig.enemies)-1]
			return true
		}
	}
	return false
}

/*
ClearEnemies removes all enemies from the current room.
Useful for room transitions or level resets.
*/
func (ig *InGameState) ClearEnemies() {
	ig.enemies = ig.enemies[:0] // Clear slice but keep capacity
}

/*
GetEnemies returns a copy of the current enemies slice.
Useful for external systems that need to iterate over enemies
without modifying the internal slice.

Returns a slice containing all current enemies (Enemy interface).
*/
func (ig *InGameState) GetEnemies() []entities.Enemy {
	enemies := make([]entities.Enemy, len(ig.enemies))
	copy(enemies, ig.enemies)
	return enemies
}
