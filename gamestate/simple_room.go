package gamestate

import (
	"image"

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
	// Create an even larger room for the zoomed out view (120x60 tiles = 1920x960 pixels at 1x scale)
	// With 800x450 viewport, this gives plenty of room to explore
	room := &SimpleRoom{
		BaseRoom: NewBaseRoom(zoneID, 120, 60),
		tileSize: TILE_SIZE,
		tilesPerRow: 8, // Forest tilemap has 8 tiles per row
		forestTiles: make(map[int]*ebiten.Image),
	}

	room.initializeForestTiles()
	room.buildRoom()
	return room
}

// initializeForestTiles extracts individual tiles from the forest tilemap
func (sr *SimpleRoom) initializeForestTiles() {
	if globalTileSprite == nil {
		return
	}

	// The forest tilemap has 8 tiles per row and 3 rows (24 tiles total, indices 0-23)
	for i := 0; i <= 23; i++ {
		x := (i % sr.tilesPerRow) * sr.tileSize
		y := (i / sr.tilesPerRow) * sr.tileSize
		
		// Use SubImage to extract the tile from the tilemap
		subImg := globalTileSprite.SubImage(image.Rect(x, y, x+sr.tileSize, y+sr.tileSize)).(*ebiten.Image)
		sr.forestTiles[i] = subImg
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

// buildRoom sets up the forest tile layout for this larger metroidvania-style room
func (sr *SimpleRoom) buildRoom() {
	if globalTileSprite == nil {
		return
	}

	// Create a larger room with multiple platforms and areas to explore
	// The room is 120x60 tiles (1920x960 pixels at 1x scale)
	
	// Fill the bottom with ground tiles
	for y := 45; y < 60; y++ {
		for x := 0; x < 120; x++ {
			sr.tileMap.SetTile(x, y, TILE_DIRT)
		}
	}
	
	// Create the main ground level with proper edges
	groundY := 44
	// Left edge
	sr.tileMap.SetTile(0, groundY, TILE_TOP_LEFT_CORNER)
	// Top surface
	for x := 1; x < 119; x++ {
		sr.tileMap.SetTile(x, groundY, TILE_CEILING_1)
	}
	// Right edge
	sr.tileMap.SetTile(119, groundY, TILE_TOP_RIGHT_CORNER)
	
	// Add more platforms spread across the larger room
	// Lower platforms
	sr.createPlatform(10, 35, 10)
	sr.createPlatform(30, 38, 8)
	sr.createPlatform(50, 36, 12)
	sr.createPlatform(75, 39, 10)
	sr.createPlatform(95, 37, 15)
	
	// Mid-level platforms
	sr.createPlatform(15, 28, 8)
	sr.createPlatform(40, 25, 10)
	sr.createPlatform(65, 27, 12)
	sr.createPlatform(85, 26, 10)
	sr.createPlatform(105, 30, 10)
	
	// High platforms
	sr.createPlatform(20, 18, 10)
	sr.createPlatform(45, 15, 12)
	sr.createPlatform(70, 17, 8)
	sr.createPlatform(90, 14, 14)
	
	// Very high platforms for skilled players
	sr.createPlatform(35, 8, 8)
	sr.createPlatform(60, 7, 10)
	sr.createPlatform(80, 9, 8)
	
	// Add walls and structures throughout the room
	// Left area structures
	sr.createWall(25, 35, 44)
	sr.createWall(26, 35, 44)
	
	// Center-left structures
	sr.createWall(45, 30, 44)
	sr.createWall(46, 30, 44)
	
	// Center structures
	sr.createWall(60, 25, 44)
	sr.createWall(61, 25, 44)
	
	// Center-right structures
	sr.createWall(80, 32, 44)
	sr.createWall(81, 32, 44)
	
	// Right area structures
	sr.createWall(100, 28, 44)
	sr.createWall(101, 28, 44)
	
	// Add some floating single tiles for decoration
	sr.tileMap.SetTile(55, 20, TILE_FLOATING)
	sr.tileMap.SetTile(35, 22, TILE_FLOATING)
	sr.tileMap.SetTile(75, 21, TILE_FLOATING)
	sr.tileMap.SetTile(95, 19, TILE_FLOATING)
}

// createPlatform creates a floating platform at the specified position
func (sr *SimpleRoom) createPlatform(x, y, width int) {
	if width < 2 {
		return
	}
	
	// Left edge
	sr.tileMap.SetTile(x, y, TILE_SINGLE_LEFT)
	
	// Middle tiles
	for i := 1; i < width-1; i++ {
		sr.tileMap.SetTile(x+i, y, TILE_SINGLE_HORIZONTAL)
	}
	
	// Right edge
	sr.tileMap.SetTile(x+width-1, y, TILE_SINGLE_RIGHT)
}

// createWall creates a vertical wall from startY to endY
func (sr *SimpleRoom) createWall(x, startY, endY int) {
	for y := startY; y <= endY; y++ {
		if y == startY {
			sr.tileMap.SetTile(x, y, TILE_SINGLE_TOP)
		} else if y == endY {
			sr.tileMap.SetTile(x, y, TILE_SINGLE_BOTTOM)
		} else {
			sr.tileMap.SetTile(x, y, TILE_SINGLE_VERTICAL)
		}
	}
}

// initializeLayout sets up the forest tile layout for this room using a simple tile array
func (sr *SimpleRoom) initializeLayout() {
	if globalTileSprite == nil {
		return
	}

	// Define the level layout as a 2D array of tile indices
	// -1 = empty, 0+ = tile index from the tilemap
	levelLayout := [][]int{
		// Row 0-22: Sky area (mostly empty with some floating platforms)
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 13, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 13, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 11, 20, 21, 20, 21, 20, 21, 20, 21, 20, 12, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 11, 20, 21, 20, 21, 20, 21, 20, 21, 20, 12, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 11, 20, 21, 20, 12, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, 13, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		// Ground level starts here (around row 23)
		{-1, -1, -1, -1, -1, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5, -1, -1, -1, -1},
		{20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21, 20, 21},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// Apply the layout to the tile map
	sr.loadFromLayout(levelLayout)
}

// loadFromLayout applies a 2D tile index array directly to the tile map
func (sr *SimpleRoom) loadFromLayout(layout [][]int) {
	// Simply copy the layout to the tilemap - much simpler!
	for y, row := range layout {
		if y >= sr.tileMap.Height {
			break
		}
		for x, tileIndex := range row {
			if x >= sr.tileMap.Width {
				break
			}
			sr.tileMap.SetTile(x, y, tileIndex)
		}
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
	charTileX := playerX / PHYSICS_UNIT
	charTileY := playerY / PHYSICS_UNIT

	// Check collision with ground tiles
	if charTileY >= 0 && charTileY < sr.tileMap.Height {
		for checkY := charTileY; checkY < sr.tileMap.Height; checkY++ {
			if charTileX >= 0 && charTileX < sr.tileMap.Width {
				tileIndex := sr.tileMap.GetTileIndex(charTileX, checkY)
				if IsSolidTile(tileIndex) {
					// Found solid ground, stop falling
					if playerY > checkY*PHYSICS_UNIT {
						player.SetPosition(playerX, checkY*PHYSICS_UNIT)
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
	// Draw background first (if enabled)
	if GetBackgroundVisible() && globalBackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		screen.DrawImage(globalBackgroundImage, op)
	}

	// Draw tiles using sprite provider function
	sr.DrawTiles(screen, sr.getTileSprite)
	
	// Draw debug grid overlay (if enabled)
	DrawGrid(screen)
}

// DrawWithCamera renders the room with camera offset
func (sr *SimpleRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Draw background with parallax effect (slower movement)
	if GetBackgroundVisible() && globalBackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		// Apply parallax scrolling to background (moves slower than foreground)
		bgOffsetX := cameraOffsetX * 0.3 // 30% of camera movement for parallax
		bgOffsetY := cameraOffsetY * 0.3
		op.GeoM.Translate(bgOffsetX, bgOffsetY)
		screen.DrawImage(globalBackgroundImage, op)
	}

	// Draw tiles with camera offset
	sr.DrawTilesWithCamera(screen, sr.getTileSprite, cameraOffsetX, cameraOffsetY)
	
	// Draw debug grid overlay (if enabled) - grid moves with camera
	if GetGridVisible() {
		DrawGridWithCamera(screen, cameraOffsetX, cameraOffsetY)
	}
}
