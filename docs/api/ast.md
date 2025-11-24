# AST API Reference

Package: `github.com/krizos/php-go/pkg/ast`

## Overview

The ast package defines all Abstract Syntax Tree node types for representing PHP code. The AST is a tree structure that captures the syntactic structure of PHP code after parsing. This package includes 65+ node types covering all PHP 8.4 language features.

## Base Interfaces

All AST nodes implement one of these base interfaces:

### Node

Base interface for all AST nodes.

```go
type Node interface {
    TokenLiteral() string
    String() string
}
```

### Stmt

Interface for statement nodes.

```go
type Stmt interface {
    Node
    statementNode()
}
```

### Expr

Interface for expression nodes.

```go
type Expr interface {
    Node
    expressionNode()
}
```

## Root Node

### Program

The root node of any PHP AST.

```go
type Program struct {
    Statements []Stmt
}
```

## Statement Nodes

### Control Flow

**IfStatement** - If/elseif/else conditional
```go
type IfStatement struct {
    Token       lexer.Token
    Condition   Expr
    Consequence *BlockStatement
    ElseIfs     []*ElseIfClause
    Alternative *BlockStatement
}
```

**WhileStatement** - While loop
```go
type WhileStatement struct {
    Token     lexer.Token
    Condition Expr
    Body      *BlockStatement
}
```

**DoWhileStatement** - Do-while loop
```go
type DoWhileStatement struct {
    Token     lexer.Token
    Body      *BlockStatement
    Condition Expr
}
```

**ForStatement** - For loop
```go
type ForStatement struct {
    Token      lexer.Token
    Init       []Expr
    Condition  []Expr
    Increment  []Expr
    Body       *BlockStatement
}
```

**ForeachStatement** - Foreach loop
```go
type ForeachStatement struct {
    Token  lexer.Token
    Array  Expr
    Key    Expr
    Value  Expr
    ByRef  bool
    Body   *BlockStatement
}
```

**SwitchStatement** - Switch statement
```go
type SwitchStatement struct {
    Token   lexer.Token
    Subject Expr
    Cases   []*SwitchCase
}
```

**TryStatement** - Try-catch-finally
```go
type TryStatement struct {
    Token        lexer.Token
    Body         *BlockStatement
    CatchClauses []*CatchClause
    Finally      *BlockStatement
}
```

**ThrowStatement** - Throw exception
```go
type ThrowStatement struct {
    Token      lexer.Token
    Expression Expr
}
```

### Jump Statements

**ReturnStatement** - Return from function
```go
type ReturnStatement struct {
    Token       lexer.Token
    ReturnValue Expr
}
```

**BreakStatement** - Break from loop
```go
type BreakStatement struct {
    Token lexer.Token
    Depth Expr
}
```

**ContinueStatement** - Continue to next iteration
```go
type ContinueStatement struct {
    Token lexer.Token
    Depth Expr
}
```

### Other Statements

**ExpressionStatement** - Expression used as statement
```go
type ExpressionStatement struct {
    Token      lexer.Token
    Expression Expr
}
```

**BlockStatement** - Block of statements
```go
type BlockStatement struct {
    Token      lexer.Token
    Statements []Stmt
}
```

**EchoStatement** - Echo output
```go
type EchoStatement struct {
    Token       lexer.Token
    Expressions []Expr
}
```

**GlobalStatement** - Declare global variables
```go
type GlobalStatement struct {
    Token     lexer.Token
    Variables []*Identifier
}
```

**UnsetStatement** - Unset variables
```go
type UnsetStatement struct {
    Token     lexer.Token
    Variables []Expr
}
```

**NamespaceStatement** - Namespace declaration
```go
type NamespaceStatement struct {
    Token      lexer.Token
    Name       *NamespaceName
    Body       *BlockStatement
    Statements []Stmt
}
```

**UseStatement** - Import declaration
```go
type UseStatement struct {
    Token  lexer.Token
    Uses   []*UseImport
    Type   string
    Prefix string
}
```

**DeclareStatement** - Declare directives
```go
type DeclareStatement struct {
    Token      lexer.Token
    Directives map[string]interface{}
    Body       Stmt
}
```

## Expression Nodes

### Literals

**IntegerLiteral** - Integer value
```go
type IntegerLiteral struct {
    Token lexer.Token
    Value int64
}
```

**FloatLiteral** - Float value
```go
type FloatLiteral struct {
    Token lexer.Token
    Value float64
}
```

**StringLiteral** - String value
```go
type StringLiteral struct {
    Token lexer.Token
    Value string
}
```

**BooleanLiteral** - Boolean value
```go
type BooleanLiteral struct {
    Token lexer.Token
    Value bool
}
```

**NullLiteral** - Null value
```go
type NullLiteral struct {
    Token lexer.Token
}
```

