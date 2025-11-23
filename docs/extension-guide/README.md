# PHP-Go Extension Development Guide

This guide shows you how to create PHP extensions in Go using the PHP-Go extension system.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Extension Basics](#extension-basics)
3. [Type Marshaling](#type-marshaling)
4. [FFI System](#ffi-system)
5. [Best Practices](#best-practices)
6. [Examples](#examples)

## Quick Start

### Creating Your First Extension

```go
package myextension

import (
    "github.com/krizos/php-go/pkg/goext"
    "github.com/krizos/php-go/pkg/types"
    "github.com/krizos/php-go/pkg/vm"
)

// Create extension using BaseExtension
func NewMyExtension() *goext.BaseExtension {
    ext := goext.NewBaseExtension("my_extension", "1.0.0")

    // Add a function
    ext.AddFunction("my_func", myFunction)

    // Add a constant
    ext.AddConstant("MY_CONST", types.NewInt(42))

    return ext
}

// Function handler
func myFunction(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("my_func() expects 1 argument")
    }

    input := args[0].ToString()
    return types.NewString("Hello, " + input), nil
}
```

### Registering and Loading

```go
// Create extension manager
registry := goext.NewFunctionRegistry()
manager := goext.NewExtensionManager(registry)

// Register extension
ext := NewMyExtension()
manager.Register(ext)

// Load into VM
vm := vm.New()
manager.LoadIntoVM("my_extension", vm)

// Now the extension functions are available in the VM
```

## Extension Basics

### Extension Interface

All extensions must implement the `Extension` interface:

```go
type Extension interface {
    Name() string                              // Extension name
    Version() string                           // Version string
    Init(*vm.VM) error                        // Called when loaded
    Functions() map[string]FunctionHandler     // PHP functions
    Constants() map[string]*types.Value        // PHP constants
}
```

### BaseExtension Helper

For simple extensions, use `BaseExtension`:

```go
ext := goext.NewBaseExtension("my_ext", "1.0.0")
ext.AddFunction("func_name", handler)
ext.AddConstant("CONST_NAME", value)
```

### Custom Extension

For complex extensions, implement the interface directly:

```go
type MyExtension struct {
    config Config
}

func (e *MyExtension) Name() string {
    return "my_extension"
}

func (e *MyExtension) Version() string {
    return "1.0.0"
}

func (e *MyExtension) Init(v *vm.VM) error {
    // Perform initialization
    return nil
}

func (e *MyExtension) Functions() map[string]FunctionHandler {
    return map[string]FunctionHandler{
        "my_func": e.myFunction,
    }
}

func (e *MyExtension) Constants() map[string]*types.Value {
    return map[string]*types.Value{
        "MY_CONST": types.NewInt(42),
    }
}
```

## Type Marshaling

### PHP to Go Conversion

```go
m := goext.NewMarshaler()

// PHP → Go
phpVal := types.NewInt(42)
goVal, err := m.ToGo(phpVal)
// goVal is int64(42)

// Type-specific conversion
phpStr := types.NewString("hello")
str, err := m.ToGoString(phpStr)
// str is "hello"
```

### Go to PHP Conversion

```go
m := goext.NewMarshaler()

// Go → PHP
phpVal, err := m.ToPHP(int64(42))
// phpVal is types.Value with TypeInt

// Slices become PHP arrays
phpArr, err := m.ToPHP([]int64{1, 2, 3})
// phpArr is PHP array [1, 2, 3]

// Maps become associative arrays
phpMap, err := m.ToPHP(map[string]int{"a": 1, "b": 2})
// phpMap is PHP array ["a" => 1, "b" => 2]
```

### Type Mapping

| Go Type | PHP Type |
|---------|----------|
| nil | null |
| bool | bool |
| int, int64 | int |
| float64 | float |
| string | string |
| []interface{} | array (list) |
| map[string]interface{} | array (assoc) |
| struct | array or object |

### Custom Type Marshaling

```go
// Implement CustomMarshaler interface
type MyType struct {
    Value int
}

func (m *MyType) MarshalPHP() (*types.Value, error) {
    return types.NewInt(int64(m.Value)), nil
}

// Or register a custom converter
am := goext.NewAdvancedMarshaler()
am.RegisterConverter(reflect.TypeOf(MyType{}), &goext.TypeConverter{
    ToPHP: func(v interface{}) (*types.Value, error) {
        t := v.(MyType)
        return types.NewInt(int64(t.Value)), nil
    },
    ToGo: func(v *types.Value) (interface{}, error) {
        return MyType{Value: int(v.ToInt())}, nil
    },
})
```

## FFI System

### Registering Go Functions

```go
registry := goext.NewFunctionRegistry()

// Simple function
registry.Register("math.Add", func(a, b int64) int64 {
    return a + b
})

// Variadic function
registry.Register("string.Concat", func(strs ...string) string {
    return strings.Join(strs, "")
})

// Function with error
registry.Register("file.Read", func(path string) (string, error) {
    data, err := os.ReadFile(path)
    return string(data), err
})
```

### Calling from PHP

In PHP, use `go_call()` to invoke registered functions:

```php
<?php
// Simple call
$result = go_call('math.Add', 10, 20);  // 30

// Variadic
$joined = go_call('string.Concat', 'hello', ' ', 'world');  // "hello world"

// Error handling
$content = go_call('file.Read', 'file.txt');
```

### FFI Helper Functions

```php
<?php
// Check if function exists
$exists = go_has_function('math.Add');  // true

// Get function info
$info = go_function_info('math.Add');
// Returns: ['name' => 'math.Add', 'num_params' => 2, ...]

// List all functions
$functions = go_list_functions();
// Returns array of function names
```

## Best Practices

### 1. Error Handling

Always return errors as the second return value:

```go
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) < 1 {
        return nil, fmt.Errorf("myFunc() requires 1 argument, got %d", len(args))
    }

    // Your logic here
    return result, nil
}
```

### 2. Argument Validation

Validate arguments at the start:

```go
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
    }

    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf("argument 1 must be string, got %s",
            args[0].TypeString())
    }

    // Use the arguments
    str := args[0].ToString()
    num := args[1].ToInt()

    return result, nil
}
```

### 3. Thread Safety

Use mutexes for shared state:

```go
type MyExtension struct {
    *goext.BaseExtension
    mu    sync.RWMutex
    cache map[string]string
}

func (e *MyExtension) cachedFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    key := args[0].ToString()

    e.mu.RLock()
    val, ok := e.cache[key]
    e.mu.RUnlock()

    if ok {
        return types.NewString(val), nil
    }

    // Compute value
    result := computeValue(key)

    e.mu.Lock()
    e.cache[key] = result
    e.mu.Unlock()

    return types.NewString(result), nil
}
```

### 4. Resource Cleanup

Use Init() for setup and cleanup:

```go
type DBExtension struct {
    *goext.BaseExtension
    db *sql.DB
}

func (e *DBExtension) Init(v *vm.VM) error {
    var err error
    e.db, err = sql.Open("postgres", connectionString)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    // Register cleanup handler
    runtime.SetFinalizer(e, func(e *DBExtension) {
        if e.db != nil {
            e.db.Close()
        }
    })

    return nil
}
```

### 5. Documentation

Document your functions:

```go
// greet returns a greeting message for the given name.
//
// Usage in PHP:
//   $message = my_greet("John");  // "Hello, John!"
//
// Parameters:
//   - name (string): The name to greet
//
// Returns:
//   - string: The greeting message
func greet(v *vm.VM, args []*types.Value) (*types.Value, error) {
    name := args[0].ToString()
    return types.NewString("Hello, " + name + "!"), nil
}
```

## Examples

### Example 1: Math Extension

```go
func NewMathExtension() *goext.BaseExtension {
    ext := goext.NewBaseExtension("math_extra", "1.0.0")

    ext.AddFunction("factorial", factorial)
    ext.AddFunction("fibonacci", fibonacci)
    ext.AddConstant("GOLDEN_RATIO", types.NewFloat(1.618033988749))

    return ext
}

func factorial(v *vm.VM, args []*types.Value) (*types.Value, error) {
    n := args[0].ToInt()
    result := int64(1)
    for i := int64(2); i <= n; i++ {
        result *= i
    }
    return types.NewInt(result), nil
}
```

### Example 2: HTTP Client Extension

See `pkg/goext/bindings/http.go` for a complete HTTP client implementation.

### Example 3: Caching Extension

```go
type CacheExtension struct {
    *goext.BaseExtension
    mu    sync.RWMutex
    cache map[string]*types.Value
}

func NewCacheExtension() *CacheExtension {
    ext := &CacheExtension{
        BaseExtension: goext.NewBaseExtension("cache", "1.0.0"),
        cache:        make(map[string]*types.Value),
    }

    ext.AddFunction("cache_set", ext.set)
    ext.AddFunction("cache_get", ext.get)
    ext.AddFunction("cache_has", ext.has)
    ext.AddFunction("cache_clear", ext.clear)

    return ext
}

func (e *CacheExtension) set(v *vm.VM, args []*types.Value) (*types.Value, error) {
    key := args[0].ToString()
    value := args[1]

    e.mu.Lock()
    e.cache[key] = value
    e.mu.Unlock()

    return types.NewBool(true), nil
}

func (e *CacheExtension) get(v *vm.VM, args []*types.Value) (*types.Value, error) {
    key := args[0].ToString()

    e.mu.RLock()
    value, ok := e.cache[key]
    e.mu.RUnlock()

    if !ok {
        return types.NewNull(), nil
    }

    return value, nil
}
```

## Testing Your Extension

```go
func TestMyExtension(t *testing.T) {
    // Create extension
    ext := NewMyExtension()

    // Create VM and manager
    testVM := vm.New()
    registry := goext.NewFunctionRegistry()
    manager := goext.NewExtensionManager(registry)

    // Register and load
    if err := manager.Register(ext); err != nil {
        t.Fatalf("Failed to register: %v", err)
    }

    if err := manager.LoadIntoVM("my_extension", testVM); err != nil {
        t.Fatalf("Failed to load: %v", err)
    }

    // Test function
    funcs := ext.Functions()
    result, err := funcs["my_func"](testVM, []*types.Value{
        types.NewString("World"),
    })

    if err != nil {
        t.Fatalf("Function call failed: %v", err)
    }

    if result.ToString() != "Hello, World" {
        t.Errorf("Expected 'Hello, World', got '%s'", result.ToString())
    }
}
```

## Next Steps

- Review the [FFI Usage Guide](ffi-guide.md) for advanced FFI patterns
- Check [Type Marshaling Guide](marshaling-guide.md) for complex type conversions
- See [Example Extensions](examples/) for more complete examples
- Read the [API Reference](../api-reference.md) for detailed API documentation

## Support

For questions or issues:
- File an issue on GitHub
- Check existing extensions in `pkg/goext/bindings/`
- Review example tests in `pkg/goext/example_test.go`
