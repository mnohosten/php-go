# VM API Reference

Package: `github.com/krizos/php-go/pkg/vm`

## Overview

The vm package implements the PHP Virtual Machine that executes compiled bytecode. It provides a stack-based execution engine with support for functions, objects, arrays, exceptions, generators, and all PHP 8.4 features. The VM manages call frames, globals, output buffering, and dispatches opcodes to their handlers.

## Main Types

### VM

The PHP Virtual Machine that executes bytecode.

```go
type VM struct {
    // Private fields
}
```

**Constructor:**

```go
func New() *VM
```

Creates a new virtual machine with initialized state (empty globals, function registry, class registry, frame stack).

```go
func NewWithBytecode(instructions Instructions, constants []interface{}) *VM
```

Creates a new VM and loads bytecode, creating a main frame ready for execution.

**Methods:**

```go
func (vm *VM) Execute(instructions Instructions) error
```

Executes bytecode starting from the main program. This is the main entry point for running compiled PHP code.

```go
func (vm *VM) ExecuteWithCVs(instructions Instructions, numCVs int) error
```

Executes bytecode with a specified number of compiled variables.

```go
func (vm *VM) LoadConstants(constants []interface{})
```

Loads the constant table from compiled bytecode.

**Global Variables:**

```go
func (vm *VM) SetGlobal(name string, value *types.Value)
func (vm *VM) GetGlobal(name string) (*types.Value, bool)
```

Set and get global variables.

**Functions:**

```go
func (vm *VM) RegisterFunction(name string, fn *CompiledFunction)
func (vm *VM) GetFunction(name string) (*CompiledFunction, bool)
```

Register and retrieve compiled functions.

**Classes:**

```go
func (vm *VM) RegisterClass(name string, class *CompiledClass)
func (vm *VM) GetClass(name string) (*CompiledClass, bool)
```

Register and retrieve compiled classes.

**Output:**

```go
func (vm *VM) Output() []byte
func (vm *VM) OutputString() string
```

Get the accumulated output from echo/print statements.

**State:**

```go
func (vm *VM) Reset()
```

Resets the VM state (clears frames, output, exit flag).

### CompiledFunction

Represents a compiled PHP function ready for execution.

```go
type CompiledFunction struct {
    Name         string
    Instructions Instructions
    NumLocals    int
    NumParams    int
    NumCVs       int
}
```

**Fields:**
- `Name`: Function name
- `Instructions`: Bytecode instructions for the function body
- `NumLocals`: Number of local variables (CVs + temps)
- `NumParams`: Number of parameters the function accepts
- `NumCVs`: Number of compiled variables (user-defined variables)

### Closure

Represents a PHP closure/anonymous function with captured variables.

```go
type Closure struct {
    Function     *CompiledFunction
    CapturedVars map[string]*types.Value
    Static       bool
    ReturnByRef  bool
}
```

**Fields:**
- `Function`: The underlying compiled function
- `CapturedVars`: Variables captured from the parent scope via `use`
- `Static`: If true, closure cannot access `$this`
- `ReturnByRef`: If true, function returns by reference

### Frame

Represents a call frame (activation record) on the call stack.

```go
type Frame struct {
    // Private fields
}
```

**Constructor:**

```go
func NewFrame(fn *CompiledFunction) *Frame
```

Creates a new call frame for a function.

**Frame Management (internal methods):**

```go
func (vm *VM) pushFrame(frame *Frame)
func (vm *VM) popFrame() *Frame
func (vm *VM) currentFrame() *Frame
```

### Instructions

A slice of Instruction structs representing bytecode.

```go
type Instructions []Instruction
```

**Methods:**

```go
func (ins Instructions) String() string
```

Returns a human-readable disassembly of the instructions.

### Instruction

A single bytecode instruction (24 bytes, fixed size).

```go
type Instruction struct {
    Opcode        Opcode
    Lineno        uint32
    ExtendedValue uint32
    Op1           Operand
    Op2           Operand
    Result        Operand
}
```

**Fields:**
- `Opcode`: The operation to perform (see Opcode section)
- `Lineno`: Source line number for debugging/stack traces
- `ExtendedValue`: Additional flags/data for certain opcodes
- `Op1`, `Op2`: Input operands
- `Result`: Output operand (where result is stored)

**Methods:**

```go
func (i Instruction) String() string
```

Returns disassembled instruction as string.

### Opcode

Represents a VM opcode (operation code).

```go
type Opcode uint8
```

**Opcode Categories:**

