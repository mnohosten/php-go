package math

import (
	"math"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Basic Math Functions Tests
// ============================================================================

func TestAbs(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected interface{}
	}{
		{types.NewInt(5), int64(5)},
		{types.NewInt(-5), int64(5)},
		{types.NewInt(0), int64(0)},
		{types.NewFloat(3.14), 3.14},
		{types.NewFloat(-3.14), 3.14},
	}

	for _, tt := range tests {
		result := Abs(tt.input)

		switch exp := tt.expected.(type) {
		case int64:
			if result.Type() != types.TypeInt || result.ToInt() != exp {
				t.Errorf("Abs(%v) = %v, want %v", tt.input, result, exp)
			}
		case float64:
			if result.Type() != types.TypeFloat || result.ToFloat() != exp {
				t.Errorf("Abs(%v) = %v, want %v", tt.input, result, exp)
			}
		}
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected float64
	}{
		{types.NewFloat(4.3), 5.0},
		{types.NewFloat(9.999), 10.0},
		{types.NewFloat(-3.14), -3.0},
		{types.NewInt(5), 5.0},
	}

	for _, tt := range tests {
		result := Ceil(tt.input)
		if result.ToFloat() != tt.expected {
			t.Errorf("Ceil(%v) = %v, want %v", tt.input, result.ToFloat(), tt.expected)
		}
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected float64
	}{
		{types.NewFloat(4.3), 4.0},
		{types.NewFloat(9.999), 9.0},
		{types.NewFloat(-3.14), -4.0},
		{types.NewInt(5), 5.0},
	}

	for _, tt := range tests {
		result := Floor(tt.input)
		if result.ToFloat() != tt.expected {
			t.Errorf("Floor(%v) = %v, want %v", tt.input, result.ToFloat(), tt.expected)
		}
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input     *types.Value
		precision *types.Value
		expected  float64
	}{
		{types.NewFloat(3.4), nil, 3.0},
		{types.NewFloat(3.5), nil, 4.0},
		{types.NewFloat(3.6), nil, 4.0},
		{types.NewFloat(3.14159), types.NewInt(2), 3.14},
		{types.NewFloat(1234.5678), types.NewInt(2), 1234.57},
	}

	for _, tt := range tests {
		var result *types.Value
		if tt.precision == nil {
			result = Round(tt.input)
		} else {
			result = Round(tt.input, tt.precision)
		}

		if result.ToFloat() != tt.expected {
			t.Errorf("Round(%v, %v) = %v, want %v", tt.input, tt.precision, result.ToFloat(), tt.expected)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		values   []*types.Value
		expected interface{}
	}{
		{[]*types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)}, int64(1)},
		{[]*types.Value{types.NewInt(5), types.NewInt(2), types.NewInt(8)}, int64(2)},
		{[]*types.Value{types.NewFloat(1.5), types.NewFloat(0.5), types.NewFloat(2.5)}, 0.5},
	}

	for _, tt := range tests {
		result := Min(tt.values...)

		switch exp := tt.expected.(type) {
		case int64:
			if result.ToInt() != exp {
				t.Errorf("Min(%v) = %v, want %v", tt.values, result, exp)
			}
		case float64:
			if result.ToFloat() != exp {
				t.Errorf("Min(%v) = %v, want %v", tt.values, result, exp)
			}
		}
	}
}

func TestMinWithArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(5))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(8))

	result := Min(types.NewArray(arr))
	if result.ToInt() != 2 {
		t.Errorf("Min(array) = %v, want 2", result.ToInt())
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		values   []*types.Value
		expected interface{}
	}{
		{[]*types.Value{types.NewInt(1), types.NewInt(2), types.NewInt(3)}, int64(3)},
		{[]*types.Value{types.NewInt(5), types.NewInt(2), types.NewInt(8)}, int64(8)},
		{[]*types.Value{types.NewFloat(1.5), types.NewFloat(0.5), types.NewFloat(2.5)}, 2.5},
	}

	for _, tt := range tests {
		result := Max(tt.values...)

		switch exp := tt.expected.(type) {
		case int64:
			if result.ToInt() != exp {
				t.Errorf("Max(%v) = %v, want %v", tt.values, result, exp)
			}
		case float64:
			if result.ToFloat() != exp {
				t.Errorf("Max(%v) = %v, want %v", tt.values, result, exp)
			}
		}
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		base     *types.Value
		exp      *types.Value
		expected interface{}
	}{
		{types.NewInt(2), types.NewInt(3), int64(8)},
		{types.NewInt(5), types.NewInt(2), int64(25)},
		{types.NewFloat(2.0), types.NewFloat(0.5), 1.4142135623730951},
	}

	for _, tt := range tests {
		result := Pow(tt.base, tt.exp)

		switch exp := tt.expected.(type) {
		case int64:
			if result.Type() != types.TypeInt || result.ToInt() != exp {
				t.Errorf("Pow(%v, %v) = %v, want %v", tt.base, tt.exp, result, exp)
			}
		case float64:
			if math.Abs(result.ToFloat()-exp) > 0.0001 {
				t.Errorf("Pow(%v, %v) = %v, want %v", tt.base, tt.exp, result.ToFloat(), exp)
			}
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input    *types.Value
		expected float64
	}{
		{types.NewFloat(4.0), 2.0},
		{types.NewFloat(9.0), 3.0},
		{types.NewFloat(2.0), 1.4142135623730951},
		{types.NewInt(16), 4.0},
	}

	for _, tt := range tests {
		result := Sqrt(tt.input)
		if math.Abs(result.ToFloat()-tt.expected) > 0.0001 {
			t.Errorf("Sqrt(%v) = %v, want %v", tt.input, result.ToFloat(), tt.expected)
		}
	}
}

