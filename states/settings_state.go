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

// SettingsTab represents different sections in the settings menu
type SettingsTab int

const (
	TabGeneral SettingsTab = iota
	TabKeybindings
	TabDeveloper
	TabTileViewer
)

// KeyBinding represents a single keybinding configuration
type KeyBinding struct {
	Action      string
	Primary     ebiten.Key
	Secondary   ebiten.Key
	Description string
}

// DeveloperOption represents a toggleable debug option
type DeveloperOption struct {
	Name        string
	Description string
	Value       *bool
	OnToggle    func(bool)
}

/*
SettingsState represents the settings menu screen with multiple tabs.
Provides access to game configuration, keybindings, developer options,
and the tile viewer for debugging purposes.

Features:
  - Tabbed interface for organized settings
  - Keybinding configuration
  - Developer/debug options
  - Tile viewer for room inspection
  - Return to main menu or pause state
*/
type SettingsState struct {
	stateManager  *engine.StateManager // Reference to state manager for transitions
	currentTab    SettingsTab          // Currently selected tab
	selectedIndex int                  // Selected item in current tab
	scrollY       int                  // Vertical scroll position for current tab
	tileSize      int                  // Size of each tile display in pixels (for tile viewer)
	currentRoom   world.Room           // Reference to current room for tile data
	returnToPause bool                 // Whether to return to pause state or main menu
	pauseState    *PauseState          // Reference to pause state for return navigation
	
	// Keybindings configuration
	keybindings []KeyBinding
	editingKey  int // Index of keybinding being edited, -1 if none
	
	// Developer options
	developerOptions []DeveloperOption
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
	s := &SettingsState{
		stateManager:  sm,
		currentTab:    TabGeneral,
		selectedIndex: 0,
		scrollY:       0,
		tileSize:      40,
		currentRoom:   room,
		returnToPause: false,
		pauseState:    nil,
		keybindings:   make([]KeyBinding, 0),
		editingKey:    -1,
		developerOptions: make([]DeveloperOption, 0),
	}
	s.initializeKeybindings()
	s.initializeDeveloperOptions()
	return s
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
	s := &SettingsState{
		stateManager:  sm,
		currentTab:    TabGeneral,
		selectedIndex: 0,
		scrollY:       0,
		tileSize:      40,
		currentRoom:   ingameState.GetCurrentRoom(),
		returnToPause: true,
		pauseState:    pauseState,
		keybindings:   make([]KeyBinding, 0),
		editingKey:    -1,
		developerOptions: make([]DeveloperOption, 0),
	}
	s.initializeKeybindings()
	s.initializeDeveloperOptions()
	return s
}

/*
initializeKeybindings sets up the default keybindings configuration.
*/
func (s *SettingsState) initializeKeybindings() {
	s.keybindings = []KeyBinding{
		// Movement
		{Action: "Move Left", Primary: ebiten.KeyA, Secondary: ebiten.KeyArrowLeft, Description: "Move character left"},
		{Action: "Move Right", Primary: ebiten.KeyD, Secondary: ebiten.KeyArrowRight, Description: "Move character right"},
		{Action: "Jump", Primary: ebiten.KeySpace, Secondary: ebiten.KeyW, Description: "Make character jump"},
		{Action: "Move Down", Primary: ebiten.KeyS, Secondary: ebiten.KeyArrowDown, Description: "Move down/crouch"},
		
		// Menu
		{Action: "Pause", Primary: ebiten.KeyP, Secondary: ebiten.KeyEscape, Description: "Pause/unpause game"},
		{Action: "Settings", Primary: ebiten.KeyS, Secondary: -1, Description: "Open settings menu"},
		{Action: "Quit to Menu", Primary: ebiten.KeyQ, Secondary: -1, Description: "Return to main menu"},
		
		// Debug
		{Action: "Toggle Debug", Primary: ebiten.KeyF3, Secondary: -1, Description: "Toggle debug overlay"},
		{Action: "Toggle Grid", Primary: ebiten.KeyG, Secondary: -1, Description: "Toggle grid display"},
		{Action: "Toggle Background", Primary: ebiten.KeyB, Secondary: -1, Description: "Toggle background rendering"},
		{Action: "Toggle Depth", Primary: ebiten.KeyH, Secondary: -1, Description: "Toggle depth visualization"},
		
		// Camera
		{Action: "Zoom In", Primary: ebiten.KeyEqual, Secondary: -1, Description: "Increase tile scale"},
		{Action: "Zoom Out", Primary: ebiten.KeyMinus, Secondary: -1, Description: "Decrease tile scale"},
	}
}

