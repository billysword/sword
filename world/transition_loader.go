package world

import (
	"encoding/json"
	"os"
)

// transitionFile represents the JSON structure for room transitions and spawn points
// It maps room IDs to their transition and spawn point definitions
// Example JSON structure:
// {
//   "rooms": {
//     "room1": {
//       "spawn_points": [ {"id":"spawn","x":0,"y":0} ],
//       "transitions": [ {"type":0,"direction":2,"trigger_bounds":{"x":0,"y":0,"width":10,"height":10},"target_room_id":"room2","target_spawn_id":"spawn"} ]
//     }
//   }
// }

type transitionFile struct {
	Rooms map[string]struct {
		SpawnPoints []SpawnPoint      `json:"spawn_points"`
		Transitions []TransitionPoint `json:"transitions"`
	} `json:"rooms"`
}

// LoadTransitionsFromFile loads room transitions and spawn points from a JSON file
func LoadTransitionsFromFile(rtm *RoomTransitionManager, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return LoadTransitionsFromBytes(rtm, data)
}

// LoadTransitionsFromBytes loads room transitions and spawn points from a JSON byte slice
func LoadTransitionsFromBytes(rtm *RoomTransitionManager, data []byte) error {
	var cfg transitionFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	for roomID, roomData := range cfg.Rooms {
		for _, sp := range roomData.SpawnPoints {
			if err := rtm.AddSpawnPoint(roomID, sp); err != nil {
				return err
			}
		}
		for _, tp := range roomData.Transitions {
			if err := rtm.AddTransitionPoint(roomID, tp); err != nil {
				return err
			}
		}
	}
	return nil
}
