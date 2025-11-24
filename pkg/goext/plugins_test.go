package goext

import (
	"path/filepath"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// TestPluginManager tests the basic plugin manager operations
func TestPluginManager(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	if pm == nil {
		t.Fatal("NewPluginManager returned nil")
	}

	if pm.Count() != 0 {
		t.Errorf("Expected 0 plugins initially, got %d", pm.Count())
	}

	if len(pm.List()) != 0 {
		t.Errorf("Expected empty list initially, got %v", pm.List())
	}
}

// TestGetGlobalPluginManager tests the global plugin manager singleton
func TestGetGlobalPluginManager(t *testing.T) {
	pm1 := GetGlobalPluginManager()
	pm2 := GetGlobalPluginManager()

	if pm1 != pm2 {
		t.Error("GetGlobalPluginManager should return same instance")
	}

	if pm1 == nil {
		t.Fatal("GetGlobalPluginManager returned nil")
	}
}

// TestPluginInfo tests the PluginInfo structure
func TestPluginInfo(t *testing.T) {
	info := &PluginInfo{
		Path:   "/path/to/plugin.so",
		Loaded: true,
	}

	if info.Path != "/path/to/plugin.so" {
		t.Errorf("Expected path '/path/to/plugin.so', got '%s'", info.Path)
	}

	if !info.Loaded {
		t.Error("Expected Loaded to be true")
	}
}

// TestLoadPluginNonexistent tests loading a nonexistent plugin
func TestLoadPluginNonexistent(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	err := pm.LoadPlugin("/nonexistent/plugin.so")
	if err == nil {
		t.Error("Expected error loading nonexistent plugin")
	}

	if pm.Count() != 0 {
		t.Errorf("Expected 0 plugins after failed load, got %d", pm.Count())
	}
}

// TestLoadPluginInvalidPath tests loading with invalid path
func TestLoadPluginInvalidPath(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Empty path
	err := pm.LoadPlugin("")
	if err == nil {
		t.Error("Expected error with empty path")
	}
}

// TestIsLoaded tests the IsLoaded method
func TestIsLoaded(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	if pm.IsLoaded("/nonexistent.so") {
		t.Error("Expected false for nonexistent plugin")
	}

	// Create a failed load entry
	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:   "/test.so",
		Loaded: false,
	}
	pm.mu.Unlock()

	if pm.IsLoaded("/test.so") {
		t.Error("Expected false for failed plugin load")
	}
}

// TestGetPlugin tests retrieving plugin info
func TestGetPlugin(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Non-existent plugin
	_, ok := pm.GetPlugin("/nonexistent.so")
	if ok {
		t.Error("Expected false for nonexistent plugin")
	}

	// Add a plugin manually
	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:   "/test.so",
		Loaded: true,
	}
	pm.mu.Unlock()

	info, ok := pm.GetPlugin("/test.so")
	if !ok {
		t.Error("Expected true for existing plugin")
	}
	if info.Path != "/test.so" {
		t.Errorf("Expected path '/test.so', got '%s'", info.Path)
	}
}

// TestGetPluginByExtension tests retrieving plugin by extension name
func TestGetPluginByExtension(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Create a mock extension
	ext := NewBaseExtension("testextension", "1.0.0")

	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:      "/test.so",
		Extension: ext,
		Loaded:    true,
	}
	pm.mu.Unlock()

	// Find by extension name
	info, ok := pm.GetPluginByExtension("testextension")
	if !ok {
		t.Error("Expected to find plugin by extension name")
	}
	if info.Path != "/test.so" {
		t.Errorf("Expected path '/test.so', got '%s'", info.Path)
	}

	// Non-existent extension
	_, ok = pm.GetPluginByExtension("nonexistent")
	if ok {
		t.Error("Expected false for nonexistent extension")
	}
}