**Arithmetic Operations:**
- `OpAdd`, `OpSub`, `OpMul`, `OpDiv`, `OpMod`, `OpPow`

**Comparison Operations:**
- `OpIsEqual`, `OpIsNotEqual`, `OpIsIdentical`, `OpIsNotIdentical`
- `OpIsSmaller`, `OpIsSmallerOrEqual`, `OpSpaceship`

**Bitwise Operations:**
- `OpBWAnd`, `OpBWOr`, `OpBWXor`, `OpBWNot`
- `OpSL` (shift left), `OpSR` (shift right)

**Logical Operations:**
- `OpBoolNot`

**Constants:**
- `OpFetchConstant` - Load constant from constant table

**Variables:**
- `OpAssign` - Assign value to variable
- `OpFetchR` - Read variable value
- `OpFree` - Free temporary variable
- `OpBindGlobal` - Bind local variable to global scope
- `OpUnsetVar` - Unset variable
- `OpIssetIsemptyVar` - Check if variable is set/empty
- `OpQMAssign` - Null coalescing assignment

**Control Flow:**
- `OpJmp` - Unconditional jump
- `OpJmpZ` - Jump if zero (false)
- `OpJmpNZ` - Jump if not zero (true)

**Functions:**
- `OpInitFcall` - Initialize function call
- `OpInitFcallByName` - Initialize call by name
- `OpSendVal` - Send argument value
- `OpDoFcall` - Execute function call
- `OpDoUcall` - Execute user function call
- `OpDoIcall` - Execute internal function call
- `OpReturn` - Return from function
- `OpRecv` - Receive parameter
- `OpRecvInit` - Receive parameter with default value

**Arrays:**
- `OpInitArray` - Create empty array
- `OpAddArrayElement` - Add element to array
- `OpFetchDimR` - Read array element
- `OpFetchDimW` - Fetch for write
- `OpFetchDimRW` - Fetch for read-write
- `OpFetchDimIs` - Fetch with isset check
- `OpFetchDimFuncArg` - Fetch for function argument
- `OpFetchDimUnset` - Fetch for unset
- `OpAssignDim` - Assign to array element
- `OpAssignDimOp` - Compound assignment to array element
- `OpUnsetDim` - Unset array element
- `OpIssetIsemptyDimObj` - Check if array element is set/empty
- `OpCount` - Count array elements
- `OpInArray` - Check if value is in array
- `OpArrayKeyExists` - Check if key exists

**Objects:**
- `OpNew` - Create new object
- `OpClone` - Clone object
- `OpInstanceof` - Type check
- `OpGetClass` - Get object class name
- `OpFetchClassName` - Fetch ::class constant
- `OpFetchThis` - Fetch $this
- Property operations: `OpFetchObjR`, `OpFetchObjW`, `OpFetchObjRW`, etc.
- Assignment: `OpAssignObj`, `OpAssignObjOp`, `OpAssignObjRef`
- Methods: `OpInitMethodCall`, `OpInitStaticMethodCall`
- Increment/decrement: `OpPreIncObj`, `OpPostIncObj`, `OpPreDecObj`, `OpPostDecObj`
- Unset/Isset: `OpUnsetObj`, `OpIssetIsemptyPropObj`

**Strings:**
- `OpConcat` - String concatenation
- `OpFastConcat` - Optimized concatenation

**I/O:**
- `OpEcho` - Output to stdout

**Foreach:**
- `OpFeResetR` - Reset foreach iterator (read)
- `OpFeResetRW` - Reset foreach iterator (read-write)
- `OpFeFetchR` - Fetch next iteration (read)
- `OpFeFetchRW` - Fetch next iteration (read-write)
- `OpFeFree` - Free foreach iterator

**Closures:**
- `OpDeclareFunction` - Declare regular function
- `OpDeclareLambdaFunction` - Declare closure
- `OpBindLexical` - Bind closure variable

**Generators:**
- `OpGeneratorCreate` - Create generator
- `OpYield` - Yield value
- `OpYieldFrom` - Yield from another generator
- `OpGeneratorReturn` - Return from generator

**Exceptions:**
- `OpThrow` - Throw exception
- `OpBeginSilence` - Begin @ error suppression
- `OpEndSilence` - End @ error suppression

**Include:**
- `OpIncludeOrEval` - Include/require files or eval code

Total: 210 opcodes defined

### Operand

Represents an operand in an instruction.

```go
type Operand uint32
```

