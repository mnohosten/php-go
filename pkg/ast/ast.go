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

// ClosureExpression represents an anonymous function (closure)
// Example: function($x) use ($y) { return $x + $y; }
type ClosureExpression struct {
	Token      lexer.Token   // The FUNCTION token
	Parameters []*Parameter  // Function parameters
	Use        []*UseClause  // Variables captured from parent scope
	ReturnType Expr          // Return type hint (can be nil)
	Body       *BlockStatement
	ByRef      bool          // Returns reference (&function)
	Static     bool          // static function() (cannot use $this)
}

func (ce *ClosureExpression) expressionNode()      {}
func (ce *ClosureExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ClosureExpression) String() string {
	return "function(...) { ... }"
}

// UseClause represents a variable captured in a closure's use clause
// Example: use ($x, &$y)
type UseClause struct {
	Variable *Variable
	ByRef    bool // Capture by reference (&$var)
}

// ArrowFunctionExpression represents an arrow function (PHP 7.4+)
// Example: fn($x) => $x * 2
type ArrowFunctionExpression struct {
	Token      lexer.Token  // The FN token
	Parameters []*Parameter // Function parameters
	ReturnType Expr         // Return type hint (can be nil)
	Body       Expr         // Single expression (not a block)
	ByRef      bool         // Returns reference
	Static     bool         // static fn()
}

func (af *ArrowFunctionExpression) expressionNode()      {}
func (af *ArrowFunctionExpression) TokenLiteral() string { return af.Token.Literal }
func (af *ArrowFunctionExpression) String() string {
	return "fn(...) => ..."
}

// EchoStatement represents echo statement
type EchoStatement struct {
	Token       lexer.Token // The ECHO token
	Expressions []Expr
}

func (es *EchoStatement) statementNode()       {}
func (es *EchoStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EchoStatement) String() string {
	return "echo ..."
}

// ReturnStatement represents return statement
type ReturnStatement struct {
	Token       lexer.Token // The RETURN token
	ReturnValue Expr        // Can be nil
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return "return " + rs.ReturnValue.String()
	}
	return "return"
}

