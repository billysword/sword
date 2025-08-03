package world

import (
	"fmt"
	"strings"
)

// MiniMapRenderer handles rendering the mini-map overlay
// TODO: Replace ASCII rendering with proper Ebiten rendering system
type MiniMapRenderer struct {
	worldMap *WorldMap
	size     int    // Size of the mini-map in characters for ASCII
	visible  bool
}

// NewMiniMapRenderer creates a new mini-map renderer
func NewMiniMapRenderer(worldMap *WorldMap, size int, x, y int) *MiniMapRenderer {
	return &MiniMapRenderer{
		worldMap: worldMap,
		size:     size,
		visible:  true,
	}
}

// SetVisible toggles mini-map visibility
func (mmr *MiniMapRenderer) SetVisible(visible bool) {
	mmr.visible = visible
}

// IsVisible returns whether the mini-map is currently visible
func (mmr *MiniMapRenderer) IsVisible() bool {
	return mmr.visible
}

// ToggleVisible toggles mini-map visibility
func (mmr *MiniMapRenderer) ToggleVisible() {
	mmr.visible = !mmr.visible
}

// SetPosition updates the mini-map position on screen
// TODO: Implement proper positioning when integrating with Ebiten rendering
func (mmr *MiniMapRenderer) SetPosition(x, y int) {
	// Placeholder for future Ebiten integration
}

// Draw renders the mini-map as ASCII text for now
// TODO: Replace with proper Ebiten image rendering on screen overlay
func (mmr *MiniMapRenderer) Draw(screen interface{}, player interface{}) {
	if !mmr.visible {
		return
	}

	// For now, just print ASCII representation to console
	// TODO: Render to actual screen surface with Ebiten
	asciiMap := mmr.generateASCIIMap(player)
	if asciiMap != "" {
		fmt.Println("=== MINI-MAP DEBUG ===")
		fmt.Println(asciiMap)
		fmt.Println("=====================")
	}
}

// generateASCIIMap creates a simple ASCII representation of the current room
// TODO: Replace with pixel-based rendering for actual game display
func (mmr *MiniMapRenderer) generateASCIIMap(player interface{}) string {
	currentRoomID := mmr.worldMap.GetCurrentRoom()
	if currentRoomID == "" {
		return ""
	}

	discoveredRooms := mmr.worldMap.GetDiscoveredRooms()
	currentRoom, exists := discoveredRooms[currentRoomID]
	if !exists {
		return ""
	}

	// Create ASCII grid
	mapSize := 20 // 20x20 character map
	grid := make([][]rune, mapSize)
	for i := range grid {
		grid[i] = make([]rune, mapSize)
		for j := range grid[i] {
			grid[i][j] = ' ' // Empty space
		}
	}

	// Draw room boundary
	for i := 0; i < mapSize; i++ {
		for j := 0; j < mapSize; j++ {
			if i == 0 || i == mapSize-1 || j == 0 || j == mapSize-1 {
				grid[i][j] = '#' // Room walls
			}
		}
	}

	// Draw room thumbnail if available
	if len(currentRoom.ThumbnailData) > 0 {
		mmr.drawThumbnailASCII(grid, currentRoom.ThumbnailData, mapSize)
	}

	// Draw player position (center for now)
	// TODO: Calculate actual player position relative to room when player interface is available
	centerX, centerY := mapSize/2, mapSize/2
	grid[centerY][centerX] = '@' // Player marker

	// Draw exits
	connections := mmr.worldMap.GetRoomConnections(currentRoom.ZoneID)
	for direction := range connections {
		mmr.drawExitASCII(grid, direction, mapSize)
	}

	// Convert grid to string
	var result strings.Builder
	for _, row := range grid {
		result.WriteString(string(row))
		result.WriteString("\n")
	}

	return result.String()
}

// drawThumbnailASCII draws room thumbnail data as ASCII characters
// TODO: Replace with actual tile sprite rendering
func (mmr *MiniMapRenderer) drawThumbnailASCII(grid [][]rune, thumbnail [][]int, mapSize int) {
	if len(thumbnail) == 0 {
		return
	}

	thumbnailHeight := len(thumbnail)
	thumbnailWidth := len(thumbnail[0])

	// Map thumbnail to grid (skip borders)
	for y := 1; y < mapSize-1 && y-1 < thumbnailHeight; y++ {
		for x := 1; x < mapSize-1 && x-1 < thumbnailWidth; x++ {
			tileIndex := thumbnail[y-1][x-1]
			
			// Convert tile index to ASCII character
			// TODO: Use actual tile sprites when rendering to Ebiten
			switch {
			case tileIndex < 0:
				grid[y][x] = '.' // Empty space
			case tileIndex == 0:
				grid[y][x] = '█' // Solid tile
			default:
				grid[y][x] = '▓' // Other tile types
			}
		}
	}
}

// drawExitASCII draws exit indicators in ASCII
// TODO: Replace with proper exit rendering (colored dots, arrows, etc.)
func (mmr *MiniMapRenderer) drawExitASCII(grid [][]rune, direction Direction, mapSize int) {
	midPoint := mapSize / 2
	
	switch direction {
	case North:
		if midPoint < len(grid[0]) {
			grid[0][midPoint] = '^'
		}
	case South:
		if midPoint < len(grid[mapSize-1]) {
			grid[mapSize-1][midPoint] = 'v'
		}
	case East:
		if midPoint < len(grid) {
			grid[midPoint][mapSize-1] = '>'
		}
	case West:
		if midPoint < len(grid) {
			grid[midPoint][0] = '<'
		}
	// TODO: Add diagonal and vertical direction indicators
	}
}

// GetMapData returns the current map data for external rendering systems
// This is the key method for providing data to proper rendering implementations
func (mmr *MiniMapRenderer) GetMapData() *MapDisplayData {
	currentRoomID := mmr.worldMap.GetCurrentRoom()
	if currentRoomID == "" {
		return nil
	}

	discoveredRooms := mmr.worldMap.GetDiscoveredRooms()
	currentRoom, exists := discoveredRooms[currentRoomID]
	if !exists {
		return nil
	}

	connections := mmr.worldMap.GetRoomConnections(currentRoom.ZoneID)
	playerTrail := mmr.worldMap.GetPlayerTrail()

	return &MapDisplayData{
		CurrentRoom:    currentRoom,
		Connections:    connections,
		PlayerTrail:    playerTrail,
		DiscoveredRooms: discoveredRooms,
		MapBounds:      mmr.worldMap.GetMapBounds(),
	}
}

// MapDisplayData contains all the data needed for rendering the map
// This struct should be used by proper rendering implementations
type MapDisplayData struct {
	CurrentRoom     *DiscoveredRoom
	Connections     map[Direction]string
	PlayerTrail     []Point
	DiscoveredRooms map[string]*DiscoveredRoom
	MapBounds       interface{} // image.Rectangle when imported
}

// TODO: Implement proper Ebiten rendering methods:
// - DrawToEbitenImage(screen *ebiten.Image, player *entities.Player)
// - RenderRoomThumbnail(screen *ebiten.Image, room *DiscoveredRoom)
// - RenderPlayerIndicator(screen *ebiten.Image, player *entities.Player, room *DiscoveredRoom)
// - RenderPlayerTrail(screen *ebiten.Image, trail []Point, room *DiscoveredRoom)
// - RenderExitIndicators(screen *ebiten.Image, connections map[Direction]string)
// - RenderAdjacentRooms(screen *ebiten.Image, rooms map[string]*DiscoveredRoom)