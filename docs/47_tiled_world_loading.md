# Tiled World Loading (Initial Integration)

## Overview

The game can now load rooms directly from Tiled maps exported as JSON (`.tmj`) in `data/zones/<zone>/`. This is the first step toward migrating from hardcoded demo rooms to data-driven content.

## What’s Implemented

- New `TiledRoom` adapter implements the `world.Room` interface from a parsed Tiled map
  - File: `world/tiled_room.go`
  - Builds a `world.TileMap` from the TMJ `render` layer (converts GIDs to local tile indices)
- Zone loader to scan and register rooms and transitions
  - File: `world/zone_loader.go`
  - Function: `LoadZoneRoomsFromData(rtm, zoneName, baseDir)`
  - Loads all `.tmj` files under `data/zones/<zoneName>`
  - Parses external `.tsx` tilesets referenced by the maps
  - Extracts portal objects from the `portals` object layer and registers transitions
    - Portal object properties: `toZone`, `toRoom`, `toPortal`
    - Portal object `name` is used to infer direction: `left|right|up|down`
- Game state integration
  - `states/ingame_state.go` attempts to load the `cradle` zone on startup
  - If data is missing, it falls back to the existing simple demo rooms

## Data Layout

- Zone maps: `data/zones/<zoneName>/*.tmj`
- Tilesets: `data/tilesets/*.tsx`
- Example: `data/zones/cradle/r01.tmj` referencing `../../tilesets/cavern.tsx`

## How It Works at Runtime

1. Startup tries to load `cradle` via `LoadZoneRoomsFromData(...)`
2. A starting room is selected (prefers `<zone>/r01`, otherwise the first found)
3. All rooms from the zone are registered in the `RoomTransitionManager`
4. Transitions are built from `portals` objects
5. Spawn falls back to `main_spawn` from legacy JSON or room center if not defined

## Notes and Limitations

- Collision: Tiled collision layer is parsed by the internal loader and can be used for future physics improvements, but room-level collision still relies on existing helpers.
- Spawns: Spawn points are currently still sourced from embedded JSON (fallback). Tiled-based spawns can be added in a follow-up.
- Directions: Portal `name` infers direction for the world map graph. Use `left|right|up|down` names in Tiled for best results.

## Next Steps

- Read spawn points from a Tiled object layer (e.g., `spawns`) and set `facing_id`.
- Add per-tile collision using Tiled properties (already parsed in `internal/tiled/loader.go`).
- Add zone selection/config for startup instead of hardcoding `cradle`.