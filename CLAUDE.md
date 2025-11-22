# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PHP-Go is a complete rewrite of the PHP 8.4 interpreter in Go, featuring automatic parallelization and native Go library integration. This is a multi-phase implementation project currently at Phase 5 completion (53% complete, 552/1050 hours).

**Current Status**: Phases 0-5 complete (Foundation, Compiler, Runtime/VM, Data Structures, Object System). Phase 6+ pending.

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

**Phase 6-10 (Pending)**: Standard Library, Parallelization, Go Integration, Advanced Features, Testing

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

### Opcode Handlers
Each opcode has a handler in `pkg/vm/handlers_*.go`:
- `handlers_arithmetic.go` - Math operations with type juggling
- `handlers_comparison.go` - Comparisons (==, ===, <, >, <=>)
- `handlers_logic.go` - Boolean and bitwise operations
- `handlers_variables.go` - Variable operations (fetch, assign)
- `handlers_control.go` - Jumps (JMP, JMPZ, JMPNZ)
- `handlers_functions.go` - Function calls (INIT_FCALL, SEND_VAL, DO_FCALL)
- `handlers_object.go` - Object operations (NEW, FETCH_OBJ, ASSIGN_OBJ, method calls)
- `handlers_array.go` - Array operations (INIT_ARRAY, FETCH_DIM, ASSIGN_DIM)
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
- Named arguments (pending Phase 9)
- Attributes (pending Phase 9)
- Readonly properties (PHP 8.1+)
- Readonly classes (PHP 8.2+)
- Enums (PHP 8.1+)
- First-class callables (pending Phase 9)

### Not Yet Implemented
- Generators (Phase 9)
- Closures with variable capture (Phase 9)
- Arrow functions (Phase 9)
- Full reflection API (partial in Phase 5, complete in Phase 9)
- Standard library functions (Phase 6)
- Parallelization (Phase 7)
- Go integration (Phase 8)

## Performance Considerations

### Benchmarking Baselines (Phase 1)
- Lexer: ~1-18μs per operation (simple to complex)
- Parser: ~10-92μs per operation
- Target: Competitive with PHP 8.4 + opcache

### Optimization Strategies
- Constant folding at compile time (implemented)
- Dead code elimination (implemented)
- Packed arrays for sequential integer keys (Phase 4 pending)
- Copy-on-write for arrays and strings (Phase 7)
- Strength reduction (deferred)

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

**Next Tasks (Phase 6 starting)**:
- Standard library implementation (~300+ functions)
- Array manipulation functions
- String processing functions
- File I/O operations
- JSON extension
- PCRE (regex) support

**Project Timeline**: 12-17 months to v1.0 (currently 53% through planned hours)
