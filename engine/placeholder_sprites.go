package engine

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
PlaceholderType defines the type of placeholder sprite to generate.
*/
type PlaceholderType int

const (
	PlaceholderPlayer PlaceholderType = iota
	PlaceholderEnemy
	PlaceholderTileGround
	PlaceholderTileWall
	PlaceholderTilePlatform
	PlaceholderTileSpike
	PlaceholderTileDecoration
	PlaceholderProjectile
	PlaceholderItem
	PlaceholderBackground
)

/*
PlaceholderGenerator creates simple geometric placeholder sprites.
These are low-fidelity sprites that can be used during development
while waiting for final art assets.
*/
type PlaceholderGenerator struct {
	cache map[string]*ebiten.Image
}

// Global placeholder generator instance
var placeholderGen *PlaceholderGenerator

/*
InitPlaceholderGenerator initializes the global placeholder generator.
*/
func InitPlaceholderGenerator() {
	placeholderGen = &PlaceholderGenerator{
		cache: make(map[string]*ebiten.Image),
	}
}

/*
GetPlaceholderGenerator returns the global placeholder generator instance.
*/
func GetPlaceholderGenerator() *PlaceholderGenerator {
	if placeholderGen == nil {
		InitPlaceholderGenerator()
	}
	return placeholderGen
}

/*
GeneratePlaceholder creates a placeholder sprite of the specified type and size.
Uses caching to avoid regenerating the same sprites.
*/
func (pg *PlaceholderGenerator) GeneratePlaceholder(pType PlaceholderType, width, height int) *ebiten.Image {
	// Create cache key
	cacheKey := fmt.Sprintf("%d_%dx%d", pType, width, height)
	
	// Check cache
	if sprite, exists := pg.cache[cacheKey]; exists {
		return sprite
	}
	
	// Create new image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Generate based on type
	switch pType {
	case PlaceholderPlayer:
		pg.drawPlayer(img)
	case PlaceholderEnemy:
		pg.drawEnemy(img)
	case PlaceholderTileGround:
		pg.drawGroundTile(img)
	case PlaceholderTileWall:
		pg.drawWallTile(img)
	case PlaceholderTilePlatform:
		pg.drawPlatformTile(img)
	case PlaceholderTileSpike:
		pg.drawSpikeTile(img)
	case PlaceholderTileDecoration:
		pg.drawDecorationTile(img)
	case PlaceholderProjectile:
		pg.drawProjectile(img)
	case PlaceholderItem:
		pg.drawItem(img)
	case PlaceholderBackground:
		pg.drawBackground(img)
	default:
		pg.drawDefault(img)
	}
	
	// Convert to ebiten image
	sprite := ebiten.NewImageFromImage(img)
	
	// Cache and return
	pg.cache[cacheKey] = sprite
	return sprite
}

/*
drawPlayer draws a simple player character placeholder.
*/
func (pg *PlaceholderGenerator) drawPlayer(img *image.RGBA) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Body (blue rectangle)
	bodyColor := color.RGBA{64, 128, 255, 255}
	bodyRect := image.Rect(width/4, height/3, 3*width/4, 4*height/5)
	draw.Draw(img, bodyRect, &image.Uniform{bodyColor}, image.Point{}, draw.Src)
	
	// Head (circle-ish square)
	headColor := color.RGBA{255, 220, 177, 255}
	headRect := image.Rect(3*width/8, height/6, 5*width/8, height/3)
	draw.Draw(img, headRect, &image.Uniform{headColor}, image.Point{}, draw.Src)
	
	// Eyes (two small dots)
	eyeColor := color.RGBA{0, 0, 0, 255}
	leftEye := image.Rect(5*width/12, height/5, 5*width/12+2, height/5+2)
	rightEye := image.Rect(7*width/12-2, height/5, 7*width/12, height/5+2)
	draw.Draw(img, leftEye, &image.Uniform{eyeColor}, image.Point{}, draw.Src)
	draw.Draw(img, rightEye, &image.Uniform{eyeColor}, image.Point{}, draw.Src)
}

