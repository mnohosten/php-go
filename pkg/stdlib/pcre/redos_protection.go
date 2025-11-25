package pcre

import (
	"context"
	"regexp"
	"time"
)

// ============================================================================
// ReDoS (Regular Expression Denial of Service) Protection
// ============================================================================

// DefaultRegexTimeout is the default timeout for regex operations (100ms)
// This prevents catastrophic backtracking in malicious or poorly written regex patterns
const DefaultRegexTimeout = 100 * time.Millisecond

// MaxRegexTimeout is the maximum allowed timeout (10 seconds)
// This prevents users from setting excessively long timeouts
const MaxRegexTimeout = 10 * time.Second

// regexTimeout stores the current timeout setting
var regexTimeout = DefaultRegexTimeout

// regexTimeoutEnabled controls whether timeout protection is active
var regexTimeoutEnabled = true

// SetRegexTimeout sets the timeout for regex operations
// If timeout > MaxRegexTimeout, it will be capped at MaxRegexTimeout
func SetRegexTimeout(timeout time.Duration) {
	if timeout > MaxRegexTimeout {
		regexTimeout = MaxRegexTimeout
	} else if timeout > 0 {
		regexTimeout = timeout
	} else {
		regexTimeout = DefaultRegexTimeout
	}
}

// GetRegexTimeout returns the current regex timeout setting
func GetRegexTimeout() time.Duration {
	return regexTimeout
}

// EnableRegexTimeout enables regex timeout protection
func EnableRegexTimeout() {
	regexTimeoutEnabled = true
}

// DisableRegexTimeout disables regex timeout protection
// WARNING: Only use this in trusted environments. Disabling timeout protection
// makes the system vulnerable to ReDoS attacks.
func DisableRegexTimeout() {
	regexTimeoutEnabled = false
}

// IsRegexTimeoutEnabled returns whether regex timeout protection is enabled
func IsRegexTimeoutEnabled() bool {
	return regexTimeoutEnabled
}

// ============================================================================
// Timeout-Protected Regex Operations
// ============================================================================

// matchWithTimeout performs regex matching with timeout protection
func matchWithTimeout(re *regexp.Regexp, subject string, operation func() interface{}) (interface{}, bool) {
	if !regexTimeoutEnabled {
		// Timeout disabled, run directly
		return operation(), true
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), regexTimeout)
	defer cancel()

	// Result channel
	resultChan := make(chan interface{}, 1)

	// Run regex operation in goroutine
	go func() {
		defer func() {
			// Recover from any panics in regex operation
			if r := recover(); r != nil {
				resultChan <- nil
			}
		}()
		resultChan <- operation()
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result, true
	case <-ctx.Done():
		// Timeout occurred
		return nil, false
	}
}

// FindStringSubmatchWithTimeout performs FindStringSubmatch with timeout protection
func FindStringSubmatchWithTimeout(re *regexp.Regexp, subject string) ([]string, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.FindStringSubmatch(subject)
	})

	if !ok {
		return nil, false
	}

	if result == nil {
		return nil, true
	}

	return result.([]string), true
}

// FindAllStringSubmatchWithTimeout performs FindAllStringSubmatch with timeout protection
func FindAllStringSubmatchWithTimeout(re *regexp.Regexp, subject string, n int) ([][]string, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.FindAllStringSubmatch(subject, n)
	})

	if !ok {
		return nil, false
	}

	if result == nil {
		return nil, true
	}

	return result.([][]string), true
}

// ReplaceAllStringWithTimeout performs ReplaceAllString with timeout protection
func ReplaceAllStringWithTimeout(re *regexp.Regexp, subject, replacement string) (string, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.ReplaceAllString(subject, replacement)
	})

	if !ok {
		return "", false
	}

	return result.(string), true
}

// ReplaceAllStringFuncWithTimeout performs ReplaceAllStringFunc with timeout protection
func ReplaceAllStringFuncWithTimeout(re *regexp.Regexp, subject string, repl func(string) string) (string, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.ReplaceAllStringFunc(subject, repl)
	})

	if !ok {
		return "", false
	}

	return result.(string), true
}

// SplitWithTimeout performs Split with timeout protection
func SplitWithTimeout(re *regexp.Regexp, subject string, n int) ([]string, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.Split(subject, n)
	})

	if !ok {
		return nil, false
	}

	if result == nil {
		return nil, true
	}

	return result.([]string), true
}

// MatchStringWithTimeout performs MatchString with timeout protection
func MatchStringWithTimeout(re *regexp.Regexp, subject string) (bool, bool) {
	result, ok := matchWithTimeout(re, subject, func() interface{} {
		return re.MatchString(subject)
	})

	if !ok {
		return false, false
	}

	return result.(bool), true
}
