package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/*
Player represents the player character with all its functionality.
Handles movement, physics, input processing, and rendering for the
main character. Uses physics units for all position and velocity
calculations to ensure consistent behavior across different scale factors.

Position and velocity are stored in physics units (see GetPhysicsUnit()).
*/
type Player struct {
	x  int
	y  int
	vx int
	vy int
	onGround bool
}

/*
NewPlayer creates a new player at the specified position.
Initializes the player with zero velocity and not on ground.
Position coordinates should be in physics units.

Parameters:
  - x: Initial horizontal position in physics units
  - y: Initial vertical position in physics units

Returns a pointer to the new Player instance.
*/
func NewPlayer(x, y int) *Player {
	return &Player{
		x: x,
		y: y,
		vx: 0,
		vy: 0,
		onGround: false,
	}
}

/*
HandleInput processes player input for movement and actions.
Reads keyboard input and updates player velocity accordingly.
Uses GameConfig values for movement speeds and physics calculations.

Input mapping:
  - A/Left Arrow: Move left
  - D/Right Arrow: Move right  
  - Space: Jump
*/
func (p *Player) HandleInput() {
	physicsUnit := GetPhysicsUnit()
	
	// Horizontal movement - using config values
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.vx = -GameConfig.PlayerMoveSpeed * physicsUnit
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.vx = GameConfig.PlayerMoveSpeed * physicsUnit
	}

	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.tryJump()
	}
}

// tryJump makes the player jump if possible
func (p *Player) tryJump() {
	physicsUnit := GetPhysicsUnit()
	// Allow jumping even if not on ground (mid-air jumping as per original design)
	p.vy = -GameConfig.PlayerJumpPower * physicsUnit
}

/*
Update handles player physics and movement.
Applies velocity to position, handles ground collision, applies friction
to horizontal movement, and applies gravity. Should be called once per frame.

Uses values from GameConfig for all physics calculations including friction,
gravity, and ground level.
*/
func (p *Player) Update() {
	physicsUnit := GetPhysicsUnit()
	
	// Apply movement
	p.x += p.vx
	p.y += p.vy

	// Ground collision - using config ground level
	groundY := GameConfig.GroundLevel * physicsUnit
	if p.y > groundY {
		p.y = groundY
		p.onGround = true
	} else {
		p.onGround = false
	}

	// Apply friction to horizontal movement
	if p.vx > 0 {
		p.vx -= GameConfig.PlayerFriction
		if p.vx < 0 {
			p.vx = 0
		}
	} else if p.vx < 0 {
		p.vx += GameConfig.PlayerFriction
		if p.vx > 0 {
			p.vx = 0
		}
	}

	// Apply gravity
	if p.vy < GameConfig.MaxFallSpeed*physicsUnit {
		p.vy += GameConfig.Gravity
	}
}

/*
Draw renders the player character.
Chooses the appropriate sprite based on movement direction and renders
it at the player's current position. Uses GameConfig.CharScaleFactor
for consistent scaling.

Parameters:
  - screen: The target screen/image to render the player to
*/
func (p *Player) Draw(screen *ebiten.Image) {
	// Choose sprite based on movement direction
	sprite := globalIdleSprite
	switch {
	case p.vx > 0:
		sprite = globalRightSprite
	case p.vx < 0:
		sprite = globalLeftSprite
	}

	// Set up drawing options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GameConfig.CharScaleFactor, GameConfig.CharScaleFactor)
	op.GeoM.Translate(float64(p.x)/float64(GetPhysicsUnit()), float64(p.y)/float64(GetPhysicsUnit()))
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
DrawWithCamera renders the player character with camera offset.
Similar to Draw() but applies camera transformation for scrolling worlds.
The camera offset is applied in addition to the player's world position.

Parameters:
  - screen: The target screen/image to render the player to
  - cameraOffsetX: Horizontal camera offset in pixels
  - cameraOffsetY: Vertical camera offset in pixels
*/
func (p *Player) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Choose sprite based on movement direction
	sprite := globalIdleSprite
	switch {
	case p.vx > 0:
		sprite = globalRightSprite
	case p.vx < 0:
		sprite = globalLeftSprite
	}

	// Set up drawing options with camera offset
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GameConfig.CharScaleFactor, GameConfig.CharScaleFactor)
	// Convert player position from physics units to pixels and apply camera offset
	renderX := float64(p.x)/float64(GetPhysicsUnit()) + cameraOffsetX
	renderY := float64(p.y)/float64(GetPhysicsUnit()) + cameraOffsetY
	op.GeoM.Translate(renderX, renderY)
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
GetPosition returns the player's current position.
Position values are in physics units and represent the player's
world coordinates.

Returns:
  - x: Horizontal position in physics units
  - y: Vertical position in physics units
*/
func (p *Player) GetPosition() (int, int) {
	return p.x, p.y
}

/*
SetPosition sets the player's position.
Directly updates the player's world coordinates. Useful for
teleporting, respawning, or room transitions.

Parameters:
  - x: New horizontal position in physics units
  - y: New vertical position in physics units
*/
func (p *Player) SetPosition(x, y int) {
	p.x = x
	p.y = y
}

/*
GetVelocity returns the player's current velocity.
Velocity values are in physics units per frame and indicate
the player's movement speed and direction.

Returns:
  - vx: Horizontal velocity in physics units per frame
  - vy: Vertical velocity in physics units per frame
*/
func (p *Player) GetVelocity() (int, int) {
	return p.vx, p.vy
}

/*
SetVelocity sets the player's velocity.
Directly modifies the player's movement speed and direction.
Useful for special abilities, knockback effects, or movement overrides.

Parameters:
  - vx: New horizontal velocity in physics units per frame
  - vy: New vertical velocity in physics units per frame
*/
func (p *Player) SetVelocity(vx, vy int) {
	p.vx = vx
	p.vy = vy
}

/*
IsOnGround returns whether the player is on the ground.
Useful for determining if the player can jump, applying different
physics rules, or triggering ground-based effects.

Returns true if the player is currently touching the ground.
*/
func (p *Player) IsOnGround() bool {
	return p.onGround
}

/*
Reset resets the player to initial state at given position.
Clears all velocity and sets the player to the specified position.
Useful for respawning, level restarts, or room transitions.

Parameters:
  - x: Reset position horizontal coordinate in physics units
  - y: Reset position vertical coordinate in physics units
*/
func (p *Player) Reset(x, y int) {
	p.x = x
	p.y = y
	p.vx = 0
	p.vy = 0
	p.onGround = false
}