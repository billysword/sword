#!/usr/bin/env python3
"""
World Map System Demo

This script demonstrates how to use the world map system to track
discovered rooms and their spatial relationships.
"""

from world_map import WorldMap, Direction
from map_visualizer import MapVisualizer

def demo_basic_mapping():
    """Demonstrate basic room mapping functionality"""
    print("=== Basic World Mapping Demo ===\n")
    
    # Create a new world map
    world_map = WorldMap()
    
    # Player starts in the tavern
    tavern = world_map.add_room(
        room_id="tavern_001",
        name="The Prancing Pony Tavern",
        description="A cozy tavern with a warm fireplace and friendly atmosphere.",
        width=8.0,
        height=6.0,
        is_starting_room=True
    )
    
    # Player discovers there are exits to the north and east
    world_map.discover_exit("tavern_001", Direction.NORTH)
    world_map.discover_exit("tavern_001", Direction.EAST)
    
    print("Player starts in the tavern and discovers exits...")
    visualizer = MapVisualizer(world_map)
    print(visualizer.render_ascii_map(40, 15, "tavern_001"))
    print()
    
    # Player goes north and discovers a market square
    market = world_map.add_room(
        room_id="market_001",
        name="Market Square",
        description="A bustling marketplace filled with merchants and shoppers.",
        width=12.0,
        height=10.0
    )
    
    # Connect the rooms
    world_map.connect_rooms("tavern_001", Direction.NORTH, "market_001")
    
    # Player discovers more exits from the market
    world_map.discover_exit("market_001", Direction.WEST)
    world_map.discover_exit("market_001", Direction.EAST)
    
    print("Player goes north to the market square...")
    print(visualizer.render_ascii_map(50, 20, "market_001"))
    print()
    
    # Player goes east from tavern to discover a small alley
    alley = world_map.add_room(
        room_id="alley_001",
        name="Narrow Alley",
        description="A cramped alley between buildings, dimly lit.",
        width=3.0,
        height=8.0
    )
    
    world_map.connect_rooms("tavern_001", Direction.EAST, "alley_001")
    world_map.discover_exit("alley_001", Direction.NORTH)
    
    print("Player explores east to find a narrow alley...")
    print(visualizer.render_ascii_map(50, 20, "alley_001"))
    print()
    
    # Show map summary
    print(visualizer.print_map_summary())
    print("\n" + "="*50 + "\n")

def demo_complex_mapping():
    """Demonstrate more complex room relationships and conflict resolution"""
    print("=== Complex World Mapping Demo ===\n")
    
    world_map = WorldMap()
    
    # Create a more complex dungeon layout
    rooms_data = [
        ("entrance", "Dungeon Entrance", "A stone archway leading into darkness.", 6, 6, True),
        ("corridor_n", "North Corridor", "A long stone corridor.", 4, 12),
        ("chamber_nw", "Crystal Chamber", "A chamber filled with glowing crystals.", 8, 8),
        ("chamber_ne", "Armory", "An old armory with rusty weapons.", 6, 10),
        ("corridor_e", "East Corridor", "A winding corridor.", 10, 4),
        ("treasure", "Treasure Room", "A room filled with gold and gems.", 12, 8),
        ("secret", "Secret Passage", "A hidden passage behind a bookshelf.", 3, 15),
    ]
    
    # Add all rooms
    for room_id, name, desc, width, height, *is_start in rooms_data:
        is_starting = len(is_start) > 0 and is_start[0]
        world_map.add_room(room_id, name, desc, width, height, is_starting)
    
    # Create complex connections
    connections = [
        ("entrance", Direction.NORTH, "corridor_n"),
        ("entrance", Direction.EAST, "corridor_e"),
        ("corridor_n", Direction.NORTHWEST, "chamber_nw"),
        ("corridor_n", Direction.NORTHEAST, "chamber_ne"),
        ("corridor_e", Direction.NORTH, "treasure"),
        ("chamber_nw", Direction.UP, "secret"),  # Secret passage above
        ("chamber_ne", Direction.EAST, "corridor_e"),
    ]
    
    for from_room, direction, to_room in connections:
        world_map.connect_rooms(from_room, direction, to_room)
    
    # Add some unexplored exits
    world_map.discover_exit("treasure", Direction.SOUTH)
    world_map.discover_exit("secret", Direction.WEST)
    world_map.discover_exit("entrance", Direction.DOWN)
    
    visualizer = MapVisualizer(world_map)
    
    print("Complex dungeon layout discovered:")
    print(visualizer.render_ascii_map(70, 25))
    print()
    
    # Show room details
    print("Detailed room information:")
    for room in visualizer.get_room_list('discovery'):
        print(f"\n{visualizer.print_room_details(room.id)}")
    
    print("\n" + "="*50 + "\n")

