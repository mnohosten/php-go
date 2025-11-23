package phptest

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseSectionName tests parsing section names
func TestParseSectionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TestSection
	}{
		{"TEST", "TEST", SectionTest},
		{"SKIPIF", "SKIPIF", SectionSkipIf},
		{"FILE", "FILE", SectionFile},
		{"EXPECT", "EXPECT", SectionExpect},
		{"EXPECTF", "EXPECTF", SectionExpectF},
		{"EXPECTREGEX", "EXPECTREGEX", SectionExpectRegex},
		{"EXPECTHEADERS", "EXPECTHEADERS", SectionExpectHeaders},
		{"ENV", "ENV", SectionEnv},
		{"ARGS", "ARGS", SectionArgs},
		{"INI", "INI", SectionIni},
		{"CLEAN", "CLEAN", SectionClean},
		{"CREDITS", "CREDITS", SectionCredits},
		{"EXTENSIONS", "EXTENSIONS", SectionExtensions},
		{"FILE_EXTERNAL", "FILE_EXTERNAL", SectionFileExternal},
		{"REDIRECTTEST", "REDIRECTTEST", SectionRedirectTest},
		{"Unknown", "UNKNOWN", SectionUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSectionName(tt.input)
			if result != tt.expected {
				t.Errorf("parseSectionName(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTestSection_String tests section string representation
func TestTestSection_String(t *testing.T) {
	tests := []struct {
		section  TestSection
		expected string
	}{
		{SectionTest, "TEST"},
		{SectionSkipIf, "SKIPIF"},
		{SectionFile, "FILE"},
		{SectionExpect, "EXPECT"},
		{SectionExpectF, "EXPECTF"},
		{SectionExpectRegex, "EXPECTREGEX"},
		{SectionUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.section.String()
			if result != tt.expected {
				t.Errorf("section.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestParseKeyValueSection tests parsing key=value sections
func TestParseKeyValueSection(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name:    "single pair",
			content: "key=value",
			expected: map[string]string{
				"key": "value",
			},
		},
		{
			name:    "multiple pairs",
			content: "key1=value1\nkey2=value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:    "with spaces",
			content: "key = value",
			expected: map[string]string{
				"key": "value",
			},
		},
		{
			name:     "empty",
			content:  "",
			expected: map[string]string{},
		},
		{
			name:    "value with equals",
			content: "key=value=with=equals",
			expected: map[string]string{
				"key": "value=with=equals",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseKeyValueSection(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("parseKeyValueSection() returned %d items, want %d", len(result), len(tt.expected))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("parseKeyValueSection()[%s] = %s, want %s", k, result[k], v)
				}
			}
		})
	}
}

// TestParseArgs tests parsing command line arguments
func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "single arg",
			content:  "arg1",
			expected: []string{"arg1"},
		},
		{
			name:     "multiple args",
			content:  "arg1 arg2 arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "empty",
			content:  "",
			expected: []string{},
		},
		{
			name:     "with newlines",
			content:  "arg1\narg2",
			expected: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArgs(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("parseArgs() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("parseArgs()[%d] = %s, want %s", i, result[i], v)
				}
			}
		})
	}
}

// TestParseExtensions tests parsing required extensions
func TestParseExtensions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "single extension",
			content:  "json",
			expected: []string{"json"},
		},
		{
			name:     "multiple extensions",
			content:  "json pcre mbstring",
			expected: []string{"json", "pcre", "mbstring"},
		},
		{
			name:     "empty",
			content:  "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExtensions(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("parseExtensions() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("parseExtensions()[%d] = %s, want %s", i, result[i], v)
				}
			}
		})
	}
}

