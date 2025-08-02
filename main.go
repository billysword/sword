package main

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/resources/images"
	"sword/resources/images/platformer"
	"sword/states"
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

/*
Game represents the main game application.
Implements the ebiten.Game interface to provide the core game loop
functionality. Manages the state manager and delegates all game
logic to the appropriate game state.

The Game struct serves as the entry point for Ebitengine and handles:
  - State manager initialization and updates
  - Global sprite setup and management
  - Window lifecycle management
*/
type Game struct {
	stateManager *engine.StateManager  // Manages all game states and transitions
}

/*
Update implements ebiten.Game.Update() for the main game loop.
Called once per frame to handle game logic updates. Initializes
the state manager and game states on first call, then delegates
all subsequent updates to the current game state.

Returns any error from the current game state's update logic.
*/
func (g *Game) Update() error {
	if g.stateManager == nil {
		g.stateManager = engine.NewStateManager()
		startState := states.NewStartState(g.stateManager)
		
		// Initialize sprite manager and load tile sheets
		engine.InitSpriteManager()
		sm := engine.GetSpriteManager()
		
		// Load the forest tile sheet with proper configuration
		err := sm.LoadSpriteSheet("forest", tileSprite, 16, 16)
		if err != nil {
			panic(err)
		}
		
		// Pass sprites to the state manager for use by game states
		engine.SetGlobalSprites(leftSprite, rightSprite, idleSprite, backgroundImage)
		engine.SetGlobalTileSprites(tileSprite, tilesSprite)
		g.stateManager.ChangeState(startState)
	}
	return g.stateManager.Update()
}

/*
Draw implements ebiten.Game.Draw() for rendering.
Called once per frame to render the game to the screen.
Delegates all rendering to the current game state through
the state manager.

Parameters:
  - screen: The target screen/image to render the game to
*/
func (g *Game) Draw(screen *ebiten.Image) {
	if g.stateManager != nil {
		g.stateManager.Draw(screen)
	}
}

/*
Layout implements ebiten.Game.Layout() for window sizing.
Determines the logical screen size for the game. Currently
returns the actual window size, allowing the game to adapt
to different window sizes and resolutions.

Parameters:
  - outsideWidth: The width of the outside area (not used)
  - outsideHeight: The height of the outside area (not used)

Returns the logical screen width and height in pixels.
*/
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.WindowSize()
}

/*
main is the application entry point.
Sets up the game window and starts the main game loop using Ebitengine.
Configures window properties from the game configuration and handles
any startup errors.

The function:
  1. Retrieves window settings from gamestate.GameConfig
  2. Configures Ebitengine window properties
  3. Starts the game loop with a new Game instance
  4. Panics on any startup errors
*/
func main() {
	// Get config for window settings
	config := engine.GameConfig
	
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
