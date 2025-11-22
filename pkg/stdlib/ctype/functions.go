package ctype

import (
	"unicode"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Ctype Functions
// PHP's ctype functions check character types in strings
// ============================================================================

// CtypeAlnum checks for alphanumeric character(s)
// ctype_alnum(mixed $text): bool
func CtypeAlnum(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeAlpha checks for alphabetic character(s)
// ctype_alpha(mixed $text): bool
func CtypeAlpha(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsLetter(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeCntrl checks for control character(s)
// ctype_cntrl(mixed $text): bool
func CtypeCntrl(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsControl(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeDigit checks for numeric character(s)
// ctype_digit(mixed $text): bool
func CtypeDigit(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsDigit(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeGraph checks for any printable character(s) except space
// ctype_graph(mixed $text): bool
func CtypeGraph(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeLower checks for lowercase character(s)
// ctype_lower(mixed $text): bool
func CtypeLower(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	hasLetter := false
	for _, r := range str {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsLower(r) {
				return types.NewBool(false)
			}
		} else {
			// Non-letters make it false
			return types.NewBool(false)
		}
	}

	return types.NewBool(hasLetter)
}

// CtypePrint checks for printable character(s)
// ctype_print(mixed $text): bool
func CtypePrint(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsPrint(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypePunct checks for any printable character which is not whitespace or an alphanumeric character
// ctype_punct(mixed $text): bool
func CtypePunct(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeSpace checks for whitespace character(s)
// ctype_space(mixed $text): bool
func CtypeSpace(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !unicode.IsSpace(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// CtypeUpper checks for uppercase character(s)
// ctype_upper(mixed $text): bool
func CtypeUpper(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	hasLetter := false
	for _, r := range str {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return types.NewBool(false)
			}
		} else {
			// Non-letters make it false
			return types.NewBool(false)
		}
	}

	return types.NewBool(hasLetter)
}

// CtypeXdigit checks for hexadecimal digit character(s)
// ctype_xdigit(mixed $text): bool
func CtypeXdigit(text *types.Value) *types.Value {
	str := text.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	for _, r := range str {
		if !isHexDigit(r) {
			return types.NewBool(false)
		}
	}

	return types.NewBool(true)
}

// isHexDigit checks if a rune is a hexadecimal digit (0-9, a-f, A-F)
func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