**ArrayExpression** - Array literal
```go
type ArrayExpression struct {
    Token    lexer.Token
    Elements []ArrayElement
}

type ArrayElement struct {
    Key   Expr
    Value Expr
}
```

### Variables and Identifiers

**Variable** - PHP variable
```go
type Variable struct {
    Token lexer.Token
    Name  string
}
```

**Identifier** - Name/identifier
```go
type Identifier struct {
    Token lexer.Token
    Value string
}
```

**NamespaceName** - Namespaced name
```go
type NamespaceName struct {
    Token lexer.Token
    Parts []string
}
```

### Operators

**PrefixExpression** - Prefix operator
```go
type PrefixExpression struct {
    Token    lexer.Token
    Operator string
    Right    Expr
}
```

**InfixExpression** - Binary operator
```go
type InfixExpression struct {
    Token    lexer.Token
    Left     Expr
    Operator string
    Right    Expr
}
```

**AssignmentExpression** - Assignment
```go
type AssignmentExpression struct {
    Token    lexer.Token
    Left     Expr
    Operator string
    Right    Expr
}
```

**TernaryExpression** - Ternary conditional
```go
type TernaryExpression struct {
    Token       lexer.Token
    Condition   Expr
    Consequence Expr
    Alternative Expr
}
```

### Array and Object Access

**IndexExpression** - Array/string access
```go
type IndexExpression struct {
    Token lexer.Token
    Left  Expr
    Index Expr
}
```

**PropertyExpression** - Property access
```go
type PropertyExpression struct {
    Token    lexer.Token
    Object   Expr
    Property Expr
}
```

**NullsafePropertyExpression** - Nullsafe property access
```go
type NullsafePropertyExpression struct {
    Token    lexer.Token
    Object   Expr
    Property Expr
}
```

**StaticPropertyExpression** - Static property access
```go
type StaticPropertyExpression struct {
    Token    lexer.Token
    Class    Expr
    Property Expr
}
```

### Function and Method Calls

**CallExpression** - Function call
```go
type CallExpression struct {
    Token     lexer.Token
    Function  Expr
    Arguments []*Argument
}
```

**MethodCallExpression** - Method call
```go
type MethodCallExpression struct {
    Token     lexer.Token
    Object    Expr
    Method    Expr
    Arguments []*Argument
}
```

**StaticCallExpression** - Static method call
```go
type StaticCallExpression struct {
    Token     lexer.Token
    Class     Expr
    Method    Expr
    Arguments []*Argument
}
```

**Argument** - Function/method argument
```go
type Argument struct {
    Token  lexer.Token
    Name   string  // For named arguments
    Value  Expr
    Unpack bool    // For ... spread operator
}
```

### Object Operations

**NewExpression** - Object creation
```go
type NewExpression struct {
    Token     lexer.Token
    Class     Expr
    Arguments []*Argument
}
```

**CloneExpression** - Object cloning
```go
type CloneExpression struct {
    Token  lexer.Token
    Object Expr
}
```

**InstanceofExpression** - Type check
```go
type InstanceofExpression struct {
    Token lexer.Token
    Left  Expr
    Right Expr
}
```

### Language Constructs

**IssetExpression** - isset() check
```go
type IssetExpression struct {
    Token     lexer.Token
    Variables []Expr
}
```

**EmptyExpression** - empty() check
```go
type EmptyExpression struct {
    Token    lexer.Token
    Variable Expr
}
```

**ListExpression** - list() destructuring
```go
type ListExpression struct {
    Token    lexer.Token
    Elements []*ListElement
}
```

**IncludeExpression** - include/require
```go
type IncludeExpression struct {
    Token lexer.Token
    Path  Expr
    Type  string
}
```

### Closures and Functions

**ClosureExpression** - Anonymous function
```go
type ClosureExpression struct {
    Token      lexer.Token
    Parameters []*Parameter
    Use        []*UseClause
    ReturnType Expr
    Body       *BlockStatement
    ByRef      bool
    Static     bool
}
```

**ArrowFunctionExpression** - Arrow function
```go
type ArrowFunctionExpression struct {
    Token      lexer.Token
    Parameters []*Parameter
    ReturnType Expr
    Body       Expr
    ByRef      bool
    Static     bool
}
```

**FirstClassCallableExpression** - First-class callable
```go
type FirstClassCallableExpression struct {
    Token    lexer.Token
    Callable Expr
}
```

### Other Expressions

**CastExpression** - Type casting
```go
type CastExpression struct {
    Token lexer.Token
    Type  string
    Expr  Expr
}
```

**GroupedExpression** - Parenthesized expression
```go
type GroupedExpression struct {
    Token lexer.Token
    Expr  Expr
}
```

