// Package goext provides Go extension and FFI integration for PHP-Go.
// It enables bidirectional type conversion and function calls between PHP and Go.
package goext

import (
	"fmt"
	"reflect"

	"github.com/krizos/php-go/pkg/types"
)

// Marshaler handles conversion between PHP and Go types.
// It provides bidirectional marshaling with proper type handling.
type Marshaler struct{}

// NewMarshaler creates a new Marshaler instance.
func NewMarshaler() *Marshaler {
	return &Marshaler{}
}

// ToGo converts a PHP Value to a native Go type.
// Returns interface{} that can be type-asserted to the actual Go type.
// Type mapping:
//   - Null -> nil
//   - Bool -> bool
//   - Int -> int64
//   - Float -> float64
//   - String -> string
//   - Array (list) -> []interface{}
//   - Array (assoc) -> map[string]interface{}
//   - Object -> map[string]interface{}
func (m *Marshaler) ToGo(phpVal *types.Value) (interface{}, error) {
	if phpVal == nil {
		return nil, nil
	}

	switch phpVal.Type() {
	case types.TypeNull, types.TypeUndef:
		return nil, nil

	case types.TypeBool:
		return phpVal.ToBool(), nil

	case types.TypeInt:
		return phpVal.ToInt(), nil

	case types.TypeFloat:
		return phpVal.ToFloat(), nil

	case types.TypeString:
		return phpVal.ToString(), nil

	case types.TypeArray:
		arr := phpVal.ToArray()
		if arr == nil {
			return []interface{}{}, nil
		}

		// Check if array is a list (sequential integer keys starting from 0)
		if m.isListArray(phpVal) {
			return m.arrayToSlice(arr)
		}
		// Otherwise, convert to map
		return m.arrayToMap(arr)

	case types.TypeObject:
		obj := phpVal.ToObject()
		if obj == nil {
			return map[string]interface{}{}, nil
		}

		// Convert object properties to map
		result := make(map[string]interface{})
		for key, prop := range obj.Properties {
			goVal, err := m.ToGo(prop.Value)
			if err != nil {
				return nil, fmt.Errorf("converting object property '%s': %w", key, err)
			}
			result[key] = goVal
		}
		return result, nil

	case types.TypeReference:
		// Dereference and convert
		deref := phpVal.Deref()
		return m.ToGo(deref)

	default:
		return nil, fmt.Errorf("unsupported PHP type for Go conversion: %s", phpVal.TypeString())
	}
}

// ToGoInt converts a PHP Value to int64.
func (m *Marshaler) ToGoInt(phpVal *types.Value) (int64, error) {
	if phpVal == nil {
		return 0, fmt.Errorf("cannot convert nil to int64")
	}
	return phpVal.ToInt(), nil
}

// ToGoFloat converts a PHP Value to float64.
func (m *Marshaler) ToGoFloat(phpVal *types.Value) (float64, error) {
	if phpVal == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
	}
	return phpVal.ToFloat(), nil
}

// ToGoBool converts a PHP Value to bool.
func (m *Marshaler) ToGoBool(phpVal *types.Value) (bool, error) {
	if phpVal == nil {
		return false, nil
	}
	return phpVal.ToBool(), nil
}

// ToGoString converts a PHP Value to string.
func (m *Marshaler) ToGoString(phpVal *types.Value) (string, error) {
	if phpVal == nil {
		return "", nil
	}
	return phpVal.ToString(), nil
}

// ToGoSlice converts a PHP array to []interface{}.
func (m *Marshaler) ToGoSlice(phpVal *types.Value) ([]interface{}, error) {
	if phpVal == nil || phpVal.Type() != types.TypeArray {
		return nil, fmt.Errorf("cannot convert non-array to slice")
	}

	arr := phpVal.ToArray()
	if arr == nil {
		return []interface{}{}, nil
	}

	return m.arrayToSlice(arr)
}

// ToGoMap converts a PHP array to map[string]interface{}.
func (m *Marshaler) ToGoMap(phpVal *types.Value) (map[string]interface{}, error) {
	if phpVal == nil || phpVal.Type() != types.TypeArray {
		return nil, fmt.Errorf("cannot convert non-array to map")
	}

	arr := phpVal.ToArray()
	if arr == nil {
		return map[string]interface{}{}, nil
	}

	return m.arrayToMap(arr)
}

// ToPHP converts a native Go value to a PHP Value.
// Supports: nil, bool, int, int64, float64, string, []interface{}, map[string]interface{}, structs
func (m *Marshaler) ToPHP(goVal interface{}) (*types.Value, error) {
	if goVal == nil {
		return types.NewNull(), nil
	}

	v := reflect.ValueOf(goVal)
	return m.reflectToPHP(v)
}

// IntToPHP converts an int64 to PHP Value.
func (m *Marshaler) IntToPHP(val int64) *types.Value {
	return types.NewInt(val)
}

// FloatToPHP converts a float64 to PHP Value.
func (m *Marshaler) FloatToPHP(val float64) *types.Value {
	return types.NewFloat(val)
}

