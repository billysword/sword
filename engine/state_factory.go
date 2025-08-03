package engine

import (
	"fmt"
)

/*
StateType represents the different types of game states.
Used by the StateFactory to identify which state to create.
*/
type StateType int

const (
	StateTypeStart StateType = iota
	StateTypeInGame
	StateTypePause
)

/*
StateConfig holds configuration for creating states.
Allows customization of state initialization without
modifying individual state constructors.
*/
type StateConfig struct {
	// General settings
	StateType StateType
	
	// InGame state specific settings
	PlayerStartX      int
	PlayerStartY      int
	StartingRoom      string
	EnableEnemies     bool
	EnemyCount        int
	
	// Pause state specific settings
	BackgroundState   State
	
	// Additional metadata
	SaveData          map[string]interface{}
	CustomProperties  map[string]interface{}
}

/*
State constructor function types for avoiding circular imports.
These allow the states package to register its constructors with the engine.
*/
type StartStateConstructor func(*StateManager) State
type InGameStateConstructor func(*StateManager) State  
type PauseStateConstructor func(*StateManager, State) State

/*
StateFactory provides centralized state creation and management.
Handles state initialization with proper configuration and
provides a clean interface for state transitions.
*/
type StateFactory struct {
	stateManager *StateManager
	configs      map[StateType]StateConfig
	
	// Constructor function registry to avoid circular imports
	startStateConstructor   StartStateConstructor
	inGameStateConstructor  InGameStateConstructor
	pauseStateConstructor   PauseStateConstructor
}

/*
NewStateFactory creates a new state factory.
Initializes the factory with default configurations for all state types.

Parameters:
  - sm: StateManager instance for handling state transitions

Returns a pointer to the new StateFactory instance.
*/
func NewStateFactory(sm *StateManager) *StateFactory {
	factory := &StateFactory{
		stateManager: sm,
		configs:      make(map[StateType]StateConfig),
	}
	
	// Initialize default configurations
	factory.initDefaultConfigs()
	
	return factory
}

/*
RegisterStartStateConstructor registers the constructor function for start states.
This avoids circular imports between engine and states packages.
*/
func (sf *StateFactory) RegisterStartStateConstructor(constructor StartStateConstructor) {
	sf.startStateConstructor = constructor
}

/*
RegisterInGameStateConstructor registers the constructor function for in-game states.
This avoids circular imports between engine and states packages.
*/
func (sf *StateFactory) RegisterInGameStateConstructor(constructor InGameStateConstructor) {
	sf.inGameStateConstructor = constructor
}

/*
RegisterPauseStateConstructor registers the constructor function for pause states.
This avoids circular imports between engine and states packages.
*/
func (sf *StateFactory) RegisterPauseStateConstructor(constructor PauseStateConstructor) {
	sf.pauseStateConstructor = constructor
}

/*
initDefaultConfigs sets up default state configurations.
These can be overridden later using SetStateConfig.
*/
func (sf *StateFactory) initDefaultConfigs() {
	physicsUnit := GetPhysicsUnit()
	groundY := GameConfig.GroundLevel * physicsUnit
	
	// Default start state config
	sf.configs[StateTypeStart] = StateConfig{
		StateType: StateTypeStart,
	}
	
	// Default ingame state config
	sf.configs[StateTypeInGame] = StateConfig{
		StateType:     StateTypeInGame,
		PlayerStartX:  50 * physicsUnit,
		PlayerStartY:  groundY,
		StartingRoom:  "main",
		EnableEnemies: true,
		EnemyCount:    4,
	}
	
	// Default pause state config (will be updated when created)
	sf.configs[StateTypePause] = StateConfig{
		StateType: StateTypePause,
	}
}

/*
SetStateConfig updates the configuration for a specific state type.
Allows customization of state initialization parameters.

Parameters:
  - stateType: The type of state to configure
  - config: The new configuration to apply
*/
func (sf *StateFactory) SetStateConfig(stateType StateType, config StateConfig) {
	config.StateType = stateType // Ensure consistency
	sf.configs[stateType] = config
}

/*
GetStateConfig returns the current configuration for a state type.
Useful for inspecting or modifying existing configurations.

Parameters:
  - stateType: The type of state to get configuration for

Returns the current StateConfig for the specified type.
*/
func (sf *StateFactory) GetStateConfig(stateType StateType) StateConfig {
	if config, exists := sf.configs[stateType]; exists {
		return config
	}
	// Return empty config if not found
	return StateConfig{StateType: stateType}
}

/*
CreateState creates a new state instance based on the specified type.
Uses the current configuration for the state type to initialize
the state with appropriate parameters.

Parameters:
  - stateType: The type of state to create

Returns the created State instance and any error.
*/
func (sf *StateFactory) CreateState(stateType StateType) (State, error) {
	config, exists := sf.configs[stateType]
	if !exists {
		return nil, fmt.Errorf("no configuration found for state type %d", stateType)
	}
	
	switch stateType {
	case StateTypeStart:
		return sf.createStartState(config)
	case StateTypeInGame:
		return sf.createInGameState(config)
	case StateTypePause:
		return sf.createPauseState(config)
	default:
		return nil, fmt.Errorf("unknown state type: %d", stateType)
	}
}

/*
TransitionTo creates and transitions to a new state.
Convenience method that combines state creation and transition.

Parameters:
  - stateType: The type of state to transition to

Returns any error from state creation or transition.
*/
func (sf *StateFactory) TransitionTo(stateType StateType) error {
	state, err := sf.CreateState(stateType)
	if err != nil {
		return err
	}
	
	sf.stateManager.ChangeState(state)
	return nil
}

/*
TransitionToPause creates a pause state with the current state as background.
Special method for pause transitions that need the current state reference.

Returns any error from state creation or transition.
*/
func (sf *StateFactory) TransitionToPause() error {
	currentState := sf.stateManager.GetCurrentState()
	if currentState == nil {
		return fmt.Errorf("no current state to pause")
	}
	
	// Update pause config with current state
	pauseConfig := sf.configs[StateTypePause]
	pauseConfig.BackgroundState = currentState
	sf.configs[StateTypePause] = pauseConfig
	
	return sf.TransitionTo(StateTypePause)
}

// Internal state creation methods - use registered constructors
func (sf *StateFactory) createStartState(config StateConfig) (State, error) {
	if sf.startStateConstructor == nil {
		return nil, fmt.Errorf("start state constructor not registered")
	}
	return sf.startStateConstructor(sf.stateManager), nil
}

func (sf *StateFactory) createInGameState(config StateConfig) (State, error) {
	if sf.inGameStateConstructor == nil {
		return nil, fmt.Errorf("in-game state constructor not registered")
	}
	return sf.inGameStateConstructor(sf.stateManager), nil
}

func (sf *StateFactory) createPauseState(config StateConfig) (State, error) {
	if sf.pauseStateConstructor == nil {
		return nil, fmt.Errorf("pause state constructor not registered")
	}
	// For pause state, we need the background state from config
	if config.BackgroundState == nil {
		return nil, fmt.Errorf("pause state requires background state in config")
	}
	return sf.pauseStateConstructor(sf.stateManager, config.BackgroundState), nil
}