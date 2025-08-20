package world

import (
	"fmt"
	"strings"
	"sword/engine"
	"sword/entities"
)

// TransitionType defines how players move between rooms
type TransitionType int

const (
	TransitionWalk     TransitionType = iota // Player walks to edge and transitions
	TransitionDoor                           // Player interacts with door/portal
	TransitionTeleport                       // Instant teleportation
	TransitionStairs                         // Vertical movement between levels
)

// String returns the string representation of a transition type
func (tt TransitionType) String() string {
	switch tt {
	case TransitionWalk:
		return "Walk"
	case TransitionDoor:
		return "Door"
	case TransitionTeleport:
		return "Teleport"
	case TransitionStairs:
		return "Stairs"
	default:
		return "Unknown"
	}
}

// TransitionPoint represents a connection point between rooms
type TransitionPoint struct {
	Type           TransitionType `json:"type"`            // How the transition works
	TriggerBounds  Rectangle      `json:"trigger_bounds"`  // Area that triggers the transition
	TargetRoomID   string         `json:"target_room_id"`  // ID of the room to transition to
	TargetSpawnID  string         `json:"target_spawn_id"` // Spawn point in target room
	RequiresAction bool           `json:"requires_action"` // Does player need to press a key?
	IsEnabled      bool           `json:"is_enabled"`      // Can this transition be used?
	Direction      Direction      `json:"direction"`       // Which direction triggers this
}

// Rectangle represents a rectangular area
type Rectangle struct {
	X, Y, Width, Height int
}

// Contains checks if a point is within the rectangle
func (r Rectangle) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

// Intersects checks if this rectangle overlaps another rectangle (AABB)
func (r Rectangle) Intersects(o Rectangle) bool {
	return r.X < o.X+o.Width && r.X+r.Width > o.X && r.Y < o.Y+o.Height && r.Y+r.Height > o.Y
}

// SpawnPoint represents a place where players can appear in a room
type SpawnPoint struct {
	ID       string `json:"id"`        // Unique identifier for this spawn point
	X        int    `json:"x"`         // Position X in physics units
	Y        int    `json:"y"`         // Position Y in physics units
	FacingID string `json:"facing_id"` // Direction player should face when spawning
}

// RoomTransitionManager handles room transitions and connections
type RoomTransitionManager struct {
	currentRoomID     string
	rooms             map[string]Room
	transitionPoints  map[string][]TransitionPoint // room_id -> transition points
	spawnPoints       map[string][]SpawnPoint      // room_id -> spawn points
	pendingTransition *PendingTransition           // Queued room transition
	worldMap          *WorldMap                    // Reference to world map for auto connections
}

// PendingTransition represents a room change that will happen at the end of the frame
type PendingTransition struct {
	TargetRoomID   string
	TargetSpawnID  string
	TransitionType TransitionType
}

// NewRoomTransitionManager creates a new room transition manager
func NewRoomTransitionManager(worldMap *WorldMap) *RoomTransitionManager {
	return &RoomTransitionManager{
		rooms:            make(map[string]Room),
		transitionPoints: make(map[string][]TransitionPoint),
		spawnPoints:      make(map[string][]SpawnPoint),
		worldMap:         worldMap,
	}
}

// RegisterRoom adds a room to the transition system
func (rtm *RoomTransitionManager) RegisterRoom(room Room) {
	rtm.rooms[room.GetZoneID()] = room

	// Initialize empty transition and spawn point slices if they don't exist
	if _, exists := rtm.transitionPoints[room.GetZoneID()]; !exists {
		rtm.transitionPoints[room.GetZoneID()] = make([]TransitionPoint, 0)
	}
	if _, exists := rtm.spawnPoints[room.GetZoneID()]; !exists {
		rtm.spawnPoints[room.GetZoneID()] = make([]SpawnPoint, 0)
	}

	if rtm.worldMap != nil {
		rtm.worldMap.DiscoverRoom(room)
	}
}

