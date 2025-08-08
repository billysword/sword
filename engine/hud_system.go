package engine

// HUDComponent represents a UI component that can be updated and drawn
type HUDComponent interface {
	Update() error
	Draw(screen interface{}) error
	IsVisible() bool
	SetVisible(visible bool)
	GetName() string
}

// HUDManager manages all HUD components using the same Update/Draw paradigm as game states
type HUDManager struct {
	components map[string]HUDComponent
	enabled    bool
}

// NewHUDManager creates a new HUD manager
func NewHUDManager() *HUDManager {
	return &HUDManager{
		components: make(map[string]HUDComponent),
		enabled:    true,
	}
}

// AddComponent adds a HUD component to the manager
func (hm *HUDManager) AddComponent(component HUDComponent) {
	hm.components[component.GetName()] = component
}

// RemoveComponent removes a HUD component from the manager
func (hm *HUDManager) RemoveComponent(name string) {
	delete(hm.components, name)
}

// GetComponent retrieves a HUD component by name
func (hm *HUDManager) GetComponent(name string) HUDComponent {
	return hm.components[name]
}

// Update updates all HUD components
func (hm *HUDManager) Update() error {
	if !hm.enabled {
		return nil
	}

	for _, component := range hm.components {
		if component.IsVisible() {
			if err := component.Update(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Draw draws all HUD components
func (hm *HUDManager) Draw(screen interface{}) error {
	if !hm.enabled {
		return nil
	}

	// If world map overlay is visible, render only it to avoid other HUD elements on top
	if overlay, exists := hm.components["world_map"]; exists && overlay.IsVisible() {
		LogDebug("DRAW_LAYER: HUDComponent(" + overlay.GetName() + ")")
		return overlay.Draw(screen)
	}

	for _, component := range hm.components {
		if component.IsVisible() {
			LogDebug("DRAW_LAYER: HUDComponent(" + component.GetName() + ")")
			if err := component.Draw(screen); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetEnabled enables or disables the entire HUD system
func (hm *HUDManager) SetEnabled(enabled bool) {
	hm.enabled = enabled
}

// IsEnabled returns whether the HUD system is enabled
func (hm *HUDManager) IsEnabled() bool {
	return hm.enabled
}

// ToggleComponent toggles visibility of a specific component
func (hm *HUDManager) ToggleComponent(name string) {
	if component, exists := hm.components[name]; exists {
		component.SetVisible(!component.IsVisible())
	}
}

// SetComponentVisible sets visibility of a specific component
func (hm *HUDManager) SetComponentVisible(name string, visible bool) {
	if component, exists := hm.components[name]; exists {
		component.SetVisible(visible)
	}
}

// GetVisibleComponents returns a list of currently visible component names
func (hm *HUDManager) GetVisibleComponents() []string {
	var visible []string
	for name, component := range hm.components {
		if component.IsVisible() {
			visible = append(visible, name)
		}
	}
	return visible
}
