# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PHP-Go is a complete rewrite of the PHP 8.4 interpreter in Go, featuring automatic parallelization and native Go library integration. This is a multi-phase implementation project currently at Phase 10 (Testing & Production Readiness) in progress (93.6% complete, 1340/1430 hours).

**Current Status**: Phases 0-9 complete (Foundation through Advanced Features). Phase 10 in progress (Security Audit complete, Stress Testing complete, Framework Testing complete, Critical Parser Features 84% complete). Only 1 major blocker remaining: Closures/Anonymous Functions (15-25h).

## Essential Commands

### Building and Testing
```bash
# Build the CLI
go build -o php-go ./cmd/php-go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run single package tests
go test ./pkg/lexer/
go test ./pkg/parser/
go test ./pkg/compiler/

# Run benchmarks
go test -bench=. -benchmem ./pkg/lexer/
go test -bench=. ./pkg/parser/
```

### Development Tools
```bash
# Tokenize PHP file (show lexer output)
./php-go lex test.php
./php-go lex --json test.php

# Parse PHP file (show AST)
./php-go parse test.php
./php-go parse --json test.php

# Format code
gofmt -w .
```

### Coverage Analysis
```bash
# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Check specific package coverage
go test -cover ./pkg/types/
```

### Example Testing
```bash
# Run all example files and compare with expected output
go test ./tests/

# Run examples with PHP 8.4 validation
./test_with_php.sh

# Run examples with php-go and show detailed output
./test_all_examples.sh

# Run examples with php-go and show summary only
./test_summary.sh
```

## Architecture Overview

### Package Structure

**Phase 1 (Complete)**: Foundation - Lexer, Parser, AST
- `pkg/lexer/` - Tokenization (~780 lines core lexer, ~630 lines token definitions)
- `pkg/parser/` - Parsing with Pratt precedence climbing (~2200 lines)
- `pkg/ast/` - Abstract Syntax Tree (65+ node types, ~800 lines)

**Phase 2 (Complete)**: Compiler - AST to Opcodes
- `pkg/compiler/` - Compiles AST to PHP VM opcodes (~1200 lines)
- `pkg/vm/opcodes.go` - 210 PHP VM opcode definitions
- `pkg/vm/instruction.go` - 24-byte instruction encoding

**Phase 3 (Complete)**: Runtime & Virtual Machine
- `pkg/vm/` - Bytecode execution engine with opcode handlers
- `pkg/types/value.go` - PHP value type system (10 types: Undef, Null, Bool, Int, Float, String, Array, Object, Resource, Reference)
- `pkg/runtime/` - Runtime support (globals, superglobals, constants, errors, output buffering)

**Phase 4 (Partial)**: Data Structures
- `pkg/types/array.go` - PHP associative arrays (order-preserving)
- `pkg/types/string.go` - Binary-safe strings

**Phase 5 (Complete)**: Object System
- `pkg/types/object.go` - Full OOP implementation (437 lines)
  - Classes with inheritance, interfaces, traits, enums
  - Properties (visibility, readonly, hooks)
  - Methods (visibility, static, abstract, final)
  - Magic methods, reflection, late static binding

**Phase 6 (In Progress)**: Example Compatibility & Bug Fixes
- `pkg/compiler/regression_test.go` - Comprehensive regression tests (56 test cases)
- `tests/examples_test.go` - Example file integration tests
- Phase 6A-6C complete: Critical bug fixes, PHP 7+ language features, testing & validation

**Phase 7-10 (Pending)**: Standard Library, Parallelization, Go Integration, Advanced Features

### Key Design Patterns

**Lexer**: Character stream → Token stream
- Handles PHP tags (`<?php`, `?>`), strings (single/double-quoted, heredoc/nowdoc), interpolation
- Position tracking for error reporting
- Keywords vs identifiers differentiation

**Parser**: Token stream → AST
- Pratt parsing for expressions (operator precedence)
- Recursive descent for statements and declarations
- Error recovery mechanisms for better diagnostics

**Compiler**: AST → Bytecode
- Symbol tables with nested scopes (global, local, free variables)
- Jump patching for control flow (if/else, loops, try/catch)
- Constant folding and dead code elimination
- Loop context stack for break/continue

