package engine

import (
	"fmt"
	"reflect"
	"sync"
)

// ServiceContainer manages dependency injection
type ServiceContainer struct {
	services map[string]interface{}
	factories map[string]func() interface{}
	singletons map[string]interface{}
	mutex sync.RWMutex
}

// NewServiceContainer creates a new service container
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services:   make(map[string]interface{}),
		factories:  make(map[string]func() interface{}),
		singletons: make(map[string]interface{}),
	}
}

// Register registers a service instance
func (sc *ServiceContainer) Register(name string, service interface{}) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.services[name] = service
}

// RegisterFactory registers a factory function for creating services
func (sc *ServiceContainer) RegisterFactory(name string, factory func() interface{}) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.factories[name] = factory
}

// RegisterSingleton registers a singleton service that will be created once
func (sc *ServiceContainer) RegisterSingleton(name string, factory func() interface{}) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.factories[name] = factory
}

// Get retrieves a service by name
func (sc *ServiceContainer) Get(name string) (interface{}, error) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	
	// Check for direct service registration
	if service, exists := sc.services[name]; exists {
		return service, nil
	}
	
	// Check for singleton
	if singleton, exists := sc.singletons[name]; exists {
		return singleton, nil
	}
	
	// Check for factory
	if factory, exists := sc.factories[name]; exists {
		sc.mutex.RUnlock()
		sc.mutex.Lock()
		
		// Double-check for singleton creation
		if singleton, exists := sc.singletons[name]; exists {
			sc.mutex.Unlock()
			sc.mutex.RLock()
			return singleton, nil
		}
		
		// Create new instance
		service := factory()
		sc.singletons[name] = service
		sc.mutex.Unlock()
		sc.mutex.RLock()
		return service, nil
	}
	
	return nil, fmt.Errorf("service %s not found", name)
}

// MustGet retrieves a service and panics if not found
func (sc *ServiceContainer) MustGet(name string) interface{} {
	service, err := sc.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// GetTyped retrieves a service and casts it to the specified type
func (sc *ServiceContainer) GetTyped(name string, target interface{}) error {
	service, err := sc.Get(name)
	if err != nil {
		return err
	}
	
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	
	serviceValue := reflect.ValueOf(service)
	targetType := targetValue.Elem().Type()
	
	if !serviceValue.Type().ConvertibleTo(targetType) {
		return fmt.Errorf("service %s cannot be converted to %s", name, targetType)
	}
	
	targetValue.Elem().Set(serviceValue.Convert(targetType))
	return nil
}

// Exists checks if a service is registered
func (sc *ServiceContainer) Exists(name string) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	
	_, hasService := sc.services[name]
	_, hasFactory := sc.factories[name]
	_, hasSingleton := sc.singletons[name]
	
	return hasService || hasFactory || hasSingleton
}

// GameContext holds all game dependencies
type GameContext struct {
	Container    *ServiceContainer
	Config       *Config
	Logger       Logger
	SpriteManager *SpriteManager
	Camera       *Camera
	HUDManager   *HUDManager
}

// NewGameContext creates a new game context with default services
func NewGameContext() *GameContext {
	container := NewServiceContainer()
	config := DefaultConfig()
	
	ctx := &GameContext{
		Container: container,
		Config:    config,
	}
	
	// Register core services
	ctx.registerCoreServices()
	
	return ctx
}

// registerCoreServices registers the core game services
func (gc *GameContext) registerCoreServices() {
	// Register config
	gc.Container.Register("config", gc.Config)
	
	// Register logger factory
	gc.Container.RegisterSingleton("logger", func() interface{} {
		// Create logger if not already exists
		if err := InitLogger("game.log"); err != nil {
			panic(fmt.Sprintf("Failed to initialize logger: %v", err))
		}
		return &defaultLogger{} // Implement proper logger interface
	})
	
	// Register sprite manager factory
	gc.Container.RegisterSingleton("sprite_manager", func() interface{} {
		InitSpriteManager()
		return GetSpriteManager()
	})
	
	// Register camera factory
	gc.Container.RegisterFactory("camera", func() interface{} {
		return NewCamera(gc.Config.WindowWidth, gc.Config.WindowHeight)
	})
	
	// Register HUD manager factory
	gc.Container.RegisterFactory("hud_manager", func() interface{} {
		return NewHUDManager()
	})
	
	// Register viewport renderer factory
	gc.Container.RegisterFactory("viewport_renderer", func() interface{} {
		return NewViewportRenderer(gc.Config.WindowWidth, gc.Config.WindowHeight)
	})
}

