package world

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/internal/tiled"
)

// tileSpriteResolver resolves tile indices to sprites using engine helpers and Tiled tileset info.
type tileSpriteResolver struct {
	loaded *tiled.LoadedMap
}

func newTileSpriteResolver(lm *tiled.LoadedMap) *tileSpriteResolver {
	return &tileSpriteResolver{loaded: lm}
}

func (r *tileSpriteResolver) Resolve(tileIndex int) *ebiten.Image {
	if engine.GameConfig.UsePlaceholderSprites {
		return engine.GetTileSpriteByType(tileIndex)
	}
	if sprite := engine.LoadSpriteByHex(tileIndex); sprite != nil {
		engine.LogDebug(fmt.Sprintf("Sprite resolved via default LoadSpriteByHex idx=%d", tileIndex))
		return sprite
	}
	// If default loading fails, attempt tileset-name-based mapping to sheet
	if r.loaded != nil && len(r.loaded.Tilesets) > 0 {
		// Prefer last tileset by firstGID proximity to indices we used
		tsName := r.loaded.Tilesets[len(r.loaded.Tilesets)-1].TSX.Name
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