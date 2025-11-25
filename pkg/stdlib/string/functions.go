package string

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Security - String Length Limits
// ============================================================================

// String length limits to prevent memory exhaustion attacks
const (
	// MaxStringLength is the maximum allowed length for string operations
	// Set to 100MB to prevent DoS attacks via excessive memory allocation
	// This limit applies to operations like str_repeat, str_pad, etc.
	MaxStringLength = 1024 * 1024 * 100 // 100MB

	// MaxRepeatCount is the maximum number of repetitions for str_repeat
	// Even with small strings, huge repeat counts can cause issues
	MaxRepeatCount = 1000000 // 1 million repetitions
)

// ============================================================================
// String Length and Substring
// ============================================================================

// Strlen returns the length of a string in bytes
// strlen(string $string): int
func Strlen(str *types.Value) *types.Value {
	if str == nil {
		return types.NewInt(0)
	}
	s := str.ToString()
	return types.NewInt(int64(len(s)))
}

// Substr returns a portion of a string
// substr(string $string, int $offset, ?int $length = null): string
func Substr(str *types.Value, offset *types.Value, length ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	start := int(offset.ToInt())
	strLen := len(s)

	// Handle negative offset (count from end)
	if start < 0 {
		start = strLen + start
		if start < 0 {
			start = 0
		}
	}

	// Start beyond string length
	if start >= strLen {
		return types.NewString("")
	}

	// Calculate end position
	end := strLen
	if len(length) > 0 && length[0] != nil {
		lengthInt := int(length[0].ToInt())

		if lengthInt < 0 {
			// Negative length means "all except last N"
			end = strLen + lengthInt
			if end < start {
				return types.NewString("")
			}
		} else {
			end = start + lengthInt
			if end > strLen {
				end = strLen
			}
		}
	}

	return types.NewString(s[start:end])
}

// ============================================================================
// String Searching
// ============================================================================

// Strpos finds the position of first occurrence of a substring
// strpos(string $haystack, string $needle, int $offset = 0): int|false
func Strpos(haystack *types.Value, needle *types.Value, offset ...*types.Value) *types.Value {
	h := haystack.ToString()
	n := needle.ToString()

	// Empty needle
	if n == "" {
		return types.NewBool(false)
	}

	start := 0
	if len(offset) > 0 && offset[0] != nil {
		start = int(offset[0].ToInt())
		if start < 0 || start >= len(h) {
			return types.NewBool(false)
		}
	}

	// Search from offset
	index := strings.Index(h[start:], n)
	if index == -1 {
		return types.NewBool(false)
	}

	return types.NewInt(int64(start + index))
}

// Strrpos finds the position of last occurrence of a substring
// strrpos(string $haystack, string $needle, int $offset = 0): int|false
func Strrpos(haystack *types.Value, needle *types.Value, offset ...*types.Value) *types.Value {
	h := haystack.ToString()
	n := needle.ToString()

	// Empty needle
	if n == "" {
		return types.NewBool(false)
	}

	start := 0
	if len(offset) > 0 && offset[0] != nil {
		start = int(offset[0].ToInt())
		if start < 0 {
			// Negative offset means search before that position from end
			start = len(h) + start
		}
		if start < 0 || start >= len(h) {
			return types.NewBool(false)
		}
	}

	// Search for last occurrence
	index := strings.LastIndex(h[start:], n)
	if index == -1 {
		return types.NewBool(false)
	}

	return types.NewInt(int64(start + index))
}

// Stripos finds position of first occurrence (case-insensitive)
// stripos(string $haystack, string $needle, int $offset = 0): int|false
func Stripos(haystack *types.Value, needle *types.Value, offset ...*types.Value) *types.Value {
	h := strings.ToLower(haystack.ToString())
	n := strings.ToLower(needle.ToString())

	if n == "" {
		return types.NewBool(false)
	}

	start := 0
	if len(offset) > 0 && offset[0] != nil {
		start = int(offset[0].ToInt())
		if start < 0 || start >= len(h) {
			return types.NewBool(false)
		}
	}

	index := strings.Index(h[start:], n)
	if index == -1 {
		return types.NewBool(false)
	}

	return types.NewInt(int64(start + index))
}

// Strripos finds position of last occurrence (case-insensitive)
// strripos(string $haystack, string $needle, int $offset = 0): int|false
func Strripos(haystack *types.Value, needle *types.Value, offset ...*types.Value) *types.Value {
	h := strings.ToLower(haystack.ToString())
	n := strings.ToLower(needle.ToString())

	if n == "" {
		return types.NewBool(false)
	}

	start := 0
	if len(offset) > 0 && offset[0] != nil {
		start = int(offset[0].ToInt())
		if start < 0 {
			start = len(h) + start
		}
		if start < 0 || start >= len(h) {
			return types.NewBool(false)
		}
	}

	index := strings.LastIndex(h[start:], n)
	if index == -1 {
		return types.NewBool(false)
	}

	return types.NewInt(int64(start + index))
}

// ============================================================================
// String Replacement
// ============================================================================

// StrReplace replaces all occurrences of search with replace
// str_replace(mixed $search, mixed $replace, mixed $subject): string|array
func StrReplace(search *types.Value, replace *types.Value, subject *types.Value) *types.Value {
	// For simplicity, handle string-to-string replacement
	// PHP supports arrays for all three parameters, but we'll implement the basic case
	s := search.ToString()
	r := replace.ToString()
	subj := subject.ToString()

	// PHP returns the original string if search is empty
	if s == "" {
		return types.NewString(subj)
	}

	// Security: Estimate result size to prevent excessive memory allocation
	// Count occurrences and calculate estimated result length
	occurrences := strings.Count(subj, s)
	if occurrences > 0 {
		// Calculate: original length - (search length * occurrences) + (replace length * occurrences)
		estimatedLength := int64(len(subj)) - int64(len(s)*occurrences) + int64(len(r)*occurrences)
		if estimatedLength > MaxStringLength {
			// Result would exceed maximum length
			return types.NewBool(false)
		}
	}

	result := strings.ReplaceAll(subj, s, r)
	return types.NewString(result)
}

// StrIreplace replaces all occurrences (case-insensitive)
// str_ireplace(mixed $search, mixed $replace, mixed $subject): string|array
func StrIreplace(search *types.Value, replace *types.Value, subject *types.Value) *types.Value {
	s := search.ToString()
	r := replace.ToString()
	subj := subject.ToString()

	// Case-insensitive replacement
	// We'll use a simple approach: find and replace manually
	lowerSubj := strings.ToLower(subj)
	lowerSearch := strings.ToLower(s)

	result := ""
	lastIdx := 0

	for {
		idx := strings.Index(lowerSubj[lastIdx:], lowerSearch)
		if idx == -1 {
			result += subj[lastIdx:]
			break
		}

		realIdx := lastIdx + idx
		result += subj[lastIdx:realIdx]
		result += r
		lastIdx = realIdx + len(s)
	}

	return types.NewString(result)
}

// ============================================================================
// Case Conversion
// ============================================================================

// Strtolower converts string to lowercase
// strtolower(string $string): string
func Strtolower(str *types.Value) *types.Value {
	s := str.ToString()
	return types.NewString(strings.ToLower(s))
}

// Strtoupper converts string to uppercase
// strtoupper(string $string): string
func Strtoupper(str *types.Value) *types.Value {
	s := str.ToString()
	return types.NewString(strings.ToUpper(s))
}

// Ucfirst makes the first character uppercase
// ucfirst(string $string): string
func Ucfirst(str *types.Value) *types.Value {
	s := str.ToString()
	if len(s) == 0 {
		return types.NewString("")
	}

	return types.NewString(strings.ToUpper(s[:1]) + s[1:])
}

// Lcfirst makes the first character lowercase
// lcfirst(string $string): string
func Lcfirst(str *types.Value) *types.Value {
	s := str.ToString()
	if len(s) == 0 {
		return types.NewString("")
	}

	return types.NewString(strings.ToLower(s[:1]) + s[1:])
}

// Ucwords makes the first character of each word uppercase
// ucwords(string $string): string
func Ucwords(str *types.Value) *types.Value {
	s := str.ToString()
	return types.NewString(strings.Title(s))
}

// ============================================================================
// Trimming
// ============================================================================

// Trim strips whitespace from beginning and end
// trim(string $string, string $characters = " \t\n\r\0\x0B"): string
func Trim(str *types.Value, characters ...*types.Value) *types.Value {
	s := str.ToString()

	if len(characters) > 0 && characters[0] != nil {
		cutset := characters[0].ToString()
		return types.NewString(strings.Trim(s, cutset))
	}

	return types.NewString(strings.TrimSpace(s))
}

