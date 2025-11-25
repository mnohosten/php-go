package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/krizos/php-go/tests/phptest"
)

func main() {
	// Command line flags
	testDir := flag.String("dir", "php-src/Zend/tests", "Directory containing .phpt test files")
	phpBinary := flag.String("php", "./php-go", "Path to PHP interpreter (php-go or php)")
	category := flag.String("category", "", "Test category filter (language, stdlib, syntax, oop, etc.)")
	pattern := flag.String("pattern", "", "Test file pattern (substring match)")
	verbose := flag.Bool("verbose", false, "Verbose output")
	timeout := flag.Duration("timeout", 30*time.Second, "Test timeout")
	maxTests := flag.Int("max", 100, "Maximum number of tests to run (0 = all)")
	continueOnError := flag.Bool("continue", true, "Continue running tests after failures")

	flag.Parse()

	// Ensure test directory exists
	if _, err := os.Stat(*testDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Test directory does not exist: %s\n", *testDir)
		os.Exit(1)
	}

	// Ensure PHP binary exists
	phpPath := *phpBinary
	if !filepath.IsAbs(phpPath) {
		// Resolve relative path
		absPath, err := filepath.Abs(phpPath)
		if err == nil {
			phpPath = absPath
		}
	}
	if _, err := os.Stat(phpPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: PHP binary does not exist: %s\n", phpPath)
		os.Exit(1)
	}

	fmt.Printf("PHP-Go PHPT Test Runner\n")
	fmt.Printf("========================\n\n")
	fmt.Printf("Test Directory: %s\n", *testDir)
	fmt.Printf("PHP Binary:     %s\n", phpPath)
	if *category != "" {
		fmt.Printf("Category:       %s\n", *category)
	}
	if *pattern != "" {
		fmt.Printf("Pattern:        %s\n", *pattern)
	}
	if *maxTests > 0 {
		fmt.Printf("Max Tests:      %d\n", *maxTests)
	}
	fmt.Printf("\n")

	// Find all tests
	fmt.Printf("Finding tests...\n")
	tests, err := phptest.FindTests(*testDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding tests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d test files\n", len(tests))

	// Filter by category if requested
	if *category != "" {
		var filtered []*phptest.PHPTest
		targetCat := phptest.TestCategory(*category)
		for _, test := range tests {
			cat := phptest.CategorizeTest(test)
			if cat == targetCat {
				filtered = append(filtered, test)
			}
		}
		tests = filtered
		fmt.Printf("Filtered to %d tests in category '%s'\n", len(tests), *category)
	}

	// Filter by pattern if requested
	if *pattern != "" {
		var filtered []*phptest.PHPTest
		for _, test := range tests {
			if strings.Contains(test.FilePath, *pattern) {
				filtered = append(filtered, test)
			}
		}
		tests = filtered
		fmt.Printf("Filtered to %d tests matching pattern '%s'\n", len(tests), *pattern)
	}

	// Limit number of tests if requested
	if *maxTests > 0 && len(tests) > *maxTests {
		tests = tests[:*maxTests]
		fmt.Printf("Limited to first %d tests\n", *maxTests)
	}

	if len(tests) == 0 {
		fmt.Printf("No tests to run\n")
		os.Exit(0)
	}

	fmt.Printf("\nRunning %d tests...\n\n", len(tests))

	// Create executor
	executor := phptest.NewTestExecutor(phpPath)
	executor.Timeout = *timeout
	executor.Verbose = *verbose

	// Run tests
	var results []*phptest.TestResult
	startTime := time.Now()

	var pass, fail, skip, errCount int

	for i, test := range tests {
		// Run test
		result := executor.Execute(test)
		results = append(results, result)

		// Update counters
		switch result.Status {
		case phptest.StatusPass:
			pass++
		case phptest.StatusFail:
			fail++
		case phptest.StatusSkip:
			skip++
		case phptest.StatusError:
			errCount++
		}

		// Print progress
		if *verbose || result.Status != phptest.StatusPass {
			status := result.Status.String()
			symbol := " "
			if result.Status == phptest.StatusPass {
				symbol = "✓"
			} else if result.Status == phptest.StatusFail {
				symbol = "✗"
			} else if result.Status == phptest.StatusSkip {
				symbol = "⊘"
			} else if result.Status == phptest.StatusError {
				symbol = "⚠"
			}

			fmt.Printf("[%d/%d] %s %s: %s", i+1, len(tests), symbol, status, test.Description)
			if result.Status == phptest.StatusSkip && result.SkipReason != "" {
				fmt.Printf(" (%s)", result.SkipReason)
			}
			fmt.Printf("\n")

			if result.Status == phptest.StatusFail && *verbose {
				fmt.Printf("  File: %s\n", test.FilePath)
				fmt.Printf("  Expected:\n%s\n", result.ExpectedOutput)
				fmt.Printf("  Actual:\n%s\n", result.ActualOutput)
			}
			if result.Status == phptest.StatusError {
				fmt.Printf("  File: %s\n", test.FilePath)
				fmt.Printf("  Error: %s\n", result.ErrorMessage)
				if result.ErrorOutput != "" && *verbose {
					fmt.Printf("  Output: %s\n", result.ErrorOutput)
				}
			}
		} else if (i+1)%10 == 0 || i+1 == len(tests) {
			fmt.Printf("Progress: %d/%d (%d✓ %d✗ %d⊘ %d⚠)\n", i+1, len(tests), pass, fail, skip, errCount)
		}

		// Stop on error if requested
		if !*continueOnError && (result.Status == phptest.StatusError || result.Status == phptest.StatusFail) {
			fmt.Printf("\nStopping due to error/failure (use -continue to continue)\n")
			break
		}
	}

	totalDuration := time.Since(startTime)

	// Generate report
	fmt.Printf("\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Test Results\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n")

	total := len(results)
	passRate := 0.0
	if total > 0 {
		passRate = float64(pass) / float64(total) * 100
	}

	fmt.Printf("Total:    %d tests\n", total)
	fmt.Printf("Passed:   %d (%.1f%%)\n", pass, passRate)
	fmt.Printf("Failed:   %d\n", fail)
	fmt.Printf("Skipped:  %d\n", skip)
	fmt.Printf("Errors:   %d\n", errCount)
	fmt.Printf("Duration: %s\n", totalDuration.Round(time.Millisecond))
	if total > 0 {
		avgDuration := totalDuration / time.Duration(total)
		fmt.Printf("Avg/test: %s\n", avgDuration.Round(time.Millisecond))
	}
	fmt.Printf("\n")

	// List failed tests if not verbose
	if !*verbose && (fail > 0 || errCount > 0) {
		fmt.Printf("Failed/Error tests:\n")
		for _, r := range results {
			if r.Status == phptest.StatusFail || r.Status == phptest.StatusError {
				fmt.Printf("  - %s: %s\n", r.Status.String(), r.Test.Description)
				fmt.Printf("    File: %s\n", r.Test.FilePath)
				if r.Status == phptest.StatusError {
					fmt.Printf("    Error: %s\n", r.ErrorMessage)
				}
			}
		}
		fmt.Printf("\n")
	}

	// Exit with appropriate code
	if fail > 0 || errCount > 0 {
		os.Exit(1)
	}
}
