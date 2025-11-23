package parallel

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// Phase 7 Integration Tests - Test all parallel components together
// ============================================================================

// TestIntegrationPoolWithContext tests worker pool with request context
func TestIntegrationPoolWithContext(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	ctx := NewRequestContext("test")

	var counter atomic.Int64
	var wg sync.WaitGroup

	// Submit 100 tasks
	for i := 0; i < 100; i++ {
		wg.Add(1)
		taskNum := i
		task := NewTask(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			counter.Add(1)
			ctx.SetGlobal(fmt.Sprintf("task-%d", taskNum), true)
			wg.Done()
			return nil, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	if counter.Load() != 100 {
		t.Errorf("Expected 100 tasks executed, got %d", counter.Load())
	}

	if ctx.GlobalCount() != 100 {
		t.Errorf("Expected 100 globals in context, got %d", ctx.GlobalCount())
	}
}

// TestIntegrationPoolWithCOW tests worker pool with COW optimization
func TestIntegrationPoolWithCOW(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	// Create shared COW array
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	// Multiple workers read from shared array
	var wg sync.WaitGroup
	results := make([]int, 4)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		workerIdx := i
		task := NewTask(fmt.Sprintf("worker-%d", i), func() (interface{}, error) {
			clone := cowArr.Clone()
			// Read-only access should not trigger copy
			sum := 0
			for j := 0; j < clone.Len(); j++ {
				sum += clone.Get(j).(int)
			}
			results[workerIdx] = sum
			wg.Done()
			return sum, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	// All workers should get same sum
	expected := 4950 // sum of 0..99
	for i, result := range results {
		if result != expected {
			t.Errorf("Worker %d got sum %d, expected %d", i, result, expected)
		}
	}

	// Array should still be shared (only reads occurred)
	if !cowArr.IsShared() {
		t.Error("COWArray should still be shared after read-only access")
	}
}

// TestIntegrationParallelArrayWithMetrics tests array operations with metrics
func TestIntegrationParallelArrayWithMetrics(t *testing.T) {
	// Create large array
	arr := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}

	// Configure parallel map
	config := DefaultArrayMapConfig()
	config.MinSize = 500
	config.NumWorkers = 4

	// Map with metrics tracking
	start := time.Now()
	result, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, config)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("ParallelArrayMap failed: %v", err)
	}

	if len(result) != 1000 {
		t.Errorf("Expected 1000 results, got %d", len(result))
	}

	// Verify results
	for i, val := range result {
		expected := i * 2
		if val.(int) != expected {
			t.Errorf("result[%d] = %d, expected %d", i, val, expected)
		}
	}

	t.Logf("Parallel map of 1000 elements took %v", duration)
}

// TestIntegrationContextWithSync tests context with sync primitives
func TestIntegrationContextWithSync(t *testing.T) {
	barrier := NewBarrier(4)

	var wg sync.WaitGroup
	var counter atomic.Int64

	// Launch 4 workers that synchronize at barrier
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Phase 1: increment counter
			counter.Add(1)

			// Wait at barrier
			barrier.Wait()

			// Phase 2: all counters should be 4
			if counter.Load() != 4 {
				t.Errorf("Worker %d: expected counter=4, got %d", id, counter.Load())
			}
		}(i)
	}

	wg.Wait()
}

// TestIntegrationFullPipeline tests complete pipeline with all components
func TestIntegrationFullPipeline(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	cowMgr := NewCOWManager()

	// Create input data with COW
	inputData := make([]interface{}, 200)
	for i := 0; i < 200; i++ {
		inputData[i] = i
	}
	cowInput := NewCOWArray(inputData)

	// Track completion
	var wg sync.WaitGroup
	results := make([]int, 4)

	// Stage 1: Multiple workers clone and filter
	for i := 0; i < 4; i++ {
		wg.Add(1)
		workerIdx := i
		task := NewTask(fmt.Sprintf("worker-%d", i), func() (interface{}, error) {
			// Clone data (COW)
			clone := cowInput.Clone()
			cowMgr.RecordShare(int64(EstimateArraySize(inputData)))

			// Filter even numbers
			filtered := make([]interface{}, 0)
			for j := 0; j < clone.Len(); j++ {
				val := clone.Get(j).(int)
				if val%2 == 0 {
					filtered = append(filtered, val)
				}
			}

			// Reduce to sum
			sum := 0
			for _, val := range filtered {
				sum += val.(int)
			}

			results[workerIdx] = sum
			wg.Done()
			return sum, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	// Verify all workers got same result
	expected := 9900 // sum of even numbers 0..198
	for i, result := range results {
		if result != expected {
			t.Errorf("Worker %d got sum %d, expected %d", i, result, expected)
		}
	}

	// Check COW stats
	cowStats := cowMgr.GetStats()
	if cowStats.TotalShares != 4 {
		t.Errorf("Expected 4 COW shares, got %d", cowStats.TotalShares)
	}
}

// TestIntegrationErrorHandling tests error propagation across components
func TestIntegrationErrorHandling(t *testing.T) {
	// Create array with some elements that will cause errors
	arr := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		arr[i] = i
	}

	config := DefaultArrayMapConfig()
	config.MinSize = 50
	config.NumWorkers = 4

	// Map function that errors on specific values
	_, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		val := v.(int)
		if val == 50 {
			return nil, errors.New("error at value 50")
		}
		return val * 2, nil
	}, config)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "error at value 50" {
		t.Errorf("Expected 'error at value 50', got '%s'", err.Error())
	}
}

