package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TileType represents different types of tiles
type TileType int

const (
	TileEmpty TileType = iota
	TileGround
	TilePlatform
	TileWall
	TileBackground
)

// Tile represents a single tile in the tile map
type Tile struct {
	Type   TileType
	X, Y   int
	Sprite *ebiten.Image
}

// TileMap represents a 2D grid of tiles for a zone
type TileMap struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

// NewTileMap creates a new tile map with specified dimensions
func NewTileMap(width, height int) *TileMap {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}

	return &TileMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

// SetTile sets a tile at the specified position
func (tm *TileMap) SetTile(x, y int, tileType TileType, sprite *ebiten.Image) {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		tm.Tiles[y][x] = Tile{
			Type:   tileType,
			X:      x,
			Y:      y,
			Sprite: sprite,
		}
	}
}

// GetTile returns the tile at the specified position
func (tm *TileMap) GetTile(x, y int) *Tile {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		return &tm.Tiles[y][x]
	}
	return nil
}

// Room represents a modular game area with its own tile map and logic
type Room interface {
	// Core room functionality
	GetTileMap() *TileMap
	GetZoneID() string

	// Game logic that can be extracted from main loop
	Update(character *Character) error
	HandleCollisions(character *Character)

	// Room-specific events
	OnEnter(character *Character)
	OnExit(character *Character)

	// Rendering
	Draw(screen *ebiten.Image)
	DrawTiles(screen *ebiten.Image)
}

// BaseRoom provides default implementation for common room functionality
type BaseRoom struct {
	zoneID  string
	tileMap *TileMap
}

// NewBaseRoom creates a new base room
func NewBaseRoom(zoneID string, width, height int) *BaseRoom {
	return &BaseRoom{
		zoneID:  zoneID,
		tileMap: NewTileMap(width, height),
	}
}

// GetTileMap returns the room's tile map
func (br *BaseRoom) GetTileMap() *TileMap {
	return br.tileMap
}

// GetZoneID returns the room's zone identifier
func (br *BaseRoom) GetZoneID() string {
	return br.zoneID
}

// Update provides default room update logic
func (br *BaseRoom) Update(character *Character) error {
	// Default: no special room logic
	return nil
}

// HandleCollisions provides default collision handling
func (br *BaseRoom) HandleCollisions(character *Character) {
	// Default: basic ground collision using existing groundY
	if character.y > groundY*unit {
		character.y = groundY * unit
	}
}

// OnEnter is called when entering the room
func (br *BaseRoom) OnEnter(character *Character) {
	// Default: no special entry logic
}

// OnExit is called when leaving the room
func (br *BaseRoom) OnExit(character *Character) {
	// Default: no special exit logic
}

// Draw renders the room
func (br *BaseRoom) Draw(screen *ebiten.Image) {
	// Draw tiles first
	br.DrawTiles(screen)
}

// DrawTiles renders the room's tile map
func (br *BaseRoom) DrawTiles(screen *ebiten.Image) {
	if br.tileMap == nil {
		return
	}

	for y := 0; y < br.tileMap.Height; y++ {
		for x := 0; x < br.tileMap.Width; x++ {
			tile := &br.tileMap.Tiles[y][x]
			if tile.Type != TileEmpty && tile.Sprite != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*unit), float64(y*unit))
				screen.DrawImage(tile.Sprite, op)
			}
		}
	}
}