func TestSqrtNegative(t *testing.T) {
	result := Sqrt(types.NewFloat(-1.0))
	if !math.IsNaN(result.ToFloat()) {
		t.Errorf("Sqrt(-1) should return NaN, got %v", result.ToFloat())
	}
}

// ============================================================================
// Trigonometric Functions Tests
// ============================================================================

func TestSin(t *testing.T) {
	result := Sin(types.NewFloat(math.Pi / 2))
	if math.Abs(result.ToFloat()-1.0) > 0.0001 {
		t.Errorf("Sin(π/2) = %v, want 1.0", result.ToFloat())
	}
}

func TestCos(t *testing.T) {
	result := Cos(types.NewFloat(0))
	if math.Abs(result.ToFloat()-1.0) > 0.0001 {
		t.Errorf("Cos(0) = %v, want 1.0", result.ToFloat())
	}
}

func TestTan(t *testing.T) {
	result := Tan(types.NewFloat(math.Pi / 4))
	if math.Abs(result.ToFloat()-1.0) > 0.0001 {
		t.Errorf("Tan(π/4) = %v, want 1.0", result.ToFloat())
	}
}

func TestAsin(t *testing.T) {
	result := Asin(types.NewFloat(1.0))
	if math.Abs(result.ToFloat()-math.Pi/2) > 0.0001 {
		t.Errorf("Asin(1) = %v, want π/2", result.ToFloat())
	}
}

func TestAcos(t *testing.T) {
	result := Acos(types.NewFloat(1.0))
	if math.Abs(result.ToFloat()-0.0) > 0.0001 {
		t.Errorf("Acos(1) = %v, want 0", result.ToFloat())
	}
}

func TestAtan(t *testing.T) {
	result := Atan(types.NewFloat(1.0))
	if math.Abs(result.ToFloat()-math.Pi/4) > 0.0001 {
		t.Errorf("Atan(1) = %v, want π/4", result.ToFloat())
	}
}

func TestAtan2(t *testing.T) {
	result := Atan2(types.NewFloat(1.0), types.NewFloat(1.0))
	if math.Abs(result.ToFloat()-math.Pi/4) > 0.0001 {
		t.Errorf("Atan2(1, 1) = %v, want π/4", result.ToFloat())
	}
}

func TestDeg2rad(t *testing.T) {
	result := Deg2rad(types.NewFloat(180.0))
	if math.Abs(result.ToFloat()-math.Pi) > 0.0001 {
		t.Errorf("Deg2rad(180) = %v, want π", result.ToFloat())
	}
}

func TestRad2deg(t *testing.T) {
	result := Rad2deg(types.NewFloat(math.Pi))
	if math.Abs(result.ToFloat()-180.0) > 0.0001 {
		t.Errorf("Rad2deg(π) = %v, want 180", result.ToFloat())
	}
}

// ============================================================================
// Exponential and Logarithmic Functions Tests
// ============================================================================

func TestExp(t *testing.T) {
	result := Exp(types.NewFloat(1.0))
	if math.Abs(result.ToFloat()-math.E) > 0.0001 {
		t.Errorf("Exp(1) = %v, want e", result.ToFloat())
	}
}

func TestLog(t *testing.T) {
	result := Log(types.NewFloat(math.E))
	if math.Abs(result.ToFloat()-1.0) > 0.0001 {
		t.Errorf("Log(e) = %v, want 1.0", result.ToFloat())
	}
}

func TestLogWithBase(t *testing.T) {
	result := Log(types.NewFloat(100.0), types.NewFloat(10.0))
	if math.Abs(result.ToFloat()-2.0) > 0.0001 {
		t.Errorf("Log(100, 10) = %v, want 2.0", result.ToFloat())
	}
}

func TestLog10(t *testing.T) {
	result := Log10(types.NewFloat(100.0))
	if math.Abs(result.ToFloat()-2.0) > 0.0001 {
		t.Errorf("Log10(100) = %v, want 2.0", result.ToFloat())
	}
}

// ============================================================================
// Random Number Tests
// ============================================================================

func TestRand(t *testing.T) {
	result := Rand()
	if result.Type() != types.TypeInt {
		t.Errorf("Rand() should return int, got %v", result.Type())
	}
}

func TestRandWithLimits(t *testing.T) {
	min := types.NewInt(1)
	max := types.NewInt(10)
	result := Rand(min, max)

	val := result.ToInt()
	if val < 1 || val > 10 {
		t.Errorf("Rand(1, 10) = %v, should be between 1 and 10", val)
	}
}

