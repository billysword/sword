package world

// ApplyLayout applies a 2D layout array to a room's tilemap
// This is the simple way to use layouts from the room_layouts package
func ApplyLayout(room *BaseRoom, layout [][]int) {
	tileMap := room.GetTileMap()
	for y, row := range layout {
		if y >= tileMap.Height {
			break
		}
		for x, tileIndex := range row {
			if x >= tileMap.Width {
				break
			}
			tileMap.SetTile(x, y, tileIndex)
		}
	}
}

// GetLayoutDimensions returns the intended width and height (in tiles)
// for a given 2D layout array.
func GetLayoutDimensions(layout [][]int) (width, height int) {
	height = len(layout)
	width = 0
	for _, row := range layout {
		if len(row) > width {
			width = len(row)
		}
	}
	return
}