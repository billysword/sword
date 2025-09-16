# Sword

## Project Goals
- Build a playable 2D sword-platformer prototype using [Ebiten](https://ebitengine.org/) with responsive movement and combat-ready systems.
- Experiment with tile-based rooms, world traversal, and state-driven gameplay loops suitable for future content expansion.
- Provide a flexible codebase where rendering, input, and debug tooling are easy to iterate on during early development.

## Prerequisites
- [Go](https://go.dev/dl/) **1.24** or newer (the module targets Go 1.24.1).
- Git for cloning and managing your workspace.
- (Optional) [Tiled](https://www.mapeditor.org/) if you plan to edit the bundled `.tmx` maps or the `sword.tiled-project` configuration.

## Build & Run
1. Clone the repository and install Go dependencies:
   ```bash
   git clone <repo-url>
   cd sword
   go mod download
   ```
2. Run the game straight from source:
   ```bash
   go run main.go
   ```
   Add `--placeholders` to swap in the debug sprites defined in `engine/placeholder_sprites.go` while iterating on gameplay art.
3. (Optional) Build a standalone binary:
   ```bash
   go build -o sword
   ./sword
   ```
4. Execute the automated checks before submitting changes:
   ```bash
   go test ./...
   ```

## Module Layout
- **`engine/`** – Core engine services such as configuration, sprite management, camera control, HUD/debug tooling, and the state manager that orchestrates gameplay loops.
- **`entities/`** – Game actors including the player implementation, reusable enemy logic, and shared collision handling used by characters and hazards.
- **`world/`** – Tile-map orchestration, room transitions, minimap overlays, and helpers for loading `tmx` data into runtime structures.
- **`states/`** – High-level game states (start menu, in-game action, pause/settings, world map, debug views) that drive scene transitions through the engine’s state manager.

## Contributor Guide
1. Fork the repository and work on a feature branch; keep commits focused and descriptive.
2. Follow Go conventions—run `go fmt ./...` before committing and ensure files compile without warnings.
3. Validate behaviour with `go test ./...` and, when relevant, exercise the game via `go run main.go` to manually verify changes.
4. Update documentation (including this README) when you add new systems or assets so newcomers can stay oriented.
5. Open a pull request that explains **what** changed and **why**, noting any assets or tooling requirements reviewers need to reproduce the results.
