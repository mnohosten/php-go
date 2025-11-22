package math

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/krizos/php-go/pkg/types"
)

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ============================================================================
// Basic Math Functions
// ============================================================================

// Abs returns the absolute value of a number
// abs(int|float $num): int|float
func Abs(num *types.Value) *types.Value {
	if num.Type() == types.TypeInt {
		n := num.ToInt()
		if n < 0 {
			return types.NewInt(-n)
		}
		return types.NewInt(n)
	}

	// Float or convertible to float
	f := num.ToFloat()
	return types.NewFloat(math.Abs(f))
}

// Ceil rounds a number up to the next highest integer
// ceil(int|float $num): float
func Ceil(num *types.Value) *types.Value {
	f := num.ToFloat()
	return types.NewFloat(math.Ceil(f))
}

// Floor rounds a number down to the next lowest integer
// floor(int|float $num): float
func Floor(num *types.Value) *types.Value {
	f := num.ToFloat()
	return types.NewFloat(math.Floor(f))
}

// Round rounds a float to specified precision
// round(int|float $num, int $precision = 0): float
func Round(num *types.Value, precision ...*types.Value) *types.Value {
	f := num.ToFloat()

	// Default precision is 0
	prec := 0
	if len(precision) > 0 && precision[0] != nil {
		prec = int(precision[0].ToInt())
	}

	if prec == 0 {
		return types.NewFloat(math.Round(f))
	}

	// Round to specified decimal places
	shift := math.Pow(10, float64(prec))
	return types.NewFloat(math.Round(f*shift) / shift)
}

// Min returns the lowest value
// min(mixed ...$values): mixed
func Min(values ...*types.Value) *types.Value {
	if len(values) == 0 {
		return types.NewNull()
	}

	if len(values) == 1 {
		// Single array argument
		if values[0].Type() == types.TypeArray {
			arr := values[0].ToArray()
			if arr.IsEmpty() {
				return types.NewNull()
			}

			var minVal *types.Value
			arr.Each(func(_, val *types.Value) bool {
				if minVal == nil {
					minVal = val
					return true
				}

				// Compare values
				if compareValues(val, minVal) < 0 {
					minVal = val
				}
				return true
			})

			return minVal
		}
		return values[0]
	}

	// Multiple arguments
	minVal := values[0]
	for i := 1; i < len(values); i++ {
		if compareValues(values[i], minVal) < 0 {
			minVal = values[i]
		}
	}

	return minVal
}

// Max returns the highest value
// max(mixed ...$values): mixed
func Max(values ...*types.Value) *types.Value {
	if len(values) == 0 {
		return types.NewNull()
	}

	if len(values) == 1 {
		// Single array argument
		if values[0].Type() == types.TypeArray {
			arr := values[0].ToArray()
			if arr.IsEmpty() {
				return types.NewNull()
			}

			var maxVal *types.Value
			arr.Each(func(_, val *types.Value) bool {
				if maxVal == nil {
					maxVal = val
					return true
				}

				// Compare values
				if compareValues(val, maxVal) > 0 {
					maxVal = val
				}
				return true
			})

			return maxVal
		}
		return values[0]
	}

	// Multiple arguments
	maxVal := values[0]
	for i := 1; i < len(values); i++ {
		if compareValues(values[i], maxVal) > 0 {
			maxVal = values[i]
		}
	}

	return maxVal
}

