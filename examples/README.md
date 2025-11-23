# PHP-Go Examples

This directory contains example code demonstrating PHP-Go features and functionality.

## Directory Structure

### basic/
Working examples of basic PHP language features:
- `hello.php` - Hello world with for loop
- `variables.php` - Variable types and assignment
- `expressions.php` - Arithmetic and logical expressions
- `strings.php` - String operations and strlen()
- `control_flow.php` - If/else, for, while loops
- `functions.php` - Built-in functions (strlen, exit, die)
- `arrays.php` - Array operations (partial, in development)

### oop/
Object-oriented programming examples (Phase 5 complete):
- Coming soon: classes, inheritance, interfaces, traits, enums

### parallel/
Parallelization examples (Phase 7 complete, stdlib integration pending):
- `parallel_examples.php` - Comprehensive parallelization demos
- `ecommerce_parallel.php` - Real-world e-commerce example
- `README.md` - Status and notes about future functionality

### advanced/
Advanced language features (Phase 9):
- Coming soon: generators, closures, arrow functions, attributes

## Current Implementation Status

✅ **Working Now:**
- Basic variables and expressions (Phase 1-3)
- Control flow (if/else, for, while)
- Built-in functions: strlen(), exit(), die()
- Object system (classes, inheritance, traits - Phase 5)

🚧 **In Progress:**
- Standard library functions (Phase 6)
- String concatenation operator
- Full array support

⏳ **Future:**
- Parallelization integration (Phase 7 + Phase 6)
- Go integration (Phase 8)
- Generators and closures (Phase 9)

## Running Examples

```bash
# Build php-go first
go build -o php-go ./cmd/php-go

# Run an example
./php-go examples/basic/hello.php
./php-go examples/basic/variables.php

# See output
./php-go examples/basic/strings.php
```

## Learning Path

1. **Start here**: `basic/hello.php` - Simple hello world
2. **Then try**: `basic/variables.php` - Learn about variables
3. **Next**: `basic/expressions.php` - Arithmetic and logic
4. **Continue**: `basic/control_flow.php` - If/else and loops
5. **Explore**: Other basic/ examples

## For More Information

- See `docs/examples/README.md` for planned example categories
- See `docs/phases/` for implementation phase details
- See project README for overall architecture

## Contributing

When adding examples:
1. Place in appropriate subdirectory (basic, oop, advanced, parallel)
2. Add comments explaining what it demonstrates
3. Only include features that currently work (unless marked as future)
4. Update this README with new files

---

**Last Updated**: 2025-11-23
