package goext

import (
	"errors"
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Test FFI Manager creation

func TestNewFFIManager(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)

	if ffi == nil {
		t.Fatal("NewFFIManager returned nil")
	}

	if ffi.registry != reg {
		t.Error("FFI manager has wrong registry")
	}
}

func TestGetGlobalFFI(t *testing.T) {
	ffi1 := GetGlobalFFI()
	ffi2 := GetGlobalFFI()

	if ffi1.registry != ffi2.registry {
		t.Error("GetGlobalFFI returned different registries")
	}
}

// Test GoCall function

func TestGoCallBasic(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	// Register a simple function
	add := func(a, b int64) int64 {
		return a + b
	}

	err := reg.Register("math.Add", add)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Call via FFI
	args := []*types.Value{
		types.NewString("math.Add"),
		types.NewInt(10),
		types.NewInt(20),
	}

	result, err := ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}

	if result.Type() != types.TypeInt {
		t.Errorf("Expected int result, got %s", result.TypeString())
	}

	if result.ToInt() != 30 {
		t.Errorf("Expected 30, got %d", result.ToInt())
	}
}

func TestGoCallNoArgs(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	called := false
	fn := func() {
		called = true
	}

	err := reg.Register("test.NoArgs", fn)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Call with just function name, no additional args
	args := []*types.Value{
		types.NewString("test.NoArgs"),
	}

	result, err := ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}

	if !called {
		t.Error("Function was not called")
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null result, got %s", result.TypeString())
	}
}

func TestGoCallWithReturn(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	greet := func(name string) string {
		return "Hello, " + name + "!"
	}

	err := reg.Register("string.Greet", greet)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	args := []*types.Value{
		types.NewString("string.Greet"),
		types.NewString("World"),
	}

	result, err := ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}

	expected := "Hello, World!"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}
}

func TestGoCallWithError(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	divide := func(a, b int64) (int64, error) {
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}

	err := reg.Register("math.Divide", divide)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Test success case
	args := []*types.Value{
		types.NewString("math.Divide"),
		types.NewInt(10),
		types.NewInt(2),
	}

	result, err := ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}

	if result.ToInt() != 5 {
		t.Errorf("Expected 5, got %d", result.ToInt())
	}

	// Test error case
	args = []*types.Value{
		types.NewString("math.Divide"),
		types.NewInt(10),
		types.NewInt(0),
	}

	_, err = ffi.GoCall(testVM, args)
	if err == nil {
		t.Error("Expected error for division by zero")
	}

	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("Expected 'division by zero' in error, got: %v", err)
	}
}

func TestGoCallVariadic(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	sum := func(nums ...int64) int64 {
		total := int64(0)
		for _, n := range nums {
			total += n
		}
		return total
	}

	err := reg.Register("math.Sum", sum)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Test with multiple args
	args := []*types.Value{
		types.NewString("math.Sum"),
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
		types.NewInt(4),
		types.NewInt(5),
	}

	result, err := ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}

	if result.ToInt() != 15 {
		t.Errorf("Expected 15, got %d", result.ToInt())
	}

	// Test with no variadic args
	args = []*types.Value{
		types.NewString("math.Sum"),
	}

	result, err = ffi.GoCall(testVM, args)
	if err != nil {
		t.Fatalf("GoCall failed with no args: %v", err)
	}

	if result.ToInt() != 0 {
		t.Errorf("Expected 0, got %d", result.ToInt())
	}
}

// Test error cases

func TestGoCallNoArguments(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	_, err := ffi.GoCall(testVM, []*types.Value{})
	if err == nil {
		t.Error("Expected error when calling go_call with no arguments")
	}

	if !strings.Contains(err.Error(), "at least 1 argument") {
		t.Errorf("Expected 'at least 1 argument' error, got: %v", err)
	}
}

func TestGoCallNonStringFunctionName(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	args := []*types.Value{
		types.NewInt(123), // Not a string
	}

	_, err := ffi.GoCall(testVM, args)
	if err == nil {
		t.Error("Expected error when function name is not a string")
	}

	if !strings.Contains(err.Error(), "string (function name)") {
		t.Errorf("Expected 'string (function name)' error, got: %v", err)
	}
}

