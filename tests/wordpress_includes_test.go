package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/runtime"
)

// TestWordPressStyleRequireOnce tests require_once with relative paths
// This pattern is extremely common in WordPress:
// require_once __DIR__ . '/wp-load.php';
func TestWordPressStyleRequireOnce(t *testing.T) {
	// Create temporary directory structure mimicking WordPress
	tmpDir, err := os.MkdirTemp("", "wp_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create wp-load.php
	wpLoad := filepath.Join(tmpDir, "wp-load.php")
	wpLoadContent := `<?php
	define('WP_LOADED', true);
	?>`
	if err := os.WriteFile(wpLoad, []byte(wpLoadContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create main file that includes wp-load.php
	mainFile := filepath.Join(tmpDir, "index.php")
	mainContent := `<?php
	require_once __DIR__ . '/wp-load.php';
	?>`
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Parse the main file
	content, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(content), mainFile)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	// Compile
	c := compiler.New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := c.Bytecode()
	if len(bytecode.Instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Verify the include manager can resolve the path
	im := runtime.GetGlobalIncludeManager()
	im.ClearIncluded()
	im.SetCurrentScriptDir(tmpDir)

	resolved, err := im.ResolvePath("./wp-load.php")
	if err != nil {
		t.Fatalf("Failed to resolve wp-load.php: %v", err)
	}

	expected, _ := filepath.EvalSymlinks(wpLoad)
	if expected == "" {
		expected, _ = filepath.Abs(wpLoad)
	}

	if resolved != expected {
		t.Errorf("Path resolution failed: expected %s, got %s", expected, resolved)
	}
}

// TestMultipleIncludesWithOnceVariants tests that _once variants prevent duplicate includes
// WordPress uses this pattern extensively to prevent redefinition errors
func TestMultipleIncludesWithOnceVariants(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wp_once_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a library file
	libFile := filepath.Join(tmpDir, "functions.php")
	libContent := `<?php
	function my_function() {
		return 'Hello';
	}
	?>`
	if err := os.WriteFile(libFile, []byte(libContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test include manager _once behavior
	im := runtime.NewIncludeManager()
	im.SetCurrentScriptDir(tmpDir)

	// First include
	wasIncluded := im.MarkIncluded(libFile)
	if wasIncluded {
		t.Error("First MarkIncluded should return false (not previously included)")
	}

	// Second include (should be skipped)
	wasIncluded = im.MarkIncluded(libFile)
	if !wasIncluded {
		t.Error("Second MarkIncluded should return true (already included)")
	}

	// Verify file is tracked
	if !im.IsIncluded(libFile) {
		t.Error("File should be marked as included")
	}

	// Clear and test again
	im.ClearIncluded()
	if im.IsIncluded(libFile) {
		t.Error("After clear, file should not be marked as included")
	}
}

// TestNestedIncludes tests WordPress-style nested includes
// WordPress has deep include chains: index.php → wp-blog-header.php → wp-load.php → wp-config.php
func TestNestedIncludes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wp_nested_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested file structure
	// wp-config.php (deepest level)
	wpConfig := filepath.Join(tmpDir, "wp-config.php")
	wpConfigContent := `<?php
	define('DB_NAME', 'wordpress');
	?>`
	if err := os.WriteFile(wpConfig, []byte(wpConfigContent), 0644); err != nil {
		t.Fatal(err)
	}

	// wp-load.php (includes wp-config.php)
	wpLoad := filepath.Join(tmpDir, "wp-load.php")
	wpLoadContent := `<?php
	require_once __DIR__ . '/wp-config.php';
	define('WP_LOADED', true);
	?>`
	if err := os.WriteFile(wpLoad, []byte(wpLoadContent), 0644); err != nil {
		t.Fatal(err)
	}

	// wp-blog-header.php (includes wp-load.php)
	wpBlogHeader := filepath.Join(tmpDir, "wp-blog-header.php")
	wpBlogHeaderContent := `<?php
	require_once __DIR__ . '/wp-load.php';
	?>`
	if err := os.WriteFile(wpBlogHeader, []byte(wpBlogHeaderContent), 0644); err != nil {
		t.Fatal(err)
	}

	// index.php (top level, includes wp-blog-header.php)
	indexFile := filepath.Join(tmpDir, "index.php")
	indexContent := `<?php
	require_once __DIR__ . '/wp-blog-header.php';
	?>`
	if err := os.WriteFile(indexFile, []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Parse and compile each file to verify they all parse correctly
	files := []string{indexFile, wpBlogHeader, wpLoad, wpConfig}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", file, err)
		}

		l := lexer.New(string(content), file)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			t.Fatalf("Parser errors in %s: %v", file, p.Errors())
		}

		c := compiler.New()
		if err := c.Compile(program); err != nil {
			t.Fatalf("Compilation error in %s: %v", file, err)
		}
	}

	// Test include manager can track nested includes
	im := runtime.NewIncludeManager()
	im.SetCurrentScriptDir(tmpDir)

	// Simulate the include chain
	im.MarkIncluded(indexFile)
	im.MarkIncluded(wpBlogHeader)
	im.MarkIncluded(wpLoad)
	im.MarkIncluded(wpConfig)

	includedFiles := im.GetIncludedFiles()
	if len(includedFiles) != 4 {
		t.Errorf("Expected 4 included files, got %d", len(includedFiles))
	}
}

// TestIncludeWithRelativePaths tests various relative path patterns
func TestIncludeWithRelativePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wp_relative_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure:
	// tmpDir/
	//   wp-admin/
	//     admin.php
	//   wp-includes/
	//     functions.php
	//   wp-load.php

	wpAdmin := filepath.Join(tmpDir, "wp-admin")
	wpIncludes := filepath.Join(tmpDir, "wp-includes")

	if err := os.Mkdir(wpAdmin, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(wpIncludes, 0755); err != nil {
		t.Fatal(err)
	}

	// Create wp-load.php
	wpLoad := filepath.Join(tmpDir, "wp-load.php")
	if err := os.WriteFile(wpLoad, []byte("<?php ?>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create wp-includes/functions.php
	functions := filepath.Join(wpIncludes, "functions.php")
	if err := os.WriteFile(functions, []byte("<?php ?>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create wp-admin/admin.php that includes files from parent directories
	admin := filepath.Join(wpAdmin, "admin.php")
	adminContent := `<?php
	require_once __DIR__ . '/../wp-load.php';
	require_once __DIR__ . '/../wp-includes/functions.php';
	?>`
	if err := os.WriteFile(admin, []byte(adminContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Parse admin.php
	content, err := os.ReadFile(admin)
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(content), admin)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	// Compile
	c := compiler.New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	// Test include manager resolves parent directory references
	im := runtime.NewIncludeManager()
	im.SetCurrentScriptDir(wpAdmin)

	// Resolve ../wp-load.php from wp-admin/
	resolved, err := im.ResolvePath("../wp-load.php")
	if err != nil {
		t.Fatalf("Failed to resolve ../wp-load.php: %v", err)
	}

	expected, _ := filepath.EvalSymlinks(wpLoad)
	if expected == "" {
		expected, _ = filepath.Abs(wpLoad)
	}

	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}

	// Resolve ../wp-includes/functions.php from wp-admin/
	resolved, err = im.ResolvePath("../wp-includes/functions.php")
	if err != nil {
		t.Fatalf("Failed to resolve ../wp-includes/functions.php: %v", err)
	}

	expected, _ = filepath.EvalSymlinks(functions)
	if expected == "" {
		expected, _ = filepath.Abs(functions)
	}

	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

// TestReturnValuesFromIncludes tests that include expressions can be used as values
// This is a valid PHP pattern: $config = require 'config.php';
func TestReturnValuesFromIncludes(t *testing.T) {
	code := `<?php
	$config = require 'config.php';
	$data = include 'data.php';
	$result = require_once 'functions.php';
	?>`

	l := lexer.New(code, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	// Compile to verify it generates correct opcodes
	c := compiler.New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := c.Bytecode()
	if len(bytecode.Instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Note: Full execution would require breaking the VM/compiler circular dependency
	// For now, we verify that parsing and compilation work correctly
}

// TestIncludePathSearching tests WordPress-style include path searching
// WordPress uses multiple include paths for plugins, themes, etc.
func TestIncludePathSearching(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wp_paths_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure:
	// tmpDir/
	//   wp-content/
	//     plugins/
	//       my-plugin/
	//         plugin.php
	//     themes/
	//       my-theme/
	//         functions.php

	wpContent := filepath.Join(tmpDir, "wp-content")
	plugins := filepath.Join(wpContent, "plugins", "my-plugin")
	themes := filepath.Join(wpContent, "themes", "my-theme")

	if err := os.MkdirAll(plugins, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(themes, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files
	pluginFile := filepath.Join(plugins, "plugin.php")
	themeFile := filepath.Join(themes, "functions.php")

	if err := os.WriteFile(pluginFile, []byte("<?php ?>"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(themeFile, []byte("<?php ?>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set up include manager with multiple paths
	im := runtime.NewIncludeManager()
	im.SetIncludePaths([]string{
		plugins,
		themes,
	})

	// Should find plugin.php in plugins directory
	resolved, err := im.ResolvePath("plugin.php")
	if err != nil {
		t.Fatalf("Failed to resolve plugin.php: %v", err)
	}

	expected, _ := filepath.EvalSymlinks(pluginFile)
	if expected == "" {
		expected, _ = filepath.Abs(pluginFile)
	}

	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}

	// Should find functions.php in themes directory (second path)
	resolved, err = im.ResolvePath("functions.php")
	if err != nil {
		t.Fatalf("Failed to resolve functions.php: %v", err)
	}

	expected, _ = filepath.EvalSymlinks(themeFile)
	if expected == "" {
		expected, _ = filepath.Abs(themeFile)
	}

	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

// TestConditionalIncludes tests WordPress pattern of conditional includes
func TestConditionalIncludes(t *testing.T) {
	code := `<?php
	if (file_exists(__DIR__ . '/config.php')) {
		require_once __DIR__ . '/config.php';
	}

	if (!defined('WP_LOADED')) {
		require __DIR__ . '/wp-load.php';
	}
	?>`

	l := lexer.New(code, "test.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	// Compile to verify
	c := compiler.New()
	if err := c.Compile(program); err != nil {
		t.Fatalf("Compilation error: %v", err)
	}

	bytecode := c.Bytecode()
	if len(bytecode.Instructions) == 0 {
		t.Fatal("No instructions generated")
	}
}
