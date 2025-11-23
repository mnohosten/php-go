package goext

import (
	"errors"
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Mock extension for testing

type MockExtension struct {
	name         string
	version      string
	initCalled   bool
	initError    error
	functions    map[string]FunctionHandler
	constants    map[string]*types.Value
}

func NewMockExtension(name, version string) *MockExtension {
	return &MockExtension{
		name:      name,
		version:   version,
		functions: make(map[string]FunctionHandler),
		constants: make(map[string]*types.Value),
	}
}

func (e *MockExtension) Name() string {
	return e.name
}

func (e *MockExtension) Version() string {
	return e.version
}

func (e *MockExtension) Init(v *vm.VM) error {
	e.initCalled = true
	return e.initError
}

func (e *MockExtension) Functions() map[string]FunctionHandler {
	return e.functions
}

func (e *MockExtension) Constants() map[string]*types.Value {
	return e.constants
}

// Test ExtensionManager creation

func TestNewExtensionManager(t *testing.T) {
	reg := NewFunctionRegistry()
	em := NewExtensionManager(reg)

	if em == nil {
		t.Fatal("NewExtensionManager returned nil")
	}

	if em.registry != reg {
		t.Error("ExtensionManager has wrong registry")
	}

	if em.Count() != 0 {
		t.Errorf("Expected 0 extensions, got %d", em.Count())
	}
}

func TestGetGlobalExtensionManager(t *testing.T) {
	em1 := GetGlobalExtensionManager()
	em2 := GetGlobalExtensionManager()

	if em1 != em2 {
		t.Error("GetGlobalExtensionManager returned different instances")
	}
}

// Test extension registration

func TestRegisterExtension(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	ext := NewMockExtension("test_ext", "1.0.0")
	err := em.Register(ext)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !em.Has("test_ext") {
		t.Error("Extension not found after registration")
	}

	info, ok := em.Get("test_ext")
	if !ok {
		t.Fatal("Get returned false for registered extension")
	}

	if info.Extension.Name() != "test_ext" {
		t.Errorf("Expected name 'test_ext', got '%s'", info.Extension.Name())
	}

	if info.Loaded {
		t.Error("Extension should not be loaded yet")
	}
}

func TestRegisterExtensionDuplicate(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	ext := NewMockExtension("test_ext", "1.0.0")
	err := em.Register(ext)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	err = em.Register(ext)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}

	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("Expected 'already registered' error, got: %v", err)
	}
}

func TestRegisterExtensionNil(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	err := em.Register(nil)
	if err == nil {
		t.Error("Expected error for nil extension")
	}
}

func TestRegisterExtensionEmptyName(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	ext := NewMockExtension("", "1.0.0")
	err := em.Register(ext)
	if err == nil {
		t.Error("Expected error for empty name")
	}
}

// Test extension loading

func TestLoadIntoVM(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	ext := NewMockExtension("test_ext", "1.0.0")

	// Add a function
	ext.functions["test_func"] = func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		return types.NewString("test_result"), nil
	}

	// Add a constant
	ext.constants["TEST_CONST"] = types.NewInt(42)

	err := em.Register(ext)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = em.LoadIntoVM("test_ext", testVM)
	if err != nil {
		t.Fatalf("LoadIntoVM failed: %v", err)
	}

	// Check that Init was called
	if !ext.initCalled {
		t.Error("Init was not called")
	}

	// Check that extension is marked as loaded
	info, _ := em.Get("test_ext")
	if !info.Loaded {
		t.Error("Extension not marked as loaded")
	}

	// Check that constant was registered
	constVal, ok := testVM.GetGlobal("TEST_CONST")
	if !ok {
		t.Error("Constant not registered in VM")
	}
	if constVal.ToInt() != 42 {
		t.Errorf("Expected constant value 42, got %d", constVal.ToInt())
	}
}

func TestLoadIntoVMNotRegistered(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	err := em.LoadIntoVM("nonexistent", testVM)
	if err == nil {
		t.Error("Expected error for non-existent extension")
	}

	if !strings.Contains(err.Error(), "not registered") {
		t.Errorf("Expected 'not registered' error, got: %v", err)
	}
}

func TestLoadIntoVMInitError(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	ext := NewMockExtension("test_ext", "1.0.0")
	ext.initError = errors.New("init failed")

	err := em.Register(ext)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = em.LoadIntoVM("test_ext", testVM)
	if err == nil {
		t.Error("Expected error from Init")
	}

	if !strings.Contains(err.Error(), "init failed") {
		t.Errorf("Expected 'init failed' error, got: %v", err)
	}

	// Check that InitError is stored
	info, _ := em.Get("test_ext")
	if info.InitError == nil {
		t.Error("InitError not stored")
	}
}

