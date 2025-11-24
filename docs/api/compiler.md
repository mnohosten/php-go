# Compiler API Reference

Package: `github.com/krizos/php-go/pkg/compiler`

## Overview

The compiler package compiles PHP Abstract Syntax Trees (AST) into bytecode for execution by the virtual machine. It performs constant folding, dead code elimination, and manages symbol tables for variable scoping.

## Main Types

### Compiler

The main compiler that transforms AST nodes into VM bytecode.

```go
type Compiler struct {
    // Private fields
}
```

**Constructor:**

```go
func New() *Compiler
```

Creates a new compiler instance with initialized symbol tables and empty instruction stream.

**Methods:**

```go
func (c *Compiler) Compile(node ast.Node) error
```

Compiles an AST node into bytecode. Returns an error if compilation fails. This is the main entry point for compilation.

```go
func (c *Compiler) Bytecode() *Bytecode
```

Returns the compiled bytecode including instructions and constants table.

```go
func (c *Compiler) AddConstant(value interface{}) int
```

Adds a constant to the constant table and returns its index. Deduplicates identical constants.

```go
func (c *Compiler) GetConstant(idx int) (interface{}, error)
```

Retrieves a constant by its index.

```go
func (c *Compiler) Constants() []interface{}
```

Returns a copy of the constant table.

**Instruction Emission:**

```go
func (c *Compiler) Emit(opcode vm.Opcode, operands ...vm.Operand) int
```

Emits an instruction with the given opcode and operands. Returns the position of the emitted instruction.

```go
func (c *Compiler) EmitWithLine(opcode vm.Opcode, lineno uint32, operands ...vm.Operand) int
```

Emits an instruction with line number information for debugging/stack traces.

```go
func (c *Compiler) EmitWithExtended(opcode vm.Opcode, lineno uint32, extended uint32, operands ...vm.Operand) int
```

Emits an instruction with extended value field for additional data.

**Instruction Manipulation:**

```go
func (c *Compiler) ReplaceInstruction(pos int, instr vm.Instruction) error
```

Replaces an instruction at a specific position (used for jump patching).

```go
func (c *Compiler) ChangeOperand(pos int, operandNum int, operand vm.Operand) error
```

Changes a specific operand in an instruction.

```go
func (c *Compiler) CurrentPosition() int
```

Returns the position where the next instruction will be emitted.

```go
func (c *Compiler) LastInstructionIs(opcode vm.Opcode) bool
```

Checks if the last emitted instruction has the given opcode.

```go
func (c *Compiler) RemoveLastInstruction()
```

Removes the last emitted instruction (used for optimization).

**Temporary Variables:**

```go
func (c *Compiler) AllocTemp() vm.Operand
```

Allocates a new temporary variable and returns its operand.

```go
func (c *Compiler) FreeTemp()
```

Deallocates the most recently allocated temporary variable.

```go
func (c *Compiler) CurrentTemp() vm.Operand
```

Returns the currently active temporary variable operand.

### Bytecode

Represents the compiled bytecode program.

```go
type Bytecode struct {
    Instructions vm.Instructions
    Constants    []interface{}
    NumCVs       int
}
```

**Fields:**
- `Instructions`: The compiled bytecode instructions
- `Constants`: Literal values referenced by the bytecode
- `NumCVs`: Number of compiled variables (for proper TMPVAR offset)

### SymbolTable

Manages variable scoping and resolution.

```go
type SymbolTable struct {
    // Private fields
}
```

**Methods:**

```go
func (c *Compiler) InitSymbolTable()
```

Initializes the symbol table with built-in functions.

```go
func (st *SymbolTable) Define(name string) Symbol
```

Defines a new symbol in the current scope.

```go
func (st *SymbolTable) Resolve(name string) (Symbol, bool)
```

Resolves a symbol by name, searching from current scope upward.

```go
func (st *SymbolTable) NumDefinitions() int
```

Returns the number of symbols defined in the current scope.

```go
func (c *Compiler) ResolveVariable(name string) (Symbol, bool)
```

Helper method to resolve a variable symbol.

### Symbol

Represents a variable or function symbol.

```go
type Symbol struct {
    Name  string
    Scope SymbolScope
    Index int
}
```

**Symbol Scopes:**
- `GlobalScope` - Global variables
- `LocalScope` - Function/method local variables
- `BuiltinScope` - Built-in functions
- `FreeScope` - Closure captured variables
- `FunctionScope` - Function names

### LoopContext

Tracks information about loops for break/continue compilation.

```go
type LoopContext struct {
    StartPos        int
    BreakJumps      []int
    ContinueJumps   []int
}
```

### EmittedInstruction

Tracks metadata about emitted instructions.

```go
type EmittedInstruction struct {
    Opcode   vm.Opcode
    Position int
}
```

## Compilation Process

The compiler works in several phases:

### 1. Symbol Registration

Before compiling, built-in functions are registered in the symbol table:
- `echo`, `print`, `var_dump`
- `isset`, `empty`, `unset`
- `count`, `strlen`
- And more as implemented

