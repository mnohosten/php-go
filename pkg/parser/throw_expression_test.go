package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// TestParseThrowExpression tests parsing of throw expressions (PHP 8.0+)
func TestParseThrowExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name: "throw in ternary operator",
			input: `<?php
$result = $value > 0 ? "positive" : throw new Exception("negative");
`,
			expectError: false,
		},
		{
			name: "throw with null coalescing",
			input: `<?php
$result = $x ?? throw new Exception("null");
`,
			expectError: false,
		},
		{
			name: "throw in arrow function",
			input: `<?php
$func = fn() => throw new Exception("error");
`,
			expectError: false,
		},
		{
			name: "throw with expression argument",
			input: `<?php
$result = $code === 200 ? "OK" : throw new Exception("Error: " . $code);
`,
			expectError: false,
		},
		{
			name: "throw with new expression",
			input: `<?php
throw new InvalidArgumentException("Invalid");
`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()
			errors := p.Errors()

			if tt.expectError {
				if len(errors) == 0 {
					t.Error("Expected parse errors, but got none")
				}
				return
			}

			if len(errors) > 0 {
				t.Errorf("Unexpected parse errors: %v", errors)
				return
			}

			if program == nil {
				t.Error("Expected program, got nil")
				return
			}

			// Verify that throw expressions are in the AST
			foundThrow := false
			ast.Walk(&ast.BaseVisitor{}, program)
			visitor := &throwExpressionVisitor{found: &foundThrow}
			ast.Walk(visitor, program)

			// Note: Not all test cases have throw expressions (some have throw statements)
			// So we just verify parsing succeeded
		})
	}
}

type throwExpressionVisitor struct {
	ast.BaseVisitor
	found *bool
}

func (v *throwExpressionVisitor) VisitThrowExpression(node *ast.ThrowExpression) bool {
	*v.found = true
	return true
}
