# Debug HUD Improvements

## Overview
Fixed flashing issues in the debug HUD and added comprehensive player information display.

## Changes Made

### 1. Fixed Flashing Issues
- **Problem**: Toggle states (Background, Grid, Depth of Field) were being recalculated every frame in the Draw method, causing inconsistent display
- **Solution**: Moved toggle state calculations to the Update method and stored them in the DebugInfo struct
- **Result**: Consistent display without flashing

### 2. Added Player Information
- **Velocity Display**: Shows both physics units and pixels/frame for X and Y velocity
- **Status Display**: Shows whether player is on ground or in air
- **Facing Direction**: Shows which direction the player is facing (Left/Right)

### 3. Improved Organization
- Reorganized DebugInfo struct with clear sections:
  - Player info (position, velocity, status)
  - World info (room, camera)
  - System info (window size, performance)
  - Toggle states (stored to prevent flashing)
  - Custom info (extensible)

### 4. Better Performance Tracking
- Performance info (FPS/TPS) now stored in struct and updated consistently
- Window size stored separately to avoid recalculation

## Debug HUD Sections

1. **PERFORMANCE**: FPS and TPS
2. **PLAYER**: Position, velocity, and status
3. **PHYSICS**: Physics unit, movement settings
4. **RENDERING**: Tile and character scaling info
5. **CAMERA**: Camera position and smoothing
6. **ROOM**: Current room info and ground level
7. **TOGGLES**: Background, Grid, Depth of Field states
8. **CUSTOM**: Any additional debug info
9. **HOTKEYS**: Control reference (right side)

## Usage
- Press F3 to toggle the debug HUD
- All information updates in real-time
- No more flashing or inconsistent display