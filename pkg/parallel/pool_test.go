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
// Task Tests
// ============================================================================

func TestNewTask(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return "result", nil
	})

	if task.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", task.ID)
	}

	if task.Function == nil {
		t.Error("Task function should not be nil")
	}

	if task.Done == nil {
		t.Error("Task done channel should not be nil")
	}
}

func TestTaskWait(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return 42, nil
	})

	// Run task in goroutine
	go func() {
		task.Result, task.Error = task.Function()
		close(task.Done)
	}()

	result, err := task.Wait()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.(int) != 42 {
		t.Errorf("Expected result 42, got %v", result)
	}
}

// ============================================================================
// Future Tests
// ============================================================================

func TestFutureGet(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return "success", nil
	})

	future := NewFuture(task)

	// Complete task in background
	go func() {
		time.Sleep(50 * time.Millisecond)
		task.Result = "success"
		close(task.Done)
	}()

	result, err := future.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.(string) != "success" {
		t.Errorf("Expected 'success', got '%v'", result)
	}
}

func TestFutureGetWithTimeout_Success(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return "quick", nil
	})

	future := NewFuture(task)

	// Complete task quickly
	go func() {
		time.Sleep(10 * time.Millisecond)
		task.Result = "quick"
		close(task.Done)
	}()

	result, err := future.GetWithTimeout(100 * time.Millisecond)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.(string) != "quick" {
		t.Errorf("Expected 'quick', got '%v'", result)
	}
}

func TestFutureGetWithTimeout_Timeout(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return "slow", nil
	})

	future := NewFuture(task)

	// Never complete the task
	_, err := future.GetWithTimeout(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error")
	}

	if err.Error() != "timeout waiting for task test-1" {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

func TestFutureIsDone(t *testing.T) {
	task := NewTask("test-1", func() (interface{}, error) {
		return nil, nil
	})

	future := NewFuture(task)

	// Initially not done
	if future.IsDone() {
		t.Error("Future should not be done initially")
	}

	// Complete task
	close(task.Done)

	// Now should be done
	if !future.IsDone() {
		t.Error("Future should be done after completion")
	}
}

// ============================================================================
// Worker Pool Tests
// ============================================================================

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4)

	if pool.maxWorkers != 4 {
		t.Errorf("Expected 4 workers, got %d", pool.maxWorkers)
	}

	if pool.workers != nil && len(pool.workers) != 0 {
		t.Error("Workers should not be initialized until Start()")
	}
}

func TestWorkerPoolStart(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()
	defer pool.Shutdown()

	if !pool.started {
		t.Error("Pool should be marked as started")
	}

	if len(pool.workers) != 2 {
		t.Errorf("Expected 2 workers, got %d", len(pool.workers))
	}
}

func TestWorkerPoolSubmit(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()
	defer pool.Shutdown()

	counter := int32(0)

	// Submit tasks
	future1 := pool.SubmitFunc("task-1", func() (interface{}, error) {
		atomic.AddInt32(&counter, 1)
		return "result1", nil
	})

	future2 := pool.SubmitFunc("task-2", func() (interface{}, error) {
		atomic.AddInt32(&counter, 1)
		return "result2", nil
	})

	// Wait for results
	result1, err1 := future1.Get()
	if err1 != nil || result1.(string) != "result1" {
		t.Errorf("Task 1 failed: %v, %v", result1, err1)
	}

	result2, err2 := future2.Get()
	if err2 != nil || result2.(string) != "result2" {
		t.Errorf("Task 2 failed: %v, %v", result2, err2)
	}

	if atomic.LoadInt32(&counter) != 2 {
		t.Errorf("Expected counter to be 2, got %d", counter)
	}
}

func TestWorkerPoolParallel(t *testing.T) {
	pool := NewWorkerPool(4)
	pool.Start()
	defer pool.Shutdown()

	numTasks := 10
	results := make([]int, numTasks)
	var mu sync.Mutex

	futures := make([]*Future, numTasks)

	// Submit multiple tasks
	for i := 0; i < numTasks; i++ {
		idx := i
		futures[i] = pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond) // Simulate work
			mu.Lock()
			results[idx] = idx * 2
			mu.Unlock()
			return idx * 2, nil
		})
	}

	// Wait for all tasks
	for i, future := range futures {
		result, err := future.Get()
		if err != nil {
			t.Errorf("Task %d failed: %v", i, err)
		}
		if result.(int) != i*2 {
			t.Errorf("Task %d: expected %d, got %v", i, i*2, result)
		}
	}

	// Verify all results
	mu.Lock()
	for i := 0; i < numTasks; i++ {
		if results[i] != i*2 {
			t.Errorf("Result %d: expected %d, got %d", i, i*2, results[i])
		}
	}
	mu.Unlock()
}

func TestWorkerPoolWait(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()
	defer pool.Shutdown()

	completed := int32(0)

	// Submit tasks
	for i := 0; i < 5; i++ {
		pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil, nil
		})
	}

	// Wait for all tasks
	pool.Wait()

	if atomic.LoadInt32(&completed) != 5 {
		t.Errorf("Expected 5 completed tasks, got %d", completed)
	}
}

