package world

import (
	"path/filepath"
	"strings"
	
	"sword/internal/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
)

// TiledRoom adapts a Tiled map to our Room interface
// It embeds BaseRoom and stores the parsed loaded map for reference
type TiledRoom struct {
	*BaseRoom
	loaded *tiled.LoadedMap
}

// NewTiledRoomFromLoadedMap creates a Room from a parsed Tiled LoadedMap
func NewTiledRoomFromLoadedMap(zoneID string, lm *tiled.LoadedMap) *TiledRoom {
	// Initialize base room and copy render layer into our TileMap as indices
	width := lm.TMJ.Width
	height := lm.TMJ.Height
	room := &TiledRoom{
		BaseRoom: NewBaseRoom(zoneID, width, height),
		loaded:   lm,
	}

	// Populate tiles from render layer data, converting GIDs to 0-based indices
	if lm.RenderLayer != nil && len(lm.RenderLayer.Data) == width*height {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				gid := lm.RenderLayer.Data[idx]
				if gid == 0 {
					room.tileMap.Tiles[y][x] = -1
					continue
				}
				baseIndex := gidToTilesetLocalIndex(lm, gid)
				room.tileMap.Tiles[y][x] = baseIndex
			}
		}
	}

	return room
}

// gidToTilesetLocalIndex converts a global id to a 0-based tile index relative to its tileset
// If tileset cannot be determined, returns -1
func gidToTilesetLocalIndex(lm *tiled.LoadedMap, gid uint32) int {
	bestFirst := -1
	bestIdx := -1
	for _, ts := range lm.Tilesets {
		if int(gid) >= ts.FirstGID {
			if ts.FirstGID > bestFirst {
				bestFirst = ts.FirstGID
				bestIdx = int(gid) - ts.FirstGID
			}
		}
	}
	return bestIdx
}

// Draw overrides to render using our tile-based renderer and sprite provider
func (tr *TiledRoom) Draw(screen *ebiten.Image) {
	tr.DrawTiles(screen, tr.getTileSprite)
}

func (tr *TiledRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	tr.DrawTilesWithCamera(screen, tr.getTileSprite, cameraOffsetX, cameraOffsetY)
}

// getTileSprite resolves a tile index to an ebiten image using existing engine helpers
func (tr *TiledRoom) getTileSprite(tileIndex int) *ebiten.Image {
	if engine.GameConfig.UsePlaceholderSprites {
		return engine.GetTileSpriteByType(tileIndex)
	}
	if sprite := engine.LoadSpriteByHex(tileIndex); sprite != nil {
		return sprite
	}
	return engine.GetTileSprite()
}

// Utility to create a stable room id from zone and file path like r01.tmj -> "zone/r01"
func RoomIDFromPath(zoneName, path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return zoneName + "/" + base
}