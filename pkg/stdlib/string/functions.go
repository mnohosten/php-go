package string

import (
	"strings"

	"github.com/krizos/php-go/pkg/types"
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

	if n < 0 {
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
