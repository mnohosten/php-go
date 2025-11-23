package goext

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// FunctionHandler is the signature for Go functions callable from PHP.
// It receives the VM context and PHP arguments, returns a PHP value and error.
type FunctionHandler func(*vm.VM, []*types.Value) (*types.Value, error)

// RegisteredFunction holds metadata about a registered Go function.
type RegisteredFunction struct {
	Name        string          // Fully qualified name (e.g., "math.Add")
	Handler     FunctionHandler // The wrapped handler function
	GoFunc      interface{}     // Original Go function
	Signature   *FuncSignature  // Parsed function signature
	IsVariadic  bool            // Whether function accepts variadic arguments
}

// FuncSignature describes a Go function's signature.
type FuncSignature struct {
	Name       string
	NumIn      int           // Number of input parameters
	NumOut     int           // Number of return values
	InTypes    []reflect.Type // Input parameter types
	OutTypes   []reflect.Type // Output return types
	IsVariadic bool          // Variadic function
}

// FunctionRegistry manages registered Go functions.
type FunctionRegistry struct {
	mu        sync.RWMutex
	functions map[string]*RegisteredFunction
	marshaler *Marshaler
}

// Global function registry
var (
	globalRegistry     *FunctionRegistry
	globalRegistryOnce sync.Once
)

// GetGlobalRegistry returns the singleton function registry.
func GetGlobalRegistry() *FunctionRegistry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewFunctionRegistry()
	})
	return globalRegistry
}

// NewFunctionRegistry creates a new function registry.
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]*RegisteredFunction),
		marshaler: NewMarshaler(),
	}
}

// RegisterFunction registers a Go function to be callable from PHP.
// The function can have various signatures, and will be automatically wrapped.
//
// Supported function signatures:
//   - func()
//   - func() T
//   - func() (T, error)
//   - func(T1) T2
//   - func(T1, T2) (T3, error)
//   - func(T1, T2, ...TN) (TR, error)
//
// Where T can be: int, int64, float64, bool, string, []interface{}, map[string]interface{}
//
// Examples:
//   RegisterFunction("math.Add", func(a, b int64) int64 { return a + b })
//   RegisterFunction("string.Upper", func(s string) string { return strings.ToUpper(s) })
func RegisterFunction(name string, fn interface{}) error {
	return GetGlobalRegistry().Register(name, fn)
}

// Register registers a function in this registry.
func (r *FunctionRegistry) Register(name string, fn interface{}) error {
	if name == "" {
		return fmt.Errorf("function name cannot be empty")
	}
	if fn == nil {
		return fmt.Errorf("function cannot be nil")
	}

	// Parse function signature
	sig, err := parseFunctionSignature(name, fn)
	if err != nil {
		return fmt.Errorf("parsing function signature: %w", err)
	}

	// Validate signature
	if err := validateSignature(sig); err != nil {
		return fmt.Errorf("invalid function signature: %w", err)
	}

	// Wrap the function
	handler := r.wrapFunction(fn, sig)

	// Store registered function
	regFunc := &RegisteredFunction{
		Name:       name,
		Handler:    handler,
		GoFunc:     fn,
		Signature:  sig,
		IsVariadic: sig.IsVariadic,
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.functions[name]; exists {
		return fmt.Errorf("function '%s' already registered", name)
	}

	r.functions[name] = regFunc
	return nil
}

// Get retrieves a registered function by name.
func (r *FunctionRegistry) Get(name string) (*RegisteredFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.functions[name]
	return fn, ok
}

// Has checks if a function is registered.
func (r *FunctionRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.functions[name]
	return ok
}

// List returns all registered function names.
func (r *FunctionRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered functions.
func (r *FunctionRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.functions)
}

// Clear removes all registered functions.
func (r *FunctionRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.functions = make(map[string]*RegisteredFunction)
}

// parseFunctionSignature uses reflection to parse a function's signature.
func parseFunctionSignature(name string, fn interface{}) (*FuncSignature, error) {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function: %T", fn)
	}

	sig := &FuncSignature{
		Name:       name,
		NumIn:      fnType.NumIn(),
		NumOut:     fnType.NumOut(),
		InTypes:    make([]reflect.Type, fnType.NumIn()),
		OutTypes:   make([]reflect.Type, fnType.NumOut()),
		IsVariadic: fnType.IsVariadic(),
	}

	// Parse input parameters
	for i := 0; i < fnType.NumIn(); i++ {
		sig.InTypes[i] = fnType.In(i)
	}

	// Parse output parameters
	for i := 0; i < fnType.NumOut(); i++ {
		sig.OutTypes[i] = fnType.Out(i)
	}

	return sig, nil
}

// validateSignature validates that a function signature is supported.
func validateSignature(sig *FuncSignature) error {
	// Check output: must be (), (T), or (T, error)
	if sig.NumOut > 2 {
		return fmt.Errorf("too many return values: %d (max 2)", sig.NumOut)
	}

	// If 2 return values, second must be error
	if sig.NumOut == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !sig.OutTypes[1].Implements(errorType) {
			return fmt.Errorf("second return value must be error, got %s", sig.OutTypes[1])
		}
	}

	// Check that input/output types are marshallable
	for i, t := range sig.InTypes {
		if !isMarshallableType(t) {
			return fmt.Errorf("parameter %d: unsupported type %s", i, t)
		}
	}

	if sig.NumOut > 0 {
		if !isMarshallableType(sig.OutTypes[0]) {
			return fmt.Errorf("return value: unsupported type %s", sig.OutTypes[0])
		}
	}

	return nil
}

