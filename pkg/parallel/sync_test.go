package parallel

import (
	"sync"
	"testing"
	"time"
)

// ============================================================================
// Mutex Tests
// ============================================================================

func TestNewMutex(t *testing.T) {
	mu := NewMutex()
	if mu == nil {
		t.Error("NewMutex returned nil")
	}
	if mu.IsLocked() {
		t.Error("New mutex should not be locked")
	}
}

func TestMutexLockUnlock(t *testing.T) {
	mu := NewMutex()

	mu.Lock()
	if !mu.IsLocked() {
		t.Error("Mutex should be locked after Lock()")
	}

	mu.Unlock()
	if mu.IsLocked() {
		t.Error("Mutex should not be locked after Unlock()")
	}
}

func TestMutexTryLock(t *testing.T) {
	mu := NewMutex()

	// Should succeed when unlocked
	if !mu.TryLock() {
		t.Error("TryLock should succeed on unlocked mutex")
	}

	// Should fail when locked
	if mu.TryLock() {
		t.Error("TryLock should fail on locked mutex")
	}

	mu.Unlock()

	// Should succeed again
	if !mu.TryLock() {
		t.Error("TryLock should succeed after unlock")
	}

	mu.Unlock()
}

func TestMutexConcurrent(t *testing.T) {
	mu := NewMutex()
	counter := 0
	iterations := 1000

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := 10 * iterations
	if counter != expected {
		t.Errorf("Counter = %d, expected %d (race condition detected)", counter, expected)
	}
}

// ============================================================================
// RWMutex Tests
// ============================================================================

func TestNewRWMutex(t *testing.T) {
	rw := NewRWMutex()
	if rw == nil {
		t.Error("NewRWMutex returned nil")
	}
	if rw.HasWriter() {
		t.Error("New RWMutex should not have writer")
	}
	if rw.ReaderCount() != 0 {
		t.Error("New RWMutex should have 0 readers")
	}
}

func TestRWMutexReadLock(t *testing.T) {
	rw := NewRWMutex()

	rw.RLock()
	if rw.ReaderCount() != 1 {
		t.Errorf("Reader count = %d, expected 1", rw.ReaderCount())
	}

	rw.RLock()
	if rw.ReaderCount() != 2 {
		t.Errorf("Reader count = %d, expected 2", rw.ReaderCount())
	}

	rw.RUnlock()
	if rw.ReaderCount() != 1 {
		t.Errorf("Reader count = %d, expected 1", rw.ReaderCount())
	}

	rw.RUnlock()
	if rw.ReaderCount() != 0 {
		t.Errorf("Reader count = %d, expected 0", rw.ReaderCount())
	}
}

func TestRWMutexWriteLock(t *testing.T) {
	rw := NewRWMutex()

	rw.Lock()
	if !rw.HasWriter() {
		t.Error("RWMutex should have writer after Lock()")
	}

	rw.Unlock()
	if rw.HasWriter() {
		t.Error("RWMutex should not have writer after Unlock()")
	}
}

func TestRWMutexTryRLock(t *testing.T) {
	rw := NewRWMutex()

	if !rw.TryRLock() {
		t.Error("TryRLock should succeed")
	}

	if rw.ReaderCount() != 1 {
		t.Error("Reader count should be 1")
	}

	rw.RUnlock()
}

func TestRWMutexTryLock(t *testing.T) {
	rw := NewRWMutex()

	if !rw.TryLock() {
		t.Error("TryLock should succeed on unlocked RWMutex")
	}

	if !rw.HasWriter() {
		t.Error("Should have writer")
	}

	rw.Unlock()
}

func TestRWMutexConcurrentReads(t *testing.T) {
	rw := NewRWMutex()
	value := 0

	var wg sync.WaitGroup
	readers := 10

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rw.RLock()
			_ = value
			time.Sleep(10 * time.Millisecond)
			rw.RUnlock()
		}()
	}

	// Give readers time to acquire locks
	time.Sleep(5 * time.Millisecond)

	// All readers should be active
	if rw.ReaderCount() == 0 {
		t.Error("Should have multiple readers")
	}

	wg.Wait()

	if rw.ReaderCount() != 0 {
		t.Error("Should have no readers after completion")
	}
}

// ============================================================================
// Semaphore Tests
// ============================================================================

func TestNewSemaphore(t *testing.T) {
	sem := NewSemaphore(5)
	if sem == nil {
		t.Error("NewSemaphore returned nil")
	}
	if sem.Capacity() != 5 {
		t.Errorf("Capacity = %d, expected 5", sem.Capacity())
	}
	if sem.Available() != 5 {
		t.Errorf("Available = %d, expected 5", sem.Available())
	}
}

