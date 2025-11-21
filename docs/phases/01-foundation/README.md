# Phase 1: Foundation - Lexer, Parser, and AST

## Overview

Phase 1 establishes the foundation of the PHP-Go interpreter by implementing the lexer (tokenizer), parser, and Abstract Syntax Tree (AST) representation. This phase transforms PHP source code into a structured tree representation that can be compiled.

## Goals

1. **Lexer**: Tokenize PHP source code into a stream of tokens
2. **Parser**: Parse token stream into an Abstract Syntax Tree
3. **AST**: Define all AST node types for PHP 8.4 syntax
4. **CLI**: Basic command-line interface to test the lexer and parser

## Dependencies

- **Inputs**: PHP source code (strings or files)
- **Outputs**: Abstract Syntax Tree (AST)
- **External**: None (pure Go implementation)

## Success Criteria

- [ ] Lexer can tokenize all PHP 8.4 syntax
- [ ] Parser can parse all PHP 8.4 language constructs
- [ ] AST accurately represents PHP code structure
- [ ] Can parse real-world PHP files without errors
- [ ] Comprehensive test coverage (>90%)
- [ ] Clear error messages for syntax errors

## Components

### 1. Lexer (Tokenizer)

**File**: `pkg/lexer/lexer.go`

**Responsibilities**:
- Read PHP source code character by character
- Recognize PHP tokens (keywords, operators, identifiers, literals, etc.)
- Handle string interpolation
- Support heredoc/nowdoc syntax
- Track source positions (file, line, column)
- Handle PHP tags (<?php, ?>, <?=)

**Token Types** (see PHP source: `Zend/zend_language_scanner.l`):
- Keywords: `if`, `else`, `while`, `for`, `function`, `class`, etc.
- Operators: `+`, `-`, `*`, `/`, `==`, `===`, `=>`, `->`, `::`, etc.
- Delimiters: `{`, `}`, `(`, `)`, `[`, `]`, `;`, `,`
- Literals: integers, floats, strings, heredoc, nowdoc
- Identifiers: variable names, function names, class names
- Special: T_VARIABLE (`$var`), T_STRING, T_WHITESPACE, T_COMMENT

**Key Challenges**:
- String interpolation: `"Hello $name"` or `"Hello {$obj->prop}"`
- Heredoc/Nowdoc: Complex multi-line string syntax
- PHP tags: Switching between PHP and non-PHP content
- Backslash escaping in strings
- Binary/octal/hex number literals
- Float literals with exponents

### 2. Parser

**File**: `pkg/parser/parser.go`

**Responsibilities**:
- Consume token stream from lexer
- Build Abstract Syntax Tree
- Enforce PHP grammar rules
- Generate helpful error messages
- Handle operator precedence
- Support PHP 8.4 syntax features

**Parser Approach**:

Option 1: **Recursive Descent Parser** (Recommended)
- Hand-written parser
- Full control over error messages
- Easy to understand and maintain
- Can handle PHP's complex grammar
- Reference: Go's parser package

Option 2: **Parser Generator** (Alternative)
- Use a tool like goyacc or participle
- Faster initial development
- Based on PHP's Bison grammar
- Less control over errors

**Recommendation**: Start with recursive descent for Phase 1, consider generator later if needed.

**Grammar Rules** (see PHP source: `Zend/zend_language_parser.y`):

Major categories:
- Top-level: namespace, use, class, function, const declarations
- Statements: if, while, for, foreach, switch, try-catch, etc.
- Expressions: binary ops, unary ops, function calls, array access, etc.
- Types: scalar types, union types, intersection types, nullable types
- Classes: properties, methods, constants, traits
- Functions: parameters, return types, references

**Operator Precedence** (lowest to highest):
1. `or`
2. `xor`
3. `and`
4. Assignment: `=`, `+=`, `-=`, etc.
5. Ternary: `? :`
6. Null coalesce: `??`
7. Logical OR: `||`
8. Logical AND: `&&`
9. Bitwise OR: `|`
10. Bitwise XOR: `^`
11. Bitwise AND: `&`
12. Equality: `==`, `===`, `!=`, `!==`
13. Comparison: `<`, `>`, `<=`, `>=`, `<=>`, `instanceof`
14. Shift: `<<`, `>>`
15. Addition: `+`, `-`, `.` (concatenation)
16. Multiplication: `*`, `/`, `%`
17. Unary: `!`, `~`, `++`, `--`, `@`, cast
18. Power: `**`
19. Array/Object access: `[]`, `->`
20. Clone, new

### 3. Abstract Syntax Tree (AST)

**File**: `pkg/ast/node.go`

**AST Node Hierarchy**:

