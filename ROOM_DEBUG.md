# Room Construction Debugging

This game includes a room debugging system that automatically logs ASCII representations of rooms when they are first rendered, with support for both decimal and hexadecimal formats suitable for easy copy-paste into code.

## Features

- **Automatic Logging**: Rooms are automatically logged the first time they are constructed
- **Multiple Formats**: Supports both decimal (legacy) and hexadecimal output formats
- **Copy-Paste Ready**: Generates Go array declarations ready for direct code integration
- **Hexadecimal Support**: Uses 0xFF (255) as maximum value with proper hex formatting
- **Layout File Generation**: Creates standalone .go files with formatted room layouts
- **Rotating Log Files**: Daily log rotation with files named `room_debug_YYYY-MM-DD.log`
- **One-time Logging**: Each room is only logged once per session to avoid spam
- **Thread-safe**: Safe for concurrent room creation

## Log Format

Each room log entry includes:
- Room name/zone ID
- Timestamp of first render
- Room dimensions (width x height in tiles)
- ASCII grid representation using 2-digit tile indices (decimal format)
- ASCII grid representation using 2-digit tile indices (hexadecimal format)
- Go array declaration ready for copy-paste
- Empty tiles are represented as `99` (decimal) or `FF` (hexadecimal)

## Example Log Entry

```
=== ROOM DEBUG: main ===
Timestamp: 2025-08-01 23:32:33
Room Dimensions: 80x60 tiles

ASCII Representation (2-digit tile indices - decimal):
99,99,99,99,99,99,99,99,99,99
99,99,99,99,99,99,99,99,99,99
99,99,05,06,07,99,99,99,99,99
01,02,03,01,02,99,99,99,99,99

ASCII Representation (2-digit tile indices - hexadecimal):
FF,FF,FF,FF,FF,FF,FF,FF,FF,FF
FF,FF,FF,FF,FF,FF,FF,FF,FF,FF
FF,FF,05,06,07,FF,FF,FF,FF,FF
01,02,03,01,02,FF,FF,FF,FF,FF

Go Array Format (ready for copy-paste):
levelLayout := [][]int{
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, 0x5, 0x6, 0x7, -1, -1, -1, -1, -1},
	{0x1, 0x2, 0x3, 0x1, 0x2, -1, -1, -1, -1, -1},
}
=== END ROOM DEBUG ===
```

## Hexadecimal Format Benefits

- **Maximum Value**: Uses 0xFF (255) as the maximum tile index instead of 99
- **Compact Representation**: Single-byte values map perfectly to hex format
- **Easy Editing**: Monospace-friendly format that's easy to edit in text editors
- **Industry Standard**: Hexadecimal is commonly used in game development
- **Copy-Paste Ready**: Generated arrays can be directly copied into Go code

## Log Location

- Logs are stored in the `log/` directory
- Files are named `room_debug_YYYY-MM-DD.log`
- Layout files are named `room_layout_<room_name>.go`
- Each day gets its own log file

## Implementation

The debugging system consists of:
- `world/debug.go` - Core debugging functionality with hex support
- `RoomDebugger` singleton - Manages logging state and file operations
- `BaseRoom.LogRoomDebug()` - Helper method for any room type
- `BaseRoom.GenerateHexLayoutFile()` - Creates standalone layout files
- `BaseRoom.GetHexLayoutArray()` - Returns formatted array string
- Automatic integration in `SimpleRoom.buildRoom()`

## Usage for Custom Rooms

### Basic Debug Logging

If you create custom room types, add this call at the end of your room construction:

```go
func (r *MyCustomRoom) buildRoom() {
    // ... room construction logic ...
    
    // Debug: Log ASCII representation on first render
    debugger := world.GetRoomDebugger()
    debugger.LogRoomFirstRender(r.GetZoneID(), r.tileMap)
}
```

Or use the helper method from BaseRoom:

```go
func (r *MyCustomRoom) buildRoom() {
    // ... room construction logic ...
    
    // Debug: Log ASCII representation on first render
    r.LogRoomDebug()
}
```

### Generate Layout Files

To create standalone Go files with your room layouts:

```go
func (r *MyCustomRoom) buildRoom() {
    // ... room construction logic ...
    
    // Generate a .go file with the layout
    r.GenerateHexLayoutFile()
}
```

### Get Layout as String

To get the layout as a formatted string for immediate use:

```go
func (r *MyCustomRoom) exportLayout() {
    layoutString := r.GetHexLayoutArray()
    fmt.Println(layoutString)
    // Copy-paste this output directly into your code
}
```

## Editing Workflows

### Method 1: Direct Array Editing
1. Generate the hex layout using `GenerateHexLayoutFile()`
2. Open the generated `.go` file in your editor
3. Edit the hex values directly (0x00 to 0xFF range)
4. Copy the array declaration back to your room code

### Method 2: Log-based Editing
1. Run your game to generate debug logs
2. Copy the "Go Array Format" section from the log
3. Paste into your room construction code
4. Edit the hex values as needed

### Method 3: Programmatic Generation
1. Use `GetHexLayoutArray()` to get the current layout
2. Modify your room programmatically
3. Export the new layout for manual fine-tuning

## Maintenance

The system includes automatic cleanup functionality to remove old log files. This can be called periodically:

```go
debugger := world.GetRoomDebugger()
debugger.CleanupOldLogs(7) // Keep logs for 7 days
```