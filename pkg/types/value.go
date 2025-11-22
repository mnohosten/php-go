package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ValueType represents the type of a PHP value
type ValueType uint8

const (
	TypeUndef    ValueType = iota // Undefined/uninitialized
	TypeNull                       // null
	TypeBool                       // bool (true/false)
	TypeInt                        // int64
	TypeFloat                      // float64
	TypeString                     // string
	TypeArray                      // PHP array
	TypeObject                     // PHP object
	TypeResource                   // Resource handle
	TypeReference                  // Reference to another value
)

// ValueFlags represents metadata about a value
type ValueFlags uint8

const (
	FlagNone       ValueFlags = 0
	FlagIsRef      ValueFlags = 1 << 0 // Value is a reference
	FlagImmutable  ValueFlags = 1 << 1 // Value is immutable (constant)
	FlagPersistent ValueFlags = 1 << 2 // Persistent allocation
)

// Value is PHP's universal value container (zval equivalent)
// This is the core data structure that represents any PHP value.
type Value struct {
	typ   ValueType
	flags ValueFlags
	data  interface{}
}

// ============================================================================
// Type Constructors
// ============================================================================

// NewUndef creates an undefined value
func NewUndef() *Value {
	return &Value{typ: TypeUndef}
}

// NewNull creates a null value
func NewNull() *Value {
	return &Value{typ: TypeNull}
}

// NewBool creates a boolean value
func NewBool(v bool) *Value {
	return &Value{typ: TypeBool, data: v}
}

// NewInt creates an integer value
func NewInt(v int64) *Value {
	return &Value{typ: TypeInt, data: v}
}

// NewFloat creates a float value
func NewFloat(v float64) *Value {
	return &Value{typ: TypeFloat, data: v}
}

// NewString creates a string value
func NewString(v string) *Value {
	return &Value{typ: TypeString, data: v}
}

// NewArray creates an array value
func NewArray(v *Array) *Value {
	return &Value{typ: TypeArray, data: v}
}

// NewObject creates an object value
func NewObject(v *Object) *Value {
	return &Value{typ: TypeObject, data: v}
}

// NewResource creates a resource value
func NewResource(v *Resource) *Value {
	return &Value{typ: TypeResource, data: v}
}

// NewReference creates a reference to another value
func NewReference(v *Value) *Value {
	return &Value{typ: TypeReference, flags: FlagIsRef, data: v}
}

// ============================================================================
// Type Queries
// ============================================================================

// Type returns the type of the value
func (v *Value) Type() ValueType {
	if v == nil {
		return TypeNull
	}
	return v.typ
}

// IsUndef returns true if the value is undefined
func (v *Value) IsUndef() bool {
	return v == nil || v.typ == TypeUndef
}

// IsNull returns true if the value is null
func (v *Value) IsNull() bool {
	return v == nil || v.typ == TypeNull
}

// IsBool returns true if the value is a boolean
func (v *Value) IsBool() bool {
	return v != nil && v.typ == TypeBool
}

// IsInt returns true if the value is an integer
func (v *Value) IsInt() bool {
	return v != nil && v.typ == TypeInt
}

// IsFloat returns true if the value is a float
func (v *Value) IsFloat() bool {
	return v != nil && v.typ == TypeFloat
}

// IsString returns true if the value is a string
func (v *Value) IsString() bool {
	return v != nil && v.typ == TypeString
}

// IsArray returns true if the value is an array
func (v *Value) IsArray() bool {
	return v != nil && v.typ == TypeArray
}

// IsObject returns true if the value is an object
func (v *Value) IsObject() bool {
	return v != nil && v.typ == TypeObject
}

// IsResource returns true if the value is a resource
func (v *Value) IsResource() bool {
	return v != nil && v.typ == TypeResource
}

// IsReference returns true if the value is a reference
func (v *Value) IsReference() bool {
	return v != nil && (v.typ == TypeReference || v.flags&FlagIsRef != 0)
}

// IsScalar returns true if the value is a scalar (bool, int, float, string)
func (v *Value) IsScalar() bool {
	if v == nil {
		return false
	}
	return v.typ == TypeBool || v.typ == TypeInt || v.typ == TypeFloat || v.typ == TypeString
}

// IsCallable returns true if the value can be called as a function
func (v *Value) IsCallable() bool {
	if v == nil {
		return false
	}
	// TODO: Implement in Phase 5 (Objects) - check for __invoke magic method
	return false
}