func TestWorkerPoolErrorHandling(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()
	defer pool.Shutdown()

	expectedError := errors.New("task failed")

	future := pool.SubmitFunc("error-task", func() (interface{}, error) {
		return nil, expectedError
	})

	result, err := future.Get()
	if err != expectedError {
		t.Errorf("Expected error '%v', got '%v'", expectedError, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestWorkerPoolPanicHandling(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()
	defer pool.Shutdown()

	future := pool.SubmitFunc("panic-task", func() (interface{}, error) {
		panic("something went wrong")
	})

	result, err := future.Get()
	if err == nil {
		t.Error("Expected error from panic")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Pool should still be functional
	future2 := pool.SubmitFunc("normal-task", func() (interface{}, error) {
		return "ok", nil
	})

	result2, err2 := future2.Get()
	if err2 != nil || result2.(string) != "ok" {
		t.Error("Pool should recover from panic and continue working")
	}
}

func TestWorkerPoolShutdown(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()

	// Submit some tasks
	for i := 0; i < 3; i++ {
		pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		})
	}

	// Shutdown the pool
	pool.Shutdown()

	if pool.started {
		t.Error("Pool should not be marked as started after shutdown")
	}
}

func TestWorkerPoolShutdownWithTimeout_Success(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()

	// Submit a quick task
	pool.SubmitFunc("quick-task", func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, nil
	})

	// Shutdown with generous timeout
	err := pool.ShutdownWithTimeout(1 * time.Second)
	if err != nil {
		t.Errorf("Shutdown should succeed, got error: %v", err)
	}
}

func TestWorkerPoolShutdownWithTimeout_Timeout(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()

	// Submit a slow task
	pool.SubmitFunc("slow-task", func() (interface{}, error) {
		time.Sleep(2 * time.Second)
		return nil, nil
	})

	// Give task time to start
	time.Sleep(10 * time.Millisecond)

	// Shutdown with short timeout (should timeout)
	err := pool.ShutdownWithTimeout(100 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestWorkerPoolStats(t *testing.T) {
	pool := NewWorkerPool(2)
	pool.Start()
	defer pool.Shutdown()

	// Submit tasks
	futures := make([]*Future, 5)
	for i := 0; i < 5; i++ {
		futures[i] = pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			return nil, nil
		})
	}

	// Wait for all tasks
	for _, future := range futures {
		future.Get()
	}

	stats := pool.Stats()

	if stats.TotalTasks != 5 {
		t.Errorf("Expected 5 total tasks, got %d", stats.TotalTasks)
	}

	if stats.CompletedTasks != 5 {
		t.Errorf("Expected 5 completed tasks, got %d", stats.CompletedTasks)
	}

	if stats.WorkerCount != 2 {
		t.Errorf("Expected 2 workers, got %d", stats.WorkerCount)
	}

	if len(stats.WorkerStats) != 2 {
		t.Errorf("Expected 2 worker stats, got %d", len(stats.WorkerStats))
	}
}

func TestWorkerPoolSize(t *testing.T) {
	pool := NewWorkerPool(8)

	if pool.Size() != 8 {
		t.Errorf("Expected size 8, got %d", pool.Size())
	}
}

func TestWorkerPoolZeroWorkers(t *testing.T) {
	// Should default to 1 worker
	pool := NewWorkerPool(0)

	if pool.maxWorkers != 1 {
		t.Errorf("Expected 1 worker for 0 input, got %d", pool.maxWorkers)
	}
}

// ============================================================================
// Worker Tests
// ============================================================================

func TestWorkerTaskCount(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()
	defer pool.Shutdown()

	// Submit multiple tasks
	for i := 0; i < 5; i++ {
		future := pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			return nil, nil
		})
		future.Get() // Wait for completion
	}

	// Check worker task count
	worker := pool.workers[0]
	if worker.TaskCount() != 5 {
		t.Errorf("Expected worker to have processed 5 tasks, got %d", worker.TaskCount())
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestWorkerPoolConcurrency(t *testing.T) {
	pool := NewWorkerPool(4)
	pool.Start()
	defer pool.Shutdown()

	numTasks := 100
	counter := int32(0)

	futures := make([]*Future, numTasks)
	for i := 0; i < numTasks; i++ {
		futures[i] = pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			atomic.AddInt32(&counter, 1)
			time.Sleep(1 * time.Millisecond)
			return nil, nil
		})
	}

	// Wait for all
	for _, future := range futures {
		future.Get()
	}

	if atomic.LoadInt32(&counter) != int32(numTasks) {
		t.Errorf("Expected %d tasks executed, got %d", numTasks, counter)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkWorkerPool_Submit(b *testing.B) {
	pool := NewWorkerPool(4)
	pool.Start()
	defer pool.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		future := pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			return i, nil
		})
		future.Get()
	}
}

func BenchmarkWorkerPool_Parallel(b *testing.B) {
	pool := NewWorkerPool(8)
	pool.Start()
	defer pool.Shutdown()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			future := pool.SubmitFunc(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
				return i, nil
			})
			future.Get()
			i++
		}
	})
}
