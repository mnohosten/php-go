// Package migrate provides tools for analyzing PHP applications and migrating them to PHP-Go.
package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
)

// CompatibilityAnalyzer analyzes PHP code for compatibility with PHP-Go
type CompatibilityAnalyzer struct {
	// Configuration
	RootPath string
	Exclude  []string // Patterns to exclude (e.g., "vendor", "tests")

	// Analysis results
	Files              []FileAnalysis
	TotalFiles         int
	SuccessfulParses   int
	FailedParses       int
	ExtensionUsage     map[string]int // extension name -> count
	FunctionUsage      map[string]int // function name -> count
	LanguageFeatures   map[string]int // feature name -> count
	Issues             []Issue
	ParseErrors        []ParseError
	UnsupportedExts    map[string]int // extension -> count
	UnsupportedFuncs   map[string]int // function -> count
}

// FileAnalysis contains analysis results for a single file
type FileAnalysis struct {
	Path             string
	Size             int64
	Lines            int
	ParseSuccess     bool
	ParseError       string
	ExtensionsUsed   []string
	FunctionsUsed    []string
	LanguageFeatures []string
	Issues           []Issue
}

// Issue represents a compatibility issue
type Issue struct {
	File        string
	Line        int
	Column      int
	Severity    Severity
	Category    Category
	Message     string
	Suggestion  string
}

// Severity levels
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Category types
type Category int

const (
	CategoryExtension Category = iota
	CategoryFunction
	CategoryLanguageFeature
	CategorySyntax
	CategoryPerformance
	CategorySecurity
)

func (c Category) String() string {
	switch c {
	case CategoryExtension:
		return "Extension"
	case CategoryFunction:
		return "Function"
	case CategoryLanguageFeature:
		return "Language Feature"
	case CategorySyntax:
		return "Syntax"
	case CategoryPerformance:
		return "Performance"
	case CategorySecurity:
		return "Security"
	default:
		return "Unknown"
	}
}

// ParseError represents a parsing error
type ParseError struct {
	File    string
	Line    int
	Column  int
	Message string
}

// NewCompatibilityAnalyzer creates a new analyzer
func NewCompatibilityAnalyzer(rootPath string) *CompatibilityAnalyzer {
	return &CompatibilityAnalyzer{
		RootPath:         rootPath,
		Exclude:          []string{},
		Files:            []FileAnalysis{},
		ExtensionUsage:   make(map[string]int),
		FunctionUsage:    make(map[string]int),
		LanguageFeatures: make(map[string]int),
		Issues:           []Issue{},
		ParseErrors:      []ParseError{},
		UnsupportedExts:  make(map[string]int),
		UnsupportedFuncs: make(map[string]int),
	}
}

// SetExcludePatterns sets patterns to exclude from analysis
func (ca *CompatibilityAnalyzer) SetExcludePatterns(patterns []string) {
	ca.Exclude = patterns
}

// Analyze performs the compatibility analysis
func (ca *CompatibilityAnalyzer) Analyze() error {
	return filepath.Walk(ca.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Check if this directory should be excluded
			for _, pattern := range ca.Exclude {
				if strings.Contains(path, pattern) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Only analyze PHP files
		if !strings.HasSuffix(path, ".php") {
			return nil
		}

		// Check if file should be excluded
		for _, pattern := range ca.Exclude {
			if strings.Contains(path, pattern) {
				return nil
			}
		}

		// Analyze the file
		ca.analyzeFile(path, info.Size())
		return nil
	})
}

// analyzeFile analyzes a single PHP file
func (ca *CompatibilityAnalyzer) analyzeFile(path string, size int64) {
	ca.TotalFiles++

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		ca.FailedParses++
		ca.ParseErrors = append(ca.ParseErrors, ParseError{
			File:    path,
			Message: fmt.Sprintf("Failed to read file: %v", err),
		})
		return
	}

	lines := strings.Count(string(content), "\n") + 1

	// Parse file
	l := lexer.New(string(content), path)
	p := parser.New(l)
	_ = p.ParseProgram()

	analysis := FileAnalysis{
		Path:             path,
		Size:             size,
		Lines:            lines,
		ParseSuccess:     len(p.Errors()) == 0,
		ExtensionsUsed:   []string{},
		FunctionsUsed:    []string{},
		LanguageFeatures: []string{},
		Issues:           []Issue{},
	}

	if len(p.Errors()) > 0 {
		ca.FailedParses++
		analysis.ParseError = strings.Join(p.Errors(), "; ")
		for _, errMsg := range p.Errors() {
			ca.ParseErrors = append(ca.ParseErrors, ParseError{
				File:    path,
				Message: errMsg,
			})
		}
	} else {
		ca.SuccessfulParses++
		// Analyze the file content for function calls and features
		ca.analyzeContent(string(content), path, &analysis)
	}

	ca.Files = append(ca.Files, analysis)
}

