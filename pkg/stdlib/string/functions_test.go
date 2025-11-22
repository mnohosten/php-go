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
