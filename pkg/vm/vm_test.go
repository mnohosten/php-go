package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// VM Core Tests
// ============================================================================

func TestNew(t *testing.T) {
	vm := New()

	if vm == nil {
		t.Fatal("New() returned nil")
	}

	if vm.globals == nil {
		t.Error("globals not initialized")
	}

	if vm.functions == nil {
		t.Error("functions not initialized")
	}

	if vm.classes == nil {
		t.Error("classes not initialized")
	}

	if vm.maxStackDepth != 1000 {
		t.Errorf("Expected maxStackDepth 1000, got %d", vm.maxStackDepth)
	}
}

func TestNewWithBytecode(t *testing.T) {
	instructions := Instructions{
		*NewInstruction(OpReturn, 1).WithOp1(OpConst, 0),
	}
	constants := []interface{}{int64(42)}

	vm := NewWithBytecode(instructions, constants)

	if vm == nil {
		t.Fatal("NewWithBytecode() returned nil")
	}

	if len(vm.constants) != 1 {
		t.Errorf("Expected 1 constant, got %d", len(vm.constants))
	}

	if vm.frameIndex != 0 {
		t.Errorf("Expected frameIndex 0, got %d", vm.frameIndex)
	}
}

func TestGetConstant(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{
		int64(42),
		float64(3.14),
		"hello",
		true,
		nil,
	}

	tests := []struct {
		index    int
		expected types.ValueType
		value    interface{}
	}{
		{0, types.TypeInt, int64(42)},
		{1, types.TypeFloat, 3.14},
		{2, types.TypeString, "hello"},
		{3, types.TypeBool, true},
		{4, types.TypeNull, nil},
	}

	for _, tt := range tests {
		val, err := vm.GetConstant(tt.index)
		if err != nil {
			t.Errorf("GetConstant(%d) error: %v", tt.index, err)
			continue
		}

		if val.Type() != tt.expected {
			t.Errorf("GetConstant(%d): expected type %v, got %v",
				tt.index, tt.expected, val.Type())
		}
	}
}

func TestGetConstant_OutOfRange(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(42)}

	_, err := vm.GetConstant(10)
	if err == nil {
		t.Error("Expected error for out of range constant")
	}

	_, err = vm.GetConstant(-1)
	if err == nil {
		t.Error("Expected error for negative constant index")
	}
}

// ============================================================================
// Global Variables Tests
// ============================================================================

func TestSetGetGlobal(t *testing.T) {
	vm := New()

	val := types.NewInt(42)
	vm.SetGlobal("x", val)

	retrieved, ok := vm.GetGlobal("x")
	if !ok {
		t.Error("GetGlobal() returned false for existing variable")
	}

	if retrieved.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", retrieved.ToInt())
	}
}

func TestGetGlobal_NotExists(t *testing.T) {
	vm := New()

	_, ok := vm.GetGlobal("nonexistent")
	if ok {
		t.Error("GetGlobal() returned true for non-existent variable")
	}
}

// ============================================================================
// Function Registry Tests
// ============================================================================

func TestRegisterFunction(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: Instructions{},
		NumLocals:    5,
		NumParams:    2,
	}

	vm.RegisterFunction("test", fn)

	retrieved, ok := vm.GetFunction("test")
	if !ok {
		t.Error("GetFunction() returned false for registered function")
	}

	if retrieved.Name != "test" {
		t.Errorf("Expected function name 'test', got '%s'", retrieved.Name)
	}
}

// ============================================================================
// Frame Management Tests
// ============================================================================

func TestPushPopFrame(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: Instructions{},
		NumLocals:    5,
	}

	frame := NewFrame(fn)

	// Push frame
	err := vm.pushFrame(frame)
	if err != nil {
		t.Errorf("pushFrame() error: %v", err)
	}

	if vm.frameIndex != 0 {
		t.Errorf("Expected frameIndex 0, got %d", vm.frameIndex)
	}

	// Get current frame
	current := vm.currentFrame()
	if current != frame {
		t.Error("currentFrame() returned wrong frame")
	}

	// Pop frame
	popped := vm.popFrame()
	if popped != frame {
		t.Error("popFrame() returned wrong frame")
	}

	if vm.frameIndex != -1 {
		t.Errorf("Expected frameIndex -1 after pop, got %d", vm.frameIndex)
	}
}

