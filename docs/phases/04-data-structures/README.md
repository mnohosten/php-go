# Phase 4: Core Data Structures

## Overview

Phase 4 implements PHP's core data structures: strings, arrays (hash tables), and resources. These are fundamental to PHP's operation and must match PHP's exact semantics.

## Goals

1. **PHP Strings**: Binary-safe strings with PHP semantics
2. **PHP Arrays**: Ordered associative arrays (hash tables)
3. **Resources**: External resource handling
4. **Array Opcodes**: Array access and manipulation opcodes
5. **String Operations**: String manipulation opcodes

## Success Criteria

- [ ] Strings are binary-safe
- [ ] Arrays maintain insertion order
- [ ] Arrays support both integer and string keys
- [ ] Packed array optimization for sequential keys
- [ ] All array opcodes working
- [ ] All string opcodes working
- [ ] Test coverage >90%

## Components

### 1. PHP Strings

**File**: `pkg/types/string.go`

PHP strings are:
- Binary-safe (can contain null bytes)
- Mutable in some contexts
- Byte-oriented (not necessarily UTF-8)
- Length-tracked

```go
type String struct {
    data []byte
    len  int
    hash uint32  // Cached hash
}

func NewString(s string) *String
func (s *String) Len() int
func (s *String) Bytes() []byte
func (s *String) String() string
func (s *String) Concat(other *String) *String
func (s *String) Substr(start, length int) *String
func (s *String) IndexOf(needle *String) int
func (s *String) Compare(other *String) int
```

### 2. PHP Arrays

**File**: `pkg/types/array.go`

PHP arrays are:
- Ordered associative arrays
- Maintain insertion order
- Support integer and string keys
- Auto-incrementing integer keys
- Mixed key types in same array

```go
type Array struct {
    elements map[ArrayKey]*Value
    order    []ArrayKey
    nextKey  int64
}

type ArrayKey struct {
    typ KeyType  // Integer or String
    i   int64
    s   string
}

func NewArray() *Array
func (a *Array) Get(key ArrayKey) (*Value, bool)
func (a *Array) Set(key ArrayKey, value *Value)
func (a *Array) Append(value *Value)  // $arr[] = value
func (a *Array) Delete(key ArrayKey)
func (a *Array) Count() int
func (a *Array) Keys() []ArrayKey
func (a *Array) Values() []*Value
func (a *Array) Exists(key ArrayKey) bool
```

**Packed Array Optimization**:
```go
// For arrays with sequential integer keys starting at 0
type PackedArray struct {
    elements []*Value
}
```

### 3. Resources

**File**: `pkg/types/resource.go`

Resources represent external resources (files, database connections, etc.)

```go
type Resource struct {
    id   int
    typ  string
    data interface{}
}
```

## Implementation Tasks

### Task 4.1: String Implementation
**File**: `pkg/types/string.go`
**Effort**: 10 hours

- [ ] String struct
- [ ] String creation
- [ ] String concatenation
- [ ] Substring operations
- [ ] String comparison
- [ ] String hashing
- [ ] Binary-safe operations
- [ ] Conversion to/from Go strings

### Task 4.2: Array Implementation
**File**: `pkg/types/array.go`
**Effort**: 16 hours

- [ ] Array struct with map + order slice
- [ ] ArrayKey type (integer or string)
- [ ] Get/Set operations
- [ ] Append operation ($arr[] = val)
- [ ] Delete operation (unset)
- [ ] Count operation
- [ ] Key/Value iteration
- [ ] Array copying (copy-on-write semantics)
- [ ] Exists check (isset)

**Critical Component!**

### Task 4.3: Packed Array Optimization
**File**: `pkg/types/packed_array.go`
**Effort**: 8 hours

- [ ] Detect sequential integer keys
- [ ] PackedArray implementation
- [ ] Automatic conversion to/from regular array
- [ ] Performance optimization for common case

### Task 4.4: Resource Implementation
**File**: `pkg/types/resource.go`
**Effort**: 4 hours