def demo_persistence():
    """Demonstrate saving and loading world maps"""
    print("=== Map Persistence Demo ===\n")
    
    # Create and populate a world map
    world_map = WorldMap()
    
    # Simple 3-room layout
    world_map.add_room("start", "Starting Room", "Where it all begins.", 5, 5, True)
    world_map.add_room("middle", "Middle Room", "A crossroads.", 8, 6)
    world_map.add_room("end", "End Room", "The final destination.", 6, 7)
    
    world_map.connect_rooms("start", Direction.NORTH, "middle")
    world_map.connect_rooms("middle", Direction.EAST, "end")
    world_map.discover_exit("end", Direction.SOUTH)
    
    print("Original map:")
    visualizer = MapVisualizer(world_map)
    print(visualizer.render_ascii_map(40, 15))
    
    # Save to file
    filename = "demo_world_map.json"
    world_map.save_to_file(filename)
    print(f"\nMap saved to {filename}")
    
    # Load from file
    loaded_map = WorldMap.load_from_file(filename)
    loaded_visualizer = MapVisualizer(loaded_map)
    
    print("\nLoaded map:")
    print(loaded_visualizer.render_ascii_map(40, 15))
    
    print("\nMap data matches:", world_map.to_dict() == loaded_map.to_dict())
    print("\n" + "="*50 + "\n")

def demo_graphical_data():
    """Demonstrate generating data for graphical rendering"""
    print("=== Graphical Data Export Demo ===\n")
    
    world_map = WorldMap()
    
    # Create a circular room layout
    center = world_map.add_room("center", "Central Plaza", "The heart of the city.", 10, 10, True)
    
    # Add surrounding rooms
    surrounding_rooms = [
        ("north_gate", "North Gate", "The northern entrance.", 8, 4, Direction.NORTH),
        ("south_market", "South Market", "Trading district.", 12, 6, Direction.SOUTH),
        ("east_temple", "Temple of Light", "A sacred temple.", 6, 8, Direction.EAST),
        ("west_inn", "Wayfarer's Inn", "A place to rest.", 7, 5, Direction.WEST),
    ]
    
    for room_id, name, desc, width, height, direction in surrounding_rooms:
        world_map.add_room(room_id, name, desc, width, height)
        world_map.connect_rooms("center", direction, room_id)
    
    visualizer = MapVisualizer(world_map)
    
    print("City layout:")
    print(visualizer.render_ascii_map(50, 20))
    print()
    
    # Export graphical data
    graphical_data = visualizer.get_graphical_data()
    
    print("Graphical data structure:")
    print(f"Total rooms: {graphical_data['total_rooms']}")
    print(f"Map bounds: {graphical_data['bounds']}")
    print(f"Starting room: {graphical_data['starting_room_id']}")
    
    print("\nRoom data:")
    for room_data in graphical_data['rooms']:
        print(f"  {room_data['name']}: pos=({room_data['x']:.1f}, {room_data['y']:.1f}), "
              f"size={room_data['width']}x{room_data['height']}")
    
    print(f"\nConnections: {len(graphical_data['connections'])} found")
    for conn in graphical_data['connections']:
        from_room = next(r for r in graphical_data['rooms'] if r['id'] == conn['from_room'])
        to_room = next(r for r in graphical_data['rooms'] if r['id'] == conn['to_room'])
        print(f"  {from_room['name']} -> {to_room['name']} ({conn['direction']})")
    
    print("\n" + "="*50 + "\n")