// TestList tests listing loaded plugins
func TestList(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Add some plugins
	pm.mu.Lock()
	pm.plugins["/test1.so"] = &PluginInfo{
		Path:   "/test1.so",
		Loaded: true,
	}
	pm.plugins["/test2.so"] = &PluginInfo{
		Path:   "/test2.so",
		Loaded: true,
	}
	pm.plugins["/test3.so"] = &PluginInfo{
		Path:   "/test3.so",
		Loaded: false, // Not loaded
	}
	pm.mu.Unlock()

	list := pm.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 loaded plugins, got %d", len(list))
	}

	// Check that only loaded plugins are returned
	found := make(map[string]bool)
	for _, path := range list {
		found[path] = true
	}

	if !found["/test1.so"] || !found["/test2.so"] {
		t.Error("Expected to find both loaded plugins")
	}
	if found["/test3.so"] {
		t.Error("Did not expect to find unloaded plugin")
	}
}

// TestListAll tests listing all plugins including failed loads
func TestListAll(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	pm.mu.Lock()
	pm.plugins["/test1.so"] = &PluginInfo{Path: "/test1.so", Loaded: true}
	pm.plugins["/test2.so"] = &PluginInfo{Path: "/test2.so", Loaded: false}
	pm.mu.Unlock()

	list := pm.ListAll()
	if len(list) != 2 {
		t.Errorf("Expected 2 plugins total, got %d", len(list))
	}
}

// TestCount tests counting loaded plugins
func TestCount(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	if pm.Count() != 0 {
		t.Errorf("Expected 0 plugins, got %d", pm.Count())
	}

	pm.mu.Lock()
	pm.plugins["/test1.so"] = &PluginInfo{Path: "/test1.so", Loaded: true}
	pm.plugins["/test2.so"] = &PluginInfo{Path: "/test2.so", Loaded: true}
	pm.plugins["/test3.so"] = &PluginInfo{Path: "/test3.so", Loaded: false}
	pm.mu.Unlock()

	if pm.Count() != 2 {
		t.Errorf("Expected 2 loaded plugins, got %d", pm.Count())
	}
}

// TestUnloadPlugin tests unloading a plugin
func TestUnloadPlugin(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	pm := NewPluginManager(em)

	ext := NewBaseExtension("testextension", "1.0.0")
	em.Register(ext)

	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:      "/test.so",
		Extension: ext,
		Loaded:    true,
	}
	pm.mu.Unlock()

	if pm.Count() != 1 {
		t.Errorf("Expected 1 plugin before unload, got %d", pm.Count())
	}

	err := pm.UnloadPlugin("/test.so")
	if err != nil {
		t.Errorf("Unexpected error unloading plugin: %v", err)
	}

	if pm.Count() != 0 {
		t.Errorf("Expected 0 plugins after unload, got %d", pm.Count())
	}

	// Verify extension was unregistered
	if em.Has("testextension") {
		t.Error("Extension should have been unregistered")
	}
}

// TestUnloadPluginNonexistent tests unloading a nonexistent plugin
func TestUnloadPluginNonexistent(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	err := pm.UnloadPlugin("/nonexistent.so")
	if err == nil {
		t.Error("Expected error unloading nonexistent plugin")
	}
}

// TestClear tests clearing all plugins
func TestClear(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	pm.mu.Lock()
	pm.plugins["/test1.so"] = &PluginInfo{Path: "/test1.so", Loaded: true}
	pm.plugins["/test2.so"] = &PluginInfo{Path: "/test2.so", Loaded: true}
	pm.mu.Unlock()

	if pm.Count() != 2 {
		t.Errorf("Expected 2 plugins before clear, got %d", pm.Count())
	}

	pm.Clear()

	if pm.Count() != 0 {
		t.Errorf("Expected 0 plugins after clear, got %d", pm.Count())
	}

	if len(pm.ListAll()) != 0 {
		t.Error("Expected empty list after clear")
	}
}

// TestGlobalFunctions tests the global plugin functions
func TestGlobalFunctions(t *testing.T) {
	// These should not panic
	plugins := ListPlugins()
	if plugins == nil {
		t.Error("ListPlugins returned nil")
	}

	_, ok := GetPluginInfo("/nonexistent.so")
	if ok {
		t.Error("Expected false for nonexistent plugin")
	}

	err := LoadPlugin("/nonexistent.so")
	if err == nil {
		t.Error("Expected error loading nonexistent plugin")
	}
}

// TestLoadPluginByName tests the LoadPluginByName convenience method
func TestLoadPluginByName(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Try to load nonexistent plugin
	_, err := pm.LoadPluginByName("/nonexistent.so")
	if err == nil {
		t.Error("Expected error loading nonexistent plugin")
	}
}

