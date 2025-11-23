package goext

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Test basic function registration

func TestRegisterFunction(t *testing.T) {
	reg := NewFunctionRegistry()

	add := func(a, b int64) int64 {
		return a + b
	}

	err := reg.Register("math.Add", add)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !reg.Has("math.Add") {
		t.Error("Function not found after registration")
	}

	fn, ok := reg.Get("math.Add")
	if !ok {
		t.Fatal("Get returned false for registered function")
	}

	if fn.Name != "math.Add" {
		t.Errorf("Expected name 'math.Add', got '%s'", fn.Name)
	}
}

func TestRegisterFunctionDuplicate(t *testing.T) {
	reg := NewFunctionRegistry()

	fn := func() {}

	err := reg.Register("test.Func", fn)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	err = reg.Register("test.Func", fn)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}
}

func TestRegisterFunctionInvalid(t *testing.T) {
	reg := NewFunctionRegistry()

	// Nil function
	err := reg.Register("test.Nil", nil)
	if err == nil {
		t.Error("Expected error for nil function")
	}

	// Empty name
	err = reg.Register("", func() {})
	if err == nil {
		t.Error("Expected error for empty name")
	}

	// Not a function
	err = reg.Register("test.NotFunc", 42)
	if err == nil {
		t.Error("Expected error for non-function")
	}
}

// Test function signature parsing

func TestParseFunctionSignature(t *testing.T) {
	tests := []struct {
		name       string
		fn         interface{}
		expectNumIn  int
		expectNumOut int
		expectVariadic bool
	}{
		{"no params, no return", func() {}, 0, 0, false},
		{"one param, no return", func(int64) {}, 1, 0, false},
		{"two params, one return", func(int64, string) bool { return true }, 2, 1, false},
		{"variadic", func(args ...int64) {}, 1, 0, true},
		{"return with error", func() (string, error) { return "", nil }, 0, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := parseFunctionSignature(tt.name, tt.fn)
			if err != nil {
				t.Fatalf("parseFunctionSignature failed: %v", err)
			}

			if sig.NumIn != tt.expectNumIn {
				t.Errorf("NumIn: expected %d, got %d", tt.expectNumIn, sig.NumIn)
			}
			if sig.NumOut != tt.expectNumOut {
				t.Errorf("NumOut: expected %d, got %d", tt.expectNumOut, sig.NumOut)
			}
			if sig.IsVariadic != tt.expectVariadic {
				t.Errorf("IsVariadic: expected %v, got %v", tt.expectVariadic, sig.IsVariadic)
			}
		})
	}
}

// Test signature validation

func TestValidateSignature(t *testing.T) {
	validCases := []interface{}{
		func() {},
		func() int64 { return 0 },
		func() (int64, error) { return 0, nil },
		func(int64) {},
		func(int64, string) bool { return true },
		func(a, b int64) (int64, error) { return 0, nil },
	}

	for i, fn := range validCases {
		sig, err := parseFunctionSignature("test", fn)
		if err != nil {
			t.Fatalf("Case %d: parse failed: %v", i, err)
		}

		err = validateSignature(sig)
		if err != nil {
			t.Errorf("Case %d: expected valid, got error: %v", i, err)
		}
	}

	invalidCases := []interface{}{
		func() (int, int, int) { return 0, 0, 0 }, // Too many returns
		func() (int, string) { return 0, "" },      // Second return not error
	}

	for i, fn := range invalidCases {
		sig, err := parseFunctionSignature("test", fn)
		if err != nil {
			t.Fatalf("Invalid case %d: parse failed: %v", i, err)
		}

		err = validateSignature(sig)
		if err == nil {
			t.Errorf("Invalid case %d: expected error, got nil", i)
		}
	}
}

// Test function invocation

