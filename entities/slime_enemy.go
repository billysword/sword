package entities

import (
	"math/rand"
	"sword/engine"
)

/*
SlimeEnemy represents a basic slime enemy with patrol behavior.
Embeds BaseEnemy for common functionality and implements specific
AI logic for a simple patrolling slime that moves back and forth
within a defined area.

This serves as an example enemy type that can be used as a template
for creating other enemy types with different behaviors.
*/
type SlimeEnemy struct {
	*BaseEnemy // Embed BaseEnemy for common functionality
	
	// Slime-specific AI properties
	moveDirection int // -1 for left, 1 for right, 0 for stationary
	moveTimer     int // Frames until direction change
	patrolRange   int // Maximum distance to patrol from spawn point
	spawnX        int // Original spawn position for patrol limits
}

/*
NewSlimeEnemy creates a new slime enemy at the specified position.
Initializes the slime with patrol behavior and appropriate settings
for a basic ground-based enemy.

Parameters:
  - x: Initial horizontal position in physics units
  - y: Initial vertical position in physics units

Returns a pointer to the new SlimeEnemy instance that implements the Enemy interface.
*/
func NewSlimeEnemy(x, y int) *SlimeEnemy {
	physicsUnit := engine.GetPhysicsUnit()
	
	slime := &SlimeEnemy{
		BaseEnemy: NewBaseEnemy(x, y),
		
		// Slime-specific AI properties
		moveDirection: []int{-1, 1}[rand.Intn(2)], // Random starting direction
		moveTimer:     60 + rand.Intn(120),        // 1-3 seconds at 60fps
		patrolRange:   200 * physicsUnit,          // 200 physics units patrol range
		spawnX:        x,                          // Remember spawn position
	}
	
	// Configure slime-specific properties
	slime.SetMoveSpeed(engine.GameConfig.PlayerPhysics.MoveSpeed / 2) // Half player speed
	slime.SetScale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
	
	return slime
}

/*
HandleAI implements the Enemy interface's AI method.
Provides patrol behavior where the slime moves back and forth
within a certain range of its spawn point. Changes direction when
reaching patrol boundaries or after a random time interval.

This is the core AI logic specific to slime enemies.
*/
func (s *SlimeEnemy) HandleAI() {
	physicsUnit := engine.GetPhysicsUnit()
	
	// Decrease move timer
	s.moveTimer--
	
	// Check if we've moved too far from spawn point
	distanceFromSpawn := s.x - s.spawnX
	if distanceFromSpawn > s.patrolRange {
		s.moveDirection = -1 // Move back toward spawn
		s.moveTimer = 60 + rand.Intn(60) // Reset timer
	} else if distanceFromSpawn < -s.patrolRange {
		s.moveDirection = 1 // Move back toward spawn
		s.moveTimer = 60 + rand.Intn(60) // Reset timer
	}
	
	// Change direction randomly when timer expires
	if s.moveTimer <= 0 {
		// 70% chance to change direction, 30% chance to keep moving
		if rand.Float32() < 0.7 {
			directions := []int{-1, 0, 1} // left, stop, right
			s.moveDirection = directions[rand.Intn(len(directions))]
		}
		s.moveTimer = 60 + rand.Intn(180) // 1-4 seconds
	}
	
	// Apply movement based on direction
	if s.moveDirection != 0 {
		s.vx = s.moveDirection * s.GetMoveSpeed() * physicsUnit
	}
}

/*
Reset resets the slime enemy to initial state at given position.
Overrides BaseEnemy.Reset() to also reset slime-specific AI state.

Parameters:
  - x: Reset position horizontal coordinate in physics units
  - y: Reset position vertical coordinate in physics units
*/
func (s *SlimeEnemy) Reset(x, y int) {
	// Call base reset
	s.BaseEnemy.Reset(x, y)
	
	// Reset slime-specific AI state
	s.spawnX = x
	s.moveDirection = []int{-1, 1}[rand.Intn(2)]
	s.moveTimer = 60 + rand.Intn(120)
}

/*
GetEnemyType returns the type identifier for slime enemies.
Implements the Enemy interface type identification method.

Returns "slime" as the enemy type identifier.
*/
func (s *SlimeEnemy) GetEnemyType() string {
	return "slime"
}

/*
SetPatrolRange sets the patrol range for this slime.
Allows customization of how far the slime will wander from its spawn point.

Parameters:
  - range: Maximum patrol distance in physics units
*/
func (s *SlimeEnemy) SetPatrolRange(patrolRange int) {
	s.patrolRange = patrolRange
}

/*
GetPatrolRange returns the current patrol range.
Useful for debugging or AI behavior analysis.

Returns the patrol range in physics units.
*/
func (s *SlimeEnemy) GetPatrolRange() int {
	return s.patrolRange
}

/*
GetDistanceFromSpawn returns how far the slime is from its spawn point.
Useful for debugging patrol behavior or implementing more complex AI.

Returns the distance from spawn in physics units (positive = right, negative = left).
*/
func (s *SlimeEnemy) GetDistanceFromSpawn() int {
	return s.x - s.spawnX
}