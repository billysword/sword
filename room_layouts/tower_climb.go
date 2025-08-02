package room_layouts

// TowerClimb is a vertical room with platforms for climbing
var TowerClimb = [][]int{
	{-1, -1, -1, -1, 0x5, 0x6, -1, -1, -1, -1}, // Top platform
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // Sky
	{-1, -1, 0x5, 0x6, -1, -1, -1, -1, -1, -1}, // Left platform
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // Sky
	{-1, -1, -1, -1, -1, -1, 0x5, 0x6, -1, -1}, // Right platform
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // Sky
	{0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1}, // Ground
	{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}, // Underground
}