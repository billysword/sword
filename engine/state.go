package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Global debug rendering settings
var (
	showBackground = true
	showGrid       = true
)

/*
GetBackgroundImage retrieves the background image from the sprite manager.
Returns nil if the background sheet is not loaded.
*/
func GetBackgroundImage() *ebiten.Image {
	sm := GetSpriteManager()
	if sheet, ok := sm.sheets["background"]; ok {
		return sheet.Image
	}
	return nil
}

/*
GetTileSprite retrieves the main tile sprite sheet from the sprite manager.
Returns nil if the sheet is not loaded.
*/
func GetTileSprite() *ebiten.Image {
	sm := GetSpriteManager()
	if sheet, ok := sm.sheets["forest"]; ok {
		return sheet.Image
	}
	return nil
}

/*
ToggleBackground toggles background rendering on/off.
Useful for debugging and performance testing by removing
background rendering overhead.
*/
func ToggleBackground() {
	showBackground = !showBackground
}

/*
GetBackgroundVisible returns whether background is visible.
Returns true if backgrounds should be rendered, false otherwise.
*/
func GetBackgroundVisible() bool {
	return showBackground
}

/*
ToggleGrid toggles grid overlay on/off.
The debug grid helps visualize tile boundaries and positioning
during development and debugging.
*/
func ToggleGrid() {
	showGrid = !showGrid
}

/*
GetGridVisible returns whether grid is visible.
Returns true if the debug grid overlay should be rendered.
*/
func GetGridVisible() bool {
	return showGrid
}

/*
IsGridEnabled returns whether grid is enabled.
Alias for GetGridVisible for consistency with settings interface.
*/
func IsGridEnabled() bool {
	return showGrid
}

/*
EnableGrid enables the debug grid overlay.
*/
func EnableGrid() {
	showGrid = true
}

/*
DisableGrid disables the debug grid overlay.
*/
func DisableGrid() {
	showGrid = false
}

/*
DrawGrid renders a faint grid overlay for debugging tile positions.
The grid uses the current GameConfig settings for color and spacing.
This version draws a static grid that doesn't move with the camera.

Parameters:
  - screen: The target screen/image to draw the grid on
*/
func DrawGrid(screen *ebiten.Image) {
	if !showGrid {
		return
	}

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	gridColor := color.RGBA{
		GameConfig.GridColor[0],
		GameConfig.GridColor[1],
		GameConfig.GridColor[2],
		GameConfig.GridColor[3],
	}

	physicsUnit := GetPhysicsUnit()

	// Draw vertical lines
	for x := 0; x < screenWidth; x += physicsUnit {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screenHeight), 1, gridColor, false)
	}

	// Draw horizontal lines
	for y := 0; y < screenHeight; y += physicsUnit {
		vector.StrokeLine(screen, 0, float32(y), float32(screenWidth), float32(y), 1, gridColor, false)
	}
}

/*
DrawGridWithCamera renders a faint grid overlay that moves with the camera.
This version of the grid rendering adjusts for camera movement, ensuring
grid lines always align with tile boundaries regardless of camera position.

The grid is essential for debugging collision detection in the center of rooms.
Each grid cell represents exactly one tile (16x16 pixels by default), making it
easy to see:
- Where collision boundaries should be
- How player position aligns with tiles
- Whether collision detection is working properly at tile boundaries
- Room layout and tile positioning accuracy

To use for collision debugging:
1. Enable grid overlay in settings or call engine.EnableGrid()
2. The grid moves with the camera to always show tile boundaries
3. Check that player collision box aligns with grid when colliding
4. Verify that collision tiles match the visual grid layout

Parameters:
  - screen: The target screen/image to draw the grid on
  - cameraOffsetX: Horizontal camera offset in pixels
  - cameraOffsetY: Vertical camera offset in pixels
*/
func DrawGridWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	if !showGrid {
		return
	}

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	gridColor := color.RGBA{
		GameConfig.GridColor[0],
		GameConfig.GridColor[1],
		GameConfig.GridColor[2],
		GameConfig.GridColor[3],
	}

	physicsUnit := GetPhysicsUnit()

	// Calculate grid offset to ensure grid lines align with tiles
	gridOffsetX := int(cameraOffsetX) % physicsUnit
	gridOffsetY := int(cameraOffsetY) % physicsUnit

	// Draw vertical lines
	for x := gridOffsetX; x < screenWidth+physicsUnit; x += physicsUnit {
		if x >= 0 {
			vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screenHeight), 1, gridColor, false)
		}
	}

	// Draw horizontal lines
	for y := gridOffsetY; y < screenHeight+physicsUnit; y += physicsUnit {
		if y >= 0 {
			vector.StrokeLine(screen, 0, float32(y), float32(screenWidth), float32(y), 1, gridColor, false)
		}
	}
}

/*
State represents a game state that can handle input, update logic, and rendering.
This interface defines the core contract that all game states must implement
for proper integration with the state management system.

All states should implement:
  - Update(): Handle input processing and game logic per frame
  - Draw(): Render the state's visual elements to screen
  - OnEnter(): Initialize state when transitioning into it
  - OnExit(): Cleanup when transitioning out of the state
*/
type State interface {
	Update() error
	Draw(screen *ebiten.Image)
	OnEnter()
	OnExit()
}

/*
StateManager manages game states and transitions.
Provides a centralized system for handling state changes, ensuring proper
cleanup of old states and initialization of new states. Supports deferred
state transitions to avoid mid-update state changes.
*/
type StateManager struct {
	currentState State
	nextState    State
}

/*
NewStateManager creates a new state manager.
Returns an initialized StateManager ready to handle state transitions.
The manager starts with no active state.
*/
func NewStateManager() *StateManager {
	return &StateManager{}
}

/*
ChangeState transitions to a new state.
Queues a state transition that will be processed during the next Update() call.
This deferred approach prevents issues with changing states mid-update.

Parameters:
  - newState: The state to transition to
*/
func (sm *StateManager) ChangeState(newState State) {
	sm.nextState = newState
}

/*
Update processes state transitions and updates the current state.
Handles any queued state transitions first (calling OnExit/OnEnter as needed),
then updates the current active state. Should be called once per frame.

Returns any error from the current state's Update() method.
*/
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

/*
Draw renders the current state.
Delegates rendering to the currently active state's Draw() method.
Should be called once per frame after Update().

Parameters:
  - screen: The target screen/image to render to
*/
func (sm *StateManager) Draw(screen *ebiten.Image) {
	if sm.currentState != nil {
		sm.currentState.Draw(screen)
	}
}

/*
GetCurrentState returns the current state (useful for debugging).
Provides access to the currently active state for inspection or
debugging purposes. Returns nil if no state is active.
*/
func (sm *StateManager) GetCurrentState() State {
	return sm.currentState
}
