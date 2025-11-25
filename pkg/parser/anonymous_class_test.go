package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestAnonymousClassParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		hasExtends     bool
		extendsClass   string
		implementsCount int
		bodyCount      int
		argsCount      int
	}{
		{
			name:      "basic anonymous class",
			input:     `<?php $obj = new class {};`,
			bodyCount: 0,
		},
		{
			name:      "anonymous class with method",
			input:     `<?php $obj = new class { public function foo() {} };`,
			bodyCount: 1,
		},
		{
			name:      "anonymous class with constructor args",
			input:     `<?php $obj = new class("value") { private $data; };`,
			argsCount: 1,
			bodyCount: 1,
		},
		{
			name:         "anonymous class extends",
			input:        `<?php $obj = new class extends Base {};`,
			hasExtends:   true,
			extendsClass: "Base",
		},
		{
			name:            "anonymous class implements",
			input:           `<?php $obj = new class implements Iface {};`,
			implementsCount: 1,
		},
		{
			name:            "anonymous class implements multiple",
			input:           `<?php $obj = new class implements Iface1, Iface2 {};`,
			implementsCount: 2,
		},
		{
			name:            "anonymous class extends and implements",
			input:           `<?php $obj = new class extends Base implements Iface {};`,
			hasExtends:      true,
			extendsClass:    "Base",
			implementsCount: 1,
		},
		{
			name:            "anonymous class full syntax",
			input:           `<?php $obj = new class("arg1", "arg2") extends Base implements Iface1, Iface2 { public function test() {} };`,
			hasExtends:      true,
			extendsClass:    "Base",
			implementsCount: 2,
			argsCount:       2,
			bodyCount:       1,
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

			exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
			}

			assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression)
			if !ok {
				t.Fatalf("expected AssignmentExpression, got %T", exprStmt.Expression)
			}

			anonClass, ok := assignExpr.Right.(*ast.AnonymousClassExpression)
			if !ok {
				t.Fatalf("expected AnonymousClassExpression, got %T", assignExpr.Right)
			}

			// Check extends
			if tt.hasExtends {
				if anonClass.Extends == nil {
					t.Error("expected Extends to be set")
				} else if anonClass.Extends.Value != tt.extendsClass {
					t.Errorf("expected Extends to be %q, got %q", tt.extendsClass, anonClass.Extends.Value)
				}
			} else {
				if anonClass.Extends != nil {
					t.Errorf("expected Extends to be nil, got %v", anonClass.Extends)
				}
			}

			// Check implements
			if len(anonClass.Implements) != tt.implementsCount {
				t.Errorf("expected %d implements, got %d", tt.implementsCount, len(anonClass.Implements))
			}

			// Check body
			if tt.bodyCount > 0 && len(anonClass.Body) != tt.bodyCount {
				t.Errorf("expected %d body members, got %d", tt.bodyCount, len(anonClass.Body))
			}

			// Check constructor args
			if tt.argsCount > 0 && len(anonClass.ConstructorArgs) != tt.argsCount {
				t.Errorf("expected %d constructor args, got %d", tt.argsCount, len(anonClass.ConstructorArgs))
			}
		})
	}
}

func TestAnonymousClassAsArgument(t *testing.T) {
	input := `<?php
foo(new class {
    public function bar() { return 42; }
});`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected CallExpression, got %T", exprStmt.Expression)
	}

	if len(callExpr.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(callExpr.Arguments))
	}

	_, ok = callExpr.Arguments[0].Value.(*ast.AnonymousClassExpression)
	if !ok {
		t.Fatalf("expected AnonymousClassExpression as argument, got %T", callExpr.Arguments[0].Value)
	}
}

func TestAnonymousClassWithProperties(t *testing.T) {
	input := `<?php
$obj = new class {
    public int $id = 1;
    private string $name = "test";

    public function getId(): int {
        return $this->id;
    }
};`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("expected AssignmentExpression, got %T", exprStmt.Expression)
	}

	anonClass, ok := assignExpr.Right.(*ast.AnonymousClassExpression)
	if !ok {
		t.Fatalf("expected AnonymousClassExpression, got %T", assignExpr.Right)
	}

	// Should have 2 properties + 1 method = 3 body members
	if len(anonClass.Body) != 3 {
		t.Errorf("expected 3 body members, got %d", len(anonClass.Body))
	}
}
