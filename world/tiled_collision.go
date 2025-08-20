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
	
	// Prefer explicit collision layer when present; fallback to tile properties from render layer
	return ts.loaded.IsSolidAt(index)
}

// findFloorAtX finds the floor Y position at the given X coordinate using the provided solidity function.
// Returns the Y position in physics units where entities should stand (just above solid ground).
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

	// Scan from top to bottom to find air space above solid ground
	for tileY := 0; tileY < tm.Height-1; tileY++ {
		currentIndex := tileY*tm.Width + tileX
		nextIndex := (tileY+1)*tm.Width + tileX
		
		// Check if current tile is air and next tile is solid
		currentSolid := isSolidAt != nil && isSolidAt(currentIndex)
		nextSolid := isSolidAt != nil && isSolidAt(nextIndex)
		
		if !currentSolid && nextSolid {
			// Found air above solid ground - return the Y position of the air tile (where player stands)
			return tileY * u
		}
	}

	// Fallback: if no proper floor found, return near bottom of map
	return (tm.Height - 2) * u
}