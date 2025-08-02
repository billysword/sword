package world

import (
	"fmt"
	"strings"
)

// PrintRoomLayout prints a simple ASCII representation of the room layout to console
func PrintRoomLayout(roomName string, tileMap *TileMap) {
	fmt.Printf("\n=== ROOM LAYOUT: %s ===\n", roomName)
	fmt.Printf("Dimensions: %dx%d tiles\n\n", tileMap.Width, tileMap.Height)
	
	// Print hex format for easy copying
	fmt.Println("Hex format (copy-paste ready):")
	printHexLayout(tileMap)
	
	fmt.Printf("\n=== END ROOM LAYOUT ===\n\n")
}

// printHexLayout prints the layout in hex format ready for copying
func printHexLayout(tileMap *TileMap) {
	fmt.Println("var layout = [][]int{")
	for y := 0; y < tileMap.Height; y++ {
		fmt.Print("\t{")
		var values []string
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]
			if tileIndex == -1 {
				values = append(values, "-1")
			} else {
				if tileIndex > 255 {
					tileIndex = 255
				}
				if tileIndex < 16 {
					values = append(values, fmt.Sprintf("0x%X", tileIndex))
				} else {
					values = append(values, fmt.Sprintf("0x%02X", tileIndex))
				}
			}
		}
		fmt.Print(strings.Join(values, ", "))
		if y < tileMap.Height-1 {
			fmt.Println("},")
		} else {
			fmt.Println("}")
		}
	}
	fmt.Println("}")
}