// analyzeContent analyzes PHP file content for function calls and features using pattern matching
func (ca *CompatibilityAnalyzer) analyzeContent(content string, filePath string, analysis *FileAnalysis) {
	// Check for common function patterns
	funcs := []string{
		// Database
		"mysqli_connect", "mysqli_query", "mysqli_close", "mysqli_fetch_assoc",
		"pg_connect", "pg_query", "pg_close",
		// XML
		"simplexml_load_file", "simplexml_load_string",
		"xml_parser_create", "xml_parse",
		// DOM
		"DOMDocument", "dom_import_simplexml",
		// Images
		"imagecreatetruecolor", "imagecreate", "imagejpeg", "imagepng",
		// Curl
		"curl_init", "curl_exec", "curl_close", "curl_setopt",
		// Compression
		"gzopen", "gzwrite", "gzclose", "gzcompress",
		// Sockets
		"socket_create", "socket_bind", "socket_listen",
		// Multi-byte
		"mb_strlen", "mb_substr", "mb_convert_encoding",
		// Character type
		"ctype_alpha", "ctype_digit", "ctype_alnum",
		// Filter
		"filter_var", "filter_input",
		// OpenSSL
		"openssl_encrypt", "openssl_decrypt",
		// Cache
		"apcu_fetch", "apcu_store",
		"opcache_reset", "opcache_compile_file",
	}

	for _, funcName := range funcs {
		if strings.Contains(content, funcName) {
			ca.FunctionUsage[funcName]++
			analysis.FunctionsUsed = append(analysis.FunctionsUsed, funcName)

			ext := getExtensionForFunction(funcName)
			if ext != "" {
				ca.ExtensionUsage[ext]++
				if !contains(analysis.ExtensionsUsed, ext) {
					analysis.ExtensionsUsed = append(analysis.ExtensionsUsed, ext)
				}

				if !isExtensionSupported(ext) {
					ca.UnsupportedExts[ext]++
					ca.UnsupportedFuncs[funcName]++
					issue := Issue{
						File:       filePath,
						Severity:   SeverityError,
						Category:   CategoryExtension,
						Message:    fmt.Sprintf("Function '%s' from unsupported extension '%s'", funcName, ext),
						Suggestion: fmt.Sprintf("Extension '%s' is not yet available in PHP-Go. Consider using alternative approaches or implementing a custom extension.", ext),
					}
					analysis.Issues = append(analysis.Issues, issue)
					ca.Issues = append(ca.Issues, issue)
				}
			}
		}
	}

	// Check for language features
	features := map[string]struct {
		pattern    string
		supported  bool
		suggestion string
	}{
		"generators": {
			pattern:    "yield ",
			supported:  false,
			suggestion: "Generators are planned for Phase 9. Consider refactoring to use arrays or iterators.",
		},
		"arrow_functions": {
			pattern:    " => ",
			supported:  false,
			suggestion: "Arrow functions are planned for Phase 9. Use regular closures instead.",
		},
		"attributes": {
			pattern:    "#[",
			supported:  false,
			suggestion: "Attributes are planned for Phase 9. Use docblock comments for metadata.",
		},
		"match_expressions": {
			pattern:    "match(",
			supported:  false,
			suggestion: "Match expressions are planned for Phase 9. Use switch statements instead.",
		},
		"enums": {
			pattern:    "enum ",
			supported:  true,
			suggestion: "",
		},
		"readonly": {
			pattern:    "readonly ",
			supported:  true,
			suggestion: "",
		},
	}

	for feature, info := range features {
		if strings.Contains(content, info.pattern) {
			ca.LanguageFeatures[feature]++
			if !contains(analysis.LanguageFeatures, feature) {
				analysis.LanguageFeatures = append(analysis.LanguageFeatures, feature)
			}

			if !info.supported {
				issue := Issue{
					File:       filePath,
					Severity:   SeverityWarning,
					Category:   CategoryLanguageFeature,
					Message:    fmt.Sprintf("Potentially unsupported language feature: %s", feature),
					Suggestion: info.suggestion,
				}
				analysis.Issues = append(analysis.Issues, issue)
				ca.Issues = append(ca.Issues, issue)
			}
		}
	}
}