// BoolToPHP converts a bool to PHP Value.
func (m *Marshaler) BoolToPHP(val bool) *types.Value {
	return types.NewBool(val)
}

// StringToPHP converts a string to PHP Value.
func (m *Marshaler) StringToPHP(val string) *types.Value {
	return types.NewString(val)
}

// SliceToPHP converts a []interface{} to PHP array.
func (m *Marshaler) SliceToPHP(val []interface{}) (*types.Value, error) {
	arr := types.NewEmptyArray()
	for _, item := range val {
		phpVal, err := m.ToPHP(item)
		if err != nil {
			return nil, err
		}
		arr.Append(phpVal)
	}
	return types.NewArray(arr), nil
}

// MapToPHP converts a map[string]interface{} to PHP array.
func (m *Marshaler) MapToPHP(val map[string]interface{}) (*types.Value, error) {
	arr := types.NewEmptyArray()
	for key, item := range val {
		phpVal, err := m.ToPHP(item)
		if err != nil {
			return nil, fmt.Errorf("converting map value for key '%s': %w", key, err)
		}
		arr.Set(types.NewString(key), phpVal)
	}
	return types.NewArray(arr), nil
}

// Helper methods

// isListArray checks if a PHP array is a sequential list (0, 1, 2, ...)
func (m *Marshaler) isListArray(phpVal *types.Value) bool {
	if phpVal == nil || phpVal.Type() != types.TypeArray {
		return false
	}

	arr := phpVal.ToArray()
	if arr == nil || arr.Len() == 0 {
		return true // Empty array is a list
	}

	// Check if array is packed (already optimized as sequential)
	if arr.IsPacked() {
		return true
	}

	keys := arr.Keys()
	if keys.Len() == 0 {
		return true
	}

	// Check if all keys are sequential integers starting from 0
	keysArray := keys.Values()
	idx := 0
	keysArray.Each(func(_, key *types.Value) bool {
		if key.Type() != types.TypeInt {
			return false // Stop iteration
		}
		if key.ToInt() != int64(idx) {
			return false
		}
		idx++
		return true // Continue iteration
	})

	return idx == keys.Len()
}

// arrayToSlice converts a PHP array to Go slice
func (m *Marshaler) arrayToSlice(arr *types.Array) ([]interface{}, error) {
	if arr == nil {
		return []interface{}{}, nil
	}

	result := make([]interface{}, 0, arr.Len())
	values := arr.Values()

	values.Each(func(_, val *types.Value) bool {
		goVal, err := m.ToGo(val)
		if err != nil {
			// Can't return error from Each callback
			return false
		}
		result = append(result, goVal)
		return true
	})

	return result, nil
}

// arrayToMap converts a PHP array to Go map
func (m *Marshaler) arrayToMap(arr *types.Array) (map[string]interface{}, error) {
	if arr == nil {
		return map[string]interface{}{}, nil
	}

	result := make(map[string]interface{})

	arr.Each(func(key, val *types.Value) bool {
		// Convert key to string
		keyStr := key.ToString()

		goVal, err := m.ToGo(val)
		if err != nil {
			// Can't return error from Each callback
			return false
		}
		result[keyStr] = goVal
		return true
	})

	return result, nil
}

// reflectToPHP uses reflection to convert arbitrary Go values to PHP
func (m *Marshaler) reflectToPHP(v reflect.Value) (*types.Value, error) {
	// Handle invalid/nil values
	if !v.IsValid() {
		return types.NewNull(), nil
	}

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return types.NewNull(), nil
		}
		v = v.Elem()
	}

	// Handle interfaces
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return types.NewNull(), nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		return types.NewBool(v.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.NewInt(v.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.NewInt(int64(v.Uint())), nil

	case reflect.Float32, reflect.Float64:
		return types.NewFloat(v.Float()), nil

	case reflect.String:
		return types.NewString(v.String()), nil

	case reflect.Slice, reflect.Array:
		arr := types.NewEmptyArray()
		for i := 0; i < v.Len(); i++ {
			phpVal, err := m.reflectToPHP(v.Index(i))
			if err != nil {
				return nil, err
			}
			arr.Append(phpVal)
		}
		return types.NewArray(arr), nil

	case reflect.Map:
		arr := types.NewEmptyArray()
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()

			// Convert key to PHP value
			phpKey, err := m.reflectToPHP(key)
			if err != nil {
				return nil, err
			}

			// Convert value to PHP value
			phpVal, err := m.reflectToPHP(val)
			if err != nil {
				return nil, err
			}

			arr.Set(phpKey, phpVal)
		}
		return types.NewArray(arr), nil

	case reflect.Struct:
		// Convert struct to associative array
		arr := types.NewEmptyArray()
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			fieldVal, err := m.reflectToPHP(v.Field(i))
			if err != nil {
				return nil, err
			}

			arr.Set(types.NewString(field.Name), fieldVal)
		}
		return types.NewArray(arr), nil

	default:
		return nil, fmt.Errorf("unsupported Go type for PHP conversion: %s", v.Kind())
	}
}
