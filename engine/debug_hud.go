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
	
	// Draw debug information
	y := 10
	lineHeight := 15
	
	if dh.debugInfo.RoomInfo != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.RoomInfo, 10, y)
		y += lineHeight
	}
	
	if dh.debugInfo.PlayerPos != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.PlayerPos, 10, y)
		y += lineHeight
	}
	
	if dh.debugInfo.CameraPos != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.CameraPos, 10, y)
		y += lineHeight
	}
	
	if dh.debugInfo.WindowSize != "" {
		ebitenutil.DebugPrintAt(ebitenScreen, dh.debugInfo.WindowSize, 10, y)
		y += lineHeight
	}
	
	// Draw custom debug info
	for key, value := range dh.debugInfo.CustomInfo {
		debugText := fmt.Sprintf("%s: %s", key, value)
		ebitenutil.DebugPrintAt(ebitenScreen, debugText, 10, y)
		y += lineHeight
	}
	
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