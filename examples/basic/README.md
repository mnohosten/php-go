# Basic PHP-Go Examples

This directory contains basic PHP examples that demonstrate the fundamental features of PHP-Go.

## Example Files

### hello.php
**Purpose**: Simple "Hello World" example
**Features**: Basic echo statement
**Expected Output**: See `hello.expected`

### variables.php
**Purpose**: Variable assignment and basic types
**Features**:
- Integer, string, float, and boolean variables
- Variable reassignment
- Multiple variable declarations
**Expected Output**: See `variables.expected`

### expressions.php
**Purpose**: Arithmetic and comparison expressions
**Features**:
- Basic arithmetic operators (+, -, *, /, %)
- Operator precedence
- Comparison operators (==, >, <)
- Boolean output representation
**Expected Output**: See `expressions.expected`

### control_flow.php
**Purpose**: Control flow structures
**Features**:
- if/else statements
- for loops with increment
- while loops
- Nested conditionals
**Expected Output**: See `control_flow.expected`

### functions.php
**Purpose**: Function usage examples
**Features**:
- Built-in function calls (strlen)
- Comments about user-defined functions
**Expected Output**: See `functions.expected`

### arrays.php
**Purpose**: Array creation and operations
**Features**:
- Array literal syntax [1, 2, 3]
- Comments about future associative array support
**Expected Output**: See `arrays.expected`

### strings.php
**Purpose**: String operations
**Features**:
- String literals and variables
- Escape sequences (\n)
- Built-in string functions (strlen)
- Empty strings
**Expected Output**: See `strings.expected`

## Expected Output Files

Each `.expected` file contains the exact output that should be produced when the corresponding PHP file is executed with PHP 8.4. These files serve as:

1. **Test Oracles**: Automated tests can compare php-go output with these expected outputs
2. **Documentation**: Developers can see exactly what behavior to expect
3. **Regression Testing**: Changes to the implementation can be validated against these baselines

## Testing

To verify that php-go produces the correct output for these examples:

```bash
# Run all example tests
go test ./tests -run TestBasicExamples -v

# Run comprehensive example test suite
go test ./tests -run TestExampleFiles -v

# Run a specific example with php-go
./php-go examples/basic/hello.php

# Compare with PHP 8.4 reference
php examples/basic/hello.php
```

## Adding New Examples

When adding a new example file:

1. Create the `.php` file with clear comments
2. Run it with PHP 8.4: `php examples/basic/yourfile.php > examples/basic/yourfile.expected`
3. Verify the output is correct
4. Add a test case in `tests/examples_test.go` if needed
5. Update this README with the new example

## Status

All basic examples currently pass with php-go (7/7 - 100%).
