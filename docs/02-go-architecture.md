# Go Architecture for PHP-Go

This document describes the Go-specific architecture and design decisions for the PHP-Go interpreter.

## Design Philosophy

1. **Idiomatic Go** - Follow Go conventions and best practices
2. **Type Safety** - Leverage Go's type system where possible
3. **Performance** - Optimize for Go's strengths (goroutines, GC, etc.)
4. **Maintainability** - Clean, well-documented, testable code
5. **Compatibility** - Full PHP 8.4 semantic compatibility

## Package Structure

```
pkg/
├── lexer/              # Tokenization
│   ├── lexer.go        # Main lexer
│   ├── token.go        # Token definitions
│   └── position.go     # Source position tracking
│
├── parser/             # Parsing
│   ├── parser.go       # Main parser
│   ├── grammar.go      # Grammar rules
│   └── errors.go       # Parse errors
│
├── ast/                # Abstract Syntax Tree
│   ├── node.go         # Base node interface
│   ├── expr.go         # Expression nodes
│   ├── stmt.go         # Statement nodes
│   ├── visitor.go      # Visitor pattern
│   └── printer.go      # AST printer (debugging)
│
├── compiler/           # Compiler (AST → Opcodes)
│   ├── compiler.go     # Main compiler
│   ├── context.go      # Compilation context
│   ├── symbols.go      # Symbol table
│   ├── jumps.go        # Jump resolution
│   └── optimizer.go    # Basic optimizations
│
├── vm/                 # Virtual Machine
│   ├── vm.go           # VM executor
│   ├── opcodes.go      # Opcode definitions
│   ├── instruction.go  # Instruction encoding
│   ├── frame.go        # Call frames
│   ├── stack.go        # Operand stack
│   └── handlers.go     # Opcode handlers
│
├── types/              # Type System
│   ├── value.go        # PHP value type (Zval equivalent)
│   ├── string.go       # PHP strings
│   ├── array.go        # PHP arrays
│   ├── object.go       # PHP objects
│   ├── resource.go     # PHP resources
│   ├── callable.go     # Callables (functions, methods)
│   ├── class.go        # Class definitions
│   ├── interface.go    # Interface definitions
│   └── conversions.go  # Type conversions & juggling
│
├── runtime/            # Runtime Support
│   ├── runtime.go      # Runtime environment
│   ├── globals.go      # Superglobals ($_GET, $_POST, etc.)
│   ├── functions.go    # Function registry
│   ├── classes.go      # Class registry
│   ├── constants.go    # Constant registry
│   ├── errors.go       # Error handling
│   ├── exceptions.go   # Exception handling
│   └── output.go       # Output buffering
│
├── stdlib/             # Standard Library
│   ├── array/          # Array functions
│   ├── string/         # String functions
│   ├── file/           # File I/O
│   ├── math/           # Math functions
│   ├── var/            # Variable functions
│   ├── json/           # JSON extension
│   ├── pcre/           # Regular expressions
│   ├── date/           # Date/time
│   └── ... /           # Other extensions
│
├── parallel/           # Parallelization Engine
│   ├── analyzer.go     # Safety analysis
│   ├── executor.go     # Parallel executor
│   ├── pool.go         # Goroutine pool
│   └── context.go      # Parallel context
│
└── goext/              # Go Integration
    ├── ffi.go          # FFI system
    ├── bindings.go     # Type marshaling
    ├── extensions.go   # Native extension API
    └── loader.go       # Dynamic loading
```

## Core Type System

### 1. PHP Value Representation

```go
// Value is the equivalent of PHP's zval
type Value struct {
    typ      ValueType    // Type tag
    flags    ValueFlags   // Reference, immutable, etc.
    data     interface{}  // Actual data (type-specific)
    refcount *int32       // Optional: for explicit sharing
}

// ValueType tags
type ValueType uint8

const (
    TypeUndef ValueType = iota
    TypeNull
    TypeBool
    TypeInt
    TypeFloat
    TypeString
    TypeArray
    TypeObject
    TypeResource
    TypeCallable
    TypeReference  // PHP reference (&$var)
)

// Alternative: Interface-based approach
type PHPValue interface {
    Type() ValueType
    String() string
    ToInt() int64
    ToFloat() float64
    ToBool() bool
    IsTrue() bool
    Copy() PHPValue
}
```

