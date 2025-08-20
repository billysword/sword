package world

// CollisionType represents different types of collision behaviors
type CollisionType uint32

const (
	CollisionNone     CollisionType = 0 // No collision (air/empty)
	CollisionSolid    CollisionType = 1 // Fully solid tile
	CollisionPlatform CollisionType = 2 // One-way platform or special collision
)

// GetCollisionType returns the collision type for a given collision value
func GetCollisionType(value uint32) CollisionType {
	return CollisionType(value)
}

// IsSolid returns true if the collision type blocks all movement
func (ct CollisionType) IsSolid() bool {
	return ct == CollisionSolid
}

// IsPlatform returns true if this is a one-way platform
func (ct CollisionType) IsPlatform() bool {
	return ct == CollisionPlatform
}

// BlocksMovement returns true if this collision type blocks movement in the given direction
// For platforms, only blocks downward movement (when falling onto them)
func (ct CollisionType) BlocksMovement(fromAbove bool) bool {
	switch ct {
	case CollisionSolid:
		return true
	case CollisionPlatform:
		// Platforms only block when moving down onto them
		return fromAbove
	default:
		return false
	}
}