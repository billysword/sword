package world

import (
	"encoding/json"
	"fmt"
	"image"
	"sync"
	"sword/engine"
)

// Direction represents movement directions between rooms
type Direction int

const (
	North Direction = iota
	South
	East
	West
	Up   // For vertical level connections
	Down
)

// String returns the string representation of a direction
func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	case Up:
		return "Up"
	case Down:
		return "Down"
	default:
		return "Unknown"
	}
}

// Opposite returns the opposite direction
func (d Direction) Opposite() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	case Up:
		return Down
	case Down:
		return Up
	default:
		return d
	}
}

// Point represents a 2D coordinate
type Point struct {
	X, Y int
}

// DiscoveredRoom represents a room that the player has discovered
type DiscoveredRoom struct {
	ZoneID        string                 `json:"zone_id"`        // Unique room identifier
	Name          string                 `json:"name"`           // Display name for the room
	Bounds        image.Rectangle        `json:"bounds"`         // World-space boundaries (pixel coordinates)
	ExitPoints    map[Direction]Point    `json:"exit_points"`    // Known exit locations
	IsExplored    bool                   `json:"is_explored"`    // Has player visited this room
	ThumbnailData [][]int                `json:"thumbnail_data"` // Simplified tile data for mini-map rendering
	WorldPos      Point                  `json:"world_pos"`      // Position in world map coordinate system
}

// NewDiscoveredRoom creates a new discovered room from a Room interface
func NewDiscoveredRoom(room Room) *DiscoveredRoom {
	tileMap := room.GetTileMap()
	u := engine.GetPhysicsUnit()
	
	bounds := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: tileMap.Width * u, Y: tileMap.Height * u},
	}
	
	return &DiscoveredRoom{
		ZoneID:        room.GetZoneID(),
		Name:          room.GetZoneID(), // TODO: Get display name when Room interface is extended
		Bounds:        bounds,
		ExitPoints:    make(map[Direction]Point),
		IsExplored:    true,
		ThumbnailData: generateThumbnail(tileMap),
		WorldPos:      Point{X: 0, Y: 0}, // Will be positioned when connected to other rooms
	}
}

// generateThumbnail creates a simplified representation of the room for mini-map display
func generateThumbnail(tileMap *TileMap) [][]int {
	// Create a downscaled version of the tile map for thumbnail display
	thumbnailWidth := min(tileMap.Width, 32)   // Max 32x32 thumbnail
	thumbnailHeight := min(tileMap.Height, 32)
	
	scaleX := float64(tileMap.Width) / float64(thumbnailWidth)
	scaleY := float64(tileMap.Height) / float64(thumbnailHeight)
	
	thumbnail := make([][]int, thumbnailHeight)
	for y := 0; y < thumbnailHeight; y++ {
		thumbnail[y] = make([]int, thumbnailWidth)
		for x := 0; x < thumbnailWidth; x++ {
			// Sample the original tile map
			origX := int(float64(x) * scaleX)
			origY := int(float64(y) * scaleY)
			
			if origX < tileMap.Width && origY < tileMap.Height {
				thumbnail[y][x] = tileMap.GetTileIndex(origX, origY)
			} else {
				thumbnail[y][x] = -1 // Empty tile
			}
		}
	}
	
	return thumbnail
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// WorldMap manages the discovered world map and spatial relationships between rooms
type WorldMap struct {
	discoveredRooms map[string]*DiscoveredRoom      // discovered rooms by ID
	currentRoomID   string                          // current room ID
	roomConnections map[string]map[Direction]string // room_id -> direction -> connected_room_id
	playerTrail     []Point                         // recent player movement for mini-map
	mutex           sync.RWMutex                    // thread safety
}

// NewWorldMap creates a new world map manager
func NewWorldMap() *WorldMap {
	return &WorldMap{
		discoveredRooms: make(map[string]*DiscoveredRoom),
		roomConnections: make(map[string]map[Direction]string),
		playerTrail:     make([]Point, 0, 100), // Keep last 100 positions
		mutex:           sync.RWMutex{},
	}
}

// DiscoverRoom adds a new room to the discovered map or updates an existing one
func (wm *WorldMap) DiscoverRoom(room Room) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	
	zoneID := room.GetZoneID()
	
	// Check if room is already discovered
	if existingRoom, exists := wm.discoveredRooms[zoneID]; exists {
		existingRoom.IsExplored = true
		return
	}
	
	// Create new discovered room
	discoveredRoom := NewDiscoveredRoom(room)
	wm.discoveredRooms[zoneID] = discoveredRoom
	
	// Initialize connections map for this room
	wm.roomConnections[zoneID] = make(map[Direction]string)
	
	// If this is the first room, position it at origin
	if len(wm.discoveredRooms) == 1 {
		discoveredRoom.WorldPos = Point{X: 0, Y: 0}
	}
}

