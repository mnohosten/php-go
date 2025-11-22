package date

import (
	"strings"
	"testing"
	"time"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Time Functions Tests
// ============================================================================

func TestTime(t *testing.T) {
	before := time.Now().Unix()
	result := Time()
	after := time.Now().Unix()

	timestamp := result.ToInt()
	if timestamp < before || timestamp > after {
		t.Errorf("Time() = %v, should be between %v and %v", timestamp, before, after)
	}
}

func TestMicrotimeString(t *testing.T) {
	result := Microtime()

	if result.Type() != types.TypeString {
		t.Errorf("Microtime() should return string by default")
	}

	// Should be in format "0.MICROSECONDS SECONDS"
	parts := strings.Split(result.ToString(), " ")
	if len(parts) != 2 {
		t.Errorf("Microtime() should return 'usec sec' format, got %v", result.ToString())
	}
}

func TestMicrotimeFloat(t *testing.T) {
	result := Microtime(types.NewBool(true))

	if result.Type() != types.TypeFloat {
		t.Errorf("Microtime(true) should return float")
	}

	now := float64(time.Now().Unix())
	mt := result.ToFloat()

	// Should be close to current time (within 1 second)
	if mt < now-1 || mt > now+1 {
		t.Errorf("Microtime(true) = %v, expected around %v", mt, now)
	}
}

// ============================================================================
// Date Formatting Tests
// ============================================================================

func TestDateBasicFormats(t *testing.T) {
	// Use a known timestamp in local time: 2024-03-15 14:30:45
	timestamp := time.Date(2024, 3, 15, 14, 30, 45, 0, time.Local).Unix()

	tests := []struct {
		format   string
		expected string
	}{
		{"Y-m-d", "2024-03-15"},
		{"Y", "2024"},
		{"y", "24"},
		{"m", "03"},
		{"n", "3"},
		{"d", "15"},
		{"j", "15"},
		{"H:i:s", "14:30:45"},
		{"H", "14"},
		{"i", "30"},
		{"s", "45"},
	}

	for _, tt := range tests {
		result := Date(types.NewString(tt.format), types.NewInt(timestamp))
		if result.ToString() != tt.expected {
			t.Errorf("date(%q, timestamp) = %v, want %v", tt.format, result.ToString(), tt.expected)
		}
	}
}

func TestDateMonthFormats(t *testing.T) {
	// January
	timestamp := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	tests := []struct {
		format   string
		expected string
	}{
		{"F", "January"},
		{"M", "Jan"},
		{"m", "01"},
		{"n", "1"},
	}

	for _, tt := range tests {
		result := Date(types.NewString(tt.format), types.NewInt(timestamp))
		if result.ToString() != tt.expected {
			t.Errorf("date(%q, Jan timestamp) = %v, want %v", tt.format, result.ToString(), tt.expected)
		}
	}
}

func TestDateDayFormats(t *testing.T) {
	// Monday, March 4, 2024
	timestamp := time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC).Unix()

	tests := []struct {
		format   string
		expected string
	}{
		{"l", "Monday"},
		{"D", "Mon"},
		{"N", "1"}, // Monday is 1 in ISO-8601
		{"w", "1"}, // Monday is 1 in PHP (0=Sunday)
	}

	for _, tt := range tests {
		result := Date(types.NewString(tt.format), types.NewInt(timestamp))
		if result.ToString() != tt.expected {
			t.Errorf("date(%q, Monday timestamp) = %v, want %v", tt.format, result.ToString(), tt.expected)
		}
	}
}

func TestDate12HourFormats(t *testing.T) {
	// 2PM (14:00) in local time
	timestamp := time.Date(2024, 1, 1, 14, 0, 0, 0, time.Local).Unix()

	tests := []struct {
		format   string
		expected string
	}{
		{"g", "2"},
		{"h", "02"},
		{"a", "pm"},
		{"A", "PM"},
	}

	for _, tt := range tests {
		result := Date(types.NewString(tt.format), types.NewInt(timestamp))
		if result.ToString() != tt.expected {
			t.Errorf("date(%q, 2PM timestamp) = %v, want %v", tt.format, result.ToString(), tt.expected)
		}
	}
}

func TestDateLeapYear(t *testing.T) {
	// 2024 is a leap year
	timestamp2024 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	result := Date(types.NewString("L"), types.NewInt(timestamp2024))
	if result.ToString() != "1" {
		t.Errorf("2024 should be a leap year, got L=%v", result.ToString())
	}

	// 2023 is not a leap year
	timestamp2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	result = Date(types.NewString("L"), types.NewInt(timestamp2023))
	if result.ToString() != "0" {
		t.Errorf("2023 should not be a leap year, got L=%v", result.ToString())
	}
}

