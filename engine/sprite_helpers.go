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
	
	switch facing {
	case "left":
		return globalLeftSprite
	case "right":
		return globalRightSprite
	default:
		return globalIdleSprite
	}
}

/*
GetEnemySprite returns the appropriate enemy sprite based on config.
If UsePlaceholderSprites is true, returns a placeholder sprite.
*/
func GetEnemySprite() *ebiten.Image {
	if GameConfig.UsePlaceholderSprites {
		return GenerateEnemyPlaceholder()
	}
	
	// Return default enemy sprite (for now, use idle sprite)
	// This should be replaced with actual enemy sprites when available
	return globalIdleSprite
}

/*
GetTileSpriteByType returns the appropriate tile sprite based on config.
If UsePlaceholderSprites is true, returns a placeholder sprite.
*/
func GetTileSpriteByType(tileType int) *ebiten.Image {
	if GameConfig.UsePlaceholderSprites {
		// Map tile types to placeholder types
		switch tileType {
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
	
	// Use existing tile sprite system
	return GetTileSprite()
}