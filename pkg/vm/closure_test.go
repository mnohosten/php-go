package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestClosureCreation tests basic closure creation
func TestClosureCreation(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    2,
	}

	closure := NewClosure(fn)

	if closure.Function != fn {
		t.Error("Closure function not set correctly")
	}

	if closure.Static {
		t.Error("Closure should not be static by default")
	}

	if closure.CapturedVars == nil {
		t.Error("CapturedVars map should be initialized")
	}
}

// TestStaticClosureCreation tests static closure creation
func TestStaticClosureCreation(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewStaticClosure(fn)

	if !closure.Static {
		t.Error("Closure should be static")
	}

	if !closure.IsStatic() {
		t.Error("IsStatic() should return true")
	}
}

// TestBindVariableByValue tests binding a variable by value
func TestBindVariableByValue(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	// Bind a variable by value
	value := types.NewInt(42)
	closure.BindVariable("x", value, false)

	// Check if variable is captured
	if !closure.HasCapturedVariable("x") {
		t.Error("Variable 'x' should be captured")
	}

	// Get captured variable
	capturedVal, ok := closure.GetCapturedVariable("x")
	if !ok {
		t.Fatal("Failed to get captured variable 'x'")
	}

	if capturedVal.ToInt() != 42 {
		t.Errorf("Expected captured value 42, got %d", capturedVal.ToInt())
	}
}

// TestBindVariableByReference tests binding a variable by reference
func TestBindVariableByReference(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	// Create a value
	value := types.NewInt(100)

	// Bind by reference
	closure.BindVariable("y", value, true)

	// Check if variable is captured
	if !closure.HasCapturedVariable("y") {
		t.Error("Variable 'y' should be captured")
	}

	// The captured variable should be a reference
	capturedRef := closure.CapturedVars["y"]
	if capturedRef.Type() != types.TypeReference {
		t.Errorf("Expected reference type, got %s", capturedRef.TypeString())
	}

	// Dereferencing should give us the original value
	derefVal := capturedRef.Deref()
	if derefVal.ToInt() != 100 {
		t.Errorf("Expected dereferenced value 100, got %d", derefVal.ToInt())
	}
}

// TestBindMultipleVariables tests binding multiple variables
func TestBindMultipleVariables(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	// Bind multiple variables
	closure.BindVariable("a", types.NewInt(1), false)
	closure.BindVariable("b", types.NewString("hello"), false)
	closure.BindVariable("c", types.NewBool(true), false)

	// Check all are captured
	if !closure.HasCapturedVariable("a") {
		t.Error("Variable 'a' should be captured")
	}
	if !closure.HasCapturedVariable("b") {
		t.Error("Variable 'b' should be captured")
	}
	if !closure.HasCapturedVariable("c") {
		t.Error("Variable 'c' should be captured")
	}

	// Verify values
	if val, ok := closure.GetCapturedVariable("a"); !ok || val.ToInt() != 1 {
		t.Error("Variable 'a' has wrong value")
	}
	if val, ok := closure.GetCapturedVariable("b"); !ok || val.ToString() != "hello" {
		t.Error("Variable 'b' has wrong value")
	}
	if val, ok := closure.GetCapturedVariable("c"); !ok || !val.ToBool() {
		t.Error("Variable 'c' has wrong value")
	}
}

// TestClosureBind tests Closure::bindTo()
func TestClosureBind(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	// Bind a variable
	closure.BindVariable("x", types.NewInt(42), false)

	// Create new $this
	newThis := &types.Object{
		ClassName: "TestClass",
	}

	// Bind to new $this
	newClosure, err := closure.Bind(newThis, nil)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// New closure should have same function
	if newClosure.Function != closure.Function {
		t.Error("Bound closure should have same function")
	}

	// New closure should have copied captured variables
	if !newClosure.HasCapturedVariable("x") {
		t.Error("Bound closure should have captured variable 'x'")
	}

	val, _ := newClosure.GetCapturedVariable("x")
	if val.ToInt() != 42 {
		t.Error("Bound closure should have same captured value")
	}
}

// TestClosureBindStaticError tests that static closures cannot be rebound
func TestClosureBindStaticError(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewStaticClosure(fn)

	newThis := &types.Object{
		ClassName: "TestClass",
	}

	// Try to bind static closure
	_, err := closure.Bind(newThis, nil)
	if err == nil {
		t.Error("Expected error when binding static closure")
	}
}

// TestClosureBindStatic tests creating a static version
func TestClosureBindStatic(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)
	closure.BindVariable("x", types.NewInt(42), false)

	// Create static version
	staticClosure, err := closure.BindStatic(nil)
	if err != nil {
		t.Fatalf("BindStatic failed: %v", err)
	}

	if !staticClosure.IsStatic() {
		t.Error("New closure should be static")
	}

	// Should have same function
	if staticClosure.Function != closure.Function {
		t.Error("Static closure should have same function")
	}

	// Should have copied captured variables
	if !staticClosure.HasCapturedVariable("x") {
		t.Error("Static closure should have captured variable 'x'")
	}
}

// TestClosureClone tests cloning a closure
func TestClosureClone(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)
	closure.BindVariable("x", types.NewInt(42), false)
	closure.BindVariable("y", types.NewString("test"), false)

	// Clone the closure
	cloned := closure.Clone()

	// Cloned closure should have same function
	if cloned.Function != closure.Function {
		t.Error("Cloned closure should have same function")
	}

	// Cloned closure should have same static flag
	if cloned.Static != closure.Static {
		t.Error("Cloned closure should have same static flag")
	}

	// Cloned closure should have same captured variables
	if !cloned.HasCapturedVariable("x") {
		t.Error("Cloned closure should have captured variable 'x'")
	}
	if !cloned.HasCapturedVariable("y") {
		t.Error("Cloned closure should have captured variable 'y'")
	}

	// Verify values
	if val, ok := cloned.GetCapturedVariable("x"); !ok || val.ToInt() != 42 {
		t.Error("Cloned closure has wrong value for 'x'")
	}
	if val, ok := cloned.GetCapturedVariable("y"); !ok || val.ToString() != "test" {
		t.Error("Cloned closure has wrong value for 'y'")
	}
}

// TestGetCapturedVariables tests getting all captured variables
func TestGetCapturedVariables(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)
	closure.BindVariable("a", types.NewInt(1), false)
	closure.BindVariable("b", types.NewInt(2), false)
	closure.BindVariable("c", types.NewInt(3), false)

	// Get all captured variables
	vars := closure.GetCapturedVariables()

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

// TestClosureGetFunction tests GetFunction method
func TestClosureGetFunction(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	retrievedFn := closure.GetFunction()
	if retrievedFn != fn {
		t.Error("GetFunction should return the original function")
	}
}

// TestGetClosureFromValue tests extracting closure from a Value
func TestGetClosureFromValue(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testClosure",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	closure := NewClosure(fn)

	// Wrap in resource
	resource := types.NewResourceHandle("Closure", closure)
	closureVal := types.NewResource(resource)

	// Extract closure
	extractedClosure, err := GetClosureFromValue(closureVal)
	if err != nil {
		t.Fatalf("GetClosureFromValue failed: %v", err)
	}

	if extractedClosure != closure {
		t.Error("Extracted closure should be the same as original")
	}
}

// TestGetClosureFromValueError tests error when not a closure
func TestGetClosureFromValueError(t *testing.T) {
	// Try with non-resource value
	intVal := types.NewInt(42)
	_, err := GetClosureFromValue(intVal)
	if err == nil {
		t.Error("Expected error when extracting closure from non-resource")
	}
}
