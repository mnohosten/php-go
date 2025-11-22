package array

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Count/Sizeof Tests
// ============================================================================

func TestCount(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	result := Count(arrVal)
	if result.ToInt() != 3 {
		t.Errorf("Expected count 3, got %d", result.ToInt())
	}
}

func TestCountEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	arrVal := types.NewArray(arr)

	result := Count(arrVal)
	if result.ToInt() != 0 {
		t.Errorf("Expected count 0, got %d", result.ToInt())
	}
}

func TestSizeof(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2))
	arrVal := types.NewArray(arr)

	result := Sizeof(arrVal)
	if result.ToInt() != 2 {
		t.Errorf("Expected sizeof 2, got %d", result.ToInt())
	}
}

// ============================================================================
// ArrayKeys/ArrayValues Tests
// ============================================================================

func TestArrayKeys(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(2))
	arr.Set(types.NewInt(0), types.NewInt(3))
	arrVal := types.NewArray(arr)

	keys := ArrayKeys(arrVal)
	keysArray := keys.ToArray()

	if keysArray.Len() != 3 {
		t.Errorf("Expected 3 keys, got %d", keysArray.Len())
	}
}

func TestArrayValues(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(2))
	arrVal := types.NewArray(arr)

	values := ArrayValues(arrVal)
	valuesArray := values.ToArray()

	if valuesArray.Len() != 2 {
		t.Errorf("Expected 2 values, got %d", valuesArray.Len())
	}

	// Values should be reindexed starting from 0
	val, _ := valuesArray.Get(types.NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first value 1, got %d", val.ToInt())
	}
}

// ============================================================================
// Stack Operations Tests
// ============================================================================

func TestArrayPush(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arrVal := types.NewArray(arr)

	length := ArrayPush(arrVal, types.NewInt(2), types.NewInt(3))

	if length.ToInt() != 3 {
		t.Errorf("Expected length 3, got %d", length.ToInt())
	}

	val, _ := arr.Get(types.NewInt(2))
	if val.ToInt() != 3 {
		t.Errorf("Expected value 3, got %d", val.ToInt())
	}
}

func TestArrayPop(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	popped := ArrayPop(arrVal)

	if popped.ToInt() != 3 {
		t.Errorf("Expected popped value 3, got %d", popped.ToInt())
	}

	if arr.Len() != 2 {
		t.Errorf("Expected length 2 after pop, got %d", arr.Len())
	}
}

func TestArrayShift(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	shifted := ArrayShift(arrVal)

	if shifted.ToInt() != 1 {
		t.Errorf("Expected shifted value 1, got %d", shifted.ToInt())
	}

	if arr.Len() != 2 {
		t.Errorf("Expected length 2 after shift, got %d", arr.Len())
	}
}

func TestArrayUnshift(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(3), types.NewInt(4))
	arrVal := types.NewArray(arr)

	length := ArrayUnshift(arrVal, types.NewInt(1), types.NewInt(2))

	if length.ToInt() != 4 {
		t.Errorf("Expected length 4, got %d", length.ToInt())
	}

	val, _ := arr.Get(types.NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first value 1, got %d", val.ToInt())
	}
}

// ============================================================================
// ArrayMerge Tests
// ============================================================================

func TestArrayMerge(t *testing.T) {
	arr1 := types.NewEmptyArray()
	arr1.Push(types.NewInt(1), types.NewInt(2))
	arrVal1 := types.NewArray(arr1)

	arr2 := types.NewEmptyArray()
	arr2.Push(types.NewInt(3), types.NewInt(4))
	arrVal2 := types.NewArray(arr2)

	merged := ArrayMerge(arrVal1, arrVal2)
	mergedArray := merged.ToArray()

	if mergedArray.Len() != 4 {
		t.Errorf("Expected merged length 4, got %d", mergedArray.Len())
	}

	val, _ := mergedArray.Get(types.NewInt(3))
	if val.ToInt() != 4 {
		t.Errorf("Expected value at index 3 to be 4, got %d", val.ToInt())
	}
}

func TestArrayMergeEmpty(t *testing.T) {
	merged := ArrayMerge()
	mergedArray := merged.ToArray()

	if mergedArray.Len() != 0 {
		t.Errorf("Expected empty array, got length %d", mergedArray.Len())
	}
}

