package entities

import (
	"reflect"
	"sword/engine"
)

/*
CollisionService provides a clean interface for collision detection
with better separation of concerns. This encapsulates collision logic
and reduces coupling between player physics and world representation.
*/
type CollisionService struct {
	solidityProvider TileSolidityProvider
	roomWidth        int
	roomHeight       int
	// tiles holds a flattened snapshot of the provider's tiles for fallback collision
	// when no explicit TileSolidityProvider is available
	tiles       []int
	providerPtr uintptr
}

/*
NewCollisionService creates a new collision service for the given room.
Parameters:
  - tileProvider: Provides basic tile information (width, height)
  - solidityProvider: Provides collision detection functionality
*/
func NewCollisionService(tileProvider TileProvider) *CollisionService {
	cs := &CollisionService{
		roomWidth:   tileProvider.GetWidth(),
		roomHeight:  tileProvider.GetHeight(),
		tiles:       tileProvider.GetTiles(),
		providerPtr: reflect.ValueOf(tileProvider).Pointer(),
	}
	
	// Check if room supports advanced collision detection
	if sp, ok := any(tileProvider).(TileSolidityProvider); ok {
		cs.solidityProvider = sp
		engine.LogInfo("CollisionService: Using advanced collision detection")
	} else {
		engine.LogInfo("CollisionService: Room doesn't support TileSolidityProvider")
	}
	
	return cs
}

// IsForProvider returns true if this collision service matches the given tile provider
func (cs *CollisionService) IsForProvider(tileProvider TileProvider) bool {
	if cs == nil || tileProvider == nil {
		return false
	}
	return cs.providerPtr == reflect.ValueOf(tileProvider).Pointer() &&
		cs.roomWidth == tileProvider.GetWidth() &&
		cs.roomHeight == tileProvider.GetHeight()
}

/*
CheckBoxCollision checks if a collision box would collide with solid tiles.
Parameters:
  - box: The collision box to test
Returns true if the box would collide with any solid tiles.
*/
func (cs *CollisionService) CheckBoxCollision(box CollisionBox) bool {
	// Convert collision box to tile coordinates
	physicsUnit := engine.GetPhysicsUnit()
	leftTile := box.X / physicsUnit
	rightTile := (box.X + box.Width - 1) / physicsUnit
	topTile := box.Y / physicsUnit
	bottomTile := (box.Y + box.Height - 1) / physicsUnit
	
	// Check all tiles the collision box overlaps
	for y := topTile; y <= bottomTile; y++ {
		for x := leftTile; x <= rightTile; x++ {
			// Skip out-of-bounds tiles
			if x < 0 || x >= cs.roomWidth || y < 0 || y >= cs.roomHeight {
				continue
			}
			
			// Check if this tile is solid
			tileIndex := y*cs.roomWidth + x
			if cs.solidityProvider != nil {
				if cs.solidityProvider.IsSolidAtFlatIndex(tileIndex) {
					return true
				}
			} else if len(cs.tiles) == cs.roomWidth*cs.roomHeight {
				// Fallback to tile index based solidity if available
				if tileIndex >= 0 && tileIndex < len(cs.tiles) {
					if IsSolidTile(cs.tiles[tileIndex]) {
						return true
					}
				}
			}
		}
	}
	
	return false
}

/*
CheckPositionCollision checks if a player at a given position would collide.
Parameters:
  - player: The player instance (for collision box calculation)
  - x, y: Position to test in physics units
Returns true if the player would collide at that position.
*/
func (cs *CollisionService) CheckPositionCollision(player *Player, x, y int) bool {
	// Save current position
	oldX, oldY := player.x, player.y
	
	// Temporarily move to test position
	player.x, player.y = x, y
	box := player.GetCollisionBox()
	
	// Restore position
	player.x, player.y = oldX, oldY
	
	// Check collision
	return cs.CheckBoxCollision(box)
}

/*
IsAvailable returns true if the collision service is properly initialized.
*/
func (cs *CollisionService) IsAvailable() bool {
	return cs.solidityProvider != nil
}

/*
GetRoomDimensions returns the room dimensions in tiles.
*/
func (cs *CollisionService) GetRoomDimensions() (width, height int) {
	return cs.roomWidth, cs.roomHeight
}