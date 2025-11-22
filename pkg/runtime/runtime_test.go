package runtime

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Runtime Creation Tests
// ============================================================================

func TestNew(t *testing.T) {
	rt := New()

	if rt == nil {
		t.Fatal("New() returned nil")
	}

	// Check superglobals are initialized
	if rt.GET == nil {
		t.Error("$_GET not initialized")
	}
	if rt.POST == nil {
		t.Error("$_POST not initialized")
	}
	if rt.SERVER == nil {
		t.Error("$_SERVER not initialized")
	}

	// Check constants map
	if rt.constants == nil {
		t.Error("constants not initialized")
	}

	// Check built-in constants exist
	if !rt.ConstantExists("PHP_VERSION") {
		t.Error("PHP_VERSION constant not defined")
	}
	if !rt.ConstantExists("TRUE") {
		t.Error("TRUE constant not defined")
	}
	if !rt.ConstantExists("FALSE") {
		t.Error("FALSE constant not defined")
	}
	if !rt.ConstantExists("NULL") {
		t.Error("NULL constant not defined")
	}
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestDefineConstant(t *testing.T) {
	rt := New()

	err := rt.DefineConstant("MY_CONST", types.NewInt(42))
	if err != nil {
		t.Errorf("DefineConstant() error: %v", err)
	}

	val, ok := rt.GetConstant("MY_CONST")
	if !ok {
		t.Error("Constant not found after definition")
	}

	if val.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", val.ToInt())
	}
}

func TestDefineConstant_Duplicate(t *testing.T) {
	rt := New()

	rt.DefineConstant("MY_CONST", types.NewInt(42))

	// Try to redefine
	err := rt.DefineConstant("MY_CONST", types.NewInt(99))
	if err == nil {
		t.Error("Expected error when redefining constant")
	}
}

func TestGetConstant_NotExists(t *testing.T) {
	rt := New()

	_, ok := rt.GetConstant("NONEXISTENT")
	if ok {
		t.Error("GetConstant() returned true for non-existent constant")
	}
}

func TestBuiltinConstants(t *testing.T) {
	rt := New()

	tests := []struct {
		name     string
		expected types.ValueType
	}{
		{"PHP_VERSION", types.TypeString},
		{"PHP_MAJOR_VERSION", types.TypeInt},
		{"PHP_MINOR_VERSION", types.TypeInt},
		{"TRUE", types.TypeBool},
		{"FALSE", types.TypeBool},
		{"NULL", types.TypeNull},
		{"PHP_INT_MAX", types.TypeInt},
		{"PHP_FLOAT_MAX", types.TypeFloat},
	}

	for _, tt := range tests {
		val, ok := rt.GetConstant(tt.name)
		if !ok {
			t.Errorf("Built-in constant '%s' not found", tt.name)
			continue
		}

		if val.Type() != tt.expected {
			t.Errorf("Constant '%s': expected type %v, got %v",
				tt.name, tt.expected, val.Type())
		}
	}
}

// ============================================================================
// Superglobals Tests
// ============================================================================

func TestGetSuperglobal(t *testing.T) {
	rt := New()

	tests := []string{
		"_GET", "_POST", "_REQUEST", "_SERVER",
		"_ENV", "_COOKIE", "_FILES", "_SESSION", "GLOBALS",
	}

	for _, name := range tests {
		val, ok := rt.GetSuperglobal(name)
		if !ok {
			t.Errorf("GetSuperglobal('%s') returned false", name)
			continue
		}

		if !val.IsArray() {
			t.Errorf("Superglobal '%s' is not an array", name)
		}
	}
}

func TestGetSuperglobal_NotExists(t *testing.T) {
	rt := New()

	_, ok := rt.GetSuperglobal("_INVALID")
	if ok {
		t.Error("GetSuperglobal() returned true for invalid superglobal")
	}
}

func TestServerSuperglobal(t *testing.T) {
	rt := New()

	server := rt.SERVER.ToArray()

	// Check for expected keys
	keys := []string{"SERVER_NAME", "SERVER_SOFTWARE", "REQUEST_METHOD"}

	for _, key := range keys {
		val, exists := server.Get(types.NewString(key))
		if !exists {
			t.Errorf("$_SERVER['%s'] not set", key)
			continue
		}

		if val.IsNull() {
			t.Errorf("$_SERVER['%s'] is null", key)
		}
	}
}

