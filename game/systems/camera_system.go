package systems

import (
	"sword/engine"
	"sword/entities"
	"sword/world"
)

// CameraSystem manages the game camera and viewport.
// Handles camera following, boundaries, and smooth transitions.
type CameraSystem struct {
	camera  *engine.Camera
	player  *entities.Player
	room    world.Room
	enabled bool
}

func NewCameraSystem(camera *engine.Camera, player *entities.Player) *CameraSystem {
	return &CameraSystem{
		camera:  camera,
		player:  player,
		enabled: true,
	}
}

func (cs *CameraSystem) GetName() string {
	return "Camera"
}

func (cs *CameraSystem) SetRoom(room world.Room) {
	cs.room = room
	// Update camera bounds when room changes
	if cs.camera != nil && room != nil {
		if tileMap := room.GetTileMap(); tileMap != nil {
			u := engine.GetPhysicsUnit()
			scale := engine.GameConfig.TileScaleFactor
			cs.camera.SetWorldBounds(int(float64(tileMap.Width*u)*scale), int(float64(tileMap.Height*u)*scale))
		}
	}
}

func (cs *CameraSystem) SetCurrentRoom(room world.Room) {
	cs.SetRoom(room)
}

// Update updates the camera to follow the player.
func (cs *CameraSystem) Update() error {
	// Get player position for camera tracking
	playerX, playerY := cs.player.GetPosition()

	// Update camera position in scaled screen pixels
	s := engine.GameConfig.TileScaleFactor
	cs.camera.Update(int(float64(playerX)*s), int(float64(playerY)*s))

	return nil
}

func (cs *CameraSystem) SetEnabled(enabled bool) {
	cs.enabled = enabled
}
