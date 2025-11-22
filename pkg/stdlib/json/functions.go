package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// JSON Constants (PHP JSON flags)
// ============================================================================

const (
	JSON_HEX_TAG             = 1 << 0  // 1
	JSON_HEX_AMP             = 1 << 1  // 2
	JSON_HEX_APOS            = 1 << 2  // 4
	JSON_HEX_QUOT            = 1 << 3  // 8
	JSON_FORCE_OBJECT        = 1 << 4  // 16
	JSON_NUMERIC_CHECK       = 1 << 5  // 32
	JSON_UNESCAPED_SLASHES   = 1 << 6  // 64
	JSON_PRETTY_PRINT        = 1 << 7  // 128
	JSON_UNESCAPED_UNICODE   = 1 << 8  // 256
	JSON_PARTIAL_OUTPUT_ON_ERROR = 1 << 9  // 512
	JSON_PRESERVE_ZERO_FRACTION = 1 << 10 // 1024
	JSON_UNESCAPED_LINE_TERMINATORS = 1 << 11 // 2048
	JSON_OBJECT_AS_ARRAY     = 1 << 0  // For decode
	JSON_BIGINT_AS_STRING    = 1 << 1  // For decode
	JSON_THROW_ON_ERROR      = 1 << 22 // 4194304
)

// JSON error constants
const (
	JSON_ERROR_NONE = iota
	JSON_ERROR_DEPTH
	JSON_ERROR_STATE_MISMATCH
	JSON_ERROR_CTRL_CHAR
	JSON_ERROR_SYNTAX
	JSON_ERROR_UTF8
	JSON_ERROR_RECURSION
	JSON_ERROR_INF_OR_NAN
	JSON_ERROR_UNSUPPORTED_TYPE
)

var lastJsonError = JSON_ERROR_NONE

// ============================================================================
// JSON Encode
// ============================================================================

// JsonEncode returns the JSON representation of a value
// json_encode(mixed $value, int $flags = 0, int $depth = 512): string|false
func JsonEncode(value *types.Value, flags ...*types.Value) *types.Value {
	options := 0
	depth := 512

	if len(flags) > 0 && flags[0] != nil {
		options = int(flags[0].ToInt())
	}

	if len(flags) > 1 && flags[1] != nil {
		depth = int(flags[1].ToInt())
	}

	// Convert PHP value to Go interface for JSON encoding
	result, err := encodeValue(value, options, depth, 0)
	if err != nil {
		lastJsonError = JSON_ERROR_SYNTAX
		return types.NewBool(false)
	}

	lastJsonError = JSON_ERROR_NONE
	return types.NewString(result)
}

// encodeValue recursively encodes a PHP value to JSON
func encodeValue(value *types.Value, options int, maxDepth int, currentDepth int) (string, error) {
	if currentDepth > maxDepth {
		return "", fmt.Errorf("maximum depth exceeded")
	}

	switch value.Type() {
	case types.TypeNull:
		return "null", nil

	case types.TypeBool:
		if value.ToBool() {
			return "true", nil
		}
		return "false", nil

	case types.TypeInt:
		return strconv.FormatInt(value.ToInt(), 10), nil

	case types.TypeFloat:
		f := value.ToFloat()
		// Check for infinity or NaN
		if !isFinite(f) {
			if options&JSON_PARTIAL_OUTPUT_ON_ERROR != 0 {
				return "0", nil
			}
			return "", fmt.Errorf("inf or nan")
		}

		str := strconv.FormatFloat(f, 'f', -1, 64)
		// JSON_PRESERVE_ZERO_FRACTION: ensure .0 for whole numbers
		if options&JSON_PRESERVE_ZERO_FRACTION != 0 {
			if !strings.Contains(str, ".") {
				str += ".0"
			}
		}
		return str, nil

	case types.TypeString:
		return encodeString(value.ToString(), options), nil

	case types.TypeArray:
		return encodeArray(value.ToArray(), options, maxDepth, currentDepth)

	case types.TypeObject:
		return encodeObject(value.ToObject(), options, maxDepth, currentDepth)

	default:
		return "", fmt.Errorf("unsupported type")
	}
}

