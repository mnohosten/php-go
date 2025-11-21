# API Reference

This directory contains API reference documentation for PHP-Go packages.

## Organization

Documentation is organized by package:

- **lexer/** - Lexer/tokenizer API
- **parser/** - Parser API
- **ast/** - Abstract Syntax Tree node types
- **compiler/** - Compiler API
- **vm/** - Virtual Machine API
- **types/** - Type system API
- **runtime/** - Runtime support API
- **stdlib/** - Standard library functions
- **parallel/** - Parallelization API
- **goext/** - Go integration API

## Documentation Format

Each package includes:

1. **Overview** - Package purpose and high-level design
2. **Types** - Structs, interfaces, and type definitions
3. **Functions** - Public functions and methods
4. **Examples** - Usage examples
5. **Notes** - Important considerations

## Generation

API documentation is generated from Go source code using:

```bash
# Generate godoc
godoc -http=:6060

# Or use pkgsite
pkgsite -http=:8080
```

## For Users

If you're using PHP-Go, you primarily need:

- **stdlib/** - PHP functions available in PHP-Go
- **goext/** - How to integrate Go code with PHP

## For Developers

If you're contributing to PHP-Go, review:

- All package documentation
- Internal implementation details
- Design decisions and trade-offs

---

**Last Updated**: 2025-11-21
