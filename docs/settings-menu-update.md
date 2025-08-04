# Settings Menu Update

## Overview
The settings menu has been completely redesigned with a tabbed interface to organize different configuration options. The new design includes sections for general settings, keybindings, developer options, and the existing tile viewer.

## New Features

### 1. Tabbed Interface
The settings menu now features four tabs:
- **General**: Basic game settings like window mode and scaling
- **Keybindings**: View and customize keyboard controls
- **Developer**: Debug and development options
- **Tile Viewer**: The existing tile inspection tool

### 2. General Tab
- **Window Mode**: Toggle between fullscreen and windowed mode
- **Tile Scale**: Adjust the scale of tiles (0.5x - 3.0x)
- **Character Scale**: Adjust the scale of character sprites (0.5x - 2.0x)

### 3. Keybindings Tab
Displays all game controls organized by category:
- **Movement**: WASD/Arrow keys for movement and jumping
- **Menu**: Pause, settings, and quit controls
- **Debug**: Debug overlay toggles (F3, G, B, H)
- **Camera**: Zoom controls (+/-)

Features:
- View primary and secondary key bindings
- Edit keybindings by pressing Enter on a binding
- Clear bindings with Delete/Backspace
- ESC cancels key editing

### 4. Developer Tab
Toggle various debug and development options:
- **Show Debug Info**: Display FPS and performance metrics
- **Show Debug Overlay**: Show collision boxes and physics info
- **Show Grid**: Display tile grid overlay
- **Enable Depth of Field**: Toggle depth blur effects
- **Variable Jump Height**: Control jump height by release timing
- **Smooth Camera**: Enable/disable camera smoothing

### 5. Navigation
- **A/D or Left/Right arrows**: Switch between tabs
- **W/S or Up/Down arrows**: Navigate within tabs
- **Enter/Space**: Select or toggle options
- **ESC/Q**: Return to previous menu (main menu or pause)

## Implementation Details

### New Types
- `SettingsTab`: Enum for tab selection
- `KeyBinding`: Stores action name, primary/secondary keys, and description
- `DeveloperOption`: Stores debug option name, description, value pointer, and toggle handler

### Key Methods
- `initializeKeybindings()`: Sets up default keybindings
- `initializeDeveloperOptions()`: Sets up developer options
- `updateGeneralTab()`, `updateKeybindingsTab()`, etc.: Handle input for each tab
- `drawGeneralTab()`, `drawKeybindingsTab()`, etc.: Render each tab's content

### Integration with Engine
- Added `EnableGrid()`, `DisableGrid()`, and `IsGridEnabled()` functions to the engine package
- Developer options directly modify `engine.GameConfig` values
- Some options use custom toggle handlers for special behavior

## Usage
Access the settings menu from:
1. Main menu: Select "Settings" option
2. Pause menu: Press 'S' while paused

The settings state remembers where it was accessed from and returns to the appropriate state when closed.