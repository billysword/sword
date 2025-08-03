# CLAUDE.md - Demo/Placeholder Code Cleanup Locations

## Overview

This document systematically lists all locations in the codebase that contain demo/placeholder code, hardcoded values, or incomplete implementations that need to be cleaned up or properly implemented.

## 🎯 Priority Categories

### 🔴 **HIGH PRIORITY - Broken/Incomplete Functionality**

#### 1. State Factory - Unimplemented Methods
**Location:** `engine/state_factory.go:201-214`
```go
func (sf *StateFactory) createStartState(config StateConfig) (State, error) {
    return nil, fmt.Errorf("StartState creation not yet implemented in factory")
}
func (sf *StateFactory) createInGameState(config StateConfig) (State, error) {
    return nil, fmt.Errorf("InGameState creation not yet implemented in factory")
}
func (sf *StateFactory) createPauseState(config StateConfig) (State, error) {
    return nil, fmt.Errorf("PauseState creation not yet implemented in factory")
}
```
**Issue:** Core state creation methods return errors instead of creating states
**Impact:** State factory system is non-functional

#### 2. Base Enemy - Empty AI Implementation
**Location:** `entities/base_enemy.go:61-64`
```go
func (be *BaseEnemy) HandleAI() {
    // Empty stub - concrete enemy types must implement their own AI logic
}
```
**Issue:** Base class has empty AI method that should be abstract/interface-based
**Impact:** Breaks polymorphism expectations for enemy system

### 🟡 **MEDIUM PRIORITY - Demo Content & Hardcoded Values**

#### 3. Hardcoded Window Resolution
**Location:** `engine/config.go:75-76`
```go
WindowWidth:  1920,
WindowHeight: 1080,
```
**Issue:** Fixed 1920x1080 resolution doesn't adapt to different displays
**Impact:** Poor experience on non-1080p displays

#### 4. Demo Room Layouts
**Location:** `room_layouts/` directory
- `example_platform.go` - Demo platformer layout
- `tower_climb.go` - Demo tower layout  
- `empty_room.go` - Demo empty layout

**Issue:** These are clearly demo/example layouts with hardcoded hex values
**Impact:** Takes up space, not production-ready content

#### 5. Examples Directory - Demo Code
**Location:** `examples/simple_room_usage.go`
```go
fmt.Println("=== Simple Room Layout Usage ===\n")
fmt.Println("1. Using predefined layout (ExamplePlatform):")
// ... more demo output
fmt.Println("That's it! Much simpler than the complex file generation system.")
```
**Issue:** Pure demo/tutorial code with console output
**Impact:** Not needed for production game

#### 6. Legacy Config Constants
**Location:** `engine/config.go:191-196`
```go
const (
    TILE_SIZE         = 16   // Deprecated: use GameConfig.TileSize
    TILE_SCALE_FACTOR = 1.0  // Deprecated: use GameConfig.TileScaleFactor
    CHAR_SCALE_FACTOR = 0.4  // Deprecated: use GameConfig.CharScaleFactor
    PHYSICS_UNIT      = 16   // Deprecated: use GetPhysicsUnit()
)
```
**Issue:** Deprecated constants maintained for "backward compatibility"
**Impact:** Code confusion, technical debt

### 🟢 **LOW PRIORITY - Documentation & Polish**

#### 7. Demo References in Documentation
**Locations:**
- `docs/README.md:1` - "Forest Tilemap Platformer **Demo**"
- `docs/README.md:17` - "Room Switch**: R key (**demo feature**)"
- `docs/README.md:36` - `go build -o **demo**`
- `docs/README.md:50` - "Forest-themed **demo** room implementation"

**Issue:** Documentation treats this as a demo rather than a real game
**Impact:** Confusing for users/developers

#### 8. Demo-style Debug Output
**Location:** Multiple files contain excessive debug logging and console output
- `examples/simple_room_usage.go` - Multiple `fmt.Println` calls
- Various debug logging scattered throughout

**Issue:** Debug code should be configurable/removable for production
**Impact:** Performance and log noise

#### 9. Hardcoded Magic Numbers in Configurations
**Locations:**
- `CLAUDE.md:15-19` - Room sizes, ground levels in architectural proposals
- `engine/config.go` - Various magic numbers like `80x60 tiles`, `25 ground level`

**Issue:** Configuration values are hardcoded without semantic meaning
**Impact:** Difficult to understand and modify

## 🔧 **Recommended Actions**

### Phase 1: Fix Broken Functionality
1. **Implement State Factory Methods**
   - Replace error returns with actual state creation
   - Integrate with existing state constructors in `states/` directory

2. **Refactor Enemy AI System**
   - Make `HandleAI()` abstract or use proper interface pattern
   - Remove empty stub implementation

### Phase 2: Clean Demo Content
3. **Remove Demo Room Layouts**
   - Delete `room_layouts/example_platform.go`
   - Delete `room_layouts/tower_climb.go`
   - Keep only `empty_room.go` as a template if needed

4. **Remove Examples Directory**
   - Delete `examples/simple_room_usage.go`
   - Move any useful patterns to actual game code or tests

5. **Remove Legacy Constants**
   - Delete deprecated constants in `engine/config.go`
   - Update any remaining usage to use `GameConfig`

### Phase 3: Production Polish
6. **Update Documentation Language**
   - Remove "demo" references
   - Rebrand as production game documentation
   - Update build targets from "demo" to "game"

7. **Make Debug Output Configurable**
   - Add debug level configuration
   - Remove hardcoded `fmt.Println` statements
   - Use proper logging with levels

8. **Create Semantic Configuration**
   - Replace magic numbers with named constants
   - Add configuration validation
   - Document what each value means

## 🗂️ **File Cleanup Summary**

### Files to Delete:
- `examples/simple_room_usage.go` - Pure demo code
- `room_layouts/example_platform.go` - Demo layout
- `room_layouts/tower_climb.go` - Demo layout

### Files to Significantly Modify:
- `engine/state_factory.go` - Implement missing methods
- `entities/base_enemy.go` - Fix AI interface
- `engine/config.go` - Remove legacy constants
- `docs/README.md` - Remove demo language
- `CLAUDE.md` - Replace with production roadmap

### Files to Clean (Remove Debug/Demo Code):
- All files with `fmt.Println` debug output
- Documentation with "demo" references
- Hardcoded configuration values

## 🎯 **Next Steps**

1. **Start with Phase 1** - Fix the broken state factory and enemy AI
2. **Remove demo content** - Clean up examples and demo layouts  
3. **Polish documentation** - Update language to reflect production status
4. **Create configuration system** - Replace hardcoded values with semantic configs

This cleanup will transform the project from a demo/prototype into a production-ready game foundation.