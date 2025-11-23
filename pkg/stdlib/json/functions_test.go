package json

import (
	"strings"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// JSON Encode Tests
// ============================================================================

func TestJsonEncodeNull(t *testing.T) {
	result := JsonEncode(types.NewNull())
	if result.ToString() != "null" {
		t.Errorf("JsonEncode(null) = %v, want 'null'", result.ToString())
	}
}

func TestJsonEncodeBool(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected string
	}{
		{types.NewBool(true), "true"},
		{types.NewBool(false), "false"},
	}

	for _, tt := range tests {
		result := JsonEncode(tt.input)
		if result.ToString() != tt.expected {
			t.Errorf("JsonEncode(%v) = %v, want %v", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestJsonEncodeInt(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected string
	}{
		{types.NewInt(0), "0"},
		{types.NewInt(42), "42"},
		{types.NewInt(-100), "-100"},
	}

	for _, tt := range tests {
		result := JsonEncode(tt.input)
		if result.ToString() != tt.expected {
			t.Errorf("JsonEncode(%v) = %v, want %v", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestJsonEncodeFloat(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected string
	}{
		{types.NewFloat(3.14), "3.14"},
		{types.NewFloat(0.0), "0"},
		{types.NewFloat(-2.5), "-2.5"},
	}

	for _, tt := range tests {
		result := JsonEncode(tt.input)
		if result.ToString() != tt.expected {
			t.Errorf("JsonEncode(%v) = %v, want %v", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestJsonEncodeString(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected string
	}{
		{types.NewString("hello"), `"hello"`},
		{types.NewString(""), `""`},
		{types.NewString("hello\nworld"), `"hello\nworld"`},
		{types.NewString(`hello"world`), `"hello\"world"`},
	}

	for _, tt := range tests {
		result := JsonEncode(tt.input)
		if result.ToString() != tt.expected {
			t.Errorf("JsonEncode(%v) = %v, want %v", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestJsonEncodeArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))

	result := JsonEncode(types.NewArray(arr))
	expected := "[1,2,3]"

	if result.ToString() != expected {
		t.Errorf("JsonEncode(array) = %v, want %v", result.ToString(), expected)
	}
}

func TestJsonEncodeAssociativeArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))

	result := JsonEncode(types.NewArray(arr))
	str := result.ToString()

	// Check it's an object
	if !strings.HasPrefix(str, "{") || !strings.HasSuffix(str, "}") {
		t.Errorf("JsonEncode(associative array) should return object, got %v", str)
	}

	// Check it contains the keys
	if !strings.Contains(str, `"name"`) || !strings.Contains(str, `"John"`) {
		t.Errorf("JsonEncode(associative array) should contain name:John, got %v", str)
	}
}

func TestJsonEncodeEmptyArray(t *testing.T) {
	arr := types.NewEmptyArray()
	result := JsonEncode(types.NewArray(arr))

	if result.ToString() != "[]" {
		t.Errorf("JsonEncode(empty array) = %v, want '[]'", result.ToString())
	}
}

func TestJsonEncodeNestedArray(t *testing.T) {
	inner := types.NewEmptyArray()
	inner.Append(types.NewInt(1))
	inner.Append(types.NewInt(2))

	outer := types.NewEmptyArray()
	outer.Append(types.NewArray(inner))
	outer.Append(types.NewInt(3))

	result := JsonEncode(types.NewArray(outer))
	expected := "[[1,2],3]"

	if result.ToString() != expected {
		t.Errorf("JsonEncode(nested array) = %v, want %v", result.ToString(), expected)
	}
}

func TestJsonEncodeForceObject(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))

	flags := types.NewInt(JSON_FORCE_OBJECT)
	result := JsonEncode(types.NewArray(arr), flags)

	str := result.ToString()
	if !strings.HasPrefix(str, "{") {
		t.Errorf("JsonEncode with JSON_FORCE_OBJECT should return object, got %v", str)
	}
}

func TestJsonEncodeUnescapedSlashes(t *testing.T) {
	str := types.NewString("http://example.com")
	flags := types.NewInt(JSON_UNESCAPED_SLASHES)

	result := JsonEncode(str, flags)
	expected := `"http://example.com"`

	if result.ToString() != expected {
		t.Errorf("JsonEncode with JSON_UNESCAPED_SLASHES = %v, want %v", result.ToString(), expected)
	}
}

func TestJsonEncodeObject(t *testing.T) {
	class := types.NewClassEntry("TestClass")

	// Add property definition
	class.Properties["name"] = &types.PropertyDef{
		Name:       "name",
		Visibility: types.VisibilityPublic,
		IsStatic:   false,
	}

	obj := types.NewObjectFromClass(class)
	obj.Properties["name"] = &types.Property{
		Value:      types.NewString("test"),
		Visibility: types.VisibilityPublic,
	}

	result := JsonEncode(types.NewObject(obj))
	str := result.ToString()

	if !strings.Contains(str, `"name"`) || !strings.Contains(str, `"test"`) {
		t.Errorf("JsonEncode(object) should contain name:test, got %v", str)
	}
}

// ============================================================================
// JSON Decode Tests
// ============================================================================

func TestJsonDecodeNull(t *testing.T) {
	result := JsonDecode(types.NewString("null"))
	if result.Type() != types.TypeNull {
		t.Errorf("JsonDecode('null') should return NULL, got %v", result.Type())
	}
}

func TestJsonDecodeBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		result := JsonDecode(types.NewString(tt.input))
		if result.Type() != types.TypeBool || result.ToBool() != tt.expected {
			t.Errorf("JsonDecode(%v) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

func TestJsonDecodeInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"42", 42},
		{"-100", -100},
	}

	for _, tt := range tests {
		result := JsonDecode(types.NewString(tt.input))
		if result.Type() != types.TypeInt || result.ToInt() != tt.expected {
			t.Errorf("JsonDecode(%v) = %v, want %v", tt.input, result.ToInt(), tt.expected)
		}
	}
}

func TestJsonDecodeFloat(t *testing.T) {
	result := JsonDecode(types.NewString("3.14"))
	if result.Type() != types.TypeFloat || result.ToFloat() != 3.14 {
		t.Errorf("JsonDecode('3.14') = %v, want 3.14", result.ToFloat())
	}
}

func TestJsonDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`""`, ""},
		{`"hello\nworld"`, "hello\nworld"},
	}

	for _, tt := range tests {
		result := JsonDecode(types.NewString(tt.input))
		if result.Type() != types.TypeString || result.ToString() != tt.expected {
			t.Errorf("JsonDecode(%v) = %v, want %v", tt.input, result.ToString(), tt.expected)
		}
	}
}

func TestJsonDecodeArray(t *testing.T) {
	result := JsonDecode(types.NewString("[1,2,3]"))

	if result.Type() != types.TypeArray {
		t.Errorf("JsonDecode('[1,2,3]') should return array, got %v", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() != 3 {
		t.Errorf("JsonDecode('[1,2,3]') should return array of length 3, got %d", arr.Len())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("First element should be 1, got %d", val.ToInt())
	}
}

func TestJsonDecodeObject(t *testing.T) {
	result := JsonDecode(types.NewString(`{"name":"John","age":30}`))

	if result.Type() != types.TypeObject {
		t.Errorf("JsonDecode(object) should return object, got %v", result.Type())
	}

	obj := result.ToObject()
	nameVal, exists := obj.GetProperty("name", nil)
	if !exists || nameVal.ToString() != "John" {
		t.Errorf("Object should have property 'name' = 'John'")
	}
}

func TestJsonDecodeObjectAsArray(t *testing.T) {
	assoc := types.NewBool(true)
	result := JsonDecode(types.NewString(`{"name":"John","age":30}`), assoc)

	if result.Type() != types.TypeArray {
		t.Errorf("JsonDecode(object, true) should return array, got %v", result.Type())
	}

	arr := result.ToArray()
	val, exists := arr.Get(types.NewString("name"))
	if !exists || val.ToString() != "John" {
		t.Errorf("Array should have key 'name' = 'John'")
	}
}

func TestJsonDecodeInvalidJson(t *testing.T) {
	result := JsonDecode(types.NewString(`{"invalid json`))

	if result.Type() != types.TypeNull {
		t.Errorf("JsonDecode(invalid) should return NULL, got %v", result.Type())
	}

	// Check error was set
	errCode := JsonLastError()
	if errCode.ToInt() == JSON_ERROR_NONE {
		t.Errorf("JsonLastError should not be JSON_ERROR_NONE after decode error")
	}
}

func TestJsonDecodeNestedStructure(t *testing.T) {
	json := `{
		"name": "John",
		"age": 30,
		"addresses": [
			{"city": "New York", "zip": "10001"},
			{"city": "Boston", "zip": "02101"}
		]
	}`

	result := JsonDecode(types.NewString(json))

	if result.Type() != types.TypeObject {
		t.Errorf("JsonDecode should return object, got %v", result.Type())
	}

	obj := result.ToObject()
	addresses, exists := obj.GetProperty("addresses", nil)
	if !exists || addresses.Type() != types.TypeArray {
		t.Errorf("Object should have 'addresses' array property")
	}

	arr := addresses.ToArray()
	if arr.Len() != 2 {
		t.Errorf("Addresses array should have 2 elements, got %d", arr.Len())
	}
}

// ============================================================================
// JSON Error Handling Tests
// ============================================================================

func TestJsonLastError(t *testing.T) {
	// Reset error state with valid encode
	JsonEncode(types.NewInt(42))

	result := JsonLastError()
	if result.ToInt() != JSON_ERROR_NONE {
		t.Errorf("JsonLastError after successful encode should be JSON_ERROR_NONE")
	}
}

func TestJsonLastErrorMsg(t *testing.T) {
	// Reset error state
	JsonEncode(types.NewInt(42))

	result := JsonLastErrorMsg()
	if result.ToString() != "No error" {
		t.Errorf("JsonLastErrorMsg should return 'No error', got %v", result.ToString())
	}
}

// ============================================================================
// Round-trip Tests
// ============================================================================

func TestJsonRoundTripSimple(t *testing.T) {
	original := types.NewInt(42)
	encoded := JsonEncode(original)
	decoded := JsonDecode(encoded)

	if decoded.ToInt() != 42 {
		t.Errorf("Round-trip failed: got %v, want 42", decoded.ToInt())
	}
}

func TestJsonRoundTripArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))

	encoded := JsonEncode(types.NewArray(arr))
	decoded := JsonDecode(encoded)

	if decoded.Type() != types.TypeArray {
		t.Errorf("Round-trip should return array")
	}

	decodedArr := decoded.ToArray()
	if decodedArr.Len() != 3 {
		t.Errorf("Round-trip array should have 3 elements, got %d", decodedArr.Len())
	}
}

func TestJsonRoundTripAssociativeArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("name"), types.NewString("John"))
	arr.Set(types.NewString("age"), types.NewInt(30))

	encoded := JsonEncode(types.NewArray(arr))
	assoc := types.NewBool(true)
	decoded := JsonDecode(encoded, assoc)

	if decoded.Type() != types.TypeArray {
		t.Errorf("Round-trip should return array")
	}

	decodedArr := decoded.ToArray()
	nameVal, exists := decodedArr.Get(types.NewString("name"))
	if !exists || nameVal.ToString() != "John" {
		t.Errorf("Round-trip should preserve associative array")
	}
}

// ============================================================================
// Special Cases Tests
// ============================================================================

func TestJsonEncodeEscaping(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{`hello"world`, `\"`},
		{"hello\nworld", `\n`},
		{"hello\tworld", `\t`},
	}

	for _, tt := range tests {
		result := JsonEncode(types.NewString(tt.input))
		if !strings.Contains(result.ToString(), tt.contains) {
			t.Errorf("JsonEncode(%q) should contain %q, got %v", tt.input, tt.contains, result.ToString())
		}
	}
}

func TestJsonDecodeUnicode(t *testing.T) {
	result := JsonDecode(types.NewString(`"hello \u0041 world"`))

	if result.Type() != types.TypeString {
		t.Errorf("JsonDecode unicode should return string")
	}

	// Go's JSON decoder handles unicode automatically
	str := result.ToString()
	if !strings.Contains(str, "A") {
		t.Errorf("JsonDecode should decode unicode escape, got %v", str)
	}
}

