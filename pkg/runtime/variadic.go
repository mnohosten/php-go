package runtime

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Variadic Function Support
// ============================================================================

// FunctionArguments stores the arguments passed to a function call
// This is used by func_get_args(), func_num_args(), and func_get_arg()
type FunctionArguments struct {
	// All arguments passed to the function (including named and unpacked)
	Args []*types.Value

	// Number of named parameters in the function signature
	NumNamedParams int

	// Whether the function has a variadic parameter
	HasVariadic bool
}

// NewFunctionArguments creates a new function arguments tracker
func NewFunctionArguments(args []*types.Value, numNamedParams int, hasVariadic bool) *FunctionArguments {
	return &FunctionArguments{
		Args:           args,
		NumNamedParams: numNamedParams,
		HasVariadic:    hasVariadic,
	}
}

// GetNumArgs returns the number of arguments passed to the function
// Implements func_num_args()
func (fa *FunctionArguments) GetNumArgs() int {
	if fa == nil {
		return 0
	}
	return len(fa.Args)
}

// GetArgs returns all arguments as an array
// Implements func_get_args()
func (fa *FunctionArguments) GetArgs() *types.Value {
	if fa == nil {
		return types.NewArray(types.NewArrayWithCapacity(0))
	}

	arr := types.NewArrayWithCapacity(len(fa.Args))
	for i, arg := range fa.Args {
		arr.Set(types.NewInt(int64(i)), arg)
	}

	return types.NewArray(arr)
}

// GetArg returns a specific argument by index
// Implements func_get_arg(int $position)
func (fa *FunctionArguments) GetArg(index int) (*types.Value, error) {
	if fa == nil {
		return nil, fmt.Errorf("func_get_arg(): Called from the global scope - no function context")
	}

	if index < 0 {
		return nil, fmt.Errorf("func_get_arg(): Argument number must be non-negative")
	}

	if index >= len(fa.Args) {
		return nil, fmt.Errorf("func_get_arg(): Argument %d not passed to function", index)
	}

	return fa.Args[index], nil
}

// GetVariadicArgs returns only the variadic arguments (after named parameters)
// This is useful for functions that need to process only the extra arguments
func (fa *FunctionArguments) GetVariadicArgs() []*types.Value {
	if fa == nil || !fa.HasVariadic {
		return nil
	}

	if len(fa.Args) <= fa.NumNamedParams {
		return nil
	}

	return fa.Args[fa.NumNamedParams:]
}

// GetVariadicArray returns variadic arguments as a PHP array
func (fa *FunctionArguments) GetVariadicArray() *types.Value {
	variadicArgs := fa.GetVariadicArgs()
	if variadicArgs == nil {
		return types.NewArray(types.NewArrayWithCapacity(0))
	}

	arr := types.NewArrayWithCapacity(len(variadicArgs))
	for i, arg := range variadicArgs {
		arr.Set(types.NewInt(int64(i)), arg)
	}

	return types.NewArray(arr)
}

// ============================================================================
// Argument Unpacking Support
// ============================================================================

// UnpackArgument unpacks an array or traversable into individual arguments
// Implements the ... operator in function calls: foo(...$array)
func UnpackArgument(value *types.Value) ([]*types.Value, error) {
	if value == nil {
		return nil, fmt.Errorf("Cannot unpack null value")
	}

	switch value.Type() {
	case types.TypeArray:
		arr := value.ToArray()
		if arr == nil {
			return nil, fmt.Errorf("Cannot unpack array")
		}

		// Get all values from the array
		result := make([]*types.Value, 0, arr.Len())
		arr.Each(func(key, val *types.Value) bool {
			result = append(result, val)
			return true
		})
		return result, nil

	default:
		return nil, fmt.Errorf("Only arrays and Traversables can be unpacked")
	}
}

// ExpandArguments processes a list of arguments, expanding any that are marked for unpacking
// This is called by the compiler when it encounters unpacked arguments
func ExpandArguments(args []*types.Value, unpackFlags []bool) ([]*types.Value, error) {
	if len(args) != len(unpackFlags) {
		return nil, fmt.Errorf("Argument count mismatch in unpacking")
	}

	result := make([]*types.Value, 0, len(args))

	for i, arg := range args {
		if unpackFlags[i] {
			// Unpack this argument
			unpacked, err := UnpackArgument(arg)
			if err != nil {
				return nil, err
			}
			result = append(result, unpacked...)
		} else {
			// Regular argument
			result = append(result, arg)
		}
	}

	return result, nil
}

// ============================================================================
// Variadic Parameter Handling
// ============================================================================

// CollectVariadicArgs collects variadic arguments into an array
// This is called when a function with a variadic parameter is invoked
//
// For example: function foo($a, $b, ...$rest)
// When called as: foo(1, 2, 3, 4, 5)
// The variadic collector will:
// - Bind 1 to $a
// - Bind 2 to $b
// - Collect [3, 4, 5] into $rest
func CollectVariadicArgs(args []*types.Value, numNamedParams int) (*types.Value, []*types.Value) {
	if len(args) <= numNamedParams {
		// No variadic arguments provided
		emptyArray := types.NewArray(types.NewArrayWithCapacity(0))
		return emptyArray, args
	}

	// Split arguments into named and variadic
	namedArgs := args[:numNamedParams]
	variadicArgs := args[numNamedParams:]

	// Collect variadic args into an array
	arr := types.NewArrayWithCapacity(len(variadicArgs))
	for i, arg := range variadicArgs {
		arr.Set(types.NewInt(int64(i)), arg)
	}

	return types.NewArray(arr), namedArgs
}

// ============================================================================
// Validation
// ============================================================================

// ValidateVariadicCall validates that a variadic function call is legal
// Rules:
// 1. Variadic parameter must be the last parameter
// 2. Unpacked arguments can appear anywhere but are expanded inline
// 3. Named arguments cannot be unpacked
func ValidateVariadicCall(hasVariadic bool, namedArgIndex int, unpackIndex int) error {
	// Named arguments cannot be combined with unpacking in the same call
	// (PHP allows this, but it's complex - we can add support later if needed)
	if namedArgIndex >= 0 && unpackIndex >= 0 && unpackIndex >= namedArgIndex {
		return fmt.Errorf("Cannot use argument unpacking with named arguments")
	}

	return nil
}
