# PHP-Go: PHP 8.4 Interpreter in Go with Automatic Parallelization

A complete rewrite of the PHP 8.4 interpreter in Go, featuring automatic parallelization and native Go library integration.

## Features

- **Full PHP 8.4 Compatibility** - Execute existing PHP code without modifications
- **Automatic Parallelization** - Leverage Go's goroutines for concurrent execution
- **Multi-threaded** - Unlike traditional PHP, PHP-Go can utilize all CPU cores
- **Native Go Integration** - Call Go functions and use Go libraries from PHP
- **Modern Architecture** - Clean Go codebase with comprehensive testing
- **Performance** - Competitive performance with PHP 8.4 + opcache

## Status

**Project Status**: Planning & Initial Development

This is an ambitious project to completely rewrite PHP in Go. See [docs/00-project-overview.md](docs/00-project-overview.md) for detailed information.

### Current Phase

**Phase 0**: Documentation and Planning (In Progress)

- [x] Project overview and architecture
- [x] PHP 8.4 source code analysis
- [x] Go architecture design
- [x] Phase 1-10 implementation plans
- [ ] Phase 1: Foundation (Lexer, Parser, AST) - Next

See [docs/phases/](docs/phases/) for detailed phase documentation.

## Quick Start

```bash
# Install (not yet available)
go install github.com/krizos/php-go/cmd/php-go@latest

# Run PHP script
php-go script.php

# Interactive mode
php-go -a

# Built-in web server
php-go -S localhost:8000
```

## Project Goals

1. **PHP 8.4 Compatibility** - Support all language features and standard library
2. **Automatic Parallelization** - Transparently parallelize safe operations
3. **Go Integration** - Seamless interop with Go ecosystem
4. **Production Ready** - Stable, tested, and documented

## Architecture

PHP-Go consists of several major components:

1. **Lexer & Parser** - Tokenization and parsing (Phase 1)
2. **Compiler** - AST to bytecode compilation (Phase 2)
3. **Virtual Machine** - Bytecode execution engine (Phase 3)
4. **Type System** - PHP value types in Go (Phase 3-4)
5. **Standard Library** - ~300+ built-in functions (Phase 6)
6. **Parallelization Engine** - Automatic concurrency (Phase 7)
7. **Go Integration** - FFI and native extensions (Phase 8)

See [docs/02-go-architecture.md](docs/02-go-architecture.md) for detailed architecture.

## Documentation

### For Users

- [Installation Guide](docs/user-guide/installation.md) (coming soon)
- [Getting Started](docs/user-guide/getting-started.md) (coming soon)
- [Configuration](docs/user-guide/configuration.md) (coming soon)
- [Migration from PHP](docs/migration-guide/from-php.md) (coming soon)

### For Developers

- [Project Overview](docs/00-project-overview.md) ✓
- [PHP 8.4 Analysis](docs/01-php-analysis.md) ✓
- [Go Architecture](docs/02-go-architecture.md) ✓
- [Implementation Phases](docs/phases/) ✓
- [Contributing Guide](CONTRIBUTING.md) (coming soon)

### API Reference

- [Standard Library](docs/api/stdlib/) (coming soon)
- [Go Integration](docs/api/goext/) (coming soon)
- [Extension Development](docs/extension-guide/) (coming soon)

## Examples

```php
<?php
// Traditional PHP code works as-is
$numbers = range(1, 1000);
$sum = array_sum($numbers);
echo "Sum: $sum\n";

// Automatic parallelization for independent requests
// (Runs in parallel goroutines automatically)

// Explicit parallelism with new APIs
$futures = go_parallel([
    fn() => heavyComputation1(),
    fn() => heavyComputation2(),
    fn() => heavyComputation3(),
]);
$results = go_wait($futures);

// Call Go functions directly
$response = go_call('http.Get', 'https://api.example.com');
$hash = go_call('crypto.SHA256', $data);
```

See [docs/examples/](docs/examples/) for more examples.

## Performance

**Target Performance (v1.0)**:
- Simple scripts: 0.5-2x PHP 8.4
- Web requests: 1-2x PHP 8.4
- Parallel requests: 2-4x PHP 8.4 (multi-core benefit)

**Future Performance (v2.0+ with JIT)**:
- All scenarios: 0.5-1x PHP 8.4 (potentially faster)

## Scope

### Phase 1-6: Core Implementation (~30 weeks)
- Complete interpreter
- Standard library
- Object system
- ~600K LOC equivalent

