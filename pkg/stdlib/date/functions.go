package date

import (
	"fmt"
	"strings"
	"time"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Time Functions
// ============================================================================

// Time returns current Unix timestamp
// time(): int
func Time() *types.Value {
	return types.NewInt(time.Now().Unix())
}

// Microtime returns current Unix timestamp with microseconds
// microtime(bool $as_float = false): string|float
func Microtime(args ...*types.Value) *types.Value {
	asFloat := false
	if len(args) > 0 {
		asFloat = args[0].ToBool()
	}

	now := time.Now()
	sec := now.Unix()
	usec := now.UnixMicro() % 1000000

	if asFloat {
		return types.NewFloat(float64(sec) + float64(usec)/1000000.0)
	}

	return types.NewString(fmt.Sprintf("0.%06d %d", usec, sec))
}

// ============================================================================
// Date Formatting Functions
// ============================================================================

// Date formats a Unix timestamp
// date(string $format, int $timestamp = null): string
func Date(format *types.Value, args ...*types.Value) *types.Value {
	var t time.Time
	if len(args) > 0 {
		timestamp := args[0].ToInt()
		t = time.Unix(timestamp, 0)
	} else {
		t = time.Now()
	}

	return types.NewString(formatDate(format.ToString(), t))
}

// Gmdate formats a GMT/UTC date/time
// gmdate(string $format, int $timestamp = null): string
func Gmdate(format *types.Value, args ...*types.Value) *types.Value {
	var t time.Time
	if len(args) > 0 {
		timestamp := args[0].ToInt()
		t = time.Unix(timestamp, 0).UTC()
	} else {
		t = time.Now().UTC()
	}

	return types.NewString(formatDate(format.ToString(), t))
}

// formatDate formats a time.Time according to PHP date format
func formatDate(format string, t time.Time) string {
	var result strings.Builder

	for i := 0; i < len(format); i++ {
		ch := format[i]

		switch ch {
		// Day
		case 'd': // Day of month, 2 digits with leading zeros
			result.WriteString(fmt.Sprintf("%02d", t.Day()))
		case 'D': // Textual day, 3 letters
			result.WriteString(t.Weekday().String()[:3])
		case 'j': // Day of month without leading zeros
			result.WriteString(fmt.Sprintf("%d", t.Day()))
		case 'l': // Full textual day
			result.WriteString(t.Weekday().String())
		case 'N': // ISO-8601 day of week (1=Monday, 7=Sunday)
			day := int(t.Weekday())
			if day == 0 {
				day = 7
			}
			result.WriteString(fmt.Sprintf("%d", day))
		case 'w': // Day of week (0=Sunday, 6=Saturday)
			result.WriteString(fmt.Sprintf("%d", int(t.Weekday())))
		case 'z': // Day of year (0-365)
			result.WriteString(fmt.Sprintf("%d", t.YearDay()-1))

		// Month
		case 'F': // Full textual month
			result.WriteString(t.Month().String())
		case 'm': // Month, 2 digits with leading zeros
			result.WriteString(fmt.Sprintf("%02d", int(t.Month())))
		case 'M': // Textual month, 3 letters
			result.WriteString(t.Month().String()[:3])
		case 'n': // Month without leading zeros
			result.WriteString(fmt.Sprintf("%d", int(t.Month())))
		case 't': // Number of days in month
			result.WriteString(fmt.Sprintf("%d", daysInMonth(t)))

		// Year
		case 'Y': // Full year, 4 digits
			result.WriteString(fmt.Sprintf("%04d", t.Year()))
		case 'y': // Year, 2 digits
			result.WriteString(fmt.Sprintf("%02d", t.Year()%100))
		case 'L': // Leap year (1 or 0)
			if isLeapYear(t.Year()) {
				result.WriteString("1")
			} else {
				result.WriteString("0")
			}

		// Time
		case 'a': // am or pm
			if t.Hour() < 12 {
				result.WriteString("am")
			} else {
				result.WriteString("pm")
			}
		case 'A': // AM or PM
			if t.Hour() < 12 {
				result.WriteString("AM")
			} else {
				result.WriteString("PM")
			}
		case 'g': // 12-hour format without leading zeros
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result.WriteString(fmt.Sprintf("%d", h))
		case 'G': // 24-hour format without leading zeros
			result.WriteString(fmt.Sprintf("%d", t.Hour()))
		case 'h': // 12-hour format with leading zeros
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result.WriteString(fmt.Sprintf("%02d", h))
		case 'H': // 24-hour format with leading zeros
			result.WriteString(fmt.Sprintf("%02d", t.Hour()))
		case 'i': // Minutes with leading zeros
			result.WriteString(fmt.Sprintf("%02d", t.Minute()))
		case 's': // Seconds with leading zeros
			result.WriteString(fmt.Sprintf("%02d", t.Second()))
		case 'u': // Microseconds
			result.WriteString(fmt.Sprintf("%06d", t.Nanosecond()/1000))

		// Timezone
		case 'e': // Timezone identifier
			zone, _ := t.Zone()
			result.WriteString(zone)
		case 'O': // Difference to GMT in hours (+0200)
			_, offset := t.Zone()
			hours := offset / 3600
			minutes := (offset % 3600) / 60
			result.WriteString(fmt.Sprintf("%+03d%02d", hours, abs(minutes)))
		case 'P': // Difference to GMT with colon (+02:00)
			_, offset := t.Zone()
			hours := offset / 3600
			minutes := (offset % 3600) / 60
			result.WriteString(fmt.Sprintf("%+03d:%02d", hours, abs(minutes)))
		case 'T': // Timezone abbreviation
			zone, _ := t.Zone()
			result.WriteString(zone)
		case 'Z': // Timezone offset in seconds
			_, offset := t.Zone()
			result.WriteString(fmt.Sprintf("%d", offset))

		// Full date/time
		case 'c': // ISO 8601 date
			result.WriteString(t.Format("2006-01-02T15:04:05-07:00"))
		case 'r': // RFC 2822 formatted date
			result.WriteString(t.Format(time.RFC1123Z))
		case 'U': // Unix timestamp
			result.WriteString(fmt.Sprintf("%d", t.Unix()))

		// Escape character
		case '\\':
			if i+1 < len(format) {
				i++
				result.WriteByte(format[i])
			}

		default:
			result.WriteByte(ch)
		}
	}

	return result.String()
}

// Helper functions
func daysInMonth(t time.Time) int {
	// Get first day of next month, then go back one day
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth.Day()
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ============================================================================
// Mktime and Strtotime
// ============================================================================

// Mktime returns Unix timestamp for a date
// mktime(int $hour, int $minute = 0, int $second = 0, int $month = 1, int $day = 1, int $year = 0): int|false
func Mktime(args ...*types.Value) *types.Value {
	now := time.Now()

	hour := 0
	minute := 0
	second := 0
	month := int(now.Month())
	day := now.Day()
	year := now.Year()

	if len(args) > 0 {
		hour = int(args[0].ToInt())
	}
	if len(args) > 1 {
		minute = int(args[1].ToInt())
	}
	if len(args) > 2 {
		second = int(args[2].ToInt())
	}
	if len(args) > 3 {
		month = int(args[3].ToInt())
	}
	if len(args) > 4 {
		day = int(args[4].ToInt())
	}
	if len(args) > 5 {
		year = int(args[5].ToInt())
	}

	// Handle 2-digit years
	if year >= 0 && year < 70 {
		year += 2000
	} else if year >= 70 && year < 100 {
		year += 1900
	}

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return types.NewInt(t.Unix())
}

// Gmmktime returns Unix timestamp for a GMT date
// gmmktime(int $hour, int $minute = 0, int $second = 0, int $month = 1, int $day = 1, int $year = 0): int|false
func Gmmktime(args ...*types.Value) *types.Value {
	now := time.Now().UTC()

	hour := 0
	minute := 0
	second := 0
	month := int(now.Month())
	day := now.Day()
	year := now.Year()

	if len(args) > 0 {
		hour = int(args[0].ToInt())
	}
	if len(args) > 1 {
		minute = int(args[1].ToInt())
	}
	if len(args) > 2 {
		second = int(args[2].ToInt())
	}
	if len(args) > 3 {
		month = int(args[3].ToInt())
	}
	if len(args) > 4 {
		day = int(args[4].ToInt())
	}
	if len(args) > 5 {
		year = int(args[5].ToInt())
	}

	// Handle 2-digit years
	if year >= 0 && year < 70 {
		year += 2000
	} else if year >= 70 && year < 100 {
		year += 1900
	}

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return types.NewInt(t.Unix())
}

// Strtotime parses an English textual datetime into a Unix timestamp
// strtotime(string $datetime, int $baseTimestamp = null): int|false
func Strtotime(datetime *types.Value, args ...*types.Value) *types.Value {
	dateStr := strings.TrimSpace(datetime.ToString())

	var baseTime time.Time
	if len(args) > 0 {
		baseTime = time.Unix(args[0].ToInt(), 0)
	} else {
		baseTime = time.Now()
	}

	// Try common formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return types.NewInt(t.Unix())
		}
	}

	// Handle relative times (simplified)
	lower := strings.ToLower(dateStr)

	switch lower {
	case "now":
		return types.NewInt(baseTime.Unix())
	case "today":
		t := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), 0, 0, 0, 0, baseTime.Location())
		return types.NewInt(t.Unix())
	case "tomorrow":
		t := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+1, 0, 0, 0, 0, baseTime.Location())
		return types.NewInt(t.Unix())
	case "yesterday":
		t := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()-1, 0, 0, 0, 0, baseTime.Location())
		return types.NewInt(t.Unix())
	}

	// Relative time with + or -
	if strings.HasPrefix(lower, "+") || strings.HasPrefix(lower, "-") {
		// Simplified relative time parsing
		if strings.Contains(lower, "day") {
			if strings.HasPrefix(lower, "+") {
				days := 1
				fmt.Sscanf(lower, "+%d", &days)
				return types.NewInt(baseTime.AddDate(0, 0, days).Unix())
			} else {
				days := 1
				fmt.Sscanf(lower, "-%d", &days)
				return types.NewInt(baseTime.AddDate(0, 0, -days).Unix())
			}
		}
	}

	// Failed to parse
	return types.NewBool(false)
}