// SetCurrentRoom updates the current room the player is in
func (wm *WorldMap) SetCurrentRoom(zoneID string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	
	wm.currentRoomID = zoneID
}

// GetCurrentRoom returns the current room ID
func (wm *WorldMap) GetCurrentRoom() string {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	
	return wm.currentRoomID
}

// ConnectRooms establishes a connection between two rooms
func (wm *WorldMap) ConnectRooms(fromRoomID string, direction Direction, toRoomID string) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	
	// Verify both rooms exist
	fromRoom, fromExists := wm.discoveredRooms[fromRoomID]
	toRoom, toExists := wm.discoveredRooms[toRoomID]
	
	if !fromExists {
		return fmt.Errorf("source room %s not found", fromRoomID)
	}
	if !toExists {
		return fmt.Errorf("destination room %s not found", toRoomID)
	}
	
	// Establish bidirectional connection
	wm.roomConnections[fromRoomID][direction] = toRoomID
	if wm.roomConnections[toRoomID] == nil {
		wm.roomConnections[toRoomID] = make(map[Direction]string)
	}
	wm.roomConnections[toRoomID][direction.Opposite()] = fromRoomID
	
	// Position the destination room relative to the source room
	wm.positionRoom(toRoom, fromRoom, direction)
	
	return nil
}

// positionRoom positions a room relative to another room based on connection direction
func (wm *WorldMap) positionRoom(targetRoom, sourceRoom *DiscoveredRoom, direction Direction) {
	// Calculate spacing based on room sizes
	sourceWidth := sourceRoom.Bounds.Dx()
	sourceHeight := sourceRoom.Bounds.Dy()
	targetWidth := targetRoom.Bounds.Dx()
	targetHeight := targetRoom.Bounds.Dy()
	
	// Add some padding between rooms
	padding := 64
	
	switch direction {
	case North:
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y - (sourceHeight/2 + targetHeight/2 + padding)
	case South:
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y + (sourceHeight/2 + targetHeight/2 + padding)
	case East:
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X + (sourceWidth/2 + targetWidth/2 + padding)
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y
	case West:
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X - (sourceWidth/2 + targetWidth/2 + padding)
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y
	case Up:
		// For vertical connections, keep same position but mark as different level
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y
	case Down:
		// For vertical connections, keep same position but mark as different level
		targetRoom.WorldPos.X = sourceRoom.WorldPos.X
		targetRoom.WorldPos.Y = sourceRoom.WorldPos.Y
	}
}

// GetDiscoveredRooms returns a copy of all discovered rooms
func (wm *WorldMap) GetDiscoveredRooms() map[string]*DiscoveredRoom {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	rooms := make(map[string]*DiscoveredRoom)
	for id, room := range wm.discoveredRooms {
		// Create a copy of the room
		roomCopy := *room
		rooms[id] = &roomCopy
	}
	
	return rooms
}

// GetRoomConnections returns the connections for a specific room
func (wm *WorldMap) GetRoomConnections(roomID string) map[Direction]string {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	
	if connections, exists := wm.roomConnections[roomID]; exists {
		// Return a copy to prevent external modification
		connectionsCopy := make(map[Direction]string)
		for dir, connectedRoom := range connections {
			connectionsCopy[dir] = connectedRoom
		}
		return connectionsCopy
	}
	
	return make(map[Direction]string)
}

// AddPlayerPosition adds a player position to the movement trail
func (wm *WorldMap) AddPlayerPosition(x, y int) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	
	position := Point{X: x, Y: y}
	wm.playerTrail = append(wm.playerTrail, position)
	
	// Keep trail size limited
	if len(wm.playerTrail) > 100 {
		wm.playerTrail = wm.playerTrail[1:]
	}
}

