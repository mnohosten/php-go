package parallel

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Explicit Parallelism APIs - New PHP functions for parallel execution
// ============================================================================

// These are new PHP functions that don't exist in standard PHP
// They provide explicit control over parallelism in PHP-Go

// ============================================================================
// Channel - Go-style channels for PHP
// ============================================================================

// Channel represents a typed communication channel
type Channel struct {
	ch       chan interface{}
	capacity int
	closed   bool
	mu       sync.RWMutex
}

// NewChannel creates a new channel with the specified capacity
// Capacity 0 = unbuffered, >0 = buffered
func NewChannel(capacity int) *Channel {
	if capacity < 0 {
		capacity = 0
	}

	return &Channel{
		ch:       make(chan interface{}, capacity),
		capacity: capacity,
		closed:   false,
	}
}

// Send sends a value to the channel
func (c *Channel) Send(value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cannot send on closed channel")
	}

	c.ch <- value
	return nil
}

// Receive receives a value from the channel
// Returns (value, true) if successful, (nil, false) if channel is closed
func (c *Channel) Receive() (interface{}, bool) {
	value, ok := <-c.ch
	return value, ok
}

// TryReceive attempts to receive without blocking
// Returns (value, true) if successful, (nil, false) if no value available
func (c *Channel) TryReceive() (interface{}, bool) {
	select {
	case value, ok := <-c.ch:
		return value, ok
	default:
		return nil, false
	}
}

// ReceiveWithTimeout receives with a timeout
func (c *Channel) ReceiveWithTimeout(timeout time.Duration) (interface{}, bool) {
	select {
	case value, ok := <-c.ch:
		return value, ok
	case <-time.After(timeout):
		return nil, false
	}
}

// Close closes the channel
func (c *Channel) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		close(c.ch)
		c.closed = true
	}
}

// IsClosed returns true if the channel is closed
func (c *Channel) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Capacity returns the channel's capacity
func (c *Channel) Capacity() int {
	return c.capacity
}

// ============================================================================
// Goroutine - Execute function in parallel
// ============================================================================

// Goroutine represents a running goroutine with result tracking
type Goroutine struct {
	id     string
	done   chan struct{}
	result interface{}
	err    error
	mu     sync.RWMutex
}

// NewGoroutine creates and starts a goroutine
func NewGoroutine(id string, fn func() (interface{}, error)) *Goroutine {
	gr := &Goroutine{
		id:   id,
		done: make(chan struct{}),
	}

	go func() {
		defer close(gr.done)

		// Handle panics
		defer func() {
			if r := recover(); r != nil {
				gr.mu.Lock()
				gr.err = fmt.Errorf("panic in goroutine %s: %v", id, r)
				gr.mu.Unlock()
			}
		}()

		result, err := fn()

		gr.mu.Lock()
		gr.result = result
		gr.err = err
		gr.mu.Unlock()
	}()

	return gr
}

// Wait blocks until the goroutine completes
func (gr *Goroutine) Wait() (interface{}, error) {
	<-gr.done

	gr.mu.RLock()
	defer gr.mu.RUnlock()

	return gr.result, gr.err
}

// IsDone returns true if the goroutine has completed
func (gr *Goroutine) IsDone() bool {
	select {
	case <-gr.done:
		return true
	default:
		return false
	}
}

// ID returns the goroutine's ID
func (gr *Goroutine) ID() string {
	return gr.id
}

// ============================================================================
// Parallel Execution - Run multiple functions in parallel
// ============================================================================

// ParallelResult represents the result of a parallel execution
type ParallelResult struct {
	Index  int
	Result interface{}
	Error  error
}

// Parallel executes multiple functions in parallel and returns all results
func Parallel(functions []func() (interface{}, error)) []ParallelResult {
	if len(functions) == 0 {
		return []ParallelResult{}
	}

	results := make([]ParallelResult, len(functions))
	var wg sync.WaitGroup

	for i, fn := range functions {
		wg.Add(1)
		go func(index int, function func() (interface{}, error)) {
			defer wg.Done()

			// Handle panics
			defer func() {
				if r := recover(); r != nil {
					results[index] = ParallelResult{
						Index: index,
						Error: fmt.Errorf("panic at index %d: %v", index, r),
					}
				}
			}()

			result, err := function()
			results[index] = ParallelResult{
				Index:  index,
				Result: result,
				Error:  err,
			}
		}(i, fn)
	}

	wg.Wait()
	return results
}

// ParallelWithLimit executes functions in parallel with a concurrency limit
func ParallelWithLimit(functions []func() (interface{}, error), limit int) []ParallelResult {
	if len(functions) == 0 {
		return []ParallelResult{}
	}

	if limit <= 0 {
		limit = len(functions)
	}

	results := make([]ParallelResult, len(functions))
	semaphore := make(chan struct{}, limit)
	var wg sync.WaitGroup

	for i, fn := range functions {
		wg.Add(1)
		go func(index int, function func() (interface{}, error)) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Handle panics
			defer func() {
				if r := recover(); r != nil {
					results[index] = ParallelResult{
						Index: index,
						Error: fmt.Errorf("panic at index %d: %v", index, r),
					}
				}
			}()

			result, err := function()
			results[index] = ParallelResult{
				Index:  index,
				Result: result,
				Error:  err,
			}
		}(i, fn)
	}

	wg.Wait()
	return results
}

// ============================================================================
// WaitGroup - Synchronization primitive
// ============================================================================