func TestLoadAllIntoVM(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	// Register multiple extensions
	ext1 := NewMockExtension("ext1", "1.0.0")
	ext2 := NewMockExtension("ext2", "2.0.0")
	ext3 := NewMockExtension("ext3", "3.0.0")

	em.Register(ext1)
	em.Register(ext2)
	em.Register(ext3)

	err := em.LoadAllIntoVM(testVM)
	if err != nil {
		t.Fatalf("LoadAllIntoVM failed: %v", err)
	}

	// Check all are loaded
	for _, name := range []string{"ext1", "ext2", "ext3"} {
		info, ok := em.Get(name)
		if !ok {
			t.Errorf("Extension '%s' not found", name)
			continue
		}
		if !info.Loaded {
			t.Errorf("Extension '%s' not loaded", name)
		}
	}
}

func TestLoadAllIntoVMWithError(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	ext1 := NewMockExtension("ext1", "1.0.0")
	ext2 := NewMockExtension("ext2", "2.0.0")
	ext2.initError = errors.New("ext2 init failed")

	em.Register(ext1)
	em.Register(ext2)

	err := em.LoadAllIntoVM(testVM)
	if err == nil {
		t.Error("Expected error from LoadAllIntoVM")
	}
}

// Test extension manager operations

func TestExtensionManagerList(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	em.Register(NewMockExtension("ext1", "1.0.0"))
	em.Register(NewMockExtension("ext2", "2.0.0"))
	em.Register(NewMockExtension("ext3", "3.0.0"))

	names := em.List()
	if len(names) != 3 {
		t.Errorf("Expected 3 extensions, got %d", len(names))
	}

	// Check all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"ext1", "ext2", "ext3"} {
		if !nameMap[expected] {
			t.Errorf("Expected '%s' in list", expected)
		}
	}
}

func TestExtensionManagerCount(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	if em.Count() != 0 {
		t.Errorf("Expected 0 extensions, got %d", em.Count())
	}

	em.Register(NewMockExtension("ext1", "1.0.0"))
	if em.Count() != 1 {
		t.Errorf("Expected 1 extension, got %d", em.Count())
	}

	em.Register(NewMockExtension("ext2", "2.0.0"))
	em.Register(NewMockExtension("ext3", "3.0.0"))
	if em.Count() != 3 {
		t.Errorf("Expected 3 extensions, got %d", em.Count())
	}
}

func TestExtensionManagerUnregister(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	em.Register(NewMockExtension("ext1", "1.0.0"))
	em.Register(NewMockExtension("ext2", "2.0.0"))

	if em.Count() != 2 {
		t.Fatalf("Expected 2 extensions before unregister")
	}

	err := em.Unregister("ext1")
	if err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	if em.Has("ext1") {
		t.Error("ext1 still exists after unregister")
	}

	if em.Count() != 1 {
		t.Errorf("Expected 1 extension after unregister, got %d", em.Count())
	}
}

func TestExtensionManagerUnregisterNonExistent(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	err := em.Unregister("nonexistent")
	if err == nil {
		t.Error("Expected error for unregistering non-existent extension")
	}
}

func TestExtensionManagerClear(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())

	em.Register(NewMockExtension("ext1", "1.0.0"))
	em.Register(NewMockExtension("ext2", "2.0.0"))

	if em.Count() != 2 {
		t.Fatalf("Expected 2 extensions before clear")
	}

	em.Clear()

	if em.Count() != 0 {
		t.Errorf("Expected 0 extensions after clear, got %d", em.Count())
	}

	if em.Has("ext1") {
		t.Error("ext1 still exists after clear")
	}
}

// Test BaseExtension

func TestBaseExtension(t *testing.T) {
	ext := NewBaseExtension("base_ext", "1.0.0")

	if ext.Name() != "base_ext" {
		t.Errorf("Expected name 'base_ext', got '%s'", ext.Name())
	}

	if ext.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", ext.Version())
	}

	// Test Init (should be no-op)
	testVM := vm.New()
	err := ext.Init(testVM)
	if err != nil {
		t.Errorf("BaseExtension Init failed: %v", err)
	}

	// Test empty functions and constants
	if len(ext.Functions()) != 0 {
		t.Errorf("Expected 0 functions, got %d", len(ext.Functions()))
	}

	if len(ext.Constants()) != 0 {
		t.Errorf("Expected 0 constants, got %d", len(ext.Constants()))
	}
}

