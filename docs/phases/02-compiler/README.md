# Phase 2: Compiler - AST to Opcodes

## Overview

Phase 2 implements the compiler that transforms the Abstract Syntax Tree (AST) from Phase 1 into bytecode instructions (opcodes) that the virtual machine can execute. This phase bridges the gap between parsing and execution.

## Goals

1. **Opcode Definitions**: Define all 210 PHP opcodes
2. **Compilation**: Transform AST nodes into opcode sequences
3. **Symbol Tables**: Track variables, functions, and classes
4. **Control Flow**: Handle jumps, loops, and conditionals
5. **Optimizations**: Basic compile-time optimizations

## Dependencies

- **Phase 1**: Complete AST from parser
- **Inputs**: AST tree
- **Outputs**: Compiled bytecode (opcode array + metadata)

## Success Criteria

- [ ] All 210 PHP opcodes defined
- [ ] Can compile all AST node types
- [ ] Symbol table correctly tracks scope
- [ ] Control flow (jumps) correctly resolved
- [ ] Functions and classes compiled
- [ ] Basic optimizations (constant folding)
- [ ] Test coverage >85%

## Components

### 1. Opcode Definitions

**File**: `pkg/vm/opcodes.go`

**210 PHP Opcodes** (from `Zend/zend_vm_opcodes.h`):

```go
type Opcode uint8

const (
    OpNOP Opcode = iota
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpMod
    OpSL              // Shift left
    OpSR              // Shift right
    OpConcat          // String concatenation
    OpBWOr            // Bitwise OR
    OpBWAnd           // Bitwise AND
    OpBWXor           // Bitwise XOR
    OpBoolNot         // Boolean NOT
    OpBWNot           // Bitwise NOT
    OpIsEqual         // ==
    OpIsNotEqual      // !=
    OpIsIdentical     // ===
    OpIsNotIdentical  // !==
    OpIsSmaller       // <
    OpIsSmallerOrEqual // <=
    // ... 190 more opcodes
    OpEcho
    OpReturn
    OpAssign
    OpJmp             // Unconditional jump
    OpJmpZ            // Jump if zero (false)
    OpJmpNZ           // Jump if not zero (true)
    OpInitFcall       // Initialize function call
    OpDoFcall         // Execute function call
    OpFetchDim        // Array access $a[x]
    OpAssignDim       // Array assignment $a[x] = y
    OpFetchObj        // Property access $o->x
    OpAssignObj       // Property assignment $o->x = y
    OpNew             // New object
    OpInstanceOf
    OpThrow
    OpCatch
    // ... etc
)
```

**Opcode Categories**:
1. Arithmetic: ADD, SUB, MUL, DIV, MOD, POW
2. Bitwise: BW_OR, BW_AND, BW_XOR, SL, SR
3. Comparison: IS_EQUAL, IS_SMALLER, etc.
4. Logic: BOOL_NOT, BOOL_AND, BOOL_OR, JMPZ, JMPNZ
5. Variables: ASSIGN, FETCH, UNSET, ISSET
6. Arrays: INIT_ARRAY, ADD_ARRAY_ELEMENT, FETCH_DIM
7. Objects: NEW, FETCH_OBJ, INIT_METHOD_CALL
8. Functions: INIT_FCALL, DO_FCALL, RETURN
9. Control: JMP, SWITCH, MATCH
10. Special: ECHO, INCLUDE, EVAL, YIELD

### 2. Instruction Encoding

**File**: `pkg/vm/instruction.go`

```go
// Instruction represents a single VM instruction
type Instruction struct {
    Opcode  Opcode      // The operation
    Op1     Operand     // First operand
    Op2     Operand     // Second operand
    Result  Operand     // Result operand
    Lineno  uint32      // Source line number
}

// Operand represents an instruction operand
type Operand struct {
    Type OperandType
    Num  uint32       // Index or value
}

type OperandType uint8

const (
    OpUnused OperandType = iota
    OpConst     // Constant from literal table
    OpTmpVar    // Temporary variable
    OpVar       // Runtime variable
    OpCV        // Compiled variable (optimized local)
    OpJmpAddr   // Jump address
)
```

### 3. Compiler Core

**File**: `pkg/compiler/compiler.go`

