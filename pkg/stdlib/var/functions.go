package varfuncs

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// outputWriter is the global output writer for debugging functions
// Defaults to os.Stdout but can be set by the VM for output capturing
var outputWriter io.Writer = os.Stdout

// SetOutputWriter sets the global output writer for debugging functions
func SetOutputWriter(w io.Writer) {
	outputWriter = w
}

// ============================================================================
// Variable Dumping Functions
// ============================================================================

// VarDump prints structured information about one or more variables
// var_dump(mixed ...$vars): void
func VarDump(values ...*types.Value) *types.Value {
	for _, val := range values {
		dumpValue(val, 0, make(map[interface{}]bool))
	}
	return types.NewNull()
}

// dumpValue recursively dumps a value with indentation
func dumpValue(val *types.Value, indent int, visited map[interface{}]bool) {
	prefix := strings.Repeat("  ", indent)

	switch val.Type() {
	case types.TypeNull:
		fmt.Fprintf(outputWriter, "%sNULL\n", prefix)

	case types.TypeBool:
		if val.ToBool() {
			fmt.Fprintf(outputWriter, "%sbool(true)\n", prefix)
		} else {
			fmt.Fprintf(outputWriter, "%sbool(false)\n", prefix)
		}

	case types.TypeInt:
		fmt.Fprintf(outputWriter, "%sint(%d)\n", prefix, val.ToInt())

	case types.TypeFloat:
		fmt.Fprintf(outputWriter, "%sfloat(%g)\n", prefix, val.ToFloat())

	case types.TypeString:
		str := val.ToString()
		fmt.Fprintf(outputWriter, "%sstring(%d) \"%s\"\n", prefix, len(str), str)

	case types.TypeArray:
		arr := val.ToArray()
		fmt.Fprintf(outputWriter, "%sarray(%d) {\n", prefix, arr.Len())

		// Collect and sort keys for consistent output
		var keys []*types.Value
		arr.Each(func(key, _ *types.Value) bool {
			keys = append(keys, key)
			return true
		})

		// Dump each element
		idx := 0
		arr.Each(func(key, value *types.Value) bool {
			// Print key
			if key.Type() == types.TypeInt {
				fmt.Fprintf(outputWriter, "%s  [%d]=>\n", prefix, key.ToInt())
			} else {
				fmt.Fprintf(outputWriter, "%s  [\"%s\"]=>\n", prefix, key.ToString())
			}

			// Check for circular reference
			if value.Type() == types.TypeArray || value.Type() == types.TypeObject {
				// Use pointer as identifier for visited tracking
				ptr := fmt.Sprintf("%p", value)
				if visited[ptr] {
					fmt.Fprintf(outputWriter, "%s    *RECURSION*\n", prefix)
					return true
				}
				visited[ptr] = true
			}

			// Print value
			dumpValue(value, indent+1, visited)

			idx++
			return true
		})

		fmt.Fprintf(outputWriter, "%s}\n", prefix)

	case types.TypeObject:
		obj := val.ToObject()
		className := obj.ClassName

		// Check for circular reference
		ptr := fmt.Sprintf("%p", obj)
		if visited[ptr] {
			fmt.Fprintf(outputWriter, "%sobject(%s)#%d *RECURSION*\n", prefix, className, obj.ObjectID)
			return
		}
		visited[ptr] = true

		// Get property count
		propCount := len(obj.ClassEntry.Properties)
		fmt.Fprintf(outputWriter, "%sobject(%s)#%d (%d) {\n", prefix, className, obj.ObjectID, propCount)

		// Dump properties - access directly from instance to bypass visibility checks
		for name, propDef := range obj.ClassEntry.Properties {
			visibility := propDef.Visibility.String()

			// Get property value directly from instance (var_dump shows all properties)
			var value *types.Value
			if prop, exists := obj.Properties[name]; exists && prop.Value != nil {
				value = prop.Value
			} else if propDef.HasDefault {
				value = propDef.Default
			} else {
				value = types.NewNull()
			}
			fmt.Fprintf(outputWriter, "%s  [\"%s\":%s]=>\n", prefix, name, visibility)
			dumpValue(value, indent+1, visited)
		}

		fmt.Fprintf(outputWriter, "%s}\n", prefix)

	case types.TypeResource:
		res := val.ToResource()
		if res.IsValid() {
			fmt.Fprintf(outputWriter, "%sresource(%d) of type (%s)\n", prefix, res.ID(), res.Type())
		} else {
			fmt.Fprintf(outputWriter, "%sresource(%d) of type (%s) (closed)\n", prefix, res.ID(), res.Type())
		}

	default:
		fmt.Fprintf(outputWriter, "%sunknown type\n", prefix)
	}
}