// TestParseTest_SimpleExact tests parsing a simple test with exact match
func TestParseTest_SimpleExact(t *testing.T) {
	content := `--TEST--
Simple test
--FILE--
<?php
echo "Hello World";
?>
--EXPECT--
Hello World`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if test.Description != "Simple test" {
		t.Errorf("Description = %s, want Simple test", test.Description)
	}

	expectedCode := `<?php
echo "Hello World";
?>`
	if test.Code != expectedCode {
		t.Errorf("Code = %s, want %s", test.Code, expectedCode)
	}

	if test.Expect != "Hello World" {
		t.Errorf("Expect = %s, want Hello World", test.Expect)
	}

	if test.GetExpectType() != "exact" {
		t.Errorf("GetExpectType() = %s, want exact", test.GetExpectType())
	}
}

// TestParseTest_WithFormat tests parsing test with format placeholders
func TestParseTest_WithFormat(t *testing.T) {
	content := `--TEST--
Format test
--FILE--
<?php
echo 123;
?>
--EXPECTF--
%d`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if test.ExpectF != "%d" {
		t.Errorf("ExpectF = %s, want %%d", test.ExpectF)
	}

	if test.GetExpectType() != "format" {
		t.Errorf("GetExpectType() = %s, want format", test.GetExpectType())
	}
}

// TestParseTest_WithRegex tests parsing test with regex pattern
func TestParseTest_WithRegex(t *testing.T) {
	content := `--TEST--
Regex test
--FILE--
<?php
echo "Hello";
?>
--EXPECTREGEX--
^[A-Z][a-z]+$`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if test.ExpectRegex != "^[A-Z][a-z]+$" {
		t.Errorf("ExpectRegex = %s, want ^[A-Z][a-z]+$", test.ExpectRegex)
	}

	if test.GetExpectType() != "regex" {
		t.Errorf("GetExpectType() = %s, want regex", test.GetExpectType())
	}
}

// TestParseTest_WithSkipIf tests parsing test with skip condition
func TestParseTest_WithSkipIf(t *testing.T) {
	content := `--TEST--
Skip test
--SKIPIF--
<?php if (!extension_loaded('json')) die('skip json required'); ?>
--FILE--
<?php
echo "test";
?>
--EXPECT--
test`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if !test.HasSkipCondition() {
		t.Error("HasSkipCondition() = false, want true")
	}

	expectedSkipIf := `<?php if (!extension_loaded('json')) die('skip json required'); ?>`
	if test.SkipIf != expectedSkipIf {
		t.Errorf("SkipIf = %s, want %s", test.SkipIf, expectedSkipIf)
	}
}

// TestParseTest_WithIni tests parsing test with INI settings
func TestParseTest_WithIni(t *testing.T) {
	content := `--TEST--
INI test
--INI--
display_errors=1
error_reporting=E_ALL
--FILE--
<?php
echo ini_get('display_errors');
?>
--EXPECT--
1`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if test.Ini["display_errors"] != "1" {
		t.Errorf("Ini[display_errors] = %s, want 1", test.Ini["display_errors"])
	}

	if test.Ini["error_reporting"] != "E_ALL" {
		t.Errorf("Ini[error_reporting] = %s, want E_ALL", test.Ini["error_reporting"])
	}
}

// TestParseTest_WithEnv tests parsing test with environment variables
func TestParseTest_WithEnv(t *testing.T) {
	content := `--TEST--
ENV test
--ENV--
TEST_VAR=value123
PATH=/usr/bin
--FILE--
<?php
echo getenv('TEST_VAR');
?>
--EXPECT--
value123`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	if test.Env["TEST_VAR"] != "value123" {
		t.Errorf("Env[TEST_VAR] = %s, want value123", test.Env["TEST_VAR"])
	}

	if test.Env["PATH"] != "/usr/bin" {
		t.Errorf("Env[PATH] = %s, want /usr/bin", test.Env["PATH"])
	}
}