// isMarshallableType checks if a type can be marshaled between PHP and Go.
func isMarshallableType(t reflect.Type) bool {
	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true

	case reflect.Slice, reflect.Array:
		// Check element type
		return isMarshallableType(t.Elem())

	case reflect.Map:
		// Maps must have string keys
		if t.Key().Kind() != reflect.String {
			return false
		}
		return isMarshallableType(t.Elem())

	case reflect.Struct:
		// Structs are marshallable (converted to PHP arrays)
		return true

	case reflect.Interface:
		// interface{} is allowed
		return t.NumMethod() == 0

	default:
		return false
	}
}

// wrapFunction wraps a Go function to be callable from PHP.
// It handles argument marshaling, function invocation, and result marshaling.
func (r *FunctionRegistry) wrapFunction(fn interface{}, sig *FuncSignature) FunctionHandler {
	fnValue := reflect.ValueOf(fn)

	return func(vm *vm.VM, args []*types.Value) (*types.Value, error) {
		// Determine expected argument count
		minArgs := sig.NumIn
		if sig.IsVariadic {
			minArgs-- // Variadic functions require at least numIn-1 args
		}

		// Validate argument count
		if !sig.IsVariadic && len(args) != sig.NumIn {
			return nil, fmt.Errorf("%s() expects %d arguments, got %d", sig.Name, sig.NumIn, len(args))
		}

		if sig.IsVariadic && len(args) < minArgs {
			return nil, fmt.Errorf("%s() expects at least %d arguments, got %d", sig.Name, minArgs, len(args))
		}

		// Build arguments for the Go function
		goArgs := make([]reflect.Value, 0)

		if sig.IsVariadic {
			// For variadic functions: marshal regular args, then build slice for variadic args
			numRegular := sig.NumIn - 1 // All params except the variadic one

			// Marshal regular (non-variadic) parameters
			for i := 0; i < numRegular; i++ {
				goVal, err := r.marshalToGo(args[i], sig.InTypes[i])
				if err != nil {
					return nil, fmt.Errorf("argument %d: %w", i, err)
				}
				goArgs = append(goArgs, goVal)
			}

			// Marshal variadic parameters into a slice
			varArgs := args[numRegular:] // Remaining PHP args
			elemType := sig.InTypes[numRegular].Elem() // Element type of the variadic slice

			// Build variadic args slice
			sliceType := sig.InTypes[numRegular]
			sliceVal := reflect.MakeSlice(sliceType, len(varArgs), len(varArgs))

			for j, phpArg := range varArgs {
				goVal, err := r.marshalToGo(phpArg, elemType)
				if err != nil {
					return nil, fmt.Errorf("variadic argument %d: %w", j, err)
				}
				sliceVal.Index(j).Set(goVal)
			}

			// Add the variadic slice as final argument
			goArgs = append(goArgs, sliceVal)
		} else {
			// Non-variadic: marshal all args normally
			for i := 0; i < sig.NumIn; i++ {
				goVal, err := r.marshalToGo(args[i], sig.InTypes[i])
				if err != nil {
					return nil, fmt.Errorf("argument %d: %w", i, err)
				}
				goArgs = append(goArgs, goVal)
			}
		}

		// Call the function
		var results []reflect.Value
		if sig.IsVariadic {
			// Variadic functions must use CallSlice
			results = fnValue.CallSlice(goArgs)
		} else {
			results = fnValue.Call(goArgs)
		}

		// Handle return values
		switch sig.NumOut {
		case 0:
			// No return value
			return types.NewNull(), nil

		case 1:
			// Single return value
			phpVal, err := r.marshaler.ToPHP(results[0].Interface())
			if err != nil {
				return nil, fmt.Errorf("marshaling return value: %w", err)
			}
			return phpVal, nil

		case 2:
			// (T, error) return
			// Check if error is non-nil
			if !results[1].IsNil() {
				err := results[1].Interface().(error)
				return nil, err
			}
			// Marshal the value
			phpVal, err := r.marshaler.ToPHP(results[0].Interface())
			if err != nil {
				return nil, fmt.Errorf("marshaling return value: %w", err)
			}
			return phpVal, nil

		default:
			return nil, fmt.Errorf("unexpected number of return values: %d", sig.NumOut)
		}
	}
}

// marshalToGo converts a PHP value to a Go value of the specified type.
func (r *FunctionRegistry) marshalToGo(phpVal *types.Value, targetType reflect.Type) (reflect.Value, error) {
	// Handle pointer types
	if targetType.Kind() == reflect.Ptr {
		// Create pointer to value
		goVal, err := r.marshalToGo(phpVal, targetType.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptr := reflect.New(targetType.Elem())
		ptr.Elem().Set(goVal)
		return ptr, nil
	}

	// Convert PHP value to Go interface{}
	goIface, err := r.marshaler.ToGo(phpVal)
	if err != nil {
		return reflect.Value{}, err
	}

	// Handle nil
	if goIface == nil {
		return reflect.Zero(targetType), nil
	}

	// Convert to target type
	goVal := reflect.ValueOf(goIface)

	// Try direct assignment first
	if goVal.Type().AssignableTo(targetType) {
		return goVal, nil
	}

	// Try conversion
	if goVal.Type().ConvertibleTo(targetType) {
		return goVal.Convert(targetType), nil
	}

	// Handle interface{} target
	if targetType.Kind() == reflect.Interface && targetType.NumMethod() == 0 {
		return goVal, nil
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", goVal.Type(), targetType)
}
