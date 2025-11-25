package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// TestSimpleNamespaceDeclaration tests basic namespace declaration
func TestSimpleNamespaceDeclaration(t *testing.T) {
	input := `<?php
	namespace MyNamespace;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Name == nil {
		t.Fatal("Namespace name is nil")
	}

	if len(stmt.Name.Parts) != 1 {
		t.Fatalf("Expected 1 namespace part, got %d", len(stmt.Name.Parts))
	}

	if stmt.Name.Parts[0] != "MyNamespace" {
		t.Errorf("Expected namespace part 'MyNamespace', got '%s'", stmt.Name.Parts[0])
	}
}

// TestNestedNamespaceDeclaration tests namespace with multiple parts
func TestNestedNamespaceDeclaration(t *testing.T) {
	input := `<?php
	namespace Vendor\Package\SubPackage;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Name == nil {
		t.Fatal("Namespace name is nil")
	}

	expectedParts := []string{"Vendor", "Package", "SubPackage"}
	if len(stmt.Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d namespace parts, got %d", len(expectedParts), len(stmt.Name.Parts))
	}

	for i, expected := range expectedParts {
		if stmt.Name.Parts[i] != expected {
			t.Errorf("Part %d: expected '%s', got '%s'", i, expected, stmt.Name.Parts[i])
		}
	}
}

// TestBracketedNamespace tests namespace with bracketed syntax
func TestBracketedNamespace(t *testing.T) {
	input := `<?php
	namespace MyNamespace {
		function myFunc() {
			return 42;
		}
	}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Name == nil {
		t.Fatal("Namespace name is nil")
	}

	if stmt.Name.Parts[0] != "MyNamespace" {
		t.Errorf("Expected namespace 'MyNamespace', got '%s'", stmt.Name.Parts[0])
	}

	if stmt.Body == nil {
		t.Fatal("Namespace body is nil")
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in namespace body, got %d", len(stmt.Body.Statements))
	}

	// Check that the function is in the body
	_, ok = stmt.Body.Statements[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Errorf("Expected FunctionDeclaration in body, got %T", stmt.Body.Statements[0])
	}
}

// TestGlobalNamespace tests global namespace syntax
func TestGlobalNamespace(t *testing.T) {
	input := `<?php
	namespace {
		$x = 5;
	}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Name != nil {
		t.Error("Global namespace should have nil name")
	}

	if stmt.Body == nil {
		t.Fatal("Namespace body is nil")
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(stmt.Body.Statements))
	}
}

// TestUnbracketedNamespaceWithStatements tests namespace with statements in unbracketed syntax
func TestUnbracketedNamespaceWithStatements(t *testing.T) {
	input := `<?php
	namespace MyApp;

	class MyClass {
		public $value;
	}

	function myFunction() {
		return 'hello';
	}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement (namespace), got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Name == nil {
		t.Fatal("Namespace name is nil")
	}

	if stmt.Name.Parts[0] != "MyApp" {
		t.Errorf("Expected namespace 'MyApp', got '%s'", stmt.Name.Parts[0])
	}

	// In unbracketed syntax, statements should be in stmt.Statements, not stmt.Body
	if stmt.Body != nil {
		t.Error("Unbracketed namespace should not have Body")
	}

	if len(stmt.Statements) != 2 {
		t.Fatalf("Expected 2 statements in namespace, got %d", len(stmt.Statements))
	}

	// First should be class
	_, ok = stmt.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Errorf("First statement should be ClassStatement, got %T", stmt.Statements[0])
	}

	// Second should be function
	_, ok = stmt.Statements[1].(*ast.FunctionDeclaration)
	if !ok {
		t.Errorf("Second statement should be FunctionStatement, got %T", stmt.Statements[1])
	}
}

// TestMultipleNamespaces tests multiple namespace declarations
func TestMultipleNamespaces(t *testing.T) {
	input := `<?php
	namespace First;
	class A {}

	namespace Second;
	class B {}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 namespace statements, got %d", len(program.Statements))
	}

	// Check first namespace
	stmt1, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("First statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt1.Name.Parts[0] != "First" {
		t.Errorf("Expected first namespace 'First', got '%s'", stmt1.Name.Parts[0])
	}

	if len(stmt1.Statements) != 1 {
		t.Fatalf("Expected 1 statement in first namespace, got %d", len(stmt1.Statements))
	}

	// Check second namespace
	stmt2, ok := program.Statements[1].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Second statement is not *ast.NamespaceStatement, got %T", program.Statements[1])
	}

	if stmt2.Name.Parts[0] != "Second" {
		t.Errorf("Expected second namespace 'Second', got '%s'", stmt2.Name.Parts[0])
	}

	if len(stmt2.Statements) != 1 {
		t.Fatalf("Expected 1 statement in second namespace, got %d", len(stmt2.Statements))
	}
}