/*
initializeDeveloperOptions sets up the developer/debug options.
*/
func (s *SettingsState) initializeDeveloperOptions() {
	s.developerOptions = []DeveloperOption{
		// Debug Display
		{
			Name:        "Show Debug Info",
			Description: "Display FPS, performance metrics, and debug information",
			Value:       &engine.GameConfig.ShowDebugInfo,
			OnToggle:    nil,
		},
		{
			Name:        "Show Debug Overlay",
			Description: "Display collision boxes, physics info, and other overlays",
			Value:       &engine.GameConfig.ShowDebugOverlay,
			OnToggle:    nil,
		},
		{
			Name:        "Show Grid",
			Description: "Display tile grid overlay",
			Value:       nil, // Will be handled specially since it's in engine.state
			OnToggle: func(enabled bool) {
				if enabled {
					engine.EnableGrid()
				} else {
					engine.DisableGrid()
				}
			},
		},
		{
			Name:        "Enable Depth of Field",
			Description: "Enable depth blur and transparency effects",
			Value:       &engine.GameConfig.EnableDepthOfField,
			OnToggle:    nil,
		},
		
		// Physics Tweaking
		{
			Name:        "Variable Jump Height",
			Description: "Allow controlling jump height by release timing",
			Value:       &engine.GameConfig.PlayerPhysics.VariableJumpHeight,
			OnToggle:    nil,
		},
		
		// Camera
		{
			Name:        "Smooth Camera",
			Description: "Enable camera smoothing (may cause motion sickness)",
			Value:       nil,
			OnToggle: func(enabled bool) {
				if enabled {
					engine.GameConfig.CameraSmoothing = 0.05
				} else {
					engine.GameConfig.CameraSmoothing = 0.0
				}
			},
		},
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
Processes tab navigation, item selection, and tab-specific controls.

Input handling:
  - Left/Right arrows or A/D: Switch between tabs
  - Up/Down arrows or W/S: Navigate items in current tab
  - Enter/Space: Select/toggle current item
  - ESC/Q: Return to previous state
  - Tab-specific controls

Returns any error from state transitions.
*/
func (s *SettingsState) Update() error {
	// Check for forced quit first (Alt+F4)
	if (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyF4)) {
		return ebiten.Termination
	}

	// Handle navigation back to previous state
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || (inpututil.IsKeyJustPressed(ebiten.KeyQ) && s.editingKey < 0) {
		if s.returnToPause && s.pauseState != nil {
			s.stateManager.ChangeState(s.pauseState)
		} else {
			s.stateManager.ChangeState(NewStartState(s.stateManager))
		}
		return nil
	}

	// Handle tab switching
	if s.editingKey < 0 { // Only allow tab switching when not editing a key
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			s.currentTab--
			if s.currentTab < TabGeneral {
				s.currentTab = TabTileViewer
			}
			s.selectedIndex = 0
			s.scrollY = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			s.currentTab++
			if s.currentTab > TabTileViewer {
				s.currentTab = TabGeneral
			}
			s.selectedIndex = 0
			s.scrollY = 0
		}
	}

	// Handle tab-specific input
	switch s.currentTab {
	case TabGeneral:
		return s.updateGeneralTab()
	case TabKeybindings:
		return s.updateKeybindingsTab()
	case TabDeveloper:
		return s.updateDeveloperTab()
	case TabTileViewer:
		return s.updateTileViewerTab()
	}

	return nil
}

