package array

import (
	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Array Size Functions
// ============================================================================

// Count returns the number of elements in an array
// count(array $array, int $mode = COUNT_NORMAL): int
func Count(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewInt(0)
	}

	arrayData := arr.ToArray()
	return types.NewInt(int64(arrayData.Len()))
}

// Sizeof is an alias for Count
func Sizeof(arr *types.Value) *types.Value {
	return Count(arr)
}

// ============================================================================
// Array Keys and Values
// ============================================================================

// ArrayKeys returns an array of keys from the input array
// array_keys(array $array): array
func ArrayKeys(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	keys := arrayData.Keys()
	return types.NewArray(keys)
}

// ArrayValues returns an array of values from the input array (reindexed)
// array_values(array $array): array
func ArrayValues(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	values := arrayData.Values()
	return types.NewArray(values)
}

// ============================================================================
// Array Stack Operations
// ============================================================================

// ArrayPush appends one or more elements to the end of an array
// array_push(array &$array, mixed ...$values): int
func ArrayPush(arr *types.Value, values ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewInt(0)
	}

	arrayData := arr.ToArray()
	length := arrayData.Push(values...)
	return types.NewInt(int64(length))
}

// ArrayPop removes and returns the last element of an array
// array_pop(array &$array): mixed
func ArrayPop(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewNull()
	}

	arrayData := arr.ToArray()
	value, exists := arrayData.Pop()
	if !exists {
		return types.NewNull()
	}
	return value
}

// ArrayShift removes and returns the first element of an array
// array_shift(array &$array): mixed
func ArrayShift(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewNull()
	}

	arrayData := arr.ToArray()
	value, exists := arrayData.Shift()
	if !exists {
		return types.NewNull()
	}
	return value
}

// ArrayUnshift prepends one or more elements to the beginning of an array
// array_unshift(array &$array, mixed ...$values): int
func ArrayUnshift(arr *types.Value, values ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewInt(0)
	}

	arrayData := arr.ToArray()
	length := arrayData.Unshift(values...)
	return types.NewInt(int64(length))
}

// ============================================================================
// Array Merging
// ============================================================================

// ArrayMerge merges one or more arrays
// array_merge(array ...$arrays): array
func ArrayMerge(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	// Start with the first array
	var result *types.Array
	if arrays[0] != nil && arrays[0].Type() == types.TypeArray {
		result = arrays[0].ToArray().DeepCopy()
	} else {
		result = types.NewEmptyArray()
	}

	// Merge the rest
	for i := 1; i < len(arrays); i++ {
		if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
			other := arrays[i].ToArray()
			result = result.Merge(other)
		}
	}

	return types.NewArray(result)
}

// ============================================================================
// Array Searching
// ============================================================================

// InArray checks if a value exists in an array
// in_array(mixed $needle, array $haystack, bool $strict = false): bool
func InArray(needle *types.Value, haystack *types.Value, strict ...*types.Value) *types.Value {
	if haystack == nil || haystack.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := haystack.ToArray()

	// TODO: strict comparison mode when strict is true
	// For now, we use loose comparison (Equals)
	found := arrayData.Contains(needle)
	return types.NewBool(found)
}

// ArraySearch searches for a value in an array and returns the first key
// array_search(mixed $needle, array $haystack, bool $strict = false): int|string|false
func ArraySearch(needle *types.Value, haystack *types.Value, strict ...*types.Value) *types.Value {
	if haystack == nil || haystack.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := haystack.ToArray()

	// TODO: strict comparison mode when strict is true
	key, found := arrayData.Search(needle)
	if !found {
		return types.NewBool(false)
	}
	return key
}

// ============================================================================
// Array Slicing
// ============================================================================

