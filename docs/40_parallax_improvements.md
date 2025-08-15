# Parallax and Depth Improvements

## 🐛 Bug Found: Camera Viewport Resizing Issue

### Problem Description
The main bug identified was in the camera system's handling of window resizing. While the game enables `ebiten.WindowResizingModeEnabled`, the camera viewport dimensions were set once during initialization and never updated when users resized the window.

### Impact
- **Incorrect camera bounds calculation**: Camera thought the viewport was smaller/larger than actual window
- **Misaligned dead zone positioning**: Dead zones were calculated based on the original window size
- **HUD element positioning errors**: Elements drawn using `ebiten.WindowSize()` didn't match camera expectations
- **Parallax calculation misalignment**: Background scrolling became incorrect relative to the viewport

### Solution Implemented
Added `updateCameraViewport()` method in `states/ingame_state.go` that:
1. Checks for window size changes each frame
2. Creates a new camera with updated viewport dimensions
3. Restores world bounds and player following behavior
4. Ensures proper synchronization between window size and camera viewport

**Files Modified:**
- `states/ingame_state.go`: Added viewport update logic

## 🎨 Enhanced Parallax/Depth System

### Previous Implementation Limitations
The original system had only a single background layer with fixed parallax speed (`ParallaxFactor: 0.5`), providing minimal depth perception.

### New Multi-Layer Parallax System

#### Key Features
1. **Multiple Background/Foreground Layers**: Support for unlimited parallax layers with different properties
2. **Bi-directional Parallax**: Background layers (speed < 1.0) and foreground layers (speed > 1.0)
3. **Depth-Based Effects**: Automatic transparency, scaling, and color adjustments based on depth
4. **Depth of Field Simulation**: Blur effects, desaturation, and contrast reduction for distant layers
5. **Interactive Controls**: Runtime toggling and configuration cycling for testing
6. **Unified Rendering**: No fallback mechanisms - all rendering through parallax system

#### Implementation Details

**New Configuration Structure:**
```go
type ParallaxLayer struct {
    Speed       float64 // Scroll speed relative to camera (0-1)
    Depth       float64 // Depth for visual effects (0-1)
    Alpha       float64 // Transparency (0-1)
    Scale       float64 // Scale factor for the layer
    OffsetX     float64 // Static horizontal offset
    OffsetY     float64 // Static vertical offset
    Repeatable  bool    // Whether the layer should tile/repeat
}
```

**Files Created/Modified:**
- `engine/parallax_renderer.go`: New multi-layer parallax renderer
- `engine/config.go`: Extended configuration for parallax layers and depth effects
- `states/ingame_state.go`: Added debug controls and viewport management

#### Depth of Field Effects
The system simulates depth of field through:
- **Color desaturation**: Distant layers become less saturated
- **Contrast reduction**: Background layers have reduced contrast
- **Subtle blue tinting**: Atmospheric perspective effect
- **Position jitter**: Very subtle blur simulation
- **Alpha blending**: Depth-based transparency

### Debug Controls Added

| Key | Function | Description |
|-----|----------|-------------|
| `D` | Toggle Depth of Field | Enable/disable depth blur and color effects |
| `L` | Cycle Parallax Layers | Switch between different layer configurations |

### Demo Configurations

1. **Minimal (1 layer)**: Single background layer for performance testing
2. **Balanced (3 layers)**: Background, midground, and foreground with speed > 1.0 for foreground
3. **Dramatic (6 layers)**: Maximum depth effect with multiple background and foreground layers

**Note**: Layers with speed > 1.0 move faster than the camera, creating foreground parallax effects.

## 🎯 Visual Improvements Achieved

### Before
- Single static background with fixed 50% parallax speed
- No depth perception or layering
- Static visual experience

### After
- Multiple moving layers creating depth illusion in both directions
- Background parallax (slower than camera) and foreground parallax (faster than camera)
- Automatic depth-based visual effects (transparency, scale, color)
- Interactive depth of field simulation
- Dynamic layer configuration for different visual styles
- Proper window resizing support
- Unified rendering system with no fallback mechanisms

## 🔧 Technical Benefits

1. **Modular Design**: ParallaxRenderer can be reused across different room types
2. **Performance Conscious**: Efficient rendering with proper culling and effect application
3. **Configurable**: Easy to adjust depth effects and layer properties
4. **Unified Rendering**: All background/foreground rendering goes through the parallax system
5. **Debug Friendly**: Runtime controls for testing different configurations

## 🚀 Future Enhancement Opportunities

1. **Asset Loading**: Load different images for each parallax layer
2. **Parallax Tiles**: Support for tiled/repeating backgrounds
3. **Particle Integration**: Depth-sorted particle effects
4. **Dynamic Weather**: Layer-based weather effects (rain, fog, etc.)
5. **Lighting Effects**: Depth-aware lighting and shadows
6. **Motion Blur**: Enhanced blur effects for fast camera movement

## 📊 Performance Considerations

The enhanced system is designed to be performant:
- Layers are only drawn if visible and have valid images
- Color matrix operations are only applied when needed
- Depth effects use efficient ColorM transformations
- Fallback to simple rendering if complex effects fail

## 🎮 User Experience

Players now experience:
- **Enhanced Immersion**: Multiple depth layers create realistic environment feel
- **Visual Feedback**: Depth effects provide spatial awareness
- **Customizable Experience**: Debug controls allow real-time adjustment
- **Smooth Resizing**: No visual artifacts when changing window size
- **Responsive Camera**: Proper dead zones and following behavior regardless of window size

This implementation transforms the game from a flat 2D experience into a layered pseudo-3D environment with proper depth perception and atmospheric effects.