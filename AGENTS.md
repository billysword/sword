# AGENTS.md

## Purpose

This repository is evolving from a single sword-platformer prototype into a small workspace for reusable 2D exploration/platforming systems. The immediate target is a **Metroid-style / metroidvania prototype** built on the current Go + Ebiten codebase.

Agents working in this repo should optimize for:

- shipping the next playable vertical slice quickly
- extracting reusable systems only after they prove useful
- keeping engine code separate from game content and assets
- improving test coverage around pure gameplay/world logic

## Current product direction

Build toward a metroidvania loop with:

- interconnected rooms
- strong traversal feel
- ability-gated progression
- combat-capable player movement
- room/world map discovery
- save/checkpoint flow

The current codebase already contains promising foundations in:

- `world/` for rooms, transitions, map discovery, and tiled loading
- `entities/` for player/enemy runtime behavior
- `states/` for screen/state flow
- `engine/` for rendering helpers, sprite management, config, and state management
- `internal/tiled/` for map parsing

## Architectural stance

Do **not** turn this into a large generic engine up front.

Prefer:

1. one repo
2. multiple game-capable workflows over time
3. a shared internal engine boundary
4. extraction driven by real reuse

Avoid:

- adding game assets or content definitions to shared engine packages
- hardcoding game-specific sprite sheet names inside reusable packages
- introducing new abstractions before at least one concrete gameplay use exists
- coupling pure world logic to Ebiten when it can remain data-oriented

## Dependency direction

When adding or refactoring code, move toward this shape:

- `engine/` contains reusable systems and helpers
- `internal/tiled/` stays format/data-loader focused
- `game/`, `states/`, `entities/`, and future game-specific content own concrete game behavior
- asset registration happens outside reusable engine packages

Prefer these rules:

- pure logic should not require `*ebiten.Image`
- room graph and progression logic should be testable without rendering
- logging should not be the reason a pure package imports `engine/`
- content and assets should be referenced by the game layer, not by shared engine code

## Primary metroidvania outcome

The next major outcome is a playable vertical slice with:

1. at least a small connected room set
2. room-to-room traversal with stable camera behavior
3. one progression gate tied to a newly unlocked ability
4. one checkpoint/save loop
5. basic enemy/combat pressure
6. usable minimap/world-map feedback

Agents should prioritize work that gets the project closer to that loop.

## Ordered workstreams

Complete work in roughly this order unless a task explicitly depends on another path being finished first.

### 1. Architecture stabilization

Goal: make the existing codebase easier to extend safely.

Focus areas:

- move sprite sheet registration/config out of `engine/`
- reduce engine assumptions about concrete asset names
- separate room graph / transition data from room rendering/runtime where practical
- add tests around world map, transitions, spawn resolution, and other pure logic

### 2. Traversal foundation

Goal: make movement feel good enough for a metroidvania.

Focus areas:

- tune player acceleration, air control, jump arc, wall interaction, and landing feel
- support traversal-specific collision cases
- improve camera follow behavior, dead zones, and room transition camera handoff

### 3. Room and world progression

Goal: support non-linear exploration with explicit gating.

Focus areas:

- define room metadata and progression locks
- add ability gates and lock/key or unlock-based transitions
- support spawn points, checkpoints, and respawn behavior

### 4. Combat and encounter baseline

Goal: make exploration spaces carry gameplay pressure.

Focus areas:

- reliable attack/hurtbox interactions
- enemy reset/spawn rules by room
- health, damage, invulnerability, and death/restart flow

### 5. Player progression systems

Goal: make new abilities meaningful and persistent.

Focus areas:

- unlockable abilities
- save/load integration for progression state
- UI feedback for acquired upgrades

### 6. Content pipeline and iteration tools

Goal: enable rapid addition of rooms and content.

Focus areas:

- Tiled conventions for portals, spawn points, and gates
- room validation tooling
- map authoring docs
- debug views for transitions, collisions, and progression blockers

## Parallelization guidance

Multiple agents can work in parallel if they avoid stepping on the same files or APIs.

Good parallel task groupings:

- **world graph/data**: transitions, locks, checkpoints, serialization
- **player traversal**: movement tuning, camera, ability physics
- **combat runtime**: enemies, damage, combat state
- **content pipeline**: Tiled metadata, docs, validation helpers
- **UI/map**: minimap, world-map reveal, progression indicators

Before extracting shared interfaces, let one agent land the primary vertical-slice implementation first. Follow-up agents can then generalize repeated patterns.

## Task selection rules for agents

When choosing the next task:

1. prefer the smallest task that unlocks downstream work
2. prefer tasks with visible gameplay impact
3. prefer pure-logic seams for tests
4. avoid wide refactors unless they unblock several queued tasks

If two approaches are possible:

- choose the one that keeps metroidvania-specific behavior in game/runtime packages
- choose the one that keeps reusable data structures testable without Ebiten

## Definition of done for gameplay work

A task is not done until it includes the relevant combination of:

- code changes
- tests for pure logic when feasible
- doc updates when conventions or workflows changed
- manual verification notes for runtime behavior

## Definition of done for refactor work

A refactor is not done unless it clearly improves one of:

- dependency direction
- testability
- asset/content separation
- metroidvania feature delivery speed

Do not perform large “cleanup only” refactors that do not help the target game.

## Files agents should read first

Start with:

- `README.md`
- `docs/metroidvania-roadmap.md`
- `docs/agent-task-breakdown.md`

Then inspect the relevant code:

- `world/` for room graph, transitions, map logic, Tiled integration
- `entities/` for player/enemy runtime behavior
- `game/systems/` for update ordering and game-loop coordination
- `states/` for current game flow
- `engine/` only for reusable runtime helpers and boundaries

## Prompting guidance for future agents

When launching future coding agents, give them:

- the specific metroidvania sub-goal
- the files/packages they should start from
- explicit non-goals
- the expected tests or manual verification steps

Good prompt pattern:

> Implement [small gameplay or architecture task] for the metroidvania vertical slice. Start by reading [specific files]. Keep reusable logic decoupled from game assets. Add tests for pure logic where practical. Do not broaden the abstraction unless the task requires it.