```go
// Base node interface
type Node interface {
    Pos() Position       // Start position
    End() Position       // End position
    String() string      // String representation
}

// Statement nodes
type Stmt interface {
    Node
    stmtNode()
}

// Expression nodes
type Expr interface {
    Node
    exprNode()
}
```

**Statement Types** (`pkg/ast/stmt.go`):
- `EchoStmt` - echo statement
- `ExprStmt` - expression statement
- `BlockStmt` - block of statements
- `IfStmt` - if/elseif/else
- `WhileStmt` - while loop
- `DoWhileStmt` - do-while loop
- `ForStmt` - for loop
- `ForeachStmt` - foreach loop
- `SwitchStmt` - switch statement
- `MatchStmt` - match expression (PHP 8.0+)
- `BreakStmt` - break
- `ContinueStmt` - continue
- `ReturnStmt` - return
- `ThrowStmt` - throw exception
- `TryCatchStmt` - try-catch-finally
- `FunctionDecl` - function declaration
- `ClassDecl` - class declaration
- `InterfaceDecl` - interface declaration
- `TraitDecl` - trait declaration
- `EnumDecl` - enum declaration (PHP 8.1+)
- `NamespaceDecl` - namespace declaration
- `UseDecl` - use statement
- `ConstDecl` - constant declaration
- `GlobalStmt` - global statement
- `StaticStmt` - static statement
- `UnsetStmt` - unset statement

**Expression Types** (`pkg/ast/expr.go`):
- `IntLit` - integer literal
- `FloatLit` - float literal
- `StringLit` - string literal
- `ArrayLit` - array literal
- `NullLit` - null
- `BoolLit` - true/false
- `Variable` - $var
- `BinaryExpr` - binary operations (a + b)
- `UnaryExpr` - unary operations (!a, -a)
- `AssignExpr` - assignment (=, +=, etc.)
- `TernaryExpr` - ternary operator (a ? b : c)
- `CallExpr` - function call
- `MethodCallExpr` - method call ($obj->method())
- `StaticCallExpr` - static call (Class::method())
- `NewExpr` - new Class()
- `ArrayAccessExpr` - $arr[0]
- `PropertyAccessExpr` - $obj->prop
- `StaticPropertyExpr` - Class::$prop
- `ClosureExpr` - anonymous function
- `ArrowFunctionExpr` - arrow function (PHP 7.4+)
- `YieldExpr` - yield
- `YieldFromExpr` - yield from
- `InstanceOfExpr` - instanceof
- `CastExpr` - type casting
- `CloneExpr` - clone
- `IssetExpr` - isset()
- `EmptyExpr` - empty()
- `EvalExpr` - eval()
- `IncludeExpr` - include/require
- `ShellExecExpr` - backtick execution
- `ListExpr` - list() construct

**Type Annotations**:
- `TypeInfo` - type information
- `UnionType` - union types (int|string)
- `IntersectionType` - intersection types (A&B)
- `NullableType` - ?int

### 4. CLI Tool

**File**: `cmd/php-go/main.go`

**Features**:
- Parse PHP file and print AST
- Tokenize PHP file and print tokens
- Basic error reporting
- Debug mode for development

**Usage**:
```bash
# Parse and show AST
php-go parse script.php

# Tokenize and show tokens
php-go lex script.php

# Parse and print formatted
php-go fmt script.php
```

## Implementation Tasks

### Task 1.1: Token Definitions
**File**: `pkg/lexer/token.go`
**Estimated Effort**: 4 hours

- [ ] Define Token struct (type, literal, position)
- [ ] Define all token types as constants
- [ ] Implement token String() method
- [ ] Create token lookup tables for keywords
- [ ] Add token precedence information

**References**:
- PHP: `Zend/zend_language_scanner.l` (lines 1-500)
- Go: `go/token` package

### Task 1.2: Position Tracking
**File**: `pkg/lexer/position.go`
**Estimated Effort**: 2 hours

- [ ] Define Position struct (filename, line, column, offset)
- [ ] Implement position advancement
- [ ] Add position formatting for error messages
- [ ] Support position ranges

### Task 1.3: Basic Lexer
**File**: `pkg/lexer/lexer.go`
**Estimated Effort**: 16 hours

- [ ] Create Lexer struct with input buffer
- [ ] Implement character reading (peek, advance)
- [ ] Scan identifiers and keywords
- [ ] Scan numbers (int, float, hex, octal, binary)
- [ ] Scan operators and delimiters
- [ ] Scan variables ($var)
- [ ] Handle whitespace and comments
- [ ] Scan PHP tags (<?php, ?>, <?=)
- [ ] Basic error reporting

**Test Coverage**: 80%+

### Task 1.4: String Lexing
**File**: `pkg/lexer/strings.go`
**Estimated Effort**: 12 hours

