package main

import (
	"fmt"
	"sword/room_layouts"
	"sword/world"
)

func main() {
	fmt.Println("=== Simple Room Layout Usage ===\n")
	
	// Create a room
	room := world.NewBaseRoom("test_room", 10, 8)
	
	// Option 1: Use a predefined layout
	fmt.Println("1. Using predefined layout (ExamplePlatform):")
	world.ApplyLayout(room, room_layouts.ExamplePlatform)
	world.PrintRoomLayout("ExamplePlatform", room.GetTileMap())
	
	// Option 2: Use a different predefined layout
	fmt.Println("2. Using different layout (TowerClimb):")
	room2 := world.NewBaseRoom("tower_room", 10, 8)
	world.ApplyLayout(room2, room_layouts.TowerClimb)
	world.PrintRoomLayout("TowerClimb", room2.GetTileMap())
	
	// Option 3: Create your own layout inline
	fmt.Println("3. Using custom inline layout:")
	customLayout := [][]int{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, 0x8, 0x9, 0xA, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1},
	}
	
	room3 := world.NewBaseRoom("custom_room", 10, 4)
	world.ApplyLayout(room3, customLayout)
	world.PrintRoomLayout("Custom", room3.GetTileMap())
	
	fmt.Println("That's it! Much simpler than the complex file generation system.")
	fmt.Println("Just import room_layouts and use world.ApplyLayout().")
}