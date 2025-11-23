package phptest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ============================================================================
// Test Executor
// ============================================================================
//
// Executes PHPT tests by running PHP code and comparing output

// TestResult represents the result of a test execution
type TestResult struct {
	Test     *PHPTest
	Status   TestStatus
	Duration time.Duration

	// Output
	ActualOutput string
	ExpectedOutput string

	// Error information
	ErrorMessage string
	ErrorOutput  string

	// Skip information
	SkipReason string
}

// TestStatus represents the outcome of a test
type TestStatus int

const (
	StatusUnknown TestStatus = iota
	StatusPass
	StatusFail
	StatusSkip
	StatusError
)

// String returns the string representation of a test status
func (s TestStatus) String() string {
	switch s {
	case StatusPass:
		return "PASS"
	case StatusFail:
		return "FAIL"
	case StatusSkip:
		return "SKIP"
	case StatusError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// TestExecutor executes PHPT tests
type TestExecutor struct {
	// PHPBinary is the path to the PHP interpreter
	PHPBinary string

	// Timeout for test execution (default: 30 seconds)
	Timeout time.Duration

	// WorkDir is the working directory for test execution
	WorkDir string

	// Verbose enables detailed output
	Verbose bool
}

// NewTestExecutor creates a new test executor
func NewTestExecutor(phpBinary string) *TestExecutor {
	if phpBinary == "" {
		phpBinary = "php-go" // Default to php-go binary
	}

	return &TestExecutor{
		PHPBinary: phpBinary,
		Timeout:   30 * time.Second,
		WorkDir:   "",
		Verbose:   false,
	}
}

// Execute runs a single PHPT test
func (e *TestExecutor) Execute(test *PHPTest) *TestResult {
	result := &TestResult{
		Test:   test,
		Status: StatusUnknown,
	}

	startTime := time.Now()
	defer func() {
		result.Duration = time.Since(startTime)
	}()

	// Check skip condition
	if test.HasSkipCondition() {
		if shouldSkip, reason := e.checkSkipCondition(test); shouldSkip {
			result.Status = StatusSkip
			result.SkipReason = reason
			return result
		}
	}

	// Get code to execute
	code, err := test.GetCodeToExecute()
	if err != nil {
		result.Status = StatusError
		result.ErrorMessage = fmt.Sprintf("Failed to get code: %v", err)
		return result
	}

	// Execute the test
	output, exitCode, err := e.runPHP(code, test)
	if err != nil {
		result.Status = StatusError
		result.ErrorMessage = fmt.Sprintf("Execution failed: %v", err)
		result.ErrorOutput = output
		return result
	}

	// Store actual output
	result.ActualOutput = output
	result.ExpectedOutput = test.GetExpectedOutput()

	// Check for runtime errors (non-zero exit code without EXPECT_EXIT)
	if exitCode != 0 && !e.isExpectedError(test) {
		result.Status = StatusError
		result.ErrorMessage = fmt.Sprintf("Process exited with code %d", exitCode)
		result.ErrorOutput = output
		return result
	}

	// Compare output
	matched, err := test.MatchOutput(output)
	if err != nil {
		result.Status = StatusError
		result.ErrorMessage = fmt.Sprintf("Output matching failed: %v", err)
		return result
	}

	if matched {
		result.Status = StatusPass
	} else {
		result.Status = StatusFail
		result.ErrorMessage = "Output mismatch"
	}

	// Run cleanup if present
	if test.Clean != "" {
		e.runCleanup(test)
	}

	return result
}

// checkSkipCondition evaluates the SKIPIF section
func (e *TestExecutor) checkSkipCondition(test *PHPTest) (bool, string) {
	output, exitCode, err := e.runPHP(test.SkipIf, test)
	if err != nil {
		return true, fmt.Sprintf("SKIPIF execution error: %v", err)
	}

	// If SKIPIF exits with 0 and outputs 'skip', the test is skipped
	if exitCode == 0 {
		output = strings.TrimSpace(output)
		if strings.HasPrefix(output, "skip") {
			reason := strings.TrimPrefix(output, "skip")
			reason = strings.TrimSpace(reason)
			if reason == "" {
				reason = "SKIPIF condition met"
			}
			return true, reason
		}
	}

	return false, ""
}

// runPHP executes PHP code and returns output and exit code
func (e *TestExecutor) runPHP(code string, test *PHPTest) (string, int, error) {
	// Create temporary file for code
	tmpfile, err := os.CreateTemp(e.getWorkDir(), "phpt-*.php")
	if err != nil {
		return "", 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(code)); err != nil {
		tmpfile.Close()
		return "", 0, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpfile.Close()

	// Prepare command
	args := []string{tmpfile.Name()}

	// Add command line arguments from test
	if len(test.Args) > 0 {
		args = append(args, test.Args...)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.PHPBinary, args...)

	// Set working directory
	workDir := e.getWorkDir()
	if test.FilePath != "" {
		// Use test file's directory as working directory
		workDir = filepath.Dir(test.FilePath)
	}
	cmd.Dir = workDir

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range test.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add INI settings as environment variables
	for key, value := range test.Ini {
		// PHP recognizes -d flag for INI settings
		args = []string{"-d", fmt.Sprintf("%s=%s", key, value), tmpfile.Name()}
		if len(test.Args) > 0 {
			args = append(args, test.Args...)
		}
	}

	// Rebuild command with INI flags
	if len(test.Ini) > 0 {
		cmd = exec.CommandContext(ctx, e.PHPBinary, args...)
		cmd.Dir = workDir
		cmd.Env = os.Environ()
		for key, value := range test.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err = cmd.Run()

	// Combine stdout and stderr
	output := stdout.String()
	if stderr.Len() > 0 {
		// PHP often outputs errors to stderr
		output += stderr.String()
	}

	// Get exit code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Non-exit error (e.g., command not found, timeout)
			return output, -1, err
		}
	}

	return output, exitCode, nil
}

// runCleanup runs the cleanup code
func (e *TestExecutor) runCleanup(test *PHPTest) {
	// Create a temporary test with just the cleanup code
	cleanupTest := &PHPTest{
		FilePath: test.FilePath,
		Env:      test.Env,
		Ini:      test.Ini,
	}

	_, _, _ = e.runPHP(test.Clean, cleanupTest)
	// Ignore cleanup errors
}

// isExpectedError checks if a non-zero exit code is expected
func (e *TestExecutor) isExpectedError(test *PHPTest) bool {
	// Some tests expect errors/warnings which cause non-zero exit codes
	// Check if expected output contains error patterns
	expected := test.GetExpectedOutput()
	return strings.Contains(expected, "Fatal error") ||
		strings.Contains(expected, "Parse error") ||
		strings.Contains(expected, "Warning:") ||
		strings.Contains(expected, "Notice:")
}

// getWorkDir returns the working directory for test execution
func (e *TestExecutor) getWorkDir() string {
	if e.WorkDir != "" {
		return e.WorkDir
	}
	return os.TempDir()
}

// ExecuteTests runs multiple PHPT tests
func (e *TestExecutor) ExecuteTests(tests []*PHPTest) []*TestResult {
	results := make([]*TestResult, 0, len(tests))

	for _, test := range tests {
		if e.Verbose {
			fmt.Printf("Running: %s\n", test.Description)
		}

		result := e.Execute(test)
		results = append(results, result)

		if e.Verbose {
			fmt.Printf("  Status: %s", result.Status)
			if result.Duration > 0 {
				fmt.Printf(" (%.3fs)", result.Duration.Seconds())
			}
			if result.Status == StatusSkip {
				fmt.Printf(" - %s", result.SkipReason)
			} else if result.Status == StatusFail {
				fmt.Printf(" - %s", result.ErrorMessage)
			}
			fmt.Println()
		}
	}

	return results
}

// FindTests recursively finds all .phpt files in a directory
func FindTests(dir string) ([]*PHPTest, error) {
	var tests []*PHPTest

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a .phpt file
		if info.IsDir() || filepath.Ext(path) != ".phpt" {
			return nil
		}

		// Parse test
		test, parseErr := ParseTest(path)
		if parseErr != nil {
			// Skip invalid tests but continue walking
			fmt.Fprintf(os.Stderr, "Warning: Failed to parse %s: %v\n", path, parseErr)
			return nil
		}

		tests = append(tests, test)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tests, nil
}