- [ ] Scan single-quoted strings (no interpolation)
- [ ] Scan double-quoted strings (with interpolation)
- [ ] Handle escape sequences (\n, \t, \$, etc.)
- [ ] Scan heredoc syntax
- [ ] Scan nowdoc syntax
- [ ] Handle string interpolation tokens
- [ ] Support complex interpolation {$obj->prop}

**Test Coverage**: 90%+

**This is the most complex part of the lexer!**

### Task 1.5: Parser Foundation
**File**: `pkg/parser/parser.go`
**Estimated Effort**: 8 hours

- [ ] Create Parser struct
- [ ] Token buffer management (peek, advance, expect)
- [ ] Error recovery mechanisms
- [ ] Error message formatting
- [ ] Parse top-level structure (<?php ... ?>)
- [ ] Entry point: ParseFile() and ParseString()

### Task 1.6: Expression Parsing
**File**: `pkg/parser/expr.go`
**Estimated Effort**: 20 hours

- [ ] Parse primary expressions (literals, variables)
- [ ] Parse binary expressions with precedence
- [ ] Parse unary expressions
- [ ] Parse assignment expressions
- [ ] Parse ternary operator
- [ ] Parse function calls
- [ ] Parse method calls
- [ ] Parse array access
- [ ] Parse property access
- [ ] Parse new expressions
- [ ] Parse instanceof
- [ ] Parse closures and arrow functions
- [ ] Parse array literals
- [ ] Parse string interpolation

**Test Coverage**: 85%+

**Critical for all subsequent work!**

### Task 1.7: Statement Parsing
**File**: `pkg/parser/stmt.go`
**Estimated Effort**: 16 hours

- [ ] Parse echo statement
- [ ] Parse if/elseif/else
- [ ] Parse while loop
- [ ] Parse do-while loop
- [ ] Parse for loop
- [ ] Parse foreach loop
- [ ] Parse switch statement
- [ ] Parse match expression (PHP 8.0+)
- [ ] Parse break/continue/return
- [ ] Parse try-catch-finally
- [ ] Parse throw
- [ ] Parse global/static/unset
- [ ] Parse expression statements

**Test Coverage**: 85%+

### Task 1.8: Declaration Parsing
**File**: `pkg/parser/decl.go`
**Estimated Effort**: 16 hours

- [ ] Parse function declarations
- [ ] Parse function parameters (types, defaults, variadic)
- [ ] Parse return type hints
- [ ] Parse class declarations
- [ ] Parse class properties
- [ ] Parse class methods
- [ ] Parse class constants
- [ ] Parse traits and use statements
- [ ] Parse interfaces
- [ ] Parse enums (PHP 8.1+)
- [ ] Parse attributes (PHP 8.0+)
- [ ] Parse namespace declarations
- [ ] Parse use statements

**Test Coverage**: 85%+

### Task 1.9: Type Parsing
**File**: `pkg/parser/types.go`
**Estimated Effort**: 8 hours

- [ ] Parse scalar types (int, string, bool, float)
- [ ] Parse array/callable/iterable/object
- [ ] Parse class/interface type names
- [ ] Parse nullable types (?int)
- [ ] Parse union types (int|string)
- [ ] Parse intersection types (A&B)
- [ ] Parse DNF types (Disjunctive Normal Form)
- [ ] Parse mixed/never/void/static types

**Test Coverage**: 90%+

### Task 1.10: AST Node Definitions
**File**: `pkg/ast/*.go`
**Estimated Effort**: 12 hours

- [ ] Define Node interface
- [ ] Define all statement node types
- [ ] Define all expression node types
- [ ] Implement String() methods for debugging
- [ ] Add visitor pattern support
- [ ] Implement AST printer for debugging

**Test Coverage**: 70%+

### Task 1.11: CLI Tool
**File**: `cmd/php-go/main.go`
**Estimated Effort**: 6 hours

- [ ] Set up CLI framework (cobra or flag)
- [ ] Implement `lex` command (tokenize)
- [ ] Implement `parse` command (show AST)
- [ ] Implement `fmt` command (format PHP)
- [ ] Add file reading
- [ ] Add error reporting
- [ ] Add debug/verbose modes

### Task 1.12: Testing
**Estimated Effort**: 16 hours

- [ ] Unit tests for all token types
- [ ] Unit tests for lexer edge cases
- [ ] Unit tests for parser - expressions
- [ ] Unit tests for parser - statements
- [ ] Unit tests for parser - declarations
- [ ] Integration tests with real PHP files
- [ ] Error handling tests
- [ ] Test with PHP 8.4 syntax features
- [ ] Performance benchmarks

**Target**: 85%+ code coverage

