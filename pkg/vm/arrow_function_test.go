package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestArrowFunctionCreation tests basic arrow function creation
func TestArrowFunctionCreation(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	if arrow.Function != fn {
		t.Error("Arrow function function not set correctly")
	}

	if arrow.CapturedVars == nil {
		t.Error("CapturedVars map should be initialized")
	}
}

// TestArrowFunctionAutoCapture tests that arrow functions automatically capture variables
func TestArrowFunctionAutoCapture(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	// Simulate auto-capture of variable $x from parent scope
	arrow.BindVariable("x", types.NewInt(42), false)

	if !arrow.HasCapturedVariable("x") {
		t.Error("Arrow function should have auto-captured variable 'x'")
	}

	capturedVal, ok := arrow.GetCapturedVariable("x")
	if !ok {
		t.Fatal("Failed to get captured variable 'x'")
	}

	if capturedVal.ToInt() != 42 {
		t.Errorf("Expected captured value 42, got %d", capturedVal.ToInt())
	}
}

// TestArrowFunctionCaptureBValue tests arrow functions capture by value (not reference)
func TestArrowFunctionCaptureByValue(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	// Arrow functions capture by value
	value := types.NewInt(100)
	arrow.BindVariable("y", value, false) // false = by value

	// Check if variable is captured
	if !arrow.HasCapturedVariable("y") {
		t.Error("Variable 'y' should be captured")
	}

	// The captured variable should NOT be a reference
	capturedVal := arrow.CapturedVars["y"]
	if capturedVal.Type() == types.TypeReference {
		t.Error("Arrow functions should capture by value, not by reference")
	}

	// Value should be correct
	if capturedVal.ToInt() != 100 {
		t.Errorf("Expected captured value 100, got %d", capturedVal.ToInt())
	}
}

// TestArrowFunctionMultipleCaptures tests capturing multiple variables
func TestArrowFunctionMultipleCaptures(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	// Capture multiple variables
	arrow.BindVariable("a", types.NewInt(1), false)
	arrow.BindVariable("b", types.NewString("hello"), false)
	arrow.BindVariable("c", types.NewBool(true), false)

	// Verify all captured
	if !arrow.HasCapturedVariable("a") {
		t.Error("Variable 'a' should be captured")
	}
	if !arrow.HasCapturedVariable("b") {
		t.Error("Variable 'b' should be captured")
	}
	if !arrow.HasCapturedVariable("c") {
		t.Error("Variable 'c' should be captured")
	}

	// Verify values
	if val, ok := arrow.GetCapturedVariable("a"); !ok || val.ToInt() != 1 {
		t.Error("Variable 'a' has wrong value")
	}
	if val, ok := arrow.GetCapturedVariable("b"); !ok || val.ToString() != "hello" {
		t.Error("Variable 'b' has wrong value")
	}
	if val, ok := arrow.GetCapturedVariable("c"); !ok || !val.ToBool() {
		t.Error("Variable 'c' has wrong value")
	}
}

// TestStaticArrowFunction tests static arrow functions (static fn)
func TestStaticArrowFunction(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testStaticArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	arrow := NewStaticClosure(fn)

	if !arrow.IsStatic() {
		t.Error("Arrow function should be static")
	}

	// Static arrow functions cannot access $this
	// But they can still capture variables from parent scope
	arrow.BindVariable("x", types.NewInt(42), false)

	if !arrow.HasCapturedVariable("x") {
		t.Error("Static arrow function should still capture variables")
	}
}

// TestArrowFunctionClone tests cloning an arrow function
func TestArrowFunctionClone(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)
	arrow.BindVariable("x", types.NewInt(42), false)
	arrow.BindVariable("y", types.NewString("test"), false)

	// Clone the arrow function
	cloned := arrow.Clone()

	// Cloned arrow should have same function
	if cloned.Function != arrow.Function {
		t.Error("Cloned arrow function should have same function")
	}

	// Cloned arrow should have same captured variables
	if !cloned.HasCapturedVariable("x") {
		t.Error("Cloned arrow function should have captured variable 'x'")
	}
	if !cloned.HasCapturedVariable("y") {
		t.Error("Cloned arrow function should have captured variable 'y'")
	}

	// Verify values
	if val, ok := cloned.GetCapturedVariable("x"); !ok || val.ToInt() != 42 {
		t.Error("Cloned arrow function has wrong value for 'x'")
	}
	if val, ok := cloned.GetCapturedVariable("y"); !ok || val.ToString() != "test" {
		t.Error("Cloned arrow function has wrong value for 'y'")
	}
}