func TestStackOverflow(t *testing.T) {
	vm := New()
	vm.maxStackDepth = 2

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: Instructions{},
		NumLocals:    5,
	}

	// Push to limit
	vm.pushFrame(NewFrame(fn))
	vm.pushFrame(NewFrame(fn))

	// This should fail
	err := vm.pushFrame(NewFrame(fn))
	if err == nil {
		t.Error("Expected stack overflow error")
	}
}

func TestStackDepthEnforcementInPushFrame(t *testing.T) {
	// Test that stack depth is properly enforced when pushing frames
	vm := New()
	vm.maxStackDepth = 5 // Set low limit for testing

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: Instructions{},
		NumLocals:    5,
	}

	// Push frames up to the limit
	for i := 0; i < 5; i++ {
		err := vm.pushFrame(NewFrame(fn))
		if err != nil {
			t.Errorf("Unexpected error at depth %d: %v", i, err)
		}
	}

	// This should fail - exceeds maxStackDepth
	err := vm.pushFrame(NewFrame(fn))
	if err == nil {
		t.Error("Expected stack overflow error when exceeding maxStackDepth")
	}

	// Verify error message
	expectedMsg := "stack overflow"
	if err != nil && !containsSubstring(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

// Helper function to check if string contains substring
func containsSubstring(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Output Tests
// ============================================================================

func TestOutput(t *testing.T) {
	vm := New()

	vm.writeOutput([]byte("Hello"))
	vm.writeOutput([]byte(" "))
	vm.writeOutput([]byte("World"))

	output := vm.GetOutput()
	if output != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", output)
	}
}

func TestClearOutput(t *testing.T) {
	vm := New()

	vm.writeOutput([]byte("test"))
	vm.ClearOutput()

	output := vm.GetOutput()
	if output != "" {
		t.Errorf("Expected empty output after clear, got '%s'", output)
	}
}

// ============================================================================
// Execution Tests
// ============================================================================

func TestExecute_SimpleReturn(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(42)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).
			WithOp1(OpConst, 0).
			WithResult(OpCV, 0),
		*NewInstruction(OpReturn, 2).
			WithOp1(OpCV, 0),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Arithmetic(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), int64(5)}

	// Calculate: 10 + 5
	instructions := Instructions{
		// Load 10 into CV0
		*NewInstruction(OpFetchConstant, 1).
			WithOp1(OpConst, 0).
			WithResult(OpCV, 0),
		// Load 5 into CV1
		*NewInstruction(OpFetchConstant, 2).
			WithOp1(OpConst, 1).
			WithResult(OpCV, 1),
		// Add CV0 + CV1 -> CV2
		*NewInstruction(OpAdd, 3).
			WithOp1(OpCV, 0).
			WithOp2(OpCV, 1).
			WithResult(OpCV, 2),
		// Return CV2
		*NewInstruction(OpReturn, 4).
			WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	// Check that no error occurred
	// (Full integration test would check return value)
}

func TestExecute_Echo(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"Hello, World!"}

	instructions := Instructions{
		// Load string constant
		*NewInstruction(OpFetchConstant, 1).
			WithOp1(OpConst, 0).
			WithResult(OpCV, 0),
		// Echo it
		*NewInstruction(OpEcho, 2).
			WithOp1(OpCV, 0),
		// Return null
		*NewInstruction(OpReturn, 3).
			WithOp1(OpUnused, 0),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	output := vm.GetOutput()
	if output != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", output)
	}
}

