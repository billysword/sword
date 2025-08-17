package world

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/entities"
	"sword/internal/tiled"
)

// TiledRoom adapts a Tiled map to our Room interface
// It embeds BaseRoom and stores the parsed loaded map for reference
type TiledRoom struct {
	*BaseRoom
	loaded         *tiled.LoadedMap
	spriteResolver *tileSpriteResolver
	solidity       *tiledSolidity
}

// Ensure TiledRoom implements entities.TileSolidityProvider
var _ entities.TileSolidityProvider = (*TiledRoom)(nil)

// NewTiledRoomFromLoadedMap creates a Room from a parsed Tiled LoadedMap
func NewTiledRoomFromLoadedMap(zoneID string, lm *tiled.LoadedMap) *TiledRoom {
	// Debug layer comparison for the main room
	if zoneID == "cradle/r01" {
		DebugCompareLayers(lm)
	}
	// Initialize base room and copy render layer into our TileMap as indices
	width := lm.TMJ.Width
	height := lm.TMJ.Height
	room := &TiledRoom{
		BaseRoom:       NewBaseRoom(zoneID, width, height),
		loaded:         lm,
		spriteResolver: newTileSpriteResolver(lm),
		solidity:       newTiledSolidity(lm),
	}

	// Integration logging: tilesets and mapping
	u := engine.GetPhysicsUnit()
	for _, ts := range lm.Tilesets {
		img := ts.TSX.Image.Source
		sz := fmt.Sprintf("%dx%d", ts.TSX.TileWidth, ts.TSX.TileHeight)
		engine.LogInfo("Tileset '" + ts.TSX.Name + "' img='" + img + "' tile=" + sz +
			fmt.Sprintf(" cols=%d count=%d firstGID=%d", ts.TSX.Columns, ts.TSX.TileCount, ts.FirstGID))
		if ts.TSX.TileWidth != u || ts.TSX.TileHeight != u {
			engine.LogWarn(fmt.Sprintf("Tileset '%s' tile size (%dx%d) differs from physics unit %d; verify sprite indices and scaling.", ts.TSX.Name, ts.TSX.TileWidth, ts.TSX.TileHeight, u))
		}
	}

	// Build the tile map using a focused helper and capture coverage stats
	tm, mapped, nonEmpty := buildTileMapFromRenderLayer(lm)
	room.tileMap = tm

	pct := 100.0
	if nonEmpty > 0 {
		pct = float64(mapped) * 100.0 / float64(nonEmpty)
	}
	engine.LogInfo(fmt.Sprintf("Tiled room '%s' tiles mapped: %d/%d (%.1f%%)", zoneID, mapped, nonEmpty, pct))

	return room
}

// Draw overrides to render using our tile-based renderer and sprite provider
func (tr *TiledRoom) Draw(screen *ebiten.Image) {
	tr.DrawTiles(screen, tr.spriteResolver.Resolve)
}

func (tr *TiledRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	tr.DrawTilesWithCamera(screen, tr.spriteResolver.Resolve, cameraOffsetX, cameraOffsetY)
}

// IsSolidAtFlatIndex uses the collision layer if present; otherwise falls back to tile properties from render layer
func (tr *TiledRoom) IsSolidAtFlatIndex(index int) bool {
	if tr.solidity == nil {
		return false
	}
	return tr.solidity.IsSolidAtFlatIndex(index)
}

// FindFloorAtX finds the floor Y position at the given X coordinate.
// Returns the Y position in physics units where the player should stand.
func (tr *TiledRoom) FindFloorAtX(x int) int {
	return findFloorAtX(tr.tileMap, tr.IsSolidAtFlatIndex, x)
}
