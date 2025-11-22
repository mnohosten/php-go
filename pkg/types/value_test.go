package types

import (
	"math"
	"testing"
)

// ============================================================================
// Type Constructor Tests
// ============================================================================

func TestNewUndef(t *testing.T) {
	v := NewUndef()
	if v.Type() != TypeUndef {
		t.Errorf("Expected TypeUndef, got %v", v.Type())
	}
	if !v.IsUndef() {
		t.Error("Expected IsUndef() to be true")
	}
}

func TestNewNull(t *testing.T) {
	v := NewNull()
	if v.Type() != TypeNull {
		t.Errorf("Expected TypeNull, got %v", v.Type())
	}
	if !v.IsNull() {
		t.Error("Expected IsNull() to be true")
	}
}

func TestNewBool(t *testing.T) {
	vTrue := NewBool(true)
	if vTrue.Type() != TypeBool {
		t.Errorf("Expected TypeBool, got %v", vTrue.Type())
	}
	if !vTrue.IsBool() {
		t.Error("Expected IsBool() to be true")
	}
	if !vTrue.ToBool() {
		t.Error("Expected ToBool() to be true")
	}

	vFalse := NewBool(false)
	if vFalse.ToBool() {
		t.Error("Expected ToBool() to be false")
	}
}

func TestNewInt(t *testing.T) {
	v := NewInt(42)
	if v.Type() != TypeInt {
		t.Errorf("Expected TypeInt, got %v", v.Type())
	}
	if !v.IsInt() {
		t.Error("Expected IsInt() to be true")
	}
	if v.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", v.ToInt())
	}
}

func TestNewFloat(t *testing.T) {
	v := NewFloat(3.14)
	if v.Type() != TypeFloat {
		t.Errorf("Expected TypeFloat, got %v", v.Type())
	}
	if !v.IsFloat() {
		t.Error("Expected IsFloat() to be true")
	}
	if v.ToFloat() != 3.14 {
		t.Errorf("Expected 3.14, got %f", v.ToFloat())
	}
}

func TestNewString(t *testing.T) {
	v := NewString("hello")
	if v.Type() != TypeString {
		t.Errorf("Expected TypeString, got %v", v.Type())
	}
	if !v.IsString() {
		t.Error("Expected IsString() to be true")
	}
	if v.ToString() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", v.ToString())
	}
}

func TestNewArray(t *testing.T) {
	arr := NewEmptyArray()
	v := NewArray(arr)
	if v.Type() != TypeArray {
		t.Errorf("Expected TypeArray, got %v", v.Type())
	}
	if !v.IsArray() {
		t.Error("Expected IsArray() to be true")
	}
}

func TestNewObject(t *testing.T) {
	obj := NewObjectInstance("TestClass")
	v := NewObject(obj)
	if v.Type() != TypeObject {
		t.Errorf("Expected TypeObject, got %v", v.Type())
	}
	if !v.IsObject() {
		t.Error("Expected IsObject() to be true")
	}
}

func TestNewResource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)
	if v.Type() != TypeResource {
		t.Errorf("Expected TypeResource, got %v", v.Type())
	}
	if !v.IsResource() {
		t.Error("Expected IsResource() to be true")
	}
}

func TestNewReference(t *testing.T) {
	inner := NewInt(42)
	v := NewReference(inner)
	if v.Type() != TypeReference {
		t.Errorf("Expected TypeReference, got %v", v.Type())
	}
	if !v.IsReference() {
		t.Error("Expected IsReference() to be true")
	}
	if v.Deref().ToInt() != 42 {
		t.Errorf("Expected 42, got %d", v.Deref().ToInt())
	}
}

// ============================================================================
// Type Conversion Tests
// ============================================================================

func TestToInt_FromNull(t *testing.T) {
	v := NewNull()
	if v.ToInt() != 0 {
		t.Errorf("Expected 0, got %d", v.ToInt())
	}
}

func TestToInt_FromBool(t *testing.T) {
	if NewBool(true).ToInt() != 1 {
		t.Error("Expected true -> 1")
	}
	if NewBool(false).ToInt() != 0 {
		t.Error("Expected false -> 0")
	}
}

func TestToInt_FromFloat(t *testing.T) {
	if NewFloat(3.14).ToInt() != 3 {
		t.Error("Expected 3.14 -> 3")
	}
	if NewFloat(-2.9).ToInt() != -2 {
		t.Error("Expected -2.9 -> -2")
	}
}