// ============================================================================
// Search Tests
// ============================================================================

func TestInArray(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	result := InArray(types.NewInt(20), arrVal)
	if !result.ToBool() {
		t.Error("Expected in_array to find value 20")
	}

	result = InArray(types.NewInt(99), arrVal)
	if result.ToBool() {
		t.Error("Expected in_array to not find value 99")
	}
}

func TestArraySearch(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	result := ArraySearch(types.NewInt(20), arrVal)
	if result.Type() == types.TypeBool {
		t.Error("Expected array_search to find value")
	}
	if result.ToInt() != 1 {
		t.Errorf("Expected index 1, got %d", result.ToInt())
	}

	result = ArraySearch(types.NewInt(99), arrVal)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected array_search to return false for missing value")
	}
}

// ============================================================================
// Slice Tests
// ============================================================================

func TestArraySlice(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(0), types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal := types.NewArray(arr)

	sliced := ArraySlice(arrVal, types.NewInt(1), types.NewInt(3))
	slicedArray := sliced.ToArray()

	if slicedArray.Len() != 3 {
		t.Errorf("Expected sliced length 3, got %d", slicedArray.Len())
	}

	val, _ := slicedArray.Get(types.NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first element 1, got %d", val.ToInt())
	}
}

func TestArraySliceNoLength(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(0), types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal := types.NewArray(arr)

	sliced := ArraySlice(arrVal, types.NewInt(2))
	slicedArray := sliced.ToArray()

	if slicedArray.Len() != 3 {
		t.Errorf("Expected sliced length 3, got %d", slicedArray.Len())
	}
}

// ============================================================================
// ArrayReverse Tests
// ============================================================================

func TestArrayReverse(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	reversed := ArrayReverse(arrVal)
	reversedArray := reversed.ToArray()

	val, _ := reversedArray.Get(types.NewInt(0))
	if val.ToInt() != 3 {
		t.Errorf("Expected first element 3, got %d", val.ToInt())
	}

	val, _ = reversedArray.Get(types.NewInt(2))
	if val.ToInt() != 1 {
		t.Errorf("Expected last element 1, got %d", val.ToInt())
	}
}

// ============================================================================
// ArrayUnique Tests
// ============================================================================

func TestArrayUnique(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(1), types.NewInt(3), types.NewInt(2))
	arrVal := types.NewArray(arr)

	unique := ArrayUnique(arrVal)
	uniqueArray := unique.ToArray()

	if uniqueArray.Len() != 3 {
		t.Errorf("Expected unique length 3, got %d", uniqueArray.Len())
	}
}

// ============================================================================
// ArrayCombine Tests
// ============================================================================

func TestArrayCombine(t *testing.T) {
	keys := types.NewEmptyArray()
	keys.Push(types.NewString("a"), types.NewString("b"), types.NewString("c"))
	keysVal := types.NewArray(keys)

	values := types.NewEmptyArray()
	values.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	valuesVal := types.NewArray(values)

	combined := ArrayCombine(keysVal, valuesVal)
	if combined.Type() == types.TypeBool {
		t.Fatal("Expected array, got false")
	}

	combinedArray := combined.ToArray()
	val, _ := combinedArray.Get(types.NewString("b"))
	if val.ToInt() != 2 {
		t.Errorf("Expected value 2 for key 'b', got %d", val.ToInt())
	}
}

func TestArrayCombineMismatchedLength(t *testing.T) {
	keys := types.NewEmptyArray()
	keys.Push(types.NewString("a"))
	keysVal := types.NewArray(keys)

	values := types.NewEmptyArray()
	values.Push(types.NewInt(1), types.NewInt(2))
	valuesVal := types.NewArray(values)

	combined := ArrayCombine(keysVal, valuesVal)
	if combined.Type() != types.TypeBool || combined.ToBool() != false {
		t.Error("Expected false for mismatched length arrays")
	}
}

// ============================================================================
// ArrayFlip Tests
// ============================================================================