// ============================================================================
// Type Conversions
// ============================================================================

// ToInt converts the value to an integer
func (v *Value) ToInt() int64 {
	if v == nil || v.typ == TypeNull || v.typ == TypeUndef {
		return 0
	}

	switch v.typ {
	case TypeBool:
		if v.data.(bool) {
			return 1
		}
		return 0
	case TypeInt:
		return v.data.(int64)
	case TypeFloat:
		return int64(v.data.(float64))
	case TypeString:
		return stringToInt(v.data.(string))
	case TypeArray:
		// Empty array -> 0, non-empty array -> 1
		arr := v.data.(*Array)
		if arr.Len() == 0 {
			return 0
		}
		return 1
	case TypeObject:
		// Objects convert to 1
		return 1
	case TypeResource:
		// Resources convert to their ID
		res := v.data.(*Resource)
		return int64(res.ID())
	case TypeReference:
		// Dereference and convert
		return v.data.(*Value).ToInt()
	default:
		return 0
	}
}

// ToFloat converts the value to a float
func (v *Value) ToFloat() float64 {
	if v == nil || v.typ == TypeNull || v.typ == TypeUndef {
		return 0.0
	}

	switch v.typ {
	case TypeBool:
		if v.data.(bool) {
			return 1.0
		}
		return 0.0
	case TypeInt:
		return float64(v.data.(int64))
	case TypeFloat:
		return v.data.(float64)
	case TypeString:
		return stringToFloat(v.data.(string))
	case TypeArray:
		// Empty array -> 0.0, non-empty array -> 1.0
		arr := v.data.(*Array)
		if arr.Len() == 0 {
			return 0.0
		}
		return 1.0
	case TypeObject:
		return 1.0
	case TypeResource:
		res := v.data.(*Resource)
		return float64(res.ID())
	case TypeReference:
		return v.data.(*Value).ToFloat()
	default:
		return 0.0
	}
}

// ToBool converts the value to a boolean (PHP truthiness rules)
func (v *Value) ToBool() bool {
	if v == nil || v.typ == TypeNull || v.typ == TypeUndef {
		return false
	}

	switch v.typ {
	case TypeBool:
		return v.data.(bool)
	case TypeInt:
		return v.data.(int64) != 0
	case TypeFloat:
		f := v.data.(float64)
		return f != 0.0 && !math.IsNaN(f)
	case TypeString:
		s := v.data.(string)
		// Empty string and "0" are false
		return s != "" && s != "0"
	case TypeArray:
		// Empty arrays are false
		arr := v.data.(*Array)
		return arr.Len() > 0
	case TypeObject:
		// Objects are always true
		return true
	case TypeResource:
		// Resources are always true
		return true
	case TypeReference:
		return v.data.(*Value).ToBool()
	default:
		return false
	}
}

// ToString converts the value to a string
func (v *Value) ToString() string {
	if v == nil || v.typ == TypeNull {
		return ""
	}

	switch v.typ {
	case TypeUndef:
		return ""
	case TypeBool:
		if v.data.(bool) {
			return "1"
		}
		return ""
	case TypeInt:
		return strconv.FormatInt(v.data.(int64), 10)
	case TypeFloat:
		f := v.data.(float64)
		// Format float like PHP does
		if math.IsNaN(f) {
			return "NAN"
		}
		if math.IsInf(f, 1) {
			return "INF"
		}
		if math.IsInf(f, -1) {
			return "-INF"
		}
		// Remove trailing zeros
		s := strconv.FormatFloat(f, 'f', -1, 64)
		return s
	case TypeString:
		return v.data.(string)
	case TypeArray:
		return "Array"
	case TypeObject:
		// TODO: Call __toString() magic method in Phase 5
		return "Object"
	case TypeResource:
		res := v.data.(*Resource)
		return fmt.Sprintf("Resource id #%d", res.ID())
	case TypeReference:
		return v.data.(*Value).ToString()
	default:
		return ""
	}
}

// ToArray converts the value to an array
func (v *Value) ToArray() *Array {
	if v == nil || v.typ == TypeNull || v.typ == TypeUndef {
		return NewEmptyArray()
	}

	switch v.typ {
	case TypeArray:
		return v.data.(*Array)
	case TypeObject:
		// TODO: Convert object properties to array in Phase 5
		return NewEmptyArray()
	default:
		// Scalar types become single-element array with key 0
		arr := NewEmptyArray()
		arr.Set(NewInt(0), v.Copy())
		return arr
	}
}

