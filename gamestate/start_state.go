package gamestate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// StartState represents the game start/menu screen
type StartState struct {
	stateManager *StateManager
}

// NewStartState creates a new start state
func NewStartState(sm *StateManager) *StartState {
	return &StartState{
		stateManager: sm,
	}
}

// Update handles input and logic for the start screen
func (s *StartState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.stateManager.ChangeState(NewInGameState(s.stateManager))
	}
	return nil
}

// Draw renders the start screen
func (s *StartState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x11, 0x22, 0xff})
	
	msg := "SWORD PLATFORMER\n\nPress ENTER or SPACE to start\n\nControls:\nWASD/Arrow Keys - Move\nSpace - Jump"
	ebitenutil.DebugPrintAt(screen, msg, 50, 200)
}

// OnEnter is called when entering this state
func (s *StartState) OnEnter() {
	// Initialize start state resources if needed
}

// OnExit is called when leaving this state
func (s *StartState) OnExit() {
	// Cleanup start state resources if needed
}