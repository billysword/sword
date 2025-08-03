# Forest Tilemap Platformer

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
- **Room Switch**: R key
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
go build -o game
```

### Run

```bash
./game
```

## Project Structure

- `main.go` - Entry point and sprite loading
- `gamestate/` - Game state management and room system
  - `room.go` - Base room interface and tilemap system
  - `simple_room.go` - Forest-themed room implementation
  - `ingame_state.go` - Main gameplay state
  - `state.go` - State management framework
- `resources/images/` - Game assets
  - `forest-tiles.png` - Main forest tilemap (16x16 tiles)
  - `platformer/` - Character sprites and background

## Tilemap Details

The forest tilemap uses a grid-based system with 24 different tile types (indices 0-23):

**Basic Tiles:**
- **0: Dirt** - Underground fill and base material
- **20: Floor 1** - Primary ground surface  
- **21: Floor 2** - Alternate ground surface

**Corners:**
- **1: Top left corner** - Platform top-left edge
- **4: Bottom left corner** - Platform bottom-left edge  
- **5: Top right corner** - Platform top-right edge
- **23: Bottom right corner** - Platform bottom-right edge

**Walls:**
- **2: Right wall 1** - Right edge variation 1
- **3: Right wall 2** - Right edge variation 2
- **6: Left wall 1** - Left edge variation 1
- **22: Left wall 2** - Left edge variation 2

**Ceiling/Top Elements:**
- **7: Ceiling 1** - Ceiling variation 1
- **8: Ceiling 2** - Ceiling variation 2

**Single Tiles:**
- **9: Single tile top** - Isolated top element
- **10: Single tile bottom** - Isolated bottom element
- **11: Single tile left** - Platform left edge
- **12: Single tile right** - Platform right edge
- **13: Floating tile** - Independent floating platform
- **14: Single tile horizontal** - Single horizontal platform
- **15: Single tile vertical** - Single vertical element

**Inner Corners:**
- **16: Inner corner top left**
- **17: Inner corner top right** 
- **18: Inner corner bottom right**
- **19: Inner corner bottom left**

Each tile is 16x16 pixels, extracted from the main tilemap for individual use.

## Recent Changes

- Replaced generic tile system with `forest-tiles.png` tilemap
- Implemented individual tile extraction for all 24 tile types (indices 0-23)
- Created proper platform generation with edges using single tile left/right
- Added floating platforms with appropriate tile types
- Enhanced terrain generation using Floor 1/Floor 2 tiles for surface variety
- Underground areas filled with Dirt tiles (index 0)
- Platforms now use proper edge tiles for visual consistency
- **Added interface-based enemy system with customizable AI behaviors**
- **Organized all documentation into [docs/](docs/) folder with proper indexing**

## Additional Documentation

📚 **All documentation has been organized in the [docs/](docs/) folder.**

For comprehensive documentation, visit the **[Documentation Index](docs/index.md)**.

### Quick Links:
- [Configuration System](docs/01_config_usage.md) - Configuration options and usage
- [Metroidvania Features](docs/02_metroidvania_changes.md) - Game features and mechanics  
- [Enemy System](docs/enemies_interface_02.md) - Interface-based enemy architecture
- [Development Notes](docs/03_claude_notes.md) - AI assistance and development history