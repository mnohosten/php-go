package util

import (
	"math"
	"testing"
)

func TestSafeAddSize(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		wantErr   bool
		wantValue int
	}{
		{
			name:      "Normal addition",
			a:         10,
			b:         20,
			wantErr:   false,
			wantValue: 30,
		},
		{
			name:      "Zero addition",
			a:         0,
			b:         10,
			wantErr:   false,
			wantValue: 10,
		},
		{
			name:      "Both zero",
			a:         0,
			b:         0,
			wantErr:   false,
			wantValue: 0,
		},
		{
			name:      "Negative first value",
			a:         -1,
			b:         10,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Negative second value",
			a:         10,
			b:         -1,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Overflow at MaxInt",
			a:         math.MaxInt,
			b:         1,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Large but safe addition",
			a:         math.MaxInt - 100,
			b:         50,
			wantErr:   false,
			wantValue: math.MaxInt - 50,
		},
		{
			name:      "Overflow with large values",
			a:         math.MaxInt / 2 + 1,
			b:         math.MaxInt / 2 + 1,
			wantErr:   true,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeAddSize(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeAddSize(%d, %d) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantValue {
				t.Errorf("SafeAddSize(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.wantValue)
			}
		})
	}
}

func TestSafeMultiplySize(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		wantErr   bool
		wantValue int
	}{
		{
			name:      "Normal multiplication",
			a:         10,
			b:         20,
			wantErr:   false,
			wantValue: 200,
		},
		{
			name:      "Zero first value",
			a:         0,
			b:         100,
			wantErr:   false,
			wantValue: 0,
		},
		{
			name:      "Zero second value",
			a:         100,
			b:         0,
			wantErr:   false,
			wantValue: 0,
		},
		{
			name:      "Both zero",
			a:         0,
			b:         0,
			wantErr:   false,
			wantValue: 0,
		},
		{
			name:      "Negative first value",
			a:         -1,
			b:         10,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Negative second value",
			a:         10,
			b:         -1,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Overflow at MaxInt",
			a:         math.MaxInt,
			b:         2,
			wantErr:   true,
			wantValue: 0,
		},
		{
			name:      "Large but safe multiplication",
			a:         math.MaxInt / 1000,
			b:         100,
			wantErr:   false,
			wantValue: (math.MaxInt / 1000) * 100,
		},
		{
			name:      "Multiply by 1",
			a:         12345,
			b:         1,
			wantErr:   false,
			wantValue: 12345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeMultiplySize(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeMultiplySize(%d, %d) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantValue {
				t.Errorf("SafeMultiplySize(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.wantValue)
			}
		})
	}
}

func TestSafeConvertToInt(t *testing.T) {
	tests := []struct {
		name      string
		val       int64
		fieldName string
		wantErr   bool
		wantValue int
	}{
		{
			name:      "Normal conversion",
			val:       12345,
			fieldName: "offset",
			wantErr:   false,
			wantValue: 12345,
		},
		{
			name:      "Zero conversion",
			val:       0,
			fieldName: "length",
			wantErr:   false,
			wantValue: 0,
		},
		{
			name:      "Negative conversion",
			val:       -100,
			fieldName: "offset",
			wantErr:   false,
			wantValue: -100,
		},
		{
			name:      "MaxInt conversion",
			val:       int64(math.MaxInt),
			fieldName: "size",
			wantErr:   false,
			wantValue: math.MaxInt,
		},
		{
			name:      "MinInt conversion",
			val:       int64(math.MinInt),
			fieldName: "size",
			wantErr:   false,
			wantValue: math.MinInt,
		},
	}

	// Only test overflow on 32-bit systems
	// On 64-bit systems where int == int64, SafeConvertToInt can never overflow
	// since int64 cannot hold values outside the int range
	const is32Bit = math.MaxInt == (1<<31 - 1)
	if is32Bit {
		// On 32-bit systems, test values beyond int32 range
		tests = append(tests, []struct {
			name      string
			val       int64
			fieldName string
			wantErr   bool
			wantValue int
		}{
			{
				name:      "Overflow positive (32-bit)",
				val:       2147483648, // MaxInt32 + 1
				fieldName: "offset",
				wantErr:   true,
				wantValue: 0,
			},
			{
				name:      "Overflow negative (32-bit)",
				val:       -2147483649, // MinInt32 - 1
				fieldName: "offset",
				wantErr:   true,
				wantValue: 0,
			},
			{
				name:      "Large positive overflow",
				val:       math.MaxInt64,
				fieldName: "length",
				wantErr:   true,
				wantValue: 0,
			},
		}...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeConvertToInt(tt.val, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeConvertToInt(%d, %q) error = %v, wantErr %v", tt.val, tt.fieldName, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantValue {
				t.Errorf("SafeConvertToInt(%d, %q) = %v, want %v", tt.val, tt.fieldName, got, tt.wantValue)
			}
		})
	}
}

func TestValidateIntRange(t *testing.T) {
	tests := []struct {
		name      string
		val       int64
		min       int64
		max       int64
		fieldName string
		wantErr   bool
	}{
		{
			name:      "Within range",
			val:       50,
			min:       0,
			max:       100,
			fieldName: "value",
			wantErr:   false,
		},
		{
			name:      "At minimum",
			val:       0,
			min:       0,
			max:       100,
			fieldName: "value",
			wantErr:   false,
		},
		{
			name:      "At maximum",
			val:       100,
			min:       0,
			max:       100,
			fieldName: "value",
			wantErr:   false,
		},
		{
			name:      "Below minimum",
			val:       -1,
			min:       0,
			max:       100,
			fieldName: "hour",
			wantErr:   true,
		},
		{
			name:      "Above maximum",
			val:       101,
			min:       0,
			max:       100,
			fieldName: "hour",
			wantErr:   true,
		},
		{
			name:      "Negative range",
			val:       -50,
			min:       -100,
			max:       0,
			fieldName: "temperature",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIntRange(tt.val, tt.min, tt.max, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIntRange(%d, %d, %d, %q) error = %v, wantErr %v",
					tt.val, tt.min, tt.max, tt.fieldName, err, tt.wantErr)
			}
		})
	}
}
