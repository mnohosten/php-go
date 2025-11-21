# Phase 8: Go Integration

## Overview

Phase 8 implements bidirectional integration between PHP and Go - allowing PHP code to call Go functions and Go code to embed PHP, plus a native extension API.

## Goals

1. **FFI System**: Call Go functions from PHP code
2. **Type Marshaling**: Convert between PHP and Go types
3. **Native Extension API**: Write PHP extensions in Go
4. **Go Library Bindings**: Expose Go stdlib to PHP
5. **Plugin System**: Dynamic loading of Go modules

## Components

### 1. FFI System

**File**: `pkg/goext/ffi.go`

```go
// Register Go function for PHP use
func RegisterFunction(name string, fn interface{})

// Call Go function from PHP
// PHP: $result = go_call('mypackage.MyFunc', $arg1, $arg2);
func goCall(vm *VM, args []*Value) (*Value, error)
```

**Example**:
```go
// Go code
func Add(a, b int64) int64 {
    return a + b
}

func init() {
    goext.RegisterFunction("math.Add", Add)
}

// PHP code
<?php
$result = go_call('math.Add', 10, 20); // 30
```

### 2. Type Marshaling

**File**: `pkg/goext/marshal.go`

```go
type Marshaler struct{}

// PHP → Go
func (m *Marshaler) ToGo(phpVal *Value) (interface{}, error)
func (m *Marshaler) ToGoInt(phpVal *Value) (int64, error)
func (m *Marshaler) ToGoString(phpVal *Value) (string, error)
func (m *Marshaler) ToGoSlice(phpVal *Value) ([]interface{}, error)
func (m *Marshaler) ToGoMap(phpVal *Value) (map[string]interface{}, error)

// Go → PHP
func (m *Marshaler) ToPHP(goVal interface{}) (*Value, error)
func (m *Marshaler) IntToPHP(val int64) *Value
func (m *Marshaler) StringToPHP(val string) *Value
func (m *Marshaler) SliceToPHP(val []interface{}) *Value
func (m *Marshaler) MapToPHP(val map[string]interface{}) *Value
```

**Type Mappings**:
```
PHP              Go
---              --
null             nil
bool             bool
int              int64
float            float64
string           string
array (list)     []interface{}
array (assoc)    map[string]interface{}
object           struct or map
resource         interface{}
```

### 3. Native Extension API

**File**: `pkg/goext/extension.go`

```go
type Extension interface {
    Name() string
    Version() string
    Init(*VM) error
    Functions() map[string]FunctionHandler
    Classes() map[string]ClassDefinition
    Constants() map[string]*Value
}

type FunctionHandler func(*VM, []*Value) (*Value, error)
```

**Example Extension**:
```go
type MyExtension struct{}

func (e *MyExtension) Name() string {
    return "my_extension"
}

func (e *MyExtension) Functions() map[string]FunctionHandler {
    return map[string]FunctionHandler{
        "my_func": myFunc,
    }
}

func myFunc(vm *VM, args []*Value) (*Value, error) {
    // Implementation
    return NewString("Hello from Go!"), nil
}

// Register
func init() {
    RegisterExtension(&MyExtension{})
}
```

### 4. Go Library Bindings

**File**: `pkg/goext/bindings/`

Expose Go standard library to PHP:

**http binding** (`pkg/goext/bindings/http.go`):
```php
<?php
// Call Go's net/http
$response = go_http_get('https://api.example.com');
$data = json_decode($response['body']);
```

**crypto binding** (`pkg/goext/bindings/crypto.go`):
```php
<?php
// Use Go's crypto
$hash = go_crypto_sha256('data');
$encrypted = go_crypto_aes_encrypt($data, $key);
```

**database binding** (`pkg/goext/bindings/database.go`):
```php
<?php
// Use Go's database/sql
$db = go_db_connect('postgres', $dsn);
$rows = go_db_query($db, 'SELECT * FROM users');
```

### 5. Plugin System

**File**: `pkg/goext/plugins.go`

Dynamic loading of Go plugins:

```go
type PluginManager struct {
    plugins map[string]*plugin.Plugin
}

func (pm *PluginManager) Load(path string) error
func (pm *PluginManager) LoadExtension(path string) (Extension, error)
```

