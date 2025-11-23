package phptest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ============================================================================
// PHPT Test Runner
// ============================================================================
//
// This package implements a parser and runner for PHP's PHPT test format.
// PHPT is the standard test format used by PHP's own test suite.
//
// PHPT Format:
// - --TEST--: Test description (required)
// - --SKIPIF--: PHP code that determines if test should be skipped
// - --FILE--: PHP code to execute (required)
// - --EXPECT--: Expected exact output
// - --EXPECTF--: Expected output with format placeholders (%s, %d, etc.)
// - --EXPECTREGEX--: Expected output as regex pattern
// - --EXPECTHEADERS--: Expected HTTP headers
// - --ENV--: Environment variables
// - --ARGS--: Command line arguments
// - --INI--: INI settings
// - --CLEAN--: Cleanup code to run after test
// - --CREDITS--: Test credits/author
// - --EXTENSIONS--: Required extensions
//
// Reference: https://qa.php.net/phpt_details.php

// TestSection represents a section in a PHPT test file
type TestSection int

const (
	SectionUnknown TestSection = iota
	SectionTest
	SectionSkipIf
	SectionFile
	SectionExpect
	SectionExpectF
	SectionExpectRegex
	SectionExpectHeaders
	SectionEnv
	SectionArgs
	SectionIni
	SectionClean
	SectionCredits
	SectionExtensions
	SectionFileExternal
	SectionRedirectTest
)

// String returns the string representation of a section
func (s TestSection) String() string {
	switch s {
	case SectionTest:
		return "TEST"
	case SectionSkipIf:
		return "SKIPIF"
	case SectionFile:
		return "FILE"
	case SectionExpect:
		return "EXPECT"
	case SectionExpectF:
		return "EXPECTF"
	case SectionExpectRegex:
		return "EXPECTREGEX"
	case SectionExpectHeaders:
		return "EXPECTHEADERS"
	case SectionEnv:
		return "ENV"
	case SectionArgs:
		return "ARGS"
	case SectionIni:
		return "INI"
	case SectionClean:
		return "CLEAN"
	case SectionCredits:
		return "CREDITS"
	case SectionExtensions:
		return "EXTENSIONS"
	case SectionFileExternal:
		return "FILE_EXTERNAL"
	case SectionRedirectTest:
		return "REDIRECTTEST"
	default:
		return "UNKNOWN"
	}
}

// PHPTest represents a parsed PHPT test
type PHPTest struct {
	FilePath string

	// Required sections
	Description string // --TEST--
	Code        string // --FILE--

	// Optional execution modifiers
	SkipIf     string            // --SKIPIF--
	Args       []string          // --ARGS--
	Env        map[string]string // --ENV--
	Ini        map[string]string // --INI--
	Extensions []string          // --EXTENSIONS--

	// Expected output (exactly one required)
	Expect      string // --EXPECT-- (exact match)
	ExpectF     string // --EXPECTF-- (format placeholders)
	ExpectRegex string // --EXPECTREGEX-- (regex pattern)

	// Optional headers
	ExpectHeaders map[string]string // --EXPECTHEADERS--

	// Cleanup
	Clean string // --CLEAN--

	// Metadata
	Credits string // --CREDITS--

	// External file reference
	FileExternal  string // --FILE_EXTERNAL--
	RedirectTest  string // --REDIRECTTEST--
}

// ParseTest parses a PHPT test file
func ParseTest(filePath string) (*PHPTest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open test file: %w", err)
	}
	defer file.Close()

	test := &PHPTest{
		FilePath: filePath,
		Env:      make(map[string]string),
		Ini:      make(map[string]string),
		ExpectHeaders: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)

	var currentSection TestSection
	var sectionContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a section header
		if strings.HasPrefix(line, "--") && strings.HasSuffix(line, "--") {
			// Save previous section
			if currentSection != SectionUnknown {
				if err := test.setSection(currentSection, strings.TrimSpace(sectionContent.String())); err != nil {
					return nil, err
				}
				sectionContent.Reset()
			}

			// Parse new section
			sectionName := strings.TrimPrefix(strings.TrimSuffix(line, "--"), "--")
			currentSection = parseSectionName(sectionName)

			if currentSection == SectionUnknown {
				return nil, fmt.Errorf("unknown section: %s", sectionName)
			}
			continue
		}

		// Accumulate section content
		if currentSection != SectionUnknown {
			sectionContent.WriteString(line)
			sectionContent.WriteString("\n")
		}
	}

	// Save final section
	if currentSection != SectionUnknown {
		if err := test.setSection(currentSection, strings.TrimSpace(sectionContent.String())); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading test file: %w", err)
	}

	// Validate required sections
	if err := test.Validate(); err != nil {
		return nil, err
	}

	return test, nil
}

