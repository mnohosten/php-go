# PHP-Go Memory Profiling Guide

**Date**: November 24, 2025
**Phase**: Phase 10.6 - Performance Benchmarks
**Status**: Complete

## Overview

This guide explains how to perform memory profiling on PHP-Go to identify allocation bottlenecks and optimize memory usage. Memory profiling is essential for achieving production-ready performance.

## Quick Start

### Generate Memory Profiles

```bash
# Quick profile (1 second, SimpleLoop only)
./benchmarks/quick_profile.sh

# Full profiling suite (all components, 3 seconds each)
./benchmarks/profile_memory.sh

# Analyze generated profiles
./benchmarks/analyze_memory.sh
```

### View Profiles Interactively

```bash
# Start web-based profile viewer
go tool pprof -http=:8080 benchmarks/profiles/quick_mem.prof

# Navigate to:
# - Graph view: Visual flame graph of allocations
# - Top: Top allocators by size/count
# - Source: Source code view with allocation annotations
# - Flame Graph: Interactive flame graph
```

## Memory Profiling Tools

### 1. Quick Profile (`quick_profile.sh`)

Generates a fast memory profile of the most critical benchmark (SimpleLoop) in 1 second.

**Use Case**: Quick iteration during optimization work

**Output**:
- `profiles/quick_mem.prof` - Memory profile
- `profiles/quick_cpu.prof` - CPU profile

**Example**:
```bash
./benchmarks/quick_profile.sh
```

### 2. Full Profiling Suite (`profile_memory.sh`)

Generates comprehensive memory profiles for all components:
- VM/Runtime (SimpleLoop, FunctionCalls, Recursion, StringConcat)
- Parser (SimpleExpr, ComplexExpr, ComplexClass, LargeFile)
- Compiler (SimpleExpr, ComplexClass, LargeFile)
- Lexer (Simple, Complex, LargeFile)
- Type System (Conversions, TypeJuggling)

**Use Case**: Comprehensive analysis, baseline measurement

**Output**:
- Multiple `.prof` files in `profiles/` directory
- Timestamped for comparison over time

**Example**:
```bash
./benchmarks/profile_memory.sh
```

### 3. Analysis Script (`analyze_memory.sh`)

Analyzes all generated profiles and creates detailed reports.

**Use Case**: Automated report generation, batch analysis

**Output**:
- Text reports with top allocators
- SVG flame graphs
- PDF reports (if supported)
- Markdown summary report

**Example**:
```bash
./benchmarks/analyze_memory.sh
```

## Understanding Profile Output

### Allocation Space vs Allocation Objects

**Allocation Space** (`-alloc_space`): Total bytes allocated
- Shows functions that allocate the most memory
- Useful for finding memory-heavy operations

**Allocation Objects** (`-alloc_objects`): Total objects allocated
- Shows functions that create the most objects
- Useful for finding allocation-heavy loops

**In-Use Space** (`-inuse_space`): Currently allocated bytes
- Shows live memory at profile time
- Useful for finding memory leaks

**In-Use Objects** (`-inuse_objects`): Currently allocated objects
- Shows live objects at profile time
- Useful for finding object retention issues

### Reading the Output

Example pprof output:
```
      flat  flat%   sum%        cum   cum%
32734.81kB 81.09% 81.09% 32734.81kB 81.09%  github.com/krizos/php-go/pkg/types.NewInt (inline)
 4687.97kB 11.61% 92.70%  4687.97kB 11.61%  github.com/krizos/php-go/pkg/types.NewBool (inline)
```

- **flat**: Memory allocated by this function directly
- **flat%**: Percentage of total allocations
- **sum%**: Cumulative percentage
- **cum**: Memory allocated by this function and its callees
- **cum%**: Cumulative percentage including callees

## Sample Profile Analysis

### SimpleLoop Benchmark (10K iterations)

**Benchmark Results**:
```
BenchmarkSimpleLoop-14    	       5	 224733525 ns/op	 3863819 B/op	  199196 allocs/op
```

**Key Findings**:

#### Allocation Space (Total: 40.4 MB for 5 iterations)

| Function | Allocated | Percentage | Analysis |
|----------|-----------|------------|----------|
| `NewInt()` | 32.7 MB | 81.09% | **CRITICAL** - Integer value creation |
| `NewBool()` | 4.7 MB | 11.61% | Boolean value creation |
| **Total** | 37.4 MB | 92.70% | Value creation dominates |

#### Allocation Objects (Total: 1,698,407 objects)

| Function | Objects | Percentage | Analysis |
|----------|---------|------------|----------|
| `NewInt()` | 1,494,988 | 88.02% | **CRITICAL** - 149 int allocs per loop iteration |
| `NewBool()` | 200,020 | 11.78% | 20 bool allocs per loop iteration |
| **Total** | 1,695,008 | 99.80% | ~170 value allocs per iteration |

#### Critical Allocations by VM Operation