// PrintR prints human-readable information about a variable
// print_r(mixed $value, bool $return = false): mixed
func PrintR(val *types.Value, returnOutput ...*types.Value) *types.Value {
	shouldReturn := false
	if len(returnOutput) > 0 && returnOutput[0] != nil {
		shouldReturn = returnOutput[0].ToBool()
	}

	var output strings.Builder
	printValue(&output, val, 0, make(map[interface{}]bool))

	result := output.String()
	if shouldReturn {
		return types.NewString(result)
	}

	fmt.Fprint(outputWriter, result)
	return types.NewBool(true)
}

// printValue recursively prints a value in print_r format
func printValue(out *strings.Builder, val *types.Value, indent int, visited map[interface{}]bool) {
	prefix := strings.Repeat("    ", indent)

	switch val.Type() {
	case types.TypeNull:
		// print_r doesn't show anything for null

	case types.TypeBool:
		if val.ToBool() {
			out.WriteString("1")
		} else {
			// print_r shows nothing for false (PHP quirk)
		}

	case types.TypeInt:
		out.WriteString(fmt.Sprintf("%d", val.ToInt()))

	case types.TypeFloat:
		out.WriteString(fmt.Sprintf("%g", val.ToFloat()))

	case types.TypeString:
		out.WriteString(val.ToString())

	case types.TypeArray:
		arr := val.ToArray()

		// Check for circular reference
		ptr := fmt.Sprintf("%p", arr)
		if visited[ptr] {
			out.WriteString(" *RECURSION*")
			return
		}
		visited[ptr] = true

		out.WriteString("Array\n")
		out.WriteString(prefix + "(\n")

		arr.Each(func(key, value *types.Value) bool {
			// Print key
			if key.Type() == types.TypeInt {
				out.WriteString(fmt.Sprintf("%s    [%d] => ", prefix, key.ToInt()))
			} else {
				out.WriteString(fmt.Sprintf("%s    [%s] => ", prefix, key.ToString()))
			}

			// Print value
			if value.Type() == types.TypeArray || value.Type() == types.TypeObject {
				printValue(out, value, indent+1, visited)
			} else {
				printValue(out, value, 0, visited)
				out.WriteString("\n")
			}

			return true
		})

		out.WriteString(prefix + ")\n")

	case types.TypeObject:
		obj := val.ToObject()
		className := obj.ClassName

		// Check for circular reference
		ptr := fmt.Sprintf("%p", obj)
		if visited[ptr] {
			out.WriteString(" *RECURSION*")
			return
		}
		visited[ptr] = true

		out.WriteString(fmt.Sprintf("%s Object\n", className))
		out.WriteString(prefix + "(\n")

		// Print properties - access directly from instance to bypass visibility checks
		for name, propDef := range obj.ClassEntry.Properties {
			var value *types.Value
			if prop, exists := obj.Properties[name]; exists && prop.Value != nil {
				value = prop.Value
			} else if propDef.HasDefault {
				value = propDef.Default
			} else {
				value = types.NewNull()
			}
			out.WriteString(fmt.Sprintf("%s    [%s] => ", prefix, name))

			if value.Type() == types.TypeArray || value.Type() == types.TypeObject {
				printValue(out, value, indent+1, visited)
			} else {
				printValue(out, value, 0, visited)
				out.WriteString("\n")
			}
		}

		out.WriteString(prefix + ")\n")

	case types.TypeResource:
		res := val.ToResource()
		out.WriteString(fmt.Sprintf("Resource id #%d", res.ID()))

	default:
		out.WriteString("unknown")
	}
}

// VarExport outputs or returns a parsable string representation of a variable
// var_export(mixed $value, bool $return = false): mixed
func VarExport(val *types.Value, returnOutput ...*types.Value) *types.Value {
	shouldReturn := false
	if len(returnOutput) > 0 && returnOutput[0] != nil {
		shouldReturn = returnOutput[0].ToBool()
	}

	var output strings.Builder
	exportValue(&output, val, 0, make(map[interface{}]bool))

	result := output.String()
	if shouldReturn {
		return types.NewString(result)
	}

	fmt.Print(result)
	return types.NewNull()
}