- [ ] Resource struct
- [ ] Resource registry
- [ ] Resource creation
- [ ] Resource cleanup
- [ ] Resource type checking

### Task 4.5: Array Opcodes
**File**: `pkg/vm/handlers_array.go`
**Effort**: 12 hours

- [ ] OpInitArray - Initialize empty array
- [ ] OpAddArrayElement - Add element to array
- [ ] OpFetchDim - $arr[$key] read
- [ ] OpAssignDim - $arr[$key] = val write
- [ ] OpUnsetDim - unset($arr[$key])
- [ ] OpIssetDim - isset($arr[$key])
- [ ] OpEmptyDim - empty($arr[$key])
- [ ] OpFetchDimR/W/IS/UNSET variants
- [ ] Handle nested array access

### Task 4.6: String Opcodes
**File**: `pkg/vm/handlers_string.go`
**Effort**: 6 hours

- [ ] OpConcat - String concatenation
- [ ] OpFastConcat - Optimized concatenation
- [ ] OpRopeInit/Add/End - Rope optimization
- [ ] String offset access ($str[0])
- [ ] String interpolation support

### Task 4.7: Array Functions (Basic)
**File**: `pkg/stdlib/array/functions.go`
**Effort**: 12 hours

- [ ] count() / sizeof()
- [ ] array_keys()
- [ ] array_values()
- [ ] array_merge()
- [ ] array_push()
- [ ] array_pop()
- [ ] array_shift()
- [ ] array_unshift()
- [ ] in_array()
- [ ] array_search()

**Note**: Full array functions in Phase 6

### Task 4.8: String Functions (Basic)
**File**: `pkg/stdlib/string/functions.go`
**Effort**: 10 hours

- [ ] strlen()
- [ ] substr()
- [ ] strpos()
- [ ] strrpos()
- [ ] str_replace()
- [ ] strtolower()
- [ ] strtoupper()
- [ ] trim() / ltrim() / rtrim()
- [ ] explode()
- [ ] implode()

**Note**: Full string functions in Phase 6

### Task 4.9: Testing
**Effort**: 12 hours

- [ ] String operation tests
- [ ] Binary-safe string tests
- [ ] Array operation tests
- [ ] Array order tests
- [ ] Mixed key type tests
- [ ] Packed array tests
- [ ] Array opcode tests
- [ ] String opcode tests
- [ ] Performance benchmarks

## Milestones

### Milestone 4.1: Strings (Week 1)
- String type complete
- String operations working
- Binary-safe verified

### Milestone 4.2: Arrays (Week 2-3)
- Array type complete
- Array operations working
- Order maintenance verified
- Mixed keys working

### Milestone 4.3: Optimizations (Week 4)
- Packed arrays working
- Performance improvements measured

### Milestone 4.4: Opcodes (Week 5)
- Array opcodes complete
- String opcodes complete
- Integration with VM

### Milestone 4.5: Complete Phase 4 (Week 6)
- All data structures working
- Basic stdlib functions
- Full test coverage

## Estimated Timeline

**Total Effort**: ~90 hours (5-6 weeks)

- Strings: ~10 hours
- Arrays: ~16 hours
- Packed arrays: ~8 hours
- Resources: ~4 hours
- Array opcodes: ~12 hours
- String opcodes: ~6 hours
- Basic functions: ~22 hours
- Testing: ~12 hours

## Progress Tracking

- [ ] Task 4.1: String Implementation (0%)
- [ ] Task 4.2: Array Implementation (0%)
- [ ] Task 4.3: Packed Array Optimization (0%)
- [ ] Task 4.4: Resource Implementation (0%)
- [ ] Task 4.5: Array Opcodes (0%)
- [ ] Task 4.6: String Opcodes (0%)
- [ ] Task 4.7: Array Functions (0%)
- [ ] Task 4.8: String Functions (0%)
- [ ] Task 4.9: Testing (0%)

**Overall Phase 4 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
