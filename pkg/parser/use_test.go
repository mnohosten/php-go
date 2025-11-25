package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/lexer"
)

// TestSimpleUseStatement tests basic use statement parsing
func TestSimpleUseStatement(t *testing.T) {
	input := `<?php
	use MyNamespace\MyClass;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if stmt.Type != "" {
		t.Errorf("Expected empty type for normal use, got %q", stmt.Type)
	}

	if len(stmt.Uses) != 1 {
		t.Fatalf("Expected 1 use import, got %d", len(stmt.Uses))
	}

	useImport := stmt.Uses[0]
	if useImport.Name == nil {
		t.Fatal("Expected Name to be non-nil")
	}

	expectedParts := []string{"MyNamespace", "MyClass"}
	if len(useImport.Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(useImport.Name.Parts))
	}

	for i, part := range useImport.Name.Parts {
		if part != expectedParts[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedParts[i], part)
		}
	}

	if useImport.Alias != "" {
		t.Errorf("Expected no alias, got %q", useImport.Alias)
	}
}

// TestUseStatementWithAlias tests use with alias
func TestUseStatementWithAlias(t *testing.T) {
	input := `<?php
	use MyNamespace\MyClass as MyAlias;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if len(stmt.Uses) != 1 {
		t.Fatalf("Expected 1 use import, got %d", len(stmt.Uses))
	}

	useImport := stmt.Uses[0]
	if useImport.Alias != "MyAlias" {
		t.Errorf("Expected alias 'MyAlias', got %q", useImport.Alias)
	}

	expectedParts := []string{"MyNamespace", "MyClass"}
	if len(useImport.Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(useImport.Name.Parts))
	}

	for i, part := range useImport.Name.Parts {
		if part != expectedParts[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedParts[i], part)
		}
	}
}

// TestMultipleUseStatements tests multiple use declarations separated by commas
func TestMultipleUseStatements(t *testing.T) {
	input := `<?php
	use Namespace\ClassA, Namespace\ClassB, Namespace\ClassC;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if len(stmt.Uses) != 3 {
		t.Fatalf("Expected 3 use imports, got %d", len(stmt.Uses))
	}

	expectedClasses := []string{"ClassA", "ClassB", "ClassC"}
	for i, useImport := range stmt.Uses {
		if len(useImport.Name.Parts) != 2 {
			t.Errorf("Import %d: expected 2 parts, got %d", i, len(useImport.Name.Parts))
			continue
		}

		if useImport.Name.Parts[0] != "Namespace" {
			t.Errorf("Import %d: expected namespace 'Namespace', got %q", i, useImport.Name.Parts[0])
		}

		if useImport.Name.Parts[1] != expectedClasses[i] {
			t.Errorf("Import %d: expected class %q, got %q", i, expectedClasses[i], useImport.Name.Parts[1])
		}
	}
}

// TestGroupUseStatement tests grouped use syntax
func TestGroupUseStatement(t *testing.T) {
	input := `<?php
	use Namespace\{ClassA, ClassB, ClassC};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if stmt.Prefix != "Namespace" {
		t.Errorf("Expected prefix 'Namespace', got %q", stmt.Prefix)
	}

	if len(stmt.Uses) != 3 {
		t.Fatalf("Expected 3 use imports, got %d", len(stmt.Uses))
	}

	expectedClasses := []string{"ClassA", "ClassB", "ClassC"}
	for i, useImport := range stmt.Uses {
		if len(useImport.Name.Parts) != 1 {
			t.Errorf("Import %d: expected 1 part (class name), got %d", i, len(useImport.Name.Parts))
			continue
		}

		if useImport.Name.Parts[0] != expectedClasses[i] {
			t.Errorf("Import %d: expected class %q, got %q", i, expectedClasses[i], useImport.Name.Parts[0])
		}
	}
}