// getExtensionForFunction returns the extension name for a PHP function
func getExtensionForFunction(funcName string) string {
	funcName = strings.ToLower(funcName)

	// Database extensions
	if strings.HasPrefix(funcName, "mysqli_") {
		return "mysqli"
	}
	if strings.HasPrefix(funcName, "pg_") {
		return "pgsql"
	}
	if strings.HasPrefix(funcName, "sqlite_") || strings.Contains(funcName, "sqlite") {
		return "sqlite3"
	}

	// XML extensions
	if strings.HasPrefix(funcName, "xml_") {
		return "xml"
	}
	if strings.HasPrefix(funcName, "simplexml_") {
		return "simplexml"
	}
	if strings.Contains(funcName, "dom") || strings.HasPrefix(funcName, "dom_") {
		return "dom"
	}

	// Compression
	if strings.HasPrefix(funcName, "gz") || strings.Contains(funcName, "gzip") {
		return "zlib"
	}
	if strings.Contains(funcName, "bz2") {
		return "bz2"
	}
	if strings.HasPrefix(funcName, "zip_") {
		return "zip"
	}

	// Networking
	if strings.HasPrefix(funcName, "curl_") {
		return "curl"
	}
	if strings.HasPrefix(funcName, "socket_") {
		return "sockets"
	}
	if strings.HasPrefix(funcName, "ftp_") {
		return "ftp"
	}

	// Images
	if strings.HasPrefix(funcName, "image") || strings.Contains(funcName, "gd_") {
		return "gd"
	}
	if strings.Contains(funcName, "imagick") {
		return "imagick"
	}

	// Caching
	if strings.HasPrefix(funcName, "apcu_") {
		return "apcu"
	}
	if strings.HasPrefix(funcName, "opcache_") {
		return "opcache"
	}
	if strings.HasPrefix(funcName, "memcache") {
		return "memcache"
	}

	// Encryption
	if strings.HasPrefix(funcName, "openssl_") {
		return "openssl"
	}
	if strings.Contains(funcName, "mcrypt") {
		return "mcrypt"
	}

	// Character type
	if strings.HasPrefix(funcName, "ctype_") {
		return "ctype"
	}

	// Multi-byte string
	if strings.HasPrefix(funcName, "mb_") {
		return "mbstring"
	}

	// Internationalization
	if strings.HasPrefix(funcName, "iconv") {
		return "iconv"
	}
	if strings.HasPrefix(funcName, "intl_") {
		return "intl"
	}

	// Filter
	if strings.HasPrefix(funcName, "filter_") {
		return "filter"
	}

	return ""
}