func (s *SettingsState) updateGeneralTab() error {
	// General settings navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.selectedIndex--
		if s.selectedIndex < 0 {
			s.selectedIndex = 2 // Wrap to bottom
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.selectedIndex++
		if s.selectedIndex > 2 {
			s.selectedIndex = 0 // Wrap to top
		}
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch s.selectedIndex {
		case 0: // Window Mode
			// Toggle fullscreen
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		case 1: // Tile Scale
			// Handled with +/- keys
		case 2: // Character Scale
			// Handled with +/- keys
		}
	}

	// Handle tile scale adjustment
	if s.selectedIndex == 1 {
		if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
			engine.GameConfig.TileScaleFactor -= 0.1
			if engine.GameConfig.TileScaleFactor < 0.5 {
				engine.GameConfig.TileScaleFactor = 0.5
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
			engine.GameConfig.TileScaleFactor += 0.1
			if engine.GameConfig.TileScaleFactor > 3.0 {
				engine.GameConfig.TileScaleFactor = 3.0
			}
		}
	}

	// Handle character scale adjustment
	if s.selectedIndex == 2 {
		if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
			engine.GameConfig.CharScaleFactor -= 0.1
			if engine.GameConfig.CharScaleFactor < 0.5 {
				engine.GameConfig.CharScaleFactor = 0.5
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
			engine.GameConfig.CharScaleFactor += 0.1
			if engine.GameConfig.CharScaleFactor > 2.0 {
				engine.GameConfig.CharScaleFactor = 2.0
			}
		}
	}

	return nil
}

func (s *SettingsState) updateKeybindingsTab() error {
	// If editing a key, wait for input
	if s.editingKey >= 0 {
		// Check all possible keys
		for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
			if inpututil.IsKeyJustPressed(key) {
				// Don't allow ESC to be bound
				if key != ebiten.KeyEscape {
					// Determine if this is primary or secondary binding
					if s.editingKey < len(s.keybindings) {
						s.keybindings[s.editingKey].Primary = key
					} else {
						s.keybindings[s.editingKey-len(s.keybindings)].Secondary = key
					}
				}
				s.editingKey = -1
				break
			}
		}
		// Allow canceling with ESC
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.editingKey = -1
		}
		return nil
	}

	// Normal navigation
	maxIndex := len(s.keybindings)*2 - 1 // Primary and secondary for each binding
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.selectedIndex--
		if s.selectedIndex < 0 {
			s.selectedIndex = maxIndex
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.selectedIndex++
		if s.selectedIndex > maxIndex {
			s.selectedIndex = 0
		}
	}

	// Handle selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.editingKey = s.selectedIndex
	}

	// Handle clearing a binding
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if s.selectedIndex < len(s.keybindings) {
			s.keybindings[s.selectedIndex].Primary = -1
		} else {
			s.keybindings[s.selectedIndex-len(s.keybindings)].Secondary = -1
		}
	}

	return nil
}

func (s *SettingsState) updateDeveloperTab() error {
	// Navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		s.selectedIndex--
		if s.selectedIndex < 0 {
			s.selectedIndex = len(s.developerOptions) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		s.selectedIndex++
		if s.selectedIndex >= len(s.developerOptions) {
			s.selectedIndex = 0
		}
	}

	// Handle toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if s.selectedIndex < len(s.developerOptions) {
			opt := &s.developerOptions[s.selectedIndex]
			
			// Toggle the value
			newValue := false
			if opt.Value != nil {
				newValue = !(*opt.Value)
				*opt.Value = newValue
			} else if opt.OnToggle != nil {
				// For options without a direct bool pointer, we need to determine current state
				// This is a bit hacky but works for our special cases
				if opt.Name == "Show Grid" {
					newValue = !engine.IsGridEnabled()
				} else if opt.Name == "Smooth Camera" {
					newValue = engine.GameConfig.CameraSmoothing == 0.0
				}
			}
			
			// Call the toggle handler if present
			if opt.OnToggle != nil {
				opt.OnToggle(newValue)
			}
		}
	}

	return nil
}

func (s *SettingsState) updateTileViewerTab() error {
	// Handle tile debug viewer
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		s.stateManager.ChangeState(NewTileDebugState(s.stateManager, s))
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
Displays the tabbed interface with the current tab's content.

Parameters:
  - screen: The target screen/image to render to
*/
func (s *SettingsState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x11, 0x11, 0x22, 0xff})

	_, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Draw title
	title := "SETTINGS"
	ebitenutil.DebugPrintAt(screen, title, 10, 10)

	// Draw tabs
	tabY := 40
	tabSpacing := 150
	tabs := []string{"General", "Keybindings", "Developer", "Tile Viewer"}
	
	for i, tabName := range tabs {
		x := 10 + i*tabSpacing
		if SettingsTab(i) == s.currentTab {
			// Highlight current tab
			tabBg := ebiten.NewImage(tabSpacing-10, 25)
			tabBg.Fill(color.RGBA{0x44, 0x44, 0x66, 0xff})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x-5), float64(tabY-5))
			screen.DrawImage(tabBg, opts)
		}
		ebitenutil.DebugPrintAt(screen, tabName, x, tabY)
	}

	// Draw navigation help
	var navHelp string
	if s.returnToPause {
		navHelp = "ESC/Q - Back to Pause | A/D - Switch Tabs"
	} else {
		navHelp = "ESC/Q - Back to Menu | A/D - Switch Tabs"
	}
	ebitenutil.DebugPrintAt(screen, navHelp, 10, screenHeight-30)

	// Draw current tab content
	contentY := tabY + 40
	switch s.currentTab {
	case TabGeneral:
		s.drawGeneralTab(screen, contentY)
	case TabKeybindings:
		s.drawKeybindingsTab(screen, contentY)
	case TabDeveloper:
		s.drawDeveloperTab(screen, contentY)
	case TabTileViewer:
		s.drawTileViewerTab(screen, contentY)
	}
}

