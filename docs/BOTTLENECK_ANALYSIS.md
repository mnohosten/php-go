# PHP-Go Performance Bottleneck Analysis

**Date**: November 24, 2025
**Platform**: Apple M4 Max (darwin/arm64, 14 cores)
**Analysis Phase**: Phase 10.6 - Performance Benchmarks

## Executive Summary

This document identifies performance bottlenecks in the PHP-Go interpreter based on comprehensive benchmark analysis across all core components (lexer, parser, compiler, types, and macro-level operations).

### Key Findings

1. **Memory Allocation is the Primary Bottleneck** - Excessive allocations across all components
2. **Parser Has Highest Per-Operation Overhead** - 10-25x slower than lexer
3. **String Operations Need Optimization** - High allocation rate in string handling
4. **Loop Operations Have Allocation Issues** - ~200K allocations for 10K loop iterations

---

## 1. Lexer Performance Analysis

### Benchmark Results (3-second runs)

| Benchmark | Time/op | Bytes/op | Allocs/op | Throughput |
|-----------|---------|----------|-----------|------------|
| SimpleTokens | 316ns | 32 B | 8 | 10.9M ops/s |
| VariablesAndOperators | 1.2μs | 128 B | 35 | 3.0M ops/s |
| StringLiterals | 763ns | 352 B | 25 | 4.8M ops/s |
| HeredocTokenization | 392ns | 256 B | 7 | 9.3M ops/s |
| ComplexExpression | 1.5μs | 208 B | 49 | 2.4M ops/s |
| FunctionDeclaration | 1.0μs | 96 B | 24 | 3.5M ops/s |
| ClassDeclaration | 1.3μs | 104 B | 27 | 2.8M ops/s |
| ControlFlow | 1.4μs | 160 B | 35 | 2.6M ops/s |
| LargeFile | 6.0μs | 704 B | 153 | 599K ops/s |
| NumericLiterals | 598ns | 56 B | 14 | 6.0M ops/s |
| Comments | 844ns | 44 B | 11 | 4.4M ops/s |
| StringInterpolation | 166ns | 416 B | 9 | 20.9M ops/s |

### Identified Bottlenecks

#### 1.1 High Allocation Count in Complex Expressions
- **Issue**: ComplexExpression shows 49 allocations for 1.5μs operation
- **Impact**: ~33 allocs/μs - very high allocation rate
- **Root Cause**: Token struct creation, position tracking, keyword lookups
- **Severity**: **MEDIUM** - Impacts parsing speed for expression-heavy code

#### 1.2 Large File Processing Inefficiency
- **Issue**: LargeFile benchmark shows 153 allocations (704B) for 6.0μs
- **Impact**: ~25.5 allocs/μs with 4x time increase vs simple tokens
- **Root Cause**: Linear scaling of token creation without pooling
- **Severity**: **MEDIUM** - Affects large file processing

#### 1.3 String Literal Allocation Overhead
- **Issue**: StringLiterals allocates 352B across 25 allocations
- **Impact**: ~14B per allocation - small but frequent
- **Root Cause**: String escaping, buffer allocation, interpolation detection
- **Severity**: **LOW** - Acceptable for string processing

### Recommendations

1. **Implement Token Pooling**
   - Create sync.Pool for Token structs to reduce allocations
   - Expected improvement: 30-40% reduction in allocs for large files
   - Priority: **HIGH**

2. **Optimize Position Tracking**
   - Consider embedding Position in Token instead of pointer
   - Expected improvement: 1 allocation saved per token
   - Priority: **MEDIUM**

3. **Buffer Reuse for String Processing**
   - Implement buffer pooling for string escape handling
   - Expected improvement: 20% reduction in string literal allocations
   - Priority: **MEDIUM**

---

## 2. Parser Performance Analysis

### Benchmark Results (3-second runs)