// TestPluginManagerConcurrency tests concurrent access to plugin manager
func TestPluginManagerConcurrency(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Add a test plugin
	ext := NewBaseExtension("testextension", "1.0.0")
	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:      "/test.so",
		Extension: ext,
		Loaded:    true,
	}
	pm.mu.Unlock()

	// Run concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			pm.List()
			pm.Count()
			pm.IsLoaded("/test.so")
			pm.GetPlugin("/test.so")
			pm.GetPluginByExtension("testextension")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestPluginPathNormalization tests that plugin paths are normalized
func TestPluginPathNormalization(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Test with relative path (will fail to load but path should be normalized)
	relPath := "./test.so"
	err := pm.LoadPlugin(relPath)
	if err == nil {
		t.Error("Expected error loading nonexistent plugin")
	}

	// Check that the path was normalized in the error tracking
	absRelPath, _ := filepath.Abs(relPath)
	info, ok := pm.GetPlugin(absRelPath)
	if !ok {
		t.Error("Expected plugin info to be stored with normalized path")
	}
	if info.Loaded {
		t.Error("Expected plugin to not be loaded")
	}
}

// TestPluginLoadError tests that load errors are properly tracked
func TestPluginLoadError(t *testing.T) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Try to load a nonexistent plugin
	err := pm.LoadPlugin("/nonexistent/test.so")
	if err == nil {
		t.Error("Expected error loading nonexistent plugin")
	}

	// Get the plugin info and check that error is stored
	absPath, _ := filepath.Abs("/nonexistent/test.so")
	info, ok := pm.GetPlugin(absPath)
	if !ok {
		t.Error("Expected plugin info to be stored even on error")
	}
	if info.LoadError == nil {
		t.Error("Expected LoadError to be set")
	}
	if info.Loaded {
		t.Error("Expected Loaded to be false")
	}
}

// TestPluginIntegrationWithExtensionManager tests that plugins integrate correctly with extension manager
func TestPluginIntegrationWithExtensionManager(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	pm := NewPluginManager(em)

	// Create a mock extension
	ext := NewMockExtension("mockext", "1.0.0")
	ext.functions["test_func"] = func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		return types.NewInt(42), nil
	}

	// Register it manually with the extension manager
	if err := em.Register(ext); err != nil {
		t.Fatalf("Failed to register extension: %v", err)
	}

	// Simulate plugin loading by adding to plugin manager
	pm.mu.Lock()
	pm.plugins["/mock.so"] = &PluginInfo{
		Path:      "/mock.so",
		Extension: ext,
		Loaded:    true,
	}
	pm.mu.Unlock()

	// Verify plugin is tracked
	if pm.Count() != 1 {
		t.Errorf("Expected 1 plugin, got %d", pm.Count())
	}

	// Verify we can find it by extension name
	info, ok := pm.GetPluginByExtension("mockext")
	if !ok {
		t.Error("Expected to find plugin by extension name")
	}
	if info.Extension.Name() != "mockext" {
		t.Errorf("Expected extension name 'mockext', got '%s'", info.Extension.Name())
	}
}

// BenchmarkPluginManagerList benchmarks listing plugins
func BenchmarkPluginManagerList(b *testing.B) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	// Add some test plugins
	for i := 0; i < 100; i++ {
		ext := NewBaseExtension("ext"+string(rune(i)), "1.0.0")
		pm.mu.Lock()
		pm.plugins["/test"+string(rune(i))+".so"] = &PluginInfo{
			Path:      "/test" + string(rune(i)) + ".so",
			Extension: ext,
			Loaded:    true,
		}
		pm.mu.Unlock()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.List()
	}
}

// BenchmarkPluginManagerGetPlugin benchmarks retrieving plugin info
func BenchmarkPluginManagerGetPlugin(b *testing.B) {
	pm := NewPluginManager(NewExtensionManager(NewFunctionRegistry()))

	ext := NewBaseExtension("testext", "1.0.0")
	pm.mu.Lock()
	pm.plugins["/test.so"] = &PluginInfo{
		Path:      "/test.so",
		Extension: ext,
		Loaded:    true,
	}
	pm.mu.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.GetPlugin("/test.so")
	}
}