func TestExecute_Comparison(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(5)}

	// Test 5 == 5
	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).
			WithOp1(OpConst, 0).
			WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).
			WithOp1(OpConst, 1).
			WithResult(OpCV, 1),
		*NewInstruction(OpIsEqual, 3).
			WithOp1(OpCV, 0).
			WithOp2(OpCV, 1).
			WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).
			WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Jump(t *testing.T) {
	vm := New()
	// Add jump target to constants pool
	vm.constants = []interface{}{int64(3)} // Jump target at index 3

	// Unconditional jump over echo
	instructions := Instructions{
		// JMP to instruction 3 (stored in constants[0])
		*NewInstruction(OpJmp, 1).
			WithOp1(OpConst, 0),
		// This should be skipped
		*NewInstruction(OpEcho, 2).
			WithOp1(OpConst, 0),
		*NewInstruction(OpEcho, 3).
			WithOp1(OpConst, 0),
		// Jump target
		*NewInstruction(OpReturn, 4).
			WithOp1(OpUnused, 0),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	// Output should be empty since we jumped over the echo
	output := vm.GetOutput()
	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
}

// ============================================================================
// Arithmetic Opcode Tests
// ============================================================================

func TestExecute_Sub(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), int64(3)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSub, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Mul(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(6), int64(7)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpMul, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Div(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(20), int64(4)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpDiv, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Mod(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(17), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpMod, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Pow(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(2), int64(8)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpPow, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_ArithmeticFloats(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{3.14, 2.0}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpAdd, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_DivisionByZero(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), int64(0)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpDiv, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err == nil {
		t.Error("Expected division by zero error")
	}
}

func TestExecute_ModByZero(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), int64(0)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpMod, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err == nil {
		t.Error("Expected modulo by zero error")
	}
}

func TestExecute_MixedTypeArithmetic(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), 3.5}

	// Test int + float
	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSub, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_MultiplyFloats(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{2.5, 4.0}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpMul, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

// ============================================================================
// Comparison Opcode Tests
// ============================================================================

func TestExecute_IsNotEqual(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(3)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsNotEqual, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_IsIdentical(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsIdentical, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_IsNotIdentical(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), "5"}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsNotIdentical, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_IsSmaller(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(3), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsSmaller, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_IsSmallerOrEqual(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsSmallerOrEqual, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Spaceship(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(3), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSpaceship, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Spaceship_Equal(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSpaceship, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_Spaceship_Greater(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(10), int64(3)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSpaceship, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_ComparisonFloats(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{3.5, 2.5}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpIsSmaller, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

// ============================================================================
// Logic & Bitwise Opcode Tests
// ============================================================================

func TestExecute_BoolNot(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{true}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpBoolNot, 2).WithOp1(OpCV, 0).WithResult(OpCV, 1),
		*NewInstruction(OpReturn, 3).WithOp1(OpCV, 1),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_BWNot(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpBWNot, 2).WithOp1(OpCV, 0).WithResult(OpCV, 1),
		*NewInstruction(OpReturn, 3).WithOp1(OpCV, 1),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_BWAnd(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(12), int64(10)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpBWAnd, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_BWOr(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(12), int64(10)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpBWOr, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_BWXor(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(12), int64(10)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpBWXor, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_ShiftLeft(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(5), int64(2)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSL, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

func TestExecute_ShiftRight(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{int64(20), int64(2)}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpSR, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

// ============================================================================
// String Opcode Tests
// ============================================================================

func TestExecute_Concat(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{"Hello, ", "World!"}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpFetchConstant, 2).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpConcat, 3).WithOp1(OpCV, 0).WithOp2(OpCV, 1).WithResult(OpCV, 2),
		*NewInstruction(OpReturn, 4).WithOp1(OpCV, 2),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
}

// ============================================================================
// Control Flow Opcode Tests
// ============================================================================

func TestExecute_JmpZ(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{false, "Skipped", "Executed"}

	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0),
		*NewInstruction(OpJmpZ, 2).WithOp1(OpCV, 0).WithOp2(OpConst, 4),
		*NewInstruction(OpFetchConstant, 3).WithOp1(OpConst, 1).WithResult(OpCV, 1),
		*NewInstruction(OpEcho, 4).WithOp1(OpCV, 1),
		*NewInstruction(OpFetchConstant, 5).WithOp1(OpConst, 2).WithResult(OpCV, 2),
		*NewInstruction(OpEcho, 6).WithOp1(OpCV, 2),
		*NewInstruction(OpReturn, 7).WithOp1(OpUnused, 0),
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	output := vm.GetOutput()
	if output != "Executed" {
		t.Errorf("Expected 'Executed', got '%s'", output)
	}
}

func TestExecute_JmpNZ(t *testing.T) {
	vm := New()
	vm.constants = []interface{}{true, "Executed"}

	// Test: if (true) { jump to echo }
	instructions := Instructions{
		*NewInstruction(OpFetchConstant, 1).WithOp1(OpConst, 0).WithResult(OpCV, 0), // 0: Load true
		*NewInstruction(OpJmpNZ, 2).WithOp1(OpCV, 0).WithOp2(OpConst, 3),           // 1: If true, jump to 3
		*NewInstruction(OpReturn, 3).WithOp1(OpUnused, 0),                          // 2: Return (skipped)
		*NewInstruction(OpFetchConstant, 4).WithOp1(OpConst, 1).WithResult(OpCV, 1), // 3: Load "Executed"
		*NewInstruction(OpEcho, 5).WithOp1(OpCV, 1),                                // 4: Echo "Executed"
		*NewInstruction(OpReturn, 6).WithOp1(OpUnused, 0),                          // 5: Return
	}

	err := vm.Execute(instructions)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	output := vm.GetOutput()
	if output != "Executed" {
		t.Errorf("Expected 'Executed', got '%s'", output)
	}
}

// ============================================================================
// Global Variable Tests
// ============================================================================

// TestGlobalVariableBinding tests that global statement binds local to global variable
func TestGlobalVariableBinding(t *testing.T) {
	// Set up a VM with a global variable
	vm := New()
	vm.SetGlobal("myGlobal", types.NewInt(42))

	// Constants: 0 = "myGlobal" (variable name)
	constants := []interface{}{"myGlobal"}
	vm.constants = constants

	// Create a function that uses global statement
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 2,
		NumParams: 0,
		NumCVs:    1,
	}

	frame := NewFrame(fn)

	// Simulate BIND_GLOBAL instruction
	instr := Instruction{
		Opcode: OpBindGlobal,
		Op1:    Operand{Type: OpCV, Value: 0},    // local variable index
		Op2:    Operand{Type: OpConst, Value: 0}, // constant index for "myGlobal"
	}

	err := vm.opBindGlobal(frame, instr)
	if err != nil {
		t.Fatalf("opBindGlobal failed: %v", err)
	}

	// Verify the binding was recorded
	if frame.globalBindings == nil {
		t.Fatal("globalBindings was not initialized")
	}

	if name, ok := frame.globalBindings[0]; !ok {
		t.Error("local index 0 not bound to global")
	} else if name != "myGlobal" {
		t.Errorf("Expected binding to 'myGlobal', got '%s'", name)
	}
}

// TestGlobalVariableRead tests reading a global variable through binding
func TestGlobalVariableRead(t *testing.T) {
	vm := New()
	vm.SetGlobal("counter", types.NewInt(100))
	vm.constants = []interface{}{"counter"}

	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 2,
		NumParams: 0,
		NumCVs:    1,
	}

	frame := NewFrame(fn)

	// Bind local 0 to global "counter"
	bindInstr := Instruction{
		Opcode: OpBindGlobal,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpConst, Value: 0},
	}
	err := vm.opBindGlobal(frame, bindInstr)
	if err != nil {
		t.Fatalf("opBindGlobal failed: %v", err)
	}

	// Now read the value using getOperandValue
	val, err := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 0})
	if err != nil {
		t.Fatalf("getOperandValue failed: %v", err)
	}

	if val.ToInt() != 100 {
		t.Errorf("Expected 100, got %d", val.ToInt())
	}
}

// TestGlobalVariableWrite tests writing to a global variable through binding
func TestGlobalVariableWrite(t *testing.T) {
	vm := New()
	vm.SetGlobal("counter", types.NewInt(100))
	vm.constants = []interface{}{"counter"}

	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 2,
		NumParams: 0,
		NumCVs:    1,
	}

	frame := NewFrame(fn)

	// Bind local 0 to global "counter"
	bindInstr := Instruction{
		Opcode: OpBindGlobal,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpConst, Value: 0},
	}
	err := vm.opBindGlobal(frame, bindInstr)
	if err != nil {
		t.Fatalf("opBindGlobal failed: %v", err)
	}

	// Write a new value using setOperandValue
	err = vm.setOperandValue(frame, Operand{Type: OpCV, Value: 0}, types.NewInt(200))
	if err != nil {
		t.Fatalf("setOperandValue failed: %v", err)
	}

	// Verify the global was updated
	globalVal, ok := vm.GetGlobal("counter")
	if !ok {
		t.Fatal("Global 'counter' not found")
	}
	if globalVal.ToInt() != 200 {
		t.Errorf("Expected global value 200, got %d", globalVal.ToInt())
	}
}

// TestMultipleGlobalBindings tests multiple global variable bindings
func TestMultipleGlobalBindings(t *testing.T) {
	vm := New()
	vm.SetGlobal("var1", types.NewInt(10))
	vm.SetGlobal("var2", types.NewString("hello"))
	vm.SetGlobal("var3", types.NewBool(true))
	vm.constants = []interface{}{"var1", "var2", "var3"}

	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 5,
		NumParams: 0,
		NumCVs:    3,
	}

	frame := NewFrame(fn)

	// Bind three local variables to three globals
	bindings := []struct {
		localIdx int
		constIdx int
		name     string
	}{
		{0, 0, "var1"},
		{1, 1, "var2"},
		{2, 2, "var3"},
	}

	for _, b := range bindings {
		instr := Instruction{
			Opcode: OpBindGlobal,
			Op1:    Operand{Type: OpCV, Value: uint32(b.localIdx)},
			Op2:    Operand{Type: OpConst, Value: uint32(b.constIdx)},
		}
		err := vm.opBindGlobal(frame, instr)
		if err != nil {
			t.Fatalf("opBindGlobal for %s failed: %v", b.name, err)
		}
	}

	// Verify all bindings
	if len(frame.globalBindings) != 3 {
		t.Errorf("Expected 3 bindings, got %d", len(frame.globalBindings))
	}

	// Read all values
	val1, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 0})
	val2, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 1})
	val3, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 2})

	if val1.ToInt() != 10 {
		t.Errorf("var1: expected 10, got %d", val1.ToInt())
	}
	if val2.ToString() != "hello" {
		t.Errorf("var2: expected 'hello', got '%s'", val2.ToString())
	}
	if !val3.ToBool() {
		t.Error("var3: expected true")
	}
}

