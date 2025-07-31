package gamestate

import "github.com/hajimehoshi/ebiten/v2"

// Global sprite storage
var (
	globalLeftSprite      *ebiten.Image
	globalRightSprite     *ebiten.Image
	globalIdleSprite      *ebiten.Image
	globalBackgroundImage *ebiten.Image
)

// SetGlobalSprites sets the global sprite references for use by all states
func SetGlobalSprites(left, right, idle, background *ebiten.Image) {
	globalLeftSprite = left
	globalRightSprite = right
	globalIdleSprite = idle
	globalBackgroundImage = background
}

// State represents a game state that can handle input, update logic, and rendering
type State interface {
	Update() error
	Draw(screen *ebiten.Image)
	OnEnter()
	OnExit()
}

// StateManager manages the current game state and transitions
type StateManager struct {
	currentState State
	nextState    State
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{}
}

// ChangeState queues a state change for the next update cycle
func (sm *StateManager) ChangeState(newState State) {
	sm.nextState = newState
}

// Update handles state transitions and updates the current state
func (sm *StateManager) Update() error {
	// Handle state transition
	if sm.nextState != nil {
		if sm.currentState != nil {
			sm.currentState.OnExit()
		}
		sm.currentState = sm.nextState
		sm.currentState.OnEnter()
		sm.nextState = nil
	}

	// Update current state
	if sm.currentState != nil {
		return sm.currentState.Update()
	}
	return nil
}

// Draw renders the current state
func (sm *StateManager) Draw(screen *ebiten.Image) {
	if sm.currentState != nil {
		sm.currentState.Draw(screen)
	}
}

// GetCurrentState returns the current state (useful for debugging)
func (sm *StateManager) GetCurrentState() State {
	return sm.currentState
}