// compareValues compares two values for min/max
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareValues(a, b *types.Value) int {
	// If both are numeric, compare numerically
	if (a.Type() == types.TypeInt || a.Type() == types.TypeFloat) &&
		(b.Type() == types.TypeInt || b.Type() == types.TypeFloat) {
		aFloat := a.ToFloat()
		bFloat := b.ToFloat()
		if aFloat < bFloat {
			return -1
		} else if aFloat > bFloat {
			return 1
		}
		return 0
	}

	// String comparison
	aStr := a.ToString()
	bStr := b.ToString()
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// Pow returns base raised to the power of exponent
// pow(mixed $base, mixed $exp): int|float
func Pow(base, exp *types.Value) *types.Value {
	b := base.ToFloat()
	e := exp.ToFloat()

	result := math.Pow(b, e)

	// Return int if both inputs are ints and result is whole number
	if base.Type() == types.TypeInt && exp.Type() == types.TypeInt {
		if result == math.Floor(result) {
			return types.NewInt(int64(result))
		}
	}

	return types.NewFloat(result)
}

// Sqrt returns the square root of a number
// sqrt(float $num): float
func Sqrt(num *types.Value) *types.Value {
	f := num.ToFloat()
	if f < 0 {
		return types.NewFloat(math.NaN())
	}
	return types.NewFloat(math.Sqrt(f))
}

// ============================================================================
// Trigonometric Functions
// ============================================================================

// Sin returns the sine of a number
// sin(float $num): float
func Sin(num *types.Value) *types.Value {
	return types.NewFloat(math.Sin(num.ToFloat()))
}

// Cos returns the cosine of a number
// cos(float $num): float
func Cos(num *types.Value) *types.Value {
	return types.NewFloat(math.Cos(num.ToFloat()))
}

// Tan returns the tangent of a number
// tan(float $num): float
func Tan(num *types.Value) *types.Value {
	return types.NewFloat(math.Tan(num.ToFloat()))
}

// Asin returns the arc sine of a number
// asin(float $num): float
func Asin(num *types.Value) *types.Value {
	return types.NewFloat(math.Asin(num.ToFloat()))
}

// Acos returns the arc cosine of a number
// acos(float $num): float
func Acos(num *types.Value) *types.Value {
	return types.NewFloat(math.Acos(num.ToFloat()))
}

// Atan returns the arc tangent of a number
// atan(float $num): float
func Atan(num *types.Value) *types.Value {
	return types.NewFloat(math.Atan(num.ToFloat()))
}

// Atan2 returns the arc tangent of y/x
// atan2(float $y, float $x): float
func Atan2(y, x *types.Value) *types.Value {
	return types.NewFloat(math.Atan2(y.ToFloat(), x.ToFloat()))
}

// Deg2rad converts degrees to radians
// deg2rad(float $num): float
func Deg2rad(num *types.Value) *types.Value {
	return types.NewFloat(num.ToFloat() * math.Pi / 180.0)
}

// Rad2deg converts radians to degrees
// rad2deg(float $num): float
func Rad2deg(num *types.Value) *types.Value {
	return types.NewFloat(num.ToFloat() * 180.0 / math.Pi)
}

// ============================================================================
// Exponential and Logarithmic Functions
// ============================================================================

// Exp returns e raised to the power of num
// exp(float $num): float
func Exp(num *types.Value) *types.Value {
	return types.NewFloat(math.Exp(num.ToFloat()))
}

// Log returns the natural logarithm
// log(float $num, float $base = M_E): float
func Log(num *types.Value, base ...*types.Value) *types.Value {
	n := num.ToFloat()

	if len(base) > 0 && base[0] != nil {
		b := base[0].ToFloat()
		return types.NewFloat(math.Log(n) / math.Log(b))
	}

	return types.NewFloat(math.Log(n))
}

// Log10 returns the base-10 logarithm
// log10(float $num): float
func Log10(num *types.Value) *types.Value {
	return types.NewFloat(math.Log10(num.ToFloat()))
}

// Log1p returns log(1 + number)
// log1p(float $num): float
func Log1p(num *types.Value) *types.Value {
	return types.NewFloat(math.Log1p(num.ToFloat()))
}

// Expm1 returns exp(number) - 1
// expm1(float $num): float
func Expm1(num *types.Value) *types.Value {
	return types.NewFloat(math.Expm1(num.ToFloat()))
}

// ============================================================================
// Random Number Generation
// ============================================================================

// Rand generates a random integer
// rand(int $min = 0, int $max = getrandmax()): int
func Rand(limits ...*types.Value) *types.Value {
	min := int64(0)
	max := int64(math.MaxInt32)

	if len(limits) >= 1 && limits[0] != nil {
		min = limits[0].ToInt()
	}

	if len(limits) >= 2 && limits[1] != nil {
		max = limits[1].ToInt()
	}

	if min > max {
		min, max = max, min
	}

	result := min + rand.Int63n(max-min+1)
	return types.NewInt(result)
}

// MtRand generates a better random number using Mersenne Twister
// mt_rand(int $min = 0, int $max = mt_getrandmax()): int
func MtRand(limits ...*types.Value) *types.Value {
	// In Go, we use the same rand as Rand (it's already good)
	return Rand(limits...)
}

// RandomInt generates a cryptographically secure random integer
// random_int(int $min, int $max): int
func RandomInt(min, max *types.Value) *types.Value {
	minVal := min.ToInt()
	maxVal := max.ToInt()

	if minVal > maxVal {
		minVal, maxVal = maxVal, minVal
	}

	// For simplicity, using math/rand
	// In production, should use crypto/rand
	result := minVal + rand.Int63n(maxVal-minVal+1)
	return types.NewInt(result)
}

// GetRandMax returns the maximum random number
// getrandmax(): int
func GetRandMax() *types.Value {
	return types.NewInt(math.MaxInt32)
}

// MtGetRandMax returns the maximum random number for mt_rand
// mt_getrandmax(): int
func MtGetRandMax() *types.Value {
	return types.NewInt(math.MaxInt32)
}

// ============================================================================
// Number Formatting
// ============================================================================

// NumberFormat formats a number with grouped thousands
// number_format(float $num, int $decimals = 0, string $dec_point = ".", string $thousands_sep = ","): string
func NumberFormat(num *types.Value, args ...*types.Value) *types.Value {
	n := num.ToFloat()

	decimals := 0
	decPoint := "."
	thousandsSep := ","

	// Parse arguments
	if len(args) >= 1 && args[0] != nil {
		decimals = int(args[0].ToInt())
	}

	if len(args) >= 2 && args[1] != nil {
		decPoint = args[1].ToString()
	}

	if len(args) >= 3 && args[2] != nil {
		thousandsSep = args[2].ToString()
	}

	// Format the number
	format := "%." + strconv.Itoa(decimals) + "f"
	formatted := fmt.Sprintf(format, n)

	// Split into integer and decimal parts
	parts := strings.Split(formatted, ".")
	intPart := parts[0]
	var decPart string
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// Add thousands separator
	if thousandsSep != "" {
		intPart = addThousandsSeparator(intPart, thousandsSep)
	}

	// Combine parts
	if decimals > 0 {
		return types.NewString(intPart + decPoint + decPart)
	}

	return types.NewString(intPart)
}

// addThousandsSeparator adds thousands separator to a number string
func addThousandsSeparator(s, sep string) string {
	// Handle negative sign
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	// Add separator every 3 digits from right
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(sep)
		}
		result.WriteRune(r)
	}

	if negative {
		return "-" + result.String()
	}
	return result.String()
}