### Phase 7-9: Advanced Features (~16 weeks)
- Parallelization
- Go integration
- Generators, closures, reflection

### Phase 10: Production Ready (Ongoing)
- Testing & compatibility
- Performance optimization
- Documentation

**Total Timeline**: ~12-17 months for v1.0

## Technology Stack

- **Language**: Go 1.21+
- **Parser**: Hand-written recursive descent
- **VM**: Bytecode interpreter (210 opcodes)
- **GC**: Go's built-in garbage collector
- **Concurrency**: Goroutines and channels
- **Testing**: Go's testing framework + PHP test suite

## Comparison with PHP

| Feature | PHP 8.4 | PHP-Go |
|---------|---------|--------|
| Threading | Single-threaded | Multi-threaded |
| Parallelization | None | Automatic |
| Extensions | C-based | Go-based |
| Memory | Reference counting + GC | Go GC |
| Performance | Baseline | Competitive |
| Deployment | Interpreter + extensions | Single binary |
| Go Integration | Via FFI/C | Native |

## Contributing

We welcome contributions! This is a large project with many opportunities to help.

### Areas to Contribute

- **Core Implementation** - Lexer, parser, compiler, VM
- **Standard Library** - Implement PHP functions in Go
- **Extensions** - Port PHP extensions to Go
- **Testing** - Write tests, run PHP test suite
- **Documentation** - Improve docs and examples
- **Performance** - Optimize hot paths

See [CONTRIBUTING.md](CONTRIBUTING.md) (coming soon) for guidelines.

## Roadmap

### Current Milestone: Phase 1 (Foundation)
- [ ] Lexer implementation
- [ ] Parser implementation
- [ ] AST definition
- [ ] Basic CLI tool

### Upcoming Milestones
- Phase 2: Compiler (AST → Opcodes)
- Phase 3: Runtime & VM
- Phase 4: Data Structures
- Phase 5: Object System
- Phase 6: Standard Library

See [docs/phases/](docs/phases/) for detailed plans.

## Project Structure

```
php-go/
├── cmd/
│   └── php-go/          # Main CLI executable
├── pkg/
│   ├── lexer/           # Tokenization
│   ├── parser/          # Parsing
│   ├── ast/             # Abstract Syntax Tree
│   ├── compiler/        # Compilation
│   ├── vm/              # Virtual machine
│   ├── types/           # Type system
│   ├── runtime/         # Runtime support
│   ├── stdlib/          # Standard library
│   ├── parallel/        # Parallelization
│   └── goext/           # Go integration
├── docs/                # Documentation
├── tests/               # Test suite
├── benchmarks/          # Performance benchmarks
└── examples/            # Example code
```

## References

- **PHP Source**: [php/php-src](https://github.com/php/php-src) (Reference implementation)
- **PHP Language Spec**: [php/php-langspec](https://github.com/php/php-langspec)
- **PHP Internals Book**: [phpinternalsbook.com](http://www.phpinternalsbook.com/)
- **Go**: [go.dev](https://go.dev/)

## License

[To be determined - likely MIT or Apache 2.0]

## FAQ

### Why rewrite PHP in Go?

1. **Multi-threading** - PHP is single-threaded; Go enables true parallelism
2. **Modern language** - Go's simplicity and performance
3. **Easy deployment** - Single binary, no extensions to compile
4. **Go ecosystem** - Access Go libraries from PHP
5. **Learning** - Deep dive into language implementation

### Will it be faster than PHP?

Initial focus is on correctness and compatibility. Performance will be competitive with PHP 8.4. With future optimizations (JIT, etc.), it could potentially be faster.

### Is it compatible with existing PHP code?

Yes! The goal is 100% PHP 8.4 compatibility. Existing PHP applications should run without modifications.

### Can I use PHP extensions?

Existing C-based PHP extensions won't work. However:
- Core extensions are built-in
- Common extensions will be ported
- New extensions can be written in Go (much easier!)

### When will it be ready?

- **Alpha (v0.1)**: ~6 months - Basic execution
- **Beta (v0.5)**: ~12 months - Most features
- **Production (v1.0)**: ~17 months - Full compatibility

### How can I help?

Even if you're not ready to contribute code, you can:
- Star the repository
- Provide feedback on design
- Test alpha/beta releases
- Write documentation
- Spread the word

---

**Project Started**: 2025-11-21
**Current Status**: Planning & Documentation
**Maintainer**: @krizos

**Note**: This is a personal/educational project currently in the planning phase. Contributions and feedback welcome!
