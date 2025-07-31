# Forest Tilemap Platformer Demo

A 2D platformer game built with Ebitengine featuring a forest-themed tilemap system.

## Features

- **Forest Tilemap System**: Uses `forest-tiles.png` as the main tilemap with individual tile extraction
- **Varied Terrain**: Multiple tile types including grass, dirt, stone, trees, flowers, and moss stones
- **Room-based Architecture**: Modular room system for easy level design
- **Character Controls**: Smooth character movement with jumping mechanics
- **Background Trees**: Decorative background elements using tree tiles

## Controls

- **Movement**: A/D or Arrow Keys (Left/Right)
- **Jump**: Spacebar
- **Room Switch**: R key (demo feature)
- **Pause**: P or Escape

## Building and Running

### Prerequisites

Install required system dependencies:

```bash
sudo apt-get update
sudo apt-get install -y libx11-dev libxrandr-dev libxi-dev libgl1-mesa-dev \
                       libglu1-mesa-dev libasound2-dev libxcursor-dev \
                       libxinerama-dev libxxf86vm-dev
```

### Build

```bash
go build -o demo
```

### Run

```bash
./demo
```

## Project Structure

- `main.go` - Entry point and sprite loading
- `gamestate/` - Game state management and room system
  - `room.go` - Base room interface and tilemap system
  - `simple_room.go` - Forest-themed demo room implementation
  - `ingame_state.go` - Main gameplay state
  - `state.go` - State management framework
- `resources/images/` - Game assets
  - `forest-tiles.png` - Main forest tilemap (16x16 tiles)
  - `platformer/` - Character sprites and background

## Tilemap Details

The forest tilemap uses a grid-based system with the following tile types:

- **TILE_GRASS**: Basic ground surface
- **TILE_DIRT**: Underground layers
- **TILE_STONE**: Deep underground and platform tiles
- **TILE_TREE_TOP**: Tree canopy elements
- **TILE_TREE_TRUNK**: Tree trunk sections
- **TILE_FLOWER**: Decorative surface elements
- **TILE_MOSS_STONE**: Textured platform tiles
- **TILE_DARK_DIRT**: Underground variation

Each tile is 16x16 pixels, extracted from the main tilemap for individual use.

## Recent Changes

- Replaced generic tile system with forest-themed tilemap
- Implemented individual tile extraction from `forest-tiles.png`
- Enhanced terrain generation with varied tile types
- Added background tree elements for depth
- Improved visual variety in ground and underground layers