**Design Decision**: Use struct with interface{} for flexibility, with type assertions for access.

**Rationale**:
- Go's type system makes pure interface approach slow
- Struct approach allows efficient access
- interface{} allows any Go type as backing storage
- Type tag enables fast type checking

### 2. PHP Strings

```go
// String represents a PHP string
type String struct {
    data []byte  // UTF-8 or binary data
    len  int     // Length in bytes
    hash uint32  // Cached hash value
}

// Methods
func (s *String) Len() int
func (s *String) Bytes() []byte
func (s *String) String() string  // Convert to Go string
func (s *String) Hash() uint32
func (s *String) Concat(other *String) *String
func (s *String) Substr(start, length int) *String
```

**Design Decision**: Custom String type, not Go string

**Rationale**:
- PHP strings can contain null bytes (binary-safe)
- Go strings are immutable, PHP strings need mutability semantics
- Need to track binary vs text separately
- PHP string functions operate on bytes, not runes

### 3. PHP Arrays

```go
// Array represents PHP's ordered associative array
type Array struct {
    elements map[ArrayKey]*Value  // Hash map
    order    []ArrayKey            // Insertion order
    nextKey  int64                 // Next integer key
}

type ArrayKey struct {
    typ KeyType  // Integer or string
    i   int64
    s   string
}

// Methods
func (a *Array) Get(key ArrayKey) (*Value, bool)
func (a *Array) Set(key ArrayKey, value *Value)
func (a *Array) Append(value *Value)
func (a *Array) Delete(key ArrayKey)
func (a *Array) Count() int
func (a *Array) Keys() []ArrayKey
func (a *Array) Values() []*Value

// Optimized packed array for sequential integer keys
type PackedArray struct {
    elements []*Value
}
```

**Design Decision**: Map + slice for order, separate packed array optimization

**Rationale**:
- PHP arrays maintain insertion order
- Need both integer and string keys
- Common case (sequential integers) needs optimization
- Go maps don't maintain order

### 4. PHP Objects

```go
// Object represents a PHP object instance
type Object struct {
    class      *Class
    properties map[string]*Value
    handlers   *ObjectHandlers  // Magic methods
}

// Class represents a PHP class definition
type Class struct {
    name       string
    parent     *Class
    interfaces []*Interface
    traits     []*Trait
    constants  map[string]*Value
    properties map[string]*Property
    methods    map[string]*Method
    flags      ClassFlags  // Abstract, final, etc.
}

// Property with visibility
type Property struct {
    name       string
    visibility Visibility  // public, protected, private
    static     bool
    readonly   bool
    typ        *TypeInfo   // Type hint
    default    *Value
}

// Method definition
type Method struct {
    name       string
    visibility Visibility
    static     bool
    abstract   bool
    final      bool
    params     []*Parameter
    returnType *TypeInfo
    body       *Function   // Compiled function
}
```

**Design Decision**: Separate Class and Object types

**Rationale**:
- Classes are templates, objects are instances
- Allows efficient class sharing between objects
- Mirrors PHP's internal structure
- Enables compile-time class analysis

### 5. Type Conversions

```go
// Type juggling and conversions
type Converter interface {
    ToInt(v *Value) int64
    ToFloat(v *Value) float64
    ToBool(v *Value) bool
    ToString(v *Value) *String
    ToArray(v *Value) *Array
}

// Comparison with PHP semantics
func Compare(a, b *Value) int        // <, ==, >
func IdenticalTo(a, b *Value) bool   // ===
func EqualTo(a, b *Value) bool       // ==
```

**Design Decision**: Explicit conversion functions

**Rationale**:
- PHP has complex type juggling rules
- Need to implement exact PHP semantics
- Can't rely on Go's type conversion

## Virtual Machine Architecture

### 1. Opcode Design

