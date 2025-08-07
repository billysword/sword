# Placeholder Sprite System

## Overview

The placeholder sprite system provides low-fidelity geometric sprites that can be used during development while waiting for final art assets. This system generates simple, colorful shapes to represent different game elements.

## Features

- **Automatic Generation**: Sprites are generated programmatically using simple shapes and colors
- **Cached**: Generated sprites are cached to avoid regenerating the same sprites
- **Easy Integration**: Works seamlessly with the existing sprite system
- **Runtime Toggle**: Can be enabled/disabled at runtime through settings

## Placeholder Types

The system supports the following placeholder types:

### Characters
- **Player**: Blue rectangle body with a head and eyes
- **Enemy**: Red diamond shape with yellow eyes

### Tiles
- **Ground**: Brown rectangle with horizontal texture lines
- **Wall**: Gray bricks pattern
- **Platform**: Light brown platform with highlight edge
- **Spike**: Dark gray triangular spikes
- **Decoration**: Simple flower/plant sprite

### Other
- **Projectile**: Yellow circle
- **Item**: Gold star shape
- **Background**: Blue gradient

## Usage

### Command Line

Run the game with placeholder sprites:
```bash
./game -placeholders
```

### In-Game Toggle

1. Press `Escape` to open the pause menu
2. Select "Settings"
3. Navigate to the "Developer" tab
4. Toggle "Use Placeholder Sprites"

### Programmatic Usage

```go
// Enable placeholder sprites
engine.GameConfig.UsePlaceholderSprites = true

// Generate specific placeholders
playerSprite := engine.GeneratePlayerPlaceholder()
enemySprite := engine.GenerateEnemyPlaceholder()
groundTile := engine.GenerateTilePlaceholder(engine.PlaceholderTileGround)

// Generate custom size placeholders
customSprite := engine.GetPlaceholderGenerator().GeneratePlaceholder(
    engine.PlaceholderPlayer, 
    64, 64, // width, height
)
```

## Implementation Details

The placeholder system is implemented in `/engine/placeholder_sprites.go` and integrates with:

- **Sprite Manager**: Provides placeholder sprites when enabled
- **Player Entity**: Uses `GetPlayerSprite()` which checks the placeholder flag
- **Enemy Entities**: Uses `GetEnemySprite()` which returns placeholder when enabled
- **Tile Rendering**: Uses `GetTileSpriteByType()` which maps tile IDs to placeholder types

## Benefits

1. **Rapid Prototyping**: Start testing gameplay immediately without waiting for art
2. **Clear Visual Distinction**: Different colors and shapes make it easy to identify game elements
3. **Performance Testing**: Lightweight sprites help test performance without texture overhead
4. **Parallel Development**: Artists can work on final assets while programmers implement features

## Future Enhancements

- Add more placeholder types (NPCs, power-ups, etc.)
- Support for animated placeholders
- Customizable color schemes
- Export placeholder sprites as PNG files for reference