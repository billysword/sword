package world

import (
	"testing"
	"image"
	"github.com/hajimehoshi/ebiten/v2"
	"sword/entities"
)

// MockRoom implements the Room interface for testing
type MockRoom struct {
	zoneID   string
	tileMap  *TileMap
	name     string
}

func (mr *MockRoom) GetTileMap() *TileMap {
	return mr.tileMap
}

func (mr *MockRoom) GetZoneID() string {
	return mr.zoneID
}

func (mr *MockRoom) Update(player *entities.Player) error {
	return nil
}

func (mr *MockRoom) HandleCollisions(player *entities.Player) {
	// Mock implementation
}

func (mr *MockRoom) OnEnter(player *entities.Player) {
	// Mock implementation
}

func (mr *MockRoom) OnExit(player *entities.Player) {
	// Mock implementation
}

func (mr *MockRoom) FindFloorAtX(x int) int {
	return 0
}

func (mr *MockRoom) Draw(screen *ebiten.Image) {
	// Mock implementation
}

func (mr *MockRoom) DrawTiles(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image) {
	// Mock implementation
}

func (mr *MockRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Mock implementation
}

// createMockRoom creates a mock room for testing
func createMockRoom(zoneID string, width, height int) *MockRoom {
	tileMap := NewTileMap(width, height)
	// Add some test tiles
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if y == height-1 || x == 0 || x == width-1 {
				tileMap.SetTile(x, y, 0) // Border tiles
			} else {
				tileMap.SetTile(x, y, -1) // Empty space
			}
		}
	}
	
	return &MockRoom{
		zoneID:  zoneID,
		tileMap: tileMap,
		name:    zoneID,
	}
}

func TestWorldMapCreation(t *testing.T) {
	worldMap := NewWorldMap()
	
	if worldMap == nil {
		t.Fatal("NewWorldMap() returned nil")
	}
	
	if worldMap.discoveredRooms == nil {
		t.Error("discoveredRooms map not initialized")
	}
	
	if worldMap.roomConnections == nil {
		t.Error("roomConnections map not initialized")
	}
	
	if len(worldMap.playerTrail) != 0 {
		t.Error("playerTrail should be empty initially")
	}
}

func TestRoomDiscovery(t *testing.T) {
	worldMap := NewWorldMap()
	mockRoom := createMockRoom("test_room", 10, 8)
	
	// Discover the room
	worldMap.DiscoverRoom(mockRoom)
	
	// Check if room was added
	discoveredRooms := worldMap.GetDiscoveredRooms()
	if len(discoveredRooms) != 1 {
		t.Errorf("Expected 1 discovered room, got %d", len(discoveredRooms))
	}
	
	room, exists := discoveredRooms["test_room"]
	if !exists {
		t.Error("Room 'test_room' not found in discovered rooms")
	}
	
	if room.ZoneID != "test_room" {
		t.Errorf("Expected room ID 'test_room', got '%s'", room.ZoneID)
	}
	
	if !room.IsExplored {
		t.Error("Room should be marked as explored")
	}
}

func TestRoomConnections(t *testing.T) {
	worldMap := NewWorldMap()
	room1 := createMockRoom("room1", 5, 5)
	room2 := createMockRoom("room2", 5, 5)
	
	// Discover both rooms
	worldMap.DiscoverRoom(room1)
	worldMap.DiscoverRoom(room2)
	
	// Connect rooms
	err := worldMap.ConnectRooms("room1", East, "room2")
	if err != nil {
		t.Errorf("Failed to connect rooms: %v", err)
	}
	
	// Check connections
	connections1 := worldMap.GetRoomConnections("room1")
	if len(connections1) != 1 {
		t.Errorf("Expected 1 connection from room1, got %d", len(connections1))
	}
	
	if connections1[East] != "room2" {
		t.Errorf("Expected connection to room2 in East direction, got %s", connections1[East])
	}
	
	// Check reverse connection
	connections2 := worldMap.GetRoomConnections("room2")
	if len(connections2) != 1 {
		t.Errorf("Expected 1 connection from room2, got %d", len(connections2))
	}
	
	if connections2[West] != "room1" {
		t.Errorf("Expected connection to room1 in West direction, got %s", connections2[West])
	}
}

