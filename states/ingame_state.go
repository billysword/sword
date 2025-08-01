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
	stateManager *engine.StateManager  // Reference to state manager for transitions
	player       *entities.Player        // The player character instance
	currentRoom  world.Room           // Current room/level being played
	camera       *engine.Camera        // Camera for world scrolling
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
	// Check for pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ig.stateManager.ChangeState(NewPauseState(ig.stateManager, ig))
		return nil
	}

	// Debug toggle keys
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		engine.ToggleBackground()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		engine.ToggleGrid()
	}

	ig.player.HandleInput()
	ig.player.Update()

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

	// Draw player on top of room with camera offset
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
	
	msg := fmt.Sprintf("TPS: %0.2f\nRoom: %s\nCamera: (%.0f, %.0f)\nPress SPACE to jump\nR - Switch Room\nP/ESC - Pause\nB - Background: %s\nG - Grid: %s", 
		ebiten.ActualTPS(), roomInfo, camX, camY, backgroundStatus, gridStatus)
	ebitenutil.DebugPrint(screen, msg)
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

