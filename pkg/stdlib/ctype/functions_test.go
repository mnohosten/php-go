package ctype

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Ctype Alnum Tests
// ============================================================================

func TestCtypeAlnum(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC123", true},
		{"abcDEF123", true},
		{"abc", true},
		{"123", true},
		{"abc 123", false},     // space
		{"abc-123", false},     // hyphen
		{"abc_123", false},     // underscore
		{"", false},            // empty
		{"hello@world", false}, // special char
	}

	for _, tt := range tests {
		result := CtypeAlnum(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeAlnum(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Alpha Tests
// ============================================================================

func TestCtypeAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"ABC", true},
		{"abcDEF", true},
		{"abc123", false},  // has digits
		{"abc ", false},    // has space
		{"abc-def", false}, // has hyphen
		{"", false},        // empty
		{"123", false},     // only digits
	}

	for _, tt := range tests {
		result := CtypeAlpha(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeAlpha(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Cntrl Tests
// ============================================================================

func TestCtypeCntrl(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"\n", true},
		{"\t", true},
		{"\r", true},
		{"\n\t\r", true},
		{"abc", false},
		{"\nabc", false}, // mixed
		{"", false},
	}

	for _, tt := range tests {
		result := CtypeCntrl(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeCntrl(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Digit Tests
// ============================================================================

func TestCtypeDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"999", true},
		{"12345", true},
		{"abc", false},
		{"12a34", false},
		{"12 34", false}, // space
		{"", false},
		{"-123", false}, // negative sign
		{"12.34", false}, // decimal point
	}

	for _, tt := range tests {
		result := CtypeDigit(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeDigit(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Graph Tests
// ============================================================================

func TestCtypeGraph(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"ABC123", true},
		{"!@#$%", true},
		{"abc123!@#", true},
		{"abc ", false},   // has space
		{" abc", false},   // has space
		{"abc def", false}, // has space
		{"", false},
		{"\n", false},
		{"\t", false},
	}

	for _, tt := range tests {
		result := CtypeGraph(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeGraph(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Lower Tests
// ============================================================================

func TestCtypeLower(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"xyz", true},
		{"abcdef", true},
		{"ABC", false},
		{"Abc", false},
		{"abc123", false}, // has digits
		{"abc ", false},   // has space
		{"", false},
		{"123", false},
	}

	for _, tt := range tests {
		result := CtypeLower(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeLower(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Print Tests
// ============================================================================

func TestCtypePrint(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"ABC 123", true},
		{"hello world", true},
		{"!@#$%", true},
		{"\n", false},
		{"\t", false},
		{"abc\n", false},
		{"", false},
	}

	for _, tt := range tests {
		result := CtypePrint(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypePrint(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Punct Tests
// ============================================================================

func TestCtypePunct(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!@#$%", true},
		{".,;:", true},
		{"!?", true},
		{"abc", false},
		{"123", false},
		{"!a@", false},
		{"! @", false}, // space
		{"", false},
	}

	for _, tt := range tests {
		result := CtypePunct(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypePunct(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Space Tests
// ============================================================================

func TestCtypeSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{" ", true},
		{"  ", true},
		{"\t", true},
		{"\n", true},
		{"\r", true},
		{" \t\n\r", true},
		{"abc", false},
		{" abc", false},
		{"", false},
	}

	for _, tt := range tests {
		result := CtypeSpace(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeSpace(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Upper Tests
// ============================================================================

func TestCtypeUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"ABC", true},
		{"XYZ", true},
		{"ABCDEF", true},
		{"abc", false},
		{"Abc", false},
		{"ABC123", false}, // has digits
		{"ABC ", false},   // has space
		{"", false},
		{"123", false},
	}

	for _, tt := range tests {
		result := CtypeUpper(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeUpper(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Ctype Xdigit Tests
// ============================================================================

func TestCtypeXdigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0123456789", true},
		{"abcdef", true},
		{"ABCDEF", true},
		{"0123456789abcdefABCDEF", true},
		{"abc123", true},
		{"ABC123", true},
		{"0xff", false}, // 'x' is not a hex digit
		{"g", false},
		{"xyz", false},
		{"12g34", false},
		{"", false},
		{" ", false},
	}

	for _, tt := range tests {
		result := CtypeXdigit(types.NewString(tt.input))
		if result.ToBool() != tt.expected {
			t.Errorf("CtypeXdigit(%q) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

func TestCtypeEmptyString(t *testing.T) {
	empty := types.NewString("")

	// All ctype functions should return false for empty string
	funcs := []struct {
		name string
		fn   func(*types.Value) *types.Value
	}{
		{"CtypeAlnum", CtypeAlnum},
		{"CtypeAlpha", CtypeAlpha},
		{"CtypeCntrl", CtypeCntrl},
		{"CtypeDigit", CtypeDigit},
		{"CtypeGraph", CtypeGraph},
		{"CtypeLower", CtypeLower},
		{"CtypePrint", CtypePrint},
		{"CtypePunct", CtypePunct},
		{"CtypeSpace", CtypeSpace},
		{"CtypeUpper", CtypeUpper},
		{"CtypeXdigit", CtypeXdigit},
	}

	for _, tt := range funcs {
		result := tt.fn(empty)
		if result.ToBool() != false {
			t.Errorf("%s(\"\") should return false", tt.name)
		}
	}
}

func TestCtypeSingleCharacter(t *testing.T) {
	tests := []struct {
		char  string
		alnum bool
		alpha bool
		digit bool
		lower bool
		upper bool
	}{
		{"a", true, true, false, true, false},
		{"A", true, true, false, false, true},
		{"0", true, false, true, false, false},
		{"5", true, false, true, false, false},
		{"!", false, false, false, false, false},
		{" ", false, false, false, false, false},
	}

	for _, tt := range tests {
		input := types.NewString(tt.char)

		if CtypeAlnum(input).ToBool() != tt.alnum {
			t.Errorf("CtypeAlnum(%q) = %v, want %v", tt.char, !tt.alnum, tt.alnum)
		}
		if CtypeAlpha(input).ToBool() != tt.alpha {
			t.Errorf("CtypeAlpha(%q) = %v, want %v", tt.char, !tt.alpha, tt.alpha)
		}
		if CtypeDigit(input).ToBool() != tt.digit {
			t.Errorf("CtypeDigit(%q) = %v, want %v", tt.char, !tt.digit, tt.digit)
		}
		if CtypeLower(input).ToBool() != tt.lower {
			t.Errorf("CtypeLower(%q) = %v, want %v", tt.char, !tt.lower, tt.lower)
		}
		if CtypeUpper(input).ToBool() != tt.upper {
			t.Errorf("CtypeUpper(%q) = %v, want %v", tt.char, !tt.upper, tt.upper)
		}
	}
}

// ============================================================================
// Unicode Tests
// ============================================================================

func TestCtypeUnicode(t *testing.T) {
	// Test with unicode characters
	tests := []struct {
		input string
		alpha bool
		alnum bool
	}{
		{"café", true, true},
		{"hello世界", true, true},
		{"Москва", true, true},
		{"日本語", true, true},
	}

	for _, tt := range tests {
		input := types.NewString(tt.input)

		result := CtypeAlpha(input)
		if result.ToBool() != tt.alpha {
			t.Errorf("CtypeAlpha(%q) = %v, want %v", tt.input, result.ToBool(), tt.alpha)
		}

		result = CtypeAlnum(input)
		if result.ToBool() != tt.alnum {
			t.Errorf("CtypeAlnum(%q) = %v, want %v", tt.input, result.ToBool(), tt.alnum)
		}
	}
}