/*
drawEnemy draws a simple enemy placeholder.
*/
func (pg *PlaceholderGenerator) drawEnemy(img *image.RGBA) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Body (red diamond shape)
	bodyColor := color.RGBA{255, 64, 64, 255}
	cx, cy := width/2, height/2
	
	// Draw diamond using triangles
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Check if point is inside diamond
			dx := float64(x - cx)
			dy := float64(y - cy)
			if math.Abs(dx) + math.Abs(dy) <= float64(min(width, height))/2.5 {
				img.Set(x, y, bodyColor)
			}
		}
	}
	
	// Eyes (angry looking)
	eyeColor := color.RGBA{255, 255, 0, 255}
	leftEye := image.Rect(width/3, height/3, width/3+3, height/3+3)
	rightEye := image.Rect(2*width/3-3, height/3, 2*width/3, height/3+3)
	draw.Draw(img, leftEye, &image.Uniform{eyeColor}, image.Point{}, draw.Src)
	draw.Draw(img, rightEye, &image.Uniform{eyeColor}, image.Point{}, draw.Src)
}

/*
drawGroundTile draws a simple ground tile placeholder.
*/
func (pg *PlaceholderGenerator) drawGroundTile(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Base color (brown)
	baseColor := color.RGBA{139, 90, 43, 255}
	draw.Draw(img, bounds, &image.Uniform{baseColor}, image.Point{}, draw.Src)
	
	// Add some texture lines
	lineColor := color.RGBA{101, 67, 33, 255}
	for i := 0; i < 3; i++ {
		y := bounds.Min.Y + (i+1)*bounds.Dy()/4
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, lineColor)
		}
	}
}

/*
drawWallTile draws a simple wall tile placeholder.
*/
func (pg *PlaceholderGenerator) drawWallTile(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Base color (gray)
	baseColor := color.RGBA{128, 128, 128, 255}
	draw.Draw(img, bounds, &image.Uniform{baseColor}, image.Point{}, draw.Src)
	
	// Add brick pattern
	brickColor := color.RGBA{96, 96, 96, 255}
	brickHeight := bounds.Dy() / 4
	brickWidth := bounds.Dx() / 2
	
	for row := 0; row < 4; row++ {
		offset := 0
		if row%2 == 1 {
			offset = brickWidth / 2
		}
		for col := -1; col < 3; col++ {
			x := bounds.Min.X + col*brickWidth + offset
			y := bounds.Min.Y + row*brickHeight
			// Draw vertical lines
			for dy := 0; dy < brickHeight; dy++ {
				if x >= bounds.Min.X && x < bounds.Max.X {
					img.Set(x, y+dy, brickColor)
				}
			}
		}
		// Draw horizontal lines
		y := bounds.Min.Y + row*brickHeight
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, brickColor)
		}
	}
}

/*
drawPlatformTile draws a simple platform tile placeholder.
*/
func (pg *PlaceholderGenerator) drawPlatformTile(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Platform color (light brown)
	platformColor := color.RGBA{205, 133, 63, 255}
	platformRect := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+bounds.Dy()/3)
	draw.Draw(img, platformRect, &image.Uniform{platformColor}, image.Point{}, draw.Src)
	
	// Edge highlight
	highlightColor := color.RGBA{222, 184, 135, 255}
	highlightRect := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+2)
	draw.Draw(img, highlightRect, &image.Uniform{highlightColor}, image.Point{}, draw.Src)
}

/*
drawSpikeTile draws a simple spike tile placeholder.
*/
func (pg *PlaceholderGenerator) drawSpikeTile(img *image.RGBA) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Spike color (dark gray)
	spikeColor := color.RGBA{64, 64, 64, 255}
	
	// Draw three triangular spikes
	spikes := 3
	spikeWidth := width / spikes
	
	for i := 0; i < spikes; i++ {
		baseX := i * spikeWidth
		
		// Draw triangle
		for y := 0; y < height; y++ {
			for x := baseX; x < baseX+spikeWidth; x++ {
				// Check if point is inside triangle
				relX := x - baseX
				maxWidth := (height - y) * spikeWidth / height / 2
				if relX >= spikeWidth/2-maxWidth && relX <= spikeWidth/2+maxWidth {
					img.Set(x, y, spikeColor)
				}
			}
		}
	}
}

