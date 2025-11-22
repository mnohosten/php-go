package types

import (
	"testing"
)

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNewEmptyArray(t *testing.T) {
	arr := NewEmptyArray()

	if arr == nil {
		t.Fatal("NewEmptyArray returned nil")
	}

	if arr.Len() != 0 {
		t.Errorf("Expected length 0, got %d", arr.Len())
	}

	if !arr.IsEmpty() {
		t.Error("New array should be empty")
	}

	if !arr.IsPacked() {
		t.Error("New array should start as packed")
	}
}

func TestNewArrayWithCapacity(t *testing.T) {
	arr := NewArrayWithCapacity(100)

	if arr.Len() != 0 {
		t.Errorf("Expected length 0, got %d", arr.Len())
	}

	if !arr.IsPacked() {
		t.Error("Array should be packed")
	}
}

func TestNewArrayFromSlice(t *testing.T) {
	values := []*Value{
		NewInt(10),
		NewInt(20),
		NewInt(30),
	}

	arr := NewArrayFromSlice(values)

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	if !arr.IsPacked() {
		t.Error("Array from slice should be packed")
	}

	val, exists := arr.Get(NewInt(0))
	if !exists || val.ToInt() != 10 {
		t.Error("Expected first element to be 10")
	}
}

func TestNewArrayFromMap(t *testing.T) {
	data := map[interface{}]*Value{
		"name": NewString("John"),
		"age":  NewInt(30),
		int64(0): NewString("zero"),
	}

	arr := NewArrayFromMap(data)

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	if arr.IsPacked() {
		t.Error("Array from map with string keys should not be packed")
	}
}

// ============================================================================
// Get/Set Tests
// ============================================================================

func TestArrayGet(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewString("first"))
	arr.Set(NewInt(1), NewString("second"))

	val, exists := arr.Get(NewInt(0))
	if !exists {
		t.Error("Expected key 0 to exist")
	}
	if val.ToString() != "first" {
		t.Errorf("Expected 'first', got '%s'", val.ToString())
	}

	val, exists = arr.Get(NewInt(1))
	if !exists {
		t.Error("Expected key 1 to exist")
	}
	if val.ToString() != "second" {
		t.Errorf("Expected 'second', got '%s'", val.ToString())
	}

	// Non-existent key
	val, exists = arr.Get(NewInt(5))
	if exists {
		t.Error("Key 5 should not exist")
	}
	if val.Type() != TypeNull {
		t.Error("Non-existent key should return null")
	}
}

func TestArraySetPacked(t *testing.T) {
	arr := NewEmptyArray()

	// Sequential inserts should maintain packed mode
	arr.Set(NewInt(0), NewInt(10))
	arr.Set(NewInt(1), NewInt(20))
	arr.Set(NewInt(2), NewInt(30))

	if !arr.IsPacked() {
		t.Error("Array should still be packed")
	}

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}
}

func TestArraySetConvertsToHash(t *testing.T) {
	arr := NewEmptyArray()

	arr.Set(NewInt(0), NewInt(10))
	arr.Set(NewInt(1), NewInt(20))

	// Non-sequential key should convert to hash
	arr.Set(NewInt(5), NewInt(50))

	if arr.IsPacked() {
		t.Error("Array should have converted to hash table")
	}

	val, _ := arr.Get(NewInt(5))
	if val.ToInt() != 50 {
		t.Errorf("Expected 50, got %d", val.ToInt())
	}
}

func TestArraySetStringKey(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewString("name"), NewString("Alice"))

	if arr.IsPacked() {
		t.Error("Array with string keys should not be packed")
	}

	val, exists := arr.Get(NewString("name"))
	if !exists || val.ToString() != "Alice" {
		t.Error("String key not working correctly")
	}
}

func TestArrayKeyNormalization(t *testing.T) {
	arr := NewEmptyArray()

	// Numeric string should become integer key
	arr.Set(NewString("42"), NewString("value"))
	val, exists := arr.Get(NewInt(42))
	if !exists || val.ToString() != "value" {
		t.Error("Numeric string should be converted to integer key")
	}

	// Float should be truncated to integer
	arr2 := NewEmptyArray()
	arr2.Set(NewFloat(3.7), NewString("floatkey"))
	val, exists = arr2.Get(NewInt(3))
	if !exists || val.ToString() != "floatkey" {
		t.Error("Float key should be truncated to integer")
	}

	// Boolean keys
	arr3 := NewEmptyArray()
	arr3.Set(NewBool(true), NewString("true"))
	arr3.Set(NewBool(false), NewString("false"))

	val, _ = arr3.Get(NewInt(1))
	if val.ToString() != "true" {
		t.Error("true should become key 1")
	}

	val, _ = arr3.Get(NewInt(0))
	if val.ToString() != "false" {
		t.Error("false should become key 0")
	}

	// Null key becomes empty string
	arr4 := NewEmptyArray()
	arr4.Set(NewNull(), NewString("nullkey"))
	val, exists = arr4.Get(NewString(""))
	if !exists || val.ToString() != "nullkey" {
		t.Error("Null key should become empty string")
	}
}

