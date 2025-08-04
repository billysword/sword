package states

import (
	"fmt"
	"image/color"
	"sword/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/*
TileDebugState provides a comprehensive tile viewer for debugging and development.
Displays all available tiles from the tileset using the same rendering logic as the game.
Shows tile indices, hex values, and names in an organized grid layout.
*/
type TileDebugState struct {
	stateManager *engine.StateManager
	scrollY      int
	tileSize     int
	displayScale float64
	returnState  engine.State // State to return to (settings or pause)
	selectedTile int          // Currently highlighted tile (-2 means none selected)
}

// Tile names for the forest tileset
var tileNames = map[int]string{
	-1: "EMPTY",
	0:  "DIRT",
	1:  "TOP_LEFT_CORNER",
	2:  "RIGHT_WALL_1",
	3:  "RIGHT_WALL_2",
	4:  "BOTTOM_LEFT_CORNER",
	5:  "TOP_RIGHT_CORNER",
	6:  "LEFT_WALL_1",
	7:  "CEILING_1",
	8:  "CEILING_2",
	9:  "SINGLE_TOP",
	10: "SINGLE_BOTTOM",
	11: "SINGLE_LEFT",
	12: "SINGLE_RIGHT",
	13: "FLOATING",
	14: "SINGLE_HORIZONTAL",
	15: "SINGLE_VERTICAL",
	16: "INNER_CORNER_TOP_LEFT",
	17: "INNER_CORNER_TOP_RIGHT",
	18: "INNER_CORNER_BOTTOM_RIGHT",
	19: "INNER_CORNER_BOTTOM_LEFT",
	20: "FLOOR_1",
	21: "FLOOR_2",
	22: "LEFT_WALL_2",
	23: "BOTTOM_RIGHT_CORNER",
}

/*
NewTileDebugState creates a new tile debug state.
*/
func NewTileDebugState(sm *engine.StateManager, returnState engine.State) *TileDebugState {
	return &TileDebugState{
		stateManager: sm,
		scrollY:      0,
		tileSize:     64,  // Larger tiles for better visibility
		displayScale: 2.0, // Scale factor for tile rendering
		returnState:  returnState,
		selectedTile: -2, // No selection
	}
}

/*
Update handles input for the tile debug view.
*/
func (s *TileDebugState) Update() error {
	// Check for forced quit
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		return ebiten.Termination
	}

	// Return to previous state
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		s.stateManager.ChangeState(s.returnState)
		return nil
	}

	// Handle scrolling
	scrollSpeed := 30
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.scrollY -= scrollSpeed
		if s.scrollY < 0 {
			s.scrollY = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.scrollY += scrollSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		s.scrollY -= scrollSpeed * 5
		if s.scrollY < 0 {
			s.scrollY = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		s.scrollY += scrollSpeed * 5
	}

	// Handle mouse selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.handleMouseClick()
	}

	// Zoom controls
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		s.displayScale += 0.5
		if s.displayScale > 4.0 {
			s.displayScale = 4.0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		s.displayScale -= 0.5
		if s.displayScale < 0.5 {
			s.displayScale = 0.5
		}
	}

	return nil
}

