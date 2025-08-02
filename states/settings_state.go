package states

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"sword/world"
)

/*
SettingsState represents the settings menu screen.
Displays the currently loaded room tiles with their hex values in a scrollable grid.
Allows players to inspect the tile layout and values for debugging or modding purposes.

Features:
  - Display current room tile map with hex values
  - Scrollable view for large tile maps
  - Color-coded tiles based on type
  - Return to main menu functionality
*/
type SettingsState struct {
	stateManager  *engine.StateManager // Reference to state manager for transitions
	scrollY       int                  // Vertical scroll position
	tileSize      int                  // Size of each tile display in pixels
	currentRoom   world.Room           // Reference to current room for tile data
	returnToPause bool                 // Whether to return to pause state or main menu
	pauseState    *PauseState          // Reference to pause state for return navigation
}

/*
NewSettingsState creates a new settings state.
Initializes the settings state with default scroll position and tile display size.

Parameters:
  - sm: StateManager instance for handling state transitions
  - room: Current room to display tile data from

Returns a pointer to the new SettingsState instance.
*/
func NewSettingsState(sm *engine.StateManager, room world.Room) *SettingsState {
	return &SettingsState{
		stateManager:  sm,
		scrollY:       0,
		tileSize:      40,
		currentRoom:   room,
		returnToPause: false,
		pauseState:    nil,
	}
}

/*
NewSettingsStateFromPause creates a new settings state accessible from pause menu.
Initializes the settings state with the current room from the ingame state.

Parameters:
  - sm: StateManager instance for handling state transitions
  - ingameState: The ingame state to get current room data from
  - pauseState: The pause state to return to when exiting settings

Returns a pointer to the new SettingsState instance.
*/
func NewSettingsStateFromPause(sm *engine.StateManager, ingameState *InGameState, pauseState *PauseState) *SettingsState {
	return &SettingsState{
		stateManager:  sm,
		scrollY:       0,
		tileSize:      40,
		currentRoom:   ingameState.GetCurrentRoom(),
		returnToPause: true,
		pauseState:    pauseState,
	}
}

/*
getCurrentRoom returns the current room to display.
Prioritizes the dynamic room from ingame state if available, 
otherwise falls back to the static room reference.

Returns the current room instance, or nil if no room is available.
*/
func (s *SettingsState) getCurrentRoom() world.Room {
	return s.currentRoom
}

/*
Update handles input for the settings menu.
Processes scrolling controls and navigation back to the main menu.

Input handling:
  - Up/Down arrows or W/S: Scroll through the tile grid
  - ESC/Q: Return to main menu
  - Page Up/Page Down: Fast scroll

Returns any error from state transitions.
*/
func (s *SettingsState) Update() error {
	// Check for forced quit first (Ctrl+Q or Alt+F4)
	if (ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyQ)) ||
		(ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyF4)) {
		return ebiten.Termination
	}

	// Handle navigation back to previous state
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if s.returnToPause && s.pauseState != nil {
			s.stateManager.ChangeState(s.pauseState)
		} else {
			s.stateManager.ChangeState(NewStartState(s.stateManager))
		}
		return nil
	}

	// Handle scrolling
	scrollSpeed := 20
	fastScrollSpeed := 100

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
		s.scrollY -= fastScrollSpeed
		if s.scrollY < 0 {
			s.scrollY = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		s.scrollY += fastScrollSpeed
	}

	return nil
}