def interactive_demo():
    """Interactive demo where user can explore mapping commands"""
    print("=== Interactive World Mapping Demo ===\n")
    print("Commands:")
    print("  add <id> <name> <width> <height> - Add a new room")
    print("  connect <from> <direction> <to> - Connect two rooms")
    print("  exit <room> <direction> - Discover an exit")
    print("  map [room_id] - Show the map (optionally highlight a room)")
    print("  room <id> - Show room details")
    print("  summary - Show map summary")
    print("  save <filename> - Save map to file")
    print("  load <filename> - Load map from file")
    print("  quit - Exit demo")
    print()
    
    world_map = WorldMap()
    
    # Add a starting room
    world_map.add_room("start", "Starting Chamber", "You begin your journey here.", 5, 5, True)
    print("Created starting room 'start'")
    
    while True:
        try:
            command = input("\n> ").strip().split()
            if not command:
                continue
                
            cmd = command[0].lower()
            
            if cmd == "quit":
                break
            elif cmd == "add" and len(command) >= 5:
                room_id, name, width, height = command[1], command[2], float(command[3]), float(command[4])
                world_map.add_room(room_id, name, f"A room called {name}.", width, height)
                print(f"Added room '{room_id}'")
            elif cmd == "connect" and len(command) >= 4:
                from_room, direction_str, to_room = command[1], command[2], command[3]
                try:
                    direction = Direction(direction_str.lower())
                    world_map.connect_rooms(from_room, direction, to_room)
                    print(f"Connected {from_room} -> {to_room} ({direction_str})")
                except ValueError:
                    print(f"Invalid direction: {direction_str}")
            elif cmd == "exit" and len(command) >= 3:
                room_id, direction_str = command[1], command[2]
                try:
                    direction = Direction(direction_str.lower())
                    world_map.discover_exit(room_id, direction)
                    print(f"Discovered exit {direction_str} from {room_id}")
                except ValueError:
                    print(f"Invalid direction: {direction_str}")
            elif cmd == "map":
                visualizer = MapVisualizer(world_map)
                current_room = command[1] if len(command) > 1 else None
                print(visualizer.render_ascii_map(60, 20, current_room))
            elif cmd == "room" and len(command) >= 2:
                visualizer = MapVisualizer(world_map)
                print(visualizer.print_room_details(command[1]))
            elif cmd == "summary":
                visualizer = MapVisualizer(world_map)
                print(visualizer.print_map_summary())
            elif cmd == "save" and len(command) >= 2:
                world_map.save_to_file(command[1])
                print(f"Saved map to {command[1]}")
            elif cmd == "load" and len(command) >= 2:
                world_map = WorldMap.load_from_file(command[1])
                print(f"Loaded map from {command[1]}")
            else:
                print("Invalid command or missing parameters")
                
        except Exception as e:
            print(f"Error: {e}")
    
    print("Interactive demo ended.")

def main():
    """Run all demos"""
    demo_basic_mapping()
    demo_complex_mapping()
    demo_persistence()
    demo_graphical_data()
    
    # Ask if user wants to try interactive demo
    try:
        response = input("Would you like to try the interactive demo? (y/n): ").strip().lower()
        if response.startswith('y'):
            interactive_demo()
    except KeyboardInterrupt:
        print("\nDemo interrupted.")
    except EOFError:
        print("\nDemo ended.")

if __name__ == "__main__":
    main()