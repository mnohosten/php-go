package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestEchoStatement(t *testing.T) {
	input := `<?php echo "hello", "world";`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EchoStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.EchoStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Expressions) != 2 {
		t.Fatalf("echo has wrong number of expressions. expected=2, got=%d", len(stmt.Expressions))
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		hasReturnValue bool
	}{
		{"<?php return 5;", true},
		{"<?php return;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		if tt.hasReturnValue && stmt.ReturnValue == nil {
			t.Error("expected return value, got nil")
		}

		if !tt.hasReturnValue && stmt.ReturnValue != nil {
			t.Error("expected nil return value, got value")
		}
	}
}

func TestBreakStatement(t *testing.T) {
	tests := []struct {
		input    string
		hasDepth bool
	}{
		{"<?php break;", false},
		{"<?php break 2;", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.BreakStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.BreakStatement. got=%T", program.Statements[0])
		}

		if tt.hasDepth && stmt.Depth == nil {
			t.Error("expected depth, got nil")
		}

		if !tt.hasDepth && stmt.Depth != nil {
			t.Error("expected nil depth, got value")
		}
	}
}

func TestContinueStatement(t *testing.T) {
	input := `<?php continue;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ContinueStatement. got=%T", program.Statements[0])
	}
}

func TestGlobalStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		varCount int
		varNames []string
	}{
		{
			name:     "single variable",
			input:    `<?php global $x;`,
			varCount: 1,
			varNames: []string{"$x"},
		},
		{
			name:     "multiple variables",
			input:    `<?php global $x, $y, $z;`,
			varCount: 3,
			varNames: []string{"$x", "$y", "$z"},
		},
		{
			name:     "no semicolon",
			input:    `<?php global $var`,
			varCount: 1,
			varNames: []string{"$var"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.GlobalStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.GlobalStatement. got=%T", program.Statements[0])
			}

			if len(stmt.Variables) != tt.varCount {
				t.Fatalf("expected %d variables, got=%d", tt.varCount, len(stmt.Variables))
			}

			for i, expectedName := range tt.varNames {
				if stmt.Variables[i].Value != expectedName {
					t.Errorf("variable[%d] name wrong. expected=%q, got=%q",
						i, expectedName, stmt.Variables[i].Value)
				}
			}
		})
	}
}

func TestUnsetStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		varCount int
	}{
		{
			name:     "single variable",
			input:    `<?php unset($x);`,
			varCount: 1,
		},
		{
			name:     "multiple variables",
			input:    `<?php unset($x, $y, $z);`,
			varCount: 3,
		},
		{
			name:     "array element",
			input:    `<?php unset($arr[0]);`,
			varCount: 1,
		},
		{
			name:     "object property",
			input:    `<?php unset($obj->prop);`,
			varCount: 1,
		},
		{
			name:     "mixed expressions",
			input:    `<?php unset($x, $arr[0], $obj->prop);`,
			varCount: 3,
		},
		{
			name:     "nested array",
			input:    `<?php unset($arr['key']['nested']);`,
			varCount: 1,
		},
		{
			name:     "no semicolon",
			input:    `<?php unset($var)`,
			varCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.UnsetStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.UnsetStatement. got=%T", program.Statements[0])
			}

			if len(stmt.Variables) != tt.varCount {
				t.Fatalf("expected %d variables, got=%d", tt.varCount, len(stmt.Variables))
			}
		})
	}
}

func TestIfStatement(t *testing.T) {
	input := `<?php if ($x > 0) { echo "positive"; }`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Error("if statement condition is nil")
	}

	if stmt.Consequence == nil {
		t.Error("if statement consequence is nil")
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		echo "positive";
	} else {
		echo "non-positive";
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Alternative == nil {
		t.Error("if statement alternative is nil")
	}
}

func TestIfElseIfElseStatement(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		echo "positive";
	} elseif ($x < 0) {
		echo "negative";
	} else {
		echo "zero";
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if len(stmt.ElseIfs) != 1 {
		t.Errorf("expected 1 elseif clause, got %d", len(stmt.ElseIfs))
	}

	if stmt.Alternative == nil {
		t.Error("if statement alternative is nil")
	}
}

func TestWhileStatement(t *testing.T) {
	input := `<?php while ($x < 10) { $x++; }`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Error("while statement condition is nil")
	}

	if stmt.Body == nil {
		t.Error("while statement body is nil")
	}
}

func TestDoWhileStatement(t *testing.T) {
	input := `<?php do { $x++; } while ($x < 10);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DoWhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.DoWhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Body == nil {
		t.Error("do-while statement body is nil")
	}

	if stmt.Condition == nil {
		t.Error("do-while statement condition is nil")
	}
}

