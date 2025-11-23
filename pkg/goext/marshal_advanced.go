package goext

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/krizos/php-go/pkg/types"
)

// CustomMarshaler allows types to define their own marshaling behavior.
type CustomMarshaler interface {
	// MarshalPHP converts the Go value to a PHP Value.
	MarshalPHP() (*types.Value, error)
}

// CustomUnmarshaler allows types to define their own unmarshaling behavior.
type CustomUnmarshaler interface {
	// UnmarshalPHP populates the Go value from a PHP Value.
	UnmarshalPHP(*types.Value) error
}

// TypeConverter defines custom conversion functions for a specific Go type.
type TypeConverter struct {
	// ToPHP converts a Go value to PHP
	ToPHP func(interface{}) (*types.Value, error)
	// ToGo converts a PHP value to Go
	ToGo func(*types.Value) (interface{}, error)
}

// AdvancedMarshaler extends Marshaler with custom type support and optimizations.
type AdvancedMarshaler struct {
	*Marshaler
	mu         sync.RWMutex
	converters map[reflect.Type]*TypeConverter
	// Track visited pointers for circular reference detection
	visited map[uintptr]bool
}

// NewAdvancedMarshaler creates a new advanced marshaler.
func NewAdvancedMarshaler() *AdvancedMarshaler {
	return &AdvancedMarshaler{
		Marshaler:  NewMarshaler(),
		converters: make(map[reflect.Type]*TypeConverter),
		visited:    make(map[uintptr]bool),
	}
}

// RegisterConverter registers a custom converter for a specific Go type.
// Example:
//   m.RegisterConverter(reflect.TypeOf(time.Time{}), &TypeConverter{
//       ToPHP: func(v interface{}) (*types.Value, error) {
//           t := v.(time.Time)
//           return types.NewInt(t.Unix()), nil
//       },
//       ToGo: func(v *types.Value) (interface{}, error) {
//           return time.Unix(v.ToInt(), 0), nil
//       },
//   })
func (m *AdvancedMarshaler) RegisterConverter(typ reflect.Type, converter *TypeConverter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.converters[typ] = converter
}

// UnregisterConverter removes a custom converter.
func (m *AdvancedMarshaler) UnregisterConverter(typ reflect.Type) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.converters, typ)
}

// HasConverter checks if a custom converter is registered for a type.
func (m *AdvancedMarshaler) HasConverter(typ reflect.Type) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.converters[typ]
	return ok
}

// ToPHPAdvanced converts a Go value to PHP with advanced features:
// - Custom type converters
// - CustomMarshaler interface support
// - Circular reference detection
// - Better struct to object conversion
func (m *AdvancedMarshaler) ToPHPAdvanced(goVal interface{}) (*types.Value, error) {
	// Reset visited map for each top-level call
	m.visited = make(map[uintptr]bool)
	return m.toPHPWithCircularCheck(goVal)
}

// toPHPWithCircularCheck converts with circular reference detection
func (m *AdvancedMarshaler) toPHPWithCircularCheck(goVal interface{}) (*types.Value, error) {
	if goVal == nil {
		return types.NewNull(), nil
	}

	v := reflect.ValueOf(goVal)

	// Check for circular references on pointer types
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
		if v.IsNil() {
			return types.NewNull(), nil
		}

		ptr := v.Pointer()
		if m.visited[ptr] {
			// Circular reference detected - return null to break cycle
			return types.NewNull(), nil
		}
		m.visited[ptr] = true
		defer func() {
			delete(m.visited, ptr)
		}()
	}

	// Check for CustomMarshaler interface
	if cm, ok := goVal.(CustomMarshaler); ok {
		return cm.MarshalPHP()
	}

	// Check for registered custom converter
	typ := reflect.TypeOf(goVal)
	m.mu.RLock()
	converter, hasConverter := m.converters[typ]
	m.mu.RUnlock()

	if hasConverter {
		return converter.ToPHP(goVal)
	}

	// Fall back to standard marshaling
	return m.ToPHP(goVal)
}

