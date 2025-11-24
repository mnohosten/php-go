package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestClassNameExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple class name",
			input:    `<?php $x = MyClass::class;`,
			expected: "MyClass",
		},
		{
			name:     "self keyword",
			input:    `<?php $x = self::class;`,
			expected: "self",
		},
		{
			name:     "parent keyword",
			input:    `<?php $x = parent::class;`,
			expected: "parent",
		},
		{
			name:     "static keyword",
			input:    `<?php $x = static::class;`,
			expected: "static",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			if len(program.Statements) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
			}

			exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
			}

			assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression)
			if !ok {
				t.Fatalf("Expected AssignmentExpression, got %T", exprStmt.Expression)
			}

			classNameExpr, ok := assignExpr.Right.(*ast.ClassNameExpression)
			if !ok {
				t.Fatalf("Expected ClassNameExpression, got %T", assignExpr.Right)
			}

			// Check the class name
			identExpr, ok := classNameExpr.Class.(*ast.Identifier)
			if !ok {
				t.Fatalf("Expected Identifier for class, got %T", classNameExpr.Class)
			}

			if identExpr.Value != tt.expected {
				t.Errorf("Expected class name %q, got %q", tt.expected, identExpr.Value)
			}
		})
	}
}
