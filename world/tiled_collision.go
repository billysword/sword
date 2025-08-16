package world

import (
	"sword/engine"
	"sword/internal/tiled"
)

// tiledSolidity provides per-cell solidity derived from the Tiled collision layer or tile properties.
type tiledSolidity struct {
	loaded *tiled.LoadedMap
}

func newTiledSolidity(lm *tiled.LoadedMap) *tiledSolidity {
	return &tiledSolidity{loaded: lm}
}

func (ts *tiledSolidity) IsSolidAtFlatIndex(index int) bool {
	if ts.loaded == nil {
		return false
	}
	// Prefer explicit collision layer
	if ts.loaded.CollisionLayer != nil && index >= 0 && index < len(ts.loaded.CollisionLayer.Data) {
		return ts.loaded.CollisionLayer.Data[index] != 0
	}
	// Fallback: derive from render layer tile properties
	if ts.loaded.RenderLayer != nil && index >= 0 && index < len(ts.loaded.RenderLayer.Data) {
		gid := tiled.NormalizeGID(ts.loaded.RenderLayer.Data[index])
		if props, ok := ts.loaded.PropertiesForGID(gid); ok {
			return props.Solid
		}
	}
	return false
}

// findFloorAtX finds the floor Y position at the given X coordinate using the provided solidity function.
// Returns the Y position in physics units where entities should stand.
func findFloorAtX(tm *TileMap, isSolidAt func(int) bool, x int) int {
	if tm == nil {
		return 0
	}

	u := engine.GetPhysicsUnit()
	tileX := x / u

	// Clamp X to valid tile range
	if tileX < 0 {
		tileX = 0
	}
	if tileX >= tm.Width {
		tileX = tm.Width - 1
	}

	// Scan from top to bottom to find first solid tile
	for tileY := 0; tileY < tm.Height; tileY++ {
		index := tileY*tm.Width + tileX
		if isSolidAt != nil && isSolidAt(index) {
			return tileY * u
		}
	}

	// Fallback: bottom of map if no solid tile found
	return (tm.Height - 1) * u
}