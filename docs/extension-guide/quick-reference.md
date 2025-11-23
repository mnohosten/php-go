# Quick Reference Guide

Quick reference for common PHP-Go extension and FFI operations.

## Extension Creation

```go
// Simple extension
ext := goext.NewBaseExtension("my_ext", "1.0.0")
ext.AddFunction("my_func", handler)
ext.AddConstant("MY_CONST", types.NewInt(42))

// Register and load
manager := goext.NewExtensionManager(goext.NewFunctionRegistry())
manager.Register(ext)
manager.LoadIntoVM("my_ext", vm)
```

## Function Handlers

```go
// Basic handler
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    input := args[0].ToString()
    return types.NewString("Hello, " + input), nil
}

// With validation
func validated(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("expected 2 arguments")
    }
    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf("arg 1 must be string")
    }
    // ... process
    return result, nil
}
```

## FFI Registration

```go
registry := goext.NewFunctionRegistry()

// Simple function
registry.Register("math.Add", func(a, b int64) int64 {
    return a + b
})

// With error
registry.Register("file.Read", func(path string) (string, error) {
    return os.ReadFile(path)
})

// Variadic
registry.Register("string.Join", func(sep string, parts ...string) string {
    return strings.Join(parts, sep)
})
```

## Type Marshaling

```go
m := goext.NewMarshaler()

// PHP → Go
goVal, err := m.ToGo(phpVal)
str, err := m.ToGoString(phpVal)
num, err := m.ToGoInt(phpVal)

// Go → PHP
phpVal, err := m.ToPHP(goVal)
phpInt, err := m.ToPHP(int64(42))
phpStr, err := m.ToPHP("hello")
```

## Advanced Marshaling

```go
am := goext.NewAdvancedMarshaler()

// Struct → Object
phpObj, err := am.StructToObject(myStruct, "ClassName")

// Custom converter
am.RegisterConverter(reflect.TypeOf(MyType{}), &goext.TypeConverter{
    ToPHP: func(v interface{}) (*types.Value, error) { /* ... */ },
    ToGo:  func(v *types.Value) (interface{}, error) { /* ... */ },
})
```

## PHP Value Creation

```go
// Primitives
types.NewNull()
types.NewBool(true)
types.NewInt(42)
types.NewFloat(3.14)
types.NewString("hello")

// Arrays
arr := types.NewEmptyArray()
arr.Append(types.NewInt(1))
arr.Set(types.NewString("key"), types.NewString("value"))
phpArray := types.NewArray(arr)

// Objects
obj := types.NewObjectInstance("MyClass")
obj.SetProperty("prop", types.NewInt(42), nil)
phpObj := types.NewObject(obj)
```

## PHP Value Access

```go
// Type checking
if val.Type() == types.TypeString {
    str := val.ToString()
}

// Conversion
val.ToBool()    // → bool
val.ToInt()     // → int64
val.ToFloat()   // → float64
val.ToString()  // → string
val.ToArray()   // → *types.Array
val.ToObject()  // → *types.Object

// Type name
typeName := val.TypeString()
```

## Array Operations

```go
arr := types.NewEmptyArray()

// Add elements
arr.Append(value)                              // Append to end
arr.Set(types.NewString("key"), value)         // Set by key
arr.Push(val1, val2)                          // Push multiple

// Access elements
val, ok := arr.Get(types.NewString("key"))    // Get by key
val, ok := arr.Pop()                          // Remove from end
val, ok := arr.Shift()                        // Remove from start

// Info
arr.Len()                                      // Length
arr.HasKey(types.NewString("key"))            // Check key exists
arr.Keys()                                     // Get all keys
arr.Values()                                   // Get all values

// Iteration
arr.Each(func(key, val *types.Value) bool {
    // Process each element
    return true  // continue
})
```

## Common Patterns

### Error Handling
```go
func handler(v *vm.VM, args []*types.Value) (*types.Value, error) {
    result, err := doSomething()
    if err != nil {
        return nil, fmt.Errorf("operation failed: %w", err)
    }
    return result, nil
}
```

### Variadic Arguments
```go
func handler(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) == 0 {
        return nil, fmt.Errorf("expected at least 1 argument")
    }

    // Process all arguments
    for _, arg := range args {
        process(arg)
    }

    return result, nil
}
```

### Stateful Extension
```go
type MyExtension struct {
    *goext.BaseExtension
    mu    sync.RWMutex
    state map[string]string
}

func (e *MyExtension) setState(v *vm.VM, args []*types.Value) (*types.Value, error) {
    key := args[0].ToString()
    val := args[1].ToString()

    e.mu.Lock()
    e.state[key] = val
    e.mu.Unlock()

    return types.NewBool(true), nil
}
```

## PHP Usage Examples

```php
<?php
// FFI calls
$result = go_call('math.Add', 10, 20);
$exists = go_has_function('math.Add');
$info = go_function_info('math.Add');
$functions = go_list_functions();

// Extension constants
echo MY_CONST;  // 42

// Extension functions
$greeting = my_func("World");
```

## Testing

```go
func TestMyExtension(t *testing.T) {
    ext := NewMyExtension()
    testVM := vm.New()
    registry := goext.NewFunctionRegistry()
    manager := goext.NewExtensionManager(registry)

    manager.Register(ext)
    manager.LoadIntoVM("my_ext", testVM)

    // Test function
    funcs := ext.Functions()
    result, err := funcs["my_func"](testVM, []*types.Value{
        types.NewString("test"),
    })

    assert.NoError(t, err)
    assert.Equal(t, "expected", result.ToString())
}
```

## Type Mapping Table

| Go Type | PHP Type | Notes |
|---------|----------|-------|
| `nil` | `null` | |
| `bool` | `bool` | |
| `int`, `int64` | `int` | |
| `float32`, `float64` | `float` | |
| `string` | `string` | |
| `[]T` | `array` | List if sequential keys |
| `map[string]T` | `array` | Associative array |
| `struct` | `object` or `array` | Depends on conversion method |

## Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `function not found` | Not registered | Check `go_has_function()` |
| `type mismatch` | Wrong argument type | Validate with `.Type()` |
| `nil pointer` | Accessing nil value | Check for null first |
| `index out of range` | Array access | Check `.Len()` first |

## Links

- [Extension Development Guide](README.md)
- [FFI Usage Guide](ffi-guide.md)
- [Type Marshaling Guide](marshaling-guide.md)
- [Example Extensions](../goext/bindings/)
