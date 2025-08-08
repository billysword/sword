package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
		"sword/world"
	"sword/room_layouts"
 )

/*
StartState represents the game start/menu screen.
Provides the main menu interface where players can choose to start the game,
access settings, or quit. Handles menu navigation and transitions to other states.

Features:
  - Menu navigation with keyboard controls
  - Visual feedback for selected options
  - Control instructions display
  - Smooth transition to gameplay and settings
*/
type StartState struct {
	stateManager   *engine.StateManager // Reference to state manager for transitions
	selectedOption int                  // Currently selected menu option (0-based)
	totalOptions   int                  // Total number of menu options available
}

/*
NewStartState creates a new start state.
Initializes the start state with default menu selection and sets up
the available menu options. The state manager reference is required
for transitioning to other game states.

Parameters:
  - sm: StateManager instance for handling state transitions

Returns a pointer to the new StartState instance.
*/
func NewStartState(sm *engine.StateManager) *StartState {
	return &StartState{
		stateManager:   sm,
		selectedOption: 0, // Default to "Continue"
		totalOptions:   3, // Continue, Settings, and Quit
	}
}

/*
Update handles input and logic for the start screen.
Processes menu navigation input (up/down arrows, WASD) and selection
input (Enter, Space). Also handles the ESC key as a quick-start option.

Menu options:
  - 0: Continue (start game)
  - 1: Settings (view tile data)
  - 2: Quit (exit application)

Returns ebiten.Termination if quit is selected, nil otherwise.
*/
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
		case 1: // Settings
			// Create a preview room sized from curated layout
			defaultRoom := world.NewSimpleRoomFromLayout("main", room_layouts.EmptyRoom)
			s.stateManager.ChangeState(NewSettingsState(s.stateManager, defaultRoom))
		case 2: // Quit
			return ebiten.Termination
		}
	}

	// ESC key to continue (default action)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.stateManager.ChangeState(NewInGameState(s.stateManager))
	}

	return nil
}

/*
Draw renders the start screen.
Displays the game title, menu options with selection indicators,
and control instructions. Uses a dark blue background and positions
text elements for good readability.

Parameters:
  - screen: The target screen/image to render to
*/
func (s *StartState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x11, 0x22, 0xff})

	// Title
	title := "SWORD PLATFORMER"
	ebitenutil.DebugPrintAt(screen, title, 50, 150)

	// Menu options
	continueText := "Continue"
	settingsText := "Settings"
	quitText := "Quit"

	if s.selectedOption == 0 {
		continueText = "> " + continueText + " <"
	}
	if s.selectedOption == 1 {
		settingsText = "> " + settingsText + " <"
	}
	if s.selectedOption == 2 {
		quitText = "> " + quitText + " <"
	}

	ebitenutil.DebugPrintAt(screen, continueText, 50, 220)
	ebitenutil.DebugPrintAt(screen, settingsText, 50, 250)
	ebitenutil.DebugPrintAt(screen, quitText, 50, 280)

	// Instructions
	instructions := "\nControls:\nW/S or Arrow Keys - Navigate menu\nENTER/SPACE - Select option\nESC - Continue (default)\n\nGame Controls:\nWASD/Arrow Keys - Move\nSpace - Jump"
	ebitenutil.DebugPrintAt(screen, instructions, 50, 330)
}

/*
OnEnter is called when entering this state.
Performs any necessary initialization when the start state becomes active.
Currently used for potential resource setup or state reset.
*/
func (s *StartState) OnEnter() {
	// Initialize start state resources if needed
}

/*
OnExit is called when leaving this state.
Handles cleanup when transitioning away from the start state.
Currently used for potential resource cleanup or final preparations.
*/
func (s *StartState) OnExit() {
	// Cleanup start state resources if needed
}