| Benchmark | Time/op | Bytes/op | Allocs/op | Throughput |
|-----------|---------|----------|-----------|------------|
| SimpleExpression | 3.6μs | 7,796 B | 137 | 961K ops/s |
| ComplexExpression | 5.8μs | 11,592 B | 214 | 625K ops/s |
| ArrayLiteral | 5.5μs | 11,304 B | 210 | 650K ops/s |
| IfStatement | 4.2μs | 8,200 B | 156 | 856K ops/s |
| ForLoop | 4.5μs | 9,000 B | 171 | 789K ops/s |
| WhileLoop | 4.2μs | 8,472 B | 160 | 811K ops/s |
| SwitchStatement | 4.7μs | 8,632 B | 167 | 803K ops/s |
| TryCatchFinally | 4.4μs | 8,608 B | 163 | 820K ops/s |
| FunctionDeclaration | 6.6μs | 13,968 B | 196 | 587K ops/s |
| SimpleClass | 6.5μs | 11,176 B | 199 | 555K ops/s |
| ComplexClass | 11.5μs | 18,768 B | 351 | 303K ops/s |
| Interface | 4.2μs | 8,128 B | 150 | 864K ops/s |
| Trait | 5.4μs | 10,344 B | 186 | 670K ops/s |
| MethodChaining | 4.7μs | 9,360 B | 173 | 759K ops/s |
| LargeFile | 25.8μs | 44,488 B | 840 | 139K ops/s |
| TypeHints | 8.0μs | 19,932 B | 226 | 457K ops/s |

### Identified Bottlenecks

#### 2.1 Excessive Memory Allocation Per Operation
- **Issue**: Even SimpleExpression allocates 7.8KB across 137 allocations
- **Impact**: Parser is 10-25x slower than lexer due to AST node creation
- **Root Cause**: Every expression creates multiple AST nodes + slice allocations
- **Severity**: **CRITICAL** - Primary bottleneck in parsing pipeline

#### 2.2 Complex Class Parsing Performance
- **Issue**: ComplexClass takes 11.5μs with 351 allocations (18.7KB)
- **Impact**: 3.5x slower than simple expressions
- **Root Cause**: Method parsing, property parsing, inheritance tracking
- **Severity**: **HIGH** - Impacts OOP-heavy codebases

#### 2.3 Large File Scalability
- **Issue**: LargeFile shows 840 allocations (44.5KB) for 25.8μs
- **Impact**: ~32.5 allocs/μs - high allocation rate
- **Root Cause**: No AST node pooling, deep recursion, slice growth
- **Severity**: **HIGH** - Affects application startup time

#### 2.4 Type Hint Processing Overhead
- **Issue**: TypeHints benchmark shows 8.0μs with 226 allocations (19.9KB)
- **Impact**: 2x slower than simple expressions
- **Root Cause**: Complex type parsing (union, nullable, intersection types)
- **Severity**: **MEDIUM** - Modern PHP uses types extensively

### Recommendations

1. **Implement AST Node Pooling**
   - Create pools for frequently-used AST nodes (Identifier, BinaryExpr, etc.)
   - Expected improvement: 40-50% reduction in allocations
   - Priority: **CRITICAL**

2. **Optimize Slice Pre-allocation**
   - Pre-allocate slices for common cases (parameter lists, statement blocks)
   - Expected improvement: 15-20% reduction in allocations
   - Priority: **HIGH**

3. **Reduce AST Node Depth**
   - Flatten AST structure where possible (e.g., binary expression chains)
   - Expected improvement: 10-15% reduction in allocations
   - Priority: **MEDIUM**

4. **Implement Incremental Parsing**
   - For large files, parse function bodies lazily
   - Expected improvement: 50%+ for large files with unused code
   - Priority: **LOW** (future optimization)

---

## 3. Compiler Performance Analysis

### Benchmark Results (3-second runs)

