package pcre

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// preg_match Tests
// ============================================================================

func TestPregMatchBasic(t *testing.T) {
	pattern := types.NewString("/world/")
	subject := types.NewString("hello world")

	result := PregMatch(pattern, subject)
	if result.ToInt() != 1 {
		t.Error("Expected 1 match")
	}
}

func TestPregMatchNoMatch(t *testing.T) {
	pattern := types.NewString("/xyz/")
	subject := types.NewString("hello world")

	result := PregMatch(pattern, subject)
	if result.ToInt() != 0 {
		t.Error("Expected 0 matches")
	}
}

func TestPregMatchWithCaptures(t *testing.T) {
	pattern := types.NewString("/h(\\w+)o/")
	subject := types.NewString("hello")

	matches := types.NewArray(types.NewEmptyArray())
	result := PregMatch(pattern, subject, matches)

	if result.ToInt() != 1 {
		t.Error("Expected 1 match")
	}

	matchesArray := matches.ToArray()

	// Full match
	fullMatch, _ := matchesArray.Get(types.NewInt(0))
	if fullMatch.ToString() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", fullMatch.ToString())
	}

	// First capture group
	capture1, _ := matchesArray.Get(types.NewInt(1))
	if capture1.ToString() != "ell" {
		t.Errorf("Expected 'ell', got '%s'", capture1.ToString())
	}
}

func TestPregMatchCaseInsensitive(t *testing.T) {
	pattern := types.NewString("/HELLO/i")
	subject := types.NewString("hello world")

	result := PregMatch(pattern, subject)
	if result.ToInt() != 1 {
		t.Error("Expected case-insensitive match")
	}
}

func TestPregMatchInvalidPattern(t *testing.T) {
	pattern := types.NewString("/[/") // Invalid pattern
	subject := types.NewString("test")

	result := PregMatch(pattern, subject)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for invalid pattern")
	}
}

// ============================================================================
// preg_match_all Tests
// ============================================================================

func TestPregMatchAllBasic(t *testing.T) {
	pattern := types.NewString("/\\d+/")
	subject := types.NewString("There are 123 apples and 456 oranges")

	matches := types.NewArray(types.NewEmptyArray())
	result := PregMatchAll(pattern, subject, matches)

	if result.ToInt() != 2 {
		t.Errorf("Expected 2 matches, got %d", result.ToInt())
	}

	matchesArray := matches.ToArray()
	fullMatches, _ := matchesArray.Get(types.NewInt(0))
	fullMatchesArray := fullMatches.ToArray()

	if fullMatchesArray.Len() != 2 {
		t.Errorf("Expected 2 full matches, got %d", fullMatchesArray.Len())
	}

	first, _ := fullMatchesArray.Get(types.NewInt(0))
	if first.ToString() != "123" {
		t.Errorf("Expected '123', got '%s'", first.ToString())
	}

	second, _ := fullMatchesArray.Get(types.NewInt(1))
	if second.ToString() != "456" {
		t.Errorf("Expected '456', got '%s'", second.ToString())
	}
}

func TestPregMatchAllNoMatches(t *testing.T) {
	pattern := types.NewString("/xyz/")
	subject := types.NewString("hello world")

	matches := types.NewArray(types.NewEmptyArray())
	result := PregMatchAll(pattern, subject, matches)

	if result.ToInt() != 0 {
		t.Error("Expected 0 matches")
	}
}

func TestPregMatchAllWithCaptures(t *testing.T) {
	pattern := types.NewString("/(\\w+)@(\\w+)/")
	subject := types.NewString("user@example.com and admin@test.org")

	matches := types.NewArray(types.NewEmptyArray())
	result := PregMatchAll(pattern, subject, matches)

	if result.ToInt() != 2 {
		t.Errorf("Expected 2 matches, got %d", result.ToInt())
	}

	matchesArray := matches.ToArray()

	// First capture group (usernames)
	group1, _ := matchesArray.Get(types.NewInt(1))
	group1Array := group1.ToArray()

	first, _ := group1Array.Get(types.NewInt(0))
	if first.ToString() != "user" {
		t.Errorf("Expected 'user', got '%s'", first.ToString())
	}

	second, _ := group1Array.Get(types.NewInt(1))
	if second.ToString() != "admin" {
		t.Errorf("Expected 'admin', got '%s'", second.ToString())
	}
}

