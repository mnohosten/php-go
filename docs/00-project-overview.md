# PHP-Go Project Overview

## Project Goal

Create a complete rewrite of the PHP 8.4 interpreter in Go with the following key features:

1. **Full PHP 8.4 Compatibility** - Support all language features and standard library functions
2. **Automatic Parallelization** - Leverage Go's concurrency to automatically parallelize safe operations
3. **Native Go Integration** - Allow seamless interoperability between PHP code and Go libraries
4. **Modern Architecture** - Clean, maintainable Go code with comprehensive testing

## Motivation

PHP is a powerful language but suffers from several limitations:

- **Single-threaded execution** - Cannot utilize modern multi-core processors effectively
- **Limited concurrency** - No built-in parallel execution primitives
- **Extension complexity** - Writing PHP extensions in C is difficult and error-prone
- **Performance limitations** - Even with JIT, PHP lags behind compiled languages

By rewriting PHP in Go, we can:

- Leverage Go's goroutines for automatic parallelization
- Use Go's garbage collector instead of manual reference counting
- Provide a clean Go API for writing extensions
- Enable direct integration with the vast Go ecosystem
- Improve performance through native compilation and efficient concurrency

## Target Specification

- **PHP Version**: 8.4 (latest development version from php-src)
- **Scope**: Core interpreter + standard library (~600K LOC equivalent)
- **Threading Model**: Automatic parallelization with concurrency safety analysis
- **Go Integration**: Both FFI-style bindings and native extension API

## Architecture Decisions

### 1. Virtual Machine vs Direct Execution

**Decision**: Implement a bytecode VM (similar to Zend Engine)

**Rationale**:
- Full compatibility with PHP semantics requires VM
- Enables dynamic features (eval, variable variables, etc.)
- Allows for optimization passes
- Future JIT compiler can optimize hot paths

### 2. Type System

**Decision**: Go structs for PHP types, interface-based polymorphism

**Rationale**:
- Type-safe Go implementation
- Leverage Go's interface system for PHP's dynamic typing
- Efficient memory layout
- Clear separation between type representation and operations

### 3. Memory Management

**Decision**: Use Go's garbage collector, eliminate reference counting

**Rationale**:
- Simplifies implementation significantly
- Eliminates circular reference issues
- Better performance for allocation-heavy workloads
- Trade-off: Less deterministic memory usage

### 4. Parallelization Strategy

**Decision**: Automatic parallelization with static analysis + runtime safety checks

**Approach**:
- Analyze code for shared state and side effects
- Request-level parallelism (like PHP-FPM, but in-process)
- Automatic parallelization of array operations (map, filter, reduce)
- Explicit parallel constructs for opt-in parallelism
- Copy-on-write semantics for safety

### 5. Go Integration

**Decision**: Dual approach - FFI + Native Extensions

**FFI System**:
- Call Go functions from PHP code
- Type marshaling between PHP and Go
- Dynamic loading of Go plugins

**Native Extensions**:
- Write PHP extensions entirely in Go
- Expose Go functions/classes to PHP
- Full access to runtime internals

## Project Scope

### Phase 1-4: Core Interpreter (~4-6 months)
- Lexer and parser
- Compiler (AST → opcodes)
- Virtual machine executor
- Type system and data structures
- **Milestone**: Execute basic PHP scripts

### Phase 5-6: Language Features (~3-4 months)
- Object system (classes, inheritance, traits)
- Standard library (arrays, strings, files, etc.)
- Common extensions (JSON, PCRE, date/time)
- **Milestone**: Run real-world PHP applications

### Phase 7-8: Parallelization & Integration (~3-4 months)
- Automatic parallelization engine
- Go integration (FFI + native extensions)
- Thread-safe runtime
- **Milestone**: Parallel PHP execution, Go library usage

### Phase 9-10: Advanced Features & Testing (~2-3 months)
- Generators, closures, reflection
- PHP test suite integration
- Performance optimization
- Documentation and migration guide
- **Milestone**: Production-ready release

### Total Estimated Timeline: 12-17 months

## Success Criteria

### Must Have (v1.0)
- [ ] Pass 95%+ of PHP 8.4 language tests
- [ ] Implement all standard library functions
- [ ] Support common extensions (JSON, PCRE, date, SPL, hash)
- [ ] Automatic request-level parallelization
- [ ] Basic Go integration (call Go from PHP)
- [ ] Performance within 2x of standard PHP