func TestSemaphoreZeroCapacity(t *testing.T) {
	sem := NewSemaphore(0)
	if sem.Capacity() != 1 {
		t.Error("Zero capacity should become 1")
	}
}

func TestSemaphoreAcquireRelease(t *testing.T) {
	sem := NewSemaphore(2)

	sem.Acquire()
	if sem.Available() != 1 {
		t.Errorf("Available = %d, expected 1", sem.Available())
	}

	sem.Acquire()
	if sem.Available() != 0 {
		t.Errorf("Available = %d, expected 0", sem.Available())
	}

	sem.Release()
	if sem.Available() != 1 {
		t.Errorf("Available = %d, expected 1", sem.Available())
	}

	sem.Release()
	if sem.Available() != 2 {
		t.Errorf("Available = %d, expected 2", sem.Available())
	}
}

func TestSemaphoreTryAcquire(t *testing.T) {
	sem := NewSemaphore(1)

	if !sem.TryAcquire() {
		t.Error("TryAcquire should succeed")
	}

	if sem.TryAcquire() {
		t.Error("TryAcquire should fail when at capacity")
	}

	sem.Release()

	if !sem.TryAcquire() {
		t.Error("TryAcquire should succeed after release")
	}

	sem.Release()
}

func TestSemaphoreAcquireWithTimeout(t *testing.T) {
	sem := NewSemaphore(1)

	// Should succeed immediately
	if !sem.AcquireWithTimeout(10 * time.Millisecond) {
		t.Error("AcquireWithTimeout should succeed")
	}

	// Should timeout
	if sem.AcquireWithTimeout(10 * time.Millisecond) {
		t.Error("AcquireWithTimeout should timeout")
	}

	sem.Release()
}

func TestSemaphoreConcurrency(t *testing.T) {
	sem := NewSemaphore(3)
	maxConcurrent := 0
	current := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem.Acquire()

			mu.Lock()
			current++
			if current > maxConcurrent {
				maxConcurrent = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			current--
			mu.Unlock()

			sem.Release()
		}()
	}

	wg.Wait()

	if maxConcurrent > 3 {
		t.Errorf("Max concurrent = %d, expected <= 3", maxConcurrent)
	}
}

// ============================================================================
// AtomicInt32 Tests
// ============================================================================

func TestAtomicInt32(t *testing.T) {
	a := NewAtomicInt32(10)

	if a.Load() != 10 {
		t.Errorf("Load = %d, expected 10", a.Load())
	}

	a.Store(20)
	if a.Load() != 20 {
		t.Errorf("Load = %d, expected 20", a.Load())
	}

	result := a.Add(5)
	if result != 25 {
		t.Errorf("Add returned %d, expected 25", result)
	}
	if a.Load() != 25 {
		t.Errorf("Load = %d, expected 25", a.Load())
	}
}

func TestAtomicInt32Swap(t *testing.T) {
	a := NewAtomicInt32(10)

	old := a.Swap(20)
	if old != 10 {
		t.Errorf("Swap returned %d, expected 10", old)
	}
	if a.Load() != 20 {
		t.Errorf("Load = %d, expected 20", a.Load())
	}
}

func TestAtomicInt32CompareAndSwap(t *testing.T) {
	a := NewAtomicInt32(10)

	// Should succeed
	if !a.CompareAndSwap(10, 20) {
		t.Error("CompareAndSwap should succeed")
	}
	if a.Load() != 20 {
		t.Errorf("Load = %d, expected 20", a.Load())
	}

	// Should fail
	if a.CompareAndSwap(10, 30) {
		t.Error("CompareAndSwap should fail")
	}
	if a.Load() != 20 {
		t.Errorf("Load = %d, expected 20", a.Load())
	}
}

func TestAtomicInt32IncrementDecrement(t *testing.T) {
	a := NewAtomicInt32(10)

	if a.Increment() != 11 {
		t.Error("Increment should return 11")
	}
	if a.Increment() != 12 {
		t.Error("Increment should return 12")
	}
	if a.Decrement() != 11 {
		t.Error("Decrement should return 11")
	}
	if a.Load() != 11 {
		t.Error("Load should be 11")
	}
}

func TestAtomicInt32Concurrent(t *testing.T) {
	a := NewAtomicInt32(0)
	iterations := 1000

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				a.Increment()
			}
		}()
	}

	wg.Wait()

	expected := int32(10 * iterations)
	if a.Load() != expected {
		t.Errorf("Load = %d, expected %d", a.Load(), expected)
	}
}

// ============================================================================
// AtomicInt64 Tests
// ============================================================================

