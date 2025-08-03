import math
from typing import Dict, List, Optional, Tuple, Set
from dataclasses import dataclass, field
from enum import Enum
import json

class Direction(Enum):
    """Cardinal and intercardinal directions for exits"""
    NORTH = "north"
    SOUTH = "south"
    EAST = "east"
    WEST = "west"
    NORTHEAST = "northeast"
    NORTHWEST = "northwest"
    SOUTHEAST = "southeast"
    SOUTHWEST = "southwest"
    UP = "up"
    DOWN = "down"

@dataclass
class RoomExit:
    """Represents an exit from a room"""
    direction: Direction
    target_room_id: Optional[str] = None  # None if exit leads to unexplored area
    is_discovered: bool = False
    
@dataclass
class Room:
    """Represents a discovered room in the world"""
    id: str
    name: str
    description: str
    width: float  # Room dimensions in world units
    height: float
    x: float = 0.0  # World coordinates (center of room)
    y: float = 0.0
    exits: Dict[Direction, RoomExit] = field(default_factory=dict)
    is_starting_room: bool = False
    discovery_order: int = 0  # Order in which room was discovered
    
    def add_exit(self, direction: Direction, target_room_id: Optional[str] = None):
        """Add an exit to this room"""
        self.exits[direction] = RoomExit(direction, target_room_id, target_room_id is not None)
    
    def get_exit_position(self, direction: Direction) -> Tuple[float, float]:
        """Get the world coordinates where an exit is located"""
        half_width = self.width / 2
        half_height = self.height / 2
        
        exit_positions = {
            Direction.NORTH: (self.x, self.y + half_height),
            Direction.SOUTH: (self.x, self.y - half_height),
            Direction.EAST: (self.x + half_width, self.y),
            Direction.WEST: (self.x - half_width, self.y),
            Direction.NORTHEAST: (self.x + half_width * 0.7, self.y + half_height * 0.7),
            Direction.NORTHWEST: (self.x - half_width * 0.7, self.y + half_height * 0.7),
            Direction.SOUTHEAST: (self.x + half_width * 0.7, self.y - half_height * 0.7),
            Direction.SOUTHWEST: (self.x - half_width * 0.7, self.y - half_height * 0.7),
            # UP/DOWN exits are at room center for now
            Direction.UP: (self.x, self.y),
            Direction.DOWN: (self.x, self.y),
        }
        
        return exit_positions.get(direction, (self.x, self.y))