func TestArrayFlip(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(2))
	arrVal := types.NewArray(arr)

	flipped := ArrayFlip(arrVal)
	flippedArray := flipped.ToArray()

	val, exists := flippedArray.Get(types.NewInt(1))
	if !exists {
		t.Error("Expected key 1 to exist")
	}
	if val.ToString() != "a" {
		t.Errorf("Expected value 'a' for key 1, got '%s'", val.ToString())
	}
}

// ============================================================================
// ArrayFill Tests
// ============================================================================

func TestArrayFill(t *testing.T) {
	filled := ArrayFill(types.NewInt(0), types.NewInt(5), types.NewString("test"))
	filledArray := filled.ToArray()

	if filledArray.Len() != 5 {
		t.Errorf("Expected length 5, got %d", filledArray.Len())
	}

	val, _ := filledArray.Get(types.NewInt(3))
	if val.ToString() != "test" {
		t.Errorf("Expected value 'test', got '%s'", val.ToString())
	}
}

func TestArrayFillNegativeCount(t *testing.T) {
	filled := ArrayFill(types.NewInt(0), types.NewInt(-5), types.NewString("test"))
	if filled.Type() != types.TypeBool || filled.ToBool() != false {
		t.Error("Expected false for negative count")
	}
}

// ============================================================================
// ArrayChunk Tests
// ============================================================================

func TestArrayChunk(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4), types.NewInt(5))
	arrVal := types.NewArray(arr)

	chunked := ArrayChunk(arrVal, types.NewInt(2))
	chunkedArray := chunked.ToArray()

	if chunkedArray.Len() != 3 {
		t.Errorf("Expected 3 chunks, got %d", chunkedArray.Len())
	}

	// Check first chunk
	firstChunk, _ := chunkedArray.Get(types.NewInt(0))
	firstChunkArray := firstChunk.ToArray()
	if firstChunkArray.Len() != 2 {
		t.Errorf("Expected first chunk length 2, got %d", firstChunkArray.Len())
	}

	// Check last chunk (should have 1 element)
	lastChunk, _ := chunkedArray.Get(types.NewInt(2))
	lastChunkArray := lastChunk.ToArray()
	if lastChunkArray.Len() != 1 {
		t.Errorf("Expected last chunk length 1, got %d", lastChunkArray.Len())
	}
}

func TestArrayChunkInvalidLength(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2))
	arrVal := types.NewArray(arr)

	chunked := ArrayChunk(arrVal, types.NewInt(0))
	if chunked.Type() != types.TypeBool || chunked.ToBool() != false {
		t.Error("Expected false for invalid chunk length")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestCountNil(t *testing.T) {
	result := Count(nil)
	if result.ToInt() != 0 {
		t.Errorf("Expected count 0 for nil, got %d", result.ToInt())
	}
}

func TestCountNonArray(t *testing.T) {
	result := Count(types.NewInt(42))
	if result.ToInt() != 0 {
		t.Errorf("Expected count 0 for non-array, got %d", result.ToInt())
	}
}

func TestArrayPopEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	arrVal := types.NewArray(arr)

	popped := ArrayPop(arrVal)
	if popped.Type() != types.TypeNull {
		t.Error("Expected null when popping from empty array")
	}
}

func TestArrayShiftEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	arrVal := types.NewArray(arr)

	shifted := ArrayShift(arrVal)
	if shifted.Type() != types.TypeNull {
		t.Error("Expected null when shifting from empty array")
	}
}

// ============================================================================
// Sorting Tests
// ============================================================================

func TestSort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(3), types.NewInt(1), types.NewInt(2))
	arrVal := types.NewArray(arr)

	result := Sort(arrVal)
	if !result.ToBool() {
		t.Error("Expected sort to return true")
	}

	// Check values are sorted
	val, _ := arr.Get(types.NewInt(0))
	if val.ToInt() != 1 {
		t.Errorf("Expected first element 1, got %d", val.ToInt())
	}

	val, _ = arr.Get(types.NewInt(1))
	if val.ToInt() != 2 {
		t.Errorf("Expected second element 2, got %d", val.ToInt())
	}

	val, _ = arr.Get(types.NewInt(2))
	if val.ToInt() != 3 {
		t.Errorf("Expected third element 3, got %d", val.ToInt())
	}
}

