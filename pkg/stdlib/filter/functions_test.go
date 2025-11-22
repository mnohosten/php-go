package filter

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Validate Boolean Tests
// ============================================================================

func TestFilterVarValidateBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"on", true},
		{"yes", true},
		{"0", true},
		{"false", true},
		{"off", true},
		{"no", true},
		{"", true},
		{"invalid", false},
		{"2", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_BOOLEAN))
		if result.ToBool() != tt.expected {
			t.Errorf("FilterVar(%q, FILTER_VALIDATE_BOOLEAN) = %v, want %v", tt.input, result.ToBool(), tt.expected)
		}
	}
}

// ============================================================================
// Validate Email Tests
// ============================================================================

func TestFilterVarValidateEmail(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_EMAIL))
		if tt.valid {
			if result.Type() != types.TypeString || result.ToString() != tt.input {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_EMAIL) should return string, got %v", tt.input, result)
			}
		} else {
			if result.Type() != types.TypeBool || result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_EMAIL) should return false, got %v", tt.input, result)
			}
		}
	}
}

// ============================================================================
// Validate Float Tests
// ============================================================================

func TestFilterVarValidateFloat(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"3.14", true},
		{"0.5", true},
		{"-2.5", true},
		{"123", true},
		{"invalid", false},
		{"12.34.56", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_FLOAT))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_FLOAT) should return string, got %v", tt.input, result.Type())
			}
		} else {
			if result.Type() != types.TypeBool || result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_FLOAT) should return false, got %v", tt.input, result)
			}
		}
	}
}

// ============================================================================
// Validate Int Tests
// ============================================================================

func TestFilterVarValidateInt(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"123", true},
		{"0", true},
		{"-456", true},
		{"3.14", false},
		{"invalid", false},
		{"12a34", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_INT))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_INT) should return string, got %v", tt.input, result.Type())
			}
		} else {
			if result.Type() != types.TypeBool || result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_INT) should return false, got %v", tt.input, result)
			}
		}
	}
}

func TestFilterVarValidateIntHex(t *testing.T) {
	result := FilterVar(
		types.NewString("0xFF"),
		types.NewInt(FILTER_VALIDATE_INT),
		types.NewInt(FILTER_FLAG_ALLOW_HEX),
	)

	if result.Type() != types.TypeString {
		t.Errorf("FilterVar('0xFF', FILTER_VALIDATE_INT, FILTER_FLAG_ALLOW_HEX) should return string")
	}
}

// ============================================================================
// Validate IP Tests
// ============================================================================

func TestFilterVarValidateIP(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"192.168.1.1", true},
		{"127.0.0.1", true},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"::1", true},
		{"invalid", false},
		{"256.1.1.1", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_IP))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_IP) should return string, got %v", tt.input, result.Type())
			}
		} else {
			if result.Type() != types.TypeBool || result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_IP) should return false", tt.input)
			}
		}
	}
}

func TestFilterVarValidateIPv4(t *testing.T) {
	// IPv4 only
	result := FilterVar(
		types.NewString("192.168.1.1"),
		types.NewInt(FILTER_VALIDATE_IP),
		types.NewInt(FILTER_FLAG_IPV4),
	)
	if result.Type() != types.TypeString {
		t.Errorf("IPv4 address should validate with FILTER_FLAG_IPV4")
	}

	// IPv6 with IPv4 flag should fail
	result = FilterVar(
		types.NewString("::1"),
		types.NewInt(FILTER_VALIDATE_IP),
		types.NewInt(FILTER_FLAG_IPV4),
	)
	if result.ToBool() != false {
		t.Errorf("IPv6 address should not validate with FILTER_FLAG_IPV4")
	}
}

// ============================================================================
// Validate MAC Tests
// ============================================================================

func TestFilterVarValidateMAC(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"00:11:22:33:44:55", true},
		{"00-11-22-33-44-55", true},
		{"invalid", false},
		{"00:11:22:33:44", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_MAC))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_MAC) should return string", tt.input)
			}
		} else {
			if result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_MAC) should return false", tt.input)
			}
		}
	}
}

// ============================================================================
// Validate URL Tests
// ============================================================================

func TestFilterVarValidateURL(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"http://example.com", true},
		{"https://www.example.com/path", true},
		{"ftp://ftp.example.com", true},
		{"invalid", false},
		{"http://", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_URL))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_URL) should return string", tt.input)
			}
		} else {
			if result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_URL) should return false", tt.input)
			}
		}
	}
}

// ============================================================================
// Validate Domain Tests
// ============================================================================

