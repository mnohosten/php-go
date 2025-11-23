package parallel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Task represents a unit of work to be executed by the worker pool
// ============================================================================

// Task represents a function to execute in parallel
type Task struct {
	ID       string
	Function func() (interface{}, error)
	Result   interface{}
	Error    error
	Done     chan struct{}
}

// NewTask creates a new task
func NewTask(id string, fn func() (interface{}, error)) *Task {
	return &Task{
		ID:       id,
		Function: fn,
		Done:     make(chan struct{}),
	}
}

// Wait blocks until the task is complete and returns the result
func (t *Task) Wait() (interface{}, error) {
	<-t.Done
	return t.Result, t.Error
}

// ============================================================================
// Future represents a value that will be available in the future
// ============================================================================

// Future represents an asynchronous result
type Future struct {
	task *Task
}

// NewFuture creates a new future from a task
func NewFuture(task *Task) *Future {
	return &Future{task: task}
}

// Get waits for and returns the future's result
func (f *Future) Get() (interface{}, error) {
	return f.task.Wait()
}

// GetWithTimeout waits for the result with a timeout
func (f *Future) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case <-f.task.Done:
		return f.task.Result, f.task.Error
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for task %s", f.task.ID)
	}
}

// IsDone returns true if the future is complete
func (f *Future) IsDone() bool {
	select {
	case <-f.task.Done:
		return true
	default:
		return false
	}
}

// ============================================================================
// Worker represents a single worker goroutine
// ============================================================================

// Worker processes tasks from a queue
type Worker struct {
	ID         int
	pool       *WorkerPool
	taskChan   chan *Task
	quitChan   chan struct{}
	wg         *sync.WaitGroup
	tasksCount int64
	mu         sync.Mutex
}

// NewWorker creates a new worker
func NewWorker(id int, pool *WorkerPool, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:       id,
		pool:     pool,
		taskChan: make(chan *Task),
		quitChan: make(chan struct{}),
		wg:       wg,
	}
}

// Start begins the worker's processing loop
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			// Register worker as available
			w.pool.workerChan <- w

			select {
			case task := <-w.taskChan:
				// Execute the task
				w.executeTask(task)

			case <-w.quitChan:
				// Worker shutdown requested
				return
			}
		}
	}()
}

// executeTask runs a task and captures the result
func (w *Worker) executeTask(task *Task) {
	defer func() {
		// Handle panics in task execution
		if r := recover(); r != nil {
			task.Error = fmt.Errorf("panic in task %s: %v", task.ID, r)
		}
		close(task.Done)
	}()

	// Increment task counter
	w.mu.Lock()
	w.tasksCount++
	w.mu.Unlock()

	// Execute the task function
	task.Result, task.Error = task.Function()
}

// Stop signals the worker to shut down
func (w *Worker) Stop() {
	close(w.quitChan)
}

// TaskCount returns the number of tasks this worker has processed
func (w *Worker) TaskCount() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.tasksCount
}

// ============================================================================
// WorkerPool manages a pool of worker goroutines
// ============================================================================

// WorkerPool manages a fixed pool of workers for parallel task execution
type WorkerPool struct {
	maxWorkers int
	workers    []*Worker
	workerChan chan *Worker
	taskQueue  chan *Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	started    bool
	mu         sync.Mutex

	// Statistics
	totalTasks    int64
	completedTasks int64
	statsMu       sync.Mutex
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		maxWorkers: maxWorkers,
		workers:    make([]*Worker, 0, maxWorkers),
		workerChan: make(chan *Worker, maxWorkers),
		taskQueue:  make(chan *Task, maxWorkers*2), // Buffered task queue
		ctx:        ctx,
		cancel:     cancel,
	}

	return pool
}

// Start initializes and starts all workers in the pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return
	}

	// Create and start workers
	for i := 0; i < wp.maxWorkers; i++ {
		worker := NewWorker(i, wp, &wp.wg)
		wp.workers = append(wp.workers, worker)
		worker.Start()
	}

	// Start the dispatcher
	go wp.dispatch()

	wp.started = true
}