/*
drawDecorationTile draws a simple decoration tile placeholder.
*/
func (pg *PlaceholderGenerator) drawDecorationTile(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Draw a simple flower/plant
	stemColor := color.RGBA{34, 139, 34, 255}
	flowerColor := color.RGBA{255, 192, 203, 255}
	
	// Stem
	stemX := bounds.Min.X + bounds.Dx()/2
	for y := bounds.Min.Y + bounds.Dy()/2; y < bounds.Max.Y; y++ {
		img.Set(stemX, y, stemColor)
		img.Set(stemX-1, y, stemColor)
	}
	
	// Flower (simple circle)
	cx, cy := bounds.Min.X+bounds.Dx()/2, bounds.Min.Y+bounds.Dy()/3
	radius := bounds.Dx() / 3
	for y := cy - radius; y <= cy + radius; y++ {
		for x := cx - radius; x <= cx + radius; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx + dy*dy <= radius*radius {
				img.Set(x, y, flowerColor)
			}
		}
	}
}

/*
drawProjectile draws a simple projectile placeholder.
*/
func (pg *PlaceholderGenerator) drawProjectile(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Projectile color (yellow)
	projectileColor := color.RGBA{255, 255, 0, 255}
	
	// Draw a circle
	cx, cy := bounds.Dx()/2, bounds.Dy()/2
	radius := min(bounds.Dx(), bounds.Dy()) / 3
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx + dy*dy <= radius*radius {
				img.Set(x, y, projectileColor)
			}
		}
	}
}

/*
drawItem draws a simple item placeholder.
*/
func (pg *PlaceholderGenerator) drawItem(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Clear with transparent
	draw.Draw(img, bounds, &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	// Item color (gold/yellow)
	itemColor := color.RGBA{255, 215, 0, 255}
	
	// Draw a star shape
	cx, cy := bounds.Dx()/2, bounds.Dy()/2
	radius := min(bounds.Dx(), bounds.Dy()) / 3
	
	// Simple 4-pointed star
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			// Check if point is in star shape (simplified)
			if (math.Abs(dx) < float64(radius)/3 && math.Abs(dy) < float64(radius)) ||
			   (math.Abs(dy) < float64(radius)/3 && math.Abs(dx) < float64(radius)) {
				img.Set(x, y, itemColor)
			}
		}
	}
}

/*
drawBackground draws a simple background placeholder.
*/
func (pg *PlaceholderGenerator) drawBackground(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Gradient from light blue to darker blue
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		ratio := float64(y-bounds.Min.Y) / float64(bounds.Dy())
		r := uint8(135 - 35*ratio)
		g := uint8(206 - 56*ratio) 
		b := uint8(235 - 35*ratio)
		lineColor := color.RGBA{r, g, b, 255}
		
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, lineColor)
		}
	}
}

/*
drawDefault draws a simple default placeholder.
*/
func (pg *PlaceholderGenerator) drawDefault(img *image.RGBA) {
	bounds := img.Bounds()
	
	// Purple checkerboard pattern
	color1 := color.RGBA{255, 0, 255, 255}
	color2 := color.RGBA{128, 0, 128, 255}
	
	checkSize := 8
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if ((x/checkSize)+(y/checkSize))%2 == 0 {
				img.Set(x, y, color1)
			} else {
				img.Set(x, y, color2)
			}
		}
	}
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Convenience functions for common sizes

/*
GeneratePlayerPlaceholder creates a standard player placeholder sprite.
*/
func GeneratePlayerPlaceholder() *ebiten.Image {
	return GetPlaceholderGenerator().GeneratePlaceholder(PlaceholderPlayer, 32, 32)
}

/*
GenerateEnemyPlaceholder creates a standard enemy placeholder sprite.
*/
func GenerateEnemyPlaceholder() *ebiten.Image {
	return GetPlaceholderGenerator().GeneratePlaceholder(PlaceholderEnemy, 32, 32)
}

/*
GenerateTilePlaceholder creates a standard tile placeholder sprite.
*/
func GenerateTilePlaceholder(tileType PlaceholderType) *ebiten.Image {
	return GetPlaceholderGenerator().GeneratePlaceholder(tileType, 16, 16)
}