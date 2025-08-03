package engine

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
and different game modes (zoomed in vs zoomed out, different difficulties, etc.).
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
	
	// Physics settings
	PlayerMoveSpeed   int // Horizontal movement speed in physics units
	PlayerJumpPower   int // Initial jump velocity in physics units
	PlayerFriction    int // Friction applied each frame
	Gravity           int // Gravity acceleration per frame
	MaxFallSpeed      int // Terminal velocity in physics units
	
	// Room settings
	RoomWidthTiles   int // Room width in tiles
	RoomHeightTiles  int // Room height in tiles
	GroundLevel      int // Y position of main ground (in tiles)
	
	// Debug settings
	ShowDebugInfo    bool // Show FPS and other debug text
	ShowDebugOverlay bool // Show visual debug overlays (bounding boxes, etc)
	GridColor        [4]uint8 // RGBA color for debug grid
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
		CharScaleFactor:  0.7,   // Better proportional match to tiles
		
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
		
		// Physics settings (lower values for zoomed-out feel)
		PlayerMoveSpeed:  2,     // Moderate movement speed
		PlayerJumpPower:  8,     // Good jump height
		PlayerFriction:   1,     // Quick stopping
		Gravity:          1,     // Moderate gravity
		MaxFallSpeed:     12,    // Terminal velocity
		
		// Room settings
		RoomWidthTiles:   80,    // Wide rooms for exploration
		RoomHeightTiles:  60,    // Tall rooms for vertical gameplay
		GroundLevel:      45,    // Better vertical layout for platforming
		
		// Debug settings
		ShowDebugInfo:    true,  // Enable debug info by default for development
		ShowDebugOverlay: false, // Debug overlay off by default
		GridColor:        [4]uint8{128, 128, 128, 64}, // Faint gray grid
	}
}

/*
ZoomedInConfig returns a configuration optimized for close-up gameplay.
This configuration provides larger sprites and tighter camera work
for detailed platformer action with more intimate level design.
*/
func ZoomedInConfig() Config {
	config := DefaultConfig()
	
	// Rendering settings for close-up view
	config.TileScaleFactor = 2.0   // Double-size tiles
	config.CharScaleFactor = 0.8   // Larger character
	
	// Camera settings for closer following
	config.CameraSmoothing = 0.1   // Slightly smoother camera
	config.CameraDeadZoneX = 0.05  // Smaller dead zone for tighter following
	config.CameraDeadZoneY = 0.1
	
	// Enhanced Parallax/Depth Settings for closer view
	config.ParallaxLayers = []ParallaxLayer{
		{Speed: 0.3, Depth: 0.8, Image: "assets/parallax/background.png"},
		{Speed: 0.6, Depth: 0.6, Image: "assets/parallax/midground.png"},
		{Speed: 1.0, Depth: 0.4, Image: "assets/parallax/foreground.png"},
	}
	config.EnableDepthOfField = true
	config.DepthBlurStrength = 0.5
	
	// Physics settings (higher values for zoomed-in responsiveness)
	config.PlayerMoveSpeed = 4     // Faster movement
	config.PlayerJumpPower = 12    // Higher jumps
	config.PlayerFriction = 2      // More responsive stopping
	config.Gravity = 1             // Same gravity
	config.MaxFallSpeed = 16       // Faster falling
	
	// Room settings for tighter level design
	config.RoomWidthTiles = 40     // Smaller, more focused rooms
	config.RoomHeightTiles = 30
	config.GroundLevel = 25        // Ground closer to middle
	
	return config
}

/*
SmallRoomConfig returns a configuration for small rooms with centered viewport.
This configuration is designed for testing room layouts with black dead areas
when the room doesn't fill the entire window. The actual room size is determined
by the tilemap dimensions, not the configuration.
*/
func SmallRoomConfig() Config {
	config := DefaultConfig()
	
	// Room dimensions will be set by the tilemap
	// These are just defaults that can be overridden
	config.RoomWidthTiles = 10
	config.RoomHeightTiles = 10
	config.GroundLevel = 8
	
	// Larger tiles for better visibility in small rooms
	config.TileScaleFactor = 2.0
	config.CharScaleFactor = 1.0
	
	// Disable camera smoothing for precise control
	config.CameraSmoothing = 0.1
	config.CameraDeadZoneX = 0.0
	config.CameraDeadZoneY = 0.0
	
	return config
}

// GameConfig is the global configuration instance
// Initialize with default config, can be changed at runtime
var GameConfig = DefaultConfig()

// Legacy constants removed - use GameConfig and GetPhysicsUnit() instead

/*
GetPhysicsUnit returns the current physics unit size in pixels.
This is the fundamental unit for all physics calculations, derived from
the tile size and scale factor.

Returns the physics unit size as an integer number of pixels.
*/
func GetPhysicsUnit() int {
	return int(float64(GameConfig.TileSize) * GameConfig.TileScaleFactor)
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