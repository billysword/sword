package entities

import (
	"math/rand"
	"sword/engine"
)

/*
WandererEnemy represents an enemy with random wandering behavior.
Embeds BaseEnemy for common functionality and implements AI logic
for an enemy that moves randomly without respect to patrol boundaries.

This demonstrates a different AI behavior compared to SlimeEnemy's
patrol behavior, showing the flexibility of the Enemy interface system.
*/
type WandererEnemy struct {
	*BaseEnemy // Embed BaseEnemy for common functionality
	
	// Wanderer-specific AI properties
	moveDirection   int  // -1 for left, 1 for right, 0 for stationary
	moveTimer       int  // Frames until direction change
	randomnessFactor float32 // How random the movement is (0.0 to 1.0)
	pauseChance     float32 // Chance to pause during direction changes
}

/*
NewWandererEnemy creates a new wanderer enemy at the specified position.
Initializes the wanderer with random movement behavior and different
characteristics than the slime enemy.

Parameters:
  - x: Initial horizontal position in physics units
  - y: Initial vertical position in physics units

Returns a pointer to the new WandererEnemy instance that implements the Enemy interface.
*/
func NewWandererEnemy(x, y int) *WandererEnemy {
	wanderer := &WandererEnemy{
		BaseEnemy: NewBaseEnemy(x, y),
		
		// Wanderer-specific AI properties
		moveDirection:   []int{-1, 0, 1}[rand.Intn(3)], // Random starting direction (including stationary)
		moveTimer:       30 + rand.Intn(90),             // 0.5-2 seconds at 60fps (shorter than slime)
		randomnessFactor: 0.8 + rand.Float32()*0.2,     // High randomness (0.8-1.0)
		pauseChance:     0.3,                            // 30% chance to pause
	}
	
	// Configure wanderer-specific properties
	wanderer.SetMoveSpeed(engine.GameConfig.PlayerPhysics.MoveSpeed * 3 / 4) // 75% of player speed (faster than slime)
	wanderer.SetScale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
	
	return wanderer
}

/*
HandleAI implements the Enemy interface's AI method.
Provides random wandering behavior where the wanderer changes direction
frequently and randomly, without respect to spawn point or boundaries.

This demonstrates a completely different AI pattern from the SlimeEnemy.
*/
func (w *WandererEnemy) HandleAI() {
		// Decrease move timer
	w.moveTimer--
	
	// Change direction when timer expires
	if w.moveTimer <= 0 {
		// High chance to change direction due to randomness factor
		if rand.Float32() < w.randomnessFactor {
			// Include pause option based on pause chance
			if rand.Float32() < w.pauseChance {
				w.moveDirection = 0 // Pause
			} else {
				w.moveDirection = []int{-1, 1}[rand.Intn(2)] // Pick random direction (left or right)
			}
		}
		// Reset timer with shorter, more random intervals
		w.moveTimer = 20 + rand.Intn(80) // 0.33-1.66 seconds (more erratic than slime)
	}
	
	// Apply movement based on direction
	if w.moveDirection != 0 {
		w.vx = w.moveDirection * w.GetMoveSpeed()
	}
}

/*
Reset resets the wanderer enemy to initial state at given position.
Overrides BaseEnemy.Reset() to also reset wanderer-specific AI state.

Parameters:
  - x: Reset position horizontal coordinate in physics units
  - y: Reset position vertical coordinate in physics units
*/
func (w *WandererEnemy) Reset(x, y int) {
	// Call base reset
	w.BaseEnemy.Reset(x, y)
	
	// Reset wanderer-specific AI state
	w.moveDirection = []int{-1, 0, 1}[rand.Intn(3)]
	w.moveTimer = 30 + rand.Intn(90)
	w.randomnessFactor = 0.8 + rand.Float32()*0.2
}

/*
GetEnemyType returns the type identifier for wanderer enemies.
Implements the Enemy interface type identification method.

Returns "wanderer" as the enemy type identifier.
*/
func (w *WandererEnemy) GetEnemyType() string {
	return "wanderer"
}

/*
SetRandomnessFactor sets how random the wanderer's movement is.
Higher values mean more frequent direction changes.

Parameters:
  - factor: Randomness factor from 0.0 (predictable) to 1.0 (very random)
*/
func (w *WandererEnemy) SetRandomnessFactor(factor float32) {
	if factor < 0.0 {
		factor = 0.0
	} else if factor > 1.0 {
		factor = 1.0
	}
	w.randomnessFactor = factor
}

/*
SetPauseChance sets the probability that the wanderer will pause when changing direction.

Parameters:
  - chance: Pause probability from 0.0 (never pause) to 1.0 (always pause)
*/
func (w *WandererEnemy) SetPauseChance(chance float32) {
	if chance < 0.0 {
		chance = 0.0
	} else if chance > 1.0 {
		chance = 1.0
	}
	w.pauseChance = chance
}

/*
GetCurrentDirection returns the wanderer's current movement direction.
Useful for debugging or AI behavior analysis.

Returns the current direction (-1 = left, 0 = stationary, 1 = right).
*/
func (w *WandererEnemy) GetCurrentDirection() int {
	return w.moveDirection
}