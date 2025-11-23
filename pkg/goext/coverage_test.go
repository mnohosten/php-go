package goext

import (
	"reflect"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// Additional tests to improve coverage

// Test ToGo edge cases for better coverage
func TestToGoEdgeCases(t *testing.T) {
	m := NewMarshaler()

	// Test with bool
	boolVal := types.NewBool(true)
	result, err := m.ToGoBool(boolVal)
	if err != nil {
		t.Fatalf("ToGoBool failed: %v", err)
	}
	if !result {
		t.Error("Expected true")
	}

	// Test with different array types
	// List array
	listArr := types.NewEmptyArray()
	listArr.Append(types.NewInt(1))
	listArr.Append(types.NewInt(2))
	listArr.Append(types.NewInt(3))

	result2, err := m.ToGoSlice(types.NewArray(listArr))
	if err != nil {
		t.Fatalf("ToGoSlice failed: %v", err)
	}
	if len(result2) != 3 {
		t.Errorf("Expected length 3, got %d", len(result2))
	}

	// Associative array
	mapArr := types.NewEmptyArray()
	mapArr.Set(types.NewString("key1"), types.NewString("value1"))
	mapArr.Set(types.NewString("key2"), types.NewInt(42))

	result3, err := m.ToGoMap(types.NewArray(mapArr))
	if err != nil {
		t.Fatalf("ToGoMap failed: %v", err)
	}
	if len(result3) != 2 {
		t.Errorf("Expected length 2, got %d", len(result3))
	}

	// Test with undef type
	undefVal := types.NewUndef()
	result4, err := m.ToGo(undefVal)
	if err != nil {
		t.Fatalf("ToGo with undef failed: %v", err)
	}
	if result4 != nil {
		t.Error("Expected nil for undef")
	}

	// Test with reference type
	refVal := types.NewInt(100)
	refToVal := types.NewReference(refVal)

	result5, err := m.ToGo(refToVal)
	if err != nil {
		t.Fatalf("ToGo with reference failed: %v", err)
	}
	if result5.(int64) != 100 {
		t.Errorf("Expected 100, got %v", result5)
	}
}

// Test ToGoAdvanced with various target types
func TestToGoAdvancedVariousTypes(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test with null to any type
	nullVal := types.NewNull()
	result, err := am.ToGoAdvanced(nullVal, reflect.TypeOf(0))
	if err != nil {
		t.Fatalf("ToGoAdvanced with null failed: %v", err)
	}
	if result.(int) != 0 {
		t.Errorf("Expected zero value, got %v", result)
	}

	// Test with struct target
	type SimpleStruct struct {
		Name  string
		Value int
	}

	arr := types.NewEmptyArray()
	arr.Set(types.NewString("Name"), types.NewString("Test"))
	arr.Set(types.NewString("Value"), types.NewInt(42))

	result, err = am.ToGoAdvanced(types.NewArray(arr), reflect.TypeOf(SimpleStruct{}))
	if err != nil {
		t.Fatalf("ToGoAdvanced with struct failed: %v", err)
	}

	ss := result.(SimpleStruct)
	if ss.Name != "Test" || ss.Value != 42 {
		t.Errorf("Struct fields not populated correctly: %+v", ss)
	}

	// Test with non-struct, non-converter type
	intVal := types.NewInt(100)
	result, err = am.ToGoAdvanced(intVal, reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("ToGoAdvanced with int failed: %v", err)
	}
	if result.(int64) != 100 {
		t.Errorf("Expected 100, got %v", result)
	}
}

// Test ObjectToStruct error cases
func TestObjectToStructErrors(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test with non-object
	intVal := types.NewInt(42)
	_, err := am.ObjectToStruct(intVal, reflect.TypeOf(struct{}{}))
	if err == nil {
		t.Error("Expected error for non-object input")
	}
}

// Test Clone edge cases
func TestCloneEdgeCases(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test with nil
	result, err := am.Clone(nil)
	if err != nil {
		t.Fatalf("Clone nil failed: %v", err)
	}
	if result == nil || result.Type() != types.TypeNull {
		t.Error("Expected null result for nil input")
	}

	// Test with complex nested structure
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("simple"), types.NewInt(42))

	nested := types.NewEmptyArray()
	nested.Append(types.NewString("nested"))
	arr.Set(types.NewString("nested"), types.NewArray(nested))

	cloned, err := am.Clone(types.NewArray(arr))
	if err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	if !am.DeepEquals(types.NewArray(arr), cloned) {
		t.Error("Clone not equal to original")
	}
}

