# Room Layouts

Simple room layouts that can be imported and used in your rooms.

## Usage

```go
import "sword/room_layouts"
import "sword/world"

// Create a room
room := world.NewBaseRoom("my_room", 10, 8)

// Apply a predefined layout
world.ApplyLayout(room, room_layouts.ExamplePlatform)

// Optional: Print layout to console for debugging
world.PrintRoomLayout("my_room", room.GetTileMap())
```

## Available Layouts

- **ExamplePlatform** - Basic platformer room with platforms and ground
- **EmptyRoom** - Simple empty room with just ground
- **TowerClimb** - Vertical climbing room with staggered platforms

## Creating New Layouts

Just add a new `.go` file in this folder:

```go
package room_layouts

// MyLayout is a custom room layout
var MyLayout = [][]int{
    {-1, -1, -1, -1, -1}, // Sky
    {0x1, 0x2, 0x3, 0x1, 0x2}, // Ground
}
```

## Hex Values

- `-1` = Empty space
- `0x1-0xFF` = Tile indices (0-255)
- Use hex notation (`0x5`) for clarity

That's it! Much simpler than complex file generation systems.