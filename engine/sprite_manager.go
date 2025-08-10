package engine

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
SpriteSheet represents a loaded sprite sheet with configuration.
Contains the source image and metadata about how tiles are arranged.
*/
type SpriteSheet struct {
	Image        *ebiten.Image // The loaded sprite sheet image
	TileWidth    int           // Width of each tile in pixels
	TileHeight   int           // Height of each tile in pixels
	TilesPerRow  int           // Number of tiles per row
	TotalTiles   int           // Total number of tiles in the sheet
	Name         string        // Identifier for this sprite sheet
}

/*
SpriteManager manages multiple sprite sheets and provides easy access to individual tiles.
Centralizes sprite loading and caching for efficient tile-based rendering.
*/
type SpriteManager struct {
	sheets map[string]*SpriteSheet           // Loaded sprite sheets by name
	cache  map[string]map[int]*ebiten.Image  // Cached individual tiles [sheetName][tileIndex]
}

// Global sprite manager instance
var globalSpriteManager *SpriteManager

/*
InitSpriteManager initializes the global sprite manager.
Should be called once during game startup.
*/
func InitSpriteManager() {
	globalSpriteManager = &SpriteManager{
		sheets: make(map[string]*SpriteSheet),
		cache:  make(map[string]map[int]*ebiten.Image),
	}
}

/*
GetSpriteManager returns the global sprite manager instance.
Creates a new one if not initialized.
*/
func GetSpriteManager() *SpriteManager {
	if globalSpriteManager == nil {
		InitSpriteManager()
	}
	return globalSpriteManager
}

/*
LoadSpriteSheet loads a sprite sheet and registers it with the manager.
Automatically calculates tile layout based on the provided configuration.

Parameters:
  - name: Unique identifier for this sprite sheet
  - image: The sprite sheet image
  - tileWidth: Width of each individual tile
  - tileHeight: Height of each individual tile

Returns error if loading fails.
*/
func (sm *SpriteManager) LoadSpriteSheet(name string, image *ebiten.Image, tileWidth, tileHeight int) error {
	if image == nil {
		return fmt.Errorf("sprite sheet image is nil for %s", name)
	}

	bounds := image.Bounds()
	tilesPerRow := bounds.Dx() / tileWidth
	tilesPerCol := bounds.Dy() / tileHeight
	totalTiles := tilesPerRow * tilesPerCol

	sheet := &SpriteSheet{
		Image:       image,
		TileWidth:   tileWidth,
		TileHeight:  tileHeight,
		TilesPerRow: tilesPerRow,
		TotalTiles:  totalTiles,
		Name:        name,
	}

	sm.sheets[name] = sheet
	sm.cache[name] = make(map[int]*ebiten.Image)

	LogSprite(fmt.Sprintf("Loaded sprite sheet '%s': %dx%d tiles (%d total)", 
		name, tilesPerRow, tilesPerCol, totalTiles))

	return nil
}

/*
GetTileByIndex returns a tile from a sprite sheet by its index.
Tiles are indexed left-to-right, top-to-bottom starting from 0.
Uses caching for efficient repeated access.

Parameters:
  - sheetName: Name of the sprite sheet
  - index: Tile index (0-based)

Returns the tile image, or nil if not found.
*/
func (sm *SpriteManager) GetTileByIndex(sheetName string, index int) *ebiten.Image {
	sheet, exists := sm.sheets[sheetName]
	if !exists {
		LogSprite(fmt.Sprintf("Sprite sheet '%s' not found", sheetName))
		return nil
	}

	if index < 0 || index >= sheet.TotalTiles {
		LogSprite(fmt.Sprintf("Tile index %d out of range for sheet '%s' (0-%d)", 
			index, sheetName, sheet.TotalTiles-1))
		return nil
	}

	// Check cache first
	if tile, cached := sm.cache[sheetName][index]; cached {
		return tile
	}

	// Calculate tile position
	col := index % sheet.TilesPerRow
	row := index / sheet.TilesPerRow
	x := col * sheet.TileWidth
	y := row * sheet.TileHeight

	// Extract tile from sheet
	rect := image.Rect(x, y, x+sheet.TileWidth, y+sheet.TileHeight)
	tile := sheet.Image.SubImage(rect).(*ebiten.Image)

	// Cache the tile
	sm.cache[sheetName][index] = tile

	return tile
}

