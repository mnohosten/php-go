package array

import (
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/util"
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

	// Security: Safe conversion from int64 to int
	offsetInt, err := util.SafeConvertToInt(offset.ToInt(), "offset")
	if err != nil {
		// Invalid offset - return empty array
		return types.NewArray(types.NewEmptyArray())
	}

	var lengthInt int
	if len(length) > 0 && length[0] != nil {
		// Security: Safe conversion from int64 to int
		lengthInt, err = util.SafeConvertToInt(length[0].ToInt(), "length")
		if err != nil {
			// Invalid length - use remaining length
			lengthInt = arrayData.Len()
		}
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

	// Security: Safe conversion from int64 to int
	offsetInt, err := util.SafeConvertToInt(offset.ToInt(), "offset")
	if err != nil {
		// Invalid offset - return empty array
		return types.NewArray(types.NewEmptyArray())
	}

	lengthInt, err := util.SafeConvertToInt(length.ToInt(), "length")
	if err != nil {
		// Invalid length - use remaining length
		lengthInt = arrayData.Len()
	}

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

	// Security: Safe conversion from int64 to int
	chunkSize, err := util.SafeConvertToInt(length.ToInt(), "chunk_size")
	if err != nil || chunkSize < 1 {
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

// ArrayWalkRecursive applies a user function recursively to every member of an array
// array_walk_recursive(array &$array, callable $callback, mixed $arg = null): true
func ArrayWalkRecursive(arr *types.Value, callback *types.Value, arg ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	// TODO: Implement callback invocation when we have callable support
	// This function differs from array_walk by recursively descending into nested arrays
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

// ============================================================================
// Array Padding and Computation Functions
// ============================================================================

// ArrayPad pads an array to the specified length with a value
// array_pad(array $array, int $length, mixed $value): array
func ArrayPad(arr *types.Value, length *types.Value, value *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	padSize, err := util.SafeConvertToInt(length.ToInt(), "length")
	if err != nil {
		return types.NewArray(types.NewEmptyArray())
	}

	currentLen := arrayData.Len()

	// If the absolute value of padSize is less than or equal to current length,
	// no padding is done
	absSize := padSize
	if absSize < 0 {
		absSize = -absSize
	}
	if absSize <= currentLen {
		// Return a copy of the original array
		result := types.NewEmptyArray()
		arrayData.Each(func(key, val *types.Value) bool {
			result.Set(key, val)
			return true
		})
		return types.NewArray(result)
	}

	padCount := absSize - currentLen
	result := types.NewEmptyArray()

	if padSize < 0 {
		// Pad at the beginning
		for i := 0; i < padCount; i++ {
			result.Append(value)
		}
		arrayData.Each(func(_, val *types.Value) bool {
			result.Append(val)
			return true
		})
	} else {
		// Pad at the end
		arrayData.Each(func(_, val *types.Value) bool {
			result.Append(val)
			return true
		})
		for i := 0; i < padCount; i++ {
			result.Append(value)
		}
	}

	return types.NewArray(result)
}

// ArraySum calculates the sum of values in an array
// array_sum(array $array): int|float
func ArraySum(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewInt(0)
	}

	arrayData := arr.ToArray()
	var intSum int64 = 0
	var floatSum float64 = 0
	hasFloat := false

	arrayData.Each(func(_, value *types.Value) bool {
		if value == nil {
			return true
		}

		switch value.Type() {
		case types.TypeInt:
			intSum += value.ToInt()
		case types.TypeFloat:
			hasFloat = true
			floatSum += value.ToFloat()
		case types.TypeString:
			// PHP converts strings to numbers using ToInt/ToFloat
			// which already handles numeric string detection
			floatVal := value.ToFloat()
			if floatVal != float64(int64(floatVal)) {
				hasFloat = true
				floatSum += floatVal
			} else {
				intSum += value.ToInt()
			}
		case types.TypeBool:
			if value.ToBool() {
				intSum += 1
			}
		}
		return true
	})

	if hasFloat {
		return types.NewFloat(float64(intSum) + floatSum)
	}
	return types.NewInt(intSum)
}

// ArrayProduct calculates the product of values in an array
// array_product(array $array): int|float
func ArrayProduct(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewInt(0)
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewInt(1) // PHP returns 1 for empty array
	}

	var intProduct int64 = 1
	var floatProduct float64 = 1
	hasFloat := false

	arrayData.Each(func(_, value *types.Value) bool {
		if value == nil {
			return true
		}

		switch value.Type() {
		case types.TypeInt:
			intProduct *= value.ToInt()
		case types.TypeFloat:
			hasFloat = true
			floatProduct *= value.ToFloat()
		case types.TypeString:
			// PHP converts strings to numbers using ToInt/ToFloat
			floatVal := value.ToFloat()
			if floatVal != float64(int64(floatVal)) {
				hasFloat = true
				floatProduct *= floatVal
			} else {
				intVal := value.ToInt()
				if intVal == 0 && value.ToString() != "0" {
					// Non-numeric strings are treated as 0
					intProduct = 0
				} else {
					intProduct *= intVal
				}
			}
		case types.TypeBool:
			if value.ToBool() {
				// true = 1, doesn't change product
			} else {
				intProduct = 0
			}
		default:
			// Other types treated as 0
			intProduct = 0
		}
		return true
	})

	if hasFloat {
		return types.NewFloat(float64(intProduct) * floatProduct)
	}
	return types.NewInt(intProduct)
}

// ArrayColumn returns the values from a single column in a multi-dimensional array
// array_column(array $array, int|string|null $column_key, int|string|null $index_key = null): array
func ArrayColumn(arr *types.Value, columnKey *types.Value, indexKey ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	var idxKey *types.Value
	if len(indexKey) > 0 && indexKey[0] != nil && indexKey[0].Type() != types.TypeNull {
		idxKey = indexKey[0]
	}

	arrayData.Each(func(_, row *types.Value) bool {
		if row == nil || row.Type() != types.TypeArray {
			return true // Skip non-array rows
		}

		rowArray := row.ToArray()

		// Get the value for this row
		var value *types.Value
		if columnKey == nil || columnKey.Type() == types.TypeNull {
			// If column_key is null, return the entire row
			value = row
		} else {
			// Get the specific column
			value, _ = rowArray.Get(columnKey)
			if value == nil {
				return true // Skip if column doesn't exist
			}
		}

		// Determine the key for the result
		if idxKey != nil {
			// Use a specific column as the index
			keyValue, exists := rowArray.Get(idxKey)
			if exists && keyValue != nil {
				result.Set(keyValue, value)
			}
		} else {
			// Use sequential numeric keys
			result.Append(value)
		}

		return true
	})

	return types.NewArray(result)
}

// CASE_LOWER and CASE_UPPER constants for array_change_key_case
const (
	CASE_LOWER = 0
	CASE_UPPER = 1
)

// ArrayChangeKeyCase changes the case of all keys in an array
// array_change_key_case(array $array, int $case = CASE_LOWER): array
func ArrayChangeKeyCase(arr *types.Value, caseType ...*types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	toUpper := false
	if len(caseType) > 0 && caseType[0] != nil {
		toUpper = caseType[0].ToInt() == CASE_UPPER
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	arrayData.Each(func(key, value *types.Value) bool {
		if key.Type() == types.TypeString {
			keyStr := key.ToString()
			if toUpper {
				keyStr = toUpperString(keyStr)
			} else {
				keyStr = toLowerString(keyStr)
			}
			result.Set(types.NewString(keyStr), value)
		} else {
			// Non-string keys (integers) are kept as-is
			result.Set(key, value)
		}
		return true
	})

	return types.NewArray(result)
}

// ArrayReplace replaces elements from passed arrays into the first array
// array_replace(array $array, array ...$replacements): array
func ArrayReplace(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	// Start with a copy of the first array
	var result *types.Array
	if arrays[0] != nil && arrays[0].Type() == types.TypeArray {
		result = arrays[0].ToArray().DeepCopy()
	} else {
		result = types.NewEmptyArray()
	}

	// Replace with subsequent arrays
	for i := 1; i < len(arrays); i++ {
		if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
			replacement := arrays[i].ToArray()
			replacement.Each(func(key, value *types.Value) bool {
				result.Set(key, value)
				return true
			})
		}
	}

	return types.NewArray(result)
}

// ArrayReplaceRecursive replaces elements from passed arrays into the first array recursively
// array_replace_recursive(array $array, array ...$replacements): array
func ArrayReplaceRecursive(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	// Start with a deep copy of the first array
	var result *types.Array
	if arrays[0] != nil && arrays[0].Type() == types.TypeArray {
		result = arrays[0].ToArray().DeepCopy()
	} else {
		result = types.NewEmptyArray()
	}

	// Replace with subsequent arrays recursively
	for i := 1; i < len(arrays); i++ {
		if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
			replacement := arrays[i].ToArray()
			result = arrayReplaceRecursiveHelper(result, replacement)
		}
	}

	return types.NewArray(result)
}

// arrayReplaceRecursiveHelper merges arrays recursively
func arrayReplaceRecursiveHelper(base, replacement *types.Array) *types.Array {
	replacement.Each(func(key, value *types.Value) bool {
		existingValue, exists := base.Get(key)

		if exists && existingValue != nil &&
			existingValue.Type() == types.TypeArray &&
			value != nil && value.Type() == types.TypeArray {
			// Both are arrays, recurse
			merged := arrayReplaceRecursiveHelper(
				existingValue.ToArray().DeepCopy(),
				value.ToArray(),
			)
			base.Set(key, types.NewArray(merged))
		} else {
			// Replace the value
			base.Set(key, value)
		}
		return true
	})

	return base
}

// ============================================================================
// String Case Helpers
// ============================================================================

// toLowerString converts a string to lowercase
func toLowerString(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// toUpperString converts a string to uppercase
func toUpperString(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// ============================================================================
// Array Search/Filter Functions
// ============================================================================

// ArrayKeyExists checks if the given key or index exists in the array
// array_key_exists(int|string $key, array $array): bool
func ArrayKeyExists(key *types.Value, arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewBool(false)
	}

	if key == nil {
		return types.NewBool(false)
	}

	arrayData := arr.ToArray()
	_, exists := arrayData.Get(key)
	return types.NewBool(exists)
}

// KeyExists is an alias for ArrayKeyExists
// key_exists(int|string $key, array $array): bool
func KeyExists(key *types.Value, arr *types.Value) *types.Value {
	return ArrayKeyExists(key, arr)
}

// ArrayDiffKey computes the difference of arrays using keys for comparison
// array_diff_key(array $array, array ...$arrays): array
func ArrayDiffKey(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add key-value pairs from base whose keys don't appear in other arrays
	base.Each(func(key, value *types.Value) bool {
		found := false

		// Check if key exists in any of the other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
				other := arrays[i].ToArray()
				_, exists := other.Get(key)
				if exists {
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

// ArrayIntersectKey computes the intersection of arrays using keys for comparison
// array_intersect_key(array $array, array ...$arrays): array
func ArrayIntersectKey(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add key-value pairs from base whose keys appear in all other arrays
	base.Each(func(key, value *types.Value) bool {
		inAll := true

		// Check if key exists in all other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] == nil || arrays[i].Type() != types.TypeArray {
				inAll = false
				break
			}

			other := arrays[i].ToArray()
			_, exists := other.Get(key)
			if !exists {
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

// ArrayDiffAssoc computes the difference of arrays with additional index check
// array_diff_assoc(array $array, array ...$arrays): array
func ArrayDiffAssoc(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add key-value pairs from base where key=>value doesn't appear in other arrays
	base.Each(func(key, value *types.Value) bool {
		found := false

		// Check if key=>value exists in any of the other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] != nil && arrays[i].Type() == types.TypeArray {
				other := arrays[i].ToArray()
				otherValue, exists := other.Get(key)
				// Both key must exist and values must be equal
				if exists && otherValue != nil && value.Equals(otherValue) {
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

// ArrayIntersectAssoc computes the intersection of arrays with additional index check
// array_intersect_assoc(array $array, array ...$arrays): array
func ArrayIntersectAssoc(arrays ...*types.Value) *types.Value {
	if len(arrays) == 0 {
		return types.NewArray(types.NewEmptyArray())
	}

	if arrays[0] == nil || arrays[0].Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	base := arrays[0].ToArray()
	result := types.NewEmptyArray()

	// Add key-value pairs from base where key=>value appears in all other arrays
	base.Each(func(key, value *types.Value) bool {
		inAll := true

		// Check if key=>value exists in all other arrays
		for i := 1; i < len(arrays); i++ {
			if arrays[i] == nil || arrays[i].Type() != types.TypeArray {
				inAll = false
				break
			}

			other := arrays[i].ToArray()
			otherValue, exists := other.Get(key)
			// Both key must exist and values must be equal
			if !exists || otherValue == nil || !value.Equals(otherValue) {
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

// ArrayKeyFirst gets the first key of an array
// array_key_first(array $array): int|string|null
func ArrayKeyFirst(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewNull()
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewNull()
	}

	var firstKey *types.Value
	arrayData.Each(func(key, _ *types.Value) bool {
		firstKey = key
		return false // Stop after first
	})

	if firstKey == nil {
		return types.NewNull()
	}
	return firstKey
}

// ArrayKeyLast gets the last key of an array
// array_key_last(array $array): int|string|null
func ArrayKeyLast(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewNull()
	}

	arrayData := arr.ToArray()
	if arrayData.Len() == 0 {
		return types.NewNull()
	}

	var lastKey *types.Value
	arrayData.Each(func(key, _ *types.Value) bool {
		lastKey = key
		return true // Continue to get last
	})

	if lastKey == nil {
		return types.NewNull()
	}
	return lastKey
}

// ArrayCountValues counts all the values of an array
// array_count_values(array $array): array
func ArrayCountValues(arr *types.Value) *types.Value {
	if arr == nil || arr.Type() != types.TypeArray {
		return types.NewArray(types.NewEmptyArray())
	}

	arrayData := arr.ToArray()
	result := types.NewEmptyArray()

	arrayData.Each(func(_, value *types.Value) bool {
		if value == nil {
			return true
		}

		// Only count string and integer values (PHP behavior)
		if value.Type() != types.TypeString && value.Type() != types.TypeInt {
			return true
		}

		// Get current count for this value
		currentCount, exists := result.Get(value)
		if exists {
			result.Set(value, types.NewInt(currentCount.ToInt()+1))
		} else {
			result.Set(value, types.NewInt(1))
		}

		return true
	})

	return types.NewArray(result)
}