func TestInvokeSimpleFunction(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	add := func(a, b int64) int64 {
		return a + b
	}

	err := reg.Register("math.Add", add)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("math.Add")

	args := []*types.Value{
		types.NewInt(10),
		types.NewInt(20),
	}

	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.Type() != types.TypeInt {
		t.Errorf("Expected int result, got %s", result.TypeString())
	}

	if result.ToInt() != 30 {
		t.Errorf("Expected 30, got %d", result.ToInt())
	}
}

func TestInvokeNoReturn(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	called := false
	fn := func() {
		called = true
	}

	err := reg.Register("test.NoReturn", fn)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	regFn, _ := reg.Get("test.NoReturn")
	result, err := regFn.Handler(testVM, []*types.Value{})
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if !called {
		t.Error("Function was not called")
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null result, got %s", result.TypeString())
	}
}

func TestInvokeStringFunction(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	upper := func(s string) string {
		return strings.ToUpper(s)
	}

	err := reg.Register("string.Upper", upper)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("string.Upper")

	args := []*types.Value{types.NewString("hello")}
	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.Type() != types.TypeString {
		t.Errorf("Expected string result, got %s", result.TypeString())
	}

	if result.ToString() != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%s'", result.ToString())
	}
}

func TestInvokeWithError(t *testing.T) {
	reg := NewFunctionRegistry()
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

	fn, _ := reg.Get("math.Divide")

	// Test success case
	args := []*types.Value{types.NewInt(10), types.NewInt(2)}
	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}
	if result.ToInt() != 5 {
		t.Errorf("Expected 5, got %d", result.ToInt())
	}

	// Test error case
	args = []*types.Value{types.NewInt(10), types.NewInt(0)}
	_, err = fn.Handler(testVM, args)
	if err == nil {
		t.Error("Expected error for division by zero")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("Expected 'division by zero' error, got: %v", err)
	}
}

func TestInvokeWrongArgCount(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	add := func(a, b int64) int64 {
		return a + b
	}

	err := reg.Register("math.Add", add)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("math.Add")

	// Too few arguments
	args := []*types.Value{types.NewInt(10)}
	_, err = fn.Handler(testVM, args)
	if err == nil {
		t.Error("Expected error for too few arguments")
	}

	// Too many arguments
	args = []*types.Value{types.NewInt(10), types.NewInt(20), types.NewInt(30)}
	_, err = fn.Handler(testVM, args)
	if err == nil {
		t.Error("Expected error for too many arguments")
	}
}

func TestInvokeMultipleTypes(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	concat := func(s string, n int64, b bool) string {
		return s + ":" + string(rune(n)) + ":" + string(rune(map[bool]byte{true: 'T', false: 'F'}[b]))
	}

	err := reg.Register("test.Concat", concat)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("test.Concat")

	args := []*types.Value{
		types.NewString("test"),
		types.NewInt(65), // 'A'
		types.NewBool(true),
	}

	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	expected := "test:A:T"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}
}

func TestInvokeVariadic(t *testing.T) {
	reg := NewFunctionRegistry()
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

	fn, _ := reg.Get("math.Sum")

	// No arguments
	result, err := fn.Handler(testVM, []*types.Value{})
	if err != nil {
		t.Fatalf("Handler failed with no args: %v", err)
	}
	if result.ToInt() != 0 {
		t.Errorf("Expected 0, got %d", result.ToInt())
	}

	// Multiple arguments
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
		types.NewInt(4),
		types.NewInt(5),
	}

	result, err = fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.ToInt() != 15 {
		t.Errorf("Expected 15, got %d", result.ToInt())
	}
}

// Test registry operations

func TestRegistryList(t *testing.T) {
	reg := NewFunctionRegistry()

	reg.Register("func1", func() {})
	reg.Register("func2", func() {})
	reg.Register("func3", func() {})

	names := reg.List()
	if len(names) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(names))
	}

	// Check all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"func1", "func2", "func3"} {
		if !nameMap[expected] {
			t.Errorf("Expected '%s' in list", expected)
		}
	}
}

