# Agent Task Breakdown for the Metroidvania Goal

Use this document to split work into small, reviewable tasks that move the repository toward the metroidvania vertical slice.

## Working rules

- keep tasks narrow enough to land in a focused PR
- prefer tasks with one primary package owner
- add tests when the logic can run without rendering
- document new content conventions as they are introduced
- avoid broad engine refactors unless they unblock multiple downstream tasks

## Recommended task sequence

### Wave 1: unblock future work

These tasks should happen first because they reduce coupling or clarify gameplay structure.

#### Task 1. Asset registration boundary

Goal:

- move sprite-sheet/content registration out of `engine/`

Likely files:

- `engine/sprite_config.go`
- `game/loop.go`
- asset embedding packages under `assets/` and `resources/`

Definition of done:

- reusable engine code no longer owns concrete game asset registration
- the game still boots with the same assets

#### Task 2. Transition and spawn test coverage

Goal:

- add tests around room transitions, spawn resolution, and related pure logic

Likely files:

- `world/room_transition.go`
- new `world/*_test.go` files

Definition of done:

- key transition behaviors have regression coverage without depending on full gameplay rendering

#### Task 3. Room graph versus runtime seam

Goal:

- separate room connectivity/progression data from rendering-heavy room logic where feasible

Likely files:

- `world/room_transition.go`
- `world/worldmap.go`
- `world/zone_loader.go`
- `world/tiled_room.go`

Definition of done:

- room connectivity and gate logic can evolve with less dependence on drawing/runtime details

### Wave 2: traversal and feel

#### Task 4. Camera tuning for connected rooms

Goal:

- support metroidvania-style camera behavior across room movement and transitions

Likely files:

- `game/systems/camera_system.go`
- `engine/camera.go`
- room transition integration points in `world/` or `game/systems/`

Definition of done:

- the camera behavior feels stable while traversing connected rooms

#### Task 5. Player traversal baseline

Goal:

- tune movement for exploration-heavy gameplay

Likely files:

- `entities/player.go`
- `entities/player_collision.go`
- `game/systems/physics_system.go`
- `engine/config.go`

Definition of done:

- movement supports the target slice and exposes the next missing traversal capabilities clearly

### Wave 3: progression loop

#### Task 6. Checkpoints and respawn flow

Goal:

- add consistent checkpoint ownership and respawn behavior

Likely files:

- `world/room_transition.go`
- `states/ingame_state.go`
- save/progression helpers if introduced

Definition of done:

- death or reset returns the player to a known checkpoint with predictable state

#### Task 7. Ability gates

Goal:

- support routes that open after the player gains an ability

Likely files:

- `world/`
- `entities/player.go`
- progression-related state and serialization code

Definition of done:

- at least one route changes from blocked to reachable after unlock

#### Task 8. Ability unlock persistence

Goal:

- store and reload progression state cleanly

Likely files:

- progression/save packages or files introduced for this work
- `states/`
- world/progression integration points

Definition of done:

- ability state survives restart/load and affects world traversal

### Wave 4: combat loop

#### Task 9. Combat reliability pass

Goal:

- make combat interactions dependable enough for authored encounters

Likely files:

- `entities/`
- `game/systems/`

Definition of done:

- basic enemy encounters behave predictably under room transitions and player damage rules

### Wave 5: content pipeline and map UX

#### Task 10. Tiled authoring conventions

Goal:

- document how to author portals, checkpoints, gates, and room metadata

Likely files:

- `README.md`
- new docs under `docs/`
- `world/zone_loader.go`
- `internal/tiled/`

Definition of done:

- a contributor can add a room and required metadata without reverse engineering the runtime

#### Task 11. Minimap/world-map progression feedback

Goal:

- improve map usefulness during exploration

Likely files:

- `world/minimap.go`
- `world/map_overlay.go`
- `states/world_map_state.go`

Definition of done:

- map UI reflects visited/current/locked-or-interesting room state clearly enough to guide play

## Parallel task lanes

If multiple agents are running at once, prefer these combinations:

### Lane A: architecture and tests

- Task 1
- Task 2
- Task 3

### Lane B: traversal feel

- Task 4
- Task 5

### Lane C: progression systems

- Task 6
- Task 7
- Task 8

### Lane D: combat and map UX

- Task 9
- Task 11

### Lane E: content pipeline

- Task 10

## Example prompts for future agents

### Prompt: asset boundary

> Refactor asset registration so concrete sprite-sheet wiring is owned by the game layer instead of shared engine code. Start with `engine/sprite_config.go` and `game/loop.go`. Preserve current runtime behavior. Add or update docs if the startup flow changes.

### Prompt: world transition tests

> Add regression coverage for room transitions and spawn placement. Start with `world/room_transition.go` and existing tests in `world/worldmap_test.go`. Prefer pure logic coverage and avoid introducing rendering dependencies when possible.

### Prompt: camera and traversal

> Improve the traversal feel for the metroidvania vertical slice. Inspect `entities/player.go`, `game/systems/physics_system.go`, `game/systems/camera_system.go`, and related config. Focus on movement and room-to-room camera behavior, not combat or map UI.

### Prompt: checkpoints and progression

> Implement checkpoint and ability-gated progression support for the metroidvania vertical slice. Start in `world/`, then integrate with the in-game state flow. Keep progression data structured so it can be serialized later.

## Review checklist for task owners

Before shipping a task, verify:

- does this directly help the metroidvania slice?
- is new reusable code free of game-asset assumptions?
- could the logic be tested without Ebiten, and if so, did we test it?
- did we avoid widening the abstraction unnecessarily?
