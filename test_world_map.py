#!/usr/bin/env python3
"""
Test script for the World Map System

This script tests all major functionality of the world map system
to ensure it works correctly.
"""

import os
import tempfile
from world_map import WorldMap, Direction, Room, RoomExit
from map_visualizer import MapVisualizer

def test_basic_room_creation():
    """Test basic room creation and properties"""
    print("Testing basic room creation...")
    
    world_map = WorldMap()
    
    # Test creating a starting room
    room = world_map.add_room("test_001", "Test Room", "A test room.", 10.0, 8.0, True)
    
    assert room.id == "test_001"
    assert room.name == "Test Room"
    assert room.width == 10.0
    assert room.height == 8.0
    assert room.is_starting_room == True
    assert room.x == 0.0  # Starting room should be at origin
    assert room.y == 0.0
    assert world_map.starting_room_id == "test_001"
    
    # Test creating a second room
    room2 = world_map.add_room("test_002", "Second Room", "Another test room.", 6.0, 6.0)
    
    assert room2.id == "test_002"
    assert room2.is_starting_room == False
    assert len(world_map.rooms) == 2
    
    print("✓ Basic room creation works correctly")

def test_room_connections():
    """Test connecting rooms with exits"""
    print("Testing room connections...")
    
    world_map = WorldMap()
    
    # Create two rooms
    room1 = world_map.add_room("room1", "Room 1", "First room.", 5.0, 5.0, True)
    room2 = world_map.add_room("room2", "Room 2", "Second room.", 5.0, 5.0)
    
    # Connect them
    world_map.connect_rooms("room1", Direction.NORTH, "room2")
    
    # Check that exits were created correctly
    assert Direction.NORTH in room1.exits
    assert Direction.SOUTH in room2.exits
    
    assert room1.exits[Direction.NORTH].target_room_id == "room2"
    assert room2.exits[Direction.SOUTH].target_room_id == "room1"
    
    assert room1.exits[Direction.NORTH].is_discovered == True
    assert room2.exits[Direction.SOUTH].is_discovered == True
    
    # Check positioning
    assert room2.x == 0.0  # Should be directly north
    assert room2.y == 5.0  # Distance should be sum of half-heights
    
    print("✓ Room connections work correctly")

def test_exit_discovery():
    """Test discovering exits without connections"""
    print("Testing exit discovery...")
    
    world_map = WorldMap()
    room = world_map.add_room("room1", "Room 1", "Test room.", 5.0, 5.0, True)
    
    # Discover an exit
    world_map.discover_exit("room1", Direction.EAST)
    
    assert Direction.EAST in room.exits
    assert room.exits[Direction.EAST].target_room_id is None
    assert room.exits[Direction.EAST].is_discovered == False
    
    print("✓ Exit discovery works correctly")

def test_spatial_positioning():
    """Test that rooms are positioned correctly relative to each other"""
    print("Testing spatial positioning...")
    
    world_map = WorldMap()
    
    # Create a center room
    center = world_map.add_room("center", "Center", "Center room.", 4.0, 4.0, True)
    
    # Add rooms in all directions
    north = world_map.add_room("north", "North", "North room.", 6.0, 4.0)
    south = world_map.add_room("south", "South", "South room.", 6.0, 4.0)
    east = world_map.add_room("east", "East", "East room.", 4.0, 6.0)
    west = world_map.add_room("west", "West", "West room.", 4.0, 6.0)
    
    # Connect them
    world_map.connect_rooms("center", Direction.NORTH, "north")
    world_map.connect_rooms("center", Direction.SOUTH, "south")
    world_map.connect_rooms("center", Direction.EAST, "east")
    world_map.connect_rooms("center", Direction.WEST, "west")
    
    # Check positions
    assert center.x == 0.0 and center.y == 0.0
    assert north.x == 0.0 and north.y > 0  # North should be positive Y
    assert south.x == 0.0 and south.y < 0  # South should be negative Y
    assert east.x > 0 and east.y == 0.0    # East should be positive X
    assert west.x < 0 and west.y == 0.0    # West should be negative X
    
    print("✓ Spatial positioning works correctly")

