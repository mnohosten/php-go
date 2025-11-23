package phptest

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestComputeSummary tests summary computation
func TestComputeSummary(t *testing.T) {
	results := []*TestResult{
		{Status: StatusPass, Duration: 100 * time.Millisecond},
		{Status: StatusPass, Duration: 200 * time.Millisecond},
		{Status: StatusFail, Duration: 50 * time.Millisecond},
		{Status: StatusSkip, Duration: 0},
		{Status: StatusError, Duration: 30 * time.Millisecond},
	}

	summary := ComputeSummary(results)

	if summary.TotalTests != 5 {
		t.Errorf("TotalTests = %d, want 5", summary.TotalTests)
	}

	if summary.PassedTests != 2 {
		t.Errorf("PassedTests = %d, want 2", summary.PassedTests)
	}

	if summary.FailedTests != 1 {
		t.Errorf("FailedTests = %d, want 1", summary.FailedTests)
	}

	if summary.SkippedTests != 1 {
		t.Errorf("SkippedTests = %d, want 1", summary.SkippedTests)
	}

	if summary.ErrorTests != 1 {
		t.Errorf("ErrorTests = %d, want 1", summary.ErrorTests)
	}

	expectedDuration := 380 * time.Millisecond
	if summary.TotalDuration != expectedDuration {
		t.Errorf("TotalDuration = %v, want %v", summary.TotalDuration, expectedDuration)
	}

	// Pass rate should be 50% (2 passed out of 4 non-skipped)
	expectedPassRate := 50.0
	if summary.PassRate != expectedPassRate {
		t.Errorf("PassRate = %.2f%%, want %.2f%%", summary.PassRate, expectedPassRate)
	}
}

// TestComputeSummary_AllPassed tests summary with all passing tests
func TestComputeSummary_AllPassed(t *testing.T) {
	results := []*TestResult{
		{Status: StatusPass},
		{Status: StatusPass},
		{Status: StatusPass},
	}

	summary := ComputeSummary(results)

	if summary.PassRate != 100.0 {
		t.Errorf("PassRate = %.2f%%, want 100.00%%", summary.PassRate)
	}
}

// TestComputeSummary_AllSkipped tests summary with all skipped tests
func TestComputeSummary_AllSkipped(t *testing.T) {
	results := []*TestResult{
		{Status: StatusSkip},
		{Status: StatusSkip},
	}

	summary := ComputeSummary(results)

	if summary.PassRate != 0.0 {
		t.Errorf("PassRate = %.2f%%, want 0.00%%", summary.PassRate)
	}
}

// TestNewTestReporter tests creating a test reporter
func TestNewTestReporter(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, false)

	if reporter == nil {
		t.Fatal("NewTestReporter() returned nil")
	}

	if reporter.verbose {
		t.Error("verbose should be false")
	}
}

// TestReportResults tests generating a report
func TestReportResults(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, false)

	results := []*TestResult{
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test1.phpt",
				Description: "Test 1",
			},
			Status:   StatusPass,
			Duration: 100 * time.Millisecond,
		},
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test2.phpt",
				Description: "Test 2",
			},
			Status:   StatusFail,
			Duration: 200 * time.Millisecond,
		},
	}

	reporter.ReportResults(results)

	output := buf.String()

	// Check that output contains key sections
	if !strings.Contains(output, "SUMMARY") {
		t.Error("Report should contain SUMMARY section")
	}

	if !strings.Contains(output, "CATEGORY BREAKDOWN") {
		t.Error("Report should contain CATEGORY BREAKDOWN section")
	}

	if !strings.Contains(output, "FAILURES") {
		t.Error("Report should contain FAILURES section when there are failures")
	}

	if !strings.Contains(output, "Total Tests:") {
		t.Error("Report should contain total tests count")
	}
}

// TestReportResults_AllPass tests report with all passing tests
func TestReportResults_AllPass(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, false)

	results := []*TestResult{
		{
			Test: &PHPTest{
				Description: "Test 1",
			},
			Status: StatusPass,
		},
	}

	reporter.ReportResults(results)

	output := buf.String()

	if !strings.Contains(output, "ALL TESTS PASSED") {
		t.Error("Report should indicate all tests passed")
	}

	if strings.Contains(output, "FAILURES") {
		t.Error("Report should not contain FAILURES section when all pass")
	}
}

// TestReportResults_Verbose tests verbose reporting
func TestReportResults_Verbose(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, true)

	results := []*TestResult{
		{
			Test: &PHPTest{
				Description: "Test 1",
			},
			Status:   StatusPass,
			Duration: 50 * time.Millisecond,
		},
	}

	reporter.ReportResults(results)

	output := buf.String()

	if !strings.Contains(output, "DETAILED RESULTS") {
		t.Error("Verbose report should contain DETAILED RESULTS section")
	}
}

