package goext

import (
	"reflect"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// Test basic type conversions: PHP → Go

func TestToGoNull(t *testing.T) {
	m := NewMarshaler()

	// Test null
	phpVal := types.NewNull()
	goVal, err := m.ToGo(phpVal)
	if err != nil {
		t.Fatalf("ToGo(null) failed: %v", err)
	}
	if goVal != nil {
		t.Errorf("Expected nil, got %v", goVal)
	}

	// Test nil value
	goVal, err = m.ToGo(nil)
	if err != nil {
		t.Fatalf("ToGo(nil) failed: %v", err)
	}
	if goVal != nil {
		t.Errorf("Expected nil, got %v", goVal)
	}
}

func TestToGoBool(t *testing.T) {
	m := NewMarshaler()

	tests := []struct {
		input    bool
		expected bool
	}{
		{true, true},
		{false, false},
	}

	for _, tt := range tests {
		phpVal := types.NewBool(tt.input)
		goVal, err := m.ToGo(phpVal)
		if err != nil {
			t.Fatalf("ToGo(bool) failed: %v", err)
		}
		result, ok := goVal.(bool)
		if !ok {
			t.Errorf("Expected bool, got %T", goVal)
		}
		if result != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, result)
		}
	}
}

func TestToGoInt(t *testing.T) {
	m := NewMarshaler()

	tests := []int64{0, 1, -1, 42, 1234567890}

	for _, expected := range tests {
		phpVal := types.NewInt(expected)
		goVal, err := m.ToGo(phpVal)
		if err != nil {
			t.Fatalf("ToGo(int) failed: %v", err)
		}
		result, ok := goVal.(int64)
		if !ok {
			t.Errorf("Expected int64, got %T", goVal)
		}
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	}
}

func TestToGoFloat(t *testing.T) {
	m := NewMarshaler()

	tests := []float64{0.0, 1.5, -3.14, 123.456}

	for _, expected := range tests {
		phpVal := types.NewFloat(expected)
		goVal, err := m.ToGo(phpVal)
		if err != nil {
			t.Fatalf("ToGo(float) failed: %v", err)
		}
		result, ok := goVal.(float64)
		if !ok {
			t.Errorf("Expected float64, got %T", goVal)
		}
		if result != expected {
			t.Errorf("Expected %f, got %f", expected, result)
		}
	}
}

func TestToGoString(t *testing.T) {
	m := NewMarshaler()

	tests := []string{"", "hello", "with spaces", "unicode: 你好"}

	for _, expected := range tests {
		phpVal := types.NewString(expected)
		goVal, err := m.ToGo(phpVal)
		if err != nil {
			t.Fatalf("ToGo(string) failed: %v", err)
		}
		result, ok := goVal.(string)
		if !ok {
			t.Errorf("Expected string, got %T", goVal)
		}
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	}
}

func TestToGoArrayList(t *testing.T) {
	m := NewMarshaler()

	// Create PHP array [1, 2, 3]
	phpArr := types.NewEmptyArray()
	phpArr.Append(types.NewInt(1))
	phpArr.Append(types.NewInt(2))
	phpArr.Append(types.NewInt(3))

	goVal, err := m.ToGo(types.NewArray(phpArr))
	if err != nil {
		t.Fatalf("ToGo(array) failed: %v", err)
	}

	slice, ok := goVal.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", goVal)
	}

	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}

	expected := []int64{1, 2, 3}
	for i, exp := range expected {
		val, ok := slice[i].(int64)
		if !ok {
			t.Errorf("Element %d: expected int64, got %T", i, slice[i])
			continue
		}
		if val != exp {
			t.Errorf("Element %d: expected %d, got %d", i, exp, val)
		}
	}
}