// ============================================================================
// preg_replace Tests
// ============================================================================

func TestPregReplaceBasic(t *testing.T) {
	pattern := types.NewString("/world/")
	replacement := types.NewString("PHP")
	subject := types.NewString("hello world")

	result := PregReplace(pattern, replacement, subject)
	if result.ToString() != "hello PHP" {
		t.Errorf("Expected 'hello PHP', got '%s'", result.ToString())
	}
}

func TestPregReplaceMultiple(t *testing.T) {
	pattern := types.NewString("/\\d+/")
	replacement := types.NewString("X")
	subject := types.NewString("I have 10 apples and 20 oranges")

	result := PregReplace(pattern, replacement, subject)
	if result.ToString() != "I have X apples and X oranges" {
		t.Errorf("Expected 'I have X apples and X oranges', got '%s'", result.ToString())
	}
}

func TestPregReplaceWithLimit(t *testing.T) {
	pattern := types.NewString("/\\d+/")
	replacement := types.NewString("X")
	subject := types.NewString("1 2 3 4 5")
	limit := types.NewInt(2)

	result := PregReplace(pattern, replacement, subject, limit)
	// Should replace first 2 occurrences
	if result.ToString() != "X X 3 4 5" {
		t.Errorf("Expected 'X X 3 4 5', got '%s'", result.ToString())
	}
}

func TestPregReplaceCaseInsensitive(t *testing.T) {
	pattern := types.NewString("/HELLO/i")
	replacement := types.NewString("hi")
	subject := types.NewString("Hello world, HELLO again")

	result := PregReplace(pattern, replacement, subject)
	if result.ToString() != "hi world, hi again" {
		t.Errorf("Expected 'hi world, hi again', got '%s'", result.ToString())
	}
}

func TestPregReplaceInvalidPattern(t *testing.T) {
	pattern := types.NewString("/[/")
	replacement := types.NewString("X")
	subject := types.NewString("test")

	result := PregReplace(pattern, replacement, subject)
	if result.Type() != types.TypeNull {
		t.Error("Expected null for invalid pattern")
	}
}

// ============================================================================
// preg_split Tests
// ============================================================================

func TestPregSplitBasic(t *testing.T) {
	pattern := types.NewString("/\\s+/")
	subject := types.NewString("hello   world  test")

	result := PregSplit(pattern, subject)
	if result.Type() != types.TypeArray {
		t.Fatal("Expected array result")
	}

	parts := result.ToArray()
	if parts.Len() != 3 {
		t.Errorf("Expected 3 parts, got %d", parts.Len())
	}

	first, _ := parts.Get(types.NewInt(0))
	if first.ToString() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", first.ToString())
	}
}

func TestPregSplitWithLimit(t *testing.T) {
	pattern := types.NewString("/,/")
	subject := types.NewString("a,b,c,d,e")
	limit := types.NewInt(3)

	result := PregSplit(pattern, subject, limit)
	parts := result.ToArray()

	if parts.Len() != 3 {
		t.Errorf("Expected 3 parts with limit, got %d", parts.Len())
	}

	// Last part should contain remaining string
	last, _ := parts.Get(types.NewInt(2))
	if last.ToString() != "c,d,e" {
		t.Errorf("Expected 'c,d,e', got '%s'", last.ToString())
	}
}

func TestPregSplitInvalidPattern(t *testing.T) {
	pattern := types.NewString("/[/")
	subject := types.NewString("test")

	result := PregSplit(pattern, subject)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for invalid pattern")
	}
}

// ============================================================================
// preg_grep Tests
// ============================================================================

func TestPregGrepBasic(t *testing.T) {
	pattern := types.NewString("/^a/")

	arr := types.NewEmptyArray()
	arr.Append(types.NewString("apple"))
	arr.Append(types.NewString("banana"))
	arr.Append(types.NewString("avocado"))
	arr.Append(types.NewString("cherry"))

	result := PregGrep(pattern, types.NewArray(arr))
	if result.Type() != types.TypeArray {
		t.Fatal("Expected array result")
	}

	filtered := result.ToArray()
	if filtered.Len() != 2 {
		t.Errorf("Expected 2 matches, got %d", filtered.Len())
	}

	// Check that correct items were kept
	first, _ := filtered.Get(types.NewInt(0))
	if first.ToString() != "apple" {
		t.Errorf("Expected 'apple', got '%s'", first.ToString())
	}
}

