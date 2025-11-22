package types

// Array represents a PHP array (ordered map)
// This will be fully implemented in Phase 4.
// For now, we provide a minimal implementation to satisfy the Value system.
type Array struct {
	// elements stores the array data as an ordered map
	// Key can be int64 or string
	elements map[interface{}]*Value
	// order tracks insertion order for iteration
	order []interface{}
}

// NewEmptyArray creates a new empty array
func NewEmptyArray() *Array {
	return &Array{
		elements: make(map[interface{}]*Value),
		order:    make([]interface{}, 0),
	}
}

// Len returns the number of elements in the array
func (a *Array) Len() int {
	if a == nil {
		return 0
	}
	return len(a.elements)
}

// Get retrieves a value from the array
func (a *Array) Get(key *Value) (*Value, bool) {
	if a == nil {
		return NewNull(), false
	}

	// Convert key to appropriate type
	var k interface{}
	switch key.Type() {
	case TypeInt:
		k = key.ToInt()
	case TypeString:
		k = key.ToString()
	default:
		// PHP converts other types to string for array keys
		k = key.ToString()
	}

	val, exists := a.elements[k]
	if !exists {
		return NewNull(), false
	}
	return val, true
}

// Set adds or updates a value in the array
func (a *Array) Set(key *Value, value *Value) {
	if a == nil {
		return
	}

	// Convert key to appropriate type
	var k interface{}
	switch key.Type() {
	case TypeInt:
		k = key.ToInt()
	case TypeString:
		k = key.ToString()
	default:
		k = key.ToString()
	}

	// Add to order if new key
	if _, exists := a.elements[k]; !exists {
		a.order = append(a.order, k)
	}

	a.elements[k] = value
}

// DeepCopy creates a deep copy of the array
func (a *Array) DeepCopy() *Array {
	if a == nil {
		return NewEmptyArray()
	}

	copied := NewEmptyArray()
	for _, key := range a.order {
		val := a.elements[key]

		// Convert key back to Value
		var keyVal *Value
		switch k := key.(type) {
		case int64:
			keyVal = NewInt(k)
		case string:
			keyVal = NewString(k)
		default:
			keyVal = NewString("")
		}

		copied.Set(keyVal, val.DeepCopy())
	}

	return copied
}