class WorldMap:
    """Manages the discovered world map and spatial relationships between rooms"""
    
    def __init__(self):
        self.rooms: Dict[str, Room] = {}
        self.discovery_count = 0
        self.starting_room_id: Optional[str] = None
        
    def add_room(self, room_id: str, name: str, description: str, 
                 width: float, height: float, is_starting_room: bool = False) -> Room:
        """Add a new room to the world map"""
        if room_id in self.rooms:
            return self.rooms[room_id]
            
        self.discovery_count += 1
        room = Room(
            id=room_id,
            name=name,
            description=description,
            width=width,
            height=height,
            is_starting_room=is_starting_room,
            discovery_order=self.discovery_count
        )
        
        # Position the starting room at origin
        if is_starting_room or len(self.rooms) == 0:
            room.x = 0.0
            room.y = 0.0
            self.starting_room_id = room_id
        
        self.rooms[room_id] = room
        return room
    
    def connect_rooms(self, from_room_id: str, direction: Direction, to_room_id: str):
        """Connect two rooms with an exit in the specified direction"""
        if from_room_id not in self.rooms or to_room_id not in self.rooms:
            raise ValueError("Both rooms must exist before connecting them")
            
        from_room = self.rooms[from_room_id]
        to_room = self.rooms[to_room_id]
        
        # Add the exit to the from_room
        from_room.add_exit(direction, to_room_id)
        
        # Add reverse exit to to_room
        reverse_direction = self._get_reverse_direction(direction)
        if reverse_direction:
            to_room.add_exit(reverse_direction, from_room_id)
        
        # Position the target room if it hasn't been positioned yet
        if to_room.x == 0.0 and to_room.y == 0.0 and not to_room.is_starting_room:
            self._position_room_relative_to(to_room, from_room, direction)
    
    def discover_exit(self, room_id: str, direction: Direction):
        """Mark that an exit exists in a room but doesn't lead to a known room yet"""
        if room_id in self.rooms:
            self.rooms[room_id].add_exit(direction)
    
    def get_room(self, room_id: str) -> Optional[Room]:
        """Get a room by ID"""
        return self.rooms.get(room_id)
    
    def get_rooms_in_area(self, center_x: float, center_y: float, radius: float) -> List[Room]:
        """Get all rooms within a circular area"""
        nearby_rooms = []
        for room in self.rooms.values():
            distance = math.sqrt((room.x - center_x)**2 + (room.y - center_y)**2)
            if distance <= radius:
                nearby_rooms.append(room)
        return nearby_rooms
    
    def get_map_bounds(self) -> Tuple[float, float, float, float]:
        """Get the bounding box of all discovered rooms (min_x, min_y, max_x, max_y)"""
        if not self.rooms:
            return (0, 0, 0, 0)
            
        min_x = min(room.x - room.width/2 for room in self.rooms.values())
        max_x = max(room.x + room.width/2 for room in self.rooms.values())
        min_y = min(room.y - room.height/2 for room in self.rooms.values())
        max_y = max(room.y + room.height/2 for room in self.rooms.values())
        
        return (min_x, min_y, max_x, max_y)
    
    def _position_room_relative_to(self, target_room: Room, reference_room: Room, direction: Direction):
        """Position a room relative to another room based on the direction of connection"""
        # Calculate the base distance (room edges should be touching)
        base_distance_x = (reference_room.width + target_room.width) / 2
        base_distance_y = (reference_room.height + target_room.height) / 2
        
        # Direction vectors for positioning
        direction_vectors = {
            Direction.NORTH: (0, base_distance_y),
            Direction.SOUTH: (0, -base_distance_y),
            Direction.EAST: (base_distance_x, 0),
            Direction.WEST: (-base_distance_x, 0),
            Direction.NORTHEAST: (base_distance_x * 0.7, base_distance_y * 0.7),
            Direction.NORTHWEST: (-base_distance_x * 0.7, base_distance_y * 0.7),
            Direction.SOUTHEAST: (base_distance_x * 0.7, -base_distance_y * 0.7),
            Direction.SOUTHWEST: (-base_distance_x * 0.7, -base_distance_y * 0.7),
            Direction.UP: (0, 0),  # Same position for vertical connections
            Direction.DOWN: (0, 0),
        }
        
        offset_x, offset_y = direction_vectors.get(direction, (0, 0))
        target_room.x = reference_room.x + offset_x
        target_room.y = reference_room.y + offset_y
        
        # Check for overlaps and resolve conflicts
        self._resolve_spatial_conflicts(target_room)
    
    def _resolve_spatial_conflicts(self, new_room: Room):
        """Resolve spatial conflicts when rooms would overlap"""
        max_iterations = 10
        iteration = 0
        
        while iteration < max_iterations:
            conflict_found = False
            
            for existing_room in self.rooms.values():
                if existing_room.id == new_room.id:
                    continue
                    
                if self._rooms_overlap(new_room, existing_room):
                    # Move the new room away from the conflict
                    self._move_room_away_from_conflict(new_room, existing_room)
                    conflict_found = True
                    break
            
            if not conflict_found:
                break
                
            iteration += 1
    
    def _rooms_overlap(self, room1: Room, room2: Room) -> bool:
        """Check if two rooms overlap"""
        # Calculate room boundaries
        r1_left = room1.x - room1.width / 2
        r1_right = room1.x + room1.width / 2
        r1_bottom = room1.y - room1.height / 2
        r1_top = room1.y + room1.height / 2
        
        r2_left = room2.x - room2.width / 2
        r2_right = room2.x + room2.width / 2
        r2_bottom = room2.y - room2.height / 2
        r2_top = room2.y + room2.height / 2
        
        # Check for overlap
        return not (r1_right <= r2_left or r1_left >= r2_right or 
                   r1_top <= r2_bottom or r1_bottom >= r2_top)
    
    def _move_room_away_from_conflict(self, moving_room: Room, blocking_room: Room):
        """Move a room away from another room that it's conflicting with"""
        # Calculate direction vector from blocking room to moving room
        dx = moving_room.x - blocking_room.x
        dy = moving_room.y - blocking_room.y
        
        # Normalize the direction
        distance = math.sqrt(dx*dx + dy*dy)
        if distance == 0:
            # Rooms are at same position, move arbitrarily
            dx, dy = 1, 0
            distance = 1
        else:
            dx /= distance
            dy /= distance
        
        # Calculate minimum distance to avoid overlap
        min_distance_x = (moving_room.width + blocking_room.width) / 2 + 1  # +1 for buffer
        min_distance_y = (moving_room.height + blocking_room.height) / 2 + 1
        min_distance = max(min_distance_x, min_distance_y)
        
        # Move the room
        moving_room.x = blocking_room.x + dx * min_distance
        moving_room.y = blocking_room.y + dy * min_distance
    
    def _get_reverse_direction(self, direction: Direction) -> Optional[Direction]:
        """Get the reverse of a direction"""
        reverse_map = {
            Direction.NORTH: Direction.SOUTH,
            Direction.SOUTH: Direction.NORTH,
            Direction.EAST: Direction.WEST,
            Direction.WEST: Direction.EAST,
            Direction.NORTHEAST: Direction.SOUTHWEST,
            Direction.SOUTHWEST: Direction.NORTHEAST,
            Direction.NORTHWEST: Direction.SOUTHEAST,
            Direction.SOUTHEAST: Direction.NORTHWEST,
            Direction.UP: Direction.DOWN,
            Direction.DOWN: Direction.UP,
        }
        return reverse_map.get(direction)
    
    def to_dict(self) -> dict:
        """Serialize the world map to a dictionary"""
        return {
            'discovery_count': self.discovery_count,
            'starting_room_id': self.starting_room_id,
            'rooms': {
                room_id: {
                    'id': room.id,
                    'name': room.name,
                    'description': room.description,
                    'width': room.width,
                    'height': room.height,
                    'x': room.x,
                    'y': room.y,
                    'is_starting_room': room.is_starting_room,
                    'discovery_order': room.discovery_order,
                    'exits': {
                        direction.value: {
                            'direction': exit_info.direction.value,
                            'target_room_id': exit_info.target_room_id,
                            'is_discovered': exit_info.is_discovered
                        }
                        for direction, exit_info in room.exits.items()
                    }
                }
                for room_id, room in self.rooms.items()
            }
        }
    
    @classmethod
    def from_dict(cls, data: dict) -> 'WorldMap':
        """Deserialize a world map from a dictionary"""
        world_map = cls()
        world_map.discovery_count = data.get('discovery_count', 0)
        world_map.starting_room_id = data.get('starting_room_id')
        
        # First pass: create all rooms
        for room_id, room_data in data.get('rooms', {}).items():
            room = Room(
                id=room_data['id'],
                name=room_data['name'],
                description=room_data['description'],
                width=room_data['width'],
                height=room_data['height'],
                x=room_data['x'],
                y=room_data['y'],
                is_starting_room=room_data.get('is_starting_room', False),
                discovery_order=room_data.get('discovery_order', 0)
            )
            world_map.rooms[room_id] = room
        
        # Second pass: restore exits
        for room_id, room_data in data.get('rooms', {}).items():
            room = world_map.rooms[room_id]
            for direction_str, exit_data in room_data.get('exits', {}).items():
                direction = Direction(direction_str)
                room_exit = RoomExit(
                    direction=Direction(exit_data['direction']),
                    target_room_id=exit_data.get('target_room_id'),
                    is_discovered=exit_data.get('is_discovered', False)
                )
                room.exits[direction] = room_exit
        
        return world_map
    
    def save_to_file(self, filename: str):
        """Save the world map to a JSON file"""
        with open(filename, 'w') as f:
            json.dump(self.to_dict(), f, indent=2)
    
    @classmethod
    def load_from_file(cls, filename: str) -> 'WorldMap':
        """Load a world map from a JSON file"""
        with open(filename, 'r') as f:
            data = json.load(f)
        return cls.from_dict(data)