**Usage**:
```bash
# Build plugin
go build -buildmode=plugin -o myplugin.so myplugin.go

# Load in PHP
php-go --load-plugin=myplugin.so script.php
```

## Implementation Tasks

### Task 8.1: Type Marshaling Foundation
**Effort**: 12 hours

- [ ] PHP → Go conversion
- [ ] Go → PHP conversion
- [ ] Type mapping rules
- [ ] Error handling
- [ ] Edge case handling

### Task 8.2: Function Registration
**Effort**: 8 hours

- [ ] RegisterFunction() implementation
- [ ] Function signature parsing
- [ ] Reflection-based wrapping
- [ ] Function registry

### Task 8.3: FFI Call Implementation
**Effort**: 10 hours

- [ ] go_call() function
- [ ] Function lookup
- [ ] Argument marshaling
- [ ] Return value marshaling
- [ ] Error propagation

### Task 8.4: Extension API
**Effort**: 12 hours

- [ ] Extension interface
- [ ] Extension registration
- [ ] Extension initialization
- [ ] Function/class/constant loading
- [ ] Extension manager

### Task 8.5: Go Standard Library Bindings
**Effort**: 20 hours

- [ ] HTTP client bindings
- [ ] Crypto bindings
- [ ] Database bindings
- [ ] File system bindings
- [ ] JSON bindings (native Go)
- [ ] Time/date bindings

### Task 8.6: Plugin System
**Effort**: 10 hours

- [ ] Plugin loading
- [ ] Symbol resolution
- [ ] Plugin manager
- [ ] Hot reloading (optional)

### Task 8.7: Advanced Marshaling
**Effort**: 12 hours

- [ ] Custom type marshaling
- [ ] Struct ↔ Object conversion
- [ ] Interface{} handling
- [ ] Circular reference handling
- [ ] Performance optimization

### Task 8.8: Documentation & Examples
**Effort**: 10 hours

- [ ] Extension development guide
- [ ] Example extensions
- [ ] FFI usage guide
- [ ] Type marshaling guide
- [ ] Best practices

### Task 8.9: Testing
**Effort**: 12 hours

- [ ] Marshaling tests
- [ ] FFI tests
- [ ] Extension tests
- [ ] Integration tests
- [ ] Performance benchmarks

## Estimated Timeline

**Total Effort**: ~105 hours (5-6 weeks)

## Example Use Cases

### Use Case 1: High-Performance HTTP Client
```php
<?php
// Use Go's HTTP client instead of cURL
$response = go_http_get('https://api.example.com', [
    'timeout' => 10,
    'headers' => ['Authorization' => 'Bearer ' . $token]
]);
```

### Use Case 2: Native Encryption
```php
<?php
// Use Go's crypto (faster than PHP's)
$encrypted = go_crypto_aes_encrypt($data, $key, $iv);
$decrypted = go_crypto_aes_decrypt($encrypted, $key, $iv);
```

### Use Case 3: Custom Extension
```go
// Go extension
package myext

import "github.com/yourusername/php-go/pkg/goext"

func ProcessImage(vm *VM, args []*Value) (*Value, error) {
    // High-performance image processing in Go
    imageData := args[0].ToBytes()
    processed := processImageWithGo(imageData)
    return NewString(processed), nil
}

func init() {
    goext.RegisterFunction("image.process", ProcessImage)
}
```

```php
<?php
// PHP usage
$processed = go_call('image.process', file_get_contents('photo.jpg'));
```

## Success Criteria

- [ ] Can call Go functions from PHP
- [ ] Type marshaling works bidirectionally
- [ ] Can write extensions in pure Go
- [ ] Go stdlib accessible from PHP
- [ ] Plugin system functional
- [ ] Performance better than native PHP
- [ ] Test coverage >85%

## Progress Tracking

- [ ] Task 8.1: Type Marshaling (0%)
- [ ] Task 8.2: Function Registration (0%)
- [ ] Task 8.3: FFI Implementation (0%)
- [ ] Task 8.4: Extension API (0%)
- [ ] Task 8.5: Go Stdlib Bindings (0%)
- [ ] Task 8.6: Plugin System (0%)
- [ ] Task 8.7: Advanced Marshaling (0%)
- [ ] Task 8.8: Documentation (0%)
- [ ] Task 8.9: Testing (0%)

**Overall Phase 8 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
