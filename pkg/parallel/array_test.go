package parallel

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// ============================================================================
// ParallelArrayMap Tests
// ============================================================================

func TestParallelArrayMapEmpty(t *testing.T) {
	arr := []interface{}{}
	config := DefaultArrayMapConfig()

	result, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Result length = %d, expected 0", len(result))
	}
}

func TestParallelArrayMapSequential(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayMapConfig()
	config.MinSize = 100 // Force sequential

	result, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 5 {
		t.Fatalf("Result length = %d, expected 5", len(result))
	}

	expected := []int{2, 4, 6, 8, 10}
	for i, val := range result {
		if val.(int) != expected[i] {
			t.Errorf("result[%d] = %d, expected %d", i, val, expected[i])
		}
	}
}

func TestParallelArrayMapParallel(t *testing.T) {
	arr := make([]interface{}, 200)
	for i := 0; i < 200; i++ {
		arr[i] = i
	}

	config := DefaultArrayMapConfig()
	config.MinSize = 100
	config.NumWorkers = 4

	result, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 200 {
		t.Fatalf("Result length = %d, expected 200", len(result))
	}

	// Verify all results
	for i, val := range result {
		expected := i * 2
		if val.(int) != expected {
			t.Errorf("result[%d] = %d, expected %d", i, val, expected)
		}
	}
}

func TestParallelArrayMapError(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayMapConfig()
	config.MinSize = 2

	_, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		if v.(int) == 3 {
			return nil, errors.New("error at 3")
		}
		return v, nil
	}, config)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestParallelArrayMapConcurrency(t *testing.T) {
	arr := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}

	config := DefaultArrayMapConfig()
	config.MinSize = 100
	config.NumWorkers = 8

	var counter sync.Map
	var maxConcurrent int32
	var currentConcurrent int32
	var mu sync.Mutex

	result, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		// Track concurrency
		mu.Lock()
		currentConcurrent++
		if currentConcurrent > maxConcurrent {
			maxConcurrent = currentConcurrent
		}
		mu.Unlock()

		time.Sleep(1 * time.Millisecond)
		counter.Store(v.(int), true)

		mu.Lock()
		currentConcurrent--
		mu.Unlock()

		return v.(int) * 2, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 1000 {
		t.Error("Result length mismatch")
	}

	// Should have some concurrency
	if maxConcurrent < 2 {
		t.Errorf("maxConcurrent = %d, expected >= 2", maxConcurrent)
	}
}

// ============================================================================
// ParallelArrayFilter Tests
// ============================================================================

func TestParallelArrayFilterEmpty(t *testing.T) {
	arr := []interface{}{}
	config := DefaultArrayFilterConfig()

	result, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		return true, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Result length = %d, expected 0", len(result))
	}
}

func TestParallelArrayFilterSequential(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	config := DefaultArrayFilterConfig()
	config.MinSize = 100 // Force sequential

	result, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		return v.(int)%2 == 0, nil // Even numbers only
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 5 {
		t.Fatalf("Result length = %d, expected 5", len(result))
	}

	expected := []int{2, 4, 6, 8, 10}
	for i, val := range result {
		if val.(int) != expected[i] {
			t.Errorf("result[%d] = %d, expected %d", i, val, expected[i])
		}
	}
}

func TestParallelArrayFilterParallel(t *testing.T) {
	arr := make([]interface{}, 200)
	for i := 0; i < 200; i++ {
		arr[i] = i
	}

	config := DefaultArrayFilterConfig()
	config.MinSize = 100
	config.NumWorkers = 4

	result, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		return v.(int)%2 == 0, nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 100 {
		t.Errorf("Result length = %d, expected 100", len(result))
	}

	// Verify all results are even
	for _, val := range result {
		if val.(int)%2 != 0 {
			t.Errorf("Found odd number: %d", val)
		}
	}
}

func TestParallelArrayFilterNone(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayFilterConfig()

	result, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		return false, nil // Filter out everything
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Result length = %d, expected 0", len(result))
	}
}

