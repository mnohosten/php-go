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