// ToObject converts value to Object
// Returns nil if value is not an object
func (v *Value) ToObject() *Object {
	if v == nil || v.typ != TypeObject {
		return nil
	}
	return v.data.(*Object)
}

// ToResource converts value to Resource
// Returns nil if value is not a resource
func (v *Value) ToResource() *Resource {
	if v == nil || v.typ != TypeResource {
		return nil
	}
	return v.data.(*Resource)
}

// ============================================================================
// PHP Truthiness (IsTrue)
// ============================================================================

// IsTrue returns true if the value is truthy in PHP
// This is the same as ToBool() but more explicit for readability
func (v *Value) IsTrue() bool {
	return v.ToBool()
}

// IsFalse returns true if the value is falsy in PHP
func (v *Value) IsFalse() bool {
	return !v.ToBool()
}

// ============================================================================
// Value Operations
// ============================================================================

// Copy creates a shallow copy of the value
func (v *Value) Copy() *Value {
	if v == nil {
		return NewNull()
	}

	// Create new value with same type and flags
	copied := &Value{
		typ:   v.typ,
		flags: v.flags,
	}

	// Copy data
	switch v.typ {
	case TypeUndef, TypeNull:
		// No data to copy
	case TypeBool:
		copied.data = v.data.(bool)
	case TypeInt:
		copied.data = v.data.(int64)
	case TypeFloat:
		copied.data = v.data.(float64)
	case TypeString:
		copied.data = v.data.(string)
	case TypeArray:
		// Arrays use copy-on-write (COW) semantics
		// For now, just reference the same array
		// TODO: Implement proper COW in Phase 7
		copied.data = v.data.(*Array)
	case TypeObject:
		// Objects are passed by reference in PHP
		copied.data = v.data.(*Object)
	case TypeResource:
		// Resources are passed by reference
		copied.data = v.data.(*Resource)
	case TypeReference:
		// Reference to the same value
		copied.data = v.data.(*Value)
	}

	return copied
}

// DeepCopy creates a deep copy of the value
func (v *Value) DeepCopy() *Value {
	if v == nil {
		return NewNull()
	}

	copied := &Value{
		typ:   v.typ,
		flags: v.flags,
	}

	switch v.typ {
	case TypeUndef, TypeNull:
		// No data to copy
	case TypeBool, TypeInt, TypeFloat, TypeString:
		// Scalars can be copied directly
		copied.data = v.data
	case TypeArray:
		// Deep copy the array
		arr := v.data.(*Array)
		copied.data = arr.DeepCopy()
	case TypeObject:
		// Deep copy the object
		// TODO: Implement in Phase 5
		copied.data = v.data.(*Object)
	case TypeResource:
		// Resources can't be deep copied
		copied.data = v.data.(*Resource)
	case TypeReference:
		// Deep copy the referenced value
		copied.data = v.data.(*Value).DeepCopy()
	}

	return copied
}

// Deref dereferences the value if it's a reference, otherwise returns self
func (v *Value) Deref() *Value {
	if v == nil || v.typ != TypeReference {
		return v
	}
	return v.data.(*Value).Deref()
}

// ============================================================================
// Equality and Comparison
// ============================================================================

// Equals checks loose equality (==)
func (v *Value) Equals(other *Value) bool {
	if v == nil && other == nil {
		return true
	}
	if v == nil || other == nil {
		return v.IsNull() && other.IsNull()
	}

	// Dereference if needed
	v = v.Deref()
	other = other.Deref()

	// Same type comparison
	if v.typ == other.typ {
		switch v.typ {
		case TypeUndef, TypeNull:
			return true
		case TypeBool:
			return v.data.(bool) == other.data.(bool)
		case TypeInt:
			return v.data.(int64) == other.data.(int64)
		case TypeFloat:
			return v.data.(float64) == other.data.(float64)
		case TypeString:
			return v.data.(string) == other.data.(string)
		case TypeArray:
			// Array equality is complex - implement later
			return v.data.(*Array) == other.data.(*Array)
		case TypeObject:
			// Object equality is identity
			return v.data.(*Object) == other.data.(*Object)
		case TypeResource:
			return v.data.(*Resource) == other.data.(*Resource)
		}
	}

	// Type juggling for loose comparison
	// PHP converts to numeric if one side is numeric
	if v.IsScalar() && other.IsScalar() {
		// Try numeric comparison
		if (v.typ == TypeInt || v.typ == TypeFloat) || (other.typ == TypeInt || other.typ == TypeFloat) {
			return v.ToFloat() == other.ToFloat()
		}
		// String comparison
		return v.ToString() == other.ToString()
	}

	return false
}

