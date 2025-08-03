# Player Physics Configuration Guide

## Overview

The player physics system has been completely overhauled to make it easy to adjust and work with. All physics parameters are now centralized in the `PlayerPhysicsConfig` struct, with runtime adjustment capabilities and visual debugging tools.

## Key Features

### 1. Centralized Physics Configuration

All player physics settings are now in `engine.GameConfig.PlayerPhysics`:

```go
type PlayerPhysicsConfig struct {
    // Sprite dimensions
    SpriteWidth  int     // Base sprite width in pixels (default: 32)
    SpriteHeight int     // Base sprite height in pixels (default: 32)
    
    // Collision box (relative to sprite dimensions)
    CollisionBoxOffsetX float64 // Offset from sprite left edge (0-1)
    CollisionBoxOffsetY float64 // Offset from sprite top edge (0-1)
    CollisionBoxWidth   float64 // Width as fraction of sprite width (0-1)
    CollisionBoxHeight  float64 // Height as fraction of sprite height (0-1)
    
    // Movement physics
    MoveSpeed   int     // Horizontal movement speed
    JumpPower   int     // Initial jump velocity
    AirControl  float64 // Movement control while airborne (0-1)
    Friction    int     // Ground friction
    AirFriction int     // Air friction
    
    // Advanced jump mechanics
    CoyoteTime         int  // Frames after leaving ground where jump is allowed
    JumpBufferTime     int  // Frames to buffer jump input before landing
    VariableJumpHeight bool // Allow controlling jump height by release timing
    MinJumpHeight      float64 // Minimum jump height as fraction of full jump
    
    // Gravity and falling
    Gravity            int     // Gravity acceleration per frame
    MaxFallSpeed       int     // Terminal velocity
    FastFallMultiplier float64 // Gravity multiplier when holding down
}
```

### 2. Visual Debug Tools

Press **F4** to toggle debug overlay, which shows:
- Player collision box (green when grounded, yellow when airborne, cyan when jumping)
- Ground detection area (magenta box below collision box)
- Sprite bounds (gray outline)
- Velocity vectors
- Physics state information

### 3. Runtime Physics Tuner

Press **F9** to activate the Physics Tuner:
- **Tab**: Switch between parameters
- **Up/Down Arrows**: Adjust values
- **Shift**: Fine adjustment (0.1 step)
- **Ctrl**: Coarse adjustment (10.0 step)

Adjustable parameters:
- Move Speed
- Jump Power
- Gravity
- Max Fall Speed
- Friction
- Air Control
- Collision Box dimensions and offset
- Coyote Time
- Jump Buffer Time
- Character Scale Factor
- Tile Scale Factor

### 4. Tile-Based Collision Detection

The player now uses proper tile-based collision detection:
- Collision box can be smaller than sprite for precise platforming
- Separate horizontal and vertical collision resolution
- Ground detection uses a small area below the collision box
- Works with the tileset's solid tiles

### 5. Advanced Jump Mechanics

The new system includes modern platformer features:
- **Coyote Time**: Can still jump for a few frames after leaving a platform
- **Jump Buffering**: Jump input is remembered if pressed just before landing
- **Variable Jump Height**: Release jump early for shorter jumps
- **Fast Fall**: Hold down to fall faster

## Usage Examples

### Adjusting Player Size Relative to Tiles

```go
// Make player smaller relative to tiles
config := &engine.GameConfig.PlayerPhysics
config.CollisionBoxWidth = 0.4   // 40% of sprite width
config.CollisionBoxHeight = 0.6  // 60% of sprite height
config.CollisionBoxOffsetX = 0.3 // Center the narrower collision box
config.CollisionBoxOffsetY = 0.4 // Position collision box at feet
```

### Creating a "Floaty" Jump Feel

```go
config.Gravity = 0.5            // Reduced gravity
config.MaxFallSpeed = 8         // Lower terminal velocity
config.JumpPower = 10           // Higher initial jump
config.AirControl = 0.9         // More air control
config.VariableJumpHeight = true
config.MinJumpHeight = 0.3      // Very short tap jumps possible
```

### Creating a "Heavy" Character Feel

```go
config.Gravity = 2              // Higher gravity
config.MaxFallSpeed = 20        // Faster falling
config.JumpPower = 12           // Need more power to jump
config.Friction = 3             // Stops quickly
config.AirControl = 0.3         // Limited air control
config.FastFallMultiplier = 2.0 // Very fast when holding down
```

## Debug HUD Information

The debug HUD (F3) now shows a dedicated "PLAYER PHYSICS" section with:
- Sprite dimensions
- Collision box size (as percentage)
- Movement speed and jump power
- Gravity and max fall speed
- Coyote time and jump buffer frames

## Best Practices

1. **Start with visual debugging**: Enable F4 to see collision boxes
2. **Use the Physics Tuner**: F9 for real-time adjustments
3. **Test with different room layouts**: Collision box size affects platforming feel
4. **Consider tile size**: Adjust collision box to match your level design
5. **Save good configurations**: Note down values that feel good

## Integration with Existing Code

The system maintains backward compatibility:
- `player.Update()` still works but uses simple ground collision
- `player.UpdateWithTileCollision(room)` uses the new tile-based system
- All existing player methods work unchanged

To use tile-based collision in your game state:
```go
// Instead of:
player.Update()

// Use:
player.UpdateWithTileCollision(currentRoom)
```