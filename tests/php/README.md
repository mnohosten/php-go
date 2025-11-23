# PHP Test Scripts

This directory contains PHP test scripts used for development and integration testing of the PHP-Go interpreter.

## Purpose

These are **PHP source files** (`.php`) used to test the interpreter's functionality, NOT `.phpt` format test files (those go in `tests/phpt/` when we integrate PHP's official test suite).

## Directory Structure

- **basic/** - Tests for basic language features (variables, expressions, control flow, strings, arrays)
- **vm/** - Tests for VM-specific features and opcodes
- **stdlib/** - Tests for standard library functions (strlen, array_map, etc.)
- **integration/** - End-to-end integration tests

## Usage

Run tests with the php-go CLI:

```bash
# Run a single test
./php-go tests/php/basic/variables.php

# Run all tests in a directory
for f in tests/php/basic/*.php; do ./php-go "$f"; done

# Compare output with PHP
php tests/php/basic/variables.php
./php-go tests/php/basic/variables.php
```

## Test Format

Tests are simple PHP scripts that demonstrate specific features. They should:

1. Be focused on testing one specific feature or function
2. Include comments explaining what's being tested
3. Use `echo` to produce verifiable output
4. Note any known issues or limitations

Example:

```php
<?php
// Test: Variable assignment and retrieval
$x = 5;
echo $x; // Should output: 5
```

## Difference from .phpt Files

- **These files** (`.php`): Simple PHP scripts for quick testing during development
- **PHPT files** (`.phpt`): PHP's official test format with TEST, FILE, EXPECT sections
  - Located in `tests/phpt/` (to be added in Phase 10.2)
  - Run with the PHPT test runner in `tests/phptest/`

## Contributing

When adding new features to PHP-Go, add corresponding test files here to ensure the feature works correctly and to prevent regressions.
