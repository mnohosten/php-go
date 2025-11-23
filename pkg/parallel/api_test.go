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
// Channel Tests
// ============================================================================

func TestNewChannel(t *testing.T) {
	ch := NewChannel(10)

	if ch.Capacity() != 10 {
		t.Errorf("Expected capacity 10, got %d", ch.Capacity())
	}

	if ch.IsClosed() {
		t.Error("New channel should not be closed")
	}
}

func TestChannelSendReceive(t *testing.T) {
	ch := NewChannel(0) // Unbuffered

	// Send in goroutine
	go func() {
		err := ch.Send("hello")
		if err != nil {
			t.Errorf("Failed to send: %v", err)
		}
	}()

	// Receive
	value, ok := ch.Receive()
	if !ok {
		t.Error("Receive failed")
	}

	if value.(string) != "hello" {
		t.Errorf("Expected 'hello', got '%v'", value)
	}
}

func TestChannelBuffered(t *testing.T) {
	ch := NewChannel(3)

	// Send without blocking
	ch.Send(1)
	ch.Send(2)
	ch.Send(3)

	// Receive
	v1, _ := ch.Receive()
	v2, _ := ch.Receive()
	v3, _ := ch.Receive()

	if v1.(int) != 1 || v2.(int) != 2 || v3.(int) != 3 {
		t.Error("Values don't match")
	}
}

func TestChannelClose(t *testing.T) {
	ch := NewChannel(1)

	ch.Send("test")
	ch.Close()

	if !ch.IsClosed() {
		t.Error("Channel should be closed")
	}

	// Sending on closed channel should error
	err := ch.Send("fail")
	if err == nil {
		t.Error("Should error when sending on closed channel")
	}

	// Receiving from closed channel should work for buffered data
	value, ok := ch.Receive()
	if !ok || value.(string) != "test" {
		t.Error("Should receive buffered data from closed channel")
	}

	// Second receive should return closed
	_, ok = ch.Receive()
	if ok {
		t.Error("Receive from empty closed channel should return false")
	}
}

func TestChannelTryReceive(t *testing.T) {
	ch := NewChannel(1)

	// Try receive from empty channel
	_, ok := ch.TryReceive()
	if ok {
		t.Error("TryReceive should return false for empty channel")
	}

	// Send and try receive
	ch.Send("data")
	value, ok := ch.TryReceive()
	if !ok || value.(string) != "data" {
		t.Error("TryReceive should succeed with data")
	}
}

func TestChannelReceiveWithTimeout(t *testing.T) {
	ch := NewChannel(0)

	// Timeout
	_, ok := ch.ReceiveWithTimeout(50 * time.Millisecond)
	if ok {
		t.Error("Should timeout")
	}

	// Send and receive within timeout
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch.Send("delayed")
	}()

	value, ok := ch.ReceiveWithTimeout(100 * time.Millisecond)
	if !ok || value.(string) != "delayed" {
		t.Error("Should receive within timeout")
	}
}

// ============================================================================
// Goroutine Tests
// ============================================================================

func TestNewGoroutine(t *testing.T) {
	gr := NewGoroutine("test-1", func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	})

	if gr.ID() != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", gr.ID())
	}

	result, err := gr.Wait()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.(int) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestGoroutineError(t *testing.T) {
	expectedError := errors.New("test error")

	gr := NewGoroutine("test-1", func() (interface{}, error) {
		return nil, expectedError
	})

	result, err := gr.Wait()
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestGoroutinePanic(t *testing.T) {
	gr := NewGoroutine("test-1", func() (interface{}, error) {
		panic("something went wrong")
	})

	_, err := gr.Wait()
	if err == nil {
		t.Error("Expected error from panic")
	}
}

func TestGoroutineIsDone(t *testing.T) {
	gr := NewGoroutine("test-1", func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return nil, nil
	})

	if gr.IsDone() {
		t.Error("Goroutine should not be done immediately")
	}

	gr.Wait()

	if !gr.IsDone() {
		t.Error("Goroutine should be done after Wait()")
	}
}

// ============================================================================
// Parallel Execution Tests
// ============================================================================

func TestParallel(t *testing.T) {
	functions := []func() (interface{}, error){
		func() (interface{}, error) { return 1, nil },
		func() (interface{}, error) { return 2, nil },
		func() (interface{}, error) { return 3, nil },
	}

	results := Parallel(functions)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("Result %d has error: %v", i, result.Error)
		}
		if result.Result.(int) != i+1 {
			t.Errorf("Result %d: expected %d, got %v", i, i+1, result.Result)
		}
	}
}