// exportValue recursively exports a value in valid PHP code format
func exportValue(out *strings.Builder, val *types.Value, indent int, visited map[interface{}]bool) {
	prefix := strings.Repeat("  ", indent)

	switch val.Type() {
	case types.TypeNull:
		out.WriteString("NULL")

	case types.TypeBool:
		if val.ToBool() {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}

	case types.TypeInt:
		out.WriteString(fmt.Sprintf("%d", val.ToInt()))

	case types.TypeFloat:
		f := val.ToFloat()
		// Format float to ensure it always has decimal point
		str := fmt.Sprintf("%g", f)
		if !strings.Contains(str, ".") && !strings.Contains(str, "e") {
			str += ".0"
		}
		out.WriteString(str)

	case types.TypeString:
		// Escape string for PHP
		str := val.ToString()
		escaped := strings.ReplaceAll(str, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "\\'")
		out.WriteString(fmt.Sprintf("'%s'", escaped))

	case types.TypeArray:
		arr := val.ToArray()

		// Check for circular reference
		ptr := fmt.Sprintf("%p", arr)
		if visited[ptr] {
			out.WriteString("/* RECURSION */")
			return
		}
		visited[ptr] = true

		out.WriteString("array (\n")

		// Collect keys to sort them (for consistent output)
		type keyValuePair struct {
			key   *types.Value
			value *types.Value
		}
		var pairs []keyValuePair

		arr.Each(func(key, value *types.Value) bool {
			pairs = append(pairs, keyValuePair{key, value})
			return true
		})

		// Sort by key if all integer keys
		allInt := true
		for _, pair := range pairs {
			if pair.key.Type() != types.TypeInt {
				allInt = false
				break
			}
		}

		if allInt {
			sort.Slice(pairs, func(i, j int) bool {
				return pairs[i].key.ToInt() < pairs[j].key.ToInt()
			})
		}

		for _, pair := range pairs {
			out.WriteString(prefix + "  ")

			// Print key
			if pair.key.Type() == types.TypeInt {
				out.WriteString(fmt.Sprintf("%d => ", pair.key.ToInt()))
			} else {
				out.WriteString(fmt.Sprintf("'%s' => ", pair.key.ToString()))
			}

			// Print value
			exportValue(out, pair.value, indent+1, visited)
			out.WriteString(",\n")
		}

		out.WriteString(prefix + ")")

	case types.TypeObject:
		// var_export shows stdClass objects as stdClass::__set_state(array(...))
		// For now, we'll use a simplified representation
		obj := val.ToObject()
		className := obj.ClassName

		ptr := fmt.Sprintf("%p", obj)
		if visited[ptr] {
			out.WriteString("/* RECURSION */")
			return
		}
		visited[ptr] = true

		out.WriteString(fmt.Sprintf("%s::__set_state(array(\n", className))

		// Access properties directly from instance to bypass visibility checks
		for name, propDef := range obj.ClassEntry.Properties {
			var value *types.Value
			if prop, exists := obj.Properties[name]; exists && prop.Value != nil {
				value = prop.Value
			} else if propDef.HasDefault {
				value = propDef.Default
			} else {
				value = types.NewNull()
			}
			out.WriteString(fmt.Sprintf("%s  '%s' => ", prefix, name))
			exportValue(out, value, indent+1, visited)
			out.WriteString(",\n")
		}

		out.WriteString(prefix + "))")

	case types.TypeResource:
		// Resources can't be exported
		out.WriteString("NULL")

	default:
		out.WriteString("NULL")
	}
}

// ============================================================================
// Type Checking Functions
// ============================================================================

// IsNull checks if a variable is null
// is_null(mixed $value): bool
func IsNull(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeNull)
}

// IsBool checks if a variable is a boolean
// is_bool(mixed $value): bool
func IsBool(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeBool)
}

// IsInt checks if a variable is an integer
// is_int(mixed $value): bool
func IsInt(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeInt)
}

// IsLong is an alias for IsInt
func IsLong(val *types.Value) *types.Value {
	return IsInt(val)
}