```go
// Opcode represents a VM instruction
type Opcode uint8

// 210 opcodes (same as Zend)
const (
    OpNop Opcode = iota
    OpAdd
    OpSub
    OpMul
    OpDiv
    // ... all 210 opcodes
)

// Instruction encoding
type Instruction struct {
    opcode  Opcode
    op1     Operand
    op2     Operand
    result  Operand
    ext     uint32     // Extended info
}

// Operand types
type OperandType uint8

const (
    OpTypeUnused OperandType = iota
    OpTypeConst    // Constant from literal table
    OpTypeVar      // Variable
    OpTypeTmp      // Temporary
    OpTypeCV       // Compiled variable (optimized)
)

type Operand struct {
    typ OperandType
    num uint32      // Index/offset
}
```

**Design Decision**: Similar to Zend VM encoding

**Rationale**:
- Proven design
- Efficient encoding
- Allows optimization
- Familiar to PHP developers

### 2. VM Execution

```go
// VM is the virtual machine executor
type VM struct {
    frames    []*Frame       // Call stack
    globals   *GlobalScope   // Global variables
    functions map[string]*Function
    classes   map[string]*Class
    constants map[string]*Value
    output    *OutputBuffer
}

// Frame represents a call frame
type Frame struct {
    function  *Function
    locals    []*Value      // Local variables
    temps     []*Value      // Temporary values
    stack     []*Value      // Operand stack
    ip        int           // Instruction pointer
    returnTo  int           // Return instruction
}

// Execute runs the VM
func (vm *VM) Execute(fn *Function, args []*Value) (*Value, error) {
    frame := vm.newFrame(fn, args)
    vm.frames = append(vm.frames, frame)
    defer func() { vm.frames = vm.frames[:len(vm.frames)-1] }()

    for frame.ip < len(fn.Instructions) {
        instr := fn.Instructions[frame.ip]
        if err := vm.executeInstruction(frame, instr); err != nil {
            return nil, err
        }
        frame.ip++
    }

    return frame.returnValue, nil
}
```

**Design Decision**: Stack-based VM with call frames

**Rationale**:
- Matches Zend VM architecture
- Efficient for function calls
- Natural for PHP's execution model
- Supports generators/coroutines

### 3. Dispatch Method

```go
// executeInstruction dispatches opcodes
func (vm *VM) executeInstruction(frame *Frame, instr Instruction) error {
    switch instr.opcode {
    case OpAdd:
        return vm.opAdd(frame, instr)
    case OpSub:
        return vm.opSub(frame, instr)
    // ... handle all 210 opcodes
    default:
        return fmt.Errorf("unknown opcode: %d", instr.opcode)
    }
}

// Alternative: Jump table for performance
var opHandlers = [256]func(*VM, *Frame, Instruction) error{
    OpAdd: (*VM).opAdd,
    OpSub: (*VM).opSub,
    // ...
}

func (vm *VM) executeInstruction(frame *Frame, instr Instruction) error {
    handler := opHandlers[instr.opcode]
    if handler == nil {
        return fmt.Errorf("unknown opcode: %d", instr.opcode)
    }
    return handler(vm, frame, instr)
}
```

**Design Decision**: Start with switch, optimize to jump table if needed

**Rationale**:
- Switch is simple and fast enough initially
- Go compiler optimizes switches well
- Jump table optimization available if needed
- Profile-guided optimization later

## Memory Management

### 1. Garbage Collection

**Design Decision**: Use Go's GC exclusively

**Benefits**:
- No manual reference counting
- No circular reference issues
- Automatic memory management
- Simpler implementation

**Trade-offs**:
- Less deterministic memory usage
- Can't have __destruct run immediately
- May need tuning for PHP workloads

**Implementation**:
```go
// No explicit memory management needed!
// Go GC handles everything

// For __destruct, use finalizers sparingly
func (obj *Object) setDestructor() {
    runtime.SetFinalizer(obj, func(o *Object) {
        // Call __destruct if defined
        // NOTE: Not guaranteed to run immediately!
    })
}
```

### 2. String Interning

```go
// String pool for deduplication
type StringPool struct {
    mu      sync.RWMutex
    strings map[string]*String
}

func (sp *StringPool) Intern(s string) *String {
    sp.mu.RLock()
    if str, ok := sp.strings[s]; ok {
        sp.mu.RUnlock()
        return str
    }
    sp.mu.RUnlock()

    sp.mu.Lock()
    defer sp.mu.Unlock()

    // Double-check after acquiring write lock
    if str, ok := sp.strings[s]; ok {
        return str
    }

    str := &String{data: []byte(s), len: len(s)}
    sp.strings[s] = str
    return str
}
```