func TestRsort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(3), types.NewInt(2))
	arrVal := types.NewArray(arr)

	result := Rsort(arrVal)
	if !result.ToBool() {
		t.Error("Expected rsort to return true")
	}

	// Check values are sorted in reverse
	val, _ := arr.Get(types.NewInt(0))
	if val.ToInt() != 3 {
		t.Errorf("Expected first element 3, got %d", val.ToInt())
	}

	val, _ = arr.Get(types.NewInt(2))
	if val.ToInt() != 1 {
		t.Errorf("Expected third element 1, got %d", val.ToInt())
	}
}

func TestAsort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(3))
	arr.Set(types.NewString("b"), types.NewInt(1))
	arr.Set(types.NewString("c"), types.NewInt(2))
	arrVal := types.NewArray(arr)

	result := Asort(arrVal)
	if !result.ToBool() {
		t.Error("Expected asort to return true")
	}

	// Keys should be preserved
	val, exists := arr.Get(types.NewString("b"))
	if !exists || val.ToInt() != 1 {
		t.Error("Expected key 'b' to still exist with value 1")
	}
}

func TestArsort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(3))
	arr.Set(types.NewString("c"), types.NewInt(2))
	arrVal := types.NewArray(arr)

	result := Arsort(arrVal)
	if !result.ToBool() {
		t.Error("Expected arsort to return true")
	}

	// Keys should be preserved
	val, exists := arr.Get(types.NewString("b"))
	if !exists || val.ToInt() != 3 {
		t.Error("Expected key 'b' to still exist with value 3")
	}
}

func TestKsort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("c"), types.NewInt(1))
	arr.Set(types.NewString("a"), types.NewInt(2))
	arr.Set(types.NewString("b"), types.NewInt(3))
	arrVal := types.NewArray(arr)

	result := Ksort(arrVal)
	if !result.ToBool() {
		t.Error("Expected ksort to return true")
	}

	// Values should be preserved under their keys
	val, exists := arr.Get(types.NewString("a"))
	if !exists || val.ToInt() != 2 {
		t.Error("Expected key 'a' to have value 2")
	}
}

func TestKrsort(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(2))
	arr.Set(types.NewString("c"), types.NewInt(3))
	arrVal := types.NewArray(arr)

	result := Krsort(arrVal)
	if !result.ToBool() {
		t.Error("Expected krsort to return true")
	}

	// Values should be preserved under their keys
	val, exists := arr.Get(types.NewString("c"))
	if !exists || val.ToInt() != 3 {
		t.Error("Expected key 'c' to have value 3")
	}
}

func TestSortStrings(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewString("banana"), types.NewString("apple"), types.NewString("cherry"))
	arrVal := types.NewArray(arr)

	Sort(arrVal)

	val, _ := arr.Get(types.NewInt(0))
	if val.ToString() != "apple" {
		t.Errorf("Expected first element 'apple', got '%s'", val.ToString())
	}

	val, _ = arr.Get(types.NewInt(2))
	if val.ToString() != "cherry" {
		t.Errorf("Expected last element 'cherry', got '%s'", val.ToString())
	}
}

// ============================================================================
// Functional Array Tests
// ============================================================================

func TestArrayMap(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	// For now, just test that it returns an array
	result := ArrayMap(types.NewNull(), arrVal)
	if result.Type() != types.TypeArray {
		t.Error("Expected array_map to return an array")
	}

	resultArray := result.ToArray()
	if resultArray.Len() != 3 {
		t.Errorf("Expected length 3, got %d", resultArray.Len())
	}
}

func TestArrayFilter(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(0), types.NewInt(1), types.NewInt(0), types.NewInt(2), types.NewInt(0))
	arrVal := types.NewArray(arr)

	result := ArrayFilter(arrVal)
	resultArray := result.ToArray()

	// Should filter out false-y values (0s)
	if resultArray.Len() != 2 {
		t.Errorf("Expected length 2 after filtering, got %d", resultArray.Len())
	}
}