// TestParseTest_WithArgs tests parsing test with command line arguments
func TestParseTest_WithArgs(t *testing.T) {
	content := `--TEST--
Args test
--ARGS--
arg1 arg2 arg3
--FILE--
<?php
print_r($argv);
?>
--EXPECTF--
Array%A`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	expectedArgs := []string{"arg1", "arg2", "arg3"}
	if len(test.Args) != len(expectedArgs) {
		t.Errorf("len(Args) = %d, want %d", len(test.Args), len(expectedArgs))
		return
	}

	for i, arg := range expectedArgs {
		if test.Args[i] != arg {
			t.Errorf("Args[%d] = %s, want %s", i, test.Args[i], arg)
		}
	}
}

// TestParseTest_WithExtensions tests parsing test with required extensions
func TestParseTest_WithExtensions(t *testing.T) {
	content := `--TEST--
Extensions test
--EXTENSIONS--
json pcre
--FILE--
<?php
echo "test";
?>
--EXPECT--
test`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	expectedExts := []string{"json", "pcre"}
	if len(test.Extensions) != len(expectedExts) {
		t.Errorf("len(Extensions) = %d, want %d", len(test.Extensions), len(expectedExts))
		return
	}

	for i, ext := range expectedExts {
		if test.Extensions[i] != ext {
			t.Errorf("Extensions[%d] = %s, want %s", i, test.Extensions[i], ext)
		}
	}
}

// TestParseTest_WithClean tests parsing test with cleanup code
func TestParseTest_WithClean(t *testing.T) {
	content := `--TEST--
Clean test
--FILE--
<?php
file_put_contents('test.txt', 'data');
echo "created";
?>
--CLEAN--
<?php
unlink('test.txt');
?>
--EXPECT--
created`

	tmpfile := createTempFile(t, content)
	defer os.Remove(tmpfile)

	test, err := ParseTest(tmpfile)
	if err != nil {
		t.Fatalf("ParseTest() error = %v", err)
	}

	expectedClean := `<?php
unlink('test.txt');
?>`
	if test.Clean != expectedClean {
		t.Errorf("Clean = %s, want %s", test.Clean, expectedClean)
	}
}

// TestValidate_MissingTest tests validation with missing TEST section
func TestValidate_MissingTest(t *testing.T) {
	test := &PHPTest{
		Code:   "<?php echo 'test'; ?>",
		Expect: "test",
	}

	err := test.Validate()
	if err == nil {
		t.Error("Validate() should return error for missing TEST section")
	}
}

// TestValidate_MissingFile tests validation with missing FILE section
func TestValidate_MissingFile(t *testing.T) {
	test := &PHPTest{
		Description: "Test",
		Expect:      "test",
	}

	err := test.Validate()
	if err == nil {
		t.Error("Validate() should return error for missing FILE section")
	}
}

// TestValidate_MissingExpect tests validation with missing EXPECT section
func TestValidate_MissingExpect(t *testing.T) {
	test := &PHPTest{
		Description: "Test",
		Code:        "<?php echo 'test'; ?>",
	}

	err := test.Validate()
	if err == nil {
		t.Error("Validate() should return error for missing EXPECT section")
	}
}

// TestValidate_MultipleExpect tests validation with multiple EXPECT sections
func TestValidate_MultipleExpect(t *testing.T) {
	test := &PHPTest{
		Description: "Test",
		Code:        "<?php echo 'test'; ?>",
		Expect:      "test",
		ExpectF:     "%s",
	}

	err := test.Validate()
	if err == nil {
		t.Error("Validate() should return error for multiple EXPECT sections")
	}
}

// TestMatchOutput_Exact tests exact output matching
func TestMatchOutput_Exact(t *testing.T) {
	test := &PHPTest{
		Expect: "Hello World",
	}

	tests := []struct {
		name     string
		actual   string
		expected bool
	}{
		{"exact match", "Hello World", true},
		{"mismatch", "Hello", false},
		{"case sensitive", "hello world", false},
		{"with newline", "Hello World\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, err := test.MatchOutput(tt.actual)
			if err != nil {
				t.Fatalf("MatchOutput() error = %v", err)
			}
			if matched != tt.expected {
				t.Errorf("MatchOutput(%s) = %v, want %v", tt.actual, matched, tt.expected)
			}
		})
	}
}