// TestIntegrationConcurrentCOWWrites tests COW with concurrent writes
func TestIntegrationConcurrentCOWWrites(t *testing.T) {
	// Create shared COW array
	data := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		data[i] = 0
	}
	cowArr := NewCOWArray(data)

	// Create 4 clones
	clones := make([]*COWArray, 4)
	for i := 0; i < 4; i++ {
		clones[i] = cowArr.Clone()
	}

	// Each clone writes to its own copy
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		workerIdx := i
		go func() {
			defer wg.Done()

			clone := clones[workerIdx]
			// Write to all elements
			for j := 0; j < clone.Len(); j++ {
				clone.Set(j, workerIdx)
			}
		}()
	}

	wg.Wait()

	// Each clone should have its own values
	for i := 0; i < 4; i++ {
		for j := 0; j < clones[i].Len(); j++ {
			val := clones[i].Get(j).(int)
			if val != i {
				t.Errorf("Clone %d index %d = %d, expected %d", i, j, val, i)
			}
		}
	}

	// Original should still be 0s
	for j := 0; j < cowArr.Len(); j++ {
		val := cowArr.Get(j).(int)
		if val != 0 {
			t.Errorf("Original index %d = %d, expected 0", j, val)
		}
	}
}

// TestIntegrationParallelFilterReduce tests filter then reduce pipeline
func TestIntegrationParallelFilterReduce(t *testing.T) {
	// Create large array
	arr := make([]interface{}, 2000)
	for i := 0; i < 2000; i++ {
		arr[i] = i
	}

	filterConfig := DefaultArrayFilterConfig()
	filterConfig.MinSize = 500
	filterConfig.NumWorkers = 4

	// Filter even numbers
	filtered, err := ParallelArrayFilter(arr, func(v interface{}) (bool, error) {
		return v.(int)%2 == 0, nil
	}, filterConfig)

	if err != nil {
		t.Fatalf("ParallelArrayFilter failed: %v", err)
	}

	if len(filtered) != 1000 {
		t.Errorf("Expected 1000 filtered elements, got %d", len(filtered))
	}

	reduceConfig := DefaultArrayReduceConfig()
	reduceConfig.MinSize = 500
	reduceConfig.NumWorkers = 4

	// Reduce to sum
	sum, err := ParallelArrayReduce(filtered, func(acc, val interface{}) (interface{}, error) {
		return acc.(int) + val.(int), nil
	}, 0, reduceConfig)

	if err != nil {
		t.Fatalf("ParallelArrayReduce failed: %v", err)
	}

	// Sum of even numbers 0..1998 = 999000
	expected := 999000
	if sum.(int) != expected {
		t.Errorf("Sum = %d, expected %d", sum, expected)
	}
}

// TestIntegrationStressTest stress tests all components
func TestIntegrationStressTest(t *testing.T) {
	t.Skip("Skipping stress test - takes too long")

	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	pool := NewWorkerPool(8)
	defer pool.Shutdown()

	cowMgr := NewCOWManager()

	// Create large dataset
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	var wg sync.WaitGroup
	iterations := 100

	// Run 100 iterations of parallel operations
	for iter := 0; iter < iterations; iter++ {
		wg.Add(1)
		iterNum := iter
		task := NewTask(fmt.Sprintf("iter-%d", iter), func() (interface{}, error) {
			// Create COW array
			cowArr := NewCOWArray(data)
			cowMgr.RecordShare(int64(EstimateArraySize(data)))

			// Perform parallel map
			config := DefaultArrayMapConfig()
			config.MinSize = 1000
			config.NumWorkers = 4

			result, err := ParallelArrayMap(cowArr.ToSlice(), func(v interface{}) (interface{}, error) {
				return v.(int) * 2, nil
			}, config)

			if err != nil {
				t.Errorf("Iteration %d: ParallelArrayMap failed: %v", iterNum, err)
			}

			if len(result) != 10000 {
				t.Errorf("Iteration %d: Expected 10000 results, got %d", iterNum, len(result))
			}

			wg.Done()
			return result, err
		})
		pool.Submit(task)
	}

	wg.Wait()

	// Check stats
	cowStats := cowMgr.GetStats()
	t.Logf("Stress test: %d COW shares, %d bytes saved",
		cowStats.TotalShares, cowStats.BytesSaved)
}

