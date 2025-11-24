# API Reference Documentation

This directory contains comprehensive API reference documentation for all PHP-Go packages.

## Available Documentation

### Core Packages (Implemented)

1. **[lexer.md](lexer.md)** - Lexer/Tokenizer API
   - Token types and scanning
   - PHP 8.4 syntax support
   - Position tracking for error reporting
   - ~6.2 KB

2. **[parser.md](parser.md)** - Parser API
   - AST generation from tokens
   - Pratt parsing for expressions
   - Error recovery mechanisms
   - ~6.2 KB

3. **[ast.md](ast.md)** - Abstract Syntax Tree API
   - 65+ node types
   - Statement and expression nodes
   - Declaration nodes
   - ~14 KB

4. **[compiler.md](compiler.md)** - Compiler API
   - AST to bytecode compilation
   - Symbol table management
   - Optimization passes
   - ~9.7 KB

5. **[vm.md](vm.md)** - Virtual Machine API
   - Bytecode execution engine
   - 210 opcodes
   - Frame management
   - ~13 KB

6. **[types.md](types.md)** - Type System API
   - Universal Value container
   - PHP arrays and objects
   - Type conversions and comparisons
   - Value pooling optimizations
   - ~12 KB

7. **[runtime.md](runtime.md)** - Runtime Support API
   - Superglobals ($_GET, $_POST, etc.)
   - Constants management
   - Error handling
   - Output buffering
   - ~11 KB

### Future Packages (Planned)

8. **[stdlib.md](stdlib.md)** - Standard Library (Phase 6)
   - ~300+ PHP built-in functions
   - Array, string, math, file I/O functions
   - JSON, PCRE, date/time support
   - ~8.9 KB

9. **[parallel.md](parallel.md)** - Parallelization (Phase 7)
   - Automatic parallelization
   - Worker pools
   - Copy-on-write data structures
   - Safety analysis
   - ~8.2 KB

10. **[goext.md](goext.md)** - Go Integration (Phase 8)
    - FFI (Foreign Function Interface)
    - Type marshaling
    - Extension API
    - Go standard library bindings
    - ~9.7 KB

## Documentation Structure

Each API reference document follows this structure:

1. **Overview** - Package purpose and key features
2. **Main Types** - Primary structs and interfaces with their methods
3. **Functions** - Package-level functions
4. **Usage Examples** - Practical code examples
5. **Implementation Notes** - Important details and considerations
6. **Performance** - Performance characteristics and optimizations
7. **Related Packages** - Cross-references to related packages

## Quick Navigation

### For PHP-Go Users

If you're **using** PHP-Go to run PHP code:

- Start with **[runtime.md](runtime.md)** - Understanding the runtime environment
- Read **[stdlib.md](stdlib.md)** - Available PHP functions (Phase 6+)
- Check **[goext.md](goext.md)** - How to integrate Go code (Phase 8+)

### For PHP-Go Contributors

If you're **developing** PHP-Go itself:

**Phase 1-3 (Lexer, Parser, Compiler, VM):**
- **[lexer.md](lexer.md)** → **[parser.md](parser.md)** → **[ast.md](ast.md)** → **[compiler.md](compiler.md)** → **[vm.md](vm.md)**

**Phase 4-5 (Data Structures, Objects):**
- **[types.md](types.md)** - Value system, arrays, objects

**Phase 6+ (Future Development):**
- **[stdlib.md](stdlib.md)** - Standard library implementation
- **[parallel.md](parallel.md)** - Parallelization features
- **[goext.md](goext.md)** - Go integration

## Code Examples

### Running PHP Code

```go
package main

import (
    "github.com/krizos/php-go/pkg/lexer"
    "github.com/krizos/php-go/pkg/parser"
    "github.com/krizos/php-go/pkg/compiler"
    "github.com/krizos/php-go/pkg/vm"
)

func main() {
    // Lex → Parse → Compile → Execute
    l := lexer.New("<?php echo 'Hello, World!';", "test.php")
    p := parser.New(l)
    program := p.ParseProgram()

    c := compiler.New()
    c.Compile(program)
    bytecode := c.Bytecode()

    machine := vm.New()
    machine.LoadConstants(bytecode.Constants)
    machine.Execute(bytecode.Instructions)

    println(machine.OutputString())
}
```

## Package Dependencies

```
lexer → parser → ast → compiler → vm
                         ↓         ↓
                       types ← runtime
                         ↓
                      stdlib (Phase 6)
                         ↓
                     parallel (Phase 7)
                         ↓
                      goext (Phase 8)
```

## Implementation Status

| Package | Status | Coverage | Size |
|---------|--------|----------|------|
| lexer | ✅ Complete | 82.8% | ~1,400 lines |
| parser | ✅ Complete | 85.0% | ~2,200 lines |
| ast | ✅ Complete | N/A | ~1,150 lines |
| compiler | ✅ Complete | 85.1% | ~1,200 lines |
| vm | ✅ Complete | 79.1% | ~3,500 lines |
| types | ✅ Complete | 89.2% | ~2,000 lines |
| runtime | ✅ Complete | 99.2% | ~800 lines |
| stdlib | 🚧 Phase 6 | N/A | Planned |
| parallel | 🚧 Phase 7 | N/A | Planned |
| goext | 🚧 Phase 8 | N/A | Planned |

## Additional Resources

- **Project Overview**: `../00-project-overview.md`
- **PHP Analysis**: `../01-php-analysis.md`
- **Go Architecture**: `../02-go-architecture.md`
- **Phase Documentation**: `../phases/*/README.md`
- **User Guide**: `../user-guide/README.md`

## Generating API Documentation

To generate HTML documentation from Go code:

```bash
# Install godoc
go install golang.org/x/tools/cmd/godoc@latest

# Start documentation server
godoc -http=:6060

# Open browser to http://localhost:6060/pkg/github.com/krizos/php-go/
```

Or use pkgsite:

```bash
# Install pkgsite
go install golang.org/x/pkgsite/cmd/pkgsite@latest

# Start documentation server
pkgsite -http=:8080

# Open browser to http://localhost:8080/github.com/krizos/php-go
```

## Contributing to Documentation

When adding or updating packages:

1. Update the corresponding API reference document
2. Follow the established structure
3. Include practical examples
4. Document performance characteristics
5. Add cross-references to related packages
6. Update this README with the new package

## Documentation Standards

- **Clarity**: Write for both users and developers
- **Examples**: Include working code examples
- **Completeness**: Document all exported types and functions
- **Accuracy**: Keep documentation in sync with code
- **Format**: Use Markdown with code blocks and tables

## Version Information

- **Documentation Version**: 1.0
- **PHP-Go Version**: 0.5.3 (Phase 5 complete)
- **Last Updated**: 2025-11-24
- **Total Documentation**: ~100 KB across 10 files

## Support

For questions about the API:
- Check the relevant API reference document
- Review code examples in `/examples/`
- See phase documentation in `/docs/phases/`
- Consult the main project README

---

**Note**: This documentation covers the public API only. For internal implementation details, see the Go source code in `/pkg/` with inline comments.
