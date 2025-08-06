package world

import (
	"fmt"
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/entities"
)

/*
ExitTrigger represents a trigger zone that initiates room transitions.
Can be placed on tiles or defined as rectangular areas in a room.
When the player enters the trigger area, it initiates a transition to another room.
*/
type ExitTrigger struct {
	ID           string  // Unique identifier for this exit
	X, Y         int     // Position in tile coordinates
	Width, Height int    // Size of the trigger area (default 1x1 for single tile)
	TargetRoomID string  // ID of the room to transition to
	TargetEntranceID string // ID of the entrance in the target room to spawn at
	
	// Optional transition data
	TransitionType string      // Type of transition (e.g., "fade", "slide", "instant")
	RequiredItems  []string    // Items/keys required to use this exit
	IsLocked       bool        // Whether this exit is currently locked
	Message        string      // Optional message to display when triggered/locked
}

/*
Entrance represents a spawn point in a room where players can enter.
Each entrance has a unique ID that can be referenced by exit triggers.
*/
type Entrance struct {
	ID           string  // Unique identifier for this entrance
	X, Y         int     // Position in tile coordinates where player spawns
	Direction    string  // Direction player should face on spawn ("left", "right", "up", "down")
	SpawnOffsetX int     // Fine-tuning spawn position (pixels from tile center)
	SpawnOffsetY int     // Fine-tuning spawn position (pixels from tile center)
}

/*
RoomTransitionData contains all information needed to perform a room transition.
This data is passed between the transition system and game state.
*/
type RoomTransitionData struct {
	SourceRoomID     string  // Room we're transitioning from
	TargetRoomID     string  // Room we're transitioning to
	TargetEntranceID string  // Entrance ID in target room
	ExitTriggerID    string  // ID of the exit trigger that initiated transition
	TransitionType   string  // How to perform the transition
	PlayerPosition   [2]int  // Player's current position when transition triggered
}

/*
RoomTransitionSystem manages room transitions and maintains room data.
Handles loading/unloading rooms, managing exit triggers, and coordinating transitions.
*/
type RoomTransitionSystem struct {
	currentRoomID string
	loadedRooms   map[string]Room                 // Cache of loaded rooms
	roomExits     map[string][]ExitTrigger        // Exit triggers per room
	roomEntrances map[string][]Entrance           // Entrances per room
	roomFactory   func(string) Room               // Function to create rooms by ID
}

/*
NewRoomTransitionSystem creates a new room transition system.
The roomFactory function should create and return a room instance for the given room ID.
*/
func NewRoomTransitionSystem(roomFactory func(string) Room) *RoomTransitionSystem {
	return &RoomTransitionSystem{
		loadedRooms:   make(map[string]Room),
		roomExits:     make(map[string][]ExitTrigger),
		roomEntrances: make(map[string][]Entrance),
		roomFactory:   roomFactory,
	}
}

/*
RegisterExitTrigger adds an exit trigger to a room.
*/
func (rts *RoomTransitionSystem) RegisterExitTrigger(roomID string, trigger ExitTrigger) {
	if rts.roomExits[roomID] == nil {
		rts.roomExits[roomID] = make([]ExitTrigger, 0)
	}
	rts.roomExits[roomID] = append(rts.roomExits[roomID], trigger)
}

/*
RegisterEntrance adds an entrance to a room.
*/
func (rts *RoomTransitionSystem) RegisterEntrance(roomID string, entrance Entrance) {
	if rts.roomEntrances[roomID] == nil {
		rts.roomEntrances[roomID] = make([]Entrance, 0)
	}
	rts.roomEntrances[roomID] = append(rts.roomEntrances[roomID], entrance)
}