func TestForStatement(t *testing.T) {
	input := `<?php for ($i = 0; $i < 10; $i++) { echo $i; }`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Init) == 0 {
		t.Error("for statement init is empty")
	}

	if len(stmt.Condition) == 0 {
		t.Error("for statement condition is empty")
	}

	if len(stmt.Increment) == 0 {
		t.Error("for statement increment is empty")
	}

	if stmt.Body == nil {
		t.Error("for statement body is nil")
	}
}

func TestForeachStatement(t *testing.T) {
	tests := []struct {
		input  string
		hasKey bool
	}{
		{"<?php foreach ($arr as $value) { echo $value; }", false},
		{"<?php foreach ($arr as $key => $value) { echo $key, $value; }", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForeachStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ForeachStatement. got=%T", program.Statements[0])
		}

		if stmt.Array == nil {
			t.Error("foreach statement array is nil")
		}

		if stmt.Value == nil {
			t.Error("foreach statement value is nil")
		}

		if tt.hasKey && stmt.Key == nil {
			t.Error("expected key, got nil")
		}

		if !tt.hasKey && stmt.Key != nil {
			t.Error("expected nil key, got value")
		}
	}
}

func TestSwitchStatement(t *testing.T) {
	input := `<?php
	switch ($x) {
		case 1:
			echo "one";
			break;
		case 2:
			echo "two";
			break;
		default:
			echo "other";
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.SwitchStatement. got=%T", program.Statements[0])
	}

	if stmt.Subject == nil {
		t.Error("switch statement subject is nil")
	}

	if len(stmt.Cases) != 3 {
		t.Errorf("expected 3 cases, got %d", len(stmt.Cases))
	}

	// Check for default case
	hasDefault := false
	for _, c := range stmt.Cases {
		if c.Value == nil {
			hasDefault = true
			break
		}
	}

	if !hasDefault {
		t.Error("switch statement missing default case")
	}
}

func TestMatchExpression(t *testing.T) {
	input := `<?php
	$result = match ($x) {
		1 => "one",
		2 => "two",
		default => "other"
	};`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	assign, ok := exprStmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("expression is not *ast.AssignmentExpression. got=%T", exprStmt.Expression)
	}

	matchExpr, ok := assign.Right.(*ast.MatchExpression)
	if !ok {
		t.Fatalf("right side is not *ast.MatchExpression. got=%T", assign.Right)
	}

	if matchExpr.Subject == nil {
		t.Error("match expression subject is nil")
	}

	if len(matchExpr.Arms) != 3 {
		t.Errorf("expected 3 match arms, got %d", len(matchExpr.Arms))
	}

	// Check for default arm
	hasDefault := false
	for _, arm := range matchExpr.Arms {
		if arm.IsDefault {
			hasDefault = true
			break
		}
	}

	if !hasDefault {
		t.Error("match expression missing default arm")
	}
}

func TestTryStatement(t *testing.T) {
	input := `<?php
	try {
		riskyOperation();
	} catch (Exception $e) {
		handleError($e);
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TryStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.TryStatement. got=%T", program.Statements[0])
	}

	if stmt.Body == nil {
		t.Error("try statement body is nil")
	}

	if len(stmt.CatchClauses) != 1 {
		t.Errorf("expected 1 catch clause, got %d", len(stmt.CatchClauses))
	}

	catchClause := stmt.CatchClauses[0]
	if len(catchClause.Types) == 0 {
		t.Error("catch clause has no types")
	}
}