func TestAtomicInt64(t *testing.T) {
	a := NewAtomicInt64(100)

	if a.Load() != 100 {
		t.Errorf("Load = %d, expected 100", a.Load())
	}

	a.Store(200)
	if a.Load() != 200 {
		t.Errorf("Load = %d, expected 200", a.Load())
	}

	result := a.Add(50)
	if result != 250 {
		t.Errorf("Add returned %d, expected 250", result)
	}
}

func TestAtomicInt64IncrementDecrement(t *testing.T) {
	a := NewAtomicInt64(100)

	if a.Increment() != 101 {
		t.Error("Increment should return 101")
	}
	if a.Decrement() != 100 {
		t.Error("Decrement should return 100")
	}
}

// ============================================================================
// AtomicBool Tests
// ============================================================================

func TestAtomicBool(t *testing.T) {
	a := NewAtomicBool(false)

	if a.Load() {
		t.Error("Load should be false")
	}

	a.Store(true)
	if !a.Load() {
		t.Error("Load should be true")
	}
}

func TestAtomicBoolSwap(t *testing.T) {
	a := NewAtomicBool(false)

	old := a.Swap(true)
	if old {
		t.Error("Swap should return false")
	}
	if !a.Load() {
		t.Error("Load should be true")
	}
}

func TestAtomicBoolCompareAndSwap(t *testing.T) {
	a := NewAtomicBool(false)

	if !a.CompareAndSwap(false, true) {
		t.Error("CompareAndSwap should succeed")
	}
	if !a.Load() {
		t.Error("Load should be true")
	}

	if a.CompareAndSwap(false, false) {
		t.Error("CompareAndSwap should fail")
	}
}

func TestAtomicBoolToggle(t *testing.T) {
	a := NewAtomicBool(false)

	if !a.Toggle() {
		t.Error("Toggle should return true")
	}
	if a.Toggle() {
		t.Error("Toggle should return false")
	}
	if !a.Toggle() {
		t.Error("Toggle should return true")
	}
}

// ============================================================================
// Once Tests
// ============================================================================

func TestOnce(t *testing.T) {
	once := NewOnce()
	count := 0

	once.Do(func() {
		count++
	})

	if count != 1 {
		t.Errorf("Count = %d, expected 1", count)
	}

	// Should not execute again
	once.Do(func() {
		count++
	})

	if count != 1 {
		t.Errorf("Count = %d, expected 1", count)
	}

	if !once.Done() {
		t.Error("Done should be true")
	}
}

func TestOnceConcurrent(t *testing.T) {
	once := NewOnce()
	count := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				mu.Lock()
				count++
				mu.Unlock()
			})
		}()
	}

	wg.Wait()

	if count != 1 {
		t.Errorf("Count = %d, expected 1", count)
	}
}

func TestOnceReset(t *testing.T) {
	once := NewOnce()
	count := 0

	once.Do(func() {
		count++
	})

	once.Reset()

	once.Do(func() {
		count++
	})

	if count != 2 {
		t.Errorf("Count = %d, expected 2", count)
	}
}

// ============================================================================
// Cond Tests
// ============================================================================

func TestCond(t *testing.T) {
	mu := NewMutex()
	cond := NewCond(mu)

	if cond == nil {
		t.Error("NewCond returned nil")
	}

	ready := false
	var wg sync.WaitGroup

	// Waiter
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		for !ready {
			cond.Wait()
		}
		mu.Unlock()
	}()

	// Give waiter time to wait
	time.Sleep(10 * time.Millisecond)

	// Signal
	mu.Lock()
	ready = true
	cond.Signal()
	mu.Unlock()

	wg.Wait()
}

func TestCondBroadcast(t *testing.T) {
	mu := NewMutex()
	cond := NewCond(mu)

	ready := false
	count := 0
	var countMu sync.Mutex
	var wg sync.WaitGroup

	// Multiple waiters
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			for !ready {
				cond.Wait()
			}
			mu.Unlock()

			countMu.Lock()
			count++
			countMu.Unlock()
		}()
	}

	// Give waiters time to wait
	time.Sleep(20 * time.Millisecond)

	// Broadcast
	mu.Lock()
	ready = true
	cond.Broadcast()
	mu.Unlock()

	wg.Wait()

	if count != 5 {
		t.Errorf("Count = %d, expected 5", count)
	}
}

// ============================================================================
// Barrier Tests
// ============================================================================

func TestNewBarrier(t *testing.T) {
	b := NewBarrier(5)
	if b == nil {
		t.Error("NewBarrier returned nil")
	}
	if b.Size() != 5 {
		t.Errorf("Size = %d, expected 5", b.Size())
	}
	if b.Waiting() != 0 {
		t.Error("Waiting should be 0")
	}
}

