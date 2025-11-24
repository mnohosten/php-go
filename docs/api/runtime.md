# Runtime API Reference

Package: `github.com/krizos/php-go/pkg/runtime`

## Overview

The runtime package manages PHP runtime state including superglobals ($_GET, $_POST, etc.), constants, error handling, output buffering, and execution context. It provides the infrastructure needed for PHP script execution outside of the VM itself.

## Main Types

### Runtime

The main runtime manager that holds PHP execution state.

```go
type Runtime struct {
    // Private fields
}
```

**Constructor:**

```go
func New() *Runtime
```

Creates a new runtime instance with initialized superglobals and built-in constants.

**Superglobals:**

```go
func (rt *Runtime) GetSuperglobal(name string) (*types.Value, bool)
```

Retrieves a superglobal by name. Supported names:
- `_GET` - HTTP GET variables
- `_POST` - HTTP POST variables
- `_REQUEST` - HTTP request variables
- `_SERVER` - Server and execution environment information
- `_ENV` - Environment variables
- `_COOKIE` - HTTP cookies
- `_FILES` - HTTP file uploads
- `_SESSION` - Session variables
- `GLOBALS` - All global variables

**Direct Superglobal Access:**

```go
rt.GET     *types.Value
rt.POST    *types.Value
rt.REQUEST *types.Value
rt.SERVER  *types.Value
rt.ENV     *types.Value
rt.COOKIE  *types.Value
rt.FILES   *types.Value
rt.SESSION *types.Value
rt.GLOBALS *types.Value
```

### Constants Management

**Define Constants:**

```go
func (rt *Runtime) DefineConstant(name string, value *types.Value) error
```

Defines a constant. Returns error if constant already exists.

```go
func (rt *Runtime) GetConstant(name string) (*types.Value, bool)
```

Retrieves a constant value.

```go
func (rt *Runtime) ConstantExists(name string) bool
```

Checks if a constant is defined.

**Built-in Constants:**

The runtime automatically defines these constants:

**PHP Version:**
- `PHP_VERSION` - "8.4.0-dev"
- `PHP_MAJOR_VERSION` - 8
- `PHP_MINOR_VERSION` - 4
- `PHP_RELEASE_VERSION` - 0

**Boolean Constants:**
- `TRUE` - true
- `FALSE` - false
- `NULL` - null

**Path Constants:**
- `PHP_EOL` - Line ending ("\n")
- `DIRECTORY_SEPARATOR` - Path separator

**Math Constants:**
- `PHP_INT_MAX` - 9223372036854775807
- `PHP_INT_MIN` - -9223372036854775808
- `PHP_FLOAT_MAX` - Maximum float value
- `PHP_FLOAT_MIN` - Minimum positive float value

### Error Handling

**Error Types:**

```go
type ErrorType int

const (
    E_ERROR             ErrorType = 1 << 0
    E_WARNING           ErrorType = 1 << 1
    E_PARSE             ErrorType = 1 << 2
    E_NOTICE            ErrorType = 1 << 3
    E_CORE_ERROR        ErrorType = 1 << 4
    E_CORE_WARNING      ErrorType = 1 << 5
    E_COMPILE_ERROR     ErrorType = 1 << 6
    E_COMPILE_WARNING   ErrorType = 1 << 7
    E_USER_ERROR        ErrorType = 1 << 8
    E_USER_WARNING      ErrorType = 1 << 9
    E_USER_NOTICE       ErrorType = 1 << 10
    E_STRICT            ErrorType = 1 << 11
    E_RECOVERABLE_ERROR ErrorType = 1 << 12
    E_DEPRECATED        ErrorType = 1 << 13
    E_USER_DEPRECATED   ErrorType = 1 << 14
    E_ALL               ErrorType = 0x7FFF
)
```

**Error Handling Methods:**

```go
func (rt *Runtime) SetErrorReporting(level int)
```

Sets the error reporting level (bitmask of ErrorType constants).

```go
func (rt *Runtime) GetErrorReporting() int
```

Gets the current error reporting level.

