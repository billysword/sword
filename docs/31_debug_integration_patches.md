# Debug Logging Integration Patches

This document shows specific code changes to add debug logging to your existing systems for troubleshooting camera, window sizes, tile maps, scale and viewport issues.

## 1. Camera System Debug Logging

### In `engine/camera.go` - Update() method:

```go
func (c *Camera) Update(targetX, targetY int) {
	// Log camera state before update (add this)
	LogCameraDebug(c.x, c.y, c.targetX, c.targetY, c.width, c.height, c.worldWidth, c.worldHeight)
	
	// Calculate ideal camera position to center the target
	idealX := float64(targetX) - float64(c.width)/2
	idealY := float64(targetY) - float64(c.height)/2
	
	// Apply dead zone logic
	currentCenterX := c.x + float64(c.width)/2
	currentCenterY := c.y + float64(c.height)/2
	
	// Only move camera if target is outside dead zone
	if math.Abs(float64(targetX)-currentCenterX) > float64(c.deadZoneX) {
		c.targetX = idealX
	}
	if math.Abs(float64(targetY)-currentCenterY) > float64(c.deadZoneY) {
		c.targetY = idealY
	}
	
	// Smooth camera movement using linear interpolation
	c.x += (c.targetX - c.x) * c.smoothing
	c.y += (c.targetY - c.y) * c.smoothing
	
	// Constrain camera to world bounds with margins
	c.constrainToWorld()
	
	// Log camera state after update (add this)
	LogCameraDebug(c.x, c.y, c.targetX, c.targetY, c.width, c.height, c.worldWidth, c.worldHeight)
}
```

### In `engine/camera.go` - Coordinate conversion methods:

```go
func (c *Camera) WorldToScreen(worldX, worldY int) (int, int) {
	screenX := worldX - int(c.x)
	screenY := worldY - int(c.y)
	
	// Add debug logging for coordinate conversion (add this)
	LogCoordinateConversion("WorldToScreen", worldX, worldY, screenX, screenY)
	
	return screenX, screenY
}

func (c *Camera) ScreenToWorld(screenX, screenY int) (int, int) {
	worldX := int(c.x) + screenX
	worldY := int(c.y) + screenY
	
	// Add debug logging for coordinate conversion (add this)
	LogCoordinateConversion("ScreenToWorld", screenX, screenY, worldX, worldY)
	
	return worldX, worldY
}
```

## 2. InGame State Debug Logging

### In `states/ingame_state.go` - NewInGameState():

```go
func NewInGameState(sm *engine.StateManager) *InGameState {
	// Get the actual window size for camera viewport
	windowWidth, windowHeight := ebiten.WindowSize()
	
	// Add viewport debug logging (add this)
	engine.LogViewportDebug(windowWidth, windowHeight, 
		engine.GameConfig.TileScaleFactor, 
		engine.GameConfig.CharScaleFactor, 
		engine.GetPhysicsUnit())

	physicsUnit := engine.GetPhysicsUnit()
	groundY := engine.GameConfig.GroundLevel * physicsUnit

	return &InGameState{
		stateManager: sm,
		player:       entities.NewPlayer(50*physicsUnit, groundY),
		enemies:      make([]entities.Enemy, 0),
                // loadedMap represents a parsed Tiled map
                currentRoom:  world.NewTiledRoomFromLoadedMap("main", loadedMap),
		camera:       engine.NewCamera(windowWidth, windowHeight),
	}
}
```

### In `states/ingame_state.go` - SwitchToRoom() method around line 240:

```go
// Set up camera bounds based on room size
if ig.camera != nil && ig.currentRoom != nil {
	tileMap := ig.currentRoom.GetTileMap()
	if tileMap != nil {
		// Convert tile dimensions to pixel dimensions
		physicsUnit := engine.GetPhysicsUnit()
		worldWidth := tileMap.Width * physicsUnit
		worldHeight := tileMap.Height * physicsUnit
		
		// Add tilemap debug logging (add this)
		engine.LogTileMapDebug(ig.currentRoom.GetZoneID(), 
			tileMap.Width, tileMap.Height, physicsUnit, worldWidth, worldHeight)
		
		ig.camera.SetWorldBounds(worldWidth, worldHeight)

		// Center camera on player initially
		px, py := ig.player.GetPosition()
		ig.camera.CenterOn(px/physicsUnit, py/physicsUnit)
	}
}
```

## 3. Room Rendering Debug Logging

### In `world/room.go` - DrawWithCamera() method around line 365:

```go
func (br *BaseRoom) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64, spriteProvider func(int) *ebiten.Image) {
	physicsUnit := engine.GetPhysicsUnit()
	
	// Add tilemap debug logging when first rendering (add this)
	static bool logged = false
	if !logged {
		engine.LogTileMapDebug(br.zoneID, br.tileMap.Width, br.tileMap.Height, 
			physicsUnit, br.tileMap.Width*physicsUnit, br.tileMap.Height*physicsUnit)
		logged = true
	}
	
	for y := 0; y < br.tileMap.Height; y++ {
		for x := 0; x < br.tileMap.Width; x++ {
			tileIndex := br.tileMap.Tiles[y][x]
			if tileIndex != -1 {
				sprite := spriteProvider(tileIndex)
				if sprite != nil {
					op := &ebiten.DrawImageOptions{}
					
					// Scale tiles using global scale factor
					op.GeoM.Scale(engine.GameConfig.TileScaleFactor, engine.GameConfig.TileScaleFactor)
					renderX := float64(x * physicsUnit) + cameraOffsetX
					renderY := float64(y * physicsUnit) + cameraOffsetY
					op.GeoM.Translate(renderX, renderY)
					
					// Add rendering debug logging for first few tiles (add this)
					if x < 3 && y < 3 {
						worldX := float64(x * physicsUnit)
						worldY := float64(y * physicsUnit)
						objectType := fmt.Sprintf("Tile[%d]@(%d,%d)", tileIndex, x, y)
						engine.LogRenderingDebug(objectType, worldX, worldY, renderX, renderY, 
							engine.GameConfig.TileScaleFactor)
					}
					
					screen.DrawImage(sprite, op)
				}
			}
		}
	}
}
```

## 4. Player/Entity Debug Logging

### In `entities/player.go` - DrawWithCamera() method around line 190:

```go
func (p *Player) DrawWithCamera(screen *ebiten.Image, cameraOffsetX, cameraOffsetY float64) {
	if p.sprite != nil {
		// Create drawing options
		op := &ebiten.DrawImageOptions{}
		
		// Set up drawing options with camera offset
		op.GeoM.Scale(engine.GameConfig.CharScaleFactor, engine.GameConfig.CharScaleFactor)
		
		// Convert player position from physics units to pixels and apply camera offset
		renderX := float64(p.x)/float64(engine.GetPhysicsUnit()) + cameraOffsetX
		renderY := float64(p.y)/float64(engine.GetPhysicsUnit()) + cameraOffsetY
		
		// Add player rendering debug logging (add this)
		worldX := float64(p.x) / float64(engine.GetPhysicsUnit())
		worldY := float64(p.y) / float64(engine.GetPhysicsUnit())
		engine.LogRenderingDebug("Player", worldX, worldY, renderX, renderY, 
			engine.GameConfig.CharScaleFactor)
		
		op.GeoM.Translate(renderX, renderY)
		
		screen.DrawImage(p.sprite, op)
	}
}
```

## 5. Configuration Changes Debug Logging

### In `engine/config.go` - SetConfig() method:

```go
func SetConfig(config Config) {
	// Log before config change (add this)
	windowW, windowH := ebiten.WindowSize()
	oldPhysicsUnit := GetPhysicsUnit()
	LogViewportDebug(windowW, windowH, GameConfig.TileScaleFactor, GameConfig.CharScaleFactor, oldPhysicsUnit)
	
	GameConfig = config
	
	// Log after config change (add this)
	newPhysicsUnit := GetPhysicsUnit()
	LogViewportDebug(windowW, windowH, GameConfig.TileScaleFactor, GameConfig.CharScaleFactor, newPhysicsUnit)
	
	LogInfo(fmt.Sprintf("Configuration changed - PhysicsUnit: %d->%d", oldPhysicsUnit, newPhysicsUnit))
}
```

## Quick Integration Summary

Add these imports to any file using the debug logging:
```go
import "sword/engine"
```

### Most Critical Debug Points:

1. **Camera Updates**: Add `LogCameraDebug()` calls in `camera.Update()`
2. **Room Loading**: Add `LogTileMapDebug()` calls when setting camera bounds
3. **Viewport Changes**: Add `LogViewportDebug()` calls on window resize or config changes
4. **Coordinate Issues**: Add `LogCoordinateConversion()` calls in conversion methods
5. **Rendering Problems**: Add `LogRenderingDebug()` calls for key objects

### Expected Log Files:

- `game_room_*.log` - Room layouts, tile map dimensions, tile rendering
- `game_player_*.log` - Camera movements, coordinate conversions, player inputs
- `game_game_*.log` - Viewport info, general rendering, sprite scaling

Use these logs to track down integration issues by comparing expected vs actual values for camera positions, world bounds, rendering coordinates, and scale factors.