func TestBaseExtensionAddFunction(t *testing.T) {
	ext := NewBaseExtension("base_ext", "1.0.0")

	handler := func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		return types.NewString("test"), nil
	}

	ext.AddFunction("test_func", handler)

	funcs := ext.Functions()
	if len(funcs) != 1 {
		t.Errorf("Expected 1 function, got %d", len(funcs))
	}

	if _, ok := funcs["test_func"]; !ok {
		t.Error("test_func not found in functions")
	}
}

func TestBaseExtensionAddConstant(t *testing.T) {
	ext := NewBaseExtension("base_ext", "1.0.0")

	ext.AddConstant("MY_CONST", types.NewInt(123))

	consts := ext.Constants()
	if len(consts) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(consts))
	}

	val, ok := consts["MY_CONST"]
	if !ok {
		t.Fatal("MY_CONST not found in constants")
	}

	if val.ToInt() != 123 {
		t.Errorf("Expected constant value 123, got %d", val.ToInt())
	}
}

// Test global registration

func TestRegisterExtensionGlobal(t *testing.T) {
	// Clear global manager for clean test
	GetGlobalExtensionManager().Clear()

	ext := NewMockExtension("global_ext", "1.0.0")
	err := RegisterExtension(ext)
	if err != nil {
		t.Fatalf("RegisterExtension failed: %v", err)
	}

	em := GetGlobalExtensionManager()
	if !em.Has("global_ext") {
		t.Error("Extension not registered in global manager")
	}
}

// Test integration scenario

func TestExtensionIntegrationComplete(t *testing.T) {
	em := NewExtensionManager(NewFunctionRegistry())
	testVM := vm.New()

	// Create a complete extension
	ext := NewBaseExtension("complete_ext", "1.0.0")

	// Add functions
	ext.AddFunction("greet", func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		if len(args) != 1 {
			return nil, errors.New("greet() expects 1 argument")
		}
		name := args[0].ToString()
		return types.NewString("Hello, " + name + "!"), nil
	})

	ext.AddFunction("add", func(v *vm.VM, args []*types.Value) (*types.Value, error) {
		if len(args) != 2 {
			return nil, errors.New("add() expects 2 arguments")
		}
		result := args[0].ToInt() + args[1].ToInt()
		return types.NewInt(result), nil
	})

	// Add constants
	ext.AddConstant("EXT_VERSION", types.NewString("1.0.0"))
	ext.AddConstant("MAX_VALUE", types.NewInt(1000))

	// Register and load
	err := em.Register(ext)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = em.LoadIntoVM("complete_ext", testVM)
	if err != nil {
		t.Fatalf("LoadIntoVM failed: %v", err)
	}

	// Verify constants are available
	version, ok := testVM.GetGlobal("EXT_VERSION")
	if !ok {
		t.Error("EXT_VERSION constant not found")
	} else if version.ToString() != "1.0.0" {
		t.Errorf("Expected EXT_VERSION '1.0.0', got '%s'", version.ToString())
	}

	maxVal, ok := testVM.GetGlobal("MAX_VALUE")
	if !ok {
		t.Error("MAX_VALUE constant not found")
	} else if maxVal.ToInt() != 1000 {
		t.Errorf("Expected MAX_VALUE 1000, got %d", maxVal.ToInt())
	}

	// Verify extension is loaded
	info, ok := em.Get("complete_ext")
	if !ok {
		t.Fatal("Extension info not found")
	}

	if !info.Loaded {
		t.Error("Extension not marked as loaded")
	}

	if info.InitError != nil {
		t.Errorf("Unexpected init error: %v", info.InitError)
	}
}

// Test helper functions

func TestSetGlobalGetGlobal(t *testing.T) {
	testVM := vm.New()

	SetGlobal(testVM, "test_var", types.NewString("test_value"))

	val, ok := GetGlobal(testVM, "test_var")
	if !ok {
		t.Fatal("GetGlobal returned false")
	}

	if val.ToString() != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", val.ToString())
	}
}

func TestGetGlobalNonExistent(t *testing.T) {
	testVM := vm.New()

	_, ok := GetGlobal(testVM, "nonexistent")
	if ok {
		t.Error("GetGlobal returned true for non-existent global")
	}
}
