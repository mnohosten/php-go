# Go Extension API Reference

Package: `github.com/krizos/php-go/pkg/goext`

## Overview

The goext package will implement Go integration features, allowing PHP code to call Go functions and Go code to interact with PHP. This includes FFI (Foreign Function Interface), type marshaling between Go and PHP, and a native extension API. This package is planned for Phase 8 and is currently a placeholder.

**Status**: Not yet implemented (Phase 8 - Planned)

**Implementation Timeline**: Phase 8 of the project roadmap

**Goal**: Seamless interoperability between PHP and Go code

## Planned Features

### 1. FFI (Foreign Function Interface)
Call Go functions from PHP and vice versa.

### 2. Type Marshaling
Automatic conversion between PHP and Go types.

### 3. Extension API
Register Go functions as PHP functions.

### 4. Struct Mapping
Map Go structs to PHP objects.

### 5. Go Bindings
Pre-built bindings for common Go packages.

## Planned API

### Extension Registration

```go
// Register a Go function as a PHP function
func RegisterFunction(name string, fn interface{}) error

// Register a Go struct as a PHP class
func RegisterClass(name string, structType reflect.Type) error

// Register a Go package
func RegisterPackage(pkgPath string) error
```

### Type Marshaling

```go
// Convert PHP value to Go value
func ToGo(value *types.Value, target interface{}) error

// Convert Go value to PHP value
func ToPhp(value interface{}) (*types.Value, error)

// Marshal struct to PHP object
func MarshalStruct(s interface{}) (*types.Object, error)

// Unmarshal PHP object to struct
func UnmarshalStruct(obj *types.Object, target interface{}) error
```

### FFI Interface

```go
// Call Go function from PHP
type GoFunction struct {
    Name      string
    Function  interface{}
    Signature *FunctionSignature
}

type FunctionSignature struct {
    Params     []ParamSpec
    ReturnType TypeSpec
    Variadic   bool
}

type ParamSpec struct {
    Name string
    Type TypeSpec
}

type TypeSpec struct {
    Kind TypeKind
    Name string
}

type TypeKind int

const (
    TypeInt TypeKind = iota
    TypeFloat
    TypeString
    TypeBool
    TypeArray
    TypeObject
    TypeInterface
    TypeStruct
)
```

### Extension Definition

```go
// Define a PHP extension in Go
type Extension struct {
    Name         string
    Version      string
    Functions    []GoFunction
    Classes      []GoClass
    Constants    map[string]interface{}
    Dependencies []string
}

// Register extension
func RegisterExtension(ext *Extension) error

// Get registered extension
func GetExtension(name string) (*Extension, error)
```

### Go Class Binding

```go
// Bind Go struct as PHP class
type GoClass struct {
    Name       string
    Type       reflect.Type
    Methods    []GoMethod
    Properties []GoProperty
    Constructor func() interface{}
}

type GoMethod struct {
    Name     string
    Method   reflect.Method
    Static   bool
    Public   bool
}

type GoProperty struct {
    Name   string
    Field  reflect.StructField
    Public bool
    Getter func(interface{}) interface{}
    Setter func(interface{}, interface{})
}
```

### Callback Support

```go
// Allow Go to call PHP functions
type PHPCallback struct {
    Function *vm.CompiledFunction
    Closure  *vm.Closure
}

func NewCallback(fn interface{}) (*PHPCallback, error)
func (cb *PHPCallback) Call(args ...*types.Value) (*types.Value, error)
```

### Plugin System

```go
// Load Go plugin
func LoadPlugin(path string) (*Plugin, error)

type Plugin struct {
    Name      string
    Path      string
    Module    *plugin.Plugin
    Extension *Extension
}

func (p *Plugin) GetFunction(name string) (interface{}, error)
func (p *Plugin) GetClass(name string) (reflect.Type, error)
```

## Usage Examples (Future)

### Example 1: Register Go Function

```go
package main

import (
    "github.com/krizos/php-go/pkg/goext"
    "math"
)

func main() {
    // Register Go's math.Sqrt as PHP function
    goext.RegisterFunction("sqrt", math.Sqrt)

    // Register custom function
    goext.RegisterFunction("greet", func(name string) string {
        return "Hello, " + name + "!"
    })
}
```

Then in PHP:
```php
<?php
echo sqrt(16);        // 4.0
echo greet("World");  // Hello, World!
```

### Example 2: Register Go Struct as PHP Class

```go
package main

import "github.com/krizos/php-go/pkg/goext"

type Person struct {
    Name string
    Age  int
}

func (p *Person) Greet() string {
    return "Hello, I'm " + p.Name
}

func main() {
    goext.RegisterClass("Person", reflect.TypeOf(Person{}))
}
```

Then in PHP:
```php
<?php
$person = new Person();
$person->Name = "Alice";
$person->Age = 30;
echo $person->Greet();  // Hello, I'm Alice
```

### Example 3: Use Go Standard Library

```go
package main

import (
    "github.com/krizos/php-go/pkg/goext"
    "time"
    "net/http"
)

func main() {
    // Register time functions
    goext.RegisterFunction("time_now", time.Now)
    goext.RegisterFunction("time_sleep", time.Sleep)

    // Register HTTP client
    goext.RegisterFunction("http_get", http.Get)
}
```

