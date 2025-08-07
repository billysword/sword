package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	_ "image/png"
	"os"
	"os/signal"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/states"
)

/*
Game represents the main game application.
Implements the ebiten.Game interface to provide the core game loop
functionality. Manages the state manager and delegates all game
logic to the appropriate game state.

The Game struct serves as the entry point for Ebitengine and handles:
  - State manager initialization and updates
  - Global sprite setup and management
  - Window lifecycle management
  - Graceful shutdown handling
*/
type Game struct {
	stateManager *engine.StateManager // Manages all game states and transitions
	ctx          context.Context      // Context for cancellation support
	cancel       context.CancelFunc   // Cancel function for graceful shutdown
}

/*
NewGame creates a new game instance with proper context setup.
*/
func NewGame() *Game {
	ctx, cancel := context.WithCancel(context.Background())
	return &Game{
		ctx:    ctx,
		cancel: cancel,
	}
}

/*
Update implements ebiten.Game.Update() for the main game loop.
Called once per frame to handle game logic updates. Initializes
the state manager and game states on first call, then delegates
all subsequent updates to the current game state.

Returns any error from the current game state's update logic.
*/
func (g *Game) Update() error {
	// Check if context is cancelled (shutdown requested)
	select {
	case <-g.ctx.Done():
		return ebiten.Termination
	default:
	}

	if g.stateManager == nil {
		g.stateManager = engine.NewStateManager()
		startState := states.NewStartState(g.stateManager)

		// Initialize sprite manager and load sheets from configuration
		engine.InitSpriteManager()
		sm := engine.GetSpriteManager()

		for _, cfg := range engine.SpriteSheetConfigs {
			img, _, err := image.Decode(bytes.NewReader(cfg.ImageData))
			if err != nil {
				panic(err)
			}
			ebImg := ebiten.NewImageFromImage(img)
			if err := sm.LoadSpriteSheet(cfg.Name, ebImg, cfg.TileWidth, cfg.TileHeight); err != nil {
				panic(err)
			}
		}

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
Shutdown gracefully shuts down the game.
*/
func (g *Game) Shutdown() {
	if g.cancel != nil {
		g.cancel()
	}
}

/*
main is the application entry point.
Sets up the game window and starts the main game loop using Ebitengine.
Configures window properties from the game configuration and handles
any startup errors. Now includes proper signal handling for graceful shutdown.

The function:
 1. Sets up signal handling for graceful shutdown
 2. Initializes game logger with proper cleanup
 3. Retrieves window settings from gamestate.GameConfig
 4. Configures Ebitengine window properties
 5. Starts the game loop with a new Game instance
 6. Handles cleanup on exit or error
*/
func main() {
	// Parse command-line flags
	usePlaceholders := flag.Bool("placeholders", false, "Use placeholder sprites instead of actual sprites")
	flag.Parse()

	// Initialize game logger
	if err := engine.InitLogger("game.log"); err != nil {
		panic(err)
	}

	// Create game instance
	game := NewGame()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start signal handler in a goroutine
	go func() {
		<-sigChan
		engine.LogInfo("Received interrupt signal, shutting down gracefully...")
		game.Shutdown()
	}()

	// Ensure cleanup happens no matter how we exit
	defer func() {
		engine.LogInfo("Game cleanup starting...")
		game.Shutdown()

		// Close logger last
		if err := engine.CloseLogger(); err != nil {
			// Can't log this since logger is closing, just print to stderr
			os.Stderr.WriteString("Error closing logger: " + err.Error() + "\n")
		}
	}()

	// Log game startup
	engine.LogInfo("Game starting up...")

	// Get config for window settings
	config := engine.DefaultConfig()

	// Apply command-line flags
	if *usePlaceholders {
		config.UsePlaceholderSprites = true
		engine.LogInfo("Using placeholder sprites")
	}

	engine.SetConfig(config)

	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle(config.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	engine.LogInfo("Window configured and starting game loop...")

	if err := ebiten.RunGame(game); err != nil {
		if err == ebiten.Termination {
			engine.LogInfo("Game terminated normally")
		} else {
			engine.LogInfo("Game ended with error: " + err.Error())
			panic(err)
		}
	}

	engine.LogInfo("Game ended normally")
}