**Operand Types:**
- Constant: `ConstOperand(idx)` - References constant table
- Compiled Variable (CV): `CVOperand(idx)` - User-defined variable
- Temporary Variable: `TmpVarOperand(idx)` - Compiler-generated temporary
- Unused: `UnusedOperand()` - No operand

**Helper Functions:**

```go
func ConstOperand(idx uint32) Operand
func CVOperand(idx uint32) Operand
func TmpVarOperand(idx uint32) Operand
func UnusedOperand() Operand
```

**Methods:**

```go
func (o Operand) Type() OperandType
func (o Operand) Value() uint32
func (o Operand) String() string
```

## Execution Model

### Stack-Based Execution

The VM uses a stack-based execution model:

1. **Call Frames**: Each function call creates a new frame
2. **Local Variables**: Stored in frame's local variable array
3. **Temporary Variables**: Used for expression evaluation
4. **Stack Depth Limit**: Default 1000 frames (configurable)

### Instruction Dispatch

The VM executes instructions in a loop:

```go
for {
    instr := currentFrame.instructions[ip]
    dispatch(instr)
    ip++
}
```

Each opcode has a dedicated handler function in `handlers_*.go`:
- `handlers_arithmetic.go` - Math operations
- `handlers_comparison.go` - Comparisons
- `handlers_logic.go` - Boolean/bitwise operations
- `handlers_variables.go` - Variable operations
- `handlers_control.go` - Jumps
- `handlers_functions.go` - Function calls
- `handlers_array.go` - Array operations
- `handlers_object.go` - Object operations
- `handlers_strings.go` - String concatenation
- `handlers_io.go` - Output operations
- `handlers_foreach.go` - Foreach loops
- `handlers_closure.go` - Closures
- `handlers_generator.go` - Generators
- `handlers_exception.go` - Exception handling

### Frame Management

Function calls:
1. Push new frame onto frame stack
2. Copy arguments to new frame's locals
3. Execute function body
4. Pop frame and return value

### Error Handling

VM errors include:
- Stack overflow (too many nested calls)
- Undefined variable access
- Type errors in operations
- Division by zero
- Array key not found
- Call to undefined function

## Built-in Functions

The VM includes built-in functions registered in `builtins.go`:

**Variable Functions:**
- `isset()`, `empty()`, `unset()`
- `var_dump()`, `print_r()`
- `gettype()`, `get_defined_vars()`

**Array Functions:**
- `count()`, `in_array()`, `array_key_exists()`

**String Functions:**
- `strlen()`, `str_repeat()`, `str_pad()`

**Output Functions:**
- `echo`, `print`

More functions are provided by the stdlib package (Phase 6).

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/compiler"
    "github.com/krizos/php-go/pkg/lexer"
    "github.com/krizos/php-go/pkg/parser"
    "github.com/krizos/php-go/pkg/vm"
)

func main() {
    input := `<?php
    function greet($name) {
        return "Hello, " . $name;
    }

    echo greet("World");
    `

    // Parse
    l := lexer.New(input, "test.php")
    p := parser.New(l)
    program := p.ParseProgram()

    // Compile
    c := compiler.New()
    c.Compile(program)
    bytecode := c.Bytecode()

    // Execute
    machine := vm.New()
    machine.LoadConstants(bytecode.Constants)
    err := machine.ExecuteWithCVs(bytecode.Instructions, bytecode.NumCVs)

    if err != nil {
        fmt.Println("Runtime error:", err)
        return
    }

    fmt.Println("Output:", machine.OutputString())
}
```

## Performance Considerations

- **Fixed-size Instructions**: 24 bytes per instruction for cache efficiency
- **Pre-allocated Frame Stack**: 1024 frames pre-allocated to avoid allocations
- **Efficient Dispatch**: Direct switch on opcode (no table lookup)
- **Value Pooling**: Reuses Value structs to reduce GC pressure (in types package)
- **Optimized Handlers**: Specialized handlers for common operations

## Implementation Notes

- The VM is **not** thread-safe (each goroutine needs its own VM instance)
- Frame stack grows dynamically if needed
- Output is buffered in memory (use runtime package for output buffering control)
- Exit/die terminates execution immediately
- Generators use special frame management for resumable execution
- Exception handling uses try/catch frame tracking

## Debugging Support

- Line number tracking in every instruction
- Stack trace generation on errors
- Instruction disassembly via `Instructions.String()`
- Frame inspection capabilities

## Coverage

The VM package has ~79% code coverage with comprehensive tests for all opcode handlers.