func TestRegistryCount(t *testing.T) {
	reg := NewFunctionRegistry()

	if reg.Count() != 0 {
		t.Errorf("Expected 0 functions, got %d", reg.Count())
	}

	reg.Register("func1", func() {})
	if reg.Count() != 1 {
		t.Errorf("Expected 1 function, got %d", reg.Count())
	}

	reg.Register("func2", func() {})
	reg.Register("func3", func() {})
	if reg.Count() != 3 {
		t.Errorf("Expected 3 functions, got %d", reg.Count())
	}
}

func TestRegistryClear(t *testing.T) {
	reg := NewFunctionRegistry()

	reg.Register("func1", func() {})
	reg.Register("func2", func() {})

	if reg.Count() != 2 {
		t.Fatalf("Expected 2 functions before clear")
	}

	reg.Clear()

	if reg.Count() != 0 {
		t.Errorf("Expected 0 functions after clear, got %d", reg.Count())
	}

	if reg.Has("func1") {
		t.Error("func1 still exists after clear")
	}
}

func TestGlobalRegistry(t *testing.T) {
	// Get global registry
	reg1 := GetGlobalRegistry()
	reg2 := GetGlobalRegistry()

	// Should be same instance
	if reg1 != reg2 {
		t.Error("GetGlobalRegistry returned different instances")
	}

	// Clear for clean test
	reg1.Clear()

	// Register via global function
	err := RegisterFunction("global.Test", func() {})
	if err != nil {
		t.Fatalf("RegisterFunction failed: %v", err)
	}

	if !reg1.Has("global.Test") {
		t.Error("Function not in global registry")
	}
}

// Test complex type marshaling

func TestInvokeSliceArg(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	sumSlice := func(nums []interface{}) int64 {
		total := int64(0)
		for _, n := range nums {
			if val, ok := n.(int64); ok {
				total += val
			}
		}
		return total
	}

	err := reg.Register("test.SumSlice", sumSlice)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("test.SumSlice")

	// Create PHP array [1, 2, 3, 4, 5]
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))
	arr.Append(types.NewInt(4))
	arr.Append(types.NewInt(5))

	args := []*types.Value{types.NewArray(arr)}
	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.ToInt() != 15 {
		t.Errorf("Expected 15, got %d", result.ToInt())
	}
}

func TestInvokeMapArg(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	getKey := func(m map[string]interface{}, key string) string {
		if val, ok := m[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	err := reg.Register("test.GetKey", getKey)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("test.GetKey")

	// Create PHP array ['name' => 'John', 'age' => '30']
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewString("30"))

	args := []*types.Value{
		types.NewArray(arr),
		types.NewString("name"),
	}

	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.ToString() != "John" {
		t.Errorf("Expected 'John', got '%s'", result.ToString())
	}
}

func TestInvokeMapReturn(t *testing.T) {
	reg := NewFunctionRegistry()
	testVM := vm.New()

	makeMap := func(key, value string) map[string]interface{} {
		return map[string]interface{}{
			key: value,
		}
	}

	err := reg.Register("test.MakeMap", makeMap)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	fn, _ := reg.Get("test.MakeMap")

	args := []*types.Value{
		types.NewString("greeting"),
		types.NewString("hello"),
	}

	result, err := fn.Handler(testVM, args)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array result, got %s", result.TypeString())
	}

	arr := result.ToArray()
	val, ok := arr.Get(types.NewString("greeting"))
	if !ok {
		t.Fatal("Key 'greeting' not found in result")
	}

	if val.ToString() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", val.ToString())
	}
}

// Test marshallable type checking

func TestIsMarshallableType(t *testing.T) {
	validTypes := []interface{}{
		true,
		int(0),
		int64(0),
		float64(0),
		"string",
		[]int{},
		[]interface{}{},
		map[string]interface{}{},
		struct{ Name string }{},
	}

	for _, val := range validTypes {
		t.Run(typeof(val), func(t *testing.T) {
			typ := reflect.TypeOf(val)
			if !isMarshallableType(typ) {
				t.Errorf("Type %s should be marshallable", typ)
			}
		})
	}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
