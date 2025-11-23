package phptest

import (
	"testing"
)

// TestCategorizeTest tests test categorization
func TestCategorizeTest(t *testing.T) {
	tests := []struct {
		name     string
		test     *PHPTest
		expected TestCategory
	}{
		{
			name: "language test by path",
			test: &PHPTest{
				FilePath:    "/tests/lang/basic.phpt",
				Description: "Basic language test",
			},
			expected: CategoryLanguage,
		},
		{
			name: "OOP test by path",
			test: &PHPTest{
				FilePath:    "/tests/oop/classes/basic.phpt",
				Description: "Basic class test",
			},
			expected: CategoryOOP,
		},
		{
			name: "array test by path",
			test: &PHPTest{
				FilePath:    "/tests/array/functions.phpt",
				Description: "Array functions",
			},
			expected: CategoryArray,
		},
		{
			name: "string test by path",
			test: &PHPTest{
				FilePath:    "/tests/string/basic.phpt",
				Description: "String test",
			},
			expected: CategoryString,
		},
		{
			name: "array test by description",
			test: &PHPTest{
				FilePath:    "/tests/other/test.phpt",
				Description: "Array manipulation test",
			},
			expected: CategoryArray,
		},
		{
			name: "OOP test by description",
			test: &PHPTest{
				FilePath:    "/tests/other/test.phpt",
				Description: "Class inheritance test",
			},
			expected: CategoryOOP,
		},
		{
			name: "other test",
			test: &PHPTest{
				FilePath:    "/tests/other/misc.phpt",
				Description: "Miscellaneous test",
			},
			expected: CategoryOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategorizeTest(tt.test)
			if result != tt.expected {
				t.Errorf("CategorizeTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTestCategory_String tests category string representation
func TestTestCategory_String(t *testing.T) {
	tests := []struct {
		category TestCategory
		expected string
	}{
		{CategoryLanguage, "language"},
		{CategoryStdLib, "stdlib"},
		{CategoryOOP, "oop"},
		{CategoryArray, "array"},
		{CategoryOther, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.category.String()
			if result != tt.expected {
				t.Errorf("category.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestNewTestFilter tests creating a test filter
func TestNewTestFilter(t *testing.T) {
	filter := NewTestFilter()

	if filter == nil {
		t.Fatal("NewTestFilter() returned nil")
	}

	if !filter.IncludeSkipped {
		t.Error("IncludeSkipped should default to true")
	}
}

// TestTestFilter_Matches tests filter matching
func TestTestFilter_Matches(t *testing.T) {
	tests := []struct {
		name     string
		filter   *TestFilter
		test     *PHPTest
		expected bool
	}{
		{
			name: "category match",
			filter: &TestFilter{
				Categories: []TestCategory{CategoryLanguage},
			},
			test: &PHPTest{
				FilePath:    "/tests/lang/basic.phpt",
				Description: "Language test",
			},
			expected: true,
		},
		{
			name: "category no match",
			filter: &TestFilter{
				Categories: []TestCategory{CategoryOOP},
			},
			test: &PHPTest{
				FilePath:    "/tests/lang/basic.phpt",
				Description: "Language test",
			},
			expected: false,
		},
		{
			name: "required extension match",
			filter: &TestFilter{
				RequiredExtensions: []string{"json"},
			},
			test: &PHPTest{
				Extensions: []string{"json", "pcre"},
			},
			expected: true,
		},
		{
			name: "required extension no match",
			filter: &TestFilter{
				RequiredExtensions: []string{"mbstring"},
			},
			test: &PHPTest{
				Extensions: []string{"json", "pcre"},
			},
			expected: false,
		},
		{
			name: "excluded extension",
			filter: &TestFilter{
				ExcludedExtensions: []string{"json"},
			},
			test: &PHPTest{
				Extensions: []string{"json", "pcre"},
			},
			expected: false,
		},
		{
			name: "no filters matches all",
			filter: &TestFilter{},
			test: &PHPTest{
				FilePath:    "/tests/test.phpt",
				Description: "Test",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.Matches(tt.test)
			if result != tt.expected {
				t.Errorf("filter.Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTestFilter_FilterTests tests filtering a list of tests
func TestTestFilter_FilterTests(t *testing.T) {
	tests := []*PHPTest{
		{
			FilePath:    "/tests/lang/basic.phpt",
			Description: "Language test",
		},
		{
			FilePath:    "/tests/oop/class.phpt",
			Description: "OOP test",
		},
		{
			FilePath:    "/tests/array/test.phpt",
			Description: "Array test",
		},
	}

	filter := &TestFilter{
		Categories: []TestCategory{CategoryLanguage, CategoryOOP},
	}

	filtered := filter.FilterTests(tests)

	if len(filtered) != 2 {
		t.Errorf("FilterTests() returned %d tests, want 2", len(filtered))
	}

	// Check that array test was filtered out
	for _, test := range filtered {
		category := CategorizeTest(test)
		if category == CategoryArray {
			t.Error("FilterTests() should have filtered out array test")
		}
	}
}

// TestGroupTestsByCategory tests grouping tests by category
func TestGroupTestsByCategory(t *testing.T) {
	tests := []*PHPTest{
		{
			FilePath:    "/tests/lang/test1.phpt",
			Description: "Language test 1",
		},
		{
			FilePath:    "/tests/lang/test2.phpt",
			Description: "Language test 2",
		},
		{
			FilePath:    "/tests/oop/test1.phpt",
			Description: "OOP test 1",
		},
	}

	grouped := GroupTestsByCategory(tests)

	if len(grouped[CategoryLanguage]) != 2 {
		t.Errorf("GroupTestsByCategory() language count = %d, want 2", len(grouped[CategoryLanguage]))
	}

	if len(grouped[CategoryOOP]) != 1 {
		t.Errorf("GroupTestsByCategory() OOP count = %d, want 1", len(grouped[CategoryOOP]))
	}
}

// TestGroupResultsByCategory tests grouping results by category
func TestGroupResultsByCategory(t *testing.T) {
	results := []*TestResult{
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test1.phpt",
				Description: "Language test 1",
			},
			Status: StatusPass,
		},
		{
			Test: &PHPTest{
				FilePath:    "/tests/lang/test2.phpt",
				Description: "Language test 2",
			},
			Status: StatusFail,
		},
		{
			Test: &PHPTest{
				FilePath:    "/tests/oop/test1.phpt",
				Description: "OOP test 1",
			},
			Status: StatusPass,
		},
	}

	grouped := GroupResultsByCategory(results)

	if len(grouped[CategoryLanguage]) != 2 {
		t.Errorf("GroupResultsByCategory() language count = %d, want 2", len(grouped[CategoryLanguage]))
	}

	if len(grouped[CategoryOOP]) != 1 {
		t.Errorf("GroupResultsByCategory() OOP count = %d, want 1", len(grouped[CategoryOOP]))
	}
}

// TestCountByCategory tests counting tests by category
func TestCountByCategory(t *testing.T) {
	tests := []*PHPTest{
		{
			FilePath:    "/tests/lang/test1.phpt",
			Description: "Language test 1",
		},
		{
			FilePath:    "/tests/lang/test2.phpt",
			Description: "Language test 2",
		},
		{
			FilePath:    "/tests/oop/test1.phpt",
			Description: "OOP test 1",
		},
		{
			FilePath:    "/tests/array/test1.phpt",
			Description: "Array test 1",
		},
	}

	counts := CountByCategory(tests)

	if counts[CategoryLanguage] != 2 {
		t.Errorf("CountByCategory() language = %d, want 2", counts[CategoryLanguage])
	}

	if counts[CategoryOOP] != 1 {
		t.Errorf("CountByCategory() OOP = %d, want 1", counts[CategoryOOP])
	}

	if counts[CategoryArray] != 1 {
		t.Errorf("CountByCategory() array = %d, want 1", counts[CategoryArray])
	}
}