func TestToInt_FromString(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"123", 123},
		{"-456", -456},
		{"123abc", 123},
		{"abc", 0},
		{"  789  ", 789},
		{"", 0},
		{"0", 0},
	}

	for _, tt := range tests {
		v := NewString(tt.input)
		if v.ToInt() != tt.expected {
			t.Errorf("String '%s' -> expected %d, got %d", tt.input, tt.expected, v.ToInt())
		}
	}
}

func TestToInt_FromArray(t *testing.T) {
	// Empty array -> 0
	emptyArr := NewEmptyArray()
	if NewArray(emptyArr).ToInt() != 0 {
		t.Error("Expected empty array -> 0")
	}

	// Non-empty array -> 1
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewInt(42))
	if NewArray(arr).ToInt() != 1 {
		t.Error("Expected non-empty array -> 1")
	}
}

func TestToFloat_FromInt(t *testing.T) {
	v := NewInt(42)
	if v.ToFloat() != 42.0 {
		t.Errorf("Expected 42.0, got %f", v.ToFloat())
	}
}

func TestToFloat_FromString(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"-2.5", -2.5},
		{"3.14abc", 3.14},
		{"abc", 0.0},
		{"", 0.0},
	}

	for _, tt := range tests {
		v := NewString(tt.input)
		if v.ToFloat() != tt.expected {
			t.Errorf("String '%s' -> expected %f, got %f", tt.input, tt.expected, v.ToFloat())
		}
	}
}

func TestToBool_PHPTruthiness(t *testing.T) {
	tests := []struct {
		value    *Value
		expected bool
		name     string
	}{
		// Falsy values
		{NewNull(), false, "null"},
		{NewBool(false), false, "false"},
		{NewInt(0), false, "int(0)"},
		{NewFloat(0.0), false, "float(0.0)"},
		{NewString(""), false, "empty string"},
		{NewString("0"), false, "string '0'"},
		{NewArray(NewEmptyArray()), false, "empty array"},

		// Truthy values
		{NewBool(true), true, "true"},
		{NewInt(1), true, "int(1)"},
		{NewInt(-1), true, "int(-1)"},
		{NewFloat(0.1), true, "float(0.1)"},
		{NewFloat(-3.14), true, "float(-3.14)"},
		{NewString("0.0"), true, "string '0.0'"},
		{NewString("hello"), true, "string 'hello'"},
		{NewString("false"), true, "string 'false'"},
	}

	for _, tt := range tests {
		if tt.value.ToBool() != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, tt.value.ToBool())
		}
	}
}

func TestToBool_NaN(t *testing.T) {
	v := NewFloat(math.NaN())
	if v.ToBool() {
		t.Error("NaN should be falsy")
	}
}

func TestToString_FromBool(t *testing.T) {
	if NewBool(true).ToString() != "1" {
		t.Error("Expected true -> '1'")
	}
	if NewBool(false).ToString() != "" {
		t.Error("Expected false -> ''")
	}
}

func TestToString_FromInt(t *testing.T) {
	if NewInt(42).ToString() != "42" {
		t.Error("Expected '42'")
	}
	if NewInt(-123).ToString() != "-123" {
		t.Error("Expected '-123'")
	}
}

func TestToString_FromFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{3.14, "3.14"},
		{math.NaN(), "NAN"},
		{math.Inf(1), "INF"},
		{math.Inf(-1), "-INF"},
	}

	for _, tt := range tests {
		v := NewFloat(tt.input)
		if v.ToString() != tt.expected {
			t.Errorf("Float %f -> expected '%s', got '%s'", tt.input, tt.expected, v.ToString())
		}
	}
}

func TestToString_FromArray(t *testing.T) {
	arr := NewEmptyArray()
	v := NewArray(arr)
	if v.ToString() != "Array" {
		t.Errorf("Expected 'Array', got '%s'", v.ToString())
	}
}

func TestToString_FromObject(t *testing.T) {
	obj := NewObjectInstance("TestClass")
	v := NewObject(obj)
	if v.ToString() != "Object" {
		t.Errorf("Expected 'Object', got '%s'", v.ToString())
	}
}

func TestToString_FromResource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)
	s := v.ToString()
	if s == "" {
		t.Error("Expected non-empty resource string")
	}
}

// ============================================================================
// IsTrue/IsFalse Tests
// ============================================================================

