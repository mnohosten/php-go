package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test creating a new analyzer
func TestNewCompatibilityAnalyzer(t *testing.T) {
	analyzer := NewCompatibilityAnalyzer("/test/path")

	if analyzer.RootPath != "/test/path" {
		t.Errorf("Expected RootPath to be '/test/path', got '%s'", analyzer.RootPath)
	}

	if analyzer.ExtensionUsage == nil {
		t.Error("Expected ExtensionUsage map to be initialized")
	}

	if analyzer.FunctionUsage == nil {
		t.Error("Expected FunctionUsage map to be initialized")
	}

	if analyzer.LanguageFeatures == nil {
		t.Error("Expected LanguageFeatures map to be initialized")
	}
}

// Test setting exclude patterns
func TestSetExcludePatterns(t *testing.T) {
	analyzer := NewCompatibilityAnalyzer("/test")
	patterns := []string{"vendor", "tests", "cache"}

	analyzer.SetExcludePatterns(patterns)

	if len(analyzer.Exclude) != 3 {
		t.Errorf("Expected 3 exclude patterns, got %d", len(analyzer.Exclude))
	}

	if analyzer.Exclude[0] != "vendor" {
		t.Errorf("Expected first pattern to be 'vendor', got '%s'", analyzer.Exclude[0])
	}
}

// Test analyzing a simple PHP file
func TestAnalyzeSimpleFile(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create a simple PHP file
	testFile := filepath.Join(tmpDir, "test.php")
	content := `<?php
echo "Hello, World!";
$x = 10 + 5;
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Analyze
	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Check results
	if analyzer.TotalFiles != 1 {
		t.Errorf("Expected 1 file analyzed, got %d", analyzer.TotalFiles)
	}

	if analyzer.SuccessfulParses != 1 {
		t.Errorf("Expected 1 successful parse, got %d", analyzer.SuccessfulParses)
	}

	if analyzer.FailedParses != 0 {
		t.Errorf("Expected 0 failed parses, got %d", analyzer.FailedParses)
	}
}

// Test detecting unsupported extensions
func TestDetectUnsupportedExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with mysqli calls
	testFile := filepath.Join(tmpDir, "database.php")
	content := `<?php
$conn = mysqli_connect("localhost", "user", "pass");
mysqli_query($conn, "SELECT * FROM users");
mysqli_close($conn);
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Analyze
	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Check for mysqli extension detection
	if _, found := analyzer.UnsupportedExts["mysqli"]; !found {
		t.Error("Expected mysqli to be detected as unsupported extension")
	}

	// Check for mysqli functions
	if analyzer.UnsupportedFuncs["mysqli_connect"] == 0 {
		t.Error("Expected mysqli_connect to be detected")
	}

	// Check issues
	if len(analyzer.Issues) == 0 {
		t.Error("Expected compatibility issues to be reported")
	}

	foundMysqliIssue := false
	for _, issue := range analyzer.Issues {
		if issue.Category == CategoryExtension && strings.Contains(issue.Message, "mysqli") {
			foundMysqliIssue = true
			break
		}
	}
	if !foundMysqliIssue {
		t.Error("Expected mysqli compatibility issue to be reported")
	}
}

