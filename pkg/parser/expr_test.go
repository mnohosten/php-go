package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestIntegerLiteralExpression(t *testing.T) {
	input := `<?php 5;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := `<?php 3.14;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 3.14 {
		t.Errorf("literal.Value not %f. got=%f", 3.14, literal.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"<?php !5;", "!", 5},
		{"<?php -15;", "-", 15},
		{"<?php +10;", "+", 10},
		{"<?php ~5;", "~", 5},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"<?php 5 + 5;", 5, "+", 5},
		{"<?php 5 - 5;", 5, "-", 5},
		{"<?php 5 * 5;", 5, "*", 5},
		{"<?php 5 / 5;", 5, "/", 5},
		{"<?php 5 % 5;", 5, "%", 5},
		{"<?php 5 ** 2;", 5, "**", 2},
		{"<?php 5 > 5;", 5, ">", 5},
		{"<?php 5 < 5;", 5, "<", 5},
		{"<?php 5 == 5;", 5, "==", 5},
		{"<?php 5 === 5;", 5, "===", 5},
		{"<?php 5 != 5;", 5, "!=", 5},
		{"<?php 5 !== 5;", 5, "!==", 5},
		{"<?php 5 && 5;", 5, "&&", 5},
		{"<?php 5 || 5;", 5, "||", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"<?php -a * b",
			"((-a) * b)",
		},
		{
			"<?php !-a",
			"(!(-a))",
		},
		{
			"<?php a + b + c",
			"((a + b) + c)",
		},
		{
			"<?php a + b - c",
			"((a + b) - c)",
		},
		{
			"<?php a * b * c",
			"((a * b) * c)",
		},
		{
			"<?php a * b / c",
			"((a * b) / c)",
		},
		{
			"<?php a + b / c",
			"(a + (b / c))",
		},
		{
			"<?php a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"<?php 3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"<?php 5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"<?php 5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"<?php 3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"<?php 2 ** 3 ** 2",
			"(2 ** (3 ** 2))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<?php true;", true},
		{"<?php false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expression)
		}

		if boolean.Value != tt.expected {
			t.Errorf("boolean.Value not %t. got=%t", tt.expected, boolean.Value)
		}
	}
}

func TestParsingVariableExpressions(t *testing.T) {
	input := `<?php $foo;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	variable, ok := stmt.Expression.(*ast.Variable)
	if !ok {
		t.Fatalf("exp not *ast.Variable. got=%T", stmt.Expression)
	}

	if variable.Name != "foo" {
		t.Errorf("variable.Name not %s. got=%s", "foo", variable.Name)
	}
}

func TestAssignmentExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"<?php $x = 5;", "="},
		{"<?php $x += 5;", "+="},
		{"<?php $x -= 5;", "-="},
		{"<?php $x *= 5;", "*="},
		{"<?php $x /= 5;", "/="},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d\n",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		assign, ok := stmt.Expression.(*ast.AssignmentExpression)
		if !ok {
			t.Fatalf("exp not *ast.AssignmentExpression. got=%T", stmt.Expression)
		}

		if assign.Operator != tt.operator {
			t.Errorf("assign.Operator not '%s'. got=%s", tt.operator, assign.Operator)
		}
	}
}

func TestTernaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<?php $a ? $b : $c", "($a ? $b : $c)"},
		{"<?php $a ?: $b", "($a ?: $b)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d\n",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		tern, ok := stmt.Expression.(*ast.TernaryExpression)
		if !ok {
			t.Fatalf("exp not *ast.TernaryExpression. got=%T", stmt.Expression)
		}

		if tern.String() != tt.expected {
			t.Errorf("ternary string wrong. expected=%s, got=%s", tt.expected, tern.String())
		}
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	input := `<?php [1, 2, 3];`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayExpression)
	if !ok {
		t.Fatalf("exp not *ast.ArrayExpression. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array.Elements does not contain 3 elements. got=%d", len(array.Elements))
	}
}

func TestArrayConstructorExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`<?php array();`, 0},
		{`<?php array(1, 2, 3);`, 3},
		{`<?php array('key' => 'value');`, 1},
		{`<?php array('a', 'b' => 'c', 'd');`, 3},
		{`<?php array(1, 2, 3,);`, 3}, // trailing comma
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		array, ok := stmt.Expression.(*ast.ArrayExpression)
		if !ok {
			t.Fatalf("exp not *ast.ArrayExpression. got=%T", stmt.Expression)
		}

		if len(array.Elements) != tt.expected {
			t.Fatalf("array.Elements does not contain %d elements. got=%d", tt.expected, len(array.Elements))
		}
	}
}

func TestArrayConstructorWithNestedArrays(t *testing.T) {
	input := `<?php array('nested' => array(1, 2, 3));`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayExpression)
	if !ok {
		t.Fatalf("exp not *ast.ArrayExpression. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 1 {
		t.Fatalf("array.Elements does not contain 1 element. got=%d", len(array.Elements))
	}

	// Check nested array
	nestedArray, ok := array.Elements[0].Value.(*ast.ArrayExpression)
	if !ok {
		t.Fatalf("nested value not *ast.ArrayExpression. got=%T", array.Elements[0].Value)
	}

	if len(nestedArray.Elements) != 3 {
		t.Fatalf("nested array.Elements does not contain 3 elements. got=%d", len(nestedArray.Elements))
	}
}

func TestIndexExpression(t *testing.T) {
	input := `<?php $myArray[1 + 1];`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testVariable(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestPropertyAccessExpression(t *testing.T) {
	input := `<?php $obj->prop;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	propExp, ok := stmt.Expression.(*ast.PropertyExpression)
	if !ok {
		t.Fatalf("exp not *ast.PropertyExpression. got=%T", stmt.Expression)
	}

	if !testVariable(t, propExp.Object, "obj") {
		return
	}
}

