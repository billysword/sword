package world

import (
	"fmt"
	"strings"
	"sword/engine"
)

// PrintRoomLayout logs a simple ASCII representation of the room layout to file
func PrintRoomLayout(roomName string, tileMap *TileMap) {
	// Log basic room info
	engine.LogRoomTile(roomName, fmt.Sprintf("Dimensions: %dx%d tiles", tileMap.Width, tileMap.Height))

	// Generate hex layout string
	layoutStr := generateHexLayoutString(tileMap)

	// Log the complete layout
	engine.LogRoomLayout(roomName, tileMap.Width, tileMap.Height, layoutStr)
}

// generateHexLayoutString generates the layout in hex format ready for copying
func generateHexLayoutString(tileMap *TileMap) string {
	var result strings.Builder
	result.WriteString("var layout = [][]int{\n")
	for y := 0; y < tileMap.Height; y++ {
		result.WriteString("\t{")
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
		result.WriteString(strings.Join(values, ", "))
		if y < tileMap.Height-1 {
			result.WriteString("},\n")
		} else {
			result.WriteString("}\n")
		}
	}
	result.WriteString("}")
	return result.String()
}
