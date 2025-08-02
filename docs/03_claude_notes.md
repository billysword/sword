# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based 2D platformer game built using the Ebitengine (formerly Ebiten) game engine. The project follows a modular game state delegate pattern:

- `main.go`: Minimal entry point that initializes sprite loading and delegates to StateManager
- `gamestate/`: Package containing all game state logic and the state management system
- `resources/`: Asset directory with embedded images organized by category
- `go.mod`: Module definition with Ebitengine v2.8.8 as the primary dependency

## Build and Run Commands

```bash
# Run the game
go run main.go

# Build the executable
go build -o sword main.go

# Get dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Architecture

The game uses a **State Delegate Pattern** for modular, maintainable code:

### Core Components

- **`main.Game`**: Minimal struct that delegates all logic to `gamestate.StateManager`
- **`gamestate.StateManager`**: Handles state transitions and delegates Update/Draw calls
- **`gamestate.State` interface**: Defines `Update()`, `Draw()`, `OnEnter()`, `OnExit()` methods
- **Global Sprite Management**: Sprites loaded once in main, accessed via `gamestate.SetGlobalSprites()`

### Game States

1. **`StartState`** (`gamestate/start_state.go`): Main menu with controls and title
2. **`InGameState`** (`gamestate/ingame_state.go`): Core gameplay with character physics
3. **`PauseState`** (`gamestate/pause_state.go`): Overlay pause menu with resume/quit options

### State Flow
- Start → InGame (ENTER/SPACE)
- InGame → Pause (P/ESC)  
- Pause → InGame (P/ESC) or Pause → Start (Q)

### Key Constants
- Screen dimensions: 960x540 pixels
- Physics unit: 16 pixels  
- Ground level: y=380 (in units)

## Asset Management

Images are embedded using Go's embed package in two locations:
- `resources/images/embed.go`: General game assets
- `resources/images/platformer/embed.go`: Platformer-specific character sprites

The embed pattern uses exported byte slice variables that need to be loaded into `*ebiten.Image` objects during initialization.

## Game Mechanics

### Character Physics (`gamestate.Character`)
- WASD/Arrow key movement with momentum-based physics
- Space bar for jumping (can jump mid-air)
- Gravity and friction applied each frame in `Character.update()`
- Character sprite changes based on movement direction

### Controls
- **Start Screen**: ENTER/SPACE to begin game
- **In-Game**: WASD/Arrows to move, SPACE to jump, P/ESC to pause
- **Pause Menu**: P/ESC to resume, Q to return to main menu

## Development Guidelines

### Adding New Game States
1. Create new file in `gamestate/` package
2. Implement the `State` interface (`Update`, `Draw`, `OnEnter`, `OnExit`)
3. Use `stateManager.ChangeState()` for transitions
4. Access sprites via global variables (`globalLeftSprite`, etc.)

### Modular Design Principles
- Keep states small and focused on single responsibilities
- Each state handles its own input, logic, and rendering
- Use `OnEnter`/`OnExit` for state-specific setup/cleanup
- Global sprite access prevents redundant asset loading