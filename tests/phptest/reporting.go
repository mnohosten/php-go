package phptest

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// ============================================================================
// Test Results Reporting
// ============================================================================
//
// Generates human-readable and machine-readable test reports

// TestSummary contains aggregated test statistics
type TestSummary struct {
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	ErrorTests   int

	TotalDuration time.Duration

	PassRate float64 // Percentage of passed tests
}

// ComputeSummary calculates summary statistics from test results
func ComputeSummary(results []*TestResult) *TestSummary {
	summary := &TestSummary{}

	for _, result := range results {
		summary.TotalTests++
		summary.TotalDuration += result.Duration

		switch result.Status {
		case StatusPass:
			summary.PassedTests++
		case StatusFail:
			summary.FailedTests++
		case StatusSkip:
			summary.SkippedTests++
		case StatusError:
			summary.ErrorTests++
		}
	}

	// Calculate pass rate (excluding skipped tests)
	testsRun := summary.TotalTests - summary.SkippedTests
	if testsRun > 0 {
		summary.PassRate = float64(summary.PassedTests) / float64(testsRun) * 100.0
	}

	return summary
}

// TestReporter generates test reports in various formats
type TestReporter struct {
	writer  io.Writer
	verbose bool
}

// NewTestReporter creates a new test reporter
func NewTestReporter(writer io.Writer, verbose bool) *TestReporter {
	return &TestReporter{
		writer:  writer,
		verbose: verbose,
	}
}

// ReportResults generates a comprehensive test report
func (r *TestReporter) ReportResults(results []*TestResult) {
	summary := ComputeSummary(results)

	// Print header
	r.printHeader()

	// Print summary
	r.printSummary(summary)

	// Print category breakdown
	r.printCategoryBreakdown(results)

	// Print failures if any
	if summary.FailedTests > 0 || summary.ErrorTests > 0 {
		r.printFailures(results)
	}

	// Print verbose details if requested
	if r.verbose {
		r.printVerboseResults(results)
	}

	// Print footer
	r.printFooter(summary)
}