func TestToGoArrayAssoc(t *testing.T) {
	m := NewMarshaler()

	// Create PHP array ['a' => 1, 'b' => 2]
	phpArr := types.NewEmptyArray()
	phpArr.Set(types.NewString("a"), types.NewInt(1))
	phpArr.Set(types.NewString("b"), types.NewInt(2))

	goVal, err := m.ToGo(types.NewArray(phpArr))
	if err != nil {
		t.Fatalf("ToGo(array) failed: %v", err)
	}

	m2, ok := goVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", goVal)
	}

	if len(m2) != 2 {
		t.Errorf("Expected length 2, got %d", len(m2))
	}

	if val, ok := m2["a"].(int64); !ok || val != 1 {
		t.Errorf("Expected m['a'] = 1, got %v", m2["a"])
	}
	if val, ok := m2["b"].(int64); !ok || val != 2 {
		t.Errorf("Expected m['b'] = 2, got %v", m2["b"])
	}
}

func TestToGoObject(t *testing.T) {
	m := NewMarshaler()

	// Create PHP object
	obj := &types.Object{
		ClassName:  "TestClass",
		Properties: make(map[string]*types.Property),
	}
	obj.Properties["name"] = &types.Property{Value: types.NewString("John")}
	obj.Properties["age"] = &types.Property{Value: types.NewInt(30)}

	phpVal := types.NewObject(obj)

	goVal, err := m.ToGo(phpVal)
	if err != nil {
		t.Fatalf("ToGo(object) failed: %v", err)
	}

	result, ok := goVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", goVal)
	}

	if name, ok := result["name"].(string); !ok || name != "John" {
		t.Errorf("Expected name='John', got %v", result["name"])
	}
	if age, ok := result["age"].(int64); !ok || age != 30 {
		t.Errorf("Expected age=30, got %v", result["age"])
	}
}

// Test typed conversion helpers

