package filter

import (
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Filter Constants
// ============================================================================

const (
	// Validate filters
	FILTER_VALIDATE_BOOLEAN = 258
	FILTER_VALIDATE_EMAIL   = 274
	FILTER_VALIDATE_FLOAT   = 259
	FILTER_VALIDATE_INT     = 257
	FILTER_VALIDATE_IP      = 275
	FILTER_VALIDATE_MAC     = 276
	FILTER_VALIDATE_REGEXP  = 272
	FILTER_VALIDATE_URL     = 273
	FILTER_VALIDATE_DOMAIN  = 277

	// Sanitize filters
	FILTER_SANITIZE_EMAIL           = 517
	FILTER_SANITIZE_ENCODED         = 514
	FILTER_SANITIZE_NUMBER_FLOAT    = 520
	FILTER_SANITIZE_NUMBER_INT      = 519
	FILTER_SANITIZE_SPECIAL_CHARS   = 515
	FILTER_SANITIZE_STRING          = 513
	FILTER_SANITIZE_STRIPPED        = 513
	FILTER_SANITIZE_URL             = 518
	FILTER_SANITIZE_ADD_SLASHES     = 523

	// Other filters
	FILTER_UNSAFE_RAW = 516
	FILTER_CALLBACK   = 1024

	// Flags
	FILTER_FLAG_ALLOW_OCTAL       = 1
	FILTER_FLAG_ALLOW_HEX         = 2
	FILTER_FLAG_STRIP_LOW         = 4
	FILTER_FLAG_STRIP_HIGH        = 8
	FILTER_FLAG_ENCODE_LOW        = 16
	FILTER_FLAG_ENCODE_HIGH       = 32
	FILTER_FLAG_ENCODE_AMP        = 64
	FILTER_FLAG_NO_ENCODE_QUOTES  = 128
	FILTER_FLAG_EMPTY_STRING_NULL = 256
	FILTER_FLAG_ALLOW_FRACTION    = 4096
	FILTER_FLAG_ALLOW_THOUSAND    = 8192
	FILTER_FLAG_ALLOW_SCIENTIFIC  = 16384
	FILTER_FLAG_STRIP_BACKTICK    = 512

	// IP flags
	FILTER_FLAG_IPV4          = 1048576
	FILTER_FLAG_IPV6          = 2097152
	FILTER_FLAG_NO_PRIV_RANGE = 8388608
	FILTER_FLAG_NO_RES_RANGE  = 4194304

	// URL flags
	FILTER_FLAG_SCHEME_REQUIRED = 65536
	FILTER_FLAG_HOST_REQUIRED   = 131072
	FILTER_FLAG_PATH_REQUIRED   = 262144
	FILTER_FLAG_QUERY_REQUIRED  = 524288

	// Regex flags
	FILTER_REQUIRE_SCALAR = 33554432
	FILTER_REQUIRE_ARRAY  = 16777216

	// Other
	FILTER_NULL_ON_FAILURE = 134217728
)

// ============================================================================
// Filter Var
// ============================================================================

// FilterVar filters a variable with a specified filter
// filter_var(mixed $value, int $filter = FILTER_DEFAULT, array|int $options = []): mixed
func FilterVar(value *types.Value, args ...*types.Value) *types.Value {
	filter := FILTER_UNSAFE_RAW
	flags := 0

	if len(args) > 0 && args[0] != nil {
		filter = int(args[0].ToInt())
	}

	if len(args) > 1 && args[1] != nil {
		// Can be int (flags) or array (options)
		if args[1].Type() == types.TypeInt {
			flags = int(args[1].ToInt())
		}
		// TODO: Handle array options
	}

	return applyFilter(value, filter, flags)
}

// applyFilter applies the specified filter to a value
func applyFilter(value *types.Value, filter int, flags int) *types.Value {
	str := value.ToString()

	switch filter {
	// Validation filters
	case FILTER_VALIDATE_BOOLEAN:
		return validateBoolean(str, flags)
	case FILTER_VALIDATE_EMAIL:
		return validateEmail(str)
	case FILTER_VALIDATE_FLOAT:
		return validateFloat(str, flags)
	case FILTER_VALIDATE_INT:
		return validateInt(str, flags)
	case FILTER_VALIDATE_IP:
		return validateIP(str, flags)
	case FILTER_VALIDATE_MAC:
		return validateMAC(str)
	case FILTER_VALIDATE_URL:
		return validateURL(str, flags)
	case FILTER_VALIDATE_DOMAIN:
		return validateDomain(str)

	// Sanitize filters
	case FILTER_SANITIZE_EMAIL:
		return sanitizeEmail(str)
	case FILTER_SANITIZE_NUMBER_FLOAT:
		return sanitizeNumberFloat(str, flags)
	case FILTER_SANITIZE_NUMBER_INT:
		return sanitizeNumberInt(str)
	case FILTER_SANITIZE_STRING:
		return sanitizeString(str, flags)
	case FILTER_SANITIZE_URL:
		return sanitizeURL(str)

	// Raw/default
	case FILTER_UNSAFE_RAW:
		return value

	default:
		return types.NewBool(false)
	}
}

// ============================================================================
// Validation Functions
// ============================================================================

func validateBoolean(str string, flags int) *types.Value {
	str = strings.TrimSpace(strings.ToLower(str))

	truthy := []string{"1", "true", "on", "yes"}
	falsy := []string{"0", "false", "off", "no", ""}

	for _, t := range truthy {
		if str == t {
			return types.NewBool(true)
		}
	}

	for _, f := range falsy {
		if str == f {
			return types.NewBool(true)
		}
	}

	return types.NewBool(false)
}

func validateEmail(str string) *types.Value {
	// Basic email validation
	_, err := mail.ParseAddress(str)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewString(str)
}

func validateFloat(str string, flags int) *types.Value {
	str = strings.TrimSpace(str)

	// Handle thousand separator if allowed
	if flags&FILTER_FLAG_ALLOW_THOUSAND != 0 {
		str = strings.ReplaceAll(str, ",", "")
	}

	// Try to parse as float
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(str)
}

func validateInt(str string, flags int) *types.Value {
	str = strings.TrimSpace(str)

	// Check for different bases
	base := 10

	if flags&FILTER_FLAG_ALLOW_HEX != 0 {
		if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
			base = 16
			str = str[2:]
		}
	}

	if flags&FILTER_FLAG_ALLOW_OCTAL != 0 {
		if strings.HasPrefix(str, "0") && len(str) > 1 {
			base = 8
			str = str[1:]
		}
	}

	// Try to parse as int
	_, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(str)
}