// IsInteger is an alias for IsInt
func IsInteger(val *types.Value) *types.Value {
	return IsInt(val)
}

// IsFloat checks if a variable is a float
// is_float(mixed $value): bool
func IsFloat(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeFloat)
}

// IsDouble is an alias for IsFloat
func IsDouble(val *types.Value) *types.Value {
	return IsFloat(val)
}

// IsReal is an alias for IsFloat (deprecated but still in PHP)
func IsReal(val *types.Value) *types.Value {
	return IsFloat(val)
}

// IsString checks if a variable is a string
// is_string(mixed $value): bool
func IsString(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeString)
}

// IsArray checks if a variable is an array
// is_array(mixed $value): bool
func IsArray(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeArray)
}

// IsObject checks if a variable is an object
// is_object(mixed $value): bool
func IsObject(val *types.Value) *types.Value {
	return types.NewBool(val.Type() == types.TypeObject)
}

// IsResource checks if a variable is a resource
// is_resource(mixed $value): bool
func IsResource(val *types.Value) *types.Value {
	if val.Type() != types.TypeResource {
		return types.NewBool(false)
	}
	// Resource must also be valid (not closed)
	res := val.ToResource()
	return types.NewBool(res.IsValid())
}

// IsNumeric checks if a variable is a number or numeric string
// is_numeric(mixed $value): bool
func IsNumeric(val *types.Value) *types.Value {
	switch val.Type() {
	case types.TypeInt, types.TypeFloat:
		return types.NewBool(true)

	case types.TypeString:
		str := val.ToString()
		// Try to parse as int or float
		_, errInt := fmt.Sscanf(str, "%d", new(int64))
		if errInt == nil {
			return types.NewBool(true)
		}
		_, errFloat := fmt.Sscanf(str, "%f", new(float64))
		return types.NewBool(errFloat == nil)

	default:
		return types.NewBool(false)
	}
}

// IsScalar checks if a variable is a scalar (int, float, string, or bool)
// is_scalar(mixed $value): bool
func IsScalar(val *types.Value) *types.Value {
	return types.NewBool(val.IsScalar())
}

// IsCallable checks if a variable can be called as a function
// is_callable(mixed $value): bool
func IsCallable(val *types.Value) *types.Value {
	// TODO: Implement proper callable checking when closures are added
	switch val.Type() {
	case types.TypeString:
		// String could be a function name
		return types.NewBool(true)

	case types.TypeArray:
		// Array like ['ClassName', 'methodName'] or [$object, 'methodName']
		arr := val.ToArray()
		if arr.Len() == 2 {
			return types.NewBool(true)
		}
		return types.NewBool(false)

	case types.TypeObject:
		// Objects with __invoke are callable
		obj := val.ToObject()
		_, hasInvoke := obj.ClassEntry.Methods["__invoke"]
		return types.NewBool(hasInvoke)

	default:
		return types.NewBool(false)
	}
}

// IsIterable checks if a variable can be iterated over
// is_iterable(mixed $value): bool
func IsIterable(val *types.Value) *types.Value {
	switch val.Type() {
	case types.TypeArray:
		return types.NewBool(true)

	case types.TypeObject:
		// Objects that implement Traversable are iterable
		// For now, we'll just check if it's an object
		// TODO: Implement proper interface checking
		return types.NewBool(true)

	default:
		return types.NewBool(false)
	}
}

// IsCountable checks if a variable can be counted
// is_countable(mixed $value): bool
func IsCountable(val *types.Value) *types.Value {
	switch val.Type() {
	case types.TypeArray:
		return types.NewBool(true)

	case types.TypeObject:
		// Objects that implement Countable are countable
		// TODO: Implement proper interface checking
		obj := val.ToObject()
		_, hasCount := obj.ClassEntry.Methods["count"]
		return types.NewBool(hasCount)

	default:
		return types.NewBool(false)
	}
}

// GetType returns the type of a variable
// gettype(mixed $value): string
func GetType(val *types.Value) *types.Value {
	switch val.Type() {
	case types.TypeNull:
		return types.NewString("NULL")
	case types.TypeBool:
		return types.NewString("boolean")
	case types.TypeInt:
		return types.NewString("integer")
	case types.TypeFloat:
		return types.NewString("double")
	case types.TypeString:
		return types.NewString("string")
	case types.TypeArray:
		return types.NewString("array")
	case types.TypeObject:
		return types.NewString("object")
	case types.TypeResource:
		return types.NewString("resource")
	default:
		return types.NewString("unknown type")
	}
}

