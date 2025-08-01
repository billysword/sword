package gamestate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Legacy constants for backward compatibility - use GameConfig instead
const (
	TILE_SIZE         = 16   // Deprecated: use GameConfig.TileSize
	TILE_SCALE_FACTOR = 1.0  // Deprecated: use GameConfig.TileScaleFactor
	CHAR_SCALE_FACTOR = 0.4  // Deprecated: use GameConfig.CharScaleFactor
	PHYSICS_UNIT      = 16   // Deprecated: use GetPhysicsUnit()
)

// Global sprite storage
var (
	globalLeftSprite      *ebiten.Image
	globalRightSprite     *ebiten.Image
	globalIdleSprite      *ebiten.Image
	globalBackgroundImage *ebiten.Image
	globalTileSprite      *ebiten.Image
	globalTilesSprite     *ebiten.Image
)

// Global debug rendering settings
var (
	showBackground = true
	showGrid      = false
)

// SetGlobalSprites sets the global sprite references for use by all states
func SetGlobalSprites(left, right, idle, background *ebiten.Image) {
	globalLeftSprite = left
	globalRightSprite = right
	globalIdleSprite = idle
	globalBackgroundImage = background
}

// SetGlobalTileSprites sets the global tile sprite references
func SetGlobalTileSprites(tile, tiles *ebiten.Image) {
	globalTileSprite = tile
	globalTilesSprite = tiles
}

// ToggleBackground toggles background rendering on/off
func ToggleBackground() {
	showBackground = !showBackground
}

// GetBackgroundVisible returns whether background is visible
func GetBackgroundVisible() bool {
	return showBackground
}

// ToggleGrid toggles grid overlay on/off
func ToggleGrid() {
	showGrid = !showGrid
}

// GetGridVisible returns whether grid is visible
func GetGridVisible() bool {
	return showGrid
}

// DrawGrid renders a faint grid overlay for debugging tile positions
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

// DrawGridWithCamera renders a faint grid overlay that moves with the camera
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

// State represents a game state that can handle input, update logic, and rendering
type State interface {
	Update() error
	Draw(screen *ebiten.Image)
	OnEnter()
	OnExit()
}

// StateManager manages game states and transitions
type StateManager struct {
	currentState State
	nextState    State
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{}
}

// ChangeState transitions to a new state
func (sm *StateManager) ChangeState(newState State) {
	sm.nextState = newState
}

// Update processes state transitions and updates the current state
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