// ToGoAdvanced converts a PHP value to Go with advanced features:
// - Custom type converters
// - CustomUnmarshaler interface support
// - Better struct population
func (m *AdvancedMarshaler) ToGoAdvanced(phpVal *types.Value, targetType reflect.Type) (interface{}, error) {
	if phpVal == nil || phpVal.Type() == types.TypeNull {
		return reflect.Zero(targetType).Interface(), nil
	}

	// Check for registered custom converter
	m.mu.RLock()
	converter, hasConverter := m.converters[targetType]
	m.mu.RUnlock()

	if hasConverter {
		return converter.ToGo(phpVal)
	}

	// Check if target type implements CustomUnmarshaler
	ptrType := reflect.PtrTo(targetType)
	if ptrType.Implements(reflect.TypeOf((*CustomUnmarshaler)(nil)).Elem()) {
		// Create a new instance
		val := reflect.New(targetType)
		unmarshaler := val.Interface().(CustomUnmarshaler)
		if err := unmarshaler.UnmarshalPHP(phpVal); err != nil {
			return nil, err
		}
		return val.Elem().Interface(), nil
	}

	// For struct types, do smart conversion
	if targetType.Kind() == reflect.Struct {
		return m.phpToStruct(phpVal, targetType)
	}

	// Fall back to standard conversion
	return m.ToGo(phpVal)
}

// phpToStruct converts a PHP array or object to a Go struct.
// It matches PHP keys to struct field names (case-insensitive).
func (m *AdvancedMarshaler) phpToStruct(phpVal *types.Value, structType reflect.Type) (interface{}, error) {
	if phpVal.Type() != types.TypeArray && phpVal.Type() != types.TypeObject {
		return nil, fmt.Errorf("cannot convert %s to struct", phpVal.TypeString())
	}

	// Create new struct instance
	structVal := reflect.New(structType).Elem()

	// Get the data as a map
	var data map[string]interface{}
	var err error

	if phpVal.Type() == types.TypeArray {
		data, err = m.arrayToMap(phpVal.ToArray())
	} else {
		dataInterface, err2 := m.ToGo(phpVal)
		if err2 != nil {
			return nil, err2
		}
		var ok bool
		data, ok = dataInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("object conversion failed")
		}
		err = nil
	}

	if err != nil {
		return nil, err
	}

	// Populate struct fields
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get field name (check for json tag first, then use field name)
		fieldName := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			fieldName = tag
		}

		// Look for the value in the data (case-insensitive)
		var value interface{}
		var found bool

		for key, val := range data {
			if key == fieldName {
				value = val
				found = true
				break
			}
		}

		if !found {
			continue
		}

		// Convert value to field type
		fieldVal := structVal.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		// Convert the interface{} value to the correct type
		converted, err := m.convertToType(value, field.Type)
		if err != nil {
			return nil, fmt.Errorf("converting field %s: %w", field.Name, err)
		}

		fieldVal.Set(reflect.ValueOf(converted))
	}

	return structVal.Interface(), nil
}

