package ast

import (
	"github.com/krizos/php-go/pkg/lexer"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string // Returns the literal value of the token
	String() string       // Returns a string representation for debugging
}

// Stmt represents a statement node
type Stmt interface {
	Node
	statementNode()
}

// Expr represents an expression node
type Expr interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	out := ""
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

// Identifier represents an identifier (variable name, function name, etc.)
type Identifier struct {
	Token lexer.Token // The IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Placeholder statement types (will be implemented in Task 1.7)

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Token      lexer.Token // The first token of the expression
	Expression Expr
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      lexer.Token // The { token
	Statements []Stmt
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	out := "{"
	for _, s := range bs.Statements {
		out += s.String()
	}
	out += "}"
	return out
}

// Placeholder expression types (will be implemented in Task 1.6)

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// NullLiteral represents a null literal
type NullLiteral struct {
	Token lexer.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NullLiteral) String() string       { return "null" }

// Variable represents a PHP variable ($var)
type Variable struct {
	Token lexer.Token
	Name  string // Without the $ prefix
}

func (v *Variable) expressionNode()      {}
func (v *Variable) TokenLiteral() string { return v.Token.Literal }
func (v *Variable) String() string       { return v.Token.Literal }

// Additional node types will be added in Tasks 1.6-1.10
