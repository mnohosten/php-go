package util

import (
	"fmt"
	"math"
)

// SafeAddSize safely adds two sizes, checking for overflow.
// Returns an error if the addition would overflow or if either value is negative.
func SafeAddSize(a, b int) (int, error) {
	// Check for negative values
	if a < 0 || b < 0 {
		return 0, fmt.Errorf("negative size not allowed: %d, %d", a, b)
	}

	// Check if addition would overflow
	// If a > 0 and b > 0, then a + b overflows if a > MaxInt - b
	if a > 0 && b > 0 && a > (math.MaxInt-b) {
		return 0, fmt.Errorf("size overflow: %d + %d exceeds maximum", a, b)
	}

	return a + b, nil
}

// SafeMultiplySize safely multiplies two sizes, checking for overflow.
// Returns an error if the multiplication would overflow or if either value is negative.
func SafeMultiplySize(a, b int) (int, error) {
	// Check for negative values
	if a < 0 || b < 0 {
		return 0, fmt.Errorf("negative size not allowed: %d, %d", a, b)
	}

	// Zero multiplication is always safe
	if a == 0 || b == 0 {
		return 0, nil
	}

	// Check if multiplication would overflow
	// If a > 0 and b > 0, then a * b overflows if a > MaxInt / b
	if a > (math.MaxInt / b) {
		return 0, fmt.Errorf("size overflow: %d * %d exceeds maximum", a, b)
	}

	return a * b, nil
}

// SafeConvertToInt converts int64 to int with overflow checking.
// Returns an error if the value is outside the range of int.
func SafeConvertToInt(val int64, name string) (int, error) {
	if val > int64(math.MaxInt) {
		return 0, fmt.Errorf("%s too large: %d exceeds maximum int", name, val)
	}
	if val < int64(math.MinInt) {
		return 0, fmt.Errorf("%s too small: %d below minimum int", name, val)
	}
	return int(val), nil
}

// ValidateIntRange checks if an integer is within expected bounds.
// Returns an error if the value is outside the specified range.
func ValidateIntRange(val int64, min, max int64, name string) error {
	if val < min || val > max {
		return fmt.Errorf("%s out of range: %d not in [%d, %d]", name, val, min, max)
	}
	return nil
}
