package world

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"sword/entities"
)

// WorldMapOverlay renders a large overlay of the world map
// Implements engine.HUDComponent
type WorldMapOverlay struct {
	worldMap *WorldMap
	player   *entities.Player
	visible  bool
	name     string
}

func NewWorldMapOverlay(worldMap *WorldMap, player *entities.Player) *WorldMapOverlay {
	return &WorldMapOverlay{
		worldMap: worldMap,
		player:   player,
		visible:  false,
		name:     "world_map",
	}
}

func (wmo *WorldMapOverlay) GetName() string { return wmo.name }
func (wmo *WorldMapOverlay) SetVisible(v bool) { wmo.visible = v }
func (wmo *WorldMapOverlay) IsVisible() bool { return wmo.visible }

func (wmo *WorldMapOverlay) Update() error { return nil }

func (wmo *WorldMapOverlay) Draw(screen interface{}) error {
	if !wmo.visible {
		return nil
	}
	img, ok := screen.(*ebiten.Image)
	if !ok {
		return nil
	}
	mapData := wmo.getMapData()
	if mapData == nil || len(mapData.DiscoveredRooms) == 0 {
		return nil
	}

	// Darken background
	w, h := ebiten.WindowSize()
	vector.DrawFilledRect(img, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 180}, false)

	// Map drawing area (padding)
	pad := float32(40)
	drawX := pad
	drawY := pad
	drawW := float32(w) - pad*2
	drawH := float32(h) - pad*2
	vector.StrokeRect(img, drawX-2, drawY-2, drawW+4, drawH+4, 2, color.RGBA{255, 255, 255, 50}, false)

	bounds := mapData.MapBounds
	bw := bounds.Dx()
	bh := bounds.Dy()
	if bw <= 0 || bh <= 0 {
		return nil
	}

	scaleX := drawW / float32(bw)
	scaleY := drawH / float32(bh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	mapDrawW := float32(bw) * scale
	mapDrawH := float32(bh) * scale
	offsetX := drawX + (drawW-mapDrawW)/2
	offsetY := drawY + (drawH-mapDrawH)/2

	toMini := func(wx, wy int) (float32, float32) {
		x := float32(wx-bounds.Min.X) * scale
		y := float32(wy-bounds.Min.Y) * scale
		return offsetX + x, offsetY + y
	}

	roomColor := color.RGBA{60, 120, 180, 200}
	currentRoomColor := color.RGBA{220, 240, 255, 255}
	borderColor := color.RGBA{255, 255, 255, 180}

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
		vector.DrawFilledRect(img, x1, y1, w, h, fill, false)
		vector.StrokeRect(img, x1, y1, w, h, 2, borderColor, false)

		// Label room
		ebitenutil.DebugPrintAt(img, room.Name, int(x1)+4, int(y1)+4)
	}

	// Connections
	connColor := color.RGBA{255, 255, 255, 180}
	for fromID, fromRoom := range mapData.DiscoveredRooms {
		conns := wmo.worldMap.GetRoomConnections(fromID)
		for _, toID := range conns {
			toRoom, ok := mapData.DiscoveredRooms[toID]
			if !ok {
				continue
			}
			x1, y1 := toMini(fromRoom.WorldPos.X, fromRoom.WorldPos.Y)
			x2, y2 := toMini(toRoom.WorldPos.X, toRoom.WorldPos.Y)
			vector.StrokeLine(img, x1, y1, x2, y2, 2, connColor, false)
		}
	}

	// Player
	if wmo.player != nil && mapData.CurrentRoom != nil {
		px, py := wmo.player.GetPosition()
		room := mapData.CurrentRoom
		roomLeft := room.WorldPos.X - room.Bounds.Dx()/2
		roomTop := room.WorldPos.Y - room.Bounds.Dy()/2
		worldX := roomLeft + px
		worldY := roomTop + py
		x, y := toMini(worldX, worldY)
		size := float32(6)
		c := color.RGBA{255, 80, 80, 255}
		if wmo.player.IsFacingRight() {
			vector.DrawFilledTriangle(img, x+size, y, x-size, y-size, x-size, y+size, c, false)
		} else {
			vector.DrawFilledTriangle(img, x-size, y, x+size, y-size, x+size, y+size, c, false)
		}
	}

	// Help text
	ebitenutil.DebugPrintAt(img, "World Map (N to close)", int(drawX), int(drawY)-20)

	return nil
}

// Internal data provider (copy of minimap with concrete bounds type)
func (wmo *WorldMapOverlay) getMapData() *struct {
	CurrentRoom     *DiscoveredRoom
	DiscoveredRooms map[string]*DiscoveredRoom
	MapBounds       image.Rectangle
} {
	currentRoomID := wmo.worldMap.GetCurrentRoom()
	if currentRoomID == "" {
		return nil
	}
	discoveredRooms := wmo.worldMap.GetDiscoveredRooms()
	currentRoom, ok := discoveredRooms[currentRoomID]
	if !ok {
		return nil
	}
	return &struct {
		CurrentRoom     *DiscoveredRoom
		DiscoveredRooms map[string]*DiscoveredRoom
		MapBounds       image.Rectangle
	}{
		CurrentRoom:     currentRoom,
		DiscoveredRooms: discoveredRooms,
		MapBounds:       wmo.worldMap.GetMapBounds(),
	}
}