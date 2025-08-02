# Enemies System Implementation

> **📁 Documentation Moved**: This file has been moved to [docs/enemies_01.md](docs/enemies_01.md). Please use the organized documentation in the [docs/](docs/) folder.
> 
> **⚠️ Note**: This describes the original enemy implementation. The current system uses an interface-based architecture documented in [docs/enemies_interface_02.md](docs/enemies_interface_02.md).

This document describes the implementation of the enemies system added to the game.

## Overview

An enemies slice has been added to track other entities in the game that behave similar to the player character but with AI-controlled movement instead of input-controlled movement.

## Components Added

### 1. Enemy Entity (`entities/enemy.go`)

**Key Features:**
- Similar physics system to the Player (gravity, friction, ground collision)
- AI-controlled movement with patrol behavior
- Same sprite rendering system as the Player
- Support for camera-offset rendering

**AI Behavior:**
- Enemies patrol back and forth within a configurable range from their spawn point
- Random direction changes every 1-4 seconds
- Movement speed is half the player's movement speed
- 70% chance to change direction when timer expires

**Movement Properties:**
- `moveDirection`: -1 (left), 0 (stationary), 1 (right)
- `moveTimer`: Frames until next direction change
- `moveSpeed`: Movement speed in physics units (half player speed)
- `patrolRange`: Maximum distance from spawn point (200 physics units)
- `spawnX`: Original spawn position for patrol boundary checking

### 2. Game State Integration (`states/ingame_state.go`)

**Enemies Slice:**
- Added `enemies []*entities.Enemy` field to `InGameState`
- Initialized as empty slice in `NewInGameState()`

**Game Loop Integration:**
- Enemy updates in the main game loop (after player update)
- Enemy rendering with camera offset (before player to render behind)
- Enemy count displayed in debug information

**Test Enemy Spawning:**
- 3 enemies spawned at different positions on game state entry
- Spawned at positions: 300, 600, and 900 physics units horizontally
- All spawned at ground level

### 3. Enemy Management Methods

**Helper Methods Added:**
- `AddEnemy(x, y int)`: Spawn a new enemy at specified position
- `RemoveEnemy(enemy *Enemy)`: Remove specific enemy from slice
- `ClearEnemies()`: Remove all enemies (for room transitions)
- `GetEnemies()`: Get copy of enemies slice for external systems

## Usage Example

```go
// Spawn a new enemy
enemy := gameState.AddEnemy(500*physicsUnit, groundY)

// Remove an enemy
gameState.RemoveEnemy(enemy)

// Clear all enemies (room transition)
gameState.ClearEnemies()

// Get all current enemies
enemies := gameState.GetEnemies()
```

## Physics Consistency

Enemies use the same physics system as the player:
- Same gravity and friction values from `engine.GameConfig`
- Same ground collision detection
- Same physics unit conversions
- Compatible with the existing camera and rendering systems

## Framework for Future Development

This implementation provides a solid framework for:
- Adding different enemy types with different AI behaviors
- Implementing collision detection between player and enemies
- Adding health systems and combat mechanics
- Creating enemy-specific sprites (currently uses player sprites)
- Room-specific enemy spawning and management

## Current Visual Representation

Enemies currently use the same sprites as the player character:
- `idleSprite` when stationary
- `leftSprite` when moving left
- `rightSprite` when moving right

This makes it easy to see the AI behavior in action, and can be easily changed to enemy-specific sprites later.

## Testing

The system can be tested by running the game - you should see:
- 3 enemies patrolling back and forth on the ground
- Enemies moving at half the player's speed
- Debug info showing "Enemies: 3"
- Enemies rendered behind the player character