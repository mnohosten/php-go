package phptest

import (
	"path/filepath"
	"strings"
)

// ============================================================================
// Test Categorization
// ============================================================================
//
// Categorizes tests by type and provides filtering capabilities

// TestCategory represents a test category
type TestCategory string

const (
	CategoryLanguage TestCategory = "language"
	CategoryStdLib   TestCategory = "stdlib"
	CategorySyntax   TestCategory = "syntax"
	CategoryOOP      TestCategory = "oop"
	CategoryError    TestCategory = "error"
	CategoryFunction TestCategory = "function"
	CategoryArray    TestCategory = "array"
	CategoryString   TestCategory = "string"
	CategoryMath     TestCategory = "math"
	CategoryFile     TestCategory = "file"
	CategoryJSON     TestCategory = "json"
	CategoryPCRE     TestCategory = "pcre"
	CategoryMB       TestCategory = "mbstring"
	CategoryDate     TestCategory = "date"
	CategoryOther    TestCategory = "other"
)

// String returns the string representation of a category
func (c TestCategory) String() string {
	return string(c)
}

// TestFilter provides filtering capabilities for tests
type TestFilter struct {
	// Category filters
	Categories []TestCategory

	// Path filters
	PathPattern string

	// Extension filters
	RequiredExtensions []string
	ExcludedExtensions []string

	// Status filters
	IncludeSkipped bool
}

// NewTestFilter creates a new test filter
func NewTestFilter() *TestFilter {
	return &TestFilter{
		Categories:         []TestCategory{},
		RequiredExtensions: []string{},
		ExcludedExtensions: []string{},
		IncludeSkipped:     true,
	}
}

// Matches checks if a test matches the filter criteria
func (f *TestFilter) Matches(test *PHPTest) bool {
	// Category filter
	if len(f.Categories) > 0 {
		category := CategorizeTest(test)
		found := false
		for _, c := range f.Categories {
			if c == category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Path pattern filter
	if f.PathPattern != "" {
		matched, _ := filepath.Match(f.PathPattern, test.FilePath)
		if !matched {
			return false
		}
	}

	// Required extensions filter
	if len(f.RequiredExtensions) > 0 {
		for _, reqExt := range f.RequiredExtensions {
			found := false
			for _, testExt := range test.Extensions {
				if testExt == reqExt {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Excluded extensions filter
	if len(f.ExcludedExtensions) > 0 {
		for _, exclExt := range f.ExcludedExtensions {
			for _, testExt := range test.Extensions {
				if testExt == exclExt {
					return false
				}
			}
		}
	}

	return true
}

// FilterTests filters a list of tests
func (f *TestFilter) FilterTests(tests []*PHPTest) []*PHPTest {
	filtered := make([]*PHPTest, 0)

	for _, test := range tests {
		if f.Matches(test) {
			filtered = append(filtered, test)
		}
	}

	return filtered
}

// CategorizeTest determines the category of a test based on its path and content
func CategorizeTest(test *PHPTest) TestCategory {
	if test.FilePath == "" {
		return CategoryOther
	}

	// Normalize path separators
	path := filepath.ToSlash(test.FilePath)
	dir := filepath.Dir(path)

	// Check directory structure for PHP test suite organization
	if strings.Contains(path, "/lang/") || strings.Contains(path, "/language/") {
		return CategoryLanguage
	}

	if strings.Contains(path, "/oop/") || strings.Contains(dir, "/classes/") {
		return CategoryOOP
	}

	if strings.Contains(path, "/syntax/") {
		return CategorySyntax
	}

	if strings.Contains(path, "/error/") || strings.Contains(dir, "/errors/") {
		return CategoryError
	}

	// Check for standard library categories
	if strings.Contains(path, "/array/") {
		return CategoryArray
	}

	if strings.Contains(path, "/string/") {
		return CategoryString
	}

	if strings.Contains(path, "/math/") {
		return CategoryMath
	}

	if strings.Contains(path, "/file/") || strings.Contains(dir, "/filesystem/") {
		return CategoryFile
	}

	if strings.Contains(path, "/json/") {
		return CategoryJSON
	}

	if strings.Contains(path, "/pcre/") || strings.Contains(dir, "/regex/") {
		return CategoryPCRE
	}

	if strings.Contains(path, "/mbstring/") {
		return CategoryMB
	}

	if strings.Contains(path, "/date/") || strings.Contains(dir, "/datetime/") {
		return CategoryDate
	}

	if strings.Contains(path, "/func/") || strings.Contains(dir, "/functions/") {
		return CategoryFunction
	}

	// Check description for category hints
	desc := strings.ToLower(test.Description)

	if strings.Contains(desc, "array") {
		return CategoryArray
	}

	if strings.Contains(desc, "string") {
		return CategoryString
	}

	if strings.Contains(desc, "class") || strings.Contains(desc, "object") ||
		strings.Contains(desc, "interface") || strings.Contains(desc, "trait") {
		return CategoryOOP
	}

	if strings.Contains(desc, "function") {
		return CategoryFunction
	}

	if strings.Contains(desc, "syntax") {
		return CategorySyntax
	}

	if strings.Contains(desc, "error") || strings.Contains(desc, "exception") {
		return CategoryError
	}

	// Default to stdlib or other
	if strings.Contains(path, "/ext/") {
		return CategoryStdLib
	}

	return CategoryOther
}

// GroupTestsByCategory groups tests by category
func GroupTestsByCategory(tests []*PHPTest) map[TestCategory][]*PHPTest {
	grouped := make(map[TestCategory][]*PHPTest)

	for _, test := range tests {
		category := CategorizeTest(test)
		grouped[category] = append(grouped[category], test)
	}

	return grouped
}

// GroupResultsByCategory groups test results by category
func GroupResultsByCategory(results []*TestResult) map[TestCategory][]*TestResult {
	grouped := make(map[TestCategory][]*TestResult)

	for _, result := range results {
		category := CategorizeTest(result.Test)
		grouped[category] = append(grouped[category], result)
	}

	return grouped
}

// CountByCategory counts tests in each category
func CountByCategory(tests []*PHPTest) map[TestCategory]int {
	counts := make(map[TestCategory]int)

	for _, test := range tests {
		category := CategorizeTest(test)
		counts[category]++
	}

	return counts
}
