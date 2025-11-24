# Parser API Reference

Package: `github.com/krizos/php-go/pkg/parser`

## Overview

The parser package converts a stream of tokens from the lexer into an Abstract Syntax Tree (AST). It implements a Pratt parser for expressions with operator precedence climbing, and recursive descent for statements and declarations. The parser supports full PHP 8.4 syntax.

## Main Types

### Parser

The main parser that converts tokens into an AST.

```go
type Parser struct {
    // Private fields
}
```

**Constructor:**

```go
func New(l *lexer.Lexer) *Parser
```

Creates a new parser from a lexer instance.

**Methods:**

```go
func (p *Parser) ParseProgram() *ast.Program
```

Parses the entire PHP program and returns the root AST node. This is the main entry point for parsing.

```go
func (p *Parser) Errors() []string
```

Returns all parsing errors encountered during parsing.

```go
func (p *Parser) HasErrors() bool
```

Returns true if any parsing errors occurred.

## Precedence Levels

The parser uses the following precedence levels for expressions (lowest to highest):

1. `LOWEST` - Lowest precedence
2. `LOGICAL_OR` - or, ||
3. `LOGICAL_XOR` - xor
4. `LOGICAL_AND` - and, &&
5. `ASSIGNMENT` - =, +=, -=, etc.
6. `TERNARY` - ? :
7. `COALESCE` - ??
8. `BITWISE_OR` - |
9. `BITWISE_XOR` - ^
10. `BITWISE_AND` - &
11. `EQUALITY` - ==, ===, !=, !==
12. `COMPARISON` - <, >, <=, >=, <=>, instanceof
13. `SHIFT` - <<, >>
14. `CONCAT` - .
15. `SUM` - +, -
16. `PRODUCT` - *, /, %
17. `POWER` - **
18. `UNARY` - !, ~, ++, --, @, cast
19. `POSTFIX` - [], ->, ::, ()
20. `NEW_CLONE` - new, clone

## Functions

```go
func ParseString(input string) (*ast.Program, []string)
```

Convenience function to parse PHP source code from a string. Returns the AST and any parsing errors.

```go
func ParseFile(filename string) (*ast.Program, error)
```

Convenience function to parse a PHP file (not yet fully implemented).

## Supported PHP Syntax

The parser supports all PHP 8.4 syntax including:

### Statements
- Expression statements
- Block statements (`{ ... }`)
- Control flow: `if`, `else`, `elseif`, `while`, `do-while`, `for`, `foreach`, `switch`
- Jump statements: `break`, `continue`, `return`
- Exception handling: `try`, `catch`, `finally`, `throw`
- `echo` statement
- `global`, `unset` statements
- `namespace` and `use` declarations
- `declare` directives

### Expressions
- Literals: integers, floats, strings, booleans, null, arrays
- Variables: `$var`
- Binary operators: arithmetic, comparison, logical, bitwise, string concatenation
- Unary operators: `!`, `~`, `+`, `-`, `++`, `--`, `@`
- Ternary operator: `? :`
- Null coalescing: `??`, `??=`
- Assignment: `=`, `+=`, `-=`, etc.
- Array access: `$arr[$key]`
- Property access: `$obj->prop`, `$obj?->prop`
- Static access: `Class::$prop`, `Class::method()`
- Function calls: `func($args)`
- Method calls: `$obj->method($args)`, `Class::method($args)`
- Object creation: `new Class($args)`
- Clone: `clone $obj`
- `instanceof` checks
- `isset()`, `empty()` language constructs
- `list()` destructuring
- Closures: `function() use ($var) { ... }`
- Arrow functions: `fn($x) => $x * 2`
- First-class callables: `strlen(...)`
- Type casting: `(int)$var`, `(string)$var`, etc.
- `match` expressions (PHP 8.0+)

### Declarations
- Functions: `function name($params): Type { ... }`
- Classes: `class Name extends Parent implements Interface { ... }`
- Interfaces: `interface Name { ... }`
- Traits: `trait Name { ... }`
- Enums (PHP 8.1+)
- Properties with visibility, types, readonly
- Methods with visibility, modifiers (static, abstract, final)
- Class constants
- Trait usage with adaptations

### Type System
- Scalar types: `int`, `float`, `bool`, `string`, `array`
- Special types: `void`, `never`, `mixed`, `object`, `callable`, `iterable`
- Nullable types: `?Type`
- Union types: `Type1|Type2|Type3` (PHP 8.0+)
- Intersection types: `Type1&Type2` (PHP 8.1+)

### PHP 8+ Features
- Named arguments: `func(name: $value)`
- Attributes: `#[AttributeName(args)]`
- Match expressions: `match ($x) { ... }`
- Nullsafe operator: `$obj?->method()`
- Constructor property promotion

## Error Handling

The parser includes error recovery mechanisms:

- Synchronization after errors to continue parsing
- Detailed error messages with position information
- Multiple errors can be collected in a single parse
- Errors are accessible via `Parser.Errors()`

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/lexer"
    "github.com/krizos/php-go/pkg/parser"
)

func main() {
    input := `<?php
    function add($a, $b) {
        return $a + $b;
    }

    echo add(5, 10);
    `

    l := lexer.New(input, "example.php")
    p := parser.New(l)

    program := p.ParseProgram()

    if p.HasErrors() {
        for _, err := range p.Errors() {
            fmt.Println(err)
        }
        return
    }

    fmt.Printf("Parsed %d statements\n", len(program.Statements))
}
```

## Parsing Strategies

### Expression Parsing

The parser uses Pratt parsing (precedence climbing) for expressions:
- Efficient handling of operator precedence
- Supports prefix and infix operators
- Easy to extend with new operators
- No need for complex grammar rules

### Statement Parsing

The parser uses recursive descent for statements:
- Top-down parsing approach
- Each statement type has its own parsing function
- Easy to understand and maintain

### Error Recovery

When a parse error occurs:
1. Error is recorded in the error list
2. Parser attempts to synchronize to a statement boundary
3. Parsing continues from the next valid statement
4. Multiple errors can be reported in one pass

## Implementation Details

- **Token Management**: Parser maintains `curToken` and `peekToken` for lookahead
- **Comment Skipping**: Comments are automatically skipped during parsing
- **PHP Tags**: Opening tags (`<?php`, `<?=`) are handled automatically
- **Inline HTML**: HTML between `?>` and `<?php` is converted to echo statements

## Performance

Parser benchmark results (Phase 1):
- Simple expressions: ~10μs
- Complex statements: ~50-90μs
- Full programs: scales linearly with code size
- Memory efficient: minimal allocations during parsing
