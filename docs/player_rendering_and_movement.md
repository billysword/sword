# Player Character Rendering and Movement Guide

## Overview

This document explains how the Player character rendering, movement, and spawn point system works in the sword game engine. This information is essential for adjusting spawn points and avoiding collision issues like getting stuck in walls.

## Player Rendering System

### Core Components

The Player rendering system is located in `entities/player.go:251-294` and handles:

1. **Sprite Selection**: Chooses sprites based on movement state
   - `idle`: When player is stationary
   - `left`: When moving left  
   - `right`: When moving right

2. **Scaling and Positioning**: Applies world and character scaling
   - Uses `engine.GameConfig.TileScaleFactor` for world scaling
   - Uses `engine.GameConfig.CharScaleFactor` for character scaling
   - Converts physics units to screen pixels

3. **Orientation Handling**: Flips sprite horizontally when facing left
   - `facingRight` boolean determines sprite orientation
   - Adjusts positioning when flipped to maintain correct placement

### Rendering Pipeline

```go
// Sprite selection based on movement
switch {
case p.vx > 0:
    sprite = engine.GetPlayerSprite("right")
case p.vx < 0:
    sprite = engine.GetPlayerSprite("left")
default:
    sprite = engine.GetPlayerSprite("idle")
}

// Position calculation
renderX := float64(p.x)*s + cameraOffsetX
renderY := float64(p.y)*s + cameraOffsetY
```

## Movement System

### Physics Architecture

The Player uses a physics-based movement system with two update methods:

1. **`Update()`**: Simple ground collision (legacy)
2. **`UpdateWithTileCollision(room)`**: Advanced tile-based collision system

### Key Physics Properties

Located in `entities/player.go:21-34`:

```go
type Player struct {
    x, y     int  // Position in physics units
    vx, vy   int  // Velocity in physics units per frame
    onGround bool // Whether the player is currently on the ground
    
    // Jump mechanics state
    coyoteTimer     int  // Frames since leaving ground
    jumpBufferTimer int  // Frames since jump was pressed
    isJumping       bool // Currently in a jump
    jumpHeldFrames  int  // How long jump has been held
    
    // Direction state
    facingRight bool // Whether the player is facing right
}
```

### Movement Physics Configuration

All physics settings are centralized in `engine.GameConfig.PlayerPhysics`:

- **Movement**: `MoveSpeed`, `AirControl`, `Friction`, `AirFriction`
- **Jumping**: `JumpPower`, `CoyoteTime`, `JumpBufferTime`, `VariableJumpHeight`
- **Gravity**: `Gravity`, `MaxFallSpeed`, `FastFallMultiplier`
- **Collision**: `CollisionBoxOffsetX/Y`, `CollisionBoxWidth/Height`

### Collision Detection

The advanced collision system (`entities/player_collision.go`) provides:

1. **Tile-Based Collision**: Checks against solid tiles in the room
2. **Stepped Movement**: Prevents tunneling through walls
3. **Separate Axis Resolution**: Handles horizontal and vertical collisions independently

```go
// Collision box calculation
func (p *Player) GetCollisionBox() CollisionBox {
    config := &engine.GameConfig.PlayerPhysics
    spriteWidth := int(float64(config.SpriteWidth) * engine.GameConfig.CharScaleFactor)
    spriteHeight := int(float64(config.SpriteHeight) * engine.GameConfig.CharScaleFactor)
    
    offsetX := int(float64(spriteWidth) * config.CollisionBoxOffsetX)
    offsetY := int(float64(spriteHeight) * config.CollisionBoxOffsetY)
    width := int(float64(spriteWidth) * config.CollisionBoxWidth)
    height := int(float64(spriteHeight) * config.CollisionBoxHeight)
    
    return CollisionBox{
        X: p.x + offsetX, Y: p.y + offsetY,
        Width: width, Height: height,
    }
}
```

## Spawn Point System

### Architecture

The spawn point system is managed by `world.RoomTransitionManager` and handles:

