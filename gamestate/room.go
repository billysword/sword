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

// CollisionInfo represents collision data for physics resolution
type CollisionInfo struct {
	HasCollision bool
	CollisionX   int  // X position where collision occurs
	CollisionY   int  // Y position where collision occurs
	SurfaceType  TileType // Type of surface collided with
}

// Tile represents a single tile in the tile map
type Tile struct {
	Type   TileType
	X, Y   int
	Sprite *ebiten.Image
}

// TileMap represents a 2D grid of tile indices for a zone
type TileMap struct {
	Width  int
	Height int
	Tiles  [][]int  // Just store tile indices, -1 for empty
}

// NewTileMap creates a new tile map with specified dimensions
func NewTileMap(width, height int) *TileMap {
	tiles := make([][]int, height)
	for i := range tiles {
		tiles[i] = make([]int, width)
		// Initialize with -1 (empty)
		for j := range tiles[i] {
			tiles[i][j] = -1
		}
	}

	return &TileMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

// SetTile sets a tile index at the specified position
func (tm *TileMap) SetTile(x, y, tileIndex int) {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		tm.Tiles[y][x] = tileIndex
	}
}

// GetTileIndex returns the tile index at the specified position
func (tm *TileMap) GetTileIndex(x, y int) int {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		return tm.Tiles[y][x]
	}
	return -1
}

// Room represents a modular game area with its own tile map and logic
type Room interface {
	// Core room functionality
	GetTileMap() *TileMap
	GetZoneID() string

	// Game logic that can be extracted from main loop
	Update(player *Player) error
	HandleCollisions(player *Player)

	// Room-specific events
	OnEnter(player *Player)
	OnExit(player *Player)

	// Rendering
	Draw(screen *ebiten.Image)
	DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64)
	DrawTiles(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image)
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
func (br *BaseRoom) Update(player *Player) error {
	// Default: no special room logic
	return nil
}

// HandleCollisions provides default collision handling
func (br *BaseRoom) HandleCollisions(player *Player) {
	// Default: basic ground collision using existing groundY
	x, y := player.GetPosition()
	if y > groundY*PHYSICS_UNIT {
		player.SetPosition(x, groundY*PHYSICS_UNIT)
	}
}

// IsSolidTile checks if a tile index represents a solid tile for collision
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
		return false
	}
}

// OnEnter is called when entering the room
func (br *BaseRoom) OnEnter(player *Player) {
	// Default: no special entry logic
}

// OnExit is called when leaving the room
func (br *BaseRoom) OnExit(player *Player) {
	// Default: no special exit logic
}

// Draw renders the room (base implementation - rooms should override this)
func (br *BaseRoom) Draw(screen *ebiten.Image) {
	// Base rooms need a sprite provider, so this is just a placeholder
	// Individual room implementations should override this method
}

// DrawWithCamera renders the room with camera offset
func (br *BaseRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Base implementation - rooms should override this
	br.Draw(screen)
}

// DrawTiles renders the room's tile map using a sprite provider function
func (br *BaseRoom) DrawTiles(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image) {
	br.DrawTilesWithCamera(screen, spriteProvider, 0, 0)
}

// DrawTilesWithCamera renders the room's tile map with camera offset
func (br *BaseRoom) DrawTilesWithCamera(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	if br.tileMap == nil {
		return
	}

	for y := 0; y < br.tileMap.Height; y++ {
		for x := 0; x < br.tileMap.Width; x++ {
			tileIndex := br.tileMap.Tiles[y][x]
			if tileIndex != -1 {
				sprite := spriteProvider(tileIndex)
				if sprite != nil {
					op := &ebiten.DrawImageOptions{}
					// Scale tiles using global scale factor
					op.GeoM.Scale(TILE_SCALE_FACTOR, TILE_SCALE_FACTOR)
					renderX := float64(x * PHYSICS_UNIT) + cameraOffsetX
					renderY := float64(y * PHYSICS_UNIT) + cameraOffsetY
					op.GeoM.Translate(renderX, renderY)
					
					screen.DrawImage(sprite, op)
				}
			}
		}
	}
}

