package engine

import (
	"math"
)

/*
Camera represents the game camera that follows the player.
Provides smooth camera movement with dead zones, world boundary constraints,
and coordinate conversion utilities. The camera uses pixel coordinates for
its position and viewport calculations.

Key features:
  - Smooth following with configurable interpolation
  - Dead zone to prevent camera jitter during small movements  
  - World boundary constraints with configurable margins
  - Screen/world coordinate conversion utilities
*/
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

/*
NewCamera creates a new camera with specified viewport dimensions.
Initializes camera settings from GameConfig and sets up dead zones
and margins based on the viewport size and configuration values.

Parameters:
  - viewportWidth: Width of the camera viewport in pixels
  - viewportHeight: Height of the camera viewport in pixels

Returns a pointer to the new Camera instance.
*/
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

/*
SetWorldBounds sets the total world size for the camera to constrain to.
The camera will not move beyond these boundaries, accounting for the
configured margins. Call this when loading a new level or room.

Parameters:
  - worldWidth: Total world width in pixels
  - worldHeight: Total world height in pixels
*/
func (c *Camera) SetWorldBounds(worldWidth, worldHeight int) {
	c.worldWidth = worldWidth
	c.worldHeight = worldHeight
}

/*
Update updates the camera position to follow the target (usually the player).
Implements smooth camera movement with dead zone logic and world boundary
constraints. Should be called once per frame with the target's position.

The camera will only move if the target is outside the dead zone, and
movement is smoothed using linear interpolation based on the smoothing factor.

Parameters:
  - targetX: Target horizontal position in pixels
  - targetY: Target vertical position in pixels
*/
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

/*
GetOffset returns the camera offset for rendering.
The offset values can be applied to world coordinates to transform
them into screen coordinates for rendering. These are the negative
of the camera's position.

Returns:
  - offsetX: Horizontal offset to apply when rendering
  - offsetY: Vertical offset to apply when rendering
*/
func (c *Camera) GetOffset() (float64, float64) {
	return -c.x, -c.y
}

/*
GetPosition returns the current camera position.
The position represents the top-left corner of the camera viewport
in world coordinates.

Returns:
  - x: Camera horizontal position in pixels
  - y: Camera vertical position in pixels
*/
func (c *Camera) GetPosition() (float64, float64) {
	return c.x, c.y
}

/*
SetPosition directly sets the camera position (useful for room transitions).
Immediately moves the camera to the specified position without smoothing.
Also updates the target position to prevent the camera from smoothing
back to the previous target.

Parameters:
  - x: New camera horizontal position in pixels
  - y: New camera vertical position in pixels
*/
func (c *Camera) SetPosition(x, y float64) {
	c.x = x
	c.y = y
	c.targetX = x
	c.targetY = y
}

/*
CenterOn immediately centers the camera on a position without smoothing.
Calculates the camera position needed to center the viewport on the
specified coordinates and sets the camera there instantly.

Parameters:
  - x: Target horizontal position to center on in pixels
  - y: Target vertical position to center on in pixels
*/
func (c *Camera) CenterOn(x, y int) {
	newX := float64(x) - float64(c.width)/2
	newY := float64(y) - float64(c.height)/2
	c.SetPosition(newX, newY)
}

/*
IsVisible checks if a rectangle is visible in the current viewport.
Useful for culling off-screen objects to improve performance.
Uses the current camera position to determine visibility.

Parameters:
  - x: Rectangle left edge in world coordinates (pixels)
  - y: Rectangle top edge in world coordinates (pixels)  
  - width: Rectangle width in pixels
  - height: Rectangle height in pixels

Returns true if any part of the rectangle is visible in the viewport.
*/
func (c *Camera) IsVisible(x, y, width, height int) bool {
	return float64(x+width) > c.x && 
	       float64(x) < c.x+float64(c.width) &&
	       float64(y+height) > c.y && 
	       float64(y) < c.y+float64(c.height)
}

/*
ScreenToWorld converts screen coordinates to world coordinates.
Takes a position on the screen and calculates the corresponding
position in the game world, accounting for camera position.

Parameters:
  - screenX: Horizontal position on screen in pixels
  - screenY: Vertical position on screen in pixels

Returns:
  - worldX: Corresponding horizontal world position in pixels
  - worldY: Corresponding vertical world position in pixels
*/
func (c *Camera) ScreenToWorld(screenX, screenY int) (int, int) {
	worldX := int(c.x) + screenX
	worldY := int(c.y) + screenY
	return worldX, worldY
}

/*
WorldToScreen converts world coordinates to screen coordinates.
Takes a position in the game world and calculates where it should
appear on the screen, accounting for camera position.

Parameters:
  - worldX: Horizontal position in world coordinates (pixels)
  - worldY: Vertical position in world coordinates (pixels)

Returns:
  - screenX: Corresponding horizontal screen position in pixels
  - screenY: Corresponding vertical screen position in pixels
*/
func (c *Camera) WorldToScreen(worldX, worldY int) (int, int) {
	screenX := worldX - int(c.x)
	screenY := worldY - int(c.y)
	return screenX, screenY
}

// GetViewportSize returns the camera viewport dimensions
func (c *Camera) GetViewportSize() (int, int) {
	return c.width, c.height
}

// GetWorldBounds returns the world boundary dimensions
func (c *Camera) GetWorldBounds() (int, int) {
	return c.worldWidth, c.worldHeight
}

// GetTargetPosition returns the current target position for smooth following
func (c *Camera) GetTargetPosition() (float64, float64) {
	return c.targetX, c.targetY
}

// GetDeadZone returns the camera dead zone dimensions
func (c *Camera) GetDeadZone() (int, int) {
	return c.deadZoneX, c.deadZoneY
}

// GetMargins returns the camera margin settings
func (c *Camera) GetMargins() (int, int, int, int) {
	return c.marginLeft, c.marginRight, c.marginTop, c.marginBottom
}

/*
GetCenteredViewport calculates the centered viewport position for small rooms.
When the room is smaller than the camera viewport, this calculates the offset
needed to center the room in the screen with black dead areas around it.

Returns:
  - offsetX: Horizontal offset to center the room
  - offsetY: Vertical offset to center the room
  - isSmaller: Whether the room is smaller than the viewport
*/
func (c *Camera) GetCenteredViewport() (int, int, bool) {
	// Check if room is smaller than viewport
	roomTooSmallX := c.worldWidth < c.width
	roomTooSmallY := c.worldHeight < c.height
	
	offsetX := 0
	offsetY := 0
	
	if roomTooSmallX {
		// Center horizontally
		offsetX = (c.width - c.worldWidth) / 2
	}
	
	if roomTooSmallY {
		// Center vertically
		offsetY = (c.height - c.worldHeight) / 2
	}
	
	return offsetX, offsetY, roomTooSmallX || roomTooSmallY
}

/*
UpdateForSmallRoom updates camera for rooms smaller than the viewport.
When the room is smaller than the viewport, the camera position is fixed
to show the entire room centered in the viewport.
*/
func (c *Camera) UpdateForSmallRoom() {
	// For small rooms, camera should be fixed at 0,0
	// The centering is handled during rendering
	c.x = 0
	c.y = 0
	c.targetX = 0
	c.targetY = 0
}