package string

import (
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Strlen/Substr Tests
// ============================================================================

func TestStrlen(t *testing.T) {
	str := types.NewString("Hello World")
	result := Strlen(str)

	if result.ToInt() != 11 {
		t.Errorf("Expected length 11, got %d", result.ToInt())
	}
}

func TestStrlenEmpty(t *testing.T) {
	str := types.NewString("")
	result := Strlen(str)

	if result.ToInt() != 0 {
		t.Errorf("Expected length 0, got %d", result.ToInt())
	}
}

func TestSubstr(t *testing.T) {
	str := types.NewString("Hello World")

	// Positive offset
	result := Substr(str, types.NewInt(0), types.NewInt(5))
	if result.ToString() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result.ToString())
	}

	// Negative offset
	result = Substr(str, types.NewInt(-5), types.NewInt(5))
	if result.ToString() != "World" {
		t.Errorf("Expected 'World', got '%s'", result.ToString())
	}

	// No length specified
	result = Substr(str, types.NewInt(6))
	if result.ToString() != "World" {
		t.Errorf("Expected 'World', got '%s'", result.ToString())
	}
}

// ============================================================================
// Search Tests
// ============================================================================

func TestStrpos(t *testing.T) {
	haystack := types.NewString("Hello World")
	needle := types.NewString("World")

	result := Strpos(haystack, needle)
	if result.ToInt() != 6 {
		t.Errorf("Expected position 6, got %d", result.ToInt())
	}

	// Not found
	result = Strpos(haystack, types.NewString("xyz"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for not found")
	}
}

func TestStrrpos(t *testing.T) {
	haystack := types.NewString("Hello World World")
	needle := types.NewString("World")

	result := Strrpos(haystack, needle)
	if result.ToInt() != 12 {
		t.Errorf("Expected position 12, got %d", result.ToInt())
	}
}

func TestStripos(t *testing.T) {
	haystack := types.NewString("Hello World")
	needle := types.NewString("world")

	result := Stripos(haystack, needle)
	if result.ToInt() != 6 {
		t.Errorf("Expected position 6 (case-insensitive), got %d", result.ToInt())
	}
}

func TestStrripos(t *testing.T) {
	haystack := types.NewString("Hello World WORLD")
	needle := types.NewString("world")

	result := Strripos(haystack, needle)
	if result.ToInt() != 12 {
		t.Errorf("Expected position 12 (case-insensitive), got %d", result.ToInt())
	}
}

// ============================================================================
// Replace Tests
// ============================================================================

func TestStrReplace(t *testing.T) {
	subject := types.NewString("Hello World")
	search := types.NewString("World")
	replace := types.NewString("PHP")

	result := StrReplace(search, replace, subject)
	if result.ToString() != "Hello PHP" {
		t.Errorf("Expected 'Hello PHP', got '%s'", result.ToString())
	}
}

func TestStrIreplace(t *testing.T) {
	subject := types.NewString("Hello World")
	search := types.NewString("world")
	replace := types.NewString("PHP")

	result := StrIreplace(search, replace, subject)
	if result.ToString() != "Hello PHP" {
		t.Errorf("Expected 'Hello PHP' (case-insensitive), got '%s'", result.ToString())
	}
}

// ============================================================================
// Case Conversion Tests
// ============================================================================

func TestStrtolower(t *testing.T) {
	str := types.NewString("Hello World")
	result := Strtolower(str)

	if result.ToString() != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result.ToString())
	}
}

func TestStrtoupper(t *testing.T) {
	str := types.NewString("Hello World")
	result := Strtoupper(str)

	if result.ToString() != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD', got '%s'", result.ToString())
	}
}

func TestUcfirst(t *testing.T) {
	str := types.NewString("hello world")
	result := Ucfirst(str)

	if result.ToString() != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", result.ToString())
	}
}

func TestLcfirst(t *testing.T) {
	str := types.NewString("Hello World")
	result := Lcfirst(str)

	if result.ToString() != "hello World" {
		t.Errorf("Expected 'hello World', got '%s'", result.ToString())
	}
}

func TestUcwords(t *testing.T) {
	str := types.NewString("hello world")
	result := Ucwords(str)

	if result.ToString() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.ToString())
	}
}

// ============================================================================
// Trim Tests
// ============================================================================

func TestTrim(t *testing.T) {
	str := types.NewString("  Hello World  ")
	result := Trim(str)

	if result.ToString() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.ToString())
	}
}

func TestLtrim(t *testing.T) {
	str := types.NewString("  Hello World  ")
	result := Ltrim(str)

	if result.ToString() != "Hello World  " {
		t.Errorf("Expected 'Hello World  ', got '%s'", result.ToString())
	}
}

func TestRtrim(t *testing.T) {
	str := types.NewString("  Hello World  ")
	result := Rtrim(str)

	if result.ToString() != "  Hello World" {
		t.Errorf("Expected '  Hello World', got '%s'", result.ToString())
	}
}

func TestTrimWithCharacters(t *testing.T) {
	str := types.NewString("##Hello##")
	result := Trim(str, types.NewString("#"))

	if result.ToString() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result.ToString())
	}
}

// ============================================================================
// Explode/Implode Tests
// ============================================================================

func TestExplode(t *testing.T) {
	str := types.NewString("one,two,three")
	delim := types.NewString(",")

	result := Explode(delim, str)
	arr := result.ToArray()

	if arr.Len() != 3 {
		t.Errorf("Expected 3 parts, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToString() != "one" {
		t.Errorf("Expected 'one', got '%s'", val.ToString())
	}
}

func TestExplodeWithLimit(t *testing.T) {
	str := types.NewString("one,two,three,four")
	delim := types.NewString(",")

	result := Explode(delim, str, types.NewInt(2))
	arr := result.ToArray()

	if arr.Len() != 2 {
		t.Errorf("Expected 2 parts with limit, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(1))
	if val.ToString() != "two,three,four" {
		t.Errorf("Expected 'two,three,four', got '%s'", val.ToString())
	}
}

func TestExplodeEmptyDelimiter(t *testing.T) {
	str := types.NewString("hello")
	delim := types.NewString("")

	result := Explode(delim, str)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for empty delimiter")
	}
}

func TestImplode(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewString("one"), types.NewString("two"), types.NewString("three"))
	arrVal := types.NewArray(arr)

	sep := types.NewString(",")
	result := Implode(sep, arrVal)

	if result.ToString() != "one,two,three" {
		t.Errorf("Expected 'one,two,three', got '%s'", result.ToString())
	}
}

func TestImplodeNoSeparator(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewString("one"), types.NewString("two"))
	arrVal := types.NewArray(arr)

	result := Implode(arrVal)

	if result.ToString() != "onetwo" {
		t.Errorf("Expected 'onetwo', got '%s'", result.ToString())
	}
}

// ============================================================================
// Additional Function Tests
// ============================================================================

