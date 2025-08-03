from typing import List, Tuple, Optional
import math
from world_map import WorldMap, Room, Direction

class MapVisualizer:
    """Provides various ways to visualize the discovered world map"""
    
    def __init__(self, world_map: WorldMap):
        self.world_map = world_map
        
    def render_ascii_map(self, width: int = 80, height: int = 24, 
                        current_room_id: Optional[str] = None) -> str:
        """Render the world map as ASCII art"""
        if not self.world_map.rooms:
            return "No rooms discovered yet."
            
        # Get map bounds and calculate scale
        min_x, min_y, max_x, max_y = self.world_map.get_map_bounds()
        
        # Add padding
        padding = 2.0
        min_x -= padding
        min_y -= padding
        max_x += padding
        max_y += padding
        
        # Calculate scale factors
        world_width = max_x - min_x
        world_height = max_y - min_y
        
        if world_width == 0:
            world_width = 1
        if world_height == 0:
            world_height = 1
            
        scale_x = (width - 2) / world_width
        scale_y = (height - 2) / world_height
        scale = min(scale_x, scale_y)  # Use smaller scale to fit both dimensions
        
        # Create the grid
        grid = [[' ' for _ in range(width)] for _ in range(height)]
        
        # Draw rooms
        for room in self.world_map.rooms.values():
            self._draw_room_on_grid(grid, room, min_x, min_y, scale, 
                                  width, height, current_room_id)
        
        # Convert grid to string
        result = []
        for row in reversed(grid):  # Reverse to have y=0 at bottom
            result.append(''.join(row))
        
        return '\n'.join(result)
    
    def _draw_room_on_grid(self, grid: List[List[str]], room: Room, 
                          min_x: float, min_y: float, scale: float,
                          grid_width: int, grid_height: int, 
                          current_room_id: Optional[str]):
        """Draw a single room on the ASCII grid"""
        # Convert world coordinates to grid coordinates
        room_left = int((room.x - room.width/2 - min_x) * scale)
        room_right = int((room.x + room.width/2 - min_x) * scale)
        room_bottom = int((room.y - room.height/2 - min_y) * scale)
        room_top = int((room.y + room.height/2 - min_y) * scale)
        
        # Clamp to grid bounds
        room_left = max(0, min(room_left, grid_width - 1))
        room_right = max(0, min(room_right, grid_width - 1))
        room_bottom = max(0, min(room_bottom, grid_height - 1))
        room_top = max(0, min(room_top, grid_height - 1))
        
        # Choose room character
        if room.id == current_room_id:
            room_char = '@'  # Current room
        elif room.is_starting_room:
            room_char = 'S'  # Starting room
        else:
            room_char = '#'  # Regular room
        
        # Draw room borders
        for y in range(room_bottom, room_top + 1):
            for x in range(room_left, room_right + 1):
                if (y == room_bottom or y == room_top or 
                    x == room_left or x == room_right):
                    if grid[y][x] == ' ':
                        grid[y][x] = '+'
                elif grid[y][x] == ' ':
                    grid[y][x] = '.'
        
        # Draw room center
        center_x = int((room.x - min_x) * scale)
        center_y = int((room.y - min_y) * scale)
        
        if (0 <= center_x < grid_width and 0 <= center_y < grid_height):
            grid[center_y][center_x] = room_char
        
        # Draw exits
        self._draw_exits_on_grid(grid, room, min_x, min_y, scale, 
                               grid_width, grid_height)
    
    def _draw_exits_on_grid(self, grid: List[List[str]], room: Room,
                           min_x: float, min_y: float, scale: float,
                           grid_width: int, grid_height: int):
        """Draw room exits on the ASCII grid"""
        for direction, room_exit in room.exits.items():
            exit_x, exit_y = room.get_exit_position(direction)
            
            # Convert to grid coordinates
            grid_x = int((exit_x - min_x) * scale)
            grid_y = int((exit_y - min_y) * scale)
            
            if (0 <= grid_x < grid_width and 0 <= grid_y < grid_height):
                if room_exit.is_discovered and room_exit.target_room_id:
                    # Connected exit
                    exit_char = self._get_direction_char(direction)
                else:
                    # Unconnected exit
                    exit_char = '?'
                
                if grid[grid_y][grid_x] in [' ', '.', '+']:
                    grid[grid_y][grid_x] = exit_char
    
    def _get_direction_char(self, direction: Direction) -> str:
        """Get ASCII character representing a direction"""
        direction_chars = {
            Direction.NORTH: '^',
            Direction.SOUTH: 'v',
            Direction.EAST: '>',
            Direction.WEST: '<',
            Direction.NORTHEAST: '/',
            Direction.NORTHWEST: '\\',
            Direction.SOUTHEAST: '\\',
            Direction.SOUTHWEST: '/',
            Direction.UP: '^',
            Direction.DOWN: 'v',
        }
        return direction_chars.get(direction, '?')
    
    def get_room_list(self, sort_by: str = 'discovery') -> List[Room]:
        """Get a list of rooms sorted by various criteria"""
        rooms = list(self.world_map.rooms.values())
        
        if sort_by == 'discovery':
            rooms.sort(key=lambda r: r.discovery_order)
        elif sort_by == 'name':
            rooms.sort(key=lambda r: r.name.lower())
        elif sort_by == 'size':
            rooms.sort(key=lambda r: r.width * r.height, reverse=True)
        elif sort_by == 'distance_from_start':
            if self.world_map.starting_room_id:
                start_room = self.world_map.get_room(self.world_map.starting_room_id)
                if start_room:
                    rooms.sort(key=lambda r: math.sqrt(
                        (r.x - start_room.x)**2 + (r.y - start_room.y)**2
                    ))
        
        return rooms
    
    def print_room_details(self, room_id: str) -> str:
        """Generate detailed text description of a room"""
        room = self.world_map.get_room(room_id)
        if not room:
            return f"Room '{room_id}' not found."
        
        lines = []
        lines.append(f"Room: {room.name} ({room.id})")
        lines.append(f"Description: {room.description}")
        lines.append(f"Size: {room.width} x {room.height}")
        lines.append(f"Position: ({room.x:.1f}, {room.y:.1f})")
        lines.append(f"Discovery Order: {room.discovery_order}")
        
        if room.is_starting_room:
            lines.append("* Starting Room *")
        
        if room.exits:
            lines.append("\nExits:")
            for direction, room_exit in room.exits.items():
                if room_exit.is_discovered and room_exit.target_room_id:
                    target_room = self.world_map.get_room(room_exit.target_room_id)
                    target_name = target_room.name if target_room else "Unknown"
                    lines.append(f"  {direction.value}: {target_name} ({room_exit.target_room_id})")
                else:
                    lines.append(f"  {direction.value}: Unexplored")
        else:
            lines.append("\nNo exits discovered.")
        
        return '\n'.join(lines)
    
    def print_map_summary(self) -> str:
        """Generate a summary of the discovered map"""
        if not self.world_map.rooms:
            return "No rooms discovered yet."
        
        lines = []
        lines.append(f"World Map Summary")
        lines.append(f"================")
        lines.append(f"Total Rooms Discovered: {len(self.world_map.rooms)}")
        
        min_x, min_y, max_x, max_y = self.world_map.get_map_bounds()
        lines.append(f"Map Bounds: ({min_x:.1f}, {min_y:.1f}) to ({max_x:.1f}, {max_y:.1f})")
        lines.append(f"Map Size: {max_x - min_x:.1f} x {max_y - min_y:.1f}")
        
        # Room statistics
        total_area = sum(room.width * room.height for room in self.world_map.rooms.values())
        avg_area = total_area / len(self.world_map.rooms)
        lines.append(f"Total Room Area: {total_area:.1f}")
        lines.append(f"Average Room Size: {avg_area:.1f}")
        
        # Exit statistics
        total_exits = sum(len(room.exits) for room in self.world_map.rooms.values())
        connected_exits = sum(
            1 for room in self.world_map.rooms.values()
            for exit in room.exits.values()
            if exit.is_discovered and exit.target_room_id
        )
        lines.append(f"Total Exits: {total_exits}")
        lines.append(f"Connected Exits: {connected_exits}")
        lines.append(f"Unexplored Exits: {total_exits - connected_exits}")
        
        # Starting room info
        if self.world_map.starting_room_id:
            start_room = self.world_map.get_room(self.world_map.starting_room_id)
            if start_room:
                lines.append(f"Starting Room: {start_room.name} ({start_room.id})")
        
        return '\n'.join(lines)
    
    def get_graphical_data(self) -> dict:
        """Get data structure suitable for graphical rendering"""
        rooms_data = []
        connections_data = []
        
        for room in self.world_map.rooms.values():
            room_data = {
                'id': room.id,
                'name': room.name,
                'description': room.description,
                'x': room.x,
                'y': room.y,
                'width': room.width,
                'height': room.height,
                'is_starting_room': room.is_starting_room,
                'discovery_order': room.discovery_order
            }
            rooms_data.append(room_data)
            
            # Add connections
            for direction, room_exit in room.exits.items():
                if room_exit.is_discovered and room_exit.target_room_id:
                    exit_x, exit_y = room.get_exit_position(direction)
                    target_room = self.world_map.get_room(room_exit.target_room_id)
                    
                    if target_room:
                        target_exit_x, target_exit_y = target_room.get_exit_position(
                            self.world_map._get_reverse_direction(direction)
                        )
                        
                        connection_data = {
                            'from_room': room.id,
                            'to_room': room_exit.target_room_id,
                            'from_x': exit_x,
                            'from_y': exit_y,
                            'to_x': target_exit_x,
                            'to_y': target_exit_y,
                            'direction': direction.value
                        }
                        connections_data.append(connection_data)
        
        min_x, min_y, max_x, max_y = self.world_map.get_map_bounds()
        
        return {
            'rooms': rooms_data,
            'connections': connections_data,
            'bounds': {
                'min_x': min_x,
                'min_y': min_y,
                'max_x': max_x,
                'max_y': max_y,
                'width': max_x - min_x,
                'height': max_y - min_y
            },
            'starting_room_id': self.world_map.starting_room_id,
            'total_rooms': len(self.world_map.rooms)
        }