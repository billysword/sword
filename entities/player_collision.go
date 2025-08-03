package entities

import (
	"sword/engine"
	"sword/world"
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
	physicsUnit := engine.GetPhysicsUnit()
	
	// Calculate sprite dimensions in physics units
	spriteWidth := int(float64(config.SpriteWidth) * engine.GameConfig.CharScaleFactor)
	spriteHeight := int(float64(config.SpriteHeight) * engine.GameConfig.CharScaleFactor)
	
	// Calculate collision box based on configuration
	offsetX := int(float64(spriteWidth) * config.CollisionBoxOffsetX)
	offsetY := int(float64(spriteHeight) * config.CollisionBoxOffsetY)
	width := int(float64(spriteWidth) * config.CollisionBoxWidth)
	height := int(float64(spriteHeight) * config.CollisionBoxHeight)
	
	// Convert to physics units
	return CollisionBox{
		X:      p.x + offsetX*physicsUnit/spriteWidth,
		Y:      p.y + offsetY*physicsUnit/spriteHeight,
		Width:  width * physicsUnit / spriteWidth,
		Height: height * physicsUnit / spriteHeight,
	}
}

/*
CheckTileCollision checks for collision between the player and solid tiles.
Returns true if there's a collision at the given position.

Parameters:
  - room: The current room to check tiles in
  - testX, testY: Position to test in physics units
*/
func (p *Player) CheckTileCollision(room world.Room, testX, testY int) bool {
	// Save current position
	oldX, oldY := p.x, p.y
	
	// Temporarily move to test position
	p.x, p.y = testX, testY
	box := p.GetCollisionBox()
	
	// Restore position
	p.x, p.y = oldX, oldY
	
	// Get physics unit for conversion
	physicsUnit := engine.GetPhysicsUnit()
	tileSize := int(float64(engine.GameConfig.TileSize) * engine.GameConfig.TileScaleFactor)
	
	// Convert collision box to tile coordinates
	leftTile := box.X / physicsUnit / tileSize
	rightTile := (box.X + box.Width) / physicsUnit / tileSize
	topTile := box.Y / physicsUnit / tileSize
	bottomTile := (box.Y + box.Height) / physicsUnit / tileSize
	
	// Check all tiles the collision box overlaps
	tiles := room.GetTiles()
	roomWidth := room.GetWidth()
	roomHeight := room.GetHeight()
	
	for y := topTile; y <= bottomTile; y++ {
		for x := leftTile; x <= rightTile; x++ {
			// Skip out-of-bounds tiles
			if x < 0 || x >= roomWidth || y < 0 || y >= roomHeight {
				continue
			}
			
			// Get tile at this position
			tileIndex := y*roomWidth + x
			if tileIndex >= 0 && tileIndex < len(tiles) {
				if world.IsSolidTile(tiles[tileIndex]) {
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
  - room: The current room for collision detection
*/
func (p *Player) UpdateWithTileCollision(room world.Room) {
	physicsUnit := engine.GetPhysicsUnit()
	config := &engine.GameConfig.PlayerPhysics
	
	// Store initial position
	startX, startY := p.x, p.y
	
	// Update coyote time
	if p.onGround {
		p.coyoteTimer = config.CoyoteTime
	} else if p.coyoteTimer > 0 {
		p.coyoteTimer--
	}
	
	// Try horizontal movement first
	targetX := p.x + p.vx
	if !p.CheckTileCollision(room, targetX, p.y) {
		p.x = targetX
	} else {
		// Hit a wall, stop horizontal movement
		p.vx = 0
	}
	
	// Try vertical movement
	targetY := p.y + p.vy
	wasOnGround := p.onGround
	
	if !p.CheckTileCollision(room, p.x, targetY) {
		p.y = targetY
		p.onGround = false
	} else {
		// Hit something vertically
		if p.vy > 0 {
			// Falling - hit ground
			p.onGround = true
			p.isJumping = false
			
			// Reset coyote timer when landing
			if !wasOnGround {
				p.coyoteTimer = config.CoyoteTime
			}
		}
		p.vy = 0
	}
	
	// Ground check - look slightly below collision box
	box := p.GetCollisionBox()
	groundCheckY := box.Y + box.Height + config.GroundCheckOffset*physicsUnit/engine.GameConfig.TileSize
	if p.CheckTileCollision(room, p.x, p.y+config.GroundCheckOffset*physicsUnit/engine.GameConfig.TileSize) {
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
	if !p.onGround && p.vy < config.MaxFallSpeed*physicsUnit {
		p.vy += config.Gravity
	}
	
	// Update jump held frames
	if p.isJumping && p.vy < 0 {
		p.jumpHeldFrames++
	}
	
	// Keep player in room bounds
	if p.x < 0 {
		p.x = 0
		p.vx = 0
	}
	maxX := room.GetWidth() * int(float64(engine.GameConfig.TileSize)*engine.GameConfig.TileScaleFactor) * physicsUnit
	if p.x > maxX {
		p.x = maxX
		p.vx = 0
	}
	
	if p.y < 0 {
		p.y = 0
		p.vy = 0
	}
	maxY := room.GetHeight() * int(float64(engine.GameConfig.TileSize)*engine.GameConfig.TileScaleFactor) * physicsUnit
	if p.y > maxY {
		p.y = maxY
		p.vy = 0
	}
}