// Ltrim strips whitespace from beginning
// ltrim(string $string, string $characters = " \t\n\r\0\x0B"): string
func Ltrim(str *types.Value, characters ...*types.Value) *types.Value {
	s := str.ToString()

	if len(characters) > 0 && characters[0] != nil {
		cutset := characters[0].ToString()
		return types.NewString(strings.TrimLeft(s, cutset))
	}

	return types.NewString(strings.TrimLeft(s, " \t\n\r\x00\x0B"))
}

// Rtrim strips whitespace from end
// rtrim(string $string, string $characters = " \t\n\r\0\x0B"): string
func Rtrim(str *types.Value, characters ...*types.Value) *types.Value {
	s := str.ToString()

	if len(characters) > 0 && characters[0] != nil {
		cutset := characters[0].ToString()
		return types.NewString(strings.TrimRight(s, cutset))
	}

	return types.NewString(strings.TrimRight(s, " \t\n\r\x00\x0B"))
}

// ============================================================================
// Explode/Implode
// ============================================================================

// Explode splits a string by delimiter
// explode(string $delimiter, string $string, int $limit = PHP_INT_MAX): array
func Explode(delimiter *types.Value, str *types.Value, limit ...*types.Value) *types.Value {
	delim := delimiter.ToString()
	s := str.ToString()

	// Empty delimiter is not allowed in PHP
	if delim == "" {
		return types.NewBool(false)
	}

	var parts []string
	if len(limit) > 0 && limit[0] != nil {
		limitInt := int(limit[0].ToInt())
		if limitInt == 1 {
			parts = []string{s}
		} else if limitInt > 1 {
			parts = strings.SplitN(s, delim, limitInt)
		} else {
			// Negative limit: return all except last |limit| elements
			allParts := strings.Split(s, delim)
			if limitInt < 0 && len(allParts)+limitInt > 0 {
				parts = allParts[:len(allParts)+limitInt]
			} else {
				parts = strings.Split(s, delim)
			}
		}
	} else {
		parts = strings.Split(s, delim)
	}

	// Convert to PHP array
	arr := types.NewEmptyArray()
	for _, part := range parts {
		arr.Append(types.NewString(part))
	}

	return types.NewArray(arr)
}

// Implode joins array elements with a string
// implode(string $separator, array $array): string
// Also: implode(array $array): string (with empty separator)
func Implode(separator *types.Value, arr ...*types.Value) *types.Value {
	var sep string
	var array *types.Array

	// Handle both signatures: implode(sep, arr) and implode(arr)
	if len(arr) == 0 {
		// implode($array) - separator is actually the array
		if separator.Type() != types.TypeArray {
			return types.NewString("")
		}
		sep = ""
		array = separator.ToArray()
	} else {
		// implode($sep, $array)
		sep = separator.ToString()
		if arr[0] == nil || arr[0].Type() != types.TypeArray {
			return types.NewString("")
		}
		array = arr[0].ToArray()
	}

	// Collect strings
	var parts []string
	array.Each(func(_, value *types.Value) bool {
		parts = append(parts, value.ToString())
		return true
	})

	return types.NewString(strings.Join(parts, sep))
}

// Join is an alias for Implode
func Join(separator *types.Value, arr ...*types.Value) *types.Value {
	return Implode(separator, arr...)
}

// ============================================================================
// Additional String Functions
// ============================================================================

// StrSplit converts a string to an array
// str_split(string $string, int $length = 1): array
func StrSplit(str *types.Value, length ...*types.Value) *types.Value {
	s := str.ToString()

	chunkLen := 1
	if len(length) > 0 && length[0] != nil {
		chunkLen = int(length[0].ToInt())
		if chunkLen < 1 {
			return types.NewBool(false)
		}
	}

	arr := types.NewEmptyArray()
	for i := 0; i < len(s); i += chunkLen {
		end := i + chunkLen
		if end > len(s) {
			end = len(s)
		}
		arr.Append(types.NewString(s[i:end]))
	}

	return types.NewArray(arr)
}

// ChunkSplit splits a string into chunks
// chunk_split(string $string, int $length = 76, string $end = "\r\n"): string
func ChunkSplit(str *types.Value, length *types.Value, end ...*types.Value) *types.Value {
	s := str.ToString()
	chunkLen := int(length.ToInt())

	if chunkLen < 1 {
		return types.NewBool(false)
	}

	ending := "\r\n"
	if len(end) > 0 && end[0] != nil {
		ending = end[0].ToString()
	}

	result := ""
	for i := 0; i < len(s); i += chunkLen {
		e := i + chunkLen
		if e > len(s) {
			e = len(s)
		}
		result += s[i:e] + ending
	}

	return types.NewString(result)
}

// StrRepeat repeats a string
// str_repeat(string $string, int $times): string
func StrRepeat(str *types.Value, times *types.Value) *types.Value {
	s := str.ToString()
	n := int(times.ToInt())

	// Validate repeat count
	if n < 0 {
		return types.NewBool(false)
	}

	// Security: Prevent excessive memory allocation
	// Check if repeat count exceeds maximum
	if n > MaxRepeatCount {
		// PHP behavior: trigger warning and return false
		return types.NewBool(false)
	}

	// Security: Check if resulting string would exceed maximum length
	// Use int64 to prevent integer overflow in multiplication
	resultLength := int64(len(s)) * int64(n)
	if resultLength > MaxStringLength {
		// PHP behavior: trigger warning and return false
		return types.NewBool(false)
	}

	return types.NewString(strings.Repeat(s, n))
}

// StrPad pads a string to a certain length
// str_pad(string $string, int $length, string $pad_string = " ", int $pad_type = STR_PAD_RIGHT): string
func StrPad(str *types.Value, length *types.Value, padString *types.Value, padType ...*types.Value) *types.Value {
	s := str.ToString()
	targetLen := int(length.ToInt())
	pad := " "

	if padString != nil {
		pad = padString.ToString()
	}

	// Security: Validate target length to prevent excessive memory allocation
	if targetLen < 0 || int64(targetLen) > MaxStringLength {
		// PHP behavior: return false for invalid length
		return types.NewBool(false)
	}

	if pad == "" || len(s) >= targetLen {
		return types.NewString(s)
	}

	padLen := targetLen - len(s)

	// Simplified padding (right pad only for now)
	// TODO: implement pad type (left, right, both) when padType parameter is provided
	padding := strings.Repeat(pad, (padLen/len(pad))+1)[:padLen]
	return types.NewString(s + padding)
}

// StrRev reverses a string
// strrev(string $string): string
func StrRev(str *types.Value) *types.Value {
	s := str.ToString()
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return types.NewString(string(runes))
}

// StrShuffle randomly shuffles a string
// str_shuffle(string $string): string
func StrShuffle(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	if len(s) <= 1 {
		return types.NewString(s)
	}

	runes := []rune(s)

	// Fisher-Yates shuffle using a simple time-based seed
	// For deterministic testing, this uses the string length as part of the seed
	seed := uint64(len(runes)) * 1103515245
	for i := len(runes) - 1; i > 0; i-- {
		seed = seed*1103515245 + 12345
		j := int(seed % uint64(i+1))
		runes[i], runes[j] = runes[j], runes[i]
	}

	return types.NewString(string(runes))
}

// StrWordCount returns information about words used in a string
// str_word_count(string $string, int $format = 0, ?string $characters = null): array|int
func StrWordCount(str *types.Value, format ...*types.Value) *types.Value {
	if str == nil {
		return types.NewInt(0)
	}

	s := str.ToString()

	// Default format is 0 (return count)
	formatVal := 0
	if len(format) > 0 && format[0] != nil {
		formatVal = int(format[0].ToInt())
	}

	// Additional characters that should be considered part of a word
	additionalChars := ""
	if len(format) > 1 && format[1] != nil {
		additionalChars = format[1].ToString()
	}

	// Helper to check if a character is a word character
	isWordChar := func(r rune) bool {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '\'' || r == '-' {
			return true
		}
		// Check additional characters
		for _, c := range additionalChars {
			if r == c {
				return true
			}
		}
		return false
	}

	words := make([]string, 0)
	positions := make([]int, 0)
	wordStart := -1

	for i, r := range s {
		if isWordChar(r) {
			if wordStart == -1 {
				wordStart = i
			}
		} else {
			if wordStart != -1 {
				words = append(words, s[wordStart:i])
				positions = append(positions, wordStart)
				wordStart = -1
			}
		}
	}
	// Handle last word
	if wordStart != -1 {
		words = append(words, s[wordStart:])
		positions = append(positions, wordStart)
	}

	switch formatVal {
	case 0:
		// Return the number of words
		return types.NewInt(int64(len(words)))
	case 1:
		// Return an array of words (indexed numerically)
		arr := types.NewEmptyArray()
		for _, word := range words {
			arr.Append(types.NewString(word))
		}
		return types.NewArray(arr)
	case 2:
		// Return an associative array with position as key
		arr := types.NewEmptyArray()
		for i, word := range words {
			arr.Set(types.NewInt(int64(positions[i])), types.NewString(word))
		}
		return types.NewArray(arr)
	default:
		return types.NewInt(int64(len(words)))
	}
}

