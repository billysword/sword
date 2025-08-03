package world

// MiniMapRenderer handles rendering the mini-map overlay
// TODO: Implement proper Ebiten rendering system
type MiniMapRenderer struct {
	worldMap *WorldMap
	visible  bool
}

// NewMiniMapRenderer creates a new mini-map renderer
func NewMiniMapRenderer(worldMap *WorldMap, size int, x, y int) *MiniMapRenderer {
	return &MiniMapRenderer{
		worldMap: worldMap,
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

// Update handles mini-map logic updates
// TODO: Implement mini-map specific update logic if needed
func (mmr *MiniMapRenderer) Update() {
	// Empty for now - placeholder for future mini-map update logic
}

// Draw renders the mini-map
// TODO: Implement proper Ebiten image rendering on screen overlay
func (mmr *MiniMapRenderer) Draw(screen interface{}, player interface{}) {
	if !mmr.visible {
		return
	}
	
	// Empty for now - placeholder for future mini-map rendering
	// TODO: Render mini-map overlay to screen using Ebiten
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
		CurrentRoom:     currentRoom,
		Connections:     connections,
		PlayerTrail:     playerTrail,
		DiscoveredRooms: discoveredRooms,
		MapBounds:       mmr.worldMap.GetMapBounds(),
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