// TestGlobalVariableNewCreation tests that binding to non-existent global creates it
func TestGlobalVariableNewCreation(t *testing.T) {
	vm := New()
	// Don't set "newVar" - it should be created by BIND_GLOBAL
	vm.constants = []interface{}{"newVar"}

	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 2,
		NumParams: 0,
		NumCVs:    1,
	}

	frame := NewFrame(fn)

	// Bind local 0 to non-existent global "newVar"
	instr := Instruction{
		Opcode: OpBindGlobal,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpConst, Value: 0},
	}
	err := vm.opBindGlobal(frame, instr)
	if err != nil {
		t.Fatalf("opBindGlobal failed: %v", err)
	}

	// Verify the global was created (with null value)
	if vm.globals["newVar"] == nil {
		t.Fatal("Global 'newVar' was not created")
	}

	// Write a value through the binding
	err = vm.setOperandValue(frame, Operand{Type: OpCV, Value: 0}, types.NewString("created"))
	if err != nil {
		t.Fatalf("setOperandValue failed: %v", err)
	}

	// Verify the global has the new value
	if vm.globals["newVar"].ToString() != "created" {
		t.Errorf("Expected 'created', got '%s'", vm.globals["newVar"].ToString())
	}
}

// TestGlobalVariableUnboundLocalAccess tests that unbound locals work normally
func TestGlobalVariableUnboundLocalAccess(t *testing.T) {
	vm := New()
	vm.SetGlobal("myGlobal", types.NewInt(999))
	vm.constants = []interface{}{"myGlobal"}

	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 3,
		NumParams: 0,
		NumCVs:    2,
	}

	frame := NewFrame(fn)

	// Only bind local 0 to global
	instr := Instruction{
		Opcode: OpBindGlobal,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpConst, Value: 0},
	}
	_ = vm.opBindGlobal(frame, instr)

	// Set local 1 directly (not bound to global)
	frame.setLocal(1, types.NewInt(123))

	// Read local 0 (should be from global)
	val0, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 0})
	if val0.ToInt() != 999 {
		t.Errorf("local 0: expected 999 from global, got %d", val0.ToInt())
	}

	// Read local 1 (should be from local storage)
	val1, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 1})
	if val1.ToInt() != 123 {
		t.Errorf("local 1: expected 123 from local, got %d", val1.ToInt())
	}

	// Write to local 1 should not affect globals
	_ = vm.setOperandValue(frame, Operand{Type: OpCV, Value: 1}, types.NewInt(456))
	val1After, _ := vm.getOperandValue(frame, Operand{Type: OpCV, Value: 1})
	if val1After.ToInt() != 456 {
		t.Errorf("local 1 after write: expected 456, got %d", val1After.ToInt())
	}

	// Global should still be 999 (local 1 write didn't affect it)
	if vm.globals["myGlobal"].ToInt() != 999 {
		t.Errorf("Global should still be 999, got %d", vm.globals["myGlobal"].ToInt())
	}
}