func TestPlayerTrail(t *testing.T) {
	worldMap := NewWorldMap()
	
	// Add some positions
	worldMap.AddPlayerPosition(100, 200)
	worldMap.AddPlayerPosition(150, 250)
	worldMap.AddPlayerPosition(200, 300)
	
	trail := worldMap.GetPlayerTrail()
	if len(trail) != 3 {
		t.Errorf("Expected 3 trail positions, got %d", len(trail))
	}
	
	// Check first position
	if trail[0].X != 100 || trail[0].Y != 200 {
		t.Errorf("Expected first position (100, 200), got (%d, %d)", trail[0].X, trail[0].Y)
	}
	
	// Check last position
	if trail[2].X != 200 || trail[2].Y != 300 {
		t.Errorf("Expected last position (200, 300), got (%d, %d)", trail[2].X, trail[2].Y)
	}
}

func TestMapBounds(t *testing.T) {
	worldMap := NewWorldMap()
	
	// Empty map should return zero bounds
	bounds := worldMap.GetMapBounds()
	if bounds != (image.Rectangle{}) {
		t.Errorf("Expected empty bounds for empty map, got %v", bounds)
	}
	
	// Add a room
	room := createMockRoom("test_room", 10, 8)
	worldMap.DiscoverRoom(room)
	
	// Check bounds
	bounds = worldMap.GetMapBounds()
	if bounds.Empty() {
		t.Error("Map bounds should not be empty after adding a room")
	}
}

func TestDirectionHelpers(t *testing.T) {
	// Test opposite directions
	testCases := []struct {
		direction Direction
		opposite  Direction
	}{
		{North, South},
		{South, North},
		{East, West},
		{West, East},
		{Up, Down},
		{Down, Up},
	}
	
	for _, tc := range testCases {
		opposite := tc.direction.Opposite()
		if opposite != tc.opposite {
			t.Errorf("Expected %s.Opposite() to be %s, got %s", 
				tc.direction.String(), tc.opposite.String(), opposite.String())
		}
	}
}

func TestThumbnailGeneration(t *testing.T) {
	// Create a test tile map
	tileMap := NewTileMap(4, 4)
	
	// Set some tiles
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if y == 3 { // Bottom row
				tileMap.SetTile(x, y, 0) // Ground tiles
			} else {
				tileMap.SetTile(x, y, -1) // Empty
			}
		}
	}
	
	thumbnail := generateThumbnail(tileMap)
	
	if len(thumbnail) == 0 {
		t.Error("Thumbnail should not be empty")
	}
	
	if len(thumbnail[0]) == 0 {
		t.Error("Thumbnail rows should not be empty")
	}
	
	// Thumbnail should be no larger than 32x32
	if len(thumbnail) > 32 {
		t.Errorf("Thumbnail height should be <= 32, got %d", len(thumbnail))
	}
	
	if len(thumbnail[0]) > 32 {
		t.Errorf("Thumbnail width should be <= 32, got %d", len(thumbnail[0]))
	}
}

func TestJSONSerialization(t *testing.T) {
	worldMap := NewWorldMap()
	room := createMockRoom("test_room", 5, 5)
	
	worldMap.DiscoverRoom(room)
	worldMap.SetCurrentRoom("test_room")
	worldMap.AddPlayerPosition(100, 200)
	
	// Serialize to JSON
	jsonData, err := worldMap.ToJSON()
	if err != nil {
		t.Errorf("Failed to serialize to JSON: %v", err)
	}
	
	// Deserialize from JSON
	newWorldMap := NewWorldMap()
	err = newWorldMap.FromJSON(jsonData)
	if err != nil {
		t.Errorf("Failed to deserialize from JSON: %v", err)
	}
	
	// Check that data was preserved
	if newWorldMap.GetCurrentRoom() != "test_room" {
		t.Errorf("Current room not preserved after JSON round trip")
	}
	
	discoveredRooms := newWorldMap.GetDiscoveredRooms()
	if len(discoveredRooms) != 1 {
		t.Errorf("Expected 1 room after deserialization, got %d", len(discoveredRooms))
	}
	
	trail := newWorldMap.GetPlayerTrail()
	if len(trail) != 1 {
		t.Errorf("Expected 1 trail position after deserialization, got %d", len(trail))
	}
}