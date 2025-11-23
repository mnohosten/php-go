package varfuncs

import (
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// var_dump Tests
// ============================================================================

func TestVarDump_Null(t *testing.T) {
	result := VarDump(types.NewNull())
	if result.Type() != types.TypeNull {
		t.Errorf("VarDump should return NULL, got %v", result.Type())
	}
}

func TestVarDump_Bool(t *testing.T) {
	VarDump(types.NewBool(true))
	VarDump(types.NewBool(false))
	// Just ensure it doesn't crash
}

func TestVarDump_Int(t *testing.T) {
	VarDump(types.NewInt(42))
	VarDump(types.NewInt(-100))
	VarDump(types.NewInt(0))
}

func TestVarDump_Float(t *testing.T) {
	VarDump(types.NewFloat(3.14))
	VarDump(types.NewFloat(-2.5))
	VarDump(types.NewFloat(0.0))
}

func TestVarDump_String(t *testing.T) {
	VarDump(types.NewString("hello"))
	VarDump(types.NewString(""))
	VarDump(types.NewString("multi\nline"))
}

func TestVarDump_Array(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))

	VarDump(types.NewArray(arr))
}

func TestVarDump_AssociativeArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))

	VarDump(types.NewArray(arr))
}

// ============================================================================
// print_r Tests
// ============================================================================

func TestPrintR_Null(t *testing.T) {
	result := PrintR(types.NewNull())
	if result.Type() != types.TypeBool {
		t.Errorf("PrintR should return bool, got %v", result.Type())
	}
	if !result.ToBool() {
		t.Errorf("PrintR should return true")
	}
}

func TestPrintR_Int(t *testing.T) {
	result := PrintR(types.NewInt(42))
	if !result.ToBool() {
		t.Errorf("PrintR should return true")
	}
}

func TestPrintR_String(t *testing.T) {
	result := PrintR(types.NewString("hello"))
	if !result.ToBool() {
		t.Errorf("PrintR should return true")
	}
}

func TestPrintR_Array(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))

	result := PrintR(types.NewArray(arr))
	if !result.ToBool() {
		t.Errorf("PrintR should return true")
	}
}

func TestPrintR_WithReturn(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))

	result := PrintR(types.NewArray(arr), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Errorf("PrintR with return should return string, got %v", result.Type())
	}

	output := result.ToString()
	if !strings.Contains(output, "Array") {
		t.Errorf("PrintR output should contain 'Array', got: %s", output)
	}
}

// ============================================================================
// var_export Tests
// ============================================================================

func TestVarExport_Null(t *testing.T) {
	result := VarExport(types.NewNull())
	if result.Type() != types.TypeNull {
		t.Errorf("VarExport should return NULL, got %v", result.Type())
	}
}

func TestVarExport_Bool(t *testing.T) {
	VarExport(types.NewBool(true))
	VarExport(types.NewBool(false))
}

func TestVarExport_Int(t *testing.T) {
	VarExport(types.NewInt(42))
}

func TestVarExport_Float(t *testing.T) {
	VarExport(types.NewFloat(3.14))
}

func TestVarExport_String(t *testing.T) {
	VarExport(types.NewString("hello"))
}

func TestVarExport_Array(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))

	VarExport(types.NewArray(arr))
}

func TestVarExport_WithReturn(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))

	result := VarExport(types.NewArray(arr), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Errorf("VarExport with return should return string, got %v", result.Type())
	}

	output := result.ToString()
	if !strings.Contains(output, "array") {
		t.Errorf("VarExport output should contain 'array', got: %s", output)
	}
}

// ============================================================================
// Type Checking Tests
// ============================================================================

