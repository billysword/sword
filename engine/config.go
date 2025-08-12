package engine

/*
PlayerPhysicsConfig holds all player-specific physics and collision settings.
This allows fine-tuning of player movement, collision detection, and sprite sizing
relative to the tileset without modifying code.
*/
type PlayerPhysicsConfig struct {
	// Sprite dimensions
	SpriteWidth  int     // Base sprite width in pixels
	SpriteHeight int     // Base sprite height in pixels
	
	// Collision box (relative to sprite dimensions)
	CollisionBoxOffsetX float64 // Offset from sprite left edge (0-1)
	CollisionBoxOffsetY float64 // Offset from sprite top edge (0-1)
	CollisionBoxWidth   float64 // Width as fraction of sprite width (0-1)
	CollisionBoxHeight  float64 // Height as fraction of sprite height (0-1)
	
	// Ground detection
	GroundCheckOffset   int     // Pixels below collision box to check for ground
	GroundCheckWidth    float64 // Width of ground check as fraction of collision box
	
	// Movement physics
	MoveSpeed          int     // Horizontal movement speed in physics units
	JumpPower          int     // Initial jump velocity in physics units
	AirControl         float64 // Movement control while airborne (0-1)
	Friction           int     // Ground friction applied each frame
	AirFriction        int     // Air friction applied each frame
	
	// Jump mechanics
	CoyoteTime         int     // Frames after leaving ground where jump is still allowed
	JumpBufferTime     int     // Frames to buffer jump input before landing
	VariableJumpHeight bool    // Allow controlling jump height by release timing
	MinJumpHeight      float64 // Minimum jump height as fraction of full jump
	
	// Gravity and falling
	Gravity            int     // Gravity acceleration per frame
	MaxFallSpeed       int     // Terminal velocity in physics units
	FastFallMultiplier float64 // Gravity multiplier when holding down
}

/*
ParallaxLayer represents a single layer in the parallax background system.
Each layer can have different scroll speeds, depths, and visual effects
to create immersive layered depth in the game world.
*/
type ParallaxLayer struct {
	Speed       float64 // Scroll speed relative to camera (0-1, where 1 = same as camera)
	Depth       float64 // Depth for visual effects (0-1, where 0 = background, 1 = foreground)
	Image       string  // Path to the image file for this layer
	Alpha       float64 // Transparency (0-1, where 1 = opaque)
	Scale       float64 // Scale factor for the layer
	OffsetX     float64 // Static horizontal offset
	OffsetY     float64 // Static vertical offset
	Repeatable  bool    // Whether the layer should tile/repeat
}

/*
Config holds all the adjustable game settings in one place.
This struct centralizes all configuration values for easy tweaking
and different game modes or difficulties.
All values use consistent units and naming conventions.
*/
type Config struct {
	// Window settings
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
	
	// Rendering settings
	TileSize         int     // Base tile size from tilemap (usually 16)
	TileScaleFactor  float64 // How much to scale tiles (1.0 = 16px, 2.0 = 32px)
	CharScaleFactor  float64 // How much to scale character sprites
	
	// Camera settings
	CameraSmoothing  float64 // 0-1, higher = smoother/slower camera
	CameraDeadZoneX  float64 // Percentage of screen width for dead zone
	CameraDeadZoneY  float64 // Percentage of screen height for dead zone
	CameraMarginLeft int     // Pixels reserved on left for HUD
	CameraMarginRight int    // Pixels reserved on right for HUD
	CameraMarginTop   int    // Pixels reserved on top for HUD
	CameraMarginBottom int   // Pixels reserved on bottom for HUD
	ParallaxFactor    float64 // Background scroll speed relative to camera (0-1)
	
	// Enhanced Parallax/Depth Settings
	ParallaxLayers    []ParallaxLayer // Multiple background/foreground layers
	EnableDepthOfField bool           // Enable blur/transparency effects
	DepthBlurStrength  float64        // Strength of depth blur (0-1)
	
	// Player physics configuration
	PlayerPhysics PlayerPhysicsConfig
	
	// Enemy physics settings (kept separate for now)
	Gravity      int // Gravity for enemies
	MaxFallSpeed int // Terminal velocity for enemies
	
	// Room settings
	RoomWidthTiles   int // Room width in tiles
	RoomHeightTiles  int // Room height in tiles
	GroundLevel      int // Y position of main ground (in tiles)
	
	// Debug settings
	ShowDebugInfo    bool // Show FPS and other debug text
	ShowDebugOverlay bool // Show visual debug overlays (bounding boxes, etc)
	GridColor        [4]uint8 // RGBA color for debug grid
	UsePlaceholderSprites bool // Use placeholder sprites instead of actual sprites

}

