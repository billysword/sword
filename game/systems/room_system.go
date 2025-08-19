package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

// RoomSystem manages room transitions and world state changes.
// Handles checking for and processing room transitions and updates
// other systems when the room changes.
type RoomSystem struct {
	transitionManager *world.RoomTransitionManager
	worldMap          *world.WorldMap
	player            *entities.Player
	physicsSystem     *PhysicsSystem
	cameraSystem      *CameraSystem
}

func NewRoomSystem(transitionManager *world.RoomTransitionManager, worldMap *world.WorldMap, player *entities.Player) *RoomSystem {
	return &RoomSystem{
		transitionManager: transitionManager,
		worldMap:          worldMap,
		player:            player,
		physicsSystem:     nil, // Will be initialized later
		cameraSystem:      nil, // Will be initialized later
	}
}

func (rs *RoomSystem) GetName() string {
	return "Room"
}

// SetPhysicsSystem sets the physics system reference.
func (rs *RoomSystem) SetPhysicsSystem(ps *PhysicsSystem) {
	rs.physicsSystem = ps
}

// SetCameraSystem sets the camera system reference.
func (rs *RoomSystem) SetCameraSystem(cs *CameraSystem) {
	rs.cameraSystem = cs
}

// Update checks for room transitions.
func (rs *RoomSystem) Update() error {
	if rs.transitionManager != nil && rs.player != nil {
		// Safety: if player fell below current room bounds, portal to safety room
		if current := rs.transitionManager.GetCurrentRoom(); current != nil {
			tm := current.GetTileMap()
			u := engine.GetPhysicsUnit()
			_, py := rs.player.GetPosition()
			if tm != nil {
				maxY := tm.Height * u
				if py > maxY+u { // allow small margin
					// Queue a transition if safety room exists
					if len(rs.transitionManager.GetSpawnPoints("safety")) > 0 {
						// Create a pending transition to safety room's default spawn
						// We use CheckTransitions/ProcessPendingTransition pathway by directly setting pending
						// Not exposed: fallback to direct spawn after SetCurrentRoom
						rs.transitionManager.SetCurrentRoom("safety")
						// Try spawn id "entry" then first spawn
						if err := rs.transitionManager.SpawnPlayerInRoom(rs.player, "safety", "entry"); err != nil {
							// ignore error, player may be placed later by fallback
						}
						// Update camera/physics to new room immediately
						if rs.physicsSystem != nil {
							rs.physicsSystem.SetRoom(rs.transitionManager.GetCurrentRoom())
						}
						if rs.cameraSystem != nil {
							rs.cameraSystem.SetRoom(rs.transitionManager.GetCurrentRoom())
						}
					}
				}
			}
		}

		// Process any pending transitions
		if rs.transitionManager.HasPendingTransition() {
			newRoom, err := rs.transitionManager.ProcessPendingTransition(rs.player)
			if err != nil {
				return fmt.Errorf("failed to process room transition: %w", err)
			}

			if newRoom != nil {
				// Recompute tile scale to fit the new room to the current window (zoom in small rooms)
				if tm := newRoom.GetTileMap(); tm != nil {
					u := engine.GetPhysicsUnit()
					winW, winH := ebiten.WindowSize()
					roomPxW := tm.Width * u
					roomPxH := tm.Height * u
					fitScaleW := float64(winW) / float64(roomPxW)
					fitScaleH := float64(winH) / float64(roomPxH)
					fitScale := math.Min(fitScaleW, fitScaleH)
					// Clamp between 1x and 4x
					if fitScale < 1.0 {
						fitScale = 1.0
					}
					if fitScale > 4.0 {
						fitScale = 4.0
					}
					engine.GameConfig.TileScaleFactor = fitScale
				}

				// Notify other systems about the room change
				if rs.physicsSystem != nil {
					rs.physicsSystem.SetRoom(newRoom)
				}
				if rs.cameraSystem != nil {
					// After updating scale, update camera bounds for new room
					rs.cameraSystem.SetRoom(newRoom)
				}
			}
		}
	}

	return nil
}

func (rs *RoomSystem) GetCurrentRoom() world.Room {
	return rs.transitionManager.GetCurrentRoom()
}