**VM**: Bytecode executor
- Stack-based execution with call frames
- Opcode dispatch via switch statement
- Value type system with PHP-compatible type juggling
- Frame management with $this binding and class context

**Type System**:
- `Value` struct wraps all PHP types (similar to PHP's zval)
- Type conversions follow PHP semantics (type juggling)
- Reference counting handled by Go's GC
- Arrays are map+slice for order preservation

## Critical Implementation Details

### Type Juggling and Conversions
PHP has complex type coercion rules. The `pkg/types/value.go` file implements:
- `Equals()` for loose comparison (==) with type juggling
- `Identical()` for strict comparison (===) without coercion
- `ToInt()`, `ToFloat()`, `ToBool()`, `ToString()`, `ToArray()` follow PHP rules
- Truthiness: empty string `""` and `"0"` are false, NaN is false

### Array Implementation
PHP arrays are ordered associative maps:
- Keys can be integers or strings
- Insertion order is preserved
- Implementation uses map + slice for order tracking
- Numeric string keys are not converted to integers (PHP quirk)

### Object System (Phase 5)
Complete PHP 8.4 OOP with:
- **Inheritance**: `InheritFrom()` copies properties/methods from parent
- **Interfaces**: `ValidateInterfaceImplementation()` checks compliance
- **Traits**: `ApplyTraits()` handles composition and conflict resolution
- **Enums**: Pure and backed enums (PHP 8.1+)
- **Magic Methods**: All 14 magic methods supported (__get, __set, __call, etc.)
- **Reflection**: Full metadata access for classes, methods, properties
- **Late Static Binding**: static:: vs self:: vs parent::
- **Class Name Resolution**: `ClassName::class` returns fully qualified name

### Phase 6 Enhancements (Compatibility & Bug Fixes)

**Phase 6A - Critical Bug Fixes**:
- **Variable Naming**: Variables can now use builtin function names (`$count`, `$empty`, etc.) without conflicts
- **Foreach Iterator Management**: Fixed temp variable allocation to prevent iterator corruption during loop body execution
- **DECLARE_CLASS Handler**: Implemented opcode handler for runtime class registration with inheritance support
- **If Statement Conditions**: Fixed condition evaluation bug where hardcoded temps caused incorrect branching
- **Increment/Decrement**: Verified all pre/post increment/decrement operators work in all contexts

**Phase 6B - PHP 7+ Language Features**:
- **Null Coalescing Operator** (`??`): Full support for `$x ?? "default"` with proper null/undef checking
- **Do-While Statements**: Complete implementation with correct body-first execution order
- **Array Destructuring in Foreach**: Support for `foreach ($arr as [$a, $b])` and `foreach ($arr as ['key' => $v])`

**Phase 6C - Testing & Validation**:
- **Example Test Suite**: 7/7 basic examples passing (100%), automated testing in `tests/examples_test.go`
- **Regression Tests**: 56 comprehensive test cases in `pkg/compiler/regression_test.go` covering all Phase 6A-6B fixes
- **Output Validation**: All basic examples validated against expected output files

### Opcode Handlers
Each opcode has a handler in `pkg/vm/handlers_*.go`:
- `handlers_arithmetic.go` - Math operations with type juggling
- `handlers_comparison.go` - Comparisons (==, ===, <, >, <=>), null coalescing (??)
- `handlers_logic.go` - Boolean and bitwise operations
- `handlers_variables.go` - Variable operations (fetch, assign)
- `handlers_control.go` - Jumps (JMP, JMPZ, JMPNZ)
- `handlers_functions.go` - Function calls (INIT_FCALL, SEND_VAL, DO_FCALL)
- `handlers_object.go` - Object operations (NEW, FETCH_OBJ, ASSIGN_OBJ, DECLARE_CLASS, method calls)
- `handlers_array.go` - Array operations (INIT_ARRAY, FETCH_DIM, ASSIGN_DIM, foreach iterators)
- `handlers_strings.go` - String concatenation
- `handlers_io.go` - Output operations (ECHO)

### Symbol Tables and Scopes
The compiler maintains symbol tables for variable tracking:
- **Global scope**: Top-level variables
- **Local scope**: Function/method parameters and locals
- **Free variables**: Captured by closures (for future implementation)
- Built-in functions pre-registered: echo, print, var_dump, isset, empty, count, strlen

### Testing Standards
All phases target 85%+ code coverage:
- Phase 1 (Lexer/Parser): 82.8% lexer, 85.0% parser
- Phase 2 (Compiler): 85.1%
- Phase 3 (Runtime/VM): ~89% average (types: 89.2%, runtime: 99.2%, vm: 79.1%)
- Phase 5 (Objects): 78.2%
- Phase 6 (Compatibility): 56 regression tests, 7/7 basic examples passing (100%)

## Common Development Workflows

### Adding a New Opcode Handler
1. Add opcode constant to `pkg/vm/opcodes.go`
2. Add case to dispatch switch in `pkg/vm/vm.go` Execute()
3. Implement handler function in appropriate `pkg/vm/handlers_*.go`
4. Add compiler support in `pkg/compiler/compiler.go`
5. Write tests in `pkg/vm/vm_test.go` or specific handler test file

### Adding a New Statement Type
1. Define AST node in `pkg/ast/ast.go`
2. Add parsing logic to `pkg/parser/stmt.go`
3. Add compilation logic to `pkg/compiler/compiler.go`
4. Write parser tests in `pkg/parser/stmt_test.go`
5. Write compiler tests in `pkg/compiler/compiler_test.go`

### Adding a New Expression Type
1. Define AST node in `pkg/ast/ast.go`
2. Add prefix/infix parser in `pkg/parser/expr.go`
3. Register in prefix/infix function maps
4. Add compilation logic to `pkg/compiler/compiler.go`
5. Write tests for parser and compiler

### Debugging Compilation
```bash
# Use lex command to check tokenization
./php-go lex test.php

# Use parse command to check AST
./php-go parse test.php

# Add debug output in compiler
# See pkg/compiler/compiler.go Compile() for instruction dump
```

## Task Tracking

**Master Reference**: `TODO.md` contains all tasks for Phases 0-10 with:
- Checkboxes for completion tracking
- Effort estimates (hours)
- File references
- Links to detailed phase documentation in `docs/phases/`

**Daily Workflow**:
1. Check `TODO.md` for next unchecked task
2. Read referenced phase doc in `docs/phases/*/README.md`
3. Implement the task
4. Write tests (maintain 85%+ coverage)
5. Mark task complete in `TODO.md`
6. Commit with format: `feat(phaseN): <description>`

## PHP Compatibility Notes

### PHP Quirks to Preserve
- Empty string `""` and `"0"` are falsy in boolean context
- Numeric strings in comparisons are converted to numbers
- Array keys: numeric strings stay as strings (not converted)
- Variable-variables: `$$var` requires runtime evaluation
- Magic quotes handling (deprecated but may need support)

### PHP 8.4 Features Implemented
- Union types (int|string)
- Nullable types (?int)
- Null coalescing operator (??) - Phase 6B ✅
- Do-while statements - Phase 6B ✅
- Class name resolution (ClassName::class) - Phase 6B ✅
- Array destructuring in foreach - Phase 6B ✅
- Array append syntax ($arr[] = value) - Phase 10 ✅
- Single-line control structures (if ($x) return 1;) - Phase 10 ✅
- Readonly properties (PHP 8.1+) ✅
- Readonly classes (PHP 8.2+) ✅
- Enums (PHP 8.1+) ✅
- Named arguments (PHP 8.0+) - Phase 10 ✅
- First-class callables (PHP 8.1) - Phase 10 ✅
- Match expressions (PHP 8.0+) - Phase 10 ✅
- Throw expressions (PHP 8.0) - Phase 10 ✅
- Array spread operator (PHP 7.4+) - Phase 10 ✅
- Generators/yield (PHP 5.5+) - Phase 10 ✅
- Alternative control structures (:endif;) - Phase 10 ✅
- array() constructor - Phase 10 ✅
- Constructor property promotion (PHP 8.0) - Phase 10 ✅
- Standard library (50+ functions) - Phase 6D/10 ✅

### Not Yet Implemented
- Closures with variable capture (15-25h remaining - LAST MAJOR BLOCKER)
- Arrow functions (depends on closures)
- Attributes (deferred)
- Additional standard library functions (~300+ remaining)
- Parallelization features (architecture ready, Phase 7)
- Go integration features (architecture ready, Phase 8)

## Performance Considerations

### Actual Performance Results (Phase 10)
- **Overall Performance**: 5.0x faster than PHP 8.4 + opcache (target exceeded!)
- **Benchmark Success**: 91% (10/11 benchmarks passing)
- **Simple Loop**: 9.0x faster than PHP 8.4
- **Function Calls**: 10.4x faster than PHP 8.4
- **Recursion**: Fibonacci(15) in 13.19μs with only 35.5KB memory (3552x faster)
- **String Operations**: 181x faster than PHP 8.4
- **Array Operations**: 5.6x faster than PHP 8.4
- **OOP Patterns**: 4.6-5.2x faster than PHP 8.4
- **Stability**: All benchmarks < 5% coefficient of variation

### Optimization Strategies (Implemented)
- Constant folding at compile time ✅
- Dead code elimination ✅
- Proper temp variable management ✅
- Jump target optimization ✅
- Security features with < 1% overhead ✅

## Error Messages and Debugging

### Position Tracking
All tokens carry position information (file, line, column, offset) for accurate error messages.

### Error Format
```
Parse error: Unexpected token 'EOF', expected ';'
  at /path/to/file.php:15:8
```

### Adding New Errors
Use `parser.error()` or `compiler.errorf()` methods with position info.

## Contributing Guidelines

### Code Style
- Follow Go conventions (gofmt, goimports)
- Document all exported items
- Keep functions focused and small
- Use descriptive names

### Commit Messages
Format: `<type>(phase<N>): <description>`
- `feat`: New feature
- `fix`: Bug fix
- `test`: Add/update tests
- `docs`: Documentation
- `refactor`: Code refactoring
- `perf`: Performance improvement

### Testing Requirements
- All new code must have tests
- Maintain 85%+ coverage for new packages
- Integration tests for complex features
- Benchmarks for performance-critical code

## Reference Documentation

**Internal Docs**:
- `docs/00-project-overview.md` - High-level architecture
- `docs/01-php-analysis.md` - PHP 8.4 internals analysis
- `docs/02-go-architecture.md` - Go implementation design
- `docs/phases/*/README.md` - Detailed phase implementation guides

**External References**:
- PHP source: `php-src/` directory (reference implementation)
- [PHP Language Spec](https://github.com/php/php-langspec)
- [PHP Internals Book](http://www.phpinternalsbook.com/)
- PHP VM opcodes: `php-src/Zend/zend_vm_opcodes.h`
- PHP compilation: `php-src/Zend/zend_compile.c`

## Current Implementation Status

**Completed (Phases 0-5, 552 hours)**:
- Lexer with full PHP tokenization
- Parser with complete PHP 8.4 syntax support
- Compiler with opcode generation and optimization
- Virtual machine with opcode execution
- Type system with PHP-compatible conversions
- Runtime support (globals, errors, output buffering)
- Complete object system (classes, interfaces, traits, enums)

**Completed (Phase 6A-6C, 52 hours)**:
- Critical bug fixes: variable naming, foreach iterators, if conditions, class declarations
- PHP 7+ language features: null coalescing (??), do-while, array destructuring
- Comprehensive testing: 56 regression tests, 7/7 basic examples passing

**Next Tasks (Phase 6D starting)**:
- Standard library implementation (~300+ functions)
- Array manipulation functions (array_push, array_pop, array_map, array_filter, array_merge)
- String processing functions (substr, str_replace, explode, implode)
- Type checking functions (is_null, is_array, is_string, is_int, is_bool, gettype)
- Utility functions (print_r, microtime, date)
- Math functions (abs, round, floor, ceil, min, max)

**Project Timeline**: 12-17 months to v1.0 (currently 58% through planned hours, 604/1050)