func validateIP(str string, flags int) *types.Value {
	ip := net.ParseIP(str)
	if ip == nil {
		return types.NewBool(false)
	}

	// Check IPv4/IPv6 flags
	if flags&FILTER_FLAG_IPV4 != 0 {
		if ip.To4() == nil {
			return types.NewBool(false)
		}
	}

	if flags&FILTER_FLAG_IPV6 != 0 {
		if ip.To4() != nil {
			return types.NewBool(false)
		}
	}

	// Check private range
	if flags&FILTER_FLAG_NO_PRIV_RANGE != 0 {
		if ip.IsPrivate() {
			return types.NewBool(false)
		}
	}

	// Check reserved range
	if flags&FILTER_FLAG_NO_RES_RANGE != 0 {
		if ip.IsLoopback() || ip.IsMulticast() {
			return types.NewBool(false)
		}
	}

	return types.NewString(str)
}

func validateMAC(str string) *types.Value {
	// MAC address validation
	_, err := net.ParseMAC(str)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewString(str)
}

func validateURL(str string, flags int) *types.Value {
	u, err := url.Parse(str)
	if err != nil {
		return types.NewBool(false)
	}

	// URL must have a scheme and host by default
	if u.Scheme == "" || u.Host == "" {
		return types.NewBool(false)
	}

	// Check required components
	if flags&FILTER_FLAG_SCHEME_REQUIRED != 0 && u.Scheme == "" {
		return types.NewBool(false)
	}

	if flags&FILTER_FLAG_HOST_REQUIRED != 0 && u.Host == "" {
		return types.NewBool(false)
	}

	if flags&FILTER_FLAG_PATH_REQUIRED != 0 && u.Path == "" {
		return types.NewBool(false)
	}

	if flags&FILTER_FLAG_QUERY_REQUIRED != 0 && u.RawQuery == "" {
		return types.NewBool(false)
	}

	return types.NewString(str)
}