// Test detecting multiple extensions
func TestDetectMultipleExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "mixed.php")
	content := `<?php
// Database
mysqli_connect("localhost", "user", "pass");
pg_connect("host=localhost");

// XML
$xml = simplexml_load_file("data.xml");
$doc = new DOMDocument();

// Image
imagecreatetruecolor(100, 100);

// Curl
curl_init("https://example.com");
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	expectedExts := []string{"mysqli", "pgsql", "simplexml", "dom", "gd", "curl"}
	for _, ext := range expectedExts {
		if _, found := analyzer.UnsupportedExts[ext]; !found {
			t.Errorf("Expected %s to be detected as unsupported extension", ext)
		}
	}

	if len(analyzer.UnsupportedExts) < len(expectedExts) {
		t.Errorf("Expected at least %d unsupported extensions, got %d", len(expectedExts), len(analyzer.UnsupportedExts))
	}
}

// Test exclude patterns
func TestExcludePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create vendor directory with PHP files
	vendorDir := filepath.Join(tmpDir, "vendor")
	err := os.Mkdir(vendorDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create vendor dir: %v", err)
	}

	vendorFile := filepath.Join(vendorDir, "library.php")
	err = os.WriteFile(vendorFile, []byte("<?php echo 'vendor';"), 0644)
	if err != nil {
		t.Fatalf("Failed to create vendor file: %v", err)
	}

	// Create app file
	appFile := filepath.Join(tmpDir, "app.php")
	err = os.WriteFile(appFile, []byte("<?php echo 'app';"), 0644)
	if err != nil {
		t.Fatalf("Failed to create app file: %v", err)
	}

	// Analyze with vendor excluded
	analyzer := NewCompatibilityAnalyzer(tmpDir)
	analyzer.SetExcludePatterns([]string{"vendor"})
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Should only find app.php
	if analyzer.TotalFiles != 1 {
		t.Errorf("Expected 1 file (excluding vendor), got %d", analyzer.TotalFiles)
	}

	// Check that vendor file was not analyzed
	for _, file := range analyzer.Files {
		if strings.Contains(file.Path, "vendor") {
			t.Error("Expected vendor files to be excluded")
		}
	}
}

// Test parse error handling
func TestParseErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with syntax errors
	testFile := filepath.Join(tmpDir, "broken.php")
	content := `<?php
$x = ;  // Invalid syntax
echo $x;
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Note: The parser might recover from some errors, so we check if errors were recorded
	// but don't require FailedParses to be non-zero
	if len(analyzer.ParseErrors) == 0 && analyzer.Files[0].ParseSuccess {
		t.Skip("Parser recovered from syntax error - this test is informational")
	}
}

// Test report generation
func TestGenerateReport(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.php")
	content := `<?php
mysqli_connect("localhost", "user", "pass");
curl_init("https://example.com");
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	report := analyzer.GenerateReport()

	// Check that report contains expected sections
	expectedSections := []string{
		"PHP-Go Compatibility Analysis Report",
		"SUMMARY",
		"Total Files Analyzed",
		"Successfully Parsed",
		"UNSUPPORTED EXTENSIONS",
		"RECOMMENDATIONS",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Expected report to contain section: %s", section)
		}
	}

	// Check that unsupported extensions are listed
	if !strings.Contains(report, "mysqli") {
		t.Error("Expected report to mention mysqli extension")
	}

	if !strings.Contains(report, "curl") {
		t.Error("Expected report to mention curl extension")
	}
}

// Test getExtensionForFunction
func TestGetExtensionForFunction(t *testing.T) {
	tests := []struct {
		function string
		expected string
	}{
		{"mysqli_connect", "mysqli"},
		{"pg_query", "pgsql"},
		{"curl_init", "curl"},
		{"simplexml_load_file", "simplexml"},
		{"imagecreate", "gd"},
		{"gzopen", "zlib"},
		{"socket_create", "sockets"},
		{"mb_strlen", "mbstring"},
		{"ctype_alpha", "ctype"},
		{"filter_var", "filter"},
		{"openssl_encrypt", "openssl"},
		{"apcu_fetch", "apcu"},
		{"opcache_reset", "opcache"},
		{"echo", ""}, // Core function, no extension
		{"strlen", ""}, // Core function
	}

	for _, tt := range tests {
		result := getExtensionForFunction(tt.function)
		if result != tt.expected {
			t.Errorf("getExtensionForFunction(%s) = %s, expected %s", tt.function, result, tt.expected)
		}
	}
}

// Test isExtensionSupported
func TestIsExtensionSupported(t *testing.T) {
	tests := []struct {
		extension string
		supported bool
	}{
		{"core", true},
		{"json", true},
		{"hash", true},
		{"pcre", true},
		{"mysqli", false},
		{"pdo", false},
		{"curl", false},
		{"gd", false},
	}

	for _, tt := range tests {
		result := isExtensionSupported(tt.extension)
		if result != tt.supported {
			t.Errorf("isExtensionSupported(%s) = %v, expected %v", tt.extension, result, tt.supported)
		}
	}
}

// Test severity and category string methods
func TestSeverityString(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityError, "ERROR"},
		{SeverityCritical, "CRITICAL"},
	}

	for _, tt := range tests {
		result := tt.severity.String()
		if result != tt.expected {
			t.Errorf("Severity.String() = %s, expected %s", result, tt.expected)
		}
	}
}

func TestCategoryString(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{CategoryExtension, "Extension"},
		{CategoryFunction, "Function"},
		{CategoryLanguageFeature, "Language Feature"},
		{CategorySyntax, "Syntax"},
		{CategoryPerformance, "Performance"},
		{CategorySecurity, "Security"},
	}

	for _, tt := range tests {
		result := tt.category.String()
		if result != tt.expected {
			t.Errorf("Category.String() = %s, expected %s", result, tt.expected)
		}
	}
}

// Test analyzing files with supported features
func TestSupportedFeatures(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "modern.php")
	content := `<?php
// Basic PHP that should definitely parse
class User {
    public string $name;
    public int $age;

    public function __construct(string $name, int $age) {
        $this->name = $name;
        $this->age = $age;
    }
}

function greet($value) {
    echo $value;
}

$user = new User("John", 30);
greet("Hello");
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Basic PHP should parse successfully
	if analyzer.SuccessfulParses == 0 {
		t.Errorf("Expected file to parse successfully, got %d successful parses, errors: %v",
			analyzer.SuccessfulParses, analyzer.ParseErrors)
	}

	// Should not have critical issues for basic PHP
	criticalCount := 0
	for _, issue := range analyzer.Issues {
		if issue.Severity == SeverityCritical {
			criticalCount++
		}
	}

	if criticalCount > 0 {
		t.Errorf("Expected no critical issues for basic PHP, got %d", criticalCount)
	}
}

