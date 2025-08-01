package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Player represents the player character with all its functionality
type Player struct {
	x  int
	y  int
	vx int
	vy int
	onGround bool
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y int) *Player {
	return &Player{
		x: x,
		y: y,
		vx: 0,
		vy: 0,
		onGround: false,
	}
}

// HandleInput processes player input
func (p *Player) HandleInput() {
	// Horizontal movement
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.vx = -4 * PHYSICS_UNIT
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.vx = 4 * PHYSICS_UNIT
	}

	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.tryJump()
	}
}

// tryJump makes the player jump if possible
func (p *Player) tryJump() {
	// Allow jumping even if not on ground (mid-air jumping as per original design)
	p.vy = -10 * PHYSICS_UNIT
}

// Update handles player physics and movement
func (p *Player) Update() {
	// Apply movement
	p.x += p.vx
	p.y += p.vy

	// Ground collision
	if p.y > groundY*PHYSICS_UNIT {
		p.y = groundY * PHYSICS_UNIT
		p.onGround = true
	} else {
		p.onGround = false
	}

	// Apply friction to horizontal movement
	if p.vx > 0 {
		p.vx -= 4
		if p.vx < 0 {
			p.vx = 0
		}
	} else if p.vx < 0 {
		p.vx += 4
		if p.vx > 0 {
			p.vx = 0
		}
	}

	// Apply gravity
	if p.vy < 20*PHYSICS_UNIT {
		p.vy += 8
	}
}

// Draw renders the player character
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
	op.GeoM.Scale(CHAR_SCALE_FACTOR, CHAR_SCALE_FACTOR)
	op.GeoM.Translate(float64(p.x)/float64(PHYSICS_UNIT), float64(p.y)/float64(PHYSICS_UNIT))
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

// DrawWithCamera renders the player character with camera offset
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
	op.GeoM.Scale(CHAR_SCALE_FACTOR, CHAR_SCALE_FACTOR)
	// Convert player position from physics units to pixels and apply camera offset
	renderX := float64(p.x)/float64(PHYSICS_UNIT) + cameraOffsetX
	renderY := float64(p.y)/float64(PHYSICS_UNIT) + cameraOffsetY
	op.GeoM.Translate(renderX, renderY)
	
	// Draw the sprite
	screen.DrawImage(sprite, op)
}

// GetPosition returns the player's current position
func (p *Player) GetPosition() (int, int) {
	return p.x, p.y
}

// SetPosition sets the player's position
func (p *Player) SetPosition(x, y int) {
	p.x = x
	p.y = y
}

// GetVelocity returns the player's current velocity
func (p *Player) GetVelocity() (int, int) {
	return p.vx, p.vy
}

// SetVelocity sets the player's velocity
func (p *Player) SetVelocity(vx, vy int) {
	p.vx = vx
	p.vy = vy
}

// IsOnGround returns whether the player is on the ground
func (p *Player) IsOnGround() bool {
	return p.onGround
}

// Reset resets the player to initial state at given position
func (p *Player) Reset(x, y int) {
	p.x = x
	p.y = y
	p.vx = 0
	p.vy = 0
	p.onGround = false
}