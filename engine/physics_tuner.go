package engine

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/*
PhysicsTuner provides runtime adjustment of physics parameters.
Allows developers to tweak physics values in real-time using keyboard shortcuts
to find the perfect feel for player movement and collision.
*/
type PhysicsTuner struct {
	active         bool
	selectedParam  int
	paramNames     []string
	adjustmentStep float64
}

// Global physics tuner instance
var globalPhysicsTuner *PhysicsTuner

/*
InitPhysicsTuner initializes the global physics tuner.
*/
func InitPhysicsTuner() {
	globalPhysicsTuner = &PhysicsTuner{
		active:         false,
		selectedParam:  0,
		adjustmentStep: 1.0,
		paramNames: []string{
			"Move Speed",
			"Jump Power",
			"Gravity",
			"Max Fall Speed",
			"Friction",
			"Air Control",
			"Collision Width",
			"Collision Height",
			"Collision Offset X",
			"Collision Offset Y",
			"Coyote Time",
			"Jump Buffer Time",
			"Char Scale Factor",
			"Tile Scale Factor",
		},
	}
}

/*
GetPhysicsTuner returns the global physics tuner instance.
*/
func GetPhysicsTuner() *PhysicsTuner {
	if globalPhysicsTuner == nil {
		InitPhysicsTuner()
	}
	return globalPhysicsTuner
}

/*
Update processes input for the physics tuner.
Should be called once per frame when the tuner is active.
*/
func (pt *PhysicsTuner) Update() {
	// Toggle tuner with F9
	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		pt.active = !pt.active
		if pt.active {
			LogDebug("Physics Tuner activated - Use arrow keys to adjust, Tab to switch parameter")
		} else {
			LogDebug("Physics Tuner deactivated")
		}
	}
	
	if !pt.active {
		return
	}
	
	// Switch parameter with Tab
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		pt.selectedParam = (pt.selectedParam + 1) % len(pt.paramNames)
	}
	
	// Adjust step size with Shift
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		pt.adjustmentStep = 0.1
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) {
		pt.adjustmentStep = 10.0
	} else {
		pt.adjustmentStep = 1.0
	}
	
	// Adjust values with arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		delta := pt.adjustmentStep
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			delta = -delta
		}
		pt.adjustParameter(delta)
	}
}

/*
adjustParameter adjusts the currently selected parameter by the given delta.
*/
func (pt *PhysicsTuner) adjustParameter(delta float64) {
	config := &GameConfig.PlayerPhysics
	
	switch pt.selectedParam {
	case 0: // Move Speed
		config.MoveSpeed = max(1, config.MoveSpeed+int(delta))
	case 1: // Jump Power
		config.JumpPower = max(1, config.JumpPower+int(delta))
	case 2: // Gravity
		config.Gravity = max(0, config.Gravity+int(delta))
	case 3: // Max Fall Speed
		config.MaxFallSpeed = max(1, config.MaxFallSpeed+int(delta))
	case 4: // Friction
		config.Friction = max(0, config.Friction+int(delta))
	case 5: // Air Control
		config.AirControl = max(0.0, min(1.0, config.AirControl+delta*0.1))
	case 6: // Collision Width
		config.CollisionBoxWidth = max(0.1, min(1.0, config.CollisionBoxWidth+delta*0.01))
	case 7: // Collision Height
		config.CollisionBoxHeight = max(0.1, min(1.0, config.CollisionBoxHeight+delta*0.01))
	case 8: // Collision Offset X
		config.CollisionBoxOffsetX = max(0.0, min(1.0, config.CollisionBoxOffsetX+delta*0.01))
	case 9: // Collision Offset Y
		config.CollisionBoxOffsetY = max(0.0, min(1.0, config.CollisionBoxOffsetY+delta*0.01))
	case 10: // Coyote Time
		config.CoyoteTime = max(0, config.CoyoteTime+int(delta))
	case 11: // Jump Buffer Time
		config.JumpBufferTime = max(0, config.JumpBufferTime+int(delta))
	case 12: // Char Scale Factor
		GameConfig.CharScaleFactor = max(0.1, GameConfig.CharScaleFactor+delta*0.1)
	case 13: // Tile Scale Factor
		GameConfig.TileScaleFactor = max(0.5, GameConfig.TileScaleFactor+delta*0.1)
	}
	
	LogDebug(fmt.Sprintf("Adjusted %s by %.2f", pt.paramNames[pt.selectedParam], delta))
}

/*
GetStatusText returns the current status of the physics tuner for display.
*/
func (pt *PhysicsTuner) GetStatusText() string {
	if !pt.active {
		return ""
	}
	
	config := &GameConfig.PlayerPhysics
	var currentValue string
	
	switch pt.selectedParam {
	case 0:
		currentValue = fmt.Sprintf("%d", config.MoveSpeed)
	case 1:
		currentValue = fmt.Sprintf("%d", config.JumpPower)
	case 2:
		currentValue = fmt.Sprintf("%d", config.Gravity)
	case 3:
		currentValue = fmt.Sprintf("%d", config.MaxFallSpeed)
	case 4:
		currentValue = fmt.Sprintf("%d", config.Friction)
	case 5:
		currentValue = fmt.Sprintf("%.2f", config.AirControl)
	case 6:
		currentValue = fmt.Sprintf("%.2f", config.CollisionBoxWidth)
	case 7:
		currentValue = fmt.Sprintf("%.2f", config.CollisionBoxHeight)
	case 8:
		currentValue = fmt.Sprintf("%.2f", config.CollisionBoxOffsetX)
	case 9:
		currentValue = fmt.Sprintf("%.2f", config.CollisionBoxOffsetY)
	case 10:
		currentValue = fmt.Sprintf("%d", config.CoyoteTime)
	case 11:
		currentValue = fmt.Sprintf("%d", config.JumpBufferTime)
	case 12:
		currentValue = fmt.Sprintf("%.2f", GameConfig.CharScaleFactor)
	case 13:
		currentValue = fmt.Sprintf("%.2f", GameConfig.TileScaleFactor)
	}
	
	return fmt.Sprintf("PHYSICS TUNER: %s = %s (Step: %.1f)", 
		pt.paramNames[pt.selectedParam], currentValue, pt.adjustmentStep)
}

/*
IsActive returns whether the physics tuner is currently active.
*/
func (pt *PhysicsTuner) IsActive() bool {
	return pt.active
}

// Helper functions for min/max
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}