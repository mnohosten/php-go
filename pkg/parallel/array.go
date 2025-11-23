package parallel

import (
	"sync"
)

// ============================================================================
// Automatic Array Parallelization - Parallel versions of array functions
// ============================================================================

// These provide parallel implementations of common PHP array functions that
// automatically parallelize when the array size is above a threshold.

// ArrayMapConfig configures parallel array_map behavior
type ArrayMapConfig struct {
	MinSize    int  // Minimum array size for parallelization
	NumWorkers int  // Number of parallel workers (0 = auto)
	Ordered    bool // Maintain order of results
}

// DefaultArrayMapConfig returns default configuration
func DefaultArrayMapConfig() ArrayMapConfig {
	return ArrayMapConfig{
		MinSize:    100,  // Parallelize arrays with 100+ elements
		NumWorkers: 0,    // Auto-detect (runtime.NumCPU())
		Ordered:    true, // Maintain order by default
	}
}

// ParallelArrayMap applies a function to each element in parallel
func ParallelArrayMap(arr []interface{}, fn func(interface{}) (interface{}, error), config ArrayMapConfig) ([]interface{}, error) {
	if len(arr) == 0 {
		return []interface{}{}, nil
	}

	// Use sequential if below threshold
	if len(arr) < config.MinSize {
		return sequentialArrayMap(arr, fn)
	}

	// Determine worker count
	numWorkers := config.NumWorkers
	if numWorkers <= 0 {
		numWorkers = 4 // Default
	}
	if numWorkers > len(arr) {
		numWorkers = len(arr)
	}

	// Create result array
	results := make([]interface{}, len(arr))
	errors := make([]error, len(arr))
	var firstError error
	var errorMu sync.Mutex

	// Parallel execution
	var wg sync.WaitGroup
	chunkSize := (len(arr) + numWorkers - 1) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}

		wg.Add(1)
		go func(startIdx, endIdx int) {
			defer wg.Done()

			for j := startIdx; j < endIdx; j++ {
				result, err := fn(arr[j])
				results[j] = result
				errors[j] = err

				if err != nil {
					errorMu.Lock()
					if firstError == nil {
						firstError = err
					}
					errorMu.Unlock()
				}
			}
		}(start, end)
	}

	wg.Wait()

	// Return first error if any
	if firstError != nil {
		return nil, firstError
	}

	return results, nil
}

// sequentialArrayMap performs sequential mapping
func sequentialArrayMap(arr []interface{}, fn func(interface{}) (interface{}, error)) ([]interface{}, error) {
	results := make([]interface{}, len(arr))

	for i, val := range arr {
		result, err := fn(val)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ============================================================================
// Parallel Array Filter
// ============================================================================

// ArrayFilterConfig configures parallel array_filter behavior
type ArrayFilterConfig struct {
	MinSize    int // Minimum array size for parallelization
	NumWorkers int // Number of parallel workers (0 = auto)
}

// DefaultArrayFilterConfig returns default configuration
func DefaultArrayFilterConfig() ArrayFilterConfig {
	return ArrayFilterConfig{
		MinSize:    100,
		NumWorkers: 0,
	}
}

// ParallelArrayFilter filters array elements in parallel
func ParallelArrayFilter(arr []interface{}, fn func(interface{}) (bool, error), config ArrayFilterConfig) ([]interface{}, error) {
	if len(arr) == 0 {
		return []interface{}{}, nil
	}

	// Use sequential if below threshold
	if len(arr) < config.MinSize {
		return sequentialArrayFilter(arr, fn)
	}

	// Determine worker count
	numWorkers := config.NumWorkers
	if numWorkers <= 0 {
		numWorkers = 4
	}
	if numWorkers > len(arr) {
		numWorkers = len(arr)
	}

	// Track which elements pass the filter
	passes := make([]bool, len(arr))
	errors := make([]error, len(arr))
	var firstError error
	var errorMu sync.Mutex

	// Parallel execution
	var wg sync.WaitGroup
	chunkSize := (len(arr) + numWorkers - 1) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}

		wg.Add(1)
		go func(startIdx, endIdx int) {
			defer wg.Done()

			for j := startIdx; j < endIdx; j++ {
				pass, err := fn(arr[j])
				passes[j] = pass
				errors[j] = err

				if err != nil {
					errorMu.Lock()
					if firstError == nil {
						firstError = err
					}
					errorMu.Unlock()
				}
			}
		}(start, end)
	}

	wg.Wait()

	// Return first error if any
	if firstError != nil {
		return nil, firstError
	}

	// Collect passing elements
	results := make([]interface{}, 0)
	for i, pass := range passes {
		if pass {
			results = append(results, arr[i])
		}
	}

	return results, nil
}

