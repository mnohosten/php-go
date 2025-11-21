# Phase 5: Object System

## Overview

Phase 5 implements PHP's object-oriented programming features: classes, objects, inheritance, interfaces, traits, and enums. This is one of the most complex phases.

## Goals

1. **Classes**: Class definitions and instantiation
2. **Objects**: Object creation and property access
3. **Inheritance**: Single inheritance and interfaces
4. **Traits**: Horizontal code reuse
5. **Enums**: Enumeration types (PHP 8.1+)
6. **Magic Methods**: __construct, __get, __set, etc.
7. **Visibility**: public, protected, private
8. **Static Members**: Static properties and methods

## Success Criteria

- [ ] Can define and instantiate classes
- [ ] Inheritance works correctly
- [ ] Interfaces enforced
- [ ] Traits compose correctly
- [ ] Enums functional
- [ ] Magic methods work
- [ ] Visibility rules enforced
- [ ] Static binding correct
- [ ] Test coverage >85%

## Components

### 1. Class Definition

**File**: `pkg/types/class.go`

```go
type Class struct {
    name         string
    parent       *Class
    interfaces   []*Interface
    traits       []*Trait
    constants    map[string]*Value
    properties   map[string]*Property
    methods      map[string]*Method
    constructor  *Method
    destructor   *Method
    flags        ClassFlags
}

type Property struct {
    name       string
    visibility Visibility
    static     bool
    readonly   bool
    typ        *TypeInfo
    default    *Value
}

type Method struct {
    name       string
    visibility Visibility
    static     bool
    abstract   bool
    final      bool
    params     []*Parameter
    returnType *TypeInfo
    body       *Function
}
```

### 2. Object Instance

**File**: `pkg/types/object.go`

```go
type Object struct {
    class      *Class
    properties map[string]*Value
    handlers   *ObjectHandlers
}

type ObjectHandlers struct {
    get        *Method  // __get
    set        *Method  // __set
    call       *Method  // __call
    toString   *Method  // __toString
    // ... other magic methods
}
```

## Implementation Tasks

### Task 5.1: Class Structure
**Effort**: 10 hours

- [ ] Class struct definition
- [ ] Property definitions
- [ ] Method definitions
- [ ] Class registry
- [ ] Constant handling

### Task 5.2: Object Creation
**Effort**: 8 hours

- [ ] Object struct
- [ ] OpNew - Create object
- [ ] OpInitMethodCall - Method call setup
- [ ] OpDoMethodCall - Execute method
- [ ] Constructor invocation

### Task 5.3: Property Access
**Effort**: 10 hours

- [ ] OpFetchObj - Read property
- [ ] OpAssignObj - Write property
- [ ] Visibility checking
- [ ] Static property access
- [ ] Dynamic property names

### Task 5.4: Method Calls
**Effort**: 10 hours

- [ ] Instance method calls
- [ ] Static method calls (::)
- [ ] Method lookup
- [ ] $this binding
- [ ] self/parent/static resolution

### Task 5.5: Inheritance
**Effort**: 14 hours

- [ ] Class extension
- [ ] Method override checking
- [ ] Property inheritance
- [ ] Parent method calls (parent::)
- [ ] Abstract class enforcement
- [ ] Final class/method enforcement

**Reference**: `php-src/Zend/zend_inheritance.c` (142KB - complex!)

### Task 5.6: Interfaces
**Effort**: 8 hours

- [ ] Interface definitions
- [ ] Interface implementation
- [ ] Multiple interfaces
- [ ] Interface compliance checking
- [ ] Interface constants

### Task 5.7: Traits
**Effort**: 12 hours

- [ ] Trait definitions
- [ ] Trait composition
- [ ] Trait method conflicts
- [ ] Trait precedence
- [ ] Trait aliasing
- [ ] Trait properties

### Task 5.8: Enums
**Effort**: 8 hours

- [ ] Enum declarations
- [ ] Backed enums
- [ ] Enum cases
- [ ] Enum methods
- [ ] Pattern matching with enums

### Task 5.9: Magic Methods
**Effort**: 12 hours

- [ ] __construct
- [ ] __destruct (with finalizer)
- [ ] __get / __set
- [ ] __isset / __unset
- [ ] __call / __callStatic
- [ ] __toString
- [ ] __invoke
- [ ] __clone
- [ ] __debugInfo
- [ ] __serialize / __unserialize

### Task 5.10: Type Checking
**Effort**: 8 hours

- [ ] Property type hints
- [ ] Parameter type checking
- [ ] Return type checking
- [ ] Type covariance
- [ ] Type contravariance
- [ ] Readonly properties

### Task 5.11: Late Static Binding
**Effort**: 6 hours

- [ ] static:: resolution
- [ ] get_called_class()
- [ ] Late static binding in inheritance

### Task 5.12: Reflection
**Effort**: 10 hours

- [ ] ReflectionClass
- [ ] ReflectionMethod
- [ ] ReflectionProperty
- [ ] Class metadata access

### Task 5.13: Testing
**Effort**: 14 hours

- [ ] Class definition tests
- [ ] Inheritance tests
- [ ] Interface tests
- [ ] Trait tests
- [ ] Enum tests
- [ ] Magic method tests
- [ ] Visibility tests
- [ ] Complex inheritance hierarchies

## Estimated Timeline

**Total Effort**: ~130 hours (7-8 weeks)

This is a complex phase!

## Progress Tracking

- [ ] Task 5.1: Class Structure (0%)
- [ ] Task 5.2: Object Creation (0%)
- [ ] Task 5.3: Property Access (0%)
- [ ] Task 5.4: Method Calls (0%)
- [ ] Task 5.5: Inheritance (0%)
- [ ] Task 5.6: Interfaces (0%)
- [ ] Task 5.7: Traits (0%)
- [ ] Task 5.8: Enums (0%)
- [ ] Task 5.9: Magic Methods (0%)
- [ ] Task 5.10: Type Checking (0%)
- [ ] Task 5.11: Late Static Binding (0%)
- [ ] Task 5.12: Reflection (0%)
- [ ] Task 5.13: Testing (0%)

**Overall Phase 5 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
