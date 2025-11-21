# Phase 3: Runtime & Virtual Machine

## Overview

Phase 3 implements the Virtual Machine (VM) that executes the bytecode generated in Phase 2. This includes the execution engine, stack management, function calls, and basic runtime support.

## Goals

1. **VM Executor**: Interpret and execute opcodes
2. **Execution Context**: Call frames, stacks, and state management
3. **Type System**: Implement PHP value types (zval equivalent)
4. **Function Calls**: Handle function invocation and returns
5. **Error Handling**: Runtime errors and exceptions
6. **Output Buffering**: Capture PHP output

## Dependencies

- **Phase 2**: Compiled bytecode
- **Inputs**: Bytecode program
- **Outputs**: Execution results, output

## Success Criteria

- [ ] VM can execute all 210 opcodes
- [ ] Type system handles all PHP types
- [ ] Function calls work correctly
- [ ] Stack management is correct
- [ ] Can run simple PHP scripts end-to-end
- [ ] Error handling functional
- [ ] Test coverage >85%

## Components

### 1. Value System (Zval)

**File**: `pkg/types/value.go`

```go
// Value is PHP's universal value container (zval equivalent)
type Value struct {
    typ   ValueType
    flags ValueFlags
    data  interface{}
}

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
    TypeReference
)

// Core methods
func (v *Value) Type() ValueType
func (v *Value) ToInt() int64
func (v *Value) ToFloat() float64
func (v *Value) ToBool() bool
func (v *Value) ToString() string
func (v *Value) Copy() *Value
func (v *Value) IsTrue() bool
```

### 2. Virtual Machine

**File**: `pkg/vm/vm.go`

```go
type VM struct {
    program   *Program
    frames    []*Frame
    globals   map[string]*Value
    functions map[string]*Function
    classes   map[string]*Class
    constants map[string]*Value
    output    *OutputBuffer
}

func (vm *VM) Execute(fn *Function, args []*Value) (*Value, error) {
    frame := vm.newFrame(fn, args)
    return vm.executeFrame(frame)
}

func (vm *VM) executeFrame(frame *Frame) (*Value, error) {
    for frame.ip < len(frame.function.Instructions) {
        instr := frame.function.Instructions[frame.ip]
        if err := vm.dispatch(frame, instr); err != nil {
            return nil, err
        }
        frame.ip++
    }
    return frame.returnValue, nil
}
```

### 3. Execution Frame

**File**: `pkg/vm/frame.go`

```go
type Frame struct {
    function    *Function
    locals      []*Value      // Local variables
    stack       []*Value      // Operand stack
    ip          int           // Instruction pointer
    returnValue *Value
    returnTo    int
    this        *Object       // $this object
}

func (f *Frame) push(v *Value)
func (f *Frame) pop() *Value
func (f *Frame) getLocal(index int) *Value
func (f *Frame) setLocal(index int, v *Value)
```

## Implementation Tasks

### Task 3.1: Value Type System
**File**: `pkg/types/value.go`
**Effort**: 12 hours

- [ ] Define Value struct
- [ ] Implement type constructors (NewInt, NewString, etc.)
- [ ] Implement type conversions (ToInt, ToString, etc.)
- [ ] Implement IsTrue() for truthiness
- [ ] Implement Copy() for value copying
- [ ] Handle reference semantics
- [ ] Add debugging String() method

### Task 3.2: Type Conversions & Juggling
**File**: `pkg/types/conversions.go`
**Effort**: 10 hours

- [ ] Int to other types
- [ ] Float to other types
- [ ] String to numeric
- [ ] Array to scalar
- [ ] Object to scalar
- [ ] Comparison rules (==, ===)
- [ ] Type coercion for operators

**Critical for PHP compatibility!**

### Task 3.3: VM Core Structure
**File**: `pkg/vm/vm.go`
**Effort**: 8 hours

- [ ] Create VM struct
- [ ] Initialize VM state
- [ ] Load program
- [ ] Create global scope
- [ ] Register built-in functions
- [ ] Implement Execute() entry point

### Task 3.4: Execution Frame
**File**: `pkg/vm/frame.go`
**Effort**: 6 hours

- [ ] Define Frame struct
- [ ] Stack operations (push/pop)
- [ ] Local variable access
- [ ] Frame creation and destruction
- [ ] Stack overflow protection

### Task 3.5: Opcode Handlers - Arithmetic
**File**: `pkg/vm/handlers_arithmetic.go`
**Effort**: 8 hours

- [ ] OpAdd - Addition
- [ ] OpSub - Subtraction
- [ ] OpMul - Multiplication
- [ ] OpDiv - Division
- [ ] OpMod - Modulo
- [ ] OpPow - Power
- [ ] OpNegate - Unary minus
- [ ] Handle type juggling for each

### Task 3.6: Opcode Handlers - Comparison
**File**: `pkg/vm/handlers_comparison.go`
**Effort**: 8 hours

- [ ] OpIsEqual (==)
- [ ] OpIsIdentical (===)
- [ ] OpIsNotEqual (!=)
- [ ] OpIsNotIdentical (!==)
- [ ] OpIsSmaller (<)
- [ ] OpIsGreater (>)
- [ ] OpIsSmallerOrEqual (<=)
- [ ] OpIsGreaterOrEqual (>=)
- [ ] OpSpaceship (<=>)

### Task 3.7: Opcode Handlers - Logic & Bitwise
**File**: `pkg/vm/handlers_logic.go`
**Effort**: 6 hours

