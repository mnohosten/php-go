# Phase 9: Advanced Features

## Overview

Phase 9 implements advanced PHP features: generators, closures, exceptions, reflection, attributes, and other modern PHP functionality.

## Goals

1. **Generators**: yield and yield from
2. **Closures**: Anonymous functions with variable capture
3. **Arrow Functions**: Short closure syntax
4. **Exceptions**: Complete exception handling
5. **Reflection**: Runtime introspection
6. **Attributes**: PHP 8.0+ attributes/annotations
7. **Weak References**: Weak reference support
8. **Named Arguments**: PHP 8.0+ named parameters

## Components

### 1. Generators

**File**: `pkg/runtime/generator.go`

```go
type Generator struct {
    function *Function
    frame    *Frame
    state    GeneratorState
    key      *Value
    value    *Value
}

type GeneratorState int

const (
    GenStateStart GeneratorState = iota
    GenStateRunning
    GenStateYielded
    GenStateDone
)

func (g *Generator) Current() *Value
func (g *Generator) Key() *Value
func (g *Generator) Next()
func (g *Generator) Valid() bool
func (g *Generator) Rewind()
func (g *Generator) Send(val *Value) *Value
```

**Opcodes**:
- OpYield
- OpYieldFrom
- OpGeneratorReturn

### 2. Closures

**File**: `pkg/runtime/closure.go`

```go
type Closure struct {
    function *Function
    bindings map[string]*Value  // Captured variables
    useRef   map[string]bool    // By-reference capture
}

func NewClosure(fn *Function, bindings map[string]*Value) *Closure
func (c *Closure) Call(args []*Value) (*Value, error)
func (c *Closure) Bind($this *Object, newScope *Class) *Closure
```

### 3. Exceptions

**File**: `pkg/runtime/exception.go`

```go
type Exception struct {
    class    *Class
    message  string
    code     int64
    file     string
    line     int
    trace    *StackTrace
    previous *Exception
}

type StackTrace struct {
    frames []*StackFrame
}

type StackFrame struct {
    function string
    file     string
    line     int
    args     []*Value
}
```

**Opcodes**:
- OpThrow
- OpCatch
- OpFinally

### 4. Reflection

**File**: `pkg/stdlib/reflection/`

Classes to implement:
- ReflectionClass
- ReflectionMethod
- ReflectionProperty
- ReflectionParameter
- ReflectionFunction
- ReflectionType
- ReflectionAttribute

```php
<?php
$class = new ReflectionClass('MyClass');
$methods = $class->getMethods();
$properties = $class->getProperties();
$attrs = $class->getAttributes();
```

### 5. Attributes

**File**: `pkg/runtime/attribute.go`

```go
type Attribute struct {
    name string
    args []*Value
}

type AttributeTarget int

const (
    TargetClass AttributeTarget = 1 << iota
    TargetFunction
    TargetMethod
    TargetProperty
    TargetClassConstant
    TargetParameter
)
```

**PHP Example**:
```php
<?php
#[Route('/api/users', methods: ['GET', 'POST'])]
class UserController {
    #[Authorize(role: 'admin')]
    public function index() { }
}
```

## Implementation Tasks

### Task 9.1: Generator Implementation
**Effort**: 16 hours

- [ ] Generator struct
- [ ] Generator state machine
- [ ] OpYield implementation
- [ ] OpYieldFrom implementation
- [ ] Generator iterator interface
- [ ] send() and throw() methods
- [ ] Generator finalization

### Task 9.2: Closure Implementation
**Effort**: 12 hours

- [ ] Closure struct
- [ ] Variable capture (by value)
- [ ] Variable capture (by reference)
- [ ] Closure call mechanism
- [ ] $this binding
- [ ] Closure rebinding (bindTo)
- [ ] Static closures

### Task 9.3: Arrow Functions
**Effort**: 6 hours

- [ ] Arrow function parsing
- [ ] Implicit variable capture
- [ ] Arrow function compilation
- [ ] Single-expression body

### Task 9.4: Exception System
**Effort**: 14 hours

- [ ] Exception class hierarchy
- [ ] OpThrow implementation
- [ ] OpCatch implementation
- [ ] OpFinally implementation
- [ ] Stack trace generation
- [ ] Exception chaining (previous)
- [ ] Error conversion to exceptions

### Task 9.5: Reflection - Classes
**Effort**: 12 hours

- [ ] ReflectionClass implementation
- [ ] Class metadata access
- [ ] Method enumeration
- [ ] Property enumeration
- [ ] Constructor access
- [ ] Parent/interface access

### Task 9.6: Reflection - Functions & Methods
**Effort**: 10 hours

- [ ] ReflectionFunction
- [ ] ReflectionMethod
- [ ] Parameter reflection
- [ ] Return type reflection
- [ ] Invocation through reflection

### Task 9.7: Reflection - Properties & Parameters
**Effort**: 8 hours

- [ ] ReflectionProperty
- [ ] ReflectionParameter
- [ ] Type information
- [ ] Default values
- [ ] Attributes access

### Task 9.8: Attributes
**Effort**: 12 hours

- [ ] Attribute parsing
- [ ] Attribute compilation
- [ ] Attribute storage
- [ ] Attribute reflection
- [ ] Built-in attributes
- [ ] Attribute validation

### Task 9.9: Weak References
**Effort**: 6 hours

- [ ] WeakReference class
- [ ] WeakMap class
- [ ] Reference tracking
- [ ] GC integration

### Task 9.10: Named Arguments
**Effort**: 8 hours

- [ ] Named argument parsing
- [ ] Named argument calling
- [ ] Argument order independence
- [ ] Mixed positional/named

### Task 9.11: Variadic Functions
**Effort**: 6 hours

- [ ] ... operator for parameters
- [ ] ... operator for arguments (unpacking)
- [ ] func_get_args() compatibility

### Task 9.12: First-Class Callables
**Effort**: 4 hours

- [ ] Callable syntax (PHP 8.1+)
- [ ] strlen(...) creates callable
- [ ] $obj->method(...) creates callable

### Task 9.13: Testing
**Effort**: 16 hours

- [ ] Generator tests
- [ ] Closure tests
- [ ] Exception tests
- [ ] Reflection tests
- [ ] Attribute tests
- [ ] Integration tests

## Estimated Timeline

**Total Effort**: ~130 hours (7 weeks)

## Success Criteria

- [ ] Generators work with yield/yield from
- [ ] Closures capture variables correctly
- [ ] Arrow functions functional
- [ ] Exceptions throw/catch/finally working
- [ ] Reflection provides full introspection
- [ ] Attributes parse and reflect
- [ ] Test coverage >85%

## Progress Tracking

- [ ] Task 9.1: Generators (0%)
- [ ] Task 9.2: Closures (0%)
- [ ] Task 9.3: Arrow Functions (0%)
- [ ] Task 9.4: Exception System (0%)
- [ ] Task 9.5: Reflection - Classes (0%)
- [ ] Task 9.6: Reflection - Functions (0%)
- [ ] Task 9.7: Reflection - Properties (0%)
- [ ] Task 9.8: Attributes (0%)
- [ ] Task 9.9: Weak References (0%)
- [ ] Task 9.10: Named Arguments (0%)
- [ ] Task 9.11: Variadic Functions (0%)
- [ ] Task 9.12: First-Class Callables (0%)
- [ ] Task 9.13: Testing (0%)

**Overall Phase 9 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