// SetCurrentRoom sets the active room
func (rtm *RoomTransitionManager) SetCurrentRoom(roomID string) error {
	if _, exists := rtm.rooms[roomID]; !exists {
		return fmt.Errorf("room %s not found", roomID)
	}
	rtm.currentRoomID = roomID
	return nil
}

// GetCurrentRoom returns the current room
func (rtm *RoomTransitionManager) GetCurrentRoom() Room {
	if room, exists := rtm.rooms[rtm.currentRoomID]; exists {
		return room
	}
	return nil
}

// GetCurrentRoomID returns the current room ID
func (rtm *RoomTransitionManager) GetCurrentRoomID() string {
	return rtm.currentRoomID
}

// AddTransitionPoint adds a transition point to a room
func (rtm *RoomTransitionManager) AddTransitionPoint(roomID string, transition TransitionPoint) error {
	if _, exists := rtm.rooms[roomID]; !exists {
		return fmt.Errorf("room %s not found", roomID)
	}

	rtm.transitionPoints[roomID] = append(rtm.transitionPoints[roomID], transition)
	engine.LogInfo(fmt.Sprintf("Added %s transition from %s to %s",
		transition.Type.String(), roomID, transition.TargetRoomID))

	if rtm.worldMap != nil {
		if err := rtm.worldMap.ConnectRooms(roomID, transition.Direction, transition.TargetRoomID); err != nil {
			engine.LogInfo(fmt.Sprintf("Failed to connect rooms: %v", err))
		}
	}
	return nil
}

// AddSpawnPoint adds a spawn point to a room
func (rtm *RoomTransitionManager) AddSpawnPoint(roomID string, spawn SpawnPoint) error {
	if _, exists := rtm.rooms[roomID]; !exists {
		return fmt.Errorf("room %s not found", roomID)
	}

	rtm.spawnPoints[roomID] = append(rtm.spawnPoints[roomID], spawn)
	engine.LogInfo(fmt.Sprintf("Added spawn point %s to room %s at (%d, %d)",
		spawn.ID, roomID, spawn.X, spawn.Y))
	return nil
}

// CheckTransitions checks if the player should transition to another room
func (rtm *RoomTransitionManager) CheckTransitions(player *entities.Player, actionPressed bool) bool {
	if rtm.currentRoomID == "" {
		return false
	}

	// Build player's collision rectangle for robust trigger detection
	box := player.GetCollisionBox()
	playerRect := Rectangle{X: box.X, Y: box.Y, Width: box.Width, Height: box.Height}

	transitions := rtm.transitionPoints[rtm.currentRoomID]

	for _, transition := range transitions {
		if !transition.IsEnabled {
			continue
		}

		// Check rectangle overlap instead of point; works even if player's origin isn't inside trigger
		if transition.TriggerBounds.Intersects(playerRect) {
			// Check if action is required and was pressed
			if transition.RequiresAction && !actionPressed {
				continue
			}

			// Queue the transition
			rtm.pendingTransition = &PendingTransition{
				TargetRoomID:   transition.TargetRoomID,
				TargetSpawnID:  transition.TargetSpawnID,
				TransitionType: transition.Type,
			}

			engine.LogInfo(fmt.Sprintf("Triggered %s transition from %s to %s",
				transition.Type.String(), rtm.currentRoomID, transition.TargetRoomID))
			return true
		}
	}

	return false
}