// SubstrCount counts the number of substring occurrences
// substr_count(string $haystack, string $needle, int $offset = 0, ?int $length = null): int
func SubstrCount(haystack *types.Value, needle *types.Value, offsetAndLength ...*types.Value) *types.Value {
	if haystack == nil || needle == nil {
		return types.NewInt(0)
	}

	h := haystack.ToString()
	n := needle.ToString()

	if n == "" {
		// Empty needle - PHP returns warning and 0
		return types.NewInt(0)
	}

	// Handle offset
	offset := 0
	if len(offsetAndLength) > 0 && offsetAndLength[0] != nil {
		offset = int(offsetAndLength[0].ToInt())
	}

	// Handle negative offset
	if offset < 0 {
		offset = len(h) + offset
		if offset < 0 {
			offset = 0
		}
	}

	// Clamp offset
	if offset >= len(h) {
		return types.NewInt(0)
	}

	// Handle length
	searchStr := h[offset:]
	if len(offsetAndLength) > 1 && offsetAndLength[1] != nil {
		length := int(offsetAndLength[1].ToInt())
		if length < 0 {
			// Negative length - count from end
			length = len(searchStr) + length
		}
		if length > 0 && length < len(searchStr) {
			searchStr = searchStr[:length]
		}
	}

	// Count occurrences
	count := 0
	for i := 0; i <= len(searchStr)-len(n); {
		idx := strings.Index(searchStr[i:], n)
		if idx == -1 {
			break
		}
		count++
		i += idx + len(n)
	}

	return types.NewInt(int64(count))
}

// SubstrReplace replaces part of a string
// substr_replace(array|string $string, array|string $replace, array|int $offset, array|int|null $length = null): string|array
func SubstrReplace(str *types.Value, replacement *types.Value, offset *types.Value, length ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	r := ""
	if replacement != nil {
		r = replacement.ToString()
	}

	start := 0
	if offset != nil {
		start = int(offset.ToInt())
	}

	// Handle negative offset
	if start < 0 {
		start = len(s) + start
		if start < 0 {
			start = 0
		}
	}

	// Clamp start to string length
	if start > len(s) {
		start = len(s)
	}

	// Default length is rest of string
	end := len(s)
	if len(length) > 0 && length[0] != nil {
		l := int(length[0].ToInt())
		if l < 0 {
			// Negative length means leave that many characters at the end
			end = len(s) + l
			if end < start {
				end = start
			}
		} else {
			end = start + l
			if end > len(s) {
				end = len(s)
			}
		}
	}

	// Build result: before + replacement + after
	result := s[:start] + r + s[end:]
	return types.NewString(result)
}

// Strstr finds the first occurrence of a string (returns substring from match)
// strstr(string $haystack, mixed $needle, bool $before_needle = false): string|false
func Strstr(haystack *types.Value, needle *types.Value, beforeNeedle ...*types.Value) *types.Value {
	h := haystack.ToString()
	n := needle.ToString()

	index := strings.Index(h, n)
	if index == -1 {
		return types.NewBool(false)
	}

	before := false
	if len(beforeNeedle) > 0 && beforeNeedle[0] != nil {
		before = beforeNeedle[0].ToBool()
	}

	if before {
		return types.NewString(h[:index])
	}

	return types.NewString(h[index:])
}

// Strchr is an alias for Strstr
func Strchr(haystack *types.Value, needle *types.Value, beforeNeedle ...*types.Value) *types.Value {
	return Strstr(haystack, needle, beforeNeedle...)
}

// ============================================================================
// String Formatting Functions
// ============================================================================

// Sprintf returns a formatted string
// sprintf(string $format, mixed ...$values): string
func Sprintf(format *types.Value, values ...*types.Value) *types.Value {
	if format == nil {
		return types.NewString("")
	}

	f := format.ToString()
	result := formatString(f, values)
	return types.NewString(result)
}

// Printf outputs a formatted string
// printf(string $format, mixed ...$values): int
func Printf(format *types.Value, values ...*types.Value) *types.Value {
	result := Sprintf(format, values...)
	output := result.ToString()
	// In a real implementation, this would write to stdout
	// For now, we just return the length
	return types.NewInt(int64(len(output)))
}

// Vsprintf returns a formatted string from an array of arguments
// vsprintf(string $format, array $values): string
func Vsprintf(format *types.Value, values *types.Value) *types.Value {
	if format == nil {
		return types.NewString("")
	}

	// Convert array to slice of values
	var args []*types.Value
	if values != nil && values.Type() == types.TypeArray {
		arr := values.ToArray()
		arr.Each(func(key *types.Value, val *types.Value) bool {
			args = append(args, val)
			return true
		})
	}

	return Sprintf(format, args...)
}

// Vprintf outputs a formatted string from an array of arguments
// vprintf(string $format, array $values): int
func Vprintf(format *types.Value, values *types.Value) *types.Value {
	result := Vsprintf(format, values)
	return types.NewInt(int64(len(result.ToString())))
}

// NumberFormat formats a number with grouped thousands
// number_format(float $num, int $decimals = 0, ?string $decimal_separator = ".", ?string $thousands_separator = ","): string
func NumberFormat(num *types.Value, optionalArgs ...*types.Value) *types.Value {
	if num == nil {
		return types.NewString("0")
	}

	number := num.ToFloat()
	decimals := 0
	decimalSep := "."
	thousandsSep := ","

	if len(optionalArgs) > 0 && optionalArgs[0] != nil {
		decimals = int(optionalArgs[0].ToInt())
		if decimals < 0 {
			decimals = 0
		}
	}
	if len(optionalArgs) > 1 && optionalArgs[1] != nil {
		decimalSep = optionalArgs[1].ToString()
	}
	if len(optionalArgs) > 2 && optionalArgs[2] != nil {
		thousandsSep = optionalArgs[2].ToString()
	}

	// Handle negative numbers
	negative := number < 0
	if negative {
		number = -number
	}

	// Round to specified decimal places
	multiplier := 1.0
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	rounded := int64(number*multiplier + 0.5)

	// Split into integer and decimal parts
	intPart := rounded
	decPart := int64(0)
	if decimals > 0 {
		divisor := int64(multiplier)
		intPart = rounded / divisor
		decPart = rounded % divisor
	}

	// Format integer part with thousands separators
	intStr := formatWithThousands(intPart, thousandsSep)

	// Build result
	var result strings.Builder
	if negative {
		result.WriteByte('-')
	}
	result.WriteString(intStr)

	// Add decimal part if needed
	if decimals > 0 {
		result.WriteString(decimalSep)
		decStr := fmt.Sprintf("%0*d", decimals, decPart)
		result.WriteString(decStr)
	}

	return types.NewString(result.String())
}

