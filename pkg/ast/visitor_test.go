package ast

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
)

// TestCountingVisitor tests the visitor pattern by counting nodes
func TestCountingVisitor(t *testing.T) {
	// Create a simple AST
	program := &Program{
		Statements: []Stmt{
			&ExpressionStatement{
				Expression: &InfixExpression{
					Left: &IntegerLiteral{
						Token: lexer.Token{Type: lexer.INTEGER, Literal: "5"},
						Value: 5,
					},
					Operator: "+",
					Right: &IntegerLiteral{
						Token: lexer.Token{Type: lexer.INTEGER, Literal: "3"},
						Value: 3,
					},
				},
			},
			&ReturnStatement{
				ReturnValue: &Variable{
					Token: lexer.Token{Type: lexer.VARIABLE, Literal: "$x"},
					Name:  "x",
				},
			},
		},
	}

	// Create a visitor that counts nodes
	counter := &CountingVisitor{}
	for _, stmt := range program.Statements {
		Walk(counter, stmt)
	}

	// We should have visited:
	// 1 ExpressionStatement
	// 1 InfixExpression
	// 2 IntegerLiterals
	// 1 ReturnStatement
	// 1 Variable
	// Total: 6 nodes
	expectedCount := 6
	if counter.Count != expectedCount {
		t.Errorf("expected to visit %d nodes, got %d", expectedCount, counter.Count)
	}
}

// CountingVisitor counts all visited nodes
type CountingVisitor struct {
	BaseVisitor
	Count int
}

func (cv *CountingVisitor) VisitExpressionStatement(node *ExpressionStatement) bool {
	cv.Count++
	return true
}

func (cv *CountingVisitor) VisitInfixExpression(node *InfixExpression) bool {
	cv.Count++
	return true
}

func (cv *CountingVisitor) VisitIntegerLiteral(node *IntegerLiteral) bool {
	cv.Count++
	return true
}

func (cv *CountingVisitor) VisitReturnStatement(node *ReturnStatement) bool {
	cv.Count++
	return true
}

func (cv *CountingVisitor) VisitVariable(node *Variable) bool {
	cv.Count++
	return true
}

// TestVariableCollector tests collecting specific node types
func TestVariableCollector(t *testing.T) {
	// Create AST with multiple variables
	program := &Program{
		Statements: []Stmt{
			&ExpressionStatement{
				Expression: &AssignmentExpression{
					Left: &Variable{
						Token: lexer.Token{Type: lexer.VARIABLE, Literal: "$x"},
						Name:  "x",
					},
					Operator: "=",
					Right: &Variable{
						Token: lexer.Token{Type: lexer.VARIABLE, Literal: "$y"},
						Name:  "y",
					},
				},
			},
		},
	}

	// Collect all variables
	collector := &VariableCollector{Variables: []string{}}
	for _, stmt := range program.Statements {
		Walk(collector, stmt)
	}

	if len(collector.Variables) != 2 {
		t.Errorf("expected 2 variables, got %d", len(collector.Variables))
	}

	if collector.Variables[0] != "x" || collector.Variables[1] != "y" {
		t.Errorf("unexpected variable names: %v", collector.Variables)
	}
}

// VariableCollector collects all variable names in the AST
type VariableCollector struct {
	BaseVisitor
	Variables []string
}

func (vc *VariableCollector) VisitVariable(node *Variable) bool {
	vc.Variables = append(vc.Variables, node.Name)
	return true
}

// TestFunctionCollector tests collecting function declarations
func TestFunctionCollector(t *testing.T) {
	program := &Program{
		Statements: []Stmt{
			&FunctionDeclaration{
				Name: &Identifier{Value: "foo"},
				Body: &BlockStatement{},
			},
			&FunctionDeclaration{
				Name: &Identifier{Value: "bar"},
				Body: &BlockStatement{},
			},
		},
	}

	collector := &FunctionCollector{Functions: []string{}}
	for _, stmt := range program.Statements {
		Walk(collector, stmt)
	}

	if len(collector.Functions) != 2 {
		t.Errorf("expected 2 functions, got %d", len(collector.Functions))
	}

	if collector.Functions[0] != "foo" || collector.Functions[1] != "bar" {
		t.Errorf("unexpected function names: %v", collector.Functions)
	}
}

// FunctionCollector collects all function declaration names
type FunctionCollector struct {
	BaseVisitor
	Functions []string
}

func (fc *FunctionCollector) VisitFunctionDeclaration(node *FunctionDeclaration) bool {
	fc.Functions = append(fc.Functions, node.Name.Value)
	return true
}

// TestVisitorSkip tests skipping child nodes
func TestVisitorSkip(t *testing.T) {
	// Create AST with nested structure
	program := &Program{
		Statements: []Stmt{
			&IfStatement{
				Condition: &BooleanLiteral{Value: true},
				Consequence: &BlockStatement{
					Statements: []Stmt{
						&ExpressionStatement{
							Expression: &Variable{Name: "x"},
						},
					},
				},
			},
		},
	}

	// Visitor that skips BlockStatements
	skipper := &SkippingVisitor{}
	for _, stmt := range program.Statements {
		Walk(skipper, stmt)
	}

	// Should visit IfStatement and BooleanLiteral, but skip BlockStatement and Variable
	if skipper.VisitedIf != 1 {
		t.Errorf("expected to visit 1 if statement, got %d", skipper.VisitedIf)
	}

	if skipper.VisitedBool != 1 {
		t.Errorf("expected to visit 1 boolean, got %d", skipper.VisitedBool)
	}

	if skipper.VisitedVar != 0 {
		t.Errorf("expected to visit 0 variables (should be skipped), got %d", skipper.VisitedVar)
	}
}

// SkippingVisitor skips BlockStatements
type SkippingVisitor struct {
	BaseVisitor
	VisitedIf   int
	VisitedBool int
	VisitedVar  int
}

func (sv *SkippingVisitor) VisitIfStatement(node *IfStatement) bool {
	sv.VisitedIf++
	return true
}

func (sv *SkippingVisitor) VisitBlockStatement(node *BlockStatement) bool {
	// Return false to skip visiting children
	return false
}

func (sv *SkippingVisitor) VisitBooleanLiteral(node *BooleanLiteral) bool {
	sv.VisitedBool++
	return true
}

func (sv *SkippingVisitor) VisitVariable(node *Variable) bool {
	sv.VisitedVar++
	return true
}
