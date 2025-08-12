package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

/*
RenderLayer represents a single rendering layer in the game.
Each layer can be drawn independently with its own transform.
*/
type RenderLayer int

const (
	LayerPrimaryBackground RenderLayer = iota
	LayerParallaxBackground
	LayerRoomTiles
	LayerEntities
	LayerForeground
	LayerViewportFrame
	LayerHUD
)

/*
ViewportRenderer handles the viewport frame/border rendering.
This creates the black borders when the room is smaller than the screen.
*/
type ViewportRenderer struct {
	screenWidth  int
	screenHeight int
	worldWidth   int
	worldHeight  int
	offsetX      float64
	offsetY      float64
	black        *ebiten.Image
}

/*
NewViewportRenderer creates a new viewport renderer.
*/
func NewViewportRenderer(screenW, screenH int) *ViewportRenderer {
	return &ViewportRenderer{
		screenWidth:  screenW,
		screenHeight: screenH,
	}
}

/*
SetWorldBounds updates the world size for viewport calculations.
*/
func (vr *ViewportRenderer) SetWorldBounds(worldW, worldH int) {
	vr.worldWidth = worldW
	vr.worldHeight = worldH
}

/*
SetOffset updates the camera offset for viewport calculations.
*/
func (vr *ViewportRenderer) SetOffset(offsetX, offsetY float64) {
	vr.offsetX = offsetX
	vr.offsetY = offsetY
}

/*
DrawFrame draws the black viewport frame/borders.
This should be called after all world rendering but before HUD.
*/
func (vr *ViewportRenderer) DrawFrame(screen *ebiten.Image) {
	// Lazily create 1x1 black image for scaling
	if vr.black == nil {
		vr.black = ebiten.NewImage(1, 1)
		vr.black.Fill(color.Black)
	}

	// Treat offset in screen pixels and world bounds in unscaled pixels (physics units):
	// convert world bounds to screen space using the current tile scale factor.
	s := GameConfig.TileScaleFactor
	worldLeft := vr.offsetX
	worldTop := vr.offsetY
	worldRight := worldLeft + float64(vr.worldWidth)*s
	worldBottom := worldTop + float64(vr.worldHeight)*s

	// Left border
	if worldLeft > 0 {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(worldLeft, float64(vr.screenHeight))
		screen.DrawImage(vr.black, opts)
	}

	// Right border
	if worldRight < float64(vr.screenWidth) {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(worldRight, 0)
		width := float64(vr.screenWidth) - worldRight
		opts.GeoM.Scale(width, float64(vr.screenHeight))
		screen.DrawImage(vr.black, opts)
	}

	// Top border
	if worldTop > 0 {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(vr.screenWidth), worldTop)
		screen.DrawImage(vr.black, opts)
	}

	// Bottom border
	if worldBottom < float64(vr.screenHeight) {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, worldBottom)
		height := float64(vr.screenHeight) - worldBottom
		opts.GeoM.Scale(float64(vr.screenWidth), height)
		screen.DrawImage(vr.black, opts)
	}
}