package world

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"sword/entities"
)

// MiniMapRenderer handles rendering the mini-map overlay
// Implements engine.HUDComponent interface
type MiniMapRenderer struct {
	worldMap *WorldMap
	player   *entities.Player
	visible  bool
	name     string
	size     int
	x        int
	y        int
}

// NewMiniMapRenderer creates a new mini-map renderer
func NewMiniMapRenderer(worldMap *WorldMap, player *entities.Player, size int, x, y int) *MiniMapRenderer {
	return &MiniMapRenderer{
		worldMap: worldMap,
		player:   player,
		visible:  true,
		name:     "minimap",
		size:     size,
		x:        x,
		y:        y,
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
	// Placeholder for future mini-map update logic
	return nil
}

// Draw renders the mini-map (required by HUDComponent interface)
func (mmr *MiniMapRenderer) Draw(screen interface{}) error {
	if !mmr.visible {
		return nil
	}

	ebScreen, ok := screen.(*ebiten.Image)
	if !ok {
		return nil
	}

	mapData := mmr.GetMapData()
	if mapData == nil || len(mapData.DiscoveredRooms) == 0 {
		return nil
	}

	// Panel background
	panelBg := color.RGBA{0, 0, 0, 160}
	vector.DrawFilledRect(ebScreen, float32(mmr.x-4), float32(mmr.y-4), float32(mmr.size+8), float32(mmr.size+8), panelBg, false)

	bounds := mapData.MapBounds
	bw := bounds.Dx()
	bh := bounds.Dy()
	if bw <= 0 || bh <= 0 {
		return nil
	}

	// Maintain aspect ratio within square minimap region
	scaleX := float32(mmr.size) / float32(bw)
	scaleY := float32(mmr.size) / float32(bh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	// Center the map within the panel if aspect ratios differ
	mapDrawW := float32(bw) * scale
	mapDrawH := float32(bh) * scale
	offsetX := float32(mmr.x) + (float32(mmr.size)-mapDrawW)/2
	offsetY := float32(mmr.y) + (float32(mmr.size)-mapDrawH)/2

	// Helper to convert world-map coordinates to minimap pixels
	toMini := func(wx, wy int) (float32, float32) {
		x := float32(wx-bounds.Min.X) * scale
		y := float32(wy-bounds.Min.Y) * scale
		return offsetX + x, offsetY + y
	}

	// Draw discovered rooms as rectangles
	roomColor := color.RGBA{80, 80, 160, 220}
	currentRoomColor := color.RGBA{200, 200, 255, 255}
	borderColor := color.RGBA{220, 220, 240, 255}

	for _, room := range mapData.DiscoveredRooms {
		left := room.WorldPos.X - room.Bounds.Dx()/2
		right := room.WorldPos.X + room.Bounds.Dx()/2
		top := room.WorldPos.Y - room.Bounds.Dy()/2
		bottom := room.WorldPos.Y + room.Bounds.Dy()/2

		x1, y1 := toMini(left, top)
		x2, y2 := toMini(right, bottom)
		w := x2 - x1
		h := y2 - y1

		fill := roomColor
		if mapData.CurrentRoom != nil && room.ZoneID == mapData.CurrentRoom.ZoneID {
			fill = currentRoomColor
		}
		vector.DrawFilledRect(ebScreen, x1, y1, w, h, fill, false)
		vector.StrokeRect(ebScreen, x1, y1, w, h, 1, borderColor, false)
	}

	// Draw connections (simple lines between room centers)
	connColor := color.RGBA{200, 200, 200, 200}
	for fromID, fromRoom := range mapData.DiscoveredRooms {
		conns := mmr.worldMap.GetRoomConnections(fromID)
		for _, toID := range conns {
			toRoom, ok := mapData.DiscoveredRooms[toID]
			if !ok {
				continue
			}
			x1, y1 := toMini(fromRoom.WorldPos.X, fromRoom.WorldPos.Y)
			x2, y2 := toMini(toRoom.WorldPos.X, toRoom.WorldPos.Y)
			vector.StrokeLine(ebScreen, x1, y1, x2, y2, 1, connColor, false)
		}
	}

	// Draw player trail (faint)
	trail := mapData.PlayerTrail
	if len(trail) > 1 && mapData.CurrentRoom != nil {
		trailColor := color.RGBA{255, 255, 255, 120}
		// Convert each trail point from room-local to world-map space by assuming they
		// were recorded in the context of the current room at the time. We render only the
		// last N points to avoid clutter.
		const maxTrail = 20
		start := 0
		if len(trail) > maxTrail {
			start = len(trail) - maxTrail
		}
		prevSet := false
		var px, py float32
		for i := start; i < len(trail); i++ {
			pt := trail[i]
			room := mapData.CurrentRoom
			roomLeft := room.WorldPos.X - room.Bounds.Dx()/2
			worldX := roomLeft + pt.X
			roomTop := room.WorldPos.Y - room.Bounds.Dy()/2
			worldY := roomTop + pt.Y
			x, y := toMini(worldX, worldY)
			if prevSet {
				vector.StrokeLine(ebScreen, px, py, x, y, 1, trailColor, false)
			}
			px, py = x, y
			prevSet = true
		}
	}

	// Draw player indicator (triangle)
	if mmr.player != nil && mapData.CurrentRoom != nil {
		px, py := mmr.player.GetPosition()
		room := mapData.CurrentRoom
		roomLeft := room.WorldPos.X - room.Bounds.Dx()/2
		roomTop := room.WorldPos.Y - room.Bounds.Dy()/2
		worldX := roomLeft + px
		worldY := roomTop + py
		x, y := toMini(worldX, worldY)
		// Simple triangle pointing to facing direction
		size := float32(4)
		c := color.RGBA{255, 80, 80, 255}
		if mmr.player.IsFacingRight() {
			vector.DrawFilledTriangle(ebScreen, x+size, y, x-size, y-size, x-size, y+size, c, false)
		} else {
			vector.DrawFilledTriangle(ebScreen, x-size, y, x+size, y-size, x+size, y+size, c, false)
		}
	}

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
	MapBounds       image.Rectangle
}

// TODO: Implement additional rendering helpers for reuse across overlays if needed.