// TestIntegrationMapFilterReduce tests complete map-filter-reduce pipeline
func TestIntegrationMapFilterReduce(t *testing.T) {
	// Input: 0..999
	arr := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}

	// Map: x -> x * 2
	mapConfig := DefaultArrayMapConfig()
	mapConfig.MinSize = 500
	mapped, err := ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, mapConfig)

	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Filter: keep only numbers > 1000
	filterConfig := DefaultArrayFilterConfig()
	filterConfig.MinSize = 500
	filtered, err := ParallelArrayFilter(mapped, func(v interface{}) (bool, error) {
		return v.(int) > 1000, nil
	}, filterConfig)

	if err != nil {
		t.Fatalf("Filter failed: %v", err)
	}

	// Should have 499 elements (1002, 1004, ..., 1998)
	if len(filtered) != 499 {
		t.Errorf("Expected 499 filtered elements, got %d", len(filtered))
	}

	// Reduce: sum
	reduceConfig := DefaultArrayReduceConfig()
	reduceConfig.MinSize = 100
	sum, err := ParallelArrayReduce(filtered, func(acc, val interface{}) (interface{}, error) {
		return acc.(int) + val.(int), nil
	}, 0, reduceConfig)

	if err != nil {
		t.Fatalf("Reduce failed: %v", err)
	}

	// Sum = 1002 + 1004 + ... + 1998
	// = 499 * (1002 + 1998) / 2 = 499 * 1500 = 748500
	expected := 748500
	if sum.(int) != expected {
		t.Errorf("Sum = %d, expected %d", sum, expected)
	}
}

// TestIntegrationCOWMapWithPool tests COW map with worker pool
func TestIntegrationCOWMapWithPool(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	// Create shared COW map
	data := map[string]interface{}{
		"name":  "test",
		"count": 0,
		"items": []int{1, 2, 3},
	}
	cowMap := NewCOWMapFromMap(data)

	var wg sync.WaitGroup
	results := make([]string, 4)

	// Each worker clones and reads
	for i := 0; i < 4; i++ {
		wg.Add(1)
		workerIdx := i
		task := NewTask(fmt.Sprintf("worker-%d", i), func() (interface{}, error) {
			clone := cowMap.Clone()
			name, _ := clone.Get("name")
			results[workerIdx] = name.(string)
			wg.Done()
			return name, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	// All should read same value
	for i, name := range results {
		if name != "test" {
			t.Errorf("Worker %d got name '%s', expected 'test'", i, name)
		}
	}
}

// TestIntegrationBarrierComplete tests barrier with all workers arriving
func TestIntegrationBarrierComplete(t *testing.T) {
	barrier := NewBarrier(3)

	var wg sync.WaitGroup
	var counter atomic.Int32

	// All three workers arrive
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
			err := barrier.Wait()
			if err != nil {
				t.Errorf("Barrier wait failed: %v", err)
			}
			// All should see counter=3 after barrier
			if counter.Load() != 3 {
				t.Errorf("Expected counter=3, got %d", counter.Load())
			}
		}()
	}

	wg.Wait()
}

// TestIntegrationMetricsCollection tests metrics collection across operations
func TestIntegrationMetricsCollection(t *testing.T) {
	cowMgr := NewCOWManager()

	// Simulate series of operations
	for i := 0; i < 10; i++ {
		// Simulate COW share
		cowMgr.RecordShare(1024)

		// Simulate some copies
		if i%3 == 0 {
			cowMgr.RecordCopy(512)
		}
	}

	// Check COW stats
	cowStats := cowMgr.GetStats()
	if cowStats.TotalShares != 10 {
		t.Errorf("Expected 10 shares, got %d", cowStats.TotalShares)
	}

	expectedCopies := int64(4) // i=0,3,6,9
	if cowStats.TotalCopies != expectedCopies {
		t.Errorf("Expected %d copies, got %d", expectedCopies, cowStats.TotalCopies)
	}

	// Check efficiency
	efficiency := cowStats.Efficiency()
	if efficiency <= 0 || efficiency >= 1 {
		t.Errorf("Efficiency should be between 0 and 1, got %f", efficiency)
	}

	t.Logf("Efficiency: %.2f%% (%d bytes saved, %d bytes allocated)",
		efficiency*100, cowStats.BytesSaved, cowStats.BytesAlloced)
}

// TestIntegrationRecovery tests error recovery in parallel operations
func TestIntegrationRecovery(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	var wg sync.WaitGroup
	var panicCount atomic.Int32

	// Submit tasks that panic
	for i := 0; i < 10; i++ {
		wg.Add(1)
		taskNum := i
		task := NewTask(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					panicCount.Add(1)
				}
				wg.Done()
			}()

			if taskNum%3 == 0 {
				panic(fmt.Sprintf("panic in task %d", taskNum))
			}
			return nil, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	// Should have 4 panics (0, 3, 6, 9)
	if panicCount.Load() != 4 {
		t.Errorf("Expected 4 panics, got %d", panicCount.Load())
	}
}