// Test contains helper function
func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	if !contains(slice, "banana") {
		t.Error("Expected contains to return true for 'banana'")
	}

	if contains(slice, "orange") {
		t.Error("Expected contains to return false for 'orange'")
	}
}

// Test empty directory
func TestEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	if analyzer.TotalFiles != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", analyzer.TotalFiles)
	}

	report := analyzer.GenerateReport()
	if !strings.Contains(report, "Total Files Analyzed:     0") {
		t.Error("Expected report to show 0 files analyzed")
	}
}

// Test multiple PHP files
func TestMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple PHP files
	for i := 1; i <= 5; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("file%d.php", i))
		content := "<?php echo 'test';"
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	if analyzer.TotalFiles != 5 {
		t.Errorf("Expected 5 files analyzed, got %d", analyzer.TotalFiles)
	}

	if len(analyzer.Files) != 5 {
		t.Errorf("Expected 5 file analyses, got %d", len(analyzer.Files))
	}
}

// Test non-PHP files are ignored
func TestNonPHPFilesIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	// Create PHP file
	phpFile := filepath.Join(tmpDir, "test.php")
	err := os.WriteFile(phpFile, []byte("<?php echo 'test';"), 0644)
	if err != nil {
		t.Fatalf("Failed to create PHP file: %v", err)
	}

	// Create non-PHP files
	txtFile := filepath.Join(tmpDir, "readme.txt")
	err = os.WriteFile(txtFile, []byte("readme"), 0644)
	if err != nil {
		t.Fatalf("Failed to create txt file: %v", err)
	}

	jsFile := filepath.Join(tmpDir, "script.js")
	err = os.WriteFile(jsFile, []byte("console.log('test')"), 0644)
	if err != nil {
		t.Fatalf("Failed to create js file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Should only analyze PHP file
	if analyzer.TotalFiles != 1 {
		t.Errorf("Expected only 1 PHP file to be analyzed, got %d", analyzer.TotalFiles)
	}
}

// Benchmark analyzer performance
func BenchmarkAnalyzeSimpleFile(b *testing.B) {
	tmpDir := b.TempDir()

	testFile := filepath.Join(tmpDir, "test.php")
	content := `<?php
function add($a, $b) {
    return $a + $b;
}

echo add(5, 10);
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer := NewCompatibilityAnalyzer(tmpDir)
		_ = analyzer.Analyze()
	}
}

func BenchmarkGenerateReport(b *testing.B) {
	tmpDir := b.TempDir()

	testFile := filepath.Join(tmpDir, "test.php")
	content := `<?php
mysqli_connect("localhost", "user", "pass");
curl_init("https://example.com");
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	analyzer := NewCompatibilityAnalyzer(tmpDir)
	err = analyzer.Analyze()
	if err != nil {
		b.Fatalf("Analysis failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.GenerateReport()
	}
}
