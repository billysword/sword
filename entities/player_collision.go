package entities

import (
	"sword/engine"
)

/*
CollisionBox represents the player's collision box in world coordinates.
All values are in physics units.
*/
type CollisionBox struct {
	X, Y          int // Top-left position
	Width, Height int // Dimensions
}

/*
GetCollisionBox returns the player's current collision box in physics units.
The collision box is calculated based on the sprite position and the
collision box configuration in PlayerPhysicsConfig.
*/
func (p *Player) GetCollisionBox() CollisionBox {
	config := &engine.GameConfig.PlayerPhysics
	
	// Calculate sprite dimensions in pixels (physics units now equal base pixels)
	spriteWidth := int(float64(config.SpriteWidth) * engine.GameConfig.CharScaleFactor)
	spriteHeight := int(float64(config.SpriteHeight) * engine.GameConfig.CharScaleFactor)
	
	// Calculate collision box based on configuration
	offsetX := int(float64(spriteWidth) * config.CollisionBoxOffsetX)
	offsetY := int(float64(spriteHeight) * config.CollisionBoxOffsetY)
	width := int(float64(spriteWidth) * config.CollisionBoxWidth)
	height := int(float64(spriteHeight) * config.CollisionBoxHeight)
	
	// Already in pixels; use directly
	return CollisionBox{
		X:      p.x + offsetX,
		Y:      p.y + offsetY,
		Width:  width,
		Height: height,
	}
}

/*
CheckTileCollision checks if the player would collide with solid tiles at a given position.
Used for movement validation and collision response.

Parameters:
  - tileProvider: The tile provider (room) to check tiles in
  - testX, testY: Position to test in physics units
*/
func (p *Player) CheckTileCollision(tileProvider TileProvider, testX, testY int) bool {
	// Save current position
	oldX, oldY := p.x, p.y
	
	// Temporarily move to test position
	p.x, p.y = testX, testY
	box := p.GetCollisionBox()
	
	// Restore position
	p.x, p.y = oldX, oldY
	
	// Get physics unit for conversion
	// Convert collision box to tile coordinates (physics units -> tile indices)
	leftTile := box.X / engine.GetPhysicsUnit()
	rightTile := (box.X + box.Width - 1) / engine.GetPhysicsUnit()
	topTile := box.Y / engine.GetPhysicsUnit()
	bottomTile := (box.Y + box.Height - 1) / engine.GetPhysicsUnit()
	
	// Check all tiles the collision box overlaps
	tiles := tileProvider.GetTiles()
	roomWidth := tileProvider.GetWidth()
	roomHeight := tileProvider.GetHeight()
	
	for y := topTile; y <= bottomTile; y++ {
		for x := leftTile; x <= rightTile; x++ {
			// Skip out-of-bounds tiles
			if x < 0 || x >= roomWidth || y < 0 || y >= roomHeight {
				continue
			}
			
			// Get tile at this position
			tileIndex := y*roomWidth + x
			if tileIndex >= 0 && tileIndex < len(tiles) {
				if IsSolidTile(tiles[tileIndex]) {
					return true
				}
			}
		}
	}
	
	return false
}

/*
UpdateWithTileCollision updates the player with proper tile collision detection.
This replaces the simple ground collision with actual tile-based physics.

Parameters:
  - tileProvider: The tile provider (room) for collision detection
*/
func (p *Player) UpdateWithTileCollision(tileProvider TileProvider) {
	config := &engine.GameConfig.PlayerPhysics
	
	// Update coyote time
	if p.onGround {
		p.coyoteTimer = config.CoyoteTime
	} else if p.coyoteTimer > 0 {
		p.coyoteTimer--
	}
	
	// Helper for stepped axis movement to prevent tunneling
	stepUnit := engine.GetPhysicsUnit() / 4
	if stepUnit < 1 {
		stepUnit = 1
	}
	
	// Horizontal movement (stepped)
	if p.vx != 0 {
		remaining := p.vx
		step := stepUnit
		if remaining < 0 {
			step = -stepUnit
		}
		for remaining != 0 {
			// Clamp step to remaining distance
			if remaining > 0 && step > remaining {
				step = remaining
			}
			if remaining < 0 && step < remaining {
				step = remaining
			}
			nextX := p.x + step
			if !p.CheckTileCollision(tileProvider, nextX, p.y) {
				p.x = nextX
				remaining -= step
				// Reset step sign in case we clamped above
				if remaining > 0 {
					step = stepUnit
				} else if remaining < 0 {
					step = -stepUnit
				}
			} else {
				// Hit a wall
				p.vx = 0
				break
			}
		}
	}
	
	// Vertical movement (stepped)
	wasOnGround := p.onGround
	p.onGround = false
	if p.vy != 0 {
		remaining := p.vy
		step := stepUnit
		if remaining < 0 {
			step = -stepUnit
		}
		for remaining != 0 {
			// Clamp step to remaining distance
			if remaining > 0 && step > remaining {
				step = remaining
			}
			if remaining < 0 && step < remaining {
				step = remaining
			}
			nextY := p.y + step
			if !p.CheckTileCollision(tileProvider, p.x, nextY) {
				p.y = nextY
				remaining -= step
				// Reset step sign in case we clamped above
				if remaining > 0 {
					step = stepUnit
				} else if remaining < 0 {
					step = -stepUnit
				}
			} else {
				// Vertical collision
				if p.vy > 0 {
					// Landed on ground
					p.onGround = true
					p.isJumping = false
					if !wasOnGround {
						p.coyoteTimer = config.CoyoteTime
					}
				}
				p.vy = 0
				break
			}
		}
	}
	
	// Ground check - look slightly below collision box
	offset := int(float64(config.GroundCheckOffset) * engine.GameConfig.CharScaleFactor)
	if p.CheckTileCollision(tileProvider, p.x, p.y+offset) {
		p.onGround = true
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
	if !p.onGround && p.vy < config.MaxFallSpeed {
		p.vy += config.Gravity
	}
	
	// Update jump held frames
	if p.isJumping && p.vy < 0 {
		p.jumpHeldFrames++
	}
	
	// Keep player in room bounds (one physicsUnit per tile)
	if p.x < 0 {
		p.x = 0
		p.vx = 0
	}
	maxX := tileProvider.GetWidth() * engine.GetPhysicsUnit()
	if p.x > maxX {
		p.x = maxX
		p.vx = 0
	}
	
	if p.y < 0 {
		p.y = 0
		p.vy = 0
	}
	maxY := tileProvider.GetHeight() * engine.GetPhysicsUnit()
	if p.y > maxY {
		p.y = maxY
		p.vy = 0
	}
}

// IsSolidTile checks if a tile index represents a solid tile for collision.
// This is a copy of the logic from world package to avoid circular dependency.
func IsSolidTile(tileIndex int) bool {
	// Define which tile indices are solid for collision
	switch tileIndex {
	case -1: // empty
		return false
	case 0: // dirt - solid
		return true
	case 1, 2, 3, 4, 5, 6, 7, 8: // walls, corners, ceilings - solid
		return true
	case 9, 10, 11, 12, 13, 14, 15: // platform tiles - solid
		return true
	case 16, 17, 18, 19: // inner corners - solid
		return true
	case 20, 21: // floor tiles - solid
		return true
	case 22, 23: // more walls - solid
		return true
	default:
		// Unknown tiles are considered non-solid by default
		return false
	}
}