func TestFilterVarValidateDomain(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"example.com", true},
		{"sub.example.com", true},
		{"example.co.uk", true},
		{"invalid", false},
		{"", false},
		{".com", false},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_VALIDATE_DOMAIN))
		if tt.valid {
			if result.Type() != types.TypeString {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_DOMAIN) should return string", tt.input)
			}
		} else {
			if result.ToBool() != false {
				t.Errorf("FilterVar(%q, FILTER_VALIDATE_DOMAIN) should return false", tt.input)
			}
		}
	}
}

// ============================================================================
// Sanitize Email Tests
// ============================================================================

func TestFilterVarSanitizeEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "test@example.com"},
		{"test+tag@example.com", "test+tag@example.com"},
		{"test<script>@example.com", "testscript@example.com"}, // '<' and '>' removed, but 'script' stays
		{"user name@example.com", "username@example.com"},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_SANITIZE_EMAIL))
		if result.ToString() != tt.expected {
			t.Errorf("FilterVar(%q, FILTER_SANITIZE_EMAIL) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Sanitize Number Float Tests
// ============================================================================

func TestFilterVarSanitizeNumberFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3.14", "3.14"},
		{"12.34abc", "12.34"},
		{"+3.14", "+3.14"},
		{"-2.5", "-2.5"},
		{"1,234.56", "1234.56"},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_SANITIZE_NUMBER_FLOAT))
		if result.ToString() != tt.expected {
			t.Errorf("FilterVar(%q, FILTER_SANITIZE_NUMBER_FLOAT) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Sanitize Number Int Tests
// ============================================================================

func TestFilterVarSanitizeNumberInt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"123abc", "123"},
		{"+456", "+456"},
		{"-789", "-789"},
		{"12.34", "1234"},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_SANITIZE_NUMBER_INT))
		if result.ToString() != tt.expected {
			t.Errorf("FilterVar(%q, FILTER_SANITIZE_NUMBER_INT) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Sanitize String Tests
// ============================================================================

func TestFilterVarSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello<script>alert('xss')</script>", "helloalert('xss')"},
		{"<b>bold</b>", "bold"},
	}

	for _, tt := range tests {
		result := FilterVar(types.NewString(tt.input), types.NewInt(FILTER_SANITIZE_STRING))
		if result.ToString() != tt.expected {
			t.Errorf("FilterVar(%q, FILTER_SANITIZE_STRING) = %q, want %q", tt.input, result.ToString(), tt.expected)
		}
	}
}

// ============================================================================
// Sanitize URL Tests
// ============================================================================

func TestFilterVarSanitizeURL(t *testing.T) {
	result := FilterVar(types.NewString("http://example.com"), types.NewInt(FILTER_SANITIZE_URL))
	if result.Type() != types.TypeString {
		t.Errorf("FILTER_SANITIZE_URL should return string")
	}

	// URL should be sanitized
	str := result.ToString()
	if str == "" {
		t.Errorf("FILTER_SANITIZE_URL should not return empty string for valid URL")
	}
}

// ============================================================================
// Filter Var Array Tests
// ============================================================================

func TestFilterVarArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("email"), types.NewString("test@example.com"))
	arr.Set(types.NewString("number"), types.NewString("123"))

	result := FilterVarArray(types.NewArray(arr), types.NewInt(FILTER_UNSAFE_RAW))

	if result.Type() != types.TypeArray {
		t.Errorf("FilterVarArray should return array, got %v", result.Type())
	}

	resultArr := result.ToArray()
	if resultArr.Len() != 2 {
		t.Errorf("FilterVarArray should return array with 2 elements, got %d", resultArr.Len())
	}
}

func TestFilterVarArrayInvalidInput(t *testing.T) {
	result := FilterVarArray(types.NewString("not an array"))

	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("FilterVarArray with non-array should return false")
	}
}

// ============================================================================
// Default Filter Tests
// ============================================================================

func TestFilterVarDefault(t *testing.T) {
	input := types.NewString("test")
	result := FilterVar(input)

	// FILTER_UNSAFE_RAW (default) should return unchanged
	if result.ToString() != "test" {
		t.Errorf("FilterVar with default filter should return unchanged value")
	}
}

func TestFilterVarUnsafeRaw(t *testing.T) {
	input := types.NewString("<script>alert('xss')</script>")
	result := FilterVar(input, types.NewInt(FILTER_UNSAFE_RAW))

	// FILTER_UNSAFE_RAW should not sanitize
	if result.ToString() != input.ToString() {
		t.Errorf("FILTER_UNSAFE_RAW should return unchanged value")
	}
}
