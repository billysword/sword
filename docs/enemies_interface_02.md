# Enemy Interface System Implementation (Current)

> **✅ Current Implementation**: This describes the active enemy system architecture using Go interfaces.

This document describes the refactored enemy system using Go interfaces, allowing for flexible enemy types with individual AI implementations.

## Architecture Overview

The enemy system now uses a protocol/interface pattern where:
- All enemies implement the `Enemy` interface
- Common functionality is provided by `BaseEnemy` struct
- Each enemy type implements its own `HandleAI()` method
- No default AI behavior is provided - each type must define its own

## Core Components

### 1. Enemy Interface (`entities/enemy_interface.go`)

**Purpose:** Defines the contract that all enemy types must implement.

**Key Methods:**
- `HandleAI()` - **Must be implemented by each enemy type**
- `Update()` - Physics and AI processing
- `Draw()` / `DrawWithCamera()` - Rendering methods
- `GetPosition()` / `SetPosition()` - Position management
- `GetVelocity()` / `SetVelocity()` - Velocity management
- `IsOnGround()` - Ground state checking
- `Reset()` - State reset functionality
- `GetEnemyType()` - Type identification

### 2. BaseEnemy Struct (`entities/base_enemy.go`)

**Purpose:** Provides common physics and rendering functionality.

**Key Features:**
- **Empty AI stub** - `HandleAI()` does nothing by default
- Complete physics system (gravity, friction, ground collision)
- Basic rendering using player sprites
- Position and velocity management
- Customizable properties (move speed, scale, friction)

**Protected Methods for Subclasses:**
- `SetMoveSpeed()` / `GetMoveSpeed()`
- `SetScale()`

### 3. Concrete Enemy Types

#### SlimeEnemy (`entities/slime_enemy.go`)
- **AI Behavior:** Patrol within range of spawn point
- **Movement:** Back and forth with random direction changes
- **Speed:** Half player speed
- **Characteristics:** Standard size, boundary-respecting

#### WandererEnemy (`entities/wanderer_enemy.go`)
- **AI Behavior:** Random wandering without boundaries
- **Movement:** Frequent random direction changes with pausing
- **Speed:** 75% of player speed (faster than slime)
- **Characteristics:** Smaller size (80% scale), highly unpredictable

## Interface Usage Pattern

### Creating New Enemy Types

```go
type MyCustomEnemy struct {
    *BaseEnemy // Embed for common functionality
    
    // Custom AI properties
    customProperty int
}

func NewMyCustomEnemy(x, y int) *MyCustomEnemy {
    enemy := &MyCustomEnemy{
        BaseEnemy: NewBaseEnemy(x, y),
        customProperty: someValue,
    }
    
    // Customize properties
    enemy.SetMoveSpeed(customSpeed)
    enemy.SetScale(customScaleX, customScaleY)
    
    return enemy
}

// REQUIRED: Implement your own AI logic
func (e *MyCustomEnemy) HandleAI() {
    // Your custom AI logic here
    // This method MUST be implemented
}

// REQUIRED: Implement type identification
func (e *MyCustomEnemy) GetEnemyType() string {
    return "mycustom"
}

// Optional: Override other methods as needed
func (e *MyCustomEnemy) Reset(x, y int) {
    e.BaseEnemy.Reset(x, y)
    // Reset custom AI state
}
```

### Game State Integration

The game state now works with the `Enemy` interface:

```go
type InGameState struct {
    enemies []entities.Enemy // Interface slice
    // ... other fields
}

// Adding enemies
func (ig *InGameState) AddEnemy(enemy entities.Enemy) {
    ig.enemies = append(ig.enemies, enemy)
}

// Convenience methods for specific types
func (ig *InGameState) AddSlimeEnemy(x, y int) *entities.SlimeEnemy {
    slime := entities.NewSlimeEnemy(x, y)
    ig.enemies = append(ig.enemies, slime)
    return slime
}
```

## AI Implementation Requirements

### Mandatory AI Implementation

Each enemy type **MUST** implement `HandleAI()` with their specific behavior:

```go
func (e *EnemyType) HandleAI() {
    // Your AI logic here
    // Set e.vx for horizontal movement
    // Set e.vy for vertical movement (if needed)
    // Access base properties via e.BaseEnemy
}
```

### No Default Behavior

- The base `HandleAI()` method is intentionally empty
- This forces each enemy type to explicitly define its behavior
- Prevents accidental use of generic/default AI
- Ensures each enemy type has purpose-built logic

## Current Enemy Types Comparison

| Enemy Type | AI Pattern | Speed | Size | Behavior |
|------------|------------|-------|------|----------|
| SlimeEnemy | Patrol | 50% player speed | Normal | Stays near spawn, boundaries |
| WandererEnemy | Random | 75% player speed | 80% scale | No boundaries, frequent changes |

## Testing Setup

The game now spawns a mix of enemy types:
- Position 300: SlimeEnemy (patrol)
- Position 600: WandererEnemy (random)
- Position 900: SlimeEnemy (patrol) 
- Position 1200: WandererEnemy (random)

## Benefits of Interface System

### For Developers:
1. **Type Safety:** Interface ensures all methods are implemented
2. **Flexibility:** Easy to add new enemy types with unique behaviors
3. **Code Reuse:** Common functionality shared via BaseEnemy
4. **Maintainability:** Clear separation of concerns

### For Game Design:
1. **Behavioral Diversity:** Each enemy can have completely unique AI
2. **Easy Extension:** Add new enemy types without modifying existing code
3. **Performance:** No unnecessary default behavior execution
4. **Debugging:** Type identification for logging and analysis

## Example Enemy Behaviors to Implement

### Suggested Enemy Types:
- **JumperEnemy:** Jumps periodically while moving
- **ChargerEnemy:** Rushes toward player when in range
- **GuardEnemy:** Stationary until disturbed
- **FlierEnemy:** Ignores gravity, flies in patterns
- **ShooterEnemy:** Stationary, shoots projectiles

### Implementation Template:
```go
func (e *NewEnemyType) HandleAI() {
    physicsUnit := engine.GetPhysicsUnit()
    
    // 1. Update AI state/timers
    e.aiTimer--
    
    // 2. Make movement decisions
    if /* condition */ {
        e.vx = direction * e.GetMoveSpeed() * physicsUnit
    }
    
    // 3. Handle special behaviors
    if /* special condition */ {
        // Jump, shoot, charge, etc.
    }
}
```

This interface-based system provides a solid foundation for creating diverse enemy types while maintaining code organization and type safety.