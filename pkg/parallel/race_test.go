package parallel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// Race Condition Tests - Run with: go test -race
// ============================================================================

// TestRaceCOWArrayConcurrentReads tests concurrent reads are safe
func TestRaceCOWArrayConcurrentReads(t *testing.T) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	var wg sync.WaitGroup
	readers := 10

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Concurrent reads
			for j := 0; j < 100; j++ {
				_ = cowArr.Get(j % cowArr.Len())
			}
		}()
	}

	wg.Wait()
}

// TestRaceCOWArrayConcurrentClones tests concurrent cloning is safe
func TestRaceCOWArrayConcurrentClones(t *testing.T) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	var wg sync.WaitGroup
	cloners := 10

	for i := 0; i < cloners; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Concurrent clones
			for j := 0; j < 10; j++ {
				clone := cowArr.Clone()
				_ = clone.Get(0)
			}
		}()
	}

	wg.Wait()
}

// TestRaceCOWArrayConcurrentWrites tests concurrent writes trigger proper copies
func TestRaceCOWArrayConcurrentWrites(t *testing.T) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = 0
	}
	cowArr := NewCOWArray(data)

	// Create clones
	clones := make([]*COWArray, 10)
	for i := 0; i < 10; i++ {
		clones[i] = cowArr.Clone()
	}

	var wg sync.WaitGroup

	// Each clone writes concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		workerIdx := i
		go func() {
			defer wg.Done()

			clone := clones[workerIdx]
			for j := 0; j < clone.Len(); j++ {
				clone.Set(j, workerIdx)
			}
		}()
	}

	wg.Wait()

	// Verify each clone has its own data
	for i := 0; i < 10; i++ {
		for j := 0; j < clones[i].Len(); j++ {
			val := clones[i].Get(j).(int)
			if val != i {
				t.Errorf("Clone %d index %d = %d, expected %d", i, j, val, i)
			}
		}
	}
}

// TestRaceCOWMapConcurrentAccess tests COW map concurrent access
func TestRaceCOWMapConcurrentAccess(t *testing.T) {
	cowMap := NewCOWMap()
	cowMap.Set("key1", "value1")
	cowMap.Set("key2", "value2")

	var wg sync.WaitGroup
	readers := 10

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				_, _ = cowMap.Get("key1")
				_ = cowMap.Has("key2")
				_ = cowMap.Len()
			}
		}()
	}

	wg.Wait()
}

// TestRaceCOWStringConcurrentAccess tests COW string concurrent access
func TestRaceCOWStringConcurrentAccess(t *testing.T) {
	cowStr := NewCOWString("hello world")

	var wg sync.WaitGroup
	readers := 10

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				_ = cowStr.String()
				_ = cowStr.Len()
				_ = cowStr.IsShared()
			}
		}()
	}

	wg.Wait()
}

// TestRaceWorkerPoolConcurrentSubmit tests concurrent task submission
func TestRaceWorkerPoolConcurrentSubmit(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	var counter atomic.Int64
	var wg sync.WaitGroup
	var taskWg sync.WaitGroup

	// Multiple goroutines submitting tasks
	submitters := 10
	tasksPerSubmitter := 100

	for i := 0; i < submitters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < tasksPerSubmitter; j++ {
				taskWg.Add(1)
				task := NewTask("bench", func() (interface{}, error) {
					counter.Add(1)
					taskWg.Done()
					return nil, nil
				})
				pool.Submit(task)
			}
		}()
	}

	wg.Wait()
	taskWg.Wait()

	expected := int64(submitters * tasksPerSubmitter)
	if counter.Load() != expected {
		t.Errorf("Expected %d tasks executed, got %d", expected, counter.Load())
	}
}

// TestRaceRequestContextConcurrentGlobals tests concurrent global access
func TestRaceRequestContextConcurrentGlobals(t *testing.T) {
	ctx := NewRequestContext("test")

	var wg sync.WaitGroup
	recorders := 10
	recordsPerRecorder := 100

	for i := 0; i < recorders; i++ {
		wg.Add(1)
		workerIdx := i
		go func() {
			defer wg.Done()

			for j := 0; j < recordsPerRecorder; j++ {
				ctx.SetGlobal(fmt.Sprintf("worker-%d-key-%d", workerIdx, j), j)
				_, _ = ctx.GetGlobal(fmt.Sprintf("worker-%d-key-%d", workerIdx, j))
			}
		}()
	}

	wg.Wait()

	expected := recorders * recordsPerRecorder
	if ctx.GlobalCount() != expected {
		t.Errorf("Expected %d globals, got %d", expected, ctx.GlobalCount())
	}
}