func TestArrayFilterTruthyValues(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Set(types.NewString("a"), types.NewInt(1))
	arr.Set(types.NewString("b"), types.NewInt(0))
	arr.Set(types.NewString("c"), types.NewString("hello"))
	arr.Set(types.NewString("d"), types.NewString(""))
	arrVal := types.NewArray(arr)

	result := ArrayFilter(arrVal)
	resultArray := result.ToArray()

	// Should keep keys 'a' and 'c' (truthy values)
	val, exists := resultArray.Get(types.NewString("a"))
	if !exists || val.ToInt() != 1 {
		t.Error("Expected key 'a' with value 1")
	}

	val, exists = resultArray.Get(types.NewString("c"))
	if !exists || val.ToString() != "hello" {
		t.Error("Expected key 'c' with value 'hello'")
	}

	// Should not have keys 'b' and 'd' (falsy values)
	_, exists = resultArray.Get(types.NewString("b"))
	if exists {
		t.Error("Expected key 'b' to be filtered out")
	}
}

func TestArrayReduce(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	// Test with initial value
	result := ArrayReduce(arrVal, types.NewNull(), types.NewInt(10))
	if result.ToInt() != 10 {
		t.Errorf("Expected initial value 10, got %d", result.ToInt())
	}

	// Test without initial value
	result = ArrayReduce(arrVal, types.NewNull())
	if result.Type() != types.TypeNull {
		t.Error("Expected null when no initial value provided")
	}
}

func TestArrayWalk(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3))
	arrVal := types.NewArray(arr)

	result := ArrayWalk(arrVal, types.NewNull())
	if !result.ToBool() {
		t.Error("Expected array_walk to return true")
	}
}

// ============================================================================
// Set Operations Tests
// ============================================================================

func TestArrayDiff(t *testing.T) {
	arr1 := types.NewEmptyArray()
	arr1.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal1 := types.NewArray(arr1)

	arr2 := types.NewEmptyArray()
	arr2.Push(types.NewInt(2), types.NewInt(4))
	arrVal2 := types.NewArray(arr2)

	result := ArrayDiff(arrVal1, arrVal2)
	resultArray := result.ToArray()

	// Should contain 1 and 3 (values in arr1 but not arr2)
	if resultArray.Len() != 2 {
		t.Errorf("Expected length 2, got %d", resultArray.Len())
	}

	// Check that 1 and 3 are in the result
	found1 := false
	found3 := false
	resultArray.Each(func(_, value *types.Value) bool {
		if value.ToInt() == 1 {
			found1 = true
		}
		if value.ToInt() == 3 {
			found3 = true
		}
		return true
	})

	if !found1 || !found3 {
		t.Error("Expected to find values 1 and 3 in diff result")
	}
}

func TestArrayDiffMultiple(t *testing.T) {
	arr1 := types.NewEmptyArray()
	arr1.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4), types.NewInt(5))
	arrVal1 := types.NewArray(arr1)

	arr2 := types.NewEmptyArray()
	arr2.Push(types.NewInt(2), types.NewInt(3))
	arrVal2 := types.NewArray(arr2)

	arr3 := types.NewEmptyArray()
	arr3.Push(types.NewInt(4))
	arrVal3 := types.NewArray(arr3)

	result := ArrayDiff(arrVal1, arrVal2, arrVal3)
	resultArray := result.ToArray()

	// Should contain only 1 and 5
	if resultArray.Len() != 2 {
		t.Errorf("Expected length 2, got %d", resultArray.Len())
	}
}

func TestArrayIntersect(t *testing.T) {
	arr1 := types.NewEmptyArray()
	arr1.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal1 := types.NewArray(arr1)

	arr2 := types.NewEmptyArray()
	arr2.Push(types.NewInt(2), types.NewInt(3), types.NewInt(5))
	arrVal2 := types.NewArray(arr2)

	result := ArrayIntersect(arrVal1, arrVal2)
	resultArray := result.ToArray()

	// Should contain 2 and 3 (values in both arrays)
	if resultArray.Len() != 2 {
		t.Errorf("Expected length 2, got %d", resultArray.Len())
	}

	found2 := false
	found3 := false
	resultArray.Each(func(_, value *types.Value) bool {
		if value.ToInt() == 2 {
			found2 = true
		}
		if value.ToInt() == 3 {
			found3 = true
		}
		return true
	})

	if !found2 || !found3 {
		t.Error("Expected to find values 2 and 3 in intersect result")
	}
}

