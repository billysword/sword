package engine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger handles file-based logging for the game
type Logger struct {
	file   *os.File
	logger *log.Logger
	mutex  sync.Mutex
}

// LoggerManager manages multiple logger instances for different categories
type LoggerManager struct {
	roomLogger   *Logger
	playerLogger *Logger
	gameLogger   *Logger
	mutex        sync.Mutex
}

var (
	loggerManager *LoggerManager
	once          sync.Once
)

// InitLoggers initializes all category-specific loggers
func InitLoggers(baseFilename string) error {
	var err error
	once.Do(func() {
		loggerManager = &LoggerManager{}

		// Create logs directory if it doesn't exist
		logsDir := "logs"
		if err = os.MkdirAll(logsDir, 0755); err != nil {
			return
		}

		// Create timestamped base name
		timestamp := time.Now().Format("2006-01-02_15-04-05")

		// Extract base name without extension
		baseName := baseFilename
		if ext := filepath.Ext(baseFilename); ext != "" {
			baseName = baseFilename[:len(baseFilename)-len(ext)]
		}

		// Initialize room/layout logger
		roomLogPath := filepath.Join(logsDir, fmt.Sprintf("%s_room_%s.log", baseName, timestamp))
		loggerManager.roomLogger, err = createLogger(roomLogPath)
		if err != nil {
			return
		}

		// Initialize player input/diagnostics logger
		playerLogPath := filepath.Join(logsDir, fmt.Sprintf("%s_player_%s.log", baseName, timestamp))
		loggerManager.playerLogger, err = createLogger(playerLogPath)
		if err != nil {
			return
		}

		// Initialize general game logger
		gameLogPath := filepath.Join(logsDir, fmt.Sprintf("%s_game_%s.log", baseName, timestamp))
		loggerManager.gameLogger, err = createLogger(gameLogPath)
		if err != nil {
			return
		}

		// Log initialization messages
		loggerManager.roomLogger.LogInfo(fmt.Sprintf("=== Room/Layout Logger Initialized - Log file: %s ===", roomLogPath))
		loggerManager.playerLogger.LogInfo(fmt.Sprintf("=== Player Input/Diagnostics Logger Initialized - Log file: %s ===", playerLogPath))
		loggerManager.gameLogger.LogInfo(fmt.Sprintf("=== General Game Logger Initialized - Log file: %s ===", gameLogPath))
	})
	return err
}

// createLogger creates a new Logger instance with the specified file path
func createLogger(logPath string) (*Logger, error) {
	logger := &Logger{}

	// Create or open log file with append mode
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger.file = file
	logger.logger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)

	return logger, nil
}

// GetLoggerManager returns the singleton logger manager instance
func GetLoggerManager() *LoggerManager {
	if loggerManager == nil {
		// Initialize with default filename if not already done
		InitLoggers("game")
	}
	return loggerManager
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[INFO] %s", message)
	}
}

// LogDebug logs a debug message
func (l *Logger) LogDebug(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[DEBUG] %s", message)
	}
}

// LogRoomTile logs room tile information - routes to room logger
func (l *Logger) LogRoomTile(roomName string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[ROOM_TILE] %s: %s", roomName, message)
	}
}

// LogPlayerInput logs player input and position - routes to player logger
func (l *Logger) LogPlayerInput(key string, playerX, playerY int, roomName string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[PLAYER_INPUT] Key=%s Position=(%d,%d) Room=%s", key, playerX, playerY, roomName)
	}
}

// LogSprite logs sprite-related information
func (l *Logger) LogSprite(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[SPRITE] %s", message)
	}
}

// LogRoomLayout logs the complete room layout - routes to room logger
func (l *Logger) LogRoomLayout(roomName string, width, height int, layout string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.logger != nil {
		l.logger.Printf("[ROOM_LAYOUT] %s (%dx%d):\n%s", roomName, width, height, layout)
	}
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.file != nil {
		if l.logger != nil {
			l.logger.Printf("[INFO] === Logger Closing ===")
		}
		// Flush any pending writes
		l.file.Sync()
		err := l.file.Close()
		l.file = nil // Prevent double-close
		l.logger = nil
		return err
	}
	return nil
}

// CloseAllLoggers closes all category-specific loggers
func (lm *LoggerManager) CloseAllLoggers() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	var lastErr error

	if lm.roomLogger != nil {
		if err := lm.roomLogger.Close(); err != nil {
			lastErr = err
		}
	}

	if lm.playerLogger != nil {
		if err := lm.playerLogger.Close(); err != nil {
			lastErr = err
		}
	}

	if lm.gameLogger != nil {
		if err := lm.gameLogger.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Convenience functions for global logger access

// LogInfo logs an informational message using the general game logger
func LogInfo(message string) {
	GetLoggerManager().gameLogger.LogInfo(message)
}

// LogDebug logs a debug message using the general game logger
func LogDebug(message string) {
	GetLoggerManager().gameLogger.LogDebug(message)
}

// LogRoomTile logs room tile information using the room logger
func LogRoomTile(roomName string, message string) {
	GetLoggerManager().roomLogger.LogRoomTile(roomName, message)
}

// LogPlayerInput logs player input and position using the player logger
func LogPlayerInput(key string, playerX, playerY int, roomName string) {
	GetLoggerManager().playerLogger.LogPlayerInput(key, playerX, playerY, roomName)
}

// LogSprite logs sprite-related information using the general game logger
func LogSprite(message string) {
	GetLoggerManager().gameLogger.LogSprite(message)
}

// LogRoomLayout logs the complete room layout using the room logger
func LogRoomLayout(roomName string, width, height int, layout string) {
	GetLoggerManager().roomLogger.LogRoomLayout(roomName, width, height, layout)
}

// CloseLogger closes all loggers
func CloseLogger() error {
	if loggerManager != nil {
		return loggerManager.CloseAllLoggers()
	}
	return nil
}

// InitLogger maintains backward compatibility with single logger initialization
func InitLogger(filename string) error {
	return InitLoggers(filename)
}