/*
Draw renders the tile debug screen.
*/
func (s *TileDebugState) Draw(screen *ebiten.Image) {
	// Dark background
	screen.Fill(color.RGBA{0x1a, 0x1a, 0x2e, 0xff})

	// Title
	title := "TILE DEBUG VIEWER"
	ebitenutil.DebugPrintAt(screen, title, 10, 10)

	// Instructions
	instructions := "ESC/Q - Back | W/S or Arrows - Scroll | +/- Zoom | Click - Select Tile"
	ebitenutil.DebugPrintAt(screen, instructions, 10, 30)

	// Sprite manager info
	sm := engine.GetSpriteManager()
	info := fmt.Sprintf("Loaded Sheets: %v | Scale: %.1fx", sm.ListSheets(), s.displayScale)
	ebitenutil.DebugPrintAt(screen, info, 10, 50)

	// Calculate grid layout
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tilesPerRow := (screenWidth - 40) / (s.tileSize + 10)
	if tilesPerRow < 1 {
		tilesPerRow = 1
	}

	startY := 80 - s.scrollY
	tileIndex := -1 // Start with empty tile

	// Draw all known tiles
	for tileIndex < 24 { // We have tiles from -1 to 23
		// Calculate position
		gridIndex := tileIndex + 1 // Adjust for -1 start
		row := gridIndex / tilesPerRow
		col := gridIndex % tilesPerRow

		x := 20 + col*(s.tileSize+10)
		y := startY + row*(s.tileSize+60) // Extra space for text

		// Skip if off-screen
		if y < -s.tileSize || y > screenHeight {
			tileIndex++
			continue
		}

		// Draw tile background
		bgColor := color.RGBA{0x2d, 0x2d, 0x3d, 0xff}
		if tileIndex == s.selectedTile {
			bgColor = color.RGBA{0x4d, 0x4d, 0x6d, 0xff}
		}

		tileRect := ebiten.NewImage(s.tileSize, s.tileSize)
		tileRect.Fill(bgColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(tileRect, opts)

		// Draw the actual tile sprite using game rendering logic
		if tileIndex >= 0 {
			sprite := engine.LoadSpriteByHex(tileIndex)
			if sprite != nil {
				spriteOpts := &ebiten.DrawImageOptions{}
				// Use the game's tile scale factor
				spriteOpts.GeoM.Scale(s.displayScale, s.displayScale)
				// Center the sprite in the display box
				spriteWidth := float64(engine.GameConfig.TileSize) * s.displayScale
				offsetX := (float64(s.tileSize) - spriteWidth) / 2
				offsetY := (float64(s.tileSize) - spriteWidth) / 2
				spriteOpts.GeoM.Translate(float64(x)+offsetX, float64(y)+offsetY)
				screen.DrawImage(sprite, spriteOpts)
			}
		} else if tileIndex == -1 {
			// Draw an X for empty tiles
			emptyText := "X"
			ebitenutil.DebugPrintAt(screen, emptyText, x+s.tileSize/2-4, y+s.tileSize/2-4)
		}

		// Draw tile info below
		name := tileNames[tileIndex]
		if name == "" {
			name = fmt.Sprintf("TILE_%d", tileIndex)
		}

		// Tile name
		ebitenutil.DebugPrintAt(screen, name, x, y+s.tileSize+2)

		// Index and hex value
		indexText := fmt.Sprintf("Index: %d", tileIndex)
		if tileIndex >= 0 {
			indexText += fmt.Sprintf(" (0x%02X)", tileIndex)
		}
		ebitenutil.DebugPrintAt(screen, indexText, x, y+s.tileSize+14)

		tileIndex++
	}

	// Draw selected tile info if any
	if s.selectedTile >= -1 && s.selectedTile < 24 {
		infoY := screenHeight - 100
		selectedInfo := fmt.Sprintf("Selected: %s (Index: %d)", tileNames[s.selectedTile], s.selectedTile)
		ebitenutil.DebugPrintAt(screen, selectedInfo, 10, infoY)

		// Additional info about the tile
		tileDesc := s.getTileDescription(s.selectedTile)
		ebitenutil.DebugPrintAt(screen, tileDesc, 10, infoY+15)
	}
}

/*
handleMouseClick processes mouse clicks to select tiles.
*/
func (s *TileDebugState) handleMouseClick() {
	mx, my := ebiten.CursorPosition()

	screenWidth, _ := ebiten.WindowSize()
	tilesPerRow := (screenWidth - 40) / (s.tileSize + 10)
	if tilesPerRow < 1 {
		tilesPerRow = 1
	}

	startY := 80 - s.scrollY

	// Check each tile position
	for tileIndex := -1; tileIndex < 24; tileIndex++ {
		gridIndex := tileIndex + 1
		row := gridIndex / tilesPerRow
		col := gridIndex % tilesPerRow

		x := 20 + col*(s.tileSize+10)
		y := startY + row*(s.tileSize+60)

		// Check if mouse is within this tile
		if mx >= x && mx <= x+s.tileSize && my >= y && my <= y+s.tileSize {
			s.selectedTile = tileIndex
			break
		}
	}
}

/*
getTileDescription returns a description of what the tile is used for.
*/
func (s *TileDebugState) getTileDescription(tileIndex int) string {
	descriptions := map[int]string{
		-1: "Empty space - no collision",
		0:  "Basic dirt/ground tile - solid collision",
		1:  "Top-left corner of platform",
		2:  "Right wall variant 1",
		3:  "Right wall variant 2",
		4:  "Bottom-left corner",
		5:  "Top-right corner",
		6:  "Left wall variant 1",
		7:  "Ceiling variant 1",
		8:  "Ceiling variant 2",
		9:  "Single platform top edge",
		10: "Single platform bottom edge",
		11: "Single platform left edge",
		12: "Single platform right edge",
		13: "Floating single tile",
		14: "Horizontal platform piece",
		15: "Vertical platform piece",
		16: "Inner corner top-left",
		17: "Inner corner top-right",
		18: "Inner corner bottom-right",
		19: "Inner corner bottom-left",
		20: "Floor variant 1",
		21: "Floor variant 2",
		22: "Left wall variant 2",
		23: "Bottom-right corner",
	}

	if desc, ok := descriptions[tileIndex]; ok {
		return desc
	}
	return "Unknown tile type"
}

/*
OnEnter is called when entering the tile debug state.
*/
func (s *TileDebugState) OnEnter() {
	s.scrollY = 0
	s.selectedTile = -2
}

/*
OnExit is called when leaving the tile debug state.
*/
func (s *TileDebugState) OnExit() {
	// Cleanup if needed
}