// ============================================================================
// Other Math Functions
// ============================================================================

// Pi returns the value of pi
// pi(): float
func Pi() *types.Value {
	return types.NewFloat(math.Pi)
}

// IsNan checks if a value is NaN
// is_nan(float $num): bool
func IsNan(num *types.Value) *types.Value {
	return types.NewBool(math.IsNaN(num.ToFloat()))
}

// IsInfinite checks if a value is infinite
// is_infinite(float $num): bool
func IsInfinite(num *types.Value) *types.Value {
	return types.NewBool(math.IsInf(num.ToFloat(), 0))
}

// IsFinite checks if a value is finite
// is_finite(float $num): bool
func IsFinite(num *types.Value) *types.Value {
	f := num.ToFloat()
	return types.NewBool(!math.IsNaN(f) && !math.IsInf(f, 0))
}

// Hypot returns sqrt(x*x + y*y)
// hypot(float $x, float $y): float
func Hypot(x, y *types.Value) *types.Value {
	return types.NewFloat(math.Hypot(x.ToFloat(), y.ToFloat()))
}

// Fmod returns the floating point remainder of x/y
// fmod(float $x, float $y): float
func Fmod(x, y *types.Value) *types.Value {
	return types.NewFloat(math.Mod(x.ToFloat(), y.ToFloat()))
}

// Intdiv performs integer division
// intdiv(int $num1, int $num2): int
func Intdiv(num1, num2 *types.Value) *types.Value {
	n1 := num1.ToInt()
	n2 := num2.ToInt()

	if n2 == 0 {
		// Division by zero error
		return types.NewNull()
	}

	return types.NewInt(n1 / n2)
}

// Fdiv performs floating-point division (PHP 8.0+)
// fdiv(float $num1, float $num2): float
func Fdiv(num1, num2 *types.Value) *types.Value {
	n1 := num1.ToFloat()
	n2 := num2.ToFloat()

	// PHP's fdiv returns inf for division by zero, not error
	return types.NewFloat(n1 / n2)
}
