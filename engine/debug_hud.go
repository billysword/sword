package engine

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DebugHUD handles debug information display
// Implements HUDComponent interface
type DebugHUD struct {
	visible   bool
	name      string
	debugInfo DebugInfo
}

// DebugInfo contains debug information to display
type DebugInfo struct {
	RoomInfo     string
	PlayerPos    string
	CameraPos    string
	WindowSize   string
	CustomInfo   map[string]string // For additional debug data
}

// NewDebugHUD creates a new debug HUD component
func NewDebugHUD() *DebugHUD {
	return &DebugHUD{
		visible: true,
		name:    "debug_hud",
		debugInfo: DebugInfo{
			CustomInfo: make(map[string]string),
		},
	}
}

// GetName returns the component name (required by HUDComponent interface)
func (dh *DebugHUD) GetName() string {
	return dh.name
}

// SetVisible toggles debug HUD visibility (required by HUDComponent interface)
func (dh *DebugHUD) SetVisible(visible bool) {
	dh.visible = visible
}

// IsVisible returns whether the debug HUD is currently visible (required by HUDComponent interface)
func (dh *DebugHUD) IsVisible() bool {
	return dh.visible
}

// Update handles debug HUD logic updates (required by HUDComponent interface)
func (dh *DebugHUD) Update() error {
	// Update debug info that can be calculated independently
	// Window size and performance info
	windowWidth, windowHeight := ebiten.WindowSize()
	dh.debugInfo.WindowSize = fmt.Sprintf("Window: %dx%d | TPS: %.2f", windowWidth, windowHeight, ebiten.ActualTPS())
	
	// Update debug toggle states
	backgroundStatus := "ON"
	if !GetBackgroundVisible() {
		backgroundStatus = "OFF"
	}
	dh.SetCustomInfo("Background", backgroundStatus)
	
	gridStatus := "OFF"
	if GetGridVisible() {
		gridStatus = "ON"
	}
	dh.SetCustomInfo("Grid", gridStatus)
	
	depthStatus := "OFF"
	if GameConfig.EnableDepthOfField {
		depthStatus = "ON"
	}
	dh.SetCustomInfo("Depth of Field", depthStatus)
	
	return nil
}

// Draw renders the debug HUD (required by HUDComponent interface)
func (dh *DebugHUD) Draw(screen interface{}) error {
	if !dh.visible {
		return nil
	}
	
	ebitenScreen, ok := screen.(*ebiten.Image)
	if !ok {
		return nil // Skip if not Ebiten screen
	}
	
	// Get physics unit for conversions
	physicsUnit := GetPhysicsUnit()
	
	// Build comprehensive debug text
	y := 10
	lineHeight := 15
	
	// Performance section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== PERFORMANCE ===", 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("FPS: %.1f | TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()), 10, y)
	y += lineHeight * 2
	
	// Player section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== PLAYER ===", 10, y)
	y += lineHeight
	if dh.debugInfo.PlayerPos != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.PlayerPos, 10, y)
		y += lineHeight
	}
	
	// Physics section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== PHYSICS ===", 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Physics Unit: %d px", physicsUnit), 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Move Speed: %d | Jump: %d", GameConfig.PlayerMoveSpeed, GameConfig.PlayerJumpPower), 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Gravity: %d | Friction: %d", GameConfig.Gravity, GameConfig.PlayerFriction), 10, y)
	y += lineHeight * 2
	
	// Rendering section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== RENDERING ===", 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Tile: %dx%.1f = %.1f px", GameConfig.TileSize, GameConfig.TileScaleFactor, float64(GameConfig.TileSize)*GameConfig.TileScaleFactor), 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Char Scale: %.2f (%.1f px)", GameConfig.CharScaleFactor, 32.0*GameConfig.CharScaleFactor), 10, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Window: %dx%d", GameConfig.WindowWidth, GameConfig.WindowHeight), 10, y)
	y += lineHeight * 2
	
	// Camera section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== CAMERA ===", 10, y)
	y += lineHeight
	if dh.debugInfo.CameraPos != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.CameraPos, 10, y)
		y += lineHeight
	}
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Smoothing: %.2f", GameConfig.CameraSmoothing), 10, y)
	y += lineHeight * 2
	
	// Room section
	ebitenutil.DebugPrintAt(ebitenScreen, "=== ROOM ===", 10, y)
	y += lineHeight
	if dh.debugInfo.RoomInfo != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Room: %s", dh.debugInfo.RoomInfo), 10, y)
		y += lineHeight
	}
	ebitenutil.DebugPrintAt(ebitenScreen, fmt.Sprintf("Ground Level: %d tiles", GameConfig.GroundLevel), 10, y)
	y += lineHeight * 2
	
	// Custom debug info
	for key, value := range dh.debugInfo.CustomInfo {
		debugText := fmt.Sprintf("%s: %s", key, value)
		ebitenutil.DebugPrintAt(ebitenScreen, debugText, 10, y)
		y += lineHeight
	}
	
	// Hotkey help (right side of screen)
	windowWidth, _ := ebiten.WindowSize()
	helpX := windowWidth - 300
	helpY := 10
	
	ebitenutil.DebugPrintAt(ebitenScreen, "=== HOTKEYS ===", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "F3: Toggle Debug HUD", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "F4: Toggle Debug Overlay", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "[ ]: Char Scale -/+", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "- =: Tile Scale -/+", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "1: Move Speed +/-", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "2: Jump Power +/-", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "3: Gravity +/-", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "(Hold Shift for fine-tune)", helpX, helpY)
	helpY += lineHeight * 2
	ebitenutil.DebugPrintAt(ebitenScreen, "G: Toggle Grid", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "B: Toggle Background", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "M: Toggle Mini-Map", helpX, helpY)
	helpY += lineHeight
	ebitenutil.DebugPrintAt(ebitenScreen, "ESC/P: Pause", helpX, helpY)
	
	return nil
}

// UpdateRoomInfo updates the room information
func (dh *DebugHUD) UpdateRoomInfo(roomInfo string) {
	dh.debugInfo.RoomInfo = roomInfo
}

// UpdatePlayerPos updates the player position information
func (dh *DebugHUD) UpdatePlayerPos(playerPos string) {
	dh.debugInfo.PlayerPos = playerPos
}

// UpdateCameraPos updates the camera position information
func (dh *DebugHUD) UpdateCameraPos(cameraPos string) {
	dh.debugInfo.CameraPos = cameraPos
}

// SetCustomInfo sets custom debug information
func (dh *DebugHUD) SetCustomInfo(key, value string) {
	dh.debugInfo.CustomInfo[key] = value
}

// ClearCustomInfo clears all custom debug information
func (dh *DebugHUD) ClearCustomInfo() {
	dh.debugInfo.CustomInfo = make(map[string]string)
}