**MatchExpression** - Match expression (PHP 8.0+)
```go
type MatchExpression struct {
    Token   lexer.Token
    Subject Expr
    Arms    []*MatchArm
}
```

**MagicConstant** - Magic constant
```go
type MagicConstant struct {
    Token lexer.Token
    Kind  lexer.TokenType
}
```

**InterpolatedStringExpression** - String with interpolation
```go
type InterpolatedStringExpression struct {
    Token lexer.Token
    Parts []Expr
}
```

**ClassNameExpression** - ::class constant
```go
type ClassNameExpression struct {
    Token lexer.Token
    Class Expr
}
```

## Declaration Nodes

### Functions

**FunctionDeclaration** - Function declaration
```go
type FunctionDeclaration struct {
    Token      lexer.Token
    Name       *Identifier
    Parameters []*Parameter
    ReturnType Expr
    Body       *BlockStatement
    ByRef      bool
    Attributes []*AttributeGroup
}
```

**Parameter** - Function/method parameter
```go
type Parameter struct {
    Name         *Variable
    Type         Expr
    DefaultValue Expr
    ByRef        bool
    Variadic     bool
    Attributes   []*AttributeGroup
    Visibility   string
    Readonly     bool
}
```

### Classes

**ClassDeclaration** - Class declaration
```go
type ClassDeclaration struct {
    Token      lexer.Token
    Name       *Identifier
    Extends    *Identifier
    Implements []*Identifier
    Body       []Stmt
    Modifiers  []string
    Attributes []*AttributeGroup
}
```

**PropertyDeclaration** - Class property
```go
type PropertyDeclaration struct {
    Token      lexer.Token
    Visibility string
    Attributes []*AttributeGroup
    Static     bool
    Readonly   bool
    Type       Expr
    Properties []*PropertyItem
}
```

**MethodDeclaration** - Class method
```go
type MethodDeclaration struct {
    Token      lexer.Token
    Visibility string
    Static     bool
    Abstract   bool
    Final      bool
    Name       *Identifier
    Parameters []*Parameter
    ReturnType Expr
    Body       *BlockStatement
    ByRef      bool
    Attributes []*AttributeGroup
}
```

**ClassConstantDeclaration** - Class constant
```go
type ClassConstantDeclaration struct {
    Token      lexer.Token
    Visibility string
    Constants  []*ConstantItem
}
```

### Interfaces and Traits

**InterfaceDeclaration** - Interface declaration
```go
type InterfaceDeclaration struct {
    Token   lexer.Token
    Name    *Identifier
    Extends []*Identifier
    Body    []*MethodSignature
}
```

**TraitDeclaration** - Trait declaration
```go
type TraitDeclaration struct {
    Token lexer.Token
    Name  *Identifier
    Body  []Stmt
}
```

**TraitUse** - Trait usage
```go
type TraitUse struct {
    Token       lexer.Token
    Traits      []*Identifier
    Adaptations []TraitAdaptation
}
```

## Type Nodes

**NullableType** - Nullable type
```go
type NullableType struct {
    Token lexer.Token
    Type  Expr
}
```

**UnionType** - Union type
```go
type UnionType struct {
    Token lexer.Token
    Types []Expr
}
```

**IntersectionType** - Intersection type
```go
type IntersectionType struct {
    Token lexer.Token
    Types []Expr
}
```

## Attribute Nodes (PHP 8.0+)

**Attribute** - Single attribute
```go
type Attribute struct {
    Token     lexer.Token
    Name      Expr
    Arguments []Expr
    Named     map[string]Expr
}
```

**AttributeGroup** - Attribute group
```go
type AttributeGroup struct {
    Token      lexer.Token
    Attributes []*Attribute
}
```

## Visitor Pattern

For traversing the AST, see `/Users/krizos/code/mnohosten/php-go/pkg/ast/visitor.go` for the Visitor interface implementation.

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/krizos/php-go/pkg/ast"
    "github.com/krizos/php-go/pkg/lexer"
    "github.com/krizos/php-go/pkg/parser"
)

func main() {
    input := `<?php $x = 5 + 10;`

    l := lexer.New(input, "test.php")
    p := parser.New(l)
    program := p.ParseProgram()

    // Walk through statements
    for _, stmt := range program.Statements {
        if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
            if assign, ok := exprStmt.Expression.(*ast.AssignmentExpression); ok {
                fmt.Printf("Assignment: %s\n", assign.String())
            }
        }
    }
}
```

## Implementation Notes

- All nodes carry position information via `Token` field
- The `String()` method provides debugging output
- Node types follow PHP's grammar structure
- Supports all PHP 8.4 features including attributes, enums, readonly properties
- Named arguments are supported for function/method calls
- Type system nodes support nullable, union, and intersection types
