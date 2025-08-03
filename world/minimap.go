package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"sword/entities"
)

// MiniMapRenderer handles rendering the mini-map overlay
type MiniMapRenderer struct {
	worldMap     *WorldMap
	size         int                    // Size of the mini-map in pixels
	position     Point                  // Position on screen (top-left corner)
	scale        float64               // Current scale factor
	borderColor  color.Color
	bgColor      color.Color
	roomColor    color.Color
	playerColor  color.Color
	trailColor   color.Color
	exitColor    color.Color
	visible      bool
}

// NewMiniMapRenderer creates a new mini-map renderer
func NewMiniMapRenderer(worldMap *WorldMap, size int, x, y int) *MiniMapRenderer {
	return &MiniMapRenderer{
		worldMap:    worldMap,
		size:        size,
		position:    Point{X: x, Y: y},
		scale:       1.0,
		borderColor: color.RGBA{255, 255, 255, 200}, // White border
		bgColor:     color.RGBA{0, 0, 0, 150},       // Semi-transparent black background
		roomColor:   color.RGBA{100, 100, 100, 200}, // Gray rooms
		playerColor: color.RGBA{255, 0, 0, 255},     // Red player dot
		trailColor:  color.RGBA{255, 255, 0, 100},   // Yellow trail
		exitColor:   color.RGBA{0, 255, 0, 150},     // Green exits
		visible:     true,
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
func (mmr *MiniMapRenderer) SetPosition(x, y int) {
	mmr.position.X = x
	mmr.position.Y = y
}

// Draw renders the mini-map overlay on the screen
func (mmr *MiniMapRenderer) Draw(screen *ebiten.Image, player *entities.Player) {
	if !mmr.visible {
		return
	}

	// Get current room
	currentRoomID := mmr.worldMap.GetCurrentRoom()
	if currentRoomID == "" {
		return
	}

	discoveredRooms := mmr.worldMap.GetDiscoveredRooms()
	currentRoom, exists := discoveredRooms[currentRoomID]
	if !exists {
		return
	}

	// Calculate scale to fit current room in mini-map
	roomWidth := float64(currentRoom.Bounds.Dx())
	roomHeight := float64(currentRoom.Bounds.Dy())
	maxDimension := math.Max(roomWidth, roomHeight)
	mmr.scale = float64(mmr.size-20) / maxDimension // Leave some border space

	// Draw background
	mmr.drawBackground(screen)

	// Draw current room
	mmr.drawCurrentRoom(screen, currentRoom)

	// Draw adjacent rooms (if connected)
	mmr.drawAdjacentRooms(screen, currentRoom, discoveredRooms)

	// Draw player trail
	mmr.drawPlayerTrail(screen, currentRoom)

	// Draw player position
	mmr.drawPlayer(screen, player, currentRoom)

	// Draw border
	mmr.drawBorder(screen)
}

// drawBackground draws the mini-map background
func (mmr *MiniMapRenderer) drawBackground(screen *ebiten.Image) {
	x := float32(mmr.position.X)
	y := float32(mmr.position.Y)
	size := float32(mmr.size)

	vector.DrawFilledRect(screen, x, y, size, size, mmr.bgColor, false)
}

// drawBorder draws the mini-map border
func (mmr *MiniMapRenderer) drawBorder(screen *ebiten.Image) {
	x := float32(mmr.position.X)
	y := float32(mmr.position.Y)
	size := float32(mmr.size)

	// Draw border lines
	vector.StrokeLine(screen, x, y, x+size, y, 2, mmr.borderColor, false)         // Top
	vector.StrokeLine(screen, x, y, x, y+size, 2, mmr.borderColor, false)         // Left
	vector.StrokeLine(screen, x+size, y, x+size, y+size, 2, mmr.borderColor, false) // Right
	vector.StrokeLine(screen, x, y+size, x+size, y+size, 2, mmr.borderColor, false) // Bottom
}

// drawCurrentRoom draws the current room layout
func (mmr *MiniMapRenderer) drawCurrentRoom(screen *ebiten.Image, room *DiscoveredRoom) {
	centerX := mmr.position.X + mmr.size/2
	centerY := mmr.position.Y + mmr.size/2

	// Draw room thumbnail
	if len(room.ThumbnailData) > 0 {
		mmr.drawRoomThumbnail(screen, room.ThumbnailData, centerX, centerY, mmr.roomColor)
	}

	// Draw exits
	mmr.drawRoomExits(screen, room, centerX, centerY)
}

// drawRoomThumbnail draws a simplified version of the room layout
func (mmr *MiniMapRenderer) drawRoomThumbnail(screen *ebiten.Image, thumbnail [][]int, centerX, centerY int, roomColor color.Color) {
	if len(thumbnail) == 0 {
		return
	}

	thumbnailHeight := len(thumbnail)
	thumbnailWidth := len(thumbnail[0])

	// Calculate pixel size for each thumbnail tile
	tileSize := float32(mmr.size-40) / float32(math.Max(float64(thumbnailWidth), float64(thumbnailHeight)))
	if tileSize < 1 {
		tileSize = 1
	}

	// Calculate starting position to center the thumbnail
	startX := float32(centerX) - float32(thumbnailWidth)*tileSize/2
	startY := float32(centerY) - float32(thumbnailHeight)*tileSize/2

	// Draw each tile
	for y := 0; y < thumbnailHeight; y++ {
		for x := 0; x < thumbnailWidth; x++ {
			tileIndex := thumbnail[y][x]
			if tileIndex >= 0 { // Only draw non-empty tiles
				tileX := startX + float32(x)*tileSize
				tileY := startY + float32(y)*tileSize

				// Use different colors for different tile types
				var tileColor color.Color = roomColor
				if tileIndex == 0 { // Assume 0 is solid ground
					tileColor = color.RGBA{80, 80, 80, 200} // Darker for solid tiles
				}

				vector.DrawFilledRect(screen, tileX, tileY, tileSize, tileSize, tileColor, false)
			}
		}
	}
}

// drawRoomExits draws exit indicators on the room
func (mmr *MiniMapRenderer) drawRoomExits(screen *ebiten.Image, room *DiscoveredRoom, centerX, centerY int) {
	connections := mmr.worldMap.GetRoomConnections(room.ZoneID)
	
	// Calculate room bounds on mini-map
	roomWidth := float64(room.Bounds.Dx()) * mmr.scale
	roomHeight := float64(room.Bounds.Dy()) * mmr.scale

	for direction := range connections {
		var exitX, exitY float32

		switch direction {
		case North:
			exitX = float32(centerX)
			exitY = float32(centerY) - float32(roomHeight/2)
		case South:
			exitX = float32(centerX)
			exitY = float32(centerY) + float32(roomHeight/2)
		case East:
			exitX = float32(centerX) + float32(roomWidth/2)
			exitY = float32(centerY)
		case West:
			exitX = float32(centerX) - float32(roomWidth/2)
			exitY = float32(centerY)
		default:
			continue // Skip Up/Down for now
		}

		// Draw exit indicator
		vector.DrawFilledCircle(screen, exitX, exitY, 3, mmr.exitColor, false)
	}
}

// drawAdjacentRooms draws outlines of connected adjacent rooms
func (mmr *MiniMapRenderer) drawAdjacentRooms(screen *ebiten.Image, currentRoom *DiscoveredRoom, allRooms map[string]*DiscoveredRoom) {
	connections := mmr.worldMap.GetRoomConnections(currentRoom.ZoneID)
	
	centerX := mmr.position.X + mmr.size/2
	centerY := mmr.position.Y + mmr.size/2

	for direction, connectedRoomID := range connections {
		if connectedRoom, exists := allRooms[connectedRoomID]; exists {
			mmr.drawAdjacentRoomOutline(screen, connectedRoom, direction, centerX, centerY)
		}
	}
}

// drawAdjacentRoomOutline draws an outline of an adjacent room
func (mmr *MiniMapRenderer) drawAdjacentRoomOutline(screen *ebiten.Image, room *DiscoveredRoom, direction Direction, currentCenterX, currentCenterY int) {
	// Calculate adjacent room position relative to current room center
	roomWidth := float64(room.Bounds.Dx()) * mmr.scale * 0.3  // Smaller scale for adjacent rooms
	roomHeight := float64(room.Bounds.Dy()) * mmr.scale * 0.3

	var roomCenterX, roomCenterY float32

	switch direction {
	case North:
		roomCenterX = float32(currentCenterX)
		roomCenterY = float32(currentCenterY) - float32(mmr.size/4)
	case South:
		roomCenterX = float32(currentCenterX)
		roomCenterY = float32(currentCenterY) + float32(mmr.size/4)
	case East:
		roomCenterX = float32(currentCenterX) + float32(mmr.size/4)
		roomCenterY = float32(currentCenterY)
	case West:
		roomCenterX = float32(currentCenterX) - float32(mmr.size/4)
		roomCenterY = float32(currentCenterY)
	default:
		return
	}

	// Draw room outline
	x := roomCenterX - float32(roomWidth/2)
	y := roomCenterY - float32(roomHeight/2)
	w := float32(roomWidth)
	h := float32(roomHeight)

	outlineColor := color.RGBA{150, 150, 150, 100} // Light gray outline
	vector.StrokeLine(screen, x, y, x+w, y, 1, outlineColor, false)         // Top
	vector.StrokeLine(screen, x, y, x, y+h, 1, outlineColor, false)         // Left
	vector.StrokeLine(screen, x+w, y, x+w, y+h, 1, outlineColor, false)     // Right
	vector.StrokeLine(screen, x, y+h, x+w, y+h, 1, outlineColor, false)     // Bottom
}

// drawPlayerTrail draws the recent player movement trail
func (mmr *MiniMapRenderer) drawPlayerTrail(screen *ebiten.Image, currentRoom *DiscoveredRoom) {
	trail := mmr.worldMap.GetPlayerTrail()
	if len(trail) < 2 {
		return
	}

	centerX := float32(mmr.position.X + mmr.size/2)
	centerY := float32(mmr.position.Y + mmr.size/2)

	// Convert room coordinates to mini-map coordinates
	roomWidth := float64(currentRoom.Bounds.Dx())
	roomHeight := float64(currentRoom.Bounds.Dy())

	for i := 1; i < len(trail); i++ {
		// Convert world coordinates to mini-map coordinates
		prevX := centerX + float32((float64(trail[i-1].X)/roomWidth-0.5)*float64(mmr.size)*0.8)
		prevY := centerY + float32((float64(trail[i-1].Y)/roomHeight-0.5)*float64(mmr.size)*0.8)
		currX := centerX + float32((float64(trail[i].X)/roomWidth-0.5)*float64(mmr.size)*0.8)
		currY := centerY + float32((float64(trail[i].Y)/roomHeight-0.5)*float64(mmr.size)*0.8)

		// Draw trail line with fading alpha
		alpha := uint8(float64(i) / float64(len(trail)) * 100)
		trailColor := color.RGBA{255, 255, 0, alpha}
		vector.StrokeLine(screen, prevX, prevY, currX, currY, 1, trailColor, false)
	}
}

// drawPlayer draws the player position indicator
func (mmr *MiniMapRenderer) drawPlayer(screen *ebiten.Image, player *entities.Player, currentRoom *DiscoveredRoom) {
	if player == nil {
		return
	}

	centerX := float32(mmr.position.X + mmr.size/2)
	centerY := float32(mmr.position.Y + mmr.size/2)

	// Get player world position
	playerWorldX, playerWorldY := player.GetPosition()

	// Convert to mini-map coordinates
	roomWidth := float64(currentRoom.Bounds.Dx())
	roomHeight := float64(currentRoom.Bounds.Dy())

	// Normalize player position within room (0.0 to 1.0)
	normalizedX := float64(playerWorldX) / roomWidth
	normalizedY := float64(playerWorldY) / roomHeight

	// Convert to mini-map pixel coordinates
	playerX := centerX + float32((normalizedX-0.5)*float64(mmr.size)*0.8)
	playerY := centerY + float32((normalizedY-0.5)*float64(mmr.size)*0.8)

	// Draw player as a circle
	vector.DrawFilledCircle(screen, playerX, playerY, 4, mmr.playerColor, false)

	// Draw player facing direction as a small line
	facingRight := player.IsFacingRight()
	dirLength := float32(6)
	var endX float32
	if facingRight {
		endX = playerX + dirLength
	} else {
		endX = playerX - dirLength
	}

	vector.StrokeLine(screen, playerX, playerY, endX, playerY, 2, mmr.playerColor, false)
}