def test_conflict_resolution():
    """Test that overlapping rooms are moved apart"""
    print("Testing conflict resolution...")
    
    world_map = WorldMap()
    
    # Create rooms that would naturally overlap
    room1 = world_map.add_room("room1", "Room 1", "Large room.", 20.0, 20.0, True)
    room2 = world_map.add_room("room2", "Room 2", "Another large room.", 20.0, 20.0)
    room3 = world_map.add_room("room3", "Room 3", "Third large room.", 20.0, 20.0)
    
    # Connect them in ways that would cause overlap
    world_map.connect_rooms("room1", Direction.EAST, "room2")
    world_map.connect_rooms("room1", Direction.NORTH, "room3")
    
    # Check that rooms don't overlap
    assert not world_map._rooms_overlap(room1, room2)
    assert not world_map._rooms_overlap(room1, room3)
    assert not world_map._rooms_overlap(room2, room3)
    
    print("✓ Conflict resolution works correctly")

def test_map_bounds():
    """Test map bounds calculation"""
    print("Testing map bounds calculation...")
    
    world_map = WorldMap()
    
    # Test empty map
    bounds = world_map.get_map_bounds()
    assert bounds == (0, 0, 0, 0)
    
    # Add some rooms
    world_map.add_room("center", "Center", "Center room.", 4.0, 4.0, True)
    world_map.add_room("far_east", "Far East", "Far room.", 2.0, 2.0)
    world_map.add_room("far_west", "Far West", "Far room.", 2.0, 2.0)
    
    # Position them manually for testing
    world_map.rooms["far_east"].x = 10.0
    world_map.rooms["far_east"].y = 0.0
    world_map.rooms["far_west"].x = -10.0
    world_map.rooms["far_west"].y = 0.0
    
    min_x, min_y, max_x, max_y = world_map.get_map_bounds()
    
    # Should include all room boundaries
    assert min_x <= -11.0  # far_west left edge
    assert max_x >= 11.0   # far_east right edge
    assert min_y <= -2.0   # room bottom edges
    assert max_y >= 2.0    # room top edges
    
    print("✓ Map bounds calculation works correctly")

def test_rooms_in_area():
    """Test spatial area queries"""
    print("Testing spatial area queries...")
    
    world_map = WorldMap()
    
    # Create rooms at known positions
    center = world_map.add_room("center", "Center", "Center room.", 2.0, 2.0, True)
    near = world_map.add_room("near", "Near", "Near room.", 2.0, 2.0)
    far = world_map.add_room("far", "Far", "Far room.", 2.0, 2.0)
    
    # Position them manually
    near.x, near.y = 3.0, 0.0   # 3 units away
    far.x, far.y = 20.0, 0.0    # 20 units away
    
    # Query for rooms within 5 units of center
    nearby_rooms = world_map.get_rooms_in_area(0.0, 0.0, 5.0)
    room_ids = [room.id for room in nearby_rooms]
    
    assert "center" in room_ids
    assert "near" in room_ids
    assert "far" not in room_ids
    
    print("✓ Spatial area queries work correctly")

def test_persistence():
    """Test saving and loading world maps"""
    print("Testing map persistence...")
    
    # Create a world map with some complexity
    world_map = WorldMap()
    
    room1 = world_map.add_room("room1", "First Room", "Starting point.", 5.0, 5.0, True)
    room2 = world_map.add_room("room2", "Second Room", "Connected room.", 7.0, 4.0)
    room3 = world_map.add_room("room3", "Third Room", "Another room.", 3.0, 8.0)
    
    world_map.connect_rooms("room1", Direction.NORTH, "room2")
    world_map.connect_rooms("room2", Direction.EAST, "room3")
    world_map.discover_exit("room3", Direction.SOUTH)
    
    # Save to temporary file
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        temp_filename = f.name
    
    try:
        world_map.save_to_file(temp_filename)
        
        # Load from file
        loaded_map = WorldMap.load_from_file(temp_filename)
        
        # Verify all data was preserved
        assert len(loaded_map.rooms) == 3
        assert loaded_map.starting_room_id == "room1"
        assert loaded_map.discovery_count == 3
        
        # Check room data
        loaded_room1 = loaded_map.get_room("room1")
        assert loaded_room1.name == "First Room"
        assert loaded_room1.width == 5.0
        assert loaded_room1.is_starting_room == True
        
        # Check connections
        assert Direction.NORTH in loaded_room1.exits
        assert loaded_room1.exits[Direction.NORTH].target_room_id == "room2"
        
        # Check unexplored exit
        loaded_room3 = loaded_map.get_room("room3")
        assert Direction.SOUTH in loaded_room3.exits
        assert loaded_room3.exits[Direction.SOUTH].target_room_id is None
        
    finally:
        # Clean up
        if os.path.exists(temp_filename):
            os.unlink(temp_filename)
    
    print("✓ Map persistence works correctly")

