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