// ============================================================================
// Function/Class Existence Checks
// ============================================================================

// Note: These functions require VM context to check defined functions/classes.
// The actual implementation will be in builtins.go with access to VM state.

// FunctionExistsStub is a placeholder - real implementation needs VM context
// function_exists(string $function_name): bool
func FunctionExistsStub(funcName *types.Value) *types.Value {
	// This is a stub - the real implementation is in builtins.go
	// which has access to the VM's function registry
	return types.NewBool(false)
}

// ClassExistsStub is a placeholder - real implementation needs VM context
// class_exists(string $class, bool $autoload = true): bool
func ClassExistsStub(className *types.Value, args ...*types.Value) *types.Value {
	// This is a stub - the real implementation is in builtins.go
	// which has access to the VM's class registry
	return types.NewBool(false)
}

// MethodExistsStub is a placeholder - real implementation needs VM context
// method_exists(object|string $object_or_class, string $method): bool
func MethodExistsStub(objectOrClass *types.Value, method *types.Value) *types.Value {
	// Check if the first argument is an object
	if objectOrClass.Type() == types.TypeObject {
		obj := objectOrClass.ToObject()
		methodName := method.ToString()
		_, exists := obj.ClassEntry.Methods[methodName]
		return types.NewBool(exists)
	}
	// For string class names, we need VM context
	return types.NewBool(false)
}

// ============================================================================
// Serialization Functions
// ============================================================================

// Serialize generates a storable representation of a value
// serialize(mixed $value): string
func Serialize(value *types.Value) *types.Value {
	// PHP serialize format implementation
	return types.NewString(serializeValue(value))
}

// serializeValue recursively serializes a PHP value
func serializeValue(value *types.Value) string {
	switch value.Type() {
	case types.TypeNull:
		return "N;"
	case types.TypeBool:
		if value.ToBool() {
			return "b:1;"
		}
		return "b:0;"
	case types.TypeInt:
		return fmt.Sprintf("i:%d;", value.ToInt())
	case types.TypeFloat:
		return fmt.Sprintf("d:%v;", value.ToFloat())
	case types.TypeString:
		s := value.ToString()
		return fmt.Sprintf("s:%d:\"%s\";", len(s), s)
	case types.TypeArray:
		arr := value.ToArray()
		var result strings.Builder
		result.WriteString(fmt.Sprintf("a:%d:{", arr.Len()))
		arr.Each(func(key, val *types.Value) bool {
			result.WriteString(serializeValue(key))
			result.WriteString(serializeValue(val))
			return true
		})
		result.WriteString("}")
		return result.String()
	case types.TypeObject:
		obj := value.ToObject()
		className := obj.ClassName
		var result strings.Builder
		propCount := len(obj.Properties)
		result.WriteString(fmt.Sprintf("O:%d:\"%s\":%d:{", len(className), className, propCount))
		for name, prop := range obj.Properties {
			// Serialize property name
			result.WriteString(fmt.Sprintf("s:%d:\"%s\";", len(name), name))
			// Serialize property value
			result.WriteString(serializeValue(prop.Value))
		}
		result.WriteString("}")
		return result.String()
	default:
		return "N;"
	}
}

// Unserialize creates a PHP value from a stored representation
// unserialize(string $data, array $options = []): mixed
func Unserialize(data *types.Value, args ...*types.Value) *types.Value {
	str := data.ToString()
	if str == "" {
		return types.NewBool(false)
	}

	result, _, err := unserializeValue(str, 0)
	if err != nil {
		return types.NewBool(false)
	}
	return result
}