// ============================================================================
// Append Tests
// ============================================================================

func TestArrayAppend(t *testing.T) {
	arr := NewEmptyArray()

	arr.Append(NewInt(10))
	arr.Append(NewInt(20))
	arr.Append(NewInt(30))

	if arr.Len() != 3 {
		t.Errorf("Expected length 3, got %d", arr.Len())
	}

	if !arr.IsPacked() {
		t.Error("Append should maintain packed mode")
	}

	val, _ := arr.Get(NewInt(2))
	if val.ToInt() != 30 {
		t.Errorf("Expected 30, got %d", val.ToInt())
	}
}

func TestArrayAppendAfterStringKey(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewString("key"), NewString("value"))
	arr.Append(NewInt(100))

	val, exists := arr.Get(NewInt(0))
	if !exists || val.ToInt() != 100 {
		t.Error("Append should use next available integer index")
	}
}

// ============================================================================
// Unset Tests
// ============================================================================

func TestArrayUnset(t *testing.T) {
	arr := NewEmptyArray()
	arr.Append(NewInt(10))
	arr.Append(NewInt(20))
	arr.Append(NewInt(30))

	arr.Unset(NewInt(1))

	if arr.Len() != 2 {
		t.Errorf("Expected length 2, got %d", arr.Len())
	}

	_, exists := arr.Get(NewInt(1))
	if exists {
		t.Error("Key 1 should have been removed")
	}

	// Unset should convert packed to hash
	if arr.IsPacked() {
		t.Error("Unset should convert to hash table")
	}
}

// ============================================================================
// Push/Pop Tests
// ============================================================================

func TestArrayPush(t *testing.T) {
	arr := NewEmptyArray()
	length := arr.Push(NewInt(1), NewInt(2), NewInt(3))

	if length != 3 {
		t.Errorf("Expected length 3, got %d", length)
	}

	val, _ := arr.Get(NewInt(2))
	if val.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", val.ToInt())
	}
}

func TestArrayPop(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(1), NewInt(2), NewInt(3))

	val, exists := arr.Pop()
	if !exists || val.ToInt() != 3 {
		t.Error("Pop should return last element (3)")
	}

	if arr.Len() != 2 {
		t.Errorf("Expected length 2 after pop, got %d", arr.Len())
	}

	// Pop from empty
	empty := NewEmptyArray()
	val, exists = empty.Pop()
	if exists {
		t.Error("Pop from empty should return false")
	}
}

// ============================================================================
// Shift/Unshift Tests
// ============================================================================

func TestArrayShift(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(1), NewInt(2), NewInt(3))

	val, exists := arr.Shift()
	if !exists || val.ToInt() != 1 {
		t.Error("Shift should return first element (1)")
	}

	if arr.Len() != 2 {
		t.Errorf("Expected length 2 after shift, got %d", arr.Len())
	}
}

func TestArrayUnshift(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(3), NewInt(4))

	length := arr.Unshift(NewInt(1), NewInt(2))

	if length != 4 {
		t.Errorf("Expected length 4, got %d", length)
	}

	val, _ := arr.Get(NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first element to be 1, got %d", val.ToInt())
	}
}

// ============================================================================
// Slice Tests
// ============================================================================

func TestArraySlice(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(0), NewInt(1), NewInt(2), NewInt(3), NewInt(4))

	// Slice [1:3]
	sliced := arr.Slice(1, 2)

	if sliced.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sliced.Len())
	}

	val, _ := sliced.Get(NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first element 1, got %d", val.ToInt())
	}
}

func TestArraySliceNegativeOffset(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(0), NewInt(1), NewInt(2), NewInt(3), NewInt(4))

	// Slice from end
	sliced := arr.Slice(-2, 2)

	if sliced.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sliced.Len())
	}

	val, _ := sliced.Get(NewInt(0))
	if val.ToInt() != 3 {
		t.Errorf("Expected 3, got %d", val.ToInt())
	}
}