// parseSectionName converts a section name string to a TestSection constant
func parseSectionName(name string) TestSection {
	switch name {
	case "TEST":
		return SectionTest
	case "SKIPIF":
		return SectionSkipIf
	case "FILE":
		return SectionFile
	case "EXPECT":
		return SectionExpect
	case "EXPECTF":
		return SectionExpectF
	case "EXPECTREGEX":
		return SectionExpectRegex
	case "EXPECTHEADERS":
		return SectionExpectHeaders
	case "ENV":
		return SectionEnv
	case "ARGS":
		return SectionArgs
	case "INI":
		return SectionIni
	case "CLEAN":
		return SectionClean
	case "CREDITS":
		return SectionCredits
	case "EXTENSIONS":
		return SectionExtensions
	case "FILE_EXTERNAL":
		return SectionFileExternal
	case "REDIRECTTEST":
		return SectionRedirectTest
	default:
		return SectionUnknown
	}
}

// setSection sets the appropriate field in PHPTest based on the section type
func (t *PHPTest) setSection(section TestSection, content string) error {
	switch section {
	case SectionTest:
		t.Description = content

	case SectionFile:
		t.Code = content

	case SectionSkipIf:
		t.SkipIf = content

	case SectionExpect:
		t.Expect = content

	case SectionExpectF:
		t.ExpectF = content

	case SectionExpectRegex:
		t.ExpectRegex = content

	case SectionExpectHeaders:
		t.ExpectHeaders = parseKeyValueSection(content)

	case SectionEnv:
		t.Env = parseKeyValueSection(content)

	case SectionArgs:
		t.Args = parseArgs(content)

	case SectionIni:
		t.Ini = parseKeyValueSection(content)

	case SectionClean:
		t.Clean = content

	case SectionCredits:
		t.Credits = content

	case SectionExtensions:
		t.Extensions = parseExtensions(content)

	case SectionFileExternal:
		t.FileExternal = content

	case SectionRedirectTest:
		t.RedirectTest = content

	default:
		return fmt.Errorf("unknown section: %s", section)
	}

	return nil
}

