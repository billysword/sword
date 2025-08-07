package entities

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"sword/engine"
)

/*
Player represents the main character entity in the game.
Handles movement, physics, input processing, and rendering for the
main character. Uses physics units for all position and velocity
calculations to maintain consistent behavior across different scales.

Position and velocity are stored in physics units (see engine.GetPhysicsUnit()).
The player's collision box can be configured through engine.GameConfig.PlayerPhysics.
*/
type Player struct {
	x, y     int  // Position in physics units
	vx, vy   int  // Velocity in physics units per frame
	onGround bool // Whether the player is currently on the ground

	// Jump mechanics state
	coyoteTimer     int  // Frames since leaving ground
	jumpBufferTimer int  // Frames since jump was pressed
	isJumping       bool // Currently in a jump (for variable height)
	jumpHeldFrames  int  // How long jump has been held

	// Direction state
	facingRight bool // Whether the player is facing right
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
		x:           x,
		y:           y,
		vx:          0,
		vy:          0,
		onGround:    false,
		facingRight: true, // Default to facing right
	}
}

/*
ProcessInput handles player input for movement and actions.
Reads keyboard input and updates player velocity accordingly.
This method should be called once per frame before Update().

Uses engine.GameConfig.PlayerPhysics values for movement speeds and physics calculations.
Implements advanced jump mechanics including coyote time and jump buffering.
*/
func (p *Player) ProcessInput() {
	physicsUnit := engine.GetPhysicsUnit()
	config := &engine.GameConfig.PlayerPhysics

	// Update jump buffer timer
	if p.jumpBufferTimer > 0 {
		p.jumpBufferTimer--
	}

	// Horizontal movement
	p.vx = 0
	moveSpeed := config.MoveSpeed
	if !p.onGround {
		// Apply air control
		moveSpeed = int(float64(moveSpeed) * config.AirControl)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.vx = -moveSpeed * physicsUnit
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.vx = moveSpeed * physicsUnit
	}

	// Jump input handling
	jumpPressed := ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW)

	if jumpPressed {
		p.jumpBufferTimer = config.JumpBufferTime
	}

	// Check if we can jump (on ground or within coyote time)
	canJump := p.onGround || (p.coyoteTimer > 0 && p.vy >= 0)

	// Execute jump if buffered and able
	if p.jumpBufferTimer > 0 && canJump {
		p.Jump()
		p.jumpBufferTimer = 0
		p.coyoteTimer = 0
	}

	// Variable jump height - reduce upward velocity if jump released early
	if config.VariableJumpHeight && p.isJumping && !jumpPressed && p.vy < 0 {
		// Calculate how much to reduce jump
		minVelocity := int(float64(-config.JumpPower*physicsUnit) * config.MinJumpHeight)
		if p.vy < minVelocity {
			p.vy = minVelocity
		}
		p.isJumping = false
	}

	// Fast fall when holding down
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if p.vy > 0 {
			p.vy = int(float64(p.vy) * config.FastFallMultiplier)
		}
	}
}

/*
Jump makes the player jump if they are on the ground.
Sets the vertical velocity to the configured jump power.
*/
func (p *Player) Jump() {
	physicsUnit := engine.GetPhysicsUnit()
	config := &engine.GameConfig.PlayerPhysics

	p.vy = -config.JumpPower * physicsUnit
	p.onGround = false
	p.isJumping = true
	p.jumpHeldFrames = 0
}

/*
HandleInput processes player input for movement and actions.
Backward compatibility wrapper for ProcessInput().
*/
func (p *Player) HandleInput() {
	p.ProcessInput()
}

/*
HandleInputWithLogging processes player input and logs keystrokes with position.
Backward compatibility wrapper that includes logging.

Parameters:
  - roomName: Name of the current room for logging context
*/
func (p *Player) HandleInputWithLogging(roomName string) {
	// Log current position before processing input
	if roomName != "" {
		playerX, playerY := p.GetPosition()
		engine.LogPlayerPosition(playerX, playerY, roomName)
	}
	p.ProcessInput()
}

