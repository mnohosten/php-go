# PHP-Go Test Suite

This directory contains comprehensive tests for the PHP-Go interpreter.

## Test Organization

### `examples_test.go`
Comprehensive test runner for all PHP example files in the project. This test suite ensures that all examples execute without errors and can be used for regression testing.

#### Test Functions

1. **TestExampleFiles** - Runs all example files and provides a summary report
   - Tests files in `examples/basic/`, `examples/parallel/`, `examples/oop/`, `examples/advanced/`
   - Reports pass/fail statistics
   - Provides detailed error messages for failures

2. **TestBasicExamples** - Focused tests for basic example files with custom output validation
   - Individual tests for each basic example file
   - Custom validation logic per file

3. **TestPhpTestFiles** - Runs PHP test files in `tests/php/` directory
   - Tests files in `tests/php/basic/`, `tests/php/stdlib/`, `tests/php/vm/`
   - Note: Some tests may timeout due to infinite loops in test files

4. **BenchmarkExampleExecution** - Benchmarks example file execution performance

## Running Tests

```bash
# Run all example tests
go test ./tests/ -v

# Run specific test
go test ./tests/ -v -run TestBasicExamples

# Run with timeout
go test ./tests/ -timeout 30s

# Run benchmarks
go test ./tests/ -bench=. -benchmem
```

## Current Status

### Example Files (as of Phase 6C.1)
- **Total**: 9 example files
- **Passing**: 7 (77.8%)
- **Failing**: 2 (22.2%)

**Passing Files**:
- ✅ examples/basic/arrays.php
- ✅ examples/basic/control_flow.php
- ✅ examples/basic/expressions.php
- ✅ examples/basic/functions.php
- ✅ examples/basic/hello.php
- ✅ examples/basic/strings.php
- ✅ examples/basic/variables.php

**Failing Files**:
- ❌ examples/parallel/ecommerce_parallel.php - Parser limitation: array destructuring with default values
- ❌ examples/parallel/parallel_examples.php - Parser limitation: array destructuring with default values

## Integration with CI

The test suite is designed to be CI-friendly:
- Fast execution for basic examples
- Configurable timeouts
- Detailed error reporting
- Summary statistics

To integrate with CI, add to your CI configuration:
```yaml
- name: Test Examples
  run: go test ./tests/ -v -timeout 1m
```

## Future Enhancements

- [ ] Add expected output files for comparison testing
- [ ] Implement output comparison logic
- [ ] Add more granular test categories
- [ ] Implement test discovery for new example files
- [ ] Add performance regression testing