// encodeString encodes a string with proper JSON escaping
func encodeString(s string, options int) string {
	var result strings.Builder
	result.WriteByte('"')

	for _, r := range s {
		switch r {
		case '"':
			if options&JSON_HEX_QUOT != 0 {
				result.WriteString("\\u0022")
			} else {
				result.WriteString("\\\"")
			}
		case '\\':
			result.WriteString("\\\\")
		case '/':
			if options&JSON_UNESCAPED_SLASHES != 0 {
				result.WriteRune(r)
			} else {
				result.WriteString("\\/")
			}
		case '\b':
			result.WriteString("\\b")
		case '\f':
			result.WriteString("\\f")
		case '\n':
			result.WriteString("\\n")
		case '\r':
			result.WriteString("\\r")
		case '\t':
			result.WriteString("\\t")
		case '\'':
			if options&JSON_HEX_APOS != 0 {
				result.WriteString("\\u0027")
			} else {
				result.WriteRune(r)
			}
		case '&':
			if options&JSON_HEX_AMP != 0 {
				result.WriteString("\\u0026")
			} else {
				result.WriteRune(r)
			}
		case '<':
			if options&JSON_HEX_TAG != 0 {
				result.WriteString("\\u003C")
			} else {
				result.WriteRune(r)
			}
		case '>':
			if options&JSON_HEX_TAG != 0 {
				result.WriteString("\\u003E")
			} else {
				result.WriteRune(r)
			}
		default:
			if r < 0x20 {
				// Control character
				result.WriteString(fmt.Sprintf("\\u%04x", r))
			} else if r > 0x7F && options&JSON_UNESCAPED_UNICODE == 0 {
				// Non-ASCII - escape unless flag set
				result.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				result.WriteRune(r)
			}
		}
	}

	result.WriteByte('"')
	return result.String()
}

// encodeArray encodes a PHP array to JSON
func encodeArray(arr *types.Array, options int, maxDepth int, currentDepth int) (string, error) {
	if arr.IsEmpty() {
		if options&JSON_FORCE_OBJECT != 0 {
			return "{}", nil
		}
		return "[]", nil
	}

	// Determine if array is sequential (list) or associative (object)
	isSequential := true
	expectedIndex := int64(0)

	arr.Each(func(key, _ *types.Value) bool {
		if key.Type() != types.TypeInt || key.ToInt() != expectedIndex {
			isSequential = false
			return false
		}
		expectedIndex++
		return true
	})

	if options&JSON_FORCE_OBJECT != 0 {
		isSequential = false
	}

	var result strings.Builder
	first := true

	if isSequential {
		result.WriteByte('[')
		arr.Each(func(_, value *types.Value) bool {
			if !first {
				result.WriteByte(',')
			}
			first = false

			encoded, err := encodeValue(value, options, maxDepth, currentDepth+1)
			if err != nil {
				return false
			}
			result.WriteString(encoded)
			return true
		})
		result.WriteByte(']')
	} else {
		result.WriteByte('{')
		arr.Each(func(key, value *types.Value) bool {
			if !first {
				result.WriteByte(',')
			}
			first = false

			// Encode key as string
			keyStr := key.ToString()
			result.WriteString(encodeString(keyStr, options))
			result.WriteByte(':')

			encoded, err := encodeValue(value, options, maxDepth, currentDepth+1)
			if err != nil {
				return false
			}
			result.WriteString(encoded)
			return true
		})
		result.WriteByte('}')
	}

	return result.String(), nil
}

// encodeObject encodes a PHP object to JSON
func encodeObject(obj *types.Object, options int, maxDepth int, currentDepth int) (string, error) {
	var result strings.Builder
	result.WriteByte('{')

	first := true
	for name, _ := range obj.Properties {
		if !first {
			result.WriteByte(',')
		}
		first = false

		// Property name
		result.WriteString(encodeString(name, options))
		result.WriteByte(':')

		// Property value (pass nil for access context to get any property)
		value, _ := obj.GetProperty(name, nil)
		encoded, err := encodeValue(value, options, maxDepth, currentDepth+1)
		if err != nil {
			return "", err
		}
		result.WriteString(encoded)
	}

	result.WriteByte('}')
	return result.String(), nil
}

// ============================================================================
// JSON Decode
// ============================================================================

// JsonDecode decodes a JSON string
// json_decode(string $json, bool $associative = false, int $depth = 512, int $flags = 0): mixed
func JsonDecode(jsonStr *types.Value, args ...*types.Value) *types.Value {
	associative := false
	depth := 512
	flags := 0

	if len(args) > 0 && args[0] != nil {
		associative = args[0].ToBool()
	}

	if len(args) > 1 && args[1] != nil {
		depth = int(args[1].ToInt())
	}

	if len(args) > 2 && args[2] != nil {
		flags = int(args[2].ToInt())
	}

	str := jsonStr.ToString()

	// Use Go's JSON decoder
	var result interface{}
	decoder := json.NewDecoder(strings.NewReader(str))

	if err := decoder.Decode(&result); err != nil {
		lastJsonError = JSON_ERROR_SYNTAX
		return types.NewNull()
	}

	lastJsonError = JSON_ERROR_NONE
	return convertFromJSON(result, associative, flags, 0, depth)
}