func TestIsNull(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewNull(), true},
		{types.NewBool(false), false},
		{types.NewInt(0), false},
		{types.NewString(""), false},
	}

	for _, tt := range tests {
		result := IsNull(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsNull(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsBool(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewBool(true), true},
		{types.NewBool(false), true},
		{types.NewInt(1), false},
		{types.NewString("true"), false},
	}

	for _, tt := range tests {
		result := IsBool(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsBool(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsInt(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewInt(42), true},
		{types.NewInt(0), true},
		{types.NewInt(-100), true},
		{types.NewFloat(42.0), false},
		{types.NewString("42"), false},
	}

	for _, tt := range tests {
		result := IsInt(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsInt(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsFloat(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewFloat(3.14), true},
		{types.NewFloat(0.0), true},
		{types.NewInt(42), false},
		{types.NewString("3.14"), false},
	}

	for _, tt := range tests {
		result := IsFloat(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsFloat(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsString(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewString("hello"), true},
		{types.NewString(""), true},
		{types.NewInt(42), false},
		{types.NewBool(false), false},
	}

	for _, tt := range tests {
		result := IsString(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsString(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsArray(t *testing.T) {
	arr := types.NewEmptyArray()

	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewArray(arr), true},
		{types.NewString("array"), false},
		{types.NewInt(0), false},
	}

	for _, tt := range tests {
		result := IsArray(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsArray(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsObject(t *testing.T) {
	// Create a simple class for testing
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)

	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewObject(obj), true},
		{types.NewString("object"), false},
		{types.NewInt(0), false},
	}

	for _, tt := range tests {
		result := IsObject(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsObject(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsResource(t *testing.T) {
	// Create a test resource
	res := types.NewResourceHandle("test", "test data")

	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewResource(res), true},
		{types.NewString("resource"), false},
		{types.NewInt(0), false},
	}

	for _, tt := range tests {
		result := IsResource(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsResource(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}

	// Test closed resource
	res.Close()
	result := IsResource(types.NewResource(res))
	if result.ToBool() {
		t.Errorf("IsResource(closed resource) should return false")
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewInt(42), true},
		{types.NewFloat(3.14), true},
		{types.NewString("42"), true},
		{types.NewString("3.14"), true},
		{types.NewString("hello"), false},
		{types.NewBool(true), false},
	}

	for _, tt := range tests {
		result := IsNumeric(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsNumeric(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsScalar(t *testing.T) {
	arr := types.NewEmptyArray()

	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewInt(42), true},
		{types.NewFloat(3.14), true},
		{types.NewString("hello"), true},
		{types.NewBool(true), true},
		{types.NewArray(arr), false},
		{types.NewNull(), false},
	}

	for _, tt := range tests {
		result := IsScalar(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsScalar(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestIsIterable(t *testing.T) {
	arr := types.NewEmptyArray()
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)

	tests := []struct {
		value    *types.Value
		expected bool
	}{
		{types.NewArray(arr), true},
		{types.NewObject(obj), true},
		{types.NewString("hello"), false},
		{types.NewInt(42), false},
	}

	for _, tt := range tests {
		result := IsIterable(tt.value)
		if result.ToBool() != tt.expected {
			t.Errorf("IsIterable(%v) = %v, want %v", tt.value, result.ToBool(), tt.expected)
		}
	}
}

func TestGetType(t *testing.T) {
	arr := types.NewEmptyArray()
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)
	res := types.NewResourceHandle("test", "data")

	tests := []struct {
		value    *types.Value
		expected string
	}{
		{types.NewNull(), "NULL"},
		{types.NewBool(true), "boolean"},
		{types.NewInt(42), "integer"},
		{types.NewFloat(3.14), "double"},
		{types.NewString("hello"), "string"},
		{types.NewArray(arr), "array"},
		{types.NewObject(obj), "object"},
		{types.NewResource(res), "resource"},
	}

	for _, tt := range tests {
		result := GetType(tt.value)
		if result.ToString() != tt.expected {
			t.Errorf("GetType(%v) = %v, want %v", tt.value, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Alias Tests
// ============================================================================

func TestIsLong(t *testing.T) {
	result := IsLong(types.NewInt(42))
	if !result.ToBool() {
		t.Errorf("IsLong(42) should return true")
	}
}

func TestIsInteger(t *testing.T) {
	result := IsInteger(types.NewInt(42))
	if !result.ToBool() {
		t.Errorf("IsInteger(42) should return true")
	}
}

func TestIsDouble(t *testing.T) {
	result := IsDouble(types.NewFloat(3.14))
	if !result.ToBool() {
		t.Errorf("IsDouble(3.14) should return true")
	}
}

func TestIsReal(t *testing.T) {
	result := IsReal(types.NewFloat(3.14))
	if !result.ToBool() {
		t.Errorf("IsReal(3.14) should return true")
	}
}

// ============================================================================
// IsCallable Tests
// ============================================================================

func TestIsCallable_String(t *testing.T) {
	// String function names are callable
	result := IsCallable(types.NewString("strlen"))
	if !result.ToBool() {
		t.Error("String should be considered callable")
	}
}

func TestIsCallable_Array(t *testing.T) {
	// Array with 2 elements is callable
	arr := types.NewEmptyArray()
	arr.Append(types.NewString("ClassName"))
	arr.Append(types.NewString("methodName"))

	result := IsCallable(types.NewArray(arr))
	if !result.ToBool() {
		t.Error("Array with 2 elements should be callable")
	}
}

func TestIsCallable_ArrayInvalid(t *testing.T) {
	// Array with 1 element is not callable
	arr := types.NewEmptyArray()
	arr.Append(types.NewString("function"))

	result := IsCallable(types.NewArray(arr))
	if result.ToBool() {
		t.Error("Array with 1 element should not be callable")
	}

	// Array with 3 elements is not callable
	arr2 := types.NewEmptyArray()
	arr2.Append(types.NewString("a"))
	arr2.Append(types.NewString("b"))
	arr2.Append(types.NewString("c"))

	result2 := IsCallable(types.NewArray(arr2))
	if result2.ToBool() {
		t.Error("Array with 3 elements should not be callable")
	}
}

func TestIsCallable_Object(t *testing.T) {
	// Object with __invoke is callable
	class := types.NewClassEntry("CallableClass")
	class.Methods["__invoke"] = &types.MethodDef{
		Name:       "__invoke",
		Visibility: types.VisibilityPublic,
	}
	obj := types.NewObjectFromClass(class)

	result := IsCallable(types.NewObject(obj))
	if !result.ToBool() {
		t.Error("Object with __invoke should be callable")
	}
}

func TestIsCallable_ObjectNoInvoke(t *testing.T) {
	// Object without __invoke is not callable
	class := types.NewClassEntry("NonCallableClass")
	obj := types.NewObjectFromClass(class)

	result := IsCallable(types.NewObject(obj))
	if result.ToBool() {
		t.Error("Object without __invoke should not be callable")
	}
}

func TestIsCallable_Other(t *testing.T) {
	// Other types are not callable
	tests := []*types.Value{
		types.NewInt(42),
		types.NewFloat(3.14),
		types.NewBool(true),
		types.NewNull(),
	}

	for _, val := range tests {
		result := IsCallable(val)
		if result.ToBool() {
			t.Errorf("Type %v should not be callable", val.Type())
		}
	}
}

// ============================================================================
// IsCountable Tests
// ============================================================================

func TestIsCountable_Array(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))

	result := IsCountable(types.NewArray(arr))
	if !result.ToBool() {
		t.Error("Array should be countable")
	}
}

func TestIsCountable_Object(t *testing.T) {
	// Object with count method is countable
	class := types.NewClassEntry("CountableClass")
	class.Methods["count"] = &types.MethodDef{
		Name:       "count",
		Visibility: types.VisibilityPublic,
	}
	obj := types.NewObjectFromClass(class)

	result := IsCountable(types.NewObject(obj))
	if !result.ToBool() {
		t.Error("Object with count method should be countable")
	}
}

func TestIsCountable_ObjectNoCount(t *testing.T) {
	// Object without count method is not countable
	class := types.NewClassEntry("NonCountableClass")
	obj := types.NewObjectFromClass(class)

	result := IsCountable(types.NewObject(obj))
	if result.ToBool() {
		t.Error("Object without count method should not be countable")
	}
}

func TestIsCountable_Other(t *testing.T) {
	// Other types are not countable
	tests := []*types.Value{
		types.NewInt(42),
		types.NewFloat(3.14),
		types.NewString("hello"),
		types.NewBool(true),
		types.NewNull(),
	}

	for _, val := range tests {
		result := IsCountable(val)
		if result.ToBool() {
			t.Errorf("Type %v should not be countable", val.Type())
		}
	}
}

// ============================================================================
// Additional var_dump, print_r, var_export Edge Cases
// ============================================================================

func TestVarDump_NestedArray(t *testing.T) {
	inner := types.NewEmptyArray()
	inner.Append(types.NewInt(1))
	inner.Append(types.NewInt(2))

	outer := types.NewEmptyArray()
	outer.Append(types.NewArray(inner))
	outer.Append(types.NewInt(3))

	result := VarDump(types.NewArray(outer))
	if result.Type() != types.TypeNull {
		t.Error("VarDump should return NULL")
	}
}

func TestVarDump_Object(t *testing.T) {
	class := types.NewClassEntry("TestClass")
	class.Properties["name"] = &types.PropertyDef{
		Name:       "name",
		Visibility: types.VisibilityPublic,
	}
	obj := types.NewObjectFromClass(class)
	obj.Properties["name"] = &types.Property{
		Value:      types.NewString("test"),
		Visibility: types.VisibilityPublic,
	}

	result := VarDump(types.NewObject(obj))
	if result.Type() != types.TypeNull {
		t.Error("VarDump should return NULL")
	}
}

func TestPrintR_NestedArray(t *testing.T) {
	inner := types.NewEmptyArray()
	inner.Append(types.NewInt(1))

	outer := types.NewEmptyArray()
	outer.Append(types.NewArray(inner))

	result := PrintR(types.NewArray(outer), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Error("PrintR with return should return string")
	}
}

func TestPrintR_Object(t *testing.T) {
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)

	result := PrintR(types.NewObject(obj), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Error("PrintR with return should return string")
	}
}

func TestVarExport_NestedArray(t *testing.T) {
	inner := types.NewEmptyArray()
	inner.Append(types.NewInt(1))

	outer := types.NewEmptyArray()
	outer.Append(types.NewArray(inner))

	result := VarExport(types.NewArray(outer), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Error("VarExport with return should return string")
	}

	output := result.ToString()
	if !strings.Contains(output, "array") {
		t.Error("VarExport should contain 'array' for nested arrays")
	}
}

func TestVarExport_Object(t *testing.T) {
	class := types.NewClassEntry("TestClass")
	obj := types.NewObjectFromClass(class)

	result := VarExport(types.NewObject(obj), types.NewBool(true))
	if result.Type() != types.TypeString {
		t.Error("VarExport with return should return string")
	}
}

func TestGetType_UnknownType(t *testing.T) {
	// Test that GetType handles all types
	res := types.NewResourceHandle("file", "data")
	result := GetType(types.NewResource(res))

	if result.ToString() != "resource" {
		t.Errorf("GetType(resource) = %v, want 'resource'", result.ToString())
	}
}