// ArraySlice extracts a slice of an array
// array_slice(array $array, int $offset, ?int $length = null): array
func ArraySlice(arr *types.Value, offset *types.Value, length ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	offsetInt := int(offset.ToInt())

	var lengthInt int
	if len(length) > 0 && length[0] != nil {
		lengthInt = int(length[0].ToInt())
	} else {
		// No length specified, go to end
		lengthInt = arrayData.Len()
	}

	sliced := arrayData.Slice(offsetInt, lengthInt)
	return types.NewArray(sliced)
}

// ArraySplice removes and replaces a portion of an array
// array_splice(array &$array, int $offset, ?int $length = null, mixed $replacement = []): array
func ArraySplice(arr *types.Value, offset *types.Value, length *types.Value, replacement ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	offsetInt := int(offset.ToInt())
	lengthInt := int(length.ToInt())

	// Extract the portion to be removed
	removed := arrayData.Slice(offsetInt, lengthInt)

	// Get arrays before and after the splice point
	before := arrayData.Slice(0, offsetInt)
	after := arrayData.Slice(offsetInt+lengthInt, arrayData.Len())

	// Build the new array: before + replacement + after
	result := before

	// Add replacement elements if provided
	if len(replacement) > 0 {
		for _, val := range replacement {
			if val.Type() == types.TypeArray {
				// If replacement is an array, merge its elements
				replArray := val.ToArray()
				result = result.Merge(replArray)
			} else {
				// If replacement is a single value, append it
				result.Append(val)
			}
		}
	}

	// Add the after portion
	result = result.Merge(after)

	// Update the original array
	// In PHP, array_splice modifies the array in place
	// We need to replace the array's contents
	arrayData.Reset()
	result.Each(func(key, value *types.Value) bool {
		arrayData.Set(key, value)
		return true
	})

	// Return the removed elements
	return types.NewArray(removed)
}

// ============================================================================
// Additional Helper Functions
// ============================================================================

// ArrayReverse reverses the order of elements in an array
// array_reverse(array $array, bool $preserve_keys = false): array
func ArrayReverse(arr *types.Value, preserveKeys ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	preserve := false
	if len(preserveKeys) > 0 && preserveKeys[0] != nil {
		preserve = preserveKeys[0].ToBool()
	}

	// Collect all key-value pairs
	var pairs [][2]*types.Value
	arrayData.Each(func(key, value *types.Value) bool {
		pairs = append(pairs, [2]*types.Value{key, value})
		return true
	})

	// Add in reverse order
	for i := len(pairs) - 1; i >= 0; i-- {
		if preserve {
			result.Set(pairs[i][0], pairs[i][1])
		} else {
			result.Append(pairs[i][1])
		}
	}

	return types.NewArray(result)
}

// ArrayUnique removes duplicate values from an array
// array_unique(array $array): array
func ArrayUnique(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()
	seen := types.NewEmptyArray()

	arrayData.Each(func(key, value *types.Value) bool {
		// Check if we've seen this value before
		if !seen.Contains(value) {
			seen.Append(value)
			result.Set(key, value)
		}
		return true
	})

	return types.NewArray(result)
}

// ArrayCombine creates an array using one array for keys and another for values
// array_combine(array $keys, array $values): array
func ArrayCombine(keys *types.Value, values *types.Value) *types.Value {
	if keys == nil || keys.Type() != types.TypeArray ||
		values == nil || values.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	keysArray := keys.ToArray()
	valuesArray := values.ToArray()

	// Arrays must be same length
	if keysArray.Len() != valuesArray.Len() {
		return types.NewBool(false)
	}

	result := types.NewEmptyArray()

	// Collect keys and values
	var keyList []*types.Value
	var valueList []*types.Value

	keysArray.Each(func(_, k *types.Value) bool {
		keyList = append(keyList, k)
		return true
	})

	valuesArray.Each(func(_, v *types.Value) bool {
		valueList = append(valueList, v)
		return true
	})

	// Combine them
	for i := 0; i < len(keyList); i++ {
		result.Set(keyList[i], valueList[i])
	}

	return types.NewArray(result)
}