// formatWithThousands adds thousand separators to an integer
func formatWithThousands(n int64, sep string) string {
	if sep == "" {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	start := len(str) % 3
	if start == 0 {
		start = 3
	}

	result.WriteString(str[:start])
	for i := start; i < len(str); i += 3 {
		result.WriteString(sep)
		result.WriteString(str[i : i+3])
	}

	return result.String()
}

// Sscanf parses input from a string according to a format
// sscanf(string $string, string $format, mixed &...$vars): array|int|null
func Sscanf(str *types.Value, format *types.Value, vars ...*types.Value) *types.Value {
	if str == nil || format == nil {
		return types.NewNull()
	}

	s := str.ToString()
	f := format.ToString()

	// Parse format and extract values
	values := parseSscanf(s, f)

	// If no vars provided, return array of values
	if len(vars) == 0 {
		arr := types.NewEmptyArray()
		for _, v := range values {
			arr.Append(v)
		}
		return types.NewArray(arr)
	}

	// Otherwise return count of assigned values
	return types.NewInt(int64(len(values)))
}

// parseSscanf extracts values from string based on format
func parseSscanf(str, format string) []*types.Value {
	var values []*types.Value
	strIdx := 0
	formatIdx := 0

	for formatIdx < len(format) && strIdx < len(str) {
		if format[formatIdx] != '%' {
			// Literal character - must match
			if format[formatIdx] == str[strIdx] {
				formatIdx++
				strIdx++
				continue
			}
			break
		}

		// Handle %%
		if formatIdx+1 < len(format) && format[formatIdx+1] == '%' {
			if str[strIdx] == '%' {
				formatIdx += 2
				strIdx++
				continue
			}
			break
		}

		formatIdx++
		if formatIdx >= len(format) {
			break
		}

		// Skip width specifier
		for formatIdx < len(format) && format[formatIdx] >= '0' && format[formatIdx] <= '9' {
			formatIdx++
		}

		if formatIdx >= len(format) {
			break
		}

		// Handle format specifier
		switch format[formatIdx] {
		case 'd', 'i': // Integer
			// Read integer
			start := strIdx
			if strIdx < len(str) && (str[strIdx] == '-' || str[strIdx] == '+') {
				strIdx++
			}
			for strIdx < len(str) && str[strIdx] >= '0' && str[strIdx] <= '9' {
				strIdx++
			}
			if strIdx > start {
				var val int64
				fmt.Sscanf(str[start:strIdx], "%d", &val)
				values = append(values, types.NewInt(val))
			}
		case 'f', 'e', 'g': // Float
			start := strIdx
			if strIdx < len(str) && (str[strIdx] == '-' || str[strIdx] == '+') {
				strIdx++
			}
			for strIdx < len(str) && ((str[strIdx] >= '0' && str[strIdx] <= '9') || str[strIdx] == '.') {
				strIdx++
			}
			// Handle exponent
			if strIdx < len(str) && (str[strIdx] == 'e' || str[strIdx] == 'E') {
				strIdx++
				if strIdx < len(str) && (str[strIdx] == '-' || str[strIdx] == '+') {
					strIdx++
				}
				for strIdx < len(str) && str[strIdx] >= '0' && str[strIdx] <= '9' {
					strIdx++
				}
			}
			if strIdx > start {
				var val float64
				fmt.Sscanf(str[start:strIdx], "%f", &val)
				values = append(values, types.NewFloat(val))
			}
		case 's': // String (non-whitespace)
			start := strIdx
			for strIdx < len(str) && str[strIdx] != ' ' && str[strIdx] != '\t' && str[strIdx] != '\n' && str[strIdx] != '\r' {
				strIdx++
			}
			if strIdx > start {
				values = append(values, types.NewString(str[start:strIdx]))
			}
		case 'c': // Character
			if strIdx < len(str) {
				values = append(values, types.NewString(string(str[strIdx])))
				strIdx++
			}
		case 'x', 'X': // Hex
			start := strIdx
			for strIdx < len(str) && ((str[strIdx] >= '0' && str[strIdx] <= '9') ||
				(str[strIdx] >= 'a' && str[strIdx] <= 'f') ||
				(str[strIdx] >= 'A' && str[strIdx] <= 'F')) {
				strIdx++
			}
			if strIdx > start {
				var val int64
				fmt.Sscanf(str[start:strIdx], "%x", &val)
				values = append(values, types.NewInt(val))
			}
		}
		formatIdx++
	}

	return values
}

// StrGetcsv parses a CSV string into an array
// str_getcsv(string $string, string $separator = ",", string $enclosure = "\"", string $escape = "\\"): array
func StrGetcsv(str *types.Value, optionalArgs ...*types.Value) *types.Value {
	if str == nil {
		return types.NewArray(types.NewEmptyArray())
	}

	s := str.ToString()
	separator := ","
	enclosure := "\""
	escape := "\\"

	if len(optionalArgs) > 0 && optionalArgs[0] != nil {
		sep := optionalArgs[0].ToString()
		if len(sep) > 0 {
			separator = string(sep[0])
		}
	}
	if len(optionalArgs) > 1 && optionalArgs[1] != nil {
		enc := optionalArgs[1].ToString()
		if len(enc) > 0 {
			enclosure = string(enc[0])
		}
	}
	if len(optionalArgs) > 2 && optionalArgs[2] != nil {
		esc := optionalArgs[2].ToString()
		if len(esc) > 0 {
			escape = string(esc[0])
		}
	}

	// Parse CSV
	fields := parseCSVLine(s, separator[0], enclosure[0], escape[0])

	// Convert to array
	arr := types.NewEmptyArray()
	for _, field := range fields {
		arr.Append(types.NewString(field))
	}

	return types.NewArray(arr)
}

// parseCSVLine parses a single CSV line
func parseCSVLine(line string, sep, enc, esc byte) []string {
	var fields []string
	var field strings.Builder
	inQuotes := false
	i := 0

	for i < len(line) {
		ch := line[i]

		if inQuotes {
			if ch == esc && i+1 < len(line) && line[i+1] == enc {
				// Escaped enclosure
				field.WriteByte(enc)
				i += 2
				continue
			}
			if ch == enc {
				// End of quoted field
				inQuotes = false
				i++
				continue
			}
			field.WriteByte(ch)
			i++
		} else {
			if ch == enc {
				// Start of quoted field
				inQuotes = true
				i++
				continue
			}
			if ch == sep {
				// End of field
				fields = append(fields, field.String())
				field.Reset()
				i++
				continue
			}
			field.WriteByte(ch)
			i++
		}
	}

	// Add last field
	fields = append(fields, field.String())

	return fields
}

// formatString implements basic sprintf-style formatting
func formatString(format string, values []*types.Value) string {
	var result strings.Builder
	valueIdx := 0

	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			result.WriteByte(format[i])
			continue
		}

		// Handle %%
		if i+1 < len(format) && format[i+1] == '%' {
			result.WriteByte('%')
			i++
			continue
		}

		// No more values
		if valueIdx >= len(values) {
			result.WriteByte('%')
			continue
		}

		// Parse format specifier
		i++
		if i >= len(format) {
			break
		}

		// Skip padding/width for now (simplified implementation)
		for i < len(format) && (format[i] == '-' || format[i] == '+' || format[i] == ' ' || format[i] == '0' || (format[i] >= '0' && format[i] <= '9') || format[i] == '.') {
			i++
		}

		if i >= len(format) {
			break
		}

		// Handle format type
		value := values[valueIdx]
		valueIdx++

		switch format[i] {
		case 's': // String
			result.WriteString(value.ToString())
		case 'd', 'i': // Integer
			result.WriteString(value.ToString())
		case 'f', 'F': // Float
			result.WriteString(value.ToString())
		case 'x': // Hex lowercase
			result.WriteString(value.ToString())
		case 'X': // Hex uppercase
			result.WriteString(strings.ToUpper(value.ToString()))
		case 'c': // Character
			if value.Type() == types.TypeInt {
				result.WriteByte(byte(value.ToInt()))
			} else {
				s := value.ToString()
				if len(s) > 0 {
					result.WriteByte(s[0])
				}
			}
		default:
			result.WriteByte('%')
			result.WriteByte(format[i])
		}
	}

	return result.String()
}

// ============================================================================
// String Comparison Functions
// ============================================================================

// Strcmp performs binary safe string comparison
// strcmp(string $string1, string $string2): int
func Strcmp(str1 *types.Value, str2 *types.Value) *types.Value {
	s1 := str1.ToString()
	s2 := str2.ToString()

	if s1 == s2 {
		return types.NewInt(0)
	}
	if s1 < s2 {
		return types.NewInt(-1)
	}
	return types.NewInt(1)
}

// Strcasecmp performs case-insensitive string comparison
// strcasecmp(string $string1, string $string2): int
func Strcasecmp(str1 *types.Value, str2 *types.Value) *types.Value {
	s1 := strings.ToLower(str1.ToString())
	s2 := strings.ToLower(str2.ToString())

	if s1 == s2 {
		return types.NewInt(0)
	}
	if s1 < s2 {
		return types.NewInt(-1)
	}
	return types.NewInt(1)
}

// Strncmp performs binary safe string comparison of first n characters
// strncmp(string $string1, string $string2, int $length): int
func Strncmp(str1 *types.Value, str2 *types.Value, length *types.Value) *types.Value {
	s1 := str1.ToString()
	s2 := str2.ToString()
	n := int(length.ToInt())

	if n <= 0 {
		return types.NewInt(0)
	}

	if len(s1) > n {
		s1 = s1[:n]
	}
	if len(s2) > n {
		s2 = s2[:n]
	}

	if s1 == s2 {
		return types.NewInt(0)
	}
	if s1 < s2 {
		return types.NewInt(-1)
	}
	return types.NewInt(1)
}

