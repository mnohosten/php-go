# Examples

This directory contains example code demonstrating PHP-Go features.

## Categories

### Basic Examples

- **hello-world.php** - Simplest PHP-Go program
- **variables.php** - Variable usage
- **functions.php** - Function definitions and calls
- **arrays.php** - Array operations
- **classes.php** - Object-oriented programming

### Advanced Examples

- **generators.php** - Generator and yield usage
- **closures.php** - Anonymous functions and closures
- **exceptions.php** - Exception handling
- **reflection.php** - Runtime introspection
- **attributes.php** - PHP 8.0+ attributes

### Parallelization Examples

- **parallel-requests.php** - Concurrent request handling
- **parallel-array-ops.php** - Automatic array parallelization
- **explicit-parallel.php** - Using go_* functions for explicit parallelism
- **channels.php** - Channel-based communication

### Go Integration Examples

- **calling-go.php** - Call Go functions from PHP
- **go-http-client.php** - Use Go's HTTP client
- **go-crypto.php** - Use Go's crypto libraries
- **custom-extension/** - Example of writing a PHP extension in Go

### Real-World Examples

- **web-server/** - Simple web server
- **rest-api/** - REST API implementation
- **database-app/** - Database interaction
- **worker-queue/** - Background job processing
- **microservice/** - Microservice example

## Running Examples

```bash
# Run an example
php-go examples/hello-world.php

# Run with debug mode
php-go --debug examples/generators.php

# Run with parallelization enabled
php-go --parallel examples/parallel-requests.php
```

## Example Structure

Each example includes:

1. **Code** - The PHP (and Go if applicable) code
2. **README** - Explanation of what it demonstrates
3. **Output** - Expected output
4. **Notes** - Important points or gotchas

## Contributing Examples

When adding examples:

1. Keep them focused on one concept
2. Include clear comments
3. Show expected output
4. Explain why it's interesting
5. Add to this index

## Learning Path

**Beginners** - Start with basic examples
**Intermediate** - Advanced examples
**Advanced** - Go integration and parallelization
**Production** - Real-world examples

---

**Last Updated**: 2025-11-21