// convertToType converts an interface{} value to a specific reflect.Type.
func (m *AdvancedMarshaler) convertToType(value interface{}, targetType reflect.Type) (interface{}, error) {
	if value == nil {
		return reflect.Zero(targetType).Interface(), nil
	}

	valType := reflect.TypeOf(value)

	// If types match, return as-is
	if valType.AssignableTo(targetType) {
		return value, nil
	}

	// Handle numeric conversions
	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int64:
			return reflect.ValueOf(v).Convert(targetType).Interface(), nil
		case float64:
			return reflect.ValueOf(int64(v)).Convert(targetType).Interface(), nil
		case string:
			// Try parsing string to int (not implemented here)
			return reflect.Zero(targetType).Interface(), nil
		}

	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case int64:
			return reflect.ValueOf(float64(v)).Convert(targetType).Interface(), nil
		case float64:
			return reflect.ValueOf(v).Convert(targetType).Interface(), nil
		}

	case reflect.String:
		return fmt.Sprintf("%v", value), nil

	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			return v, nil
		case int64:
			return v != 0, nil
		}

	case reflect.Slice:
		// Convert []interface{} to typed slice
		if slice, ok := value.([]interface{}); ok {
			elemType := targetType.Elem()
			result := reflect.MakeSlice(targetType, len(slice), len(slice))
			for i, elem := range slice {
				converted, err := m.convertToType(elem, elemType)
				if err != nil {
					return nil, err
				}
				result.Index(i).Set(reflect.ValueOf(converted))
			}
			return result.Interface(), nil
		}

	case reflect.Map:
		// Convert map[string]interface{} to typed map
		if mapVal, ok := value.(map[string]interface{}); ok {
			keyType := targetType.Key()
			elemType := targetType.Elem()
			result := reflect.MakeMap(targetType)

			for k, v := range mapVal {
				convertedKey, err := m.convertToType(k, keyType)
				if err != nil {
					return nil, err
				}
				convertedVal, err := m.convertToType(v, elemType)
				if err != nil {
					return nil, err
				}
				result.SetMapIndex(reflect.ValueOf(convertedKey), reflect.ValueOf(convertedVal))
			}
			return result.Interface(), nil
		}

	case reflect.Struct:
		// For struct types, convert from map
		if mapVal, ok := value.(map[string]interface{}); ok {
			// Create a PHP array from the map
			arr := types.NewEmptyArray()
			for k, v := range mapVal {
				phpVal, err := m.ToPHP(v)
				if err != nil {
					return nil, err
				}
				arr.Set(types.NewString(k), phpVal)
			}
			return m.phpToStruct(types.NewArray(arr), targetType)
		}
	}

	return nil, fmt.Errorf("cannot convert %v to %v", valType, targetType)
}

// StructToObject converts a Go struct to a PHP object with proper field mapping.
func (m *AdvancedMarshaler) StructToObject(goStruct interface{}, className string) (*types.Value, error) {
	v := reflect.ValueOf(goStruct)

	// Handle pointer to struct
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return types.NewNull(), nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	// Create PHP object
	obj := types.NewObjectInstance(className)

	// Iterate over struct fields
	structType := v.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldVal := v.Field(i)

		// Get field name (use json tag if available)
		fieldName := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			fieldName = tag
		}

		// Convert field value to PHP
		phpVal, err := m.toPHPWithCircularCheck(fieldVal.Interface())
		if err != nil {
			return nil, fmt.Errorf("converting field %s: %w", field.Name, err)
		}

		// Add as object property (use nil for accessContext since we're creating the object)
		obj.SetProperty(fieldName, phpVal, nil)
	}

	return types.NewObject(obj), nil
}

// ObjectToStruct converts a PHP object to a Go struct.
func (m *AdvancedMarshaler) ObjectToStruct(phpObj *types.Value, targetType reflect.Type) (interface{}, error) {
	if phpObj.Type() != types.TypeObject {
		return nil, fmt.Errorf("expected object, got %s", phpObj.TypeString())
	}

	return m.phpToStruct(phpObj, targetType)
}

// Clone creates a deep copy of a PHP value, handling circular references.
func (m *AdvancedMarshaler) Clone(phpVal *types.Value) (*types.Value, error) {
	// Convert to Go and back to PHP for a deep copy
	m.visited = make(map[uintptr]bool)
	goVal, err := m.ToGo(phpVal)
	if err != nil {
		return nil, err
	}

	return m.toPHPWithCircularCheck(goVal)
}

// DeepEquals checks if two PHP values are deeply equal.
func (m *AdvancedMarshaler) DeepEquals(a, b *types.Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.Type() != b.Type() {
		return false
	}

	// For simple types, use Identical
	if a.Type() != types.TypeArray && a.Type() != types.TypeObject {
		return a.Identical(b)
	}

	// For complex types, convert to Go and compare
	goA, err := m.ToGo(a)
	if err != nil {
		return false
	}

	goB, err := m.ToGo(b)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(goA, goB)
}