func TestDateEscapeCharacter(t *testing.T) {
	timestamp := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	result := Date(types.NewString("Y\\Y"), types.NewInt(timestamp))
	if result.ToString() != "2024Y" {
		t.Errorf("date('Y\\\\Y') = %v, want '2024Y'", result.ToString())
	}
}

func TestGmdate(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC).Unix()

	result := Gmdate(types.NewString("Y-m-d H:i:s"), types.NewInt(timestamp))
	expected := "2024-03-15 14:30:45"

	if result.ToString() != expected {
		t.Errorf("gmdate() = %v, want %v", result.ToString(), expected)
	}
}

// ============================================================================
// Mktime Tests
// ============================================================================

func TestMktime(t *testing.T) {
	// mktime(14, 30, 45, 3, 15, 2024)
	result := Mktime(
		types.NewInt(14), // hour
		types.NewInt(30), // minute
		types.NewInt(45), // second
		types.NewInt(3),  // month
		types.NewInt(15), // day
		types.NewInt(2024), // year
	)

	expected := time.Date(2024, 3, 15, 14, 30, 45, 0, time.Local).Unix()
	if result.ToInt() != expected {
		t.Errorf("mktime(14,30,45,3,15,2024) = %v, want %v", result.ToInt(), expected)
	}
}

func TestMktimeTwoDigitYear(t *testing.T) {
	// Year 24 should become 2024
	result := Mktime(
		types.NewInt(0),
		types.NewInt(0),
		types.NewInt(0),
		types.NewInt(1),
		types.NewInt(1),
		types.NewInt(24),
	)

	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	if result.ToInt() != expected {
		t.Errorf("mktime with year 24 should give 2024")
	}

	// Year 95 should become 1995
	result = Mktime(
		types.NewInt(0),
		types.NewInt(0),
		types.NewInt(0),
		types.NewInt(1),
		types.NewInt(1),
		types.NewInt(95),
	)

	expected = time.Date(1995, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	if result.ToInt() != expected {
		t.Errorf("mktime with year 95 should give 1995")
	}
}

func TestGmmktime(t *testing.T) {
	result := Gmmktime(
		types.NewInt(14),
		types.NewInt(30),
		types.NewInt(45),
		types.NewInt(3),
		types.NewInt(15),
		types.NewInt(2024),
	)

	expected := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC).Unix()
	if result.ToInt() != expected {
		t.Errorf("gmmktime() = %v, want %v", result.ToInt(), expected)
	}
}

// ============================================================================
// Strtotime Tests
// ============================================================================

func TestStrtotimeFormats(t *testing.T) {
	tests := []struct {
		input    string
		expected string // Will be parsed to verify it's valid
	}{
		{"2024-03-15", "2024-03-15"},
		{"2024-03-15 14:30:45", "2024-03-15 14:30:45"},
		{"2024-03-15T14:30:45", "2024-03-15T14:30:45"},
	}

	for _, tt := range tests {
		result := Strtotime(types.NewString(tt.input))
		if result.Type() == types.TypeBool && !result.ToBool() {
			t.Errorf("strtotime(%q) failed to parse", tt.input)
		}
	}
}

func TestStrtotimeRelative(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"now"},
		{"today"},
		{"tomorrow"},
		{"yesterday"},
	}

	for _, tt := range tests {
		result := Strtotime(types.NewString(tt.input))
		if result.Type() == types.TypeBool {
			t.Errorf("strtotime(%q) should not return false", tt.input)
		}
	}
}

func TestStrtotimeNow(t *testing.T) {
	before := time.Now().Unix()
	result := Strtotime(types.NewString("now"))
	after := time.Now().Unix()

	timestamp := result.ToInt()
	if timestamp < before || timestamp > after {
		t.Errorf("strtotime('now') = %v, should be between %v and %v", timestamp, before, after)
	}
}

func TestStrtotimeInvalid(t *testing.T) {
	result := Strtotime(types.NewString("invalid date string"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("strtotime(invalid) should return false")
	}
}

// ============================================================================
// Getdate Tests
// ============================================================================

func TestGetdate(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 14, 30, 45, 0, time.Local).Unix()
	result := Getdate(types.NewInt(timestamp))

	if result.Type() != types.TypeArray {
		t.Errorf("getdate() should return array")
	}

	arr := result.ToArray()

	// Check some key values
	seconds, _ := arr.Get(types.NewString("seconds"))
	if seconds.ToInt() != 45 {
		t.Errorf("getdate()['seconds'] = %v, want 45", seconds.ToInt())
	}

	minutes, _ := arr.Get(types.NewString("minutes"))
	if minutes.ToInt() != 30 {
		t.Errorf("getdate()['minutes'] = %v, want 30", minutes.ToInt())
	}

	hours, _ := arr.Get(types.NewString("hours"))
	if hours.ToInt() != 14 {
		t.Errorf("getdate()['hours'] = %v, want 14", hours.ToInt())
	}

	mday, _ := arr.Get(types.NewString("mday"))
	if mday.ToInt() != 15 {
		t.Errorf("getdate()['mday'] = %v, want 15", mday.ToInt())
	}

	mon, _ := arr.Get(types.NewString("mon"))
	if mon.ToInt() != 3 {
		t.Errorf("getdate()['mon'] = %v, want 3", mon.ToInt())
	}

	year, _ := arr.Get(types.NewString("year"))
	if year.ToInt() != 2024 {
		t.Errorf("getdate()['year'] = %v, want 2024", year.ToInt())
	}
}