// ============================================================================
// Unset Tests
// ============================================================================

// TestOpUnsetVar tests unsetting a local variable
func TestOpUnsetVar(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up a variable
	frame.setLocal(0, types.NewInt(42))

	// Verify it's set
	if frame.getLocal(0).Type() != types.TypeInt {
		t.Fatalf("Variable should be set to Int, got %v", frame.getLocal(0).Type())
	}

	// Unset the variable
	instr := Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 0},
	}

	err := vm.opUnsetVar(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetVar failed: %v", err)
	}

	// Verify the variable is now undefined
	if frame.getLocal(0).Type() != types.TypeUndef {
		t.Errorf("Variable should be Undef after unset, got %v", frame.getLocal(0).Type())
	}
}

// TestOpUnsetVar_String tests unsetting a string variable
func TestOpUnsetVar_String(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up a string variable
	frame.setLocal(1, types.NewString("hello"))

	// Unset
	instr := Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 1},
	}

	err := vm.opUnsetVar(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetVar failed: %v", err)
	}

	if frame.getLocal(1).Type() != types.TypeUndef {
		t.Errorf("String variable should be Undef after unset, got %v", frame.getLocal(1).Type())
	}
}

// TestOpUnsetVar_Array tests unsetting an array variable
func TestOpUnsetVar_Array(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up an array variable
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("key"), types.NewString("value"))
	frame.setLocal(2, types.NewArray(arr))

	// Unset the whole array
	instr := Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 2},
	}

	err := vm.opUnsetVar(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetVar failed: %v", err)
	}

	if frame.getLocal(2).Type() != types.TypeUndef {
		t.Errorf("Array variable should be Undef after unset, got %v", frame.getLocal(2).Type())
	}
}

