package bindings

import (
	"encoding/json"
	"fmt"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// JSONExtension provides Go JSON encoding/decoding for PHP.
type JSONExtension struct {
	*goext.BaseExtension
	marshaler *goext.Marshaler
}

// NewJSONExtension creates a new JSON extension.
func NewJSONExtension() *JSONExtension {
	ext := &JSONExtension{
		BaseExtension: goext.NewBaseExtension("go_json", "1.0.0"),
		marshaler:     goext.NewMarshaler(),
	}

	// Register functions
	ext.AddFunction("go_json_encode", ext.jsonEncode)
	ext.AddFunction("go_json_decode", ext.jsonDecode)
	ext.AddFunction("go_json_validate", ext.jsonValidate)

	// Register constants
	ext.AddConstant("JSON_PRETTY_PRINT", types.NewInt(1))
	ext.AddConstant("JSON_UNESCAPED_SLASHES", types.NewInt(2))
	ext.AddConstant("JSON_UNESCAPED_UNICODE", types.NewInt(4))

	return ext
}

// jsonEncode implements go_json_encode($value, $flags = 0)
//
// Flags:
//   - JSON_PRETTY_PRINT (1): Pretty-print JSON
//   - JSON_UNESCAPED_SLASHES (2): Don't escape forward slashes
//   - JSON_UNESCAPED_UNICODE (4): Don't escape unicode characters
//
// Returns: JSON string or false on error
func (e *JSONExtension) jsonEncode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_json_encode() expects at least 1 argument (value), got %d", len(args))
	}

	phpVal := args[0]

	// Parse flags
	var flags int64
	if len(args) >= 2 {
		flags = args[1].ToInt()
	}

	// Convert PHP value to Go
	goVal, err := e.marshaler.ToGo(phpVal)
	if err != nil {
		return types.NewBool(false), nil
	}

	// Encode to JSON
	var jsonBytes []byte
	if flags&1 != 0 { // JSON_PRETTY_PRINT
		jsonBytes, err = json.MarshalIndent(goVal, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(goVal)
	}

	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(string(jsonBytes)), nil
}

// jsonDecode implements go_json_decode($json, $assoc = false)
//
// Parameters:
//   - json: JSON string to decode
//   - assoc: if true, return arrays instead of objects
//
// Returns: decoded value or null on error
func (e *JSONExtension) jsonDecode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_json_decode() expects at least 1 argument (json), got %d", len(args))
	}

	jsonStr := args[0].ToString()

	// Parse assoc flag
	assoc := false
	if len(args) >= 2 {
		assoc = args[1].ToBool()
	}

	// Decode JSON
	var goVal interface{}
	if err := json.Unmarshal([]byte(jsonStr), &goVal); err != nil {
		return types.NewNull(), nil
	}

	// Convert to PHP value
	phpVal, err := e.marshaler.ToPHP(goVal)
	if err != nil {
		return types.NewNull(), nil
	}

	// If assoc is false and we have an object, we need to keep it as an object
	// For now, we always return arrays/maps since we don't have full object support
	_ = assoc

	return phpVal, nil
}

// jsonValidate implements go_json_validate($json)
//
// Returns: true if valid JSON, false otherwise
func (e *JSONExtension) jsonValidate(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_json_validate() expects 1 argument (json), got %d", len(args))
	}

	jsonStr := args[0].ToString()

	var tmp interface{}
	err := json.Unmarshal([]byte(jsonStr), &tmp)

	return types.NewBool(err == nil), nil
}
