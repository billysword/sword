package world

import (
	"fmt"
	"sword/engine"
	"sword/internal/tiled"
)

/*
DebugCompareLayers prints a visual comparison of render and collision layers
to help identify mismatches between what's visible and what's collidable.
*/
func DebugCompareLayers(lm *tiled.LoadedMap) {
	if lm == nil {
		engine.LogInfo("DEBUG_LAYERS: No loaded map to debug")
		return
	}
	
	engine.LogInfo("=== LAYER COMPARISON DEBUG ===")
	
	// Get dimensions
	width := 16  // From TMJ file
	height := 12 // From TMJ file
	
	if lm.RenderLayer != nil {
		engine.LogInfo(fmt.Sprintf("Render layer: %d tiles", len(lm.RenderLayer.Data)))
	}
	if lm.CollisionLayer != nil {
		engine.LogInfo(fmt.Sprintf("Collision layer: %d tiles", len(lm.CollisionLayer.Data)))
	}
	
	// Print collision layer in visual format
	if lm.CollisionLayer != nil && len(lm.CollisionLayer.Data) >= width*height {
		engine.LogInfo("COLLISION LAYER:")
		for y := 0; y < height; y++ {
			line := fmt.Sprintf("Row %2d: ", y)
			for x := 0; x < width; x++ {
				index := y*width + x
				if index < len(lm.CollisionLayer.Data) {
					val := lm.CollisionLayer.Data[index]
					if val == 0 {
						line += ". "  // Air
					} else {
						line += "# "  // Solid
					}
				} else {
					line += "? "
				}
			}
			engine.LogInfo(line)
		}
	}
	
	// Print render layer in visual format
	if lm.RenderLayer != nil && len(lm.RenderLayer.Data) >= width*height {
		engine.LogInfo("RENDER LAYER:")
		for y := 0; y < height; y++ {
			line := fmt.Sprintf("Row %2d: ", y)
			for x := 0; x < width; x++ {
				index := y*width + x
				if index < len(lm.RenderLayer.Data) {
					val := lm.RenderLayer.Data[index]
					if val == 0 {
						line += ". "  // Empty
					} else {
						line += fmt.Sprintf("%X ", val%16)  // Show tile ID (mod 16 for readability)
					}
				} else {
					line += "? "
				}
			}
			engine.LogInfo(line)
		}
	}
	
	engine.LogInfo("=== END LAYER COMPARISON ===")
}