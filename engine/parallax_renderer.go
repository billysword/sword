package engine

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

/*
ParallaxRenderer handles multi-layer background rendering with depth effects.
Supports multiple parallax layers with different scroll speeds, transparency,
and depth-of-field effects for enhanced visual depth.
*/
type ParallaxRenderer struct {
	layers       []ParallaxLayer
	layerImages  []*ebiten.Image
	depthOfField bool
	blurStrength float64
	screenWidth  int
	screenHeight int
}

/*
NewParallaxRenderer creates a new parallax renderer with the specified layers.
Initializes the renderer with the current window dimensions and depth settings.
*/
func NewParallaxRenderer(layers []ParallaxLayer, enableDepthOfField bool, blurStrength float64) *ParallaxRenderer {
	pr := &ParallaxRenderer{
		layers:       layers,
		layerImages:  make([]*ebiten.Image, len(layers)),
		depthOfField: enableDepthOfField,
		blurStrength: blurStrength,
	}

	pr.updateScreenSize()
	pr.loadLayerImages()

	return pr
}

/*
updateScreenSize updates the renderer's screen dimensions for proper scaling.
Should be called when the window is resized.
*/
func (pr *ParallaxRenderer) updateScreenSize() {
	pr.screenWidth, pr.screenHeight = ebiten.WindowSize()
}

/*
loadLayerImages loads the images for each parallax layer.
For now, uses the background image as fallback until proper assets are loaded.
*/
func (pr *ParallaxRenderer) loadLayerImages() {
	// For demo purposes, use the existing background image for all layers
	// In a real implementation, you would load different images per layer
	backgroundImage := GetBackgroundImage()

	for i := range pr.layers {
		if backgroundImage != nil {
			pr.layerImages[i] = backgroundImage
		}
	}
}

/*
DrawParallaxLayers renders all parallax layers with proper depth effects.
Layers are drawn from background to foreground with appropriate transformations.
*/
func (pr *ParallaxRenderer) DrawParallaxLayers(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	// Always render parallax layers regardless of background visibility setting
	// Individual layers can be controlled through their Alpha property

	pr.updateScreenSize()

	// Ensure we have images loaded
	if len(pr.layerImages) != len(pr.layers) {
		pr.loadLayerImages()
	}

	// Draw layers from background to foreground (lowest depth first)
	for i, layer := range pr.layers {
		if pr.layerImages[i] != nil {
			pr.drawLayer(screen, layer, pr.layerImages[i], cameraOffsetX, cameraOffsetY, i)
		}
	}
}

