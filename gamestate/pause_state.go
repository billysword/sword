package gamestate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PauseState represents the pause menu overlay
type PauseState struct {
	stateManager    *StateManager
	backgroundState State // Store the previous state to draw behind pause menu
}

// NewPauseState creates a new pause state
func NewPauseState(sm *StateManager, background State) *PauseState {
	return &PauseState{
		stateManager:    sm,
		backgroundState: background,
	}
}

// Update handles input for the pause menu
func (p *PauseState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Resume game
		p.stateManager.ChangeState(p.backgroundState)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		// Return to main menu
		p.stateManager.ChangeState(NewStartState(p.stateManager))
	}
	return nil
}

// Draw renders the pause overlay
func (p *PauseState) Draw(screen *ebiten.Image) {
	// Draw the background state (frozen game)
	if p.backgroundState != nil {
		p.backgroundState.Draw(screen)
	}
	
	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(color.RGBA{0x00, 0x00, 0x00, 0xaa}) // Semi-transparent black
	screen.DrawImage(overlay, nil)
	
	// Draw pause menu
	msg := "PAUSED\n\nP/ESC - Resume\nQ - Main Menu"
	ebitenutil.DebugPrintAt(screen, msg, 400, 250)
}

// OnEnter is called when entering pause state
func (p *PauseState) OnEnter() {
	// Pause state setup
}

// OnExit is called when leaving pause state
func (p *PauseState) OnExit() {
	// Pause state cleanup
}