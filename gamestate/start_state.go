package gamestate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// StartState represents the game start/menu screen
type StartState struct {
	stateManager     *StateManager
	selectedOption   int  // 0 = Continue, 1 = Quit
	totalOptions     int
}

// NewStartState creates a new start state
func NewStartState(sm *StateManager) *StartState {
	return &StartState{
		stateManager:   sm,
		selectedOption: 0, // Default to "Continue"
		totalOptions:   2, // Continue and Quit
	}
}

// Update handles input and logic for the start screen
func (s *StartState) Update() error {
	// Handle menu navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.selectedOption = (s.selectedOption - 1 + s.totalOptions) % s.totalOptions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.selectedOption = (s.selectedOption + 1) % s.totalOptions
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch s.selectedOption {
		case 0: // Continue
			s.stateManager.ChangeState(NewInGameState(s.stateManager))
		case 1: // Quit
			return ebiten.Termination
		}
	}

	// ESC key to continue (default action)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.stateManager.ChangeState(NewInGameState(s.stateManager))
	}

	return nil
}

// Draw renders the start screen
func (s *StartState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x11, 0x22, 0xff})
	
	// Title
	title := "SWORD PLATFORMER"
	ebitenutil.DebugPrintAt(screen, title, 50, 150)
	
	// Menu options
	continueText := "Continue"
	quitText := "Quit"
	
	if s.selectedOption == 0 {
		continueText = "> " + continueText + " <"
	}
	if s.selectedOption == 1 {
		quitText = "> " + quitText + " <"
	}
	
	ebitenutil.DebugPrintAt(screen, continueText, 50, 220)
	ebitenutil.DebugPrintAt(screen, quitText, 50, 250)
	
	// Instructions
	instructions := "\nControls:\nW/S or Arrow Keys - Navigate menu\nENTER/SPACE - Select option\nESC - Continue (default)\n\nGame Controls:\nWASD/Arrow Keys - Move\nSpace - Jump"
	ebitenutil.DebugPrintAt(screen, instructions, 50, 300)
}

// OnEnter is called when entering this state
func (s *StartState) OnEnter() {
	// Initialize start state resources if needed
}

// OnExit is called when leaving this state
func (s *StartState) OnExit() {
	// Cleanup start state resources if needed
}