/*
CheckExitTriggers checks if the player is touching any exit triggers in the current room.
Returns the transition data if a trigger is activated, nil otherwise.
*/
func (rts *RoomTransitionSystem) CheckExitTriggers(player *entities.Player, roomID string) *RoomTransitionData {
	playerX, playerY := player.GetPosition()
	physicsUnit := engine.GetPhysicsUnit()
	
	// Convert player position to tile coordinates
	playerTileX := playerX / physicsUnit
	playerTileY := playerY / physicsUnit
	
	exits := rts.roomExits[roomID]
	for _, exit := range exits {
		// Check if player is within the exit trigger area
		if playerTileX >= exit.X && playerTileX < exit.X+exit.Width &&
		   playerTileY >= exit.Y && playerTileY < exit.Y+exit.Height {
			
			// Check if exit is accessible
			if exit.IsLocked {
				// TODO: Show locked message or handle locked exits
				continue
			}
			
			// TODO: Check required items
			
			// Create transition data
			return &RoomTransitionData{
				SourceRoomID:     roomID,
				TargetRoomID:     exit.TargetRoomID,
				TargetEntranceID: exit.TargetEntranceID,
				ExitTriggerID:    exit.ID,
				TransitionType:   exit.TransitionType,
				PlayerPosition:   [2]int{playerX, playerY},
			}
		}
	}
	
	return nil
}

/*
GetEntrancePosition returns the spawn position for a specific entrance in a room.
Returns the position in physics units (pixels) where the player should spawn.
*/
func (rts *RoomTransitionSystem) GetEntrancePosition(roomID, entranceID string) (int, int, bool) {
	entrances := rts.roomEntrances[roomID]
	physicsUnit := engine.GetPhysicsUnit()
	
	for _, entrance := range entrances {
		if entrance.ID == entranceID {
			// Convert tile position to physics position and apply offset
			spawnX := entrance.X*physicsUnit + entrance.SpawnOffsetX
			spawnY := entrance.Y*physicsUnit + entrance.SpawnOffsetY
			return spawnX, spawnY, true
		}
	}
	
	return 0, 0, false
}

/*
LoadRoom loads or retrieves a room from the cache.
If the room is not cached, it uses the room factory to create it.
*/
func (rts *RoomTransitionSystem) LoadRoom(roomID string) Room {
	if room, exists := rts.loadedRooms[roomID]; exists {
		return room
	}
	
	// Create new room using factory
	room := rts.roomFactory(roomID)
	if room != nil {
		rts.loadedRooms[roomID] = room
	}
	
	return room
}

/*
UnloadRoom removes a room from the cache to free memory.
Useful for managing memory when dealing with many rooms.
*/
func (rts *RoomTransitionSystem) UnloadRoom(roomID string) {
	if room, exists := rts.loadedRooms[roomID]; exists {
		// Call OnExit if the room has cleanup to do
		if room != nil {
			// Room cleanup would go here if needed
		}
		delete(rts.loadedRooms, roomID)
	}
}

/*
GetCurrentRoomID returns the ID of the currently active room.
*/
func (rts *RoomTransitionSystem) GetCurrentRoomID() string {
	return rts.currentRoomID
}

/*
SetCurrentRoomID sets the current room ID.
*/
func (rts *RoomTransitionSystem) SetCurrentRoomID(roomID string) {
	rts.currentRoomID = roomID
}

/*
PerformTransition executes a room transition with the given transition data.
This is the main function that coordinates the entire transition process:
- Load the target room
- Position the player at the target entrance  
- Update camera bounds
- Call appropriate room callbacks
*/
func (rts *RoomTransitionSystem) PerformTransition(
	transitionData *RoomTransitionData,
	player *entities.Player,
	camera *engine.Camera,
) (Room, error) {
	// Load target room
	targetRoom := rts.LoadRoom(transitionData.TargetRoomID)
	if targetRoom == nil {
		engine.LogError(fmt.Sprintf("Failed to load room: %s", transitionData.TargetRoomID))
	return nil, fmt.Errorf("failed to load room: %s", transitionData.TargetRoomID)
	}
	
	// Get entrance position
	spawnX, spawnY, found := rts.GetEntrancePosition(transitionData.TargetRoomID, transitionData.TargetEntranceID)
	if !found {
		// Fallback to room's default spawn logic
		spawnX = targetRoom.FindFloorAtX(0)
		spawnY = targetRoom.FindFloorAtX(spawnX)
		engine.LogWarn("Entrance not found, using fallback spawn position")
	}
	
	// Position player at entrance
	player.SetPosition(spawnX, spawnY)
	
	// Update camera world bounds for new room
	tileMap := targetRoom.GetTileMap()
	if tileMap != nil {
		physicsUnit := engine.GetPhysicsUnit()
		worldWidth := tileMap.Width * physicsUnit
		worldHeight := tileMap.Height * physicsUnit
		camera.SetWorldBounds(worldWidth, worldHeight)
		
		// Immediately center camera on player (no smooth transition for room changes)
		camera.CenterOn(spawnX, spawnY)
	}
	
	// Call room transition callbacks
	targetRoom.OnEnter(player)
	
	// Update current room tracking
	rts.currentRoomID = transitionData.TargetRoomID
	
	engine.LogInfo(fmt.Sprintf("Room transition completed: %s -> %s (entrance: %s)", 
		transitionData.SourceRoomID, 
		transitionData.TargetRoomID,
		transitionData.TargetEntranceID))
	
	return targetRoom, nil
}

