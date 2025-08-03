package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
)

/*
BaseEnemy provides common functionality for all enemy types.
Contains shared physics, rendering, and state management code.
Concrete enemy types should embed this struct and implement their own HandleAI() method.

Position and velocity are stored in physics units (see engine.GetPhysicsUnit()).
*/
type BaseEnemy struct {
	x        int
	y        int
	vx       int
	vy       int
	onGround bool
	
	// Common properties that can be overridden by specific enemy types
	moveSpeed     int    // Movement speed in physics units
	friction      int    // Friction applied each frame
	scaleX        float64 // Horizontal scale factor (can be negative for flipping)
	scaleY        float64 // Vertical scale factor
}

/*
NewBaseEnemy creates a new base enemy at the specified position.
Initializes common properties with default values. Concrete enemy types
should call this and then set their specific properties.

Parameters:
  - x: Initial horizontal position in physics units
  - y: Initial vertical position in physics units

Returns a pointer to the new BaseEnemy instance.
*/
func NewBaseEnemy(x, y int) *BaseEnemy {
	return &BaseEnemy{
		x:        x,
		y:        y,
		vx:       0,
		vy:       0,
		onGround: false,
		
		// Default properties
		moveSpeed: engine.GameConfig.PlayerMoveSpeed / 2, // Half player speed by default
		friction:  engine.GameConfig.PlayerFriction,
		scaleX:    engine.GameConfig.CharScaleFactor,
		scaleY:    engine.GameConfig.CharScaleFactor,
	}
}

/*
NOTE: HandleAI() is intentionally NOT implemented in BaseEnemy.
This method must be implemented by concrete enemy types (SlimeEnemy, etc.)
to provide their specific AI behavior. BaseEnemy is not meant to be used directly.
*/

/*
Update handles common physics and movement for all enemy types.
Applies physics like velocity, ground collision, friction, and gravity.
Should be called once per frame AFTER the concrete enemy type has
handled its AI logic and set velocity values.

Uses values from engine.GameConfig for physics calculations.
*/
func (be *BaseEnemy) Update() {
	physicsUnit := engine.GetPhysicsUnit()
	
	// NOTE: AI logic should be handled by the concrete enemy type
	// before calling this Update() method
	
	// Apply movement
	be.x += be.vx
	be.y += be.vy

	// Ground collision - using same ground level as player
	groundY := engine.GameConfig.GroundLevel * physicsUnit
	if be.y > groundY {
		be.y = groundY
		be.onGround = true
	} else {
		be.onGround = false
	}

	// Apply friction to horizontal movement
	if be.vx > 0 {
		be.vx -= be.friction
		if be.vx < 0 {
			be.vx = 0
		}
	} else if be.vx < 0 {
		be.vx += be.friction
		if be.vx > 0 {
			be.vx = 0
		}
	}

	// Apply gravity
	if be.vy < engine.GameConfig.MaxFallSpeed*physicsUnit {
		be.vy += engine.GameConfig.Gravity
	}
}

