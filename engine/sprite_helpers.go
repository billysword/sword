package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/*
GetPlayerSprite returns the appropriate player sprite based on config.
If UsePlaceholderSprites is true, returns a placeholder sprite.
*/
func GetPlayerSprite(facing string) *ebiten.Image {
	if GameConfig.UsePlaceholderSprites {
		return GeneratePlayerPlaceholder()
	}

	sm := GetSpriteManager()
	var sheetName string
	switch facing {
	case "left":
		sheetName = "player_left"
	case "right":
		sheetName = "player_right"
	default:
		sheetName = "player_idle"
	}

	sprite := sm.GetTileByIndex(sheetName, 0)
	if sprite == nil {
		return GeneratePlayerPlaceholder()
	}
	return sprite
}

/*
GetEnemySprite returns the appropriate enemy sprite based on config.
If UsePlaceholderSprites is true, returns a placeholder sprite.
*/
func GetEnemySprite() *ebiten.Image {
	if GameConfig.UsePlaceholderSprites {
		return GenerateEnemyPlaceholder()
	}

	sm := GetSpriteManager()
	sprite := sm.GetTileByIndex("enemy", 0)
	if sprite == nil {
		return GenerateEnemyPlaceholder()
	}
	return sprite
}

/*
GetTileSpriteByType returns the appropriate tile sprite based on config.
If UsePlaceholderSprites is true, returns a placeholder sprite.
*/
func GetTileSpriteByType(tileType int) *ebiten.Image {
	// Map tile types to placeholder types for fallback
	placeholderForType := func(t int) *ebiten.Image {
		switch t {
		case 0x01, 0x02, 0x03: // Ground tiles
			return GenerateTilePlaceholder(PlaceholderTileGround)
		case 0x10, 0x11, 0x12: // Wall tiles
			return GenerateTilePlaceholder(PlaceholderTileWall)
		case 0x20, 0x21: // Platform tiles
			return GenerateTilePlaceholder(PlaceholderTilePlatform)
		case 0x30: // Spike tiles
			return GenerateTilePlaceholder(PlaceholderTileSpike)
		case 0x40, 0x41, 0x42: // Decoration tiles
			return GenerateTilePlaceholder(PlaceholderTileDecoration)
		default:
			return GenerateTilePlaceholder(PlaceholderTileGround)
		}
	}

	if GameConfig.UsePlaceholderSprites {
		return placeholderForType(tileType)
	}

	sm := GetSpriteManager()
	sprite := sm.GetTileByHex("forest", tileType)
	if sprite == nil {
		return placeholderForType(tileType)
	}
	return sprite
}