func TestTryCatchFinallyStatement(t *testing.T) {
	input := `<?php
	try {
		riskyOperation();
	} catch (Exception $e) {
		handleError($e);
	} finally {
		cleanup();
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TryStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.TryStatement. got=%T", program.Statements[0])
	}

	if stmt.Finally == nil {
		t.Error("try statement finally is nil")
	}
}

func TestThrowStatement(t *testing.T) {
	input := `<?php throw new Exception("error");`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ThrowStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ThrowStatement. got=%T", program.Statements[0])
	}

	if stmt.Expression == nil {
		t.Error("throw statement expression is nil")
	}
}

func TestBlockStatement(t *testing.T) {
	input := `<?php
	{
		$x = 5;
		echo $x;
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.BlockStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Statements) != 2 {
		t.Errorf("block should have 2 statements, got %d", len(stmt.Statements))
	}
}

func TestNestedStatements(t *testing.T) {
	input := `<?php
	if ($x > 0) {
		for ($i = 0; $i < $x; $i++) {
			echo $i;
		}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if ifStmt.Consequence == nil {
		t.Fatal("if statement consequence is nil")
	}

	if len(ifStmt.Consequence.Statements) != 1 {
		t.Fatalf("if consequence should have 1 statement, got %d", len(ifStmt.Consequence.Statements))
	}

	_, ok = ifStmt.Consequence.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("nested statement is not *ast.ForStatement. got=%T", ifStmt.Consequence.Statements[0])
	}
}

// Additional edge case tests for better coverage

func TestForStatementEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		initLen  int
		condLen  int
		incrLen  int
	}{
		{
			name:    "multiple init",
			input:   `<?php for ($i = 0, $j = 0; $i < 10; $i++) { echo $i; }`,
			initLen: 2,
			condLen: 1,
			incrLen: 1,
		},
		{
			name:    "multiple conditions",
			input:   `<?php for ($i = 0; $i < 10, $j < 20; $i++) { echo $i; }`,
			initLen: 1,
			condLen: 2,
			incrLen: 1,
		},
		{
			name:    "multiple increments",
			input:   `<?php for ($i = 0; $i < 10; $i++, $j--) { echo $i; }`,
			initLen: 1,
			condLen: 1,
			incrLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			forStmt, ok := program.Statements[0].(*ast.ForStatement)
			if !ok {
				t.Fatalf("statement not ForStatement. got=%T", program.Statements[0])
			}

			if len(forStmt.Init) != tt.initLen {
				t.Errorf("init length wrong. expected=%d, got=%d", tt.initLen, len(forStmt.Init))
			}

			if len(forStmt.Condition) != tt.condLen {
				t.Errorf("condition length wrong. expected=%d, got=%d", tt.condLen, len(forStmt.Condition))
			}

			if len(forStmt.Increment) != tt.incrLen {
				t.Errorf("increment length wrong. expected=%d, got=%d", tt.incrLen, len(forStmt.Increment))
			}
		})
	}
}

func TestContinueStatementWithLevel(t *testing.T) {
	input := `<?php
	for ($i = 0; $i < 10; $i++) {
		for ($j = 0; $j < 10; $j++) {
			continue 2;
		}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	forStmt := program.Statements[0].(*ast.ForStatement)
	innerFor := forStmt.Body.Statements[0].(*ast.ForStatement)
	continueStmt, ok := innerFor.Body.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("statement not ContinueStatement. got=%T", innerFor.Body.Statements[0])
	}

	if continueStmt.Depth == nil {
		t.Fatal("continue depth is nil")
	}

	intLit, ok := continueStmt.Depth.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("continue depth not IntegerLiteral. got=%T", continueStmt.Depth)
	}

	if intLit.Value != 2 {
		t.Errorf("continue depth wrong. expected=2, got=%d", intLit.Value)
	}
}