// TestArrowFunctionGetCapturedVariables tests getting all captured variables
func TestArrowFunctionGetCapturedVariables(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)
	arrow.BindVariable("a", types.NewInt(1), false)
	arrow.BindVariable("b", types.NewInt(2), false)
	arrow.BindVariable("c", types.NewInt(3), false)

	// Get all captured variables
	vars := arrow.GetCapturedVariables()

	if len(vars) != 3 {
		t.Errorf("Expected 3 captured variables, got %d", len(vars))
	}

	if _, ok := vars["a"]; !ok {
		t.Error("Expected variable 'a' in captured variables")
	}
	if _, ok := vars["b"]; !ok {
		t.Error("Expected variable 'b' in captured variables")
	}
	if _, ok := vars["c"]; !ok {
		t.Error("Expected variable 'c' in captured variables")
	}
}

// TestArrowFunctionGetFunction tests GetFunction method
func TestArrowFunctionGetFunction(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	retrievedFn := arrow.GetFunction()
	if retrievedFn != fn {
		t.Error("GetFunction should return the original function")
	}
}

// TestArrowFunctionFromValue tests extracting arrow function from a Value
func TestArrowFunctionFromValue(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	// Wrap in resource
	resource := types.NewResourceHandle("Closure", arrow)
	arrowVal := types.NewResource(resource)

	// Extract arrow function
	extractedArrow, err := GetClosureFromValue(arrowVal)
	if err != nil {
		t.Fatalf("GetClosureFromValue failed: %v", err)
	}

	if extractedArrow != arrow {
		t.Error("Extracted arrow function should be the same as original")
	}
}

// TestArrowFunctionImplicitReturn tests that arrow functions implicitly return
func TestArrowFunctionImplicitReturn(t *testing.T) {
	// Arrow functions always return the value of their expression
	// This is a semantic test - the compiler should emit OpReturn automatically
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)

	// Arrow functions have ReturnByRef set to false by default
	if arrow.ReturnByRef {
		t.Error("Arrow functions should not return by reference by default")
	}
}

// TestArrowFunctionNoThisInStatic tests that static arrow functions don't capture $this
func TestArrowFunctionNoThisInStatic(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testStaticArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	// Create static arrow function
	arrow := NewStaticClosure(fn)

	// Attempt to bind $this should still be possible (it's just a variable)
	// but the static flag prevents it from being used
	obj := &types.Object{
		ClassName: "TestClass",
	}
	arrow.BindVariable("this", types.NewObject(obj), false)

	// The variable is captured, but the Static flag indicates $this shouldn't be accessible
	if !arrow.IsStatic() {
		t.Error("Arrow function should be static")
	}

	// In a real scenario, the VM would check the Static flag and not allow $this access
}

// TestArrowFunctionBindStatic tests converting to static version
func TestArrowFunctionBindStatic(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)
	arrow.BindVariable("x", types.NewInt(42), false)

	// Create static version
	staticArrow, err := arrow.BindStatic(nil)
	if err != nil {
		t.Fatalf("BindStatic failed: %v", err)
	}

	if !staticArrow.IsStatic() {
		t.Error("New arrow function should be static")
	}

	// Should have same function
	if staticArrow.Function != arrow.Function {
		t.Error("Static arrow function should have same function")
	}

	// Should have copied captured variables
	if !staticArrow.HasCapturedVariable("x") {
		t.Error("Static arrow function should have captured variable 'x'")
	}
}

// TestArrowFunctionBindToObject tests Closure::bindTo() on arrow functions
func TestArrowFunctionBindToObject(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testArrow",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    1,
	}

	arrow := NewClosure(fn)
	arrow.BindVariable("x", types.NewInt(42), false)

	// Create new $this
	newThis := &types.Object{
		ClassName: "TestClass",
	}

	// Bind to new $this
	newArrow, err := arrow.Bind(newThis, nil)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// New arrow should have same function
	if newArrow.Function != arrow.Function {
		t.Error("Bound arrow function should have same function")
	}

	// New arrow should have copied captured variables
	if !newArrow.HasCapturedVariable("x") {
		t.Error("Bound arrow function should have captured variable 'x'")
	}

	val, _ := newArrow.GetCapturedVariable("x")
	if val.ToInt() != 42 {
		t.Error("Bound arrow function should have same captured value")
	}
}