func TestToGoIntHelper(t *testing.T) {
	m := NewMarshaler()

	phpVal := types.NewInt(42)
	result, err := m.ToGoInt(phpVal)
	if err != nil {
		t.Fatalf("ToGoInt failed: %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test conversion from float
	phpVal = types.NewFloat(3.14)
	result, err = m.ToGoInt(phpVal)
	if err != nil {
		t.Fatalf("ToGoInt(float) failed: %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}

	// Test nil
	_, err = m.ToGoInt(nil)
	if err == nil {
		t.Error("Expected error for nil, got nil")
	}
}

func TestToGoFloatHelper(t *testing.T) {
	m := NewMarshaler()

	phpVal := types.NewFloat(3.14)
	result, err := m.ToGoFloat(phpVal)
	if err != nil {
		t.Fatalf("ToGoFloat failed: %v", err)
	}
	if result != 3.14 {
		t.Errorf("Expected 3.14, got %f", result)
	}

	// Test conversion from int
	phpVal = types.NewInt(42)
	result, err = m.ToGoFloat(phpVal)
	if err != nil {
		t.Fatalf("ToGoFloat(int) failed: %v", err)
	}
	if result != 42.0 {
		t.Errorf("Expected 42.0, got %f", result)
	}
}

func TestToGoBoolHelper(t *testing.T) {
	m := NewMarshaler()

	tests := []struct {
		input    *types.Value
		expected bool
	}{
		{types.NewBool(true), true},
		{types.NewBool(false), false},
		{types.NewInt(1), true},
		{types.NewInt(0), false},
		{types.NewString("hello"), true},
		{types.NewString(""), false},
		{nil, false},
	}

	for _, tt := range tests {
		result, err := m.ToGoBool(tt.input)
		if err != nil {
			t.Errorf("ToGoBool(%v) failed: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("ToGoBool(%v): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestToGoStringHelper(t *testing.T) {
	m := NewMarshaler()

	phpVal := types.NewString("hello")
	result, err := m.ToGoString(phpVal)
	if err != nil {
		t.Fatalf("ToGoString failed: %v", err)
	}
	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}

	// Test conversion from int
	phpVal = types.NewInt(42)
	result, err = m.ToGoString(phpVal)
	if err != nil {
		t.Fatalf("ToGoString(int) failed: %v", err)
	}
	if result != "42" {
		t.Errorf("Expected '42', got %q", result)
	}
}

func TestToGoSliceHelper(t *testing.T) {
	m := NewMarshaler()

	phpArr := types.NewEmptyArray()
	phpArr.Append(types.NewInt(1))
	phpArr.Append(types.NewString("two"))
	phpArr.Append(types.NewFloat(3.0))

	result, err := m.ToGoSlice(types.NewArray(phpArr))
	if err != nil {
		t.Fatalf("ToGoSlice failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}

	// Test error on non-array
	_, err = m.ToGoSlice(types.NewInt(42))
	if err == nil {
		t.Error("Expected error for non-array, got nil")
	}
}

func TestToGoMapHelper(t *testing.T) {
	m := NewMarshaler()

	phpArr := types.NewEmptyArray()
	phpArr.Set(types.NewString("key1"), types.NewInt(1))
	phpArr.Set(types.NewString("key2"), types.NewString("value"))

	result, err := m.ToGoMap(types.NewArray(phpArr))
	if err != nil {
		t.Fatalf("ToGoMap failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}

	// Test error on non-array
	_, err = m.ToGoMap(types.NewInt(42))
	if err == nil {
		t.Error("Expected error for non-array, got nil")
	}
}

// Test conversions: Go → PHP

func TestToPHPNil(t *testing.T) {
	m := NewMarshaler()

	phpVal, err := m.ToPHP(nil)
	if err != nil {
		t.Fatalf("ToPHP(nil) failed: %v", err)
	}
	if phpVal.Type() != types.TypeNull {
		t.Errorf("Expected TypeNull, got %s", phpVal.TypeString())
	}
}

func TestToPHPBool(t *testing.T) {
	m := NewMarshaler()

	tests := []bool{true, false}

	for _, input := range tests {
		phpVal, err := m.ToPHP(input)
		if err != nil {
			t.Fatalf("ToPHP(bool) failed: %v", err)
		}
		if phpVal.Type() != types.TypeBool {
			t.Errorf("Expected TypeBool, got %s", phpVal.TypeString())
		}
		if phpVal.ToBool() != input {
			t.Errorf("Expected %v, got %v", input, phpVal.ToBool())
		}
	}
}

func TestToPHPInt(t *testing.T) {
	m := NewMarshaler()

	tests := []int64{0, 1, -1, 42, 1234567890}

	for _, input := range tests {
		phpVal, err := m.ToPHP(input)
		if err != nil {
			t.Fatalf("ToPHP(int64) failed: %v", err)
		}
		if phpVal.Type() != types.TypeInt {
			t.Errorf("Expected TypeInt, got %s", phpVal.TypeString())
		}
		if phpVal.ToInt() != input {
			t.Errorf("Expected %d, got %d", input, phpVal.ToInt())
		}
	}
}

func TestToPHPIntVariants(t *testing.T) {
	m := NewMarshaler()

	tests := []struct {
		input    interface{}
		expected int64
	}{
		{int(42), 42},
		{int8(42), 42},
		{int16(42), 42},
		{int32(42), 42},
		{uint(42), 42},
		{uint8(42), 42},
		{uint16(42), 42},
		{uint32(42), 42},
	}

	for _, tt := range tests {
		phpVal, err := m.ToPHP(tt.input)
		if err != nil {
			t.Fatalf("ToPHP(%T) failed: %v", tt.input, err)
		}
		if phpVal.Type() != types.TypeInt {
			t.Errorf("ToPHP(%T): expected TypeInt, got %s", tt.input, phpVal.TypeString())
		}
		if phpVal.ToInt() != tt.expected {
			t.Errorf("ToPHP(%T): expected %d, got %d", tt.input, tt.expected, phpVal.ToInt())
		}
	}
}

func TestToPHPFloat(t *testing.T) {
	m := NewMarshaler()

	tests := []float64{0.0, 1.5, -3.14, 123.456}

	for _, input := range tests {
		phpVal, err := m.ToPHP(input)
		if err != nil {
			t.Fatalf("ToPHP(float64) failed: %v", err)
		}
		if phpVal.Type() != types.TypeFloat {
			t.Errorf("Expected TypeFloat, got %s", phpVal.TypeString())
		}
		if phpVal.ToFloat() != input {
			t.Errorf("Expected %f, got %f", input, phpVal.ToFloat())
		}
	}
}

func TestToPHPString(t *testing.T) {
	m := NewMarshaler()

	tests := []string{"", "hello", "with spaces", "unicode: 你好"}

	for _, input := range tests {
		phpVal, err := m.ToPHP(input)
		if err != nil {
			t.Fatalf("ToPHP(string) failed: %v", err)
		}
		if phpVal.Type() != types.TypeString {
			t.Errorf("Expected TypeString, got %s", phpVal.TypeString())
		}
		if phpVal.ToString() != input {
			t.Errorf("Expected %q, got %q", input, phpVal.ToString())
		}
	}
}

func TestToPHPSlice(t *testing.T) {
	m := NewMarshaler()

	input := []interface{}{int64(1), "two", 3.0}
	phpVal, err := m.ToPHP(input)
	if err != nil {
		t.Fatalf("ToPHP(slice) failed: %v", err)
	}

	if phpVal.Type() != types.TypeArray {
		t.Fatalf("Expected TypeArray, got %s", phpVal.TypeString())
	}

	arr := phpVal.ToArray()
	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	// Check elements
	val0, _ := arr.Get(types.NewInt(0))
	if val0.ToInt() != 1 {
		t.Errorf("Element 0: expected 1, got %v", val0.ToInt())
	}

	val1, _ := arr.Get(types.NewInt(1))
	if val1.ToString() != "two" {
		t.Errorf("Element 1: expected 'two', got %v", val1.ToString())
	}

	val2, _ := arr.Get(types.NewInt(2))
	if val2.ToFloat() != 3.0 {
		t.Errorf("Element 2: expected 3.0, got %v", val2.ToFloat())
	}
}

func TestToPHPMap(t *testing.T) {
	m := NewMarshaler()

	input := map[string]interface{}{
		"name": "John",
		"age":  int64(30),
	}

	phpVal, err := m.ToPHP(input)
	if err != nil {
		t.Fatalf("ToPHP(map) failed: %v", err)
	}

	if phpVal.Type() != types.TypeArray {
		t.Fatalf("Expected TypeArray, got %s", phpVal.TypeString())
	}

	arr := phpVal.ToArray()
	if arr.Len() != 2 {
		t.Errorf("Expected length 2, got %d", arr.Len())
	}

	// Check values
	name, _ := arr.Get(types.NewString("name"))
	if name.ToString() != "John" {
		t.Errorf("Expected name='John', got %v", name.ToString())
	}

	age, _ := arr.Get(types.NewString("age"))
	if age.ToInt() != 30 {
		t.Errorf("Expected age=30, got %v", age.ToInt())
	}
}

func TestToPHPStruct(t *testing.T) {
	m := NewMarshaler()

	type Person struct {
		Name string
		Age  int
	}

	input := Person{Name: "Alice", Age: 25}
	phpVal, err := m.ToPHP(input)
	if err != nil {
		t.Fatalf("ToPHP(struct) failed: %v", err)
	}

	if phpVal.Type() != types.TypeArray {
		t.Fatalf("Expected TypeArray, got %s", phpVal.TypeString())
	}

	arr := phpVal.ToArray()

	name, _ := arr.Get(types.NewString("Name"))
	if name.ToString() != "Alice" {
		t.Errorf("Expected Name='Alice', got %v", name.ToString())
	}

	age, _ := arr.Get(types.NewString("Age"))
	if age.ToInt() != 25 {
		t.Errorf("Expected Age=25, got %v", age.ToInt())
	}
}

// Test helper functions

func TestIntToPHP(t *testing.T) {
	m := NewMarshaler()
	phpVal := m.IntToPHP(42)
	if phpVal.Type() != types.TypeInt || phpVal.ToInt() != 42 {
		t.Errorf("IntToPHP failed")
	}
}

func TestFloatToPHP(t *testing.T) {
	m := NewMarshaler()
	phpVal := m.FloatToPHP(3.14)
	if phpVal.Type() != types.TypeFloat || phpVal.ToFloat() != 3.14 {
		t.Errorf("FloatToPHP failed")
	}
}

func TestBoolToPHP(t *testing.T) {
	m := NewMarshaler()
	phpVal := m.BoolToPHP(true)
	if phpVal.Type() != types.TypeBool || !phpVal.ToBool() {
		t.Errorf("BoolToPHP failed")
	}
}

func TestStringToPHP(t *testing.T) {
	m := NewMarshaler()
	phpVal := m.StringToPHP("hello")
	if phpVal.Type() != types.TypeString || phpVal.ToString() != "hello" {
		t.Errorf("StringToPHP failed")
	}
}

func TestSliceToPHP(t *testing.T) {
	m := NewMarshaler()
	input := []interface{}{int64(1), "two", 3.0}
	phpVal, err := m.SliceToPHP(input)
	if err != nil {
		t.Fatalf("SliceToPHP failed: %v", err)
	}
	if phpVal.Type() != types.TypeArray || phpVal.ToArray().Len() != 3 {
		t.Errorf("SliceToPHP produced incorrect array")
	}
}

func TestMapToPHP(t *testing.T) {
	m := NewMarshaler()
	input := map[string]interface{}{"a": int64(1), "b": "two"}
	phpVal, err := m.MapToPHP(input)
	if err != nil {
		t.Fatalf("MapToPHP failed: %v", err)
	}
	if phpVal.Type() != types.TypeArray || phpVal.ToArray().Len() != 2 {
		t.Errorf("MapToPHP produced incorrect array")
	}
}

// Test edge cases

func TestEmptyArrayConversions(t *testing.T) {
	m := NewMarshaler()

	// Empty PHP array to Go
	phpArr := types.NewEmptyArray()
	goVal, err := m.ToGo(types.NewArray(phpArr))
	if err != nil {
		t.Fatalf("ToGo(empty array) failed: %v", err)
	}
	slice, ok := goVal.([]interface{})
	if !ok {
		t.Errorf("Expected []interface{}, got %T", goVal)
	}
	if len(slice) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(slice))
	}

	// Empty Go slice to PHP
	goSlice := []interface{}{}
	phpVal, err := m.ToPHP(goSlice)
	if err != nil {
		t.Fatalf("ToPHP(empty slice) failed: %v", err)
	}
	if phpVal.Type() != types.TypeArray {
		t.Errorf("Expected TypeArray, got %s", phpVal.TypeString())
	}
	if phpVal.ToArray().Len() != 0 {
		t.Errorf("Expected empty array, got length %d", phpVal.ToArray().Len())
	}
}

func TestNestedStructures(t *testing.T) {
	m := NewMarshaler()

	// Nested PHP array
	inner := types.NewEmptyArray()
	inner.Append(types.NewInt(1))
	inner.Append(types.NewInt(2))

	outer := types.NewEmptyArray()
	outer.Set(types.NewString("data"), types.NewArray(inner))
	outer.Set(types.NewString("count"), types.NewInt(2))

	goVal, err := m.ToGo(types.NewArray(outer))
	if err != nil {
		t.Fatalf("ToGo(nested) failed: %v", err)
	}

	result, ok := goVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", goVal)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Errorf("Expected data to be []interface{}, got %T", result["data"])
	}
	if len(data) != 2 {
		t.Errorf("Expected data length 2, got %d", len(data))
	}

	// Nested Go structure
	goNested := map[string]interface{}{
		"list": []interface{}{int64(1), int64(2), int64(3)},
		"meta": map[string]interface{}{
			"count": int64(3),
		},
	}

	phpVal, err := m.ToPHP(goNested)
	if err != nil {
		t.Fatalf("ToPHP(nested) failed: %v", err)
	}

	if phpVal.Type() != types.TypeArray {
		t.Errorf("Expected TypeArray, got %s", phpVal.TypeString())
	}
}

func TestPointerHandling(t *testing.T) {
	m := NewMarshaler()

	// Pointer to int
	val := int64(42)
	phpVal, err := m.ToPHP(&val)
	if err != nil {
		t.Fatalf("ToPHP(*int64) failed: %v", err)
	}
	if phpVal.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", phpVal.ToInt())
	}

	// Nil pointer
	var nilPtr *int64
	phpVal, err = m.ToPHP(nilPtr)
	if err != nil {
		t.Fatalf("ToPHP(nil pointer) failed: %v", err)
	}
	if phpVal.Type() != types.TypeNull {
		t.Errorf("Expected TypeNull, got %s", phpVal.TypeString())
	}
}

func TestReflectComplexTypes(t *testing.T) {
	m := NewMarshaler()

	// Test with reflection-based conversion
	type Complex struct {
		Numbers []int
		Nested  map[string]string
	}

	input := Complex{
		Numbers: []int{1, 2, 3},
		Nested:  map[string]string{"key": "value"},
	}

	phpVal, err := m.ToPHP(input)
	if err != nil {
		t.Fatalf("ToPHP(complex) failed: %v", err)
	}

	if phpVal.Type() != types.TypeArray {
		t.Errorf("Expected TypeArray, got %s", phpVal.TypeString())
	}

	arr := phpVal.ToArray()
	numbers, _ := arr.Get(types.NewString("Numbers"))
	if numbers.Type() != types.TypeArray {
		t.Errorf("Expected Numbers to be array, got %s", numbers.TypeString())
	}
}

func TestIsListArray(t *testing.T) {
	m := NewMarshaler()

	// List array [0, 1, 2]
	list := types.NewEmptyArray()
	list.Append(types.NewInt(10))
	list.Append(types.NewInt(20))
	list.Append(types.NewInt(30))

	if !m.isListArray(types.NewArray(list)) {
		t.Error("Expected list array to be detected as list")
	}

	// Associative array ['a' => 1]
	assoc := types.NewEmptyArray()
	assoc.Set(types.NewString("a"), types.NewInt(1))

	if m.isListArray(types.NewArray(assoc)) {
		t.Error("Expected assoc array not to be detected as list")
	}

	// Non-sequential [0 => 1, 2 => 3] (missing 1)
	nonSeq := types.NewEmptyArray()
	nonSeq.Set(types.NewInt(0), types.NewInt(1))
	nonSeq.Set(types.NewInt(2), types.NewInt(3))

	if m.isListArray(types.NewArray(nonSeq)) {
		t.Error("Expected non-sequential array not to be detected as list")
	}

	// Empty array
	empty := types.NewEmptyArray()
	if !m.isListArray(types.NewArray(empty)) {
		t.Error("Expected empty array to be detected as list")
	}
}

func TestRoundTrip(t *testing.T) {
	m := NewMarshaler()

	// Test round-trip: Go → PHP → Go
	original := map[string]interface{}{
		"name":   "Alice",
		"age":    int64(30),
		"active": true,
		"scores": []interface{}{int64(95), int64(87), int64(92)},
	}

	// Go → PHP
	phpVal, err := m.ToPHP(original)
	if err != nil {
		t.Fatalf("ToPHP failed: %v", err)
	}

	// PHP → Go
	result, err := m.ToGo(phpVal)
	if err != nil {
		t.Fatalf("ToGo failed: %v", err)
	}

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", result)
	}

	if resultMap["name"] != "Alice" {
		t.Errorf("name mismatch: expected 'Alice', got %v", resultMap["name"])
	}
	if resultMap["age"] != int64(30) {
		t.Errorf("age mismatch: expected 30, got %v", resultMap["age"])
	}
	if resultMap["active"] != true {
		t.Errorf("active mismatch: expected true, got %v", resultMap["active"])
	}

	scores, ok := resultMap["scores"].([]interface{})
	if !ok {
		t.Fatalf("Expected scores to be []interface{}, got %T", resultMap["scores"])
	}
	if !reflect.DeepEqual(scores, []interface{}{int64(95), int64(87), int64(92)}) {
		t.Errorf("scores mismatch: got %v", scores)
	}
}