// ProcessPendingTransition executes any queued room transition
func (rtm *RoomTransitionManager) ProcessPendingTransition(player *entities.Player) (Room, error) {
	if rtm.pendingTransition == nil {
		return nil, nil
	}

	transition := rtm.pendingTransition
	rtm.pendingTransition = nil // Clear the pending transition

	// Validate target room exists
	targetRoom, exists := rtm.rooms[transition.TargetRoomID]
	if !exists {
		return nil, fmt.Errorf("target room %s not found", transition.TargetRoomID)
	}

	// Exit current room
	if currentRoom := rtm.GetCurrentRoom(); currentRoom != nil {
		currentRoom.OnExit(player)
	}

	// Set new current room
	rtm.currentRoomID = transition.TargetRoomID
	if rtm.worldMap != nil {
		rtm.worldMap.SetCurrentRoom(transition.TargetRoomID)
	}

	// Position player at spawn point
	if err := rtm.SpawnPlayerInRoom(player, transition.TargetRoomID, transition.TargetSpawnID); err != nil {
		engine.LogInfo(fmt.Sprintf("Warning: spawn positioning failed: %v", err))
		// Continue with transition anyway, but player might be at wrong position
	}

	// Enter new room
	targetRoom.OnEnter(player)

	engine.LogInfo(fmt.Sprintf("Completed room transition to %s", transition.TargetRoomID))
	return targetRoom, nil
}

// SpawnPlayerInRoom positions the player at a specific spawn point
func (rtm *RoomTransitionManager) SpawnPlayerInRoom(player *entities.Player, roomID, spawnID string) error {
	spawnPoints := rtm.spawnPoints[roomID]

	// Find the specific spawn point
	for _, spawn := range spawnPoints {
		if spawn.ID == spawnID {
			x, y := rtm.findSafeSpawnPosition(player, roomID, spawn.X, spawn.Y)
			player.SetPosition(x, y)
			// Set player facing direction based on spawn.FacingID (optional)
			switch strings.ToLower(spawn.FacingID) {
			case "east", "right":
				player.SetFacingRight(true)
			case "west", "left":
				player.SetFacingRight(false)
			}
			engine.LogInfo(fmt.Sprintf("Spawned player at %s in room %s", spawnID, roomID))
			return nil
		}
	}

	// Fallback: use first spawn point if specific one not found
	if len(spawnPoints) > 0 {
		spawn := spawnPoints[0]
		x, y := rtm.findSafeSpawnPosition(player, roomID, spawn.X, spawn.Y)
		player.SetPosition(x, y)
		engine.LogInfo(fmt.Sprintf("Spawned player at fallback spawn %s in room %s", spawn.ID, roomID))
		return nil
	}

	// Ultimate fallback: center of room
	if room, exists := rtm.rooms[roomID]; exists {
		if tileMap := room.GetTileMap(); tileMap != nil {
			u := engine.GetPhysicsUnit()
			centerX := (tileMap.Width / 2) * u
			centerY := (tileMap.Height / 2) * u
			centerX, centerY = rtm.findSafeSpawnPosition(player, roomID, centerX, centerY)
			player.SetPosition(centerX, centerY)
			engine.LogInfo(fmt.Sprintf("Spawned player at room center (%d, %d) in %s", centerX, centerY, roomID))
			return nil
		}
	}

	return fmt.Errorf("could not determine spawn location for room %s", roomID)
}

