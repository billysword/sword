# Simple Room Layout System

## Overview

A greatly simplified room layout system that uses a dedicated `room_layouts/` package for predefined layouts and simple console debugging.

## Key Benefits

- **Simple Import**: Just `import "sword/room_layouts"` and use predefined layouts
- **Easy Creation**: Add new `.go` files to `room_layouts/` folder
- **Console Debugging**: `PrintRoomLayout()` outputs copy-paste ready code
- **No File Generation**: No complex log files or auto-generated content
- **Clean Git**: No gitignore complexity or generated file management

## Usage

### Using Predefined Layouts

```go
import "sword/room_layouts"
import "sword/world"

// Create a room
room := world.NewBaseRoom("my_room", 10, 8)

// Apply a predefined layout
world.ApplyLayout(room, room_layouts.ExamplePlatform)

// Optional: Print to console for debugging
world.PrintRoomLayout("my_room", room.GetTileMap())
```

### Available Layouts

- `room_layouts.ExamplePlatform` - Basic platformer with platforms and ground
- `room_layouts.EmptyRoom` - Simple empty room with just ground
- `room_layouts.TowerClimb` - Vertical climbing room with staggered platforms

### Creating Custom Layouts

Just add a new file to `room_layouts/`:

```go
// room_layouts/my_layout.go
package room_layouts

var MyLayout = [][]int{
    {-1, -1, -1, -1, -1}, // Sky
    {-1, 0x5, 0x6, 0x7, -1}, // Platform
    {0x1, 0x2, 0x3, 0x1, 0x2}, // Ground
}
```

### Inline Layouts

For simple cases, create layouts inline:

```go
customLayout := [][]int{
    {-1, -1, -1, -1, -1},
    {0x1, 0x2, 0x3, 0x1, 0x2},
}
world.ApplyLayout(room, customLayout)
```

## Debugging

### Console Output

Use `PrintRoomLayout()` to get copy-paste ready output:

```go
world.PrintRoomLayout("room_name", room.GetTileMap())
```

Output:
```
=== ROOM LAYOUT: room_name ===
Dimensions: 10x8 tiles

Hex format (copy-paste ready):
var layout = [][]int{
    {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
    {0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1},
}
=== END ROOM LAYOUT ===
```

### BaseRoom Helper

```go
room.PrintRoomDebug() // Same as PrintRoomLayout(room.GetZoneID(), room.GetTileMap())
```

## File Structure

```
room_layouts/
├── README.md           # Usage instructions
├── empty_room.go       # Basic empty room
├── example_platform.go # Platform example
├── tower_climb.go      # Vertical climbing
└── my_custom.go        # Your custom layouts
```

## Hex Values

- `-1` = Empty space (air)
- `0x1-0xFF` = Tile indices (0-255)
- Use hex notation for clarity: `0x5` instead of `5`

## Comparison with Old System

| Feature | Old Complex System | New Simple System |
|---------|-------------------|-------------------|
| **Usage** | Complex file generation | Simple import + apply |
| **Debugging** | Log files + standalone .go files | Console output |
| **Adding Layouts** | Copy from generated files | Add .go file to folder |
| **Git Management** | Complex gitignore patterns | No special handling needed |
| **File Count** | Many auto-generated files | Just the layouts you create |

## Migration from Complex System

1. **Remove old calls**: Replace `LogRoomDebug()` + `GenerateHexLayoutFile()` with `PrintRoomLayout()`
2. **Create layouts**: Move any useful generated layouts to `room_layouts/` folder
3. **Update imports**: Add `import "sword/room_layouts"`
4. **Apply layouts**: Use `world.ApplyLayout(room, room_layouts.LayoutName)`

## Examples

See `examples/simple_room_usage.go` for a complete working example.

## Related Files

- `room_layouts/` - All predefined layouts
- `world/simple_debug.go` - Console debugging functions
- `world/layout_helper.go` - ApplyLayout helper function
- `examples/simple_room_usage.go` - Working example

This system is much simpler and more maintainable than the previous complex file generation approach.