func (s *SettingsState) drawGeneralTab(screen *ebiten.Image, startY int) {
	y := startY
	
	// General settings options
	options := []struct {
		name  string
		value string
	}{
		{"Window Mode", fmt.Sprintf("%s (Enter to toggle)", map[bool]string{true: "Fullscreen", false: "Windowed"}[ebiten.IsFullscreen()])},
		{"Tile Scale", fmt.Sprintf("%.1fx (Use -/+ to adjust)", engine.GameConfig.TileScaleFactor)},
		{"Character Scale", fmt.Sprintf("%.1fx (Use -/+ to adjust)", engine.GameConfig.CharScaleFactor)},
	}

	for i, opt := range options {
		if i == s.selectedIndex {
			// Highlight selected option
			highlight := ebiten.NewImage(400, 20)
			highlight.Fill(color.RGBA{0x33, 0x33, 0x55, 0xff})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(10), float64(y-2))
			screen.DrawImage(highlight, opts)
		}
		
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: %s", opt.name, opt.value), 20, y)
		y += 25
	}

	// Instructions
	y += 20
	ebitenutil.DebugPrintAt(screen, "Use W/S to navigate, Enter to select/toggle", 20, y)
}

func (s *SettingsState) drawKeybindingsTab(screen *ebiten.Image, startY int) {
	y := startY
	
	if s.editingKey >= 0 {
		ebitenutil.DebugPrintAt(screen, "Press any key to bind (ESC to cancel)", 20, y)
		y += 30
	} else {
		ebitenutil.DebugPrintAt(screen, "W/S - Navigate | Enter - Edit | Delete - Clear", 20, y)
		y += 30
	}

	// Draw keybindings in two columns
	for i, kb := range s.keybindings {
		// Action name
		ebitenutil.DebugPrintAt(screen, kb.Action, 20, y)
		
		// Primary key
		primaryStr := s.getKeyName(kb.Primary)
		primaryX := 200
		if i == s.selectedIndex && s.editingKey < 0 {
			// Highlight primary binding
			highlight := ebiten.NewImage(100, 18)
			highlight.Fill(color.RGBA{0x33, 0x33, 0x55, 0xff})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(primaryX-2), float64(y-2))
			screen.DrawImage(highlight, opts)
		} else if s.editingKey == i {
			primaryStr = "[Press Key]"
		}
		ebitenutil.DebugPrintAt(screen, primaryStr, primaryX, y)
		
		// Secondary key
		secondaryStr := s.getKeyName(kb.Secondary)
		secondaryX := 320
		secondaryIndex := i + len(s.keybindings)
		if secondaryIndex == s.selectedIndex && s.editingKey < 0 {
			// Highlight secondary binding
			highlight := ebiten.NewImage(100, 18)
			highlight.Fill(color.RGBA{0x33, 0x33, 0x55, 0xff})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(secondaryX-2), float64(y-2))
			screen.DrawImage(highlight, opts)
		} else if s.editingKey == secondaryIndex {
			secondaryStr = "[Press Key]"
		}
		ebitenutil.DebugPrintAt(screen, secondaryStr, secondaryX, y)
		
		// Description
		ebitenutil.DebugPrintAt(screen, kb.Description, 450, y)
		
		y += 20
		
		// Add spacing between sections
		if i == 3 || i == 6 || i == 10 { // After movement, menu, debug sections
			y += 10
		}
	}
}