func TestMethodCallExpression(t *testing.T) {
	input := `<?php $obj->method(1, 2);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	methCall, ok := stmt.Expression.(*ast.MethodCallExpression)
	if !ok {
		t.Fatalf("exp not *ast.MethodCallExpression. got=%T", stmt.Expression)
	}

	if !testVariable(t, methCall.Object, "obj") {
		return
	}

	if len(methCall.Arguments) != 2 {
		t.Fatalf("wrong number of arguments. got=%d", len(methCall.Arguments))
	}
}

func TestCallExpression(t *testing.T) {
	input := `<?php add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testIntegerLiteral(t, exp.Arguments[0].Value, 1)
	testInfixExpression(t, exp.Arguments[1].Value, 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2].Value, 4, "+", 5)
}

func TestStaticCallExpression(t *testing.T) {
	input := `<?php MyClass::staticMethod(1, 2);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	staticCall, ok := stmt.Expression.(*ast.StaticCallExpression)
	if !ok {
		t.Fatalf("exp not *ast.StaticCallExpression. got=%T", stmt.Expression)
	}

	if len(staticCall.Arguments) != 2 {
		t.Fatalf("wrong number of arguments. got=%d", len(staticCall.Arguments))
	}
}

func TestNewExpression(t *testing.T) {
	input := `<?php new MyClass(1, 2);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	newExp, ok := stmt.Expression.(*ast.NewExpression)
	if !ok {
		t.Fatalf("exp not *ast.NewExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, newExp.Class, "MyClass") {
		return
	}

	if len(newExp.Arguments) != 2 {
		t.Fatalf("wrong number of arguments. got=%d", len(newExp.Arguments))
	}
}

func TestInstanceofExpression(t *testing.T) {
	input := `<?php $obj instanceof MyClass;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	instExp, ok := stmt.Expression.(*ast.InstanceofExpression)
	if !ok {
		t.Fatalf("exp not *ast.InstanceofExpression. got=%T", stmt.Expression)
	}

	if !testVariable(t, instExp.Left, "obj") {
		return
	}

	if !testIdentifier(t, instExp.Right, "MyClass") {
		return
	}
}

func TestCastExpression(t *testing.T) {
	tests := []struct {
		input string
		typ   string
	}{
		{"<?php (int)$x;", "int"},
		{"<?php (string)$y;", "string"},
		{"<?php (bool)$z;", "bool"},
		{"<?php (float)$a;", "float"},
		{"<?php (array)$b;", "array"},
		{"<?php (object)$c;", "object"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d\n",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		cast, ok := stmt.Expression.(*ast.CastExpression)
		if !ok {
			t.Fatalf("exp not *ast.CastExpression. got=%T", stmt.Expression)
		}

		if cast.Type != tt.typ {
			t.Errorf("cast.Type not '%s'. got=%s", tt.typ, cast.Type)
		}
	}
}

func TestGroupedExpression(t *testing.T) {
	input := `<?php (5 + 5) * 2;`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not *ast.InfixExpression. got=%T", stmt.Expression)
	}

	if exp.Operator != "*" {
		t.Fatalf("exp.Operator is not '*'. got=%s", exp.Operator)
	}

	grouped, ok := exp.Left.(*ast.GroupedExpression)
	if !ok {
		t.Fatalf("exp.Left is not *ast.GroupedExpression. got=%T", exp.Left)
	}

	if !testInfixExpression(t, grouped.Expr, 5, "+", 5) {
		return
	}
}

// Helper functions for testing

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expr, value interface{}) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	intValue, ok := value.(int)
	if !ok {
		t.Errorf("value not int. got=%T", value)
		return false
	}

	if integ.Value != int64(intValue) {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testVariable(t *testing.T, exp ast.Expr, name string) bool {
	variable, ok := exp.(*ast.Variable)
	if !ok {
		t.Errorf("exp not *ast.Variable. got=%T", exp)
		return false
	}

	if variable.Name != name {
		t.Errorf("variable.Name not %s. got=%s", name, variable.Name)
		return false
	}

	return true
}

func TestListExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedLength int
		description    string
	}{
		{"<?php list($a, $b);", 2, "basic list with two elements"},
		{"<?php list($a, $b, $c);", 3, "list with three elements"},
		{"<?php list($x);", 1, "list with single element"},
		{"<?php list($a, , $c);", 3, "list with skipped middle element"},
		{"<?php list($a, , , $d);", 4, "list with multiple skipped elements"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.php")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("%s: program has wrong number of statements. got=%d",
				tt.description, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("%s: program.Statements[0] is not ast.ExpressionStatement. got=%T",
				tt.description, program.Statements[0])
		}

		listExpr, ok := stmt.Expression.(*ast.ListExpression)
		if !ok {
			t.Fatalf("%s: exp not *ast.ListExpression. got=%T",
				tt.description, stmt.Expression)
		}

		if len(listExpr.Elements) != tt.expectedLength {
			t.Errorf("%s: wrong number of elements. got=%d, want=%d",
				tt.description, len(listExpr.Elements), tt.expectedLength)
		}
	}
}

func TestListExpressionWithKeyedElements(t *testing.T) {
	// Test PHP 7.1+ keyed list syntax: list("a" => $a, "b" => $b)
	input := `<?php list("x" => $a, "y" => $b);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	listExpr, ok := stmt.Expression.(*ast.ListExpression)
	if !ok {
		t.Fatalf("exp not *ast.ListExpression. got=%T", stmt.Expression)
	}

	if len(listExpr.Elements) != 2 {
		t.Fatalf("wrong number of elements. got=%d", len(listExpr.Elements))
	}

	// Check first element has a key
	if listExpr.Elements[0].Key == nil {
		t.Errorf("first element should have a key")
	}

	// Check second element has a key
	if listExpr.Elements[1].Key == nil {
		t.Errorf("second element should have a key")
	}

	// Check values are variables
	if _, ok := listExpr.Elements[0].Value.(*ast.Variable); !ok {
		t.Errorf("first element value should be a variable. got=%T", listExpr.Elements[0].Value)
	}

	if _, ok := listExpr.Elements[1].Value.(*ast.Variable); !ok {
		t.Errorf("second element value should be a variable. got=%T", listExpr.Elements[1].Value)
	}
}

func TestListExpressionInAssignment(t *testing.T) {
	// Test list() in assignment context
	input := `<?php list($a, $b) = array(1, 2);`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("exp not *ast.AssignmentExpression. got=%T", stmt.Expression)
	}

	// Left side should be a list expression
	listExpr, ok := assignExpr.Left.(*ast.ListExpression)
	if !ok {
		t.Fatalf("assignment left side not *ast.ListExpression. got=%T", assignExpr.Left)
	}

	if len(listExpr.Elements) != 2 {
		t.Errorf("wrong number of list elements. got=%d", len(listExpr.Elements))
	}

	// Right side should be an array expression
	_, ok = assignExpr.Right.(*ast.ArrayExpression)
	if !ok {
		t.Fatalf("assignment right side not *ast.ArrayExpression. got=%T",
			assignExpr.Right)
	}
}

