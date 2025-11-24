package types

import (
	"testing"
)

// BenchmarkNewInt benchmarks integer value creation
func BenchmarkNewInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewInt(int64(i))
	}
}

// BenchmarkNewFloat benchmarks float value creation
func BenchmarkNewFloat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFloat(3.14)
	}
}

// BenchmarkNewString benchmarks string value creation
func BenchmarkNewString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewString("Hello, World!")
	}
}

// BenchmarkNewBool benchmarks boolean value creation
func BenchmarkNewBool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewBool(true)
	}
}

// BenchmarkNewNull benchmarks null value creation
func BenchmarkNewNull(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewNull()
	}
}

// BenchmarkIntToString benchmarks integer to string conversion
func BenchmarkIntToString(b *testing.B) {
	val := NewInt(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToString()
	}
}

// BenchmarkIntToFloat benchmarks integer to float conversion
func BenchmarkIntToFloat(b *testing.B) {
	val := NewInt(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToFloat()
	}
}

// BenchmarkIntToBool benchmarks integer to boolean conversion
func BenchmarkIntToBool(b *testing.B) {
	val := NewInt(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToBool()
	}
}

// BenchmarkStringToInt benchmarks string to integer conversion
func BenchmarkStringToInt(b *testing.B) {
	val := NewString("12345")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToInt()
	}
}

// BenchmarkStringToFloat benchmarks string to float conversion
func BenchmarkStringToFloat(b *testing.B) {
	val := NewString("3.14159")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToFloat()
	}
}

// BenchmarkStringToBool benchmarks string to boolean conversion
func BenchmarkStringToBool(b *testing.B) {
	val := NewString("hello")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToBool()
	}
}

// BenchmarkFloatToInt benchmarks float to integer conversion
func BenchmarkFloatToInt(b *testing.B) {
	val := NewFloat(3.14159)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToInt()
	}
}

// BenchmarkFloatToString benchmarks float to string conversion
func BenchmarkFloatToString(b *testing.B) {
	val := NewFloat(3.14159)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToString()
	}
}

// BenchmarkValueEquals benchmarks loose equality comparison
func BenchmarkValueEquals(b *testing.B) {
	v1 := NewInt(42)
	v2 := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Equals(v2)
	}
}

// BenchmarkValueIdentical benchmarks strict equality comparison
func BenchmarkValueIdentical(b *testing.B) {
	v1 := NewInt(42)
	v2 := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Identical(v2)
	}
}

// BenchmarkStringEquality benchmarks string equality comparison
func BenchmarkStringEquality(b *testing.B) {
	v1 := NewString("Hello, World!")
	v2 := NewString("Hello, World!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Equals(v2)
	}
}

// BenchmarkTypeJugglingIntString benchmarks type juggling between int and string
func BenchmarkTypeJugglingIntString(b *testing.B) {
	v1 := NewInt(42)
	v2 := NewString("42")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Equals(v2)
	}
}

// BenchmarkTypeJugglingStringInt benchmarks type juggling string comparisons
func BenchmarkTypeJugglingStringInt(b *testing.B) {
	v1 := NewString("42")
	v2 := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Equals(v2)
	}
}

// BenchmarkArrayCreation benchmarks array creation
func BenchmarkArrayCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewArrayWithCapacity(10)
	}
}

// BenchmarkValueType benchmarks getting value type
func BenchmarkValueType(b *testing.B) {
	val := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Type()
	}
}

// BenchmarkValueIsCallable benchmarks checking if value is callable
func BenchmarkValueIsCallable(b *testing.B) {
	val := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.IsCallable()
	}
}

// BenchmarkStringConcatenation benchmarks string concatenation
func BenchmarkStringConcatenation(b *testing.B) {
	v1 := NewString("Hello")
	v2 := NewString(", ")
	v3 := NewString("World")
	v4 := NewString("!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1 := v1.ToString() + v2.ToString()
		s2 := s1 + v3.ToString()
		_ = s2 + v4.ToString()
	}
}

// BenchmarkValueCopy benchmarks copying values
func BenchmarkValueCopy(b *testing.B) {
	val := NewInt(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = *val
	}
}

// BenchmarkComplexTypeConversion benchmarks multiple type conversions
func BenchmarkComplexTypeConversion(b *testing.B) {
	val := NewString("42")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.ToInt()
		val.ToFloat()
		val.ToBool()
		val.ToString()
	}
}

// BenchmarkMultipleValueCreation benchmarks creating multiple values
func BenchmarkMultipleValueCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewInt(42)
		NewString("test")
		NewFloat(3.14)
		NewBool(true)
		NewNull()
	}
}

// BenchmarkBoolConversions benchmarks boolean conversions from different types
func BenchmarkBoolConversions(b *testing.B) {
	vInt := NewInt(1)
	vFloat := NewFloat(1.5)
	vString := NewString("hello")
	vEmpty := NewString("")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vInt.ToBool()
		vFloat.ToBool()
		vString.ToBool()
		vEmpty.ToBool()
	}
}

// BenchmarkNumericStringParsing benchmarks parsing numeric strings
func BenchmarkNumericStringParsing(b *testing.B) {
	v1 := NewString("12345")
	v2 := NewString("3.14159")
	v3 := NewString("1e10")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.ToInt()
		v2.ToFloat()
		v3.ToFloat()
	}
}
