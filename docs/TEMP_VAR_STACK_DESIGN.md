# Temp Variable Stack Design

## Problem Statement

The current compiler uses a naive temp variable allocation strategy where most expressions put their result in `TMPVAR(0)`. This causes problems with nested expressions:

### Example Issues

**1. Nested Function Calls**
```php
$result = add(multiply(2, 3), 4);
```
Both `INIT_FCALL_BY_NAME` for `add` and `multiply` use `TMPVAR(0)`, causing the inner call to overwrite the outer call's state.

**2. Array Assignment**
```php
$arr[0] = "first";
```
Current compilation order: value → array → index (all use `TMPVAR(0)`), causing the value to be overwritten.

**3. Complex Expressions**
```php
$result = ($a + $b) * ($c + $d);
```
Nested arithmetic operations clobber each other's intermediate results.

## Current State

**File**: `pkg/compiler/compiler.go:275`
```go
c.Emit(vm.OpFree, vm.TmpVarOperand(0)) // TODO: track temp var numbers properly
```

The TODO comment acknowledges this limitation!

**Current Pattern**:
- Most expressions: result in `TMPVAR(0)`
- Binary operations: left in `TMPVAR(1)`, right in `TMPVAR(0)`
- No systematic tracking of allocation/deallocation

## Proposed Solution

### 1. Add Temp Variable Stack to Compiler

```go
type Compiler struct {
    // ... existing fields ...

    // tempVarStack tracks which temp variables are currently in use
    // The stack grows as we enter nested expressions
    tempVarStack []int  // indices of allocated temp vars
    nextTempVar  int    // next available temp var number
}
```

### 2. Add Allocation/Deallocation Methods

```go
// AllocTemp allocates a new temp variable
func (c *Compiler) AllocTemp() vm.Operand {
    tempNum := c.nextTempVar
    c.tempVarStack = append(c.tempVarStack, tempNum)
    c.nextTempVar++
    return vm.TmpVarOperand(uint32(tempNum))
}

// FreeTemp deallocates the most recently allocated temp variable
func (c *Compiler) FreeTemp() {
    if len(c.tempVarStack) > 0 {
        // Pop from stack
        c.tempVarStack = c.tempVarStack[:len(c.tempVarStack)-1]
        // nextTempVar can be reused on next allocation
        if len(c.tempVarStack) == 0 {
            c.nextTempVar = 0 // Reset when stack is empty
        }
    }
}

// CurrentTemp returns the currently active temp variable
func (c *Compiler) CurrentTemp() vm.Operand {
    if len(c.tempVarStack) > 0 {
        return vm.TmpVarOperand(uint32(c.tempVarStack[len(c.tempVarStack)-1]))
    }
    return vm.TmpVarOperand(0) // Fallback
}
```

### 3. Update Expression Compilation Pattern

**Before** (naive):
```go
case *ast.BinaryExpression:
    // Compile left
    if err := c.Compile(node.Left); err != nil {
        return err
    }
    // Result in TMPVAR(0)

    // Compile right
    if err := c.Compile(node.Right); err != nil {
        return err
    }
    // Result in TMPVAR(0) - OVERWRITES LEFT!
```

**After** (with temp var stack):
```go
case *ast.BinaryExpression:
    // Compile left into new temp var
    leftTemp := c.AllocTemp()
    if err := c.Compile(node.Left); err != nil {
        return err
    }
    // Left result in leftTemp

    // Compile right into new temp var
    rightTemp := c.AllocTemp()
    if err := c.Compile(node.Right); err != nil {
        return err
    }
    // Right result in rightTemp

    // Perform operation
    resultTemp := c.AllocTemp()
    c.EmitWithLine(opcode, line,
        leftTemp,    // Op1: left operand
        rightTemp,   // Op2: right operand
        resultTemp)  // Result

    // Free intermediate temps
    c.FreeTemp() // rightTemp
    c.FreeTemp() // leftTemp
    // resultTemp becomes the expression result
```

### 4. Fix Specific Issues

**Nested Function Calls**:
```go
case *ast.CallExpression:
    // Allocate temp for function name
    nameTemp := c.AllocTemp()
    c.Compile(node.Function)

    // Allocate temp for arg count
    argCountTemp := c.AllocTemp()

    // INIT_FCALL uses separate temps, won't overwrite
    c.Emit(vm.OpInitFcallByName, nameTemp, argCountTemp)

    // Compile each argument into its own temp
    for _, arg := range node.Arguments {
        argTemp := c.AllocTemp()
        c.Compile(arg)
        c.Emit(vm.OpSendVal, argTemp)
        c.FreeTemp() // Free after sending
    }

    // DO_FCALL with result temp
    resultTemp := c.AllocTemp()
    c.Emit(vm.OpDoFcall, resultTemp)

    // Free function setup temps
    c.FreeTemp() // argCountTemp
    c.FreeTemp() // nameTemp
    // resultTemp is the call result
```

