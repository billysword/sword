# CLAUDE.md - Architectural Proposals for Game Development

## Project Overview

This is a Go-based 2D platformer/metroidvania game built using Ebitengine. The project uses a modular architecture with state management, configurable settings, and a flexible room layout system.

## Current Architecture Status

### Window & Display
- **Current Resolution**: 1920x1080 (16:9 aspect ratio)
- **Tile Size**: 16x16 pixels (standard)
- **Tile Scale Factor**: 1.0 (no scaling in default config)
- **Character Scale Factor**: 0.7 (default), 0.4 (legacy)

### Room System
- **Default Room Size**: 80x60 tiles (1280x960 pixels)
- **Zoomed-in Room Size**: 40x30 tiles (640x480 pixels)
- **Ground Level**: 45 tiles from top (default)

## Architectural Proposals

### 1. Aspect Ratio & Resolution System

**Current Issues:**
- Fixed 16:9 aspect ratio doesn't adapt to different displays
- No support for ultrawide (21:9) or legacy (4:3) monitors
- Room layouts don't adjust to different aspect ratios

**Proposed Solution:**

```go
// engine/display_config.go
type AspectRatio struct {
    Name   string
    Width  int
    Height int
    // Calculated room dimensions for this aspect
    DefaultRoomTilesX int
    DefaultRoomTilesY int
}

var SupportedAspectRatios = []AspectRatio{
    {"16:9 HD", 1920, 1080, 80, 45},
    {"16:9 FHD", 2560, 1440, 106, 60},
    {"21:9 UW", 2560, 1080, 106, 45},
    {"16:10", 1920, 1200, 80, 50},
    {"4:3", 1600, 1200, 66, 50},
}

// Auto-detect and adapt room layouts
func AdaptRoomToAspectRatio(room *Room, aspectRatio AspectRatio) {
    // Dynamically adjust room bounds and camera limits
}
```

### 2. Dynamic Tile Size System

**Current Issues:**
- Fixed 16x16 tile size limits visual fidelity
- No support for high-DPI displays
- Character sprites don't scale proportionally with tiles

**Proposed Solution:**

```go
// engine/tile_system.go
type TileConfig struct {
    BaseSize     int     // Original tile size (16)
    RenderSize   int     // Actual rendered size
    ScaleFactor  float64 // Dynamic scale based on resolution
    PixelPerfect bool    // Maintain pixel boundaries
}

// Calculate optimal tile size based on display
func CalculateOptimalTileSize(windowWidth, windowHeight int) TileConfig {
    // Target ~80-120 tiles horizontally for good visibility
    targetTilesX := 100
    optimalSize := windowWidth / targetTilesX
    
    // Snap to nearest multiple of base size
    scaleFactor := float64(optimalSize) / 16.0
    if pixelPerfect {
        scaleFactor = math.Round(scaleFactor)
    }
    
    return TileConfig{
        BaseSize:     16,
        RenderSize:   int(16 * scaleFactor),
        ScaleFactor:  scaleFactor,
        PixelPerfect: true,
    }
}
```

### 3. Character Size & Proportion System

**Current Issues:**
- Character scale is hardcoded relative to tiles
- No visual hierarchy between player, enemies, and bosses
- Sprites look too small in zoomed-out view

**Proposed Solution:**

```go
// entities/size_config.go
type EntitySizeConfig struct {
    BaseWidth    int     // Base sprite width
    BaseHeight   int     // Base sprite height
    ScaleRelativeToTile float64 // Scale relative to tile size
    
    // Size categories
    SizeClass    EntitySize
    
    // Collision box relative to sprite size
    CollisionScale float64
}

type EntitySize int
const (
    SizeTiny EntitySize = iota  // 0.5x tile (items, projectiles)
    SizeSmall                   // 1x tile (small enemies)
    SizeMedium                  // 1.5x tile (player, regular enemies)
    SizeLarge                   // 2x tile (mini-bosses)
    SizeHuge                    // 3x+ tiles (bosses)
)

// Standardized entity sizes
var EntitySizes = map[EntitySize]EntitySizeConfig{
    SizeMedium: {
        BaseWidth:  32,
        BaseHeight: 32,
        ScaleRelativeToTile: 1.5,
        CollisionScale: 0.8,
    },
    // ... other sizes
}
```

### 4. Advanced Room Layout System

**Current Issues:**
- Simple 2D array limits complex room designs
- No support for room templates or procedural generation
- Difficult to create interconnected metroidvania-style maps

**Proposed Solution:**

