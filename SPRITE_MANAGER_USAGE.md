# Sprite Manager Usage Guide

The new sprite management system provides a clean, configuration-based approach to loading and accessing sprite sheets. This replaces the previous direct coupling to game states and provides easy access to tiles by hex index.

## Basic Usage

### Loading Sprite Sheets

```go
// Initialize sprite manager (done once at startup)
engine.InitSpriteManager()
sm := engine.GetSpriteManager()

// Load a sprite sheet
err := sm.LoadSpriteSheet("forest", forestTileImage, 16, 16)
if err != nil {
    panic(err)
}
```

### Accessing Tiles

```go
// Get tile by hex index (most common usage)
sprite := engine.LoadSpriteByHex(0x05)  // Gets tile #5

// Get tile from specific sheet
sprite := engine.LoadTileFromSheet("forest", 0x0A)  // Gets tile #10 from forest sheet

// Get tile by decimal index
sm := engine.GetSpriteManager()
sprite := sm.GetTileByIndex("forest", 15)
```

## Settings Debug Menu Integration

The settings/debug tile menu now shows:

1. **Actual Sprite Visualization**: Each tile displays the real sprite from the sprite manager
2. **Hex Values**: Clear hex display (e.g., `0x05`)  
3. **Decimal Indices**: Easy reference (e.g., `#5`)
4. **Coordinates**: Tile position (e.g., `(12,8)`)
5. **Sprite Sheet Info**: Which sheets are loaded and available

## Architecture Benefits

### Before
- Settings state was tightly coupled to InGameState
- Static room references that could become stale
- No visual feedback of actual sprites
- Hard to debug tile index issues

### After
- Clean separation of concerns
- Configuration-based sprite loading
- Visual representation of actual sprites
- Easy access pattern: `LoadSpriteByHex(0x05)`
- Cached sprites for performance
- Support for multiple sprite sheets

## Configuration

The sprite manager automatically calculates tile layout based on:
- **Tile Size**: Width/height of individual tiles (usually 16x16)
- **Sheet Dimensions**: Automatically calculated from image size
- **Indexing**: Left-to-right, top-to-bottom starting from 0

## Example Usage in Game Code

```go
// In room initialization
func (room *MyRoom) setupTiles() {
    // Just reference tiles by hex - no complex setup needed
    room.setTile(x, y, 0x05)  // Dirt tile
    room.setTile(x, y, 0x12)  // Platform tile
}

// In rendering code
func (room *MyRoom) drawTile(tileIndex int) {
    sprite := engine.LoadSpriteByHex(tileIndex)
    if sprite != nil {
        // Render sprite...
    }
}
```

## Debug Features

- **Settings Menu**: Press `P` then `S` to view all tiles with hex indices
- **Live Data**: Shows current room's actual tile map
- **Sprite Sheets**: Displays which sheets are loaded
- **Visual Feedback**: See the actual sprites, not just colored rectangles

This system makes it much easier to:
1. Know which hex index you're using for easy fixes
2. See tiles in-game easily with visual representation
3. Debug tile placement issues
4. Add new sprite sheets without code changes