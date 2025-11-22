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