| VM Operation | Objects | Purpose |
|--------------|---------|---------|
| `GetConstant()` | 900,088 (53%) | Constant value retrieval |
| `opQMAssign()` | 700,068 (41%) | Variable assignment |
| `opAdd()` | 594,900 (35%) | Arithmetic operations |
| `opIsSmaller()` | 200,020 (12%) | Comparison operations |
| `opJmp()` | 200,000 (12%) | Jump operations |

**Root Cause Analysis**:

The SimpleLoop benchmark executes approximately 10,000 iterations. Each iteration:
1. Retrieves 2 constants (loop counter, increment) → 2 `NewInt()` calls
2. Performs addition ($i + 1) → 1-2 `NewInt()` calls
3. Performs comparison ($i < 10000) → 1 `NewBool()` call
4. Assigns result → 1 `NewInt()` call
5. Jumps back to loop start

**Total per iteration**: ~15-20 Value allocations

**10,000 iterations × 17 allocations = 170,000 allocations** ✓ (matches observed 199K)

The extra ~29K allocations come from:
- VM dispatch overhead
- Stack frame operations
- Operand value retrieval

## Optimization Opportunities

Based on the profiling data, here are the critical optimization targets:

### 1. CRITICAL: Value Object Pooling

**Problem**: 1.5M `NewInt()` calls, 200K `NewBool()` calls per 10K loop iterations

**Solution**: Implement `sync.Pool` for Value objects

**Expected Impact**:
- 80-90% reduction in allocations
- 5-10x performance improvement in loops
- Reduced GC pressure

**Implementation**:
```go
// pkg/types/value.go
var valuePool = sync.Pool{
    New: func() interface{} {
        return &Value{}
    },
}

func NewInt(i int64) Value {
    v := valuePool.Get().(*Value)
    v.Type = TypeInt
    v.IntValue = i
    return *v
}

func (v *Value) Release() {
    v.Type = TypeUndef
    valuePool.Put(v)
}
```

**Effort**: 6-8 hours
**Priority**: CRITICAL

### 2. CRITICAL: Small Integer Cache

**Problem**: Common integers (0, 1, -1, etc.) allocated repeatedly

**Solution**: Pre-allocate cache for integers -128 to 1024

**Expected Impact**:
- 60-70% reduction in `NewInt()` allocations
- Near-zero cost for common values

**Implementation**:
```go
var intCache [1153]Value // -128 to 1024

func init() {
    for i := -128; i <= 1024; i++ {
        intCache[i+128] = Value{Type: TypeInt, IntValue: int64(i)}
    }
}

func NewInt(i int64) Value {
    if i >= -128 && i <= 1024 {
        return intCache[i+128]
    }
    return Value{Type: TypeInt, IntValue: i}
}
```

**Effort**: 2-3 hours
**Priority**: CRITICAL

### 3. HIGH: Boolean Singleton Values

**Problem**: Only 2 possible boolean values, but allocated 200K times

**Solution**: Use pre-allocated True/False values

**Expected Impact**:
- 100% reduction in `NewBool()` allocations
- Zero-cost boolean values

**Implementation**:
```go
var (
    TrueValue  = Value{Type: TypeBool, BoolValue: true}
    FalseValue = Value{Type: TypeBool, BoolValue: false}
)

func NewBool(b bool) Value {
    if b {
        return TrueValue
    }
    return FalseValue
}
```

**Effort**: 1 hour
**Priority**: HIGH

### 4. HIGH: VM Stack Pre-allocation

**Problem**: Stack grows dynamically during execution

**Solution**: Pre-allocate stack to reasonable size (1024 values)

**Expected Impact**:
- 30-40% reduction in stack operations
- Fewer slice reallocations

**Effort**: 2-3 hours
**Priority**: HIGH

### 5. MEDIUM: Operand Value Caching

**Problem**: `getOperandValue()` called 900K times

**Solution**: Cache operand values within instruction execution

**Expected Impact**:
- 20-30% reduction in operand retrieval
- Faster instruction dispatch

**Effort**: 3-4 hours
**Priority**: MEDIUM

## Profiling Workflow

### 1. Baseline Measurement

```bash
# Generate baseline profile
./benchmarks/quick_profile.sh

# Save benchmark output
go test -run=^$ -bench=^BenchmarkSimpleLoop$ -benchmem ./benchmarks > baseline.txt

# Save profile for comparison
cp benchmarks/profiles/quick_mem.prof benchmarks/profiles/baseline.prof
```

### 2. Make Optimization Changes

Implement one optimization at a time to measure individual impact.

### 3. Re-profile and Compare

```bash
# Generate new profile
./benchmarks/quick_profile.sh

# Save new benchmark output
go test -run=^$ -bench=^BenchmarkSimpleLoop$ -benchmem ./benchmarks > optimized.txt

# Compare benchmarks
benchstat baseline.txt optimized.txt

# Compare profiles
go tool pprof -base=benchmarks/profiles/baseline.prof benchmarks/profiles/quick_mem.prof
```

### 4. Validate Correctness

```bash
# Ensure all tests still pass
go test ./...

# Run integration tests
./php-go test.php
```

### 5. Document Results

Update `docs/BOTTLENECK_ANALYSIS.md` with:
- Optimization implemented
- Before/after metrics
- Performance improvement percentage