// dispatch assigns tasks to available workers
func (wp *WorkerPool) dispatch() {
	for {
		select {
		case task := <-wp.taskQueue:
			// Wait for an available worker
			worker := <-wp.workerChan

			// Assign task to worker
			worker.taskChan <- task

		case <-wp.ctx.Done():
			// Pool shutdown requested
			return
		}
	}
}

// Submit adds a task to the pool and returns a future
func (wp *WorkerPool) Submit(task *Task) *Future {
	// Auto-start the pool if not started
	wp.mu.Lock()
	if !wp.started {
		wp.mu.Unlock()
		wp.Start()
	} else {
		wp.mu.Unlock()
	}

	wp.statsMu.Lock()
	wp.totalTasks++
	wp.statsMu.Unlock()

	// Wrap task to track completion
	originalFunction := task.Function
	task.Function = func() (interface{}, error) {
		result, err := originalFunction()
		wp.statsMu.Lock()
		wp.completedTasks++
		wp.statsMu.Unlock()
		return result, err
	}

	// Send task to queue, checking for shutdown
	select {
	case wp.taskQueue <- task:
		return NewFuture(task)
	case <-wp.ctx.Done():
		// Pool is shutting down, return a future with error
		task.Error = fmt.Errorf("pool is shutting down")
		close(task.Done)
		return NewFuture(task)
	}
}

// SubmitFunc is a convenience method to submit a function directly
func (wp *WorkerPool) SubmitFunc(id string, fn func() (interface{}, error)) *Future {
	task := NewTask(id, fn)
	return wp.Submit(task)
}

// Wait blocks until all submitted tasks are complete
func (wp *WorkerPool) Wait() {
	// Wait for all tasks in queue to be processed
	for {
		wp.statsMu.Lock()
		total := wp.totalTasks
		completed := wp.completedTasks
		wp.statsMu.Unlock()

		if total == completed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Shutdown gracefully shuts down the worker pool
// It waits for all workers to finish their current tasks
func (wp *WorkerPool) Shutdown() {
	wp.mu.Lock()
	if !wp.started {
		wp.mu.Unlock()
		return
	}
	wp.mu.Unlock()

	// Signal shutdown
	wp.cancel()

	// Stop all workers
	for _, worker := range wp.workers {
		worker.Stop()
	}

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close channels
	close(wp.taskQueue)
	close(wp.workerChan)

	wp.mu.Lock()
	wp.started = false
	wp.mu.Unlock()
}

// ShutdownWithTimeout shuts down the pool with a timeout
func (wp *WorkerPool) ShutdownWithTimeout(timeout time.Duration) error {
	done := make(chan struct{})

	go func() {
		wp.Shutdown()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

// Size returns the number of workers in the pool
func (wp *WorkerPool) Size() int {
	return wp.maxWorkers
}

// Stats returns pool statistics
func (wp *WorkerPool) Stats() PoolStats {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()

	stats := PoolStats{
		TotalTasks:     wp.totalTasks,
		CompletedTasks: wp.completedTasks,
		PendingTasks:   wp.totalTasks - wp.completedTasks,
		WorkerCount:    int64(wp.maxWorkers),
		QueueSize:      int64(len(wp.taskQueue)),
	}

	// Get per-worker stats
	stats.WorkerStats = make([]WorkerStats, len(wp.workers))
	for i, worker := range wp.workers {
		stats.WorkerStats[i] = WorkerStats{
			ID:        worker.ID,
			TaskCount: worker.TaskCount(),
		}
	}

	return stats
}

// ============================================================================
// Statistics Types
// ============================================================================

// PoolStats contains worker pool statistics
type PoolStats struct {
	TotalTasks     int64
	CompletedTasks int64
	PendingTasks   int64
	WorkerCount    int64
	QueueSize      int64
	WorkerStats    []WorkerStats
}

// WorkerStats contains per-worker statistics
type WorkerStats struct {
	ID        int
	TaskCount int64
}
