package phptest

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewTestExecutor tests creating a test executor
func TestNewTestExecutor(t *testing.T) {
	executor := NewTestExecutor("")
	if executor == nil {
		t.Fatal("NewTestExecutor() returned nil")
	}

	if executor.PHPBinary != "php-go" {
		t.Errorf("PHPBinary = %s, want php-go", executor.PHPBinary)
	}

	if executor.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", executor.Timeout)
	}
}

// TestNewTestExecutor_CustomBinary tests creating executor with custom binary
func TestNewTestExecutor_CustomBinary(t *testing.T) {
	executor := NewTestExecutor("/usr/bin/php")
	if executor.PHPBinary != "/usr/bin/php" {
		t.Errorf("PHPBinary = %s, want /usr/bin/php", executor.PHPBinary)
	}
}

// TestTestStatus_String tests status string representation
func TestTestStatus_String(t *testing.T) {
	tests := []struct {
		status   TestStatus
		expected string
	}{
		{StatusPass, "PASS"},
		{StatusFail, "FAIL"},
		{StatusSkip, "SKIP"},
		{StatusError, "ERROR"},
		{StatusUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("status.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestExecute_SimplePass tests executing a passing test
func TestExecute_SimplePass(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Simple pass test",
		Code:        "<?php\necho 'Hello World';\n?>",
		Expect:      "Hello World",
	}

	result := executor.Execute(test)

	if result.Status != StatusPass {
		t.Errorf("Execute() status = %v, want StatusPass", result.Status)
		t.Logf("Actual output: %q", result.ActualOutput)
		t.Logf("Expected output: %q", result.ExpectedOutput)
		if result.ErrorMessage != "" {
			t.Logf("Error: %s", result.ErrorMessage)
		}
	}
}

// TestExecute_SimpleFail tests executing a failing test
func TestExecute_SimpleFail(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Simple fail test",
		Code:        "<?php\necho 'Hello World';\n?>",
		Expect:      "Goodbye World",
	}

	result := executor.Execute(test)

	if result.Status != StatusFail {
		t.Errorf("Execute() status = %v, want StatusFail", result.Status)
	}

	if result.ErrorMessage == "" {
		t.Error("Execute() should set ErrorMessage for failed test")
	}
}

// TestExecute_WithFormat tests executing test with format placeholders
func TestExecute_WithFormat(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Format test",
		Code:        "<?php\necho 123;\n?>",
		ExpectF:     "%d",
	}

	result := executor.Execute(test)

	if result.Status != StatusPass {
		t.Errorf("Execute() status = %v, want StatusPass", result.Status)
		t.Logf("Actual output: %q", result.ActualOutput)
	}
}

// TestExecute_WithRegex tests executing test with regex pattern
func TestExecute_WithRegex(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Regex test",
		Code:        "<?php\necho 'Hello';\n?>",
		ExpectRegex: "^[A-Z][a-z]+$",
	}

	result := executor.Execute(test)

	if result.Status != StatusPass {
		t.Errorf("Execute() status = %v, want StatusPass", result.Status)
		t.Logf("Actual output: %q", result.ActualOutput)
	}
}

// TestExecute_InvalidCode tests executing test with code retrieval error
func TestExecute_InvalidCode(t *testing.T) {
	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description:  "Invalid code test",
		FileExternal: "nonexistent.php",
		Expect:       "test",
	}

	result := executor.Execute(test)

	if result.Status != StatusError {
		t.Errorf("Execute() status = %v, want StatusError", result.Status)
	}

	if result.ErrorMessage == "" {
		t.Error("Execute() should set ErrorMessage for code error")
	}
}

// TestExecute_Timeout tests test execution timeout
func TestExecute_Timeout(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	executor.Timeout = 100 * time.Millisecond

	test := &PHPTest{
		Description: "Timeout test",
		Code:        "<?php\nsleep(10);\necho 'Done';\n?>",
		Expect:      "Done",
	}

	result := executor.Execute(test)

	// Should either error or fail due to timeout
	if result.Status != StatusError && result.Status != StatusFail {
		t.Errorf("Execute() status = %v, want StatusError or StatusFail", result.Status)
	}
}

// TestCheckSkipCondition_Skip tests skip condition evaluation
func TestCheckSkipCondition_Skip(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		SkipIf: "<?php\ndie('skip Test skipped');\n?>",
	}

	shouldSkip, reason := executor.checkSkipCondition(test)

	if !shouldSkip {
		t.Error("checkSkipCondition() = false, want true")
	}

	if reason == "" {
		t.Error("checkSkipCondition() should return skip reason")
	}
}

// TestCheckSkipCondition_NoSkip tests no skip condition
func TestCheckSkipCondition_NoSkip(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		SkipIf: "<?php\necho 'run';\n?>",
	}

	shouldSkip, _ := executor.checkSkipCondition(test)

	if shouldSkip {
		t.Error("checkSkipCondition() = true, want false")
	}
}

// TestExecute_WithSkip tests executing a skipped test
func TestExecute_WithSkip(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Skip test",
		SkipIf:      "<?php\ndie('skip Extension not available');\n?>",
		Code:        "<?php\necho 'test';\n?>",
		Expect:      "test",
	}

	result := executor.Execute(test)

	if result.Status != StatusSkip {
		t.Errorf("Execute() status = %v, want StatusSkip", result.Status)
	}

	if result.SkipReason == "" {
		t.Error("Execute() should set SkipReason for skipped test")
	}
}

