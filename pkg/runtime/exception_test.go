package runtime

import (
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestNewException tests basic exception creation
func TestNewException(t *testing.T) {
	exc := NewException("Test error", 123, nil)

	if exc.Message != "Test error" {
		t.Errorf("Expected message 'Test error', got '%s'", exc.Message)
	}

	if exc.Code != 123 {
		t.Errorf("Expected code 123, got %d", exc.Code)
	}

	if exc.ClassName != "Exception" {
		t.Errorf("Expected className 'Exception', got '%s'", exc.ClassName)
	}

	if exc.Previous != nil {
		t.Error("Expected no previous exception")
	}

	if exc.Trace == nil {
		t.Error("Expected trace to be initialized")
	}
}

// TestExceptionChaining tests exception chaining
func TestExceptionChaining(t *testing.T) {
	previous := NewException("Previous error", 100, nil)
	exc := NewException("Current error", 200, previous)

	if exc.Previous != previous {
		t.Error("Previous exception not linked correctly")
	}

	if exc.GetPrevious() != previous {
		t.Error("GetPrevious() should return previous exception")
	}
}

// TestExceptionGetters tests exception getter methods
func TestExceptionGetters(t *testing.T) {
	exc := NewException("Test message", 42, nil)
	exc.File = "test.php"
	exc.Line = 10

	if exc.GetMessage() != "Test message" {
		t.Error("GetMessage() failed")
	}

	if exc.GetCode() != 42 {
		t.Error("GetCode() failed")
	}

	if exc.GetFile() != "test.php" {
		t.Error("GetFile() failed")
	}

	if exc.GetLine() != 10 {
		t.Error("GetLine() failed")
	}
}

// TestExceptionToObject tests converting exception to PHP object
func TestExceptionToObject(t *testing.T) {
	exc := NewException("Test error", 123, nil)
	exc.File = "test.php"
	exc.Line = 42

	obj := exc.ToObject()

	if obj.ClassName != "Exception" {
		t.Errorf("Expected className 'Exception', got '%s'", obj.ClassName)
	}

	// Check properties
	if msgProp, ok := obj.Properties["message"]; ok {
		if msgProp.Value.ToString() != "Test error" {
			t.Error("Message property incorrect")
		}
	} else {
		t.Error("Message property not found")
	}

	if codeProp, ok := obj.Properties["code"]; ok {
		if codeProp.Value.ToInt() != 123 {
			t.Error("Code property incorrect")
		}
	} else {
		t.Error("Code property not found")
	}

	if fileProp, ok := obj.Properties["file"]; ok {
		if fileProp.Value.ToString() != "test.php" {
			t.Error("File property incorrect")
		}
	} else {
		t.Error("File property not found")
	}

	if lineProp, ok := obj.Properties["line"]; ok {
		if lineProp.Value.ToInt() != 42 {
			t.Error("Line property incorrect")
		}
	} else {
		t.Error("Line property not found")
	}
}

// TestFromObject tests creating exception from PHP object
func TestFromObject(t *testing.T) {
	// Create a PHP object
	obj := &types.Object{
		ClassName:  "Exception",
		Properties: make(map[string]*types.Property),
	}

	obj.Properties["message"] = &types.Property{
		Value:      types.NewString("Test error"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["code"] = &types.Property{
		Value:      types.NewInt(123),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["file"] = &types.Property{
		Value:      types.NewString("test.php"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["line"] = &types.Property{
		Value:      types.NewInt(42),
		Visibility: types.VisibilityPublic,
	}

	// Convert to exception
	exc, err := FromObject(obj)
	if err != nil {
		t.Fatalf("FromObject failed: %v", err)
	}

	if exc.Message != "Test error" {
		t.Error("Message not extracted correctly")
	}

	if exc.Code != 123 {
		t.Error("Code not extracted correctly")
	}

	if exc.File != "test.php" {
		t.Error("File not extracted correctly")
	}

	if exc.Line != 42 {
		t.Error("Line not extracted correctly")
	}
}

// TestExceptionWithPreviousToObject tests chained exceptions to object
func TestExceptionWithPreviousToObject(t *testing.T) {
	previous := NewException("Previous error", 100, nil)
	exc := NewException("Current error", 200, previous)

	obj := exc.ToObject()

	// Check previous property
	if prevProp, ok := obj.Properties["previous"]; ok {
		if prevProp.Value.Type() != types.TypeObject {
			t.Error("Previous should be an object")
		}

		prevObj := prevProp.Value.ToObject()
		if msgProp, ok := prevObj.Properties["message"]; ok {
			if msgProp.Value.ToString() != "Previous error" {
				t.Error("Previous exception message incorrect")
			}
		} else {
			t.Error("Previous exception message not found")
		}
	} else {
		t.Error("Previous property not found")
	}
}

// TestExceptionString tests exception string representation
func TestExceptionString(t *testing.T) {
	exc := NewException("Test error", 123, nil)
	exc.File = "test.php"
	exc.Line = 42

	str := exc.String()

	if !strings.Contains(str, "Exception: Test error") {
		t.Error("String should contain exception class and message")
	}

	if !strings.Contains(str, "test.php:42") {
		t.Error("String should contain file and line")
	}
}

// TestExceptionGetTraceAsString tests stack trace string
func TestExceptionGetTraceAsString(t *testing.T) {
	exc := NewException("Test error", 123, nil)

	// Add some stack frames
	exc.Trace.AddFrame(&StackFrame{
		File:     "test.php",
		Line:     10,
		Function: "testFunc",
	})
	exc.Trace.AddFrame(&StackFrame{
		File:     "main.php",
		Line:     20,
		Function: "main",
	})

	trace := exc.GetTraceAsString()

	if !strings.Contains(trace, "test.php") {
		t.Error("Trace should contain file name")
	}

	if !strings.Contains(trace, "testFunc") {
		t.Error("Trace should contain function name")
	}
}

// TestErrorException tests ErrorException
func TestErrorException(t *testing.T) {
	exc := NewErrorException("Error", 123, E_WARNING, nil)

	if exc.ClassName != "ErrorException" {
		t.Errorf("Expected className 'ErrorException', got '%s'", exc.ClassName)
	}

	if exc.GetSeverity() != E_WARNING {
		t.Error("Severity not set correctly")
	}
}

// TestLogicException tests LogicException
func TestLogicException(t *testing.T) {
	exc := NewLogicException("Logic error", 123, nil)

	if exc.ClassName != "LogicException" {
		t.Errorf("Expected className 'LogicException', got '%s'", exc.ClassName)
	}
}

// TestRuntimeException tests RuntimeException
func TestRuntimeException(t *testing.T) {
	exc := NewRuntimeException("Runtime error", 123, nil)

	if exc.ClassName != "RuntimeException" {
		t.Errorf("Expected className 'RuntimeException', got '%s'", exc.ClassName)
	}
}

// TestInvalidArgumentException tests InvalidArgumentException
func TestInvalidArgumentException(t *testing.T) {
	exc := NewInvalidArgumentException("Invalid argument", 123, nil)

	if exc.ClassName != "InvalidArgumentException" {
		t.Errorf("Expected className 'InvalidArgumentException', got '%s'", exc.ClassName)
	}
}

// TestOutOfBoundsException tests OutOfBoundsException
func TestOutOfBoundsException(t *testing.T) {
	exc := NewOutOfBoundsException("Out of bounds", 123, nil)

	if exc.ClassName != "OutOfBoundsException" {
		t.Errorf("Expected className 'OutOfBoundsException', got '%s'", exc.ClassName)
	}
}

// TestOutOfRangeException tests OutOfRangeException
func TestOutOfRangeException(t *testing.T) {
	exc := NewOutOfRangeException("Out of range", 123, nil)

	if exc.ClassName != "OutOfRangeException" {
		t.Errorf("Expected className 'OutOfRangeException', got '%s'", exc.ClassName)
	}
}

// TestIsException tests exception detection
func TestIsException(t *testing.T) {
	// Create exception object
	exc := NewException("Test", 0, nil)
	obj := exc.ToObject()
	val := types.NewObject(obj)

	if !IsException(val) {
		t.Error("Should detect exception object")
	}

	// Non-exception
	nonExc := types.NewInt(42)
	if IsException(nonExc) {
		t.Error("Should not detect non-exception as exception")
	}
}

// TestGetExceptionFromValue tests extracting exception from value
func TestGetExceptionFromValue(t *testing.T) {
	exc := NewException("Test error", 123, nil)
	val := WrapExceptionAsValue(exc)

	extracted, err := GetExceptionFromValue(val)
	if err != nil {
		t.Fatalf("GetExceptionFromValue failed: %v", err)
	}

	if extracted.Message != "Test error" {
		t.Error("Extracted exception message incorrect")
	}

	if extracted.Code != 123 {
		t.Error("Extracted exception code incorrect")
	}
}

// TestGetExceptionFromValueError tests error case
func TestGetExceptionFromValueError(t *testing.T) {
	val := types.NewInt(42)

	_, err := GetExceptionFromValue(val)
	if err == nil {
		t.Error("Should return error for non-exception value")
	}
}

// TestWrapExceptionAsValue tests wrapping exception as value
func TestWrapExceptionAsValue(t *testing.T) {
	exc := NewException("Test", 123, nil)
	val := WrapExceptionAsValue(exc)

	if val.Type() != types.TypeResource {
		t.Error("Exception should be wrapped as resource")
	}

	// Should be able to extract it back
	extracted, err := GetExceptionFromValue(val)
	if err != nil {
		t.Error("Should be able to extract exception from wrapped value")
	}

	if extracted != exc {
		t.Error("Extracted exception should be same as original")
	}
}

// TestExceptionWithStackTrace tests exception with provided stack trace
func TestExceptionWithStackTrace(t *testing.T) {
	trace := NewStackTrace()
	trace.AddFrame(&StackFrame{
		File:     "test.php",
		Line:     10,
		Function: "testFunc",
	})

	exc := NewExceptionWithTrace("Error", 123, nil, trace)

	if exc.Trace != trace {
		t.Error("Stack trace not set correctly")
	}

	if len(exc.Trace.Frames) != 1 {
		t.Error("Stack trace should have 1 frame")
	}
}

// TestFromObjectWithPrevious tests converting object with previous exception
func TestFromObjectWithPrevious(t *testing.T) {
	// Create previous exception object
	prevExc := NewException("Previous", 100, nil)
	prevObj := prevExc.ToObject()

	// Create current exception object with previous
	obj := &types.Object{
		ClassName:  "Exception",
		Properties: make(map[string]*types.Property),
	}
	obj.Properties["message"] = &types.Property{
		Value:      types.NewString("Current"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["code"] = &types.Property{
		Value:      types.NewInt(200),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["previous"] = &types.Property{
		Value:      types.NewObject(prevObj),
		Visibility: types.VisibilityPublic,
	}

	// Convert to exception
	exc, err := FromObject(obj)
	if err != nil {
		t.Fatalf("FromObject failed: %v", err)
	}

	if exc.Previous == nil {
		t.Fatal("Previous exception should be set")
	}

	if exc.Previous.Message != "Previous" {
		t.Error("Previous exception message incorrect")
	}

	if exc.Previous.Code != 100 {
		t.Error("Previous exception code incorrect")
	}
}

// TestFromObjectNullPrevious tests converting object with null previous
func TestFromObjectNullPrevious(t *testing.T) {
	obj := &types.Object{
		ClassName:  "Exception",
		Properties: make(map[string]*types.Property),
	}
	obj.Properties["message"] = &types.Property{
		Value:      types.NewString("Test"),
		Visibility: types.VisibilityPublic,
	}
	obj.Properties["previous"] = &types.Property{
		Value:      types.NewNull(),
		Visibility: types.VisibilityPublic,
	}

	exc, err := FromObject(obj)
	if err != nil {
		t.Fatalf("FromObject failed: %v", err)
	}

	if exc.Previous != nil {
		t.Error("Previous should be nil for null value")
	}
}

// TestFromObjectNil tests error handling for nil object
func TestFromObjectNil(t *testing.T) {
	_, err := FromObject(nil)
	if err == nil {
		t.Error("Should return error for nil object")
	}
}