func TestParallelEmpty(t *testing.T) {
	results := Parallel([]func() (interface{}, error){})

	if len(results) != 0 {
		t.Error("Expected empty results for empty input")
	}
}

func TestParallelWithErrors(t *testing.T) {
	functions := []func() (interface{}, error){
		func() (interface{}, error) { return 1, nil },
		func() (interface{}, error) { return nil, errors.New("error 2") },
		func() (interface{}, error) { return 3, nil },
	}

	results := Parallel(functions)

	if results[0].Error != nil {
		t.Error("Result 0 should not have error")
	}

	if results[1].Error == nil {
		t.Error("Result 1 should have error")
	}

	if results[2].Error != nil {
		t.Error("Result 2 should not have error")
	}
}

func TestParallelWithPanic(t *testing.T) {
	functions := []func() (interface{}, error){
		func() (interface{}, error) { return 1, nil },
		func() (interface{}, error) { panic("test panic") },
		func() (interface{}, error) { return 3, nil },
	}

	results := Parallel(functions)

	if results[1].Error == nil {
		t.Error("Result 1 should have panic error")
	}

	// Other results should still succeed
	if results[0].Error != nil || results[2].Error != nil {
		t.Error("Other results should succeed despite panic")
	}
}

func TestParallelWithLimit(t *testing.T) {
	counter := int32(0)
	maxConcurrent := int32(0)
	var mu sync.Mutex

	functions := make([]func() (interface{}, error), 10)
	for i := 0; i < 10; i++ {
		functions[i] = func() (interface{}, error) {
			current := atomic.AddInt32(&counter, 1)

			mu.Lock()
			if current > maxConcurrent {
				maxConcurrent = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			atomic.AddInt32(&counter, -1)
			return nil, nil
		}
	}

	results := ParallelWithLimit(functions, 3)

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	mu.Lock()
	if maxConcurrent > 3 {
		t.Errorf("Max concurrent should be <= 3, got %d", maxConcurrent)
	}
	mu.Unlock()
}

// ============================================================================
// WaitGroup Tests
// ============================================================================

func TestWaitGroup(t *testing.T) {
	wg := NewWaitGroup()
	counter := int32(0)

	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
		}()
	}

	wg.Wait()

	if atomic.LoadInt32(&counter) != 3 {
		t.Errorf("Expected counter 3, got %d", counter)
	}
}

// ============================================================================
// ParallelExecutor Tests
// ============================================================================

func TestParallelExecutor(t *testing.T) {
	executor := NewParallelExecutor(4)
	defer executor.Shutdown()

	future := executor.Execute("task-1", func() (interface{}, error) {
		return 42, nil
	})

	result, err := future.Get()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.(int) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestParallelExecutorMany(t *testing.T) {
	executor := NewParallelExecutor(4)
	defer executor.Shutdown()

	functions := map[string]func() (interface{}, error){
		"task-1": func() (interface{}, error) { return 1, nil },
		"task-2": func() (interface{}, error) { return 2, nil },
		"task-3": func() (interface{}, error) { return 3, nil },
	}

	futures := executor.ExecuteMany(functions)

	if len(futures) != 3 {
		t.Errorf("Expected 3 futures, got %d", len(futures))
	}

	for id, future := range futures {
		_, err := future.Get()
		if err != nil {
			t.Errorf("Task %s failed: %v", id, err)
		}
	}
}

func TestParallelExecutorWait(t *testing.T) {
	executor := NewParallelExecutor(2)
	defer executor.Shutdown()

	completed := int32(0)

	for i := 0; i < 5; i++ {
		executor.Execute(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil, nil
		})
	}

	executor.Wait()

	if atomic.LoadInt32(&completed) != 5 {
		t.Errorf("Expected 5 completed, got %d", completed)
	}
}

func TestParallelExecutorStats(t *testing.T) {
	executor := NewParallelExecutor(2)
	defer executor.Shutdown()

	futures := make([]*Future, 3)
	for i := 0; i < 3; i++ {
		futures[i] = executor.Execute(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		})
	}

	for _, future := range futures {
		future.Get()
	}

	stats := executor.Stats()
	if stats.TotalTasks != 3 {
		t.Errorf("Expected 3 total tasks, got %d", stats.TotalTasks)
	}
}

// ============================================================================
// Pipeline Tests
// ============================================================================