// ArrayFlip exchanges all keys with their associated values
// array_flip(array $array): array
func ArrayFlip(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	arrayData.Each(func(key, value *types.Value) bool {
		// In PHP, only strings and integers can be keys
		// Values become keys, keys become values
		result.Set(value, key)
		return true
	})

	return types.NewArray(result)
}

// ArrayFill fills an array with values
// array_fill(int $start_index, int $count, mixed $value): array
func ArrayFill(startIndex *types.Value, count *types.Value, value *types.Value) *types.Value {
	start := startIndex.ToInt()
	num := count.ToInt()

	if num < 0 {
		return types.NewBool(false)
	}

	result := types.NewEmptyArray()

	for i := int64(0); i < num; i++ {
		result.Set(types.NewInt(start+i), value)
	}

	return types.NewArray(result)
}

// ArrayChunk splits an array into chunks
// array_chunk(array $array, int $length, bool $preserve_keys = false): array
func ArrayChunk(arr *types.Value, length *types.Value, preserveKeys ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	chunkSize := int(length.ToInt())
	if chunkSize < 1 {
		return types.NewBool(false)
	}

	preserve := false
	if len(preserveKeys) > 0 && preserveKeys[0] != nil {
		preserve = preserveKeys[0].ToBool()
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()
	currentChunk := types.NewEmptyArray()
	count := 0

	arrayData.Each(func(key, value *types.Value) bool {
		if preserve {
			currentChunk.Set(key, value)
		} else {
			currentChunk.Append(value)
		}
		count++

		if count >= chunkSize {
			result.Append(types.NewArray(currentChunk))
			currentChunk = types.NewEmptyArray()
			count = 0
		}

		return true
	})

	// Add remaining elements
	if count > 0 {
		result.Append(types.NewArray(currentChunk))
	}

	return types.NewArray(result)
}

// ============================================================================
// Array Sorting Functions
// ============================================================================

// Sort sorts an array by values in ascending order
// sort(array &$array, int $flags = SORT_REGULAR): true
func Sort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect values
	var values []*types.Value
	arrayData.Each(func(_, value *types.Value) bool {
		values = append(values, value)
		return true
	})

	// Sort values
	sortValues(values, false)

	// Reset array and add sorted values with numeric keys
	arrayData.Reset()
	for i, val := range values {
		arrayData.Set(types.NewInt(int64(i)), val)
	}

	return types.NewBool(true)
}

// Rsort sorts an array by values in descending order
// rsort(array &$array, int $flags = SORT_REGULAR): true
func Rsort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect values
	var values []*types.Value
	arrayData.Each(func(_, value *types.Value) bool {
		values = append(values, value)
		return true
	})

	// Sort values in reverse
	sortValues(values, true)

	// Reset array and add sorted values with numeric keys
	arrayData.Reset()
	for i, val := range values {
		arrayData.Set(types.NewInt(int64(i)), val)
	}

	return types.NewBool(true)
}

// Asort sorts an array by values in ascending order, preserving keys
// asort(array &$array, int $flags = SORT_REGULAR): true
func Asort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect key-value pairs
	var pairs []struct{ key, value *types.Value }
	arrayData.Each(func(key, value *types.Value) bool {
		pairs = append(pairs, struct{ key, value *types.Value }{key, value})
		return true
	})

	// Sort by values
	sortPairsByValue(pairs, false)

	// Reset array and add sorted pairs
	arrayData.Reset()
	for _, pair := range pairs {
		arrayData.Set(pair.key, pair.value)
	}

	return types.NewBool(true)
}

// Arsort sorts an array by values in descending order, preserving keys
// arsort(array &$array, int $flags = SORT_REGULAR): true
func Arsort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect key-value pairs
	var pairs []struct{ key, value *types.Value }
	arrayData.Each(func(key, value *types.Value) bool {
		pairs = append(pairs, struct{ key, value *types.Value }{key, value})
		return true
	})

	// Sort by values in reverse
	sortPairsByValue(pairs, true)

	// Reset array and add sorted pairs
	arrayData.Reset()
	for _, pair := range pairs {
		arrayData.Set(pair.key, pair.value)
	}

	return types.NewBool(true)
}