func TestPregGrepInvert(t *testing.T) {
	pattern := types.NewString("/^a/")
	flags := types.NewInt(1) // PREG_GREP_INVERT

	arr := types.NewEmptyArray()
	arr.Append(types.NewString("apple"))
	arr.Append(types.NewString("banana"))
	arr.Append(types.NewString("avocado"))

	result := PregGrep(pattern, types.NewArray(arr), flags)
	filtered := result.ToArray()

	// Should return items that DON'T start with 'a'
	if filtered.Len() != 1 {
		t.Errorf("Expected 1 match with invert, got %d", filtered.Len())
	}

	first, _ := filtered.Get(types.NewInt(1))
	if first.ToString() != "banana" {
		t.Errorf("Expected 'banana', got '%s'", first.ToString())
	}
}

// ============================================================================
// preg_quote Tests
// ============================================================================

func TestPregQuoteBasic(t *testing.T) {
	str := types.NewString("$40 for a [special] item?")
	result := PregQuote(str)

	expected := "\\$40 for a \\[special\\] item\\?"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}
}

func TestPregQuoteWithDelimiter(t *testing.T) {
	str := types.NewString("/path/to/file")
	delimiter := types.NewString("/")

	result := PregQuote(str, delimiter)

	expected := "\\/path\\/to\\/file"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}
}

func TestPregQuoteAllSpecialChars(t *testing.T) {
	str := types.NewString(".*+?^$[]{}()|\\")
	result := PregQuote(str)

	// All special chars should be escaped
	resultStr := result.ToString()
	if !containsAllEscaped(resultStr) {
		t.Error("Not all special characters were escaped")
	}
}

func containsAllEscaped(s string) bool {
	// Check that special chars are escaped
	specialChars := []string{"\\.", "\\*", "\\+", "\\?", "\\^", "\\$", "\\[", "\\]", "\\{", "\\}", "\\(", "\\)", "\\|", "\\\\"}
	for _, char := range specialChars {
		if !contains(s, char) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Pattern Compilation Tests
// ============================================================================

func TestPatternCaching(t *testing.T) {
	// Clear cache first
	ClearPatternCache()

	pattern := types.NewString("/test/")
	subject := types.NewString("test string")

	// First call - should compile and cache
	PregMatch(pattern, subject)

	// Second call - should use cached pattern
	result := PregMatch(pattern, subject)

	if result.ToInt() != 1 {
		t.Error("Cached pattern should work correctly")
	}
}

func TestMultilineFlag(t *testing.T) {
	pattern := types.NewString("/^world/m")
	subject := types.NewString("hello\nworld")

	result := PregMatch(pattern, subject)
	if result.ToInt() != 1 {
		t.Error("Multiline flag should enable ^ to match after newline")
	}
}

func TestDotAllFlag(t *testing.T) {
	pattern := types.NewString("/a.b/s")
	subject := types.NewString("a\nb")

	result := PregMatch(pattern, subject)
	if result.ToInt() != 1 {
		t.Error("Dotall flag should enable . to match newline")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestPregMatchNilInputs(t *testing.T) {
	result := PregMatch(nil, nil)
	if result.ToInt() != 0 {
		t.Error("Expected 0 for nil inputs")
	}
}

func TestPregReplaceNilInputs(t *testing.T) {
	result := PregReplace(nil, nil, nil)
	if result.Type() != types.TypeNull {
		t.Error("Expected null for nil inputs")
	}
}

func TestPregSplitNilInputs(t *testing.T) {
	result := PregSplit(nil, nil)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for nil inputs")
	}
}

func TestPregGrepNonArray(t *testing.T) {
	pattern := types.NewString("/test/")
	notArray := types.NewString("not an array")

	result := PregGrep(pattern, notArray)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for non-array input")
	}
}

func TestPregQuoteNil(t *testing.T) {
	result := PregQuote(nil)
	if result.ToString() != "" {
		t.Error("Expected empty string for nil input")
	}
}

func TestPregLastError(t *testing.T) {
	result := PregLastError()
	if result.ToInt() != 0 {
		t.Error("Expected error code 0 (PREG_NO_ERROR)")
	}
}