// ============================================================================
// Getdate
// ============================================================================

// Getdate gets date/time information
// getdate(int $timestamp = null): array
func Getdate(args ...*types.Value) *types.Value {
	var t time.Time
	if len(args) > 0 {
		t = time.Unix(args[0].ToInt(), 0)
	} else {
		t = time.Now()
	}

	arr := types.NewEmptyArray()
	arr.Set(types.NewString("seconds"), types.NewInt(int64(t.Second())))
	arr.Set(types.NewString("minutes"), types.NewInt(int64(t.Minute())))
	arr.Set(types.NewString("hours"), types.NewInt(int64(t.Hour())))
	arr.Set(types.NewString("mday"), types.NewInt(int64(t.Day())))
	arr.Set(types.NewString("wday"), types.NewInt(int64(t.Weekday())))
	arr.Set(types.NewString("mon"), types.NewInt(int64(t.Month())))
	arr.Set(types.NewString("year"), types.NewInt(int64(t.Year())))
	arr.Set(types.NewString("yday"), types.NewInt(int64(t.YearDay()-1)))
	arr.Set(types.NewString("weekday"), types.NewString(t.Weekday().String()))
	arr.Set(types.NewString("month"), types.NewString(t.Month().String()))
	arr.Set(types.NewString("0"), types.NewInt(t.Unix()))

	return types.NewArray(arr)
}