**Main Structure**:
```go
type Compiler struct {
    ast         *ast.File
    opcodes     []Instruction
    constants   []*types.Value    // Literal table
    functions   map[string]*CompiledFunction
    classes     map[string]*CompiledClass
    scope       *Scope            // Current scope
    breakStack  []int             // Break jump targets
    continueStack []int           // Continue jump targets
    loopDepth   int
}

func (c *Compiler) Compile(file *ast.File) (*Program, error) {
    // Walk AST and generate opcodes
    c.compileFile(file)
    return &Program{
        Opcodes:   c.opcodes,
        Constants: c.constants,
        Functions: c.functions,
        Classes:   c.classes,
    }, nil
}
```

### 4. Symbol Table

**File**: `pkg/compiler/symbols.go`

```go
// Scope tracks variables in current scope
type Scope struct {
    parent    *Scope
    vars      map[string]*Variable
    depth     int
}

type Variable struct {
    Name      string
    Index     int       // Variable index in local array
    Type      VarType   // CV, VAR, TMP
    IsGlobal  bool
    IsStatic  bool
}

func (s *Scope) Declare(name string) *Variable
func (s *Scope) Lookup(name string) (*Variable, bool)
func (s *Scope) EnterScope()  // Create child scope
func (s *Scope) ExitScope()   // Return to parent
```

## Implementation Tasks

### Task 2.1: Opcode Definitions
**File**: `pkg/vm/opcodes.go`
**Effort**: 6 hours

- [ ] Define all 210 opcode constants
- [ ] Group opcodes by category
- [ ] Add String() method for debugging
- [ ] Document each opcode's purpose
- [ ] Add opcode metadata (operand counts, side effects)

**Reference**: `php-src/Zend/zend_vm_opcodes.h`

### Task 2.2: Instruction Encoding
**File**: `pkg/vm/instruction.go`
**Effort**: 4 hours

- [ ] Define Instruction struct
- [ ] Define Operand types
- [ ] Implement instruction encoding/decoding
- [ ] Add instruction String() for debugging
- [ ] Implement instruction builder helpers

### Task 2.3: Compiler Core
**File**: `pkg/compiler/compiler.go`
**Effort**: 8 hours

- [ ] Create Compiler struct
- [ ] Implement AST visitor pattern
- [ ] Opcode emission methods
- [ ] Constant table management
- [ ] Program assembly
- [ ] Error reporting

### Task 2.4: Symbol Tables
**File**: `pkg/compiler/symbols.go`
**Effort**: 6 hours

- [ ] Implement Scope struct
- [ ] Variable declaration and lookup
- [ ] Scope enter/exit
- [ ] Global vs local variables
- [ ] Static variable handling
- [ ] Variable index assignment

### Task 2.5: Expression Compilation
**File**: `pkg/compiler/expr.go`
**Effort**: 16 hours

- [ ] Compile binary expressions (+, -, *, /, etc.)
- [ ] Compile unary expressions (!, -, ~, etc.)
- [ ] Compile assignment expressions
- [ ] Compile variable access
- [ ] Compile literals (int, float, string, array)
- [ ] Compile function calls
- [ ] Compile method calls
- [ ] Compile array access
- [ ] Compile property access
- [ ] Compile ternary operator
- [ ] Compile null coalesce
- [ ] Compile instanceof
- [ ] Compile type casts

**Critical Component!**

### Task 2.6: Statement Compilation
**File**: `pkg/compiler/stmt.go`
**Effort**: 12 hours

- [ ] Compile echo statement
- [ ] Compile expression statements
- [ ] Compile if/elseif/else
- [ ] Compile while loop
- [ ] Compile do-while loop
- [ ] Compile for loop
- [ ] Compile foreach loop
- [ ] Compile switch statement
- [ ] Compile match expression
- [ ] Compile break/continue
- [ ] Compile return statement
- [ ] Compile throw statement
- [ ] Compile try-catch-finally

### Task 2.7: Control Flow & Jumps
**File**: `pkg/compiler/jumps.go`
**Effort**: 10 hours

- [ ] Implement jump placeholders
- [ ] Patch jump addresses after compilation
- [ ] Track break/continue targets
- [ ] Handle nested loops
- [ ] Switch statement jump tables
- [ ] Try-catch exception tables
- [ ] Verify all jumps resolved

**Important for correctness!**

### Task 2.8: Function Compilation
**File**: `pkg/compiler/function.go`
**Effort**: 10 hours

- [ ] Compile function declarations
- [ ] Compile function parameters
- [ ] Handle default parameters
- [ ] Handle variadic parameters
- [ ] Handle by-reference parameters
- [ ] Compile function body
- [ ] Return type checking
- [ ] Closure compilation
- [ ] Arrow function compilation

### Task 2.9: Class Compilation
**File**: `pkg/compiler/class.go`
**Effort**: 12 hours