// Strncasecmp performs case-insensitive string comparison of first n characters
// strncasecmp(string $string1, string $string2, int $length): int
func Strncasecmp(str1 *types.Value, str2 *types.Value, length *types.Value) *types.Value {
	s1 := strings.ToLower(str1.ToString())
	s2 := strings.ToLower(str2.ToString())
	n := int(length.ToInt())

	if n <= 0 {
		return types.NewInt(0)
	}

	if len(s1) > n {
		s1 = s1[:n]
	}
	if len(s2) > n {
		s2 = s2[:n]
	}

	if s1 == s2 {
		return types.NewInt(0)
	}
	if s1 < s2 {
		return types.NewInt(-1)
	}
	return types.NewInt(1)
}

// Stristr finds the first occurrence of a string (case-insensitive)
// stristr(string $haystack, mixed $needle, bool $before_needle = false): string|false
func Stristr(haystack *types.Value, needle *types.Value, beforeNeedle ...*types.Value) *types.Value {
	h := strings.ToLower(haystack.ToString())
	n := strings.ToLower(needle.ToString())
	hOrig := haystack.ToString()

	index := strings.Index(h, n)
	if index == -1 {
		return types.NewBool(false)
	}

	before := false
	if len(beforeNeedle) > 0 && beforeNeedle[0] != nil {
		before = beforeNeedle[0].ToBool()
	}

	if before {
		return types.NewString(hOrig[:index])
	}

	return types.NewString(hOrig[index:])
}

// Strrchr finds the last occurrence of a character in a string
// strrchr(string $haystack, mixed $needle): string|false
func Strrchr(haystack *types.Value, needle *types.Value) *types.Value {
	h := haystack.ToString()
	n := needle.ToString()

	if n == "" {
		return types.NewBool(false)
	}

	// Use first character of needle
	char := n[0]
	index := strings.LastIndexByte(h, char)

	if index == -1 {
		return types.NewBool(false)
	}

	return types.NewString(h[index:])
}

// ============================================================================
// HTML/Special Character Functions
// ============================================================================

// PHP htmlspecialchars flags
const (
	ENT_COMPAT     = 2   // Convert double-quotes, leave single-quotes
	ENT_QUOTES     = 3   // Convert both double and single quotes
	ENT_NOQUOTES   = 0   // Don't convert any quotes
	ENT_HTML401    = 0   // HTML 4.01 (default)
	ENT_XML1       = 16  // XML 1
	ENT_XHTML      = 32  // XHTML
	ENT_HTML5      = 48  // HTML 5
	ENT_IGNORE     = 4   // Silently discard invalid encoding
	ENT_SUBSTITUTE = 8   // Replace invalid encoding with Unicode replacement char
	ENT_DISALLOWED = 128 // Replace code points that are disallowed
)

// Htmlspecialchars converts special characters to HTML entities
// htmlspecialchars(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null, bool $double_encode = true): string
func Htmlspecialchars(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Default flags: ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	flagValue := ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	if len(flags) > 0 && flags[0] != nil {
		flagValue = int(flags[0].ToInt())
	}

	// Determine quote handling
	// ENT_NOQUOTES = 0: no quotes
	// ENT_COMPAT = 2: double quotes only
	// ENT_QUOTES = 3: both quotes
	quoteBits := flagValue & ENT_QUOTES
	doubleQuotes := quoteBits == ENT_COMPAT || quoteBits == ENT_QUOTES
	singleQuotes := quoteBits == ENT_QUOTES

	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '&':
			result.WriteString("&amp;")
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '"':
			if doubleQuotes {
				result.WriteString("&quot;")
			} else {
				result.WriteByte(ch)
			}
		case '\'':
			if singleQuotes {
				result.WriteString("&#039;")
			} else {
				result.WriteByte(ch)
			}
		default:
			result.WriteByte(ch)
		}
	}

	return types.NewString(result.String())
}

// htmlEntities is a map of HTML entity names to their character values
var htmlEntities = map[string]string{
	// Basic entities (from htmlspecialchars)
	"amp": "&", "lt": "<", "gt": ">", "quot": "\"", "apos": "'",
	// Latin-1 supplement
	"nbsp": "\u00A0", "iexcl": "¡", "cent": "¢", "pound": "£", "curren": "¤",
	"yen": "¥", "brvbar": "¦", "sect": "§", "uml": "¨", "copy": "©",
	"ordf": "ª", "laquo": "«", "not": "¬", "shy": "\u00AD", "reg": "®",
	"macr": "¯", "deg": "°", "plusmn": "±", "sup2": "²", "sup3": "³",
	"acute": "´", "micro": "µ", "para": "¶", "middot": "·", "cedil": "¸",
	"sup1": "¹", "ordm": "º", "raquo": "»", "frac14": "¼", "frac12": "½",
	"frac34": "¾", "iquest": "¿",
	// Latin-1 letters
	"Agrave": "À", "Aacute": "Á", "Acirc": "Â", "Atilde": "Ã", "Auml": "Ä",
	"Aring": "Å", "AElig": "Æ", "Ccedil": "Ç", "Egrave": "È", "Eacute": "É",
	"Ecirc": "Ê", "Euml": "Ë", "Igrave": "Ì", "Iacute": "Í", "Icirc": "Î",
	"Iuml": "Ï", "ETH": "Ð", "Ntilde": "Ñ", "Ograve": "Ò", "Oacute": "Ó",
	"Ocirc": "Ô", "Otilde": "Õ", "Ouml": "Ö", "times": "×", "Oslash": "Ø",
	"Ugrave": "Ù", "Uacute": "Ú", "Ucirc": "Û", "Uuml": "Ü", "Yacute": "Ý",
	"THORN": "Þ", "szlig": "ß", "agrave": "à", "aacute": "á", "acirc": "â",
	"atilde": "ã", "auml": "ä", "aring": "å", "aelig": "æ", "ccedil": "ç",
	"egrave": "è", "eacute": "é", "ecirc": "ê", "euml": "ë", "igrave": "ì",
	"iacute": "í", "icirc": "î", "iuml": "ï", "eth": "ð", "ntilde": "ñ",
	"ograve": "ò", "oacute": "ó", "ocirc": "ô", "otilde": "õ", "ouml": "ö",
	"divide": "÷", "oslash": "ø", "ugrave": "ù", "uacute": "ú", "ucirc": "û",
	"uuml": "ü", "yacute": "ý", "thorn": "þ", "yuml": "ÿ",
	// Greek letters
	"Alpha": "Α", "Beta": "Β", "Gamma": "Γ", "Delta": "Δ", "Epsilon": "Ε",
	"Zeta": "Ζ", "Eta": "Η", "Theta": "Θ", "Iota": "Ι", "Kappa": "Κ",
	"Lambda": "Λ", "Mu": "Μ", "Nu": "Ν", "Xi": "Ξ", "Omicron": "Ο",
	"Pi": "Π", "Rho": "Ρ", "Sigma": "Σ", "Tau": "Τ", "Upsilon": "Υ",
	"Phi": "Φ", "Chi": "Χ", "Psi": "Ψ", "Omega": "Ω",
	"alpha": "α", "beta": "β", "gamma": "γ", "delta": "δ", "epsilon": "ε",
	"zeta": "ζ", "eta": "η", "theta": "θ", "iota": "ι", "kappa": "κ",
	"lambda": "λ", "mu": "μ", "nu": "ν", "xi": "ξ", "omicron": "ο",
	"pi": "π", "rho": "ρ", "sigmaf": "ς", "sigma": "σ", "tau": "τ",
	"upsilon": "υ", "phi": "φ", "chi": "χ", "psi": "ψ", "omega": "ω",
	// Math symbols
	"forall": "∀", "part": "∂", "exist": "∃", "empty": "∅", "nabla": "∇",
	"isin": "∈", "notin": "∉", "ni": "∋", "prod": "∏", "sum": "∑",
	"minus": "−", "lowast": "∗", "radic": "√", "prop": "∝", "infin": "∞",
	"ang": "∠", "and": "∧", "or": "∨", "cap": "∩", "cup": "∪",
	"int": "∫", "there4": "∴", "sim": "∼", "cong": "≅", "asymp": "≈",
	"ne": "≠", "equiv": "≡", "le": "≤", "ge": "≥", "sub": "⊂",
	"sup": "⊃", "nsub": "⊄", "sube": "⊆", "supe": "⊇", "oplus": "⊕",
	"otimes": "⊗", "perp": "⊥", "sdot": "⋅",
	// Punctuation
	"bull": "•", "hellip": "…", "prime": "′", "Prime": "″", "oline": "‾",
	"frasl": "⁄", "euro": "€", "trade": "™", "larr": "←", "uarr": "↑",
	"rarr": "→", "darr": "↓", "harr": "↔", "crarr": "↵",
	// Miscellaneous
	"loz": "◊", "spades": "♠", "clubs": "♣", "hearts": "♥", "diams": "♦",
}

