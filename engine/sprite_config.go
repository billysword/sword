package engine

import (
	"sword/assets/tilesets"
	"sword/resources/images/platformer"
)

// SpriteSheetConfig defines configuration for a sprite sheet.
// ImageData contains the raw bytes for the sheet image.
type SpriteSheetConfig struct {
	Name       string
	ImageData  []byte
	TileWidth  int
	TileHeight int
}

// SpriteSheetConfigs lists all sprite sheets to load at startup.
var SpriteSheetConfigs = []SpriteSheetConfig{
	{
		Name:       "player_left",
		ImageData:  platformer.Left_png,
		TileWidth:  165,
		TileHeight: 205,
	},
	{
		Name:       "player_right",
		ImageData:  platformer.Right_png,
		TileWidth:  165,
		TileHeight: 205,
	},
	{
		Name:       "player_idle",
		ImageData:  platformer.MainChar_png,
		TileWidth:  165,
		TileHeight: 205,
	},
	{
		Name:       "background",
		ImageData:  platformer.Background_png,
		TileWidth:  1920,
		TileHeight: 1080,
	},
	{
		Name:       "forest",
		ImageData:  tilesets.ForestTiles_png,
		TileWidth:  16,
		TileHeight: 16,
	},
}
