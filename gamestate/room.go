package gamestate

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/*
TileType represents different types of tiles.
Used to categorize tiles for collision detection, rendering, and gameplay logic.
Each tile type has different properties and behaviors in the game world.
*/
type TileType int

const (
	TileEmpty TileType = iota      // Empty space, no collision
	TileGround                     // Solid ground tile
	TilePlatform                   // Platform that can be jumped through from below
	TileWall                       // Solid wall tile
	TileBackground                 // Background decoration, no collision
)

/*
CollisionInfo represents collision data for physics resolution.
Provides detailed information about collision detection results,
including position data and surface type for proper physics response.
*/
type CollisionInfo struct {
	HasCollision bool      // Whether a collision was detected
	CollisionX   int       // X position where collision occurs
	CollisionY   int       // Y position where collision occurs
	SurfaceType  TileType  // Type of surface collided with
}

/*
Tile represents a single tile in the tile map.
Contains all data needed to render and interact with an individual tile,
including its type, position, and sprite reference.
*/
type Tile struct {
	Type   TileType         // Type of tile for collision and logic
	X, Y   int              // Position coordinates in tile units
	Sprite *ebiten.Image    // Sprite image for rendering this tile
}

/*
TileMap represents a 2D grid of tile indices for a zone.
Stores the layout of tiles in a room as a 2D array of indices.
Uses -1 to represent empty/air tiles. The indices correspond to
tile types or sprite indices depending on the room implementation.
*/
type TileMap struct {
	Width  int       // Width of the tile map in tiles
	Height int       // Height of the tile map in tiles
	Tiles  [][]int   // 2D array of tile indices, -1 for empty
}

/*
NewTileMap creates a new tile map with specified dimensions.
Initializes all tiles to empty (-1) by default. The returned TileMap
is ready to be populated with tile data.

Parameters:
  - width: Width of the tile map in tiles
  - height: Height of the tile map in tiles

Returns a pointer to the new TileMap instance.
*/
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

/*
SetTile sets a tile index at the specified position.
Updates the tile map with a new tile index at the given coordinates.
Performs bounds checking to prevent invalid array access.

Parameters:
  - x: Horizontal tile coordinate
  - y: Vertical tile coordinate  
  - tileIndex: Index of the tile to place (-1 for empty)
*/
func (tm *TileMap) SetTile(x, y, tileIndex int) {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		tm.Tiles[y][x] = tileIndex
	}
}

/*
GetTileIndex returns the tile index at the specified position.
Retrieves the tile index stored at the given coordinates.
Returns -1 for positions outside the map bounds or empty tiles.

Parameters:
  - x: Horizontal tile coordinate
  - y: Vertical tile coordinate

Returns the tile index at the position, or -1 if empty/out of bounds.
*/
func (tm *TileMap) GetTileIndex(x, y int) int {
	if x >= 0 && x < tm.Width && y >= 0 && y < tm.Height {
		return tm.Tiles[y][x]
	}
	return -1
}

/*
Room represents a modular game area with its own tile map and logic.
This interface defines the contract for all room implementations, allowing
for different room types with custom behavior while maintaining consistent
integration with the game systems.

Room implementations should handle:
  - Tile-based collision detection and response
  - Room-specific game logic and events
  - Rendering with proper camera support
  - Player interaction and room transitions
*/
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

/*
BaseRoom provides default implementation for common room functionality.
Serves as a foundation for room implementations, providing basic tile map
management and default behavior that can be overridden as needed.

Most rooms should embed BaseRoom and override specific methods for
custom behavior rather than implementing the entire Room interface.
*/
type BaseRoom struct {
	zoneID  string    // Unique identifier for this room
	tileMap *TileMap  // Tile layout for this room
}

/*
NewBaseRoom creates a new base room.
Initializes a base room with the specified ID and tile map dimensions.
The tile map is initialized with all empty tiles.

Parameters:
  - zoneID: Unique identifier for this room
  - width: Width of the room in tiles
  - height: Height of the room in tiles

Returns a pointer to the new BaseRoom instance.
*/
func NewBaseRoom(zoneID string, width, height int) *BaseRoom {
	return &BaseRoom{
		zoneID:  zoneID,
		tileMap: NewTileMap(width, height),
	}
}

/*
GetTileMap returns the room's tile map.
Provides access to the room's tile layout for collision detection,
rendering, and other tile-based operations.

Returns a pointer to the room's TileMap.
*/
func (br *BaseRoom) GetTileMap() *TileMap {
	return br.tileMap
}