// printHeader prints the report header
func (r *TestReporter) printHeader() {
	fmt.Fprintln(r.writer, "")
	fmt.Fprintln(r.writer, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(r.writer, "║              PHP-Go Test Suite Results                         ║")
	fmt.Fprintln(r.writer, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(r.writer, "")
}

// printSummary prints the test summary
func (r *TestReporter) printSummary(summary *TestSummary) {
	fmt.Fprintln(r.writer, "SUMMARY")
	fmt.Fprintln(r.writer, strings.Repeat("-", 70))
	fmt.Fprintf(r.writer, "Total Tests:     %d\n", summary.TotalTests)
	fmt.Fprintf(r.writer, "Passed:          %d (✓)\n", summary.PassedTests)
	fmt.Fprintf(r.writer, "Failed:          %d (✗)\n", summary.FailedTests)
	fmt.Fprintf(r.writer, "Errors:          %d (⚠)\n", summary.ErrorTests)
	fmt.Fprintf(r.writer, "Skipped:         %d (⊗)\n", summary.SkippedTests)
	fmt.Fprintf(r.writer, "Pass Rate:       %.2f%%\n", summary.PassRate)
	fmt.Fprintf(r.writer, "Total Duration:  %.3fs\n", summary.TotalDuration.Seconds())
	fmt.Fprintln(r.writer, "")
}

// printCategoryBreakdown prints results grouped by category
func (r *TestReporter) printCategoryBreakdown(results []*TestResult) {
	grouped := GroupResultsByCategory(results)

	// Sort categories for consistent output
	categories := make([]TestCategory, 0, len(grouped))
	for category := range grouped {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return string(categories[i]) < string(categories[j])
	})

	fmt.Fprintln(r.writer, "CATEGORY BREAKDOWN")
	fmt.Fprintln(r.writer, strings.Repeat("-", 70))

	for _, category := range categories {
		categoryResults := grouped[category]
		categorySummary := ComputeSummary(categoryResults)

		fmt.Fprintf(r.writer, "%-15s: %3d tests | %3d pass | %3d fail | %3d skip | %.1f%% pass\n",
			category,
			categorySummary.TotalTests,
			categorySummary.PassedTests,
			categorySummary.FailedTests,
			categorySummary.SkippedTests,
			categorySummary.PassRate,
		)
	}

	fmt.Fprintln(r.writer, "")
}

// printFailures prints details of failed tests
func (r *TestReporter) printFailures(results []*TestResult) {
	fmt.Fprintln(r.writer, "FAILURES")
	fmt.Fprintln(r.writer, strings.Repeat("-", 70))

	for i, result := range results {
		if result.Status == StatusFail || result.Status == StatusError {
			fmt.Fprintf(r.writer, "%d. %s [%s]\n", i+1, result.Test.Description, result.Status)
			fmt.Fprintf(r.writer, "   File: %s\n", result.Test.FilePath)

			if result.ErrorMessage != "" {
				fmt.Fprintf(r.writer, "   Error: %s\n", result.ErrorMessage)
			}

			if result.Status == StatusFail {
				// Show output mismatch
				fmt.Fprintln(r.writer, "   Expected:")
				fmt.Fprintln(r.writer, r.indentText(result.ExpectedOutput, "     "))
				fmt.Fprintln(r.writer, "   Actual:")
				fmt.Fprintln(r.writer, r.indentText(result.ActualOutput, "     "))
			} else if result.ErrorOutput != "" {
				fmt.Fprintln(r.writer, "   Output:")
				fmt.Fprintln(r.writer, r.indentText(result.ErrorOutput, "     "))
			}

			fmt.Fprintln(r.writer, "")
		}
	}
}

// printVerboseResults prints all test results
func (r *TestReporter) printVerboseResults(results []*TestResult) {
	fmt.Fprintln(r.writer, "DETAILED RESULTS")
	fmt.Fprintln(r.writer, strings.Repeat("-", 70))

	for i, result := range results {
		statusSymbol := r.getStatusSymbol(result.Status)
		fmt.Fprintf(r.writer, "%4d. [%s] %-50s (%.3fs)\n",
			i+1,
			statusSymbol,
			r.truncate(result.Test.Description, 50),
			result.Duration.Seconds(),
		)

		if result.Status == StatusSkip {
			fmt.Fprintf(r.writer, "      Reason: %s\n", result.SkipReason)
		} else if result.Status == StatusFail || result.Status == StatusError {
			fmt.Fprintf(r.writer, "      Error: %s\n", result.ErrorMessage)
		}
	}

	fmt.Fprintln(r.writer, "")
}

// printFooter prints the report footer
func (r *TestReporter) printFooter(summary *TestSummary) {
	fmt.Fprintln(r.writer, strings.Repeat("=", 70))

	if summary.FailedTests == 0 && summary.ErrorTests == 0 {
		fmt.Fprintln(r.writer, "✓ ALL TESTS PASSED!")
	} else {
		fmt.Fprintf(r.writer, "✗ %d TEST(S) FAILED\n", summary.FailedTests+summary.ErrorTests)
	}

	fmt.Fprintln(r.writer, "")
}

// getStatusSymbol returns a symbol for a test status
func (r *TestReporter) getStatusSymbol(status TestStatus) string {
	switch status {
	case StatusPass:
		return "✓"
	case StatusFail:
		return "✗"
	case StatusSkip:
		return "⊗"
	case StatusError:
		return "⚠"
	default:
		return "?"
	}
}

// indentText indents each line of text
func (r *TestReporter) indentText(text, indent string) string {
	lines := strings.Split(text, "\n")
	indented := make([]string, len(lines))
	for i, line := range lines {
		indented[i] = indent + line
	}
	return strings.Join(indented, "\n")
}

// truncate truncates a string to a maximum length
func (r *TestReporter) truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ReportProgress prints a progress indicator
func (r *TestReporter) ReportProgress(current, total int, result *TestResult) {
	if !r.verbose {
		return
	}

	symbol := r.getStatusSymbol(result.Status)
	percentage := float64(current) / float64(total) * 100.0

	fmt.Fprintf(r.writer, "[%3.0f%%] %s %s\n",
		percentage,
		symbol,
		result.Test.Description,
	)
}

// ReportJUnit generates a JUnit-compatible XML report
func ReportJUnit(results []*TestResult, writer io.Writer) {
	summary := ComputeSummary(results)

	fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(writer, `<testsuite name="PHP-Go Test Suite" tests="%d" failures="%d" errors="%d" skipped="%d" time="%.3f">`+"\n",
		summary.TotalTests,
		summary.FailedTests,
		summary.ErrorTests,
		summary.SkippedTests,
		summary.TotalDuration.Seconds(),
	)

	for _, result := range results {
		className := string(CategorizeTest(result.Test))
		testName := result.Test.Description

		fmt.Fprintf(writer, `  <testcase classname="%s" name="%s" time="%.3f">`+"\n",
			className,
			escapeXML(testName),
			result.Duration.Seconds(),
		)

		switch result.Status {
		case StatusFail:
			fmt.Fprintf(writer, `    <failure message="%s">%s</failure>`+"\n",
				escapeXML(result.ErrorMessage),
				escapeXML(result.ActualOutput),
			)
		case StatusError:
			fmt.Fprintf(writer, `    <error message="%s">%s</error>`+"\n",
				escapeXML(result.ErrorMessage),
				escapeXML(result.ErrorOutput),
			)
		case StatusSkip:
			fmt.Fprintf(writer, `    <skipped message="%s"/>`+"\n",
				escapeXML(result.SkipReason),
			)
		}

		fmt.Fprintln(writer, `  </testcase>`)
	}

	fmt.Fprintln(writer, `</testsuite>`)
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// ReportTAP generates a TAP (Test Anything Protocol) report
func ReportTAP(results []*TestResult, writer io.Writer) {
	fmt.Fprintf(writer, "1..%d\n", len(results))

	for i, result := range results {
		testNum := i + 1

		switch result.Status {
		case StatusPass:
			fmt.Fprintf(writer, "ok %d - %s\n", testNum, result.Test.Description)
		case StatusFail:
			fmt.Fprintf(writer, "not ok %d - %s\n", testNum, result.Test.Description)
			fmt.Fprintf(writer, "  ---\n")
			fmt.Fprintf(writer, "  message: %s\n", result.ErrorMessage)
			fmt.Fprintf(writer, "  ...\n")
		case StatusSkip:
			fmt.Fprintf(writer, "ok %d - %s # SKIP %s\n", testNum, result.Test.Description, result.SkipReason)
		case StatusError:
			fmt.Fprintf(writer, "not ok %d - %s\n", testNum, result.Test.Description)
			fmt.Fprintf(writer, "  ---\n")
			fmt.Fprintf(writer, "  message: %s\n", result.ErrorMessage)
			fmt.Fprintf(writer, "  severity: error\n")
			fmt.Fprintf(writer, "  ...\n")
		}
	}
}
