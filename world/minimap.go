package world

// MiniMapRenderer handles rendering the mini-map overlay
// Implements engine.HUDComponent interface
type MiniMapRenderer struct {
	worldMap *WorldMap
	visible  bool
	name     string
}

// NewMiniMapRenderer creates a new mini-map renderer
func NewMiniMapRenderer(worldMap *WorldMap, size int, x, y int) *MiniMapRenderer {
	return &MiniMapRenderer{
		worldMap: worldMap,
		visible:  true,
		name:     "minimap",
	}
}

// GetName returns the component name (required by HUDComponent interface)
func (mmr *MiniMapRenderer) GetName() string {
	return mmr.name
}

// SetVisible toggles mini-map visibility (required by HUDComponent interface)
func (mmr *MiniMapRenderer) SetVisible(visible bool) {
	mmr.visible = visible
}

// IsVisible returns whether the mini-map is currently visible (required by HUDComponent interface)
func (mmr *MiniMapRenderer) IsVisible() bool {
	return mmr.visible
}

// ToggleVisible toggles mini-map visibility
func (mmr *MiniMapRenderer) ToggleVisible() {
	mmr.visible = !mmr.visible
}

// Update handles mini-map logic updates (required by HUDComponent interface)
func (mmr *MiniMapRenderer) Update() error {
	// Empty for now - placeholder for future mini-map update logic
	// TODO: Implement mini-map specific update logic if needed
	return nil
}

// Draw renders the mini-map (required by HUDComponent interface)
func (mmr *MiniMapRenderer) Draw(screen interface{}) error {
	if !mmr.visible {
		return nil
	}
	
	// Empty for now - placeholder for future mini-map rendering
	// TODO: Render mini-map overlay to screen using Ebiten
	return nil
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