1. **Spawn Point Storage**: Maps room IDs to lists of spawn points
2. **Player Positioning**: Places player at specified spawn points during transitions
3. **Collision Safety**: Ensures spawn points don't place player inside walls

### Spawn Point Structure

Located in `world/room_transition.go:57-63`:

```go
type SpawnPoint struct {
    ID       string `json:"id"`        // Unique identifier
    X        int    `json:"x"`         // Position X in physics units
    Y        int    `json:"y"`         // Position Y in physics units
    FacingID string `json:"facing_id"` // Direction player should face
}
```

### Spawn Point Configuration

Spawn points are defined in JSON files like `resources/rooms/room_transitions.json`:

```json
{
  "rooms": {
    "main": {
      "spawn_points": [
        {"id": "main_spawn", "x": 128, "y": 96, "facing_id": "east"}
      ]
    }
  }
}
```

### Safe Spawn Positioning

The system includes automatic collision checking (`world/room_transition.go:283-310`):

```go
func (rtm *RoomTransitionManager) findNonSolidPosition(roomID string, x, y int) (int, int) {
    // Check if spawn position is inside a solid tile
    // Search nearby tiles for a safe position
    // Return adjusted coordinates or original if no solid collision
}
```

## Debug Tools

### Visual Debug Overlay (F4)

Shows collision and movement visualization:

- **Green collision box**: Player is grounded
- **Yellow collision box**: Player is airborne  
- **Cyan collision box**: Player is jumping
- **Magenta area**: Ground detection zone
- **Gray outline**: Sprite bounds
- **Red vectors**: Velocity indicators

### Debug Information

Press F3 for debug HUD showing:
- Position in tiles and physics units
- Velocity values
- Ground state
- Coyote time and jump buffer timers

## Fixing Spawn Point Issues

### Common Problems

1. **Spawning in walls**: Spawn coordinates place player inside solid tiles
2. **Spawning in air**: Player falls immediately after spawning
3. **Wrong orientation**: Player faces wrong direction after transition

### Solutions

1. **Check tile coordinates**: Convert spawn points to tile positions
   ```
   tileX = spawnX / engine.GetPhysicsUnit()
   tileY = spawnY / engine.GetPhysicsUnit()
   ```

2. **Verify collision box placement**: Account for collision box offset
   - Player collision box is typically smaller than sprite
   - Default offset: 25% from left, 50% from top
   - Default size: 50% of sprite width/height

3. **Test with debug overlay**: Use F4 to visualize collision boxes at spawn points

4. **Use safe positioning**: The system automatically adjusts spawn points if they're in solid tiles

### Best Practices

1. **Place spawn points on solid ground**: Avoid mid-air spawning
2. **Account for collision box size**: Spawn point is top-left of sprite, not collision box
3. **Test all transitions**: Verify spawn points work from all connecting rooms
4. **Use meaningful spawn IDs**: Name spawn points descriptively (e.g., "west_entrance", "upper_platform")

### Debugging Spawn Issues

1. Enable debug overlay (F4) and debug HUD (F3)
2. Check player position after spawning
3. Verify collision box is not overlapping solid tiles
4. Adjust spawn coordinates in room transition JSON files
5. Test movement immediately after spawning

## Physics Units vs Screen Pixels

**Important**: All player positions and spawn points use physics units, not screen pixels.

- **Physics Unit**: Base unit for position/velocity calculations
- **Screen Pixels**: Final rendered position after scaling
- **Conversion**: `screenPixels = physicsUnits * TileScaleFactor`

Default physics unit is typically 32 pixels, making one physics unit equal to one tile.

## Integration Example

To spawn a player safely in a room:

```go
// In your game state
err := roomTransitionMgr.SpawnPlayerInRoom(player, "main", "main_spawn")
if err != nil {
    log.Printf("Spawn failed: %v", err)
}

// Update with collision detection
player.UpdateWithTileCollision(currentRoom)
```

This system ensures the player is positioned correctly and can move without getting stuck in walls.