**Design Decision**: Optional string interning for literals

**Rationale**:
- Reduces memory for duplicate strings
- Go already interns string literals
- Useful for class/function/constant names
- Make it opt-in to avoid overhead

## Parallelization Architecture

### 1. Safety Analysis

```go
// Analyzer determines if code can be parallelized safely
type SafetyAnalyzer struct {
    // Tracks shared state access
    globalReads   map[string]bool
    globalWrites  map[string]bool
    staticAccess  map[string]bool
    fileIO        bool
    networkIO     bool
    sideEffects   bool
}

func (sa *SafetyAnalyzer) CanParallelize(fn *Function) bool {
    // Analyze function for:
    // - No global state modifications
    // - No static variable access
    // - No file/network I/O (or marked as safe)
    // - No observable side effects
    // Return true only if completely safe

    return !sa.globalWrites && !sa.staticAccess && !sa.fileIO
}
```

### 2. Parallel Execution

```go
// ParallelExecutor runs PHP code in parallel
type ParallelExecutor struct {
    vm   *VM
    pool *WorkerPool
}

// Execute multiple requests in parallel
func (pe *ParallelExecutor) ExecuteRequests(requests []*Request) []*Response {
    responses := make([]*Response, len(requests))
    var wg sync.WaitGroup

    for i, req := range requests {
        wg.Add(1)
        go func(idx int, r *Request) {
            defer wg.Done()
            responses[idx] = pe.executeRequest(r)
        }(i, req)
    }

    wg.Wait()
    return responses
}

// Automatically parallelize array operations
func (pe *ParallelExecutor) ArrayMap(arr *Array, fn *Function) *Array {
    if !pe.canParallelizeArrayOp(fn) {
        // Fall back to sequential
        return sequentialArrayMap(arr, fn)
    }

    // Parallel implementation
    result := NewArray()
    var mu sync.Mutex
    var wg sync.WaitGroup

    for key, value := range arr.Iterate() {
        wg.Add(1)
        go func(k ArrayKey, v *Value) {
            defer wg.Done()
            mapped := pe.vm.Call(fn, []*Value{v})
            mu.Lock()
            result.Set(k, mapped)
            mu.Unlock()
        }(key, value)
    }

    wg.Wait()
    return result
}
```

### 3. Request-Level Parallelism

```go
// RequestContext isolates each request
type RequestContext struct {
    vm        *VM      // Separate VM instance
    globals   *Globals // Request-specific globals
    output    *Buffer  // Output buffer
    session   *Session // Session data
}

// Server handles multiple requests concurrently
type Server struct {
    vmPool chan *VM    // Pool of VM instances
}

func (s *Server) HandleRequest(req *http.Request) *http.Response {
    vm := <-s.vmPool
    defer func() { s.vmPool <- vm }()

    ctx := NewRequestContext(vm, req)
    result := vm.Execute(ctx)
    return result.ToHTTPResponse()
}
```

**Design Decision**: Share-nothing by default, explicit sharing when safe

**Rationale**:
- Matches PHP-FPM model (familiar)
- Avoids most concurrency issues
- Easy to reason about
- Opt-in for advanced parallelism

## Go Integration

### 1. FFI System

```go
// Call Go functions from PHP
type GoFunction struct {
    name    string
    fn      interface{}    // Go function
    params  []TypeInfo     // Parameter types
    returns []TypeInfo     // Return types
}

// Register Go function
func RegisterGoFunction(name string, fn interface{}) {
    // Parse function signature using reflection
    // Create PHP-callable wrapper
    // Add to function registry
}

// PHP code can call: go_call('mypackage.MyFunc', $arg1, $arg2)
func goCall(vm *VM, args []*Value) (*Value, error) {
    funcName := args[0].ToString()
    goFn := lookupGoFunction(funcName)

    // Marshal PHP values to Go types
    goArgs := marshalToGo(args[1:], goFn.params)

    // Call Go function
    goResults := goFn.Call(goArgs)

    // Marshal Go return values to PHP
    phpResults := marshalToPHP(goResults)

    return phpResults, nil
}
```

### 2. Native Extensions