func TestStrSplit(t *testing.T) {
	str := types.NewString("Hello")
	result := StrSplit(str)
	arr := result.ToArray()

	if arr.Len() != 5 {
		t.Errorf("Expected 5 characters, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToString() != "H" {
		t.Errorf("Expected 'H', got '%s'", val.ToString())
	}
}

func TestStrSplitWithLength(t *testing.T) {
	str := types.NewString("Hello")
	result := StrSplit(str, types.NewInt(2))
	arr := result.ToArray()

	if arr.Len() != 3 {
		t.Errorf("Expected 3 chunks, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToString() != "He" {
		t.Errorf("Expected 'He', got '%s'", val.ToString())
	}
}

func TestStrRepeat(t *testing.T) {
	str := types.NewString("ab")
	result := StrRepeat(str, types.NewInt(3))

	if result.ToString() != "ababab" {
		t.Errorf("Expected 'ababab', got '%s'", result.ToString())
	}
}

func TestStrRepeatNegative(t *testing.T) {
	str := types.NewString("ab")
	result := StrRepeat(str, types.NewInt(-1))

	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for negative repeat count")
	}
}

func TestStrPad(t *testing.T) {
	str := types.NewString("Hello")
	result := StrPad(str, types.NewInt(10), types.NewString("-"))

	if result.ToString() != "Hello-----" {
		t.Errorf("Expected 'Hello-----', got '%s'", result.ToString())
	}
}

func TestStrRev(t *testing.T) {
	str := types.NewString("Hello")
	result := StrRev(str)

	if result.ToString() != "olleH" {
		t.Errorf("Expected 'olleH', got '%s'", result.ToString())
	}
}

func TestStrstr(t *testing.T) {
	haystack := types.NewString("Hello World")
	needle := types.NewString("World")

	result := Strstr(haystack, needle)
	if result.ToString() != "World" {
		t.Errorf("Expected 'World', got '%s'", result.ToString())
	}

	// With before_needle
	result = Strstr(haystack, needle, types.NewBool(true))
	if result.ToString() != "Hello " {
		t.Errorf("Expected 'Hello ', got '%s'", result.ToString())
	}

	// Not found
	result = Strstr(haystack, types.NewString("xyz"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for not found")
	}
}

func TestChunkSplit(t *testing.T) {
	str := types.NewString("Hello")
	result := ChunkSplit(str, types.NewInt(2), types.NewString("-"))

	if result.ToString() != "He-ll-o-" {
		t.Errorf("Expected 'He-ll-o-', got '%s'", result.ToString())
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestStrlenNil(t *testing.T) {
	result := Strlen(nil)
	if result.ToInt() != 0 {
		t.Errorf("Expected length 0 for nil, got %d", result.ToInt())
	}
}

func TestSubstrBeyondLength(t *testing.T) {
	str := types.NewString("Hello")
	result := Substr(str, types.NewInt(10))

	if result.ToString() != "" {
		t.Errorf("Expected empty string, got '%s'", result.ToString())
	}
}

func TestSubstrNegativeLength(t *testing.T) {
	str := types.NewString("Hello World")
	result := Substr(str, types.NewInt(0), types.NewInt(-6))

	if result.ToString() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result.ToString())
	}
}

func TestStrposOffset(t *testing.T) {
	haystack := types.NewString("Hello World World")
	needle := types.NewString("World")

	// Start search from position 7
	result := Strpos(haystack, needle, types.NewInt(7))
	if result.ToInt() != 12 {
		t.Errorf("Expected position 12 with offset, got %d", result.ToInt())
	}
}

func TestJoinAlias(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewString("a"), types.NewString("b"))
	arrVal := types.NewArray(arr)

	result := Join(types.NewString("-"), arrVal)
	if result.ToString() != "a-b" {
		t.Errorf("Expected 'a-b', got '%s'", result.ToString())
	}
}

func TestStrchrAlias(t *testing.T) {
	haystack := types.NewString("Hello World")
	needle := types.NewString("World")

	result := Strchr(haystack, needle)
	if result.ToString() != "World" {
		t.Errorf("Expected 'World', got '%s'", result.ToString())
	}
}

// ============================================================================
// String Formatting Tests
// ============================================================================

func TestSprintf(t *testing.T) {
	// String formatting
	result := Sprintf(types.NewString("Hello %s"), types.NewString("World"))
	if result.ToString() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.ToString())
	}

	// Integer formatting
	result = Sprintf(types.NewString("Number: %d"), types.NewInt(42))
	if result.ToString() != "Number: 42" {
		t.Errorf("Expected 'Number: 42', got '%s'", result.ToString())
	}

	// Multiple values
	result = Sprintf(types.NewString("%s has %d apples"), types.NewString("John"), types.NewInt(5))
	if result.ToString() != "John has 5 apples" {
		t.Errorf("Expected 'John has 5 apples', got '%s'", result.ToString())
	}

	// Percent escape
	result = Sprintf(types.NewString("100%% complete"))
	if result.ToString() != "100% complete" {
		t.Errorf("Expected '100%% complete', got '%s'", result.ToString())
	}
}

func TestSprintfCharacter(t *testing.T) {
	// Character formatting
	result := Sprintf(types.NewString("Char: %c"), types.NewInt(65))
	if result.ToString() != "Char: A" {
		t.Errorf("Expected 'Char: A', got '%s'", result.ToString())
	}
}

func TestPrintf(t *testing.T) {
	result := Printf(types.NewString("Test %s"), types.NewString("message"))
	// Printf returns the length of the output
	if result.ToInt() != 12 { // "Test message" is 12 chars
		t.Errorf("Expected length 12, got %d", result.ToInt())
	}
}

// ============================================================================
// String Comparison Tests
// ============================================================================

func TestStrcmp(t *testing.T) {
	// Equal strings
	result := Strcmp(types.NewString("test"), types.NewString("test"))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for equal strings, got %d", result.ToInt())
	}

	// First < Second
	result = Strcmp(types.NewString("abc"), types.NewString("xyz"))
	if result.ToInt() != -1 {
		t.Errorf("Expected -1, got %d", result.ToInt())
	}

	// First > Second
	result = Strcmp(types.NewString("xyz"), types.NewString("abc"))
	if result.ToInt() != 1 {
		t.Errorf("Expected 1, got %d", result.ToInt())
	}
}

func TestStrcasecmp(t *testing.T) {
	// Case-insensitive equal
	result := Strcasecmp(types.NewString("Test"), types.NewString("TEST"))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for case-insensitive equal, got %d", result.ToInt())
	}

	// Case-insensitive less than
	result = Strcasecmp(types.NewString("ABC"), types.NewString("xyz"))
	if result.ToInt() != -1 {
		t.Errorf("Expected -1, got %d", result.ToInt())
	}
}

func TestStrncmp(t *testing.T) {
	// Compare first 3 characters
	result := Strncmp(types.NewString("testing"), types.NewString("tested"), types.NewInt(4))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for first 4 chars equal, got %d", result.ToInt())
	}

	// Different within n characters
	result = Strncmp(types.NewString("apple"), types.NewString("orange"), types.NewInt(3))
	if result.ToInt() == 0 {
		t.Error("Expected non-zero for different strings")
	}

	// Zero length comparison
	result = Strncmp(types.NewString("abc"), types.NewString("xyz"), types.NewInt(0))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for zero length comparison, got %d", result.ToInt())
	}
}

func TestStrncasecmp(t *testing.T) {
	result := Strncasecmp(types.NewString("Testing"), types.NewString("TESTED"), types.NewInt(4))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for case-insensitive first 4 chars, got %d", result.ToInt())
	}
}