/*
Draw renders the enemy character using the default player sprites.
Concrete enemy types can override this method to use their own sprites.
Chooses sprite based on movement direction and applies scaling.

Parameters:
  - screen: The target screen/image to render the enemy to
*/
func (be *BaseEnemy) Draw(screen *ebiten.Image) {
	// Choose sprite based on movement direction (using player sprites by default)
	sprite := engine.GetIdleSprite()
	switch {
	case be.vx > 0:
		sprite = engine.GetRightSprite()
		be.scaleX = engine.GameConfig.CharScaleFactor // Face right
	case be.vx < 0:
		sprite = engine.GetLeftSprite()
		be.scaleX = engine.GameConfig.CharScaleFactor // Face left
	}

	// Set up drawing options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(be.scaleX, be.scaleY)
	op.GeoM.Translate(float64(be.x)/float64(engine.GetPhysicsUnit()), float64(be.y)/float64(engine.GetPhysicsUnit()))
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
DrawWithCamera renders the enemy character with camera offset.
Similar to Draw() but applies camera transformation for scrolling worlds.
The camera offset is applied in addition to the enemy's world position.

Parameters:
  - screen: The target screen/image to render the enemy to
  - cameraOffsetX: Horizontal camera offset in pixels
  - cameraOffsetY: Vertical camera offset in pixels
*/
func (be *BaseEnemy) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Choose sprite based on movement direction
	sprite := engine.GetIdleSprite()
	switch {
	case be.vx > 0:
		sprite = engine.GetRightSprite()
		be.scaleX = engine.GameConfig.CharScaleFactor // Face right
	case be.vx < 0:
		sprite = engine.GetLeftSprite()
		be.scaleX = engine.GameConfig.CharScaleFactor // Face left
	}

	// Set up drawing options with camera offset
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(be.scaleX, be.scaleY)
	// Convert enemy position from physics units to pixels and apply camera offset
	renderX := float64(be.x)/float64(engine.GetPhysicsUnit()) + cameraOffsetX
	renderY := float64(be.y)/float64(engine.GetPhysicsUnit()) + cameraOffsetY
	op.GeoM.Translate(renderX, renderY)
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
GetPosition returns the enemy's current position.
Position values are in physics units and represent the enemy's world coordinates.

Returns:
  - x: Horizontal position in physics units
  - y: Vertical position in physics units
*/
func (be *BaseEnemy) GetPosition() (int, int) {
	return be.x, be.y
}

/*
SetPosition sets the enemy's position.
Directly updates the enemy's world coordinates.

Parameters:
  - x: New horizontal position in physics units
  - y: New vertical position in physics units
*/
func (be *BaseEnemy) SetPosition(x, y int) {
	be.x = x
	be.y = y
}

/*
GetVelocity returns the enemy's current velocity.
Velocity values are in physics units per frame.

Returns:
  - vx: Horizontal velocity in physics units per frame
  - vy: Vertical velocity in physics units per frame
*/
func (be *BaseEnemy) GetVelocity() (int, int) {
	return be.vx, be.vy
}

/*
SetVelocity sets the enemy's velocity.
Directly modifies the enemy's movement speed and direction.

Parameters:
  - vx: New horizontal velocity in physics units per frame
  - vy: New vertical velocity in physics units per frame
*/
func (be *BaseEnemy) SetVelocity(vx, vy int) {
	be.vx = vx
	be.vy = vy
}

/*
IsOnGround returns whether the enemy is on the ground.
Useful for AI decisions and physics calculations.

Returns true if the enemy is currently touching the ground.
*/
func (be *BaseEnemy) IsOnGround() bool {
	return be.onGround
}

/*
Reset resets the enemy to initial state at given position.
Clears all velocity and sets the enemy to the specified position.
Concrete enemy types should override this to reset their specific AI state.

Parameters:
  - x: Reset position horizontal coordinate in physics units
  - y: Reset position vertical coordinate in physics units
*/
func (be *BaseEnemy) Reset(x, y int) {
	be.x = x
	be.y = y
	be.vx = 0
	be.vy = 0
	be.onGround = false
}

/*
GetEnemyType returns the type name of this enemy.
Base implementation returns "base" - concrete enemy types should override this.

Returns the enemy type identifier as a string.
*/
func (be *BaseEnemy) GetEnemyType() string {
	return "base"
}

// Protected methods for use by concrete enemy types

/*
SetMoveSpeed sets the movement speed for this enemy type.
Used by concrete enemy implementations to customize their speed.

Parameters:
  - speed: Movement speed in physics units
*/
func (be *BaseEnemy) SetMoveSpeed(speed int) {
	be.moveSpeed = speed
}

/*
GetMoveSpeed returns the current movement speed.
Used by concrete enemy implementations in their AI logic.

Returns the movement speed in physics units.
*/
func (be *BaseEnemy) GetMoveSpeed() int {
	return be.moveSpeed
}

/*
SetScale sets the rendering scale for this enemy type.
Used by concrete enemy implementations to customize their size.

Parameters:
  - scaleX: Horizontal scale factor
  - scaleY: Vertical scale factor
*/
func (be *BaseEnemy) SetScale(scaleX, scaleY float64) {
	be.scaleX = scaleX
	be.scaleY = scaleY
}