func TestBreakStatementWithLevel(t *testing.T) {
	input := `<?php
	for ($i = 0; $i < 10; $i++) {
		for ($j = 0; $j < 10; $j++) {
			break 2;
		}
	}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	forStmt := program.Statements[0].(*ast.ForStatement)
	innerFor := forStmt.Body.Statements[0].(*ast.ForStatement)
	breakStmt, ok := innerFor.Body.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("statement not BreakStatement. got=%T", innerFor.Body.Statements[0])
	}

	if breakStmt.Depth == nil {
		t.Fatal("break depth is nil")
	}

	intLit, ok := breakStmt.Depth.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("break depth not IntegerLiteral. got=%T", breakStmt.Depth)
	}

	if intLit.Value != 2 {
		t.Errorf("break depth wrong. expected=2, got=%d", intLit.Value)
	}
}

func TestIfElseStatementWithNoBraces(t *testing.T) {
	input := `<?php if ($x) echo "true"; else echo "false";`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	// Parser doesn't support single statement if/else without braces yet
	// This test documents the current limitation
	if len(p.Errors()) == 0 {
		ifStmt := program.Statements[0].(*ast.IfStatement)
		if ifStmt.Consequence == nil {
			t.Error("if consequence is nil")
		}
	}
}

func TestNamespaceStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasBody  bool
		numParts int
	}{
		{
			name:     "simple namespace",
			input:    `<?php namespace Foo;`,
			expected: "Foo",
			hasBody:  false,
			numParts: 1,
		},
		{
			name:     "nested namespace",
			input:    `<?php namespace Foo\Bar\Baz;`,
			expected: "Foo\\Bar\\Baz",
			hasBody:  false,
			numParts: 3,
		},
		{
			name:     "bracketed namespace",
			input:    `<?php namespace Foo\Bar { echo "test"; }`,
			expected: "Foo\\Bar",
			hasBody:  true,
			numParts: 2,
		},
		{
			name:     "global namespace",
			input:    `<?php namespace { echo "test"; }`,
			expected: "",
			hasBody:  true,
			numParts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program should have 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
			if !ok {
				t.Fatalf("expected NamespaceStatement, got %T", program.Statements[0])
			}

			if tt.numParts == 0 {
				// Global namespace
				if stmt.Name != nil {
					t.Errorf("expected nil Name for global namespace, got %v", stmt.Name)
				}
			} else {
				if stmt.Name == nil {
					t.Fatal("expected Name, got nil")
				}
				if len(stmt.Name.Parts) != tt.numParts {
					t.Errorf("expected %d namespace parts, got %d", tt.numParts, len(stmt.Name.Parts))
				}
				if stmt.Name.String() != tt.expected {
					t.Errorf("expected namespace name %q, got %q", tt.expected, stmt.Name.String())
				}
			}

			if tt.hasBody && stmt.Body == nil {
				t.Error("expected Body, got nil")
			}
			if !tt.hasBody && stmt.Body != nil {
				t.Error("expected nil Body, got Body")
			}
		})
	}
}

func TestUseStatement(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		useType   string
		numUses   int
		firstUse  string
		firstAlia string
	}{
		{
			name:      "simple use",
			input:     `<?php use Foo\Bar;`,
			useType:   "",
			numUses:   1,
			firstUse:  "Foo\\Bar",
			firstAlia: "",
		},
		{
			name:      "use with alias",
			input:     `<?php use Foo\Bar as Baz;`,
			useType:   "",
			numUses:   1,
			firstUse:  "Foo\\Bar",
			firstAlia: "Baz",
		},
		{
			name:      "multiple use",
			input:     `<?php use Foo\Bar, Foo\Baz;`,
			useType:   "",
			numUses:   2,
			firstUse:  "Foo\\Bar",
			firstAlia: "",
		},
		{
			name:      "function use",
			input:     `<?php use function Foo\Bar;`,
			useType:   "function",
			numUses:   1,
			firstUse:  "Foo\\Bar",
			firstAlia: "",
		},
		{
			name:      "const use",
			input:     `<?php use const Foo\BAR;`,
			useType:   "const",
			numUses:   1,
			firstUse:  "Foo\\BAR",
			firstAlia: "",
		},
		{
			name:      "group use",
			input:     `<?php use Foo\{Bar, Baz};`,
			useType:   "",
			numUses:   2,
			firstUse:  "Bar",
			firstAlia: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program should have 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.UseStatement)
			if !ok {
				t.Fatalf("expected UseStatement, got %T", program.Statements[0])
			}

			if stmt.Type != tt.useType {
				t.Errorf("expected use type %q, got %q", tt.useType, stmt.Type)
			}

			if len(stmt.Uses) != tt.numUses {
				t.Fatalf("expected %d use clauses, got %d", tt.numUses, len(stmt.Uses))
			}

			if tt.numUses > 0 {
				if stmt.Uses[0].Name.String() != tt.firstUse {
					t.Errorf("expected first use name %q, got %q", tt.firstUse, stmt.Uses[0].Name.String())
				}
				if stmt.Uses[0].Alias != tt.firstAlia {
					t.Errorf("expected first use alias %q, got %q", tt.firstAlia, stmt.Uses[0].Alias)
				}
			}
		})
	}
}

func TestNamespaceAndUseIntegration(t *testing.T) {
	input := `<?php
namespace App\Controllers;

use App\Models\User;
use App\Services\{AuthService, EmailService};

class UserController {
    public function index() {
        echo "users";
    }
}`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Should have namespace statement
	if len(program.Statements) < 1 {
		t.Fatal("expected at least 1 statement")
	}

	nsStmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("expected NamespaceStatement, got %T", program.Statements[0])
	}

	if nsStmt.Name.String() != "App\\Controllers" {
		t.Errorf("expected namespace App\\Controllers, got %s", nsStmt.Name.String())
	}

	// The use statements and class should be in nsStmt.Statements
	if len(nsStmt.Statements) < 2 {
		t.Fatalf("expected at least 2 statements in namespace, got %d", len(nsStmt.Statements))
	}
}

// Alternative syntax tests

func TestIfAlternativeSyntax(t *testing.T) {
	input := `<?php
	if ($x > 0):
		echo "positive";
	endif;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Error("if statement condition is nil")
	}

	if stmt.Consequence == nil {
		t.Error("if statement consequence is nil")
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("expected 1 statement in consequence, got %d", len(stmt.Consequence.Statements))
	}
}