// TestGetStatusSymbol tests status symbol generation
func TestGetStatusSymbol(t *testing.T) {
	reporter := &TestReporter{}

	tests := []struct {
		status   TestStatus
		expected string
	}{
		{StatusPass, "✓"},
		{StatusFail, "✗"},
		{StatusSkip, "⊗"},
		{StatusError, "⚠"},
		{StatusUnknown, "?"},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			result := reporter.getStatusSymbol(tt.status)
			if result != tt.expected {
				t.Errorf("getStatusSymbol(%v) = %s, want %s", tt.status, result, tt.expected)
			}
		})
	}
}

// TestIndentText tests text indentation
func TestIndentText(t *testing.T) {
	reporter := &TestReporter{}

	input := "line1\nline2\nline3"
	expected := "  line1\n  line2\n  line3"

	result := reporter.indentText(input, "  ")

	if result != expected {
		t.Errorf("indentText() = %q, want %q", result, expected)
	}
}

// TestTruncate tests string truncation
func TestTruncate(t *testing.T) {
	reporter := &TestReporter{}

	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length",
			input:    "exactly10!",
			maxLen:   10,
			expected: "exactly10!",
		},
		{
			name:     "truncation needed",
			input:    "this is a very long string",
			maxLen:   15,
			expected: "this is a ve...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reporter.truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestReportJUnit tests JUnit XML report generation
func TestReportJUnit(t *testing.T) {
	buf := &bytes.Buffer{}

	results := []*TestResult{
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test1.phpt",
				Description: "Test 1",
			},
			Status:   StatusPass,
			Duration: 100 * time.Millisecond,
		},
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test2.phpt",
				Description: "Test 2",
			},
			Status:       StatusFail,
			Duration:     200 * time.Millisecond,
			ErrorMessage: "Output mismatch",
			ActualOutput: "actual",
		},
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test3.phpt",
				Description: "Test 3",
			},
			Status:     StatusSkip,
			SkipReason: "Extension not available",
		},
	}

	ReportJUnit(results, buf)

	output := buf.String()

	// Check for XML structure
	if !strings.Contains(output, "<?xml version") {
		t.Error("JUnit report should contain XML declaration")
	}

	if !strings.Contains(output, "<testsuite") {
		t.Error("JUnit report should contain testsuite element")
	}

	if !strings.Contains(output, "<testcase") {
		t.Error("JUnit report should contain testcase elements")
	}

	if !strings.Contains(output, "<failure") {
		t.Error("JUnit report should contain failure element for failed test")
	}

	if !strings.Contains(output, "<skipped") {
		t.Error("JUnit report should contain skipped element for skipped test")
	}
}

// TestEscapeXML tests XML escaping
func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"<tag>", "&lt;tag&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"'quoted'", "&apos;quoted&apos;"},
		{"<>&\"'", "&lt;&gt;&amp;&quot;&apos;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeXML(tt.input)
			if result != tt.expected {
				t.Errorf("escapeXML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestReportTAP tests TAP format report generation
func TestReportTAP(t *testing.T) {
	buf := &bytes.Buffer{}

	results := []*TestResult{
		{
			Test: &PHPTest{
				Description: "Test 1",
			},
			Status: StatusPass,
		},
		{
			Test: &PHPTest{
				Description: "Test 2",
			},
			Status:       StatusFail,
			ErrorMessage: "Failed",
		},
		{
			Test: &PHPTest{
				Description: "Test 3",
			},
			Status:     StatusSkip,
			SkipReason: "Not available",
		},
	}

	ReportTAP(results, buf)

	output := buf.String()

	// Check TAP format
	if !strings.HasPrefix(output, "1..3") {
		t.Error("TAP report should start with test plan")
	}

	if !strings.Contains(output, "ok 1 - Test 1") {
		t.Error("TAP report should contain ok for passed test")
	}

	if !strings.Contains(output, "not ok 2 - Test 2") {
		t.Error("TAP report should contain not ok for failed test")
	}

	if !strings.Contains(output, "ok 3 - Test 3 # SKIP") {
		t.Error("TAP report should contain SKIP directive for skipped test")
	}
}

// TestReportProgress tests progress reporting
func TestReportProgress(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, true)

	result := &TestResult{
		Test: &PHPTest{
			Description: "Test 1",
		},
		Status: StatusPass,
	}

	reporter.ReportProgress(5, 10, result)

	output := buf.String()

	if !strings.Contains(output, "50%") {
		t.Error("Progress report should contain percentage")
	}

	if !strings.Contains(output, "Test 1") {
		t.Error("Progress report should contain test description")
	}
}

// TestReportProgress_NonVerbose tests that progress is not reported in non-verbose mode
func TestReportProgress_NonVerbose(t *testing.T) {
	buf := &bytes.Buffer{}
	reporter := NewTestReporter(buf, false)

	result := &TestResult{
		Test: &PHPTest{
			Description: "Test 1",
		},
		Status: StatusPass,
	}

	reporter.ReportProgress(5, 10, result)

	output := buf.String()

	if output != "" {
		t.Error("Progress should not be reported in non-verbose mode")
	}
}
