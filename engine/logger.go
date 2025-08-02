package engine

import (
	"log"
	"os"
	"sync"
)

// Logger handles file-based logging for the game
type Logger struct {
	file   *os.File
	logger *log.Logger
	mutex  sync.Mutex
}

var (
	gameLogger *Logger
	once       sync.Once
)

// InitLogger initializes the game logger with a file
func InitLogger(filename string) error {
	var err error
	once.Do(func() {
		gameLogger = &Logger{}

		// Create or open log file with append mode
		gameLogger.file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		// Create logger with custom format
		gameLogger.logger = log.New(gameLogger.file, "", log.LstdFlags|log.Lmicroseconds)

		// Log initialization
		gameLogger.LogInfo("=== Game Logger Initialized ===")
	})
	return err
}

// GetLogger returns the singleton game logger instance
func GetLogger() *Logger {
	if gameLogger == nil {
		// Initialize with default filename if not already done
		InitLogger("game.log")
	}
	return gameLogger
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[INFO] %s", message)
}

// LogDebug logs a debug message
func (l *Logger) LogDebug(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[DEBUG] %s", message)
}

// LogRoomTile logs room tile information
func (l *Logger) LogRoomTile(roomName string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[ROOM_TILE] %s: %s", roomName, message)
}

// LogPlayerInput logs player input and position
func (l *Logger) LogPlayerInput(key string, playerX, playerY int, roomName string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[PLAYER_INPUT] Key=%s Position=(%d,%d) Room=%s", key, playerX, playerY, roomName)
}

// LogSprite logs sprite-related information
func (l *Logger) LogSprite(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[SPRITE] %s", message)
}

// LogRoomLayout logs the complete room layout
func (l *Logger) LogRoomLayout(roomName string, width, height int, layout string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf("[ROOM_LAYOUT] %s (%dx%d):\n%s", roomName, width, height, layout)
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.file != nil {
		l.LogInfo("=== Game Logger Closing ===")
		// Flush any pending writes
		l.file.Sync()
		err := l.file.Close()
		l.file = nil // Prevent double-close
		return err
	}
	return nil
}

// Convenience functions for global logger access

// LogInfo logs an informational message using the global logger
func LogInfo(message string) {
	GetLogger().LogInfo(message)
}

// LogDebug logs a debug message using the global logger
func LogDebug(message string) {
	GetLogger().LogDebug(message)
}

// LogRoomTile logs room tile information using the global logger
func LogRoomTile(roomName string, message string) {
	GetLogger().LogRoomTile(roomName, message)
}

// LogPlayerInput logs player input and position using the global logger
func LogPlayerInput(key string, playerX, playerY int, roomName string) {
	GetLogger().LogPlayerInput(key, playerX, playerY, roomName)
}

// LogSprite logs sprite-related information using the global logger
func LogSprite(message string) {
	GetLogger().LogSprite(message)
}

// LogRoomLayout logs the complete room layout using the global logger
func LogRoomLayout(roomName string, width, height int, layout string) {
	GetLogger().LogRoomLayout(roomName, width, height, layout)
}

// CloseLogger closes the global logger
func CloseLogger() error {
	if gameLogger != nil {
		return gameLogger.Close()
	}
	return nil
}
