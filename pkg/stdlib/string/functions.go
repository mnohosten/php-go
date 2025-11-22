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

// Htmlspecialchars converts special characters to HTML entities
// htmlspecialchars(string $string, int $flags = ENT_COMPAT | ENT_HTML401): string
func Htmlspecialchars(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Basic entity encoding (simplified)
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#039;")

	return types.NewString(s)
}

// Htmlentities converts all applicable characters to HTML entities
// htmlentities(string $string, int $flags = ENT_COMPAT | ENT_HTML401): string
func Htmlentities(str *types.Value, flags ...*types.Value) *types.Value {
	// For simplified implementation, htmlentities behaves like htmlspecialchars
	// In full implementation, would encode more characters
	return Htmlspecialchars(str, flags...)
}

// HtmlspecialcharsDecode converts special HTML entities back to characters
// htmlspecialchars_decode(string $string, int $flags = ENT_COMPAT | ENT_HTML401): string
func HtmlspecialcharsDecode(str *types.Value, flags ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#039;", "'")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&amp;", "&")

	return types.NewString(s)
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

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '\'' || ch == '"' || ch == '\\' || ch == 0 {
			result.WriteByte('\\')
		}
		result.WriteByte(ch)
	}

	return types.NewString(result.String())
}

// Stripslashes removes backslashes from a string
// stripslashes(string $string): string
func Stripslashes(str *types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()
	var result strings.Builder

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++ // Skip the backslash
		}
		result.WriteByte(s[i])
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
