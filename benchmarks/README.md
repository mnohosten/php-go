# PHP-Go Macro Benchmarks

This directory contains macro-benchmark tests for php-go, measuring end-to-end performance of complete PHP scripts and application-level operations.

## Overview

Macro-benchmarks test full script execution including:
- Lexing
- Parsing
- Compilation
- VM execution

This is in contrast to micro-benchmarks (in `pkg/*/bench_test.go`) which test individual components.

## Structure

```
benchmarks/
├── README.md                 # This file
├── macro_bench_test.go       # Go benchmark tests
├── scripts/
│   ├── bench.php            # Full Zend benchmark suite (original syntax)
│   └── simple_bench.php     # Simplified benchmarks (php-go compatible)
└── compare.sh               # Script to compare php-go vs PHP performance
```

## Running Benchmarks

### Run all benchmarks
```bash
go test -bench=. ./benchmarks
```

### Run specific benchmark
```bash
go test -bench=BenchmarkSimpleLoop ./benchmarks
```

### Run with memory profiling
```bash
go test -bench=. -benchmem ./benchmarks
```

### Run for longer duration (more accurate)
```bash
go test -bench=. -benchtime=5s ./benchmarks
```

### Compare with PHP
```bash
# Run comparison tests (requires PHP in PATH)
go test -run=TestComparison ./benchmarks -v
```

## Available Benchmarks

### Successfully Running Benchmarks

1. **BenchmarkSimpleLoop** - Simple loop with integer increment
   - ~4.1ms per iteration
   - Tests basic loop and arithmetic performance

2. **BenchmarkFunctionCalls** - Function call overhead
   - ~3.8ms per iteration
   - Tests function call and return mechanisms

3. **BenchmarkRecursion** - Recursive Fibonacci (fibo(15))
   - ~8.3μs per iteration
   - Tests recursive function calls and stack management

4. **BenchmarkStringConcatenation** - String operations
   - ~255μs per iteration
   - Tests string concatenation and memory allocation

### Benchmarks Requiring Additional Features

The following benchmarks require stdlib functions not yet implemented:

- **BenchmarkArrayOperations** - Requires `count()`
- **BenchmarkHashOperations** - Requires `dechex()`
- **BenchmarkMatrixMultiplication** - Uses `count()` in loop conditions
- **BenchmarkFullScript** - Full suite requires multiple stdlib functions

## Benchmark Results (Apple M4 Max)

```
BenchmarkSimpleLoop-14              	     135	   4114191 ns/op	 3857935 B/op	  199195 allocs/op
BenchmarkFunctionCalls-14           	     158	   3768979 ns/op	 5783199 B/op	  179736 allocs/op
BenchmarkRecursion-14               	   73688	      8265 ns/op	   28688 B/op	     256 allocs/op
BenchmarkStringConcatenation-14     	    2368	    255196 ns/op	  255601 B/op	   12528 allocs/op
```

### Key Observations

1. **Recursion is efficient** - Fibonacci(15) runs in ~8μs with only 256 allocations
2. **Function calls have reasonable overhead** - ~3.8ms for 5000 calls
3. **String operations** - String concatenation shows room for optimization (many allocations)
4. **Memory usage** - High allocation counts in loops suggest opportunities for optimization

## Comparison with PHP

To compare php-go performance with standard PHP:

```bash
# Run comparison tests
go test -run=TestComparison -v ./benchmarks
```

This will execute the same script in both PHP and php-go, showing:
- Execution time comparison
- Performance ratio (php-go/PHP)
- Output verification

## Future Work

### Short Term
1. Implement missing stdlib functions (`count()`, `dechex()`, etc.)
2. Add more real-world application benchmarks
3. Profile and optimize hot paths
4. Reduce memory allocations in loops

### Long Term
1. Add benchmarks for parallel execution (Phase 7 features)
2. Compare with PHP + opcache
3. Add web application benchmarks (WordPress, Laravel)

## Memory Profiling

Comprehensive memory profiling tools are available to identify allocation bottlenecks:

### Quick Profile
```bash
# Fast profile (1 second, SimpleLoop only)
./benchmarks/quick_profile.sh
```

### Full Profiling Suite
```bash
# Profile all components (VM, Parser, Compiler, Lexer, Types)
./benchmarks/profile_memory.sh

# Analyze profiles and generate reports
./benchmarks/analyze_memory.sh
```

### View Profiles
```bash
# Interactive web viewer
go tool pprof -http=:8080 benchmarks/profiles/quick_mem.prof

# Text report
go tool pprof -text -alloc_space benchmarks/profiles/quick_mem.prof | head -30
```

### Documentation
See `docs/MEMORY_PROFILING.md` for:
- Detailed profiling guide
- Sample profile analysis
- Optimization opportunities
- Profiling workflow
- Integration with development

**Key Findings** (from memory profiling):
- SimpleLoop: 199K allocations per run (10K iterations)
- Primary bottleneck: Value object creation (88% of allocations)
- Critical optimizations identified:
  1. Value object pooling → 80-90% allocation reduction
  2. Integer cache (-128 to 1024) → 60-70% NewInt() reduction
  3. Boolean singletons → 100% NewBool() reduction

## Notes

- Some benchmarks use explicit `$i = $i + 1` instead of `$i++` because the increment/decrement operators are not fully supported yet
- All benchmarks use currently supported php-go syntax
- Full Zend benchmark suite (`scripts/bench.php`) preserved for future compatibility testing

## Performance Goals

Target performance metrics:
- **Within 2x of PHP 8.4 + opcache** for most operations
- **Better than PHP** for parallel operations (Phase 7)
- **Competitive memory usage** with PHP

Current status:
- Core operations functional but not yet optimized
- Focus is on correctness and feature completeness
- Performance optimization planned for later phases
