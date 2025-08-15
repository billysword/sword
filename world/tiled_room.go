package world

import (
	"fmt"
	"path/filepath"
	"strings"
	
	"sword/internal/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/entities"
)

// TiledRoom adapts a Tiled map to our Room interface
// It embeds BaseRoom and stores the parsed loaded map for reference
type TiledRoom struct {
	*BaseRoom
	loaded *tiled.LoadedMap
}

// Ensure TiledRoom implements entities.TileSolidityProvider
var _ entities.TileSolidityProvider = (*TiledRoom)(nil)

// NewTiledRoomFromLoadedMap creates a Room from a parsed Tiled LoadedMap
func NewTiledRoomFromLoadedMap(zoneID string, lm *tiled.LoadedMap) *TiledRoom {
	// Initialize base room and copy render layer into our TileMap as indices
	width := lm.TMJ.Width
	height := lm.TMJ.Height
	room := &TiledRoom{
		BaseRoom: NewBaseRoom(zoneID, width, height),
		loaded:   lm,
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

	// Track mapping coverage
	mapped := 0
	nonEmpty := 0

	// Populate tiles from render layer data, converting GIDs to 0-based indices
	if lm.RenderLayer != nil && len(lm.RenderLayer.Data) == width*height {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				gid := tiled.NormalizeGID(lm.RenderLayer.Data[idx])
				if gid == 0 {
					room.tileMap.Tiles[y][x] = -1
					continue
				}
				nonEmpty++
				baseIndex := gidToTilesetLocalIndex(lm, gid)
				if baseIndex >= 0 {
					mapped++
				}
				room.tileMap.Tiles[y][x] = baseIndex
			}
		}
	}

	engine.LogInfo(fmt.Sprintf("Tiled room '%s' tiles mapped: %d/%d (%.1f%%)", zoneID, mapped, nonEmpty, percent(mapped, nonEmpty)))

	return room
}

// gidToTilesetLocalIndex converts a global id to a 0-based tile index relative to its tileset
// If tileset cannot be determined, returns -1
func gidToTilesetLocalIndex(lm *tiled.LoadedMap, gid uint32) int {
	gid = tiled.NormalizeGID(gid)
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

func percent(a, b int) float64 {
	if b <= 0 {
		return 100
	}
	return float64(a) * 100.0 / float64(b)
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
		engine.LogDebug(fmt.Sprintf("Sprite resolved via default LoadSpriteByHex idx=%d", tileIndex))
		return sprite
	}
	// If default loading fails, attempt tileset-name-based mapping to sheet
	if tr.loaded != nil && len(tr.loaded.Tilesets) > 0 {
		// Prefer last tileset by firstGID proximity to indices we used
		tsName := tr.loaded.Tilesets[len(tr.loaded.Tilesets)-1].TSX.Name
		sheet := engine.MapTilesetToSheet(tsName)
		if sheet != "" {
			if spr := engine.LoadTileFromSheet(sheet, tileIndex); spr != nil {
				engine.LogDebug(fmt.Sprintf("Mapped Tiled tileset '%s' -> sheet '%s' for index %d", tsName, sheet, tileIndex))
				return spr
			}
		}
		engine.LogWarn(fmt.Sprintf("No sprite sheet mapping for tileset '%s' index %d; using fallback tile sheet", tsName, tileIndex))
	}
	engine.LogWarn(fmt.Sprintf("Falling back to engine.GetTileSprite for idx=%d", tileIndex))
	return engine.GetTileSprite()
}

// IsSolidAtFlatIndex uses the collision layer if present; otherwise falls back to tile properties from render layer
func (tr *TiledRoom) IsSolidAtFlatIndex(index int) bool {
	if tr.loaded == nil {
		return false
	}
	// Prefer explicit collision layer
	if tr.loaded.CollisionLayer != nil && index >= 0 && index < len(tr.loaded.CollisionLayer.Data) {
		return tr.loaded.CollisionLayer.Data[index] != 0
	}
	// Fallback: derive from render layer tile properties
	if tr.loaded.RenderLayer != nil && index >= 0 && index < len(tr.loaded.RenderLayer.Data) {
		gid := tiled.NormalizeGID(tr.loaded.RenderLayer.Data[index])
		if props, ok := tr.loaded.PropertiesForGID(gid); ok {
			return props.Solid
		}
	}
	return false
}

// Utility to create a stable room id from zone and file path like r01.tmj -> "zone/r01"
func RoomIDFromPath(zoneName, path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return zoneName + "/" + base
}