/*
drawLayer renders a single parallax layer with depth effects.
Applies parallax scrolling, depth-based transparency, and scaling.
*/
func (pr *ParallaxRenderer) drawLayer(screen *ebiten.Image, layer ParallaxLayer, layerImage *ebiten.Image, cameraOffsetX, cameraOffsetY float64, layerIndex int) {
	LogDebug(fmt.Sprintf("DRAW_LAYER: ParallaxLayer[%d]", layerIndex))
	// Respect the global background visibility setting for background layers
	// Foreground layers (depth > 0.5) are always drawn
	if !GetBackgroundVisible() && layer.Depth <= 0.5 {
		return
	}

	// Calculate parallax offset
	parallaxOffsetX := cameraOffsetX * layer.Speed
	parallaxOffsetY := cameraOffsetY * layer.Speed

	// Add static offset
	parallaxOffsetX += layer.OffsetX
	parallaxOffsetY += layer.OffsetY

	// Create transformation matrix
	op := &ebiten.DrawImageOptions{}

	// Apply depth-based scaling (closer layers appear larger)
	scale := layer.Scale
	if scale == 0 {
		scale = 0.5 + (layer.Depth * 0.5) // Scale from 0.5 to 1.0 based on depth
	}

	// Ensure minimum scale to cover viewport plus parallax movement
	// Calculate how much the image needs to be scaled to cover the screen plus parallax range
	imageW := float64(layerImage.Bounds().Dx())
	imageH := float64(layerImage.Bounds().Dy())
	screenW := float64(pr.screenWidth)
	screenH := float64(pr.screenHeight)

	// Calculate the maximum parallax movement range (assume camera can move full screen distances)
	maxParallaxOffsetX := screenW * layer.Speed
	maxParallaxOffsetY := screenH * layer.Speed

	// Calculate minimum scale needed to cover screen plus parallax movement with buffer
	buffer := 1.5 // 50% extra buffer to ensure no edges show
	minScaleX := (screenW + 2*maxParallaxOffsetX) / imageW * buffer
	minScaleY := (screenH + 2*maxParallaxOffsetY) / imageH * buffer
	minScale := math.Max(minScaleX, minScaleY)

	// Use the larger of the calculated scale or minimum required scale
	scale = math.Max(scale, minScale)

	op.GeoM.Scale(scale, scale)

	// Center the scaled image and apply parallax translation
	// Account for the increased size when centering
	scaledW := imageW * scale
	scaledH := imageH * scale
	centerOffsetX := (scaledW - screenW) / -2
	centerOffsetY := (scaledH - screenH) / -2

	op.GeoM.Translate(centerOffsetX+parallaxOffsetX, centerOffsetY+parallaxOffsetY)

	// Apply depth-of-field effects if enabled
	if pr.depthOfField {
		pr.applyDepthEffects(op, layer)
	}

	// Apply transparency based on depth and layer settings
	alpha := layer.Alpha
	if alpha == 0 {
		alpha = 0.3 + (layer.Depth * 0.7) // Alpha from 0.3 to 1.0 based on depth
	}

	if alpha < 1.0 {
		// Use ColorM to apply transparency
		var cm colorm.ColorM
		cm.Scale(1, 1, 1, alpha)
		// Convert to colorm.DrawImageOptions
		colorOp := &colorm.DrawImageOptions{}
		colorOp.GeoM = op.GeoM
		colorOp.Blend = op.Blend
		colorOp.Filter = op.Filter

		colorm.DrawImage(screen, layerImage, cm, colorOp)
	} else {
		screen.DrawImage(layerImage, op)
	}
}

/*
applyDepthEffects applies depth-of-field effects to a layer.
Includes blur simulation through transparency and color adjustments.
*/
func (pr *ParallaxRenderer) applyDepthEffects(op *ebiten.DrawImageOptions, layer ParallaxLayer) {
	// Simulate depth blur by adjusting color saturation and contrast
	// Farther layers (lower depth) get more desaturated and blurred
	depthEffect := 1.0 - layer.Depth
	blurAmount := depthEffect * pr.blurStrength

	// Reduce contrast for distant layers
	contrast := 1.0 - (blurAmount * 0.3)

	// Create color matrix for depth effects
	r := contrast
	g := contrast
	b := contrast * (1.0 - blurAmount*0.2) // Slight blue tint for distance

	// Apply subtle color shift for depth
	op.ColorM.Scale(r, g, b, 1.0)

	// Add slight position jitter for blur simulation (very subtle)
	if blurAmount > 0.1 {
		jitter := math.Sin(float64(ebiten.TPS())*0.1) * blurAmount * 0.5
		op.GeoM.Translate(jitter, 0)
	}
}

/*
SetLayers updates the parallax layers and reloads images if necessary.
*/
func (pr *ParallaxRenderer) SetLayers(layers []ParallaxLayer) {
	pr.layers = layers
	pr.layerImages = make([]*ebiten.Image, len(layers))
	pr.loadLayerImages()
}

/*
SetDepthOfField enables or disables depth-of-field effects.
*/
func (pr *ParallaxRenderer) SetDepthOfField(enabled bool, strength float64) {
	pr.depthOfField = enabled
	pr.blurStrength = strength
}