### Should Have (v1.x)
- [ ] Pass 99%+ of PHP test suite
- [ ] Implement database extensions (PDO, mysqli)
- [ ] XML extensions (dom, simplexml)
- [ ] Compression (zlib, zip)
- [ ] Advanced parallelization (array operations)
- [ ] Full native extension API
- [ ] Performance parity with PHP 8.4 + opcache

### Nice to Have (v2.0+)
- [ ] JIT compiler for hot paths
- [ ] Optimizer (SSA-based like Zend)
- [ ] Advanced Go integration features
- [ ] Performance exceeding PHP + JIT
- [ ] All PHP extensions implemented

## Technical Constraints

### Compatibility Requirements
- Must execute existing PHP 8.4 code without modifications
- Maintain identical semantics for all operations
- Support all PHP 8.4 syntax and features
- Compatible with PHP's type juggling and coercion rules

### Performance Requirements
- Startup time < 2x standard PHP
- Execution performance within 2x for v1.0
- Memory usage comparable to standard PHP
- Parallel execution should show linear speedup for independent requests

### Go Requirements
- Minimum Go 1.21 (for latest language features)
- Pure Go implementation (no CGO for core)
- Idiomatic Go code (follow Go conventions)
- Comprehensive testing (unit + integration)

## Non-Goals

The following are explicitly OUT of scope:

1. **PHP < 8.4 Compatibility** - Focus only on latest PHP version
2. **Windows Support (initially)** - Focus on Linux/macOS first
3. **JIT Compiler (v1.0)** - Defer to future versions
4. **All Extensions** - Only implement common/essential extensions
5. **ABI Compatibility** - Not compatible with existing PHP extensions
6. **Zend API Compatibility** - New Go-based extension API

## Repository Structure

```
php-go/
├── docs/                    # Comprehensive documentation
│   ├── 00-project-overview.md
│   ├── 01-php-analysis.md
│   ├── 02-go-architecture.md
│   ├── phases/              # Phase-specific documentation
│   │   ├── 01-foundation/
│   │   ├── 02-compiler/
│   │   └── ...
│   ├── api/                 # API documentation
│   ├── internals/           # Implementation details
│   ├── examples/            # Code examples
│   └── rfcs/                # Design proposals
│
├── cmd/
│   └── php-go/              # Main CLI executable
│
├── pkg/
│   ├── lexer/               # Tokenization
│   ├── parser/              # Parsing
│   ├── ast/                 # Abstract Syntax Tree
│   ├── compiler/            # Compilation (AST → opcodes)
│   ├── vm/                  # Virtual machine
│   ├── types/               # PHP type system
│   ├── runtime/             # Runtime support
│   ├── stdlib/              # Standard library
│   ├── parallel/            # Parallelization engine
│   └── goext/               # Go integration
│
├── internal/                # Internal packages
│   ├── opcodes/             # Opcode definitions
│   └── util/                # Utilities
│
├── tests/                   # Test suite
│   ├── unit/                # Unit tests
│   ├── integration/         # Integration tests
│   └── compatibility/       # PHP compatibility tests
│
├── benchmarks/              # Performance benchmarks
│
└── examples/                # Example PHP scripts
```

## Contributing

This is a large-scale project that will be developed in phases. Each phase has:

1. **Design Document** - Architecture and approach
2. **Task Breakdown** - Specific implementation tasks
3. **Test Plan** - How to verify correctness
4. **Progress Tracking** - Current status

See individual phase documents in `docs/phases/` for detailed information.

## Resources

- **PHP Source**: ../php-src (PHP 8.4 reference implementation)
- **PHP Internals Book**: https://www.phpinternalsbook.com/
- **Zend Engine Documentation**: PHP source code and internal docs
- **Go Documentation**: https://go.dev/doc/

## Next Steps

1. Review this overview and architecture decisions
2. Study `01-php-analysis.md` for PHP 8.4 implementation details
3. Review `02-go-architecture.md` for Go-specific design
4. Start with Phase 1: Foundation (lexer and parser)

---

**Last Updated**: 2025-11-21
**Status**: Planning Complete, Ready for Implementation
