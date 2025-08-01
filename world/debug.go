package world

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RoomDebugger handles room debugging output with rotating log files
type RoomDebugger struct {
	logDir       string
	mutex        sync.Mutex
	renderedRooms map[string]bool // Track which rooms have been logged
}

var debugger *RoomDebugger
var once sync.Once

// GetRoomDebugger returns a singleton instance of the room debugger
func GetRoomDebugger() *RoomDebugger {
	once.Do(func() {
		debugger = &RoomDebugger{
			logDir:        "log",
			renderedRooms: make(map[string]bool),
		}
		// Ensure log directory exists
		os.MkdirAll(debugger.logDir, 0755)
	})
	return debugger
}

// LogRoomFirstRender logs the ASCII representation of a room on its first render
func (rd *RoomDebugger) LogRoomFirstRender(roomName string, tileMap *TileMap) {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	// Check if this room has already been logged
	if rd.renderedRooms[roomName] {
		return
	}

	// Mark this room as logged
	rd.renderedRooms[roomName] = true

	// Generate ASCII representation
	asciiRep := rd.generateASCIIRepresentation(tileMap)

	// Create log entry
	logEntry := fmt.Sprintf("=== ROOM DEBUG: %s ===\n", roomName)
	logEntry += fmt.Sprintf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	logEntry += fmt.Sprintf("Room Dimensions: %dx%d tiles\n", tileMap.Width, tileMap.Height)
	logEntry += "ASCII Representation (2-digit tile indices):\n"
	logEntry += asciiRep
	logEntry += "\n=== END ROOM DEBUG ===\n\n"

	// Write to rotating log file
	rd.writeToRotatingLog(logEntry)
}

// generateASCIIRepresentation creates an ASCII representation of the room
// using 2-digit tile indices with comma separation per row
func (rd *RoomDebugger) generateASCIIRepresentation(tileMap *TileMap) string {
	var builder strings.Builder

	for y := 0; y < tileMap.Height; y++ {
		var rowValues []string
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]
			// Convert -1 (empty) to 99 for better visualization
			if tileIndex == -1 {
				rowValues = append(rowValues, "99")
			} else {
				rowValues = append(rowValues, fmt.Sprintf("%02d", tileIndex))
			}
		}
		builder.WriteString(strings.Join(rowValues, ","))
		if y < tileMap.Height-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// writeToRotatingLog writes the log entry to a rotating log file
func (rd *RoomDebugger) writeToRotatingLog(entry string) {
	// Create filename with current date for daily rotation
	filename := fmt.Sprintf("room_debug_%s.log", time.Now().Format("2006-01-02"))
	filepath := filepath.Join(rd.logDir, filename)

	// Open file in append mode, create if it doesn't exist
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the log entry
	if _, err := file.WriteString(entry); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}

// CleanupOldLogs removes log files older than the specified number of days
func (rd *RoomDebugger) CleanupOldLogs(daysToKeep int) {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	cutoff := time.Now().AddDate(0, 0, -daysToKeep)

	files, err := os.ReadDir(rd.logDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "room_debug_") && strings.HasSuffix(file.Name(), ".log") {
			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				filepath := filepath.Join(rd.logDir, file.Name())
				os.Remove(filepath)
			}
		}
	}
}