func (s *SettingsState) drawDeveloperTab(screen *ebiten.Image, startY int) {
	y := startY
	
	ebitenutil.DebugPrintAt(screen, "Developer Options - W/S to navigate, Enter to toggle", 20, y)
	y += 30

	for i, opt := range s.developerOptions {
		if i == s.selectedIndex {
			// Highlight selected option
			highlight := ebiten.NewImage(600, 20)
			highlight.Fill(color.RGBA{0x33, 0x33, 0x55, 0xff})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(20), float64(y-2))
			screen.DrawImage(highlight, opts)
		}
		
		// Get current value
		enabled := false
		if opt.Value != nil {
			enabled = *opt.Value
		} else if opt.OnToggle != nil {
			// Special cases
			if opt.Name == "Show Grid" {
				enabled = engine.IsGridEnabled()
			} else if opt.Name == "Smooth Camera" {
				enabled = engine.GameConfig.CameraSmoothing > 0.0
			}
		}
		
		status := map[bool]string{true: "[ON]", false: "[OFF]"}[enabled]
		
		// Draw option name
		ebitenutil.DebugPrintAt(screen, opt.Name, 30, y)
		
		// Draw status with color
		statusX := 250
		ebitenutil.DebugPrintAt(screen, status, statusX, y)
		
		// Draw description
		ebitenutil.DebugPrintAt(screen, opt.Description, 320, y)
		
		y += 25
	}
	
	// Additional info
	y += 20
	ebitenutil.DebugPrintAt(screen, "Note: Some options may require returning to game to see changes", 20, y)
}

func (s *SettingsState) drawTileViewerTab(screen *ebiten.Image, startY int) {
	// Instructions
	instructions := "W/S or Arrow Keys - Scroll | Page Up/Down - Fast Scroll | T - Tile Debug"
	ebitenutil.DebugPrintAt(screen, instructions, 10, startY)

	// Get tile map data
	currentRoom := s.getCurrentRoom()
	if currentRoom == nil {
		noRoomMsg := "No room loaded"
		ebitenutil.DebugPrintAt(screen, noRoomMsg, 10, startY+30)
		return
	}

	tileMap := currentRoom.GetTileMap()
	if tileMap == nil {
		noMapMsg := "No tile map available"
		ebitenutil.DebugPrintAt(screen, noMapMsg, 10, startY+30)
		return
	}

	// Room info
	roomInfo := fmt.Sprintf("Room: %s | Size: %dx%d tiles",
		currentRoom.GetZoneID(), tileMap.Width, tileMap.Height)
	ebitenutil.DebugPrintAt(screen, roomInfo, 10, startY+20)
	
	// Sprite sheet info
	sm := engine.GetSpriteManager()
	sheetsInfo := fmt.Sprintf("Loaded Sheets: %v", sm.ListSheets())
	ebitenutil.DebugPrintAt(screen, sheetsInfo, 10, startY+40)

	// Calculate display area
	displayStartY := startY + 70 - s.scrollY
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tilesPerRow := (screenWidth - 20) / s.tileSize
	if tilesPerRow < 1 {
		tilesPerRow = 1
	}
	
	// Draw tile grid with hex values and actual sprites
	for y := 0; y < tileMap.Height; y++ {
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]

			// Calculate display position
			displayX := 10 + (x%tilesPerRow)*s.tileSize
			displayY := displayStartY + ((y*tileMap.Width+x)/tilesPerRow)*(s.tileSize+10)

			// Skip if off-screen
			if displayY < startY+60 || displayY > screenHeight-150 {
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
	legendY := screenHeight - 120
	ebitenutil.DebugPrintAt(screen, "LEGEND:", 10, legendY)
	ebitenutil.DebugPrintAt(screen, "Empty (-1) - Black", 10, legendY+15)
	ebitenutil.DebugPrintAt(screen, "Ground (0x01+) - Brown", 10, legendY+30)
	ebitenutil.DebugPrintAt(screen, "Platform - Green", 10, legendY+45)
	ebitenutil.DebugPrintAt(screen, "Other - Blue", 10, legendY+60)
	ebitenutil.DebugPrintAt(screen, "Display: Hex, Decimal #, (x,y)", 10, legendY+75)
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

/*
getKeyName returns a human-readable name for a key.
*/
func (s *SettingsState) getKeyName(key ebiten.Key) string {
	if key < 0 {
		return "[None]"
	}
	// This is a simplified version - you could expand this with better names
	return key.String()
}