func TestArrayIntersectMultiple(t *testing.T) {
	arr1 := types.NewEmptyArray()
	arr1.Push(types.NewInt(1), types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal1 := types.NewArray(arr1)

	arr2 := types.NewEmptyArray()
	arr2.Push(types.NewInt(2), types.NewInt(3), types.NewInt(4))
	arrVal2 := types.NewArray(arr2)

	arr3 := types.NewEmptyArray()
	arr3.Push(types.NewInt(3), types.NewInt(4), types.NewInt(5))
	arrVal3 := types.NewArray(arr3)

	result := ArrayIntersect(arrVal1, arrVal2, arrVal3)
	resultArray := result.ToArray()

	// Should contain only 3 and 4 (values in all arrays)
	if resultArray.Len() != 2 {
		t.Errorf("Expected length 2, got %d", resultArray.Len())
	}
}

// ============================================================================
// Pointer Functions Tests
// ============================================================================

func TestCurrent(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	current := Current(arrVal)
	if current.ToInt() != 10 {
		t.Errorf("Expected current element 10, got %d", current.ToInt())
	}
}

func TestKey(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	key := Key(arrVal)
	if key.ToInt() != 0 {
		t.Errorf("Expected current key 0, got %d", key.ToInt())
	}
}

func TestReset(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	// Reset should return first element
	value := Reset(arrVal)
	if value.ToInt() != 10 {
		t.Errorf("Expected reset to return 10, got %d", value.ToInt())
	}

	// Current should also be first element
	current := Current(arrVal)
	if current.ToInt() != 10 {
		t.Errorf("Expected current to be 10 after reset, got %d", current.ToInt())
	}
}

func TestEnd(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	value := End(arrVal)
	if value.ToInt() != 30 {
		t.Errorf("Expected end to return 30, got %d", value.ToInt())
	}
}

func TestNext(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	value := Next(arrVal)
	if value.ToInt() != 20 {
		t.Errorf("Expected next to return 20, got %d", value.ToInt())
	}
}

func TestPrev(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20), types.NewInt(30))
	arrVal := types.NewArray(arr)

	// Prev without state always returns false
	value := Prev(arrVal)
	if value.Type() != types.TypeBool || value.ToBool() != false {
		t.Error("Expected prev to return false (simplified implementation)")
	}
}

func TestCurrentEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	arrVal := types.NewArray(arr)

	current := Current(arrVal)
	if current.Type() != types.TypeBool || current.ToBool() != false {
		t.Error("Expected false for current on empty array")
	}
}

func TestKeyEmpty(t *testing.T) {
	arr := types.NewEmptyArray()
	arrVal := types.NewArray(arr)

	key := Key(arrVal)
	if key.Type() != types.TypeNull {
		t.Error("Expected null for key on empty array")
	}
}

func TestNextBeyondEnd(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10))
	arrVal := types.NewArray(arr)

	// Try to go beyond the end
	next := Next(arrVal)
	if next.Type() != types.TypeBool || next.ToBool() != false {
		t.Error("Expected false when next goes beyond end")
	}
}

func TestPrevBeforeStart(t *testing.T) {
	arr := types.NewEmptyArray()
	arr.Push(types.NewInt(10), types.NewInt(20))
	arrVal := types.NewArray(arr)

	// Try to go before the start
	prev := Prev(arrVal)
	if prev.Type() != types.TypeBool || prev.ToBool() != false {
		t.Error("Expected false when prev goes before start")
	}
}

// ============================================================================
// Edge Cases for New Functions
// ============================================================================

func TestSortNil(t *testing.T) {
	result := Sort(nil)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for sort on nil")
	}
}

func TestArrayDiffEmpty(t *testing.T) {
	result := ArrayDiff()
	if result.Type() != types.TypeArray {
		t.Error("Expected empty array for diff with no arguments")
	}
}

func TestArrayIntersectEmpty(t *testing.T) {
	result := ArrayIntersect()
	if result.Type() != types.TypeArray {
		t.Error("Expected empty array for intersect with no arguments")
	}
}

func TestArrayMapEmpty(t *testing.T) {
	result := ArrayMap(types.NewNull())
	if result.Type() != types.TypeArray {
		t.Error("Expected empty array for map with no arrays")
	}
}
