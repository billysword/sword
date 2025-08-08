package world

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"sword/entities"
)

// ZoneMapOverlay renders a large overlay of the current zone (room + immediate neighbors)
// Implements engine.HUDComponent
type ZoneMapOverlay struct {
	worldMap *WorldMap
	player   *entities.Player
	visible  bool
	name     string
}

func NewZoneMapOverlay(worldMap *WorldMap, player *entities.Player) *ZoneMapOverlay {
	return &ZoneMapOverlay{
		worldMap: worldMap,
		player:   player,
		visible:  false,
		name:     "zone_map",
	}
}

func (zmo *ZoneMapOverlay) GetName() string { return zmo.name }
func (zmo *ZoneMapOverlay) SetVisible(v bool) { zmo.visible = v }
func (zmo *ZoneMapOverlay) IsVisible() bool { return zmo.visible }

func (zmo *ZoneMapOverlay) Update() error { return nil }

func (zmo *ZoneMapOverlay) Draw(screen interface{}) error {
	if !zmo.visible {
		return nil
	}
	img, ok := screen.(*ebiten.Image)
	if !ok {
		return nil
	}
	mapData := zmo.getZoneData()
	if mapData == nil || len(mapData.DiscoveredRooms) == 0 {
		return nil
	}

	// Darken background
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
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
	drawn := make(map[string]struct{})
	for fromID, fromRoom := range mapData.DiscoveredRooms {
		conns := zmo.worldMap.GetRoomConnections(fromID)
		for _, toID := range conns {
			toRoom, ok := mapData.DiscoveredRooms[toID]
			if !ok {
				continue
			}
			// Avoid drawing both directions
			a, b := fromID, toID
			if a > b {
				a, b = b, a
			}
			key := a + "|" + b
			if _, exists := drawn[key]; exists {
				continue
			}
			x1, y1 := toMini(fromRoom.WorldPos.X, fromRoom.WorldPos.Y)
			x2, y2 := toMini(toRoom.WorldPos.X, toRoom.WorldPos.Y)
			vector.StrokeLine(img, x1, y1, x2, y2, 2, connColor, false)
			drawn[key] = struct{}{}
		}
	}

	// Player
	if zmo.player != nil && mapData.CurrentRoom != nil {
		px, py := zmo.player.GetPosition()
		room := mapData.CurrentRoom
		roomLeft := room.WorldPos.X - room.Bounds.Dx()/2
		roomTop := room.WorldPos.Y - room.Bounds.Dy()/2
		worldX := roomLeft + px
		worldY := roomTop + py
		x, y := toMini(worldX, worldY)
		size := float32(6)
		c := color.RGBA{255, 80, 80, 255}
		if zmo.player.IsFacingRight() {
			vector.DrawFilledTriangle(img, x+size, y, x-size, y-size, x-size, y+size, c, false)
		} else {
			vector.DrawFilledTriangle(img, x-size, y, x+size, y-size, x+size, y+size, c, false)
		}
	}

	// Help text
	ebitenutil.DebugPrintAt(img, "Zone Map (Z to close)", int(drawX), int(drawY)-20)

	return nil
}

// Internal data provider: restrict to current room and directly connected rooms only
func (zmo *ZoneMapOverlay) getZoneData() *struct {
	CurrentRoom     *DiscoveredRoom
	DiscoveredRooms map[string]*DiscoveredRoom
	MapBounds       image.Rectangle
} {
	currentRoomID := zmo.worldMap.GetCurrentRoom()
	if currentRoomID == "" {
		return nil
	}
	allRooms := zmo.worldMap.GetDiscoveredRooms()
	currentRoom, ok := allRooms[currentRoomID]
	if !ok {
		return nil
	}
	// Build filtered set: current + immediate neighbors
	filtered := make(map[string]*DiscoveredRoom)
	filtered[currentRoomID] = currentRoom
	conns := zmo.worldMap.GetRoomConnections(currentRoomID)
	for _, neighborID := range conns {
		if r, exists := allRooms[neighborID]; exists {
			filtered[neighborID] = r
		}
	}

	// Compute bounds of filtered rooms
	if len(filtered) == 0 {
		return nil
	}
	first := true
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, room := range filtered {
		left := room.WorldPos.X - room.Bounds.Dx()/2
		right := room.WorldPos.X + room.Bounds.Dx()/2
		top := room.WorldPos.Y - room.Bounds.Dy()/2
		bottom := room.WorldPos.Y + room.Bounds.Dy()/2
		if first {
			minX, maxX = left, right
			minY, maxY = top, bottom
			first = false
		} else {
			if left < minX {
				minX = left
			}
			if right > maxX {
				maxX = right
			}
			if top < minY {
				minY = top
			}
			if bottom > maxY {
				maxY = bottom
			}
		}
	}
	bounds := image.Rect(minX, minY, maxX, maxY)
	return &struct {
		CurrentRoom     *DiscoveredRoom
		DiscoveredRooms map[string]*DiscoveredRoom
		MapBounds       image.Rectangle
	}{
		CurrentRoom:     currentRoom,
		DiscoveredRooms: filtered,
		MapBounds:       bounds,
	}
}