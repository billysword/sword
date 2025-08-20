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
	
	// Check collision layer
	if ts.loaded.CollisionLayer != nil && index >= 0 && index < len(ts.loaded.CollisionLayer.Data) {
		val := ts.loaded.CollisionLayer.Data[index]
		
		// Value 1 is always solid
		if val == 1 {
			return true
		}
		
		// Value 2 needs context - check position
		// Row 1 (index 16-31) should be passable (decorative ceiling)
		// Other value 2s should be solid (walls)
		if val == 2 {
			width := ts.loaded.TMJ.Width
			row := index / width
			// Row 1 is decorative and should be passable
			if row == 1 {
				return false
			}
			// Other rows with value 2 are solid (walls)
			return true
		}
		
		return false
	}
	
	// Fallback to tile properties from render layer
	if ts.loaded.RenderLayer != nil && index >= 0 && index < len(ts.loaded.RenderLayer.Data) {
		gid := tiled.NormalizeGID(ts.loaded.RenderLayer.Data[index])
		if props, ok := ts.loaded.PropertiesForGID(gid); ok {
			return props.Solid
		}
	}
	return false
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