// ============================================================================
// Merge Tests
// ============================================================================

func TestArrayMerge(t *testing.T) {
	arr1 := NewEmptyArray()
	arr1.Push(NewInt(1), NewInt(2))

	arr2 := NewEmptyArray()
	arr2.Push(NewInt(3), NewInt(4))

	merged := arr1.Merge(arr2)

	if merged.Len() != 4 {
		t.Errorf("Expected length 4, got %d", merged.Len())
	}

	val, _ := merged.Get(NewInt(3))
	if val.ToInt() != 4 {
		t.Errorf("Expected 4, got %d", val.ToInt())
	}
}

func TestArrayMergeWithStringKeys(t *testing.T) {
	arr1 := NewEmptyArray()
	arr1.Set(NewString("a"), NewInt(1))

	arr2 := NewEmptyArray()
	arr2.Set(NewString("b"), NewInt(2))

	merged := arr1.Merge(arr2)

	if merged.Len() != 2 {
		t.Errorf("Expected length 2, got %d", merged.Len())
	}

	val, exists := merged.Get(NewString("b"))
	if !exists || val.ToInt() != 2 {
		t.Error("String key 'b' should exist in merged array")
	}
}

// ============================================================================
// Keys/Values Tests
// ============================================================================

func TestArrayKeys(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewString("a"))
	arr.Set(NewInt(2), NewString("b"))
	arr.Set(NewString("key"), NewString("c"))

	keys := arr.Keys()

	if keys.Len() != 3 {
		t.Errorf("Expected 3 keys, got %d", keys.Len())
	}
}

func TestArrayValues(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewString("a"))
	arr.Set(NewInt(2), NewString("b"))
	arr.Set(NewString("key"), NewString("c"))

	values := arr.Values()

	if values.Len() != 3 {
		t.Errorf("Expected 3 values, got %d", values.Len())
	}

	// Values should be reindexed as 0, 1, 2
	val, _ := values.Get(NewInt(0))
	if val.ToString() != "a" {
		t.Errorf("Expected 'a', got '%s'", val.ToString())
	}
}

// ============================================================================
// Contains/Search Tests
// ============================================================================

func TestArrayContains(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(10), NewInt(20), NewInt(30))

	if !arr.Contains(NewInt(20)) {
		t.Error("Array should contain 20")
	}

	if arr.Contains(NewInt(99)) {
		t.Error("Array should not contain 99")
	}
}

func TestArraySearch(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(10), NewInt(20), NewInt(30))

	key, found := arr.Search(NewInt(20))
	if !found {
		t.Error("Should find value 20")
	}
	if key.ToInt() != 1 {
		t.Errorf("Expected key 1, got %d", key.ToInt())
	}

	// Not found
	key, found = arr.Search(NewInt(99))
	if found {
		t.Error("Should not find value 99")
	}
}

func TestArrayHasKey(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewString("name"), NewString("Alice"))
	arr.Set(NewInt(0), NewInt(42))

	if !arr.HasKey(NewString("name")) {
		t.Error("Should have key 'name'")
	}

	if !arr.HasKey(NewInt(0)) {
		t.Error("Should have key 0")
	}

	if arr.HasKey(NewString("missing")) {
		t.Error("Should not have key 'missing'")
	}
}

// ============================================================================
// Iteration Tests
// ============================================================================

func TestArrayEach(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(10), NewInt(20), NewInt(30))

	sum := int64(0)
	arr.Each(func(key, value *Value) bool {
		sum += value.ToInt()
		return true // Continue iteration
	})

	if sum != 60 {
		t.Errorf("Expected sum 60, got %d", sum)
	}
}

func TestArrayEachBreak(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(10), NewInt(20), NewInt(30))

	count := 0
	arr.Each(func(key, value *Value) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})

	if count != 2 {
		t.Errorf("Expected 2 iterations, got %d", count)
	}
}

// ============================================================================
// DeepCopy Tests
// ============================================================================

func TestArrayDeepCopy(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(1), NewInt(2), NewInt(3))

	copied := arr.DeepCopy()

	if copied.Len() != arr.Len() {
		t.Error("Copy should have same length")
	}

	// Modify copy
	copied.Set(NewInt(0), NewInt(999))

	// Original should be unchanged
	val, _ := arr.Get(NewInt(0))
	if val.ToInt() != 1 {
		t.Error("Original array should be unchanged")
	}
}