```go
// Extension interface
type Extension interface {
    Name() string
    Version() string
    Functions() map[string]*Function
    Classes() map[string]*Class
    Constants() map[string]*Value
    Init(*VM) error
}

// Example extension in Go
type JSONExtension struct{}

func (e *JSONExtension) Name() string { return "json" }

func (e *JSONExtension) Functions() map[string]*Function {
    return map[string]*Function{
        "json_encode": {
            Name: "json_encode",
            Handler: func(vm *VM, args []*Value) (*Value, error) {
                // Implementation
                return jsonEncode(args[0])
            },
        },
        "json_decode": {
            Name: "json_decode",
            Handler: func(vm *VM, args []*Value) (*Value, error) {
                return jsonDecode(args[0])
            },
        },
    }
}

// Register extension
vm.RegisterExtension(&JSONExtension{})
```

### 3. Type Marshaling

```go
// Marshal between PHP and Go types
type Marshaler struct{}

func (m *Marshaler) ToPHP(goVal interface{}) (*Value, error) {
    switch v := goVal.(type) {
    case nil:
        return NewNull(), nil
    case bool:
        return NewBool(v), nil
    case int, int64:
        return NewInt(v.(int64)), nil
    case float64:
        return NewFloat(v), nil
    case string:
        return NewString(v), nil
    case []interface{}:
        arr := NewArray()
        for _, item := range v {
            val, _ := m.ToPHP(item)
            arr.Append(val)
        }
        return NewValue(arr), nil
    case map[string]interface{}:
        arr := NewArray()
        for k, item := range v {
            val, _ := m.ToPHP(item)
            arr.Set(StringKey(k), val)
        }
        return NewValue(arr), nil
    default:
        return nil, fmt.Errorf("unsupported type: %T", goVal)
    }
}

func (m *Marshaler) ToGo(phpVal *Value) (interface{}, error) {
    switch phpVal.Type() {
    case TypeNull:
        return nil, nil
    case TypeBool:
        return phpVal.ToBool(), nil
    case TypeInt:
        return phpVal.ToInt(), nil
    case TypeFloat:
        return phpVal.ToFloat(), nil
    case TypeString:
        return phpVal.ToString(), nil
    case TypeArray:
        // Convert to map[string]interface{} or []interface{}
        // depending on array keys
        return m.arrayToGo(phpVal.ToArray()), nil
    default:
        return nil, fmt.Errorf("cannot marshal type: %s", phpVal.Type())
    }
}
```

## Error Handling

```go
// PHP errors and exceptions
type Error struct {
    level   ErrorLevel
    message string
    file    string
    line    int
    trace   *StackTrace
}

type ErrorLevel int

const (
    E_ERROR ErrorLevel = 1 << iota
    E_WARNING
    E_PARSE
    E_NOTICE
    E_CORE_ERROR
    E_CORE_WARNING
    E_COMPILE_ERROR
    E_COMPILE_WARNING
    E_USER_ERROR
    E_USER_WARNING
    E_USER_NOTICE
    E_STRICT
    E_RECOVERABLE_ERROR
    E_DEPRECATED
    E_USER_DEPRECATED
    E_ALL = 0xFFFF
)

// Exception is PHP's exception
type Exception struct {
    message  string
    code     int64
    file     string
    line     int
    trace    *StackTrace
    previous *Exception
}

// Error handler
func (vm *VM) HandleError(err *Error) {
    if vm.errorHandler != nil {
        vm.errorHandler(err)
    } else {
        // Default: log and continue or panic
        if err.level & E_ERROR != 0 {
            panic(err)  // Fatal error
        } else {
            log.Println(err)  // Warning/notice
        }
    }
}
```

## Testing Strategy

### 1. Unit Tests

```go
// Test each component in isolation
func TestArrayAppend(t *testing.T) {
    arr := NewArray()
    arr.Append(NewInt(1))
    arr.Append(NewInt(2))

    assert.Equal(t, 2, arr.Count())
    assert.Equal(t, int64(1), arr.Get(IntKey(0)).ToInt())
}
```

### 2. Integration Tests

```go
// Test complete execution
func TestExecutePHPScript(t *testing.T) {
    vm := NewVM()
    code := `<?php
    $x = 1 + 2;
    echo $x;
    `

    result, err := vm.ExecuteString(code)
    assert.NoError(t, err)
    assert.Equal(t, "3", result.Output())
}
```

