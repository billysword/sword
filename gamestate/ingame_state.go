package gamestate

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	unit    = 16
	groundY = 380
)

// Character represents the player character
type Character struct {
	x  int
	y  int
	vx int
	vy int
}

func (c *Character) tryJump() {
	c.vy = -10 * unit
}

func (c *Character) update() {
	c.x += c.vx
	c.y += c.vy
	if c.y > groundY*unit {
		c.y = groundY * unit
	}
	if c.vx > 0 {
		c.vx -= 4
	} else if c.vx < 0 {
		c.vx += 4
	}
	if c.vy < 20*unit {
		c.vy += 8
	}
}

func (c *Character) draw(screen *ebiten.Image) {
	s := globalIdleSprite
	switch {
	case c.vx > 0:
		s = globalRightSprite
	case c.vx < 0:
		s = globalLeftSprite
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(float64(c.x)/unit, float64(c.y)/unit)
	screen.DrawImage(s, op)
}

// InGameState represents the actual gameplay state
type InGameState struct {
	stateManager *StateManager
	character    *Character
}

// NewInGameState creates a new in-game state
func NewInGameState(sm *StateManager) *InGameState {
	return &InGameState{
		stateManager: sm,
		character:    &Character{x: 50 * unit, y: groundY * unit},
	}
}

// Update handles game logic and input
func (ig *InGameState) Update() error {
	// Check for pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ig.stateManager.ChangeState(NewPauseState(ig.stateManager, ig))
		return nil
	}

	// Character controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		ig.character.vx = -4 * unit
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		ig.character.vx = 4 * unit
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ig.character.tryJump()
	}

	ig.character.update()
	return nil
}

// Draw renders the game world
func (ig *InGameState) Draw(screen *ebiten.Image) {
	// Draw background
	if globalBackgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		screen.DrawImage(globalBackgroundImage, op)
	}

	// Draw character
	if ig.character != nil {
		ig.character.draw(screen)
	}

	// Show debug info
	msg := fmt.Sprintf("TPS: %0.2f\nPress SPACE to jump\nP/ESC - Pause", ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

// OnEnter is called when entering the game state
func (ig *InGameState) OnEnter() {
	// Reset character position or load level data
	if ig.character == nil {
		ig.character = &Character{x: 50 * unit, y: groundY * unit}
	}
}

// OnExit is called when leaving the game state
func (ig *InGameState) OnExit() {
	// Save game state or cleanup resources
}