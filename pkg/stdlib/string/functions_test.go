package string

import (
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