def test_visualization():
    """Test map visualization functionality"""
    print("Testing map visualization...")
    
    world_map = WorldMap()
    
    # Create a simple layout
    center = world_map.add_room("center", "Central Plaza", "The heart of town.", 6.0, 6.0, True)
    north = world_map.add_room("north", "North District", "Northern area.", 4.0, 4.0)
    
    world_map.connect_rooms("center", Direction.NORTH, "north")
    world_map.discover_exit("north", Direction.WEST)
    
    visualizer = MapVisualizer(world_map)
    
    # Test ASCII rendering
    ascii_map = visualizer.render_ascii_map(40, 15)
    assert isinstance(ascii_map, str)
    assert len(ascii_map) > 0
    assert 'S' in ascii_map  # Should contain starting room marker
    
    # Test room details
    details = visualizer.print_room_details("center")
    assert "Central Plaza" in details
    assert "6.0 x 6.0" in details
    
    # Test map summary
    summary = visualizer.print_map_summary()
    assert "Total Rooms Discovered: 2" in summary
    assert "Starting Room: Central Plaza" in summary
    
    # Test graphical data export
    graphical_data = visualizer.get_graphical_data()
    assert graphical_data['total_rooms'] == 2
    assert len(graphical_data['rooms']) == 2
    # Note: connections count is 2 because each connection is stored from both directions
    assert len(graphical_data['connections']) >= 1
    assert 'bounds' in graphical_data
    
    print("✓ Map visualization works correctly")

def test_direction_utilities():
    """Test direction-related utility functions"""
    print("Testing direction utilities...")
    
    world_map = WorldMap()
    
    # Test reverse direction mapping
    assert world_map._get_reverse_direction(Direction.NORTH) == Direction.SOUTH
    assert world_map._get_reverse_direction(Direction.EAST) == Direction.WEST
    assert world_map._get_reverse_direction(Direction.NORTHEAST) == Direction.SOUTHWEST
    assert world_map._get_reverse_direction(Direction.UP) == Direction.DOWN
    
    # Test exit position calculation
    room = Room("test", "Test", "Test room.", 10.0, 8.0, 5.0, 3.0)
    
    north_pos = room.get_exit_position(Direction.NORTH)
    assert north_pos == (5.0, 7.0)  # center_x, center_y + half_height
    
    east_pos = room.get_exit_position(Direction.EAST)
    assert east_pos == (10.0, 3.0)  # center_x + half_width, center_y
    
    print("✓ Direction utilities work correctly")

def run_all_tests():
    """Run all test functions"""
    print("Running World Map System Tests")
    print("=" * 40)
    
    tests = [
        test_basic_room_creation,
        test_room_connections,
        test_exit_discovery,
        test_spatial_positioning,
        test_conflict_resolution,
        test_map_bounds,
        test_rooms_in_area,
        test_persistence,
        test_visualization,
        test_direction_utilities,
    ]
    
    passed = 0
    failed = 0
    
    for test_func in tests:
        try:
            test_func()
            passed += 1
        except Exception as e:
            print(f"✗ {test_func.__name__} FAILED: {e}")
            failed += 1
    
    print("\n" + "=" * 40)
    print(f"Test Results: {passed} passed, {failed} failed")
    
    if failed == 0:
        print("🎉 All tests passed!")
        return True
    else:
        print("❌ Some tests failed!")
        return False

if __name__ == "__main__":
    success = run_all_tests()
    exit(0 if success else 1)