// TestNamespaceNameParsing tests various namespace name formats
func TestNamespaceNameParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single part",
			input:    `<?php namespace App;`,
			expected: []string{"App"},
		},
		{
			name:     "two parts",
			input:    `<?php namespace App\Models;`,
			expected: []string{"App", "Models"},
		},
		{
			name:     "three parts",
			input:    `<?php namespace Vendor\Package\Component;`,
			expected: []string{"Vendor", "Package", "Component"},
		},
		{
			name:     "many parts",
			input:    `<?php namespace A\B\C\D\E;`,
			expected: []string{"A", "B", "C", "D", "E"},
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

			stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
			if !ok {
				t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
			}

			if stmt.Name == nil {
				t.Fatal("Namespace name is nil")
			}

			if len(stmt.Name.Parts) != len(tt.expected) {
				t.Fatalf("Expected %d parts, got %d", len(tt.expected), len(stmt.Name.Parts))
			}

			for i, expectedPart := range tt.expected {
				if stmt.Name.Parts[i] != expectedPart {
					t.Errorf("Part %d: expected '%s', got '%s'", i, expectedPart, stmt.Name.Parts[i])
				}
			}
		})
	}
}

// TestWordPressStyleNamespaces tests WordPress-like namespace usage
func TestWordPressStyleNamespaces(t *testing.T) {
	input := `<?php
	namespace WP\REST\V1;

	class EndpointController {
		public function register() {
			// ...
		}
	}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	expectedParts := []string{"WP", "REST", "V1"}
	if len(stmt.Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d namespace parts, got %d", len(expectedParts), len(stmt.Name.Parts))
	}

	for i, expected := range expectedParts {
		if stmt.Name.Parts[i] != expected {
			t.Errorf("Part %d: expected '%s', got '%s'", i, expected, stmt.Name.Parts[i])
		}
	}

	if len(stmt.Statements) != 1 {
		t.Fatalf("Expected 1 statement in namespace, got %d", len(stmt.Statements))
	}

	_, ok = stmt.Statements[0].(*ast.ClassDeclaration)
	if !ok {
		t.Errorf("Expected ClassStatement, got %T", stmt.Statements[0])
	}
}

// TestMixedBracketedAndUnbracketed tests mixing both namespace syntaxes (not in same file)
func TestBracketedNamespaceWithMultipleClasses(t *testing.T) {
	input := `<?php
	namespace MyNamespace {
		class ClassA {}
		class ClassB {}
		function helper() {}
	}
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.NamespaceStatement)
	if !ok {
		t.Fatalf("Statement is not *ast.NamespaceStatement, got %T", program.Statements[0])
	}

	if stmt.Body == nil {
		t.Fatal("Namespace body is nil")
	}

	if len(stmt.Body.Statements) != 3 {
		t.Fatalf("Expected 3 statements in namespace body, got %d", len(stmt.Body.Statements))
	}

	// Check statement types
	_, ok1 := stmt.Body.Statements[0].(*ast.ClassDeclaration)
	_, ok2 := stmt.Body.Statements[1].(*ast.ClassDeclaration)
	_, ok3 := stmt.Body.Statements[2].(*ast.FunctionDeclaration)

	if !ok1 || !ok2 || !ok3 {
		t.Error("Expected 2 ClassStatements and 1 FunctionStatement in body")
	}
}