func TestBarrierZeroSize(t *testing.T) {
	b := NewBarrier(0)
	if b.Size() != 1 {
		t.Error("Zero size should become 1")
	}
}

func TestBarrierWait(t *testing.T) {
	b := NewBarrier(3)
	results := make([]int, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			results[id] = 1
			err := b.Wait()
			if err != nil {
				t.Errorf("Wait error: %v", err)
			}
			results[id] = 2
		}(i)
	}

	wg.Wait()

	// All should have reached 2
	for i, val := range results {
		if val != 2 {
			t.Errorf("results[%d] = %d, expected 2", i, val)
		}
	}
}

func TestBarrierWaiting(t *testing.T) {
	b := NewBarrier(5)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Wait()
		}()
	}

	// Give goroutines time to wait
	time.Sleep(20 * time.Millisecond)

	if b.Waiting() != 4 {
		t.Errorf("Waiting = %d, expected 4", b.Waiting())
	}

	// Release all
	go b.Wait()

	wg.Wait()

	// Give time for the counter to be reset
	time.Sleep(5 * time.Millisecond)

	if b.Waiting() != 0 {
		t.Errorf("Waiting = %d, should be 0 after all released", b.Waiting())
	}
}

func TestBarrierReset(t *testing.T) {
	b := NewBarrier(3)

	// First round
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Wait()
		}()
	}
	wg.Wait()

	// Reset
	b.Reset()

	// Second round
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Wait()
		}()
	}
	wg.Wait()

	if b.Waiting() != 0 {
		t.Error("Waiting should be 0 after second round")
	}
}

// ============================================================================
// LockManager Tests
// ============================================================================

func TestNewLockManager(t *testing.T) {
	lm := NewLockManager()
	if lm == nil {
		t.Error("NewLockManager returned nil")
	}
	if lm.Count() != 0 {
		t.Error("Count should be 0")
	}
}

func TestLockManagerLockUnlock(t *testing.T) {
	lm := NewLockManager()

	lm.Lock("resource1")
	if lm.Count() != 1 {
		t.Error("Count should be 1")
	}
	if !lm.IsLocked("resource1") {
		t.Error("resource1 should be locked")
	}

	err := lm.Unlock("resource1")
	if err != nil {
		t.Errorf("Unlock error: %v", err)
	}
	if lm.IsLocked("resource1") {
		t.Error("resource1 should be unlocked")
	}
}

func TestLockManagerUnlockNonexistent(t *testing.T) {
	lm := NewLockManager()

	err := lm.Unlock("nonexistent")
	if err == nil {
		t.Error("Unlock should error for nonexistent lock")
	}
}

func TestLockManagerTryLock(t *testing.T) {
	lm := NewLockManager()

	if !lm.TryLock("resource1") {
		t.Error("TryLock should succeed")
	}

	if lm.TryLock("resource1") {
		t.Error("TryLock should fail on locked resource")
	}

	lm.Unlock("resource1")

	if !lm.TryLock("resource1") {
		t.Error("TryLock should succeed after unlock")
	}

	lm.Unlock("resource1")
}

func TestLockManagerDeleteLock(t *testing.T) {
	lm := NewLockManager()

	lm.Lock("resource1")
	lm.Unlock("resource1")

	err := lm.DeleteLock("resource1")
	if err != nil {
		t.Errorf("DeleteLock error: %v", err)
	}

	if lm.Count() != 0 {
		t.Error("Count should be 0")
	}
}

func TestLockManagerDeleteLockedLock(t *testing.T) {
	lm := NewLockManager()

	lm.Lock("resource1")

	err := lm.DeleteLock("resource1")
	if err == nil {
		t.Error("DeleteLock should error for locked lock")
	}

	lm.Unlock("resource1")
}

func TestLockManagerClear(t *testing.T) {
	lm := NewLockManager()

	lm.Lock("resource1")
	lm.Unlock("resource1")
	lm.Lock("resource2")
	lm.Unlock("resource2")
	lm.Lock("resource3") // Keep locked

	removed := lm.Clear()
	if removed != 2 {
		t.Errorf("Removed = %d, expected 2", removed)
	}
	if lm.Count() != 1 {
		t.Errorf("Count = %d, expected 1", lm.Count())
	}

	lm.Unlock("resource3")
}

func TestLockManagerConcurrent(t *testing.T) {
	lm := NewLockManager()
	counter := 0

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				lm.Lock("counter")
				counter++
				lm.Unlock("counter")
			}
		}()
	}

	wg.Wait()

	if counter != 1000 {
		t.Errorf("Counter = %d, expected 1000", counter)
	}
}