// reverseHtmlEntities maps characters to their entity names (for htmlentities)
var reverseHtmlEntities = map[rune]string{
	'&': "amp", '<': "lt", '>': "gt", '"': "quot", '\'': "apos",
	'\u00A0': "nbsp", '¡': "iexcl", '¢': "cent", '£': "pound", '¤': "curren",
	'¥': "yen", '¦': "brvbar", '§': "sect", '¨': "uml", '©': "copy",
	'ª': "ordf", '«': "laquo", '¬': "not", '\u00AD': "shy", '®': "reg",
	'¯': "macr", '°': "deg", '±': "plusmn", '²': "sup2", '³': "sup3",
	'´': "acute", 'µ': "micro", '¶': "para", '·': "middot", '¸': "cedil",
	'¹': "sup1", 'º': "ordm", '»': "raquo", '¼': "frac14", '½': "frac12",
	'¾': "frac34", '¿': "iquest",
	'À': "Agrave", 'Á': "Aacute", 'Â': "Acirc", 'Ã': "Atilde", 'Ä': "Auml",
	'Å': "Aring", 'Æ': "AElig", 'Ç': "Ccedil", 'È': "Egrave", 'É': "Eacute",
	'Ê': "Ecirc", 'Ë': "Euml", 'Ì': "Igrave", 'Í': "Iacute", 'Î': "Icirc",
	'Ï': "Iuml", 'Ð': "ETH", 'Ñ': "Ntilde", 'Ò': "Ograve", 'Ó': "Oacute",
	'Ô': "Ocirc", 'Õ': "Otilde", 'Ö': "Ouml", '×': "times", 'Ø': "Oslash",
	'Ù': "Ugrave", 'Ú': "Uacute", 'Û': "Ucirc", 'Ü': "Uuml", 'Ý': "Yacute",
	'Þ': "THORN", 'ß': "szlig", 'à': "agrave", 'á': "aacute", 'â': "acirc",
	'ã': "atilde", 'ä': "auml", 'å': "aring", 'æ': "aelig", 'ç': "ccedil",
	'è': "egrave", 'é': "eacute", 'ê': "ecirc", 'ë': "euml", 'ì': "igrave",
	'í': "iacute", 'î': "icirc", 'ï': "iuml", 'ð': "eth", 'ñ': "ntilde",
	'ò': "ograve", 'ó': "oacute", 'ô': "ocirc", 'õ': "otilde", 'ö': "ouml",
	'÷': "divide", 'ø': "oslash", 'ù': "ugrave", 'ú': "uacute", 'û': "ucirc",
	'ü': "uuml", 'ý': "yacute", 'þ': "thorn", 'ÿ': "yuml",
}

// Htmlentities converts all applicable characters to HTML entities
// htmlentities(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null, bool $double_encode = true): string
func Htmlentities(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Default flags: ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	flagValue := ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	if len(flags) > 0 && flags[0] != nil {
		flagValue = int(flags[0].ToInt())
	}

	// Determine quote handling
	// ENT_NOQUOTES = 0: no quotes
	// ENT_COMPAT = 2: double quotes only
	// ENT_QUOTES = 3: both quotes
	quoteBits := flagValue & ENT_QUOTES
	doubleQuotes := quoteBits == ENT_COMPAT || quoteBits == ENT_QUOTES
	singleQuotes := quoteBits == ENT_QUOTES

	var result strings.Builder
	result.Grow(len(s) * 2) // Estimate growth

	for _, r := range s {
		// Check for special handling of quotes
		if r == '"' && !doubleQuotes {
			result.WriteRune(r)
			continue
		}
		if r == '\'' && !singleQuotes {
			result.WriteRune(r)
			continue
		}

		// Check if we have a named entity for this character
		if entity, ok := reverseHtmlEntities[r]; ok {
			result.WriteByte('&')
			result.WriteString(entity)
			result.WriteByte(';')
		} else if r > 127 {
			// For other high-value Unicode characters, use numeric entity
			result.WriteString("&#")
			// Convert rune to decimal string
			num := int(r)
			digits := make([]byte, 0, 10)
			for num > 0 {
				digits = append(digits, byte('0'+num%10))
				num /= 10
			}
			// Reverse and write
			for i := len(digits) - 1; i >= 0; i-- {
				result.WriteByte(digits[i])
			}
			result.WriteByte(';')
		} else {
			result.WriteRune(r)
		}
	}

	return types.NewString(result.String())
}

// HtmlspecialcharsDecode converts special HTML entities back to characters
// htmlspecialchars_decode(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401): string
func HtmlspecialcharsDecode(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Default flags: ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	flagValue := ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	if len(flags) > 0 && flags[0] != nil {
		flagValue = int(flags[0].ToInt())
	}

	// Determine quote handling
	// ENT_NOQUOTES = 0: no quotes
	// ENT_COMPAT = 2: double quotes only
	// ENT_QUOTES = 3: both quotes
	quoteBits := flagValue & ENT_QUOTES
	doubleQuotes := quoteBits == ENT_COMPAT || quoteBits == ENT_QUOTES
	singleQuotes := quoteBits == ENT_QUOTES

	// Replace entities in order (& must be last to avoid double-decoding)
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")

	if doubleQuotes {
		s = strings.ReplaceAll(s, "&quot;", "\"")
	}
	if singleQuotes {
		s = strings.ReplaceAll(s, "&#039;", "'")
		s = strings.ReplaceAll(s, "&#39;", "'")
		s = strings.ReplaceAll(s, "&apos;", "'")
	}

	// &amp; must be last
	s = strings.ReplaceAll(s, "&amp;", "&")

	return types.NewString(s)
}

// HtmlEntityDecode converts HTML entities to their corresponding characters
// html_entity_decode(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null): string
func HtmlEntityDecode(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Default flags: ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	flagValue := ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401
	if len(flags) > 0 && flags[0] != nil {
		flagValue = int(flags[0].ToInt())
	}

	// Determine quote handling
	// ENT_NOQUOTES = 0: no quotes
	// ENT_COMPAT = 2: double quotes only
	// ENT_QUOTES = 3: both quotes
	quoteBits := flagValue & ENT_QUOTES
	doubleQuotes := quoteBits == ENT_COMPAT || quoteBits == ENT_QUOTES
	singleQuotes := quoteBits == ENT_QUOTES

	var result strings.Builder
	result.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] == '&' {
			// Look for entity end
			end := i + 1
			for end < len(s) && end < i+32 && s[end] != ';' && s[end] != '&' {
				end++
			}

			if end < len(s) && s[end] == ';' {
				entityName := s[i+1 : end]

				// Check for numeric entity
				if len(entityName) > 0 && entityName[0] == '#' {
					var codePoint int
					if len(entityName) > 1 && (entityName[1] == 'x' || entityName[1] == 'X') {
						// Hexadecimal: &#xHHHH;
						for j := 2; j < len(entityName); j++ {
							digit := unhexChar(entityName[j])
							if digit < 0 {
								break
							}
							codePoint = codePoint*16 + digit
						}
					} else {
						// Decimal: &#NNNN;
						for j := 1; j < len(entityName); j++ {
							if entityName[j] >= '0' && entityName[j] <= '9' {
								codePoint = codePoint*10 + int(entityName[j]-'0')
							}
						}
					}
					if codePoint > 0 && codePoint <= 0x10FFFF {
						result.WriteRune(rune(codePoint))
						i = end + 1
						continue
					}
				}

				// Check for named entity
				if decoded, ok := htmlEntities[entityName]; ok {
					// Check quote handling
					if entityName == "quot" && !doubleQuotes {
						result.WriteString("&quot;")
					} else if (entityName == "apos" || entityName == "39") && !singleQuotes {
						result.WriteString("&" + entityName + ";")
					} else {
						result.WriteString(decoded)
					}
					i = end + 1
					continue
				}
			}
		}
		result.WriteByte(s[i])
		i++
	}

	return types.NewString(result.String())
}

// unhexChar converts a hex character to its value
func unhexChar(ch byte) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	if ch >= 'a' && ch <= 'f' {
		return int(ch - 'a' + 10)
	}
	if ch >= 'A' && ch <= 'F' {
		return int(ch - 'A' + 10)
	}
	return -1
}

