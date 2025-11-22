package types

import "fmt"

// Array represents a PHP array (ordered associative array)
// PHP arrays are ordered maps that can have both integer and string keys
// and preserve insertion order.
type Array struct {
	// For packed arrays (sequential integer keys starting from 0)
	packed     bool
	packedData []*Value

	// For hash tables (string keys or non-sequential integer keys)
	elements map[interface{}]*Value
	order    []interface{}

	// Next auto-index for append operations
	nextIndex int64
}

// ============================================================================
// Constructors
// ============================================================================

// NewEmptyArray creates a new empty array
func NewEmptyArray() *Array {
	return &Array{
		packed:     true, // Start as packed, demote to hash if needed
		packedData: make([]*Value, 0, 8),
		elements:   nil,
		order:      nil,
		nextIndex:  0,
	}
}

// NewArrayWithCapacity creates a new array with pre-allocated capacity
func NewArrayWithCapacity(capacity int) *Array {
	return &Array{
		packed:     true,
		packedData: make([]*Value, 0, capacity),
		elements:   nil,
		order:      nil,
		nextIndex:  0,
	}
}

// NewArrayFromSlice creates an array from a slice of values
func NewArrayFromSlice(values []*Value) *Array {
	arr := &Array{
		packed:     true,
		packedData: make([]*Value, len(values)),
		elements:   nil,
		order:      nil,
		nextIndex:  int64(len(values)),
	}
	copy(arr.packedData, values)
	return arr
}

// NewArrayFromMap creates an array from a map
func NewArrayFromMap(data map[interface{}]*Value) *Array {
	arr := &Array{
		packed:     false,
		packedData: nil,
		elements:   make(map[interface{}]*Value),
		order:      make([]interface{}, 0, len(data)),
		nextIndex:  0,
	}

	maxIndex := int64(-1)
	for k, v := range data {
		arr.order = append(arr.order, k)
		arr.elements[k] = v

		// Track next index for integer keys
		if intKey, ok := k.(int64); ok && intKey > maxIndex {
			maxIndex = intKey
		}
	}

	arr.nextIndex = maxIndex + 1
	return arr
}

// ============================================================================
// Basic Properties
// ============================================================================

// Len returns the number of elements in the array
func (a *Array) Len() int {
	if a == nil {
		return 0
	}
	if a.packed {
		return len(a.packedData)
	}
	return len(a.elements)
}

// IsEmpty returns true if the array is empty
func (a *Array) IsEmpty() bool {
	return a.Len() == 0
}

// IsPacked returns true if the array is using packed optimization
func (a *Array) IsPacked() bool {
	if a == nil {
		return false
	}
	return a.packed
}

// ============================================================================
// Key Normalization
// ============================================================================

// normalizeKey converts a Value to an appropriate array key
// Following PHP's key conversion rules:
// - Integers stay as int64
// - Numeric strings become integers ("42" -> 42)
// - Floats become integers (truncated)
// - Booleans: true->1, false->0
// - Null becomes empty string
func normalizeKey(key *Value) interface{} {
	switch key.Type() {
	case TypeInt:
		return key.ToInt()

	case TypeFloat:
		// Float keys are truncated to integers
		return int64(key.ToFloat())

	case TypeBool:
		// Boolean keys: true->1, false->0
		if key.ToBool() {
			return int64(1)
		}
		return int64(0)

	case TypeNull:
		// Null becomes empty string
		return ""

	case TypeString:
		s := key.ToString()
		// Try to convert to integer if it's a numeric string
		if i, err := parseInt(s); err == nil {
			return i
		}
		return s

	default:
		// Other types become strings
		return key.ToString()
	}
}