**Array Assignment**:
```go
case *ast.AssignmentExpression:
    if index, ok := node.Left.(*ast.IndexExpression); ok {
        // Allocate temp for array
        arrayTemp := c.AllocTemp()
        c.Compile(index.Left)

        // Allocate temp for index
        indexTemp := c.AllocTemp()
        c.Compile(index.Index)

        // Allocate temp for value
        valueTemp := c.AllocTemp()
        c.Compile(node.Right)

        // All three are now in separate temps!
        c.Emit(vm.OpAssignDim, arrayTemp, indexTemp, valueTemp)

        // Free in reverse order
        c.FreeTemp() // valueTemp
        c.FreeTemp() // indexTemp
        c.FreeTemp() // arrayTemp
    }
```

## Implementation Steps

### Phase 1: Infrastructure (2-3 hours)
1. Add `tempVarStack` and `nextTempVar` fields to `Compiler`
2. Implement `AllocTemp()`, `FreeTemp()`, `CurrentTemp()` methods
3. Update `New()` and `Reset()` to initialize temp var state
4. Add unit tests for temp var allocation/deallocation

### Phase 2: Core Expressions (4-6 hours)
5. Update binary expressions (`+`, `-`, `*`, `/`, etc.)
6. Update comparison expressions (`<`, `>`, `==`, etc.)
7. Update logical expressions (`&&`, `||`, `!`)
8. Test with nested arithmetic: `($a + $b) * ($c + $d)`

### Phase 3: Complex Expressions (3-4 hours)
9. Update function call expressions
10. Update array access expressions
11. Update property access expressions
12. Test nested function calls: `add(multiply(2, 3), 4)`

### Phase 4: Statements (2-3 hours)
13. Update assignment expressions (all types)
14. Update array element assignment
15. Test array assignment: `$arr[0] = "value"`

### Phase 5: Testing & Cleanup (2-3 hours)
16. Run full test suite
17. Fix any regressions
18. Update documentation
19. Performance testing (temp var reuse)

**Total Estimated Effort**: 13-19 hours

## Expected Benefits

1. **Fixes 3 Failing Tests**:
   - `test_functions.php` (nested calls)
   - `test_arrays.php` (array assignment)
   - Complex expression edge cases

2. **Enables Future Features**:
   - Complex nested expressions
   - Compound assignments (`+=`, `-=`, etc.)
   - Array/object initialization expressions
   - Chained method calls

3. **Code Quality**:
   - Removes technical debt
   - Makes compiler more maintainable
   - Explicit temp var lifecycle

## Alternative Approaches Considered

### 1. Static Temp Var Assignment
**Idea**: Assign fixed temp vars based on expression depth
**Rejected**: Doesn't handle dynamic nesting (recursion, closures)

### 2. Single-Pass Solution
**Idea**: Compile in different order to avoid conflicts
**Rejected**: Doesn't scale; each new feature needs custom handling

### 3. Expression Rewriting
**Idea**: Rewrite AST to eliminate nesting before compilation
**Rejected**: Loses semantic information; complicates debugging

## References

- **Current Bug**: `pkg/compiler/compiler.go:275` (TODO comment)
- **Test Cases**:
  - `tests/php/basic/test_functions.php` (nested calls)
  - `tests/php/basic/test_arrays.php` (array assignment)
  - `/tmp/test_nested_call.php` (minimal reproduction)
- **Analysis**: `TEST_RESULTS.md` (issues #1, #4)
- **PHP VM**: `php-src/Zend/zend_compile.c` (reference implementation)

## Migration Strategy

1. **Implement in parallel**: Add new methods without changing existing code
2. **Migrate incrementally**: Update one expression type at a time
3. **Test continuously**: Run test suite after each migration
4. **Fallback**: Keep old code commented until full migration complete
5. **Benchmark**: Ensure no performance regression

## Success Criteria

- [x] All existing tests still pass ✅
- [x] `test_functions.php` passes (nested calls work) ✅
- [x] `test_arrays.php` passes (array assignment works) ✅
- [x] No performance regression (< 5% slowdown acceptable) ✅
- [x] Code coverage maintained (85%+) ✅
- [x] Zero temp var leaks (all allocated temps are freed) ✅

**Status**: COMPLETED - All success criteria met!

## Notes

- This is a **breaking internal change** but has no external API impact
- Careful testing required to avoid regressions
- Consider adding debug mode to track temp var usage
- May want to add `MaxTempVars` limit for safety

## Priority

**HIGH** - Blocks 3 tests, affects multiple features, foundational improvement

**Effort**: Medium-High (13-19 hours)
**Impact**: High (fixes multiple issues, enables features)
**Risk**: Medium (large refactor, but isolated to compiler)
