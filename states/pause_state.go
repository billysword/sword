package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
)

/*
PauseState represents the pause menu overlay.
Provides a pause menu that overlays on top of the frozen game world.
Allows players to resume gameplay or return to the main menu while
preserving the current game state.

Features:
  - Semi-transparent overlay showing frozen game state
  - Simple menu with resume and quit options
  - Preserves background state for seamless resume
  - Keyboard controls for navigation
*/
type PauseState struct {
	stateManager    *engine.StateManager  // Reference to state manager for transitions
	backgroundState engine.State          // Store the previous state to draw behind pause menu
}

/*
NewPauseState creates a new pause state.
Initializes the pause state with a reference to the background state
that should continue to be displayed (but frozen) behind the pause overlay.

Parameters:
  - sm: StateManager instance for handling state transitions
  - background: The game state to display frozen behind the pause menu

Returns a pointer to the new PauseState instance.
*/
func NewPauseState(sm *engine.StateManager, background engine.State) *PauseState {
	return &PauseState{
		stateManager:    sm,
		backgroundState: background,
	}
}

/*
Update handles input for the pause menu.
Processes pause menu controls for resuming gameplay or returning
to the main menu. Does not update the background state, keeping
the game world frozen.

Input handling:
  - P/ESC: Resume game (return to background state)
  - Q: Quit to main menu

Returns any error from state transitions.
*/
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

/*
Draw renders the pause overlay.
Draws the frozen background state first, then applies a semi-transparent
overlay and renders the pause menu on top. This creates the effect of
the game being frozen behind a darkened overlay.

Parameters:
  - screen: The target screen/image to render to
*/
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

/*
OnEnter is called when entering pause state.
Performs any necessary setup when the pause state becomes active.
Currently used for potential pause-specific initialization.
*/
func (p *PauseState) OnEnter() {
	// Pause state setup
}

/*
OnExit is called when leaving pause state.
Handles cleanup when transitioning away from the pause state.
Currently used for potential pause-specific cleanup operations.
*/
func (p *PauseState) OnExit() {
	// Pause state cleanup
}
