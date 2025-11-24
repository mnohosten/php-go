# Types API Reference

Package: `github.com/krizos/php-go/pkg/types`

## Overview

The types package implements PHP's type system including the universal Value container (equivalent to PHP's zval), arrays, objects, resources, and all type conversion rules. It provides PHP-compatible type juggling, comparison semantics, and memory optimization through value pooling.

## Main Types

### Value

The universal container for all PHP values (equivalent to PHP's zval).

```go
type Value struct {
    // Private fields
}
```

**Type Constructors:**

```go
func NewUndef() *Value
func NewNull() *Value
func NewBool(v bool) *Value
func NewInt(v int64) *Value
func NewFloat(v float64) *Value
func NewString(v string) *Value
func NewArray(v *Array) *Value
func NewObject(v *Object) *Value
func NewResource(v *Resource) *Value
func NewReference(v *Value) *Value
```

**Type Queries:**

```go
func (v *Value) Type() ValueType
func (v *Value) IsUndef() bool
func (v *Value) IsNull() bool
func (v *Value) IsBool() bool
func (v *Value) IsInt() bool
func (v *Value) IsFloat() bool
func (v *Value) IsString() bool
func (v *Value) IsArray() bool
func (v *Value) IsObject() bool
func (v *Value) IsResource() bool
func (v *Value) IsReference() bool
func (v *Value) IsScalar() bool
func (v *Value) IsCallable() bool
```

**Type Conversions (PHP semantics):**

```go
func (v *Value) ToInt() int64
func (v *Value) ToFloat() float64
func (v *Value) ToBool() bool
func (v *Value) ToString() string
func (v *Value) ToArray() *Array
```

**Type Accessors:**

```go
func (v *Value) AsInt() int64
func (v *Value) AsFloat() float64
func (v *Value) AsBool() bool
func (v *Value) AsString() string
func (v *Value) AsArray() *Array
func (v *Value) AsObject() *Object
func (v *Value) AsResource() *Resource
```

**Comparisons:**

```go
func (v *Value) Equals(other *Value) bool        // Loose comparison (==)
func (v *Value) Identical(other *Value) bool     // Strict comparison (===)
func (v *Value) Compare(other *Value) int        // Spaceship operator (<=>)
func (v *Value) LessThan(other *Value) bool
func (v *Value) LessThanOrEqual(other *Value) bool
func (v *Value) GreaterThan(other *Value) bool
func (v *Value) GreaterThanOrEqual(other *Value) bool
```

**Memory Management:**

```go
func (v *Value) Release()
```

Returns a Value to the pool for reuse. Only call on values no longer referenced. Do NOT release cached values (integers -128 to 1023, booleans, null, undef).

**Other Methods:**

```go
func (v *Value) Copy() *Value
func (v *Value) Dereference() *Value
func (v *Value) String() string
```

### ValueType

Enumeration of PHP value types.

```go
type ValueType uint8

const (
    TypeUndef    ValueType = iota
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
```

**Method:**

```go
func (vt ValueType) String() string
```

### Array

PHP's associative array (ordered map).

```go
type Array struct {
    // Private fields
}
```

**Constructors:**

```go
func NewEmptyArray() *Array
func NewArrayFromMap(m map[string]*Value) *Array
func NewArrayFromSlice(slice []*Value) *Array
```

**Basic Operations:**

```go
func (a *Array) Set(key *Value, value *Value)
func (a *Array) Get(key *Value) (*Value, bool)
func (a *Array) Has(key *Value) bool
func (a *Array) Delete(key *Value)
func (a *Array) Len() int
func (a *Array) Clear()
```

**Iteration:**

```go
func (a *Array) Keys() []*Value
func (a *Array) Values() []*Value
func (a *Array) ForEach(fn func(key, value *Value) bool)
```

**Array-specific Methods:**

```go
func (a *Array) Append(value *Value)
func (a *Array) Push(values ...*Value)
func (a *Array) Pop() (*Value, bool)
func (a *Array) Shift() (*Value, bool)
func (a *Array) Unshift(values ...*Value)
func (a *Array) Slice(start, length int) *Array
func (a *Array) Merge(other *Array) *Array
```

**Utility:**

```go
func (a *Array) Copy() *Array
func (a *Array) String() string
func (a *Array) IsEmpty() bool
```

**Implementation Notes:**
- Preserves insertion order (like PHP)
- Keys can be integers or strings
- Numeric string keys are NOT converted to integers (PHP quirk)
- Uses map + slice for order tracking
- O(1) for Get/Set/Has/Delete
- O(n) for Keys/Values/ForEach

### Object

PHP object instance.

```go
type Object struct {
    // Private fields
}
```

**Constructor:**

```go
func NewObject(class *ClassEntry) *Object
```

**Property Access:**

```go
func (o *Object) GetProperty(name string) (*Value, bool)
func (o *Object) SetProperty(name string, value *Value)
func (o *Object) HasProperty(name string) bool
func (o *Object) UnsetProperty(name string)
func (o *Object) GetProperties() map[string]*Value
```

**Method Calls:**

```go
func (o *Object) CallMethod(methodName string, args []*Value) (*Value, error)
func (o *Object) HasMethod(methodName string) bool
```

**Magic Methods:**

```go
func (o *Object) CallMagicGet(propName string) (*Value, error)
func (o *Object) CallMagicSet(propName string, value *Value) error
func (o *Object) CallMagicIsset(propName string) (bool, error)
func (o *Object) CallMagicUnset(propName string) error
func (o *Object) CallMagicCall(methodName string, args []*Value) (*Value, error)
func (o *Object) CallMagicInvoke(args []*Value) (*Value, error)
func (o *Object) CallMagicToString() (string, error)
```

Supports all 14 magic methods: `__construct`, `__destruct`, `__get`, `__set`, `__isset`, `__unset`, `__call`, `__callStatic`, `__invoke`, `__toString`, `__clone`, `__sleep`, `__wakeup`, `__debugInfo`

**Class Information:**

```go
func (o *Object) GetClass() *ClassEntry
func (o *Object) GetClassName() string
func (o *Object) InstanceOf(className string) bool
```

**Cloning:**

```go
func (o *Object) Clone() *Object
```

**Utility:**

```go
func (o *Object) String() string
```

### ClassEntry

Represents a PHP class definition.

```go
type ClassEntry struct {
    Name          string
    ParentClass   *ClassEntry
    Interfaces    []*InterfaceEntry
    Traits        []*TraitEntry
    Properties    map[string]*PropertyEntry
    Methods       map[string]*MethodEntry
    Constants     map[string]*Value
    IsAbstract    bool
    IsFinal       bool
    IsReadonly    bool
    IsEnum        bool
    EnumType      string
    EnumCases     map[string]*EnumCase
}
```

**Constructor:**

```go
func NewClassEntry(name string) *ClassEntry
```

**Methods:**

```go
func (ce *ClassEntry) InheritFrom(parent *ClassEntry) error
func (ce *ClassEntry) AddInterface(iface *InterfaceEntry) error
func (ce *ClassEntry) ApplyTraits(traits []*TraitEntry) error
func (ce *ClassEntry) ValidateInterfaceImplementation() error
func (ce *ClassEntry) GetMethod(name string) (*MethodEntry, bool)
func (ce *ClassEntry) HasMethod(name string) bool
```

### PropertyEntry

Represents a class property definition.

```go
type PropertyEntry struct {
    Name         string
    Visibility   Visibility
    Static       bool
    Readonly     bool
    Type         string
    DefaultValue *Value
    Hooks        *PropertyHooks
}
```

### MethodEntry

Represents a class method definition.

```go
type MethodEntry struct {
    Name        string
    Visibility  Visibility
    Static      bool
    Abstract    bool
    Final       bool
    Function    interface{}
    IsInternal  bool
}
```

### InterfaceEntry

Represents a PHP interface.

```go
type InterfaceEntry struct {
    Name       string
    Extends    []*InterfaceEntry
    Methods    map[string]*MethodSignature
    Constants  map[string]*Value
}
```

### TraitEntry

Represents a PHP trait.

```go
type TraitEntry struct {
    Name       string
    Properties map[string]*PropertyEntry
    Methods    map[string]*MethodEntry
}
```

### EnumCase

Represents an enum case (PHP 8.1+).

```go
type EnumCase struct {
    Name  string
    Value *Value
}
```

### Visibility

Enumeration of property/method visibility.

```go
type Visibility uint8

const (
    VisibilityPublic Visibility = iota
    VisibilityProtected
    VisibilityPrivate
)
```

### Resource

Represents a PHP resource handle.

```go
type Resource struct {
    // Private fields
}
```

**Constructor:**

```go
func NewResource(data interface{}, resourceType string) *Resource
```

**Methods:**

```go
func (r *Resource) ID() int
func (r *Resource) Type() string
func (r *Resource) Data() interface{}
func (r *Resource) IsValid() bool
func (r *Resource) Close()
func (r *Resource) String() string
```

## Type Conversion Rules

### To Boolean (ToBool)

Falsy values:
- `false`
- `0` (integer)
- `0.0` (float)
- `""` (empty string)
- `"0"` (string)
- `null`
- Empty array
- `NaN`

Everything else is truthy.

### To Integer (ToInt)

- Boolean: `true` → 1, `false` → 0
- Float: Truncated to integer
- String: Leading numeric portion parsed
- Array: Empty → 0, non-empty → 1
- Object: Always 1
- Resource: Resource ID

### To Float (ToFloat)

Similar to ToInt but preserves decimal values.

### To String (ToString)

- Boolean: `true` → "1", `false` → ""
- Integer/Float: String representation
- Array: "Array"
- Object: Calls `__toString()` if defined, otherwise "Object"
- Resource: "Resource id #123"

### To Array (ToArray)

- Null/undefined: Empty array
- Scalar: Single-element array
- Array: Return as-is
- Object: Properties become array elements

## Comparison Semantics

### Loose Comparison (==)

Uses type juggling:
- `"123" == 123` → true
- `true == 1` → true
- `null == false` → true

### Strict Comparison (===)

No type juggling:
- `"123" === 123` → false
- `true === 1` → false
- `null === false` → false

### Spaceship Operator (<=>)

Returns:
- `-1` if left < right
- `0` if left == right
- `1` if left > right

## Performance Optimizations

### Value Pooling

The types package uses `sync.Pool` to reuse Value structs:

**Integer Cache**: Pre-allocated values for -128 to 1023
- Eliminates ~60-70% of NewInt() allocations
- Immutable cached values (never released)

**Boolean Singletons**: Pre-allocated true/false
- Eliminates 100% of NewBool() allocations

**Null/Undef Singletons**: Pre-allocated
- Eliminates all NewNull()/NewUndef() allocations

**Value Pooling**: Non-cached values use sync.Pool
- ~45% allocation reduction
- Call `Value.Release()` when done with non-cached values

**Measured Results** (SimpleLoop benchmark, 10K iterations):
- Before: 199,196 allocs/op, 3.86 MB/op
- After: 112,020 allocs/op, 1.82 MB/op
- Improvement: 43.7% fewer allocations, 52.9% less memory

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/types"
)

