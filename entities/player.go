package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

/*
Player represents the player character with all its functionality.
Handles movement, physics, input processing, and rendering for the
main character. Uses physics units for all position and velocity
calculations to ensure consistent behavior across different scale factors.

Position and velocity are stored in physics units (see engine.GetPhysicsUnit()).
*/
type Player struct {
	x          int
	y          int
	vx         int
	vy         int
	onGround   bool
	facingRight bool // Track which direction the player is facing
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
		x:          x,
		y:          y,
		vx:         0,
		vy:         0,
		onGround:   false,
		facingRight: true, // Default to facing right
	}
}

/*
HandleInput processes player input for movement and actions.
Reads keyboard input and updates player velocity accordingly.
Uses engine.GameConfig values for movement speeds and physics calculations.

Input mapping:
  - A/Left Arrow: Move left
  - D/Right Arrow: Move right
  - Space: Jump
*/
func (p *Player) HandleInput() {
	p.HandleInputWithLogging("")
}

/*
HandleInputWithLogging processes player input and logs keystrokes with position.
This version includes logging for debugging and analytics.

Parameters:
  - roomName: Name of the current room for logging context
*/
func (p *Player) HandleInputWithLogging(roomName string) {
	physicsUnit := engine.GetPhysicsUnit()
	playerX, playerY := p.GetPosition()

	// Horizontal movement - using config values
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if roomName != "" {
			engine.LogPlayerInput("A/Left", playerX, playerY, roomName)
		}
		p.vx = -engine.GameConfig.PlayerMoveSpeed * physicsUnit
		p.facingRight = false // Facing left
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if roomName != "" {
			engine.LogPlayerInput("D/Right", playerX, playerY, roomName)
		}
		p.vx = engine.GameConfig.PlayerMoveSpeed * physicsUnit
		p.facingRight = true // Facing right
	}

	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if roomName != "" {
			engine.LogPlayerInput("Space/Jump", playerX, playerY, roomName)
		}
		p.tryJump()
	}
}

// tryJump makes the player jump if possible
func (p *Player) tryJump() {
	physicsUnit := engine.GetPhysicsUnit()
	// Allow jumping even if not on ground (mid-air jumping as per original design)
	p.vy = -engine.GameConfig.PlayerJumpPower * physicsUnit
}

/*
Update handles player physics and movement.
Applies velocity to position, handles ground collision, applies friction
to horizontal movement, and applies gravity. Should be called once per frame.

Uses values from engine.GameConfig for all physics calculations including friction,
gravity, and ground level.
*/
func (p *Player) Update() {
	physicsUnit := engine.GetPhysicsUnit()

	// Apply movement
	p.x += p.vx
	p.y += p.vy

	// Ground collision - using config ground level
	groundY := engine.GameConfig.GroundLevel * physicsUnit
	if p.y > groundY {
		p.y = groundY
		p.onGround = true
	} else {
		p.onGround = false
	}

	// Apply friction to horizontal movement
	if p.vx > 0 {
		p.vx -= engine.GameConfig.PlayerFriction
		if p.vx < 0 {
			p.vx = 0
		}
	} else if p.vx < 0 {
		p.vx += engine.GameConfig.PlayerFriction
		if p.vx > 0 {
			p.vx = 0
		}
	}

	// Apply gravity
	if p.vy < engine.GameConfig.MaxFallSpeed*physicsUnit {
		p.vy += engine.GameConfig.Gravity
	}
}

/*
Draw renders the player character at its world position.
The player is drawn using its current sprite (idle, left, or right)
at its world coordinates. Camera transformation should be applied
at a higher level, not by the entity itself.

Parameters:
  - screen: The target screen/image to render the player to
*/
func (p *Player) Draw(screen *ebiten.Image) {
	// Choose sprite based on movement direction
	sprite := engine.GetIdleSprite()
	switch {
	case p.vx > 0:
		sprite = engine.GetRightSprite()
	case p.vx < 0:
		sprite = engine.GetLeftSprite()
	}

	// Set up drawing options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
	// Convert player position from physics units to pixels
	renderX := float64(p.x) / float64(engine.GetPhysicsUnit())
	renderY := float64(p.y) / float64(engine.GetPhysicsUnit())
	op.GeoM.Translate(renderX, renderY)

	// Draw the sprite
	screen.DrawImage(sprite, op)
}

