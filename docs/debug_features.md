# Debug Features Documentation

## Overview
This document describes the debugging features added to help analyze and adjust player size, physics, and game architecture.

## Debug HUD (F3)
The enhanced debug HUD displays comprehensive information about:

### Performance
- FPS (Frames Per Second)
- TPS (Ticks Per Second)

### Player Information
- Physics Position: Raw position in physics units
- Pixel Position: Position converted to screen pixels
- Tile Position: Grid coordinates
- Sprite Size: Actual rendered size in pixels

### Physics Settings
- Physics Unit: Base unit for physics calculations
- Move Speed, Jump Power, Gravity, Friction values

### Rendering Settings
- Tile size and scale factor
- Character scale factor
- Window dimensions

### Camera & Room Info
- Camera position and smoothing
- Room name and dimensions
- Ground level

## Runtime Configuration Hotkeys

### Visual Toggles
- **F3**: Toggle Debug HUD
- **F4**: Toggle Debug Overlay (bounding boxes, velocity vectors)
- **G**: Toggle Grid
- **B**: Toggle Background
- **M**: Toggle Mini-Map

### Scaling Adjustments
- **[ ]**: Decrease/Increase Character Scale
- **- =**: Decrease/Increase Tile Scale
- **Hold Shift**: Fine-tune adjustments (smaller increments)

### Physics Adjustments
- **1**: Adjust Move Speed (Shift to decrease)
- **2**: Adjust Jump Power (Shift to decrease)
- **3**: Adjust Gravity (Shift to decrease)

## Debug Overlay (F4)
When enabled, shows visual debug information:

### Player Debug Overlay
- **Green Box**: Bounding box (yellow when airborne)
- **Red Dot**: Center point
- **Purple Line**: Velocity vector with direction
- **Cyan Arrow**: Facing direction indicator
- **Orange Line**: Ground sensor

### Enemy Debug Overlay
- **Red Box**: Bounding box (orange when airborne)
- **Yellow Dot**: Center point
- **Pink Line**: Velocity vector
- **Red Line**: Facing direction

## Current Configuration (10x10 Room)
The game is currently using `SmallRoomConfig()` with:
- Tile Scale: 2.0 (32px tiles)
- Character Scale: 1.0 (32px sprite)
- 10x10 tile room
- Window: 1920x1080

## Architecture Notes

### Player Implementation
The player uses a physics-based system where:
- Position and velocity are stored in physics units
- Physics unit = TileSize × TileScaleFactor
- Rendering converts physics units to pixels

### Suggested Improvements
1. The player size seems large for the 10x10 grid
2. Consider reducing character scale to 0.5-0.7
3. Physics values may need adjustment for the scale

## Usage Tips
1. Press F3 to see all current values
2. Use [ ] keys to adjust character size in real-time
3. Use F4 to visualize actual bounding boxes
4. Adjust physics with number keys to fine-tune movement feel