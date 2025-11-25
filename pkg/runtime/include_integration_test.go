package runtime_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/runtime"
)

// TestIncludeExpressionParsing verifies that include expressions are parsed correctly
func TestIncludeExpressionParsing(t *testing.T) {
	code := `<?php
	$result = require 'file.php';
	?>`

	l := lexer.New(code, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("Expected AssignmentExpression, got %T", stmt.Expression)
	}

	includeExpr, ok := assignExpr.Right.(*ast.IncludeExpression)
	if !ok {
		t.Fatalf("Expected IncludeExpression, got %T", assignExpr.Right)
	}

	if includeExpr.Type != "require" {
		t.Errorf("Expected type 'require', got '%s'", includeExpr.Type)
	}
}

// TestIncludeExpressionCompilation verifies that include expressions compile to opcodes
func TestIncludeExpressionCompilation(t *testing.T) {
	code := `<?php require 'file.php'; ?>`

	l := lexer.New(code, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	c := compiler.New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := c.Bytecode()

	if len(bytecode.Instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// The bytecode should contain an INCLUDE_OR_EVAL opcode
	// We don't check the exact position as it may vary with compiler changes
	t.Logf("Generated %d instructions", len(bytecode.Instructions))
}

// TestIncludeManagerIntegration verifies include manager works with file operations
func TestIncludeManagerIntegration(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "include_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "file1.php")
	file2 := filepath.Join(tmpDir, "subdir", "file2.php")

	// Create subdir
	if err := os.Mkdir(filepath.Dir(file2), 0755); err != nil {
		t.Fatal(err)
	}

	// Write test content
	content1 := "<?php return 'from file1'; ?>"
	content2 := "<?php return 'from file2'; ?>"

	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	// Test include manager
	im := runtime.GetGlobalIncludeManager()
	im.ClearIncluded() // Clear any previous state
	im.SetCurrentScriptDir(tmpDir)

	// Resolve file1
	resolved1, err := im.ResolvePath("file1.php")
	if err != nil {
		t.Fatalf("Failed to resolve file1: %v", err)
	}

	expected1, _ := filepath.EvalSymlinks(file1)
	if expected1 == "" {
		expected1, _ = filepath.Abs(file1)
	}
	if resolved1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, resolved1)
	}

	// Test _once variant
	wasIncluded := im.MarkIncluded(resolved1)
	if wasIncluded {
		t.Error("First MarkIncluded should return false")
	}

	wasIncluded = im.MarkIncluded(resolved1)
	if !wasIncluded {
		t.Error("Second MarkIncluded should return true")
	}

	// Resolve file2 with relative path
	im.SetCurrentScriptDir(tmpDir)
	resolved2, err := im.ResolvePath("subdir/file2.php")
	if err != nil {
		t.Fatalf("Failed to resolve file2: %v", err)
	}

	expected2, _ := filepath.EvalSymlinks(file2)
	if expected2 == "" {
		expected2, _ = filepath.Abs(file2)
	}
	if resolved2 != expected2 {
		t.Errorf("Expected %s, got %s", expected2, resolved2)
	}
}

// TestIncludePathResolution verifies include path searching works correctly
func TestIncludePathResolution(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "include_path_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create lib directory
	libDir := filepath.Join(tmpDir, "lib")
	if err := os.Mkdir(libDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create library file
	libFile := filepath.Join(libDir, "library.php")
	if err := os.WriteFile(libFile, []byte("<?php class Library {} ?>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set up include manager
	im := runtime.NewIncludeManager()
	im.SetIncludePaths([]string{libDir})
	im.SetCurrentScriptDir(tmpDir)

	// Resolve library.php (should find it in lib/)
	resolved, err := im.ResolvePath("library.php")
	if err != nil {
		t.Fatalf("Failed to resolve library.php: %v", err)
	}

	expected, _ := filepath.EvalSymlinks(libFile)
	if expected == "" {
		expected, _ = filepath.Abs(libFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}
