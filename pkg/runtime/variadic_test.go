package runtime

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestFunctionArguments_Creation tests creating function arguments tracker
func TestFunctionArguments_Creation(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
	}

	fa := NewFunctionArguments(args, 2, true)

	if fa == nil {
		t.Fatal("NewFunctionArguments() returned nil")
	}

	if len(fa.Args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(fa.Args))
	}

	if fa.NumNamedParams != 2 {
		t.Errorf("Expected NumNamedParams = 2, got %d", fa.NumNamedParams)
	}

	if !fa.HasVariadic {
		t.Error("Expected HasVariadic = true")
	}
}

// TestFunctionArguments_GetNumArgs tests func_num_args()
func TestFunctionArguments_GetNumArgs(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
	}

	fa := NewFunctionArguments(args, 0, false)
	count := fa.GetNumArgs()

	if count != 3 {
		t.Errorf("GetNumArgs() = %d, want 3", count)
	}
}

// TestFunctionArguments_GetNumArgs_Nil tests nil function arguments
func TestFunctionArguments_GetNumArgs_Nil(t *testing.T) {
	var fa *FunctionArguments
	count := fa.GetNumArgs()

	if count != 0 {
		t.Errorf("GetNumArgs() on nil = %d, want 0", count)
	}
}

// TestFunctionArguments_GetArgs tests func_get_args()
func TestFunctionArguments_GetArgs(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewString("hello"),
	}

	fa := NewFunctionArguments(args, 0, false)
	result := fa.GetArgs()

	if result.Type() != types.TypeArray {
		t.Fatalf("GetArgs() type = %v, want TypeArray", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() != 3 {
		t.Fatalf("GetArgs() array length = %d, want 3", arr.Len())
	}

	// Check first argument
	val0, _ := arr.Get(types.NewInt(0))
	if val0.ToInt() != 1 {
		t.Errorf("args[0] = %d, want 1", val0.ToInt())
	}

	// Check second argument
	val1, _ := arr.Get(types.NewInt(1))
	if val1.ToInt() != 2 {
		t.Errorf("args[1] = %d, want 2", val1.ToInt())
	}

	// Check third argument
	val2, _ := arr.Get(types.NewInt(2))
	if val2.ToString() != "hello" {
		t.Errorf("args[2] = %s, want hello", val2.ToString())
	}
}

// TestFunctionArguments_GetArgs_Nil tests nil function arguments
func TestFunctionArguments_GetArgs_Nil(t *testing.T) {
	var fa *FunctionArguments
	result := fa.GetArgs()

	if result.Type() != types.TypeArray {
		t.Fatalf("GetArgs() on nil type = %v, want TypeArray", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() != 0 {
		t.Errorf("GetArgs() on nil length = %d, want 0", arr.Len())
	}
}

// TestFunctionArguments_GetArg tests func_get_arg()
func TestFunctionArguments_GetArg(t *testing.T) {
	args := []*types.Value{
		types.NewInt(10),
		types.NewInt(20),
		types.NewInt(30),
	}

	fa := NewFunctionArguments(args, 0, false)

	// Get first argument
	arg0, err := fa.GetArg(0)
	if err != nil {
		t.Fatalf("GetArg(0) error = %v", err)
	}
	if arg0.ToInt() != 10 {
		t.Errorf("GetArg(0) = %d, want 10", arg0.ToInt())
	}

	// Get second argument
	arg1, err := fa.GetArg(1)
	if err != nil {
		t.Fatalf("GetArg(1) error = %v", err)
	}
	if arg1.ToInt() != 20 {
		t.Errorf("GetArg(1) = %d, want 20", arg1.ToInt())
	}

	// Get third argument
	arg2, err := fa.GetArg(2)
	if err != nil {
		t.Fatalf("GetArg(2) error = %v", err)
	}
	if arg2.ToInt() != 30 {
		t.Errorf("GetArg(2) = %d, want 30", arg2.ToInt())
	}
}

// TestFunctionArguments_GetArg_Error_Negative tests negative index
func TestFunctionArguments_GetArg_Error_Negative(t *testing.T) {
	args := []*types.Value{types.NewInt(1)}
	fa := NewFunctionArguments(args, 0, false)

	_, err := fa.GetArg(-1)
	if err == nil {
		t.Error("GetArg(-1) should return error")
	}
}

// TestFunctionArguments_GetArg_Error_OutOfBounds tests out of bounds index
func TestFunctionArguments_GetArg_Error_OutOfBounds(t *testing.T) {
	args := []*types.Value{types.NewInt(1)}
	fa := NewFunctionArguments(args, 0, false)

	_, err := fa.GetArg(5)
	if err == nil {
		t.Error("GetArg(5) should return error when only 1 arg exists")
	}
}

// TestFunctionArguments_GetArg_Error_Nil tests nil function arguments
func TestFunctionArguments_GetArg_Error_Nil(t *testing.T) {
	var fa *FunctionArguments

	_, err := fa.GetArg(0)
	if err == nil {
		t.Error("GetArg(0) on nil should return error")
	}
}

// TestFunctionArguments_GetVariadicArgs tests extracting variadic arguments
func TestFunctionArguments_GetVariadicArgs(t *testing.T) {
	// function foo($a, $b, ...$rest)
	// Called as: foo(1, 2, 3, 4, 5)
	args := []*types.Value{
		types.NewInt(1), // $a
		types.NewInt(2), // $b
		types.NewInt(3), // $rest[0]
		types.NewInt(4), // $rest[1]
		types.NewInt(5), // $rest[2]
	}

	fa := NewFunctionArguments(args, 2, true)
	variadicArgs := fa.GetVariadicArgs()

	if len(variadicArgs) != 3 {
		t.Fatalf("Expected 3 variadic args, got %d", len(variadicArgs))
	}

	if variadicArgs[0].ToInt() != 3 {
		t.Errorf("variadicArgs[0] = %d, want 3", variadicArgs[0].ToInt())
	}
	if variadicArgs[1].ToInt() != 4 {
		t.Errorf("variadicArgs[1] = %d, want 4", variadicArgs[1].ToInt())
	}
	if variadicArgs[2].ToInt() != 5 {
		t.Errorf("variadicArgs[2] = %d, want 5", variadicArgs[2].ToInt())
	}
}

// TestFunctionArguments_GetVariadicArgs_NoExtra tests when no variadic args provided
func TestFunctionArguments_GetVariadicArgs_NoExtra(t *testing.T) {
	// function foo($a, $b, ...$rest)
	// Called as: foo(1, 2)
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
	}

	fa := NewFunctionArguments(args, 2, true)
	variadicArgs := fa.GetVariadicArgs()

	if variadicArgs != nil {
		t.Errorf("Expected nil variadic args, got %d args", len(variadicArgs))
	}
}

// TestFunctionArguments_GetVariadicArray tests variadic args as array
func TestFunctionArguments_GetVariadicArray(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
		types.NewInt(4),
	}

	fa := NewFunctionArguments(args, 1, true)
	result := fa.GetVariadicArray()

	if result.Type() != types.TypeArray {
		t.Fatalf("GetVariadicArray() type = %v, want TypeArray", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() != 3 {
		t.Fatalf("GetVariadicArray() length = %d, want 3", arr.Len())
	}
}

// TestUnpackArgument_Array tests unpacking an array
func TestUnpackArgument_Array(t *testing.T) {
	arr := types.NewArrayWithCapacity(3)
	arr.Set(types.NewInt(0), types.NewInt(10))
	arr.Set(types.NewInt(1), types.NewInt(20))
	arr.Set(types.NewInt(2), types.NewInt(30))

	value := types.NewArray(arr)

	unpacked, err := UnpackArgument(value)
	if err != nil {
		t.Fatalf("UnpackArgument() error = %v", err)
	}

	if len(unpacked) != 3 {
		t.Fatalf("Expected 3 unpacked values, got %d", len(unpacked))
	}

	if unpacked[0].ToInt() != 10 {
		t.Errorf("unpacked[0] = %d, want 10", unpacked[0].ToInt())
	}
	if unpacked[1].ToInt() != 20 {
		t.Errorf("unpacked[1] = %d, want 20", unpacked[1].ToInt())
	}
	if unpacked[2].ToInt() != 30 {
		t.Errorf("unpacked[2] = %d, want 30", unpacked[2].ToInt())
	}
}

// TestUnpackArgument_EmptyArray tests unpacking empty array
func TestUnpackArgument_EmptyArray(t *testing.T) {
	arr := types.NewArrayWithCapacity(0)
	value := types.NewArray(arr)

	unpacked, err := UnpackArgument(value)
	if err != nil {
		t.Fatalf("UnpackArgument() error = %v", err)
	}

	if len(unpacked) != 0 {
		t.Errorf("Expected 0 unpacked values from empty array, got %d", len(unpacked))
	}
}

// TestUnpackArgument_Error_Null tests unpacking null
func TestUnpackArgument_Error_Null(t *testing.T) {
	_, err := UnpackArgument(nil)
	if err == nil {
		t.Error("UnpackArgument(nil) should return error")
	}
}

// TestUnpackArgument_Error_NonArray tests unpacking non-array
func TestUnpackArgument_Error_NonArray(t *testing.T) {
	value := types.NewInt(42)

	_, err := UnpackArgument(value)
	if err == nil {
		t.Error("UnpackArgument(int) should return error")
	}
}

// TestExpandArguments tests expanding arguments with unpacking
func TestExpandArguments(t *testing.T) {
	// foo(1, ...[2, 3], 4)
	arr := types.NewArrayWithCapacity(2)
	arr.Set(types.NewInt(0), types.NewInt(2))
	arr.Set(types.NewInt(1), types.NewInt(3))

	args := []*types.Value{
		types.NewInt(1),
		types.NewArray(arr),
		types.NewInt(4),
	}

	unpackFlags := []bool{false, true, false}

	expanded, err := ExpandArguments(args, unpackFlags)
	if err != nil {
		t.Fatalf("ExpandArguments() error = %v", err)
	}

	if len(expanded) != 4 {
		t.Fatalf("Expected 4 expanded args, got %d", len(expanded))
	}

	if expanded[0].ToInt() != 1 {
		t.Errorf("expanded[0] = %d, want 1", expanded[0].ToInt())
	}
	if expanded[1].ToInt() != 2 {
		t.Errorf("expanded[1] = %d, want 2", expanded[1].ToInt())
	}
	if expanded[2].ToInt() != 3 {
		t.Errorf("expanded[2] = %d, want 3", expanded[2].ToInt())
	}
	if expanded[3].ToInt() != 4 {
		t.Errorf("expanded[3] = %d, want 4", expanded[3].ToInt())
	}
}

// TestExpandArguments_NoUnpacking tests when no arguments need unpacking
func TestExpandArguments_NoUnpacking(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
	}

	unpackFlags := []bool{false, false, false}

	expanded, err := ExpandArguments(args, unpackFlags)
	if err != nil {
		t.Fatalf("ExpandArguments() error = %v", err)
	}

	if len(expanded) != 3 {
		t.Fatalf("Expected 3 expanded args, got %d", len(expanded))
	}

	// Should be unchanged
	for i := range args {
		if expanded[i] != args[i] {
			t.Errorf("Argument %d was modified", i)
		}
	}
}

// TestExpandArguments_Error_Mismatch tests mismatched counts
func TestExpandArguments_Error_Mismatch(t *testing.T) {
	args := []*types.Value{types.NewInt(1), types.NewInt(2)}
	unpackFlags := []bool{false} // Wrong length

	_, err := ExpandArguments(args, unpackFlags)
	if err == nil {
		t.Error("ExpandArguments() with mismatched counts should return error")
	}
}

// TestCollectVariadicArgs tests collecting variadic arguments
func TestCollectVariadicArgs(t *testing.T) {
	// function foo($a, $b, ...$rest)
	// Called as: foo(1, 2, 3, 4, 5)
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
		types.NewInt(3),
		types.NewInt(4),
		types.NewInt(5),
	}

	variadicArray, namedArgs := CollectVariadicArgs(args, 2)

	// Check named args
	if len(namedArgs) != 2 {
		t.Fatalf("Expected 2 named args, got %d", len(namedArgs))
	}
	if namedArgs[0].ToInt() != 1 {
		t.Errorf("namedArgs[0] = %d, want 1", namedArgs[0].ToInt())
	}
	if namedArgs[1].ToInt() != 2 {
		t.Errorf("namedArgs[1] = %d, want 2", namedArgs[1].ToInt())
	}

	// Check variadic array
	if variadicArray.Type() != types.TypeArray {
		t.Fatalf("variadicArray type = %v, want TypeArray", variadicArray.Type())
	}

	arr := variadicArray.ToArray()
	if arr.Len() != 3 {
		t.Fatalf("variadicArray length = %d, want 3", arr.Len())
	}

	val0, _ := arr.Get(types.NewInt(0))
	if val0.ToInt() != 3 {
		t.Errorf("variadicArray[0] = %d, want 3", val0.ToInt())
	}
}

// TestCollectVariadicArgs_NoExtra tests when no extra args provided
func TestCollectVariadicArgs_NoExtra(t *testing.T) {
	args := []*types.Value{
		types.NewInt(1),
		types.NewInt(2),
	}

	variadicArray, namedArgs := CollectVariadicArgs(args, 2)

	// Check named args
	if len(namedArgs) != 2 {
		t.Fatalf("Expected 2 named args, got %d", len(namedArgs))
	}

	// Check variadic array (should be empty)
	if variadicArray.Type() != types.TypeArray {
		t.Fatalf("variadicArray type = %v, want TypeArray", variadicArray.Type())
	}

	arr := variadicArray.ToArray()
	if arr.Len() != 0 {
		t.Errorf("variadicArray length = %d, want 0", arr.Len())
	}
}

// TestValidateVariadicCall tests validation
func TestValidateVariadicCall(t *testing.T) {
	tests := []struct {
		name           string
		hasVariadic    bool
		namedArgIndex  int
		unpackIndex    int
		expectError    bool
	}{
		{"no named or unpack", true, -1, -1, false},
		{"only named", false, 2, -1, false},
		{"only unpack", true, -1, 1, false},
		{"unpack before named", false, 1, 0, false},
		{"named before unpack", false, 0, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariadicCall(tt.hasVariadic, tt.namedArgIndex, tt.unpackIndex)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
