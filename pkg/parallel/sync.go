package parallel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Synchronization Primitives - Thread-safe primitives for PHP
// ============================================================================

// These provide PHP wrappers around Go's sync primitives for safe
// concurrent programming in PHP-Go

// ============================================================================
// Mutex - Mutual exclusion lock
// ============================================================================

// Mutex is a mutual exclusion lock wrapper for PHP
type Mutex struct {
	mu     sync.Mutex
	locked atomic.Bool
	holder string // Debug: who holds the lock
}

// NewMutex creates a new mutex
func NewMutex() *Mutex {
	return &Mutex{}
}

// Lock acquires the mutex
func (m *Mutex) Lock() {
	m.mu.Lock()
	m.locked.Store(true)
}

// Unlock releases the mutex
func (m *Mutex) Unlock() {
	m.locked.Store(false)
	m.mu.Unlock()
}

// TryLock attempts to acquire the mutex without blocking
// Returns true if lock was acquired, false otherwise
func (m *Mutex) TryLock() bool {
	acquired := m.mu.TryLock()
	if acquired {
		m.locked.Store(true)
	}
	return acquired
}

// IsLocked returns true if the mutex is currently locked
// Note: This is best-effort and may be stale immediately after returning
func (m *Mutex) IsLocked() bool {
	return m.locked.Load()
}

// ============================================================================
// RWMutex - Read-write mutual exclusion lock
// ============================================================================

// RWMutex is a read-write mutex wrapper for PHP
// Allows multiple readers or a single writer
type RWMutex struct {
	mu           sync.RWMutex
	readers      atomic.Int32
	writer       atomic.Bool
	writerHolder string // Debug: who holds the write lock
}

// NewRWMutex creates a new read-write mutex
func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

// RLock acquires a read lock
func (rw *RWMutex) RLock() {
	rw.mu.RLock()
	rw.readers.Add(1)
}

// RUnlock releases a read lock
func (rw *RWMutex) RUnlock() {
	rw.readers.Add(-1)
	rw.mu.RUnlock()
}

// Lock acquires a write lock
func (rw *RWMutex) Lock() {
	rw.mu.Lock()
	rw.writer.Store(true)
}

// Unlock releases a write lock
func (rw *RWMutex) Unlock() {
	rw.writer.Store(false)
	rw.mu.Unlock()
}

// TryRLock attempts to acquire a read lock without blocking
func (rw *RWMutex) TryRLock() bool {
	acquired := rw.mu.TryRLock()
	if acquired {
		rw.readers.Add(1)
	}
	return acquired
}

// TryLock attempts to acquire a write lock without blocking
func (rw *RWMutex) TryLock() bool {
	acquired := rw.mu.TryLock()
	if acquired {
		rw.writer.Store(true)
	}
	return acquired
}

// ReaderCount returns the current number of readers (approximate)
func (rw *RWMutex) ReaderCount() int32 {
	return rw.readers.Load()
}

// HasWriter returns true if a writer holds the lock
func (rw *RWMutex) HasWriter() bool {
	return rw.writer.Load()
}

// ============================================================================
// Semaphore - Counting semaphore
// ============================================================================

// Semaphore is a counting semaphore for controlling access to resources
type Semaphore struct {
	sem     chan struct{}
	current atomic.Int32
	max     int32
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
	if capacity <= 0 {
		capacity = 1
	}

	return &Semaphore{
		sem: make(chan struct{}, capacity),
		max: int32(capacity),
	}
}

// Acquire acquires a semaphore permit, blocking if necessary
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
	s.current.Add(1)
}

// Release releases a semaphore permit
func (s *Semaphore) Release() {
	<-s.sem
	s.current.Add(-1)
}

// TryAcquire attempts to acquire a permit without blocking
// Returns true if acquired, false otherwise
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		s.current.Add(1)
		return true
	default:
		return false
	}
}

// AcquireWithTimeout attempts to acquire with a timeout
func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case s.sem <- struct{}{}:
		s.current.Add(1)
		return true
	case <-time.After(timeout):
		return false
	}
}

// Available returns the number of available permits
func (s *Semaphore) Available() int32 {
	return s.max - s.current.Load()
}

// Capacity returns the maximum capacity
func (s *Semaphore) Capacity() int32 {
	return s.max
}

// ============================================================================
// Atomic Operations - Thread-safe atomic operations
// ============================================================================

// AtomicInt32 provides atomic operations on int32 values
type AtomicInt32 struct {
	value atomic.Int32
}