// TestExecute_WithEnv tests executing test with environment variables
func TestExecute_WithEnv(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "ENV test",
		Env: map[string]string{
			"TEST_VAR": "test_value",
		},
		Code:    "<?php\necho getenv('TEST_VAR');\n?>",
		ExpectF: "%s",
	}

	result := executor.Execute(test)

	// May pass or fail depending on getenv implementation
	if result.Status == StatusError {
		t.Logf("Execute() with env variables resulted in error: %s", result.ErrorMessage)
	}
}

// TestExecute_WithArgs tests executing test with command line arguments
func TestExecute_WithArgs(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	test := &PHPTest{
		Description: "Args test",
		Args:        []string{"arg1", "arg2"},
		Code:        "<?php\necho 'test';\n?>",
		Expect:      "test",
	}

	result := executor.Execute(test)

	// Should pass regardless of args since we're not checking them
	if result.Status != StatusPass {
		t.Errorf("Execute() status = %v, want StatusPass", result.Status)
	}
}

// TestExecuteTests tests executing multiple tests
func TestExecuteTests(t *testing.T) {
	if !hasPHPGo() || !hasVMExecution() {
		t.Skip("php-go binary not available or VM execution not implemented")
	}

	executor := NewTestExecutor(getPHPGoPath())
	executor.Verbose = false

	tests := []*PHPTest{
		{
			Description: "Test 1",
			Code:        "<?php\necho '1';\n?>",
			Expect:      "1",
		},
		{
			Description: "Test 2",
			Code:        "<?php\necho '2';\n?>",
			Expect:      "2",
		},
	}

	results := executor.ExecuteTests(tests)

	if len(results) != 2 {
		t.Errorf("ExecuteTests() returned %d results, want 2", len(results))
	}

	for i, result := range results {
		if result.Status != StatusPass {
			t.Errorf("Result %d status = %v, want StatusPass", i, result.Status)
		}
	}
}

// TestFindTests tests finding test files
func TestFindTests(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test files
	test1 := filepath.Join(tmpDir, "test1.phpt")
	test1Content := `--TEST--
Test 1
--FILE--
<?php
echo 'test1';
?>
--EXPECT--
test1`
	if err := os.WriteFile(test1, []byte(test1Content), 0644); err != nil {
		t.Fatalf("Failed to create test1.phpt: %v", err)
	}

	test2 := filepath.Join(tmpDir, "test2.phpt")
	test2Content := `--TEST--
Test 2
--FILE--
<?php
echo 'test2';
?>
--EXPECT--
test2`
	if err := os.WriteFile(test2, []byte(test2Content), 0644); err != nil {
		t.Fatalf("Failed to create test2.phpt: %v", err)
	}

	// Create subdirectory with test
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	test3 := filepath.Join(subDir, "test3.phpt")
	test3Content := `--TEST--
Test 3
--FILE--
<?php
echo 'test3';
?>
--EXPECT--
test3`
	if err := os.WriteFile(test3, []byte(test3Content), 0644); err != nil {
		t.Fatalf("Failed to create test3.phpt: %v", err)
	}

	// Find tests
	tests, err := FindTests(tmpDir)
	if err != nil {
		t.Fatalf("FindTests() error = %v", err)
	}

	if len(tests) != 3 {
		t.Errorf("FindTests() returned %d tests, want 3", len(tests))
	}
}

// TestFindTests_InvalidDirectory tests finding tests in nonexistent directory
func TestFindTests_InvalidDirectory(t *testing.T) {
	_, err := FindTests("/nonexistent/directory")
	if err == nil {
		t.Error("FindTests() should return error for nonexistent directory")
	}
}

// TestIsExpectedError tests checking for expected errors
func TestIsExpectedError(t *testing.T) {
	executor := NewTestExecutor(getPHPGoPath())

	tests := []struct {
		name     string
		test     *PHPTest
		expected bool
	}{
		{
			name: "fatal error",
			test: &PHPTest{
				Expect: "Fatal error: something",
			},
			expected: true,
		},
		{
			name: "parse error",
			test: &PHPTest{
				Expect: "Parse error: syntax error",
			},
			expected: true,
		},
		{
			name: "warning",
			test: &PHPTest{
				Expect: "Warning: undefined variable",
			},
			expected: true,
		},
		{
			name: "no error",
			test: &PHPTest{
				Expect: "Normal output",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.isExpectedError(tt.test)
			if result != tt.expected {
				t.Errorf("isExpectedError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to check if php-go binary exists and return its path
func hasPHPGo() bool {
	// Check if php-go exists in project root
	path := getPHPGoPath()
	_, err := os.Stat(path)
	return err == nil
}

// Helper function to get php-go binary path
func getPHPGoPath() string {
	// Get absolute path to php-go binary
	absPath, err := filepath.Abs("../../php-go")
	if err != nil {
		return "../../php-go"
	}
	return absPath
}

// Helper function to check if VM execution is implemented
func hasVMExecution() bool {
	// VM execution is not yet integrated into the CLI (coming in Phase 10)
	// For now, we're building the test infrastructure
	return false
}
