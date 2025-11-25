package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

func TestIncludeExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // expected Type field
	}{
		{
			name:     "require with quotes",
			input:    `<?php require 'file.php'; ?>`,
			expected: "require",
		},
		{
			name:     "require with parentheses",
			input:    `<?php require('file.php'); ?>`,
			expected: "require",
		},
		{
			name:     "require_once",
			input:    `<?php require_once 'config.php'; ?>`,
			expected: "require_once",
		},
		{
			name:     "include",
			input:    `<?php include 'header.php'; ?>`,
			expected: "include",
		},
		{
			name:     "include_once",
			input:    `<?php include_once 'functions.php'; ?>`,
			expected: "include_once",
		},
		{
			name:     "require with concatenation",
			input:    `<?php require __DIR__ . '/file.php'; ?>`,
			expected: "require",
		},
		{
			name:     "require in assignment",
			input:    `<?php $result = require 'file.php'; ?>`,
			expected: "require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input, "test.php")
			p := New(l)
			program := p.ParseProgram()

			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
			}

			stmt := program.Statements[0]

			// Check if it's an expression statement
			var includeExpr *ast.IncludeExpression

			exprStmt, ok := stmt.(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statement is not *ast.ExpressionStatement, got %T", stmt)
			}

			// The expression might be an IncludeExpression directly or wrapped in an assignment
			if assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression); ok {
				// Handle assignment case: $result = require 'file.php'
				includeExpr, ok = assignExpr.Right.(*ast.IncludeExpression)
				if !ok {
					t.Fatalf("Right side of assignment is not *ast.IncludeExpression, got %T", assignExpr.Right)
				}
			} else {
				// Direct include expression
				includeExpr, ok = exprStmt.Expression.(*ast.IncludeExpression)
				if !ok {
					t.Fatalf("Expression is not *ast.IncludeExpression, got %T", exprStmt.Expression)
				}
			}

			if includeExpr.Type != tt.expected {
				t.Errorf("Expected type '%s', got '%s'", tt.expected, includeExpr.Type)
			}

			if includeExpr.Path == nil {
				t.Fatalf("Path is nil")
			}
		})
	}
}

func TestIncludeExpressionWithoutParentheses(t *testing.T) {
	input := `<?php
	require 'file1.php';
	include 'file2.php';
	require_once 'file3.php';
	include_once 'file4.php';
	?>`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("Expected 4 statements, got %d", len(program.Statements))
	}

	expectedTypes := []string{"require", "include", "require_once", "include_once"}

	for i, expectedType := range expectedTypes {
		stmt := program.Statements[i]
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement %d is not *ast.ExpressionStatement, got %T", i, stmt)
		}

		includeExpr, ok := exprStmt.Expression.(*ast.IncludeExpression)
		if !ok {
			t.Fatalf("Expression %d is not *ast.IncludeExpression, got %T", i, exprStmt.Expression)
		}

		if includeExpr.Type != expectedType {
			t.Errorf("Statement %d: Expected type '%s', got '%s'", i, expectedType, includeExpr.Type)
		}
	}
}

func TestIncludeExpressionWithConcatenation(t *testing.T) {
	input := `<?php require __DIR__ . '/config.php'; ?>`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.ExpressionStatement, got %T", stmt)
	}

	includeExpr, ok := exprStmt.Expression.(*ast.IncludeExpression)
	if !ok {
		t.Fatalf("Expression is not *ast.IncludeExpression, got %T", exprStmt.Expression)
	}

	if includeExpr.Type != "require" {
		t.Errorf("Expected type 'require', got '%s'", includeExpr.Type)
	}

	// Path should be an InfixExpression (concatenation)
	_, ok = includeExpr.Path.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Path is not *ast.InfixExpression, got %T", includeExpr.Path)
	}
}

func TestIncludeExpressionAsValue(t *testing.T) {
	input := `<?php $config = require 'config.php'; ?>`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	// This should parse as an assignment where the right side is an IncludeExpression
	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.ExpressionStatement, got %T", stmt)
	}

	assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("Expression is not *ast.AssignmentExpression, got %T", exprStmt.Expression)
	}

	includeExpr, ok := assignExpr.Right.(*ast.IncludeExpression)
	if !ok {
		t.Fatalf("Right side is not *ast.IncludeExpression, got %T", assignExpr.Right)
	}

	if includeExpr.Type != "require" {
		t.Errorf("Expected type 'require', got '%s'", includeExpr.Type)
	}
}