func TestStristr(t *testing.T) {
	// Case-insensitive search
	result := Stristr(types.NewString("Hello World"), types.NewString("WORLD"))
	if result.ToString() != "World" {
		t.Errorf("Expected 'World', got '%s'", result.ToString())
	}

	// Not found
	result = Stristr(types.NewString("Hello"), types.NewString("xyz"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for not found")
	}

	// Before needle
	result = Stristr(types.NewString("Hello World"), types.NewString("WORLD"), types.NewBool(true))
	if result.ToString() != "Hello " {
		t.Errorf("Expected 'Hello ', got '%s'", result.ToString())
	}
}

func TestStrrchr(t *testing.T) {
	// Find last occurrence
	result := Strrchr(types.NewString("hello world"), types.NewString("o"))
	if result.ToString() != "orld" {
		t.Errorf("Expected 'orld', got '%s'", result.ToString())
	}

	// Not found
	result = Strrchr(types.NewString("hello"), types.NewString("x"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for not found")
	}
}

// ============================================================================
// HTML Functions Tests
// ============================================================================

func TestHtmlspecialchars(t *testing.T) {
	input := types.NewString("<div class=\"test\">Hello & goodbye</div>")
	result := Htmlspecialchars(input)

	expected := "&lt;div class=&quot;test&quot;&gt;Hello &amp; goodbye&lt;/div&gt;"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}

	// Test single quotes
	input = types.NewString("It's a test")
	result = Htmlspecialchars(input)
	if !strings.Contains(result.ToString(), "&#039;") {
		t.Error("Expected single quote to be encoded")
	}
}

func TestHtmlentities(t *testing.T) {
	input := types.NewString("<b>Bold</b>")
	result := Htmlentities(input)

	if result.ToString() != "&lt;b&gt;Bold&lt;/b&gt;" {
		t.Errorf("Expected '&lt;b&gt;Bold&lt;/b&gt;', got '%s'", result.ToString())
	}
}

func TestHtmlspecialcharsDecode(t *testing.T) {
	input := types.NewString("&lt;div&gt;Test&amp;decode&lt;/div&gt;")
	result := HtmlspecialcharsDecode(input)

	expected := "<div>Test&decode</div>"
	if result.ToString() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.ToString())
	}

	// Test quotes
	input = types.NewString("&quot;quoted&quot; and &#039;single&#039;")
	result = HtmlspecialcharsDecode(input)
	if result.ToString() != "\"quoted\" and 'single'" {
		t.Errorf("Expected '\"quoted\" and 'single'', got '%s'", result.ToString())
	}
}

// ============================================================================
// Slashing Tests
// ============================================================================

func TestAddslashes(t *testing.T) {
	input := types.NewString("It's a \"test\"")
	result := Addslashes(input)

	if result.ToString() != "It\\'s a \\\"test\\\"" {
		t.Errorf("Expected 'It\\\\'s a \\\\\"test\\\\\"', got '%s'", result.ToString())
	}

	// Test backslashes
	input = types.NewString("path\\to\\file")
	result = Addslashes(input)
	if result.ToString() != "path\\\\to\\\\file" {
		t.Errorf("Expected 'path\\\\\\\\to\\\\\\\\file', got '%s'", result.ToString())
	}
}

func TestStripslashes(t *testing.T) {
	input := types.NewString("It\\'s a \\\"test\\\"")
	result := Stripslashes(input)

	if result.ToString() != "It's a \"test\"" {
		t.Errorf("Expected 'It's a \"test\"', got '%s'", result.ToString())
	}

	// Test backslashes
	input = types.NewString("path\\\\to\\\\file")
	result = Stripslashes(input)
	if result.ToString() != "path\\to\\file" {
		t.Errorf("Expected 'path\\to\\file', got '%s'", result.ToString())
	}
}

// ============================================================================
// Text Formatting Tests
// ============================================================================

func TestNl2br(t *testing.T) {
	// XHTML style (default)
	input := types.NewString("Line 1\nLine 2\nLine 3")
	result := Nl2br(input)

	if !strings.Contains(result.ToString(), "<br />") {
		t.Error("Expected XHTML style breaks")
	}

	// HTML style
	result = Nl2br(input, types.NewBool(false))
	if !strings.Contains(result.ToString(), "<br>") {
		t.Error("Expected HTML style breaks")
	}

	// Test \r\n
	input = types.NewString("Windows\r\nLine")
	result = Nl2br(input)
	if !strings.Contains(result.ToString(), "<br />") {
		t.Error("Expected breaks for \\r\\n")
	}
}

func TestWordwrap(t *testing.T) {
	input := types.NewString("The quick brown fox jumps over the lazy dog")
	result := Wordwrap(input, types.NewInt(15))

	// Should have line breaks
	if !strings.Contains(result.ToString(), "\n") {
		t.Error("Expected line breaks in wrapped text")
	}

	// Custom break string
	result = Wordwrap(input, types.NewInt(15), types.NewString("<br>"))
	if !strings.Contains(result.ToString(), "<br>") {
		t.Error("Expected custom break string")
	}
}

// ============================================================================
// URL Encoding Tests
// ============================================================================

func TestUrlencode(t *testing.T) {
	// Spaces become +
	input := types.NewString("Hello World")
	result := Urlencode(input)
	if result.ToString() != "Hello+World" {
		t.Errorf("Expected 'Hello+World', got '%s'", result.ToString())
	}

	// Special characters
	input = types.NewString("test@example.com")
	result = Urlencode(input)
	if !strings.Contains(result.ToString(), "%40") {
		t.Error("Expected @ to be encoded as %40")
	}

	// Safe characters
	input = types.NewString("abc123-_.~")
	result = Urlencode(input)
	if result.ToString() != "abc123-_.~" {
		t.Errorf("Expected safe chars unchanged, got '%s'", result.ToString())
	}
}

func TestUrldecode(t *testing.T) {
	// Decode +
	input := types.NewString("Hello+World")
	result := Urldecode(input)
	if result.ToString() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.ToString())
	}

	// Decode hex
	input = types.NewString("test%40example.com")
	result = Urldecode(input)
	if result.ToString() != "test@example.com" {
		t.Errorf("Expected 'test@example.com', got '%s'", result.ToString())
	}

	// Invalid hex should remain
	input = types.NewString("test%ZZ")
	result = Urldecode(input)
	if result.ToString() != "test%ZZ" {
		t.Errorf("Expected 'test%%ZZ' unchanged, got '%s'", result.ToString())
	}
}

func TestRawurlencode(t *testing.T) {
	// Spaces become %20 (not +)
	input := types.NewString("Hello World")
	result := Rawurlencode(input)
	if result.ToString() != "Hello%20World" {
		t.Errorf("Expected 'Hello%%20World', got '%s'", result.ToString())
	}

	// Special characters
	input = types.NewString("test@example.com")
	result = Rawurlencode(input)
	if !strings.Contains(result.ToString(), "%40") {
		t.Error("Expected @ to be encoded")
	}

	// Safe characters per RFC 3986
	input = types.NewString("abc123-_.~")
	result = Rawurlencode(input)
	if result.ToString() != "abc123-_.~" {
		t.Errorf("Expected safe chars unchanged, got '%s'", result.ToString())
	}
}

