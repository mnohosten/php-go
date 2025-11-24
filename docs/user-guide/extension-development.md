# Extension Development Guide

This comprehensive guide covers everything you need to know to develop PHP extensions using Go for the PHP-Go interpreter.

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Extension Architecture](#extension-architecture)
4. [Creating Extensions](#creating-extensions)
5. [Function Handlers](#function-handlers)
6. [Type System and Marshaling](#type-system-and-marshaling)
7. [FFI (Foreign Function Interface)](#ffi-foreign-function-interface)
8. [Plugin System](#plugin-system)
9. [Built-in Extensions](#built-in-extensions)
10. [Testing Extensions](#testing-extensions)
11. [Performance Optimization](#performance-optimization)
12. [Best Practices](#best-practices)
13. [Deployment](#deployment)
14. [Troubleshooting](#troubleshooting)

## Introduction

PHP-Go allows you to extend PHP functionality by writing extensions in Go. This provides several advantages:

- **Performance**: Go's compiled nature delivers high-performance extensions
- **Concurrency**: Leverage Go's goroutines for parallel operations
- **Ecosystem**: Access to Go's rich standard library and third-party packages
- **Type Safety**: Go's strong typing helps prevent runtime errors
- **Deployment**: Compile extensions as plugins or link them statically

### When to Create an Extension

Create an extension when you need to:

- Interface with Go libraries (databases, message queues, cloud services)
- Implement performance-critical operations
- Access system resources not available in PHP
- Provide concurrent/parallel processing capabilities
- Integrate with existing Go codebases

## Getting Started

### Prerequisites

- Go 1.21 or later
- PHP-Go source code
- Basic understanding of Go and PHP

### Quick Start

Here's a minimal extension:

```go
package main

import (
    "github.com/krizos/php-go/pkg/goext"
    "github.com/krizos/php-go/pkg/types"
    "github.com/krizos/php-go/pkg/vm"
)

func NewHelloExtension() *goext.BaseExtension {
    ext := goext.NewBaseExtension("hello", "1.0.0")

    ext.AddFunction("hello_world", func(v *vm.VM, args []*types.Value) (*types.Value, error) {
        return types.NewString("Hello from Go!"), nil
    })

    ext.AddConstant("HELLO_VERSION", types.NewString("1.0.0"))

    return ext
}
```

### Building Your First Extension

1. Create a new Go module:
```bash
mkdir my-extension
cd my-extension
go mod init example.com/my-extension
```

2. Add PHP-Go dependency:
```bash
go get github.com/krizos/php-go
```

3. Create your extension (see Quick Start example above)

4. Build as a plugin:
```bash
go build -buildmode=plugin -o hello.so
```

5. Load in PHP:
```php
<?php
// PHP-Go will load the plugin
$result = hello_world();  // "Hello from Go!"
echo HELLO_VERSION;        // "1.0.0"
```

## Extension Architecture

### Core Components

PHP-Go's extension system consists of several key components:

#### 1. Extension Interface

All extensions must implement the `Extension` interface:

```go
type Extension interface {
    Name() string                              // Unique extension name
    Version() string                           // Semantic version
    Init(*vm.VM) error                        // Initialization hook
    Functions() map[string]FunctionHandler     // PHP functions provided
    Constants() map[string]*types.Value        // PHP constants provided
}
```

#### 2. Extension Manager

The `ExtensionManager` handles extension lifecycle:

- **Registration**: Register extensions with the system
- **Loading**: Load extensions into a VM instance
- **Dependency Management**: Handle inter-extension dependencies
- **Initialization**: Call Init() hooks in proper order

#### 3. Function Registry

The `FunctionRegistry` maintains a global registry of Go functions:

- Maps Go functions to PHP function names
- Handles automatic type conversion
- Supports reflection-based invocation
- Enables FFI (Foreign Function Interface)

#### 4. Type Marshaler

The marshaling system converts between Go and PHP types:

- **Basic Marshaler**: Simple type conversions
- **Advanced Marshaler**: Complex types, structs, custom conversions
- **Bidirectional**: Go ↔ PHP type mapping

### Extension Lifecycle

```
1. Create Extension
   ↓
2. Register with ExtensionManager
   ↓
3. Load into VM
   ↓
4. Init() called
   ↓
5. Functions/Constants available
   ↓
6. Extension used during execution
   ↓
7. VM shutdown (cleanup via finalizers)
```

## Creating Extensions

### Using BaseExtension

For simple extensions, use the `BaseExtension` helper:

```go
func NewMathExtension() *goext.BaseExtension {
    ext := goext.NewBaseExtension("math_extra", "1.0.0")

    // Add functions
    ext.AddFunction("factorial", factorial)
    ext.AddFunction("fibonacci", fibonacci)
    ext.AddFunction("is_prime", isPrime)

    // Add constants
    ext.AddConstant("MATH_E", types.NewFloat(2.718281828459045))
    ext.AddConstant("MATH_PHI", types.NewFloat(1.618033988749895))

    return ext
}

func factorial(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("factorial() expects 1 argument, got %d", len(args))
    }

    n := args[0].ToInt()
    if n < 0 {
        return nil, fmt.Errorf("factorial() expects non-negative integer")
    }

    result := int64(1)
    for i := int64(2); i <= n; i++ {
        result *= i
    }

    return types.NewInt(result), nil
}
```

### Custom Extension Implementation

For complex extensions with state, implement the interface directly:

```go
type DatabaseExtension struct {
    name    string
    version string
    db      *sql.DB
    mu      sync.RWMutex
    config  Config
}

func NewDatabaseExtension(config Config) *DatabaseExtension {
    return &DatabaseExtension{
        name:    "database",
        version: "1.0.0",
        config:  config,
    }
}

func (e *DatabaseExtension) Name() string {
    return e.name
}

func (e *DatabaseExtension) Version() string {
    return e.version
}

func (e *DatabaseExtension) Init(v *vm.VM) error {
    var err error
    e.db, err = sql.Open(e.config.Driver, e.config.DSN)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := e.db.Ping(); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Register cleanup
    runtime.SetFinalizer(e, func(e *DatabaseExtension) {
        if e.db != nil {
            e.db.Close()
        }
    })

    return nil
}

func (e *DatabaseExtension) Functions() map[string]goext.FunctionHandler {
    return map[string]goext.FunctionHandler{
        "db_query":   e.query,
        "db_execute": e.execute,
        "db_close":   e.close,
    }
}

func (e *DatabaseExtension) Constants() map[string]*types.Value {
    return map[string]*types.Value{
        "DB_VERSION": types.NewString(e.version),
    }
}
```

### Extension with Configuration

Extensions can accept configuration:

```go
type Config struct {
    CacheSize int
    Timeout   time.Duration
    Debug     bool
}

func NewCacheExtension(config Config) *CacheExtension {
    return &CacheExtension{
        config: config,
        cache:  make(map[string]*CacheEntry, config.CacheSize),
    }
}
```

## Function Handlers

### Function Handler Signature

All PHP function handlers must follow this signature:

```go
type FunctionHandler func(v *vm.VM, args []*types.Value) (*types.Value, error)
```

- **v**: The VM instance (provides access to globals, execution context)
- **args**: Array of PHP values passed as arguments
- **Returns**: PHP value result and optional error

### Argument Handling

#### Fixed Arguments

```go
func add(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("add() expects 2 arguments, got %d", len(args))
    }

    a := args[0].ToInt()
    b := args[1].ToInt()

    return types.NewInt(a + b), nil
}
```

#### Variable Arguments

```go
func sum(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) == 0 {
        return types.NewInt(0), nil
    }

    total := int64(0)
    for _, arg := range args {
        total += arg.ToInt()
    }

    return types.NewInt(total), nil
}
```

#### Optional Arguments

```go
func greet(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) < 1 || len(args) > 2 {
        return nil, fmt.Errorf("greet() expects 1 or 2 arguments, got %d", len(args))
    }

    name := args[0].ToString()
    prefix := "Hello"

    if len(args) == 2 {
        prefix = args[1].ToString()
    }

    return types.NewString(fmt.Sprintf("%s, %s!", prefix, name)), nil
}
```

### Type Checking

Always validate argument types:

```go
func processString(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("processString() expects 1 argument")
    }

    // Type check
    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf(
            "processString() expects string, got %s",
            args[0].TypeString(),
        )
    }

    str := args[0].ToString()
    // Process string...

    return types.NewString(result), nil
}
```

### Error Handling

Return meaningful errors:

```go
func divide(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("divide() expects 2 arguments, got %d", len(args))
    }

    a := args[0].ToFloat()
    b := args[1].ToFloat()

    if b == 0 {
        return nil, fmt.Errorf("division by zero")
    }

    return types.NewFloat(a / b), nil
}
```

### Accessing VM State

The VM parameter provides access to execution context:

```go
func getCurrentFile(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Access current execution frame
    frame := v.CurrentFrame()
    if frame == nil {
        return types.NewString("unknown"), nil
    }

    // Get file information from frame
    // (Actual implementation depends on VM internals)

    return types.NewString(filename), nil
}
```

### Working with Arrays

```go
func arraySum(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("array_sum() expects 1 argument")
    }

    if args[0].Type() != types.TypeArray {
        return nil, fmt.Errorf("array_sum() expects array")
    }

    arr := args[0].Array()
    sum := int64(0)

    for _, elem := range arr.Elements() {
        sum += elem.ToInt()
    }

    return types.NewInt(sum), nil
}
```

### Working with Objects

```go
func getProperty(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("get_property() expects 2 arguments")
    }

    if args[0].Type() != types.TypeObject {
        return nil, fmt.Errorf("first argument must be object")
    }

    obj := args[0].Object()
    propName := args[1].ToString()

    val, exists := obj.GetProperty(propName)
    if !exists {
        return types.NewNull(), nil
    }

    return val, nil
}
```

## Type System and Marshaling

### PHP Value Types

PHP-Go supports all PHP types:

```go
types.NewNull()                  // NULL
types.NewBool(true)              // bool
types.NewInt(42)                 // int
types.NewFloat(3.14)             // float
types.NewString("hello")         // string
types.NewArray()                 // array
types.NewObject(class)           // object
types.NewResource(id)            // resource
types.NewReference(val)          // reference
```

### Type Conversion

PHP values provide conversion methods:

```go
val := args[0]

// Type checking
if val.Type() == types.TypeString {
    // Handle string
}

// Type conversion (follows PHP semantics)
str := val.ToString()      // Convert to string
num := val.ToInt()         // Convert to int
flt := val.ToFloat()       // Convert to float
bln := val.ToBool()        // Convert to bool
arr := val.ToArray()       // Convert to array
```

### Basic Marshaling

Use the `Marshaler` for simple Go ↔ PHP conversions:

```go
m := goext.NewMarshaler()

// Go → PHP
phpInt, _ := m.ToPHP(int64(42))
phpStr, _ := m.ToPHP("hello")
phpArr, _ := m.ToPHP([]int64{1, 2, 3})
phpMap, _ := m.ToPHP(map[string]interface{}{
    "name": "John",
    "age":  30,
})

// PHP → Go
goVal, _ := m.ToGo(phpInt)          // int64(42)
goStr, _ := m.ToGoString(phpStr)    // "hello"
goArr, _ := m.ToGoSlice(phpArr)     // []interface{}{1, 2, 3}
goMap, _ := m.ToGoMap(phpMap)       // map[string]interface{}
```

### Advanced Marshaling

Use `AdvancedMarshaler` for complex types:

```go
type User struct {
    ID        int64
    Name      string
    Email     string
    CreatedAt time.Time
}

am := goext.NewAdvancedMarshaler()

// Struct → PHP object
user := User{ID: 1, Name: "John", Email: "john@example.com"}
phpObj, err := am.StructToObject(&user)

// PHP object → Struct
var decoded User
err = am.ObjectToStruct(phpObj, &decoded)

// Struct → PHP array
phpArr, err := am.StructToArray(&user)

// PHP array → Struct
err = am.ArrayToStruct(phpArr, &decoded)
```

### Custom Type Marshaling

Implement custom marshaling for your types:

```go
type CustomType struct {
    Data []byte
}

// Implement CustomMarshaler interface
func (c *CustomType) MarshalPHP() (*types.Value, error) {
    return types.NewString(base64.StdEncoding.EncodeToString(c.Data)), nil
}

// Register custom converter
am := goext.NewAdvancedMarshaler()
am.RegisterConverter(reflect.TypeOf(CustomType{}), &goext.TypeConverter{
    ToPHP: func(v interface{}) (*types.Value, error) {
        c := v.(CustomType)
        encoded := base64.StdEncoding.EncodeToString(c.Data)
        return types.NewString(encoded), nil
    },
    ToGo: func(v *types.Value) (interface{}, error) {
        decoded, err := base64.StdEncoding.DecodeString(v.ToString())
        if err != nil {
            return nil, err
        }
        return CustomType{Data: decoded}, nil
    },
})
```

### Type Mapping Reference

| Go Type | PHP Type | Notes |
|---------|----------|-------|
| nil | NULL | |
| bool | bool | |
| int, int8, int16, int32, int64 | int | |
| uint, uint8, uint16, uint32, uint64 | int | Converted to signed |
| float32, float64 | float | |
| string | string | |
| []T | array | Indexed array |
| [N]T | array | Indexed array |
| map[string]T | array | Associative array |
| map[int]T | array | Indexed array |
| struct | array or object | Configurable |
| *T | T | Dereferenced |
| func | Not supported | Use FFI |
| chan | Not supported | |
| interface{} | Any type | Runtime type detection |

## FFI (Foreign Function Interface)

The FFI system allows PHP code to call arbitrary Go functions without manually creating function handlers.

### Registering Go Functions

```go
registry := goext.NewFunctionRegistry()

// Simple function
registry.Register("math.Add", func(a, b int64) int64 {
    return a + b
})

// With error return
registry.Register("file.Read", func(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(data), nil
})

// Variadic function
registry.Register("strings.Join", func(sep string, strs ...string) string {
    return strings.Join(strs, sep)
})

// Multiple return values
registry.Register("math.DivMod", func(a, b int64) (int64, int64, error) {
    if b == 0 {
        return 0, 0, fmt.Errorf("division by zero")
    }
    return a / b, a % b, nil
})
```

### Calling from PHP

```php
<?php
// Simple call
$sum = go_call('math.Add', 10, 20);  // 30

// With error handling
try {
    $content = go_call('file.Read', 'config.txt');
} catch (Exception $e) {
    echo "Error: " . $e->getMessage();
}

// Variadic
$joined = go_call('strings.Join', ', ', 'a', 'b', 'c');  // "a, b, c"

// Multiple returns (returned as array)
$result = go_call('math.DivMod', 17, 5);
echo $result[0];  // 3 (quotient)
echo $result[1];  // 2 (remainder)
```

### FFI Helper Functions

```php
<?php
// Check if function exists
if (go_has_function('math.Add')) {
    $result = go_call('math.Add', 1, 2);
}

// Get function information
$info = go_function_info('math.Add');
// Returns: [
//   'name' => 'math.Add',
//   'num_params' => 2,
//   'variadic' => false,
//   'returns_error' => false
// ]

// List all registered functions
$functions = go_list_functions();
// Returns: ['math.Add', 'file.Read', ...]
```

### Supported FFI Types

| Go Type | Supported | Notes |
|---------|-----------|-------|
| bool | ✓ | |
| int, int64 | ✓ | |
| float64 | ✓ | |
| string | ✓ | |
| []byte | ✓ | Converted to string |
| []T | ✓ | Converted to PHP array |
| map[string]T | ✓ | Converted to PHP array |
| struct | ✓ | Converted to PHP array/object |
| error | ✓ | Second/last return value |
| *T | ✓ | Dereferenced |
| func | ✗ | Not supported |
| chan | ✗ | Not supported |

### FFI Limitations

- **No callbacks**: Cannot pass PHP functions to Go
- **No channels**: Go channels not supported
- **No pointers**: Cannot pass pointers from PHP
- **Error handling**: Errors must be returned, not panicked

## Plugin System

The plugin system allows loading extensions dynamically at runtime.

### Building a Plugin

```bash
# Build as a shared library (.so on Linux/macOS, .dll on Windows)
go build -buildmode=plugin -o myext.so myext.go
```

### Plugin Structure

```go
package main

import (
    "github.com/krizos/php-go/pkg/goext"
)

// Extension is the required exported symbol
var Extension goext.Extension

func init() {
    Extension = NewMyExtension()
}

func NewMyExtension() *goext.BaseExtension {
    ext := goext.NewBaseExtension("my_plugin", "1.0.0")
    // ... configure extension
    return ext
}
```

### Loading Plugins

```go
// Create plugin manager
pm := goext.NewPluginManager()

// Load plugin
err := pm.LoadPlugin("/path/to/myext.so")
if err != nil {
    log.Fatalf("Failed to load plugin: %v", err)
}

// Get extension from plugin
ext, err := pm.GetExtension("my_plugin")
if err != nil {
    log.Fatalf("Failed to get extension: %v", err)
}

// Register with extension manager
manager := goext.NewExtensionManager(registry)
manager.Register(ext)

// Load into VM
vm := vm.New()
manager.LoadIntoVM("my_plugin", vm)
```

### Plugin Discovery

```go
// Load all plugins from directory
pluginDir := "/usr/local/lib/php-go/extensions"
plugins, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))

for _, path := range plugins {
    if err := pm.LoadPlugin(path); err != nil {
        log.Printf("Warning: Failed to load %s: %v", path, err)
    }
}
```

### Plugin Information

```go
// List loaded plugins
plugins := pm.ListPlugins()
for _, info := range plugins {
    fmt.Printf("Plugin: %s\n", info.Path)
    fmt.Printf("  Extension: %s\n", info.ExtensionName)
    fmt.Printf("  Loaded: %v\n", info.Loaded)
    if info.Error != nil {
        fmt.Printf("  Error: %v\n", info.Error)
    }
}

// Query by path
info, exists := pm.GetPluginInfo("/path/to/myext.so")

// Query by extension name
path, exists := pm.GetPluginPath("my_plugin")
```

### Plugin Limitations

- **No hot reloading**: Go's plugin system doesn't support unloading
- **Platform-specific**: Linux and macOS only (no Windows support)
- **Go version**: Must match PHP-Go's Go version
- **Dependencies**: Shared dependencies must be compatible

## Built-in Extensions

PHP-Go includes several built-in extensions as examples and for common functionality.

### Time Extension

High-performance time and date operations:

```php
<?php
// Current Unix timestamp with nanosecond precision
$time = go_time();  // 1234567890.123456789

// Sleep with sub-millisecond precision
go_sleep_ms(100);   // Sleep for 100 milliseconds
go_sleep_ns(1000);  // Sleep for 1000 nanoseconds

// Parse time
$parsed = go_time_parse("2006-01-02", "2023-11-24");

// Format time
$formatted = go_time_format($parsed, "Mon Jan 2 15:04:05");
```

### Filesystem Extension

Fast file operations:

```php
<?php
// Read entire file
$content = go_file_read('/path/to/file.txt');

// Write file
go_file_write('/path/to/file.txt', $content);

// Append to file
go_file_append('/path/to/file.txt', $additional);

// Check file existence
if (go_file_exists('/path/to/file.txt')) {
    // ...
}

// Get file info
$info = go_file_stat('/path/to/file.txt');
// Returns: ['size' => 1024, 'mode' => 0644, 'mtime' => 1234567890, ...]

// List directory
$files = go_dir_list('/path/to/directory');
```

### HTTP Extension

HTTP client with connection pooling:

```php
<?php
// Simple GET request
$response = go_http_get('https://api.example.com/data');
// Returns: ['status' => 200, 'body' => '...', 'headers' => [...]]

// POST request
$response = go_http_post('https://api.example.com/data', [
    'headers' => ['Content-Type' => 'application/json'],
    'body' => json_encode(['key' => 'value']),
]);

// Custom request
$response = go_http_request([
    'method' => 'PUT',
    'url' => 'https://api.example.com/resource/123',
    'headers' => ['Authorization' => 'Bearer token'],
    'body' => json_encode($data),
    'timeout' => 30,  // seconds
]);
```

### Crypto Extension

Cryptographic operations:

```php
<?php
// Hash functions (using Go's crypto packages)
$hash = go_hash('sha256', 'data');

// HMAC
$hmac = go_hmac('sha256', 'key', 'data');

// Random bytes (cryptographically secure)
$random = go_random_bytes(32);  // 32 random bytes
```

### JSON Extension

Fast JSON encoding/decoding:

```php
<?php
// Encode to JSON (using Go's encoding/json)
$json = go_json_encode(['key' => 'value']);

// Decode JSON
$data = go_json_decode($json);

// Pretty print
$pretty = go_json_encode($data, GO_JSON_PRETTY);
```

## Testing Extensions

### Unit Testing

```go
package myextension

import (
    "testing"
    "github.com/krizos/php-go/pkg/goext"
    "github.com/krizos/php-go/pkg/types"
    "github.com/krizos/php-go/pkg/vm"
)

func TestMyFunction(t *testing.T) {
    // Create extension
    ext := NewMyExtension()

    // Create test VM
    testVM := vm.New()

    // Create managers
    registry := goext.NewFunctionRegistry()
    manager := goext.NewExtensionManager(registry)

    // Register and load
    if err := manager.Register(ext); err != nil {
        t.Fatalf("Failed to register extension: %v", err)
    }

    if err := manager.LoadIntoVM(ext.Name(), testVM); err != nil {
        t.Fatalf("Failed to load extension: %v", err)
    }

    // Test function
    funcs := ext.Functions()
    result, err := funcs["my_func"](testVM, []*types.Value{
        types.NewString("test input"),
    })

    if err != nil {
        t.Fatalf("Function call failed: %v", err)
    }

    expected := "expected output"
    if result.ToString() != expected {
        t.Errorf("Expected %q, got %q", expected, result.ToString())
    }
}
```

### Table-Driven Tests

```go
func TestFactorial(t *testing.T) {
    tests := []struct {
        name     string
        input    int64
        expected int64
        wantErr  bool
    }{
        {"zero", 0, 1, false},
        {"one", 1, 1, false},
        {"five", 5, 120, false},
        {"negative", -1, 0, true},
    }

    ext := NewMathExtension()
    testVM := vm.New()
    funcs := ext.Functions()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := funcs["factorial"](testVM, []*types.Value{
                types.NewInt(tt.input),
            })

            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr = %v, got err = %v", tt.wantErr, err)
                return
            }

            if !tt.wantErr && result.ToInt() != tt.expected {
                t.Errorf("Expected %d, got %d", tt.expected, result.ToInt())
            }
        })
    }
}
```

### Integration Testing

```go
func TestExtensionIntegration(t *testing.T) {
    // Compile PHP code
    code := `<?php
    $result = my_func("test");
    echo $result;
    `

    // Parse and compile
    l := lexer.New(code)
    p := parser.New(l)
    program, err := p.ParseProgram()
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }

    c := compiler.New()
    if err := c.Compile(program); err != nil {
        t.Fatalf("Compile error: %v", err)
    }

    // Create VM with extension
    testVM := vm.New()
    registry := goext.NewFunctionRegistry()
    manager := goext.NewExtensionManager(registry)

    ext := NewMyExtension()
    manager.Register(ext)
    manager.LoadIntoVM(ext.Name(), testVM)

    // Load bytecode
    testVM.LoadBytecode(c.Bytecode())

    // Execute
    if err := testVM.Run(); err != nil {
        t.Fatalf("Runtime error: %v", err)
    }

    // Check output
    output := testVM.Output()
    expected := "expected output"
    if output != expected {
        t.Errorf("Expected %q, got %q", expected, output)
    }
}
```

### Benchmarking

```go
func BenchmarkMyFunction(b *testing.B) {
    ext := NewMyExtension()
    testVM := vm.New()
    funcs := ext.Functions()

    args := []*types.Value{types.NewString("test")}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := funcs["my_func"](testVM, args)
        if err != nil {
            b.Fatalf("Function call failed: %v", err)
        }
    }
}
```

## Performance Optimization

### Memory Optimization

#### Value Pooling

Reuse Value objects to reduce allocations:

```go
var valuePool = sync.Pool{
    New: func() interface{} {
        return &types.Value{}
    },
}

func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Get from pool
    result := valuePool.Get().(*types.Value)

    // Use result...
    result.SetString("hello")

    // Don't put back in pool - return value ownership transfers
    return result, nil
}
```

#### Pre-allocate Slices

```go
func processArray(v *vm.VM, args []*types.Value) (*types.Value, error) {
    arr := args[0].Array()
    elements := arr.Elements()

    // Pre-allocate result slice
    results := make([]*types.Value, 0, len(elements))

    for _, elem := range elements {
        // Process elem...
        results = append(results, processed)
    }

    return types.NewArrayFromSlice(results), nil
}
```

### CPU Optimization

#### Avoid Type Conversions

Cache converted values when possible:

```go
func processMultipleTimes(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Convert once
    str := args[0].ToString()

    // Use multiple times
    result1 := process1(str)
    result2 := process2(str)

    return types.NewString(result1 + result2), nil
}
```

#### Use Integer Operations

Integer operations are faster than float:

```go
// Prefer this
func calculateInt(a, b int64) int64 {
    return a * 100 / b  // Integer percentage
}

// Over this
func calculateFloat(a, b int64) float64 {
    return float64(a) / float64(b) * 100.0
}
```

### Concurrency

#### Goroutines for Parallel Work

```go
func parallelProcess(v *vm.VM, args []*types.Value) (*types.Value, error) {
    arr := args[0].Array()
    elements := arr.Elements()

    results := make([]*types.Value, len(elements))
    var wg sync.WaitGroup

    for i, elem := range elements {
        wg.Add(1)
        go func(idx int, val *types.Value) {
            defer wg.Done()
            // Process in parallel
            results[idx] = process(val)
        }(i, elem)
    }

    wg.Wait()
    return types.NewArrayFromSlice(results), nil
}
```

#### Worker Pools

```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}

func NewWorkerPool(workers int) *WorkerPool {
    wp := &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, workers*2),
        results: make(chan Result, workers*2),
    }

    for i := 0; i < workers; i++ {
        go wp.worker()
    }

    return wp
}

func (wp *WorkerPool) worker() {
    for job := range wp.jobs {
        result := processJob(job)
        wp.results <- result
    }
}
```

### Profiling

#### CPU Profiling

```go
import (
    "os"
    "runtime/pprof"
)

func BenchmarkWithProfile(b *testing.B) {
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Run benchmarks...
}
```

View profile:
```bash
go tool pprof cpu.prof
```

#### Memory Profiling

```go
func TestWithMemProfile(t *testing.T) {
    // Run tests...

    f, _ := os.Create("mem.prof")
    defer f.Close()
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

View profile:
```bash
go tool pprof mem.prof
```

## Best Practices

### Error Handling

#### Return Descriptive Errors

```go
func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf(
            "myFunc() expects exactly 2 arguments, got %d",
            len(args),
        )
    }

    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf(
            "myFunc() argument 1 must be string, got %s",
            args[0].TypeString(),
        )
    }

    // Process...
}
```

#### Wrap Errors for Context

```go
import "fmt"

func processFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    if err := validate(data); err != nil {
        return fmt.Errorf("validation failed for %s: %w", path, err)
    }

    return nil
}
```

### Thread Safety

#### Protect Shared State

```go
type CacheExtension struct {
    mu    sync.RWMutex
    cache map[string]*types.Value
}

func (e *CacheExtension) get(key string) (*types.Value, bool) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    val, ok := e.cache[key]
    return val, ok
}

func (e *CacheExtension) set(key string, val *types.Value) {
    e.mu.Lock()
    defer e.mu.Unlock()

    e.cache[key] = val
}
```

#### Use sync.Once for Initialization

```go
type MyExtension struct {
    once   sync.Once
    client *http.Client
}

func (e *MyExtension) getClient() *http.Client {
    e.once.Do(func() {
        e.client = &http.Client{
            Timeout: 30 * time.Second,
        }
    })
    return e.client
}
```

### Documentation

#### Document Functions

```go
// factorial calculates the factorial of n.
//
// In PHP:
//   $result = factorial(5);  // 120
//
// Parameters:
//   - n: Non-negative integer
//
// Returns:
//   - Factorial of n
//
// Errors:
//   - If n < 0: "factorial() expects non-negative integer"
func factorial(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Implementation...
}
```

#### Document Types

```go
// Config holds configuration for the database extension.
type Config struct {
    // Driver is the database driver name (e.g., "postgres", "mysql")
    Driver string

    // DSN is the data source name for connecting
    DSN string

    // MaxConnections is the maximum number of connections in the pool
    MaxConnections int

    // Timeout is the connection timeout duration
    Timeout time.Duration
}
```

### Naming Conventions

#### Function Names

- Use snake_case for PHP function names: `my_function`
- Use descriptive names: `parse_json` not `pj`
- Prefix with extension name for namespacing: `db_query`, `http_get`

#### Constant Names

- Use UPPER_SNAKE_CASE: `MY_CONSTANT`
- Prefix with extension name: `DB_VERSION`, `HTTP_TIMEOUT`

#### Go Code

- Follow Go conventions: `MyExtension`, `processData`
- Use meaningful names: `userCache` not `uc`

### Versioning

Use semantic versioning (semver):

```go
ext := goext.NewBaseExtension("myext", "1.2.3")
//                                      │ │ │
//                                      │ │ └─ Patch: Bug fixes
//                                      │ └─── Minor: New features (backwards compatible)
//                                      └───── Major: Breaking changes
```

## Deployment

### Static Linking

Build extensions directly into the PHP-Go binary:

```go
package main

import (
    "github.com/krizos/php-go/cmd"
    "example.com/myextension"
)

func init() {
    // Register extension at startup
    cmd.RegisterExtension(myextension.NewMyExtension())
}

func main() {
    cmd.Execute()
}
```

Build:
```bash
go build -o php-go-custom main.go
```

### Dynamic Loading (Plugins)

1. Build as plugin:
```bash
go build -buildmode=plugin -o myext.so
```

2. Install to extension directory:
```bash
sudo cp myext.so /usr/local/lib/php-go/extensions/
```

3. Load in PHP:
```php
<?php
// Extensions in /usr/local/lib/php-go/extensions/ are auto-loaded
```

### Configuration

#### Extension Configuration File

Create `/etc/php-go/extensions.ini`:

```ini
; Extension configuration
[myextension]
enabled = true
config.option1 = value1
config.option2 = value2

[database]
enabled = true
driver = postgres
dsn = postgres://localhost/mydb
max_connections = 100
```

#### Loading Configuration

```go
type Config struct {
    Option1 string
    Option2 string
}

func NewMyExtension(configFile string) (*MyExtension, error) {
    config, err := loadConfig(configFile)
    if err != nil {
        return nil, err
    }

    return &MyExtension{
        config: config,
    }, nil
}
```

### Distribution

#### Release Checklist

- [ ] Version bump (semver)
- [ ] Update CHANGELOG.md
- [ ] Run all tests: `go test ./...`
- [ ] Run benchmarks: `go test -bench=.`
- [ ] Update documentation
- [ ] Build for all platforms
- [ ] Create GitHub release
- [ ] Update examples

#### Building for Multiple Platforms

```bash
#!/bin/bash
# build-all.sh

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"

    output="myext-${GOOS}-${GOARCH}"
    if [ "$GOOS" == "windows" ]; then
        output="${output}.dll"
    else
        output="${output}.so"
    fi

    echo "Building $output..."
    GOOS=$GOOS GOARCH=$GOARCH go build -buildmode=plugin -o "dist/$output"
done
```

## Troubleshooting

### Common Issues

#### "Symbol not found" Error

**Problem**: Plugin fails to load with symbol not found.

**Solution**: Ensure the `Extension` variable is exported:

```go
var Extension goext.Extension  // Correct (uppercase)
var extension goext.Extension  // Wrong (lowercase)
```

#### "Version mismatch" Error

**Problem**: Plugin built with different Go version.

**Solution**: Rebuild plugin with same Go version as PHP-Go:

```bash
go version  # Check Go version
go build -buildmode=plugin -o myext.so
```

#### Type Conversion Panics

**Problem**: Panic when converting types.

**Solution**: Always check types before conversion:

```go
// Wrong
str := args[0].ToString()  // Panics if not string

// Correct
if args[0].Type() != types.TypeString {
    return nil, fmt.Errorf("expected string")
}
str := args[0].ToString()
```

#### Memory Leaks

**Problem**: Memory usage grows over time.

**Solution**: Use finalizers for cleanup:

```go
func (e *MyExtension) Init(v *vm.VM) error {
    e.resource = allocateResource()

    runtime.SetFinalizer(e, func(e *MyExtension) {
        if e.resource != nil {
            e.resource.Close()
        }
    })

    return nil
}
```

### Debugging

#### Enable Debug Logging

```go
import "log"

func myFunc(v *vm.VM, args []*types.Value) (*types.Value, error) {
    log.Printf("myFunc called with %d args", len(args))

    for i, arg := range args {
        log.Printf("  arg[%d]: type=%s, value=%v", i, arg.TypeString(), arg)
    }

    // Process...
}
```

#### Use Go Debugger (Delve)

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug extension tests
dlv test ./myextension
(dlv) break myFunc
(dlv) continue
```

#### Memory Profiling

```bash
# Run with memory profile
go test -memprofile=mem.prof

# Analyze
go tool pprof mem.prof
(pprof) top
(pprof) list myFunc
```

### Getting Help

- **Documentation**: Check `docs/extension-guide/` for detailed guides
- **Examples**: Review `pkg/goext/bindings/` for working examples
- **Tests**: Study `pkg/goext/*_test.go` for usage patterns
- **Issues**: Report bugs at https://github.com/krizos/php-go/issues
- **Community**: Join discussions on GitHub

## Appendix

### Complete Example Extension

Here's a complete, production-ready extension example:

```go
package weather

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/krizos/php-go/pkg/goext"
    "github.com/krizos/php-go/pkg/types"
    "github.com/krizos/php-go/pkg/vm"
)

// WeatherExtension provides weather information functions
type WeatherExtension struct {
    *goext.BaseExtension
    apiKey     string
    httpClient *http.Client
    cache      map[string]*cacheEntry
    cacheMu    sync.RWMutex
    cacheTTL   time.Duration
}

type cacheEntry struct {
    data      *types.Value
    expiresAt time.Time
}

type weatherResponse struct {
    Temperature float64 `json:"temp"`
    Humidity    int     `json:"humidity"`
    Condition   string  `json:"condition"`
}

// NewWeatherExtension creates a new weather extension
func NewWeatherExtension(apiKey string) *WeatherExtension {
    ext := &WeatherExtension{
        BaseExtension: goext.NewBaseExtension("weather", "1.0.0"),
        apiKey:        apiKey,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        cache:    make(map[string]*cacheEntry),
        cacheTTL: 5 * time.Minute,
    }

    // Register functions
    ext.AddFunction("weather_get", ext.getWeather)
    ext.AddFunction("weather_forecast", ext.getForecast)
    ext.AddFunction("weather_clear_cache", ext.clearCache)

    // Register constants
    ext.AddConstant("WEATHER_VERSION", types.NewString("1.0.0"))
    ext.AddConstant("WEATHER_CACHE_TTL", types.NewInt(int64(ext.cacheTTL.Seconds())))

    return ext
}

// getWeather retrieves current weather for a city
//
// Usage in PHP:
//   $weather = weather_get("London");
//   echo $weather['temperature'];
//
// Parameters:
//   - city (string): City name
//
// Returns:
//   - array: Weather data with keys: temperature, humidity, condition
func (e *WeatherExtension) getWeather(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Validate arguments
    if len(args) != 1 {
        return nil, fmt.Errorf("weather_get() expects 1 argument, got %d", len(args))
    }

    if args[0].Type() != types.TypeString {
        return nil, fmt.Errorf("weather_get() expects string, got %s", args[0].TypeString())
    }

    city := args[0].ToString()

    // Check cache
    if cached := e.getFromCache(city); cached != nil {
        return cached, nil
    }

    // Fetch from API
    weather, err := e.fetchWeather(city)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch weather: %w", err)
    }

    // Convert to PHP array
    result := types.NewArray()
    result.Set(types.NewString("temperature"), types.NewFloat(weather.Temperature))
    result.Set(types.NewString("humidity"), types.NewInt(int64(weather.Humidity)))
    result.Set(types.NewString("condition"), types.NewString(weather.Condition))

    // Cache result
    e.putInCache(city, result)

    return result, nil
}

// fetchWeather fetches weather data from API
func (e *WeatherExtension) fetchWeather(city string) (*weatherResponse, error) {
    url := fmt.Sprintf("https://api.weather.com/v1/current?city=%s&key=%s", city, e.apiKey)

    resp, err := e.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    var weather weatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
        return nil, err
    }

    return &weather, nil
}

// getFromCache retrieves cached weather data
func (e *WeatherExtension) getFromCache(city string) *types.Value {
    e.cacheMu.RLock()
    defer e.cacheMu.RUnlock()

    entry, ok := e.cache[city]
    if !ok {
        return nil
    }

    if time.Now().After(entry.expiresAt) {
        return nil
    }

    return entry.data
}

// putInCache stores weather data in cache
func (e *WeatherExtension) putInCache(city string, data *types.Value) {
    e.cacheMu.Lock()
    defer e.cacheMu.Unlock()

    e.cache[city] = &cacheEntry{
        data:      data,
        expiresAt: time.Now().Add(e.cacheTTL),
    }
}

// clearCache clears the weather cache
func (e *WeatherExtension) clearCache(v *vm.VM, args []*types.Value) (*types.Value, error) {
    e.cacheMu.Lock()
    defer e.cacheMu.Unlock()

    e.cache = make(map[string]*cacheEntry)

    return types.NewBool(true), nil
}

// getForecast retrieves weather forecast (stub implementation)
func (e *WeatherExtension) getForecast(v *vm.VM, args []*types.Value) (*types.Value, error) {
    // Implementation similar to getWeather...
    return types.NewArray(), nil
}
```

### Further Reading

- [Go Extension FFI Guide](ffi-guide.md)
- [Type Marshaling Guide](marshaling-guide.md)
- [Plugin System Documentation](04-plugin-system.md)
- [Quick Reference](quick-reference.md)
- [API Reference](../api-reference/) (when available)

### Changelog

- **1.0.0** (2024-11-24): Initial extension development guide
  - Complete guide covering all aspects of extension development
  - Examples for common use cases
  - Best practices and optimization techniques
  - Troubleshooting section

---

**Need help?** File an issue at https://github.com/krizos/php-go/issues
