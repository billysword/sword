package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sword/engine"
	"sword/entities"
	"sword/world"
)

// WorldMapState shows the full world map as a separate state
// It is read-only and allows returning to the game with ESC

type WorldMapState struct {
	stateManager *engine.StateManager
	player      *entities.Player
	worldMap    *world.WorldMap
}

func NewWorldMapState(sm *engine.StateManager, player *entities.Player, wm *world.WorldMap) *WorldMapState {
	return &WorldMapState{stateManager: sm, player: player, worldMap: wm}
}

func (s *WorldMapState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Return to gameplay state
		// For now, just go to a new InGameState. If we had a stack, we would pop.
		s.stateManager.ChangeState(NewInGameState(s.stateManager))
	}
	return nil
}

func (s *WorldMapState) Draw(screen *ebiten.Image) {
	// Simple full-screen rendering using a temporary ZoneMapOverlay-like approach
	// Reuse zone map implementation by creating a temporary overlay that includes all rooms
	img := screen
	w, h := ebiten.WindowSize()
	// Dim background
	ebitenutil.DrawRect(img, 0, 0, float64(w), float64(h), engine.RGBA(0, 0, 0, 200))
	// Render with a temporary helper
	world.DrawFullWorldMap(img, s.worldMap, s.player)
	// Help text
	ebitenutil.DebugPrintAt(img, "World Map (ESC to return)", 20, 20)
}

func (s *WorldMapState) OnEnter() {}
func (s *WorldMapState) OnExit()  {}