// GetPlayerTrail returns the recent player movement trail
func (wm *WorldMap) GetPlayerTrail() []Point {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	
	// Return a copy
	trail := make([]Point, len(wm.playerTrail))
	copy(trail, wm.playerTrail)
	return trail
}

// GetMapBounds returns the bounding rectangle of all discovered rooms
func (wm *WorldMap) GetMapBounds() image.Rectangle {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	
	if len(wm.discoveredRooms) == 0 {
		return image.Rectangle{}
	}
	
	var minX, minY, maxX, maxY int
	first := true
	
	for _, room := range wm.discoveredRooms {
		roomLeft := room.WorldPos.X - room.Bounds.Dx()/2
		roomRight := room.WorldPos.X + room.Bounds.Dx()/2
		roomTop := room.WorldPos.Y - room.Bounds.Dy()/2
		roomBottom := room.WorldPos.Y + room.Bounds.Dy()/2
		
		if first {
			minX, maxX = roomLeft, roomRight
			minY, maxY = roomTop, roomBottom
			first = false
		} else {
			if roomLeft < minX {
				minX = roomLeft
			}
			if roomRight > maxX {
				maxX = roomRight
			}
			if roomTop < minY {
				minY = roomTop
			}
			if roomBottom > maxY {
				maxY = roomBottom
			}
		}
	}
	
	return image.Rectangle{
		Min: image.Point{X: minX, Y: minY},
		Max: image.Point{X: maxX, Y: maxY},
	}
}

// ToJSON serializes the world map to JSON
func (wm *WorldMap) ToJSON() ([]byte, error) {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	return json.Marshal(wm)
}

// FromJSON deserializes the world map from JSON
func (wm *WorldMap) FromJSON(data []byte) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	return json.Unmarshal(data, wm)
}

// worldMapJSON is an exported representation used for JSON (since WorldMap fields are unexported)
type worldMapJSON struct {
	DiscoveredRooms map[string]*DiscoveredRoom   `json:"discovered_rooms"`
	CurrentRoomID   string                       `json:"current_room_id"`
	RoomConnections map[string]map[string]string `json:"room_connections"` // direction as string key
	PlayerTrail     []Point                      `json:"player_trail"`
}

// MarshalJSON implements custom JSON marshaling for WorldMap
func (wm *WorldMap) MarshalJSON() ([]byte, error) {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	// Convert direction-keyed maps to string-keyed maps for JSON
	connections := make(map[string]map[string]string, len(wm.roomConnections))
	for roomID, dirMap := range wm.roomConnections {
		if dirMap == nil {
			continue
		}
		strMap := make(map[string]string, len(dirMap))
		for dir, to := range dirMap {
			strMap[dir.String()] = to
		}
		connections[roomID] = strMap
	}

	payload := worldMapJSON{
		DiscoveredRooms: wm.discoveredRooms,
		CurrentRoomID:   wm.currentRoomID,
		RoomConnections: connections,
		PlayerTrail:     wm.playerTrail,
	}
	return json.Marshal(payload)
}

// UnmarshalJSON implements custom JSON unmarshaling for WorldMap
func (wm *WorldMap) UnmarshalJSON(data []byte) error {
	var payload worldMapJSON
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	// Rebuild maps
	if wm.discoveredRooms == nil {
		wm.discoveredRooms = make(map[string]*DiscoveredRoom)
	}
	for k, v := range payload.DiscoveredRooms {
		wm.discoveredRooms[k] = v
	}

	wm.roomConnections = make(map[string]map[Direction]string)
	for roomID, strMap := range payload.RoomConnections {
		if strMap == nil {
			continue
		}
		dirMap := make(map[Direction]string, len(strMap))
		for dirStr, to := range strMap {
			dirMap[parseDirection(dirStr)] = to
		}
		wm.roomConnections[roomID] = dirMap
	}

	wm.currentRoomID = payload.CurrentRoomID
	wm.playerTrail = append([]Point(nil), payload.PlayerTrail...)
	return nil
}

// parseDirection converts a direction string to Direction enum
func parseDirection(s string) Direction {
	switch s {
	case "North":
		return North
	case "South":
		return South
	case "East":
		return East
	case "West":
		return West
	case "Up":
		return Up
	case "Down":
		return Down
	default:
		return East // sensible default
	}
}