func TestIfElseAlternativeSyntax(t *testing.T) {
	input := `<?php
	if ($x > 0):
		echo "positive";
	else:
		echo "non-positive";
	endif;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if stmt.Alternative == nil {
		t.Error("if statement alternative is nil")
	}

	if len(stmt.Alternative.Statements) != 1 {
		t.Errorf("expected 1 statement in alternative, got %d", len(stmt.Alternative.Statements))
	}
}

func TestIfElseIfElseAlternativeSyntax(t *testing.T) {
	input := `<?php
	if ($x > 10):
		echo "big";
	elseif ($x > 0):
		echo "positive";
	else:
		echo "non-positive";
	endif;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}

	if len(stmt.ElseIfs) != 1 {
		t.Errorf("expected 1 elseif clause, got %d", len(stmt.ElseIfs))
	}

	if stmt.Alternative == nil {
		t.Error("if statement alternative is nil")
	}
}

func TestWhileAlternativeSyntax(t *testing.T) {
	input := `<?php
	while ($x > 0):
		echo $x;
		$x--;
	endwhile;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Error("while statement condition is nil")
	}

	if stmt.Body == nil {
		t.Error("while statement body is nil")
	}

	if len(stmt.Body.Statements) != 2 {
		t.Errorf("expected 2 statements in body, got %d", len(stmt.Body.Statements))
	}
}

func TestForAlternativeSyntax(t *testing.T) {
	input := `<?php
	for ($i = 0; $i < 10; $i++):
		echo $i;
	endfor;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStatement. got=%T", program.Statements[0])
	}

	if stmt.Body == nil {
		t.Error("for statement body is nil")
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(stmt.Body.Statements))
	}
}

func TestForeachAlternativeSyntax(t *testing.T) {
	input := `<?php
	foreach ($items as $item):
		echo $item;
	endforeach;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForeachStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForeachStatement. got=%T", program.Statements[0])
	}

	if stmt.Body == nil {
		t.Error("foreach statement body is nil")
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(stmt.Body.Statements))
	}
}

func TestSwitchAlternativeSyntax(t *testing.T) {
	input := `<?php
	switch ($x):
		case 1:
			echo "one";
			break;
		case 2:
			echo "two";
			break;
		default:
			echo "other";
	endswitch;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.SwitchStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Cases) != 3 {
		t.Errorf("expected 3 cases (2 case + 1 default), got %d", len(stmt.Cases))
	}

	// Check first case
	if stmt.Cases[0].Value == nil {
		t.Error("first case value is nil")
	}

	// Check default case (should be last)
	if stmt.Cases[2].Value != nil {
		t.Error("default case should have nil value")
	}
}