```go
type ErrorHandler func(errorType ErrorType, message string, file string, line int)

func (rt *Runtime) SetErrorHandler(handler ErrorHandler)
```

Sets a custom error handler function.

```go
func (rt *Runtime) TriggerError(errorType ErrorType, message string, file string, line int)
```

Triggers an error. Only reports if the error type is enabled in error reporting level.

**Methods:**

```go
func (et ErrorType) String() string
```

Returns the name of the error type (e.g., "E_WARNING").

### Output Buffering

The runtime supports nested output buffering (like PHP's ob_start/ob_end).

**OutputBuffer:**

```go
type OutputBuffer struct {
    // Private fields
}
```

**Constructor:**

```go
func NewOutputBuffer() *OutputBuffer
```

**Methods:**

```go
func (ob *OutputBuffer) Write(data string)
func (ob *OutputBuffer) GetContents() string
func (ob *OutputBuffer) Clear()
func (ob *OutputBuffer) Flush() string
```

**Runtime Output Buffer Management:**

```go
func (rt *Runtime) StartOutputBuffering()
```

Starts a new output buffer. Output will be captured until `EndOutputBuffering()` is called.

```go
func (rt *Runtime) EndOutputBuffering() string
```

Ends the current output buffer and returns its contents.

```go
func (rt *Runtime) CleanOutputBuffer()
```

Ends the current buffer and discards its contents.

```go
func (rt *Runtime) GetOutputBufferContents() string
```

Gets the contents of the current buffer without ending it.

```go
func (rt *Runtime) FlushOutputBuffer() string
```

Flushes the current output buffer (returns and clears contents).

```go
func (rt *Runtime) Write(data string)
```

Writes to the current output buffer, or stdout if no buffer is active.

### Execution Context

```go
func (rt *Runtime) SetScriptPath(path string)
```

Sets the current script path. This updates $_SERVER['SCRIPT_FILENAME'] and $_SERVER['SCRIPT_NAME'].

```go
func (rt *Runtime) GetScriptPath() string
```

Gets the current script path.

## Exception Support

**Exception Types:**

```go
type Exception struct {
    Message string
    Code    int
    File    string
    Line    int
    Trace   []StackFrame
}
```

**Constructor:**

```go
func NewException(message string, code int) *Exception
```

**Methods:**

```go
func (e *Exception) Error() string
func (e *Exception) String() string
func (e *Exception) SetTrace(trace []StackFrame)
```

**StackFrame:**

```go
type StackFrame struct {
    Function string
    File     string
    Line     int
    Class    string
    Type     string
}
```

## Callable Support

**CallableType:**

```go
type CallableType int

const (
    CallableFunction CallableType = iota
    CallableMethod
    CallableStaticMethod
    CallableClosure
    CallableInvokable
)
```

**Callable:**

```go
type Callable struct {
    Type       CallableType
    Name       string
    Class      string
    Object     *types.Object
    Closure    *vm.Closure
}
```

**Functions:**

```go
func NewCallable(value *types.Value) (*Callable, error)
func (c *Callable) String() string
func (c *Callable) IsValid() bool
```

## Attribute Support (PHP 8.0+)

**Attribute:**

```go
type Attribute struct {
    Name      string
    Arguments []*types.Value
    Target    AttributeTarget
}
```

**AttributeTarget:**

```go
type AttributeTarget uint8

const (
    AttributeTargetClass AttributeTarget = 1 << iota
    AttributeTargetFunction
    AttributeTargetMethod
    AttributeTargetProperty
    AttributeTargetClassConstant
    AttributeTargetParameter
    AttributeTargetAll = 0xFF
)
```

## Metrics and Logging

**Metrics:**

```go
type Metrics struct {
    ExecutionTime time.Duration
    MemoryUsage   uint64
    PeakMemory    uint64
    FileCount     int
    FunctionCalls int
}
```

**Methods:**

```go
func (rt *Runtime) GetMetrics() *Metrics
func (rt *Runtime) ResetMetrics()
```

**Logging:**

```go
type LogLevel int

const (
    LogDebug LogLevel = iota
    LogInfo
    LogWarning
    LogError
)
```

```go
func (rt *Runtime) Log(level LogLevel, message string)
func (rt *Runtime) SetLogLevel(level LogLevel)
```

## Health Checks

**HealthStatus:**

```go
type HealthStatus struct {
    Healthy     bool
    ErrorCount  int
    Warnings    []string
    LastCheck   time.Time
}
```

**Methods:**

```go
func (rt *Runtime) CheckHealth() *HealthStatus
```

## Resource Limits

**Limits:**

```go
type Limits struct {
    MaxExecutionTime time.Duration
    MaxMemory        uint64
    MaxFileSize      uint64
    MaxFileUploads   int
}
```

**Methods:**

```go
func (rt *Runtime) SetLimits(limits *Limits)
func (rt *Runtime) GetLimits() *Limits
func (rt *Runtime) CheckLimits() error
```

## Shutdown Functions

```go
func (rt *Runtime) RegisterShutdownFunction(fn func())
func (rt *Runtime) ExecuteShutdownFunctions()
```

## Variadic Arguments

**VariadicArgs:**

```go
type VariadicArgs struct {
    Args []*types.Value
}
```

**Functions:**

```go
func NewVariadicArgs(args ...*types.Value) *VariadicArgs
func (va *VariadicArgs) ToArray() *types.Array
func (va *VariadicArgs) Len() int
func (va *VariadicArgs) Get(index int) (*types.Value, bool)
```

## Weak References

**WeakReference:**

```go
type WeakReference struct {
    // Private fields
}
```

**Constructor:**

```go
func NewWeakReference(target *types.Object) *WeakReference
```

**Methods:**

```go
func (wr *WeakReference) Get() *types.Object
func (wr *WeakReference) IsValid() bool
```

## Error Recovery

**RecoveryHandler:**

```go
type RecoveryHandler func(err interface{}) error
```

**Methods:**

```go
func (rt *Runtime) SetRecoveryHandler(handler RecoveryHandler)
func (rt *Runtime) Recover() error
```

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/runtime"
    "github.com/krizos/php-go/pkg/types"
)

