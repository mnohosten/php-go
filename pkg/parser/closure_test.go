package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestParseClosure_Basic(t *testing.T) {
	input := `<?php
	$closure = function($x) {
		return $x * 2;
	};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	closure, ok := assign.Right.(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", assign.Right)
	}

	if len(closure.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(closure.Parameters))
	}

	if closure.Parameters[0].Name.Name != "x" {
		t.Errorf("Expected parameter x, got %s", closure.Parameters[0].Name.Name)
	}

	if closure.Body == nil {
		t.Error("Closure body is nil")
	}
}

func TestParseClosure_WithUseClause(t *testing.T) {
	input := `<?php
	$closure = function($x) use ($y, &$z) {
		return $x + $y + $z;
	};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	closure, ok := assign.Right.(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", assign.Right)
	}

	if len(closure.Use) != 2 {
		t.Fatalf("Expected 2 use variables, got %d", len(closure.Use))
	}

	if closure.Use[0].Variable.Name != "$y" {
		t.Errorf("Expected $y, got %s", closure.Use[0].Variable.Name)
	}

	if closure.Use[0].ByRef {
		t.Error("Expected $y to be by value, not by reference")
	}

	if closure.Use[1].Variable.Name != "$z" {
		t.Errorf("Expected $z, got %s", closure.Use[1].Variable.Name)
	}

	if !closure.Use[1].ByRef {
		t.Error("Expected $z to be by reference")
	}
}

func TestParseClosure_WithReturnType(t *testing.T) {
	input := `<?php
	$closure = function($x): int {
		return $x * 2;
	};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	closure, ok := assign.Right.(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", assign.Right)
	}

	if closure.ReturnType == nil {
		t.Fatal("Expected return type, got nil")
	}

	typeIdent, ok := closure.ReturnType.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier for return type, got %T", closure.ReturnType)
	}

	if typeIdent.Value != "int" {
		t.Errorf("Expected return type int, got %s", typeIdent.Value)
	}
}

func TestParseClosure_Static(t *testing.T) {
	input := `<?php
	$closure = static function($x) {
		return $x * 2;
	};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	closure, ok := assign.Right.(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", assign.Right)
	}

	if !closure.Static {
		t.Error("Expected static closure")
	}
}

func TestParseClosure_ByRef(t *testing.T) {
	input := `<?php
	$closure = function &($x) {
		return $x;
	};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	closure, ok := assign.Right.(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", assign.Right)
	}

	if !closure.ByRef {
		t.Error("Expected closure with reference return")
	}
}

func TestParseArrowFunction_Basic(t *testing.T) {
	input := `<?php
	$arrow = fn($x) => $x * 2;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	arrow, ok := assign.Right.(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", assign.Right)
	}

	if len(arrow.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(arrow.Parameters))
	}

	if arrow.Parameters[0].Name.Name != "x" {
		t.Errorf("Expected parameter x, got %s", arrow.Parameters[0].Name.Name)
	}

	if arrow.Body == nil {
		t.Error("Arrow function body is nil")
	}
}

func TestParseArrowFunction_WithReturnType(t *testing.T) {
	input := `<?php
	$arrow = fn($x): int => $x * 2;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	arrow, ok := assign.Right.(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", assign.Right)
	}

	if arrow.ReturnType == nil {
		t.Fatal("Expected return type, got nil")
	}

	typeIdent, ok := arrow.ReturnType.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier for return type, got %T", arrow.ReturnType)
	}

	if typeIdent.Value != "int" {
		t.Errorf("Expected return type int, got %s", typeIdent.Value)
	}
}

func TestParseArrowFunction_Static(t *testing.T) {
	input := `<?php
	$arrow = static fn($x) => $x * 2;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	arrow, ok := assign.Right.(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", assign.Right)
	}

	if !arrow.Static {
		t.Error("Expected static arrow function")
	}
}

func TestParseArrowFunction_ByRef(t *testing.T) {
	input := `<?php
	$arrow = fn &($x) => $x;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	arrow, ok := assign.Right.(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", assign.Right)
	}

	if !arrow.ByRef {
		t.Error("Expected arrow function with reference return")
	}
}

func TestParseArrowFunction_Complex(t *testing.T) {
	input := `<?php
	$arrow = fn($x, $y): float => $x * 2.5 + $y;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assign := stmt.Expression.(*ast.AssignmentExpression)
	arrow, ok := assign.Right.(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", assign.Right)
	}

	if len(arrow.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(arrow.Parameters))
	}

	// Check body is a complex expression
	_, ok = arrow.Body.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expected InfixExpression for body, got %T", arrow.Body)
	}
}

func TestParseClosureInCall(t *testing.T) {
	input := `<?php
	array_map(function($x) { return $x * 2; }, $array);
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Expected CallExpression, got %T", stmt.Expression)
	}

	if len(call.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(call.Arguments))
	}

	// First argument should be the closure
	closure, ok := call.Arguments[0].(*ast.ClosureExpression)
	if !ok {
		t.Fatalf("Expected ClosureExpression, got %T", call.Arguments[0])
	}

	if len(closure.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(closure.Parameters))
	}
}

func TestParseArrowFunctionInCall(t *testing.T) {
	input := `<?php
	array_map(fn($x) => $x * 2, $array);
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Expected CallExpression, got %T", stmt.Expression)
	}

	if len(call.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(call.Arguments))
	}

	// First argument should be the arrow function
	arrow, ok := call.Arguments[0].(*ast.ArrowFunctionExpression)
	if !ok {
		t.Fatalf("Expected ArrowFunctionExpression, got %T", call.Arguments[0])
	}

	if len(arrow.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(arrow.Parameters))
	}
}
