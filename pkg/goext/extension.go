package goext

import (
	"fmt"
	"sync"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Extension is the interface that all Go extensions must implement.
// Extensions provide functions, classes, and constants to PHP.
type Extension interface {
	// Name returns the extension name (e.g., "myextension")
	Name() string

	// Version returns the extension version (e.g., "1.0.0")
	Version() string

	// Init is called when the extension is loaded into a VM.
	// Extensions can perform initialization here.
	Init(*vm.VM) error

	// Functions returns a map of function names to handlers.
	// These functions will be available in PHP.
	Functions() map[string]FunctionHandler

	// Constants returns a map of constant names to values.
	// These constants will be available in PHP.
	Constants() map[string]*types.Value
}

// ClassDefinition represents a PHP class provided by an extension.
// This is a placeholder for future class support.
type ClassDefinition struct {
	Name       string
	Methods    map[string]FunctionHandler
	Properties map[string]*types.Value
	// Future: Parent class, interfaces, etc.
}

// ExtensionInfo holds metadata about a loaded extension.
type ExtensionInfo struct {
	Extension Extension
	Loaded    bool
	InitError error
}

// ExtensionManager manages registered extensions.
type ExtensionManager struct {
	mu         sync.RWMutex
	extensions map[string]*ExtensionInfo
	registry   *FunctionRegistry
}

// Global extension manager
var (
	globalExtManager     *ExtensionManager
	globalExtManagerOnce sync.Once
)

// GetGlobalExtensionManager returns the singleton extension manager.
func GetGlobalExtensionManager() *ExtensionManager {
	globalExtManagerOnce.Do(func() {
		globalExtManager = NewExtensionManager(GetGlobalRegistry())
	})
	return globalExtManager
}

// NewExtensionManager creates a new extension manager.
func NewExtensionManager(registry *FunctionRegistry) *ExtensionManager {
	return &ExtensionManager{
		extensions: make(map[string]*ExtensionInfo),
		registry:   registry,
	}
}

// RegisterExtension registers an extension with the global manager.
func RegisterExtension(ext Extension) error {
	return GetGlobalExtensionManager().Register(ext)
}

// Register registers an extension with this manager.
// The extension is registered but not initialized until LoadIntoVM is called.
func (em *ExtensionManager) Register(ext Extension) error {
	if ext == nil {
		return fmt.Errorf("extension cannot be nil")
	}

	name := ext.Name()
	if name == "" {
		return fmt.Errorf("extension name cannot be empty")
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.extensions[name]; exists {
		return fmt.Errorf("extension '%s' already registered", name)
	}

	em.extensions[name] = &ExtensionInfo{
		Extension: ext,
		Loaded:    false,
	}

	return nil
}

// LoadIntoVM loads an extension into a VM.
// This calls the extension's Init() and registers its functions and constants.
func (em *ExtensionManager) LoadIntoVM(name string, v *vm.VM) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	info, exists := em.extensions[name]
	if !exists {
		return fmt.Errorf("extension '%s' not registered", name)
	}

	ext := info.Extension

	// Call extension's Init
	if err := ext.Init(v); err != nil {
		info.InitError = err
		return fmt.Errorf("extension '%s' init failed: %w", name, err)
	}

	// Register functions
	functions := ext.Functions()
	for fnName, handler := range functions {
		// Create a wrapper that matches the registered function signature
		wrappedFn := func(args ...interface{}) (interface{}, error) {
			// Convert interface{} args to []*types.Value
			phpArgs := make([]*types.Value, len(args))
			marshaler := NewMarshaler()
			for i, arg := range args {
				phpVal, err := marshaler.ToPHP(arg)
				if err != nil {
					return nil, err
				}
				phpArgs[i] = phpVal
			}

			// Call handler
			result, err := handler(v, phpArgs)
			if err != nil {
				return nil, err
			}

			// Convert result back to interface{}
			return marshaler.ToGo(result)
		}

		// Register with the function registry
		if err := em.registry.Register(fnName, wrappedFn); err != nil {
			return fmt.Errorf("registering function '%s': %w", fnName, err)
		}
	}

	// Register constants
	constants := ext.Constants()
	for constName, value := range constants {
		// Constants are registered in the VM's global scope
		v.SetGlobal(constName, value)
	}

	info.Loaded = true
	return nil
}

// LoadAllIntoVM loads all registered extensions into a VM.
func (em *ExtensionManager) LoadAllIntoVM(v *vm.VM) error {
	em.mu.RLock()
	names := make([]string, 0, len(em.extensions))
	for name := range em.extensions {
		names = append(names, name)
	}
	em.mu.RUnlock()

	for _, name := range names {
		if err := em.LoadIntoVM(name, v); err != nil {
			return err
		}
	}

	return nil
}

// Get retrieves extension info by name.
func (em *ExtensionManager) Get(name string) (*ExtensionInfo, bool) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	info, ok := em.extensions[name]
	return info, ok
}

// Has checks if an extension is registered.
func (em *ExtensionManager) Has(name string) bool {
	em.mu.RLock()
	defer em.mu.RUnlock()
	_, ok := em.extensions[name]
	return ok
}

// List returns all registered extension names.
func (em *ExtensionManager) List() []string {
	em.mu.RLock()
	defer em.mu.RUnlock()

	names := make([]string, 0, len(em.extensions))
	for name := range em.extensions {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered extensions.
func (em *ExtensionManager) Count() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.extensions)
}

// Unregister removes an extension from the manager.
func (em *ExtensionManager) Unregister(name string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.extensions[name]; !exists {
		return fmt.Errorf("extension '%s' not registered", name)
	}

	delete(em.extensions, name)
	return nil
}

// Clear removes all extensions.
func (em *ExtensionManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.extensions = make(map[string]*ExtensionInfo)
}

// BaseExtension provides a base implementation of the Extension interface.
// Extensions can embed this struct and override methods as needed.
type BaseExtension struct {
	name      string
	version   string
	functions map[string]FunctionHandler
	constants map[string]*types.Value
}

// NewBaseExtension creates a new base extension.
func NewBaseExtension(name, version string) *BaseExtension {
	return &BaseExtension{
		name:      name,
		version:   version,
		functions: make(map[string]FunctionHandler),
		constants: make(map[string]*types.Value),
	}
}

// Name returns the extension name.
func (e *BaseExtension) Name() string {
	return e.name
}

// Version returns the extension version.
func (e *BaseExtension) Version() string {
	return e.version
}

// Init performs initialization (default: no-op).
func (e *BaseExtension) Init(v *vm.VM) error {
	return nil
}

// Functions returns registered functions.
func (e *BaseExtension) Functions() map[string]FunctionHandler {
	return e.functions
}

// Constants returns registered constants.
func (e *BaseExtension) Constants() map[string]*types.Value {
	return e.constants
}

// AddFunction adds a function to the extension.
func (e *BaseExtension) AddFunction(name string, handler FunctionHandler) {
	e.functions[name] = handler
}

// AddConstant adds a constant to the extension.
func (e *BaseExtension) AddConstant(name string, value *types.Value) {
	e.constants[name] = value
}

// SetGlobal is a helper to set a global variable in the VM.
func SetGlobal(v *vm.VM, name string, value *types.Value) {
	v.SetGlobal(name, value)
}

// GetGlobal is a helper to get a global variable from the VM.
func GetGlobal(v *vm.VM, name string) (*types.Value, bool) {
	return v.GetGlobal(name)
}
