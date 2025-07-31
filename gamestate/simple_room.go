package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SimpleRoom is a basic room implementation with a simple ground layout
type SimpleRoom struct {
	*BaseRoom
}

// NewSimpleRoom creates a new simple room with basic ground tiles
func NewSimpleRoom(zoneID string) *SimpleRoom {
	// Create a room that fits the screen (960x540 with 16px tiles = 60x33 tiles)
	room := &SimpleRoom{
		BaseRoom: NewBaseRoom(zoneID, 60, 34),
	}

	// Initialize with a simple ground layout
	room.initializeLayout()

	return room
}

// initializeLayout sets up the basic tile layout for this room
func (sr *SimpleRoom) initializeLayout() {
	if globalTileSprite == nil {
		return
	}

	// Create ground tiles at the bottom rows (y=24 corresponds to groundY=380 with 16px units)
	groundRow := groundY / unit

	// Add ground tiles across the width
	for x := 0; x < sr.tileMap.Width; x++ {
		// Ground layer
		sr.tileMap.SetTile(x, groundRow, TileGround, globalTileSprite)
		// Sub-ground layers for visual depth
		for y := groundRow + 1; y < sr.tileMap.Height; y++ {
			sr.tileMap.SetTile(x, y, TileGround, globalTileSprite)
		}
	}

	// Add some platform tiles for variety
	for x := 20; x < 30; x++ {
		sr.tileMap.SetTile(x, groundRow-5, TilePlatform, globalTileSprite)
	}

	for x := 35; x < 45; x++ {
		sr.tileMap.SetTile(x, groundRow-8, TilePlatform, globalTileSprite)
	}
}

// Update handles room-specific logic
func (sr *SimpleRoom) Update(character *Character) error {
	// Add any room-specific update logic here
	// For now, just use the base implementation
	return sr.BaseRoom.Update(character)
}

// HandleCollisions provides collision detection for this room
func (sr *SimpleRoom) HandleCollisions(character *Character) {
	// Convert character position to tile coordinates
	charTileX := character.x / (unit * unit)
	charTileY := character.y / (unit * unit)

	// Check collision with ground tiles
	if charTileY >= 0 && charTileY < sr.tileMap.Height {
		for checkY := charTileY; checkY < sr.tileMap.Height; checkY++ {
			if charTileX >= 0 && charTileX < sr.tileMap.Width {
				tile := sr.tileMap.GetTile(charTileX, checkY)
				if tile != nil && (tile.Type == TileGround || tile.Type == TilePlatform) {
					// Found ground, stop falling
					if character.y > checkY*unit*unit {
						character.y = checkY * unit * unit
						character.vy = 0
					}
					return
				}
			}
		}
	}

	// Fallback to original ground collision
	sr.BaseRoom.HandleCollisions(character)
}

// Draw renders the room and its tiles
func (sr *SimpleRoom) Draw(screen *ebiten.Image) {
	// Draw background first
	if globalBackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		screen.DrawImage(globalBackgroundImage, op)
	}

	// Draw tiles
	sr.DrawTiles(screen)
}