// sequentialArrayFilter performs sequential filtering
func sequentialArrayFilter(arr []interface{}, fn func(interface{}) (bool, error)) ([]interface{}, error) {
	results := make([]interface{}, 0)

	for _, val := range arr {
		pass, err := fn(val)
		if err != nil {
			return nil, err
		}
		if pass {
			results = append(results, val)
		}
	}

	return results, nil
}

// ============================================================================
// Parallel Array Reduce
// ============================================================================

// ArrayReduceConfig configures parallel array_reduce behavior
type ArrayReduceConfig struct {
	MinSize    int // Minimum array size for parallelization
	NumWorkers int // Number of parallel workers (0 = auto)
}

// DefaultArrayReduceConfig returns default configuration
func DefaultArrayReduceConfig() ArrayReduceConfig {
	return ArrayReduceConfig{
		MinSize:    1000, // Higher threshold for reduce (needs combining)
		NumWorkers: 0,
	}
}

// ParallelArrayReduce reduces array to single value in parallel
// Note: The reduce function must be associative for correct parallel execution
func ParallelArrayReduce(arr []interface{}, fn func(interface{}, interface{}) (interface{}, error), initial interface{}, config ArrayReduceConfig) (interface{}, error) {
	if len(arr) == 0 {
		return initial, nil
	}

	// Use sequential if below threshold
	if len(arr) < config.MinSize {
		return sequentialArrayReduce(arr, fn, initial)
	}

	// Determine worker count
	numWorkers := config.NumWorkers
	if numWorkers <= 0 {
		numWorkers = 4
	}
	if numWorkers > len(arr) {
		numWorkers = len(arr)
	}

	// Partial results from each worker
	partialResults := make([]interface{}, numWorkers)
	errors := make([]error, numWorkers)

	// Parallel execution - each worker reduces its chunk
	var wg sync.WaitGroup
	chunkSize := (len(arr) + numWorkers - 1) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}

		wg.Add(1)
		go func(workerIdx, startIdx, endIdx int) {
			defer wg.Done()

			// Reduce this chunk
			accumulator := initial
			for j := startIdx; j < endIdx; j++ {
				result, err := fn(accumulator, arr[j])
				if err != nil {
					errors[workerIdx] = err
					return
				}
				accumulator = result
			}

			partialResults[workerIdx] = accumulator
		}(i, start, end)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	// Combine partial results sequentially
	result := partialResults[0]
	for i := 1; i < numWorkers; i++ {
		combined, err := fn(result, partialResults[i])
		if err != nil {
			return nil, err
		}
		result = combined
	}

	return result, nil
}

// sequentialArrayReduce performs sequential reduction
func sequentialArrayReduce(arr []interface{}, fn func(interface{}, interface{}) (interface{}, error), initial interface{}) (interface{}, error) {
	accumulator := initial

	for _, val := range arr {
		result, err := fn(accumulator, val)
		if err != nil {
			return nil, err
		}
		accumulator = result
	}

	return accumulator, nil
}

// ============================================================================
// Parallel Array Walk
// ============================================================================

// ArrayWalkConfig configures parallel array_walk behavior
type ArrayWalkConfig struct {
	MinSize    int // Minimum array size for parallelization
	NumWorkers int // Number of parallel workers (0 = auto)
}

// DefaultArrayWalkConfig returns default configuration
func DefaultArrayWalkConfig() ArrayWalkConfig {
	return ArrayWalkConfig{
		MinSize:    100,
		NumWorkers: 0,
	}
}