func TestSetScriptPath(t *testing.T) {
	rt := New()

	rt.SetScriptPath("/path/to/script.php")

	if rt.scriptPath != "/path/to/script.php" {
		t.Error("scriptPath not set correctly")
	}

	// Check $_SERVER was updated
	server := rt.SERVER.ToArray()
	filename, exists := server.Get(types.NewString("SCRIPT_FILENAME"))
	if !exists {
		t.Error("SCRIPT_FILENAME not set in $_SERVER")
	}

	if filename.ToString() != "/path/to/script.php" {
		t.Errorf("Expected '/path/to/script.php', got '%s'", filename.ToString())
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestErrorReporting(t *testing.T) {
	rt := New()

	// Default should be E_ALL
	if rt.GetErrorReporting() != int(E_ALL) {
		t.Errorf("Expected default error reporting %d, got %d",
			E_ALL, rt.GetErrorReporting())
	}

	// Set custom level
	rt.SetErrorReporting(int(E_ERROR | E_WARNING))

	if rt.GetErrorReporting() != int(E_ERROR|E_WARNING) {
		t.Error("Error reporting level not set correctly")
	}
}

func TestSetErrorHandler(t *testing.T) {
	rt := New()

	called := false
	handler := func(errorType ErrorType, message string, file string, line int) {
		called = true
	}

	rt.SetErrorHandler(handler)

	// Trigger an error
	rt.TriggerError(E_WARNING, "test warning", "test.php", 10)

	if !called {
		t.Error("Custom error handler was not called")
	}
}

func TestTriggerError_Filtered(t *testing.T) {
	rt := New()

	// Only report errors
	rt.SetErrorReporting(int(E_ERROR))

	called := false
	handler := func(errorType ErrorType, message string, file string, line int) {
		called = true
	}

	rt.SetErrorHandler(handler)

	// Trigger a warning (should be filtered)
	rt.TriggerError(E_WARNING, "test warning", "test.php", 10)

	if called {
		t.Error("Filtered error should not trigger handler")
	}

	// Trigger an error (should not be filtered)
	rt.TriggerError(E_ERROR, "test error", "test.php", 10)

	if !called {
		t.Error("Non-filtered error should trigger handler")
	}
}

// ============================================================================
// Output Buffering Tests
// ============================================================================

func TestOutputBuffering(t *testing.T) {
	rt := New()

	// Start buffering
	rt.StartOutputBuffering()

	if rt.GetOutputBufferLevel() != 1 {
		t.Errorf("Expected buffer level 1, got %d", rt.GetOutputBufferLevel())
	}

	// Write some output
	rt.Write("Hello")
	rt.Write(" ")
	rt.Write("World")

	// Get contents
	contents := rt.GetOutputBufferContents()
	if contents != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", contents)
	}

	// End buffering
	final := rt.EndOutputBuffering()
	if final != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", final)
	}

	if rt.GetOutputBufferLevel() != 0 {
		t.Errorf("Expected buffer level 0 after end, got %d", rt.GetOutputBufferLevel())
	}
}

func TestOutputBuffering_Nested(t *testing.T) {
	rt := New()

	// Start first buffer
	rt.StartOutputBuffering()
	rt.Write("Level 1")

	// Start second buffer
	rt.StartOutputBuffering()
	rt.Write("Level 2")

	if rt.GetOutputBufferLevel() != 2 {
		t.Errorf("Expected buffer level 2, got %d", rt.GetOutputBufferLevel())
	}

	// End second buffer
	level2 := rt.EndOutputBuffering()
	if level2 != "Level 2" {
		t.Errorf("Expected 'Level 2', got '%s'", level2)
	}

	// End first buffer
	level1 := rt.EndOutputBuffering()
	if level1 != "Level 1" {
		t.Errorf("Expected 'Level 1', got '%s'", level1)
	}
}

func TestCleanOutputBuffer(t *testing.T) {
	rt := New()

	rt.StartOutputBuffering()
	rt.Write("This will be discarded")

	rt.CleanOutputBuffer()

	if rt.GetOutputBufferLevel() != 0 {
		t.Error("Buffer should be removed after clean")
	}
}

