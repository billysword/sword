package entities

// TileProvider interface defines methods for accessing tile data
// This interface allows the entities package to work with rooms
// without directly importing the world package
type TileProvider interface {
	GetTiles() []int
	GetWidth() int
	GetHeight() int
}

// CollisionChecker interface for checking tile collisions
type CollisionChecker interface {
	IsSolidTile(tileIndex int) bool
}