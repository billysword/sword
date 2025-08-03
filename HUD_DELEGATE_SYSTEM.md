# HUD Delegate System Architecture

## Overview

Implemented a clean delegate system for HUD and UI components that follows the same Update/Draw paradigm as game states. This removes UI logic from the critical path and makes the system much more extensible and maintainable.

## 🏗️ Architecture

### Core Interface (`engine/hud_system.go`)
```go
type HUDComponent interface {
    Update() error
    Draw(screen interface{}) error
    IsVisible() bool
    SetVisible(visible bool)
    GetName() string
}
```

### HUD Manager (`engine/hud_system.go`)
```go
type HUDManager struct {
    components map[string]HUDComponent
    enabled    bool
}

// Key Methods:
// - AddComponent(component HUDComponent)
// - Update() error  // Updates all visible components
// - Draw(screen interface{}) error  // Draws all visible components
// - ToggleComponent(name string)
```

## 🧩 Components Implemented

### 1. Debug HUD (`engine/debug_hud.go`)
- **Purpose**: Displays debug information (replaces old inline debug rendering)
- **Features**: Room info, player position, camera position, TPS, toggle states
- **Auto-updating**: Performance info, window size, debug toggle states

### 2. Mini-Map (`world/minimap.go`)
- **Purpose**: World map visualization (when rendering is implemented)
- **Features**: Implements HUDComponent interface, provides map data via GetMapData()
- **Future-ready**: Empty Update/Draw methods ready for implementation

## 🔄 Integration with Game State

### InGameState Changes
```go
type InGameState struct {
    // Existing fields...
    
    // New HUD system (replaces old inline HUD code)
    hudManager *engine.HUDManager
}
```

### Update Pipeline
```go
// In InGameState.Update():

// 1. Update game state data for HUD components
if debugHUD := ig.hudManager.GetComponent("debug_hud"); debugHUD != nil {
    if dh, ok := debugHUD.(*engine.DebugHUD); ok {
        dh.UpdateRoomInfo(roomInfo)
        dh.UpdatePlayerPos(playerPos)
        dh.UpdateCameraPos(cameraPos)
    }
}

// 2. Update all HUD components
if err := ig.hudManager.Update(); err != nil {
    return err
}
```

### Draw Pipeline
```go
// In InGameState.Draw():

// Replace old inline HUD drawing with:
if err := ig.hudManager.Draw(screen); err != nil {
    engine.LogError("HUD draw error", err)
}
```

### Input Handling
```go
// Clean component toggling:
if inpututil.IsKeyJustPressed(ebiten.KeyM) {
    ig.hudManager.ToggleComponent("minimap")
}
```

## ✅ Benefits Achieved

### 1. **Separation of Concerns**
- Game state logic separated from UI rendering
- Each HUD component manages its own state and rendering
- InGameState focuses purely on game logic

### 2. **Extensibility**
- Easy to add new HUD components (health bars, inventory, etc.)
- Components can be toggled independently
- Clean interfaces for different component types

### 3. **Same Update/Draw Paradigm**
- HUD components follow identical pattern to game states
- Consistent error handling across all components
- Predictable lifecycle management

### 4. **Performance**
- Components only update/draw when visible
- Batch processing of all HUD elements
- Easy to disable entire HUD system if needed

### 5. **Maintainability**
- HUD logic no longer clutters game state code
- Each component is self-contained and testable
- Clear component registration and management

## 🎯 Usage Examples

### Adding a New HUD Component
```go
// 1. Implement HUDComponent interface
type HealthBar struct {
    visible bool
    name    string
    health  int
}

func (hb *HealthBar) GetName() string { return hb.name }
func (hb *HealthBar) IsVisible() bool { return hb.visible }
func (hb *HealthBar) SetVisible(visible bool) { hb.visible = visible }
func (hb *HealthBar) Update() error { /* update logic */ return nil }
func (hb *HealthBar) Draw(screen interface{}) error { /* draw logic */ return nil }

// 2. Register with HUD manager
healthBar := &HealthBar{visible: true, name: "health_bar"}
hudManager.AddComponent(healthBar)

// 3. Toggle from anywhere
hudManager.ToggleComponent("health_bar")
```

### Getting Component Data
```go
// Access specific component
if miniMap := hudManager.GetComponent("minimap"); miniMap != nil {
    if mm, ok := miniMap.(*world.MiniMapRenderer); ok {
        mapData := mm.GetMapData()
        // Use map data for custom rendering
    }
}
```

## 📁 File Structure

```
engine/
├── hud_system.go     # Core HUD interface and manager
└── debug_hud.go      # Debug information component

world/
└── minimap.go        # Mini-map component (implements HUDComponent)

states/
└── ingame_state.go   # Simplified - uses HUD manager instead of inline code
```

## 🚀 Ready for Extension

The delegate system is ready for:
- **Health/Mana Bars**: Visual player stats
- **Inventory HUD**: Item management interface  
- **Chat System**: Multiplayer communication
- **Notification System**: Game events and alerts
- **Menu Overlays**: Pause menus, settings, etc.

## ✅ Success Metrics

1. **Clean Architecture**: Game logic separated from UI concerns
2. **Extensible Design**: Easy to add new HUD components
3. **Consistent Patterns**: Same Update/Draw paradigm throughout
4. **Performance Optimized**: Components only process when visible
5. **Maintainable Code**: Self-contained, testable components

The HUD delegate system provides a **solid foundation for all UI components** while keeping the critical game state logic clean and focused.