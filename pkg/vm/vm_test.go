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

	// Unconditional jump over echo
	instructions := Instructions{
		// JMP to instruction 3
		*NewInstruction(OpJmp, 1).
			WithOp1(OpConst, 3),
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