### 2. AST Traversal

The compiler recursively walks the AST, emitting instructions for each node:

**Statements:**
- Control flow (if/while/for/foreach/switch)
- Jump statements (return/break/continue)
- Echo and output
- Global/namespace declarations
- Function and class declarations

**Expressions:**
- Literals → OpFetchConstant
- Variables → OpFetchR, OpAssign
- Binary operations → OpAdd, OpSub, OpMul, etc.
- Function calls → OpInitFcall, OpSendVal, OpDoFcall
- Array access → OpFetchDimR, OpAssignDim
- Object operations → OpNew, OpFetchObjR, OpInitMethodCall

### 3. Constant Folding

The compiler performs constant folding optimization:
```php
$x = 2 + 3;  // Compiled as: $x = 5;
```

### 4. Dead Code Elimination

Code after return statements is not compiled:
```php
function foo() {
    return 42;
    echo "never executed";  // Not compiled
}
```

### 5. Jump Patching

For control flow, the compiler:
1. Emits placeholder jumps with operand 0
2. Records jump positions
3. Patches jumps with correct targets after compiling the block

## Opcodes Generated

The compiler emits VM opcodes defined in `pkg/vm/opcodes.go`. Common opcodes include:

**Constants and Variables:**
- `OpFetchConstant` - Load constant
- `OpFetchR` - Read variable
- `OpAssign` - Assign to variable
- `OpBindGlobal` - Bind local to global

**Arithmetic:**
- `OpAdd`, `OpSub`, `OpMul`, `OpDiv`, `OpMod`, `OpPow`

**Comparison:**
- `OpIsEqual`, `OpIsIdentical`, `OpIsSmaller`, `OpSpaceship`

**Logical and Bitwise:**
- `OpBoolNot`, `OpBWAnd`, `OpBWOr`, `OpBWXor`, `OpBWNot`

**Control Flow:**
- `OpJmp` - Unconditional jump
- `OpJmpZ` - Jump if zero (false)
- `OpJmpNZ` - Jump if not zero (true)

**Functions:**
- `OpInitFcall` - Initialize function call
- `OpSendVal` - Send argument
- `OpDoFcall` - Execute function call
- `OpReturn` - Return from function

**Arrays:**
- `OpInitArray` - Create array
- `OpAddArrayElement` - Add element to array
- `OpFetchDimR` - Read array element
- `OpAssignDim` - Write array element

**Objects:**
- `OpNew` - Create object
- `OpFetchObjR` - Read property
- `OpAssignObj` - Write property
- `OpInitMethodCall` - Initialize method call

**I/O:**
- `OpEcho` - Output to stdout

**Closures:**
- `OpDeclareFunction` - Declare function
- `OpDeclareLambdaFunction` - Declare closure
- `OpBindLexical` - Bind closure variable

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/compiler"
    "github.com/krizos/php-go/pkg/lexer"
    "github.com/krizos/php-go/pkg/parser"
)

func main() {
    input := `<?php
    $x = 10;
    $y = 20;
    echo $x + $y;
    `

    // Parse
    l := lexer.New(input, "test.php")
    p := parser.New(l)
    program := p.ParseProgram()

    if p.HasErrors() {
        fmt.Println("Parse errors:", p.Errors())
        return
    }

    // Compile
    c := compiler.New()
    err := c.Compile(program)
    if err != nil {
        fmt.Println("Compilation error:", err)
        return
    }

    bytecode := c.Bytecode()
    fmt.Printf("Compiled %d instructions\n", len(bytecode.Instructions))
    fmt.Printf("Constant pool: %d entries\n", len(bytecode.Constants))
}
```

## Optimizations

The compiler implements several optimizations:

### Constant Folding
Binary operations on constants are evaluated at compile time:
- `2 + 3` → `5`
- `"hello" . " world"` → `"hello world"`

### Dead Code Elimination
Unreachable code after return/throw is not compiled.

### Constant Deduplication
Identical constants share the same index in the constant table.

### Last Instruction Tracking
Maintains `lastInstruction` and `previousInstruction` for optimization passes.

## Error Handling

Compilation errors are returned with context:
```go
err := compiler.Compile(node)
if err != nil {
    // Error contains file, line, and description
    fmt.Println(err)
}
```

Common compilation errors:
- Undefined variable access
- Break/continue outside of loop
- Invalid operand types
- Symbol resolution failures

## Implementation Details

- **Instruction Size**: 24 bytes per instruction (fixed size)
- **Operand Types**: Constant, compiled variable (CV), temporary variable (TMPVAR)
- **Stack-Based**: Expression results go into temporary variables
- **Call Frames**: Function calls create new stack frames
- **Closure Support**: Captured variables use `OpBindLexical`
- **Named Arguments**: Tracked via extended value in instructions

## Performance

Compiler achieves 85%+ code coverage with:
- Fast constant folding
- Efficient symbol table lookups
- Minimal memory allocations
- Single-pass compilation (no separate optimization passes yet)

Target: Competitive with PHP 8.4 + opcache