func TestPipeline(t *testing.T) {
	pipeline := NewPipeline()

	// Stage 1: Double the value
	pipeline.AddStage("double", func(input interface{}) (interface{}, error) {
		return input.(int) * 2, nil
	}, false)

	// Stage 2: Add 10
	pipeline.AddStage("add10", func(input interface{}) (interface{}, error) {
		return input.(int) + 10, nil
	}, false)

	result, err := pipeline.Execute(5)
	if err != nil {
		t.Errorf("Pipeline error: %v", err)
	}

	if result.(int) != 20 { // (5 * 2) + 10
		t.Errorf("Expected 20, got %v", result)
	}
}

func TestPipelineError(t *testing.T) {
	pipeline := NewPipeline()

	pipeline.AddStage("fail", func(input interface{}) (interface{}, error) {
		return nil, errors.New("stage error")
	}, false)

	_, err := pipeline.Execute(5)
	if err == nil {
		t.Error("Expected error from pipeline")
	}
}

func TestPipelineExecuteMany(t *testing.T) {
	pipeline := NewPipeline()

	pipeline.AddStage("double", func(input interface{}) (interface{}, error) {
		return input.(int) * 2, nil
	}, false)

	inputs := []interface{}{1, 2, 3, 4, 5}
	results := pipeline.ExecuteMany(inputs)

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("Result %d has error: %v", i, result.Error)
		}
		expected := (i + 1) * 2
		if result.Result.(int) != expected {
			t.Errorf("Result %d: expected %d, got %v", i, expected, result.Result)
		}
	}
}

// ============================================================================
// BatchProcessor Tests
// ============================================================================

func TestBatchProcessor(t *testing.T) {
	processor := NewBatchProcessor(3, 2)

	items := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		items[i] = i
	}

	results := processor.Process(items, func(item interface{}) (interface{}, error) {
		return item.(int) * 2, nil
	})

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("Result %d has error: %v", i, result.Error)
		}
		expected := i * 2
		if result.Result.(int) != expected {
			t.Errorf("Result %d: expected %d, got %v", i, expected, result.Result)
		}
	}
}

func TestBatchProcessorEmpty(t *testing.T) {
	processor := NewBatchProcessor(3, 2)

	results := processor.Process([]interface{}{}, func(item interface{}) (interface{}, error) {
		return nil, nil
	})

	if len(results) != 0 {
		t.Error("Expected empty results for empty input")
	}
}

func TestBatchProcessorError(t *testing.T) {
	processor := NewBatchProcessor(3, 2)

	items := []interface{}{1, 2, 3}
	results := processor.Process(items, func(item interface{}) (interface{}, error) {
		if item.(int) == 2 {
			return nil, errors.New("item error")
		}
		return item, nil
	})

	if results[1].Error == nil {
		t.Error("Item 2 should have error")
	}

	if results[0].Error != nil || results[2].Error != nil {
		t.Error("Other items should succeed")
	}
}

func TestBatchProcessorPanic(t *testing.T) {
	processor := NewBatchProcessor(3, 2)

	items := []interface{}{1, 2, 3}
	results := processor.Process(items, func(item interface{}) (interface{}, error) {
		if item.(int) == 2 {
			panic("test panic")
		}
		return item, nil
	})

	if results[1].Error == nil {
		t.Error("Item 2 should have panic error")
	}

	// Other items should still process
	if results[0].Error != nil || results[2].Error != nil {
		t.Error("Other items should succeed despite panic")
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestChannelWithGoroutines(t *testing.T) {
	ch := NewChannel(0)
	wg := NewWaitGroup()

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch.Send(i)
		}
		ch.Close()
	}()

	// Consumer
	results := make([]int, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			value, ok := ch.Receive()
			if !ok {
				break
			}
			results = append(results, value.(int))
		}
	}()

	wg.Wait()

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}
}

func TestParallelExecutorWithPipeline(t *testing.T) {
	executor := NewParallelExecutor(4)
	defer executor.Shutdown()

	pipeline := NewPipeline()
	pipeline.AddStage("multiply", func(input interface{}) (interface{}, error) {
		return input.(int) * 3, nil
	}, false)

	futures := make([]*Future, 5)
	for i := 0; i < 5; i++ {
		capturedI := i
		futures[i] = executor.Execute(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			return pipeline.Execute(capturedI)
		})
	}

	for i, future := range futures {
		result, err := future.Get()
		if err != nil {
			t.Errorf("Task %d failed: %v", i, err)
		}
		if result.(int) != i*3 {
			t.Errorf("Task %d: expected %d, got %v", i, i*3, result)
		}
	}
}