// ============================================================================
// Localtime Tests
// ============================================================================

func TestLocaltimeIndexed(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 14, 30, 45, 0, time.Local).Unix()
	result := Localtime(types.NewInt(timestamp))

	if result.Type() != types.TypeArray {
		t.Errorf("localtime() should return array")
	}

	arr := result.ToArray()

	// Check array length (should be 9 elements)
	if arr.Len() != 9 {
		t.Errorf("localtime() should return 9 elements, got %d", arr.Len())
	}

	// Check seconds (index 0)
	val, _ := arr.Get(types.NewInt(0))
	if val.ToInt() != 45 {
		t.Errorf("localtime()[0] (seconds) = %v, want 45", val.ToInt())
	}
}

func TestLocaltimeAssociative(t *testing.T) {
	timestamp := time.Date(2024, 3, 15, 14, 30, 45, 0, time.Local).Unix()
	result := Localtime(types.NewInt(timestamp), types.NewBool(true))

	if result.Type() != types.TypeArray {
		t.Errorf("localtime(, true) should return array")
	}

	arr := result.ToArray()

	sec, _ := arr.Get(types.NewString("tm_sec"))
	if sec.ToInt() != 45 {
		t.Errorf("localtime()['tm_sec'] = %v, want 45", sec.ToInt())
	}

	min, _ := arr.Get(types.NewString("tm_min"))
	if min.ToInt() != 30 {
		t.Errorf("localtime()['tm_min'] = %v, want 30", min.ToInt())
	}

	hour, _ := arr.Get(types.NewString("tm_hour"))
	if hour.ToInt() != 14 {
		t.Errorf("localtime()['tm_hour'] = %v, want 14", hour.ToInt())
	}
}

// ============================================================================
// Checkdate Tests
// ============================================================================

func TestCheckdateValid(t *testing.T) {
	tests := []struct {
		month int64
		day   int64
		year  int64
	}{
		{1, 1, 2024},     // January 1, 2024
		{12, 31, 2024},   // December 31, 2024
		{2, 29, 2024},    // Leap day in leap year
		{6, 15, 2000},    // Mid-year date
	}

	for _, tt := range tests {
		result := Checkdate(
			types.NewInt(tt.month),
			types.NewInt(tt.day),
			types.NewInt(tt.year),
		)
		if !result.ToBool() {
			t.Errorf("checkdate(%d, %d, %d) should be true", tt.month, tt.day, tt.year)
		}
	}
}

func TestCheckdateInvalid(t *testing.T) {
	tests := []struct {
		month int64
		day   int64
		year  int64
		reason string
	}{
		{0, 1, 2024, "month 0"},
		{13, 1, 2024, "month 13"},
		{1, 0, 2024, "day 0"},
		{1, 32, 2024, "day 32 in January"},
		{2, 30, 2024, "day 30 in February"},
		{2, 29, 2023, "Feb 29 in non-leap year"},
		{1, 1, 0, "year 0"},
		{1, 1, 40000, "year > 32767"},
	}

	for _, tt := range tests {
		result := Checkdate(
			types.NewInt(tt.month),
			types.NewInt(tt.day),
			types.NewInt(tt.year),
		)
		if result.ToBool() {
			t.Errorf("checkdate(%d, %d, %d) should be false (%s)", tt.month, tt.day, tt.year, tt.reason)
		}
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year     int
		month    int
		expected int
	}{
		{2024, 1, 31},  // January
		{2024, 2, 29},  // February (leap year)
		{2023, 2, 28},  // February (non-leap year)
		{2024, 4, 30},  // April
		{2024, 12, 31}, // December
	}

	for _, tt := range tests {
		tm := time.Date(tt.year, time.Month(tt.month), 1, 0, 0, 0, 0, time.UTC)
		result := daysInMonth(tm)
		if result != tt.expected {
			t.Errorf("daysInMonth(%d-%d) = %d, want %d", tt.year, tt.month, result, tt.expected)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2024, true},  // Divisible by 4
		{2023, false}, // Not divisible by 4
		{2000, true},  // Divisible by 400
		{1900, false}, // Divisible by 100 but not 400
		{2004, true},  // Divisible by 4
	}

	for _, tt := range tests {
		result := isLeapYear(tt.year)
		if result != tt.expected {
			t.Errorf("isLeapYear(%d) = %v, want %v", tt.year, result, tt.expected)
		}
	}
}
