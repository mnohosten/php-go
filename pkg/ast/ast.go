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

// FloatLiteral represents a floating-point literal
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// PrefixExpression represents a prefix operator expression (!, -, +, ~, ++, --)
type PrefixExpression struct {
	Token    lexer.Token // The prefix operator token
	Operator string
	Right    Expr
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression represents a binary operator expression
type InfixExpression struct {
	Token    lexer.Token // The operator token
	Left     Expr
	Operator string
	Right    Expr
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// AssignmentExpression represents an assignment operation
type AssignmentExpression struct {
	Token    lexer.Token // The = or +=, -=, etc. token
	Left     Expr        // Variable, property, or array access
	Operator string      // =, +=, -=, *=, etc.
	Right    Expr
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return "(" + ae.Left.String() + " " + ae.Operator + " " + ae.Right.String() + ")"
}

// TernaryExpression represents a ternary conditional (? :)
type TernaryExpression struct {
	Token       lexer.Token // The ? token
	Condition   Expr
	Consequence Expr // Can be nil for short ternary (?:)
	Alternative Expr
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	if te.Consequence == nil {
		return "(" + te.Condition.String() + " ?: " + te.Alternative.String() + ")"
	}
	return "(" + te.Condition.String() + " ? " + te.Consequence.String() + " : " + te.Alternative.String() + ")"
}

// ArrayExpression represents an array literal [key => value, ...]
type ArrayExpression struct {
	Token    lexer.Token // The [ token
	Elements []ArrayElement
}

type ArrayElement struct {
	Key   Expr // nil for non-associative elements
	Value Expr
}

func (ae *ArrayExpression) expressionNode()      {}
func (ae *ArrayExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ArrayExpression) String() string {
	return "[array]"
}

// IndexExpression represents array/string access $arr[$index]
type IndexExpression struct {
	Token lexer.Token // The [ token
	Left  Expr        // The array or string
	Index Expr
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

// PropertyExpression represents property access $obj->prop
type PropertyExpression struct {
	Token    lexer.Token // The -> token
	Object   Expr
	Property Expr // Can be Identifier or dynamic expression
}

func (pe *PropertyExpression) expressionNode()      {}
func (pe *PropertyExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PropertyExpression) String() string {
	return "(" + pe.Object.String() + "->" + pe.Property.String() + ")"
}

// NullsafePropertyExpression represents nullsafe property access $obj?->prop
type NullsafePropertyExpression struct {
	Token    lexer.Token // The ?-> token
	Object   Expr
	Property Expr
}

func (npe *NullsafePropertyExpression) expressionNode()      {}
func (npe *NullsafePropertyExpression) TokenLiteral() string { return npe.Token.Literal }
func (npe *NullsafePropertyExpression) String() string {
	return "(" + npe.Object.String() + "?->" + npe.Property.String() + ")"
}

// StaticPropertyExpression represents static property access Class::$prop
type StaticPropertyExpression struct {
	Token    lexer.Token // The :: token
	Class    Expr        // Class name or expression
	Property Expr
}

func (spe *StaticPropertyExpression) expressionNode()      {}
func (spe *StaticPropertyExpression) TokenLiteral() string { return spe.Token.Literal }
func (spe *StaticPropertyExpression) String() string {
	return "(" + spe.Class.String() + "::" + spe.Property.String() + ")"
}

// CallExpression represents a function call func($args)
type CallExpression struct {
	Token     lexer.Token // The ( token
	Function  Expr        // Identifier, method call, or closure
	Arguments []Expr
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	return ce.Function.String() + "(...)"
}

// MethodCallExpression represents a method call $obj->method($args)
type MethodCallExpression struct {
	Token     lexer.Token // The -> token
	Object    Expr
	Method    Expr // Can be Identifier or dynamic expression
	Arguments []Expr
}

func (mce *MethodCallExpression) expressionNode()      {}
func (mce *MethodCallExpression) TokenLiteral() string { return mce.Token.Literal }
func (mce *MethodCallExpression) String() string {
	return mce.Object.String() + "->" + mce.Method.String() + "(...)"
}

// StaticCallExpression represents a static method call Class::method($args)
type StaticCallExpression struct {
	Token     lexer.Token // The :: token
	Class     Expr        // Class name or expression (self, parent, static)
	Method    Expr
	Arguments []Expr
}

func (sce *StaticCallExpression) expressionNode()      {}
func (sce *StaticCallExpression) TokenLiteral() string { return sce.Token.Literal }
func (sce *StaticCallExpression) String() string {
	return sce.Class.String() + "::" + sce.Method.String() + "(...)"
}

// NewExpression represents object instantiation new Class($args)
type NewExpression struct {
	Token     lexer.Token // The NEW token
	Class     Expr        // Class name or expression
	Arguments []Expr
}

func (ne *NewExpression) expressionNode()      {}
func (ne *NewExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NewExpression) String() string {
	return "new " + ne.Class.String() + "(...)"
}

// InstanceofExpression represents instanceof check
type InstanceofExpression struct {
	Token lexer.Token // The INSTANCEOF token
	Left  Expr
	Right Expr // Class name or expression
}

func (ie *InstanceofExpression) expressionNode()      {}
func (ie *InstanceofExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InstanceofExpression) String() string {
	return "(" + ie.Left.String() + " instanceof " + ie.Right.String() + ")"
}

// CastExpression represents type casting (int)$var
type CastExpression struct {
	Token lexer.Token // The opening ( token
	Type  string      // int, string, bool, etc.
	Expr  Expr
}

func (ce *CastExpression) expressionNode()      {}
func (ce *CastExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CastExpression) String() string {
	return "((" + ce.Type + ")" + ce.Expr.String() + ")"
}

// GroupedExpression represents an expression in parentheses
type GroupedExpression struct {
	Token lexer.Token // The ( token
	Expr  Expr
}

func (ge *GroupedExpression) expressionNode()      {}
func (ge *GroupedExpression) TokenLiteral() string { return ge.Token.Literal }
func (ge *GroupedExpression) String() string {
	return "(" + ge.Expr.String() + ")"
}

// Additional node types will be added in Tasks 1.7-1.10
