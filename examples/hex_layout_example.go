// Example: Using Hexadecimal Room Layouts
// This file demonstrates how to use the new hex format for room layouts

package main

import (
	"fmt"
	"sword/world"
)

func main() {
	fmt.Println("=== Hexadecimal Room Layout Example ===\n")
	
	// Example 1: Creating a room layout using hexadecimal values
	// This layout could be generated from the debug output and then manually edited
	exampleLayout := [][]int{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, 0x5, 0x6, 0x7, -1, -1, -1, -1, -1},  // Small platform
		{-1, -1, -1, -1, -1, -1, -1, 0xA, 0xB, -1},   // Another platform
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1}, // Ground level
		{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}, // Underground
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Solid bedrock (max value)
	}
	
	fmt.Printf("Example layout dimensions: %dx%d\n", len(exampleLayout[0]), len(exampleLayout))
	fmt.Println("Hex values used: 0x1-0x3 (ground), 0x5-0x7 (platform), 0xA-0xB (platform), 0xF (underground), 0xFF (bedrock)")
	
	// Example 2: Apply this layout to a room
	baseRoom := world.NewBaseRoom("hex_example", len(exampleLayout[0]), len(exampleLayout))
	applyLayoutToRoom(baseRoom, exampleLayout)
	
	// Generate debug output to see both formats
	fmt.Println("\nGenerating debug output...")
	baseRoom.LogRoomDebug()
	
	// Get the hex layout array (useful for copying back to code)
	fmt.Println("\nGenerated hex layout array (copy-paste ready):")
	hexLayout := baseRoom.GetHexLayoutArray()
	fmt.Println(hexLayout)
	
	// Generate a standalone .go file
	baseRoom.GenerateHexLayoutFile()
	
	fmt.Println("\n=== Benefits of Hex Format ===")
	fmt.Println("• Values 0x00-0xFF (0-255) instead of limited decimal range")
	fmt.Println("• Easy to edit in monospace editors")
	fmt.Println("• Industry standard for game development")
	fmt.Println("• Clear visual distinction between tile types")
	fmt.Println("• Copy-paste ready for Go code")
	
	fmt.Println("\n=== Files Generated ===")
	fmt.Println("• log/room_debug_*.log - Contains both decimal and hex representations")
	fmt.Println("• log/room_layout_hex_example.go - Standalone Go file with the layout")
}

// applyLayoutToRoom applies a 2D layout array to a room's tilemap
func applyLayoutToRoom(room *world.BaseRoom, layout [][]int) {
	tileMap := room.GetTileMap()
	for y, row := range layout {
		for x, tileIndex := range row {
			if x < tileMap.Width && y < tileMap.Height {
				tileMap.SetTile(x, y, tileIndex)
			}
		}
	}
}