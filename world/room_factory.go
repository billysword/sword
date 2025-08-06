package world

import (
	"fmt"
	"sword/engine"
)

/*
RoomFactory provides a centralized way to create rooms by ID.
This allows the transition system to dynamically load rooms
as needed without having to know the specific room types.
*/
type RoomFactory struct {
	// Future: could include room loading configuration, asset paths, etc.
}

/*
NewRoomFactory creates a new room factory instance.
*/
func NewRoomFactory() *RoomFactory {
	return &RoomFactory{}
}

/*
CreateRoom creates a room instance for the given room ID.
This function handles the mapping from room IDs to specific room implementations.
Add new room types here as they are developed.
*/
func (rf *RoomFactory) CreateRoom(roomID string) Room {
	switch roomID {
	case "main":
		return NewSimpleRoom("main")
	case "forest_entrance":
		return NewForestEntranceRoom()
	case "underground_cave":
		return NewUndergroundCaveRoom()
	case "village_square":
		return NewVillageSquareRoom()
	default:
		engine.LogWarn(fmt.Sprintf("Unknown room ID: %s, creating default room", roomID))
		return NewSimpleRoom(roomID)
	}
}

/*
SetupDefaultRoomConnections configures the standard connections between rooms.
This sets up exit triggers and entrances for the default rooms.
Call this during game initialization to establish room connections.
*/
func (rf *RoomFactory) SetupDefaultRoomConnections(transitionSystem *RoomTransitionSystem) {
	// Main room connections
	transitionSystem.RegisterExitTrigger("main", ExitTrigger{
		ID:               "main_to_forest",
		X:                9,  // Right edge of room
		Y:                8,  // Near the bottom
		Width:            1,
		Height:           1,
		TargetRoomID:     "forest_entrance",
		TargetEntranceID: "from_main",
		TransitionType:   "instant",
	})
	
	transitionSystem.RegisterEntrance("main", Entrance{
		ID:        "main_spawn",
		X:         4,  // Center of room
		Y:         8,  // Near bottom
		Direction: "right",
	})
	
	transitionSystem.RegisterEntrance("main", Entrance{
		ID:        "from_forest",
		X:         1,  // Left side
		Y:         8,  // Near bottom
		Direction: "right",
	})
	
	// Forest entrance room connections
	transitionSystem.RegisterExitTrigger("forest_entrance", ExitTrigger{
		ID:               "forest_to_main",
		X:                0,  // Left edge
		Y:                8,
		Width:            1,
		Height:           1,
		TargetRoomID:     "main",
		TargetEntranceID: "from_forest",
		TransitionType:   "instant",
	})
	
	transitionSystem.RegisterExitTrigger("forest_entrance", ExitTrigger{
		ID:               "forest_to_cave",
		X:                4,  // Center bottom
		Y:                9,  // Bottom row
		Width:            2,  // 2-tile wide entrance
		Height:           1,
		TargetRoomID:     "underground_cave",
		TargetEntranceID: "from_forest",
		TransitionType:   "instant",
	})
	
	transitionSystem.RegisterEntrance("forest_entrance", Entrance{
		ID:        "from_main",
		X:         9,  // Right side
		Y:         8,
		Direction: "left",
	})
	
	transitionSystem.RegisterEntrance("forest_entrance", Entrance{
		ID:        "from_cave",
		X:         5,  // Center
		Y:         8,
		Direction: "up",
	})
	
	// Underground cave connections
	transitionSystem.RegisterExitTrigger("underground_cave", ExitTrigger{
		ID:               "cave_to_forest",
		X:                4,  // Center
		Y:                0,  // Top of room
		Width:            2,
		Height:           1,
		TargetRoomID:     "forest_entrance",
		TargetEntranceID: "from_cave",
		TransitionType:   "instant",
	})
	
	transitionSystem.RegisterEntrance("underground_cave", Entrance{
		ID:        "from_forest",
		X:         5,  // Center
		Y:         1,  // Near top
		Direction: "down",
	})
	
	// Village square connections (for future expansion)
	transitionSystem.RegisterEntrance("village_square", Entrance{
		ID:        "main_entrance",
		X:         5,
		Y:         8,
		Direction: "up",
	})
	
	engine.LogInfo("Default room connections configured")
}

/*
NewForestEntranceRoom creates a forest-themed entrance room.
This is a placeholder implementation - could be expanded with custom layouts.
*/
func NewForestEntranceRoom() Room {
	return NewSimpleRoom("forest_entrance")
}

/*
NewUndergroundCaveRoom creates an underground cave room.
This is a placeholder implementation - could be expanded with custom layouts.
*/
func NewUndergroundCaveRoom() Room {
	return NewSimpleRoom("underground_cave")
}

/*
NewVillageSquareRoom creates a village square room.
This is a placeholder implementation - could be expanded with custom layouts.
*/
func NewVillageSquareRoom() Room {
	return NewSimpleRoom("village_square")
}

/*
GetRoomFactory returns a factory function that can be used by the transition system.
This wraps the room factory in a simple function signature.
*/
func GetRoomFactory() func(string) Room {
	factory := NewRoomFactory()
	return factory.CreateRoom
}