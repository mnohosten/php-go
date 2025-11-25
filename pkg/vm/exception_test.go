package vm

import (
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/runtime"
	"github.com/krizos/php-go/pkg/types"
)

// TestFrameExceptionMethods tests frame exception handling methods
func TestFrameExceptionMethods(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "test",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	frame := NewFrame(fn)

	// Initially no exception
	if frame.HasException() {
		t.Error("Frame should not have exception initially")
	}

	if frame.GetException() != nil {
		t.Error("GetException should return nil initially")
	}

	// Set exception
	exc := runtime.NewException("Test error", 123, nil)
	frame.SetException(exc)

	if !frame.HasException() {
		t.Error("Frame should have exception after SetException")
	}

	retrieved := frame.GetException()
	if retrieved == nil {
		t.Fatal("GetException should return exception")
	}

	if retrieved.Message != "Test error" {
		t.Error("Retrieved exception has wrong message")
	}

	// Clear exception
	frame.ClearException()

	if frame.HasException() {
		t.Error("Frame should not have exception after Clear")
	}

	if frame.GetException() != nil {
		t.Error("GetException should return nil after Clear")
	}
}

// TestGenerateStackTrace tests stack trace generation
func TestGenerateStackTrace(t *testing.T) {
	vm := New()

	// Create some nested frames
	fn1 := &CompiledFunction{
		Name:         "main",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	fn2 := &CompiledFunction{
		Name:         "foo",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	fn3 := &CompiledFunction{
		Name:         "bar",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	frame1 := NewFrame(fn1)
	frame1.lastLine = 10

	frame2 := NewFrame(fn2)
	frame2.lastLine = 20

	frame3 := NewFrame(fn3)
	frame3.lastLine = 30

	vm.pushFrame(frame1)
	vm.pushFrame(frame2)
	vm.pushFrame(frame3)

	// Generate trace
	trace := vm.generateStackTrace()

	if len(trace.Frames) != 3 {
		t.Errorf("Expected 3 frames, got %d", len(trace.Frames))
	}

	// Check frames (in reverse order - bar, foo, main)
	if trace.Frames[0].Function != "bar" {
		t.Errorf("Expected first frame to be 'bar', got '%s'", trace.Frames[0].Function)
	}

	if trace.Frames[1].Function != "foo" {
		t.Errorf("Expected second frame to be 'foo', got '%s'", trace.Frames[1].Function)
	}

	if trace.Frames[2].Function != "main" {
		t.Errorf("Expected third frame to be 'main', got '%s'", trace.Frames[2].Function)
	}

	// Check line numbers
	if trace.Frames[0].Line != 30 {
		t.Errorf("Expected line 30, got %d", trace.Frames[0].Line)
	}

	if trace.Frames[1].Line != 20 {
		t.Errorf("Expected line 20, got %d", trace.Frames[1].Line)
	}

	if trace.Frames[2].Line != 10 {
		t.Errorf("Expected line 10, got %d", trace.Frames[2].Line)
	}
}

// TestFindCatchHandler tests finding catch handlers
func TestFindCatchHandler(t *testing.T) {
	vm := New()

	// Create function with catch handler
	// OpCatch with empty type (catches all) is at index 2
	instructions := Instructions{
		{Opcode: OpFetchConstant, Lineno: 1},
		{Opcode: OpFetchConstant, Lineno: 2},
		{Opcode: OpCatch, Op1: Operand{Type: OpConst, Value: 0}, Lineno: 3}, // Catch handler with empty type
		{Opcode: OpReturn, Lineno: 4},
	}

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: instructions,
		NumLocals:    10,
		NumParams:    0,
	}

	frame := NewFrame(fn)
	frame.ip = 0
	// Set exception so findCatchHandler can work
	frame.exception = runtime.NewException("test", 0, nil)

	// Should find catch handler
	found := vm.findCatchHandler(frame)
	if !found {
		t.Error("Should find catch handler")
	}

	if frame.ip != 2 {
		t.Errorf("IP should be at catch instruction (2), got %d", frame.ip)
	}
}

// TestFindCatchHandlerNotFound tests when no catch handler exists
func TestFindCatchHandlerNotFound(t *testing.T) {
	vm := New()

	// Create function without catch handler
	instructions := Instructions{
		{Opcode: OpFetchConstant, Lineno: 1},
		{Opcode: OpFetchConstant, Lineno: 2},
		{Opcode: OpReturn, Lineno: 3}, // Return ends search
	}

	fn := &CompiledFunction{
		Name:         "test",
		Instructions: instructions,
		NumLocals:    10,
		NumParams:    0,
	}

	frame := NewFrame(fn)
	frame.ip = 0
	// Set exception so findCatchHandler can work
	frame.exception = runtime.NewException("test", 0, nil)

	// Should not find catch handler
	found := vm.findCatchHandler(frame)
	if found {
		t.Error("Should not find catch handler")
	}
}

// TestExceptionWrapping tests wrapping and unwrapping exceptions
func TestExceptionWrapping(t *testing.T) {
	exc := runtime.NewException("Test error", 123, nil)

	// Wrap as value (now wraps as object, not resource)
	val := runtime.WrapExceptionAsValue(exc)

	if val.Type() != types.TypeObject {
		t.Errorf("Exception should be wrapped as object, got %v", val.Type())
	}

	// Unwrap
	unwrapped, err := runtime.GetExceptionFromValue(val)
	if err != nil {
		t.Fatalf("Failed to unwrap exception: %v", err)
	}

	if unwrapped.Message != "Test error" {
		t.Error("Unwrapped exception has wrong message")
	}

	if unwrapped.Code != 123 {
		t.Error("Unwrapped exception has wrong code")
	}

	// When wrapped as object, FromObject creates a new Exception from the object properties
	// So we check that values match, not object identity
	if unwrapped.Message != exc.Message || unwrapped.Code != exc.Code {
		t.Error("Unwrapped exception should have same values")
	}
}

// TestExceptionFromObject tests creating exception from PHP object
func TestExceptionFromObject(t *testing.T) {
	// Create PHP object representing an exception
	obj := &types.Object{
		ClassName:  "Exception",
		Properties: make(map[string]*types.Property),
	}

	obj.Properties["message"] = &types.Property{
		Value:      types.NewString("Object error"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["code"] = &types.Property{
		Value:      types.NewInt(456),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["file"] = &types.Property{
		Value:      types.NewString("test.php"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["line"] = &types.Property{
		Value:      types.NewInt(99),
		Visibility: types.VisibilityPublic,
	}

	objVal := types.NewObject(obj)

	// Extract exception
	exc, err := runtime.GetExceptionFromValue(objVal)
	if err != nil {
		t.Fatalf("Failed to extract exception: %v", err)
	}

	if exc.Message != "Object error" {
		t.Error("Exception message incorrect")
	}

	if exc.Code != 456 {
		t.Error("Exception code incorrect")
	}

	if exc.File != "test.php" {
		t.Error("Exception file incorrect")
	}

	if exc.Line != 99 {
		t.Error("Exception line incorrect")
	}
}

// TestStackTraceWithClassContext tests stack trace with class context
func TestStackTraceWithClassContext(t *testing.T) {
	vm := New()

	// Create frame with class context
	fn := &CompiledFunction{
		Name:         "myMethod",
		Instructions: make(Instructions, 0),
		NumLocals:    10,
		NumParams:    0,
	}

	classEntry := &types.ClassEntry{
		Name: "MyClass",
	}

	frame := NewFrame(fn)
	frame.lastLine = 42
	frame.classEntry = classEntry

	vm.pushFrame(frame)

	// Generate trace
	trace := vm.generateStackTrace()

	if len(trace.Frames) != 1 {
		t.Fatalf("Expected 1 frame, got %d", len(trace.Frames))
	}

	stackFrame := trace.Frames[0]

	if stackFrame.Class != "MyClass" {
		t.Errorf("Expected class 'MyClass', got '%s'", stackFrame.Class)
	}

	if stackFrame.Type != "->" {
		t.Errorf("Expected type '->',  got '%s'", stackFrame.Type)
	}

	if stackFrame.Function != "myMethod" {
		t.Errorf("Expected function 'myMethod', got '%s'", stackFrame.Function)
	}
}

// TestExceptionChaining tests exception chaining
func TestExceptionChaining(t *testing.T) {
	// Create chain: exc3 -> exc2 -> exc1
	exc1 := runtime.NewException("First error", 1, nil)
	exc2 := runtime.NewException("Second error", 2, exc1)
	exc3 := runtime.NewException("Third error", 3, exc2)

	// Check chain
	if exc3.GetPrevious() != exc2 {
		t.Error("exc3 previous should be exc2")
	}

	if exc2.GetPrevious() != exc1 {
		t.Error("exc2 previous should be exc1")
	}

	if exc1.GetPrevious() != nil {
		t.Error("exc1 previous should be nil")
	}

	// Convert to object and check previous is preserved
	obj := exc3.ToObject()

	if prevProp, ok := obj.Properties["previous"]; ok {
		if prevProp.Value.Type() != types.TypeObject {
			t.Error("Previous should be object")
		}

		prevObj := prevProp.Value.ToObject()
		if msgProp, ok := prevObj.Properties["message"]; ok {
			if msgProp.Value.ToString() != "Second error" {
				t.Error("Previous exception message wrong")
			}
		}
	} else {
		t.Error("Previous property not found")
	}
}

// TestExceptionSubclasses tests different exception subclasses
func TestExceptionSubclasses(t *testing.T) {
	tests := []struct {
		name      string
		exception interface {
			GetMessage() string
			GetCode() int64
		}
		expectedClass string
	}{
		{"ErrorException", runtime.NewErrorException("Error", 1, runtime.E_ERROR, nil), "ErrorException"},
		{"LogicException", runtime.NewLogicException("Logic", 2, nil), "LogicException"},
		{"RuntimeException", runtime.NewRuntimeException("Runtime", 3, nil), "RuntimeException"},
		{"InvalidArgumentException", runtime.NewInvalidArgumentException("Invalid", 4, nil), "InvalidArgumentException"},
		{"OutOfBoundsException", runtime.NewOutOfBoundsException("Bounds", 5, nil), "OutOfBoundsException"},
		{"OutOfRangeException", runtime.NewOutOfRangeException("Range", 6, nil), "OutOfRangeException"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.exception.GetMessage() == "" {
				t.Error("Message should not be empty")
			}

			if tt.exception.GetCode() == 0 {
				t.Error("Code should not be zero")
			}

			// Check class name via type assertion
			switch exc := tt.exception.(type) {
			case *runtime.ErrorException:
				if exc.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.ClassName)
				}
			case *runtime.LogicException:
				if exc.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.ClassName)
				}
			case *runtime.RuntimeException:
				if exc.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.ClassName)
				}
			case *runtime.InvalidArgumentException:
				if exc.Exception.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.Exception.ClassName)
				}
			case *runtime.OutOfBoundsException:
				if exc.Exception.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.Exception.ClassName)
				}
			case *runtime.OutOfRangeException:
				if exc.Exception.ClassName != tt.expectedClass {
					t.Errorf("Expected class %s, got %s", tt.expectedClass, exc.Exception.ClassName)
				}
			}
		})
	}
}