// ParallelArrayWalk applies a function to each element (no return value)
func ParallelArrayWalk(arr []interface{}, fn func(interface{}, int) error, config ArrayWalkConfig) error {
	if len(arr) == 0 {
		return nil
	}

	// Use sequential if below threshold
	if len(arr) < config.MinSize {
		return sequentialArrayWalk(arr, fn)
	}

	// Determine worker count
	numWorkers := config.NumWorkers
	if numWorkers <= 0 {
		numWorkers = 4
	}
	if numWorkers > len(arr) {
		numWorkers = len(arr)
	}

	// Track errors
	errors := make([]error, numWorkers)

	// Parallel execution
	var wg sync.WaitGroup
	chunkSize := (len(arr) + numWorkers - 1) / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}

		wg.Add(1)
		go func(workerIdx, startIdx, endIdx int) {
			defer wg.Done()

			for j := startIdx; j < endIdx; j++ {
				err := fn(arr[j], j)
				if err != nil {
					errors[workerIdx] = err
					return
				}
			}
		}(i, start, end)
	}

	wg.Wait()

	// Return first error if any
	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

// sequentialArrayWalk performs sequential walk
func sequentialArrayWalk(arr []interface{}, fn func(interface{}, int) error) error {
	for i, val := range arr {
		if err := fn(val, i); err != nil {
			return err
		}
	}
	return nil
}

// ============================================================================
// Automatic Threshold Detection
// ============================================================================

// AutoThreshold automatically determines optimal parallelization threshold
type AutoThreshold struct {
	measurements []ThresholdMeasurement
	mu           sync.RWMutex
	maxSamples   int
}

// ThresholdMeasurement represents a performance measurement
type ThresholdMeasurement struct {
	ArraySize      int
	Sequential     int64 // nanoseconds
	Parallel       int64 // nanoseconds
	SpeedupFactor  float64
	Recommendation bool // true if parallel was faster
}

// NewAutoThreshold creates a new auto-threshold detector
func NewAutoThreshold() *AutoThreshold {
	return &AutoThreshold{
		measurements: make([]ThresholdMeasurement, 0),
		maxSamples:   100,
	}
}

// RecordMeasurement records a performance measurement
func (at *AutoThreshold) RecordMeasurement(measurement ThresholdMeasurement) {
	at.mu.Lock()
	defer at.mu.Unlock()

	if len(at.measurements) >= at.maxSamples {
		// Remove oldest
		at.measurements = at.measurements[1:]
	}

	at.measurements = append(at.measurements, measurement)
}

// GetRecommendedThreshold returns recommended threshold based on measurements
func (at *AutoThreshold) GetRecommendedThreshold() int {
	at.mu.RLock()
	defer at.mu.RUnlock()

	if len(at.measurements) == 0 {
		return 100 // Default
	}

	// Find the smallest array size where parallel was consistently faster
	threshold := 100
	for _, m := range at.measurements {
		if m.Recommendation && m.ArraySize < threshold {
			threshold = m.ArraySize
		}
	}

	return threshold
}

// GetMeasurements returns all measurements
func (at *AutoThreshold) GetMeasurements() []ThresholdMeasurement {
	at.mu.RLock()
	defer at.mu.RUnlock()

	result := make([]ThresholdMeasurement, len(at.measurements))
	copy(result, at.measurements)
	return result
}

// Clear clears all measurements
func (at *AutoThreshold) Clear() {
	at.mu.Lock()
	defer at.mu.Unlock()

	at.measurements = make([]ThresholdMeasurement, 0)
}

// ============================================================================
// Helper Functions
// ============================================================================

// ShouldParallelize determines if an array should be parallelized
func ShouldParallelize(arraySize int, minSize int) bool {
	return arraySize >= minSize
}

// OptimalWorkerCount determines optimal number of workers for array size
func OptimalWorkerCount(arraySize int) int {
	// Simple heuristic: use 4 workers for small-medium arrays
	// Use more for larger arrays
	if arraySize < 1000 {
		return 4
	} else if arraySize < 10000 {
		return 8
	} else {
		return 16
	}
}