func TestJsonEncodeEmptyObject(t *testing.T) {
	arr := types.NewEmptyArray()
	flags := types.NewInt(JSON_FORCE_OBJECT)

	result := JsonEncode(types.NewArray(arr), flags)
	if result.ToString() != "{}" {
		t.Errorf("JsonEncode(empty, JSON_FORCE_OBJECT) = %v, want '{}'", result.ToString())
	}
}

// ============================================================================
// Additional JSON Encoding Flag Tests
// ============================================================================

func TestJsonEncodeHexQuot(t *testing.T) {
	str := types.NewString(`hello"world`)
	flags := types.NewInt(JSON_HEX_QUOT)
	result := JsonEncode(str, flags)

	if !strings.Contains(result.ToString(), "\\u0022") {
		t.Errorf("JSON_HEX_QUOT should encode quotes as \\u0022, got %v", result.ToString())
	}
}

func TestJsonEncodeHexApos(t *testing.T) {
	str := types.NewString("hello'world")
	flags := types.NewInt(JSON_HEX_APOS)
	result := JsonEncode(str, flags)

	if !strings.Contains(result.ToString(), "\\u0027") {
		t.Errorf("JSON_HEX_APOS should encode apostrophe as \\u0027, got %v", result.ToString())
	}
}

func TestJsonEncodeHexAmp(t *testing.T) {
	str := types.NewString("hello&world")
	flags := types.NewInt(JSON_HEX_AMP)
	result := JsonEncode(str, flags)

	if !strings.Contains(result.ToString(), "\\u0026") {
		t.Errorf("JSON_HEX_AMP should encode & as \\u0026, got %v", result.ToString())
	}
}