// TestOpUnsetVar_AlreadyUnset tests unsetting an already unset variable (no error)
func TestOpUnsetVar_AlreadyUnset(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Variable is not set (already Undef by default)
	instr := Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 3},
	}

	// This should not error - PHP allows unsetting undefined variables
	err := vm.opUnsetVar(frame, instr)
	if err != nil {
		t.Errorf("opUnsetVar should not error on already unset variable: %v", err)
	}

	if frame.getLocal(3).Type() != types.TypeUndef {
		t.Errorf("Variable should remain Undef, got %v", frame.getLocal(3).Type())
	}
}

// TestOpUnsetDim tests unsetting array elements
func TestOpUnsetDim(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up an array
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))
	arr.Set(types.NewInt(0), types.NewString("first"))
	frame.setLocal(0, types.NewArray(arr))

	// Unset $arr["age"]
	frame.setLocal(1, types.NewString("age"))
	instr := Instruction{
		Opcode: OpUnsetDim,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opUnsetDim(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetDim failed: %v", err)
	}

	// Verify "age" was deleted
	_, exists := arr.Get(types.NewString("age"))
	if exists {
		t.Error("Key 'age' should have been deleted from array")
	}

	// Verify other elements still exist
	val, exists := arr.Get(types.NewString("name"))
	if !exists || val.ToString() != "John" {
		t.Error("Key 'name' should still exist with value 'John'")
	}

	val, exists = arr.Get(types.NewInt(0))
	if !exists || val.ToString() != "first" {
		t.Error("Key 0 should still exist with value 'first'")
	}
}