func TestIsTrue_IsFalse(t *testing.T) {
	trueVal := NewBool(true)
	if !trueVal.IsTrue() {
		t.Error("Expected IsTrue() to be true")
	}
	if trueVal.IsFalse() {
		t.Error("Expected IsFalse() to be false")
	}

	falseVal := NewBool(false)
	if falseVal.IsTrue() {
		t.Error("Expected IsTrue() to be false")
	}
	if !falseVal.IsFalse() {
		t.Error("Expected IsFalse() to be true")
	}
}

// ============================================================================
// Copy Tests
// ============================================================================

func TestCopy_Scalars(t *testing.T) {
	tests := []*Value{
		NewNull(),
		NewBool(true),
		NewInt(42),
		NewFloat(3.14),
		NewString("hello"),
	}

	for _, original := range tests {
		copied := original.Copy()

		// Should have same type
		if copied.Type() != original.Type() {
			t.Errorf("Type mismatch: expected %v, got %v", original.Type(), copied.Type())
		}

		// Should have same value
		if !copied.Identical(original) {
			t.Errorf("Value mismatch after copy")
		}

		// Should be different pointer
		if copied == original {
			t.Error("Copy should create new value")
		}
	}
}

func TestDeepCopy(t *testing.T) {
	// Create nested array
	inner := NewEmptyArray()
	inner.Set(NewInt(0), NewInt(42))

	outer := NewEmptyArray()
	outer.Set(NewInt(0), NewArray(inner))

	v := NewArray(outer)
	copied := v.DeepCopy()

	// Should be different pointer
	if copied == v {
		t.Error("DeepCopy should create new value")
	}

	// Should have same type
	if copied.Type() != v.Type() {
		t.Error("Type mismatch after deep copy")
	}
}

func TestDeref(t *testing.T) {
	// Non-reference should return self
	v := NewInt(42)
	if v.Deref() != v {
		t.Error("Deref on non-reference should return self")
	}

	// Reference should return inner value
	ref := NewReference(v)
	if ref.Deref() != v {
		t.Error("Deref on reference should return inner value")
	}

	// Nested references
	ref2 := NewReference(ref)
	if ref2.Deref() != v {
		t.Error("Deref should recursively dereference")
	}
}

// ============================================================================
// Equality Tests
// ============================================================================

func TestEquals_SameType(t *testing.T) {
	tests := []struct {
		a        *Value
		b        *Value
		expected bool
		name     string
	}{
		{NewNull(), NewNull(), true, "null == null"},
		{NewBool(true), NewBool(true), true, "true == true"},
		{NewBool(true), NewBool(false), false, "true != false"},
		{NewInt(42), NewInt(42), true, "42 == 42"},
		{NewInt(42), NewInt(43), false, "42 != 43"},
		{NewFloat(3.14), NewFloat(3.14), true, "3.14 == 3.14"},
		{NewString("hello"), NewString("hello"), true, "'hello' == 'hello'"},
		{NewString("hello"), NewString("world"), false, "'hello' != 'world'"},
	}

	for _, tt := range tests {
		if tt.a.Equals(tt.b) != tt.expected {
			t.Errorf("%s: expected %v", tt.name, tt.expected)
		}
	}
}

func TestEquals_TypeJuggling(t *testing.T) {
	tests := []struct {
		a        *Value
		b        *Value
		expected bool
		name     string
	}{
		// PHP type juggling: numeric strings
		{NewInt(42), NewString("42"), true, "42 == '42'"},
		{NewFloat(3.14), NewString("3.14"), true, "3.14 == '3.14'"},

		// Int and float comparison
		{NewInt(5), NewFloat(5.0), true, "5 == 5.0"},
		{NewInt(5), NewFloat(5.1), false, "5 != 5.1"},
	}

	for _, tt := range tests {
		if tt.a.Equals(tt.b) != tt.expected {
			t.Errorf("%s: expected %v", tt.name, tt.expected)
		}
	}
}

func TestIdentical_StrictEquality(t *testing.T) {
	tests := []struct {
		a        *Value
		b        *Value
		expected bool
		name     string
	}{
		{NewInt(42), NewInt(42), true, "42 === 42"},
		{NewInt(42), NewString("42"), false, "42 !== '42' (different types)"},
		{NewInt(5), NewFloat(5.0), false, "5 !== 5.0 (different types)"},
		{NewBool(true), NewInt(1), false, "true !== 1 (different types)"},
	}

	for _, tt := range tests {
		if tt.a.Identical(tt.b) != tt.expected {
			t.Errorf("%s: expected %v", tt.name, tt.expected)
		}
	}
}

