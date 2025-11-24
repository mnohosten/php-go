package goext

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"
	"time"
)

// PluginInfo holds information about a loaded plugin.
type PluginInfo struct {
	Path      string
	Plugin    *plugin.Plugin
	Extension Extension
	Loaded    bool
	LoadError error
	ModTime   time.Time // Last modification time of the plugin file
}

// PluginManager manages dynamically loaded Go plugins.
// Plugins are shared libraries (.so on Linux/macOS, .dll on Windows)
// that export an Extension via a well-known symbol.
type PluginManager struct {
	mu           sync.RWMutex
	plugins      map[string]*PluginInfo
	extMgr       *ExtensionManager
	watchStop    chan struct{}
	watchEnabled bool
}

// Global plugin manager
var (
	globalPluginManager     *PluginManager
	globalPluginManagerOnce sync.Once
)

// GetGlobalPluginManager returns the singleton plugin manager.
func GetGlobalPluginManager() *PluginManager {
	globalPluginManagerOnce.Do(func() {
		globalPluginManager = NewPluginManager(GetGlobalExtensionManager())
	})
	return globalPluginManager
}

// NewPluginManager creates a new plugin manager.
func NewPluginManager(extMgr *ExtensionManager) *PluginManager {
	return &PluginManager{
		plugins:   make(map[string]*PluginInfo),
		extMgr:    extMgr,
		watchStop: make(chan struct{}),
	}
}

// LoadPlugin loads a Go plugin from the specified path.
// The plugin must export a symbol named "Extension" that implements the Extension interface.
//
// Expected plugin structure:
//
//	package main
//
//	import "github.com/krizos/php-go/pkg/goext"
//
//	type MyExtension struct {
//	    *goext.BaseExtension
//	}
//
//	// Extension is the exported symbol
//	var Extension goext.Extension = &MyExtension{...}
func (pm *PluginManager) LoadPlugin(path string) error {
	// Normalize path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving plugin path: %w", err)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if already loaded
	if info, exists := pm.plugins[absPath]; exists {
		if info.Loaded {
			return fmt.Errorf("plugin '%s' already loaded", absPath)
		}
		// If previously failed, allow retry
	}

	// Get file modification time
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("stating plugin file '%s': %w", absPath, err)
	}
	modTime := fileInfo.ModTime()

	// Load the plugin
	p, err := plugin.Open(absPath)
	if err != nil {
		info := &PluginInfo{
			Path:      absPath,
			LoadError: err,
			ModTime:   modTime,
		}
		pm.plugins[absPath] = info
		return fmt.Errorf("loading plugin '%s': %w", absPath, err)
	}

	// Look up the Extension symbol
	sym, err := p.Lookup("Extension")
	if err != nil {
		info := &PluginInfo{
			Path:      absPath,
			Plugin:    p,
			LoadError: err,
			ModTime:   modTime,
		}
		pm.plugins[absPath] = info
		return fmt.Errorf("plugin '%s' does not export 'Extension' symbol: %w", absPath, err)
	}

	// Type assert to Extension interface
	ext, ok := sym.(Extension)
	if !ok {
		// Try pointer to Extension
		extPtr, ptrOk := sym.(*Extension)
		if ptrOk && extPtr != nil {
			ext = *extPtr
		} else {
			err := fmt.Errorf("symbol 'Extension' is not of type goext.Extension (got %T)", sym)
			info := &PluginInfo{
				Path:      absPath,
				Plugin:    p,
				LoadError: err,
				ModTime:   modTime,
			}
			pm.plugins[absPath] = info
			return err
		}
	}

	// Register extension with extension manager
	if err := pm.extMgr.Register(ext); err != nil {
		info := &PluginInfo{
			Path:      absPath,
			Plugin:    p,
			Extension: ext,
			LoadError: err,
			ModTime:   modTime,
		}
		pm.plugins[absPath] = info
		return fmt.Errorf("registering extension from plugin '%s': %w", absPath, err)
	}

	// Store plugin info
	pm.plugins[absPath] = &PluginInfo{
		Path:      absPath,
		Plugin:    p,
		Extension: ext,
		Loaded:    true,
		ModTime:   modTime,
	}

	return nil
}

// LoadPluginByName loads a plugin and returns its extension name.
// This is a convenience method that loads the plugin and returns the extension name.
func (pm *PluginManager) LoadPluginByName(path string) (string, error) {
	if err := pm.LoadPlugin(path); err != nil {
		return "", err
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	absPath, _ := filepath.Abs(path)
	info, exists := pm.plugins[absPath]
	if !exists || !info.Loaded {
		return "", fmt.Errorf("plugin not loaded")
	}

	return info.Extension.Name(), nil
}

// GetPlugin returns plugin info by path.
func (pm *PluginManager) GetPlugin(path string) (*PluginInfo, bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, false
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, ok := pm.plugins[absPath]
	return info, ok
}

// GetPluginByExtension returns plugin info by extension name.
func (pm *PluginManager) GetPluginByExtension(extName string) (*PluginInfo, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, info := range pm.plugins {
		if info.Loaded && info.Extension.Name() == extName {
			return info, true
		}
	}
	return nil, false
}

// IsLoaded checks if a plugin is loaded by path.
func (pm *PluginManager) IsLoaded(path string) bool {
	info, ok := pm.GetPlugin(path)
	return ok && info.Loaded
}

// List returns all loaded plugin paths.
func (pm *PluginManager) List() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	paths := make([]string, 0, len(pm.plugins))
	for path, info := range pm.plugins {
		if info.Loaded {
			paths = append(paths, path)
		}
	}
	return paths
}

