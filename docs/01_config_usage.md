# Game Configuration Usage Guide

## Overview
All game settings are now centralized in `gamestate/config.go`. This makes it easy to adjust the game's appearance and behavior without hunting through multiple files.

## Quick Examples

### 1. Changing Zoom Level (Player Size Relative to Tiles)

To make the player appear larger or smaller relative to tiles:

```go
// In gamestate/config.go, adjust these values:

// For zoomed out view (player ~3 tiles wide):
TileScaleFactor: 1.0,  // Smaller tiles
CharScaleFactor: 0.4,  // Smaller character

// For zoomed in view (player ~2 tiles wide):
TileScaleFactor: 2.0,  // Larger tiles  
CharScaleFactor: 0.5,  // Slightly larger character
```

### 2. Changing Window Size

```go
// Adjust window dimensions:
WindowWidth:  800,   // Width in pixels
WindowHeight: 450,   // Height in pixels
```

### 3. Adjusting Player Movement

```go
// Make player faster/slower:
PlayerMoveSpeed: 3,    // Horizontal speed multiplier
PlayerJumpPower: 8,    // Jump height
PlayerFriction:  2,    // How quickly player stops
```

### 4. Camera Settings

```go
// Smooth camera following:
CameraSmoothing: 0.1,   // 0.1 = responsive, 0.9 = very smooth

// Dead zones (player can move without camera following):
CameraDeadZoneX: 0.25,  // 25% of screen width
CameraDeadZoneY: 0.16,  // 16% of screen height

// Parallax background scrolling:
ParallaxFactor: 0.3,    // Background moves at 30% of camera speed
```

### 5. Room Size

```go
// Change the explorable area:
RoomWidthTiles:  120,   // Room width in tiles
RoomHeightTiles: 60,    // Room height in tiles
GroundLevel:     44,    // Where the main ground is
```

## Using Preset Configurations

The config file includes preset configurations:

```go
// For zoomed out metroidvania view:
GameConfig = DefaultConfig()

// For more zoomed in view:
GameConfig = ZoomedInConfig()
```

## Creating Custom Configurations

You can create your own configuration function:

```go
func MyCustomConfig() *Config {
    config := DefaultConfig()
    
    // Customize specific values
    config.WindowWidth = 1024
    config.WindowHeight = 576
    config.TileScaleFactor = 1.5
    config.PlayerMoveSpeed = 5
    
    return config
}
```

Then use it:
```go
GameConfig = MyCustomConfig()
```

## Live Adjustments

Since `GameConfig` is a global variable, you can even adjust values at runtime:

```go
// In any game state or update function:
if inpututil.IsKeyJustPressed(ebiten.KeyPlus) {
    GameConfig.TileScaleFactor += 0.1
}
```

## Important Notes

1. **Physics Unit**: The physics unit is calculated from `TileSize * TileScaleFactor`. Changing these affects all physics calculations.

2. **Consistency**: When changing scale factors, you may need to adjust movement speeds and physics values to maintain gameplay feel.

3. **Performance**: Larger rooms and smaller tile scales mean more tiles to render, which can impact performance.

4. **Aspect Ratio**: The default configurations maintain a 16:9 aspect ratio, which is standard for modern displays.