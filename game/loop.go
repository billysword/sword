package game

import (
    "bytes"
    "context"
    "image"
    _ "image/png"
    "os"
    "os/signal"
    "syscall"

    "github.com/hajimehoshi/ebiten/v2"
    "sword/engine"
    "sword/states"
)

// Game wraps the ebiten.Game implementation and owns the main loop delegates.
type Game struct {
    stateManager *engine.StateManager
    ctx          context.Context
    cancel       context.CancelFunc
}

// New creates a new Game with cancellation context.
func New() *Game {
    ctx, cancel := context.WithCancel(context.Background())
    return &Game{ctx: ctx, cancel: cancel}
}

// Update implements ebiten.Game.Update() and delegates to the current state.
func (g *Game) Update() error {
    select {
    case <-g.ctx.Done():
        return ebiten.Termination
    default:
    }

    if g.stateManager == nil {
        g.stateManager = engine.NewStateManager()

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
            engine.LogInfo("Loaded spritesheet: " + cfg.Name)
        }

        // Determine starting state (save/load hook ready)
        g.stateManager.ChangeState(states.NewInGameState(g.stateManager))
    }
    return g.stateManager.Update()
}

// Draw implements ebiten.Game.Draw().
func (g *Game) Draw(screen *ebiten.Image) {
    if g.stateManager != nil {
        g.stateManager.Draw(screen)
    }
}

// Layout implements ebiten.Game.Layout().
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ebiten.WindowSize()
}

// Shutdown requests the game loop to terminate.
func (g *Game) Shutdown() { if g.cancel != nil { g.cancel() } }

// Run configures the window, sets up signals, and starts the ebiten loop.
func Run(usePlaceholders bool) error {
    if err := engine.InitLogger("game.log"); err != nil {
        return err
    }

    game := New()

    // Handle OS signals for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        engine.LogInfo("Received interrupt signal, shutting down gracefully...")
        game.Shutdown()
    }()

    // Ensure cleanup regardless of exit path
    defer func() {
        engine.LogInfo("Game cleanup starting...")
        game.Shutdown()
        if err := engine.CloseLogger(); err != nil {
            os.Stderr.WriteString("Error closing logger: " + err.Error() + "\n")
        }
    }()

    // Configure window
    engine.LogInfo("Game starting up...")
    config := engine.DefaultConfig()
    if usePlaceholders {
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
            return nil
        }
        engine.LogInfo("Game ended with error: " + err.Error())
        return err
    }

    engine.LogInfo("Game ended normally")
    return nil
}

