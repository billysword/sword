# CLAUDE.md - Character Movement Logging for Collision & Physics Development

## ✅ **CLEANUP COMPLETED**

All demo/placeholder code has been successfully cleaned up:
- ✅ **State Factory** - Fixed broken unimplemented methods
- ✅ **Enemy AI System** - Removed empty stub, proper interface pattern
- ✅ **Demo Content** - Deleted example room layouts and examples directory
- ✅ **Legacy Constants** - Removed deprecated hardcoded values
- ✅ **Documentation** - Updated all demo references to production language

## 🎯 **NEXT PHASE: Character Movement Logging for Collision & Physics**

To properly design and implement collision detection and physics systems, we need comprehensive logging of character movement data. This will help us understand current behavior and identify areas for improvement.

### 📊 **Character Movement Data We Need to Track**

#### 1. **Position & Velocity Tracking**
```go
// entities/movement_logger.go
type MovementFrame struct {
    Timestamp    int64   // Frame number or time
    PlayerID     string  // For multi-entity tracking
    
    // Position data
    X, Y         float64 // Current position (sub-pixel precision)
    PrevX, PrevY float64 // Previous frame position
    
    // Velocity data  
    VelX, VelY   float64 // Current velocity
    PrevVelX, PrevVelY float64 // Previous frame velocity
    
    // Movement state
    IsMoving     bool    // Any movement this frame
    IsOnGround   bool    // Ground contact
    IsJumping    bool    // Jump in progress
    IsFalling    bool    // Falling state
    
    // Input state
    InputLeft    bool    // Left key pressed
    InputRight   bool    // Right key pressed  
    InputJump    bool    // Jump key pressed
    InputJumpHeld bool   // Jump key held (for variable jump height)
}
```

#### 2. **Collision Detection Events**
```go
type CollisionEvent struct {
    Timestamp    int64
    PlayerID     string
    
    // Collision details
    CollisionType CollisionType // Ground, Wall, Ceiling, Entity
    CollisionSide CollisionSide // Top, Bottom, Left, Right
    
    // Position at collision
    CollisionX, CollisionY float64
    
    // Velocity at collision
    PreCollisionVelX, PreCollisionVelY   float64
    PostCollisionVelX, PostCollisionVelY float64
    
    // Tile/Entity info
    TileX, TileY     int    // Tile coordinates if tile collision
    TileType         int    // Tile type ID
    EntityID         string // Entity ID if entity collision
    
    // Resolution info
    Penetration      float64 // How far entity penetrated
    CorrectionX, CorrectionY float64 // Position correction applied
}

type CollisionType int
const (
    CollisionTile CollisionType = iota
    CollisionEntity
    CollisionWorldBounds
)

type CollisionSide int  
const (
    CollisionTop CollisionSide = iota
    CollisionBottom
    CollisionLeft
    CollisionRight
)
```

#### 3. **Physics State Transitions**
```go
type PhysicsStateChange struct {
    Timestamp    int64
    PlayerID     string
    
    // State transition
    FromState    PhysicsState
    ToState      PhysicsState
    TriggerEvent string // What caused the transition
    
    // Context data
    Position     Point2D
    Velocity     Point2D
    GroundY      float64
    
    // Timing info
    StateFrameCount int // How long in previous state
}

type PhysicsState int
const (
    StateIdle PhysicsState = iota
    StateWalking
    StateJumping
    StateFalling
    StateOnGround
    StateInAir
    StateWallSliding
    StateLanding
)
```

### 🔧 **Implementation Strategy**

#### Phase 1: Basic Movement Logging
```go
// engine/movement_logger.go
type MovementLogger struct {
    enabled      bool
    frameBuffer  []MovementFrame
    bufferSize   int
    currentFrame int64
    logFile      *os.File
    
    // Performance settings
    logEveryNFrames int  // Log every N frames (default: 1)
    maxBufferSize   int  // Max frames in memory before flush
}

func NewMovementLogger(enabled bool, logFilePath string) *MovementLogger {
    if !enabled {
        return &MovementLogger{enabled: false}
    }
    
    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        engine.LogError("Failed to open movement log: " + err.Error())
        return &MovementLogger{enabled: false}
    }
    
    return &MovementLogger{
        enabled:         true,
        frameBuffer:     make([]MovementFrame, 0, 1000),
        bufferSize:      1000,
        logFile:         file,
        logEveryNFrames: 1,
        maxBufferSize:   1000,
    }
}

func (ml *MovementLogger) LogMovementFrame(frame MovementFrame) {
    if !ml.enabled || ml.currentFrame % int64(ml.logEveryNFrames) != 0 {
        return
    }
    
    frame.Timestamp = ml.currentFrame
    ml.frameBuffer = append(ml.frameBuffer, frame)
    ml.currentFrame++
    
    // Flush if buffer is full
    if len(ml.frameBuffer) >= ml.maxBufferSize {
        ml.FlushToFile()
    }
}
```

