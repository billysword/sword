package world

import (
	"testing"

	"sword/entities"
	roomsres "sword/resources/rooms"
	"sword/room_layouts"
)

// TestRoomTransitionDemo simulates moving through multiple rooms using JSON-defined transitions
func TestRoomTransitionDemo(t *testing.T) {
	worldMap := NewWorldMap()
	rtm := NewRoomTransitionManager(worldMap)

	// Create and register rooms
	mainRoom := NewSimpleRoom("main")
	ApplyLayout(mainRoom.BaseRoom, room_layouts.EmptyRoom)
	forestRight := NewSimpleRoom("forest_right")
	ApplyLayout(forestRight.BaseRoom, room_layouts.ForestRight)
	forestLeft := NewSimpleRoom("forest_left")
	ApplyLayout(forestLeft.BaseRoom, room_layouts.ForestLeft)

	rtm.RegisterRoom(mainRoom)
	rtm.RegisterRoom(forestRight)
	rtm.RegisterRoom(forestLeft)

	if err := rtm.SetCurrentRoom("main"); err != nil {
		t.Fatalf("failed to set current room: %v", err)
	}

	// Load transitions from embedded JSON
	if err := LoadTransitionsFromBytes(rtm, roomsres.RoomTransitionsJSON); err != nil {
		t.Fatalf("failed to load transitions: %v", err)
	}

	player := entities.NewPlayer(0, 0)
	if err := rtm.SpawnPlayerInRoom(player, "main", "main_spawn"); err != nil {
		t.Fatalf("failed to spawn player: %v", err)
	}

	// Transition from main -> forest_right (trigger at x>=240 for 16x10 rooms)
	player.SetPosition(248, 112)
	if !rtm.CheckTransitions(player, false) {
		t.Fatal("expected transition to forest_right")
	}
	if _, err := rtm.ProcessPendingTransition(player); err != nil {
		t.Fatalf("process transition: %v", err)
	}
	if rtm.GetCurrentRoomID() != "forest_right" {
		t.Fatalf("expected forest_right, got %s", rtm.GetCurrentRoomID())
	}

	// Transition from forest_right -> forest_left (trigger at x>=240)
	player.SetPosition(248, 112)
	if !rtm.CheckTransitions(player, false) {
		t.Fatal("expected transition to forest_left")
	}
	if _, err := rtm.ProcessPendingTransition(player); err != nil {
		t.Fatalf("process transition: %v", err)
	}
	if rtm.GetCurrentRoomID() != "forest_left" {
		t.Fatalf("expected forest_left, got %s", rtm.GetCurrentRoomID())
	}

	// Transition back to forest_right (left edge trigger at x<=16)
	player.SetPosition(8, 112)
	if !rtm.CheckTransitions(player, false) {
		t.Fatal("expected transition back to forest_right")
	}
	if _, err := rtm.ProcessPendingTransition(player); err != nil {
		t.Fatalf("process transition: %v", err)
	}
	if rtm.GetCurrentRoomID() != "forest_right" {
		t.Fatalf("expected forest_right again, got %s", rtm.GetCurrentRoomID())
	}
}