func TestJsonEncodeHexTag(t *testing.T) {
	str := types.NewString("hello<>world")
	flags := types.NewInt(JSON_HEX_TAG)
	result := JsonEncode(str, flags)

	if !strings.Contains(result.ToString(), "\\u003C") || !strings.Contains(result.ToString(), "\\u003E") {
		t.Errorf("JSON_HEX_TAG should encode < and > as unicode, got %v", result.ToString())
	}
}

func TestJsonEncodeUnescapedUnicode(t *testing.T) {
	str := types.NewString("hello世界")
	flags := types.NewInt(JSON_UNESCAPED_UNICODE)
	result := JsonEncode(str, flags)

	// Should contain the actual unicode characters, not \uXXXX
	if strings.Contains(result.ToString(), "\\u") {
		t.Errorf("JSON_UNESCAPED_UNICODE should not escape unicode, got %v", result.ToString())
	}
}

func TestJsonEncodeControlCharacters(t *testing.T) {
	// Test control characters < 0x20
	str := types.NewString("hello\x01\x02world")
	result := JsonEncode(str)

	// Should be escaped as \uXXXX
	if !strings.Contains(result.ToString(), "\\u00") {
		t.Errorf("Control characters should be escaped, got %v", result.ToString())
	}
}

func TestJsonEncodeBackslash(t *testing.T) {
	str := types.NewString("hello\\world")
	result := JsonEncode(str)

	if !strings.Contains(result.ToString(), "\\\\") {
		t.Errorf("Backslash should be escaped, got %v", result.ToString())
	}
}

func TestJsonEncodeSpecialEscapes(t *testing.T) {
	tests := []struct {
		input    string
		contains string
		desc     string
	}{
		{"\b", "\\b", "backspace"},
		{"\f", "\\f", "form feed"},
		{"\r", "\\r", "carriage return"},
	}

	for _, tt := range tests {
		result := JsonEncode(types.NewString(tt.input))
		if !strings.Contains(result.ToString(), tt.contains) {
			t.Errorf("JsonEncode(%s) should contain %s, got %v", tt.desc, tt.contains, result.ToString())
		}
	}
}

// ============================================================================
// JSON Error Handling Tests
// ============================================================================

func TestJsonLastErrorAfterSuccess(t *testing.T) {
	// Perform successful encode
	JsonEncode(types.NewInt(42))

	errCode := JsonLastError()
	if errCode.ToInt() != JSON_ERROR_NONE {
		t.Errorf("After successful encode, error should be JSON_ERROR_NONE, got %d", errCode.ToInt())
	}

	errMsg := JsonLastErrorMsg()
	if errMsg.ToString() != "No error" {
		t.Errorf("Error message should be 'No error', got %v", errMsg.ToString())
	}
}

func TestJsonLastErrorMessages(t *testing.T) {
	// We can't directly set the error code, but we can test the messages
	// by testing JsonLastErrorMsg after decoding invalid JSON
	JsonDecode(types.NewString("{invalid}"))

	// At least verify we get an error
	errCode := JsonLastError()
	if errCode.ToInt() == JSON_ERROR_NONE {
		t.Error("Should have error after invalid JSON")
	}

	// Test that error message is not empty
	errMsg := JsonLastErrorMsg()
	if errMsg.ToString() == "" || errMsg.ToString() == "No error" {
		t.Error("Should have error message after invalid JSON")
	}
}
