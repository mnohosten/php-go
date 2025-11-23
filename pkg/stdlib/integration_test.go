package stdlib_test

import (
	"testing"

	"github.com/krizos/php-go/pkg/stdlib/array"
	"github.com/krizos/php-go/pkg/stdlib/json"
	"github.com/krizos/php-go/pkg/stdlib/pcre"
	stringfuncs "github.com/krizos/php-go/pkg/stdlib/string"
	varfuncs "github.com/krizos/php-go/pkg/stdlib/var"
	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Integration Tests
// ============================================================================

// TestArrayJsonIntegration tests array and JSON functions working together
func TestArrayJsonIntegration(t *testing.T) {
	// Create an array
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))
	arr.Set(types.NewString("city"), types.NewString("New York"))

	// Encode to JSON
	jsonStr := json.JsonEncode(types.NewArray(arr))

	// Decode back
	assoc := types.NewBool(true)
	decoded := json.JsonDecode(jsonStr, assoc)

	if decoded.Type() != types.TypeArray {
		t.Error("Integration: JSON decode should return array")
	}

	// Verify array keys
	decodedArr := decoded.ToArray()
	keys := array.ArrayKeys(types.NewArray(decodedArr))

	if keys.Type() != types.TypeArray {
		t.Error("Integration: ArrayKeys should work on decoded array")
	}

	keysArr := keys.ToArray()
	if keysArr.Len() != 3 {
		t.Errorf("Integration: Expected 3 keys, got %d", keysArr.Len())
	}
}

// TestStringPcreIntegration tests string and PCRE functions working together
func TestStringPcreIntegration(t *testing.T) {
	// Create a string
	text := types.NewString("Hello World, this is a test! Contact: user@example.com")

	// Convert to uppercase
	upper := stringfuncs.Strtoupper(text)

	// Search for email pattern
	pattern := types.NewString("/[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}/i")
	matches := types.NewArray(types.NewEmptyArray())
	result := pcre.PregMatch(pattern, upper, matches)

	if result.ToInt() != 1 {
		t.Error("Integration: Should find email in uppercase string")
	}

	// Extract the email
	matchesArr := matches.ToArray()
	email, _ := matchesArr.Get(types.NewInt(0))

	// Convert back to lowercase
	lowerEmail := stringfuncs.Strtolower(email)

	if lowerEmail.ToString() != "user@example.com" {
		t.Errorf("Integration: Expected 'user@example.com', got '%s'", lowerEmail.ToString())
	}
}

// TestArrayStringIntegration tests array and string functions working together
func TestArrayStringIntegration(t *testing.T) {
	// Create an array of strings
	arr := types.NewEmptyArray()
	arr.Append(types.NewString("apple"))
	arr.Append(types.NewString("banana"))
	arr.Append(types.NewString("cherry"))

	// Join them
	glue := types.NewString(", ")
	joined := stringfuncs.Implode(glue, types.NewArray(arr))

	if joined.ToString() != "apple, banana, cherry" {
		t.Errorf("Integration: Expected 'apple, banana, cherry', got '%s'", joined.ToString())
	}

	// Split back
	exploded := stringfuncs.Explode(glue, joined)

	if exploded.Type() != types.TypeArray {
		t.Error("Integration: Explode should return array")
	}

	explodedArr := exploded.ToArray()
	if explodedArr.Len() != 3 {
		t.Errorf("Integration: Expected 3 elements, got %d", explodedArr.Len())
	}

	// Sort the array
	sorted := array.Sort(types.NewArray(explodedArr))
	if !sorted.ToBool() {
		t.Error("Integration: Sort should succeed")
	}
}

// TestVarDumpJsonIntegration tests var dump with JSON-decoded data
func TestVarDumpJsonIntegration(t *testing.T) {
	// Create complex data structure
	jsonStr := types.NewString(`{
		"users": [
			{"name": "John", "age": 30},
			{"name": "Jane", "age": 25}
		],
		"total": 2
	}`)

	// Decode
	decoded := json.JsonDecode(jsonStr)

	// Verify it's an object
	if !varfuncs.IsObject(decoded).ToBool() {
		t.Error("Integration: Decoded JSON should be object")
	}

	// Use var_export to get string representation
	exported := varfuncs.VarExport(decoded, types.NewBool(true))

	if exported.Type() != types.TypeString {
		t.Error("Integration: VarExport should return string")
	}

	exportedStr := exported.ToString()
	if exportedStr == "" {
		t.Error("Integration: VarExport should not return empty string")
	}
}

// TestPcreStringReplace tests PCRE replace with string functions
func TestPcreStringReplace(t *testing.T) {
	// Create a text with multiple numbers
	text := types.NewString("I have 10 apples and 20 oranges")

	// Replace all numbers with 'X'
	pattern := types.NewString("/\\d+/")
	replacement := types.NewString("X")
	replaced := pcre.PregReplace(pattern, replacement, text)

	expected := "I have X apples and X oranges"
	if replaced.ToString() != expected {
		t.Errorf("Integration: Expected '%s', got '%s'", expected, replaced.ToString())
	}

	// Verify replacement worked
	if !stringfuncs.Strpos(replaced, types.NewString("X")).ToBool() && stringfuncs.Strpos(replaced, types.NewString("X")).ToInt() >= 0 {
		t.Error("Integration: String should contain 'X'")
	}
}

// TestArrayFilterMap tests array operations
func TestArrayFilterMap(t *testing.T) {
	// Create array with values
	arr := types.NewEmptyArray()
	arr.Append(types.NewString("apple"))
	arr.Append(types.NewString("banana"))
	arr.Append(types.NewString("cherry"))
	arr.Append(types.NewString("date"))
	arr.Append(types.NewString("elderberry"))

	// Count elements
	count := array.Count(types.NewArray(arr))
	if count.ToInt() != 5 {
		t.Errorf("Integration: Expected count 5, got %d", count.ToInt())
	}

	// Check if value exists
	needle := types.NewString("banana")
	exists := array.InArray(needle, types.NewArray(arr))
	if !exists.ToBool() {
		t.Error("Integration: InArray should find 'banana'")
	}

	// Reverse the array
	reversed := array.ArrayReverse(types.NewArray(arr))
	if reversed.Type() != types.TypeArray {
		t.Error("Integration: ArrayReverse should return array")
	}

	reversedArr := reversed.ToArray()
	first, _ := reversedArr.Get(types.NewInt(0))
	if first.ToString() != "elderberry" {
		t.Errorf("Integration: First element of reversed array should be 'elderberry', got '%s'", first.ToString())
	}
}
