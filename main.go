package main

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/gamestate"
	"sword/resources/images"
	"sword/resources/images/platformer"
)

const (
	// Settings
	screenWidth  = 640  // Standard metroidvania width
	screenHeight = 360  // Standard metroidvania height (16:9 aspect ratio)
)

var (
	leftSprite      *ebiten.Image
	rightSprite     *ebiten.Image
	idleSprite      *ebiten.Image
	backgroundImage *ebiten.Image
	tileSprite      *ebiten.Image
	tilesSprite     *ebiten.Image
)

func init() {
	// Load background image
	img, _, err := image.Decode(bytes.NewReader(platformer.Background_png))
	if err != nil {
		panic(err)
	}
	backgroundImage = ebiten.NewImageFromImage(img)

	// Load character sprites
	img, _, err = image.Decode(bytes.NewReader(platformer.Left_png))
	if err != nil {
		panic(err)
	}
	leftSprite = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(platformer.Right_png))
	if err != nil {
		panic(err)
	}
	rightSprite = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(platformer.MainChar_png))
	if err != nil {
		panic(err)
	}
	idleSprite = ebiten.NewImageFromImage(img)

	// Load forest tile sprites
	img, _, err = image.Decode(bytes.NewReader(images.ForestTiles_png))
	if err != nil {
		panic(err)
	}
	tileSprite = ebiten.NewImageFromImage(img)
	tilesSprite = tileSprite // Use the same forest tilemap for both tile references
}

type Game struct {
	stateManager *gamestate.StateManager
}

func (g *Game) Update() error {
	if g.stateManager == nil {
		g.stateManager = gamestate.NewStateManager()
		startState := gamestate.NewStartState(g.stateManager)
		// Pass sprites to the state manager for use by game states
		gamestate.SetGlobalSprites(leftSprite, rightSprite, idleSprite, backgroundImage)
		gamestate.SetGlobalTileSprites(tileSprite, tilesSprite)
		g.stateManager.ChangeState(startState)
	}
	return g.stateManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.stateManager != nil {
		g.stateManager.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Platformer (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