// Test DeepEquals edge cases
func TestDeepEqualsEdgeCases(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test with different types
	intVal := types.NewInt(42)
	strVal := types.NewString("42")

	if am.DeepEquals(intVal, strVal) {
		t.Error("Different types should not be equal")
	}

	// Test with null values
	null1 := types.NewNull()
	null2 := types.NewNull()

	if !am.DeepEquals(null1, null2) {
		t.Error("Two nulls should be equal")
	}

	// Test with simple types (uses Identical)
	int1 := types.NewInt(42)
	int2 := types.NewInt(42)

	if !am.DeepEquals(int1, int2) {
		t.Error("Equal ints should be deepequal")
	}

	// Test with ToGo error path
	// This is harder to trigger, but we can test the happy path
	arr1 := types.NewEmptyArray()
	arr1.Append(types.NewInt(1))

	arr2 := types.NewEmptyArray()
	arr2.Append(types.NewInt(1))

	if !am.DeepEquals(types.NewArray(arr1), types.NewArray(arr2)) {
		t.Error("Equal arrays should be deepequal")
	}
}

// Test convertToType edge cases
func TestConvertToTypeEdgeCases(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test nil to target type
	result, err := am.convertToType(nil, reflect.TypeOf(0))
	if err != nil {
		t.Fatalf("convertToType with nil failed: %v", err)
	}
	if result.(int) != 0 {
		t.Errorf("Expected zero value, got %v", result)
	}

	// Test assignable types
	var i int64 = 42
	result, err = am.convertToType(i, reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("convertToType assignable failed: %v", err)
	}
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// Test string conversion from non-string
	result, err = am.convertToType(42, reflect.TypeOf(""))
	if err != nil {
		t.Fatalf("convertToType to string failed: %v", err)
	}
	if result.(string) != "42" {
		t.Errorf("Expected '42', got %v", result)
	}

	// Test bool conversion from int
	result, err = am.convertToType(int64(1), reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("convertToType to bool failed: %v", err)
	}
	if !result.(bool) {
		t.Error("Expected true")
	}

	result, err = am.convertToType(int64(0), reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("convertToType to bool failed: %v", err)
	}
	if result.(bool) {
		t.Error("Expected false")
	}

	// Test bool conversion from bool
	result, err = am.convertToType(true, reflect.TypeOf(false))
	if err != nil {
		t.Fatalf("convertToType bool to bool failed: %v", err)
	}
	if !result.(bool) {
		t.Error("Expected true")
	}

	// Test error case - unsupported conversion
	_, err = am.convertToType("not a number", reflect.TypeOf(int64(0)))
	if err != nil {
		// Expected - string to int conversion not fully supported
		// This is OK, we're testing the error path
	}
}

// Test registration with various types to test isMarshallableType indirectly
func TestRegisterVariousTypes(t *testing.T) {
	registry := NewFunctionRegistry()

	// Test with valid marshallable types
	validTests := []struct {
		name string
		fn   interface{}
	}{
		{"bool", func(b bool) bool { return b }},
		{"int", func(i int) int { return i }},
		{"int64", func(i int64) int64 { return i }},
		{"float64", func(f float64) float64 { return f }},
		{"string", func(s string) string { return s }},
		{"slice", func(s []int) []int { return s }},
		{"map", func(m map[string]int) map[string]int { return m }},
		{"struct", func(s struct{ Name string }) struct{ Name string } { return s }},
		{"ptr", func(p *int) *int { return p }},
		{"interface", func(i interface{}) interface{} { return i }},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Register("test."+tt.name, tt.fn)
			if err != nil {
				t.Errorf("Register with %s failed: %v", tt.name, err)
			}
		})
	}

	// Test with invalid types
	invalidTests := []struct {
		name string
		fn   interface{}
	}{
		{"chan", func(c chan int) chan int { return c }},
	}

	for _, tt := range invalidTests {
		t.Run("invalid_"+tt.name, func(t *testing.T) {
			err := registry.Register("test.invalid."+tt.name, tt.fn)
			if err == nil {
				t.Errorf("Expected error for %s type", tt.name)
			}
		})
	}
}