| Benchmark | Time/op | Bytes/op | Allocs/op | Throughput |
|-----------|---------|----------|-----------|------------|
| SimpleExpression | 809ns | 3,296 B | 19 | 4.4M ops/s |
| ArithmeticOperations | 1.1μs | 4,600 B | 21 | 3.2M ops/s |
| ComplexExpression | 1.7μs | 6,532 B | 24 | 2.0M ops/s |
| ArrayLiteral | 2.3μs | 8,664 B | 34 | 1.6M ops/s |
| IfStatement | 826ns | 2,664 B | 23 | 4.4M ops/s |
| ForLoop | 837ns | 2,688 B | 22 | 4.7M ops/s |
| WhileLoop | 687ns | 2,432 B | 20 | 5.2M ops/s |
| ForeachLoop | 941ns | 3,640 B | 25 | 3.8M ops/s |
| SwitchStatement | 1.3μs | 4,544 B | 33 | 2.8M ops/s |
| TryCatchFinally | 1.1μs | 4,624 B | 25 | 3.3M ops/s |
| FunctionDeclaration | 1.2μs | 4,240 B | 25 | 2.9M ops/s |
| FunctionCall | 942ns | 3,490 B | 26 | 3.8M ops/s |
| SimpleClass | 1.6μs | 5,600 B | 39 | 2.3M ops/s |
| ComplexClass | 2.6μs | 9,280 B | 57 | 1.4M ops/s |
| ObjectInstantiation | 1.6μs | 6,560 B | 24 | 2.2M ops/s |
| PropertyAccess | 1.2μs | 4,800 B | 29 | 3.0M ops/s |
| MethodCall | 1.6μs | 5,672 B | 36 | 2.3M ops/s |
| StaticAccess | 719ns | 2,769 B | 20 | 5.0M ops/s |
| ArrayAccess | 1.2μs | 4,624 B | 25 | 3.1M ops/s |
| Closure | 1.6μs | 6,400 B | 35 | 2.3M ops/s |
| LargeFile | 15.1μs | 53,784 B | 135 | 239K ops/s |
| ConstantFolding | 1.2μs | 4,712 B | 27 | 3.1M ops/s |
| SymbolTableOperations | 1.4μs | 4,832 B | 24 | 2.7M ops/s |
| ConstantTableOperations | 1.9μs | 7,440 B | 31 | 1.9M ops/s |
| JumpPatching | 2.1μs | 8,200 B | 34 | 1.7M ops/s |

### Identified Bottlenecks

#### 3.1 Instruction Buffer Allocation
- **Issue**: Even simple operations allocate 2-3KB for instruction bytecode
- **Impact**: High memory usage per compilation unit
- **Root Cause**: Instruction slice grows dynamically without pre-sizing
- **Severity**: **MEDIUM** - Impacts memory footprint

#### 3.2 Symbol Table Overhead
- **Issue**: Symbol table operations show consistent allocation overhead
- **Impact**: Every variable access/assignment touches symbol table
- **Root Cause**: Map allocations, string comparisons, scope tracking
- **Severity**: **LOW** - Necessary for correctness

#### 3.3 Array Literal Compilation
- **Issue**: Array literals take 2.3μs with 34 allocations (8.7KB)
- **Impact**: 3x slower than simple expressions
- **Root Cause**: Each array element creates instruction sequence
- **Severity**: **MEDIUM** - Arrays are common in PHP

#### 3.4 Complex Class Compilation
- **Issue**: ComplexClass compiles in 2.6μs with 57 allocations (9.3KB)
- **Impact**: Highest allocation count in compiler benchmarks
- **Root Cause**: Method compilation, property initialization, inheritance setup
- **Severity**: **MEDIUM** - OOP overhead

### Recommendations

1. **Pre-size Instruction Buffers**
   - Estimate instruction count from AST node count
   - Expected improvement: 20-30% reduction in allocations
   - Priority: **HIGH**

2. **Optimize Constant Table**
   - Use string interning for constant identifiers
   - Expected improvement: 15-20% memory savings
   - Priority: **MEDIUM**

3. **Batch Instruction Emission**
   - Emit multiple instructions at once for common patterns
   - Expected improvement: 10-15% time savings
   - Priority: **MEDIUM**

4. **Symbol Table Caching**
   - Cache symbol lookups within expression compilation
   - Expected improvement: 5-10% in variable-heavy code
   - Priority: **LOW**

---

## 4. Type System Performance Analysis

### Benchmark Results (3-second runs)

