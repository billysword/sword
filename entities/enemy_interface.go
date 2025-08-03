package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/*
Enemy defines the interface that all enemy types must implement.
This allows for different enemy types with varying AI behaviors,
sprites, and characteristics while maintaining a consistent API
for the game state to interact with them.

All enemy implementations must provide their own AI logic in HandleAI().
*/
type Enemy interface {
	// AI and Update methods
	HandleAI()                    // Process AI logic - must be implemented by each enemy type
	Update()                      // Update physics and apply AI decisions
	
	// Rendering methods
	Draw(screen *ebiten.Image)    // Render without camera offset
	DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) // Render with camera offset
	DrawDebug(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) // Render debug info
	
	// Position and movement methods
	GetPosition() (int, int)      // Get current position in physics units
	SetPosition(x, y int)         // Set position in physics units
	GetVelocity() (int, int)      // Get current velocity in physics units per frame
	SetVelocity(vx, vy int)       // Set velocity in physics units per frame
	
	// State methods
	IsOnGround() bool             // Check if enemy is touching the ground
	Reset(x, y int)               // Reset enemy to initial state at given position
	
	// Type identification
	GetEnemyType() string         // Return the type name of this enemy (e.g., "slime", "goblin")
}