// convertFromJSON converts Go JSON types to PHP types
func convertFromJSON(val interface{}, associative bool, flags int, currentDepth int, maxDepth int) *types.Value {
	if currentDepth > maxDepth {
		lastJsonError = JSON_ERROR_DEPTH
		return types.NewNull()
	}

	switch v := val.(type) {
	case nil:
		return types.NewNull()

	case bool:
		return types.NewBool(v)

	case float64:
		// JSON numbers are always floats in Go
		// Check if it's actually an integer
		if v == float64(int64(v)) && flags&JSON_BIGINT_AS_STRING == 0 {
			return types.NewInt(int64(v))
		}
		return types.NewFloat(v)

	case string:
		return types.NewString(v)

	case []interface{}:
		// JSON array
		arr := types.NewEmptyArray()
		for _, item := range v {
			arr.Append(convertFromJSON(item, associative, flags, currentDepth+1, maxDepth))
		}
		return types.NewArray(arr)

	case map[string]interface{}:
		// JSON object
		if associative || flags&JSON_OBJECT_AS_ARRAY != 0 {
			// Return as associative array
			arr := types.NewEmptyArray()
			for key, value := range v {
				arr.Set(
					types.NewString(key),
					convertFromJSON(value, associative, flags, currentDepth+1, maxDepth),
				)
			}
			return types.NewArray(arr)
		} else {
			// Return as object (stdClass)
			class := types.NewClassEntry("stdClass")
			obj := types.NewObjectFromClass(class)

			for key, value := range v {
				// Add property to object
				propDef := &types.PropertyDef{
					Name:       key,
					Visibility: types.VisibilityPublic,
					IsStatic:   false,
					Type:       "",
					HasDefault: false,
					Default:    nil,
					IsReadOnly: false,
				}
				obj.ClassEntry.Properties[key] = propDef

				prop := &types.Property{
					Value:      convertFromJSON(value, associative, flags, currentDepth+1, maxDepth),
					Visibility: types.VisibilityPublic,
					IsStatic:   false,
					Type:       "",
					HasDefault: false,
					Default:    nil,
					IsReadOnly: false,
				}
				obj.Properties[key] = prop
			}

			return types.NewObject(obj)
		}

	default:
		return types.NewNull()
	}
}

// ============================================================================
// JSON Error Handling
// ============================================================================

// JsonLastError returns the last error occurred
// json_last_error(): int
func JsonLastError() *types.Value {
	return types.NewInt(int64(lastJsonError))
}

// JsonLastErrorMsg returns the error message of the last json_encode() or json_decode() call
// json_last_error_msg(): string
func JsonLastErrorMsg() *types.Value {
	switch lastJsonError {
	case JSON_ERROR_NONE:
		return types.NewString("No error")
	case JSON_ERROR_DEPTH:
		return types.NewString("Maximum stack depth exceeded")
	case JSON_ERROR_STATE_MISMATCH:
		return types.NewString("State mismatch (invalid or malformed JSON)")
	case JSON_ERROR_CTRL_CHAR:
		return types.NewString("Control character error, possibly incorrectly encoded")
	case JSON_ERROR_SYNTAX:
		return types.NewString("Syntax error")
	case JSON_ERROR_UTF8:
		return types.NewString("Malformed UTF-8 characters, possibly incorrectly encoded")
	case JSON_ERROR_RECURSION:
		return types.NewString("Recursion detected")
	case JSON_ERROR_INF_OR_NAN:
		return types.NewString("Inf and NaN cannot be JSON encoded")
	case JSON_ERROR_UNSUPPORTED_TYPE:
		return types.NewString("Type is not supported")
	default:
		return types.NewString("Unknown error")
	}
}

// ============================================================================
// Utility Functions
// ============================================================================

// isFinite checks if a float is finite (not NaN or Inf)
func isFinite(f float64) bool {
	// Check for NaN
	if f != f {
		return false
	}
	// Check for infinity
	if f > 1.7976931348623157e+308 || f < -1.7976931348623157e+308 {
		return false
	}
	return true
}