func TestMtRand(t *testing.T) {
	result := MtRand()
	if result.Type() != types.TypeInt {
		t.Errorf("MtRand() should return int, got %v", result.Type())
	}
}

func TestRandomInt(t *testing.T) {
	min := types.NewInt(5)
	max := types.NewInt(15)
	result := RandomInt(min, max)

	val := result.ToInt()
	if val < 5 || val > 15 {
		t.Errorf("RandomInt(5, 15) = %v, should be between 5 and 15", val)
	}
}

func TestGetRandMax(t *testing.T) {
	result := GetRandMax()
	if result.ToInt() != math.MaxInt32 {
		t.Errorf("GetRandMax() = %v, want %v", result.ToInt(), math.MaxInt32)
	}
}

// ============================================================================
// Number Formatting Tests
// ============================================================================

func TestNumberFormat(t *testing.T) {
	tests := []struct {
		num      *types.Value
		args     []*types.Value
		expected string
	}{
		{types.NewFloat(1234.56), []*types.Value{}, "1,235"},
		{types.NewFloat(1234.56), []*types.Value{types.NewInt(2)}, "1,234.56"},
		{types.NewFloat(1234567.891), []*types.Value{types.NewInt(2)}, "1,234,567.89"},
		{types.NewFloat(1234.56), []*types.Value{types.NewInt(2), types.NewString(","), types.NewString(".")}, "1.234,56"},
	}

	for _, tt := range tests {
		result := NumberFormat(tt.num, tt.args...)
		if result.ToString() != tt.expected {
			t.Errorf("NumberFormat(%v, %v) = %v, want %v", tt.num, tt.args, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Other Math Functions Tests
// ============================================================================

func TestPi(t *testing.T) {
	result := Pi()
	if result.ToFloat() != math.Pi {
		t.Errorf("Pi() = %v, want %v", result.ToFloat(), math.Pi)
	}
}

func TestIsNan(t *testing.T) {
	result := IsNan(types.NewFloat(math.NaN()))
	if !result.ToBool() {
		t.Errorf("IsNan(NaN) should return true")
	}

	result = IsNan(types.NewFloat(1.0))
	if result.ToBool() {
		t.Errorf("IsNan(1.0) should return false")
	}
}

func TestIsInfinite(t *testing.T) {
	result := IsInfinite(types.NewFloat(math.Inf(1)))
	if !result.ToBool() {
		t.Errorf("IsInfinite(Inf) should return true")
	}

	result = IsInfinite(types.NewFloat(1.0))
	if result.ToBool() {
		t.Errorf("IsInfinite(1.0) should return false")
	}
}

func TestIsFinite(t *testing.T) {
	result := IsFinite(types.NewFloat(1.0))
	if !result.ToBool() {
		t.Errorf("IsFinite(1.0) should return true")
	}

	result = IsFinite(types.NewFloat(math.Inf(1)))
	if result.ToBool() {
		t.Errorf("IsFinite(Inf) should return false")
	}

	result = IsFinite(types.NewFloat(math.NaN()))
	if result.ToBool() {
		t.Errorf("IsFinite(NaN) should return false")
	}
}

func TestHypot(t *testing.T) {
	result := Hypot(types.NewFloat(3.0), types.NewFloat(4.0))
	if math.Abs(result.ToFloat()-5.0) > 0.0001 {
		t.Errorf("Hypot(3, 4) = %v, want 5.0", result.ToFloat())
	}
}

func TestFmod(t *testing.T) {
	result := Fmod(types.NewFloat(5.7), types.NewFloat(1.3))
	if math.Abs(result.ToFloat()-0.5) > 0.0001 {
		t.Errorf("Fmod(5.7, 1.3) = %v, want ~0.5", result.ToFloat())
	}
}

func TestIntdiv(t *testing.T) {
	result := Intdiv(types.NewInt(10), types.NewInt(3))
	if result.ToInt() != 3 {
		t.Errorf("Intdiv(10, 3) = %v, want 3", result.ToInt())
	}
}

func TestIntdivByZero(t *testing.T) {
	result := Intdiv(types.NewInt(10), types.NewInt(0))
	if result.Type() != types.TypeNull {
		t.Errorf("Intdiv(10, 0) should return NULL, got %v", result.Type())
	}
}

func TestFdiv(t *testing.T) {
	result := Fdiv(types.NewFloat(10.0), types.NewFloat(3.0))
	expected := 10.0 / 3.0
	if math.Abs(result.ToFloat()-expected) > 0.0001 {
		t.Errorf("Fdiv(10, 3) = %v, want %v", result.ToFloat(), expected)
	}
}

func TestFdivByZero(t *testing.T) {
	result := Fdiv(types.NewFloat(10.0), types.NewFloat(0.0))
	if !math.IsInf(result.ToFloat(), 1) {
		t.Errorf("Fdiv(10, 0) should return +Inf, got %v", result.ToFloat())
	}
}
