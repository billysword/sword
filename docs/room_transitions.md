# Room Transition System

The room transition system provides a flexible way to handle movement between different rooms/areas in the game. It consists of exit triggers, entrance spawn points, and automatic room loading/unloading.

## Components

### 1. Exit Triggers (`ExitTrigger`)

Exit triggers are areas in a room that initiate transitions to other rooms when the player enters them.

```go
type ExitTrigger struct {
    ID               string   // Unique identifier for this exit
    X, Y             int      // Position in tile coordinates
    Width, Height    int      // Size of the trigger area (default 1x1 for single tile)
    TargetRoomID     string   // ID of the room to transition to
    TargetEntranceID string   // ID of the entrance in the target room to spawn at
    TransitionType   string   // Type of transition (e.g., "fade", "slide", "instant")
    RequiredItems    []string // Items/keys required to use this exit
    IsLocked         bool     // Whether this exit is currently locked
    Message          string   // Optional message to display when triggered/locked
}
```

### 2. Entrances (`Entrance`)

Entrances define where players spawn when entering a room from an exit trigger.

```go
type Entrance struct {
    ID           string  // Unique identifier for this entrance
    X, Y         int     // Position in tile coordinates where player spawns
    Direction    string  // Direction player should face on spawn ("left", "right", "up", "down")
    SpawnOffsetX int     // Fine-tuning spawn position (pixels from tile center)
    SpawnOffsetY int     // Fine-tuning spawn position (pixels from tile center)
}
```

### 3. Room Transition System (`RoomTransitionSystem`)

The main system that manages room transitions, loading/unloading, and coordinates the transition process.

## How to Use

### Setting up Room Connections

1. **Register Exit Triggers** for each room:
```go
transitionSystem.RegisterExitTrigger("main", ExitTrigger{
    ID:               "main_to_forest",
    X:                9,  // Right edge of room
    Y:                8,  // Near the bottom
    Width:            1,
    Height:           1,
    TargetRoomID:     "forest_entrance",
    TargetEntranceID: "from_main",
    TransitionType:   "instant",
})
```

2. **Register Entrances** for each room:
```go
transitionSystem.RegisterEntrance("forest_entrance", Entrance{
    ID:        "from_main",
    X:         9,  // Right side
    Y:         8,
    Direction: "left",
})
```

### Creating New Rooms

1. **Add room to the factory** in `world/room_factory.go`:
```go
func (rf *RoomFactory) CreateRoom(roomID string) Room {
    switch roomID {
    case "your_new_room":
        return NewYourNewRoom()
    // ... other cases
    }
}
```

2. **Set up connections** in `SetupDefaultRoomConnections()`:
```go
// Add exit triggers and entrances for your new room
```

## Current Room Connections

The default setup includes these rooms and connections:

- **main** → **forest_entrance** (exit at right edge, x=9, y=8)
- **forest_entrance** → **main** (exit at left edge, x=0, y=8)
- **forest_entrance** → **underground_cave** (exit at bottom center, x=4-5, y=9)
- **underground_cave** → **forest_entrance** (exit at top center, x=4-5, y=0)

## Debug Visualization

When debug overlay is enabled (F4 key), the system displays:

- **Green rectangles**: Accessible exit triggers
- **Red rectangles**: Locked exit triggers
- **Blue markers**: Entrance spawn points

## Features

### Automatic Features
- **Room Loading**: Rooms are loaded on-demand when first accessed
- **Camera Recentering**: Camera automatically centers on player in new room
- **Player Positioning**: Player spawns at correct entrance position
- **World Bounds Update**: Camera world bounds update to new room size
- **Enemy Management**: Enemies are cleared when changing rooms (room-specific)
- **World Map Integration**: New rooms are automatically discovered and marked

### Transition Types
Currently supports "instant" transitions. Future expansion could include:
- "fade": Fade to black transition
- "slide": Sliding camera transition
- "door": Door opening animation

### Locked Exits
Exits can be locked and require specific items:
```go
ExitTrigger{
    // ... other fields
    IsLocked:      true,
    RequiredItems: []string{"key_forest", "torch"},
    Message:       "You need a forest key to enter here.",
}
```

## Integration

The room transition system is integrated into `InGameState` and automatically checks for transitions each frame after player movement and collision handling. When a player touches an exit trigger:

1. Transition data is created
2. Target room is loaded (if not already cached)
3. Player is positioned at target entrance
4. Camera bounds are updated
5. World map is updated
6. Room callbacks (OnExit/OnEnter) are called

## Performance

- **Room Caching**: Loaded rooms are cached to avoid reloading
- **On-Demand Loading**: Rooms are only loaded when first accessed
- **Memory Management**: Rooms can be unloaded if memory management is needed

## Example Usage

```go
// In your game initialization:
transitionSystem := world.NewRoomTransitionSystem(world.GetRoomFactory())
roomFactory := world.NewRoomFactory()
roomFactory.SetupDefaultRoomConnections(transitionSystem)

// The system automatically handles transitions during gameplay
// Player just needs to walk into an exit trigger area
```