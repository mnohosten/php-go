# FFI (Foreign Function Interface) Usage Guide

This guide explains how to use the PHP-Go FFI system to call Go functions from PHP code.

## Table of Contents

1. [Overview](#overview)
2. [Basic Usage](#basic-usage)
3. [Function Registration](#function-registration)
4. [Calling Functions](#calling-functions)
5. [Advanced Patterns](#advanced-patterns)
6. [Error Handling](#error-handling)
7. [Performance](#performance)

## Overview

The FFI system enables PHP code to call Go functions through the `go_call()` function. This provides:

- Direct access to Go's performance
- Use of Go's extensive standard library
- Type-safe function calls with automatic marshaling
- Support for variadic functions
- Error propagation from Go to PHP

## Basic Usage

### From PHP Side

```php
<?php
// Call a Go function
$result = go_call('function.name', $arg1, $arg2, ...);

// Check if function exists
if (go_has_function('math.Add')) {
    $sum = go_call('math.Add', 10, 20);
}

// Get function information
$info = go_function_info('math.Add');
print_r($info);
// Array (
//     [name] => math.Add
//     [num_params] => 2
//     [num_returns] => 1
//     [variadic] => false
// )

// List all available functions
$functions = go_list_functions();
```

### From Go Side

```go
// Register a function
registry := goext.NewFunctionRegistry()
registry.Register("math.Add", func(a, b int64) int64 {
    return a + b
})

// Make it available to PHP via FFI
ffi := goext.NewFFIManager(registry)
```

## Function Registration

### Simple Functions

```go
registry := goext.NewFunctionRegistry()

// No arguments
registry.Register("time.Now", func() int64 {
    return time.Now().Unix()
})

// Single argument
registry.Register("string.Upper", func(s string) string {
    return strings.ToUpper(s)
})

// Multiple arguments
registry.Register("math.Add", func(a, b int64) int64 {
    return a + b
})
```

### Functions with Errors

```go
// Return (value, error)
registry.Register("file.Read", func(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(data), nil
})

// PHP usage:
// $content = go_call('file.Read', '/path/to/file.txt');
```

### Variadic Functions

```go
// Variadic parameters
registry.Register("string.Join", func(sep string, parts ...string) string {
    return strings.Join(parts, sep)
})

// PHP usage:
// $result = go_call('string.Join', '-', 'a', 'b', 'c');
// Result: "a-b-c"
```

### Multiple Return Values

```go
// Multiple returns (non-error)
registry.Register("math.DivMod", func(a, b int64) (int64, int64) {
    return a / b, a % b
})

// Returns array in PHP:
// $result = go_call('math.DivMod', 17, 5);
// $result is [3, 2]
```

### No Return Value

```go
// Void functions
registry.Register("logger.Log", func(message string) {
    log.Println(message)
})

// PHP usage:
// go_call('logger.Log', 'Hello from PHP');
// Returns null
```

## Calling Functions

### Type Conversions

The FFI system automatically converts between PHP and Go types:

```php
<?php
// Integers
$sum = go_call('math.Add', 10, 20);  // int64

// Floats
$sqrt = go_call('math.Sqrt', 16.0);  // float64

// Strings
$upper = go_call('string.Upper', 'hello');  // string

// Booleans
$valid = go_call('string.IsValid', 'test');  // bool

// Arrays (becomes slice in Go)
$joined = go_call('array.Join', [1, 2, 3]);

// Associative arrays (becomes map in Go)
$data = go_call('json.Encode', ['name' => 'John', 'age' => 30]);
```

### Working with Arrays

```go
// Slice parameters
registry.Register("array.Sum", func(nums []int64) int64 {
    sum := int64(0)
    for _, n := range nums {
        sum += n
    }
    return sum
})
```

```php
<?php
$total = go_call('array.Sum', [1, 2, 3, 4, 5]);  // 15
```

### Working with Maps

```go
// Map parameters
registry.Register("map.Get", func(data map[string]interface{}, key string) interface{} {
    return data[key]
})
```

```php
<?php
$value = go_call('map.Get', ['name' => 'John', 'age' => 30], 'name');  // "John"
```

## Advanced Patterns

### Stateful Functions

```go
type Counter struct {
    mu    sync.Mutex
    count int64
}

func (c *Counter) Increment() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    return c.count
}

func (c *Counter) Get() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

// Register
counter := &Counter{}
registry.Register("counter.Increment", counter.Increment)
registry.Register("counter.Get", counter.Get)
```

```php
<?php
$count1 = go_call('counter.Increment');  // 1
$count2 = go_call('counter.Increment');  // 2
$current = go_call('counter.Get');       // 2
```

### Callbacks (Future Feature)

```go
// Currently not supported, but planned:
registry.Register("array.Map", func(arr []int64, callback func(int64) int64) []int64 {
    result := make([]int64, len(arr))
    for i, v := range arr {
        result[i] = callback(v)
    }
    return result
})
```

### Context and Cancellation

```go
// Access VM context in functions
registry.Register("task.Run", func(vm *vm.VM, args []*types.Value) (*types.Value, error) {
    // Can access VM state if needed
    taskName := args[0].ToString()

    // Perform task
    result := performTask(taskName)

    return types.NewString(result), nil
})
```

## Error Handling

### From Go to PHP

```go
registry.Register("divide", func(a, b int64) (int64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
})
```

```php
<?php
// In PHP, check for errors
$result = go_call('divide', 10, 2);  // 5

$result = go_call('divide', 10, 0);  // Triggers PHP error/exception
```

### Custom Error Types

```go
type ValidationError struct {
    Field string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

registry.Register("validate.Email", func(email string) (bool, error) {
    if !strings.Contains(email, "@") {
        return false, &ValidationError{
            Field: "email",
            Message: "must contain @",
        }
    }
    return true, nil
})
```

### Panic Recovery

The FFI system automatically recovers from panics:

```go
registry.Register("unsafe.Panic", func() string {
    panic("something went wrong")
})
```

```php
<?php
// This will return an error instead of crashing
$result = go_call('unsafe.Panic');  // Error: panic: something went wrong
```

## Performance

### Benchmarking

```go
func BenchmarkFFICall(b *testing.B) {
    registry := goext.NewFunctionRegistry()
    registry.Register("add", func(a, b int64) int64 {
        return a + b
    })

    ffi := goext.NewFFIManager(registry)
    vm := vm.New()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ffi.GoCall(vm, []*types.Value{
            types.NewString("add"),
            types.NewInt(10),
            types.NewInt(20),
        })
    }
}
```

### Optimization Tips

1. **Minimize Marshaling**: Keep data in Go as long as possible
```go
// Bad: Convert back and forth
registry.Register("process1", func(data string) string { ... })
registry.Register("process2", func(data string) string { ... })

// Good: Process in one call
registry.Register("processAll", func(data string) string {
    step1 := process1(data)
    step2 := process2(step1)
    return step2
})
```

2. **Batch Operations**: Process multiple items at once
```go
// Bad: Multiple calls from PHP
for ($i = 0; $i < 1000; $i++) {
    go_call('process.Item', $items[$i]);
}

// Good: Single batch call
go_call('process.Batch', $items);
```

3. **Cache Function Lookups**: The registry caches function lookups internally

4. **Use Appropriate Types**: Avoid interface{} when possible
```go
// Better: Specific types
func add(a, b int64) int64 { ... }

// Slower: Generic types
func add(a, b interface{}) interface{} { ... }
```

### Performance Characteristics

Typical overhead per FFI call:
- Function lookup: ~100ns (cached)
- Type marshaling: ~500ns per argument
- Function execution: depends on function
- Result marshaling: ~500ns

For most use cases, FFI overhead is negligible compared to the work performed.

## Complete Example

Here's a complete example showing FFI usage:

```go
// Go side - math_ext.go
package mathext

import (
    "errors"
    "math"
    "github.com/krizos/php-go/pkg/goext"
)

func RegisterMathFunctions(registry *goext.FunctionRegistry) {
    registry.Register("math.Add", Add)
    registry.Register("math.Multiply", Multiply)
    registry.Register("math.Pow", Pow)
    registry.Register("math.Sqrt", Sqrt)
    registry.Register("math.Sum", Sum)
}

func Add(a, b int64) int64 {
    return a + b
}

func Multiply(a, b int64) int64 {
    return a * b
}

func Pow(base, exp float64) float64 {
    return math.Pow(base, exp)
}

func Sqrt(x float64) (float64, error) {
    if x < 0 {
        return 0, errors.New("cannot sqrt negative number")
    }
    return math.Sqrt(x), nil
}

func Sum(numbers ...int64) int64 {
    total := int64(0)
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```php
<?php
// PHP side - math_example.php

// Basic arithmetic
$sum = go_call('math.Add', 10, 20);
echo "10 + 20 = $sum\n";  // 30

$product = go_call('math.Multiply', 5, 6);
echo "5 * 6 = $product\n";  // 30

// Floating point
$power = go_call('math.Pow', 2.0, 10.0);
echo "2^10 = $power\n";  // 1024

// With error handling
$sqrt = go_call('math.Sqrt', 16.0);
echo "√16 = $sqrt\n";  // 4

// Variadic function
$total = go_call('math.Sum', 1, 2, 3, 4, 5);
echo "Sum = $total\n";  // 15
```

## Next Steps

- See [Extension Development Guide](README.md) for creating extensions
- Check [Type Marshaling Guide](marshaling-guide.md) for complex types
- Review [examples](examples/) for more patterns

## Troubleshooting

### Function Not Found
```
Error: go_call(): function 'xxx' not found
```
- Verify function is registered: `go_has_function('xxx')`
- Check function name spelling
- Ensure registry is loaded into VM

### Type Mismatch
```
Error: expected int64, got string
```
- Check argument types match Go function signature
- Use explicit type conversion in PHP
- Review [Type Marshaling Guide](marshaling-guide.md)

### Performance Issues
- Profile with `go test -bench`
- Reduce marshaling overhead
- Batch operations when possible
- Consider creating an extension instead of using FFI