| Benchmark | Time/op | Bytes/op | Allocs/op | Notes |
|-----------|---------|----------|-----------|-------|
| NewInt | 0.23ns | 0 B | 0 | Excellent |
| NewFloat | 0.23ns | 0 B | 0 | Excellent |
| NewString | 0.24ns | 0 B | 0 | Excellent |
| NewBool | 0.24ns | 0 B | 0 | Excellent |
| NewNull | 0.23ns | 0 B | 0 | Excellent |
| IntToString | 12.3ns | 5 B | 1 | 1 alloc for string |
| IntToFloat | 1.25ns | 0 B | 0 | Excellent |
| IntToBool | 1.25ns | 0 B | 0 | Excellent |
| StringToInt | 5.24ns | 0 B | 0 | Excellent |
| StringToFloat | 17.8ns | 0 B | 0 | Good |
| StringToBool | 1.51ns | 0 B | 0 | Excellent |
| FloatToInt | 1.28ns | 0 B | 0 | Excellent |
| FloatToString | 34.2ns | 8 B | 1 | 1 alloc for string |
| ValueEquals | 1.59ns | 0 B | 0 | Excellent |
| ValueIdentical | 1.55ns | 0 B | 0 | Excellent |
| StringEquality | 2.48ns | 0 B | 0 | Excellent |
| TypeJugglingIntString | 14.7ns | 0 B | 0 | Excellent |
| TypeJugglingStringInt | 14.8ns | 0 B | 0 | Excellent |
| ArrayCreation | 1.23ns | 0 B | 0 | Excellent |
| ValueType | 0.23ns | 0 B | 0 | Excellent |
| ValueIsCallable | 0.23ns | 0 B | 0 | Excellent |
| StringConcatenation | 27.4ns | 0 B | 0 | Excellent |
| ValueCopy | 0.24ns | 0 B | 0 | Excellent |
| ComplexTypeConversion | 18.0ns | 0 B | 0 | Excellent |
| MultipleValueCreation | 0.23ns | 0 B | 0 | Excellent |
| BoolConversions | 5.01ns | 0 B | 0 | Excellent |
| NumericStringParsing | 38.7ns | 0 B | 0 | Good |

### Identified Bottlenecks

#### 4.1 String Conversion Allocations
- **Issue**: Float/Int to string conversions allocate (5-8B per conversion)
- **Impact**: Affects echo, print, string concatenation in loops
- **Root Cause**: Go's strconv allocates new strings
- **Severity**: **LOW** - Necessary allocations, but optimizable

#### 4.2 Numeric String Parsing
- **Issue**: NumericStringParsing takes 38.7ns (highest in type system)
- **Impact**: Type juggling in comparisons
- **Root Cause**: Complex parsing logic for numeric strings
- **Severity**: **LOW** - Still very fast in absolute terms

### Positive Findings

- **Value creation is essentially free** (sub-nanosecond with 0 allocations)
- **Type juggling is highly optimized** (~15ns with 0 allocations)
- **Comparisons are fast** (1-2ns with 0 allocations)
- **Type system is NOT a bottleneck** - excellent performance across the board

### Recommendations

1. **String Conversion Caching**
   - Cache common int→string conversions (0-1000)
   - Expected improvement: 50%+ for loops with string output
   - Priority: **MEDIUM**

2. **Numeric String Fast Path**
   - Add fast path for simple numeric strings ("123", "45.6")
   - Expected improvement: 30-40% for numeric string handling
   - Priority: **LOW**

---

## 5. Macro-Level Performance Analysis

### Benchmark Results (Partial - some tests failed due to missing stdlib)

| Benchmark | Time/op | Bytes/op | Allocs/op | Notes |
|-----------|---------|----------|-----------|-------|
| SimpleLoop (10K iters) | 4.1ms | 3.9MB | 199,195 | **CRITICAL ISSUE** |
| FunctionCalls (5K calls) | 3.9ms | 5.8MB | 179,736 | **HIGH ISSUE** |
| Recursion (Fib 15) | 8.3μs | 28.7KB | 256 | Good |

### Identified Bottlenecks

#### 5.1 CRITICAL: Loop Allocation Explosion
- **Issue**: 10K loop iterations cause 199,195 allocations (3.9MB)
- **Impact**: **~20 allocations per loop iteration** - completely unacceptable
- **Root Cause**: VM instruction dispatch, value creation, stack operations
- **Severity**: **CRITICAL** - Makes interpreter unusable for production
- **Expected Behavior**: Should be ~10-100 allocations for entire loop

#### 5.2 Function Call Overhead
- **Issue**: 5K function calls cause 179,736 allocations (5.8MB)
- **Impact**: **~36 allocations per function call**
- **Root Cause**: Frame creation, parameter passing, return value handling
- **Severity**: **HIGH** - Significantly impacts function-heavy code

