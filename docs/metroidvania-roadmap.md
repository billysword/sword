# Metroidvania Roadmap

## Target

Build a playable metroidvania vertical slice on top of the current repository without turning the codebase into a large monolith.

The target slice should prove:

- room-to-room exploration feels good
- progression can be gated by abilities
- combat and checkpoint loops support sustained play
- the world map/minimap gives useful navigation feedback

## Success criteria for the first vertical slice

The slice should include:

1. a connected set of authored rooms
2. at least one locked route or progression gate
3. one new traversal or exploration ability
4. one checkpoint/save interaction
5. at least one enemy type that pressures movement or combat
6. minimap/world-map feedback showing exploration progress

## Milestone order

### Milestone 0: Stabilize extension seams

Purpose: make the current repo safer to evolve.

Tasks:

- move sprite registration and concrete asset wiring out of `engine/`
- reduce hardcoded shared-engine assumptions about specific sprite sheet names
- isolate room graph / transition data from rendering-heavy room runtime where practical
- add tests for transitions, spawn placement, room discovery, and serialization

Exit criteria:

- agents can extend world/progression logic without pulling in unnecessary rendering concerns
- new assets can be registered without editing shared engine internals

### Milestone 1: Traversal core

Purpose: make movement carry the experience.

Tasks:

- tighten jump arc, landing feel, and air control
- define camera follow behavior for room-based exploration
- support traversal affordances required by the target slice
- validate room transitions while moving at gameplay speed

Exit criteria:

- moving through a room set feels stable and intentional
- camera behavior does not fight the player during traversal

### Milestone 2: World structure and gating

Purpose: create a meaningful exploration loop.

Tasks:

- formalize room metadata and progression flags
- add gateable transitions or blockers
- implement checkpoints and respawn flow
- support “visited”, “current room”, and “reachable after unlock” map state

Exit criteria:

- the player can be redirected by a gate and later return after unlocking a new route

### Milestone 3: Combat baseline

Purpose: make spaces feel inhabited and risky.

Tasks:

- stabilize attack flow and damage handling
- define a basic enemy lifecycle per room
- support combat recovery/death rules compatible with checkpoints

Exit criteria:

- player, enemy, and room systems produce a repeatable encounter loop

### Milestone 4: Ability unlock loop

Purpose: make progression tangible.

Tasks:

- add one ability unlock
- persist the unlocked state
- apply ability state to traversal and gating logic
- surface ability state in HUD or debug UI

Exit criteria:

- unlocking an ability immediately opens or simplifies at least one route

### Milestone 5: Content iteration workflow

Purpose: make it easy to keep adding rooms.

Tasks:

- document Tiled authoring conventions
- add validation helpers for required room metadata
- document how to add portals, checkpoints, gates, and spawn points
- improve debug overlays for room transitions and progression blockers

Exit criteria:

- a new room can be added with minimal code changes and predictable conventions

## Suggested package direction

This is not a mandatory restructure, but new work should lean this way:

- reusable runtime helpers remain in `engine/`
- pure map-loading logic remains in `internal/tiled/`
- room graph, transitions, gates, and serialization become increasingly data-oriented
- game-specific content and asset declarations stay outside shared engine packages

## Risk management

Avoid these traps during the metroidvania buildout:

- over-generalizing systems before they support one real gameplay loop
- making reusable packages depend on asset bundles
- mixing map authoring conventions with rendering details in the same code path
- spending large amounts of effort on cosmetic refactors with no gameplay payoff

## How to choose the next task

If several tasks are available, choose in this order:

1. tasks that unblock later milestones
2. tasks that improve the playable slice directly
3. tasks that reduce coupling in active code paths
4. docs/tooling tasks that improve content iteration