// TestExceptionTraceString tests stack trace string generation
func TestExceptionTraceString(t *testing.T) {
	exc := runtime.NewException("Test error", 123, nil)
	exc.File = "test.php"
	exc.Line = 42

	// Add stack frames
	exc.Trace.AddFrame(&runtime.StackFrame{
		File:     "foo.php",
		Line:     10,
		Function: "foo",
	})
	exc.Trace.AddFrame(&runtime.StackFrame{
		File:     "bar.php",
		Line:     20,
		Function: "bar",
		Class:    "MyClass",
		Type:     "->",
	})

	traceStr := exc.GetTraceAsString()

	if !strings.Contains(traceStr, "foo.php") {
		t.Error("Trace should contain foo.php")
	}

	if !strings.Contains(traceStr, "bar.php") {
		t.Error("Trace should contain bar.php")
	}

	if !strings.Contains(traceStr, "MyClass") {
		t.Error("Trace should contain MyClass")
	}

	// String() should include message and trace
	str := exc.String()

	if !strings.Contains(str, "Exception: Test error") {
		t.Error("String should contain exception type and message")
	}

	if !strings.Contains(str, "test.php:42") {
		t.Error("String should contain file and line")
	}

	if !strings.Contains(str, "Stack trace") {
		t.Error("String should include stack trace")
	}
}

// TestIsException tests exception detection
func TestIsException(t *testing.T) {
	// Create exception object
	exc := runtime.NewException("Test", 0, nil)
	obj := exc.ToObject()
	val := types.NewObject(obj)

	if !runtime.IsException(val) {
		t.Error("Should detect exception object")
	}

	// Test with "Error" class
	obj2 := &types.Object{
		ClassName: "Error",
	}
	val2 := types.NewObject(obj2)

	if !runtime.IsException(val2) {
		t.Error("Should detect Error as exception")
	}

	// Test with "Throwable"
	obj3 := &types.Object{
		ClassName: "Throwable",
	}
	val3 := types.NewObject(obj3)

	if !runtime.IsException(val3) {
		t.Error("Should detect Throwable as exception")
	}

	// Non-exception
	nonExc := types.NewInt(42)
	if runtime.IsException(nonExc) {
		t.Error("Should not detect int as exception")
	}
}
