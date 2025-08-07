package world

import (
	"fmt"
	"sword/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
)

// RoomConfig contains configuration for room creation
type RoomConfig struct {
	ZoneID     string
	Width      int
	Height     int
	Theme      string
	Difficulty int
	
	// Transition configuration
	EnableTransitions bool
	TransitionPoints  []TransitionPoint
	SpawnPoints       []SpawnPoint
	
	// Rendering configuration
	BackgroundEnabled bool
	ParallaxEnabled   bool
	LightingEnabled   bool
}

// RoomMetadata contains additional information about a room
type RoomMetadata struct {
	DisplayName   string
	Description   string
	Tags          []string
	RequiredItems []string
	Secrets       int
	EnemyCount    int
}

// Enhanced Room interface with better transition and connection support
type EnhancedRoom interface {
	// Core functionality (from original Room interface)
	GetTileMap() *TileMap
	GetZoneID() string
	Update(player *entities.Player) error
	HandleCollisions(player *entities.Player)
	OnEnter(player *entities.Player)
	OnExit(player *entities.Player)
	FindFloorAtX(x int) int
	Draw(screen *ebiten.Image)
	DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64)
	DrawTiles(screen *ebiten.Image, spriteProvider func(int) *ebiten.Image)

	// Enhanced functionality for better architecture
	GetMetadata() RoomMetadata
	SetMetadata(metadata RoomMetadata)
	
	// Transition point management
	GetTransitionPoints() []TransitionPoint
	AddTransitionPoint(point TransitionPoint) error
	RemoveTransitionPoint(targetRoomID string) bool
	EnableTransition(targetRoomID string, enabled bool)
	
	// Spawn point management
	GetSpawnPoints() []SpawnPoint
	AddSpawnPoint(point SpawnPoint) error
	RemoveSpawnPoint(spawnID string) bool
	GetSpawnPoint(spawnID string) (*SpawnPoint, bool)
	
	// Room connections and neighboring
	GetConnectedRooms() []string
	IsConnectedTo(roomID string) bool
	
	// Environmental and gameplay features
	GetAmbientSound() string
	SetAmbientSound(soundID string)
	GetLightLevel() float32
	SetLightLevel(level float32)
	
	// Room state management
	Save() ([]byte, error)
	Load(data []byte) error
	Reset() error
	
	// Validation and integrity
	Validate() error
	GetChecksum() string
}

// EnhancedBaseRoom provides an improved base implementation
type EnhancedBaseRoom struct {
	*BaseRoom
	metadata         RoomMetadata
	transitionPoints []TransitionPoint
	spawnPoints      []SpawnPoint
	ambientSound     string
	lightLevel       float32
	lastChecksum     string
}

// NewEnhancedBaseRoom creates a new enhanced base room
func NewEnhancedBaseRoom(config RoomConfig) *EnhancedBaseRoom {
	baseRoom := NewBaseRoom(config.ZoneID, config.Width, config.Height)
	
	enhanced := &EnhancedBaseRoom{
		BaseRoom:         baseRoom,
		metadata:         RoomMetadata{DisplayName: config.ZoneID},
		transitionPoints: make([]TransitionPoint, 0),
		spawnPoints:      make([]SpawnPoint, 0),
		lightLevel:       1.0,
	}
	
	// Copy provided transition and spawn points
	if config.EnableTransitions {
		enhanced.transitionPoints = append(enhanced.transitionPoints, config.TransitionPoints...)
		enhanced.spawnPoints = append(enhanced.spawnPoints, config.SpawnPoints...)
	}
	
	return enhanced
}

// GetMetadata returns the room metadata
func (ebr *EnhancedBaseRoom) GetMetadata() RoomMetadata {
	return ebr.metadata
}

// SetMetadata sets the room metadata
func (ebr *EnhancedBaseRoom) SetMetadata(metadata RoomMetadata) {
	ebr.metadata = metadata
}

// GetTransitionPoints returns all transition points in the room
func (ebr *EnhancedBaseRoom) GetTransitionPoints() []TransitionPoint {
	return ebr.transitionPoints
}

// AddTransitionPoint adds a new transition point to the room
func (ebr *EnhancedBaseRoom) AddTransitionPoint(point TransitionPoint) error {
	// Validate the transition point
	if point.TargetRoomID == "" {
		return fmt.Errorf("transition point must have a target room ID")
	}
	if point.TargetSpawnID == "" {
		return fmt.Errorf("transition point must have a target spawn ID")
	}
	
	// Check for duplicate transitions to the same room
	for _, existing := range ebr.transitionPoints {
		if existing.TargetRoomID == point.TargetRoomID {
			return fmt.Errorf("transition to room %s already exists", point.TargetRoomID)
		}
	}
	
	ebr.transitionPoints = append(ebr.transitionPoints, point)
	return nil
}

