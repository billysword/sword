package gamestate

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	groundY = 380
)

// InGameState represents the actual gameplay state
type InGameState struct {
	stateManager *StateManager
	player       *Player
	currentRoom  Room
	camera       *Camera
}

// NewInGameState creates a new in-game state
func NewInGameState(sm *StateManager) *InGameState {
	// Get the actual window size for camera viewport
	windowWidth, windowHeight := ebiten.WindowSize()
	
	return &InGameState{
		stateManager: sm,
		player:       NewPlayer(50*PHYSICS_UNIT, groundY*PHYSICS_UNIT),
		currentRoom:  NewSimpleRoom("main"),
		camera:       NewCamera(windowWidth, windowHeight),
	}
}

// Update handles game logic and input
func (ig *InGameState) Update() error {
	// Check for pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ig.stateManager.ChangeState(NewPauseState(ig.stateManager, ig))
		return nil
	}

	// Debug toggle keys
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		ToggleBackground()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		ToggleGrid()
	}

	ig.player.HandleInput()
	ig.player.Update()

	// Update camera to follow player
	if ig.camera != nil && ig.player != nil {
		px, py := ig.player.GetPosition()
		// Convert physics units to pixels for camera
		ig.camera.Update(px/PHYSICS_UNIT, py/PHYSICS_UNIT)
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

// Draw renders the game world
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
		if globalBackgroundImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.5, 0.5)
			op.GeoM.Translate(cameraOffsetX, cameraOffsetY)
			screen.DrawImage(globalBackgroundImage, op)
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
	if !GetBackgroundVisible() {
		backgroundStatus = "OFF"
	}
	
	gridStatus := "OFF"
	if GetGridVisible() {
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

// OnEnter is called when entering the game state
func (ig *InGameState) OnEnter() {
	// Reset player position or load level data
	if ig.player == nil {
		ig.player = NewPlayer(50*PHYSICS_UNIT, groundY*PHYSICS_UNIT)
	}

	// Initialize room if needed
	if ig.currentRoom == nil {
		ig.currentRoom = NewSimpleRoom("main")
	}

	// Set up camera bounds based on room size
	if ig.camera != nil && ig.currentRoom != nil {
		tileMap := ig.currentRoom.GetTileMap()
		if tileMap != nil {
			// Convert tile dimensions to pixel dimensions
			worldWidth := tileMap.Width * PHYSICS_UNIT
			worldHeight := tileMap.Height * PHYSICS_UNIT
			ig.camera.SetWorldBounds(worldWidth, worldHeight)
			
			// Center camera on player initially
			px, py := ig.player.GetPosition()
			ig.camera.CenterOn(px/PHYSICS_UNIT, py/PHYSICS_UNIT)
		}
	}

	// Let the room know we're entering
	ig.currentRoom.OnEnter(ig.player)
}

// OnExit is called when leaving the game state
func (ig *InGameState) OnExit() {
	// Let the room know we're leaving
	if ig.currentRoom != nil {
		ig.currentRoom.OnExit(ig.player)
	}
	// Save game state or cleanup resources
}