/*
Draw renders the settings screen.
Displays the tile grid with hex values, coordinates, and navigation instructions.
Uses color coding to distinguish different tile types.

Parameters:
  - screen: The target screen/image to render to
*/
func (s *SettingsState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x11, 0x11, 0x22, 0xff})

	// Title
	title := "SETTINGS - ROOM TILE VIEWER"
	ebitenutil.DebugPrintAt(screen, title, 10, 10)

	// Instructions
	var instructions string
	if s.returnToPause {
		instructions = "ESC/Q - Back to Pause | W/S or Arrow Keys - Scroll | Page Up/Down - Fast Scroll"
	} else {
		instructions = "ESC/Q - Back to Menu | W/S or Arrow Keys - Scroll | Page Up/Down - Fast Scroll"
	}
	ebitenutil.DebugPrintAt(screen, instructions, 10, 30)

	// Get tile map data
	currentRoom := s.getCurrentRoom()
	if currentRoom == nil {
		noRoomMsg := "No room loaded"
		ebitenutil.DebugPrintAt(screen, noRoomMsg, 10, 70)
		return
	}

	tileMap := currentRoom.GetTileMap()
	if tileMap == nil {
		noMapMsg := "No tile map available"
		ebitenutil.DebugPrintAt(screen, noMapMsg, 10, 70)
		return
	}

	// Room info - show if this is live data or static
	roomSource := "Live"
	roomInfo := fmt.Sprintf("Room: %s (%s) | Size: %dx%d tiles",
		currentRoom.GetZoneID(), roomSource, tileMap.Width, tileMap.Height)
	ebitenutil.DebugPrintAt(screen, roomInfo, 10, 50)
	
	// Sprite sheet info
	sm := engine.GetSpriteManager()
	sheetsInfo := fmt.Sprintf("Loaded Sheets: %v", sm.ListSheets())
	ebitenutil.DebugPrintAt(screen, sheetsInfo, 10, 70)

	// Calculate display area
	startY := 100 - s.scrollY
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tilesPerRow := (screenWidth - 20) / s.tileSize
	if tilesPerRow < 1 {
		tilesPerRow = 1
	}

	// Get sprite manager for tile rendering
	sm = engine.GetSpriteManager()
	
	// Draw tile grid with hex values and actual sprites
	for y := 0; y < tileMap.Height; y++ {
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]

			// Calculate display position
			displayX := 10 + (x%tilesPerRow)*s.tileSize
			displayY := startY + ((y*tileMap.Width+x)/tilesPerRow)*(s.tileSize+10)

			// Skip if off-screen
			if displayY < -s.tileSize || displayY > screenHeight {
				continue
			}

			// Draw tile background based on type
			tileColor := s.getTileColor(tileIndex)
			tileRect := ebiten.NewImage(s.tileSize-2, s.tileSize-2)
			tileRect.Fill(tileColor)

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(displayX), float64(displayY))
			screen.DrawImage(tileRect, opts)
			
			// Try to draw actual sprite from sprite manager
			if tileIndex >= 0 {
				sprite := engine.LoadSpriteByHex(tileIndex)
				if sprite != nil {
					spriteOpts := &ebiten.DrawImageOptions{}
					// Scale sprite to fit within the tile display area
					scale := float64(s.tileSize-4) / float64(engine.GameConfig.TileSize)
					spriteOpts.GeoM.Scale(scale, scale)
					spriteOpts.GeoM.Translate(float64(displayX+2), float64(displayY+2))
					screen.DrawImage(sprite, spriteOpts)
				}
			}

			// Draw hex value (larger, more prominent)
			hexText := s.formatTileValue(tileIndex)
			ebitenutil.DebugPrintAt(screen, hexText, displayX+3, displayY+5)

			// Draw coordinates (smaller text at bottom)
			coordText := fmt.Sprintf("(%d,%d)", x, y)
			ebitenutil.DebugPrintAt(screen, coordText, displayX+2, displayY+s.tileSize-12)
			
			// Draw tile index in decimal as well for easy reference
			if s.tileSize > 30 {
				decText := fmt.Sprintf("#%d", tileIndex)
				ebitenutil.DebugPrintAt(screen, decText, displayX+2, displayY+s.tileSize-25)
			}
		}
	}

	// Legend
	legendY := screenHeight - 140
	ebitenutil.DebugPrintAt(screen, "LEGEND:", 10, legendY)
	ebitenutil.DebugPrintAt(screen, "Empty (-1) - Black", 10, legendY+15)
	ebitenutil.DebugPrintAt(screen, "Ground (0x01+) - Brown", 10, legendY+30)
	ebitenutil.DebugPrintAt(screen, "Platform - Green", 10, legendY+45)
	ebitenutil.DebugPrintAt(screen, "Other - Blue", 10, legendY+60)
	ebitenutil.DebugPrintAt(screen, "Display: Hex, Decimal #, (x,y)", 10, legendY+75)
	ebitenutil.DebugPrintAt(screen, "Sprites loaded from sprite manager", 10, legendY+90)
}

/*
getTileColor returns the color to use for displaying a tile based on its type/value.
*/
func (s *SettingsState) getTileColor(tileIndex int) color.RGBA {
	switch {
	case tileIndex == -1:
		return color.RGBA{0x20, 0x20, 0x20, 0xff} // Dark gray for empty
	case tileIndex >= 0x01 && tileIndex <= 0x10:
		return color.RGBA{0x8B, 0x45, 0x13, 0xff} // Brown for ground tiles
	case tileIndex >= 0x11 && tileIndex <= 0x20:
		return color.RGBA{0x22, 0x8B, 0x22, 0xff} // Green for platform tiles
	default:
		return color.RGBA{0x41, 0x69, 0xE1, 0xff} // Blue for other tiles
	}
}

/*
formatTileValue formats a tile index as a hex string for display.
*/
func (s *SettingsState) formatTileValue(tileIndex int) string {
	if tileIndex == -1 {
		return "-1"
	}
	if tileIndex <= 0xF {
		return fmt.Sprintf("0x%X", tileIndex)
	}
	return fmt.Sprintf("0x%02X", tileIndex)
}

/*
OnEnter is called when entering settings state.
Performs any necessary setup when the settings state becomes active.
*/
func (s *SettingsState) OnEnter() {
	// Reset scroll position when entering
	s.scrollY = 0
}

/*
OnExit is called when leaving settings state.
Handles cleanup when transitioning away from the settings state.
*/
func (s *SettingsState) OnExit() {
	// Settings state cleanup
}
