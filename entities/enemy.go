package entities

import (
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
)

/*
Enemy represents an AI-controlled entity in the game.
Similar to the Player but with automated movement patterns instead
of input handling. Uses the same physics system as the player for
consistent behavior and collision detection.

Position and velocity are stored in physics units (see engine.GetPhysicsUnit()).
*/
type Enemy struct {
	x  int
	y  int
	vx int
	vy int
	onGround bool
	
	// AI movement properties
	moveDirection int // -1 for left, 1 for right, 0 for stationary
	moveTimer     int // Frames until direction change
	moveSpeed     int // Movement speed in physics units
	patrolRange   int // Maximum distance to patrol from spawn point
	spawnX        int // Original spawn position for patrol limits
}

/*
NewEnemy creates a new enemy at the specified position.
Initializes the enemy with zero velocity, basic AI movement parameters,
and not on ground. Position coordinates should be in physics units.

Parameters:
  - x: Initial horizontal position in physics units
  - y: Initial vertical position in physics units

Returns a pointer to the new Enemy instance.
*/
func NewEnemy(x, y int) *Enemy {
	physicsUnit := engine.GetPhysicsUnit()
	
	return &Enemy{
		x: x,
		y: y,
		vx: 0,
		vy: 0,
		onGround: false,
		
		// AI properties - start with random direction
		moveDirection: []int{-1, 1}[rand.Intn(2)],
		moveTimer:     60 + rand.Intn(120), // 1-3 seconds at 60fps
		moveSpeed:     engine.GameConfig.PlayerMoveSpeed / 2, // Half player speed
		patrolRange:   200 * physicsUnit, // 200 physics units patrol range
		spawnX:        x,
	}
}

/*
HandleAI processes AI logic for autonomous movement.
Implements a simple patrol behavior where the enemy moves back and forth
within a certain range of its spawn point. Changes direction when reaching
patrol boundaries or after a random time interval.
*/
func (e *Enemy) HandleAI() {
	physicsUnit := engine.GetPhysicsUnit()
	
	// Decrease move timer
	e.moveTimer--
	
	// Check if we've moved too far from spawn point
	distanceFromSpawn := e.x - e.spawnX
	if distanceFromSpawn > e.patrolRange {
		e.moveDirection = -1 // Move back toward spawn
		e.moveTimer = 60 + rand.Intn(60) // Reset timer
	} else if distanceFromSpawn < -e.patrolRange {
		e.moveDirection = 1 // Move back toward spawn
		e.moveTimer = 60 + rand.Intn(60) // Reset timer
	}
	
	// Change direction randomly when timer expires
	if e.moveTimer <= 0 {
		// 70% chance to change direction, 30% chance to keep moving
		if rand.Float32() < 0.7 {
			directions := []int{-1, 0, 1} // left, stop, right
			e.moveDirection = directions[rand.Intn(len(directions))]
		}
		e.moveTimer = 60 + rand.Intn(180) // 1-4 seconds
	}
	
	// Apply movement based on direction
	if e.moveDirection != 0 {
		e.vx = e.moveDirection * e.moveSpeed * physicsUnit
	}
}

/*
Update handles enemy physics and movement.
Applies AI logic, then updates physics similar to player:
velocity to position, ground collision, friction, and gravity.
Should be called once per frame.

Uses values from engine.GameConfig for all physics calculations.
*/
func (e *Enemy) Update() {
	physicsUnit := engine.GetPhysicsUnit()
	
	// Handle AI movement
	e.HandleAI()
	
	// Apply movement
	e.x += e.vx
	e.y += e.vy

	// Ground collision - using same ground level as player
	groundY := engine.GameConfig.GroundLevel * physicsUnit
	if e.y > groundY {
		e.y = groundY
		e.onGround = true
	} else {
		e.onGround = false
	}

	// Apply friction to horizontal movement
	if e.vx > 0 {
		e.vx -= engine.GameConfig.PlayerFriction
		if e.vx < 0 {
			e.vx = 0
		}
	} else if e.vx < 0 {
		e.vx += engine.GameConfig.PlayerFriction
		if e.vx > 0 {
			e.vx = 0
		}
	}

	// Apply gravity
	if e.vy < engine.GameConfig.MaxFallSpeed*physicsUnit {
		e.vy += engine.GameConfig.Gravity
	}
}

