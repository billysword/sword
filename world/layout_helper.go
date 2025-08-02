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