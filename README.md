# World Map System

A dynamic world mapping system for games where players discover rooms as they explore. The system creates an approximation of the world based on visited rooms and their relative positions and sizes.

## Features

### Core Functionality
- **Dynamic Room Discovery**: Players can only map rooms they've visited
- **Spatial Relationships**: Rooms are positioned relative to each other based on exit connections
- **Variable Room Sizes**: Each room has its own width and height dimensions
- **Exit Management**: Track discovered exits and connections between rooms
- **Conflict Resolution**: Automatically resolves spatial conflicts when room placements overlap

### Visualization
- **ASCII Map Rendering**: Generate text-based maps for console display
- **Room Details**: Detailed information about individual rooms
- **Map Summary Statistics**: Overview of discovered areas and connections
- **Graphical Data Export**: Export data structure for external rendering systems

### Persistence
- **Save/Load Maps**: Serialize discovered maps to JSON files
- **Cross-Session Continuity**: Maintain exploration progress between game sessions

## Core Components

### WorldMap Class
The main class that manages the collection of discovered rooms and their spatial relationships.

```python
from world_map import WorldMap, Direction

# Create a new world map
world_map = WorldMap()

# Add rooms as they're discovered
tavern = world_map.add_room(
    room_id="tavern_001",
    name="The Prancing Pony Tavern",
    description="A cozy tavern with a warm fireplace",
    width=8.0,
    height=6.0,
    is_starting_room=True
)

# Connect rooms when player moves between them
market = world_map.add_room("market_001", "Market Square", "Bustling marketplace", 12.0, 10.0)
world_map.connect_rooms("tavern_001", Direction.NORTH, "market_001")

# Discover exits before exploring them
world_map.discover_exit("market_001", Direction.EAST)
```

### Room Class
Represents a discovered room with its properties and exits.

```python
# Room properties
room.id          # Unique identifier
room.name        # Display name
room.description # Descriptive text
room.width       # Room width in world units
room.height      # Room height in world units
room.x           # World X coordinate (center)
room.y           # World Y coordinate (center)
room.exits       # Dictionary of exits by direction
```

### Direction Enum
Defines all possible exit directions including diagonals and vertical movement.

```python
Direction.NORTH, Direction.SOUTH, Direction.EAST, Direction.WEST
Direction.NORTHEAST, Direction.NORTHWEST, Direction.SOUTHEAST, Direction.SOUTHWEST
Direction.UP, Direction.DOWN
```

## Usage Examples

### Basic Room Mapping

```python
from world_map import WorldMap, Direction
from map_visualizer import MapVisualizer

# Initialize the world map
world_map = WorldMap()

# Player starts in a room
start = world_map.add_room("start", "Starting Chamber", "Where it begins", 5, 5, True)

# Player discovers exits
world_map.discover_exit("start", Direction.NORTH)
world_map.discover_exit("start", Direction.EAST)

# Player explores north
corridor = world_map.add_room("corridor", "Long Corridor", "A stone hallway", 3, 12)
world_map.connect_rooms("start", Direction.NORTH, "corridor")

# Visualize the map
visualizer = MapVisualizer(world_map)
print(visualizer.render_ascii_map(60, 20))
```

### Complex Dungeon Layout

```python
# Create multiple interconnected rooms
entrance = world_map.add_room("entrance", "Dungeon Entrance", "Stone archway", 6, 6, True)
armory = world_map.add_room("armory", "Old Armory", "Rusty weapons", 8, 10)
treasure = world_map.add_room("treasure", "Treasure Room", "Gold and gems", 12, 8)

# Create complex connections
world_map.connect_rooms("entrance", Direction.NORTH, "armory")
world_map.connect_rooms("armory", Direction.EAST, "treasure")

# Add unexplored exits
world_map.discover_exit("treasure", Direction.SOUTH)
world_map.discover_exit("entrance", Direction.DOWN)
```

### Map Persistence

```python
# Save discovered map
world_map.save_to_file("my_world_map.json")

# Load saved map
loaded_map = WorldMap.load_from_file("my_world_map.json")

# Maps can also be serialized to dictionaries
map_data = world_map.to_dict()
restored_map = WorldMap.from_dict(map_data)
```

### Visualization Options

```python
visualizer = MapVisualizer(world_map)

# ASCII map with current room highlighted
ascii_map = visualizer.render_ascii_map(80, 24, current_room_id="entrance")

# Detailed room information
room_details = visualizer.print_room_details("treasure")

# Map statistics
summary = visualizer.print_map_summary()

# Export for graphical rendering
graphical_data = visualizer.get_graphical_data()
```

## Advanced Features

### Spatial Queries

```python
# Find rooms within a circular area
nearby_rooms = world_map.get_rooms_in_area(center_x=0, center_y=0, radius=10)

# Get map boundaries
min_x, min_y, max_x, max_y = world_map.get_map_bounds()

# Sort rooms by various criteria
sorted_rooms = visualizer.get_room_list(sort_by='discovery')  # 'name', 'size', 'distance_from_start'
```

### Room Exit Positioning

```python
# Get exact coordinates of room exits
exit_x, exit_y = room.get_exit_position(Direction.NORTH)

# Exits are positioned at room edges:
# - Cardinal directions: at the center of the respective edge
# - Diagonal directions: at 70% distance from center to corner
# - UP/DOWN: at room center (same position, different level)
```

### Conflict Resolution

The system automatically handles spatial conflicts when rooms would overlap:

1. **Detection**: Checks for overlapping room boundaries
2. **Resolution**: Moves conflicting rooms apart while maintaining connections
3. **Iteration**: Repeats until all conflicts are resolved (max 10 iterations)

## File Structure

```
world_map.py        # Core world map functionality
map_visualizer.py   # Visualization and rendering
world_map_demo.py   # Comprehensive demonstration
test_world_map.py   # Test suite
README.md           # This documentation
```

## Running the System

### Run Tests
```bash
python3 test_world_map.py
```

### Run Demonstrations
```bash
python3 world_map_demo.py
```

### Interactive Demo
The demo includes an interactive mode where you can experiment with mapping commands:
- `add <id> <name> <width> <height>` - Add a new room
- `connect <from> <direction> <to>` - Connect two rooms
- `exit <room> <direction>` - Discover an exit
- `map [room_id]` - Show the map
- `room <id>` - Show room details
- `summary` - Show map summary
- `save <filename>` - Save map to file
- `load <filename>` - Load map from file

## Design Principles

### Non-Grid Based
Unlike traditional grid-based dungeon maps, this system uses actual spatial coordinates. Rooms can be any size and positioned anywhere, creating more realistic and varied layouts.

### Discovery-Driven
Players can only map what they've seen. The system supports:
- Discovered but unexplored exits (marked with '?')
- Partial knowledge of room connections
- Progressive revelation of the world

### Approximate Mapping
The system creates an approximation of the actual world layout. Room positioning is based on the order of discovery and connection information, which may not perfectly match the "real" world layout.

### Extensible Design
The modular design allows for easy extension:
- Add new room properties
- Implement custom visualization methods
- Create specialized spatial algorithms
- Integrate with different game engines

## Use Cases

- **Text-based Adventure Games**: Perfect for MUDs and interactive fiction
- **Roguelike Games**: Dynamic mapping of procedurally generated content
- **RPG Quest Tracking**: Keep track of discovered locations
- **World Building Tools**: Design and visualize interconnected areas
- **Educational Tools**: Teach spatial reasoning and graph structures

The system provides a solid foundation for any game requiring dynamic world mapping where player exploration drives map discovery.