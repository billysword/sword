package rooms

import (
	_ "embed"
)

// RoomTransitionsJSON contains the default room transitions and spawn points configuration.
//go:embed room_transitions.json
var RoomTransitionsJSON []byte