// Identical checks strict equality (===)
func (v *Value) Identical(other *Value) bool {
	if v == nil && other == nil {
		return true
	}
	if v == nil || other == nil {
		return false
	}

	// Type must match
	if v.typ != other.typ {
		return false
	}

	// Dereference
	v = v.Deref()
	other = other.Deref()

	// Check data
	switch v.typ {
	case TypeUndef, TypeNull:
		return true
	case TypeBool:
		return v.data.(bool) == other.data.(bool)
	case TypeInt:
		return v.data.(int64) == other.data.(int64)
	case TypeFloat:
		return v.data.(float64) == other.data.(float64)
	case TypeString:
		return v.data.(string) == other.data.(string)
	case TypeArray:
		// Arrays must have identical keys and values
		// TODO: Implement proper array identity check
		return v.data.(*Array) == other.data.(*Array)
	case TypeObject:
		// Object identity is reference equality
		return v.data.(*Object) == other.data.(*Object)
	case TypeResource:
		return v.data.(*Resource) == other.data.(*Resource)
	}

	return false
}

// ============================================================================
// Debugging
// ============================================================================

// String returns a string representation for debugging
func (v *Value) String() string {
	if v == nil {
		return "NULL"
	}

	switch v.typ {
	case TypeUndef:
		return "UNDEF"
	case TypeNull:
		return "NULL"
	case TypeBool:
		if v.data.(bool) {
			return "bool(true)"
		}
		return "bool(false)"
	case TypeInt:
		return fmt.Sprintf("int(%d)", v.data.(int64))
	case TypeFloat:
		return fmt.Sprintf("float(%g)", v.data.(float64))
	case TypeString:
		s := v.data.(string)
		if len(s) > 50 {
			s = s[:50] + "..."
		}
		return fmt.Sprintf("string(%d) \"%s\"", len(v.data.(string)), s)
	case TypeArray:
		arr := v.data.(*Array)
		return fmt.Sprintf("array(%d)", arr.Len())
	case TypeObject:
		// TODO: Get class name in Phase 5
		return "object"
	case TypeResource:
		res := v.data.(*Resource)
		return fmt.Sprintf("resource(%d)", res.ID())
	case TypeReference:
		return "&" + v.data.(*Value).String()
	default:
		return "UNKNOWN"
	}
}

// TypeString returns the type as a string
func (v *Value) TypeString() string {
	if v == nil {
		return "NULL"
	}

	switch v.typ {
	case TypeUndef:
		return "undefined"
	case TypeNull:
		return "NULL"
	case TypeBool:
		return "boolean"
	case TypeInt:
		return "integer"
	case TypeFloat:
		return "double"
	case TypeString:
		return "string"
	case TypeArray:
		return "array"
	case TypeObject:
		return "object"
	case TypeResource:
		return "resource"
	case TypeReference:
		return "reference"
	default:
		return "unknown"
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// stringToInt converts a string to an integer following PHP rules
func stringToInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// PHP parses leading numeric part
	// "123abc" -> 123, "abc" -> 0
	var result int64
	negative := false
	i := 0

	// Handle sign
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		if s[i] == '-' {
			negative = true
		}
		i++
	}

	// Parse digits
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		result = result*10 + int64(s[i]-'0')
		i++
	}

	if negative {
		return -result
	}
	return result
}

// stringToFloat converts a string to a float following PHP rules
func stringToFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0.0
	}

	// Try to parse the whole string
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// Parse leading numeric part (PHP behavior)
	// "3.14abc" -> 3.14
	var i int
	for i < len(s) {
		c := s[i]
		if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
			i++
		} else {
			break
		}
	}

	if i > 0 {
		if f, err := strconv.ParseFloat(s[:i], 64); err == nil {
			return f
		}
	}

	return 0.0
}