// NewAtomicInt32 creates a new atomic int32 with initial value
func NewAtomicInt32(initial int32) *AtomicInt32 {
	a := &AtomicInt32{}
	a.value.Store(initial)
	return a
}

// Load atomically loads the value
func (a *AtomicInt32) Load() int32 {
	return a.value.Load()
}

// Store atomically stores a value
func (a *AtomicInt32) Store(val int32) {
	a.value.Store(val)
}

// Add atomically adds delta and returns the new value
func (a *AtomicInt32) Add(delta int32) int32 {
	return a.value.Add(delta)
}

// Swap atomically stores new and returns the old value
func (a *AtomicInt32) Swap(new int32) int32 {
	return a.value.Swap(new)
}

// CompareAndSwap atomically compares and swaps if equal to old
// Returns true if the swap was performed
func (a *AtomicInt32) CompareAndSwap(old, new int32) bool {
	return a.value.CompareAndSwap(old, new)
}

// Increment atomically increments by 1 and returns the new value
func (a *AtomicInt32) Increment() int32 {
	return a.value.Add(1)
}

// Decrement atomically decrements by 1 and returns the new value
func (a *AtomicInt32) Decrement() int32 {
	return a.value.Add(-1)
}

// AtomicInt64 provides atomic operations on int64 values
type AtomicInt64 struct {
	value atomic.Int64
}

// NewAtomicInt64 creates a new atomic int64 with initial value
func NewAtomicInt64(initial int64) *AtomicInt64 {
	a := &AtomicInt64{}
	a.value.Store(initial)
	return a
}

// Load atomically loads the value
func (a *AtomicInt64) Load() int64 {
	return a.value.Load()
}

// Store atomically stores a value
func (a *AtomicInt64) Store(val int64) {
	a.value.Store(val)
}

// Add atomically adds delta and returns the new value
func (a *AtomicInt64) Add(delta int64) int64 {
	return a.value.Add(delta)
}

// Swap atomically stores new and returns the old value
func (a *AtomicInt64) Swap(new int64) int64 {
	return a.value.Swap(new)
}

// CompareAndSwap atomically compares and swaps if equal to old
func (a *AtomicInt64) CompareAndSwap(old, new int64) bool {
	return a.value.CompareAndSwap(old, new)
}

// Increment atomically increments by 1 and returns the new value
func (a *AtomicInt64) Increment() int64 {
	return a.value.Add(1)
}

// Decrement atomically decrements by 1 and returns the new value
func (a *AtomicInt64) Decrement() int64 {
	return a.value.Add(-1)
}

// AtomicBool provides atomic operations on bool values
type AtomicBool struct {
	value atomic.Bool
}

// NewAtomicBool creates a new atomic bool with initial value
func NewAtomicBool(initial bool) *AtomicBool {
	a := &AtomicBool{}
	a.value.Store(initial)
	return a
}

// Load atomically loads the value
func (a *AtomicBool) Load() bool {
	return a.value.Load()
}

// Store atomically stores a value
func (a *AtomicBool) Store(val bool) {
	a.value.Store(val)
}

// Swap atomically stores new and returns the old value
func (a *AtomicBool) Swap(new bool) bool {
	return a.value.Swap(new)
}

// CompareAndSwap atomically compares and swaps if equal to old
func (a *AtomicBool) CompareAndSwap(old, new bool) bool {
	return a.value.CompareAndSwap(old, new)
}

// Toggle atomically toggles the value and returns the new value
func (a *AtomicBool) Toggle() bool {
	for {
		old := a.value.Load()
		new := !old
		if a.value.CompareAndSwap(old, new) {
			return new
		}
	}
}

// ============================================================================
// Once - Execute exactly once
// ============================================================================

// Once ensures a function is executed exactly once
type Once struct {
	done atomic.Bool
	mu   sync.Mutex
}