#### 5.3 Recursion Performance (Acceptable)
- **Issue**: Fibonacci(15) shows 256 allocations but only 8.3μs
- **Impact**: Recursion is relatively efficient
- **Root Cause**: Lightweight frame creation
- **Severity**: **NONE** - This is acceptable

### Root Cause Analysis

The macro benchmarks reveal the compounding effect of micro-level allocations:

1. **Lexer** allocates 8-49 times per operation
2. **Parser** allocates 137-351 times per operation
3. **Compiler** allocates 19-57 times per operation
4. **Runtime** allocates **~20 times per loop iteration**

This creates a multiplicative effect where a simple loop becomes allocation-heavy.

### Recommendations (URGENT)

1. **VM Instruction Dispatch Optimization**
   - Pre-allocate stack frames
   - Reuse Value structs instead of creating new ones
   - Expected improvement: 80-90% reduction in loop allocations
   - Priority: **CRITICAL**

2. **Value Object Pooling**
   - Implement sync.Pool for frequently-created Values
   - Expected improvement: 60-70% reduction in allocations
   - Priority: **CRITICAL**

3. **Frame Pool Implementation**
   - Reuse call frames instead of allocating per call
   - Expected improvement: 70-80% reduction in function call allocations
   - Priority: **CRITICAL**

4. **Stack Pre-allocation**
   - Pre-allocate VM stack to reasonable size
   - Expected improvement: 30-40% reduction in stack operations
   - Priority: **HIGH**

---

## 6. Overall Bottleneck Summary

### Critical Bottlenecks (Must Fix for v1.0)

1. **VM Loop Allocation Explosion** - 20 allocs/iteration
   - Location: `pkg/vm/vm.go` Execute() loop
   - Fix: Value pooling, frame pooling, stack pre-allocation
   - Estimated effort: 8-12 hours
   - Impact: 80-90% performance improvement in real-world code

2. **Parser AST Node Allocations** - 137-351 allocs per parse operation
   - Location: `pkg/parser/` all parsing functions
   - Fix: AST node pooling, slice pre-allocation
   - Estimated effort: 12-16 hours
   - Impact: 40-50% parser performance improvement

### High-Priority Bottlenecks

3. **Function Call Overhead** - 36 allocs/call
   - Location: `pkg/vm/handlers_functions.go`
   - Fix: Frame pooling, parameter passing optimization
   - Estimated effort: 6-8 hours
   - Impact: 70% improvement in function-heavy code

4. **Compiler Instruction Buffer Growth**
   - Location: `pkg/compiler/compiler.go`
   - Fix: Pre-size instruction buffers based on AST size
   - Estimated effort: 4-6 hours
   - Impact: 20-30% memory reduction

### Medium-Priority Bottlenecks

5. **Lexer Token Creation** - 8-49 allocs per tokenization
   - Location: `pkg/lexer/lexer.go`
   - Fix: Token pooling
   - Estimated effort: 4-6 hours
   - Impact: 30-40% lexer performance improvement

6. **String Conversion Allocations** - Every ToString() allocates
   - Location: `pkg/types/value.go`
   - Fix: Common value caching (integers 0-1000)
   - Estimated effort: 3-4 hours
   - Impact: 50% improvement in output-heavy code

### Low-Priority Bottlenecks

7. **Array Literal Compilation** - 34 allocs per array
   - Location: `pkg/compiler/compiler.go`
   - Fix: Batch instruction emission for array initialization
   - Estimated effort: 4-5 hours
   - Impact: 15-20% improvement in array-heavy code

8. **Complex Class Parsing** - 351 allocs per complex class
   - Location: `pkg/parser/decl.go`
   - Fix: Incremental/lazy parsing for class bodies
   - Estimated effort: 8-10 hours
   - Impact: 25-30% improvement for OOP-heavy codebases

---

## 7. Optimization Roadmap

### Phase 1: Critical Fixes (Estimated 20-28 hours)
**Goal**: Make interpreter production-ready

1. Implement Value pooling (sync.Pool) - 6-8h
2. Implement Frame pooling - 6-8h
3. Pre-allocate VM stack - 4-6h
4. Implement AST node pooling - 4-6h

**Expected Result**: 5-10x improvement in loop-heavy code, 2-3x in general