// findSafeSpawnPosition adjusts a position to align the player's collision box with the floor
// and avoid spawning overlapping solid tiles. It prefers placing the player exactly on the
// floor at the given X, then searches nearby positions if necessary. Returns original if
// a better position cannot be found.
func (rtm *RoomTransitionManager) findSafeSpawnPosition(player *entities.Player, roomID string, startX, startY int) (int, int) {
	room, exists := rtm.rooms[roomID]
	if !exists || room == nil {
		return startX, startY
	}
	// Tile map required for collision checks
	tileMap := room.GetTileMap()
	if tileMap == nil {
		return startX, startY
	}

	u := engine.GetPhysicsUnit()

	// Helper to compute Y so that the player's collision box sits on top of floor at X
	computeYOnFloor := func(testX int) (int, bool) {
		floorY := room.FindFloorAtX(testX)
		if floorY <= 0 {
			return 0, false
		}
		cfg := &engine.GameConfig.PlayerPhysics
		// Collision box dimensions in physics units (char scale applies to sprite dims)
		spriteH := int(float64(cfg.SpriteHeight) * engine.GameConfig.CharScaleFactor)
		offsetY := int(float64(spriteH) * cfg.CollisionBoxOffsetY)
		boxH := int(float64(spriteH) * cfg.CollisionBoxHeight)
		// Position top-left of sprite such that bottom of collision box is exactly at floor
		newY := floorY - (offsetY + boxH)
		return newY, true
	}

	// Prefer aligning to floor at the requested X
	candidateX := startX
	candidateY, ok := computeYOnFloor(candidateX)
	if !ok {
		candidateY = startY
	}

	// Check collision at the candidate position using player's collision box
	if tp, okTP := any(room).(entities.TileProvider); okTP {
		if !player.CheckTileCollision(tp, candidateX, candidateY) {
			return candidateX, candidateY
		}
	}

	// Search nearby X positions (within +-2 tiles) and align to floor at each
	offsets := []int{0, 1 * u, -1 * u, 2 * u, -2 * u}
	for _, dx := range offsets {
		testX := startX + dx
		// Clamp within room bounds
		if testX < 0 {
			testX = 0
		}
		maxX := tileMap.Width*u - 1
		if testX > maxX {
			testX = maxX
		}

		testY, ok := computeYOnFloor(testX)
		if !ok {
			continue
		}
		if tp, okTP := any(room).(entities.TileProvider); okTP {
			if !player.CheckTileCollision(tp, testX, testY) {
				return testX, testY
			}
		}
	}

	// Fallback: use simple non-solid tile search near the original point
	fx, fy := rtm.findNonSolidPosition(roomID, startX, startY)
	return fx, fy
}

// findNonSolidPosition adjusts a position to a nearby non-solid tile if necessary.
// It checks the target tile and adjacent tiles to avoid spawning inside solid blocks.
// Returns the original coordinates if no better position is found.
func (rtm *RoomTransitionManager) findNonSolidPosition(roomID string, x, y int) (int, int) {
	room, exists := rtm.rooms[roomID]
	if !exists {
		return x, y
	}
	tileMap := room.GetTileMap()
	if tileMap == nil {
		return x, y
	}

	u := engine.GetPhysicsUnit()
	tileX := x / u
	tileY := y / u
	if !IsSolidTile(tileMap.GetTileIndex(tileX, tileY)) {
		return x, y
	}

	offsets := [][2]int{{1, 0}, {-1, 0}, {0, -1}, {0, 1}}
	for _, off := range offsets {
		nx := tileX + off[0]
		ny := tileY + off[1]
		if !IsSolidTile(tileMap.GetTileIndex(nx, ny)) {
			return nx * u, ny * u
		}
	}
	return x, y
}

// GetTransitionPoints returns all transition points for a room
func (rtm *RoomTransitionManager) GetTransitionPoints(roomID string) []TransitionPoint {
	return rtm.transitionPoints[roomID]
}

// GetSpawnPoints returns all spawn points for a room
func (rtm *RoomTransitionManager) GetSpawnPoints(roomID string) []SpawnPoint {
	return rtm.spawnPoints[roomID]
}

// EnableTransition enables or disables a specific transition
func (rtm *RoomTransitionManager) EnableTransition(roomID string, targetRoomID string, enabled bool) {
	transitions := rtm.transitionPoints[roomID]
	for i := range transitions {
		if transitions[i].TargetRoomID == targetRoomID {
			transitions[i].IsEnabled = enabled
			break
		}
	}
}

// HasPendingTransition returns true if there's a room transition queued
func (rtm *RoomTransitionManager) HasPendingTransition() bool {
	return rtm.pendingTransition != nil
}

func (rtm *RoomTransitionManager) GetRoom(roomID string) Room {
	if r, ok := rtm.rooms[roomID]; ok {
		return r
	}
	return nil
}

// ListRoomIDs returns all registered room IDs in undefined order
func (rtm *RoomTransitionManager) ListRoomIDs() []string {
	ids := make([]string, 0, len(rtm.rooms))
	for id := range rtm.rooms {
		ids = append(ids, id)
	}
	return ids
}