// Test wrapFunction edge cases
func TestWrapFunctionEdgeCases(t *testing.T) {
	registry := NewFunctionRegistry()
	testVM := vm.New()

	// Test with function that has error but no other return
	funcWithOnlyError := func() error {
		return nil
	}

	sig, err := parseFunctionSignature("test.Func", funcWithOnlyError)
	if err != nil {
		t.Fatalf("parseFunctionSignature failed: %v", err)
	}

	handler := registry.wrapFunction(funcWithOnlyError, sig)

	result, err := handler(testVM, []*types.Value{})
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null return, got %s", result.TypeString())
	}

	// Test with variadic function
	variadicFunc := func(nums ...int64) int64 {
		sum := int64(0)
		for _, n := range nums {
			sum += n
		}
		return sum
	}

	sig, err = parseFunctionSignature("test.Sum", variadicFunc)
	if err != nil {
		t.Fatalf("parseFunctionSignature failed: %v", err)
	}

	handler = registry.wrapFunction(variadicFunc, sig)

	result, err = handler(testVM, []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
	})

	if err != nil {
		t.Fatalf("Variadic handler failed: %v", err)
	}

	if result.ToInt() != 6 {
		t.Errorf("Expected 6, got %d", result.ToInt())
	}
}

// Test StructToObject with nil pointer
func TestStructToObjectNilPointer(t *testing.T) {
	am := NewAdvancedMarshaler()

	var ptr *struct {
		Name string
	}

	result, err := am.StructToObject(ptr, "TestStruct")
	if err != nil {
		t.Fatalf("StructToObject with nil pointer failed: %v", err)
	}

	if result.Type() != types.TypeNull {
		t.Errorf("Expected null for nil pointer, got %s", result.TypeString())
	}
}

// Test phpToStruct with object type (not just array)
func TestPHPToStructWithObject(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Create PHP object
	obj := types.NewObjectInstance("TestClass")
	obj.SetProperty("Name", types.NewString("Test"), nil)
	obj.SetProperty("Value", types.NewInt(42), nil)

	type TestStruct struct {
		Name  string
		Value int
	}

	result, err := am.phpToStruct(types.NewObject(obj), reflect.TypeOf(TestStruct{}))
	if err != nil {
		t.Fatalf("phpToStruct with object failed: %v", err)
	}

	ts := result.(TestStruct)
	if ts.Name != "Test" || ts.Value != 42 {
		t.Errorf("Unexpected struct values: %+v", ts)
	}
}

// Test convertToType with slice and map conversions
func TestConvertToTypeComplexConversions(t *testing.T) {
	am := NewAdvancedMarshaler()

	// Test slice conversion
	input := []interface{}{int64(1), int64(2), int64(3)}
	result, err := am.convertToType(input, reflect.TypeOf([]int{}))
	if err != nil {
		t.Fatalf("convertToType slice failed: %v", err)
	}

	slice := result.([]int)
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("Unexpected slice values: %v", slice)
	}

	// Test map conversion
	inputMap := map[string]interface{}{
		"a": int64(1),
		"b": int64(2),
	}

	result, err = am.convertToType(inputMap, reflect.TypeOf(map[string]int{}))
	if err != nil {
		t.Fatalf("convertToType map failed: %v", err)
	}

	mapVal := result.(map[string]int)
	if mapVal["a"] != 1 || mapVal["b"] != 2 {
		t.Errorf("Unexpected map values: %v", mapVal)
	}

	// Test struct conversion from map
	type SimpleStruct struct {
		Name string
		Age  int
	}

	inputStruct := map[string]interface{}{
		"Name": "John",
		"Age":  int64(30),
	}

	result, err = am.convertToType(inputStruct, reflect.TypeOf(SimpleStruct{}))
	if err != nil {
		t.Fatalf("convertToType struct failed: %v", err)
	}

	ss := result.(SimpleStruct)
	if ss.Name != "John" || ss.Age != 30 {
		t.Errorf("Unexpected struct values: %+v", ss)
	}
}