/*
GetZoneID returns the room's zone identifier.
The zone ID is used for room identification, save systems,
and debugging purposes.

Returns the room's unique identifier string.
*/
func (br *BaseRoom) GetZoneID() string {
	return br.zoneID
}

/*
Update provides default room update logic.
Base implementation performs no special room logic. Individual room
types should override this method to implement custom behavior like
moving platforms, environmental hazards, or interactive elements.

Parameters:
  - player: The player instance for interaction

Returns any error from room update logic.
*/
func (br *BaseRoom) Update(player *Player) error {
	// Default: no special room logic
	return nil
}

/*
HandleCollisions provides default collision handling.
Base implementation uses simple ground collision based on the configured
ground level. More complex rooms should override this for tile-based
collision detection and platform behavior.

Parameters:
  - player: The player instance to check collisions for
*/
func (br *BaseRoom) HandleCollisions(player *Player) {
	// Default: basic ground collision using config ground level
	physicsUnit := GetPhysicsUnit()
	groundY := GameConfig.GroundLevel * physicsUnit
	
	x, y := player.GetPosition()
	if y > groundY {
		player.SetPosition(x, groundY)
	}
}

/*
IsSolidTile checks if a tile index represents a solid tile for collision.
Determines whether a given tile index should block player movement.
This function defines the collision properties of different tile types.

Parameters:
  - tileIndex: The tile index to check

Returns true if the tile should block movement, false otherwise.
*/
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

/*
OnEnter is called when entering the room.
Base implementation performs no special actions. Override this method
to implement room entry effects, music changes, or setup logic.

Parameters:
  - player: The player entering the room
*/
func (br *BaseRoom) OnEnter(player *Player) {
	// Default: no special entry logic
}

/*
OnExit is called when leaving the room.
Base implementation performs no special actions. Override this method
to implement cleanup, save room state, or exit effects.

Parameters:
  - player: The player leaving the room
*/
func (br *BaseRoom) OnExit(player *Player) {
	// Default: no special exit logic
}

/*
Draw renders the room (base implementation - rooms should override this).
Base implementation provides no rendering. Individual room types should
override this method to implement proper tile rendering and background drawing.

Parameters:
  - screen: The target screen/image to render to
*/
func (br *BaseRoom) Draw(screen *ebiten.Image) {
	// Base rooms need a sprite provider, so this is just a placeholder
	// Individual room implementations should override this method
}

/*
DrawWithCamera renders the room with camera offset.
Base implementation delegates to Draw() method. Rooms should override
this to properly handle camera-based rendering for scrolling worlds.

Parameters:
  - screen: The target screen/image to render to
  - cameraOffsetX: Horizontal camera offset in pixels
  - cameraOffsetY: Vertical camera offset in pixels
*/
func (br *BaseRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Base implementation - rooms should override this
	br.Draw(screen)
}

/*
DrawTiles renders the room's tile map using a sprite provider function.
Renders all tiles in the room without camera offset. Uses the provided
sprite provider function to get the appropriate image for each tile index.

Parameters:
  - screen: The target screen/image to render to
  - spriteProvider: Function that returns sprite for a given tile index
*/
func (br *BaseRoom) DrawTiles(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image) {
	br.DrawTilesWithCamera(screen, spriteProvider, 0, 0)
}

/*
DrawTilesWithCamera renders the room's tile map with camera offset.
Core tile rendering method that handles drawing all tiles with proper
scaling and camera transformation. Uses the sprite provider to get
the appropriate image for each tile index.

Parameters:
  - screen: The target screen/image to render to
  - spriteProvider: Function that returns sprite for a given tile index  
  - cameraOffsetX: Horizontal camera offset in pixels
  - cameraOffsetY: Vertical camera offset in pixels
*/
func (br *BaseRoom) DrawTilesWithCamera(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	if br.tileMap == nil {
		return
	}

	physicsUnit := GetPhysicsUnit()
	
	for y := 0; y < br.tileMap.Height; y++ {
		for x := 0; x < br.tileMap.Width; x++ {
			tileIndex := br.tileMap.Tiles[y][x]
			if tileIndex != -1 {
				sprite := spriteProvider(tileIndex)
				if sprite != nil {
					op := &ebiten.DrawImageOptions{}
					// Scale tiles using global scale factor
					op.GeoM.Scale(GameConfig.TileScaleFactor, GameConfig.TileScaleFactor)
					renderX := float64(x * physicsUnit) + cameraOffsetX
					renderY := float64(y * physicsUnit) + cameraOffsetY
					op.GeoM.Translate(renderX, renderY)
					
					screen.DrawImage(sprite, op)
				}
			}
		}
	}
}