// BreakStatement represents break statement
type BreakStatement struct {
	Token lexer.Token // The BREAK token
	Depth Expr        // Optional depth (break 2)
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string {
	return "break"
}

// ContinueStatement represents continue statement
type ContinueStatement struct {
	Token lexer.Token // The CONTINUE token
	Depth Expr        // Optional depth (continue 2)
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string {
	return "continue"
}

// IfStatement represents if/elseif/else statement
type IfStatement struct {
	Token       lexer.Token // The IF token
	Condition   Expr
	Consequence *BlockStatement
	ElseIfs     []*ElseIfClause
	Alternative *BlockStatement // Can be nil
}

type ElseIfClause struct {
	Token       lexer.Token // The ELSEIF token
	Condition   Expr
	Consequence *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	return "if (...) { ... }"
}

// WhileStatement represents while loop
type WhileStatement struct {
	Token     lexer.Token // The WHILE token
	Condition Expr
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	return "while (...) { ... }"
}

// DoWhileStatement represents do-while loop
type DoWhileStatement struct {
	Token     lexer.Token // The DO token
	Body      *BlockStatement
	Condition Expr
}

func (dws *DoWhileStatement) statementNode()       {}
func (dws *DoWhileStatement) TokenLiteral() string { return dws.Token.Literal }
func (dws *DoWhileStatement) String() string {
	return "do { ... } while (...)"
}

// ForStatement represents for loop
type ForStatement struct {
	Token      lexer.Token // The FOR token
	Init       []Expr      // Initialization expressions
	Condition  []Expr      // Condition expressions
	Increment  []Expr      // Increment expressions
	Body       *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	return "for (...; ...; ...) { ... }"
}

// ForeachStatement represents foreach loop
type ForeachStatement struct {
	Token     lexer.Token // The FOREACH token
	Array     Expr
	Key       Expr        // Can be nil
	Value     Expr
	ByRef     bool        // true if value is by reference
	Body      *BlockStatement
}

func (fs *ForeachStatement) statementNode()       {}
func (fs *ForeachStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForeachStatement) String() string {
	return "foreach (...) { ... }"
}

// SwitchStatement represents switch statement
type SwitchStatement struct {
	Token   lexer.Token // The SWITCH token
	Subject Expr
	Cases   []*SwitchCase
}

type SwitchCase struct {
	Token   lexer.Token // The CASE or DEFAULT token
	Value   Expr        // nil for default case
	Body    []Stmt
}

func (ss *SwitchStatement) statementNode()       {}
func (ss *SwitchStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SwitchStatement) String() string {
	return "switch (...) { ... }"
}

// MatchExpression represents match expression (PHP 8.0+)
type MatchExpression struct {
	Token lexer.Token // The MATCH token
	Subject Expr
	Arms  []*MatchArm
}

type MatchArm struct {
	Conditions []Expr // Multiple conditions separated by comma
	Body       Expr
	IsDefault  bool
}

func (me *MatchExpression) expressionNode()      {}
func (me *MatchExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MatchExpression) String() string {
	return "match (...) { ... }"
}

// TryStatement represents try-catch-finally statement
type TryStatement struct {
	Token        lexer.Token // The TRY token
	Body         *BlockStatement
	CatchClauses []*CatchClause
	Finally      *BlockStatement // Can be nil
}

type CatchClause struct {
	Token     lexer.Token // The CATCH token
	Types     []Expr      // Exception types (can be multiple with |)
	Variable  *Variable   // Variable to catch exception
	Body      *BlockStatement
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	return "try { ... } catch (...) { ... }"
}

// ThrowStatement represents throw statement
type ThrowStatement struct {
	Token      lexer.Token // The THROW token
	Expression Expr
}

func (ts *ThrowStatement) statementNode()       {}
func (ts *ThrowStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *ThrowStatement) String() string {
	return "throw ..."
}

// Task 1.8: Declaration node types

// Parameter represents a function/method parameter
type Parameter struct {
	Name         *Variable
	Type         Expr // Type hint (can be nil for untyped parameters)
	DefaultValue Expr // Default value (can be nil)
	ByRef        bool // Pass by reference (&$param)
	Variadic     bool // Variadic parameter (...$param)
}

// FunctionDeclaration represents a function declaration
type FunctionDeclaration struct {
	Token      lexer.Token // The FUNCTION token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType Expr // Return type hint (can be nil)
	Body       *BlockStatement
	ByRef      bool // Returns reference (&function)
}

func (fd *FunctionDeclaration) statementNode()       {}
func (fd *FunctionDeclaration) TokenLiteral() string { return fd.Token.Literal }
func (fd *FunctionDeclaration) String() string {
	return "function " + fd.Name.Value + "(...) { ... }"
}

// ClassDeclaration represents a class declaration
type ClassDeclaration struct {
	Token      lexer.Token // The CLASS token
	Name       *Identifier
	Extends    *Identifier // Parent class (can be nil)
	Implements []*Identifier
	Body       []Stmt // Properties, methods, constants, trait uses
	Modifiers  []string   // abstract, final
}

func (cd *ClassDeclaration) statementNode()       {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ClassDeclaration) String() string {
	return "class " + cd.Name.Value + " { ... }"
}

// PropertyDeclaration represents a class property
type PropertyDeclaration struct {
	Token        lexer.Token // The first token (visibility or VAR)
	Visibility   string      // public, protected, private
	Static       bool
	Readonly     bool
	Type         Expr        // Type hint (can be nil)
	Properties   []*PropertyItem
}

type PropertyItem struct {
	Name         *Variable
	DefaultValue Expr // Can be nil
}

func (pd *PropertyDeclaration) statementNode()       {}
func (pd *PropertyDeclaration) TokenLiteral() string { return pd.Token.Literal }
func (pd *PropertyDeclaration) String() string {
	return "property declaration"
}

// MethodDeclaration represents a class method
type MethodDeclaration struct {
	Token      lexer.Token // The FUNCTION token
	Visibility string      // public, protected, private
	Static     bool
	Abstract   bool
	Final      bool
	Name       *Identifier
	Parameters []*Parameter
	ReturnType Expr // Return type hint (can be nil)
	Body       *BlockStatement // nil for abstract methods
	ByRef      bool // Returns reference
}

func (md *MethodDeclaration) statementNode()       {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }
func (md *MethodDeclaration) String() string {
	return "method " + md.Name.Value + "(...)"
}

// InterfaceDeclaration represents an interface declaration
type InterfaceDeclaration struct {
	Token   lexer.Token // The INTERFACE token
	Name    *Identifier
	Extends []*Identifier // Interfaces can extend multiple interfaces
	Body    []*MethodSignature
}

type MethodSignature struct {
	Token      lexer.Token // The FUNCTION token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType Expr // Return type hint (can be nil)
	ByRef      bool
}

func (id *InterfaceDeclaration) statementNode()       {}
func (id *InterfaceDeclaration) TokenLiteral() string { return id.Token.Literal }
func (id *InterfaceDeclaration) String() string {
	return "interface " + id.Name.Value + " { ... }"
}

// TraitDeclaration represents a trait declaration
type TraitDeclaration struct {
	Token lexer.Token // The TRAIT token
	Name  *Identifier
	Body  []Stmt // Properties and methods
}

func (td *TraitDeclaration) statementNode()       {}
func (td *TraitDeclaration) TokenLiteral() string { return td.Token.Literal }
func (td *TraitDeclaration) String() string {
	return "trait " + td.Name.Value + " { ... }"
}

// ClassConstantDeclaration represents class constants
type ClassConstantDeclaration struct {
	Token      lexer.Token // The CONST token
	Visibility string      // public, protected, private (PHP 7.1+)
	Constants  []*ConstantItem
}

type ConstantItem struct {
	Name  *Identifier
	Value Expr
}

func (ccd *ClassConstantDeclaration) statementNode()       {}
func (ccd *ClassConstantDeclaration) TokenLiteral() string { return ccd.Token.Literal }
func (ccd *ClassConstantDeclaration) String() string {
	return "const declaration"
}

// TraitUse represents trait usage in a class
type TraitUse struct {
	Token   lexer.Token // The USE token
	Traits  []*Identifier
	Adaptations []TraitAdaptation // insteadof, as
}

type TraitAdaptation interface {
	Node
	traitAdaptationNode()
}

// TraitPrecedence represents trait method precedence (insteadof)
type TraitPrecedence struct {
	Token      lexer.Token // The INSTEADOF token
	TraitName  *Identifier
	MethodName *Identifier
	Instead    []*Identifier // Traits to use instead
}

func (tp *TraitPrecedence) traitAdaptationNode()      {}
func (tp *TraitPrecedence) TokenLiteral() string      { return tp.Token.Literal }
func (tp *TraitPrecedence) String() string            { return "trait precedence" }

// TraitAlias represents trait method aliasing
type TraitAlias struct {
	Token      lexer.Token // The AS token
	TraitName  *Identifier // Can be nil
	MethodName *Identifier
	Alias      *Identifier // New name
	Visibility string      // New visibility (can be empty)
}

func (ta *TraitAlias) traitAdaptationNode()      {}
func (ta *TraitAlias) TokenLiteral() string      { return ta.Token.Literal }
func (ta *TraitAlias) String() string            { return "trait alias" }

func (tu *TraitUse) statementNode()       {}
func (tu *TraitUse) TokenLiteral() string { return tu.Token.Literal }
func (tu *TraitUse) String() string {
	return "use traits"
}

// Task 1.9: Type expression node types

// NullableType represents a nullable type (?Type)
type NullableType struct {
	Token lexer.Token // The ? token
	Type  Expr
}

func (nt *NullableType) expressionNode()      {}
func (nt *NullableType) TokenLiteral() string { return nt.Token.Literal }
func (nt *NullableType) String() string {
	return "?" + nt.Type.String()
}

// UnionType represents a union type (Type1|Type2|Type3)
type UnionType struct {
	Token lexer.Token // The first type token
	Types []Expr
}

func (ut *UnionType) expressionNode()      {}
func (ut *UnionType) TokenLiteral() string { return ut.Token.Literal }
func (ut *UnionType) String() string {
	s := ""
	for i, t := range ut.Types {
		if i > 0 {
			s += "|"
		}
		s += t.String()
	}
	return s
}

// IntersectionType represents an intersection type (Type1&Type2)
type IntersectionType struct {
	Token lexer.Token // The first type token
	Types []Expr
}

func (it *IntersectionType) expressionNode()      {}
func (it *IntersectionType) TokenLiteral() string { return it.Token.Literal }
func (it *IntersectionType) String() string {
	s := ""
	for i, t := range it.Types {
		if i > 0 {
			s += "&"
		}
		s += t.String()
	}
	return s
}

// Additional node types will be added in Task 1.10