func TestArrayDeepCopyHash(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewString("a"), NewInt(1))
	arr.Set(NewString("b"), NewInt(2))

	copied := arr.DeepCopy()

	if copied.IsPacked() == arr.IsPacked() && !arr.IsPacked() {
		// Both should be hash tables
		val, _ := copied.Get(NewString("a"))
		if val.ToInt() != 1 {
			t.Error("Copy should preserve hash table structure")
		}
	}
}

// ============================================================================
// Reset Tests
// ============================================================================

func TestArrayReset(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(1), NewInt(2), NewInt(3))

	arr.Reset()

	if arr.Len() != 0 {
		t.Error("Reset array should be empty")
	}

	if !arr.IsPacked() {
		t.Error("Reset array should be packed")
	}
}

// ============================================================================
// Packed Array Optimization Tests
// ============================================================================

func TestPackedArrayStaysPacked(t *testing.T) {
	arr := NewEmptyArray()

	// Sequential integer appends should stay packed
	for i := 0; i < 100; i++ {
		arr.Append(NewInt(int64(i)))
	}

	if !arr.IsPacked() {
		t.Error("Sequential appends should maintain packed optimization")
	}

	if arr.Len() != 100 {
		t.Errorf("Expected length 100, got %d", arr.Len())
	}
}

func TestPackedArrayConvertsOnStringKey(t *testing.T) {
	arr := NewEmptyArray()
	arr.Append(NewInt(1))
	arr.Append(NewInt(2))

	if !arr.IsPacked() {
		t.Error("Array should start packed")
	}

	// Adding string key should convert to hash
	arr.Set(NewString("key"), NewInt(3))

	if arr.IsPacked() {
		t.Error("Array should convert to hash table on string key")
	}
}

func TestPackedArrayConvertsOnNonSequential(t *testing.T) {
	arr := NewEmptyArray()
	arr.Append(NewInt(1))
	arr.Append(NewInt(2))

	// Skip index 2, set index 10
	arr.Set(NewInt(10), NewInt(999))

	if arr.IsPacked() {
		t.Error("Array should convert to hash on non-sequential index")
	}

	val, _ := arr.Get(NewInt(10))
	if val.ToInt() != 999 {
		t.Error("Non-sequential index should work")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestArrayNilOperations(t *testing.T) {
	var arr *Array

	if arr.Len() != 0 {
		t.Error("Nil array should have length 0")
	}

	if !arr.IsEmpty() {
		t.Error("Nil array should be empty")
	}

	val, exists := arr.Get(NewInt(0))
	if exists {
		t.Error("Get on nil array should return not exists")
	}
	if val.Type() != TypeNull {
		t.Error("Get on nil array should return null")
	}
}

func TestArrayEmptyOperations(t *testing.T) {
	arr := NewEmptyArray()

	_, exists := arr.Pop()
	if exists {
		t.Error("Pop on empty array should return false")
	}

	_, exists = arr.Shift()
	if exists {
		t.Error("Shift on empty array should return false")
	}

	if arr.Contains(NewInt(1)) {
		t.Error("Empty array should not contain anything")
	}
}

func TestArrayStringRepresentation(t *testing.T) {
	arr := NewEmptyArray()
	arr.Push(NewInt(1), NewInt(2))

	str := arr.String()
	if str == "" {
		t.Error("String representation should not be empty")
	}

	// Empty array
	empty := NewEmptyArray()
	if empty.String() != "[]" {
		t.Errorf("Empty array string should be '[]', got '%s'", empty.String())
	}
}

func TestArrayUpdateExisting(t *testing.T) {
	arr := NewEmptyArray()
	arr.Set(NewInt(0), NewInt(100))
	arr.Set(NewInt(0), NewInt(200)) // Update

	val, _ := arr.Get(NewInt(0))
	if val.ToInt() != 200 {
		t.Error("Update should replace existing value")
	}

	if arr.Len() != 1 {
		t.Error("Update should not increase length")
	}
}

func TestArrayMixedKeys(t *testing.T) {
	arr := NewEmptyArray()

	// Mix of integer and string keys
	arr.Set(NewInt(0), NewString("zero"))
	arr.Set(NewString("name"), NewString("Alice"))
	arr.Set(NewInt(1), NewString("one"))
	arr.Set(NewString("age"), NewInt(30))

	if arr.Len() != 4 {
		t.Errorf("Expected length 4, got %d", arr.Len())
	}

	// Check all values
	val, _ := arr.Get(NewInt(0))
	if val.ToString() != "zero" {
		t.Error("Integer key 0 failed")
	}

	val, _ = arr.Get(NewString("name"))
	if val.ToString() != "Alice" {
		t.Error("String key 'name' failed")
	}
}