/*
DefaultConfig returns the default game configuration.
This configuration provides a balanced zoomed-out view suitable for
platformer gameplay with good visibility of the surrounding environment.
*/
func DefaultConfig() Config {
	return Config{
		// Window settings
		WindowWidth:  1920,
		WindowHeight: 1080,
		WindowTitle:  "Sword",
		
		// Rendering settings
		TileSize:         16,    // Standard 16x16 tiles
		TileScaleFactor:  1.0,   // No scaling for zoomed-out view
		CharScaleFactor:  1.0,   // Render player at ~32x32 by default
		
		// Camera settings
		CameraSmoothing:    0.05,  // Fast camera for responsiveness
		CameraDeadZoneX:    0.1,   // Small dead zone
		CameraDeadZoneY:    0.15,  // Slightly larger vertical dead zone
		CameraMarginLeft:   0,     // No UI margins
		CameraMarginRight:  0,
		CameraMarginTop:    0,
		CameraMarginBottom: 0,
		ParallaxFactor:     0.5,   // Half-speed background scrolling
		
		// Enhanced Parallax/Depth Settings
		ParallaxLayers:    []ParallaxLayer{},
		EnableDepthOfField: false,
		DepthBlurStrength:  0.0,
		
		// Player physics configuration
		PlayerPhysics: PlayerPhysicsConfig{
			// Sprite dimensions (32x32 base sprite)
			SpriteWidth:  32,
			SpriteHeight: 32,
			
			// Collision box (centered horizontally, bottom-aligned)
			CollisionBoxOffsetX: 0.25,  // 25% from left = centered for 50% width
			CollisionBoxOffsetY: 0.5,   // 50% from top
			CollisionBoxWidth:   0.5,   // 50% of sprite width
			CollisionBoxHeight:  0.5,   // 50% of sprite height
			
			// Ground detection
			GroundCheckOffset: 2,       // 2 pixels below collision box
			GroundCheckWidth:  0.8,     // 80% of collision box width
			
			// Movement physics
			MoveSpeed:   2,             // Moderate movement speed
			JumpPower:   8,             // Good jump height (~2 tiles)
			AirControl:  0.85,          // 85% control in air
			Friction:    1,             // Quick stopping
			AirFriction: 0,             // No air friction
			
			// Jump mechanics
			CoyoteTime:         6,      // 6 frames (0.1 seconds at 60fps)
			JumpBufferTime:     10,     // 10 frames buffer
			VariableJumpHeight: true,   // Can control jump height
			MinJumpHeight:      0.5,    // 50% minimum jump
			
			// Gravity and falling
			Gravity:            1,      // Moderate gravity
			MaxFallSpeed:       12,     // Terminal velocity
			FastFallMultiplier: 1.75,   // 75% faster when holding down
		},
		
		// Enemy physics settings (kept separate for now)
		Gravity:      1,     // Same as player gravity
		MaxFallSpeed: 12,    // Same as player terminal velocity
		
		// Room settings
		RoomWidthTiles:   16,   // Expanded demo rooms: 16 tiles wide
		RoomHeightTiles:  10,   // Expanded demo rooms: 10 tiles tall
		GroundLevel:      9,    // Bottom row for 10-tile tall rooms
		
		// Debug settings
		ShowDebugInfo:    true,  // Enable debug info by default for development
		ShowDebugOverlay: true,  // Debug overlay ON by default
		GridColor:        [4]uint8{128, 128, 128, 64}, // Faint gray grid
		UsePlaceholderSprites: true, // Use placeholder sprites by default for debugging

	}
}



// GameConfig is the global configuration instance
// Initialize with default config, can be changed at runtime
var GameConfig = DefaultConfig()

// Legacy constants removed - use GameConfig and GetPhysicsUnit() instead

/*
GetPhysicsUnit returns the base physics unit size in pixels.
This is the fundamental unit for physics and tile math and is independent
of the render scale. It is equal to the base TileSize.

Returns the physics unit size as an integer number of pixels.
*/
func GetPhysicsUnit() int {
	return GameConfig.TileSize
}

/*
SetConfig updates the global game configuration.
Allows switching between different configuration presets or
applying custom configuration changes at runtime.

Parameters:
  - config: The new configuration to apply globally
*/
func SetConfig(config Config) {
	GameConfig = config
}