Then in PHP:
```php
<?php
$now = time_now();
echo $now->Format("2006-01-02");

$response = http_get("https://api.example.com/data");
$body = $response->Body->ReadAll();
```

### Example 4: Type Marshaling

```go
package main

import (
    "github.com/krizos/php-go/pkg/goext"
    "github.com/krizos/php-go/pkg/types"
)

func processData(phpArr *types.Value) (*types.Value, error) {
    // Convert PHP array to Go slice
    var data []int
    err := goext.ToGo(phpArr, &data)
    if err != nil {
        return nil, err
    }

    // Process in Go
    sum := 0
    for _, n := range data {
        sum += n
    }

    // Convert back to PHP
    return goext.ToPhp(sum)
}
```

## Pre-built Bindings

The goext package will include pre-built bindings for common Go packages:

### Time Package (`pkg/goext/bindings/time.go`)
```php
<?php
$now = Go\Time\Now();
$tomorrow = $now->Add(Go\Time\Hour(24));
echo $tomorrow->Format("2006-01-02");
```

### HTTP Client (`pkg/goext/bindings/http.go`)
```php
<?php
$client = new Go\Net\Http\Client();
$response = $client->Get("https://example.com");
echo $response->StatusCode;
```

### JSON (`pkg/goext/bindings/json.go`)
```php
<?php
$data = ["name" => "Alice", "age" => 30];
$json = Go\Encoding\Json\Marshal($data);
$decoded = Go\Encoding\Json\Unmarshal($json);
```

### Filesystem (`pkg/goext/bindings/filesystem.go`)
```php
<?php
$files = Go\Os\ReadDir("/path/to/dir");
foreach ($files as $file) {
    echo $file->Name();
}
```

### Crypto (`pkg/goext/bindings/crypto.go`)
```php
<?php
$hash = Go\Crypto\Sha256\Sum("hello world");
$encoded = Go\Encoding\Base64\Encode($hash);
```

## Implementation Plan (Phase 8)

### Task 8.1: Type Marshaling (40 hours)
Implement conversion between PHP and Go types.

### Task 8.2: FFI Core (50 hours)
Implement function call interface.

### Task 8.3: Extension API (40 hours)
Implement extension registration and loading.

### Task 8.4: Struct Binding (35 hours)
Implement Go struct to PHP class mapping.

### Task 8.5: Standard Bindings (60 hours)
Create bindings for common Go packages.

### Task 8.6: Plugin System (30 hours)
Implement dynamic plugin loading.

### Task 8.7: Testing and Documentation (45 hours)
Comprehensive testing and documentation.

**Total Estimated**: 300 hours for Phase 8

## Type Mapping

| PHP Type | Go Type | Notes |
|----------|---------|-------|
| null | nil | |
| bool | bool | |
| int | int64 | Always 64-bit |
| float | float64 | Always 64-bit |
| string | string | UTF-8 compatible |
| array | []interface{} or map[string]interface{} | Depends on keys |
| object | map[string]interface{} or struct | Depends on class |
| resource | unsafe.Pointer | For Go resources |
| callable | func | Function references |

## Performance Considerations

1. **Zero-copy**: Minimize data copying between PHP and Go
2. **Caching**: Cache type reflection information
3. **Native Speed**: Go code runs at native speed
4. **Overhead**: FFI calls have ~10-20ns overhead

## Safety Considerations

1. **Type Safety**: Runtime type checking for FFI calls
2. **Panic Recovery**: Go panics converted to PHP exceptions
3. **Memory Safety**: Prevent use-after-free
4. **Thread Safety**: Proper synchronization between PHP and Go goroutines

## Design Decisions

1. **Reflection-based**: Use Go reflection for automatic marshaling
2. **Opt-in**: Extensions must be explicitly registered
3. **Namespaced**: Go functions live in `Go\` namespace in PHP
4. **Versioned**: Extensions declare version compatibility

## Limitations

1. **No CGo**: Pure Go only (no C dependencies)
2. **Single VM**: One PHP VM per Go process
3. **Goroutine Isolation**: PHP code runs in specific goroutines
4. **Memory Model**: Different GC between PHP and Go

## References

- Phase 8 Documentation: `/Users/krizos/code/mnohosten/php-go/docs/phases/08-go-integration/`
- Go Plugin Package: https://pkg.go.dev/plugin
- Go Reflect Package: https://pkg.go.dev/reflect
- PHP FFI Extension: https://www.php.net/manual/en/book.ffi.php

## Future Enhancements (Beyond Phase 8)

- WebAssembly integration
- Remote procedure calls (RPC) support
- Automatic binding generation from Go code
- IDE integration for autocomplete
- Performance profiling across PHP/Go boundary

## Related Packages

- `pkg/types` - PHP type system
- `pkg/vm` - PHP VM for executing callbacks
- `pkg/runtime` - Runtime for managing execution context

## Contributing

When Phase 8 begins, contributions to Go integration are welcome. Each binding should:
- Follow Go conventions
- Be thoroughly tested
- Have clear documentation
- Include usage examples
- Consider performance implications
