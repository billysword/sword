package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"sword/entities"
)

// DrawFullWorldMap renders the full discovered world map to the given image
func DrawFullWorldMap(img *ebiten.Image, wm *WorldMap, player *entities.Player) {
	if wm == nil {
		return
	}
	rooms := wm.GetDiscoveredRooms()
	if len(rooms) == 0 {
		return
	}
	bounds := wm.GetMapBounds()
	bw, bh := bounds.Dx(), bounds.Dy()
	if bw <= 0 || bh <= 0 {
		return
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	pad := float32(40)
	drawX := pad
	drawY := pad
	drawW := float32(w) - pad*2
	drawH := float32(h) - pad*2
	vector.StrokeRect(img, drawX-2, drawY-2, drawW+4, drawH+4, 2, color.RGBA{255, 255, 255, 50}, false)

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
	borderColor := color.RGBA{255, 255, 255, 180}
	for _, room := range rooms {
		left := room.WorldPos.X - room.Bounds.Dx()/2
		right := room.WorldPos.X + room.Bounds.Dx()/2
		top := room.WorldPos.Y - room.Bounds.Dy()/2
		bottom := room.WorldPos.Y + room.Bounds.Dy()/2
		x1, y1 := toMini(left, top)
		x2, y2 := toMini(right, bottom)
		w := x2 - x1
		h := y2 - y1
		vector.DrawFilledRect(img, x1, y1, w, h, roomColor, false)
		vector.StrokeRect(img, x1, y1, w, h, 2, borderColor, false)
	}

	connColor := color.RGBA{255, 255, 255, 180}
	drawn := make(map[string]struct{})
	for fromID, fromRoom := range rooms {
		conns := wm.GetRoomConnections(fromID)
		for _, toID := range conns {
			toRoom, ok := rooms[toID]
			if !ok {
				continue
			}
			// Avoid duplicate undirected edges
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

	// Player marker in current room
	if player != nil {
		currentID := wm.GetCurrentRoom()
		if r, ok := rooms[currentID]; ok {
			px, py := player.GetPosition()
			roomLeft := r.WorldPos.X - r.Bounds.Dx()/2
			roomTop := r.WorldPos.Y - r.Bounds.Dy()/2
			x, y := toMini(roomLeft+px, roomTop+py)
			c := color.RGBA{255, 80, 80, 255}
			size := float32(6)
			if player.IsFacingRight() {
				vector.DrawFilledTriangle(img, x+size, y, x-size, y-size, x-size, y+size, c, false)
			} else {
				vector.DrawFilledTriangle(img, x-size, y, x+size, y-size, x+size, y+size, c, false)
			}
		}
	}
}