// parseInt attempts to parse a string as an integer
func parseInt(s string) (int64, error) {
	// Handle empty string
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	// Try parsing as integer
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// ============================================================================
// Get/Set Operations
// ============================================================================

// Get retrieves a value from the array
func (a *Array) Get(key *Value) (*Value, bool) {
	if a == nil {
		return NewNull(), false
	}

	k := normalizeKey(key)

	// Fast path for packed arrays with integer keys
	if a.packed {
		if intKey, ok := k.(int64); ok {
			if intKey >= 0 && intKey < int64(len(a.packedData)) {
				return a.packedData[intKey], true
			}
		}
		return NewNull(), false
	}

	// Hash table lookup
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

	k := normalizeKey(key)

	// Try to keep packed optimization
	if a.packed {
		if intKey, ok := k.(int64); ok {
			// Check if we can append to packed array
			if intKey == int64(len(a.packedData)) {
				a.packedData = append(a.packedData, value)
				a.nextIndex = intKey + 1
				return
			}

			// Check if it's an update to existing packed element
			if intKey >= 0 && intKey < int64(len(a.packedData)) {
				a.packedData[intKey] = value
				return
			}
		}

		// Need to convert to hash table
		a.convertToHash()
	}

	// Add to hash table
	if _, exists := a.elements[k]; !exists {
		a.order = append(a.order, k)
	}
	a.elements[k] = value

	// Update next index if this is an integer key
	if intKey, ok := k.(int64); ok && intKey >= a.nextIndex {
		a.nextIndex = intKey + 1
	}
}

// Append adds a value to the end of the array with auto-incrementing index
func (a *Array) Append(value *Value) {
	if a == nil {
		return
	}

	if a.packed {
		a.packedData = append(a.packedData, value)
		a.nextIndex++
		return
	}

	// Hash table append
	key := a.nextIndex
	a.order = append(a.order, key)
	a.elements[key] = value
	a.nextIndex++
}

// Unset removes a key from the array
func (a *Array) Unset(key *Value) {
	if a == nil {
		return
	}

	k := normalizeKey(key)

	if a.packed {
		// For packed arrays, unset converts to hash
		a.convertToHash()
	}

	// Remove from hash table
	if _, exists := a.elements[k]; exists {
		delete(a.elements, k)

		// Remove from order
		for i, orderKey := range a.order {
			if orderKey == k {
				a.order = append(a.order[:i], a.order[i+1:]...)
				break
			}
		}
	}
}

// ============================================================================
// Array Operations
// ============================================================================

// Push appends one or more values to the end (alias for multiple Append)
func (a *Array) Push(values ...*Value) int {
	for _, v := range values {
		a.Append(v)
	}
	return a.Len()
}

// Pop removes and returns the last element
func (a *Array) Pop() (*Value, bool) {
	if a.IsEmpty() {
		return NewNull(), false
	}

	if a.packed {
		lastIdx := len(a.packedData) - 1
		value := a.packedData[lastIdx]
		a.packedData = a.packedData[:lastIdx]
		a.nextIndex--
		return value, true
	}

	// Hash table pop
	if len(a.order) == 0 {
		return NewNull(), false
	}

	lastKey := a.order[len(a.order)-1]
	value := a.elements[lastKey]
	delete(a.elements, lastKey)
	a.order = a.order[:len(a.order)-1]

	return value, true
}

// Shift removes and returns the first element
func (a *Array) Shift() (*Value, bool) {
	if a.IsEmpty() {
		return NewNull(), false
	}

	if a.packed {
		if len(a.packedData) == 0 {
			return NewNull(), false
		}
		value := a.packedData[0]
		a.packedData = a.packedData[1:]
		// Note: This doesn't reindex, which matches PHP behavior
		return value, true
	}

	// Hash table shift
	if len(a.order) == 0 {
		return NewNull(), false
	}

	firstKey := a.order[0]
	value := a.elements[firstKey]
	delete(a.elements, firstKey)
	a.order = a.order[1:]

	return value, true
}

// Unshift prepends one or more values to the beginning
func (a *Array) Unshift(values ...*Value) int {
	if len(values) == 0 {
		return a.Len()
	}

	if a.packed {
		// Prepend to packed array
		newData := make([]*Value, len(values)+len(a.packedData))
		copy(newData, values)
		copy(newData[len(values):], a.packedData)
		a.packedData = newData
		return len(a.packedData)
	}

	// Convert to hash table for unshift
	a.convertToHash()

	// Prepend to hash table
	newOrder := make([]interface{}, 0, len(values)+len(a.order))
	for i, v := range values {
		key := int64(i)
		newOrder = append(newOrder, key)
		a.elements[key] = v
	}
	newOrder = append(newOrder, a.order...)
	a.order = newOrder

	return len(a.order)
}

// Slice returns a portion of the array
func (a *Array) Slice(offset, length int) *Array {
	if a == nil || a.IsEmpty() {
		return NewEmptyArray()
	}

	totalLen := a.Len()

	// Handle negative offset
	if offset < 0 {
		offset = totalLen + offset
		if offset < 0 {
			offset = 0
		}
	}

	// Handle offset beyond length
	if offset >= totalLen {
		return NewEmptyArray()
	}

	// Calculate end
	end := offset + length
	if length < 0 || end > totalLen {
		end = totalLen
	}

	if a.packed {
		// Slice packed array
		sliced := a.packedData[offset:end]
		return NewArrayFromSlice(sliced)
	}

	// Slice hash table
	result := NewEmptyArray()
	result.convertToHash()

	for i := offset; i < end && i < len(a.order); i++ {
		key := a.order[i]
		result.order = append(result.order, key)
		result.elements[key] = a.elements[key]
	}

	return result
}

// Merge combines this array with another
func (a *Array) Merge(other *Array) *Array {
	if a == nil {
		if other == nil {
			return NewEmptyArray()
		}
		return other.DeepCopy()
	}
	if other == nil {
		return a.DeepCopy()
	}

	result := a.DeepCopy()

	// If other is packed, append its values
	if other.packed {
		for _, v := range other.packedData {
			result.Append(v)
		}
		return result
	}

	// If other is a hash table, merge keys
	for _, key := range other.order {
		value := other.elements[key]

		// String keys are preserved, integer keys are reindexed
		switch k := key.(type) {
		case string:
			result.Set(NewString(k), value)
		case int64:
			result.Append(value)
		}
	}

	return result
}

// Keys returns an array of all keys
func (a *Array) Keys() *Array {
	if a == nil || a.IsEmpty() {
		return NewEmptyArray()
	}

	keys := NewEmptyArray()

	if a.packed {
		for i := 0; i < len(a.packedData); i++ {
			keys.Append(NewInt(int64(i)))
		}
		return keys
	}

	for _, key := range a.order {
		switch k := key.(type) {
		case int64:
			keys.Append(NewInt(k))
		case string:
			keys.Append(NewString(k))
		}
	}

	return keys
}

// Values returns an array of all values
func (a *Array) Values() *Array {
	if a == nil || a.IsEmpty() {
		return NewEmptyArray()
	}

	values := NewEmptyArray()

	if a.packed {
		for _, v := range a.packedData {
			values.Append(v)
		}
		return values
	}

	for _, key := range a.order {
		values.Append(a.elements[key])
	}

	return values
}

// Contains checks if a value exists in the array
func (a *Array) Contains(needle *Value) bool {
	if a == nil || a.IsEmpty() {
		return false
	}

	if a.packed {
		for _, v := range a.packedData {
			if v.Equals(needle) {
				return true
			}
		}
		return false
	}

	for _, key := range a.order {
		if a.elements[key].Equals(needle) {
			return true
		}
	}

	return false
}

// Search finds the first key for a given value
func (a *Array) Search(needle *Value) (*Value, bool) {
	if a == nil || a.IsEmpty() {
		return NewNull(), false
	}

	if a.packed {
		for i, v := range a.packedData {
			if v.Equals(needle) {
				return NewInt(int64(i)), true
			}
		}
		return NewBool(false), false
	}

	for _, key := range a.order {
		if a.elements[key].Equals(needle) {
			switch k := key.(type) {
			case int64:
				return NewInt(k), true
			case string:
				return NewString(k), true
			}
		}
	}

	return NewBool(false), false
}

// HasKey checks if a key exists in the array
func (a *Array) HasKey(key *Value) bool {
	if a == nil {
		return false
	}

	k := normalizeKey(key)

	if a.packed {
		if intKey, ok := k.(int64); ok {
			return intKey >= 0 && intKey < int64(len(a.packedData))
		}
		return false
	}

	_, exists := a.elements[k]
	return exists
}

// ============================================================================
// Iteration Support
// ============================================================================

// Each iterates over the array calling fn for each key-value pair
func (a *Array) Each(fn func(key, value *Value) bool) {
	if a == nil || a.IsEmpty() {
		return
	}

	if a.packed {
		for i, v := range a.packedData {
			if !fn(NewInt(int64(i)), v) {
				break
			}
		}
		return
	}

	for _, key := range a.order {
		var keyVal *Value
		switch k := key.(type) {
		case int64:
			keyVal = NewInt(k)
		case string:
			keyVal = NewString(k)
		default:
			continue
		}

		if !fn(keyVal, a.elements[key]) {
			break
		}
	}
}

// ============================================================================
// Conversion and Copying
// ============================================================================

// DeepCopy creates a deep copy of the array
func (a *Array) DeepCopy() *Array {
	if a == nil {
		return NewEmptyArray()
	}

	if a.packed {
		copied := &Array{
			packed:     true,
			packedData: make([]*Value, len(a.packedData)),
			elements:   nil,
			order:      nil,
			nextIndex:  a.nextIndex,
		}
		for i, v := range a.packedData {
			copied.packedData[i] = v.DeepCopy()
		}
		return copied
	}

	copied := &Array{
		packed:     false,
		packedData: nil,
		elements:   make(map[interface{}]*Value, len(a.elements)),
		order:      make([]interface{}, len(a.order)),
		nextIndex:  a.nextIndex,
	}

	copy(copied.order, a.order)
	for k, v := range a.elements {
		copied.elements[k] = v.DeepCopy()
	}

	return copied
}

// convertToHash converts a packed array to a hash table
func (a *Array) convertToHash() {
	if !a.packed {
		return
	}

	a.elements = make(map[interface{}]*Value, len(a.packedData))
	a.order = make([]interface{}, len(a.packedData))

	for i, v := range a.packedData {
		key := int64(i)
		a.elements[key] = v
		a.order[i] = key
	}

	a.packed = false
	a.packedData = nil
}

// Reset resets the array to empty
func (a *Array) Reset() {
	if a == nil {
		return
	}

	a.packed = true
	a.packedData = make([]*Value, 0, 8)
	a.elements = nil
	a.order = nil
	a.nextIndex = 0
}

// ============================================================================
// String Representation
// ============================================================================

// String returns a string representation of the array (for debugging)
func (a *Array) String() string {
	if a == nil || a.IsEmpty() {
		return "[]"
	}

	result := "["
	first := true

	if a.packed {
		for i, v := range a.packedData {
			if !first {
				result += ", "
			}
			result += fmt.Sprintf("%d => %v", i, v)
			first = false
		}
	} else {
		for _, key := range a.order {
			if !first {
				result += ", "
			}
			result += fmt.Sprintf("%v => %v", key, a.elements[key])
			first = false
		}
	}

	result += "]"
	return result
}
