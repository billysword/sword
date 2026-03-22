# AGENTS.md

## Repository Overview
- **Project:** `sword`
- **Purpose:** 2D sword-platformer prototype built with [Ebiten](https://ebitengine.org/).
- **Entry point:** `main.go` calls `game.Run(...)`.
- **Optional debug flag:** `go run main.go --placeholders` swaps in placeholder sprites for faster iteration on gameplay code.

## Toolchain
- Go **1.24.1** or newer
- Git
- Optional: [Tiled](https://www.mapeditor.org/) for editing `.tmx` maps and `sword.tiled-project`

## Common Commands
- Install dependencies: `go mod download`
- Run the game: `go run main.go`
- Run with placeholder sprites: `go run main.go --placeholders`
- Run tests: `go test ./...`
- Format Go code before committing: `go fmt ./...`

## Repository Layout
- `engine/` - engine services such as config, sprite management, camera logic, HUD/debug tooling, and state orchestration
- `entities/` - player, enemies, and shared collision logic
- `game/` - game bootstrap and runtime wiring
- `states/` - top-level gameplay/menu/debug states
- `world/` - tile maps, rooms, transitions, minimap/world map systems, and related helpers
- `assets/`, `resources/` - art and other runtime assets
- `room_layouts/` - reusable room layout definitions and notes

## Working Agreements
- Keep changes focused and small; avoid broad refactors unless the task requires them.
- Prefer updating the smallest relevant package instead of moving cross-package responsibilities around.
- Follow normal Go conventions for naming, package structure, and error handling.
- Update documentation when you introduce a new gameplay system, asset requirement, or workflow.
- Do not modify art, maps, or tiled project files unless the task specifically involves them.

## Validation Expectations
- Run `go test ./...` for code changes whenever practical.
- Run `go fmt ./...` before finalizing Go changes.
- For gameplay, rendering, camera, input, or state-transition work, manual verification with `go run main.go` is important because automated coverage is limited.
- The existing automated tests are light and currently centered around `world/`, so use extra care when changing runtime behavior in `engine/`, `entities/`, `game/`, or `states/`.

## Change-Specific Notes
- If you add new room layouts, keep them in `room_layouts/` and document any new conventions there when needed.
- If you add flags or startup options, reflect them in `README.md` and this file.
- If you introduce new assets or resource-loading behavior, confirm paths carefully and document any required setup for future contributors.