// parseKeyValueSection parses sections with key=value format
func parseKeyValueSection(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

// parseArgs parses command line arguments
func parseArgs(content string) []string {
	// Arguments are space-separated
	args := strings.Fields(content)
	return args
}

// parseExtensions parses required extensions
func parseExtensions(content string) []string {
	extensions := strings.Fields(content)
	return extensions
}

// Validate checks that the test has required sections
func (t *PHPTest) Validate() error {
	// TEST section is required
	if t.Description == "" {
		return fmt.Errorf("missing required --TEST-- section")
	}

	// FILE or FILE_EXTERNAL is required
	if t.Code == "" && t.FileExternal == "" && t.RedirectTest == "" {
		return fmt.Errorf("missing required --FILE-- or --FILE_EXTERNAL-- section")
	}

	// At least one EXPECT* section is required
	if t.Expect == "" && t.ExpectF == "" && t.ExpectRegex == "" {
		return fmt.Errorf("missing required --EXPECT--, --EXPECTF--, or --EXPECTREGEX-- section")
	}

	// Only one EXPECT* section is allowed
	expectCount := 0
	if t.Expect != "" {
		expectCount++
	}
	if t.ExpectF != "" {
		expectCount++
	}
	if t.ExpectRegex != "" {
		expectCount++
	}
	if expectCount > 1 {
		return fmt.Errorf("only one of --EXPECT--, --EXPECTF--, or --EXPECTREGEX-- is allowed")
	}

	return nil
}

// HasSkipCondition returns true if the test has a SKIPIF section
func (t *PHPTest) HasSkipCondition() bool {
	return t.SkipIf != ""
}

// GetExpectType returns the type of expected output
func (t *PHPTest) GetExpectType() string {
	if t.Expect != "" {
		return "exact"
	}
	if t.ExpectF != "" {
		return "format"
	}
	if t.ExpectRegex != "" {
		return "regex"
	}
	return "none"
}

// GetExpectedOutput returns the expected output based on type
func (t *PHPTest) GetExpectedOutput() string {
	if t.Expect != "" {
		return t.Expect
	}
	if t.ExpectF != "" {
		return t.ExpectF
	}
	if t.ExpectRegex != "" {
		return t.ExpectRegex
	}
	return ""
}

// GetCodeToExecute returns the PHP code to execute
func (t *PHPTest) GetCodeToExecute() (string, error) {
	if t.Code != "" {
		return t.Code, nil
	}

	if t.FileExternal != "" {
		// Read external file
		dir := filepath.Dir(t.FilePath)
		externalPath := filepath.Join(dir, t.FileExternal)

		content, err := os.ReadFile(externalPath)
		if err != nil {
			return "", fmt.Errorf("failed to read external file %s: %w", externalPath, err)
		}

		return string(content), nil
	}

	return "", fmt.Errorf("no code to execute")
}

// MatchOutput checks if the actual output matches the expected output
func (t *PHPTest) MatchOutput(actual string) (bool, error) {
	// Normalize line endings
	actual = normalizeLineEndings(actual)
	expected := normalizeLineEndings(t.GetExpectedOutput())

	switch t.GetExpectType() {
	case "exact":
		return actual == expected, nil

	case "format":
		return matchFormat(actual, expected)

	case "regex":
		return matchRegex(actual, expected)

	default:
		return false, fmt.Errorf("unknown expect type")
	}
}

// normalizeLineEndings converts all line endings to \n
func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// matchFormat checks if actual matches expected with format placeholders
// Placeholders: %s (string), %d (integer), %f (float), %x (hex), %c (char)
func matchFormat(actual, expected string) (bool, error) {
	// Convert EXPECTF format to regex
	pattern := regexp.QuoteMeta(expected)

	// Replace format placeholders with regex patterns
	pattern = strings.ReplaceAll(pattern, "%s", ".*?")       // Any string (non-greedy)
	pattern = strings.ReplaceAll(pattern, "%S", "\\S+")      // Non-whitespace string
	pattern = strings.ReplaceAll(pattern, "%a", ".+?")       // Any character (non-greedy)
	pattern = strings.ReplaceAll(pattern, "%A", ".*")        // Any character (greedy)
	pattern = strings.ReplaceAll(pattern, "%w", "\\s*")      // Optional whitespace
	pattern = strings.ReplaceAll(pattern, "%i", "[+-]?\\d+") // Integer
	pattern = strings.ReplaceAll(pattern, "%d", "\\d+")      // Unsigned integer
	pattern = strings.ReplaceAll(pattern, "%x", "[0-9a-fA-F]+") // Hex
	pattern = strings.ReplaceAll(pattern, "%f", "[+-]?(?:\\d+\\.\\d+|\\d+\\.?|\\.\\.\\d+)(?:[eE][+-]?\\d+)?") // Float
	pattern = strings.ReplaceAll(pattern, "%c", ".")         // Single character
	pattern = strings.ReplaceAll(pattern, "%e", regexp.QuoteMeta(string(filepath.Separator))) // Directory separator

	// Match with anchors
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, actual)
	if err != nil {
		return false, fmt.Errorf("regex error in format matching: %w", err)
	}

	return matched, nil
}

// matchRegex checks if actual matches expected regex pattern
func matchRegex(actual, expected string) (bool, error) {
	matched, err := regexp.MatchString(expected, actual)
	if err != nil {
		return false, fmt.Errorf("regex error: %w", err)
	}
	return matched, nil
}
