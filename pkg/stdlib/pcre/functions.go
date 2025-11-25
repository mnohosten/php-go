package pcre

import (
	"regexp"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Pattern Compilation and Caching
// ============================================================================

// patternCache stores compiled regular expressions
var patternCache = make(map[string]*regexp.Regexp)

// compilePattern compiles a PCRE pattern and caches it
func compilePattern(pattern string) (*regexp.Regexp, error) {
	// Check cache first
	if re, exists := patternCache[pattern]; exists {
		return re, nil
	}

	// Convert PCRE delimiters to Go regexp
	// PHP patterns are typically like /pattern/flags
	cleaned := pattern
	flags := ""

	// Extract pattern and flags from delimiters
	if len(pattern) > 2 {
		delimiter := pattern[0]
		lastDelim := strings.LastIndexByte(pattern, byte(delimiter))

		if lastDelim > 0 {
			cleaned = pattern[1:lastDelim]
			if lastDelim < len(pattern)-1 {
				flags = pattern[lastDelim+1:]
			}
		}
	}

	// Apply flags
	finalPattern := cleaned

	// Case-insensitive flag
	if strings.Contains(flags, "i") {
		finalPattern = "(?i)" + finalPattern
	}

	// Multiline flag
	if strings.Contains(flags, "m") {
		finalPattern = "(?m)" + finalPattern
	}

	// Dot matches newline flag
	if strings.Contains(flags, "s") {
		finalPattern = "(?s)" + finalPattern
	}

	// Compile the pattern
	re, err := regexp.Compile(finalPattern)
	if err != nil {
		return nil, err
	}

	// Cache it
	patternCache[pattern] = re
	return re, nil
}

// ============================================================================
// preg_match - Perform a regular expression match
// ============================================================================

// PregMatch performs a regular expression match
// preg_match(string $pattern, string $subject, array &$matches = null, int $flags = 0, int $offset = 0): int|false
func PregMatch(pattern *types.Value, subject *types.Value, matches ...*types.Value) *types.Value {
	if pattern == nil || subject == nil {
		return types.NewInt(0)
	}

	patternStr := pattern.ToString()
	subjectStr := subject.ToString()

	// Compile pattern
	re, err := compilePattern(patternStr)
	if err != nil {
		return types.NewBool(false)
	}

	// Find first match with timeout protection
	match, ok := FindStringSubmatchWithTimeout(re, subjectStr)
	if !ok {
		// Timeout occurred - potential ReDoS attack
		return types.NewBool(false)
	}
	if match == nil {
		// No match found
		if len(matches) > 0 && matches[0] != nil {
			// Set matches to empty array
			matches[0].ToArray().Reset()
		}
		return types.NewInt(0)
	}

	// Store matches if requested
	if len(matches) > 0 && matches[0] != nil {
		matchesArray := matches[0].ToArray()
		matchesArray.Reset()

		for i, m := range match {
			matchesArray.Set(types.NewInt(int64(i)), types.NewString(m))
		}
	}

	return types.NewInt(1)
}

// ============================================================================
// preg_match_all - Perform a global regular expression match
// ============================================================================

// PregMatchAll performs a global regular expression match
// preg_match_all(string $pattern, string $subject, array &$matches = null, int $flags = PREG_PATTERN_ORDER, int $offset = 0): int|false
func PregMatchAll(pattern *types.Value, subject *types.Value, matches ...*types.Value) *types.Value {
	if pattern == nil || subject == nil {
		return types.NewInt(0)
	}

	patternStr := pattern.ToString()
	subjectStr := subject.ToString()

	// Compile pattern
	re, err := compilePattern(patternStr)
	if err != nil {
		return types.NewBool(false)
	}

	// Find all matches with timeout protection
	allMatches, ok := FindAllStringSubmatchWithTimeout(re, subjectStr, -1)
	if !ok {
		// Timeout occurred - potential ReDoS attack
		return types.NewBool(false)
	}
	if allMatches == nil {
		// No matches found
		if len(matches) > 0 && matches[0] != nil {
			matches[0].ToArray().Reset()
		}
		return types.NewInt(0)
	}

	// Store matches if requested
	if len(matches) > 0 && matches[0] != nil {
		matchesArray := matches[0].ToArray()
		matchesArray.Reset()

		// PREG_PATTERN_ORDER format (default)
		// $matches[0] = array of full matches
		// $matches[1] = array of first capture group
		// etc.

		// Determine number of capture groups
		numGroups := 0
		if len(allMatches) > 0 {
			numGroups = len(allMatches[0])
		}

		// Initialize arrays for each capture group
		for i := 0; i < numGroups; i++ {
			groupArray := types.NewEmptyArray()
			matchesArray.Set(types.NewInt(int64(i)), types.NewArray(groupArray))
		}

		// Fill in the matches
		for _, match := range allMatches {
			for groupIdx, m := range match {
				groupArray, _ := matchesArray.Get(types.NewInt(int64(groupIdx)))
				groupArray.ToArray().Append(types.NewString(m))
			}
		}
	}

	return types.NewInt(int64(len(allMatches)))
}

// ============================================================================
// preg_replace - Perform a regular expression search and replace
// ============================================================================

// PregReplace performs a regular expression search and replace
// preg_replace(string|array $pattern, string|array $replacement, string|array $subject, int $limit = -1, int &$count = null): string|array|null
func PregReplace(pattern *types.Value, replacement *types.Value, subject *types.Value, limit ...*types.Value) *types.Value {
	if pattern == nil || replacement == nil || subject == nil {
		return types.NewNull()
	}

	patternStr := pattern.ToString()
	replacementStr := replacement.ToString()
	subjectStr := subject.ToString()

	// Get limit
	limitNum := -1
	if len(limit) > 0 && limit[0] != nil {
		limitNum = int(limit[0].ToInt())
	}

	// Compile pattern
	re, err := compilePattern(patternStr)
	if err != nil {
		return types.NewNull()
	}

	// Perform replacement with timeout protection
	var result string
	var ok bool
	if limitNum < 0 {
		// Replace all
		result, ok = ReplaceAllStringWithTimeout(re, subjectStr, replacementStr)
		if !ok {
			// Timeout occurred - potential ReDoS attack
			return types.NewNull()
		}
	} else if limitNum == 0 {
		// No replacements
		result = subjectStr
	} else {
		// Limited replacements
		count := 0
		result, ok = ReplaceAllStringFuncWithTimeout(re, subjectStr, func(match string) string {
			if count < limitNum {
				count++
				return re.ReplaceAllString(match, replacementStr)
			}
			return match
		})
		if !ok {
			// Timeout occurred - potential ReDoS attack
			return types.NewNull()
		}
	}

	return types.NewString(result)
}

// ============================================================================
// preg_split - Split string by a regular expression
// ============================================================================

// PregSplit splits a string by a regular expression
// preg_split(string $pattern, string $subject, int $limit = -1, int $flags = 0): array|false
func PregSplit(pattern *types.Value, subject *types.Value, limit ...*types.Value) *types.Value {
	if pattern == nil || subject == nil {
		return types.NewBool(false)
	}

	patternStr := pattern.ToString()
	subjectStr := subject.ToString()

	// Get limit
	limitNum := -1
	if len(limit) > 0 && limit[0] != nil {
		limitNum = int(limit[0].ToInt())
	}

	// Compile pattern
	re, err := compilePattern(patternStr)
	if err != nil {
		return types.NewBool(false)
	}

	// Split by pattern with timeout protection
	var parts []string
	var ok bool
	if limitNum < 0 {
		parts, ok = SplitWithTimeout(re, subjectStr, -1)
	} else {
		parts, ok = SplitWithTimeout(re, subjectStr, limitNum)
	}

	if !ok {
		// Timeout occurred - potential ReDoS attack
		return types.NewBool(false)
	}

	// Convert to array
	result := types.NewEmptyArray()
	for _, part := range parts {
		result.Append(types.NewString(part))
	}

	return types.NewArray(result)
}

// ============================================================================
// preg_grep - Return array entries that match the pattern
// ============================================================================

// PregGrep returns array entries that match the pattern
// preg_grep(string $pattern, array $array, int $flags = 0): array|false
func PregGrep(pattern *types.Value, input *types.Value, flags ...*types.Value) *types.Value {
	if pattern == nil || input == nil || input.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	patternStr := pattern.ToString()
	inputArray := input.ToArray()

	// Compile pattern
	re, err := compilePattern(patternStr)
	if err != nil {
		return types.NewBool(false)
	}

	// Check for PREG_GREP_INVERT flag
	invert := false
	if len(flags) > 0 && flags[0] != nil {
		// PREG_GREP_INVERT = 1
		invert = flags[0].ToInt() == 1
	}

	// Filter array with timeout protection
	result := types.NewEmptyArray()
	timedOut := false
	inputArray.Each(func(key, value *types.Value) bool {
		valueStr := value.ToString()
		matches, ok := MatchStringWithTimeout(re, valueStr)

		if !ok {
			// Timeout occurred - potential ReDoS attack
			timedOut = true
			return false // Stop iteration
		}

		// Include if matches (or doesn't match if inverted)
		if matches != invert {
			result.Set(key, value)
		}
		return true
	})

	if timedOut {
		return types.NewBool(false)
	}

	return types.NewArray(result)
}

// ============================================================================
// preg_quote - Quote regular expression characters
// ============================================================================

// PregQuote quotes regular expression characters
// preg_quote(string $string, ?string $delimiter = null): string
func PregQuote(str *types.Value, delimiter ...*types.Value) *types.Value {
	if str == nil {
		return types.NewString("")
	}

	s := str.ToString()

	// Characters that need to be escaped in regex
	specialChars := []string{
		"\\", "^", "$", ".", "[", "]", "|", "(", ")", "?", "*", "+", "{", "}",
	}

	result := s
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	// Quote delimiter if provided
	if len(delimiter) > 0 && delimiter[0] != nil {
		delim := delimiter[0].ToString()
		if delim != "" {
			result = strings.ReplaceAll(result, delim, "\\"+delim)
		}
	}

	return types.NewString(result)
}

// ============================================================================
// preg_last_error - Returns the error code of the last PCRE regex execution
// ============================================================================

// PregLastError returns the error code of the last PCRE regex execution
// preg_last_error(): int
func PregLastError() *types.Value {
	// In this simplified implementation, we always return PREG_NO_ERROR (0)
	// A full implementation would track actual errors
	return types.NewInt(0)
}

// ============================================================================
// Helper Functions
// ============================================================================

// ClearPatternCache clears the pattern compilation cache
func ClearPatternCache() {
	patternCache = make(map[string]*regexp.Regexp)
}