// GetConfig returns the configuration
func (gc *GameContext) GetConfig() *Config {
	return gc.Config
}

// GetLogger returns the logger
func (gc *GameContext) GetLogger() Logger {
	if gc.Logger == nil {
		logger, err := gc.Container.Get("logger")
		if err != nil {
			panic(err)
		}
		gc.Logger = logger.(Logger)
	}
	return gc.Logger
}

// GetSpriteManager returns the sprite manager
func (gc *GameContext) GetSpriteManager() *SpriteManager {
	if gc.SpriteManager == nil {
		sm, err := gc.Container.Get("sprite_manager")
		if err != nil {
			panic(err)
		}
		gc.SpriteManager = sm.(*SpriteManager)
	}
	return gc.SpriteManager
}

// GetCamera returns a new camera instance
func (gc *GameContext) GetCamera() *Camera {
	camera, err := gc.Container.Get("camera")
	if err != nil {
		panic(err)
	}
	return camera.(*Camera)
}

// GetHUDManager returns a new HUD manager instance
func (gc *GameContext) GetHUDManager() *HUDManager {
	hud, err := gc.Container.Get("hud_manager")
	if err != nil {
		panic(err)
	}
	return hud.(*HUDManager)
}

// GetViewportRenderer returns a new viewport renderer instance
func (gc *GameContext) GetViewportRenderer() *ViewportRenderer {
	vp, err := gc.Container.Get("viewport_renderer")
	if err != nil {
		panic(err)
	}
	return vp.(*ViewportRenderer)
}

// Logger interface for dependency injection
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	PlayerInput(action string, x, y int, room string)
}

// defaultLogger implements the Logger interface
type defaultLogger struct{}

func (dl *defaultLogger) Info(msg string) {
	LogInfo(msg)
}

func (dl *defaultLogger) Error(msg string) {
	LogInfo("ERROR: " + msg) // Using LogInfo since we don't have LogError
}

func (dl *defaultLogger) Debug(msg string) {
	LogInfo("DEBUG: " + msg)
}

func (dl *defaultLogger) PlayerInput(action string, x, y int, room string) {
	LogPlayerInput(action, x, y, room)
}

// SystemContextManager manages contexts for different systems
type SystemContextManager struct {
	gameContext *GameContext
	contexts    map[string]*ServiceContainer
}

// NewSystemContextManager creates a new system context manager
func NewSystemContextManager(gameContext *GameContext) *SystemContextManager {
	return &SystemContextManager{
		gameContext: gameContext,
		contexts:    make(map[string]*ServiceContainer),
	}
}

// CreateSystemContext creates a context for a specific system
func (scm *SystemContextManager) CreateSystemContext(systemName string) *ServiceContainer {
	context := NewServiceContainer()
	
	// Copy core services from game context
	context.Register("config", scm.gameContext.Config)
	context.Register("logger", scm.gameContext.GetLogger())
	context.Register("sprite_manager", scm.gameContext.GetSpriteManager())
	
	scm.contexts[systemName] = context
	return context
}

// GetSystemContext returns a system context by name
func (scm *SystemContextManager) GetSystemContext(systemName string) *ServiceContainer {
	if context, exists := scm.contexts[systemName]; exists {
		return context
	}
	return scm.CreateSystemContext(systemName)
}

// ConfigurableSystem represents a system that can be configured via DI
type ConfigurableSystem interface {
	Configure(context *ServiceContainer) error
	GetName() string
}