func TestTrailingCommaInFunctionCall(t *testing.T) {
	tests := []struct {
		input         string
		expectedArgs  int
		description   string
	}{
		{
			input:        `<?php foo(1, 2, 3,);`,
			expectedArgs: 3,
			description:  "trailing comma in positional arguments",
		},
		{
			input:        `<?php foo(1,);`,
			expectedArgs: 1,
			description:  "trailing comma with single argument",
		},
		{
			input:        `<?php foo(name: 'value',);`,
			expectedArgs: 1,
			description:  "trailing comma with named argument",
		},
		{
			input: `<?php foo(
				name: 'value',
				age: 30,
			);`,
			expectedArgs: 2,
			description:  "trailing comma in multi-line named arguments",
		},
		{
			input:        `<?php $obj->method(1, 2,);`,
			expectedArgs: 2,
			description:  "trailing comma in method call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			var args []*ast.Argument
			switch expr := stmt.Expression.(type) {
			case *ast.CallExpression:
				args = expr.Arguments
			case *ast.MethodCallExpression:
				args = expr.Arguments
			default:
				t.Fatalf("unexpected expression type. got=%T", stmt.Expression)
			}

			if len(args) != tt.expectedArgs {
				t.Errorf("wrong number of arguments. expected=%d, got=%d", tt.expectedArgs, len(args))
			}
		})
	}
}

func testInfixExpression(t *testing.T, exp ast.Expr, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	return true
}

func TestCloneExpression(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		objectType string // expected type of the object being cloned
	}{
		{
			name:       "clone simple variable",
			input:      `<?php clone $obj;`,
			objectType: "*ast.Variable",
		},
		{
			name:       "clone property access",
			input:      `<?php clone ($this->obj);`,
			objectType: "*ast.PropertyExpression",
		},
		{
			name:       "clone method call",
			input:      `<?php clone ($factory->createObject());`,
			objectType: "*ast.MethodCallExpression",
		},
		{
			name:       "clone new expression",
			input:      `<?php clone new MyClass();`,
			objectType: "*ast.NewExpression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			cloneExp, ok := stmt.Expression.(*ast.CloneExpression)
			if !ok {
				t.Fatalf("exp not *ast.CloneExpression. got=%T", stmt.Expression)
			}

			if cloneExp.Object == nil {
				t.Fatalf("cloneExp.Object is nil")
			}
		})
	}
}