// unserializeValue recursively unserializes a PHP value
// Returns the value, the new position, and any error
func unserializeValue(data string, pos int) (*types.Value, int, error) {
	if pos >= len(data) {
		return nil, pos, fmt.Errorf("unexpected end of data")
	}

	switch data[pos] {
	case 'N': // Null
		// Expect "N;"
		if pos+1 < len(data) && data[pos+1] == ';' {
			return types.NewNull(), pos + 2, nil
		}
		return nil, pos, fmt.Errorf("invalid null format")

	case 'b': // Boolean
		// Expect "b:0;" or "b:1;"
		if pos+3 < len(data) && data[pos+1] == ':' && data[pos+3] == ';' {
			if data[pos+2] == '1' {
				return types.NewBool(true), pos + 4, nil
			}
			return types.NewBool(false), pos + 4, nil
		}
		return nil, pos, fmt.Errorf("invalid boolean format")

	case 'i': // Integer
		// Expect "i:NUMBER;"
		if pos+2 < len(data) && data[pos+1] == ':' {
			endPos := pos + 2
			for endPos < len(data) && data[endPos] != ';' {
				endPos++
			}
			if endPos >= len(data) {
				return nil, pos, fmt.Errorf("invalid integer format")
			}
			numStr := data[pos+2 : endPos]
			var num int64
			fmt.Sscanf(numStr, "%d", &num)
			return types.NewInt(num), endPos + 1, nil
		}
		return nil, pos, fmt.Errorf("invalid integer format")

	case 'd': // Float/Double
		// Expect "d:NUMBER;"
		if pos+2 < len(data) && data[pos+1] == ':' {
			endPos := pos + 2
			for endPos < len(data) && data[endPos] != ';' {
				endPos++
			}
			if endPos >= len(data) {
				return nil, pos, fmt.Errorf("invalid float format")
			}
			numStr := data[pos+2 : endPos]
			var num float64
			fmt.Sscanf(numStr, "%f", &num)
			return types.NewFloat(num), endPos + 1, nil
		}
		return nil, pos, fmt.Errorf("invalid float format")

	case 's': // String
		// Expect "s:LENGTH:"STRING";"
		if pos+2 < len(data) && data[pos+1] == ':' {
			colonPos := pos + 2
			for colonPos < len(data) && data[colonPos] != ':' {
				colonPos++
			}
			if colonPos >= len(data) {
				return nil, pos, fmt.Errorf("invalid string format")
			}
			lengthStr := data[pos+2 : colonPos]
			var length int
			fmt.Sscanf(lengthStr, "%d", &length)

			// Expect opening quote
			if colonPos+1 >= len(data) || data[colonPos+1] != '"' {
				return nil, pos, fmt.Errorf("invalid string format: missing opening quote")
			}

			stringStart := colonPos + 2
			stringEnd := stringStart + length

			if stringEnd > len(data) {
				return nil, pos, fmt.Errorf("invalid string format: string too short")
			}

			str := data[stringStart:stringEnd]

			// Expect closing quote and semicolon
			if stringEnd+1 >= len(data) || data[stringEnd] != '"' || data[stringEnd+1] != ';' {
				return nil, pos, fmt.Errorf("invalid string format: missing closing")
			}

			return types.NewString(str), stringEnd + 2, nil
		}
		return nil, pos, fmt.Errorf("invalid string format")

	case 'a': // Array
		// Expect "a:SIZE:{...}"
		if pos+2 < len(data) && data[pos+1] == ':' {
			colonPos := pos + 2
			for colonPos < len(data) && data[colonPos] != ':' {
				colonPos++
			}
			if colonPos >= len(data) {
				return nil, pos, fmt.Errorf("invalid array format")
			}
			sizeStr := data[pos+2 : colonPos]
			var size int
			fmt.Sscanf(sizeStr, "%d", &size)

			// Expect opening brace
			if colonPos+1 >= len(data) || data[colonPos+1] != '{' {
				return nil, pos, fmt.Errorf("invalid array format: missing opening brace")
			}

			arr := types.NewEmptyArray()
			currentPos := colonPos + 2

			for i := 0; i < size; i++ {
				// Unserialize key
				key, newPos, err := unserializeValue(data, currentPos)
				if err != nil {
					return nil, pos, err
				}
				currentPos = newPos

				// Unserialize value
				val, newPos, err := unserializeValue(data, currentPos)
				if err != nil {
					return nil, pos, err
				}
				currentPos = newPos

				arr.Set(key, val)
			}

			// Expect closing brace
			if currentPos >= len(data) || data[currentPos] != '}' {
				return nil, pos, fmt.Errorf("invalid array format: missing closing brace")
			}

			return types.NewArray(arr), currentPos + 1, nil
		}
		return nil, pos, fmt.Errorf("invalid array format")

	default:
		return nil, pos, fmt.Errorf("unknown type: %c", data[pos])
	}
}