## Test Plan

### Unit Tests

**Lexer Tests**:
```go
func TestLexerKeywords(t *testing.T) {
    input := "if else while for function class"
    lexer := NewLexer(input)
    // Assert tokens match expected
}

func TestLexerStrings(t *testing.T) {
    input := `"Hello $name" 'World'`
    lexer := NewLexer(input)
    // Assert tokens match expected
}

func TestLexerHeredoc(t *testing.T) {
    input := `<<<EOT
Hello World
EOT;`
    lexer := NewLexer(input)
    // Assert tokens match expected
}
```

**Parser Tests**:
```go
func TestParseExpression(t *testing.T) {
    tests := []struct {
        input    string
        expected ast.Expr
    }{
        {"1 + 2", &ast.BinaryExpr{...}},
        {"$x = 5", &ast.AssignExpr{...}},
        // ... more cases
    }
    for _, tt := range tests {
        parser := NewParser(tt.input)
        expr, err := parser.ParseExpression()
        assert.NoError(t, err)
        assert.Equal(t, tt.expected, expr)
    }
}
```

### Integration Tests

**Real PHP Files**:
- Parse simple PHP scripts
- Parse WordPress core files
- Parse Laravel framework files
- Parse Symfony framework files
- Ensure all parse without errors

### Benchmarks

```go
func BenchmarkLexer(b *testing.B) {
    input := loadLargeFile("large.php")
    for i := 0; i < b.N; i++ {
        lexer := NewLexer(input)
        for lexer.NextToken().Type != token.EOF {}
    }
}

func BenchmarkParser(b *testing.B) {
    input := loadLargeFile("large.php")
    for i := 0; i < b.N; i++ {
        parser := NewParser(input)
        parser.Parse()
    }
}
```

## Milestones

### Milestone 1.1: Basic Lexer (Week 1)
- Token definitions complete
- Can tokenize simple PHP code
- Basic operators and keywords working

### Milestone 1.2: Full Lexer (Week 2)
- String handling complete (including interpolation)
- Heredoc/nowdoc working
- All PHP tokens supported

### Milestone 1.3: Expression Parser (Week 3)
- Can parse all PHP expressions
- Operator precedence correct
- Function/method calls working

### Milestone 1.4: Statement Parser (Week 4)
- Can parse all control flow statements
- Loops and conditionals working
- Error handling (try-catch)

### Milestone 1.5: Declaration Parser (Week 5)
- Can parse functions and classes
- Type annotations working
- Attributes and enums supported

### Milestone 1.6: Complete Phase 1 (Week 6)
- All PHP 8.4 syntax supported
- Comprehensive test coverage
- Can parse real-world PHP projects
- CLI tool functional

## Estimated Timeline

**Total Effort**: ~140 hours (6-7 weeks part-time)

- Lexer: ~35 hours
- Parser: ~60 hours
- AST: ~12 hours
- CLI: ~6 hours
- Testing: ~16 hours
- Documentation: ~10 hours

## Dependencies for Next Phase

Phase 2 (Compiler) requires:
- Complete AST representation ✓
- Parser that produces valid AST ✓
- Visitor pattern for AST traversal ✓

## References

**PHP Source Code**:
- `php-src/Zend/zend_language_scanner.l` - Lexer definition
- `php-src/Zend/zend_language_parser.y` - Parser grammar
- `php-src/Zend/zend_ast.c` - AST implementation
- `php-src/Zend/zend_ast.h` - AST definitions

**Go Resources**:
- `go/scanner` - Go lexer implementation
- `go/parser` - Go parser implementation
- `go/ast` - Go AST package

**Books & Articles**:
- "Writing An Interpreter In Go" by Thorsten Ball
- "Crafting Interpreters" by Robert Nystrom
- PHP Language Specification: https://github.com/php/php-langspec

## Progress Tracking

- [ ] Task 1.1: Token Definitions (0%)
- [ ] Task 1.2: Position Tracking (0%)
- [ ] Task 1.3: Basic Lexer (0%)
- [ ] Task 1.4: String Lexing (0%)
- [ ] Task 1.5: Parser Foundation (0%)
- [ ] Task 1.6: Expression Parsing (0%)
- [ ] Task 1.7: Statement Parsing (0%)
- [ ] Task 1.8: Declaration Parsing (0%)
- [ ] Task 1.9: Type Parsing (0%)
- [ ] Task 1.10: AST Node Definitions (0%)
- [ ] Task 1.11: CLI Tool (0%)
- [ ] Task 1.12: Testing (0%)

**Overall Phase 1 Progress**: 0%

---

**Phase Start**: TBD
**Phase End**: TBD
**Status**: Not Started
**Last Updated**: 2025-11-21
