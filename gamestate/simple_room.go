package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Forest tile indices based on the actual tilemap layout
const (
	TILE_DIRT = 0
	TILE_TOP_LEFT_CORNER = 1
	TILE_RIGHT_WALL_1 = 2
	TILE_RIGHT_WALL_2 = 3
	TILE_BOTTOM_LEFT_CORNER = 4
	TILE_TOP_RIGHT_CORNER = 5
	TILE_LEFT_WALL_1 = 6
	TILE_CEILING_1 = 7
	TILE_CEILING_2 = 8
	TILE_SINGLE_TOP = 9
	TILE_SINGLE_BOTTOM = 10
	TILE_SINGLE_LEFT = 11
	TILE_SINGLE_RIGHT = 12
	TILE_FLOATING = 13
	TILE_SINGLE_HORIZONTAL = 14
	TILE_SINGLE_VERTICAL = 15
	TILE_INNER_CORNER_TOP_LEFT = 16
	TILE_INNER_CORNER_TOP_RIGHT = 17
	TILE_INNER_CORNER_BOTTOM_RIGHT = 18
	TILE_INNER_CORNER_BOTTOM_LEFT = 19
	TILE_FLOOR_1 = 20
	TILE_FLOOR_2 = 21
	TILE_LEFT_WALL_2 = 22
	TILE_BOTTOM_RIGHT_CORNER = 23
)

// SimpleRoom is a basic room implementation with a forest theme
type SimpleRoom struct {
	*BaseRoom
	tileSize int
	tilesPerRow int
	forestTiles map[int]*ebiten.Image
}

// NewSimpleRoom creates a new simple room with forest tiles
func NewSimpleRoom(zoneID string) *SimpleRoom {
	// Create a room that fits the screen (960x540 with 16px tiles = 60x33 tiles)
	room := &SimpleRoom{
		BaseRoom: NewBaseRoom(zoneID, 60, 34),
		tileSize: 16,
		tilesPerRow: 8, // Assuming 8 tiles per row in forest-tiles.png
		forestTiles: make(map[int]*ebiten.Image),
	}

	// Extract individual tiles from the forest tilemap
	room.extractForestTiles()
	
	// Initialize with a forest layout
	room.initializeLayout()

	return room
}

// extractForestTiles extracts individual tiles from the forest tilemap
func (sr *SimpleRoom) extractForestTiles() {
	if globalTileSprite == nil {
		return
	}

	// Extract individual tiles from the forest tilemap
	// Assuming 24 tiles total (0-23) arranged in rows of 8
	for i := 0; i <= 23; i++ {
		x := (i % sr.tilesPerRow) * sr.tileSize
		y := (i / sr.tilesPerRow) * sr.tileSize
		
		tileImg := ebiten.NewImage(sr.tileSize, sr.tileSize)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(-x), float64(-y))
		tileImg.DrawImage(globalTileSprite, op)
		
		sr.forestTiles[i] = tileImg
	}
}

// getTileSprite returns the appropriate sprite for a tile type
func (sr *SimpleRoom) getTileSprite(tileIndex int) *ebiten.Image {
	if sprite, exists := sr.forestTiles[tileIndex]; exists {
		return sprite
	}
	// Fallback to dirt if not found
	if sprite, exists := sr.forestTiles[TILE_DIRT]; exists {
		return sprite
	}
	return globalTileSprite
}

// initializeLayout sets up the forest tile layout for this room
func (sr *SimpleRoom) initializeLayout() {
	if globalTileSprite == nil {
		return
	}

	// Create ground tiles at the bottom rows
	groundRow := groundY / unit

	// Create main ground platform with proper corners and edges
	platformStart := 5
	platformEnd := sr.tileMap.Width - 5

	for x := 0; x < sr.tileMap.Width; x++ {
		if x >= platformStart && x <= platformEnd {
			// Main platform surface
			tileType := TILE_FLOOR_1
			if x%3 == 0 {
				tileType = TILE_FLOOR_2 // Alternate floor tiles for variation
			}
			sr.tileMap.SetTile(x, groundRow, TileGround, sr.getTileSprite(tileType))
			
			// Underground layers with dirt
			for y := groundRow + 1; y < sr.tileMap.Height; y++ {
				sr.tileMap.SetTile(x, y, TileGround, sr.getTileSprite(TILE_DIRT))
			}
		} else {
			// Areas without platform - just dirt underground
			for y := groundRow + 1; y < sr.tileMap.Height; y++ {
				sr.tileMap.SetTile(x, y, TileGround, sr.getTileSprite(TILE_DIRT))
			}
		}
	}

	// Add floating platforms with proper edges
	sr.createPlatform(15, 25, groundRow-5)
	sr.createPlatform(35, 45, groundRow-8)
	sr.createPlatform(50, 55, groundRow-3)

	// Add some single floating tiles
	sr.tileMap.SetTile(10, groundRow-2, TilePlatform, sr.getTileSprite(TILE_FLOATING))
	sr.tileMap.SetTile(30, groundRow-6, TilePlatform, sr.getTileSprite(TILE_FLOATING))
	sr.tileMap.SetTile(48, groundRow-10, TilePlatform, sr.getTileSprite(TILE_FLOATING))
}

// createPlatform creates a platform with proper corners and edges
func (sr *SimpleRoom) createPlatform(startX, endX, y int) {
	if startX >= endX || y < 0 || y >= sr.tileMap.Height {
		return
	}

	for x := startX; x <= endX; x++ {
		var tileType int
		
		if x == startX && x == endX {
			// Single tile platform
			tileType = TILE_SINGLE_HORIZONTAL
		} else if x == startX {
			// Left edge
			tileType = TILE_SINGLE_LEFT
		} else if x == endX {
			// Right edge
			tileType = TILE_SINGLE_RIGHT
		} else {
			// Middle tiles
			if (x-startX)%2 == 0 {
				tileType = TILE_FLOOR_1
			} else {
				tileType = TILE_FLOOR_2
			}
		}
		
		sr.tileMap.SetTile(x, y, TilePlatform, sr.getTileSprite(tileType))
	}
}

// Update handles room-specific logic
func (sr *SimpleRoom) Update(player *Player) error {
	// Add any room-specific update logic here
	// For now, just use the base implementation
	return sr.BaseRoom.Update(player)
}

// HandleCollisions provides collision detection for this room
func (sr *SimpleRoom) HandleCollisions(player *Player) {
	// Get player position
	playerX, playerY := player.GetPosition()
	
	// Convert player position to tile coordinates
	charTileX := playerX / (unit * unit)
	charTileY := playerY / (unit * unit)

	// Check collision with ground tiles
	if charTileY >= 0 && charTileY < sr.tileMap.Height {
		for checkY := charTileY; checkY < sr.tileMap.Height; checkY++ {
			if charTileX >= 0 && charTileX < sr.tileMap.Width {
				tile := sr.tileMap.GetTile(charTileX, checkY)
				if tile != nil && (tile.Type == TileGround || tile.Type == TilePlatform) {
					// Found ground, stop falling
					if playerY > checkY*unit*unit {
						player.SetPosition(playerX, checkY*unit*unit)
						vx, _ := player.GetVelocity()
						player.SetVelocity(vx, 0)
					}
					return
				}
			}
		}
	}

	// Fallback to original ground collision
	sr.BaseRoom.HandleCollisions(player)
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
