# Metroidvania-Style Game Changes

## Overview
This document summarizes the changes made to transform the game into a proper metroidvania-style platformer with appropriate window size, camera system, and room design.

## Key Changes

### 1. Window Resolution (main.go)
- Changed from 960x540 to **640x360** pixels
- This is a standard metroidvania resolution with 16:9 aspect ratio
- Provides good pixel art visibility while allowing for larger explorable rooms

### 2. Camera System (gamestate/camera.go - NEW FILE)
- Implemented a smooth-following camera with:
  - Dead zones (25% width, 16% height) for player movement without camera motion
  - Smooth interpolation (10% smoothing factor)
  - World boundary constraints with margins for HUD elements
  - Parallax background scrolling (30% speed)
  - Helper methods for coordinate conversion and visibility checks

### 3. Room Size (gamestate/simple_room.go)
- Increased room size from 60x34 tiles to **80x40 tiles**
- This creates a 2560x1280 pixel world (much larger than the 640x360 viewport)
- Added multiple platforms at various heights for exploration
- Added wall structures and pillars for more interesting level design

### 4. Rendering Updates
- **InGameState** (gamestate/ingame_state.go):
  - Added camera instance and updates
  - Modified Draw method to apply camera transformations
  - Added camera position to debug info
  
- **Player** (gamestate/player.go):
  - Added `DrawWithCamera` method for camera-relative rendering
  
- **Room Interface** (gamestate/room.go):
  - Added `DrawWithCamera` method to interface
  - Added `DrawTilesWithCamera` method for tile rendering with offset
  
- **SimpleRoom** (gamestate/simple_room.go):
  - Implemented `DrawWithCamera` with parallax background
  - Updated tile rendering to use camera offset

### 5. Debug Features (gamestate/state.go)
- Added `DrawGridWithCamera` function that moves the debug grid with the camera
- Grid lines properly align with tiles regardless of camera position

## Room Layout Features
The new room design includes:
- Main ground level spanning the entire width
- 5 floating platforms at different heights
- 3 wall structures for vertical navigation
- Proper tile edges and corners for visual polish

## Camera Behavior
- Follows the player smoothly with a slight delay
- Stays within world bounds with margins:
  - 32 pixels on left/right/top
  - 48 pixels on bottom (for potential HUD)
- Background scrolls at 30% speed for parallax effect
- Debug grid moves with the camera for easier tile alignment

## Benefits
1. **Exploration**: Large rooms encourage exploration beyond the visible screen
2. **Polish**: Smooth camera and parallax create professional feel
3. **HUD Space**: Margins ensure space for UI elements
4. **Standard Size**: 640x360 is widely used in modern pixel art games
5. **Scalability**: Camera system makes it easy to add larger rooms or connected areas