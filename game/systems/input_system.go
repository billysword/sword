package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

// InputSystem handles all player input and translates it to game actions.
// Manages keyboard input for movement, jumping, and other player actions.
type InputSystem struct {
	player            *entities.Player
	roomTransitionMgr *world.RoomTransitionManager
	pauseRequested    bool
	settingsRequested bool

	// UI toggle callback to avoid coupling with HUD
	OnToggleMinimap func()
}

// NewInputSystem creates a new input system instance.
// Parameters:
//   - player: The player entity to control
//   - roomTransitionMgr: Manager for handling room transitions
func NewInputSystem(player *entities.Player, roomTransitionMgr *world.RoomTransitionManager) *InputSystem {
	return &InputSystem{
		player:            player,
		roomTransitionMgr: roomTransitionMgr,
		pauseRequested:    false,
		settingsRequested: false,
	}
}

func (is *InputSystem) GetName() string {
	return "Input"
}

func (is *InputSystem) Update() error {
	// Handle room transitions first
	if is.roomTransitionMgr != nil {
		// Check for room transitions
		is.roomTransitionMgr.CheckTransitions(is.player, ebiten.IsKeyPressed(ebiten.KeyE))
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			is.logKeyPress("E (Interact)")
		}
	}

	// Handle movement inputs - Player handles its own input
	is.player.ProcessInput()

	// Pause request handling
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		is.pauseRequested = true
		is.logKeyPress("Escape (Pause)")
	}

	// Minimap toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		is.logKeyPress("M (Toggle Minimap)")
		if is.OnToggleMinimap != nil {
			is.OnToggleMinimap()
		}
	}

	// Log movement and action keys
	keys := []ebiten.Key{
		ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp, ebiten.KeyDown,
		ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeySpace,
	}
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			is.logKeyPress(k.String())
		}
	}

	return nil
}

func (is *InputSystem) logKeyPress(desc string) {
	playerX, playerY := is.player.GetPosition()
	roomName := ""
	if is.roomTransitionMgr != nil {
		if currentRoom := is.roomTransitionMgr.GetCurrentRoom(); currentRoom != nil {
			roomName = currentRoom.GetZoneID()
		}
	}
	engine.LogPlayerInput(desc, playerX, playerY, roomName)
}

func (is *InputSystem) HasPauseRequest() bool {
	return is.pauseRequested
}

func (is *InputSystem) HasSettingsRequest() bool {
	return is.settingsRequested
}

func (is *InputSystem) ClearRequests() {
	is.pauseRequested = false
	is.settingsRequested = false
}
