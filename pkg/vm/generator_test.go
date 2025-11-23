package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestGeneratorCreation tests basic generator creation
func TestGeneratorCreation(t *testing.T) {
	vm := New()

	// Create a simple function
	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	// Create generator
	gen := NewGenerator(vm, fn)

	// Check initial state
	if gen.State() != GenStateStart {
		t.Errorf("Expected GenStateStart, got %v", gen.State())
	}

	if gen.Valid() {
		t.Error("Generator should not be valid initially")
	}

	if gen.Current().Type() != types.TypeNull {
		t.Error("Current() should return null for non-yielded generator")
	}
}

// TestGeneratorStateTransitions tests state transitions
func TestGeneratorStateTransitions(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Test state: Start
	if gen.State() != GenStateStart {
		t.Fatal("Expected GenStateStart")
	}

	// Yield a value
	gen.Yield(nil, types.NewInt(42), 0)

	// Test state: Yielded
	if gen.State() != GenStateYielded {
		t.Errorf("Expected GenStateYielded, got %v", gen.State())
	}

	if !gen.Valid() {
		t.Error("Generator should be valid after yield")
	}

	// Check yielded value
	current := gen.Current()
	if current.ToInt() != 42 {
		t.Errorf("Expected value 42, got %d", current.ToInt())
	}

	// Check auto-incremented key
	key := gen.Key()
	if key.ToInt() != 0 {
		t.Errorf("Expected key 0, got %d", key.ToInt())
	}
}

// TestGeneratorAutoIncrementKey tests auto-incrementing keys
func TestGeneratorAutoIncrementKey(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Yield multiple values without keys
	gen.Yield(nil, types.NewInt(10), 0)
	if gen.Key().ToInt() != 0 {
		t.Errorf("Expected key 0, got %d", gen.Key().ToInt())
	}

	gen.Yield(nil, types.NewInt(20), 0)
	if gen.Key().ToInt() != 1 {
		t.Errorf("Expected key 1, got %d", gen.Key().ToInt())
	}

	gen.Yield(nil, types.NewInt(30), 0)
	if gen.Key().ToInt() != 2 {
		t.Errorf("Expected key 2, got %d", gen.Key().ToInt())
	}
}

// TestGeneratorWithExplicitKey tests yielding with explicit keys
func TestGeneratorWithExplicitKey(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Yield with string key
	gen.Yield(types.NewString("foo"), types.NewInt(100), 0)

	if gen.Key().ToString() != "foo" {
		t.Errorf("Expected key 'foo', got '%s'", gen.Key().ToString())
	}

	if gen.Current().ToInt() != 100 {
		t.Errorf("Expected value 100, got %d", gen.Current().ToInt())
	}
}

// TestGeneratorReturn tests generator return value
func TestGeneratorReturn(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Return from generator
	gen.Return(types.NewString("final"))

	// Check state is done
	if gen.State() != GenStateDone {
		t.Errorf("Expected GenStateDone, got %v", gen.State())
	}

	// Get return value
	retVal, err := gen.GetReturn()
	if err != nil {
		t.Fatalf("GetReturn failed: %v", err)
	}

	if retVal.ToString() != "final" {
		t.Errorf("Expected return value 'final', got '%s'", retVal.ToString())
	}
}

// TestGeneratorGetReturnBeforeDone tests error when getting return before done
func TestGeneratorGetReturnBeforeDone(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Try to get return value before generator is done
	_, err := gen.GetReturn()
	if err == nil {
		t.Error("Expected error when getting return value before generator is done")
	}
}

// TestGeneratorClose tests closing a generator
func TestGeneratorClose(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Yield a value
	gen.Yield(nil, types.NewInt(42), 0)

	// Close generator
	gen.Close()

	// Check state is closed
	if gen.State() != GenStateClosed {
		t.Errorf("Expected GenStateClosed, got %v", gen.State())
	}

	// Try to send to closed generator
	_, err := gen.Send(types.NewInt(10))
	if err == nil {
		t.Error("Expected error when sending to closed generator")
	}
}

// TestGeneratorSendFirstIteration tests error when sending on first iteration
func TestGeneratorSendFirstIteration(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Try to send non-null value on first iteration
	_, err := gen.Send(types.NewInt(42))
	if err == nil {
		t.Error("Expected error when sending value on first iteration")
	}
}

// TestGeneratorRewindError tests error when rewinding started generator
func TestGeneratorRewindError(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Start generator by yielding
	gen.SetState(GenStateYielded)

	// Try to rewind after start
	err := gen.Rewind()
	if err == nil {
		t.Error("Expected error when rewinding already-started generator")
	}
}

// TestGeneratorFrameCreation tests that Resume creates a frame
func TestGeneratorFrameCreation(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "testGenerator",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	gen := NewGenerator(vm, fn)

	// Frame should be nil initially
	if gen.GetFrame() != nil {
		t.Error("Frame should be nil initially")
	}

	// Resume should create frame
	err := gen.Resume(types.NewNull())
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}

	// Frame should now exist
	if gen.GetFrame() == nil {
		t.Error("Frame should exist after Resume")
	}

	// State should be running
	if gen.State() != GenStateRunning {
		t.Errorf("Expected GenStateRunning, got %v", gen.State())
	}
}