// Ksort sorts an array by keys in ascending order
// ksort(array &$array, int $flags = SORT_REGULAR): true
func Ksort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect key-value pairs
	var pairs []struct{ key, value *types.Value }
	arrayData.Each(func(key, value *types.Value) bool {
		pairs = append(pairs, struct{ key, value *types.Value }{key, value})
		return true
	})

	// Sort by keys
	sortPairsByKey(pairs, false)

	// Reset array and add sorted pairs
	arrayData.Reset()
	for _, pair := range pairs {
		arrayData.Set(pair.key, pair.value)
	}

	return types.NewBool(true)
}

// Krsort sorts an array by keys in descending order
// krsort(array &$array, int $flags = SORT_REGULAR): true
func Krsort(arr *types.Value, flags ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()

	// Collect key-value pairs
	var pairs []struct{ key, value *types.Value }
	arrayData.Each(func(key, value *types.Value) bool {
		pairs = append(pairs, struct{ key, value *types.Value }{key, value})
		return true
	})

	// Sort by keys in reverse
	sortPairsByKey(pairs, true)

	// Reset array and add sorted pairs
	arrayData.Reset()
	for _, pair := range pairs {
		arrayData.Set(pair.key, pair.value)
	}

	return types.NewBool(true)
}

// ============================================================================
// Functional Array Functions
// ============================================================================

// ArrayMap applies a callback to the elements of an array
// array_map(callable $callback, array ...$arrays): array
func ArrayMap(callback *types.Value, arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	// For now, support single array
	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// TODO: Implement callback invocation when we have callable support
	// For now, just copy the array
	arrayData.Each(func(key, value *types.Value) bool {
		result.Append(value)
		return true
	})

	return types.NewArray(result)
}

// ArrayFilter filters elements of an array using a callback function
// array_filter(array $array, ?callable $callback = null, int $mode = 0): array
func ArrayFilter(arr *types.Value, callback ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	// If no callback, filter out false-y values
	arrayData.Each(func(key, value *types.Value) bool {
		if value.ToBool() {
			result.Set(key, value)
		}
		return true
	})

	return types.NewArray(result)
}

// ArrayReduce reduces an array to a single value using a callback
// array_reduce(array $array, callable $callback, mixed $initial = null): mixed
func ArrayReduce(arr *types.Value, callback *types.Value, initial ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		if len(initial) > 0 {
			return initial[0]
		}
		return types.NewNull()
	}

	// TODO: Implement callback invocation when we have callable support
	// For now, return the initial value or null
	if len(initial) > 0 {
		return initial[0]
	}
	return types.NewNull()
}

// ArrayWalk applies a user function to every member of an array
// array_walk(array &$array, callable $callback, mixed $arg = null): true
func ArrayWalk(arr *types.Value, callback *types.Value, arg ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	// TODO: Implement callback invocation when we have callable support
	// For now, just return true
	return types.NewBool(true)
}

// ============================================================================
// Array Set Operations
// ============================================================================

// ArrayDiff computes the difference of arrays
// array_diff(array $array, array ...$arrays): array
func ArrayDiff(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add values from base that don't appear in other arrays
	base.Each(func(key, value *types.Value) bool {
		found := false

		// Check if value exists in any of the other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
				other := arrays[i].ToArray()
				if other.Contains(value) {
					found = true
					break
				}
			}
		}

		if !found {
			result.Set(key, value)
		}
		return true
	})

	return types.NewArray(result)
}

// ArrayIntersect computes the intersection of arrays
// array_intersect(array $array, array ...$arrays): array
func ArrayIntersect(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add values from base that appear in all other arrays
	base.Each(func(key, value *types.Value) bool {
		inAll := true

		// Check if value exists in all other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] == nil || arrays[i].Type() != types.TypeArray {
				inAll = false
				break
			}

			other := arrays[i].ToArray()
			if !other.Contains(value) {
				inAll = false
				break
			}
		}

		if inAll {
			result.Set(key, value)
		}
		return true
	})

	return types.NewArray(result)
}

