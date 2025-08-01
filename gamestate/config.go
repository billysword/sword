package gamestate

// Config holds all the adjustable game settings in one place
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
	GridColor        [4]uint8 // RGBA color for debug grid
}

// DefaultConfig returns the default game configuration
func DefaultConfig() *Config {
	return &Config{
		// Window settings - 800x450 for zoomed out view
		WindowWidth:  800,
		WindowHeight: 450,
		WindowTitle:  "Platformer (Ebitengine Demo)",
		
		// Rendering settings - 1.0 scale for zoomed out view
		TileSize:        16,
		TileScaleFactor: 1.0,  // Makes tiles 16x16 pixels
		CharScaleFactor: 0.4,  // Makes player ~3 tiles wide
		
		// Camera settings
		CameraSmoothing:    0.1,   // 10% interpolation for smooth following
		CameraDeadZoneX:    0.25,  // 25% of screen width
		CameraDeadZoneY:    0.16,  // 16% of screen height
		CameraMarginLeft:   32,
		CameraMarginRight:  32,
		CameraMarginTop:    32,
		CameraMarginBottom: 48,    // Extra space for potential HUD
		ParallaxFactor:     0.3,   // Background moves at 30% speed
		
		// Physics settings (adjusted for 16px physics unit)
		PlayerMoveSpeed: 3,
		PlayerJumpPower: 8,
		PlayerFriction:  2,
		Gravity:         4,
		MaxFallSpeed:    15,
		
		// Room settings
		RoomWidthTiles:  120,
		RoomHeightTiles: 60,
		GroundLevel:     44,  // Near bottom of 60-tile high room
		
		// Debug settings
		ShowDebugInfo: true,
		GridColor:     [4]uint8{100, 100, 100, 80}, // Faint gray
	}
}

// ZoomedInConfig returns a configuration for a more zoomed-in view
func ZoomedInConfig() *Config {
	config := DefaultConfig()
	
	// Smaller window for zoomed in view
	config.WindowWidth = 640
	config.WindowHeight = 360
	
	// Larger tile scale
	config.TileScaleFactor = 2.0  // Makes tiles 32x32 pixels
	config.CharScaleFactor = 0.5  // Smaller character relative to tiles
	
	// Adjusted physics for larger scale
	config.PlayerMoveSpeed = 4
	config.PlayerJumpPower = 10
	config.PlayerFriction = 4
	config.Gravity = 8
	config.MaxFallSpeed = 20
	
	// Smaller room since we see less
	config.RoomWidthTiles = 80
	config.RoomHeightTiles = 40
	config.GroundLevel = 29
	
	return config
}

// Global config instance
var GameConfig = DefaultConfig()

// GetPhysicsUnit returns the physics unit size based on current config
func GetPhysicsUnit() int {
	return int(float64(GameConfig.TileSize) * GameConfig.TileScaleFactor)
}