// TestRaceCOWManagerConcurrentRecording tests COW manager concurrent recording
func TestRaceCOWManagerConcurrentRecording(t *testing.T) {
	mgr := NewCOWManager()

	var wg sync.WaitGroup
	recorders := 10
	recordsPerRecorder := 100

	for i := 0; i < recorders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < recordsPerRecorder; j++ {
				mgr.RecordShare(1024)
				if j%10 == 0 {
					mgr.RecordCopy(512)
				}
			}
		}()
	}

	wg.Wait()

	stats := mgr.GetStats()
	expectedShares := int64(recorders * recordsPerRecorder)
	if stats.TotalShares != expectedShares {
		t.Errorf("Expected %d shares, got %d", expectedShares, stats.TotalShares)
	}
}

// TestRaceBarrierConcurrentWait tests barrier with concurrent waiters
func TestRaceBarrierConcurrentWait(t *testing.T) {
	count := 10
	barrier := NewBarrier(count)

	var wg sync.WaitGroup
	var counter atomic.Int32

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Increment before barrier
			counter.Add(1)

			// Wait at barrier
			barrier.Wait()

			// All should see count after barrier
			if counter.Load() != int32(count) {
				t.Errorf("Expected counter=%d, got %d", count, counter.Load())
			}
		}()
	}

	wg.Wait()
}

// TestRaceParallelArrayMapConcurrent tests concurrent parallel map operations
func TestRaceParallelArrayMapConcurrent(t *testing.T) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	config := DefaultArrayMapConfig()
	config.MinSize = 500
	config.NumWorkers = 4

	var wg sync.WaitGroup
	operations := 10

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := ParallelArrayMap(data, func(v interface{}) (interface{}, error) {
				return v.(int) * 2, nil
			}, config)

			if err != nil {
				t.Errorf("ParallelArrayMap failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

// TestRaceParallelArrayFilterConcurrent tests concurrent filter operations
func TestRaceParallelArrayFilterConcurrent(t *testing.T) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	config := DefaultArrayFilterConfig()
	config.MinSize = 500
	config.NumWorkers = 4

	var wg sync.WaitGroup
	operations := 10

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := ParallelArrayFilter(data, func(v interface{}) (bool, error) {
				return v.(int)%2 == 0, nil
			}, config)

			if err != nil {
				t.Errorf("ParallelArrayFilter failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

// TestRaceMixedOperations tests mixed operations for race conditions
func TestRaceMixedOperations(t *testing.T) {
	pool := NewWorkerPool(8)
	defer pool.Shutdown()

	cowMgr := NewCOWManager()

	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	var wg sync.WaitGroup
	workers := 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		workerIdx := i
		task := NewTask(fmt.Sprintf("worker-%d", i), func() (interface{}, error) {
			// Clone COW array
			clone := cowArr.Clone()
			cowMgr.RecordShare(int64(EstimateArraySize(data)))

			// Read operations
			for j := 0; j < 10; j++ {
				_ = clone.Get(j)
			}

			// Write operations (triggers copy)
			if workerIdx%2 == 0 {
				clone.Set(0, workerIdx)
				cowMgr.RecordCopy(int64(EstimateArraySize(data)))
			}
			wg.Done()
			return nil, nil
		})
		pool.Submit(task)
	}

	wg.Wait()

	// Verify COW stats
	stats := cowMgr.GetStats()
	if stats.TotalShares != int64(workers) {
		t.Errorf("Expected %d shares, got %d", workers, stats.TotalShares)
	}
}

// TestRaceCOWArrayRefCount tests reference counting is race-free
func TestRaceCOWArrayRefCount(t *testing.T) {
	data := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	var wg sync.WaitGroup
	operations := 100

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Clone and check ref count
			clone := cowArr.Clone()
			_ = clone.RefCount()
			_ = cowArr.RefCount()
			_ = cowArr.IsShared()
		}()
	}

	wg.Wait()
}

// TestRaceGlobalCOWManager tests global COW manager is race-free
func TestRaceGlobalCOWManager(t *testing.T) {
	var wg sync.WaitGroup
	operations := 100

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mgr := GetGlobalCOWManager()
			mgr.RecordShare(1024)
			_ = mgr.GetStats()
		}()
	}

	wg.Wait()
}