// ============================================================================
// IsScalar Tests
// ============================================================================

func TestIsScalar(t *testing.T) {
	scalars := []*Value{
		NewBool(true),
		NewInt(42),
		NewFloat(3.14),
		NewString("hello"),
	}

	for _, v := range scalars {
		if !v.IsScalar() {
			t.Errorf("Expected %v to be scalar", v.Type())
		}
	}

	nonScalars := []*Value{
		NewNull(),
		NewArray(NewEmptyArray()),
		NewObject(NewObjectInstance("Test")),
	}

	for _, v := range nonScalars {
		if v.IsScalar() {
			t.Errorf("Expected %v to not be scalar", v.Type())
		}
	}
}

// ============================================================================
// String Representation Tests
// ============================================================================

func TestString_Debugging(t *testing.T) {
	tests := []struct {
		value    *Value
		contains string
	}{
		{NewNull(), "NULL"},
		{NewBool(true), "bool(true)"},
		{NewBool(false), "bool(false)"},
		{NewInt(42), "int(42)"},
		{NewFloat(3.14), "float(3.14)"},
		{NewString("hello"), "string"},
		{NewArray(NewEmptyArray()), "array(0)"},
	}

	for _, tt := range tests {
		s := tt.value.String()
		if s == "" {
			t.Errorf("String() should not be empty for %v", tt.value.Type())
		}
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		value    *Value
		expected string
	}{
		{NewNull(), "NULL"},
		{NewBool(true), "boolean"},
		{NewInt(42), "integer"},
		{NewFloat(3.14), "double"},
		{NewString("hello"), "string"},
		{NewArray(NewEmptyArray()), "array"},
		{NewObject(NewObjectInstance("Test")), "object"},
	}

	for _, tt := range tests {
		if tt.value.TypeString() != tt.expected {
			t.Errorf("Expected type string '%s', got '%s'", tt.expected, tt.value.TypeString())
		}
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestStringToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"123", 123},
		{"-456", -456},
		{"+789", 789},
		{"123abc", 123},
		{"abc123", 0},
		{"  42  ", 42},
		{"", 0},
		{"0", 0},
		{"-0", 0},
	}

	for _, tt := range tests {
		result := stringToInt(tt.input)
		if result != tt.expected {
			t.Errorf("stringToInt('%s'): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}

func TestStringToFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"-2.5", -2.5},
		{"3.14abc", 3.14},
		{"abc", 0.0},
		{"", 0.0},
		{"1e2", 100.0},
		{"1.5e-1", 0.15},
	}

	for _, tt := range tests {
		result := stringToFloat(tt.input)
		if result != tt.expected {
			t.Errorf("stringToFloat('%s'): expected %f, got %f", tt.input, tt.expected, result)
		}
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestNilValue(t *testing.T) {
	var v *Value

	// Nil value should behave like null
	if !v.IsNull() {
		t.Error("Nil value should be considered null")
	}
	if v.Type() != TypeNull {
		t.Error("Nil value should have TypeNull")
	}
	if v.ToInt() != 0 {
		t.Error("Nil value should convert to 0")
	}
	if v.ToString() != "" {
		t.Error("Nil value should convert to empty string")
	}
	if v.ToBool() {
		t.Error("Nil value should be falsy")
	}
}

func TestToArray_Scalar(t *testing.T) {
	v := NewInt(42)
	arr := v.ToArray()

	if arr.Len() != 1 {
		t.Errorf("Expected array length 1, got %d", arr.Len())
	}

	val, exists := arr.Get(NewInt(0))
	if !exists {
		t.Error("Expected value at index 0")
	}
	if val.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", val.ToInt())
	}
}

func TestReference_ToInt(t *testing.T) {
	v := NewInt(42)
	ref := NewReference(v)

	if ref.ToInt() != 42 {
		t.Errorf("Reference should dereference for ToInt, expected 42, got %d", ref.ToInt())
	}
}

func TestReference_ToString(t *testing.T) {
	v := NewString("hello")
	ref := NewReference(v)

	if ref.ToString() != "hello" {
		t.Errorf("Reference should dereference for ToString, expected 'hello', got '%s'", ref.ToString())
	}
}

// ============================================================================
// Additional Coverage Tests
// ============================================================================

func TestToFloat_AllTypes(t *testing.T) {
	tests := []struct {
		value    *Value
		expected float64
		name     string
	}{
		{NewNull(), 0.0, "null"},
		{NewUndef(), 0.0, "undef"},
		{NewBool(true), 1.0, "true"},
		{NewBool(false), 0.0, "false"},
		{NewInt(42), 42.0, "int"},
		{NewFloat(3.14), 3.14, "float"},
		{NewString("2.5"), 2.5, "string"},
		{NewArray(NewEmptyArray()), 0.0, "empty array"},
	}

	for _, tt := range tests {
		if tt.value.ToFloat() != tt.expected {
			t.Errorf("%s: expected %f, got %f", tt.name, tt.expected, tt.value.ToFloat())
		}
	}
}

func TestToFloat_NonEmptyArray(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewInt(42))
	v := NewArray(arr)

	if v.ToFloat() != 1.0 {
		t.Errorf("Non-empty array should convert to 1.0, got %f", v.ToFloat())
	}
}