// MapTilesetToSheet attempts to pick a sprite sheet for a Tiled tileset by name convention
// Example: tileset name 'forest' -> sheet 'forest'
func MapTilesetToSheet(tilesetName string) string {
	if tilesetName == "" {
		return ""
	}
	// Direct match
	if _, ok := GetSpriteManager().sheets[tilesetName]; ok {
		return tilesetName
	}
	// Fallbacks
	candidates := []string{"forest", "tiles", "default"}
	for _, name := range candidates {
		if _, ok := GetSpriteManager().sheets[name]; ok {
			return name
		}
	}
	return ""
}

/*
GetTileByHex returns a tile from a sprite sheet by its hexadecimal index.
Convenient wrapper for GetTileByIndex that accepts hex values.

Parameters:
  - sheetName: Name of the sprite sheet
  - hexIndex: Tile index in hexadecimal (e.g., 0x05, 0x1A)

Returns the tile image, or nil if not found.
*/
func (sm *SpriteManager) GetTileByHex(sheetName string, hexIndex int) *ebiten.Image {
	return sm.GetTileByIndex(sheetName, hexIndex)
}

/*
GetTileInfo returns information about a tile from a sprite sheet.
Useful for debugging and displaying tile metadata.

Parameters:
  - sheetName: Name of the sprite sheet
  - index: Tile index (0-based)

Returns formatted string with tile information.
*/
func (sm *SpriteManager) GetTileInfo(sheetName string, index int) string {
	sheet, exists := sm.sheets[sheetName]
	if !exists {
		return fmt.Sprintf("Sheet '%s' not found", sheetName)
	}

	if index < 0 || index >= sheet.TotalTiles {
		return fmt.Sprintf("Index %d out of range for sheet '%s'", index, sheetName)
	}

	col := index % sheet.TilesPerRow
	row := index / sheet.TilesPerRow
	
	return fmt.Sprintf("Sheet: %s, Index: %d (0x%02X), Position: (%d,%d)", 
		sheetName, index, index, col, row)
}

/*
ListSheets returns a list of all loaded sprite sheet names.
Useful for debugging and introspection.
*/
func (sm *SpriteManager) ListSheets() []string {
	var names []string
	for name := range sm.sheets {
		names = append(names, name)
	}
	return names
}

/*
GetSheetInfo returns information about a loaded sprite sheet.
*/
func (sm *SpriteManager) GetSheetInfo(sheetName string) string {
	sheet, exists := sm.sheets[sheetName]
	if !exists {
		return fmt.Sprintf("Sheet '%s' not found", sheetName)
	}

	return fmt.Sprintf("Sheet '%s': %dx%d tiles (%dx%d px each), Total: %d", 
		sheet.Name, sheet.TilesPerRow, sheet.TotalTiles/sheet.TilesPerRow,
		sheet.TileWidth, sheet.TileHeight, sheet.TotalTiles)
}

// Convenience functions for global access

/*
LoadSpriteByHex is a global convenience function for loading tiles by hex index.
Uses the default sprite sheet (first loaded or specified).
*/
func LoadSpriteByHex(hexIndex int) *ebiten.Image {
	sm := GetSpriteManager()
	
	// Use the first available sheet if only one loaded
	if len(sm.sheets) == 1 {
		for name := range sm.sheets {
			return sm.GetTileByHex(name, hexIndex)
		}
	}
	
	// Try common sheet names
	commonNames := []string{"tiles", "forest", "main", "default"}
	for _, name := range commonNames {
		if tile := sm.GetTileByHex(name, hexIndex); tile != nil {
			return tile
		}
	}
	
	LogSprite(fmt.Sprintf("No suitable sprite sheet found for hex index 0x%02X", hexIndex))
	return nil
}

/*
LoadTileFromSheet loads a tile from a specific sheet by hex index.
*/
func LoadTileFromSheet(sheetName string, hexIndex int) *ebiten.Image {
	return GetSpriteManager().GetTileByHex(sheetName, hexIndex)
}