// TestOpUnsetDim_IntKey tests unsetting array element with integer key
func TestOpUnsetDim_IntKey(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up an array with integer keys
	arr := types.NewEmptyArray()
	arr.Set(types.NewInt(0), types.NewString("zero"))
	arr.Set(types.NewInt(1), types.NewString("one"))
	arr.Set(types.NewInt(2), types.NewString("two"))
	frame.setLocal(0, types.NewArray(arr))

	// Unset $arr[1]
	frame.setLocal(1, types.NewInt(1))
	instr := Instruction{
		Opcode: OpUnsetDim,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	err := vm.opUnsetDim(frame, instr)
	if err != nil {
		t.Fatalf("opUnsetDim failed: %v", err)
	}

	// Verify index 1 was deleted
	_, exists := arr.Get(types.NewInt(1))
	if exists {
		t.Error("Index 1 should have been deleted from array")
	}

	// Verify other indices still exist
	if arr.Len() != 2 {
		t.Errorf("Array should have 2 elements, got %d", arr.Len())
	}
}

// TestOpUnsetDim_NonExistentKey tests unsetting non-existent key (no error in PHP)
func TestOpUnsetDim_NonExistentKey(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up an array
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	frame.setLocal(0, types.NewArray(arr))

	// Try to unset non-existent key
	frame.setLocal(1, types.NewString("nonexistent"))
	instr := Instruction{
		Opcode: OpUnsetDim,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	// Should not error
	err := vm.opUnsetDim(frame, instr)
	if err != nil {
		t.Errorf("opUnsetDim should not error on non-existent key: %v", err)
	}

	// Array should be unchanged
	if arr.Len() != 1 {
		t.Errorf("Array should still have 1 element, got %d", arr.Len())
	}
}

// TestOpUnsetDim_NonArray tests unsetting on non-array (no-op in PHP)
func TestOpUnsetDim_NonArray(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 5}
	frame := NewFrame(fn)

	// Set up a non-array value
	frame.setLocal(0, types.NewString("hello"))
	frame.setLocal(1, types.NewInt(0))

	instr := Instruction{
		Opcode: OpUnsetDim,
		Op1:    Operand{Type: OpCV, Value: 0},
		Op2:    Operand{Type: OpTmpVar, Value: 1},
	}

	// Should not error - it's a no-op in PHP
	err := vm.opUnsetDim(frame, instr)
	if err != nil {
		t.Errorf("opUnsetDim should not error on non-array: %v", err)
	}

	// String should be unchanged
	if frame.getLocal(0).ToString() != "hello" {
		t.Error("String value should be unchanged")
	}
}

// TestMultipleUnsets tests unsetting multiple things in sequence
func TestMultipleUnsets(t *testing.T) {
	vm := New()

	fn := &CompiledFunction{Instructions: Instructions{}, NumLocals: 10}
	frame := NewFrame(fn)

	// Set up variables
	frame.setLocal(0, types.NewInt(10))
	frame.setLocal(1, types.NewString("hello"))
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(2))
	frame.setLocal(2, types.NewArray(arr))

	// Unset $var0
	err := vm.opUnsetVar(frame, Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 0},
	})
	if err != nil {
		t.Fatalf("First unset failed: %v", err)
	}

	// Unset $var1
	err = vm.opUnsetVar(frame, Instruction{
		Opcode: OpUnsetVar,
		Op1:    Operand{Type: OpCV, Value: 1},
	})
	if err != nil {
		t.Fatalf("Second unset failed: %v", err)
	}

	// Unset $arr["a"]
	frame.setLocal(3, types.NewString("a"))
	err = vm.opUnsetDim(frame, Instruction{
		Opcode: OpUnsetDim,
		Op1:    Operand{Type: OpCV, Value: 2},
		Op2:    Operand{Type: OpTmpVar, Value: 3},
	})
	if err != nil {
		t.Fatalf("Third unset failed: %v", err)
	}

	// Verify all unsets worked
	if frame.getLocal(0).Type() != types.TypeUndef {
		t.Error("var0 should be Undef")
	}
	if frame.getLocal(1).Type() != types.TypeUndef {
		t.Error("var1 should be Undef")
	}

	// Array should only have key "b"
	_, exists := arr.Get(types.NewString("a"))
	if exists {
		t.Error("Array key 'a' should have been deleted")
	}
	val, exists := arr.Get(types.NewString("b"))
	if !exists || val.ToInt() != 2 {
		t.Error("Array key 'b' should still exist with value 2")
	}
}

// ============================================================================
// Misc Functions Tests (define, defined, constant, etc.)
// ============================================================================

func TestDefineAndDefined(t *testing.T) {
	vm := New()

	// Test defined on non-existent constant
	result := vm.funcDefined([]*types.Value{types.NewString("MY_CONST")})
	if result.ToBool() {
		t.Error("MY_CONST should not be defined yet")
	}

	// Define a constant
	result = vm.funcDefine([]*types.Value{
		types.NewString("MY_CONST"),
		types.NewInt(42),
	})
	if !result.ToBool() {
		t.Error("define() should return true for new constant")
	}

	// Test defined on existing constant
	result = vm.funcDefined([]*types.Value{types.NewString("MY_CONST")})
	if !result.ToBool() {
		t.Error("MY_CONST should be defined after define()")
	}

	// Test redefining (should fail)
	result = vm.funcDefine([]*types.Value{
		types.NewString("MY_CONST"),
		types.NewInt(100),
	})
	if result.ToBool() {
		t.Error("define() should return false when constant already exists")
	}
}