// ListAll returns all plugin paths (including failed loads).
func (pm *PluginManager) ListAll() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	paths := make([]string, 0, len(pm.plugins))
	for path := range pm.plugins {
		paths = append(paths, path)
	}
	return paths
}

// Count returns the number of successfully loaded plugins.
func (pm *PluginManager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	count := 0
	for _, info := range pm.plugins {
		if info.Loaded {
			count++
		}
	}
	return count
}

// UnloadPlugin unloads a plugin and unregisters its extension.
// Note: Go's plugin package doesn't support true unloading, so this
// only removes it from our tracking and unregisters the extension.
func (pm *PluginManager) UnloadPlugin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving plugin path: %w", err)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	info, exists := pm.plugins[absPath]
	if !exists {
		return fmt.Errorf("plugin '%s' not loaded", absPath)
	}

	if info.Loaded && info.Extension != nil {
		// Unregister extension
		if err := pm.extMgr.Unregister(info.Extension.Name()); err != nil {
			return fmt.Errorf("unregistering extension: %w", err)
		}
	}

	delete(pm.plugins, absPath)
	return nil
}

// Clear removes all plugins.
// Note: This doesn't actually unload the shared libraries (Go limitation).
func (pm *PluginManager) Clear() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins = make(map[string]*PluginInfo)
}

// LoadPlugin loads a plugin using the global plugin manager.
func LoadPlugin(path string) error {
	return GetGlobalPluginManager().LoadPlugin(path)
}

// GetPluginInfo returns plugin info by path using the global plugin manager.
func GetPluginInfo(path string) (*PluginInfo, bool) {
	return GetGlobalPluginManager().GetPlugin(path)
}

// ListPlugins returns all loaded plugins using the global plugin manager.
func ListPlugins() []string {
	return GetGlobalPluginManager().List()
}

// EnableHotReload enables automatic reloading of plugins when their files change.
// It starts a background goroutine that checks for file modifications at the specified interval.
// When a plugin file is modified, it will be unloaded and reloaded automatically.
//
// Note: Due to Go's plugin system limitations, the old plugin code remains in memory
// even after reloading. The reload operation only updates the registered extension.
//
// interval: How often to check for file changes (e.g., 1*time.Second)
func (pm *PluginManager) EnableHotReload(interval time.Duration) {
	pm.mu.Lock()
	if pm.watchEnabled {
		pm.mu.Unlock()
		return // Already enabled
	}
	pm.watchEnabled = true
	pm.mu.Unlock()

	go pm.watchPlugins(interval)
}

// DisableHotReload stops the hot reload watcher.
func (pm *PluginManager) DisableHotReload() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.watchEnabled {
		return
	}

	pm.watchEnabled = false
	close(pm.watchStop)
	pm.watchStop = make(chan struct{})
}

// watchPlugins is the internal goroutine that monitors plugin files for changes.
func (pm *PluginManager) watchPlugins(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.checkPluginsForChanges()
		case <-pm.watchStop:
			return
		}
	}
}

// checkPluginsForChanges checks all loaded plugins for file modifications.
func (pm *PluginManager) checkPluginsForChanges() {
	pm.mu.RLock()
	// Create a snapshot of plugin paths to check
	pluginsToCheck := make(map[string]*PluginInfo)
	for path, info := range pm.plugins {
		if info.Loaded {
			pluginsToCheck[path] = info
		}
	}
	pm.mu.RUnlock()

	// Check each plugin outside the lock to avoid blocking
	for path, info := range pluginsToCheck {
		fileInfo, err := os.Stat(path)
		if err != nil {
			// File no longer exists or is inaccessible
			continue
		}

		modTime := fileInfo.ModTime()
		if modTime.After(info.ModTime) {
			// File has been modified, attempt to reload
			pm.reloadPlugin(path)
		}
	}
}

// reloadPlugin attempts to reload a plugin that has been modified.
func (pm *PluginManager) reloadPlugin(path string) error {
	// Note: Due to Go's plugin limitations, we can't truly unload the old code.
	// What we do instead is unregister the old extension and register the new one.
	// The old code remains in memory but is no longer accessible.

	// Get the old extension name before unloading
	pm.mu.RLock()
	info, exists := pm.plugins[path]
	if !exists || !info.Loaded {
		pm.mu.RUnlock()
		return fmt.Errorf("plugin not loaded: %s", path)
	}
	oldExtName := info.Extension.Name()
	pm.mu.RUnlock()

	// Unload the old plugin (removes from tracking and unregisters extension)
	if err := pm.UnloadPlugin(path); err != nil {
		return fmt.Errorf("unloading plugin for reload: %w", err)
	}

	// Load the new version
	if err := pm.LoadPlugin(path); err != nil {
		return fmt.Errorf("reloading plugin: %w", err)
	}

	// Verify the extension name hasn't changed
	pm.mu.RLock()
	newInfo, _ := pm.plugins[path]
	if newInfo.Extension.Name() != oldExtName {
		pm.mu.RUnlock()
		return fmt.Errorf("plugin extension name changed from '%s' to '%s' during reload",
			oldExtName, newInfo.Extension.Name())
	}
	pm.mu.RUnlock()

	return nil
}

// ReloadPlugin manually reloads a specific plugin.
// This can be called even when hot reload is not enabled.
func (pm *PluginManager) ReloadPlugin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving plugin path: %w", err)
	}
	return pm.reloadPlugin(absPath)
}

// IsHotReloadEnabled returns whether hot reload is currently enabled.
func (pm *PluginManager) IsHotReloadEnabled() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.watchEnabled
}
