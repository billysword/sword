package systems

import "fmt"

// GameSystem defines the interface for all game systems.
// Each system manages a specific aspect of the game (input, physics, camera, etc.).
type GameSystem interface {
	GetName() string
	Update() error
}

// GameSystemManager manages all game systems and their update order.
type GameSystemManager struct {
	systems     map[string]GameSystem
	updateOrder []string
}

// NewGameSystemManager creates a new game system manager.
func NewGameSystemManager() *GameSystemManager {
	return &GameSystemManager{
		systems:     make(map[string]GameSystem),
		updateOrder: make([]string, 0),
	}
}

// AddSystem adds a new system to the manager.
func (gsm *GameSystemManager) AddSystem(name string, system GameSystem) {
	gsm.systems[name] = system
	// Add to update order if not already present
	found := false
	for _, n := range gsm.updateOrder {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		gsm.updateOrder = append(gsm.updateOrder, name)
	}
}

// GetSystem returns a system by name.
func (gsm *GameSystemManager) GetSystem(name string) GameSystem {
	return gsm.systems[name]
}

// UpdateAll updates all systems in order.
func (gsm *GameSystemManager) UpdateAll() error {
	// Update in specified order
	for _, systemName := range gsm.updateOrder {
		if system, exists := gsm.systems[systemName]; exists {
			if err := system.Update(); err != nil {
				return fmt.Errorf("system %s update failed: %w", systemName, err)
			}
		}
	}

	// Update any systems not in the order list
	for name, system := range gsm.systems {
		found := false
		for _, orderedName := range gsm.updateOrder {
			if name == orderedName {
				found = true
				break
			}
		}
		if !found {
			if err := system.Update(); err != nil {
				return fmt.Errorf("system %s update failed: %w", system.GetName(), err)
			}
		}
	}

	return nil
}

// SetUpdateOrder sets the order in which systems are updated.
func (gsm *GameSystemManager) SetUpdateOrder(order []string) {
	gsm.updateOrder = order
}