// TestGroupUseStatementWithAlias tests grouped use with alias
func TestGroupUseStatementWithAlias(t *testing.T) {
	input := `<?php
	use Namespace\{ClassA as AliasA, ClassB, ClassC as AliasC};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if len(stmt.Uses) != 3 {
		t.Fatalf("Expected 3 use imports, got %d", len(stmt.Uses))
	}

	// Check first import with alias
	if stmt.Uses[0].Name.Parts[0] != "ClassA" {
		t.Errorf("Expected ClassA, got %q", stmt.Uses[0].Name.Parts[0])
	}
	if stmt.Uses[0].Alias != "AliasA" {
		t.Errorf("Expected alias 'AliasA', got %q", stmt.Uses[0].Alias)
	}

	// Check second import without alias
	if stmt.Uses[1].Name.Parts[0] != "ClassB" {
		t.Errorf("Expected ClassB, got %q", stmt.Uses[1].Name.Parts[0])
	}
	if stmt.Uses[1].Alias != "" {
		t.Errorf("Expected no alias, got %q", stmt.Uses[1].Alias)
	}

	// Check third import with alias
	if stmt.Uses[2].Name.Parts[0] != "ClassC" {
		t.Errorf("Expected ClassC, got %q", stmt.Uses[2].Name.Parts[0])
	}
	if stmt.Uses[2].Alias != "AliasC" {
		t.Errorf("Expected alias 'AliasC', got %q", stmt.Uses[2].Alias)
	}
}

// TestFunctionUseStatement tests use function syntax
func TestFunctionUseStatement(t *testing.T) {
	input := `<?php
	use function Namespace\myFunction;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if stmt.Type != "function" {
		t.Errorf("Expected type 'function', got %q", stmt.Type)
	}

	if len(stmt.Uses) != 1 {
		t.Fatalf("Expected 1 use import, got %d", len(stmt.Uses))
	}

	expectedParts := []string{"Namespace", "myFunction"}
	if len(stmt.Uses[0].Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(stmt.Uses[0].Name.Parts))
	}

	for i, part := range stmt.Uses[0].Name.Parts {
		if part != expectedParts[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedParts[i], part)
		}
	}
}

// TestConstUseStatement tests use const syntax
func TestConstUseStatement(t *testing.T) {
	input := `<?php
	use const Namespace\MY_CONSTANT;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	if stmt.Type != "const" {
		t.Errorf("Expected type 'const', got %q", stmt.Type)
	}

	if len(stmt.Uses) != 1 {
		t.Fatalf("Expected 1 use import, got %d", len(stmt.Uses))
	}

	expectedParts := []string{"Namespace", "MY_CONSTANT"}
	if len(stmt.Uses[0].Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(stmt.Uses[0].Name.Parts))
	}

	for i, part := range stmt.Uses[0].Name.Parts {
		if part != expectedParts[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedParts[i], part)
		}
	}
}

// TestNestedNamespaceUse tests use with deeply nested namespaces
func TestNestedNamespaceUse(t *testing.T) {
	input := `<?php
	use Vendor\Package\SubPackage\DeepPackage\MyClass;
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
	}

	expectedParts := []string{"Vendor", "Package", "SubPackage", "DeepPackage", "MyClass"}
	if len(stmt.Uses[0].Name.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(stmt.Uses[0].Name.Parts))
	}

	for i, part := range stmt.Uses[0].Name.Parts {
		if part != expectedParts[i] {
			t.Errorf("Part %d: expected %q, got %q", i, expectedParts[i], part)
		}
	}
}

