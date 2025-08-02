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

	// Generate both decimal and hex representations
	asciiRepDecimal := rd.generateASCIIRepresentation(tileMap, false)
	asciiRepHex := rd.generateASCIIRepresentation(tileMap, true)
	layoutArray := rd.generateLayoutArray(tileMap)

	// Create log entry with both formats
	logEntry := fmt.Sprintf("=== ROOM DEBUG: %s ===\n", roomName)
	logEntry += fmt.Sprintf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	logEntry += fmt.Sprintf("Room Dimensions: %dx%d tiles\n", tileMap.Width, tileMap.Height)
	logEntry += "ASCII Representation (2-digit tile indices - decimal):\n"
	logEntry += asciiRepDecimal + "\n\n"
	logEntry += "ASCII Representation (2-digit tile indices - hexadecimal):\n"
	logEntry += asciiRepHex + "\n\n"
	logEntry += "Go Array Format (ready for copy-paste):\n"
	logEntry += layoutArray + "\n"
	logEntry += "=== END ROOM DEBUG ===\n\n"

	// Write to rotating log file
	rd.writeToRotatingLog(logEntry)
}

// generateASCIIRepresentation creates an ASCII representation of the room
// using 2-digit tile indices with comma separation per row
func (rd *RoomDebugger) generateASCIIRepresentation(tileMap *TileMap, useHex bool) string {
	var builder strings.Builder

	for y := 0; y < tileMap.Height; y++ {
		var rowValues []string
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]
			// Convert -1 (empty) to appropriate empty value
			if tileIndex == -1 {
				if useHex {
					rowValues = append(rowValues, "FF") // 255 in hex
				} else {
					rowValues = append(rowValues, "99") // Keep 99 for decimal compatibility
				}
			} else {
				if useHex {
					// Cap at 0xFF (255) and format as hex
					if tileIndex > 255 {
						tileIndex = 255
					}
					rowValues = append(rowValues, fmt.Sprintf("%02X", tileIndex))
				} else {
					rowValues = append(rowValues, fmt.Sprintf("%02d", tileIndex))
				}
			}
		}
		builder.WriteString(strings.Join(rowValues, ","))
		if y < tileMap.Height-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// generateLayoutArray creates a Go array declaration ready for copy-paste into code
func (rd *RoomDebugger) generateLayoutArray(tileMap *TileMap) string {
	var builder strings.Builder
	
	builder.WriteString("levelLayout := [][]int{\n")
	
	for y := 0; y < tileMap.Height; y++ {
		builder.WriteString("\t{")
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]
			
			// Convert tile indices to hex format (0x notation)
			if tileIndex == -1 {
				builder.WriteString("-1")
			} else {
				// Cap at 0xFF (255) for hex format
				if tileIndex > 255 {
					tileIndex = 255
				}
				if tileIndex < 16 {
					builder.WriteString(fmt.Sprintf("0x%01X", tileIndex))
				} else {
					builder.WriteString(fmt.Sprintf("0x%02X", tileIndex))
				}
			}
			
			if x < tileMap.Width-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString("}")
		if y < tileMap.Height-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	
	builder.WriteString("}")
	return builder.String()
}

// generateArrayBody creates just the array body (without variable declaration)
func (rd *RoomDebugger) generateArrayBody(tileMap *TileMap) string {
	var builder strings.Builder
	
	builder.WriteString("{\n")
	
	for y := 0; y < tileMap.Height; y++ {
		builder.WriteString("\t{")
		for x := 0; x < tileMap.Width; x++ {
			tileIndex := tileMap.Tiles[y][x]
			
			// Convert tile indices to hex format (0x notation)
			if tileIndex == -1 {
				builder.WriteString("-1")
			} else {
				// Cap at 0xFF (255) for hex format
				if tileIndex > 255 {
					tileIndex = 255
				}
				if tileIndex < 16 {
					builder.WriteString(fmt.Sprintf("0x%01X", tileIndex))
				} else {
					builder.WriteString(fmt.Sprintf("0x%02X", tileIndex))
				}
			}
			
			if x < tileMap.Width-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString("}")
		if y < tileMap.Height-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	
	builder.WriteString("}")
	return builder.String()
}

// GenerateHexLayoutFile creates a standalone .go file with the room layout in hex format
func (rd *RoomDebugger) GenerateHexLayoutFile(roomName string, tileMap *TileMap) {
	rd.mutex.Lock()
	defer rd.mutex.Unlock()

	// Create filename
	filename := fmt.Sprintf("room_layout_%s.go", strings.ReplaceAll(roomName, " ", "_"))
	filepath := filepath.Join(rd.logDir, filename)

	// Generate the file content
	content := fmt.Sprintf(`// Auto-generated room layout for: %s
// Generated: %s
// Dimensions: %dx%d tiles

package main

// %sLayout contains the tile layout in hexadecimal format
// -1 = empty tile, 0x00-0xFF = tile indices
var %sLayout = [][]int%s
`, roomName, time.Now().Format("2006-01-02 15:04:05"), 
   tileMap.Width, tileMap.Height,
   strings.Title(strings.ReplaceAll(roomName, " ", "")),
   strings.Title(strings.ReplaceAll(roomName, " ", "")),
   rd.generateArrayBody(tileMap))

	// Write to file
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating layout file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		fmt.Printf("Error writing layout file: %v\n", err)
	}
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