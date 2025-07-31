package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Tile indices for the forest tilemap (assuming a typical 16x16 tile grid)
const (
	TILE_GRASS = iota
	TILE_DIRT
	TILE_STONE
	TILE_TREE_TOP
	TILE_TREE_TRUNK
	TILE_FLOWER
	TILE_MOSS_STONE
	TILE_DARK_DIRT
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
	// Assuming the tilemap is 8x8 tiles of 16x16 pixels each
	for i := 0; i < 8; i++ { // Extract 8 different tile types
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
	// Fallback to the first tile if not found
	if sprite, exists := sr.forestTiles[0]; exists {
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

	// Add varied ground tiles across the width
	for x := 0; x < sr.tileMap.Width; x++ {
		// Create varied ground surface
		groundTileType := TILE_GRASS
		if x%7 == 0 {
			groundTileType = TILE_FLOWER // Add some flowers
		} else if x%11 == 0 {
			groundTileType = TILE_MOSS_STONE // Add some moss stones
		}
		
		sr.tileMap.SetTile(x, groundRow, TileGround, sr.getTileSprite(groundTileType))
		
		// Sub-ground layers with dirt and stone
		for y := groundRow + 1; y < sr.tileMap.Height; y++ {
			subGroundTile := TILE_DIRT
			if y > groundRow + 2 {
				subGroundTile = TILE_STONE // Stone deeper underground
			} else if (x+y)%3 == 0 {
				subGroundTile = TILE_DARK_DIRT // Some variation
			}
			sr.tileMap.SetTile(x, y, TileGround, sr.getTileSprite(subGroundTile))
		}
	}

	// Add forest platforms with tree elements
	for x := 15; x < 25; x++ {
		sr.tileMap.SetTile(x, groundRow-5, TilePlatform, sr.getTileSprite(TILE_MOSS_STONE))
	}

	for x := 35; x < 45; x++ {
		sr.tileMap.SetTile(x, groundRow-8, TilePlatform, sr.getTileSprite(TILE_STONE))
	}

	// Add some trees in the background
	treePositions := []int{10, 30, 50}
	for _, x := range treePositions {
		if x < sr.tileMap.Width {
			// Tree trunk
			for y := groundRow - 4; y < groundRow; y++ {
				sr.tileMap.SetTile(x, y, TileBackground, sr.getTileSprite(TILE_TREE_TRUNK))
			}
			// Tree top
			sr.tileMap.SetTile(x, groundRow-5, TileBackground, sr.getTileSprite(TILE_TREE_TOP))
			sr.tileMap.SetTile(x-1, groundRow-4, TileBackground, sr.getTileSprite(TILE_TREE_TOP))
			sr.tileMap.SetTile(x+1, groundRow-4, TileBackground, sr.getTileSprite(TILE_TREE_TOP))
		}
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