// DIGameSystem extends GameSystem with dependency injection support
type DIGameSystem interface {
	GameSystem
	ConfigurableSystem
}

// Enhanced systems with dependency injection
type DIInputSystem struct {
	*InputSystem
	context *ServiceContainer
	logger  Logger
	config  *Config
}

// NewDIInputSystem creates a new DI-enabled input system
func NewDIInputSystem(context *ServiceContainer) *DIInputSystem {
	logger, _ := context.Get("logger")
	config, _ := context.Get("config")
	
	// Create base input system (will need to be modified to accept DI)
	return &DIInputSystem{
		context: context,
		logger:  logger.(Logger),
		config:  config.(*Config),
	}
}

// Configure configures the system with dependencies
func (dis *DIInputSystem) Configure(context *ServiceContainer) error {
	dis.context = context
	
	// Get required dependencies
	logger, err := context.Get("logger")
	if err != nil {
		return fmt.Errorf("logger dependency not found: %w", err)
	}
	dis.logger = logger.(Logger)
	
	config, err := context.Get("config")
	if err != nil {
		return fmt.Errorf("config dependency not found: %w", err)
	}
	dis.config = config.(*Config)
	
	return nil
}

// DIPhysicsSystem with dependency injection
type DIPhysicsSystem struct {
	*PhysicsSystem
	context *ServiceContainer
	config  *Config
	logger  Logger
}

// NewDIPhysicsSystem creates a new DI-enabled physics system
func NewDIPhysicsSystem(context *ServiceContainer) *DIPhysicsSystem {
	config, _ := context.Get("config")
	logger, _ := context.Get("logger")
	
	return &DIPhysicsSystem{
		context: context,
		config:  config.(*Config),
		logger:  logger.(Logger),
	}
}

// Configure configures the physics system with dependencies
func (dps *DIPhysicsSystem) Configure(context *ServiceContainer) error {
	dps.context = context
	
	config, err := context.Get("config")
	if err != nil {
		return fmt.Errorf("config dependency not found: %w", err)
	}
	dps.config = config.(*Config)
	
	logger, err := context.Get("logger")
	if err != nil {
		return fmt.Errorf("logger dependency not found: %w", err)
	}
	dps.logger = logger.(Logger)
	
	return nil
}

// ServiceRegistry provides a global registry for services (transitional helper)
type ServiceRegistry struct {
	container *ServiceContainer
}

var globalRegistry *ServiceRegistry
var registryOnce sync.Once

// GetGlobalRegistry returns the global service registry
func GetGlobalRegistry() *ServiceRegistry {
	registryOnce.Do(func() {
		globalRegistry = &ServiceRegistry{
			container: NewServiceContainer(),
		}
	})
	return globalRegistry
}

// RegisterGlobalService registers a service in the global registry
func RegisterGlobalService(name string, service interface{}) {
	GetGlobalRegistry().container.Register(name, service)
}

// GetGlobalService retrieves a service from the global registry
func GetGlobalService(name string) (interface{}, error) {
	return GetGlobalRegistry().container.Get(name)
}

// Migration helpers for transitioning away from global state
func MigrateGlobalState(gameContext *GameContext) {
	// Register current global state in the DI container for gradual migration
	if globalLeftSprite != nil {
		gameContext.Container.Register("left_sprite", globalLeftSprite)
	}
	if globalRightSprite != nil {
		gameContext.Container.Register("right_sprite", globalRightSprite)
	}
	if globalIdleSprite != nil {
		gameContext.Container.Register("idle_sprite", globalIdleSprite)
	}
	if globalBackgroundImage != nil {
		gameContext.Container.Register("background_image", globalBackgroundImage)
	}
	if globalTileSprite != nil {
		gameContext.Container.Register("tile_sprite", globalTileSprite)
	}
	if globalTilesSprite != nil {
		gameContext.Container.Register("tiles_sprite", globalTilesSprite)
	}
}