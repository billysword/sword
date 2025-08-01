package gamestate

import (
	"math"
)

// Camera represents the game camera that follows the player
type Camera struct {
	x, y          float64 // Camera position (top-left corner of viewport)
	targetX, targetY float64 // Target position for smooth following
	width, height int     // Viewport dimensions
	worldWidth    int     // Total world width in pixels
	worldHeight   int     // Total world height in pixels
	smoothing     float64 // Camera smoothing factor (0-1, higher = smoother)
	
	// Dead zone for player movement (camera won't move if player is within this zone)
	deadZoneX     int
	deadZoneY     int
	
	// Border margins to ensure we don't show beyond world edges
	marginLeft    int
	marginRight   int
	marginTop     int
	marginBottom  int
}

// NewCamera creates a new camera with specified viewport dimensions
func NewCamera(viewportWidth, viewportHeight int) *Camera {
	return &Camera{
		width:        viewportWidth,
		height:       viewportHeight,
		smoothing:    GameConfig.CameraSmoothing,
		deadZoneX:    int(float64(viewportWidth) * GameConfig.CameraDeadZoneX),
		deadZoneY:    int(float64(viewportHeight) * GameConfig.CameraDeadZoneY),
		marginLeft:   GameConfig.CameraMarginLeft,
		marginRight:  GameConfig.CameraMarginRight,
		marginTop:    GameConfig.CameraMarginTop,
		marginBottom: GameConfig.CameraMarginBottom,
	}
}

// SetWorldBounds sets the total world size for the camera to constrain to
func (c *Camera) SetWorldBounds(worldWidth, worldHeight int) {
	c.worldWidth = worldWidth
	c.worldHeight = worldHeight
}

// Update updates the camera position to follow the target (usually the player)
func (c *Camera) Update(targetX, targetY int) {
	// Calculate ideal camera position to center the target
	idealX := float64(targetX) - float64(c.width)/2
	idealY := float64(targetY) - float64(c.height)/2
	
	// Apply dead zone logic
	currentCenterX := c.x + float64(c.width)/2
	currentCenterY := c.y + float64(c.height)/2
	
	// Only move camera if target is outside dead zone
	if math.Abs(float64(targetX)-currentCenterX) > float64(c.deadZoneX) {
		c.targetX = idealX
	}
	if math.Abs(float64(targetY)-currentCenterY) > float64(c.deadZoneY) {
		c.targetY = idealY
	}
	
	// Smooth camera movement using linear interpolation
	c.x += (c.targetX - c.x) * c.smoothing
	c.y += (c.targetY - c.y) * c.smoothing
	
	// Constrain camera to world bounds with margins
	c.constrainToWorld()
}

// constrainToWorld ensures the camera doesn't show beyond world boundaries
func (c *Camera) constrainToWorld() {
	// Left boundary with margin
	if c.x < float64(-c.marginLeft) {
		c.x = float64(-c.marginLeft)
		c.targetX = c.x
	}
	
	// Right boundary with margin
	maxX := float64(c.worldWidth - c.width + c.marginRight)
	if c.x > maxX {
		c.x = maxX
		c.targetX = c.x
	}
	
	// Top boundary with margin
	if c.y < float64(-c.marginTop) {
		c.y = float64(-c.marginTop)
		c.targetY = c.y
	}
	
	// Bottom boundary with margin
	maxY := float64(c.worldHeight - c.height + c.marginBottom)
	if c.y > maxY {
		c.y = maxY
		c.targetY = c.y
	}
}

// GetOffset returns the camera offset for rendering
func (c *Camera) GetOffset() (float64, float64) {
	return -c.x, -c.y
}

// GetPosition returns the current camera position
func (c *Camera) GetPosition() (float64, float64) {
	return c.x, c.y
}

// SetPosition directly sets the camera position (useful for room transitions)
func (c *Camera) SetPosition(x, y float64) {
	c.x = x
	c.y = y
	c.targetX = x
	c.targetY = y
}

// CenterOn immediately centers the camera on a position without smoothing
func (c *Camera) CenterOn(x, y int) {
	newX := float64(x) - float64(c.width)/2
	newY := float64(y) - float64(c.height)/2
	c.SetPosition(newX, newY)
}

// IsVisible checks if a rectangle is visible in the current viewport
func (c *Camera) IsVisible(x, y, width, height int) bool {
	return float64(x+width) > c.x && 
	       float64(x) < c.x+float64(c.width) &&
	       float64(y+height) > c.y && 
	       float64(y) < c.y+float64(c.height)
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY int) (int, int) {
	worldX := int(c.x) + screenX
	worldY := int(c.y) + screenY
	return worldX, worldY
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY int) (int, int) {
	screenX := worldX - int(c.x)
	screenY := worldY - int(c.y)
	return screenX, screenY
}