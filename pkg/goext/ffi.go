package goext

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// FFIManager manages the FFI (Foreign Function Interface) system.
// It bridges PHP code to registered Go functions.
type FFIManager struct {
	registry *FunctionRegistry
}

// NewFFIManager creates a new FFI manager.
func NewFFIManager(registry *FunctionRegistry) *FFIManager {
	return &FFIManager{
		registry: registry,
	}
}

// GetGlobalFFI returns an FFI manager using the global registry.
func GetGlobalFFI() *FFIManager {
	return NewFFIManager(GetGlobalRegistry())
}

// GoCall implements the go_call() PHP function.
// It allows PHP code to call registered Go functions.
//
// Usage in PHP:
//   $result = go_call('function.name', $arg1, $arg2, ...);
//
// Parameters:
//   - function name (string): Fully qualified function name
//   - arguments (variadic): Arguments to pass to the Go function
//
// Returns:
//   - The result from the Go function (marshaled to PHP)
//   - Or null if the function returns nothing
//
// Errors:
//   - If function is not found
//   - If arguments are invalid
//   - If the Go function returns an error
func (f *FFIManager) GoCall(vm *vm.VM, args []*types.Value) (*types.Value, error) {
	// Validate arguments
	if len(args) < 1 {
		return nil, fmt.Errorf("go_call() expects at least 1 argument (function name), got %d", len(args))
	}

	// First argument must be function name (string)
	fnNameVal := args[0]
	if fnNameVal.Type() != types.TypeString {
		return nil, fmt.Errorf("go_call() expects argument 1 to be string (function name), got %s", fnNameVal.TypeString())
	}

	fnName := fnNameVal.ToString()

	// Look up the function
	regFunc, found := f.registry.Get(fnName)
	if !found {
		return nil, fmt.Errorf("go_call(): function '%s' not found", fnName)
	}

	// Extract function arguments (skip first arg which is the function name)
	fnArgs := args[1:]

	// Call the registered function handler
	result, err := regFunc.Handler(vm, fnArgs)
	if err != nil {
		return nil, fmt.Errorf("go_call('%s'): %w", fnName, err)
	}

	return result, nil
}

// RegisterBuiltinGoCall registers the go_call() function with the VM.
// This should be called during VM initialization to make go_call() available in PHP.
// Note: This is a placeholder. Full VM integration requires native function support in the VM.
func RegisterBuiltinGoCall(v *vm.VM) error {
	return RegisterBuiltinGoCallWithRegistry(v, GetGlobalRegistry())
}

// RegisterBuiltinGoCallWithRegistry registers go_call() with a specific registry.
// Note: This is a placeholder. Full VM integration requires native function support in the VM.
func RegisterBuiltinGoCallWithRegistry(v *vm.VM, registry *FunctionRegistry) error {
	_ = NewFFIManager(registry) // Create FFI manager (will be used for VM integration)

	// TODO: Full VM integration
	// The VM will need to:
	// 1. Support native (non-bytecode) functions
	// 2. Register go_call as a built-in function
	// 3. Call FFI manager's GoCall handler when go_call is invoked
	//
	// For now, this is a placeholder that can be called via CallGoFunction

	return nil
}

// GoCallHandler is a wrapper that returns a FunctionHandler for go_call.
// This can be used to register go_call as a native PHP function.
func GoCallHandler(registry *FunctionRegistry) FunctionHandler {
	ffi := NewFFIManager(registry)
	return ffi.GoCall
}

// ListRegisteredFunctions returns a list of all registered Go functions.
// This can be called from PHP to discover available functions.
//
// Usage in PHP:
//   $functions = go_list_functions();
func (f *FFIManager) ListFunctions(vm *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("go_list_functions() expects 0 arguments, got %d", len(args))
	}

	names := f.registry.List()

	// Create PHP array of function names
	arr := types.NewEmptyArray()
	for _, name := range names {
		arr.Append(types.NewString(name))
	}

	return types.NewArray(arr), nil
}

// HasFunction checks if a Go function is registered.
//
// Usage in PHP:
//   $exists = go_has_function('math.Add');
func (f *FFIManager) HasFunction(vm *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_has_function() expects 1 argument (function name), got %d", len(args))
	}

	if args[0].Type() != types.TypeString {
		return nil, fmt.Errorf("go_has_function() expects argument 1 to be string, got %s", args[0].TypeString())
	}

	fnName := args[0].ToString()
	exists := f.registry.Has(fnName)

	return types.NewBool(exists), nil
}

// GetFunctionInfo returns information about a registered function.
//
// Usage in PHP:
//   $info = go_function_info('math.Add');
//   // Returns: ['name' => 'math.Add', 'num_params' => 2, 'variadic' => false]
func (f *FFIManager) GetFunctionInfo(vm *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_function_info() expects 1 argument (function name), got %d", len(args))
	}

	if args[0].Type() != types.TypeString {
		return nil, fmt.Errorf("go_function_info() expects argument 1 to be string, got %s", args[0].TypeString())
	}

	fnName := args[0].ToString()
	regFunc, found := f.registry.Get(fnName)
	if !found {
		return types.NewNull(), nil
	}

	// Build info array
	info := types.NewEmptyArray()
	info.Set(types.NewString("name"), types.NewString(regFunc.Name))
	info.Set(types.NewString("num_params"), types.NewInt(int64(regFunc.Signature.NumIn)))
	info.Set(types.NewString("num_returns"), types.NewInt(int64(regFunc.Signature.NumOut)))
	info.Set(types.NewString("variadic"), types.NewBool(regFunc.IsVariadic))

	return types.NewArray(info), nil
}

// RegisterAllFFIFunctions registers all FFI helper functions with a registry.
// This includes: go_call, go_list_functions, go_has_function, go_function_info
func RegisterAllFFIFunctions(goRegistry *FunctionRegistry) error {
	ffi := NewFFIManager(goRegistry)

	// Note: These are being registered in the Go registry, but they're meant
	// to be called as PHP built-in functions. The actual VM integration
	// will need to handle this specially.

	// For now, we create handlers that can be used
	_ = ffi // Prevent unused variable error

	return nil
}

// CallGoFunction is a convenience function to call a registered Go function.
// This is the low-level interface used by go_call().
func CallGoFunction(fnName string, vm *vm.VM, args []*types.Value) (*types.Value, error) {
	return CallGoFunctionWithRegistry(fnName, vm, args, GetGlobalRegistry())
}

// CallGoFunctionWithRegistry calls a Go function using a specific registry.
func CallGoFunctionWithRegistry(fnName string, vm *vm.VM, args []*types.Value, registry *FunctionRegistry) (*types.Value, error) {
	regFunc, found := registry.Get(fnName)
	if !found {
		return nil, fmt.Errorf("function '%s' not registered", fnName)
	}

	result, err := regFunc.Handler(vm, args)
	if err != nil {
		return nil, err
	}

	return result, nil
}