// TestRaceAutoThresholdConcurrent tests auto-threshold is race-free
func TestRaceAutoThresholdConcurrent(t *testing.T) {
	at := NewAutoThreshold()

	var wg sync.WaitGroup
	recorders := 10

	for i := 0; i < recorders; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				measurement := ThresholdMeasurement{
					ArraySize:      100 * (idx + 1),
					Sequential:     1000,
					Parallel:       500,
					SpeedupFactor:  2.0,
					Recommendation: true,
				}
				at.RecordMeasurement(measurement)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				_ = at.GetRecommendedThreshold()
				_ = at.GetMeasurements()
			}
		}()
	}

	wg.Wait()
}

// TestRaceConcurrentPoolShutdown tests pool shutdown is safe
func TestRaceConcurrentPoolShutdown(t *testing.T) {
	pool := NewWorkerPool(4)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		pool.Shutdown()
	}()

	// Submit tasks while shutdown is happening
	for i := 0; i < 100; i++ {
		task := NewTask(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			return nil, nil
		})
		pool.Submit(task)
	}

	wg.Wait()
}

// TestRaceCOWMapKeys tests concurrent key access
func TestRaceCOWMapKeys(t *testing.T) {
	cowMap := NewCOWMap()
	cowMap.Set("key1", 1)
	cowMap.Set("key2", 2)
	cowMap.Set("key3", 3)

	var wg sync.WaitGroup

	// Concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				_ = cowMap.Keys()
				_ = cowMap.ToMap()
			}
		}()
	}

	wg.Wait()
}

// TestRaceParallelReduceConcurrent tests concurrent reduce operations
func TestRaceParallelReduceConcurrent(t *testing.T) {
	data := make([]interface{}, 2000)
	for i := 0; i < 2000; i++ {
		data[i] = 1
	}

	config := DefaultArrayReduceConfig()
	config.MinSize = 1000
	config.NumWorkers = 4

	var wg sync.WaitGroup
	operations := 10

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sum, err := ParallelArrayReduce(data, func(acc, val interface{}) (interface{}, error) {
				return acc.(int) + val.(int), nil
			}, 0, config)

			if err != nil {
				t.Errorf("ParallelArrayReduce failed: %v", err)
			}

			if sum.(int) != 2000 {
				t.Errorf("Expected sum=2000, got %d", sum)
			}
		}()
	}

	wg.Wait()
}

// TestRaceCOWStringMutations tests concurrent string mutations
func TestRaceCOWStringMutations(t *testing.T) {
	cowStr := NewCOWString("hello")

	// Create clones
	clones := make([]*COWString, 10)
	for i := 0; i < 10; i++ {
		clones[i] = cowStr.Clone()
	}

	var wg sync.WaitGroup

	// Each clone mutates concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		workerIdx := i
		go func() {
			defer wg.Done()

			clone := clones[workerIdx]
			for j := 0; j < 10; j++ {
				clone.Append(" world")
			}
		}()
	}

	wg.Wait()

	// Verify original is unchanged
	if cowStr.String() != "hello" {
		t.Errorf("Original should be 'hello', got '%s'", cowStr.String())
	}
}

// TestRaceParallelWalkConcurrent tests concurrent walk operations
func TestRaceParallelWalkConcurrent(t *testing.T) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	config := DefaultArrayWalkConfig()
	config.MinSize = 500
	config.NumWorkers = 4

	var counter atomic.Int64
	var wg sync.WaitGroup
	operations := 10

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := ParallelArrayWalk(data, func(v interface{}, idx int) error {
				counter.Add(1)
				return nil
			}, config)

			if err != nil {
				t.Errorf("ParallelArrayWalk failed: %v", err)
			}
		}()
	}

	wg.Wait()

	expected := int64(operations * 1000)
	if counter.Load() != expected {
		t.Errorf("Expected %d walk calls, got %d", expected, counter.Load())
	}
}