#### Phase 2: Integration with Player Entity
```go
// entities/player.go - Add to Update() method

func (p *Player) Update() {
    // Capture pre-update state
    prevX, prevY := p.x, p.y
    prevVelX, prevVelY := p.vx, p.vy
    
    // ... existing update logic ...
    
    // Log movement data if logger is enabled
    if movementLogger := engine.GetMovementLogger(); movementLogger != nil {
        frame := MovementFrame{
            PlayerID:     "player_1",
            X:            float64(p.x),
            Y:            float64(p.y),
            PrevX:        float64(prevX),
            PrevY:        float64(prevY),
            VelX:         float64(p.vx),
            VelY:         float64(p.vy),
            PrevVelX:     float64(prevVelX),
            PrevVelY:     float64(prevVelY),
            IsMoving:     p.vx != 0 || p.vy != 0,
            IsOnGround:   p.onGround,
            IsJumping:    p.jumpKeyHeld && p.jumpTimer > 0,
            IsFalling:    p.vy > 0 && !p.onGround,
            InputLeft:    ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft),
            InputRight:   ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight),
            InputJump:    ebiten.IsKeyPressed(ebiten.KeySpace),
            InputJumpHeld: p.jumpKeyHeld,
        }
        movementLogger.LogMovementFrame(frame)
    }
}
```

#### Phase 3: Collision Logging Integration
```go
// world/collision_detector.go - New collision system with logging

type CollisionDetector struct {
    logger *CollisionLogger
}

func (cd *CollisionDetector) CheckTileCollision(entity Entity, tileMap [][]int) []CollisionEvent {
    events := make([]CollisionEvent, 0)
    
    // ... collision detection logic ...
    
    if collisionDetected {
        event := CollisionEvent{
            PlayerID:             entity.GetID(),
            CollisionType:        CollisionTile,
            CollisionSide:        determineSide(entity, tile),
            CollisionX:           float64(entity.GetX()),
            CollisionY:           float64(entity.GetY()),
            PreCollisionVelX:     float64(entity.GetVelX()),
            PreCollisionVelY:     float64(entity.GetVelY()),
            TileX:                tileX,
            TileY:                tileY,
            TileType:             tileMap[tileY][tileX],
            Penetration:          calculatePenetration(entity, tile),
        }
        
        // Apply collision resolution
        correctionX, correctionY := resolveCollision(entity, tile)
        event.CorrectionX = correctionX
        event.CorrectionY = correctionY
        event.PostCollisionVelX = float64(entity.GetVelX())
        event.PostCollisionVelY = float64(entity.GetVelY())
        
        events = append(events, event)
        
        // Log to collision logger
        if cd.logger != nil {
            cd.logger.LogCollision(event)
        }
    }
    
    return events
}
```

### 📈 **Analysis Tools for Logged Data**

#### 1. **Movement Pattern Analysis**
```go
// tools/movement_analyzer.go
func AnalyzeMovementPatterns(logFile string) MovementAnalysis {
    // Analyze for:
    // - Average movement speed
    // - Jump height consistency
    // - Landing accuracy
    // - Input responsiveness
    // - Physics stability
}
```

#### 2. **Collision Frequency Analysis**
```go
func AnalyzeCollisionFrequency(logFile string) CollisionAnalysis {
    // Analyze for:
    // - Most common collision types
    // - Collision hotspots on map
    // - Penetration depth patterns
    // - Resolution effectiveness
}
```

#### 3. **Performance Impact Analysis**
```go
func AnalyzePerformanceImpact(logFile string) PerformanceAnalysis {
    // Analyze for:
    // - Frame time impact of collision detection
    // - Memory usage patterns
    // - CPU usage spikes
    // - Optimization opportunities
}
```

### 🚀 **Configuration & Control**

#### Runtime Configuration
```go
// engine/config.go - Add to Config struct
type Config struct {
    // ... existing fields ...
    
    // Movement logging settings
    EnableMovementLogging    bool    `json:"enable_movement_logging"`
    MovementLogPath         string  `json:"movement_log_path"`
    MovementLogLevel        int     `json:"movement_log_level"` // 1=basic, 2=detailed, 3=verbose
    LogEveryNFrames         int     `json:"log_every_n_frames"`
    
    // Collision logging settings  
    EnableCollisionLogging   bool    `json:"enable_collision_logging"`
    CollisionLogPath        string  `json:"collision_log_path"`
    LogCollisionDetails     bool    `json:"log_collision_details"`
    
    // Performance settings
    MaxLogBufferSize        int     `json:"max_log_buffer_size"`
    FlushLogEveryNSeconds   int     `json:"flush_log_every_n_seconds"`
}
```

#### Development vs Production Settings
```go
func DevelopmentLoggingConfig() Config {
    config := DefaultConfig()
    config.EnableMovementLogging = true
    config.EnableCollisionLogging = true
    config.MovementLogLevel = 3 // Verbose
    config.LogEveryNFrames = 1  // Every frame
    return config
}

func ProductionLoggingConfig() Config {
    config := DefaultConfig()
    config.EnableMovementLogging = false // Disabled for performance
    config.EnableCollisionLogging = false
    return config
}
```

### 🎯 **Integration Plan**

1. **Week 1**: Implement basic movement logging framework
2. **Week 2**: Add collision detection logging  
3. **Week 3**: Create analysis tools and visualizations
4. **Week 4**: Use data to design improved collision system
5. **Week 5**: Implement physics improvements based on analysis

### 💡 **Benefits for Collision & Physics Development**

1. **Data-Driven Design**: Make decisions based on actual movement patterns
2. **Bug Identification**: Catch edge cases and inconsistencies  
3. **Performance Optimization**: Identify bottlenecks before they become problems
4. **Regression Testing**: Ensure changes don't break existing behavior
5. **Player Experience**: Optimize feel and responsiveness based on real data

This logging system will provide the foundation for building robust, data-driven collision detection and physics systems that feel great to play.