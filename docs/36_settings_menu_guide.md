# Settings Menu - Room Tile Viewer

A new settings menu has been added to the game that allows you to view the currently loaded room tiles with their hex values. This is useful for debugging, modding, and understanding the room layout structure.

## Features

### Tile Display
- **Visual Grid**: Each tile is displayed as a colored rectangle with its hex value and coordinates
- **Color Coding**: Different tile types are color-coded for easy identification:
  - **Empty tiles (-1)**: Dark gray
  - **Ground tiles (0x01-0x10)**: Brown
  - **Platform tiles (0x11-0x20)**: Green  
  - **Other tiles**: Blue
- **Hex Values**: Each tile shows its index in hexadecimal format (e.g., 0x01, 0x1F, -1)
- **Coordinates**: Each tile displays its (x,y) position in the room grid

### Navigation
- **Scrolling**: Use W/S or Arrow Keys to scroll through large tile maps
- **Fast Scroll**: Use Page Up/Page Down for faster navigation
- **Room Info**: Display shows current room name and dimensions

## How to Access

### From Main Menu
1. Start the game
2. In the main menu, use W/S or Arrow Keys to navigate to "Settings"
3. Press ENTER or SPACE to open the settings menu
4. Press ESC or Q to return to the main menu

### From In-Game (Pause Menu)
1. While playing, press P or ESC to pause the game
2. Press S to open settings with the current room's tile data
3. Press ESC or Q to return to the pause menu
4. Press P or ESC again to resume the game

## Controls

### In Settings Menu
- **W/S or ↑/↓**: Scroll up/down through the tile grid
- **Page Up/Page Down**: Fast scroll through large tile maps
- **ESC or Q**: Exit settings (returns to previous menu)

## Technical Details

### File Structure
- `states/settings_state.go`: Main settings state implementation
- `states/start_state.go`: Modified to include Settings option
- `states/pause_state.go`: Modified to allow settings access during gameplay
- `states/ingame_state.go`: Added getter method for current room access

### Integration
The settings menu integrates seamlessly with the existing state management system:
- Maintains proper navigation flow between states
- Preserves game state when accessed from pause menu
- Displays real-time room data from the currently loaded room

### Data Sources
- Accesses tile data through the `world.Room` interface
- Uses `GetTileMap()` method to retrieve tile layout information
- Displays tile indices as stored in the room's 2D tile array

## Use Cases

1. **Level Design**: Inspect room layouts and tile placement
2. **Debugging**: Verify tile indices and coordinates during development
3. **Modding**: Understand existing room structures for modification
4. **Learning**: Study how different room types are constructed

The settings menu provides a powerful tool for inspecting and understanding the game's tile-based room system, making it easier to work with and modify room layouts.