func TestConstant(t *testing.T) {
	vm := New()

	// Define a constant
	vm.funcDefine([]*types.Value{
		types.NewString("TEST_VALUE"),
		types.NewString("hello world"),
	})

	// Get the constant value
	result := vm.funcConstant([]*types.Value{types.NewString("TEST_VALUE")})
	if result.ToString() != "hello world" {
		t.Errorf("constant() should return 'hello world', got '%s'", result.ToString())
	}

	// Test non-existent constant
	result = vm.funcConstant([]*types.Value{types.NewString("NONEXISTENT")})
	if result.Type() != types.TypeNull {
		t.Error("constant() should return null for non-existent constant")
	}
}

func TestBuiltinConstants(t *testing.T) {
	vm := New()

	// Test PHP_VERSION
	result := vm.funcDefined([]*types.Value{types.NewString("PHP_VERSION")})
	if !result.ToBool() {
		t.Error("PHP_VERSION should be defined")
	}

	result = vm.funcConstant([]*types.Value{types.NewString("PHP_VERSION")})
	if result.ToString() != "8.4.0" {
		t.Errorf("PHP_VERSION should be '8.4.0', got '%s'", result.ToString())
	}

	// Test PHP_EOL
	result = vm.funcConstant([]*types.Value{types.NewString("PHP_EOL")})
	if result.ToString() != "\n" {
		t.Error("PHP_EOL should be newline character")
	}

	// Test PHP_INT_MAX
	result = vm.funcConstant([]*types.Value{types.NewString("PHP_INT_MAX")})
	if result.ToInt() != 9223372036854775807 {
		t.Errorf("PHP_INT_MAX should be 9223372036854775807, got %d", result.ToInt())
	}
}

func TestFunctionExistsWithVM(t *testing.T) {
	vm := New()

	// Test builtin function
	result := vm.funcFunctionExists([]*types.Value{types.NewString("strlen")})
	if !result.ToBool() {
		t.Error("strlen should exist as builtin")
	}

	// Test non-existent function
	result = vm.funcFunctionExists([]*types.Value{types.NewString("nonexistent_func")})
	if result.ToBool() {
		t.Error("nonexistent_func should not exist")
	}

	// Register a user function and test
	vm.RegisterFunction("my_custom_func", &CompiledFunction{
		Name:         "my_custom_func",
		Instructions: Instructions{},
		NumParams:    0,
	})
	result = vm.funcFunctionExists([]*types.Value{types.NewString("my_custom_func")})
	if !result.ToBool() {
		t.Error("my_custom_func should exist after registration")
	}
}

func TestFuncNumArgs(t *testing.T) {
	vm := New()

	// Test with no frame
	result := vm.funcNumArgs(nil)
	if result.ToInt() != 0 {
		t.Error("func_num_args with nil frame should return 0")
	}

	// Test with frame
	fn := &CompiledFunction{
		Name:      "test_func",
		NumParams: 3,
	}
	frame := NewFrame(fn)
	result = vm.funcNumArgs(frame)
	if result.ToInt() != 3 {
		t.Errorf("func_num_args should return 3, got %d", result.ToInt())
	}
}

func TestFuncGetArgs(t *testing.T) {
	vm := New()

	// Test with no frame
	result := vm.funcGetArgs(nil)
	arr := result.ToArray()
	if arr.Len() != 0 {
		t.Error("func_get_args with nil frame should return empty array")
	}

	// Test with frame that has parameters
	fn := &CompiledFunction{
		Name:      "test_func",
		NumParams: 2,
		NumLocals: 2,
	}
	frame := NewFrame(fn)
	frame.setParam(0, types.NewString("arg1"))
	frame.setParam(1, types.NewInt(42))

	result = vm.funcGetArgs(frame)
	arr = result.ToArray()
	if arr.Len() != 2 {
		t.Errorf("func_get_args should return array with 2 elements, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToString() != "arg1" {
		t.Errorf("First arg should be 'arg1', got '%s'", val.ToString())
	}

	val, _ = arr.Get(types.NewInt(1))
	if val.ToInt() != 42 {
		t.Errorf("Second arg should be 42, got %d", val.ToInt())
	}
}