func TestGoCallFunctionNotFound(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	args := []*types.Value{
		types.NewString("nonexistent.Function"),
	}

	_, err := ffi.GoCall(testVM, args)
	if err == nil {
		t.Error("Expected error when function is not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestGoCallWrongArguments(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	add := func(a, b int64) int64 {
		return a + b
	}

	err := reg.Register("math.Add", add)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Too few arguments
	args := []*types.Value{
		types.NewString("math.Add"),
		types.NewInt(10),
		// Missing second argument
	}

	_, err = ffi.GoCall(testVM, args)
	if err == nil {
		t.Error("Expected error for too few arguments")
	}
}

// Test helper functions

func TestListFunctions(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	// Register some functions
	reg.Register("func1", func() {})
	reg.Register("func2", func() {})
	reg.Register("func3", func() {})

	result, err := ffi.ListFunctions(testVM, []*types.Value{})
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array result, got %s", result.TypeString())
	}

	arr := result.ToArray()
	if arr.Len() != 3 {
		t.Errorf("Expected 3 functions, got %d", arr.Len())
	}
}

func TestListFunctionsWithArgs(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	// Should fail with arguments
	args := []*types.Value{types.NewString("invalid")}
	_, err := ffi.ListFunctions(testVM, args)
	if err == nil {
		t.Error("Expected error when passing arguments to go_list_functions")
	}
}

func TestHasFunction(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	reg.Register("test.Func", func() {})

	// Test existing function
	args := []*types.Value{types.NewString("test.Func")}
	result, err := ffi.HasFunction(testVM, args)
	if err != nil {
		t.Fatalf("HasFunction failed: %v", err)
	}

	if result.Type() != types.TypeBool {
		t.Errorf("Expected bool result, got %s", result.TypeString())
	}

	if !result.ToBool() {
		t.Error("Expected true for existing function")
	}

	// Test non-existing function
	args = []*types.Value{types.NewString("nonexistent")}
	result, err = ffi.HasFunction(testVM, args)
	if err != nil {
		t.Fatalf("HasFunction failed: %v", err)
	}

	if result.ToBool() {
		t.Error("Expected false for non-existing function")
	}
}

func TestHasFunctionInvalidArgs(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	// No arguments
	_, err := ffi.HasFunction(testVM, []*types.Value{})
	if err == nil {
		t.Error("Expected error with no arguments")
	}

	// Non-string argument
	args := []*types.Value{types.NewInt(123)}
	_, err = ffi.HasFunction(testVM, args)
	if err == nil {
		t.Error("Expected error with non-string argument")
	}
}

func TestGetFunctionInfo(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	add := func(a, b int64) int64 {
		return a + b
	}

	reg.Register("math.Add", add)

	args := []*types.Value{types.NewString("math.Add")}
	result, err := ffi.GetFunctionInfo(testVM, args)
	if err != nil {
		t.Fatalf("GetFunctionInfo failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array result, got %s", result.TypeString())
	}

	info := result.ToArray()

	// Check name
	name, _ := info.Get(types.NewString("name"))
	if name.ToString() != "math.Add" {
		t.Errorf("Expected name 'math.Add', got '%s'", name.ToString())
	}

	// Check num_params
	numParams, _ := info.Get(types.NewString("num_params"))
	if numParams.ToInt() != 2 {
		t.Errorf("Expected num_params 2, got %d", numParams.ToInt())
	}

	// Check num_returns
	numReturns, _ := info.Get(types.NewString("num_returns"))
	if numReturns.ToInt() != 1 {
		t.Errorf("Expected num_returns 1, got %d", numReturns.ToInt())
	}

	// Check variadic
	variadic, _ := info.Get(types.NewString("variadic"))
	if variadic.ToBool() {
		t.Error("Expected variadic to be false")
	}
}

func TestGetFunctionInfoNonExistent(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	args := []*types.Value{types.NewString("nonexistent")}
	result, err := ffi.GetFunctionInfo(testVM, args)
	if err != nil {
		t.Fatalf("GetFunctionInfo failed: %v", err)
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null for non-existent function, got %s", result.TypeString())
	}
}