func TestParallelArrayFilterError(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayFilterConfig()
	config.MinSize = 2

	_, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		if v.(int) == 3 {
			return false, errors.New("error at 3")
		}
		return true, nil
	}, config)

	if err == nil {
		t.Error("Expected error")
	}
}

// ============================================================================
// ParallelArrayReduce Tests
// ============================================================================

func TestParallelArrayReduceEmpty(t *testing.T) {
	arr := []interface{}{}
	config := DefaultArrayReduceConfig()

	result, err := ParallelArrayReduce(arr, func(acc, v interface{}) (interface{}, error) {
		return acc.(int) + v.(int), nil
	}, 0, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.(int) != 0 {
		t.Errorf("Result = %d, expected 0", result)
	}
}

func TestParallelArrayReduceSequential(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayReduceConfig()
	config.MinSize = 1000 // Force sequential

	result, err := ParallelArrayReduce(arr, func(acc, v interface{}) (interface{}, error) {
		return acc.(int) + v.(int), nil
	}, 0, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.(int) != 15 {
		t.Errorf("Result = %d, expected 15", result)
	}
}

func TestParallelArrayReduceParallel(t *testing.T) {
	arr := make([]interface{}, 2000)
	for i := 0; i < 2000; i++ {
		arr[i] = 1
	}

	config := DefaultArrayReduceConfig()
	config.MinSize = 1000
	config.NumWorkers = 4

	result, err := ParallelArrayReduce(arr, func(acc, v interface{}) (interface{}, error) {
		return acc.(int) + v.(int), nil
	}, 0, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.(int) != 2000 {
		t.Errorf("Result = %d, expected 2000", result)
	}
}

func TestParallelArrayReduceProduct(t *testing.T) {
	arr := []interface{}{2, 3, 4, 5}
	config := DefaultArrayReduceConfig()

	result, err := ParallelArrayReduce(arr, func(acc, v interface{}) (interface{}, error) {
		return acc.(int) * v.(int), nil
	}, 1, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.(int) != 120 {
		t.Errorf("Result = %d, expected 120 (2*3*4*5)", result)
	}
}

func TestParallelArrayReduceError(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayReduceConfig()
	config.MinSize = 2

	_, err := ParallelArrayReduce(arr, func(acc, v interface{}) (interface{}, error) {
		if v.(int) == 3 {
			return nil, errors.New("error at 3")
		}
		return acc.(int) + v.(int), nil
	}, 0, config)

	if err == nil {
		t.Error("Expected error")
	}
}

// ============================================================================
// ParallelArrayWalk Tests
// ============================================================================

func TestParallelArrayWalkEmpty(t *testing.T) {
	arr := []interface{}{}
	config := DefaultArrayWalkConfig()

	err := ParallelArrayWalk(arr, func(v interface{}, i int) error {
		return nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParallelArrayWalkSequential(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayWalkConfig()
	config.MinSize = 100 // Force sequential

	var sum int
	var mu sync.Mutex

	err := ParallelArrayWalk(arr, func(v interface{}, i int) error {
		mu.Lock()
		sum += v.(int)
		mu.Unlock()
		return nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sum != 15 {
		t.Errorf("Sum = %d, expected 15", sum)
	}
}

func TestParallelArrayWalkParallel(t *testing.T) {
	arr := make([]interface{}, 200)
	for i := 0; i < 200; i++ {
		arr[i] = 1
	}

	config := DefaultArrayWalkConfig()
	config.MinSize = 100
	config.NumWorkers = 4

	var sum int
	var mu sync.Mutex

	err := ParallelArrayWalk(arr, func(v interface{}, i int) error {
		mu.Lock()
		sum += v.(int)
		mu.Unlock()
		return nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sum != 200 {
		t.Errorf("Sum = %d, expected 200", sum)
	}
}

func TestParallelArrayWalkError(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	config := DefaultArrayWalkConfig()
	config.MinSize = 2

	err := ParallelArrayWalk(arr, func(v interface{}, i int) error {
		if v.(int) == 3 {
			return errors.New("error at 3")
		}
		return nil
	}, config)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestParallelArrayWalkIndex(t *testing.T) {
	arr := []interface{}{10, 20, 30, 40, 50}
	config := DefaultArrayWalkConfig()

	indices := make(map[int]bool)
	var mu sync.Mutex

	err := ParallelArrayWalk(arr, func(v interface{}, i int) error {
		mu.Lock()
		indices[i] = true
		mu.Unlock()
		return nil
	}, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should have all indices
	for i := 0; i < 5; i++ {
		if !indices[i] {
			t.Errorf("Missing index %d", i)
		}
	}
}

// ============================================================================
// AutoThreshold Tests
// ============================================================================

func TestNewAutoThreshold(t *testing.T) {
	at := NewAutoThreshold()
	if at == nil {
		t.Error("NewAutoThreshold returned nil")
	}

	threshold := at.GetRecommendedThreshold()
	if threshold != 100 {
		t.Errorf("Default threshold = %d, expected 100", threshold)
	}
}

func TestAutoThresholdRecordMeasurement(t *testing.T) {
	at := NewAutoThreshold()

	measurement := ThresholdMeasurement{
		ArraySize:      200,
		Sequential:     1000000,
		Parallel:       500000,
		SpeedupFactor:  2.0,
		Recommendation: true,
	}

	at.RecordMeasurement(measurement)

	measurements := at.GetMeasurements()
	if len(measurements) != 1 {
		t.Errorf("Measurement count = %d, expected 1", len(measurements))
	}
}

func TestAutoThresholdRecommendation(t *testing.T) {
	at := NewAutoThreshold()

	// Record that parallel is faster at 150
	at.RecordMeasurement(ThresholdMeasurement{
		ArraySize:      150,
		Sequential:     1000000,
		Parallel:       500000,
		Recommendation: true,
	})

	// Record that parallel is faster at 200
	at.RecordMeasurement(ThresholdMeasurement{
		ArraySize:      200,
		Sequential:     2000000,
		Parallel:       800000,
		Recommendation: true,
	})

	threshold := at.GetRecommendedThreshold()
	// Should recommend the smaller size where parallel was faster
	if threshold > 150 {
		t.Errorf("Threshold = %d, should be <= 150", threshold)
	}
}

func TestAutoThresholdMaxSamples(t *testing.T) {
	at := NewAutoThreshold()
	at.maxSamples = 10

	// Record more than max
	for i := 0; i < 15; i++ {
		at.RecordMeasurement(ThresholdMeasurement{
			ArraySize:      100 + i,
			Recommendation: true,
		})
	}

	measurements := at.GetMeasurements()
	if len(measurements) != 10 {
		t.Errorf("Measurement count = %d, expected 10", len(measurements))
	}

	// Should have kept the last 10
	if measurements[0].ArraySize != 105 {
		t.Errorf("First measurement array size = %d, expected 105", measurements[0].ArraySize)
	}
}

func TestAutoThresholdClear(t *testing.T) {
	at := NewAutoThreshold()

	at.RecordMeasurement(ThresholdMeasurement{
		ArraySize:      200,
		Recommendation: true,
	})

	at.Clear()

	measurements := at.GetMeasurements()
	if len(measurements) != 0 {
		t.Errorf("Measurement count = %d, expected 0 after clear", len(measurements))
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestShouldParallelize(t *testing.T) {
	tests := []struct {
		arraySize int
		minSize   int
		expected  bool
	}{
		{50, 100, false},
		{100, 100, true},
		{150, 100, true},
		{0, 100, false},
	}

	for _, tt := range tests {
		result := ShouldParallelize(tt.arraySize, tt.minSize)
		if result != tt.expected {
			t.Errorf("ShouldParallelize(%d, %d) = %v, expected %v",
				tt.arraySize, tt.minSize, result, tt.expected)
		}
	}
}

func TestOptimalWorkerCount(t *testing.T) {
	tests := []struct {
		arraySize int
		expected  int
	}{
		{100, 4},
		{999, 4},
		{1000, 8},
		{9999, 8},
		{10000, 16},
		{100000, 16},
	}

	for _, tt := range tests {
		result := OptimalWorkerCount(tt.arraySize)
		if result != tt.expected {
			t.Errorf("OptimalWorkerCount(%d) = %d, expected %d",
				tt.arraySize, result, tt.expected)
		}
	}
}

// ============================================================================
// Default Config Tests
// ============================================================================

func TestDefaultArrayMapConfig(t *testing.T) {
	config := DefaultArrayMapConfig()
	if config.MinSize != 100 {
		t.Errorf("MinSize = %d, expected 100", config.MinSize)
	}
	if config.NumWorkers != 0 {
		t.Errorf("NumWorkers = %d, expected 0", config.NumWorkers)
	}
	if !config.Ordered {
		t.Error("Ordered should be true")
	}
}

func TestDefaultArrayFilterConfig(t *testing.T) {
	config := DefaultArrayFilterConfig()
	if config.MinSize != 100 {
		t.Errorf("MinSize = %d, expected 100", config.MinSize)
	}
	if config.NumWorkers != 0 {
		t.Errorf("NumWorkers = %d, expected 0", config.NumWorkers)
	}
}

func TestDefaultArrayReduceConfig(t *testing.T) {
	config := DefaultArrayReduceConfig()
	if config.MinSize != 1000 {
		t.Errorf("MinSize = %d, expected 1000", config.MinSize)
	}
	if config.NumWorkers != 0 {
		t.Errorf("NumWorkers = %d, expected 0", config.NumWorkers)
	}
}

func TestDefaultArrayWalkConfig(t *testing.T) {
	config := DefaultArrayWalkConfig()
	if config.MinSize != 100 {
		t.Errorf("MinSize = %d, expected 100", config.MinSize)
	}
	if config.NumWorkers != 0 {
		t.Errorf("NumWorkers = %d, expected 0", config.NumWorkers)
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestParallelArrayIntegration(t *testing.T) {
	// Create array
	arr := make([]interface{}, 500)
	for i := 0; i < 500; i++ {
		arr[i] = i
	}

	// Map: multiply by 2
	config1 := DefaultArrayMapConfig()
	config1.MinSize = 100
	mapped, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, config1)
	if err != nil {
		t.Fatalf("Map error: %v", err)
	}

	// Filter: keep even numbers
	config2 := DefaultArrayFilterConfig()
	config2.MinSize = 100
	filtered, err := ParallelArrayFilter(mapped, func(v interface{}) (bool, error) {
		return v.(int)%4 == 0, nil // Divisible by 4
	}, config2)
	if err != nil {
		t.Fatalf("Filter error: %v", err)
	}

	// Should have 250 elements (0, 2, 4, ... 498 -> 0, 4, 8, ... 996)
	if len(filtered) != 250 {
		t.Errorf("Filtered length = %d, expected 250", len(filtered))
	}

	// Reduce: sum
	config3 := DefaultArrayReduceConfig()
	config3.MinSize = 100
	sum, err := ParallelArrayReduce(filtered, func(acc, v interface{}) (interface{}, error) {
		return acc.(int) + v.(int), nil
	}, 0, config3)
	if err != nil {
		t.Fatalf("Reduce error: %v", err)
	}

	// Sum should be 0 + 4 + 8 + ... + 996
	// = 4 * (0 + 1 + 2 + ... + 249)
	// = 4 * (249 * 250 / 2)
	// = 4 * 31125
	// = 124500
	expected := 124500
	if sum.(int) != expected {
		t.Errorf("Sum = %d, expected %d", sum, expected)
	}
}
