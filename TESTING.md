# PHP-Go Testing Guide

This guide shows you how to test PHP-Go's parallelization features.

## Quick Start

### 1. Build the Project

```bash
go build -o bin/php-go ./cmd/php-go
```

### 2. Run the Interactive Demo

```bash
./bin/php-go demo
```

This demonstrates:
- **Worker Pool**: 10 tasks processed by 4 workers (3x speedup)
- **Copy-on-Write**: Memory optimization with shared data (80% memory savings)
- **Parallel Array Map**: Automatic parallelization of array operations
- **Request Isolation**: Concurrent request handling like PHP-FPM

## Running Tests

### Run All Tests

```bash
go test ./pkg/parallel -v
```

This runs 267 tests across all parallel components:
- 100 core functionality tests
- 84 integration tests
- 23 race condition tests
- 45 benchmark tests
- 15 safety tests

### Run Specific Test Categories

**Integration Tests** (component interaction):
```bash
go test ./pkg/parallel -run "TestIntegration" -v
```

**Race Tests** (concurrency safety):
```bash
go test ./pkg/parallel -run "TestRace" -v
go test ./pkg/parallel -race  # with race detector
```

**Benchmarks** (performance measurement):
```bash
go test ./pkg/parallel -bench=. -benchmem
```

### Quick Test (Core Tests Only)

```bash
go test ./pkg/parallel -run "^Test[^IR]" -v
```

This skips Integration and Race tests, running only core functionality tests.

## Exploring Examples

### PHP Examples

See how PHP code uses parallelization:

```bash
# View the examples
cat examples/parallel_examples.php
cat examples/ecommerce_parallel.php
```

These show:
- 10 practical parallelization patterns
- Real-world e-commerce application
- Performance comparisons and metrics

### Technical Documentation

```bash
cat examples/parallel_technical.md
```

This shows how PHP code maps to Go implementation under the hood.

## Benchmark Results

Run benchmarks to see performance characteristics:

```bash
go test ./pkg/parallel -bench=BenchmarkWorkerPool -benchmem
go test ./pkg/parallel -bench=BenchmarkCOWArray -benchmem
go test ./pkg/parallel -bench=BenchmarkParallelArrayMap -benchmem
```

Expected results:
- **Worker Pool Submit**: ~1-2 µs per task
- **COW Clone**: ~50 ns (no actual copy)
- **COW Set (shared)**: ~500 ns (triggers copy)
- **Parallel vs Sequential**: 4x speedup with 4 workers (CPU-bound)

## Development Commands

### Tokenize PHP Code

```bash
./bin/php-go lex examples/parallel_examples.php
```

Shows lexer tokens from PHP file.

### Parse PHP Code

```bash
./bin/php-go parse examples/parallel_examples.php
```

Shows AST (Abstract Syntax Tree) from PHP file.

### JSON Output

```bash
./bin/php-go lex --json test.php > tokens.json
./bin/php-go parse --json test.php > ast.json
```

## Test Coverage

Current test coverage for `pkg/parallel`:

```bash
go test ./pkg/parallel -cover
```

Expected coverage: 85-90%

## Race Condition Detection

Run tests with race detector to ensure thread safety:

```bash
go test ./pkg/parallel -race
```

This checks for:
- Concurrent reads/writes to shared data
- Improper synchronization
- Data race conditions

All tests should pass with no race conditions detected.

## Performance Testing

### Worker Pool Performance

```bash
go test ./pkg/parallel -bench=BenchmarkWorkerPool -benchtime=5s
```

### COW Performance

```bash
go test ./pkg/parallel -bench=BenchmarkCOW -benchtime=5s
```

### Array Operations

```bash
go test ./pkg/parallel -bench=BenchmarkParallelArray -benchtime=5s
```

## Troubleshooting

### Tests Timeout

If tests timeout, increase timeout:
```bash
go test ./pkg/parallel -timeout=5m -v
```

### Build Fails

Ensure you have Go 1.21+:
```bash
go version
```

Clean and rebuild:
```bash
go clean
go build -o bin/php-go ./cmd/php-go
```

### Demo Fails

Check that parallel package builds:
```bash
go build ./pkg/parallel
```

## Next Steps

After testing the parallelization features:

1. Read the implementation in `pkg/parallel/`
2. Review test cases to understand usage patterns
3. Check TODO.md for upcoming features
4. Explore how parallelization will integrate with PHP execution

## Current Status

- ✅ Phase 7 Implementation: 85% complete (98h/115h)
- ✅ All core tests passing
- ✅ Race tests passing
- ✅ Benchmarks running
- ⏳ Safety Analyzer: Pending (Task 7.1)

For more details, see TODO.md in the project root.
