package gamestate

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Global scale and size constants
const (
	// Tile system constants
	TILE_SIZE           = 16    // Base size of tiles in pixels (from tilemap)
	TILE_SCALE_FACTOR   = 1.0   // Scale factor for tile rendering (16x16 -> 16x16) - zoomed out
	
	// Character scale constants
	CHAR_SCALE_FACTOR   = 0.4   // Scale factor for character sprites (adjusted for new zoom)
	
	// Physics unit (should match tile render size after scaling)
	PHYSICS_UNIT        = int(TILE_SIZE * TILE_SCALE_FACTOR) // 16 pixels
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

// ToggleBackground toggles the background visibility
func ToggleBackground() {
	showBackground = !showBackground
}

// ToggleGrid toggles the grid visibility
func ToggleGrid() {
	showGrid = !showGrid
}

// GetBackgroundVisible returns whether background is visible
func GetBackgroundVisible() bool {
	return showBackground
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
	gridColor := color.RGBA{100, 100, 100, 80} // Faint gray
	
	// Draw vertical lines
	for x := 0; x < screenWidth; x += PHYSICS_UNIT {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screenHeight), 1, gridColor, false)
	}
	
	// Draw horizontal lines
	for y := 0; y < screenHeight; y += PHYSICS_UNIT {
		vector.StrokeLine(screen, 0, float32(y), float32(screenWidth), float32(y), 1, gridColor, false)
	}
}

// DrawGridWithCamera renders a faint grid overlay that moves with the camera
func DrawGridWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	if !showGrid {
		return
	}

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	gridColor := color.RGBA{100, 100, 100, 80} // Faint gray
	
	// Calculate grid offset to ensure grid lines align with tiles
	gridOffsetX := int(cameraOffsetX) % PHYSICS_UNIT
	gridOffsetY := int(cameraOffsetY) % PHYSICS_UNIT
	
	// Draw vertical lines
	for x := gridOffsetX; x < screenWidth+PHYSICS_UNIT; x += PHYSICS_UNIT {
		if x >= 0 {
			vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screenHeight), 1, gridColor, false)
		}
	}
	
	// Draw horizontal lines
	for y := gridOffsetY; y < screenHeight+PHYSICS_UNIT; y += PHYSICS_UNIT {
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
