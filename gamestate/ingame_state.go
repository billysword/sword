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
}

// NewInGameState creates a new in-game state
func NewInGameState(sm *StateManager) *InGameState {
	return &InGameState{
		stateManager: sm,
		player:       NewPlayer(50*PHYSICS_UNIT, groundY*PHYSICS_UNIT),
		currentRoom:  NewSimpleRoom("main"),
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
	// Let the current room draw itself (includes background and tiles)
	if ig.currentRoom != nil {
		ig.currentRoom.Draw(screen)
	} else {
		// Fallback: draw background if no room
		if globalBackgroundImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.5, 0.5)
			screen.DrawImage(globalBackgroundImage, op)
		}
	}

	// Draw player on top of room
	if ig.player != nil {
		ig.player.Draw(screen)
	}

	// Show debug info
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
	
	msg := fmt.Sprintf("TPS: %0.2f\nRoom: %s\nPress SPACE to jump\nR - Switch Room\nP/ESC - Pause\nB - Background: %s\nG - Grid: %s", ebiten.ActualTPS(), roomInfo, backgroundStatus, gridStatus)
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