func main() {
    rt := runtime.New()

    // Define a constant
    rt.DefineConstant("APP_VERSION", types.NewString("1.0.0"))

    // Access superglobals
    server := rt.SERVER.ToArray()
    if scriptName, ok := server.Get(types.NewString("SCRIPT_NAME")); ok {
        fmt.Println("Script:", scriptName.AsString())
    }

    // Output buffering
    rt.StartOutputBuffering()
    rt.Write("Hello, ")
    rt.Write("World!")
    output := rt.EndOutputBuffering()
    fmt.Println(output)  // "Hello, World!"

    // Error handling
    rt.SetErrorReporting(int(runtime.E_ALL))
    rt.SetErrorHandler(func(errType runtime.ErrorType, msg, file string, line int) {
        fmt.Printf("[%s] %s in %s:%d\n", errType, msg, file, line)
    })

    rt.TriggerError(runtime.E_WARNING, "Test warning", "test.php", 10)

    // Check health
    health := rt.CheckHealth()
    fmt.Printf("Healthy: %v, Errors: %d\n", health.Healthy, health.ErrorCount)
}
```

## Implementation Notes

- **Thread Safety**: Runtime is NOT thread-safe (each execution context needs its own instance)
- **Superglobals**: Initialized as empty arrays by default
- **$_SERVER**: Auto-populated with environment variables and system information
- **Output Buffering**: Supports nested buffers (like PHP)
- **Error Reporting**: Uses bitmask for configuring which errors to report
- **Constants**: Case-sensitive (like PHP)
- **Shutdown Functions**: Execute in LIFO order (last registered, first executed)

## Coverage

The runtime package has 99%+ code coverage with comprehensive tests for all features.

## Integration with VM

The runtime package is typically used alongside the VM:

```go
rt := runtime.New()
machine := vm.New()

// Share globals between runtime and VM
machine.SetGlobal("_SERVER", rt.SERVER)
machine.SetGlobal("_GET", rt.GET)
// etc.

// Execute code
err := machine.Execute(bytecode)

// Access output
fmt.Println(rt.OutputString())
```