// ============================================================================
// Localtime
// ============================================================================

// Localtime gets the local time
// localtime(int $timestamp = null, bool $associative = false): array
func Localtime(args ...*types.Value) *types.Value {
	var t time.Time
	if len(args) > 0 {
		t = time.Unix(args[0].ToInt(), 0)
	} else {
		t = time.Now()
	}

	associative := false
	if len(args) > 1 {
		associative = args[1].ToBool()
	}

	arr := types.NewEmptyArray()

	if associative {
		arr.Set(types.NewString("tm_sec"), types.NewInt(int64(t.Second())))
		arr.Set(types.NewString("tm_min"), types.NewInt(int64(t.Minute())))
		arr.Set(types.NewString("tm_hour"), types.NewInt(int64(t.Hour())))
		arr.Set(types.NewString("tm_mday"), types.NewInt(int64(t.Day())))
		arr.Set(types.NewString("tm_mon"), types.NewInt(int64(t.Month()-1))) // 0-11
		arr.Set(types.NewString("tm_year"), types.NewInt(int64(t.Year()-1900)))
		arr.Set(types.NewString("tm_wday"), types.NewInt(int64(t.Weekday())))
		arr.Set(types.NewString("tm_yday"), types.NewInt(int64(t.YearDay()-1)))
		arr.Set(types.NewString("tm_isdst"), types.NewInt(0)) // DST info not available
	} else {
		arr.Append(types.NewInt(int64(t.Second())))
		arr.Append(types.NewInt(int64(t.Minute())))
		arr.Append(types.NewInt(int64(t.Hour())))
		arr.Append(types.NewInt(int64(t.Day())))
		arr.Append(types.NewInt(int64(t.Month() - 1))) // 0-11
		arr.Append(types.NewInt(int64(t.Year() - 1900)))
		arr.Append(types.NewInt(int64(t.Weekday())))
		arr.Append(types.NewInt(int64(t.YearDay() - 1)))
		arr.Append(types.NewInt(0)) // DST info not available
	}

	return types.NewArray(arr)
}

// ============================================================================
// Checkdate
// ============================================================================

// Checkdate validates a Gregorian date
// checkdate(int $month, int $day, int $year): bool
func Checkdate(month *types.Value, day *types.Value, year *types.Value) *types.Value {
	m := int(month.ToInt())
	d := int(day.ToInt())
	y := int(year.ToInt())

	// Month must be 1-12
	if m < 1 || m > 12 {
		return types.NewBool(false)
	}

	// Year must be 1-32767
	if y < 1 || y > 32767 {
		return types.NewBool(false)
	}

	// Check day against month
	t := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
	maxDay := daysInMonth(t)

	if d < 1 || d > maxDay {
		return types.NewBool(false)
	}

	return types.NewBool(true)
}