// StripTags removes HTML and PHP tags from a string
// strip_tags(string $string, array|string|null $allowed_tags = null): string
func StripTags(str *types.Value, allowedTags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Parse allowed tags if provided
	var allowed map[string]bool
	if len(allowedTags) > 0 && allowedTags[0] != nil && allowedTags[0].Type() != types.TypeNull {
		allowed = make(map[string]bool)
		allowStr := allowedTags[0].ToString()
		// Parse format like "<p><br><div>" or "p br div"
		i := 0
		for i < len(allowStr) {
			// Skip to start of tag
			for i < len(allowStr) && allowStr[i] != '<' && !isAlpha(allowStr[i]) {
				i++
			}
			if i >= len(allowStr) {
				break
			}

			start := i
			if allowStr[i] == '<' {
				start++
				i++
			}

			// Read tag name
			tagStart := i
			for i < len(allowStr) && (isAlnum(allowStr[i]) || allowStr[i] == '-' || allowStr[i] == '_') {
				i++
			}
			if i > tagStart {
				tagName := strings.ToLower(allowStr[tagStart:i])
				allowed[tagName] = true
			}

			// Skip to end of tag
			for i < len(allowStr) && allowStr[i] != '>' && allowStr[i] != ' ' {
				i++
			}
			if i < len(allowStr) && allowStr[i] == '>' {
				i++
			}
			_ = start // suppress unused
		}
	}

	var result strings.Builder
	result.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] == '<' {
			// Find the tag name
			start := i
			i++

			// Skip whitespace after <
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}

			// Check for closing tag
			isClosing := false
			if i < len(s) && s[i] == '/' {
				isClosing = true
				i++
			}

			// Check for PHP tag or comment
			if i < len(s) && s[i] == '?' {
				// PHP tag - skip to ?>
				for i < len(s) {
					if s[i] == '?' && i+1 < len(s) && s[i+1] == '>' {
						i += 2
						break
					}
					i++
				}
				continue
			}

			if i+2 < len(s) && s[i] == '!' && s[i+1] == '-' && s[i+2] == '-' {
				// HTML comment - skip to -->
				i += 3
				for i < len(s) {
					if s[i] == '-' && i+2 < len(s) && s[i+1] == '-' && s[i+2] == '>' {
						i += 3
						break
					}
					i++
				}
				continue
			}

			// Read tag name
			tagStart := i
			for i < len(s) && (isAlnum(s[i]) || s[i] == '-' || s[i] == '_' || s[i] == ':') {
				i++
			}
			tagName := strings.ToLower(s[tagStart:i])

			// Find end of tag
			inQuote := byte(0)
			for i < len(s) {
				if inQuote != 0 {
					if s[i] == inQuote {
						inQuote = 0
					}
				} else {
					if s[i] == '"' || s[i] == '\'' {
						inQuote = s[i]
					} else if s[i] == '>' {
						break
					}
				}
				i++
			}

			if i < len(s) && s[i] == '>' {
				i++
			}

			// Check if this tag is allowed
			if allowed != nil && allowed[tagName] {
				result.WriteString(s[start:i])
			}
			_ = isClosing // suppress unused
		} else {
			result.WriteByte(s[i])
			i++
		}
	}

	return types.NewString(result.String())
}

// Helper functions for strip_tags
func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || (ch >= '0' && ch <= '9')
}

// ============================================================================
// Slashing Functions
// ============================================================================

// Addslashes adds backslashes before characters that need to be escaped
// addslashes(string $string): string
func Addslashes(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder
	result.Grow(len(s) * 2)

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '\'', '"', '\\':
			result.WriteByte('\\')
			result.WriteByte(ch)
		case 0: // NUL byte
			result.WriteByte('\\')
			result.WriteByte('0')
		default:
			result.WriteByte(ch)
		}
	}

	return types.NewString(result.String())
}

// Stripslashes removes backslash quoting from a string
// stripslashes(string $string): string
func Stripslashes(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder
	result.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			next := s[i+1]
			switch next {
			case '\'', '"', '\\':
				result.WriteByte(next)
				i += 2
				continue
			case '0':
				result.WriteByte(0) // NUL byte
				i += 2
				continue
			}
		}
		result.WriteByte(s[i])
		i++
	}

	return types.NewString(result.String())
}

// ============================================================================
// Text Formatting Functions
// ============================================================================

// Nl2br inserts HTML line breaks before all newlines
// nl2br(string $string, bool $use_xhtml = true): string
func Nl2br(str *types.Value, useXhtml ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	xhtml := true
	if len(useXhtml) > 0 && useXhtml[0] != nil {
		xhtml = useXhtml[0].ToBool()
	}

	br := "<br />"
	if !xhtml {
		br = "<br>"
	}

	// Replace \r\n first, then \n, then \r
	s = strings.ReplaceAll(s, "\r\n", br+"\r\n")
	s = strings.ReplaceAll(s, "\n", br+"\n")
	s = strings.ReplaceAll(s, "\r", br+"\r")

	return types.NewString(s)
}

// Wordwrap wraps a string to a given number of characters
// wordwrap(string $string, int $width = 75, string $break = "\n", bool $cut_long_words = false): string
func Wordwrap(str *types.Value, width *types.Value, breakStr ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	w := int(width.ToInt())
	if w <= 0 {
		w = 75
	}

	brk := "\n"
	if len(breakStr) > 0 && breakStr[0] != nil {
		brk = breakStr[0].ToString()
	}

	// Simplified implementation: break at word boundaries
	words := strings.Fields(s)
	if len(words) == 0 {
		return types.NewString(s)
	}

	var result strings.Builder
	lineLen := 0

	for _, word := range words {
		wordLen := len(word)

		if lineLen > 0 && lineLen+1+wordLen > w {
			result.WriteString(brk)
			lineLen = 0
		} else if lineLen > 0 {
			result.WriteByte(' ')
			lineLen++
		}

		result.WriteString(word)
		lineLen += wordLen
	}

	return types.NewString(result.String())
}

// ============================================================================
// URL Encoding Functions
// ============================================================================

// Urlencode encodes a URL string
// urlencode(string $string): string
func Urlencode(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '.' || ch == '~' {
			result.WriteByte(ch)
		} else if ch == ' ' {
			result.WriteByte('+')
		} else {
			result.WriteByte('%')
			result.WriteByte(hexDigit(ch >> 4))
			result.WriteByte(hexDigit(ch & 0xF))
		}
	}

	return types.NewString(result.String())
}

// Urldecode decodes a URL-encoded string
// urldecode(string $string): string
func Urldecode(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '+' {
			result.WriteByte(' ')
		} else if ch == '%' && i+2 < len(s) {
			h1 := unhex(s[i+1])
			h2 := unhex(s[i+2])
			if h1 >= 0 && h2 >= 0 {
				result.WriteByte(byte(h1<<4 | h2))
				i += 2
			} else {
				result.WriteByte(ch)
			}
		} else {
			result.WriteByte(ch)
		}
	}

	return types.NewString(result.String())
}

// Rawurlencode encodes a URL according to RFC 3986
// rawurlencode(string $string): string
func Rawurlencode(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '.' || ch == '~' {
			result.WriteByte(ch)
		} else {
			result.WriteByte('%')
			result.WriteByte(hexDigit(ch >> 4))
			result.WriteByte(hexDigit(ch & 0xF))
		}
	}

	return types.NewString(result.String())
}

// Rawurldecode decodes a URL-encoded string
// rawurldecode(string $string): string
func Rawurldecode(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '%' && i+2 < len(s) {
			h1 := unhex(s[i+1])
			h2 := unhex(s[i+2])
			if h1 >= 0 && h2 >= 0 {
				result.WriteByte(byte(h1<<4 | h2))
				i += 2
			} else {
				result.WriteByte(ch)
			}
		} else {
			result.WriteByte(ch)
		}
	}

	return types.NewString(result.String())
}

// ============================================================================
// Helper Functions
// ============================================================================

func hexDigit(val byte) byte {
	if val < 10 {
		return '0' + val
	}
	return 'A' + val - 10
}

func unhex(ch byte) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	if ch >= 'a' && ch <= 'f' {
		return int(ch - 'a' + 10)
	}
	if ch >= 'A' && ch <= 'F' {
		return int(ch - 'A' + 10)
	}
	return -1
}

// ============================================================================
// String Search Functions
// ============================================================================

// Strpbrk searches a string for any of a set of characters
// strpbrk(string $string, string $characters): string|false
func Strpbrk(str *types.Value, characters *types.Value) *types.Value {
	if str == nil || characters == nil {
		return types.NewBool(false)
	}

	s := str.ToString()
	chars := characters.ToString()

	if chars == "" {
		return types.NewBool(false)
	}

	// Find the first occurrence of any character from chars in s
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(chars, s[i]) != -1 {
			return types.NewString(s[i:])
		}
	}

	return types.NewBool(false)
}