### Phase 2: High-Priority Optimizations (Estimated 10-14 hours)
**Goal**: Competitive performance with PHP

1. Optimize function call overhead - 6-8h
2. Pre-size compiler instruction buffers - 4-6h

**Expected Result**: 2-3x improvement in function-heavy code

### Phase 3: Medium-Priority Optimizations (Estimated 7-10 hours)
**Goal**: Polish and edge-case performance

1. Implement Token pooling - 4-6h
2. Cache common string conversions - 3-4h

**Expected Result**: 20-30% improvement across the board

### Phase 4: Low-Priority Optimizations (Estimated 12-15 hours)
**Goal**: Target specific use cases

1. Batch array instruction emission - 4-5h
2. Lazy class body parsing - 8-10h

**Expected Result**: 15-30% improvement in specific patterns

---

## 8. Performance Goals

### Current State (Before Optimization)
- **Simple Loop (10K)**: 4.1ms with 199K allocs
- **Function Calls (5K)**: 3.9ms with 180K allocs
- **Parser (SimpleExpr)**: 3.6μs with 137 allocs

### Target State (After Phase 1-2)
- **Simple Loop (10K)**: <1.0ms with <20K allocs (4x faster, 10x fewer allocs)
- **Function Calls (5K)**: <1.5ms with <50K allocs (2.5x faster, 3.5x fewer allocs)
- **Parser (SimpleExpr)**: <2.5μs with <80 allocs (1.4x faster, 1.7x fewer allocs)

### Stretch Goals (After Phase 3-4)
- **Simple Loop (10K)**: <500μs with <10K allocs
- **Function Calls (5K)**: <1.0ms with <30K allocs
- **Competitive with PHP 8.4 + opcache** (within 2-3x)

---

## 9. Next Steps

1. **Immediate Action** (This Week):
   - Start with Value pooling implementation (Critical #1)
   - Measure impact with benchmarks
   - Document results

2. **Short Term** (Next 2 Weeks):
   - Complete Phase 1 optimizations
   - Run full benchmark suite
   - Compare with PHP 8.4 baseline

3. **Medium Term** (Next Month):
   - Complete Phases 2-3
   - Real-world testing with WordPress/Laravel
   - Performance profiling with pprof

4. **Long Term** (Next Quarter):
   - Phase 4 optimizations
   - JIT compilation research
   - Advanced optimization techniques

---

## 10. Measurement & Validation

### Before/After Metrics

For each optimization, measure:
1. **Time improvement** (ns/op or μs/op)
2. **Allocation reduction** (allocs/op)
3. **Memory reduction** (B/op)
4. **Regression risk** (test coverage)

### Benchmark Commands

```bash
# Run full benchmark suite
go test -bench=. -benchmem -benchtime=3s ./...

# Compare before/after
go test -bench=. -benchmem -benchtime=3s ./pkg/vm/ > before.txt
# ... make changes ...
go test -bench=. -benchmem -benchtime=3s ./pkg/vm/ > after.txt
benchstat before.txt after.txt

# Profile with pprof
go test -bench=BenchmarkSimpleLoop -cpuprofile=cpu.prof -memprofile=mem.prof ./benchmarks/
go tool pprof -http=:8080 cpu.prof
```

### Success Criteria

Phase 1 is successful if:
- Loop allocations drop from 199K → <20K (10x reduction)
- Loop time drops from 4.1ms → <1.0ms (4x improvement)
- No test regressions
- Memory usage stays under 100MB for 10K iterations

---

## Conclusion

The primary performance bottleneck in PHP-Go is **excessive memory allocation** at the VM execution layer, compounded by allocations in the parser and compiler. The critical path for optimization is:

1. **VM Value pooling** (CRITICAL - 6-8h) → 80-90% improvement
2. **VM Frame pooling** (CRITICAL - 6-8h) → 70-80% improvement
3. **Parser AST pooling** (CRITICAL - 4-6h) → 40-50% improvement

These three optimizations alone should make PHP-Go production-viable and competitive with interpreted PHP.

**Total Critical Path Effort**: ~20-28 hours
**Expected Overall Improvement**: **5-10x faster** with **10-20x fewer allocations**

This analysis provides a clear roadmap to transform PHP-Go from a research prototype to a production-ready interpreter.
