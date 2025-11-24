# Lexer API Reference

Package: `github.com/krizos/php-go/pkg/lexer`

## Overview

The lexer package implements PHP tokenization. It converts PHP source code into a stream of tokens that can be consumed by the parser. The lexer supports full PHP 8.4 syntax including:

- All PHP operators and keywords
- String literals (single-quoted, double-quoted, heredoc, nowdoc)
- String interpolation
- PHP tags (`<?php`, `<?=`, `?>`)
- Inline HTML
- Comments (single-line and multi-line)
- Numeric literals (integers, floats, hex, binary, octal)
- All PHP 8.4 keywords and magic constants

## Main Types

### Lexer

The main tokenizer that converts source code into tokens.

```go
type Lexer struct {
    // Private fields
}
```

**Constructor:**

```go
func New(input, filename string) *Lexer
```

Creates a new lexer for the given input source code. The filename is used for error reporting.

**Methods:**

```go
func (l *Lexer) NextToken() Token
```

Returns the next token from the input. Call repeatedly until EOF token is returned.

```go
func (l *Lexer) Error(message string) Token
```

Creates an error token with the given message.

```go
func (l *Lexer) String() string
```

Returns a string representation of the lexer state for debugging.

### Token

Represents a single lexical token in PHP source code.

```go
type Token struct {
    Type    TokenType
    Literal string
    Pos     Position
}
```

**Fields:**
- `Type`: The type of token (keyword, operator, literal, etc.)
- `Literal`: The actual text of the token from source code
- `Pos`: Source code position information

**Methods:**

```go
func (t Token) String() string
```

Returns a human-readable representation of the token.

### TokenType

An integer constant representing the type of a token.

```go
type TokenType int
```

**Token Type Categories:**

**Special Tokens:**
- `ILLEGAL` - Invalid/unknown token
- `EOF` - End of file
- `COMMENT` - Single-line or multi-line comment

**Literals:**
- `INTEGER` - Integer literal (123, 0x1A, 0b1010, 0o777)
- `FLOAT` - Floating-point literal (123.45, 1.23e4)
- `STRING` - String literal
- `HEREDOC` - Heredoc string
- `NOWDOC` - Nowdoc string

**Identifiers:**
- `IDENT` - Identifier (function name, class name)
- `VARIABLE` - PHP variable ($var)

**Keywords:** All PHP keywords including:
- Control flow: `IF`, `ELSE`, `ELSEIF`, `WHILE`, `FOR`, `FOREACH`, `SWITCH`, `MATCH`, `BREAK`, `CONTINUE`
- Functions: `FUNCTION`, `RETURN`, `YIELD`, `YIELD_FROM`
- OOP: `CLASS`, `INTERFACE`, `TRAIT`, `ENUM`, `EXTENDS`, `IMPLEMENTS`, `NEW`, `CLONE`
- Modifiers: `PUBLIC`, `PROTECTED`, `PRIVATE`, `STATIC`, `ABSTRACT`, `FINAL`, `READONLY`
- Exception handling: `TRY`, `CATCH`, `FINALLY`, `THROW`
- Others: `NAMESPACE`, `USE`, `CONST`, `GLOBAL`, `ECHO`, `ISSET`, `EMPTY`, `UNSET`

**Type Keywords:**
- `INT`, `FLOAT_TYPE`, `BOOL`, `STRING_TYPE`, `ARRAY`
- `OBJECT`, `CALLABLE`, `ITERABLE`, `MIXED`, `VOID`, `NEVER`
- `TRUE`, `FALSE`, `NULL`

**Operators:**
- Arithmetic: `PLUS` (+), `MINUS` (-), `ASTERISK` (*), `SLASH` (/), `PERCENT` (%), `POWER` (**)
- Comparison: `EQ` (==), `IDENTICAL` (===), `NE` (!=), `NOT_IDENTICAL` (!==), `LT` (<), `LE` (<=), `GT` (>), `GE` (>=), `SPACESHIP` (<=>)
- Logical: `LOGICAL_AND` (&&), `LOGICAL_OR` (||), `LOGICAL_NOT` (!)
- Bitwise: `BITWISE_AND` (&), `BITWISE_OR` (|), `BITWISE_XOR` (^), `BITWISE_NOT` (~), `SL` (<<), `SR` (>>)
- Assignment: `ASSIGN` (=), `PLUS_ASSIGN` (+=), `MINUS_ASSIGN` (-=), and all compound assignments
- String: `CONCAT` (.), `CONCAT_ASSIGN` (.=)
- Other: `INC` (++), `DEC` (--), `INSTANCEOF`, `COALESCE` (??), `COALESCE_ASSIGN` (??=)

**Delimiters:**
- `LPAREN` ((), `RPAREN` ()), `LBRACE` ({), `RBRACE` (}), `LBRACKET` ([), `RBRACKET` (])
- `SEMICOLON` (;), `COMMA` (,), `COLON` (:), `QUESTION` (?)
- `DOUBLE_ARROW` (=>), `OBJECT_OPERATOR` (->), `PAAMAYIM_NEKUDOTAYIM` (::)
- `NULLSAFE_OPERATOR` (?->), `ELLIPSIS` (...)

**PHP Tags:**
- `OPEN_TAG` - `<?php` or `<?`
- `OPEN_TAG_ECHO` - `<?=`
- `CLOSE_TAG` - `?>`
- `INLINE_HTML` - HTML content between tags

**Magic Constants:**
- `LINE_CONST` (__LINE__), `FILE_CONST` (__FILE__), `DIR_CONST` (__DIR__)
- `FUNCTION_CONST` (__FUNCTION__), `CLASS_CONST` (__CLASS__), `TRAIT_CONST` (__TRAIT__)
- `METHOD_CONST` (__METHOD__), `NAMESPACE_CONST` (__NAMESPACE__)

**Methods:**

```go
func (tt TokenType) String() string
```

Returns the name of the token type.

```go
func (tt TokenType) IsKeyword() bool
func (tt TokenType) IsLiteral() bool
func (tt TokenType) IsOperator() bool
```

Type checking methods.

### Position

Represents a location in source code for error reporting.

```go
type Position struct {
    Filename string
    Offset   int
    Line     int
    Column   int
}
```

**Methods:**

```go
func (p Position) String() string
```

Returns formatted position as "filename:line:column".

```go
func (p Position) IsValid() bool
func (p Position) Before(other Position) bool
func (p Position) After(other Position) bool
```

Position comparison methods.

## Functions

```go
func LookupIdent(ident string) TokenType
```

Checks if an identifier is a PHP keyword and returns the appropriate token type. Returns `IDENT` if not a keyword.

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/lexer"
)

func main() {
    input := `<?php
    $x = 10;
    echo $x + 5;
    `

    l := lexer.New(input, "example.php")

    for {
        tok := l.NextToken()
        fmt.Printf("%s\n", tok)

        if tok.Type == lexer.EOF {
            break
        }
    }
}
```

## Implementation Notes

- The lexer maintains position tracking (line, column, offset) for accurate error reporting
- It automatically handles PHP mode switching between PHP code and inline HTML
- String interpolation is recognized but currently included as part of the string literal
- Numeric underscores are supported (PHP 7.4+): `1_000_000`
- All PHP 8.4 features are supported including attributes (`#[...]`)
- Comments are tokenized but typically skipped by the parser

## Performance

The lexer is optimized for speed:
- Single-pass scanning
- No backtracking
- Efficient character lookahead
- Benchmark results (Phase 1): 1-18μs per token depending on complexity