/*
Update handles player physics and movement.
Applies velocity to position, handles ground collision, applies friction
to horizontal movement, and applies gravity. Should be called once per frame.

Uses values from engine.GameConfig.PlayerPhysics for all physics calculations including friction,
gravity, and ground level. Also updates jump mechanics timers.
*/
func (p *Player) Update() {
	physicsUnit := engine.GetPhysicsUnit()
	config := &engine.GameConfig.PlayerPhysics

	// Update position
	p.x += p.vx
	p.y += p.vy

	// Update coyote time
	if p.onGround {
		p.coyoteTimer = config.CoyoteTime
	} else if p.coyoteTimer > 0 {
		p.coyoteTimer--
	}

	// Ground collision - using config ground level
	groundY := engine.GameConfig.GroundLevel * physicsUnit
	if p.y > groundY {
		p.y = groundY
		p.vy = 0
		wasInAir := !p.onGround
		p.onGround = true
		p.isJumping = false

		// Reset coyote timer when landing
		if wasInAir {
			p.coyoteTimer = config.CoyoteTime
		}
	} else {
		p.onGround = false
	}

	// Apply friction
	if p.onGround {
		// Ground friction
		if p.vx > 0 {
			p.vx -= config.Friction
			if p.vx < 0 {
				p.vx = 0
			}
		} else if p.vx < 0 {
			p.vx += config.Friction
			if p.vx > 0 {
				p.vx = 0
			}
		}
	} else if config.AirFriction > 0 {
		// Air friction (if configured)
		if p.vx > 0 {
			p.vx -= config.AirFriction
			if p.vx < 0 {
				p.vx = 0
			}
		} else if p.vx < 0 {
			p.vx += config.AirFriction
			if p.vx > 0 {
				p.vx = 0
			}
		}
	}

	// Apply gravity
	if p.vy < config.MaxFallSpeed*physicsUnit {
		p.vy += config.Gravity
	}

	// Update jump held frames
	if p.isJumping && p.vy < 0 {
		p.jumpHeldFrames++
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
	var sprite *ebiten.Image
	switch {
	case p.vx > 0:
		sprite = engine.GetPlayerSprite("right")
	case p.vx < 0:
		sprite = engine.GetPlayerSprite("left")
	default:
		sprite = engine.GetPlayerSprite("idle")
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
	engine.LogDebug(fmt.Sprintf("DRAW_OBJECT: Player(%d,%d)", p.x, p.y))
}

/*
DrawDebug renders debug visualization for the player.
Shows collision box, ground check area, position markers, and movement vectors.
The collision box is drawn according to the PlayerPhysicsConfig settings.

Parameters:
  - screen: The target screen/image to render debug info to
  - cameraOffsetX: Camera X offset for viewport transformation
  - cameraOffsetY: Camera Y offset for viewport transformation
*/
func (p *Player) DrawDebug(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Convert player position from physics units to render position
	physicsUnit := engine.GetPhysicsUnit()
	config := &engine.GameConfig.PlayerPhysics
	renderX := float64(p.x)/float64(physicsUnit) + cameraOffsetX
	renderY := float64(p.y)/float64(physicsUnit) + cameraOffsetY

	// Calculate sprite bounds with scaling
	spriteWidth := float64(config.SpriteWidth) * engine.GameConfig.CharScaleFactor
	spriteHeight := float64(config.SpriteHeight) * engine.GameConfig.CharScaleFactor

	// Calculate collision box based on configuration
	collisionX := renderX + (spriteWidth * config.CollisionBoxOffsetX)
	collisionY := renderY + (spriteHeight * config.CollisionBoxOffsetY)
	collisionWidth := spriteWidth * config.CollisionBoxWidth
	collisionHeight := spriteHeight * config.CollisionBoxHeight

	// Draw collision box
	boxColor := color.RGBA{0, 255, 0, 128} // Green for player
	if !p.onGround {
		boxColor = color.RGBA{255, 255, 0, 128} // Yellow when airborne
	}
	if p.isJumping {
		boxColor = color.RGBA{0, 255, 255, 128} // Cyan when jumping
	}

	vector.StrokeRect(screen, float32(collisionX), float32(collisionY),
		float32(collisionWidth), float32(collisionHeight), 2, boxColor, false)

	// Draw ground check area
	groundCheckX := collisionX + (collisionWidth * (1.0 - config.GroundCheckWidth) / 2.0)
	groundCheckY := collisionY + collisionHeight
	groundCheckWidth := collisionWidth * config.GroundCheckWidth
	groundCheckHeight := float64(config.GroundCheckOffset)

	groundCheckColor := color.RGBA{255, 0, 255, 64} // Magenta for ground check
	vector.DrawFilledRect(screen, float32(groundCheckX), float32(groundCheckY),
		float32(groundCheckWidth), float32(groundCheckHeight), groundCheckColor, false)

	// Draw sprite bounds (for reference)
	spriteColor := color.RGBA{128, 128, 128, 64} // Gray for sprite bounds
	vector.StrokeRect(screen, float32(renderX), float32(renderY),
		float32(spriteWidth), float32(spriteHeight), 1, spriteColor, false)

	// Draw center point of collision box
	centerX := collisionX + collisionWidth/2
	centerY := collisionY + collisionHeight/2
	vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), 3, color.RGBA{255, 0, 0, 255}, false)

	// Draw velocity vector
	if p.vx != 0 || p.vy != 0 {
		velScale := 0.5 // Scale factor for velocity visualization
		endX := centerX + float64(p.vx)*velScale/float64(physicsUnit)
		endY := centerY + float64(p.vy)*velScale/float64(physicsUnit)
		vector.StrokeLine(screen, float32(centerX), float32(centerY),
			float32(endX), float32(endY), 2, color.RGBA{255, 0, 0, 255}, false)
	}

	// Draw debug text
	debugY := int(renderY - 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pos: %d,%d", p.x/physicsUnit, p.y/physicsUnit),
		int(renderX), debugY)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Vel: %d,%d", p.vx, p.vy),
		int(renderX), debugY+12)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Ground: %v", p.onGround),
		int(renderX), debugY+24)
	if p.coyoteTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Coyote: %d", p.coyoteTimer),
			int(renderX), debugY+36)
	}
	if p.jumpBufferTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("JumpBuf: %d", p.jumpBufferTimer),
			int(renderX), debugY+48)
	}
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
