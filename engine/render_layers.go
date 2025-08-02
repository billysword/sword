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
	// Calculate visible world bounds
	worldLeft := vr.offsetX
	worldTop := vr.offsetY
	worldRight := worldLeft + float64(vr.worldWidth)
	worldBottom := worldTop + float64(vr.worldHeight)
	
	// Draw black borders if world is smaller than screen
	// Left border
	if worldLeft > 0 {
		opts := &ebiten.DrawImageOptions{}
		blackImg := ebiten.NewImage(int(worldLeft), vr.screenHeight)
		blackImg.Fill(color.Black)
		screen.DrawImage(blackImg, opts)
	}
	
	// Right border
	if worldRight < float64(vr.screenWidth) {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(worldRight, 0)
		width := float64(vr.screenWidth) - worldRight
		blackImg := ebiten.NewImage(int(width), vr.screenHeight)
		blackImg.Fill(color.Black)
		screen.DrawImage(blackImg, opts)
	}
	
	// Top border
	if worldTop > 0 {
		opts := &ebiten.DrawImageOptions{}
		blackImg := ebiten.NewImage(vr.screenWidth, int(worldTop))
		blackImg.Fill(color.Black)
		screen.DrawImage(blackImg, opts)
	}
	
	// Bottom border
	if worldBottom < float64(vr.screenHeight) {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(0, worldBottom)
		height := float64(vr.screenHeight) - worldBottom
		blackImg := ebiten.NewImage(vr.screenWidth, int(height))
		blackImg.Fill(color.Black)
		screen.DrawImage(blackImg, opts)
	}
}