```go
// world/room_architecture.go
type RoomTemplate struct {
    ID          string
    Name        string
    BaseLayout  [][]int
    
    // Room properties
    Type        RoomType
    Biome       BiomeType
    Difficulty  int
    
    // Connection points
    Exits       []ExitPoint
    
    // Dynamic elements
    SpawnPoints []SpawnPoint
    Triggers    []Trigger
}

type RoomType int
const (
    RoomNormal RoomType = iota
    RoomBoss
    RoomSave
    RoomSecret
    RoomTransition
    RoomHub
)

type ExitPoint struct {
    Side      Direction // North, South, East, West
    Position  int       // Position along that side (in tiles)
    Width     int       // Width of exit (in tiles)
    LeadsTo   string    // Room ID or "any" for procedural
    Required  bool      // Must connect to another room
}

// Room generation system
type RoomGenerator struct {
    templates map[string]RoomTemplate
    biomes    map[BiomeType]BiomeConfig
}

func (rg *RoomGenerator) GenerateRoom(template string, seed int64) *Room {
    // Create room from template with variations
    // Add procedural elements based on seed
    // Ensure exit compatibility
}

// Metroidvania map structure
type WorldMap struct {
    Rooms       map[string]*Room
    Connections map[string][]string // Room connections
    Regions     map[string]Region   // Named regions (e.g., "Forest", "Caves")
    
    // Player progression tracking
    VisitedRooms   map[string]bool
    UnlockedExits  map[string]bool
}
```

### 5. Responsive UI Scaling

**Current Issues:**
- UI elements don't scale with resolution
- HUD takes up fixed pixel space
- No adaptive layout for different screen sizes

**Proposed Solution:**

```go
// ui/responsive_ui.go
type UIScaler struct {
    BaseResolution  Resolution // Design resolution (e.g., 1920x1080)
    CurrentResolution Resolution
    
    // Scaling modes
    Mode UIScaleMode
    
    // Safe areas for different aspect ratios
    SafeArea Rectangle
}

type UIScaleMode int
const (
    UIScaleFit UIScaleMode = iota   // Fit to screen, may have bars
    UIScaleStretch                   // Stretch to fill (may distort)
    UIScalePixelPerfect              // Integer scaling only
    UIScaleAdaptive                  // Smart scaling with layout changes
)

// Calculate UI element positions/sizes
func (ui *UIScaler) ScaleElement(baseRect Rectangle) Rectangle {
    switch ui.Mode {
    case UIScaleAdaptive:
        // Reposition elements for optimal layout
        return ui.adaptiveScale(baseRect)
    default:
        // Simple scaling
        return ui.uniformScale(baseRect)
    }
}
```

### 6. Configuration Presets System

**Proposed Enhancement to Current Config:**

```go
// engine/config_presets.go
type ConfigPreset struct {
    Name        string
    Description string
    Config      Config
    
    // Preset categories
    Category    PresetCategory
    Recommended ForDisplayType
}

type PresetCategory int
const (
    PresetBalanced PresetCategory = iota
    PresetPerformance
    PresetQuality
    PresetAccessibility
)

var ConfigPresets = map[string]ConfigPreset{
    "pixel_perfect": {
        Name: "Pixel Perfect",
        Description: "Crisp pixels with integer scaling",
        Config: Config{
            TileScaleFactor: 2.0,
            CharScaleFactor: 2.0,
            // ... pixel-perfect settings
        },
    },
    "smooth_hd": {
        Name: "Smooth HD",
        Description: "High resolution with smooth scaling",
        // ...
    },
    "performance": {
        Name: "Performance",
        Description: "Optimized for lower-end systems",
        // ...
    },
}
```

## Implementation Roadmap

### Phase 1: Foundation (Current Sprint)
1. ✅ Clean up duplicate documentation
2. 🔄 Implement dynamic tile sizing system
3. 🔄 Create character size standardization
4. 🔄 Add aspect ratio detection

### Phase 2: Room System Enhancement
1. Implement room template system
2. Add procedural generation support
3. Create room connection validator
4. Build metroidvania map structure

### Phase 3: Visual Polish
1. Implement responsive UI system
2. Add configuration presets
3. Create smooth camera transitions
4. Implement parallax depth system

### Phase 4: Advanced Features
1. Add room-specific visual themes
2. Implement dynamic lighting system
3. Create advanced particle effects
4. Add post-processing pipeline

## Technical Considerations

### Performance
- Use object pooling for tiles and entities
- Implement frustum culling for off-screen objects
- Cache scaled sprites at common resolutions
- Use spatial partitioning for collision detection

### Compatibility
- Test on various aspect ratios and resolutions
- Ensure pixel-perfect mode works correctly
- Support both windowed and fullscreen modes
- Handle display DPI scaling properly

### Modularity
- Keep systems decoupled and configurable
- Use interfaces for extensibility
- Maintain backwards compatibility
- Document all configuration options

## Next Steps

1. **Immediate**: Update `engine/config.go` with new tile sizing options
2. **Short-term**: Implement basic aspect ratio support
3. **Medium-term**: Create room template system
4. **Long-term**: Build full metroidvania map editor

This architecture provides a solid foundation for a scalable, professional-quality 2D platformer that can adapt to different displays and player preferences while maintaining visual consistency and performance.