func main() {
    // Create values
    x := types.NewInt(42)
    y := types.NewString("42")

    // Loose comparison (type juggling)
    fmt.Println(x.Equals(y))      // true

    // Strict comparison
    fmt.Println(x.Identical(y))   // false

    // Type conversions
    fmt.Println(y.ToInt())        // 42
    fmt.Println(x.ToString())     // "42"

    // Arrays
    arr := types.NewEmptyArray()
    arr.Set(types.NewString("name"), types.NewString("Alice"))
    arr.Set(types.NewInt(0), types.NewInt(100))

    if val, ok := arr.Get(types.NewString("name")); ok {
        fmt.Println(val.AsString())  // "Alice"
    }

    fmt.Println(arr.Len())  // 2

    // Objects
    class := types.NewClassEntry("Person")
    obj := types.NewObject(class)
    obj.SetProperty("name", types.NewString("Bob"))

    if name, ok := obj.GetProperty("name"); ok {
        fmt.Println(name.AsString())  // "Bob"
    }
}
```

## String Operations

Additional string utilities in `string.go`:

```go
func StringLength(s string) int
func StringConcat(a, b string) string
func StringRepeat(s string, count int) string
func StringSlice(s string, start, length int) string
```

## Implementation Notes

- **Thread Safety**: Value, Array, and Object are NOT thread-safe
- **Reference Semantics**: Arrays and Objects are reference types
- **Copy-on-Write**: Not yet implemented (planned for Phase 7)
- **String Encoding**: All strings are binary-safe (no encoding assumptions)
- **Array Keys**: Integer and string keys use separate hash maps internally
- **Object Properties**: Dynamic properties supported
- **Weak References**: Implemented in WeakReference type (advanced feature)

## Coverage

The types package has 89%+ code coverage with comprehensive tests for:
- All type conversions
- Comparison operations
- Array operations
- Object operations (78% coverage)
- Resource handling
- Edge cases and PHP quirks
