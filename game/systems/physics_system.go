package systems

import (
	"fmt"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

// PhysicsSystem handles physics simulation for all entities.
// Manages gravity, collisions, and movement for the player and enemies.
type PhysicsSystem struct {
	player  *entities.Player
	enemies []entities.Enemy
	room    world.Room
}

// NewPhysicsSystem creates a new physics system instance.
func NewPhysicsSystem(player *entities.Player) *PhysicsSystem {
	return &PhysicsSystem{
		player:  player,
		enemies: make([]entities.Enemy, 0),
	}
}

func (ps *PhysicsSystem) GetName() string {
	return "Physics"
}

func (ps *PhysicsSystem) AddEnemy(enemy entities.Enemy) {
	ps.enemies = append(ps.enemies, enemy)
}

func (ps *PhysicsSystem) SetRoom(room world.Room) {
	ps.room = room
}

func (ps *PhysicsSystem) SetCurrentRoom(room world.Room) {
	ps.SetRoom(room)
}

func (ps *PhysicsSystem) ClearEnemies() {
	ps.enemies = ps.enemies[:0]
}

// Update updates physics for all entities.
func (ps *PhysicsSystem) Update() error {
	// Log player state before physics update
	px, py := ps.player.GetPosition()
	vx, vy := ps.player.GetVelocity()
	onGround := ps.player.IsOnGround()

	roomName := "NoRoom"
	if ps.room != nil {
		roomName = ps.room.GetZoneID()
	}

	engine.LogDebug(fmt.Sprintf("PHYSICS_BEFORE: Room=%s Pos=(%d,%d) Vel=(%d,%d) OnGround=%v",
		roomName, px, py, vx, vy, onGround))

	// Update player physics with tiles when room is present
	if ps.room != nil {
		if tileProvider, ok := ps.room.(entities.TileProvider); ok {
			engine.LogDebug("PHYSICS: Using UpdateWithTileCollision")
			ps.player.UpdateWithTileCollision(tileProvider)
		} else {
			engine.LogDebug("PHYSICS: Room doesn't implement TileProvider, using basic Update")
			ps.player.Update()
		}
	} else {
		engine.LogDebug("PHYSICS: No room, using basic Update")
		ps.player.Update()
	}

	// Log player state after physics update
	px2, py2 := ps.player.GetPosition()
	vx2, vy2 := ps.player.GetVelocity()
	onGround2 := ps.player.IsOnGround()

	engine.LogDebug(fmt.Sprintf("PHYSICS_AFTER: Room=%s Pos=(%d,%d) Vel=(%d,%d) OnGround=%v",
		roomName, px2, py2, vx2, vy2, onGround2))

	// Log movement delta if any
	if px != px2 || py != py2 {
		engine.LogDebug(fmt.Sprintf("PHYSICS_DELTA: ΔPos=(%d,%d) ΔVel=(%d,%d)",
			px2-px, py2-py, vx2-vx, vy2-vy))
	}

	// Update enemies
	for _, enemy := range ps.enemies {
		enemy.Update()
	}

	// Handle collision detection for enemies if needed (player handled above)
	return nil
}