// ============================================================================
// Array Pointer Functions
// ============================================================================

// Note: PHP's array pointer functions maintain internal state within the array.
// Since our Array implementation doesn't expose pointer methods,
// these functions provide simplified implementations that work with first/last elements.

// Current returns the first element in an array
// current(array $array): mixed
func Current(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewBool(false)
	}

	// Return first element
	var firstValue *types.Value
	arrayData.Each(func(_, value *types.Value) bool {
		firstValue = value
		return false // Stop after first element
	})

	if firstValue == nil {
		return types.NewBool(false)
	}
	return firstValue
}

// Key returns the first key of an array
// key(array $array): int|string|null
func Key(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewNull()
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewNull()
	}

	// Return first key
	var firstKey *types.Value
	arrayData.Each(func(key, _ *types.Value) bool {
		firstKey = key
		return false // Stop after first element
	})

	if firstKey == nil {
		return types.NewNull()
	}
	return firstKey
}

// Reset sets the internal pointer of an array to its first element
// reset(array &$array): mixed
func Reset(arr *types.Value) *types.Value {
	// For our implementation, this is the same as Current
	return Current(arr)
}

// End returns the last element of an array
// end(array &$array): mixed
func End(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewBool(false)
	}

	// Move to the last element
	var lastValue *types.Value
	arrayData.Each(func(_, value *types.Value) bool {
		lastValue = value
		return true
	})

	if lastValue == nil {
		return types.NewBool(false)
	}
	return lastValue
}

// Next advances the internal array pointer (simplified implementation)
// next(array &$array): mixed
func Next(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()
	if arrayData.Len() < 2 {
		return types.NewBool(false)
	}

	// Return second element as a simplified "next"
	count := 0
	var nextValue *types.Value
	arrayData.Each(func(_, value *types.Value) bool {
		count++
		if count == 2 {
			nextValue = value
			return false
		}
		return true
	})

	if nextValue == nil {
		return types.NewBool(false)
	}
	return nextValue
}

// Prev rewinds the internal array pointer (simplified implementation)
// prev(array &$array): mixed
func Prev(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	// Without pointer state, we can't implement prev meaningfully
	// Return false to indicate no previous element
	return types.NewBool(false)
}

// ============================================================================
// Helper Functions
// ============================================================================

// sortValues sorts a slice of values in place
func sortValues(values []*types.Value, reverse bool) {
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if compareValues(values[i], values[j], reverse) {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}

// sortPairsByValue sorts pairs by their value
func sortPairsByValue(pairs []struct{ key, value *types.Value }, reverse bool) {
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if compareValues(pairs[i].value, pairs[j].value, reverse) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
}

// sortPairsByKey sorts pairs by their key
func sortPairsByKey(pairs []struct{ key, value *types.Value }, reverse bool) {
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if compareValues(pairs[i].key, pairs[j].key, reverse) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
}

// compareValues compares two values for sorting
// Returns true if a should come after b
func compareValues(a, b *types.Value, reverse bool) bool {
	if a == nil || b == nil {
		return false
	}

	// Compare based on type
	aType := a.Type()
	bType := b.Type()

	// Numbers
	if aType == types.TypeInt || aType == types.TypeFloat {
		if bType == types.TypeInt || bType == types.TypeFloat {
			aNum := a.ToFloat()
			bNum := b.ToFloat()
			if reverse {
				return aNum < bNum
			}
			return aNum > bNum
		}
	}

	// Strings
	if aType == types.TypeString {
		if bType == types.TypeString {
			aStr := a.ToString()
			bStr := b.ToString()
			if reverse {
				return aStr < bStr
			}
			return aStr > bStr
		}
	}

	// Default: compare as strings
	aStr := a.ToString()
	bStr := b.ToString()
	if reverse {
		return aStr < bStr
	}
	return aStr > bStr
}
