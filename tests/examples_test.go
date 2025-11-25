package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/vm"
)

// TestExampleFiles runs all example PHP files and ensures they execute without errors
func TestExampleFiles(t *testing.T) {
	// Get project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Define example directories to test
	exampleDirs := []string{
		filepath.Join(projectRoot, "examples", "basic"),
		filepath.Join(projectRoot, "examples", "parallel"),
		filepath.Join(projectRoot, "examples", "oop"),
		filepath.Join(projectRoot, "examples", "advanced"),
	}

	var totalTests, passedTests, failedTests int
	var failedFiles []string

	for _, dir := range exampleDirs {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Find all PHP files in directory
		files, err := filepath.Glob(filepath.Join(dir, "*.php"))
		if err != nil {
			t.Errorf("Failed to list files in %s: %v", dir, err)
			continue
		}

		for _, file := range files {
			totalTests++
			relPath, _ := filepath.Rel(projectRoot, file)

			t.Run(relPath, func(t *testing.T) {
				// Read file
				code, err := os.ReadFile(file)
				if err != nil {
					t.Fatalf("Failed to read file: %v", err)
				}

				// Execute the file
				output, err := executePhpCode(string(code))
				if err != nil {
					t.Errorf("Execution failed: %v", err)
					failedTests++
					failedFiles = append(failedFiles, relPath)
					return
				}

				// For now, just check that execution completed without error
				// Future: compare with expected output files
				if output == "" && strings.Contains(string(code), "echo") {
					t.Logf("Warning: No output from file that contains echo")
				}

				passedTests++
			})
		}
	}

	// Print summary
	t.Logf("\n========================================")
	t.Logf("Example Files Test Summary")
	t.Logf("========================================")
	t.Logf("Total:  %d", totalTests)
	t.Logf("Passed: %d (%.1f%%)", passedTests, float64(passedTests)/float64(totalTests)*100)
	t.Logf("Failed: %d (%.1f%%)", failedTests, float64(failedTests)/float64(totalTests)*100)

	if len(failedFiles) > 0 {
		t.Logf("\nFailed files:")
		for _, file := range failedFiles {
			t.Logf("  - %s", file)
		}
	}
}

// TestBasicExamples tests basic example files individually with expected output validation
func TestBasicExamples(t *testing.T) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	basicDir := filepath.Join(projectRoot, "examples", "basic")

	// List of example files to test
	exampleFiles := []string{
		"hello.php",
		"variables.php",
		"expressions.php",
		"control_flow.php",
		"functions.php",
		"arrays.php",
		"strings.php",
	}

	for _, filename := range exampleFiles {
		t.Run(filename, func(t *testing.T) {
			filePath := filepath.Join(basicDir, filename)
			expectedPath := filepath.Join(basicDir, strings.TrimSuffix(filename, ".php")+".expected")

			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Skipf("File not found: %s", filename)
				return
			}

			// Read and execute
			code, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			output, err := executePhpCode(string(code))
			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			// Check for expected output file
			if _, err := os.Stat(expectedPath); err == nil {
				// Expected file exists, compare output
				expectedBytes, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Logf("Warning: Failed to read expected file: %v", err)
					return
				}

				expected := string(expectedBytes)
				if output != expected {
					t.Errorf("Output mismatch for %s:\n=== EXPECTED ===\n%s\n=== ACTUAL ===\n%s\n=== END ===",
						filename, expected, output)
				} else {
					t.Logf("✓ Output matches expected for %s", filename)
				}
			} else {
				// No expected file, just verify execution succeeded
				t.Logf("No expected file found for %s, execution succeeded with output length: %d",
					filename, len(output))
			}
		})
	}
}

// TestPhpTestFiles tests files in tests/php directory
func TestPhpTestFiles(t *testing.T) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	testDirs := []string{
		filepath.Join(projectRoot, "tests", "php", "basic"),
		filepath.Join(projectRoot, "tests", "php", "stdlib"),
		filepath.Join(projectRoot, "tests", "php", "vm"),
	}

	var totalTests, passedTests, failedTests int

	for _, dir := range testDirs {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Find all PHP files
		files, err := filepath.Glob(filepath.Join(dir, "*.php"))
		if err != nil {
			t.Errorf("Failed to list files in %s: %v", dir, err)
			continue
		}

		for _, file := range files {
			totalTests++
			relPath, _ := filepath.Rel(projectRoot, file)

			t.Run(relPath, func(t *testing.T) {
				code, err := os.ReadFile(file)
				if err != nil {
					t.Fatalf("Failed to read file: %v", err)
				}

				_, err = executePhpCode(string(code))
				if err != nil {
					t.Errorf("Execution failed: %v", err)
					failedTests++
					return
				}

				passedTests++
			})
		}
	}

	if totalTests > 0 {
		t.Logf("\n========================================")
		t.Logf("PHP Test Files Summary")
		t.Logf("========================================")
		t.Logf("Total:  %d", totalTests)
		t.Logf("Passed: %d (%.1f%%)", passedTests, float64(passedTests)/float64(totalTests)*100)
		t.Logf("Failed: %d (%.1f%%)", failedTests, float64(failedTests)/float64(totalTests)*100)
	}
}

// executePhpCode compiles and executes PHP code, returning output
func executePhpCode(code string) (string, error) {
	// Lex
	l := lexer.New(code, "test.php")

	// Parse
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parse errors
	if len(p.Errors()) > 0 {
		return "", fmt.Errorf("parse errors: %v", p.Errors())
	}

	// Compile
	c := compiler.New()
	err := c.Compile(program)
	if err != nil {
		return "", fmt.Errorf("compilation error: %w", err)
	}

	// Get bytecode
	bytecode := c.Bytecode()

	// Execute
	machine := vm.New()
	machine.LoadConstants(bytecode.Constants)
	err = machine.ExecuteWithCVs(bytecode.Instructions, bytecode.NumCVs)
	if err != nil {
		return "", fmt.Errorf("runtime error: %w", err)
	}

	// Get output
	output := machine.GetOutput()
	return output, nil
}

// getProjectRoot finds the project root directory
func getProjectRoot() (string, error) {
	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

// BenchmarkExampleExecution benchmarks example file execution
func BenchmarkExampleExecution(b *testing.B) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		b.Fatalf("Failed to find project root: %v", err)
	}

	// Benchmark a simple example file
	helloFile := filepath.Join(projectRoot, "examples", "basic", "hello.php")
	code, err := os.ReadFile(helloFile)
	if err != nil {
		b.Skipf("hello.php not found")
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executePhpCode(string(code))
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}