/*
Draw renders the enemy character.
Uses the same sprite system as the player, choosing appropriate sprite
based on movement direction. Enemies will look like the player character
for now (we can differentiate them later with different sprites).

Parameters:
  - screen: The target screen/image to render the enemy to
*/
func (e *Enemy) Draw(screen *ebiten.Image) {
	// Choose sprite based on movement direction (same as player)
	sprite := engine.GetIdleSprite()
	switch {
	case e.vx > 0:
		sprite = engine.GetRightSprite()
	case e.vx < 0:
		sprite = engine.GetLeftSprite()
	}

	// Set up drawing options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
	op.GeoM.Translate(float64(e.x)/float64(engine.GetPhysicsUnit()), float64(e.y)/float64(engine.GetPhysicsUnit()))
	
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
func (e *Enemy) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Choose sprite based on movement direction
	sprite := engine.GetIdleSprite()
	switch {
	case e.vx > 0:
		sprite = engine.GetRightSprite()
	case e.vx < 0:
		sprite = engine.GetLeftSprite()
	}

	// Set up drawing options with camera offset
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
	// Convert enemy position from physics units to pixels and apply camera offset
	renderX := float64(e.x)/float64(engine.GetPhysicsUnit()) + cameraOffsetX
	renderY := float64(e.y)/float64(engine.GetPhysicsUnit()) + cameraOffsetY
	op.GeoM.Translate(renderX, renderY)
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
GetPosition returns the enemy's current position.
Position values are in physics units and represent the enemy's
world coordinates.

Returns:
  - x: Horizontal position in physics units
  - y: Vertical position in physics units
*/
func (e *Enemy) GetPosition() (int, int) {
	return e.x, e.y
}

/*
SetPosition sets the enemy's position.
Directly updates the enemy's world coordinates. Useful for
spawning, teleporting, or room transitions.

Parameters:
  - x: New horizontal position in physics units
  - y: New vertical position in physics units
*/
func (e *Enemy) SetPosition(x, y int) {
	e.x = x
	e.y = y
}

/*
GetVelocity returns the enemy's current velocity.
Velocity values are in physics units per frame and indicate
the enemy's movement speed and direction.

Returns:
  - vx: Horizontal velocity in physics units per frame
  - vy: Vertical velocity in physics units per frame
*/
func (e *Enemy) GetVelocity() (int, int) {
	return e.vx, e.vy
}

/*
SetVelocity sets the enemy's velocity.
Directly modifies the enemy's movement speed and direction.
Useful for special abilities, knockback effects, or movement overrides.

Parameters:
  - vx: New horizontal velocity in physics units per frame
  - vy: New vertical velocity in physics units per frame
*/
func (e *Enemy) SetVelocity(vx, vy int) {
	e.vx = vx
	e.vy = vy
}

/*
IsOnGround returns whether the enemy is on the ground.
Useful for determining if the enemy can jump, applying different
physics rules, or triggering ground-based effects.

Returns true if the enemy is currently touching the ground.
*/
func (e *Enemy) IsOnGround() bool {
	return e.onGround
}

/*
Reset resets the enemy to initial state at given position.
Clears all velocity, resets AI state, and sets the enemy to the specified position.
Useful for respawning, level restarts, or room transitions.

Parameters:
  - x: Reset position horizontal coordinate in physics units
  - y: Reset position vertical coordinate in physics units
*/
func (e *Enemy) Reset(x, y int) {
	e.x = x
	e.y = y
	e.vx = 0
	e.vy = 0
	e.onGround = false
	
	// Reset AI state
	e.spawnX = x
	e.moveDirection = []int{-1, 1}[rand.Intn(2)]
	e.moveTimer = 60 + rand.Intn(120)
}