### 3. Compatibility Tests

```go
// Use PHP's .phpt test suite
func TestPHPCompatibility(t *testing.T) {
    tests := loadPHPTests("../php-src/tests")

    for _, test := range tests {
        t.Run(test.Name, func(t *testing.T) {
            result := runPHPTest(test)
            assert.Equal(t, test.Expected, result)
        })
    }
}
```

## Performance Considerations

### 1. Hot Path Optimization

- Inline simple opcodes
- Use switch for dispatch (Go optimizes well)
- Pool common objects (empty arrays, common integers)
- Cache method lookups

### 2. Escape Analysis

- Keep values on stack when possible
- Avoid unnecessary allocations
- Use value types for small objects

### 3. Profiling

```go
import _ "net/http/pprof"

// Enable profiling
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Then use: go tool pprof http://localhost:6060/debug/pprof/profile
```

## Build and Deployment

### 1. Go Modules

```go
module github.com/yourusername/php-go

go 1.21

require (
    github.com/pkg/errors v0.9.1
    // Other dependencies
)
```

### 2. CLI Tool

```bash
# Install
go install github.com/yourusername/php-go/cmd/php-go@latest

# Run
php-go script.php
php-go -S localhost:8000  # Built-in server
```

### 3. Distribution

- Single binary (no dependencies)
- Cross-compilation for multiple platforms
- Docker images available
- Homebrew/apt packages

## Migration from PHP

### 1. Drop-in Replacement

```bash
# Replace PHP binary
alias php='php-go'

# Or update scripts
#!/usr/bin/env php-go
```

### 2. Feature Flags

```go
// Enable experimental features
php-go --enable-parallel script.php

// Compatibility mode
php-go --compat=php8.4 script.php
```

### 3. Performance Monitoring

```go
// Built-in profiling
php-go --profile script.php

// Generates profile.pb.gz for analysis
```

## Phase 6 Implementation Details

### Compiler Improvements (Phase 6A-6B)

#### 1. Symbol Table Enhancements