// WaitGroup wraps sync.WaitGroup for PHP
type WaitGroup struct {
	wg sync.WaitGroup
}

// NewWaitGroup creates a new WaitGroup
func NewWaitGroup() *WaitGroup {
	return &WaitGroup{}
}

// Add adds delta to the WaitGroup counter
func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one
func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero
func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}

// ============================================================================
// ParallelExecutor - High-level parallel execution with worker pool
// ============================================================================

// ParallelExecutor provides high-level parallel execution capabilities
type ParallelExecutor struct {
	pool *WorkerPool
}

// NewParallelExecutor creates a new parallel executor with a worker pool
func NewParallelExecutor(numWorkers int) *ParallelExecutor {
	pool := NewWorkerPool(numWorkers)
	pool.Start()

	return &ParallelExecutor{
		pool: pool,
	}
}

// Execute runs a function and returns a future
func (pe *ParallelExecutor) Execute(id string, fn func() (interface{}, error)) *Future {
	return pe.pool.SubmitFunc(id, fn)
}

// ExecuteMany runs multiple functions and returns futures
func (pe *ParallelExecutor) ExecuteMany(functions map[string]func() (interface{}, error)) map[string]*Future {
	futures := make(map[string]*Future)

	for id, fn := range functions {
		futures[id] = pe.pool.SubmitFunc(id, fn)
	}

	return futures
}

// Wait waits for all pending tasks to complete
func (pe *ParallelExecutor) Wait() {
	pe.pool.Wait()
}

// Shutdown gracefully shuts down the executor
func (pe *ParallelExecutor) Shutdown() {
	pe.pool.Shutdown()
}

// ShutdownWithTimeout shuts down with a timeout
func (pe *ParallelExecutor) ShutdownWithTimeout(timeout time.Duration) error {
	return pe.pool.ShutdownWithTimeout(timeout)
}

// Stats returns executor statistics
func (pe *ParallelExecutor) Stats() PoolStats {
	return pe.pool.Stats()
}

// ============================================================================
// Pipeline - Execute functions in sequence with parallel stages
// ============================================================================

// Pipeline represents a sequence of processing stages
type Pipeline struct {
	stages []PipelineStage
}

// PipelineStage represents a single stage in the pipeline
type PipelineStage struct {
	Name     string
	Function func(interface{}) (interface{}, error)
	Parallel bool
}

// NewPipeline creates a new pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		stages: make([]PipelineStage, 0),
	}
}

// AddStage adds a stage to the pipeline
func (p *Pipeline) AddStage(name string, fn func(interface{}) (interface{}, error), parallel bool) {
	p.stages = append(p.stages, PipelineStage{
		Name:     name,
		Function: fn,
		Parallel: parallel,
	})
}

// Execute runs the pipeline on input data
func (p *Pipeline) Execute(input interface{}) (interface{}, error) {
	current := input

	for _, stage := range p.stages {
		result, err := stage.Function(current)
		if err != nil {
			return nil, fmt.Errorf("error in stage %s: %w", stage.Name, err)
		}
		current = result
	}

	return current, nil
}

// ExecuteMany runs the pipeline on multiple inputs in parallel
func (p *Pipeline) ExecuteMany(inputs []interface{}) []ParallelResult {
	functions := make([]func() (interface{}, error), len(inputs))

	for i, input := range inputs {
		capturedInput := input
		functions[i] = func() (interface{}, error) {
			return p.Execute(capturedInput)
		}
	}

	return Parallel(functions)
}

// ============================================================================
// Batch Processing
// ============================================================================

// BatchProcessor processes items in batches with parallel execution
type BatchProcessor struct {
	batchSize  int
	numWorkers int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize, numWorkers int) *BatchProcessor {
	if batchSize <= 0 {
		batchSize = 10
	}
	if numWorkers <= 0 {
		numWorkers = 4
	}

	return &BatchProcessor{
		batchSize:  batchSize,
		numWorkers: numWorkers,
	}
}

// Process processes items in batches
func (bp *BatchProcessor) Process(
	items []interface{},
	processFn func(interface{}) (interface{}, error),
) []ParallelResult {
	if len(items) == 0 {
		return []ParallelResult{}
	}

	// Create batches
	numBatches := (len(items) + bp.batchSize - 1) / bp.batchSize
	batches := make([][]interface{}, numBatches)

	for i := 0; i < numBatches; i++ {
		start := i * bp.batchSize
		end := start + bp.batchSize
		if end > len(items) {
			end = len(items)
		}
		batches[i] = items[start:end]
	}

	// Process batches in parallel
	results := make([]ParallelResult, len(items))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, bp.numWorkers)

	for batchIdx, batch := range batches {
		wg.Add(1)
		go func(idx int, batchItems []interface{}) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Process each item in batch
			for itemIdx, item := range batchItems {
				globalIdx := idx*bp.batchSize + itemIdx

				// Handle panics
				func() {
					defer func() {
						if r := recover(); r != nil {
							results[globalIdx] = ParallelResult{
								Index: globalIdx,
								Error: fmt.Errorf("panic processing item %d: %v", globalIdx, r),
							}
						}
					}()

					result, err := processFn(item)
					results[globalIdx] = ParallelResult{
						Index:  globalIdx,
						Result: result,
						Error:  err,
					}
				}()
			}
		}(batchIdx, batch)
	}

	wg.Wait()
	return results
}