func TestToFloat_Object(t *testing.T) {
	obj := NewObjectInstance("Test")
	v := NewObject(obj)

	if v.ToFloat() != 1.0 {
		t.Errorf("Object should convert to 1.0, got %f", v.ToFloat())
	}
}

func TestToFloat_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	// Should convert to resource ID
	expected := float64(res.ID())
	if v.ToFloat() != expected {
		t.Errorf("Resource should convert to ID as float, expected %f, got %f", expected, v.ToFloat())
	}
}

func TestToFloat_Reference(t *testing.T) {
	v := NewFloat(3.14)
	ref := NewReference(v)

	if ref.ToFloat() != 3.14 {
		t.Errorf("Reference should dereference for ToFloat, expected 3.14, got %f", ref.ToFloat())
	}
}

func TestToInt_Object(t *testing.T) {
	obj := NewObjectInstance("Test")
	v := NewObject(obj)

	if v.ToInt() != 1 {
		t.Errorf("Object should convert to 1, got %d", v.ToInt())
	}
}

func TestToInt_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	expected := int64(res.ID())
	if v.ToInt() != expected {
		t.Errorf("Resource should convert to ID, expected %d, got %d", expected, v.ToInt())
	}
}

func TestToBool_Object(t *testing.T) {
	obj := NewObjectInstance("Test")
	v := NewObject(obj)

	if !v.ToBool() {
		t.Error("Object should be truthy")
	}
}

func TestToBool_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	if !v.ToBool() {
		t.Error("Resource should be truthy")
	}
}

func TestToString_Reference(t *testing.T) {
	v := NewInt(42)
	ref := NewReference(v)

	if ref.ToString() != "42" {
		t.Errorf("Reference should dereference for ToString, expected '42', got '%s'", ref.ToString())
	}
}

func TestToArray_Null(t *testing.T) {
	v := NewNull()
	arr := v.ToArray()

	if arr.Len() != 0 {
		t.Errorf("Null should convert to empty array, got length %d", arr.Len())
	}
}

func TestToArray_Array(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewInt(42))
	v := NewArray(arr)

	result := v.ToArray()
	if result != arr {
		t.Error("ToArray on array should return the same array")
	}
}

func TestCopy_Array(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewInt(42))
	v := NewArray(arr)

	copied := v.Copy()

	if copied == v {
		t.Error("Copy should create new value pointer")
	}

	if copied.Type() != TypeArray {
		t.Error("Copied value should have TypeArray")
	}
}

func TestCopy_Object(t *testing.T) {
	obj := NewObjectInstance("Test")
	v := NewObject(obj)

	copied := v.Copy()

	if copied == v {
		t.Error("Copy should create new value pointer")
	}

	// Objects are passed by reference, so data should be same
	if copied.Type() != TypeObject {
		t.Error("Copied value should have TypeObject")
	}
}

func TestCopy_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	copied := v.Copy()

	if copied == v {
		t.Error("Copy should create new value pointer")
	}

	if copied.Type() != TypeResource {
		t.Error("Copied value should have TypeResource")
	}
}

func TestCopy_Reference(t *testing.T) {
	inner := NewInt(42)
	v := NewReference(inner)

	copied := v.Copy()

	if copied == v {
		t.Error("Copy should create new value pointer")
	}

	if !copied.IsReference() {
		t.Error("Copied value should still be a reference")
	}
}

func TestDeepCopy_Object(t *testing.T) {
	obj := NewObjectInstance("Test")
	v := NewObject(obj)

	copied := v.DeepCopy()

	if copied == v {
		t.Error("DeepCopy should create new value pointer")
	}

	if copied.Type() != TypeObject {
		t.Error("Deep copied value should have TypeObject")
	}
}