// RemoveTransitionPoint removes a transition point by target room ID
func (ebr *EnhancedBaseRoom) RemoveTransitionPoint(targetRoomID string) bool {
	for i, point := range ebr.transitionPoints {
		if point.TargetRoomID == targetRoomID {
			// Remove by swapping with last element
			ebr.transitionPoints[i] = ebr.transitionPoints[len(ebr.transitionPoints)-1]
			ebr.transitionPoints = ebr.transitionPoints[:len(ebr.transitionPoints)-1]
			return true
		}
	}
	return false
}

// EnableTransition enables or disables a specific transition
func (ebr *EnhancedBaseRoom) EnableTransition(targetRoomID string, enabled bool) {
	for i := range ebr.transitionPoints {
		if ebr.transitionPoints[i].TargetRoomID == targetRoomID {
			ebr.transitionPoints[i].IsEnabled = enabled
			break
		}
	}
}

// GetSpawnPoints returns all spawn points in the room
func (ebr *EnhancedBaseRoom) GetSpawnPoints() []SpawnPoint {
	return ebr.spawnPoints
}

// AddSpawnPoint adds a new spawn point to the room
func (ebr *EnhancedBaseRoom) AddSpawnPoint(point SpawnPoint) error {
	// Validate the spawn point
	if point.ID == "" {
		return fmt.Errorf("spawn point must have an ID")
	}
	
	// Check for duplicate spawn IDs
	for _, existing := range ebr.spawnPoints {
		if existing.ID == point.ID {
			return fmt.Errorf("spawn point with ID %s already exists", point.ID)
		}
	}
	
	ebr.spawnPoints = append(ebr.spawnPoints, point)
	return nil
}

// RemoveSpawnPoint removes a spawn point by ID
func (ebr *EnhancedBaseRoom) RemoveSpawnPoint(spawnID string) bool {
	for i, point := range ebr.spawnPoints {
		if point.ID == spawnID {
			// Remove by swapping with last element
			ebr.spawnPoints[i] = ebr.spawnPoints[len(ebr.spawnPoints)-1]
			ebr.spawnPoints = ebr.spawnPoints[:len(ebr.spawnPoints)-1]
			return true
		}
	}
	return false
}

// GetSpawnPoint returns a specific spawn point by ID
func (ebr *EnhancedBaseRoom) GetSpawnPoint(spawnID string) (*SpawnPoint, bool) {
	for _, point := range ebr.spawnPoints {
		if point.ID == spawnID {
			return &point, true
		}
	}
	return nil, false
}

// GetConnectedRooms returns IDs of all connected rooms
func (ebr *EnhancedBaseRoom) GetConnectedRooms() []string {
	connected := make([]string, 0, len(ebr.transitionPoints))
	for _, point := range ebr.transitionPoints {
		if point.IsEnabled {
			connected = append(connected, point.TargetRoomID)
		}
	}
	return connected
}

// IsConnectedTo checks if this room is connected to another room
func (ebr *EnhancedBaseRoom) IsConnectedTo(roomID string) bool {
	for _, point := range ebr.transitionPoints {
		if point.TargetRoomID == roomID && point.IsEnabled {
			return true
		}
	}
	return false
}

// GetAmbientSound returns the ambient sound ID
func (ebr *EnhancedBaseRoom) GetAmbientSound() string {
	return ebr.ambientSound
}

// SetAmbientSound sets the ambient sound for the room
func (ebr *EnhancedBaseRoom) SetAmbientSound(soundID string) {
	ebr.ambientSound = soundID
}

// GetLightLevel returns the current light level (0.0 = dark, 1.0 = bright)
func (ebr *EnhancedBaseRoom) GetLightLevel() float32 {
	return ebr.lightLevel
}

// SetLightLevel sets the light level for the room
func (ebr *EnhancedBaseRoom) SetLightLevel(level float32) {
	if level < 0.0 {
		level = 0.0
	}
	if level > 1.0 {
		level = 1.0
	}
	ebr.lightLevel = level
}

// Save serializes the room state (stub implementation)
func (ebr *EnhancedBaseRoom) Save() ([]byte, error) {
	// TODO: Implement room state serialization
	return []byte{}, nil
}

// Load deserializes room state (stub implementation)
func (ebr *EnhancedBaseRoom) Load(data []byte) error {
	// TODO: Implement room state deserialization
	return nil
}

