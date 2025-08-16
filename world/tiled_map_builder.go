package world

import (
	"sword/internal/tiled"
)

// buildTileMapFromRenderLayer constructs a TileMap from the Tiled render layer data.
// It returns the tile map and mapping statistics: number of mapped tiles and number of non-empty tiles.
func buildTileMapFromRenderLayer(lm *tiled.LoadedMap) (*TileMap, int, int) {
	width := lm.TMJ.Width
	height := lm.TMJ.Height
	tm := NewTileMap(width, height)

	mapped := 0
	nonEmpty := 0

	if lm.RenderLayer != nil && len(lm.RenderLayer.Data) == width*height {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := y*width + x
				gid := tiled.NormalizeGID(lm.RenderLayer.Data[idx])
				if gid == 0 {
					tm.Tiles[y][x] = -1
					continue
				}
				nonEmpty++
				baseIndex := gidToTilesetLocalIndex(lm, gid)
				if baseIndex >= 0 {
					mapped++
				}
				tm.Tiles[y][x] = baseIndex
			}
		}
	}

	return tm, mapped, nonEmpty
}

// gidToTilesetLocalIndex converts a global id to a 0-based tile index relative to its tileset.
// If tileset cannot be determined, returns -1.
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