- [ ] OpBoolNot (!)
- [ ] OpBWNot (~)
- [ ] OpBWAnd (&)
- [ ] OpBWOr (|)
- [ ] OpBWXor (^)
- [ ] OpShiftLeft (<<)
- [ ] OpShiftRight (>>)

### Task 3.8: Opcode Handlers - Variables
**File**: `pkg/vm/handlers_variables.go`
**Effort**: 8 hours

- [ ] OpAssign - Variable assignment
- [ ] OpFetch - Variable fetch
- [ ] OpUnset - Unset variable
- [ ] OpIsset - Check if set
- [ ] OpEmpty - Check if empty
- [ ] Handle global variables
- [ ] Handle static variables
- [ ] Handle superglobals ($\_GET, $\_POST, etc.)

### Task 3.9: Opcode Handlers - Control Flow
**File**: `pkg/vm/handlers_control.go`
**Effort**: 6 hours

- [ ] OpJmp - Unconditional jump
- [ ] OpJmpZ - Jump if false
- [ ] OpJmpNZ - Jump if true
- [ ] OpSwitch - Switch statement
- [ ] OpMatch - Match expression
- [ ] Verify jump targets

### Task 3.10: Opcode Handlers - Functions
**File**: `pkg/vm/handlers_functions.go`
**Effort**: 10 hours

- [ ] OpInitFcall - Initialize function call
- [ ] OpSendVal - Pass argument by value
- [ ] OpSendVar - Pass argument by variable
- [ ] OpSendRef - Pass argument by reference
- [ ] OpDoFcall - Execute function
- [ ] OpReturn - Return from function
- [ ] Handle recursion
- [ ] Handle call stack

### Task 3.11: Opcode Handlers - Strings
**File**: `pkg/vm/handlers_strings.go`
**Effort**: 4 hours

- [ ] OpConcat - String concatenation
- [ ] OpFastConcat - Optimized concat
- [ ] String interpolation handling

### Task 3.12: Opcode Handlers - I/O
**File**: `pkg/vm/handlers_io.go`
**Effort**: 4 hours

- [ ] OpEcho - Output string
- [ ] OpPrint - Print string
- [ ] Output buffering integration

### Task 3.13: Runtime Support
**File**: `pkg/runtime/runtime.go`
**Effort**: 8 hours

- [ ] Global variable management
- [ ] Superglobals ($_GET, $_POST, $_SERVER, etc.)
- [ ] Constant management
- [ ] Include/require handling
- [ ] Error reporting levels

### Task 3.14: Output Buffering
**File**: `pkg/runtime/output.go`
**Effort**: 6 hours

- [ ] OutputBuffer struct
- [ ] ob_start() / ob_end_clean()
- [ ] ob_get_contents()
- [ ] Buffer nesting
- [ ] Flush mechanisms

### Task 3.15: Error Handling
**File**: `pkg/runtime/errors.go`
**Effort**: 8 hours

- [ ] Error types (E_ERROR, E_WARNING, etc.)
- [ ] Error handler registration
- [ ] Error reporting
- [ ] @ error suppression
- [ ] Stack trace generation

### Task 3.16: Testing
**Effort**: 12 hours

- [ ] Unit tests for each value type
- [ ] Type conversion tests
- [ ] Opcode handler tests
- [ ] Stack management tests
- [ ] Function call tests
- [ ] Error handling tests
- [ ] Integration tests (end-to-end)

## Milestones

### Milestone 3.1: Type System (Week 1)
- Value types implemented
- Type conversions working
- Type juggling correct

### Milestone 3.2: VM Foundation (Week 2)
- VM structure complete
- Frame management working
- Basic opcode dispatch

### Milestone 3.3: Arithmetic & Logic (Week 3)
- Arithmetic opcodes working
- Comparison opcodes working
- Bitwise opcodes working

### Milestone 3.4: Variables & Control (Week 4)
- Variable operations working
- Control flow opcodes working
- Jump resolution correct

### Milestone 3.5: Functions (Week 5)
- Function calls working
- Parameter passing correct
- Return values working

### Milestone 3.6: Complete Phase 3 (Week 6)
- All opcodes implemented
- Can execute simple PHP scripts
- Error handling working
- Output buffering functional

## Estimated Timeline

**Total Effort**: ~120 hours (6 weeks)

- Type system: ~22 hours
- VM core: ~14 hours
- Opcode handlers: ~54 hours
- Runtime support: ~14 hours
- Output & errors: ~14 hours
- Testing: ~12 hours

## Dependencies for Next Phase

Phase 4 (Data Structures) requires:
- Working VM executor ✓
- Type system foundation ✓
- Basic opcodes working ✓

## Progress Tracking

- [ ] Task 3.1: Value Type System (0%)
- [ ] Task 3.2: Type Conversions (0%)
- [ ] Task 3.3: VM Core (0%)
- [ ] Task 3.4: Execution Frame (0%)
- [ ] Task 3.5: Arithmetic Handlers (0%)
- [ ] Task 3.6: Comparison Handlers (0%)
- [ ] Task 3.7: Logic Handlers (0%)
- [ ] Task 3.8: Variable Handlers (0%)
- [ ] Task 3.9: Control Flow Handlers (0%)
- [ ] Task 3.10: Function Handlers (0%)
- [ ] Task 3.11: String Handlers (0%)
- [ ] Task 3.12: I/O Handlers (0%)
- [ ] Task 3.13: Runtime Support (0%)
- [ ] Task 3.14: Output Buffering (0%)
- [ ] Task 3.15: Error Handling (0%)
- [ ] Task 3.16: Testing (0%)

**Overall Phase 3 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