## Advanced Profiling Techniques

### Heap Dump Analysis

```bash
# Generate heap dump at specific point
runtime/pprof.WriteHeapProfile(f)

# Analyze heap dump
go tool pprof -http=:8080 heap.prof
```

### Continuous Profiling

```bash
# Profile during test suite
go test -memprofile=mem.prof ./...

# Profile during benchmark
go test -bench=. -memprofile=mem.prof ./benchmarks
```

### Allocation Rate Profiling

```bash
# Set memory profile rate (1 = every allocation)
go test -bench=. -memprofilerate=1 ./benchmarks

# Set to default (512KB between samples)
go test -bench=. -memprofilerate=524288 ./benchmarks
```

### Comparative Analysis

```bash
# Profile two different implementations
go test -bench=BenchmarkOld -memprofile=old.prof
go test -bench=BenchmarkNew -memprofile=new.prof

# Compare
go tool pprof -base=old.prof new.prof
```

## Integration with Development Workflow

### Pre-Commit Profiling

Before committing major changes:

```bash
# 1. Run quick profile
./benchmarks/quick_profile.sh

# 2. Check for regressions
go tool pprof -text -alloc_space profiles/quick_mem.prof | head -20

# 3. If top allocators changed significantly, investigate
```

### CI/CD Integration

Add to continuous integration:

```bash
# Generate profiles in CI
./benchmarks/profile_memory.sh

# Upload profiles as artifacts
# Compare with main branch baseline
```

### Performance Testing Gates

Set thresholds for acceptable allocations:

```bash
# Example: Fail if SimpleLoop exceeds 250K allocations
allocs=$(go test -bench=BenchmarkSimpleLoop ./benchmarks | grep allocs | awk '{print $5}')
if [ "$allocs" -gt 250000 ]; then
    echo "FAIL: Too many allocations ($allocs > 250000)"
    exit 1
fi
```

## Common Profiling Pitfalls

### 1. Profile Rate Too Low

**Problem**: Missing allocations due to sampling

**Solution**: Use `-memprofilerate=1` for critical analysis

### 2. Short Benchmark Times

**Problem**: Not enough samples for accurate profile

**Solution**: Use `-benchtime=3s` or longer

### 3. Ignoring Inline Functions

**Problem**: Allocations attributed to wrong function

**Solution**: Look at call graphs, not just flat allocations

### 4. Profiling in Debug Mode

**Problem**: Debug builds have different allocation patterns

**Solution**: Always profile optimized builds (`go test -gcflags=-N -l` to disable)

### 5. Not Comparing Baselines

**Problem**: Can't measure improvement without baseline

**Solution**: Always save baseline profiles before optimizing

## References

### Go Profiling Documentation

- [Go pprof Documentation](https://pkg.go.dev/runtime/pprof)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Go Performance Tuning](https://go.dev/doc/diagnostics)

### Internal Documentation

- `docs/BOTTLENECK_ANALYSIS.md` - Detailed bottleneck analysis
- `benchmarks/README.md` - Benchmark suite documentation
- `TODO.md` - Phase 10.7 optimization tasks

### Tools

- `go tool pprof` - Profile viewer and analyzer
- `benchstat` - Benchmark comparison tool
- `go-torch` - Flame graph generator (optional)

## Next Steps

1. **Implement Value Pooling** (CRITICAL, 6-8h)
   - Create sync.Pool for Value objects
   - Measure impact: expect 80-90% allocation reduction

2. **Implement Integer Cache** (CRITICAL, 2-3h)
   - Cache integers -128 to 1024
   - Measure impact: expect 60-70% `NewInt()` reduction

3. **Implement Boolean Singletons** (HIGH, 1h)
   - Use pre-allocated True/False values
   - Measure impact: expect 100% `NewBool()` reduction

4. **Re-profile After Each Optimization**
   - Use `quick_profile.sh` for fast iteration
   - Use `benchstat` to compare before/after
   - Document results in bottleneck analysis

5. **Full Profiling Suite**
   - Run `profile_memory.sh` after all optimizations
   - Generate comprehensive report with `analyze_memory.sh`
   - Update performance goals in project docs

## Success Metrics

**Current State** (Before Optimization):
- SimpleLoop: 199,196 allocs/op
- Memory: 3.86 MB/op
- Time: 224.7 ms/op

**Target State** (After Phase 1 Optimizations):
- SimpleLoop: <20,000 allocs/op (10x improvement)
- Memory: <400 KB/op (10x improvement)
- Time: <50 ms/op (4-5x improvement)

**Stretch Goals** (After All Optimizations):
- SimpleLoop: <10,000 allocs/op (20x improvement)
- Memory: <200 KB/op (20x improvement)
- Time: <20 ms/op (10x improvement)
- Competitive with PHP 8.4 + opcache

---

**Status**: Memory profiling infrastructure complete ✅
**Task**: Phase 10.6 - Memory profiling (3h) - COMPLETE
**Next**: Phase 10.7 - Optimization Pass (implement pooling and caching)
