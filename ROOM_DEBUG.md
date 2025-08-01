# Room Construction Debugging

This game includes a room debugging system that automatically logs ASCII representations of rooms when they are first rendered.

## Features

- **Automatic Logging**: Rooms are automatically logged the first time they are constructed
- **ASCII Representation**: Each room is displayed as a 2-digit tile index grid with comma-separated values per row
- **Rotating Log Files**: Daily log rotation with files named `room_debug_YYYY-MM-DD.log`
- **One-time Logging**: Each room is only logged once per session to avoid spam
- **Thread-safe**: Safe for concurrent room creation

## Log Format

Each room log entry includes:
- Room name/zone ID
- Timestamp of first render
- Room dimensions (width x height in tiles)
- ASCII grid representation using 2-digit tile indices
- Empty tiles are represented as `99`

## Example Log Entry

```
=== ROOM DEBUG: main ===
Timestamp: 2025-08-01 23:32:33
Room Dimensions: 80x60 tiles
ASCII Representation (2-digit tile indices):
99,99,99,99,99,99,99,99,99,99
99,99,99,99,99,99,99,99,99,99
99,99,05,06,07,99,99,99,99,99
01,02,03,01,02,99,99,99,99,99
=== END ROOM DEBUG ===
```

## Log Location

- Logs are stored in the `log/` directory
- Files are named `room_debug_YYYY-MM-DD.log`
- Each day gets its own log file

## Implementation

The debugging system consists of:
- `world/debug.go` - Core debugging functionality
- `RoomDebugger` singleton - Manages logging state and file operations
- `BaseRoom.LogRoomDebug()` - Helper method for any room type
- Automatic integration in `SimpleRoom.buildRoom()`

## Usage for Custom Rooms

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

## Maintenance

The system includes automatic cleanup functionality to remove old log files. This can be called periodically:

```go
debugger := world.GetRoomDebugger()
debugger.CleanupOldLogs(7) // Keep logs for 7 days
```