- [ ] Compile class declarations
- [ ] Compile properties
- [ ] Compile methods
- [ ] Compile constructors
- [ ] Compile static members
- [ ] Compile constants
- [ ] Handle inheritance
- [ ] Handle interfaces
- [ ] Handle traits
- [ ] Compile enums

### Task 2.10: Optimizations
**File**: `pkg/compiler/optimizer.go`
**Effort**: 8 hours

- [ ] Constant folding (1 + 2 → 3)
- [ ] Dead code elimination
- [ ] Unreachable code detection
- [ ] Strength reduction
- [ ] Common subexpression elimination (basic)
- [ ] Peephole optimizations

**Optional but valuable**

### Task 2.11: Testing
**Effort**: 12 hours

- [ ] Unit tests for each opcode
- [ ] Expression compilation tests
- [ ] Statement compilation tests
- [ ] Control flow tests
- [ ] Function compilation tests
- [ ] Class compilation tests
- [ ] Optimization tests
- [ ] Integration tests

**Target**: 85%+ coverage

## Test Plan

### Unit Tests

```go
func TestCompileExpression(t *testing.T) {
    tests := []struct {
        source   string
        expected []Instruction
    }{
        {
            "$x = 1 + 2",
            []Instruction{
                {OpAdd, Const(1), Const(2), Tmp(0)},
                {OpAssign, Tmp(0), _, Var("x")},
            },
        },
        // ... more tests
    }
}

func TestCompileControlFlow(t *testing.T) {
    source := "if ($x) { echo 'yes'; } else { echo 'no'; }"
    compiler := NewCompiler()
    program := compiler.Compile(parse(source))
    // Verify jump instructions
}
```

### Integration Tests

```go
func TestCompileRealPHP(t *testing.T) {
    files := []string{
        "testdata/simple.php",
        "testdata/functions.php",
        "testdata/classes.php",
        "testdata/loops.php",
    }
    for _, file := range files {
        source := readFile(file)
        ast := parse(source)
        program, err := NewCompiler().Compile(ast)
        assert.NoError(t, err)
        assert.NotNil(t, program)
    }
}
```

## Milestones

### Milestone 2.1: Foundation (Week 1)
- Opcodes defined
- Instruction encoding working
- Compiler structure in place

### Milestone 2.2: Basic Compilation (Week 2)
- Can compile simple expressions
- Can compile basic statements
- Symbol table working

### Milestone 2.3: Control Flow (Week 3)
- If/else working
- Loops working
- Break/continue working
- Jump resolution correct

### Milestone 2.4: Functions (Week 4)
- Function declaration
- Function calls
- Parameters and returns

### Milestone 2.5: Classes (Week 5)
- Basic class support
- Methods and properties
- Inheritance basics

### Milestone 2.6: Complete Phase 2 (Week 6)
- All language constructs compile
- Optimizations working
- Full test coverage
- Ready for VM execution

## Estimated Timeline

**Total Effort**: ~110 hours (5-6 weeks)

- Opcodes: ~10 hours
- Compiler core: ~14 hours
- Expression compilation: ~16 hours
- Statement compilation: ~12 hours
- Control flow: ~10 hours
- Functions: ~10 hours
- Classes: ~12 hours
- Optimizations: ~8 hours
- Testing: ~12 hours
- Documentation: ~6 hours

## Dependencies for Next Phase

Phase 3 (Runtime & VM) requires:
- Complete opcode definitions ✓
- Working compiler producing valid bytecode ✓
- Test suite for compilation verification ✓

## References

**PHP Source**:
- `php-src/Zend/zend_compile.c` - Compiler implementation
- `php-src/Zend/zend_compile.h` - Compilation structures
- `php-src/Zend/zend_vm_def.h` - VM opcode definitions

## Progress Tracking

- [ ] Task 2.1: Opcode Definitions (0%)
- [ ] Task 2.2: Instruction Encoding (0%)
- [ ] Task 2.3: Compiler Core (0%)
- [ ] Task 2.4: Symbol Tables (0%)
- [ ] Task 2.5: Expression Compilation (0%)
- [ ] Task 2.6: Statement Compilation (0%)
- [ ] Task 2.7: Control Flow & Jumps (0%)
- [ ] Task 2.8: Function Compilation (0%)
- [ ] Task 2.9: Class Compilation (0%)
- [ ] Task 2.10: Optimizations (0%)
- [ ] Task 2.11: Testing (0%)

**Overall Phase 2 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