// NewOnce creates a new Once instance
func NewOnce() *Once {
	return &Once{}
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once
func (o *Once) Do(f func()) {
	if o.done.Load() {
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.done.Load() {
		f()
		o.done.Store(true)
	}
}

// Done returns true if Do has been called
func (o *Once) Done() bool {
	return o.done.Load()
}

// Reset resets the Once to allow Do to be called again
// WARNING: This is not thread-safe and should only be used in tests
func (o *Once) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.done.Store(false)
}

// ============================================================================
// Cond - Condition variable
// ============================================================================

// Cond implements a condition variable for waiting and signaling
type Cond struct {
	L    *Mutex
	cond *sync.Cond
}

// NewCond creates a new condition variable
func NewCond(mu *Mutex) *Cond {
	return &Cond{
		L:    mu,
		cond: sync.NewCond(&mu.mu),
	}
}

// Wait atomically unlocks L and suspends execution until awakened by Signal or Broadcast
func (c *Cond) Wait() {
	c.cond.Wait()
}

// Signal wakes one goroutine waiting on c, if there is any
func (c *Cond) Signal() {
	c.cond.Signal()
}

// Broadcast wakes all goroutines waiting on c
func (c *Cond) Broadcast() {
	c.cond.Broadcast()
}

// ============================================================================
// Barrier - Synchronization barrier
// ============================================================================

// Barrier allows multiple goroutines to wait for each other
type Barrier struct {
	size    int
	count   atomic.Int32
	waiting atomic.Int32
	mu      sync.Mutex
	cond    *sync.Cond
	epoch   atomic.Int32 // Prevents reuse issues
}

// NewBarrier creates a new barrier for n parties
func NewBarrier(n int) *Barrier {
	if n <= 0 {
		n = 1
	}

	b := &Barrier{
		size: n,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks until all parties have called Wait
// Returns error if barrier is broken
func (b *Barrier) Wait() error {
	currentEpoch := b.epoch.Load()

	b.mu.Lock()
	count := b.count.Add(1)
	b.waiting.Add(1)

	if int(count) == b.size {
		// Last one in - wake everyone and reset
		b.count.Store(0)
		b.epoch.Add(1)
		b.cond.Broadcast()

		// Wait a tiny bit for others to decrement waiting
		// Then reset waiting counter
		b.mu.Unlock()
		time.Sleep(1 * time.Millisecond)
		b.waiting.Store(0)
		return nil
	}

	// Wait for others
	for b.epoch.Load() == currentEpoch {
		b.cond.Wait()
	}

	b.waiting.Add(-1)
	b.mu.Unlock()

	return nil
}

// Waiting returns the number of parties currently waiting
func (b *Barrier) Waiting() int32 {
	return b.waiting.Load()
}

// Size returns the total number of parties
func (b *Barrier) Size() int {
	return b.size
}

// Reset resets the barrier for reuse
func (b *Barrier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.count.Store(0)
	b.waiting.Store(0)
	b.epoch.Add(1)
	b.cond.Broadcast()
}

// ============================================================================
// LockManager - Manages multiple named locks
// ============================================================================

// LockManager provides named locks for resource synchronization
type LockManager struct {
	locks map[string]*Mutex
	mu    sync.RWMutex
}

// NewLockManager creates a new lock manager
func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*Mutex),
	}
}

// Lock acquires a named lock
func (lm *LockManager) Lock(name string) {
	mu := lm.getLock(name)
	mu.Lock()
}

// Unlock releases a named lock
func (lm *LockManager) Unlock(name string) error {
	lm.mu.RLock()
	mu, exists := lm.locks[name]
	lm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("lock not found: %s", name)
	}

	mu.Unlock()
	return nil
}

// TryLock attempts to acquire a named lock without blocking
func (lm *LockManager) TryLock(name string) bool {
	mu := lm.getLock(name)
	return mu.TryLock()
}

// IsLocked returns true if a named lock is currently locked
func (lm *LockManager) IsLocked(name string) bool {
	lm.mu.RLock()
	mu, exists := lm.locks[name]
	lm.mu.RUnlock()

	if !exists {
		return false
	}

	return mu.IsLocked()
}

// DeleteLock removes a named lock (must be unlocked)
func (lm *LockManager) DeleteLock(name string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	mu, exists := lm.locks[name]
	if !exists {
		return fmt.Errorf("lock not found: %s", name)
	}

	if mu.IsLocked() {
		return fmt.Errorf("cannot delete locked lock: %s", name)
	}

	delete(lm.locks, name)
	return nil
}

// Clear removes all unlocked locks
func (lm *LockManager) Clear() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	removed := 0
	for name, mu := range lm.locks {
		if !mu.IsLocked() {
			delete(lm.locks, name)
			removed++
		}
	}

	return removed
}

// Count returns the number of named locks
func (lm *LockManager) Count() int {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return len(lm.locks)
}

// getLock gets or creates a named lock
func (lm *LockManager) getLock(name string) *Mutex {
	lm.mu.RLock()
	mu, exists := lm.locks[name]
	lm.mu.RUnlock()

	if exists {
		return mu
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Double-check after acquiring write lock
	mu, exists = lm.locks[name]
	if exists {
		return mu
	}

	mu = NewMutex()
	lm.locks[name] = mu
	return mu
}