// Strspn finds the length of the initial segment of a string consisting entirely of characters in a given mask
// strspn(string $string, string $characters, int $offset = 0, ?int $length = null): int
func Strspn(str *types.Value, characters *types.Value, offsetAndLength ...*types.Value) *types.Value {
	if str == nil || characters == nil {
		return types.NewInt(0)
	}

	s := str.ToString()
	chars := characters.ToString()

	// Handle offset
	offset := 0
	if len(offsetAndLength) > 0 && offsetAndLength[0] != nil {
		offset = int(offsetAndLength[0].ToInt())
	}

	// Handle negative offset
	if offset < 0 {
		offset = len(s) + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(s) {
		return types.NewInt(0)
	}

	// Handle length
	length := len(s) - offset
	if len(offsetAndLength) > 1 && offsetAndLength[1] != nil {
		length = int(offsetAndLength[1].ToInt())
		if length < 0 {
			// Negative length - count from end
			length = len(s) - offset + length
			if length < 0 {
				length = 0
			}
		}
	}

	// Limit to actual remaining string
	if offset+length > len(s) {
		length = len(s) - offset
	}

	// Count characters that are in the mask
	count := 0
	for i := offset; i < offset+length; i++ {
		if strings.IndexByte(chars, s[i]) != -1 {
			count++
		} else {
			break
		}
	}

	return types.NewInt(int64(count))
}

// Strcspn finds the length of the initial segment of a string consisting entirely of characters NOT in a given mask
// strcspn(string $string, string $characters, int $offset = 0, ?int $length = null): int
func Strcspn(str *types.Value, characters *types.Value, offsetAndLength ...*types.Value) *types.Value {
	if str == nil {
		return types.NewInt(0)
	}

	s := str.ToString()

	// Empty characters means count all characters until end
	chars := ""
	if characters != nil {
		chars = characters.ToString()
	}

	// Handle offset
	offset := 0
	if len(offsetAndLength) > 0 && offsetAndLength[0] != nil {
		offset = int(offsetAndLength[0].ToInt())
	}

	// Handle negative offset
	if offset < 0 {
		offset = len(s) + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(s) {
		return types.NewInt(0)
	}

	// Handle length
	length := len(s) - offset
	if len(offsetAndLength) > 1 && offsetAndLength[1] != nil {
		length = int(offsetAndLength[1].ToInt())
		if length < 0 {
			// Negative length - count from end
			length = len(s) - offset + length
			if length < 0 {
				length = 0
			}
		}
	}

	// Limit to actual remaining string
	if offset+length > len(s) {
		length = len(s) - offset
	}

	// If no characters to avoid, count up to length
	if chars == "" {
		return types.NewInt(int64(length))
	}

	// Count characters that are NOT in the mask
	count := 0
	for i := offset; i < offset+length; i++ {
		if strings.IndexByte(chars, s[i]) == -1 {
			count++
		} else {
			break
		}
	}

	return types.NewInt(int64(count))
}

// ============================================================================
// Base64 Encoding Functions
// ============================================================================

// Base64Encode encodes data with MIME base64
// base64_encode(string $string): string
func Base64Encode(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}
	data := str.ToString()
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return types.NewString(encoded)
}

// Base64Decode decodes data encoded with MIME base64
// base64_decode(string $string, bool $strict = false): string|false
func Base64Decode(str *types.Value, args ...*types.Value) *types.Value {
	if str == nil {
		return types.NewBool(false)
	}
	data := str.ToString()

	strict := false
	if len(args) > 0 && args[0] != nil {
		strict = args[0].ToBool()
	}

	var decoded []byte
	var err error

	if strict {
		// Strict mode - use StrictEncoding which requires padding
		decoded, err = base64.StdEncoding.Strict().DecodeString(data)
	} else {
		// Non-strict mode - be lenient with whitespace and padding
		decoded, err = base64.StdEncoding.DecodeString(data)
	}

	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(string(decoded))
}

// ============================================================================
// URL Parsing Functions
// ============================================================================

// PHP parse_url component constants
const (
	PHP_URL_SCHEME   = 0
	PHP_URL_HOST     = 1
	PHP_URL_PORT     = 2
	PHP_URL_USER     = 3
	PHP_URL_PASS     = 4
	PHP_URL_PATH     = 5
	PHP_URL_QUERY    = 6
	PHP_URL_FRAGMENT = 7
)

// ParseUrl parses a URL and returns its components
// parse_url(string $url, int $component = -1): array|string|int|null|false
func ParseUrl(urlStr *types.Value, args ...*types.Value) *types.Value {
	if urlStr == nil {
		return types.NewBool(false)
	}

	rawUrl := urlStr.ToString()
	if rawUrl == "" {
		return types.NewBool(false)
	}

	// Component to extract (-1 means all)
	component := -1
	if len(args) > 0 && args[0] != nil {
		component = int(args[0].ToInt())
	}

	// Parse the URL
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return types.NewBool(false)
	}

	// Helper to get specific component
	getComponent := func(comp int) *types.Value {
		switch comp {
		case PHP_URL_SCHEME:
			if parsed.Scheme != "" {
				return types.NewString(parsed.Scheme)
			}
			return types.NewNull()
		case PHP_URL_HOST:
			if parsed.Host != "" {
				host := parsed.Hostname()
				if host != "" {
					return types.NewString(host)
				}
			}
			return types.NewNull()
		case PHP_URL_PORT:
			port := parsed.Port()
			if port != "" {
				// Parse port as integer
				var portInt int
				fmt.Sscanf(port, "%d", &portInt)
				return types.NewInt(int64(portInt))
			}
			return types.NewNull()
		case PHP_URL_USER:
			if parsed.User != nil {
				return types.NewString(parsed.User.Username())
			}
			return types.NewNull()
		case PHP_URL_PASS:
			if parsed.User != nil {
				if pass, ok := parsed.User.Password(); ok {
					return types.NewString(pass)
				}
			}
			return types.NewNull()
		case PHP_URL_PATH:
			if parsed.Path != "" {
				return types.NewString(parsed.Path)
			}
			return types.NewNull()
		case PHP_URL_QUERY:
			if parsed.RawQuery != "" {
				return types.NewString(parsed.RawQuery)
			}
			return types.NewNull()
		case PHP_URL_FRAGMENT:
			if parsed.Fragment != "" {
				return types.NewString(parsed.Fragment)
			}
			return types.NewNull()
		default:
			return types.NewNull()
		}
	}

	// If specific component requested
	if component >= 0 {
		return getComponent(component)
	}

	// Return array with all components
	arr := types.NewEmptyArray()

	if parsed.Scheme != "" {
		arr.Set(types.NewString("scheme"), types.NewString(parsed.Scheme))
	}
	if parsed.Host != "" {
		host := parsed.Hostname()
		if host != "" {
			arr.Set(types.NewString("host"), types.NewString(host))
		}
		port := parsed.Port()
		if port != "" {
			var portInt int
			fmt.Sscanf(port, "%d", &portInt)
			arr.Set(types.NewString("port"), types.NewInt(int64(portInt)))
		}
	}
	if parsed.User != nil {
		arr.Set(types.NewString("user"), types.NewString(parsed.User.Username()))
		if pass, ok := parsed.User.Password(); ok {
			arr.Set(types.NewString("pass"), types.NewString(pass))
		}
	}
	if parsed.Path != "" {
		arr.Set(types.NewString("path"), types.NewString(parsed.Path))
	}
	if parsed.RawQuery != "" {
		arr.Set(types.NewString("query"), types.NewString(parsed.RawQuery))
	}
	if parsed.Fragment != "" {
		arr.Set(types.NewString("fragment"), types.NewString(parsed.Fragment))
	}

	return types.NewArray(arr)
}

// HttpBuildQuery generates a URL-encoded query string
// http_build_query(array|object $data, string $numeric_prefix = "", string $arg_separator = null, int $encoding_type = PHP_QUERY_RFC1738): string
func HttpBuildQuery(data *types.Value, args ...*types.Value) *types.Value {
	if data == nil || data.Type() != types.TypeArray {
		return types.NewString("")
	}

	arr := data.ToArray()

	numericPrefix := ""
	if len(args) > 0 && args[0] != nil {
		numericPrefix = args[0].ToString()
	}

	argSeparator := "&"
	if len(args) > 1 && args[1] != nil && args[1].Type() != types.TypeNull {
		argSeparator = args[1].ToString()
	}

	// Build query parts
	var parts []string

	arr.Each(func(key, value *types.Value) bool {
		keyStr := key.ToString()

		// Add numeric prefix for numeric keys
		if key.Type() == types.TypeInt {
			keyStr = numericPrefix + keyStr
		}

		// URL encode key and value
		encodedKey := url.QueryEscape(keyStr)
		encodedValue := url.QueryEscape(value.ToString())

		parts = append(parts, encodedKey+"="+encodedValue)
		return true
	})

	return types.NewString(strings.Join(parts, argSeparator))
}
