# Hexadecimal Room Layout System

## Overview

The room debugging system has been enhanced to support hexadecimal format output with 0xFF (255) as the maximum value instead of the previous decimal format limited to 99. This provides a more suitable format for game development with easy editing in monospace files and proper array declarations.

## Key Features

### 🔧 **Enhanced Debug Output**
- **Dual Format Support**: Both decimal (legacy) and hexadecimal formats in debug logs
- **Extended Range**: 0x00-0xFF (0-255) instead of 00-99 decimal limit
- **Copy-Paste Ready**: Generated Go array declarations ready for direct code integration

### 📁 **Standalone Layout Files**
- **Auto-generated .go files**: Complete Go files with room layouts
- **Proper Variable Names**: CamelCase variable names based on room IDs
- **Documentation**: Auto-generated comments with dimensions and generation timestamps

### 🎨 **Easy Editing Workflows**
- **Monospace Friendly**: Clean format that aligns perfectly in text editors
- **Visual Clarity**: Hex values clearly distinguish different tile types
- **Industry Standard**: Uses standard game development hex notation

## Usage Examples

### Basic Room Debug
```go
// In your room construction code
func (r *MyRoom) buildRoom() {
    // ... room construction logic ...
    
    // Generate debug output with both decimal and hex formats
    r.LogRoomDebug()
    
    // Create a standalone .go file for easy editing
    r.GenerateHexLayoutFile()
}
```

### Copy-Paste Workflow
```go
// Get the hex layout as a string for immediate use
hexLayout := room.GetHexLayoutArray()
fmt.Println(hexLayout)

// Output: Ready for copy-paste into your code
// levelLayout := [][]int{
//     {-1, -1, 0x5, 0x6, 0x7, -1},
//     {0x1, 0x2, 0x3, 0x1, 0x2, 0x3},
//     {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
// }
```

### Manual Layout Creation
```go
// Create layouts using hex notation for clarity
platformLevel := [][]int{
    {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},           // Sky
    {-1, -1, 0x5, 0x6, 0x7, -1, -1, -1, -1, -1},        // Platform tiles
    {-1, -1, -1, -1, -1, -1, -1, 0xA, 0xB, -1},         // Another platform
    {0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1}, // Ground
    {0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}, // Underground
    {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Bedrock
}
```

## Generated Output Formats

### 1. Debug Log Format
```
=== ROOM DEBUG: example_room ===
Timestamp: 2025-08-02 02:58:11
Room Dimensions: 10x8 tiles

ASCII Representation (2-digit tile indices - decimal):
99,99,05,06,07,99,99,99,99,99
01,02,03,01,02,03,01,02,03,01

ASCII Representation (2-digit tile indices - hexadecimal):
FF,FF,05,06,07,FF,FF,FF,FF,FF
01,02,03,01,02,03,01,02,03,01

Go Array Format (ready for copy-paste):
levelLayout := [][]int{
    {-1, -1, 0x5, 0x6, 0x7, -1, -1, -1, -1, -1},
    {0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1},
}
=== END ROOM DEBUG ===
```

### 2. Standalone Go File
```go
// Auto-generated room layout for: example_room
// Generated: 2025-08-02 02:58:11
// Dimensions: 10x8 tiles

package main

// Example_roomLayout contains the tile layout in hexadecimal format
// -1 = empty tile, 0x00-0xFF = tile indices
var Example_roomLayout = [][]int{
    {-1, -1, 0x5, 0x6, 0x7, -1, -1, -1, -1, -1},
    {0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1},
}
```

## Benefits Over Decimal Format

| Feature | Decimal (Old) | Hexadecimal (New) |
|---------|---------------|-------------------|
| **Range** | 0-99 (limited) | 0x00-0xFF (0-255) |
| **Empty Tiles** | 99 | FF (or -1 in arrays) |
| **Readability** | 01,02,03,99,99 | 01,02,03,FF,FF |
| **Industry Standard** | Custom notation | Standard hex notation |
| **Code Integration** | Manual conversion | Direct copy-paste |
| **Visual Alignment** | Good | Excellent |

## File Locations

- **Debug Logs**: `log/room_debug_YYYY-MM-DD.log`
- **Layout Files**: `log/room_layout_<room_name>.go`
- **Examples**: `examples/hex_layout_example.go`

## Editing Workflows

### Method 1: Direct Array Editing
1. Generate hex layout: `room.GenerateHexLayoutFile()`
2. Open generated `.go` file in editor
3. Edit hex values directly (0x00 to 0xFF range)
4. Copy array back to room code

### Method 2: Log-based Workflow
1. Run game to generate debug logs
2. Copy "Go Array Format" section from log
3. Paste into room construction code
4. Edit hex values as needed

### Method 3: Programmatic Generation
1. Use `room.GetHexLayoutArray()` to get current layout
2. Modify room programmatically
3. Export new layout for manual fine-tuning

## Integration with Existing Code

The new system is fully backward compatible:
- Old decimal debug logs still generated
- Existing room code continues to work
- New hex features are opt-in via new methods

### New BaseRoom Methods
- `LogRoomDebug()` - Enhanced logging with both formats
- `GenerateHexLayoutFile()` - Creates standalone .go files
- `GetHexLayoutArray()` - Returns formatted array string

## Example Integration

See `examples/hex_layout_example.go` for a complete working example demonstrating:
- Creating layouts with hex notation
- Generating debug output
- Copy-paste workflow
- File generation

Run with: `go run examples/hex_layout_example.go`