func TestGetFunctionInfoVariadic(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	sum := func(nums ...int64) int64 {
		return 0
	}

	reg.Register("math.Sum", sum)

	args := []*types.Value{types.NewString("math.Sum")}
	result, err := ffi.GetFunctionInfo(testVM, args)
	if err != nil {
		t.Fatalf("GetFunctionInfo failed: %v", err)
	}

	info := result.ToArray()

	// Check variadic flag
	variadic, _ := info.Get(types.NewString("variadic"))
	if !variadic.ToBool() {
		t.Error("Expected variadic to be true")
	}
}

// Test convenience functions

func TestCallGoFunction(t *testing.T) {
	// Clear global registry for clean test
	GetGlobalRegistry().Clear()

	testVM := vm.New()

	multiply := func(a, b int64) int64 {
		return a * b
	}

	err := RegisterFunction("math.Multiply", multiply)
	if err != nil {
		t.Fatalf("RegisterFunction failed: %v", err)
	}

	args := []*types.Value{
		types.NewInt(6),
		types.NewInt(7),
	}

	result, err := CallGoFunction("math.Multiply", testVM, args)
	if err != nil {
		t.Fatalf("CallGoFunction failed: %v", err)
	}

	if result.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", result.ToInt())
	}
}

func TestCallGoFunctionNotFound(t *testing.T) {
	testVM := vm.New()

	_, err := CallGoFunction("nonexistent", testVM, []*types.Value{})
	if err == nil {
		t.Error("Expected error for non-existent function")
	}

	if !strings.Contains(err.Error(), "not registered") {
		t.Errorf("Expected 'not registered' error, got: %v", err)
	}
}

func TestCallGoFunctionWithRegistry(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	subtract := func(a, b int64) int64 {
		return a - b
	}

	err := reg.Register("math.Subtract", subtract)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	args := []*types.Value{
		types.NewInt(10),
		types.NewInt(3),
	}

	result, err := CallGoFunctionWithRegistry("math.Subtract", testVM, args, reg)
	if err != nil {
		t.Fatalf("CallGoFunctionWithRegistry failed: %v", err)
	}

	if result.ToInt() != 7 {
		t.Errorf("Expected 7, got %d", result.ToInt())
	}
}

// Test integration scenarios

func TestFFIIntegrationComplete(t *testing.T) {
	reg := NewFunctionRegistry()
	ffi := NewFFIManager(reg)
	testVM := vm.New()

	// Register multiple functions
	reg.Register("string.Upper", func(s string) string {
		return strings.ToUpper(s)
	})

	reg.Register("string.Concat", func(a, b string) string {
		return a + b
	})

	reg.Register("math.Max", func(a, b int64) int64 {
		if a > b {
			return a
		}
		return b
	})

	// Test listing
	list, _ := ffi.ListFunctions(testVM, []*types.Value{})
	if list.ToArray().Len() != 3 {
		t.Errorf("Expected 3 functions, got %d", list.ToArray().Len())
	}

	// Test calling each function
	tests := []struct {
		name     string
		args     []*types.Value
		expected interface{}
	}{
		{
			"string.Upper",
			[]*types.Value{types.NewString("string.Upper"), types.NewString("hello")},
			"HELLO",
		},
		{
			"string.Concat",
			[]*types.Value{types.NewString("string.Concat"), types.NewString("Hello, "), types.NewString("World!")},
			"Hello, World!",
		},
		{
			"math.Max",
			[]*types.Value{types.NewString("math.Max"), types.NewInt(10), types.NewInt(20)},
			int64(20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ffi.GoCall(testVM, tt.args)
			if err != nil {
				t.Fatalf("GoCall failed: %v", err)
			}

			switch exp := tt.expected.(type) {
			case string:
				if result.ToString() != exp {
					t.Errorf("Expected '%s', got '%s'", exp, result.ToString())
				}
			case int64:
				if result.ToInt() != exp {
					t.Errorf("Expected %d, got %d", exp, result.ToInt())
				}
			}
		})
	}
}