func validateDomain(str string) *types.Value {
	// Basic domain validation
	if str == "" {
		return types.NewBool(false)
	}

	// Simple regex for domain validation
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(str) {
		return types.NewBool(false)
	}

	return types.NewString(str)
}

// ============================================================================
// Sanitization Functions
// ============================================================================

func sanitizeEmail(str string) *types.Value {
	// Remove all characters except letters, digits and !#$%&'*+-=?^_`{|}~@.[]
	// Note: '<' and '>' are specifically NOT allowed in email addresses
	var result strings.Builder
	for _, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '!' || r == '#' || r == '$' || r == '%' || r == '&' || r == '\'' ||
			r == '*' || r == '+' || r == '-' || r == '=' || r == '?' || r == '^' ||
			r == '_' || r == '`' || r == '{' || r == '|' || r == '}' || r == '~' ||
			r == '@' || r == '.' {
			result.WriteRune(r)
		}
		// Note: '<', '>', '[', ']', and other special chars are filtered out
	}
	return types.NewString(result.String())
}

func sanitizeNumberFloat(str string, flags int) *types.Value {
	// Keep only digits, +, -, ., and optionally , and e/E
	var result strings.Builder
	for _, r := range str {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		} else if r == '+' || r == '-' || r == '.' {
			result.WriteRune(r)
		} else if flags&FILTER_FLAG_ALLOW_FRACTION != 0 && r == ',' {
			result.WriteRune(r)
		} else if flags&FILTER_FLAG_ALLOW_SCIENTIFIC != 0 && (r == 'e' || r == 'E') {
			result.WriteRune(r)
		}
	}
	return types.NewString(result.String())
}

func sanitizeNumberInt(str string) *types.Value {
	// Keep only digits, + and -
	var result strings.Builder
	for _, r := range str {
		if (r >= '0' && r <= '9') || r == '+' || r == '-' {
			result.WriteRune(r)
		}
	}
	return types.NewString(result.String())
}

func sanitizeString(str string, flags int) *types.Value {
	// Strip tags and encode special characters
	// For now, a simplified version
	var result strings.Builder
	inTag := false

	for _, r := range str {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			// Strip low ASCII if flag set
			if flags&FILTER_FLAG_STRIP_LOW != 0 && r < 32 {
				continue
			}
			// Strip high ASCII if flag set
			if flags&FILTER_FLAG_STRIP_HIGH != 0 && r >= 128 {
				continue
			}
			result.WriteRune(r)
		}
	}

	return types.NewString(result.String())
}

func sanitizeURL(str string) *types.Value {
	// Remove all characters except letters, digits and $-_.+!*'(),{}|\\^~[]`<>#%";/?:@&=
	var result strings.Builder
	for _, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '$' || r == '-' || r == '_' || r == '.' || r == '+' ||
			r == '!' || r == '*' || r == '\'' || r == '(' || r == ')' ||
			r == ',' || r == '{' || r == '}' || r == '|' || r == '\\' ||
			r == '^' || r == '~' || r == '[' || r == ']' || r == '`' ||
			r == '<' || r == '>' || r == '#' || r == '%' || r == '"' ||
			r == ';' || r == '/' || r == '?' || r == ':' || r == '@' ||
			r == '&' || r == '=' {
			result.WriteRune(r)
		}
	}
	return types.NewString(result.String())
}

// ============================================================================
// Filter Var Array
// ============================================================================

// FilterVarArray filters multiple variables
// filter_var_array(array $array, array|int $definition = FILTER_DEFAULT): array|false
func FilterVarArray(arr *types.Value, args ...*types.Value) *types.Value {
	if arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	filter := FILTER_UNSAFE_RAW
	if len(args) > 0 && args[0] != nil {
		filter = int(args[0].ToInt())
	}

	inputArr := arr.ToArray()
	result := types.NewEmptyArray()

	inputArr.Each(func(key, value *types.Value) bool {
		filtered := applyFilter(value, filter, 0)
		result.Set(key, filtered)
		return true
	})

	return types.NewArray(result)
}