**Variable Naming Conflict Resolution**:
```go
// Modified variable resolution to distinguish between
// builtin functions and user variables
func (c *Compiler) compileVariable(node *ast.Variable) {
    sym, ok := c.symbolTable.Resolve(node.Name)

    // If resolved symbol is a builtin, treat as undefined variable
    if ok && sym.Scope == BuiltinScope {
        // Define new variable instead of referencing builtin
        sym = c.symbolTable.Define(node.Name)
    }

    // Emit FETCH_VAR with proper symbol
    c.emit(OpFetchVar, SymbolOperand(sym), ResultOperand())
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 2563-2585, 2587-2643, 2645-2680
- Applied to: variable fetch, assignment, foreach, increment/decrement, isset, empty, unset, catch

#### 2. Temporary Variable Management

**Foreach Iterator Protection**:
```go
// Save temp stack level before loop body
func (c *Compiler) compileForeach(node *ast.ForeachStatement) {
    // ... setup iterator ...

    // Save current temp level before entering loop body
    savedTempLevel := c.numTemps

    // Compile loop body
    c.compileStatement(node.Body)

    // Don't free temps allocated for iterator
    c.numTemps = savedTempLevel

    // ... continue iteration ...
}
```

**ExpressionStatement Fix**:
```go
// Only free temps allocated during this expression
func (c *Compiler) compileExpressionStatement(node *ast.ExpressionStatement) {
    savedTempLevel := c.numTemps

    c.Compile(node.Expression)

    // Free only newly allocated temps
    for i := savedTempLevel; i < c.numTemps; i++ {
        c.FreeTemp(i)
    }
    c.numTemps = savedTempLevel
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 317-335 (ExpressionStatement)
- `pkg/compiler/compiler.go` - Lines 2115-2285 (Foreach compilation)

#### 3. Control Flow Fixes

**If Statement Condition Evaluation**:
```go
// Fixed: Use condition result temp instead of hardcoded TmpVarOperand(0)
func (c *Compiler) compileIf(node *ast.IfStatement) {
    // Evaluate condition and capture its result temp
    condTemp := c.Compile(node.Condition)

    // Jump to else/end if condition is false
    // BEFORE: c.emit(OpJmpZ, vm.TmpVarOperand(0), ...)  // WRONG!
    // AFTER:
    c.emit(OpJmpZ, condTemp, JumpTarget(elseLabel))
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 1877, 1911

**Do-While Implementation**:
```go
func (c *Compiler) compileDoWhile(node *ast.DoWhileStatement) {
    startLabel := c.NewLabel()
    endLabel := c.NewLabel()

    // Enter loop context for break/continue
    c.EnterLoop(endLabel, startLabel)
    defer c.ExitLoop()

    // Mark start (body executes first in do-while)
    c.MarkLabel(startLabel)

    // Compile body
    c.compileStatement(node.Body)

    // Evaluate condition
    condResult := c.Compile(node.Condition)

    // Jump back to start if condition is true
    c.emit(OpJmpNZ, condResult, JumpTarget(startLabel))

    // Mark end
    c.MarkLabel(endLabel)
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 1954-1980

#### 4. PHP 7+ Language Features

**Null Coalescing Operator (??)**:
```go
// Compiler support
func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) vm.Operand {
    switch node.Operator {
    case "??":
        // Emit OpCoalesce opcode
        left := c.Compile(node.Left)
        right := c.Compile(node.Right)
        result := c.AllocTemp()
        c.emit(OpCoalesce, left, right, result)
        return result
    // ... other operators ...
    }
}

// VM handler
func (vm *VM) opCoalesce(frame *Frame, instr Instruction) error {
    left := vm.fetchOperand(frame, instr.op1)

    // Return left if not null/undef, otherwise right
    if left.IsNull() || left.IsUndef() {
        right := vm.fetchOperand(frame, instr.op2)
        vm.storeOperand(frame, instr.result, right)
    } else {
        vm.storeOperand(frame, instr.result, left)
    }
    return nil
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Line 1671 (compiler case)
- `pkg/vm/handlers_comparison.go` - Lines 136-152 (VM handler)
- `pkg/vm/vm.go` - Line 156 (dispatch case)

**Class Name Expression (::class)**:
```go
func (c *Compiler) compileClassNameExpression(node *ast.ClassNameExpression) vm.Operand {
    // Extract class name from identifier
    var className string
    if ident, ok := node.Class.(*ast.Identifier); ok {
        className = ident.Value
    } else if fqn, ok := node.Class.(*ast.FullyQualifiedName); ok {
        className = fqn.Parts[len(fqn.Parts)-1]
    }

    // Add as constant and return as string value
    constIdx := c.addConstant(types.NewString(className))
    result := c.AllocTemp()
    c.emit(OpAssign, vm.ConstOperand(constIdx), result)
    return result
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 1835-1863

**Array Destructuring in Foreach**:
```go
func (c *Compiler) compileForeachWithDestructuring(
    arrayExpr *ast.ArrayExpression,
    iterValue vm.Operand,
) {
    // For each element in the destructuring pattern
    for i, elem := range arrayExpr.Elements {
        // Determine key (explicit or numeric index)
        var keyOperand vm.Operand
        if elem.Key != nil {
            keyOperand = c.Compile(elem.Key)
        } else {
            keyOperand = vm.ConstOperand(c.addConstant(types.NewInt(int64(i))))
        }

        // Fetch element: $temp = $iterValue[$key]
        elementTemp := c.AllocTemp()
        c.emit(OpFetchDimR, iterValue, keyOperand, elementTemp)

        // Assign to target variable
        if varNode, ok := elem.Value.(*ast.Variable); ok {
            symbol := c.symbolTable.Define(varNode.Name)
            c.emit(OpAssign, elementTemp, vm.SymbolOperand(symbol))
        }
    }
}
```

**Implementation Files**:
- `pkg/compiler/compiler.go` - Lines 2184-2228

### Virtual Machine Enhancements (Phase 6A)

#### DECLARE_CLASS Opcode Handler

```go
func (vm *VM) opDeclareClass(frame *Frame, instr Instruction) error {
    // Fetch class name from constants
    className := vm.fetchOperand(frame, instr.op1).AsString()

    // Check for parent class (optional)
    var parentClass *types.ClassEntry
    if instr.ext != 0 {
        parentName := vm.fetchOperand(frame, instr.op2).AsString()
        parentClass = vm.classes[parentName]
    }

    // Create and register class entry
    class := &types.ClassEntry{
        Name:       className,
        Parent:     parentClass,
        Properties: make(map[string]*types.PropertyEntry),
        Methods:    make(map[string]*types.MethodEntry),
        Constants:  make(map[string]types.Value),
    }

    // Apply inheritance if parent exists
    if parentClass != nil {
        class.InheritFrom(parentClass)
    }

    vm.classes[className] = class
    return nil
}
```

**Implementation Files**:
- `pkg/vm/handlers_object.go` - Lines 289-324
- `pkg/vm/vm.go` - Line 108 (dispatch case)

### Testing Infrastructure (Phase 6C)

#### Example Test Suite

```go
// tests/examples_test.go
func TestBasicExamples(t *testing.T) {
    examples := []string{
        "arrays", "control_flow", "expressions",
        "functions", "hello", "strings", "variables",
    }

    for _, name := range examples {
        t.Run(name, func(t *testing.T) {
            // Run php-go on example
            output := runPhpGo("examples/basic/" + name + ".php")

            // Read expected output
            expected := readFile("examples/basic/" + name + ".expected")

            // Compare
            if output != expected {
                t.Errorf("Output mismatch for %s\nGot:\n%s\nExpected:\n%s",
                    name, output, expected)
            }
        })
    }
}
```

#### Regression Test Suite

```go
// pkg/compiler/regression_test.go
// 56 test cases organized by phase:

func TestRegression_Phase6A1_VariableNaming(t *testing.T) {
    // 6 tests for variables with builtin names
}

func TestRegression_Phase6A2_ForeachIncrement(t *testing.T) {
    // 4 tests for foreach with operations
}

func TestRegression_Phase6A3_DeclareClass(t *testing.T) {
    // 2 tests for class declaration
}

func TestRegression_Phase6A4_IfConditions(t *testing.T) {
    // 4 tests for if statement evaluation
}

func TestRegression_Phase6A4_IncrementDecrement(t *testing.T) {
    // 7 tests for ++ and -- operators
}

func TestRegression_Phase6B1_NullCoalescing(t *testing.T) {
    // 6 tests for ?? operator
}

func TestRegression_Phase6B2_DoWhile(t *testing.T) {
    // 4 tests for do-while loops
}

func TestRegression_Phase6B3_ClassNameExpression(t *testing.T) {
    // 4 tests for ClassName::class
}

func TestRegression_Phase6B4_ForeachDestructuring(t *testing.T) {
    // 4 tests for array destructuring
}

func TestRegression_ComplexInteractions(t *testing.T) {
    // 5 tests for combinations of features
}
```

**Implementation Files**:
- `tests/examples_test.go` - Example validation tests
- `pkg/compiler/regression_test.go` - 56 regression test cases

### Key Design Decisions

#### 1. Temporary Variable Scoping
**Decision**: Track temp allocation per expression, not globally
**Rationale**: Prevents premature freeing of iterator temps during loop bodies
**Trade-off**: Slightly more complex compiler logic, but correct semantics

#### 2. Builtin vs User Symbols
**Decision**: Check symbol scope before treating as reference
**Rationale**: PHP allows variables to shadow builtin function names
**Implementation**: Applied at all variable resolution points consistently

#### 3. Opcode Handler Organization
**Decision**: Group handlers by category in separate files
**Rationale**: Better maintainability, easier to find related opcodes
**Files**:
- `handlers_comparison.go` - Comparisons and null coalescing
- `handlers_object.go` - Object operations including DECLARE_CLASS
- `handlers_array.go` - Array operations including foreach iterators

## Summary

This architecture leverages Go's strengths:
- Goroutines for parallelization
- GC for memory management
- Strong typing for correctness
- Simple deployment model
- Excellent tooling

While maintaining PHP compatibility through:
- Accurate type juggling
- Full language feature support
- Compatible standard library
- PHP semantics preserved

Phase 6 additions demonstrate:
- Careful temp variable management
- Proper symbol table scoping
- Comprehensive regression testing
- PHP 7+ modern feature support

---

**Last Updated**: 2025-11-24
**Status**: Phase 6A-6C Complete (604/1050 hours, 58%)