// TestWordPressStyleUse tests WordPress-like use statements
func TestWordPressStyleUse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  string
		expectedParts [][]string
		expectedAlias []string
	}{
		{
			name: "WP REST API namespace",
			input: `<?php
			use WP\REST\V1\Controller;
			`,
			expectedType:  "",
			expectedParts: [][]string{{"WP", "REST", "V1", "Controller"}},
			expectedAlias: []string{""},
		},
		{
			name: "Multiple WP classes",
			input: `<?php
			use WP\Admin\Screen, WP\Admin\Menu, WP\Admin\Notice;
			`,
			expectedType: "",
			expectedParts: [][]string{
				{"WP", "Admin", "Screen"},
				{"WP", "Admin", "Menu"},
				{"WP", "Admin", "Notice"},
			},
			expectedAlias: []string{"", "", ""},
		},
		{
			name: "WP grouped use",
			input: `<?php
			use WP\Admin\{Screen, Menu, Notice};
			`,
			expectedType: "",
			expectedParts: [][]string{
				{"Screen"},
				{"Menu"},
				{"Notice"},
			},
			expectedAlias: []string{"", "", ""},
		},
		{
			name: "WP with aliases",
			input: `<?php
			use WP\REST\V1\Controller as RestController;
			`,
			expectedType:  "",
			expectedParts: [][]string{{"WP", "REST", "V1", "Controller"}},
			expectedAlias: []string{"RestController"},
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

			stmt, ok := program.Statements[0].(*ast.UseStatement)
			if !ok {
				t.Fatalf("Expected UseStatement, got %T", program.Statements[0])
			}

			if stmt.Type != tt.expectedType {
				t.Errorf("Expected type %q, got %q", tt.expectedType, stmt.Type)
			}

			if len(stmt.Uses) != len(tt.expectedParts) {
				t.Fatalf("Expected %d use imports, got %d", len(tt.expectedParts), len(stmt.Uses))
			}

			for i, useImport := range stmt.Uses {
				if len(useImport.Name.Parts) != len(tt.expectedParts[i]) {
					t.Errorf("Import %d: expected %d parts, got %d",
						i, len(tt.expectedParts[i]), len(useImport.Name.Parts))
					continue
				}

				for j, part := range useImport.Name.Parts {
					if part != tt.expectedParts[i][j] {
						t.Errorf("Import %d, part %d: expected %q, got %q",
							i, j, tt.expectedParts[i][j], part)
					}
				}

				if useImport.Alias != tt.expectedAlias[i] {
					t.Errorf("Import %d: expected alias %q, got %q",
						i, tt.expectedAlias[i], useImport.Alias)
				}
			}
		})
	}
}

// TestMixedUseStatements tests multiple use statements in one file
func TestMixedUseStatements(t *testing.T) {
	input := `<?php
	use Vendor\Package\ClassA;
	use function Vendor\Package\myFunc;
	use const Vendor\Package\MY_CONST;
	use Vendor\Other\{ClassB, ClassC as AliasC};
	`

	l := lexer.New(input, "test.php")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("Expected 4 statements, got %d", len(program.Statements))
	}

	// Check first use (normal class)
	stmt1, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Statement 0: expected UseStatement, got %T", program.Statements[0])
	}
	if stmt1.Type != "" {
		t.Errorf("Statement 0: expected empty type, got %q", stmt1.Type)
	}

	// Check second use (function)
	stmt2, ok := program.Statements[1].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Statement 1: expected UseStatement, got %T", program.Statements[1])
	}
	if stmt2.Type != "function" {
		t.Errorf("Statement 1: expected type 'function', got %q", stmt2.Type)
	}

	// Check third use (const)
	stmt3, ok := program.Statements[2].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Statement 2: expected UseStatement, got %T", program.Statements[2])
	}
	if stmt3.Type != "const" {
		t.Errorf("Statement 2: expected type 'const', got %q", stmt3.Type)
	}

	// Check fourth use (grouped)
	stmt4, ok := program.Statements[3].(*ast.UseStatement)
	if !ok {
		t.Fatalf("Statement 3: expected UseStatement, got %T", program.Statements[3])
	}
	if stmt4.Prefix == "" {
		t.Error("Statement 3: expected non-empty prefix for group use")
	}
	if len(stmt4.Uses) != 2 {
		t.Errorf("Statement 3: expected 2 use imports, got %d", len(stmt4.Uses))
	}
}
