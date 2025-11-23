package vm

import "github.com/krizos/php-go/pkg/types"

// BuiltinFunction represents a native Go function that implements a PHP built-in
type BuiltinFunction func(args []*types.Value) (*types.Value, error)

// builtinFunctions maps function names to their implementations
var builtinFunctions = map[string]BuiltinFunction{
	"strlen": builtinStrlen,
}

// IsBuiltin checks if a function name is a built-in function
func IsBuiltin(name string) bool {
	_, exists := builtinFunctions[name]
	return exists
}

// CallBuiltin calls a built-in function with the given arguments
func CallBuiltin(name string, args []*types.Value) (*types.Value, error) {
	fn, exists := builtinFunctions[name]
	if !exists {
		return nil, nil // Not a built-in
	}
	return fn(args)
}

// ============================================================================
// String Functions
// ============================================================================

// builtinStrlen implements strlen(string $string): int
func builtinStrlen(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		// PHP strlen() requires exactly 1 argument
		return types.NewInt(0), nil
	}

	str := args[0].ToString()
	return types.NewInt(int64(len(str))), nil
}