func TestDeepCopy_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	copied := v.DeepCopy()

	if copied == v {
		t.Error("DeepCopy should create new value pointer")
	}
}

func TestDeepCopy_Reference(t *testing.T) {
	inner := NewInt(42)
	v := NewReference(inner)

	copied := v.DeepCopy()

	if copied == v {
		t.Error("DeepCopy should create new value pointer")
	}

	// Should deep copy the referenced value
	if copied.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", copied.ToInt())
	}
}

func TestEquals_Null(t *testing.T) {
	v1 := NewNull()
	v2 := NewNull()

	if !v1.Equals(v2) {
		t.Error("null == null should be true")
	}
}

func TestEquals_BothNil(t *testing.T) {
	var v1 *Value
	var v2 *Value

	if !v1.Equals(v2) {
		t.Error("nil == nil should be true")
	}
}

func TestEquals_OneNil(t *testing.T) {
	var v1 *Value
	v2 := NewInt(42)

	if v1.Equals(v2) {
		t.Error("nil != 42 should be false")
	}
}

func TestIdentical_BothNil(t *testing.T) {
	var v1 *Value
	var v2 *Value

	if !v1.Identical(v2) {
		t.Error("nil === nil should be true")
	}
}

func TestIdentical_OneNil(t *testing.T) {
	var v1 *Value
	v2 := NewInt(42)

	if v1.Identical(v2) {
		t.Error("nil !== 42 should be false")
	}
}

func TestString_Reference(t *testing.T) {
	v := NewInt(42)
	ref := NewReference(v)

	s := ref.String()
	if s == "" {
		t.Error("Reference String() should not be empty")
	}
}

func TestString_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	s := v.String()
	if s == "" {
		t.Error("Resource String() should not be empty")
	}
}

func TestTypeString_Undef(t *testing.T) {
	v := NewUndef()
	if v.TypeString() != "undefined" {
		t.Errorf("Expected 'undefined', got '%s'", v.TypeString())
	}
}

func TestTypeString_Reference(t *testing.T) {
	v := NewReference(NewInt(42))
	if v.TypeString() != "reference" {
		t.Errorf("Expected 'reference', got '%s'", v.TypeString())
	}
}

func TestTypeString_Resource(t *testing.T) {
	res := NewResourceHandle("file", nil)
	v := NewResource(res)

	if v.TypeString() != "resource" {
		t.Errorf("Expected 'resource', got '%s'", v.TypeString())
	}
}

func TestIsScalar_Nil(t *testing.T) {
	var v *Value
	if v.IsScalar() {
		t.Error("Nil should not be scalar")
	}
}

func TestArray_GetNonExistent(t *testing.T) {
	arr := NewEmptyArray()
	val, exists := arr.Get(NewInt(0))

	if exists {
		t.Error("Getting non-existent key should return false")
	}

	if !val.IsNull() {
		t.Error("Getting non-existent key should return null value")
	}
}

func TestArray_GetNil(t *testing.T) {
	var arr *Array
	val, exists := arr.Get(NewInt(0))

	if exists {
		t.Error("Getting from nil array should return false")
	}

	if !val.IsNull() {
		t.Error("Getting from nil array should return null value")
	}
}

func TestArray_SetNil(t *testing.T) {
	var arr *Array
	// Should not panic
	arr.Set(NewInt(0), NewInt(42))
}

func TestArray_LenNil(t *testing.T) {
	var arr *Array
	if arr.Len() != 0 {
		t.Error("Nil array should have length 0")
	}
}

func TestArray_DeepCopyNil(t *testing.T) {
	var arr *Array
	copied := arr.DeepCopy()

	if copied.Len() != 0 {
		t.Error("Deep copy of nil array should be empty")
	}
}

func TestArray_StringKey(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewString("key"), NewString("value"))

	val, exists := arr.Get(NewString("key"))
	if !exists {
		t.Error("String key should exist")
	}

	if val.ToString() != "value" {
		t.Errorf("Expected 'value', got '%s'", val.ToString())
	}
}

func TestArray_OtherKeyType(t *testing.T) {
	arr := NewEmptyArray()
	// Bool key should convert to string
	arr.Set(NewBool(true), NewInt(42))

	val, exists := arr.Get(NewString("1"))
	if !exists {
		t.Error("Bool key converted to string should be retrievable")
	}

	if val.ToInt() != 42 {
		t.Errorf("Expected 42, got %d", val.ToInt())
	}
}
