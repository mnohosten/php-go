# Internals Documentation

This directory contains detailed documentation about PHP-Go's internal implementation.

## Contents

### Core Systems

- **lexer-internals.md** - How the lexer works internally
- **parser-internals.md** - Parser implementation details
- **compiler-internals.md** - Compilation pipeline internals
- **vm-internals.md** - Virtual machine execution details
- **memory-management.md** - Memory and GC integration
- **type-system.md** - Type representation and conversions

### Advanced Topics

- **opcode-dispatch.md** - VM dispatch mechanisms
- **optimization.md** - Optimization techniques
- **parallelization.md** - How automatic parallelization works
- **go-integration.md** - FFI and extension internals
- **error-handling.md** - Error and exception handling internals

### Design Documents

- **design-decisions.md** - Major design decisions and rationale
- **tradeoffs.md** - Performance vs compatibility trade-offs
- **php-compatibility.md** - How we maintain PHP compatibility
- **future-work.md** - Future optimization opportunities

## Purpose

These documents explain:

1. **Why** - Design rationale
2. **How** - Implementation details
3. **Trade-offs** - Decisions and alternatives
4. **Future** - Improvement opportunities

## Audience

- Contributors to PHP-Go
- Developers interested in language implementation
- Those debugging complex issues
- Researchers and educators

## Contributing

When making significant changes:

1. Update relevant internals documentation
2. Explain design decisions
3. Document trade-offs
4. Add examples where helpful

---

**Last Updated**: 2025-11-21