// isExtensionSupported checks if an extension is supported by PHP-Go
func isExtensionSupported(ext string) bool {
	supported := map[string]bool{
		"core":   true,
		"json":   true,
		"hash":   true, // Partial
		"pcre":   true, // Partial
		"date":   true, // Partial
		"spl":    true, // Partial
		"string": true,
		"array":  true,
	}
	return supported[ext]
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GenerateReport generates a comprehensive compatibility report
func (ca *CompatibilityAnalyzer) GenerateReport() string {
	var sb strings.Builder

	sb.WriteString("===========================================\n")
	sb.WriteString("PHP-Go Compatibility Analysis Report\n")
	sb.WriteString("===========================================\n\n")

	// Summary
	sb.WriteString("SUMMARY\n")
	sb.WriteString("-------\n")
	sb.WriteString(fmt.Sprintf("Total Files Analyzed:     %d\n", ca.TotalFiles))
	sb.WriteString(fmt.Sprintf("Successfully Parsed:      %d (%.2f%%)\n", ca.SuccessfulParses, float64(ca.SuccessfulParses)/float64(ca.TotalFiles)*100))
	sb.WriteString(fmt.Sprintf("Parse Failures:           %d (%.2f%%)\n", ca.FailedParses, float64(ca.FailedParses)/float64(ca.TotalFiles)*100))
	sb.WriteString(fmt.Sprintf("Total Issues Found:       %d\n", len(ca.Issues)))
	sb.WriteString(fmt.Sprintf("Unique Extensions Used:   %d\n", len(ca.ExtensionUsage)))
	sb.WriteString(fmt.Sprintf("Unsupported Extensions:   %d\n", len(ca.UnsupportedExts)))
	sb.WriteString(fmt.Sprintf("Unsupported Functions:    %d\n", len(ca.UnsupportedFuncs)))
	sb.WriteString("\n")

	// Issue severity breakdown
	severityCounts := map[Severity]int{}
	for _, issue := range ca.Issues {
		severityCounts[issue.Severity]++
	}
	sb.WriteString("ISSUE SEVERITY BREAKDOWN\n")
	sb.WriteString("------------------------\n")
	sb.WriteString(fmt.Sprintf("Critical:  %d\n", severityCounts[SeverityCritical]))
	sb.WriteString(fmt.Sprintf("Error:     %d\n", severityCounts[SeverityError]))
	sb.WriteString(fmt.Sprintf("Warning:   %d\n", severityCounts[SeverityWarning]))
	sb.WriteString(fmt.Sprintf("Info:      %d\n", severityCounts[SeverityInfo]))
	sb.WriteString("\n")

	// Unsupported extensions
	if len(ca.UnsupportedExts) > 0 {
		sb.WriteString("UNSUPPORTED EXTENSIONS\n")
		sb.WriteString("----------------------\n")

		// Sort by usage count
		type extCount struct {
			name  string
			count int
		}
		var exts []extCount
		for name, count := range ca.UnsupportedExts {
			exts = append(exts, extCount{name, count})
		}
		sort.Slice(exts, func(i, j int) bool {
			return exts[i].count > exts[j].count
		})

		for _, ext := range exts {
			sb.WriteString(fmt.Sprintf("  %-20s %d usage(s)\n", ext.name, ext.count))
		}
		sb.WriteString("\n")
	}

	// Top unsupported functions
	if len(ca.UnsupportedFuncs) > 0 {
		sb.WriteString("TOP 10 UNSUPPORTED FUNCTIONS\n")
		sb.WriteString("----------------------------\n")

		// Sort by usage count
		type funcCount struct {
			name  string
			count int
		}
		var funcs []funcCount
		for name, count := range ca.UnsupportedFuncs {
			funcs = append(funcs, funcCount{name, count})
		}
		sort.Slice(funcs, func(i, j int) bool {
			return funcs[i].count > funcs[j].count
		})

		limit := 10
		if len(funcs) < limit {
			limit = len(funcs)
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("  %-30s %d usage(s)\n", funcs[i].name, funcs[i].count))
		}
		sb.WriteString("\n")
	}

	// Language features
	if len(ca.LanguageFeatures) > 0 {
		sb.WriteString("LANGUAGE FEATURES USED\n")
		sb.WriteString("----------------------\n")

		// Sort by name
		var features []string
		for feature := range ca.LanguageFeatures {
			features = append(features, feature)
		}
		sort.Strings(features)

		for _, feature := range features {
			count := ca.LanguageFeatures[feature]
			supported := "✓"
			if isFeatureSupported(feature) {
				sb.WriteString(fmt.Sprintf("  %s %-25s %d usage(s)\n", supported, feature, count))
			} else {
				supported = "✗"
				sb.WriteString(fmt.Sprintf("  %s %-25s %d usage(s) [UNSUPPORTED]\n", supported, feature, count))
			}
		}
		sb.WriteString("\n")
	}

	// Parse errors (top 10)
	if len(ca.ParseErrors) > 0 {
		sb.WriteString("TOP 10 PARSE ERRORS\n")
		sb.WriteString("-------------------\n")
		limit := 10
		if len(ca.ParseErrors) < limit {
			limit = len(ca.ParseErrors)
		}
		for i := 0; i < limit; i++ {
			err := ca.ParseErrors[i]
			sb.WriteString(fmt.Sprintf("  %s: %s\n", err.File, err.Message))
		}
		if len(ca.ParseErrors) > 10 {
			sb.WriteString(fmt.Sprintf("  ... and %d more parse errors\n", len(ca.ParseErrors)-10))
		}
		sb.WriteString("\n")
	}

	// Recommendations
	sb.WriteString("RECOMMENDATIONS\n")
	sb.WriteString("---------------\n")

	parseRate := float64(ca.SuccessfulParses) / float64(ca.TotalFiles) * 100
	if parseRate < 50 {
		sb.WriteString("⚠ CRITICAL: Less than 50% of files parse successfully.\n")
		sb.WriteString("  Action: Review parse errors and update to PHP 8.2+ syntax before migrating.\n")
	} else if parseRate < 80 {
		sb.WriteString("⚠ WARNING: Less than 80% of files parse successfully.\n")
		sb.WriteString("  Action: Address parse errors to improve compatibility.\n")
	} else {
		sb.WriteString(fmt.Sprintf("✓ GOOD: Parse success rate is acceptable (%.2f%%).\n", parseRate))
	}

	if len(ca.UnsupportedExts) > 0 {
		sb.WriteString(fmt.Sprintf("⚠ WARNING: %d unsupported extensions detected.\n", len(ca.UnsupportedExts)))
		sb.WriteString("  Action: Review extension usage and plan for alternatives or custom implementations.\n")
	}

	if len(ca.Issues) > 100 {
		sb.WriteString(fmt.Sprintf("⚠ WARNING: %d compatibility issues found.\n", len(ca.Issues)))
		sb.WriteString("  Action: Review and prioritize issues before migration.\n")
	}

	sb.WriteString("\n")
	sb.WriteString("For detailed migration guidance, see: docs/user-guide/migration.md\n")
	sb.WriteString("===========================================\n")

	return sb.String()
}

// isFeatureSupported checks if a language feature is supported
func isFeatureSupported(feature string) bool {
	supported := map[string]bool{
		"named_arguments": true,
		"enums":           true,
		"readonly":        true,
		"union_types":     true,
		"generators":      false,
		"arrow_functions": false,
		"attributes":      false,
		"match_expressions": false,
		"throw_expressions": false,
	}
	return supported[feature]
}