// Reset restores the room to its initial state
func (ebr *EnhancedBaseRoom) Reset() error {
	// Reset room-specific state
	// This would reset enemy positions, item spawns, etc.
	return nil
}

// Validate checks room integrity and returns any issues
func (ebr *EnhancedBaseRoom) Validate() error {
	// Check tile map validity
	if ebr.tileMap == nil {
		return fmt.Errorf("room has no tile map")
	}
	
	// Validate transition points
	for i, point := range ebr.transitionPoints {
		if point.TargetRoomID == "" {
			return fmt.Errorf("transition point %d has empty target room ID", i)
		}
		if point.TargetSpawnID == "" {
			return fmt.Errorf("transition point %d has empty target spawn ID", i)
		}
		// Validate trigger bounds are within room bounds
		tileMap := ebr.GetTileMap()
		if tileMap != nil {
			u := engine.GetPhysicsUnit()
			maxX := tileMap.Width * u
			maxY := tileMap.Height * u
			
			if point.TriggerBounds.X < 0 || point.TriggerBounds.Y < 0 ||
			   point.TriggerBounds.X + point.TriggerBounds.Width > maxX ||
			   point.TriggerBounds.Y + point.TriggerBounds.Height > maxY {
				return fmt.Errorf("transition point %d has bounds outside room area", i)
			}
		}
	}
	
	// Validate spawn points
	for i, point := range ebr.spawnPoints {
		if point.ID == "" {
			return fmt.Errorf("spawn point %d has empty ID", i)
		}
		// Check for duplicate IDs
		for j, other := range ebr.spawnPoints {
			if i != j && point.ID == other.ID {
				return fmt.Errorf("duplicate spawn point ID: %s", point.ID)
			}
		}
	}
	
	return nil
}

// GetChecksum returns a checksum of the room state for change detection
func (ebr *EnhancedBaseRoom) GetChecksum() string {
	// TODO: Implement proper checksum calculation
	// This would hash the tile map, transition points, spawn points, etc.
	return "placeholder_checksum"
}

// RoomFactory creates rooms with proper configuration
type RoomFactory struct {
	defaultConfig RoomConfig
}

// NewRoomFactory creates a new room factory with default configuration
func NewRoomFactory() *RoomFactory {
	return &RoomFactory{
		defaultConfig: RoomConfig{
			Width:             20,
			Height:            15,
			Theme:             "forest",
			EnableTransitions: true,
			BackgroundEnabled: true,
			ParallaxEnabled:   true,
		},
	}
}

// CreateRoom creates a new room with the specified configuration
func (rf *RoomFactory) CreateRoom(config RoomConfig) EnhancedRoom {
	// Merge with defaults
	if config.Width == 0 {
		config.Width = rf.defaultConfig.Width
	}
	if config.Height == 0 {
		config.Height = rf.defaultConfig.Height
	}
	if config.Theme == "" {
		config.Theme = rf.defaultConfig.Theme
	}
	
	// Create enhanced base room
	room := NewEnhancedBaseRoom(config)
	
	// TODO: Apply theme-specific configuration
	// TODO: Generate procedural content if needed
	
	return room
}

// CreateConnectedRooms creates multiple rooms with automatic connections
func (rf *RoomFactory) CreateConnectedRooms(configs []RoomConfig) ([]EnhancedRoom, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no room configurations provided")
	}
	
	rooms := make([]EnhancedRoom, len(configs))
	
	// Create all rooms first
	for i, config := range configs {
		rooms[i] = rf.CreateRoom(config)
	}
	
	// Create connections based on room layout
	// This is a simple linear connection for now
	for i := 0; i < len(rooms)-1; i++ {
		currentRoom := rooms[i]
		nextRoom := rooms[i+1]
		
		// Add transition from current to next
		transition := TransitionPoint{
			Type:          TransitionWalk,
			TargetRoomID:  nextRoom.GetZoneID(),
			TargetSpawnID: "main_spawn",
			TriggerBounds: Rectangle{X: 0, Y: 0, Width: 32, Height: 32}, // TODO: calculate proper bounds
			IsEnabled:     true,
		}
		currentRoom.AddTransitionPoint(transition)
		
		// Add spawn point to next room if it doesn't have one
		if len(nextRoom.GetSpawnPoints()) == 0 {
			spawn := SpawnPoint{
				ID: "main_spawn",
				X:  100, // TODO: calculate proper position
				Y:  100,
			}
			nextRoom.AddSpawnPoint(spawn)
		}
	}
	
	return rooms, nil
}