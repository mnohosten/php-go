package vm

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

func TestNewFrame(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "test",
		Instructions: Instructions{},
		NumLocals:    20,
		NumParams:    3,
	}

	frame := NewFrame(fn)

	if frame == nil {
		t.Fatal("NewFrame() returned nil")
	}

	if frame.fn != fn {
		t.Error("Frame function not set correctly")
	}

	if frame.ip != 0 {
		t.Errorf("Expected ip 0, got %d", frame.ip)
	}

	if len(frame.locals) < 10 {
		t.Errorf("Expected at least 10 locals, got %d", len(frame.locals))
	}
}

func TestFrame_GetSetLocal(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Set a local
	val := types.NewInt(42)
	frame.setLocal(5, val)

	// Get it back
	retrieved := frame.getLocal(5)
	if retrieved.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", retrieved.ToInt())
	}
}

func TestFrame_GetLocal_OutOfBounds(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Should return null for out of bounds
	val := frame.getLocal(100)
	if !val.IsNull() {
		t.Error("Expected null for out of bounds local")
	}

	val = frame.getLocal(-1)
	if !val.IsNull() {
		t.Error("Expected null for negative index")
	}
}

func TestFrame_SetLocal_Expand(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 5,
	}

	frame := NewFrame(fn)
	initialLen := len(frame.locals)

	// Set beyond initial capacity
	frame.setLocal(50, types.NewInt(99))

	if len(frame.locals) <= initialLen {
		t.Error("Locals array should have expanded")
	}

	// Verify the value
	val := frame.getLocal(50)
	if val.ToInt() != 99 {
		t.Errorf("Expected 99, got %d", val.ToInt())
	}
}

func TestFrame_PushPopTemp(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Push a temp value
	val := types.NewInt(42)
	index := frame.pushTemp(val)

	// Pop it back
	retrieved := frame.popTemp(index)
	if retrieved.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", retrieved.ToInt())
	}

	// After pop, slot should be nil
	val2 := frame.getLocal(index)
	if val2 != nil && !val2.IsUndef() && !val2.IsNull() {
		t.Error("Expected temp slot to be cleared after pop")
	}
}

func TestFrame_SetGetParam(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumParams: 3,
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Set parameter
	frame.setParam(0, types.NewInt(10))
	frame.setParam(1, types.NewInt(20))
	frame.setParam(2, types.NewInt(30))

	// Get parameters
	if frame.getParam(0).ToInt() != 10 {
		t.Error("Parameter 0 incorrect")
	}
	if frame.getParam(1).ToInt() != 20 {
		t.Error("Parameter 1 incorrect")
	}
	if frame.getParam(2).ToInt() != 30 {
		t.Error("Parameter 2 incorrect")
	}
}

func TestFrame_SetGetParam_OutOfBounds(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumParams: 2,
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Setting out of bounds parameter should not panic
	frame.setParam(10, types.NewInt(42))

	// Getting out of bounds parameter should return null
	val := frame.getParam(10)
	if !val.IsNull() {
		t.Error("Expected null for out of bounds parameter")
	}
}

func TestFrame_ReturnValue(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "test",
		NumLocals: 10,
	}

	frame := NewFrame(fn)

	// Default return value should be null
	val := frame.getReturnValue()
	if !val.IsNull() {
		t.Error("Expected default return value to be null")
	}

	// Set return value
	frame.setReturnValue(types.NewInt(42))

	// Get it back
	val = frame.getReturnValue()
	if val.ToInt() != 42 {
		t.Errorf("Expected return value 42, got %d", val.ToInt())
	}
}

func TestFrame_String(t *testing.T) {
	fn := &CompiledFunction{
		Name:      "testFunc",
		NumLocals: 10,
	}

	frame := NewFrame(fn)
	frame.ip = 5

	str := frame.String()
	if str == "" {
		t.Error("Frame String() returned empty string")
	}

	// Should contain function name
	// (exact format may vary)
}

func TestFrame_String_Nil(t *testing.T) {
	var frame *Frame
	str := frame.String()

	if str != "<nil frame>" {
		t.Errorf("Expected '<nil frame>', got '%s'", str)
	}
}