/*
GetRoomExits returns all exit triggers for a specific room.
Useful for debugging or displaying available exits.
*/
func (rts *RoomTransitionSystem) GetRoomExits(roomID string) []ExitTrigger {
	return rts.roomExits[roomID]
}

/*
GetRoomEntrances returns all entrances for a specific room.
Useful for debugging or displaying available entrances.
*/
func (rts *RoomTransitionSystem) GetRoomEntrances(roomID string) []Entrance {
	return rts.roomEntrances[roomID]
}

/*
DrawExitTriggers renders debug visualization of exit triggers in a room.
Draws colored rectangles over exit trigger areas for debugging purposes.
*/
func (rts *RoomTransitionSystem) DrawExitTriggers(screen *ebiten.Image, roomID string, cameraOffsetX, cameraOffsetY float64) {
	if !engine.GameConfig.ShowDebugOverlay {
		return
	}
	
	exits := rts.roomExits[roomID]
	if len(exits) == 0 {
		return
	}
	
	physicsUnit := engine.GetPhysicsUnit()
	
	for _, exit := range exits {
		// Calculate screen position
		x := float64(exit.X * physicsUnit) + cameraOffsetX
		y := float64(exit.Y * physicsUnit) + cameraOffsetY
		width := float64(exit.Width * physicsUnit)
		height := float64(exit.Height * physicsUnit)
		
		// Create a semi-transparent overlay image
		overlayImg := ebiten.NewImage(int(width), int(height))
		
		// Choose color based on exit status
		var fillColor color.RGBA
		if exit.IsLocked {
			fillColor = color.RGBA{255, 0, 0, 100} // Red for locked exits
		} else {
			fillColor = color.RGBA{0, 255, 0, 100} // Green for accessible exits
		}
		
		// Fill with color
		overlayImg.Fill(fillColor)
		
		// Draw the overlay
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)
		screen.DrawImage(overlayImg, opts)
	}
}

/*
DrawEntrances renders debug visualization of entrance points in a room.
Draws markers at entrance positions for debugging purposes.
*/
func (rts *RoomTransitionSystem) DrawEntrances(screen *ebiten.Image, roomID string, cameraOffsetX, cameraOffsetY float64) {
	if !engine.GameConfig.ShowDebugOverlay {
		return
	}
	
	entrances := rts.roomEntrances[roomID]
	if len(entrances) == 0 {
		return
	}
	
	physicsUnit := engine.GetPhysicsUnit()
	
	for _, entrance := range entrances {
		// Calculate screen position
		x := float64(entrance.X * physicsUnit) + cameraOffsetX
		y := float64(entrance.Y * physicsUnit) + cameraOffsetY
		
		// Create a small marker image (cross or diamond shape)
		markerSize := 16
		markerImg := ebiten.NewImage(markerSize, markerSize)
		
		// Fill with blue color for entrance markers
		markerImg.Fill(color.RGBA{0, 0, 255, 150})
		
		// Draw the marker centered on the entrance position
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x - float64(markerSize/2), y - float64(markerSize/2))
		screen.DrawImage(markerImg, opts)
	}
}