func TestRawurldecode(t *testing.T) {
	// Decode %20
	input := types.NewString("Hello%20World")
	result := Rawurldecode(input)
	if result.ToString() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.ToString())
	}

	// Decode hex
	input = types.NewString("test%40example.com")
	result = Rawurldecode(input)
	if result.ToString() != "test@example.com" {
		t.Errorf("Expected 'test@example.com', got '%s'", result.ToString())
	}

	// Plus should remain plus
	input = types.NewString("one+two")
	result = Rawurldecode(input)
	if result.ToString() != "one+two" {
		t.Errorf("Expected 'one+two' unchanged, got '%s'", result.ToString())
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestSprintfNil(t *testing.T) {
	result := Sprintf(nil)
	if result.ToString() != "" {
		t.Error("Expected empty string for nil format")
	}
}

func TestHtmlspecialcharsNil(t *testing.T) {
	result := Htmlspecialchars(nil)
	if result.ToString() != "" {
		t.Error("Expected empty string for nil input")
	}
}

func TestUrlencodeNil(t *testing.T) {
	result := Urlencode(nil)
	if result.ToString() != "" {
		t.Error("Expected empty string for nil input")
	}
}

// ============================================================================
// String Length Limits Tests
// ============================================================================

func TestStringLengthConstants(t *testing.T) {
	// Verify constants are properly defined
	if MaxStringLength != 1024*1024*100 {
		t.Errorf("MaxStringLength = %d, expected %d", MaxStringLength, 1024*1024*100)
	}
	if MaxRepeatCount != 1000000 {
		t.Errorf("MaxRepeatCount = %d, expected %d", MaxRepeatCount, 1000000)
	}
}

func TestStrRepeatLengthLimit(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		times     int64
		shouldFail bool
		reason    string
	}{
		{
			name:      "Normal repeat",
			str:       "A",
			times:     10,
			shouldFail: false,
			reason:    "Small repeat should succeed",
		},
		{
			name:      "Large repeat count",
			str:       "A",
			times:     1000000,
			shouldFail: false,
			reason:    "Max repeat count should still work",
		},
		{
			name:      "Exceed max repeat count",
			str:       "A",
			times:     1000001,
			shouldFail: true,
			reason:    "Should fail when exceeding MaxRepeatCount",
		},
		{
			name:      "Exceed max string length",
			str:       "AAAAA", // 5 bytes
			times:     25000000, // 5 * 25M = 125MB > 100MB limit
			shouldFail: true,
			reason:    "Should fail when result would exceed MaxStringLength",
		},
		{
			name:      "Negative repeat count",
			str:       "A",
			times:     -1,
			shouldFail: true,
			reason:    "Negative count should fail",
		},
		{
			name:      "Zero repeat count",
			str:       "A",
			times:     0,
			shouldFail: false,
			reason:    "Zero count should return empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrRepeat(types.NewString(tt.str), types.NewInt(tt.times))

			if tt.shouldFail {
				if result.Type() != types.TypeBool || result.ToBool() != false {
					t.Errorf("%s: expected false, got %v", tt.reason, result)
				}
			} else {
				if result.Type() != types.TypeString {
					t.Errorf("%s: expected string type, got %v", tt.reason, result.Type())
				}
			}
		})
	}
}

func TestStrPadLengthLimit(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		targetLen  int64
		shouldFail bool
		reason     string
	}{
		{
			name:       "Normal pad",
			str:        "hello",
			targetLen:  10,
			shouldFail: false,
			reason:     "Small pad should succeed",
		},
		{
			name:       "Exceed max string length",
			str:        "hello",
			targetLen:  MaxStringLength + 1,
			shouldFail: true,
			reason:     "Should fail when target length exceeds MaxStringLength",
		},
		{
			name:       "Negative target length",
			str:        "hello",
			targetLen:  -1,
			shouldFail: true,
			reason:     "Negative length should fail",
		},
		{
			name:       "At max string length boundary",
			str:        "hello",
			targetLen:  MaxStringLength,
			shouldFail: false,
			reason:     "Exactly at limit should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrPad(types.NewString(tt.str), types.NewInt(tt.targetLen), types.NewString(" "))

			if tt.shouldFail {
				if result.Type() != types.TypeBool || result.ToBool() != false {
					t.Errorf("%s: expected false, got %v", tt.reason, result)
				}
			} else {
				if result.Type() != types.TypeString {
					t.Errorf("%s: expected string type, got %v", tt.reason, result.Type())
				}
			}
		})
	}
}

func TestStrReplaceLengthLimit(t *testing.T) {
	tests := []struct {
		name       string
		search     string
		replace    string
		subject    string
		shouldFail bool
		reason     string
	}{
		{
			name:       "Normal replace",
			search:     "a",
			replace:    "b",
			subject:    "aaa",
			shouldFail: false,
			reason:     "Simple replacement should succeed",
		},
		{
			name:       "Expand string within limits",
			search:     "x",
			replace:    "yyy",
			subject:    "xxx",
			shouldFail: false,
			reason:     "Small expansion should succeed",
		},
		{
			name:       "Empty search returns original",
			search:     "",
			replace:    "anything",
			subject:    "test",
			shouldFail: false,
			reason:     "Empty search should return original",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrReplace(
				types.NewString(tt.search),
				types.NewString(tt.replace),
				types.NewString(tt.subject),
			)

			if tt.shouldFail {
				if result.Type() != types.TypeBool || result.ToBool() != false {
					t.Errorf("%s: expected false, got %v", tt.reason, result)
				}
			} else {
				if result.Type() != types.TypeString {
					t.Errorf("%s: expected string type, got %v", tt.reason, result.Type())
				}
			}
		})
	}
}

func TestStrReplaceMemoryExhaustionPrevention(t *testing.T) {
	// Test that str_replace prevents creating excessively large strings
	// This test verifies the security check works

	// Create a string with many instances of search term
	// Each 'x' will be replaced with a long string
	subject := strings.Repeat("x", 1000000) // 1M occurrences

	// Try to replace each 'x' with 200 bytes (would result in 200MB string)
	replace := strings.Repeat("A", 200)

	result := StrReplace(
		types.NewString("x"),
		types.NewString(replace),
		types.NewString(subject),
	)

	// Should return false instead of allocating 200MB
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("str_replace should prevent excessive memory allocation (200MB result)")
	}
}

// ============================================================================
// Additional HTML Encoding Tests
// ============================================================================

func TestHtmlEntityDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "named entity",
			input:    "&copy; 2024",
			expected: "© 2024",
		},
		{
			name:     "pound entity",
			input:    "&pound;50",
			expected: "£50",
		},
		{
			name:     "decimal numeric entity",
			input:    "&#60;div&#62;",
			expected: "<div>",
		},
		{
			name:     "hex numeric entity",
			input:    "&#x3C;div&#x3E;",
			expected: "<div>",
		},
		{
			name:     "mixed entities",
			input:    "&lt;p&gt;&copy; 2024&lt;/p&gt;",
			expected: "<p>© 2024</p>",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no entities",
			input:    "Hello World",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HtmlEntityDecode(types.NewString(tt.input))
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		allowedTags *types.Value
		expected    string
	}{
		{
			name:        "remove all tags",
			input:       "<p>Hello <b>World</b></p>",
			allowedTags: nil,
			expected:    "Hello World",
		},
		{
			name:        "allow p tag",
			input:       "<p>Hello <b>World</b></p>",
			allowedTags: types.NewString("<p>"),
			expected:    "<p>Hello World</p>",
		},
		{
			name:        "allow multiple tags",
			input:       "<p>Hello <b>World</b></p>",
			allowedTags: types.NewString("<p><b>"),
			expected:    "<p>Hello <b>World</b></p>",
		},
		{
			name:        "remove PHP tags",
			input:       "<?php echo 'test'; ?> Text",
			allowedTags: nil,
			expected:    " Text",
		},
		{
			name:        "remove HTML comments",
			input:       "Hello <!-- comment --> World",
			allowedTags: nil,
			expected:    "Hello  World",
		},
		{
			name:        "empty string",
			input:       "",
			allowedTags: nil,
			expected:    "",
		},
		{
			name:        "no tags",
			input:       "Hello World",
			allowedTags: nil,
			expected:    "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *types.Value
			if tt.allowedTags == nil {
				result = StripTags(types.NewString(tt.input))
			} else {
				result = StripTags(types.NewString(tt.input), tt.allowedTags)
			}
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestAddslashesStripslashesRoundtrip(t *testing.T) {
	// Test that stripslashes reverses addslashes
	original := `It's a "test" with \backslash`
	escaped := Addslashes(types.NewString(original))
	result := Stripslashes(escaped)

	if result.ToString() != original {
		t.Errorf("Roundtrip failed: original=%q, escaped=%q, result=%q",
			original, escaped.ToString(), result.ToString())
	}
}

func TestHtmlspecialcharsWithFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		flags    int64
		expected string
	}{
		{
			name:     "ENT_NOQUOTES",
			input:    `"Hello" 'World'`,
			flags:    int64(ENT_NOQUOTES),
			expected: `"Hello" 'World'`,
		},
		{
			name:     "ENT_COMPAT (double quotes only)",
			input:    `"Hello" 'World'`,
			flags:    int64(ENT_COMPAT),
			expected: `&quot;Hello&quot; 'World'`,
		},
		{
			name:     "ENT_QUOTES (both quotes)",
			input:    `"Hello" 'World'`,
			flags:    int64(ENT_QUOTES),
			expected: `&quot;Hello&quot; &#039;World&#039;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Htmlspecialchars(types.NewString(tt.input), types.NewInt(tt.flags))
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestHtmlentitiesExtended(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "copyright symbol",
			input:    "© 2024",
			expected: "&copy; 2024",
		},
		{
			name:     "pound sign",
			input:    "£50",
			expected: "&pound;50",
		},
		{
			name:     "mixed special chars",
			input:    "<p>© 2024</p>",
			expected: "&lt;p&gt;&copy; 2024&lt;/p&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Htmlentities(types.NewString(tt.input))
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

// ============================================================================
// String Manipulation Tests (str_shuffle, str_word_count, substr_count, substr_replace)
// ============================================================================

func TestStrShuffle(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		verify func(result string) bool
	}{
		{
			name:  "empty string",
			input: "",
			verify: func(result string) bool {
				return result == ""
			},
		},
		{
			name:  "single char",
			input: "a",
			verify: func(result string) bool {
				return result == "a"
			},
		},
		{
			name:  "multi char",
			input: "abc",
			verify: func(result string) bool {
				// Should have same length and same characters
				if len(result) != 3 {
					return false
				}
				// Check all characters are present
				chars := make(map[rune]int)
				for _, c := range "abc" {
					chars[c]++
				}
				for _, c := range result {
					chars[c]--
				}
				for _, v := range chars {
					if v != 0 {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrShuffle(types.NewString(tt.input))
			if !tt.verify(result.ToString()) {
				t.Errorf("StrShuffle(%q) = %q, verification failed", tt.input, result.ToString())
			}
		})
	}
}

func TestStrShuffleNil(t *testing.T) {
	result := StrShuffle(nil)
	if result.ToString() != "" {
		t.Errorf("Expected empty string for nil, got %q", result.ToString())
	}
}

func TestStrWordCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		format   int64
		expected interface{} // int64 for format 0, []string for format 1, map[int]string for format 2
	}{
		{
			name:     "count words",
			input:    "Hello World",
			format:   0,
			expected: int64(2),
		},
		{
			name:     "count words multiple",
			input:    "The quick brown fox",
			format:   0,
			expected: int64(4),
		},
		{
			name:     "count empty string",
			input:    "",
			format:   0,
			expected: int64(0),
		},
		{
			name:     "count with punctuation",
			input:    "Hello, World!",
			format:   0,
			expected: int64(2),
		},
		{
			name:     "array of words",
			input:    "Hello World",
			format:   1,
			expected: []string{"Hello", "World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *types.Value
			if tt.format == 0 {
				result = StrWordCount(types.NewString(tt.input))
			} else {
				result = StrWordCount(types.NewString(tt.input), types.NewInt(tt.format))
			}

			switch expected := tt.expected.(type) {
			case int64:
				if result.ToInt() != expected {
					t.Errorf("Expected %d, got %d", expected, result.ToInt())
				}
			case []string:
				arr := result.ToArray()
				if arr.Len() != len(expected) {
					t.Errorf("Expected %d words, got %d", len(expected), arr.Len())
					return
				}
				for i, exp := range expected {
					val, _ := arr.Get(types.NewInt(int64(i)))
					if val.ToString() != exp {
						t.Errorf("Word %d: expected %q, got %q", i, exp, val.ToString())
					}
				}
			}
		})
	}
}

func TestStrWordCountWithPositions(t *testing.T) {
	result := StrWordCount(types.NewString("Hello World"), types.NewInt(2))
	arr := result.ToArray()

	// Should have 2 entries
	if arr.Len() != 2 {
		t.Errorf("Expected 2 entries, got %d", arr.Len())
		return
	}

	// Position 0 should be "Hello"
	val, found := arr.Get(types.NewInt(0))
	if !found || val.ToString() != "Hello" {
		t.Errorf("Position 0: expected 'Hello', got %q (found=%v)", val.ToString(), found)
	}

	// Position 6 should be "World"
	val, found = arr.Get(types.NewInt(6))
	if !found || val.ToString() != "World" {
		t.Errorf("Position 6: expected 'World', got %q (found=%v)", val.ToString(), found)
	}
}

func TestStrWordCountNil(t *testing.T) {
	result := StrWordCount(nil)
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for nil, got %d", result.ToInt())
	}
}

func TestSubstrCount(t *testing.T) {
	tests := []struct {
		name     string
		haystack string
		needle   string
		offset   *int64
		length   *int64
		expected int64
	}{
		{
			name:     "basic count",
			haystack: "aaaa",
			needle:   "a",
			expected: 4,
		},
		{
			name:     "multiple occurrences",
			haystack: "This is a test is test",
			needle:   "is",
			expected: 3,
		},
		{
			name:     "no occurrences",
			haystack: "Hello World",
			needle:   "xyz",
			expected: 0,
		},
		{
			name:     "needle longer than haystack",
			haystack: "abc",
			needle:   "abcdef",
			expected: 0,
		},
		{
			name:     "with offset",
			haystack: "aaaa",
			needle:   "a",
			offset:   ptrInt64(1),
			expected: 3,
		},
		{
			name:     "with offset and length",
			haystack: "aaaa",
			needle:   "a",
			offset:   ptrInt64(1),
			length:   ptrInt64(2),
			expected: 2,
		},
		{
			name:     "empty needle",
			haystack: "test",
			needle:   "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.offset != nil {
				args = append(args, types.NewInt(*tt.offset))
				if tt.length != nil {
					args = append(args, types.NewInt(*tt.length))
				}
			}
			result := SubstrCount(types.NewString(tt.haystack), types.NewString(tt.needle), args...)
			if result.ToInt() != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result.ToInt())
			}
		})
	}
}

func TestSubstrCountNegativeOffset(t *testing.T) {
	// Negative offset counts from end
	result := SubstrCount(types.NewString("aaaa"), types.NewString("a"), types.NewInt(-2))
	// "aa" - 2 occurrences
	if result.ToInt() != 2 {
		t.Errorf("Expected 2, got %d", result.ToInt())
	}
}

func TestSubstrReplace(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		replacement string
		offset      int64
		length      *int64
		expected    string
	}{
		{
			name:        "basic replace",
			str:         "Hello World",
			replacement: "PHP",
			offset:      6,
			length:      ptrInt64(5),
			expected:    "Hello PHP",
		},
		{
			name:        "replace at start",
			str:         "Hello World",
			replacement: "Hi",
			offset:      0,
			length:      ptrInt64(5),
			expected:    "Hi World",
		},
		{
			name:        "insert (length 0)",
			str:         "Hello World",
			replacement: " Beautiful",
			offset:      5,
			length:      ptrInt64(0),
			expected:    "Hello Beautiful World",
		},
		{
			name:        "replace to end (no length)",
			str:         "Hello World",
			replacement: "PHP",
			offset:      6,
			length:      nil,
			expected:    "Hello PHP",
		},
		{
			name:        "negative offset",
			str:         "Hello World",
			replacement: "PHP",
			offset:      -5,
			length:      ptrInt64(5),
			expected:    "Hello PHP",
		},
		{
			name:        "negative length",
			str:         "Hello World",
			replacement: "PHP",
			offset:      0,
			length:      ptrInt64(-6),
			expected:    "PHP World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.length != nil {
				args = append(args, types.NewInt(*tt.length))
			}
			result := SubstrReplace(
				types.NewString(tt.str),
				types.NewString(tt.replacement),
				types.NewInt(tt.offset),
				args...,
			)
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestSubstrReplaceNil(t *testing.T) {
	result := SubstrReplace(nil, types.NewString("replace"), types.NewInt(0))
	if result.ToString() != "" {
		t.Errorf("Expected empty string for nil, got %q", result.ToString())
	}
}

// Helper function to create pointer to int64
func ptrInt64(i int64) *int64 {
	return &i
}

// ============================================================================
// String Search Tests (strpbrk, strspn, strcspn)
// ============================================================================

func TestStrpbrk(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		characters string
		expected   interface{} // string or false
	}{
		{
			name:       "basic match",
			str:        "Hello World",
			characters: "W",
			expected:   "World",
		},
		{
			name:       "match first character",
			str:        "Hello World",
			characters: "H",
			expected:   "Hello World",
		},
		{
			name:       "match multiple chars - returns from first match",
			str:        "Hello World",
			characters: "oW",
			expected:   "o World", // 'o' at position 4 comes before 'W' at position 6
		},
		{
			name:       "no match",
			str:        "Hello World",
			characters: "xyz",
			expected:   false,
		},
		{
			name:       "empty characters",
			str:        "Hello",
			characters: "",
			expected:   false,
		},
		{
			name:       "empty string",
			str:        "",
			characters: "abc",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Strpbrk(types.NewString(tt.str), types.NewString(tt.characters))
			switch expected := tt.expected.(type) {
			case string:
				if result.ToString() != expected {
					t.Errorf("Expected %q, got %q", expected, result.ToString())
				}
			case bool:
				if result.Type() != types.TypeBool || result.ToBool() != expected {
					t.Errorf("Expected false, got %v", result)
				}
			}
		})
	}
}

func TestStrpbrkNil(t *testing.T) {
	result := Strpbrk(nil, types.NewString("abc"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("Expected false for nil string, got %v", result)
	}

	result = Strpbrk(types.NewString("test"), nil)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("Expected false for nil characters, got %v", result)
	}
}

func TestStrspn(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		characters string
		offset     *int64
		length     *int64
		expected   int64
	}{
		{
			name:       "basic count",
			str:        "42 is the answer",
			characters: "1234567890",
			expected:   2, // "42" are digits
		},
		{
			name:       "all match",
			str:        "12345",
			characters: "1234567890",
			expected:   5,
		},
		{
			name:       "no match",
			str:        "Hello",
			characters: "1234567890",
			expected:   0,
		},
		{
			name:       "with offset",
			str:        "hello world",
			characters: "worldabcdefghijklmnopqrstuvwxyz",
			offset:     ptrInt64(6),
			expected:   5, // "world" from position 6
		},
		{
			name:       "with negative offset",
			str:        "hello world",
			characters: "worldabcdefghijklmnopqrstuvwxyz",
			offset:     ptrInt64(-5),
			expected:   5, // "world" from -5
		},
		{
			name:       "with offset and length",
			str:        "aaa123",
			characters: "a",
			offset:     ptrInt64(0),
			length:     ptrInt64(2),
			expected:   2,
		},
		{
			name:       "empty characters",
			str:        "hello",
			characters: "",
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.offset != nil {
				args = append(args, types.NewInt(*tt.offset))
				if tt.length != nil {
					args = append(args, types.NewInt(*tt.length))
				}
			}
			result := Strspn(types.NewString(tt.str), types.NewString(tt.characters), args...)
			if result.ToInt() != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result.ToInt())
			}
		})
	}
}

func TestStrspnNil(t *testing.T) {
	result := Strspn(nil, types.NewString("abc"))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for nil string, got %d", result.ToInt())
	}
}

func TestStrcspn(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		characters string
		offset     *int64
		length     *int64
		expected   int64
	}{
		{
			name:       "basic count",
			str:        "hello world",
			characters: " ",
			expected:   5, // "hello" before space
		},
		{
			name:       "no match - count all",
			str:        "Hello",
			characters: "xyz",
			expected:   5,
		},
		{
			name:       "match at start",
			str:        " hello",
			characters: " ",
			expected:   0,
		},
		{
			name:       "with offset",
			str:        "hello world test",
			characters: " ",
			offset:     ptrInt64(6),
			expected:   5, // "world" from position 6
		},
		{
			name:       "with negative offset",
			str:        "hello world test",
			characters: " ",
			offset:     ptrInt64(-4),
			expected:   4, // "test" from -4
		},
		{
			name:       "with offset and length",
			str:        "hello world",
			characters: " ",
			offset:     ptrInt64(0),
			length:     ptrInt64(3),
			expected:   3,
		},
		{
			name:       "empty characters - count all",
			str:        "hello",
			characters: "",
			expected:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.offset != nil {
				args = append(args, types.NewInt(*tt.offset))
				if tt.length != nil {
					args = append(args, types.NewInt(*tt.length))
				}
			}
			result := Strcspn(types.NewString(tt.str), types.NewString(tt.characters), args...)
			if result.ToInt() != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result.ToInt())
			}
		})
	}
}

func TestStrcspnNil(t *testing.T) {
	result := Strcspn(nil, types.NewString("abc"))
	if result.ToInt() != 0 {
		t.Errorf("Expected 0 for nil string, got %d", result.ToInt())
	}
}

// ============================================================================
// String Formatting Tests (vsprintf, vprintf, number_format, sscanf, str_getcsv)
// ============================================================================

func TestVsprintf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		values   []interface{}
		expected string
	}{
		{
			name:     "string format",
			format:   "Hello %s",
			values:   []interface{}{"World"},
			expected: "Hello World",
		},
		{
			name:     "integer format",
			format:   "Number: %d",
			values:   []interface{}{int64(42)},
			expected: "Number: 42",
		},
		{
			name:     "multiple values",
			format:   "%s has %d apples",
			values:   []interface{}{"John", int64(5)},
			expected: "John has 5 apples",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := types.NewEmptyArray()
			for _, v := range tt.values {
				switch val := v.(type) {
				case string:
					arr.Append(types.NewString(val))
				case int64:
					arr.Append(types.NewInt(val))
				}
			}
			result := Vsprintf(types.NewString(tt.format), types.NewArray(arr))
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestVsprintfNil(t *testing.T) {
	result := Vsprintf(nil, types.NewArray(types.NewEmptyArray()))
	if result.ToString() != "" {
		t.Errorf("Expected empty string for nil format, got %q", result.ToString())
	}
}

func TestVprintf(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewString("World"))
	result := Vprintf(types.NewString("Hello %s"), types.NewArray(arr))
	// Returns length of formatted string
	if result.ToInt() != 11 { // "Hello World" is 11 chars
		t.Errorf("Expected 11, got %d", result.ToInt())
	}
}

func TestNumberFormat(t *testing.T) {
	tests := []struct {
		name         string
		number       float64
		decimals     int64
		decSep       string
		thousandsSep string
		expected     string
	}{
		{
			name:     "integer",
			number:   1234,
			expected: "1,234",
		},
		{
			name:     "with decimals",
			number:   1234.567,
			decimals: 2,
			expected: "1,234.57",
		},
		{
			name:         "custom separators",
			number:       1234567.89,
			decimals:     2,
			decSep:       ",",
			thousandsSep: " ",
			expected:     "1 234 567,89",
		},
		{
			name:     "small number",
			number:   123,
			expected: "123",
		},
		{
			name:     "negative",
			number:   -1234.5,
			decimals: 1,
			expected: "-1,234.5",
		},
		{
			name:         "no thousands separator",
			number:       1234567,
			decimals:     0,
			decSep:       ".",
			thousandsSep: "",
			expected:     "1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.decimals > 0 || tt.decSep != "" || tt.thousandsSep != "" {
				args = append(args, types.NewInt(tt.decimals))
				if tt.decSep != "" || tt.thousandsSep != "" {
					if tt.decSep == "" {
						tt.decSep = "."
					}
					args = append(args, types.NewString(tt.decSep))
					if tt.thousandsSep != "" || tt.thousandsSep == "" {
						args = append(args, types.NewString(tt.thousandsSep))
					}
				}
			}
			result := NumberFormat(types.NewFloat(tt.number), args...)
			if result.ToString() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.ToString())
			}
		})
	}
}

func TestNumberFormatNil(t *testing.T) {
	result := NumberFormat(nil)
	if result.ToString() != "0" {
		t.Errorf("Expected '0' for nil, got %q", result.ToString())
	}
}

func TestSscanf(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		format   string
		expected []interface{} // int64 or float64 or string
	}{
		{
			name:     "single integer",
			str:      "42",
			format:   "%d",
			expected: []interface{}{int64(42)},
		},
		{
			name:     "integer with text",
			str:      "Age: 25",
			format:   "Age: %d",
			expected: []interface{}{int64(25)},
		},
		{
			name:     "multiple values",
			str:      "John 30",
			format:   "%s %d",
			expected: []interface{}{"John", int64(30)},
		},
		{
			name:     "float value",
			str:      "3.14",
			format:   "%f",
			expected: []interface{}{float64(3.14)},
		},
		{
			name:     "hex value",
			str:      "ff",
			format:   "%x",
			expected: []interface{}{int64(255)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sscanf(types.NewString(tt.str), types.NewString(tt.format))
			arr := result.ToArray()
			if arr.Len() != len(tt.expected) {
				t.Errorf("Expected %d values, got %d", len(tt.expected), arr.Len())
				return
			}
			for i, exp := range tt.expected {
				val, _ := arr.Get(types.NewInt(int64(i)))
				switch expected := exp.(type) {
				case int64:
					if val.ToInt() != expected {
						t.Errorf("Value %d: expected %d, got %d", i, expected, val.ToInt())
					}
				case float64:
					if val.ToFloat() != expected {
						t.Errorf("Value %d: expected %f, got %f", i, expected, val.ToFloat())
					}
				case string:
					if val.ToString() != expected {
						t.Errorf("Value %d: expected %q, got %q", i, expected, val.ToString())
					}
				}
			}
		})
	}
}

func TestSscanfNil(t *testing.T) {
	result := Sscanf(nil, types.NewString("%d"))
	if result.Type() != types.TypeNull {
		t.Errorf("Expected null for nil string, got %v", result.Type())
	}
}

func TestStrGetcsv(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		separator string
		enclosure string
		escape    string
		expected  []string
	}{
		{
			name:     "simple",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with spaces",
			input:    "hello, world, test",
			expected: []string{"hello", " world", " test"},
		},
		{
			name:     "quoted fields",
			input:    `"hello","world"`,
			expected: []string{"hello", "world"},
		},
		{
			name:     "quoted with comma",
			input:    `"hello, world",test`,
			expected: []string{"hello, world", "test"},
		},
		{
			name:      "custom separator",
			input:     "a;b;c",
			separator: ";",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:     "empty fields",
			input:    "a,,c",
			expected: []string{"a", "", "c"},
		},
		{
			name:     "single field",
			input:    "hello",
			expected: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []*types.Value
			if tt.separator != "" {
				args = append(args, types.NewString(tt.separator))
			}
			if tt.enclosure != "" {
				args = append(args, types.NewString(tt.enclosure))
			}
			if tt.escape != "" {
				args = append(args, types.NewString(tt.escape))
			}
			result := StrGetcsv(types.NewString(tt.input), args...)
			arr := result.ToArray()
			if arr.Len() != len(tt.expected) {
				t.Errorf("Expected %d fields, got %d", len(tt.expected), arr.Len())
				return
			}
			for i, exp := range tt.expected {
				val, _ := arr.Get(types.NewInt(int64(i)))
				if val.ToString() != exp {
					t.Errorf("Field %d: expected %q, got %q", i, exp, val.ToString())
				}
			}
		})
	}
}

func TestStrGetcsvNil(t *testing.T) {
	result := StrGetcsv(nil)
	arr := result.ToArray()
	if arr.Len() != 0 {
		t.Errorf("Expected empty array for nil, got %d elements", arr.Len())
	}
}

// ============================================================================
// Base64 Encoding Tests
// ============================================================================

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "SGVsbG8gV29ybGQ="},
		{"", ""},
		{"a", "YQ=="},
		{"ab", "YWI="},
		{"abc", "YWJj"},
		{"PHP is great!", "UEhQIGlzIGdyZWF0IQ=="},
	}

	for _, tt := range tests {
		result := Base64Encode(types.NewString(tt.input))
		if result.ToString() != tt.expected {
			t.Errorf("Base64Encode(%q) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestBase64EncodeNil(t *testing.T) {
	result := Base64Encode(nil)
	if result.ToString() != "" {
		t.Errorf("Expected empty string for nil, got %q", result.ToString())
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SGVsbG8gV29ybGQ=", "Hello World"},
		{"", ""},
		{"YQ==", "a"},
		{"YWI=", "ab"},
		{"YWJj", "abc"},
		{"UEhQIGlzIGdyZWF0IQ==", "PHP is great!"},
	}

	for _, tt := range tests {
		result := Base64Decode(types.NewString(tt.input))
		if result.ToString() != tt.expected {
			t.Errorf("Base64Decode(%q) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestBase64DecodeInvalid(t *testing.T) {
	// Invalid base64 should return false
	result := Base64Decode(types.NewString("not-valid-base64!!!"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for invalid base64")
	}
}

func TestBase64DecodeNil(t *testing.T) {
	result := Base64Decode(nil)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for nil")
	}
}

// ============================================================================
// URL Parsing Tests
// ============================================================================

func TestParseUrl(t *testing.T) {
	url := "https://user:pass@example.com:8080/path/to/page?query=value&foo=bar#section"
	result := ParseUrl(types.NewString(url))

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array, got %v", result.Type())
	}

	arr := result.ToArray()

	// Test scheme
	scheme, _ := arr.Get(types.NewString("scheme"))
	if scheme.ToString() != "https" {
		t.Errorf("Expected scheme 'https', got %q", scheme.ToString())
	}

	// Test host
	host, _ := arr.Get(types.NewString("host"))
	if host.ToString() != "example.com" {
		t.Errorf("Expected host 'example.com', got %q", host.ToString())
	}

	// Test port
	port, _ := arr.Get(types.NewString("port"))
	if port.ToInt() != 8080 {
		t.Errorf("Expected port 8080, got %d", port.ToInt())
	}

	// Test user
	user, _ := arr.Get(types.NewString("user"))
	if user.ToString() != "user" {
		t.Errorf("Expected user 'user', got %q", user.ToString())
	}

	// Test pass
	pass, _ := arr.Get(types.NewString("pass"))
	if pass.ToString() != "pass" {
		t.Errorf("Expected pass 'pass', got %q", pass.ToString())
	}

	// Test path
	path, _ := arr.Get(types.NewString("path"))
	if path.ToString() != "/path/to/page" {
		t.Errorf("Expected path '/path/to/page', got %q", path.ToString())
	}

	// Test query
	query, _ := arr.Get(types.NewString("query"))
	if query.ToString() != "query=value&foo=bar" {
		t.Errorf("Expected query 'query=value&foo=bar', got %q", query.ToString())
	}

	// Test fragment
	fragment, _ := arr.Get(types.NewString("fragment"))
	if fragment.ToString() != "section" {
		t.Errorf("Expected fragment 'section', got %q", fragment.ToString())
	}
}

func TestParseUrlComponent(t *testing.T) {
	url := "https://example.com:8080/path"

	// Get only scheme
	result := ParseUrl(types.NewString(url), types.NewInt(PHP_URL_SCHEME))
	if result.ToString() != "https" {
		t.Errorf("Expected scheme 'https', got %q", result.ToString())
	}

	// Get only host
	result = ParseUrl(types.NewString(url), types.NewInt(PHP_URL_HOST))
	if result.ToString() != "example.com" {
		t.Errorf("Expected host 'example.com', got %q", result.ToString())
	}

	// Get only port
	result = ParseUrl(types.NewString(url), types.NewInt(PHP_URL_PORT))
	if result.ToInt() != 8080 {
		t.Errorf("Expected port 8080, got %d", result.ToInt())
	}

	// Get only path
	result = ParseUrl(types.NewString(url), types.NewInt(PHP_URL_PATH))
	if result.ToString() != "/path" {
		t.Errorf("Expected path '/path', got %q", result.ToString())
	}
}

func TestParseUrlEmpty(t *testing.T) {
	result := ParseUrl(types.NewString(""))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for empty URL")
	}
}

func TestParseUrlNil(t *testing.T) {
	result := ParseUrl(nil)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for nil URL")
	}
}

// ============================================================================
// HTTP Build Query Tests
// ============================================================================

func TestHttpBuildQuery(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John Doe"))
	arr.Set(types.NewString("age"), types.NewInt(30))

	result := HttpBuildQuery(types.NewArray(arr))
	resultStr := result.ToString()

	// Check that it contains expected parts (order may vary)
	if !containsSubstring(resultStr, "name=John+Doe") && !containsSubstring(resultStr, "name=John%20Doe") {
		t.Errorf("Expected 'name=John+Doe' or 'name=John%%20Doe' in result, got %q", resultStr)
	}
	if !containsSubstring(resultStr, "age=30") {
		t.Errorf("Expected 'age=30' in result, got %q", resultStr)
	}
}

func TestHttpBuildQueryNumericPrefix(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewInt(0), types.NewString("value0"))
	arr.Set(types.NewInt(1), types.NewString("value1"))

	result := HttpBuildQuery(types.NewArray(arr), types.NewString("item_"))
	resultStr := result.ToString()

	if !containsSubstring(resultStr, "item_0=value0") {
		t.Errorf("Expected 'item_0=value0' in result, got %q", resultStr)
	}
	if !containsSubstring(resultStr, "item_1=value1") {
		t.Errorf("Expected 'item_1=value1' in result, got %q", resultStr)
	}
}

func TestHttpBuildQueryCustomSeparator(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewString("1"))
	arr.Set(types.NewString("b"), types.NewString("2"))

	result := HttpBuildQuery(types.NewArray(arr), types.NewString(""), types.NewString(";"))
	resultStr := result.ToString()

	if !containsSubstring(resultStr, ";") {
		t.Errorf("Expected ';' separator in result, got %q", resultStr)
	}
}

func TestHttpBuildQueryEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	result := HttpBuildQuery(types.NewArray(arr))
	if result.ToString() != "" {
		t.Errorf("Expected empty string for empty array, got %q", result.ToString())
	}
}

func TestHttpBuildQueryNil(t *testing.T) {
	result := HttpBuildQuery(nil)
	if result.ToString() != "" {
		t.Errorf("Expected empty string for nil, got %q", result.ToString())
	}
}

// Helper function for tests
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || containsSubstring(s[1:], substr)))
}
