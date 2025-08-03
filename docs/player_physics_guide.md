# Player Physics Configuration Guide

## Overview

The player physics system has been completely overhauled to make it easy to adjust and work with. All physics parameters are now centralized in the `PlayerPhysicsConfig` struct, with visual debugging tools to help you find the perfect settings.

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

### 3. Tile-Based Collision Detection

The player now uses proper tile-based collision detection:
- Collision box can be smaller than sprite for precise platforming
- Separate horizontal and vertical collision resolution
- Ground detection uses a small area below the collision box
- Works with the tileset's solid tiles

### 4. Advanced Jump Mechanics

The new system includes modern platformer features:
- **Coyote Time**: Can still jump for a few frames after leaving a platform
- **Jump Buffering**: Jump input is remembered if pressed just before landing
- **Variable Jump Height**: Release jump early for shorter jumps
- **Fast Fall**: Hold down to fall faster

## Configuration Examples

### Adjusting Player Size Relative to Tiles

Edit the configuration in `engine/config.go`:

```go
// Make player smaller relative to tiles
PlayerPhysics: PlayerPhysicsConfig{
    SpriteWidth:  32,
    SpriteHeight: 32,
    CollisionBoxWidth: 0.4,   // 40% of sprite width
    CollisionBoxHeight: 0.6,  // 60% of sprite height
    CollisionBoxOffsetX: 0.3, // Center the narrower collision box
    CollisionBoxOffsetY: 0.4, // Position collision box at feet
    // ... other settings
}
```

### Creating a "Floaty" Jump Feel

```go
PlayerPhysics: PlayerPhysicsConfig{
    // ... sprite settings
    Gravity: 0.5,               // Reduced gravity
    MaxFallSpeed: 8,            // Lower terminal velocity
    JumpPower: 10,              // Higher initial jump
    AirControl: 0.9,            // More air control
    VariableJumpHeight: true,
    MinJumpHeight: 0.3,         // Very short tap jumps possible
    // ... other settings
}
```

### Creating a "Heavy" Character Feel

```go
PlayerPhysics: PlayerPhysicsConfig{
    // ... sprite settings
    Gravity: 2,                 // Higher gravity
    MaxFallSpeed: 20,           // Faster falling
    JumpPower: 12,              // Need more power to jump
    Friction: 3,                // Stops quickly
    AirControl: 0.3,            // Limited air control
    FastFallMultiplier: 2.0,    // Very fast when holding down
    // ... other settings
}
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
2. **Test with different room layouts**: Collision box size affects platforming feel
3. **Consider tile size**: Adjust collision box to match your level design
4. **Use presets**: Start with DefaultConfig() and adjust from there
5. **Test edge cases**: Make sure your collision box works well with platforms and walls

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

## Quick Reference

Common adjustments in `engine/config.go`:

```go
// In DefaultConfig() or your custom config:
PlayerPhysics: PlayerPhysicsConfig{
    // Sprite size (match your actual sprite)
    SpriteWidth:  32,
    SpriteHeight: 32,
    
    // Collision box (fractions of sprite size)
    CollisionBoxOffsetX: 0.25,  // 25% from left
    CollisionBoxOffsetY: 0.5,   // 50% from top
    CollisionBoxWidth:   0.5,   // 50% of sprite width
    CollisionBoxHeight:  0.5,   // 50% of sprite height
    
    // Movement
    MoveSpeed:   2,             // Tiles per second (roughly)
    JumpPower:   8,             // Initial jump velocity
    AirControl:  0.7,           // 70% control in air
    
    // Physics
    Gravity:            1,      // Acceleration per frame
    MaxFallSpeed:       12,     // Terminal velocity
    Friction:           1,      // Ground friction
    
    // Advanced
    CoyoteTime:         6,      // Frames of grace period
    JumpBufferTime:     10,     // Frames to remember jump input
    VariableJumpHeight: true,   // Can control jump height
    MinJumpHeight:      0.4,    // Minimum jump is 40% of full
}
```