func TestFlushOutputBuffer(t *testing.T) {
	rt := New()

	rt.StartOutputBuffering()
	rt.Write("Test")

	// Flush should return contents and clear
	contents := rt.FlushOutputBuffer()
	if contents != "Test" {
		t.Errorf("Expected 'Test', got '%s'", contents)
	}

	// Buffer should still be active but empty
	if rt.GetOutputBufferLevel() != 1 {
		t.Error("Buffer should still be active after flush")
	}

	contents = rt.GetOutputBufferContents()
	if contents != "" {
		t.Errorf("Expected empty buffer after flush, got '%s'", contents)
	}
}

func TestWrite_NoBuffer(t *testing.T) {
	rt := New()

	// Writing without buffer should not panic
	// (it will print to stdout, but we can't capture that easily in test)
	rt.Write("test")
}

// ============================================================================
// Error Type Tests
// ============================================================================

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{E_ERROR, "Fatal error"},
		{E_WARNING, "Warning"},
		{E_PARSE, "Parse error"},
		{E_NOTICE, "Notice"},
		{E_CORE_ERROR, "Core error"},
		{E_CORE_WARNING, "Core warning"},
		{E_COMPILE_ERROR, "Compile error"},
		{E_COMPILE_WARNING, "Compile warning"},
		{E_USER_ERROR, "User error"},
		{E_USER_WARNING, "User warning"},
		{E_USER_NOTICE, "User notice"},
		{E_STRICT, "Strict standards"},
		{E_RECOVERABLE_ERROR, "Recoverable error"},
		{E_DEPRECATED, "Deprecated"},
		{E_USER_DEPRECATED, "User deprecated"},
		{ErrorType(9999), "Unknown error"},
	}

	for _, tt := range tests {
		result := tt.errorType.String()
		if result != tt.expected {
			t.Errorf("ErrorType(%d).String(): expected '%s', got '%s'",
				tt.errorType, tt.expected, result)
		}
	}
}

func TestTriggerError_DefaultHandler(t *testing.T) {
	rt := New()

	// Don't set custom handler, use default
	// Just ensure it doesn't panic
	rt.TriggerError(E_WARNING, "test warning", "test.php", 10)
}

// ============================================================================
// Stack Trace Tests
// ============================================================================

func TestStackTrace(t *testing.T) {
	st := NewStackTrace()

	if st == nil {
		t.Fatal("NewStackTrace() returned nil")
	}

	if len(st.Frames) != 0 {
		t.Error("New stack trace should have no frames")
	}

	// Add a frame
	frame := &StackFrame{
		File:     "test.php",
		Line:     10,
		Function: "foo",
		Class:    "Bar",
		Type:     "->",
	}

	st.AddFrame(frame)

	if len(st.Frames) != 1 {
		t.Errorf("Expected 1 frame, got %d", len(st.Frames))
	}

	// Get string representation
	str := st.String()
	if str == "" {
		t.Error("Stack trace string is empty")
	}
}

// ============================================================================
// Output Buffer Tests
// ============================================================================

func TestOutputBuffer_Len(t *testing.T) {
	buf := NewOutputBuffer()

	if buf.Len() != 0 {
		t.Errorf("Expected length 0, got %d", buf.Len())
	}

	buf.Write("Hello")

	if buf.Len() != 5 {
		t.Errorf("Expected length 5, got %d", buf.Len())
	}
}

func TestCleanOutputBuffer_NoBuffer(t *testing.T) {
	rt := New()

	// Should not panic when no buffer exists
	rt.CleanOutputBuffer()

	if rt.GetOutputBufferLevel() != 0 {
		t.Error("Buffer level should be 0")
	}
}

func TestGetOutputBufferContents_NoBuffer(t *testing.T) {
	rt := New()

	contents := rt.GetOutputBufferContents()
	if contents != "" {
		t.Errorf("Expected empty string, got '%s'", contents)
	}
}

func TestFlushOutputBuffer_NoBuffer(t *testing.T) {
	rt := New()

	contents := rt.FlushOutputBuffer()
	if contents != "" {
		t.Errorf("Expected empty string, got '%s'", contents)
	}
}

func TestEndOutputBuffering_NoBuffer(t *testing.T) {
	rt := New()

	contents := rt.EndOutputBuffering()
	if contents != "" {
		t.Errorf("Expected empty string, got '%s'", contents)
	}
}