// TestMatchOutput_Format tests format output matching
func TestMatchOutput_Format(t *testing.T) {
	tests := []struct {
		name     string
		expectF  string
		actual   string
		expected bool
	}{
		{"string placeholder", "Hello %s", "Hello World", true},
		{"integer placeholder", "Number: %d", "Number: 123", true},
		{"float placeholder", "Pi: %f", "Pi: 3.14", true},
		{"hex placeholder", "Hex: %x", "Hex: 1a2b", true},
		{"any char", "test%amore", "test123more", true},
		{"whitespace", "hello%wworld", "hello   world", true},
		{"mismatch", "Hello %s", "Goodbye World", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &PHPTest{ExpectF: tt.expectF}
			matched, err := test.MatchOutput(tt.actual)
			if err != nil {
				t.Fatalf("MatchOutput() error = %v", err)
			}
			if matched != tt.expected {
				t.Errorf("MatchOutput(%s) with format %s = %v, want %v",
					tt.actual, tt.expectF, matched, tt.expected)
			}
		})
	}
}

// TestMatchOutput_Regex tests regex output matching
func TestMatchOutput_Regex(t *testing.T) {
	tests := []struct {
		name        string
		expectRegex string
		actual      string
		expected    bool
	}{
		{"simple pattern", "^Hello", "Hello World", true},
		{"full match", "^[A-Z][a-z]+$", "Hello", true},
		{"digit pattern", "\\d{3}", "abc123def", true},
		{"mismatch", "^Goodbye", "Hello World", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &PHPTest{ExpectRegex: tt.expectRegex}
			matched, err := test.MatchOutput(tt.actual)
			if err != nil {
				t.Fatalf("MatchOutput() error = %v", err)
			}
			if matched != tt.expected {
				t.Errorf("MatchOutput(%s) with regex %s = %v, want %v",
					tt.actual, tt.expectRegex, matched, tt.expected)
			}
		})
	}
}

// TestNormalizeLineEndings tests line ending normalization
func TestNormalizeLineEndings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"unix", "line1\nline2", "line1\nline2"},
		{"windows", "line1\r\nline2", "line1\nline2"},
		{"mac", "line1\rline2", "line1\nline2"},
		{"mixed", "line1\r\nline2\nline3\rline4", "line1\nline2\nline3\nline4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeLineEndings(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLineEndings() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestGetCodeToExecute_Direct tests getting code from FILE section
func TestGetCodeToExecute_Direct(t *testing.T) {
	test := &PHPTest{
		Code: "<?php echo 'test'; ?>",
	}

	code, err := test.GetCodeToExecute()
	if err != nil {
		t.Fatalf("GetCodeToExecute() error = %v", err)
	}

	if code != test.Code {
		t.Errorf("GetCodeToExecute() = %s, want %s", code, test.Code)
	}
}

// TestGetCodeToExecute_External tests getting code from external file
func TestGetCodeToExecute_External(t *testing.T) {
	// Create external PHP file
	dir := t.TempDir()
	externalFile := filepath.Join(dir, "external.php")
	externalContent := "<?php echo 'external'; ?>"
	if err := os.WriteFile(externalFile, []byte(externalContent), 0644); err != nil {
		t.Fatalf("Failed to create external file: %v", err)
	}

	// Create test file
	testFile := filepath.Join(dir, "test.phpt")
	test := &PHPTest{
		FilePath:     testFile,
		FileExternal: "external.php",
	}

	code, err := test.GetCodeToExecute()
	if err != nil {
		t.Fatalf("GetCodeToExecute() error = %v", err)
	}

	if code != externalContent {
		t.Errorf("GetCodeToExecute() = %s, want %s", code, externalContent)
	}
}

// Helper function to create temporary test file
func createTempFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "test-*.phpt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}