/*
DrawDebug renders debug visualization for the player.
Shows bounding box, position markers, and movement vectors.

Parameters:
  - screen: The target screen/image to render debug info to
  - cameraOffsetX: Camera X offset for viewport transformation
  - cameraOffsetY: Camera Y offset for viewport transformation
*/
func (p *Player) DrawDebug(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Convert player position from physics units to render position
	physicsUnit := engine.GetPhysicsUnit()
	renderX := float64(p.x)/float64(physicsUnit) + cameraOffsetX
	renderY := float64(p.y)/float64(physicsUnit) + cameraOffsetY
	
	// Get sprite bounds (assuming 32x32 base sprite)
	spriteWidth := 32.0 * engine.GameConfig.CharScaleFactor
	spriteHeight := 32.0 * engine.GameConfig.CharScaleFactor
	
	// Draw bounding box
	boxColor := color.RGBA{0, 255, 0, 128} // Green for player
	if !p.onGround {
		boxColor = color.RGBA{255, 255, 0, 128} // Yellow when airborne
	}
	
	// Draw the bounding box
	vector.StrokeRect(screen, float32(renderX), float32(renderY), 
		float32(spriteWidth), float32(spriteHeight), 2, boxColor, false)
	
	// Draw center point
	centerX := renderX + spriteWidth/2
	centerY := renderY + spriteHeight/2
	vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), 
		3, color.RGBA{255, 0, 0, 255}, false)
	
	// Draw velocity vector if moving
	if p.vx != 0 || p.vy != 0 {
		// Scale velocity for visualization
		velScale := 5.0
		endX := centerX + float64(p.vx)*velScale/float64(physicsUnit)
		endY := centerY + float64(p.vy)*velScale/float64(physicsUnit)
		
		vector.StrokeLine(screen, float32(centerX), float32(centerY),
			float32(endX), float32(endY), 2, color.RGBA{255, 0, 255, 200}, false)
		
		// Draw arrowhead
		vector.DrawFilledCircle(screen, float32(endX), float32(endY), 
			3, color.RGBA{255, 0, 255, 255}, false)
	}
	
	// Draw facing direction indicator
	dirIndicatorY := renderY + spriteHeight + 5
	if p.facingRight {
		vector.StrokeLine(screen, float32(centerX), float32(dirIndicatorY),
			float32(centerX+10), float32(dirIndicatorY), 3, color.RGBA{0, 255, 255, 255}, false)
		// Arrow point
		vector.StrokeLine(screen, float32(centerX+10), float32(dirIndicatorY),
			float32(centerX+7), float32(dirIndicatorY-3), 2, color.RGBA{0, 255, 255, 255}, false)
		vector.StrokeLine(screen, float32(centerX+10), float32(dirIndicatorY),
			float32(centerX+7), float32(dirIndicatorY+3), 2, color.RGBA{0, 255, 255, 255}, false)
	} else {
		vector.StrokeLine(screen, float32(centerX), float32(dirIndicatorY),
			float32(centerX-10), float32(dirIndicatorY), 3, color.RGBA{0, 255, 255, 255}, false)
		// Arrow point
		vector.StrokeLine(screen, float32(centerX-10), float32(dirIndicatorY),
			float32(centerX-7), float32(dirIndicatorY-3), 2, color.RGBA{0, 255, 255, 255}, false)
		vector.StrokeLine(screen, float32(centerX-10), float32(dirIndicatorY),
			float32(centerX-7), float32(dirIndicatorY+3), 2, color.RGBA{0, 255, 255, 255}, false)
	}
	
	// Draw ground sensor line
	groundCheckY := float32(renderY + spriteHeight)
	vector.StrokeLine(screen, float32(renderX), groundCheckY,
		float32(renderX+spriteWidth), groundCheckY, 1, color.RGBA{255, 128, 0, 128}, false)
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
	p.facingRight = true // Reset to default facing right
}

/*
IsFacingRight returns whether the player is currently facing right.
Used by rendering systems and mini-map to show player orientation.

Returns true if facing right, false if facing left.
*/
func (p *Player) IsFacingRight() bool {
	return p.facingRight
}
