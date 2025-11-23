package parallel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Request Context - Isolated execution environment for each request
// ============================================================================

// RequestContext represents an isolated execution environment for a single request
// Each request gets its own context with isolated globals, output, and state
type RequestContext struct {
	ID          string
	ctx         context.Context
	cancel      context.CancelFunc
	startTime   time.Time
	timeout     time.Duration

	// Request-specific state
	globals     map[string]interface{}
	output      []byte
	errors      []error

	// Metadata
	metadata    map[string]interface{}

	mu          sync.RWMutex
}

// NewRequestContext creates a new isolated request context
func NewRequestContext(id string) *RequestContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &RequestContext{
		ID:        id,
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
		globals:   make(map[string]interface{}),
		output:    make([]byte, 0),
		errors:    make([]error, 0),
		metadata:  make(map[string]interface{}),
	}
}

// NewRequestContextWithTimeout creates a request context with a timeout
func NewRequestContextWithTimeout(id string, timeout time.Duration) *RequestContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return &RequestContext{
		ID:        id,
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
		timeout:   timeout,
		globals:   make(map[string]interface{}),
		output:    make([]byte, 0),
		errors:    make([]error, 0),
		metadata:  make(map[string]interface{}),
	}
}

// Context returns the underlying context.Context
func (rc *RequestContext) Context() context.Context {
	return rc.ctx
}

// Cancel cancels the request context
func (rc *RequestContext) Cancel() {
	rc.cancel()
}

// IsCancelled returns true if the context has been cancelled
func (rc *RequestContext) IsCancelled() bool {
	select {
	case <-rc.ctx.Done():
		return true
	default:
		return false
	}
}

// IsTimedOut returns true if the context has timed out
func (rc *RequestContext) IsTimedOut() bool {
	if rc.timeout == 0 {
		return false
	}
	return time.Since(rc.startTime) > rc.timeout
}

// Elapsed returns the time since the request started
func (rc *RequestContext) Elapsed() time.Duration {
	return time.Since(rc.startTime)
}

// ============================================================================
// Global Variable Management
// ============================================================================

// SetGlobal sets a global variable for this request
func (rc *RequestContext) SetGlobal(name string, value interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.globals[name] = value
}

// GetGlobal retrieves a global variable for this request
func (rc *RequestContext) GetGlobal(name string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	value, exists := rc.globals[name]
	return value, exists
}

// HasGlobal checks if a global variable exists
func (rc *RequestContext) HasGlobal(name string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	_, exists := rc.globals[name]
	return exists
}

// DeleteGlobal removes a global variable
func (rc *RequestContext) DeleteGlobal(name string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.globals, name)
}

// ClearGlobals removes all global variables
func (rc *RequestContext) ClearGlobals() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.globals = make(map[string]interface{})
}

// GlobalCount returns the number of global variables
func (rc *RequestContext) GlobalCount() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.globals)
}

// ============================================================================
// Output Management
// ============================================================================

// Write appends data to the request output buffer
func (rc *RequestContext) Write(data []byte) (int, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.output = append(rc.output, data...)
	return len(data), nil
}

// WriteString appends a string to the request output buffer
func (rc *RequestContext) WriteString(s string) (int, error) {
	return rc.Write([]byte(s))
}

// GetOutput returns the current output buffer
func (rc *RequestContext) GetOutput() []byte {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	// Return a copy to prevent modification
	output := make([]byte, len(rc.output))
	copy(output, rc.output)
	return output
}

// ClearOutput clears the output buffer
func (rc *RequestContext) ClearOutput() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.output = make([]byte, 0)
}

// OutputSize returns the size of the output buffer
func (rc *RequestContext) OutputSize() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.output)
}

// ============================================================================
// Error Management
// ============================================================================

// AddError adds an error to the request context
func (rc *RequestContext) AddError(err error) {
	if err == nil {
		return
	}
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.errors = append(rc.errors, err)
}

// GetErrors returns all errors for this request
func (rc *RequestContext) GetErrors() []error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	// Return a copy
	errors := make([]error, len(rc.errors))
	copy(errors, rc.errors)
	return errors
}

// HasErrors returns true if the request has any errors
func (rc *RequestContext) HasErrors() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.errors) > 0
}

// ErrorCount returns the number of errors
func (rc *RequestContext) ErrorCount() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.errors)
}

// ClearErrors clears all errors
func (rc *RequestContext) ClearErrors() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.errors = make([]error, 0)
}

// ============================================================================
// Metadata Management
// ============================================================================

// SetMetadata sets metadata for this request
func (rc *RequestContext) SetMetadata(key string, value interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.metadata[key] = value
}

// GetMetadata retrieves metadata
func (rc *RequestContext) GetMetadata(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	value, exists := rc.metadata[key]
	return value, exists
}

// HasMetadata checks if metadata exists
func (rc *RequestContext) HasMetadata(key string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	_, exists := rc.metadata[key]
	return exists
}

// ============================================================================
// Cleanup
// ============================================================================

// Cleanup releases resources associated with this request
func (rc *RequestContext) Cleanup() {
	rc.cancel()
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Clear all state
	rc.globals = nil
	rc.output = nil
	rc.errors = nil
	rc.metadata = nil
}

// ============================================================================
// Request Manager - Manages multiple concurrent requests
// ============================================================================

// RequestManager manages multiple concurrent request contexts
type RequestManager struct {
	contexts  map[string]*RequestContext
	mu        sync.RWMutex
	maxActive int
	timeout   time.Duration

	// Statistics
	totalRequests     int64
	activeRequests    int64
	completedRequests int64
	timedOutRequests  int64
	statsMu           sync.Mutex
}

// NewRequestManager creates a new request manager
func NewRequestManager(maxActive int, defaultTimeout time.Duration) *RequestManager {
	if maxActive <= 0 {
		maxActive = 1000 // Default max active requests
	}

	return &RequestManager{
		contexts:  make(map[string]*RequestContext),
		maxActive: maxActive,
		timeout:   defaultTimeout,
	}
}

// CreateRequest creates and registers a new request context
func (rm *RequestManager) CreateRequest(id string) (*RequestContext, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Check if at capacity
	if len(rm.contexts) >= rm.maxActive {
		return nil, fmt.Errorf("max active requests reached (%d)", rm.maxActive)
	}

	// Check for duplicate ID
	if _, exists := rm.contexts[id]; exists {
		return nil, fmt.Errorf("request ID already exists: %s", id)
	}

	// Create context
	var ctx *RequestContext
	if rm.timeout > 0 {
		ctx = NewRequestContextWithTimeout(id, rm.timeout)
	} else {
		ctx = NewRequestContext(id)
	}

	rm.contexts[id] = ctx

	// Update stats
	rm.statsMu.Lock()
	rm.totalRequests++
	rm.activeRequests++
	rm.statsMu.Unlock()

	return ctx, nil
}

// GetRequest retrieves a request context by ID
func (rm *RequestManager) GetRequest(id string) (*RequestContext, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	ctx, exists := rm.contexts[id]
	return ctx, exists
}

// CompleteRequest marks a request as complete and cleans it up
func (rm *RequestManager) CompleteRequest(id string) error {
	rm.mu.Lock()
	ctx, exists := rm.contexts[id]
	if !exists {
		rm.mu.Unlock()
		return fmt.Errorf("request not found: %s", id)
	}

	// Check if timed out
	timedOut := ctx.IsTimedOut()

	// Cleanup context
	ctx.Cleanup()

	// Remove from map
	delete(rm.contexts, id)
	rm.mu.Unlock()

	// Update stats
	rm.statsMu.Lock()
	rm.activeRequests--
	rm.completedRequests++
	if timedOut {
		rm.timedOutRequests++
	}
	rm.statsMu.Unlock()

	return nil
}

// CancelRequest cancels a specific request
func (rm *RequestManager) CancelRequest(id string) error {
	rm.mu.RLock()
	ctx, exists := rm.contexts[id]
	rm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("request not found: %s", id)
	}

	ctx.Cancel()
	return nil
}

// CleanupTimedOut removes all timed-out requests
func (rm *RequestManager) CleanupTimedOut() int {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	count := 0
	for id, ctx := range rm.contexts {
		if ctx.IsTimedOut() {
			ctx.Cleanup()
			delete(rm.contexts, id)
			count++

			rm.statsMu.Lock()
			rm.activeRequests--
			rm.completedRequests++
			rm.timedOutRequests++
			rm.statsMu.Unlock()
		}
	}

	return count
}

// ActiveCount returns the number of active requests
func (rm *RequestManager) ActiveCount() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.contexts)
}

// Stats returns manager statistics
func (rm *RequestManager) Stats() RequestManagerStats {
	rm.statsMu.Lock()
	defer rm.statsMu.Unlock()

	return RequestManagerStats{
		TotalRequests:     rm.totalRequests,
		ActiveRequests:    rm.activeRequests,
		CompletedRequests: rm.completedRequests,
		TimedOutRequests:  rm.timedOutRequests,
		MaxActive:         int64(rm.maxActive),
	}
}

// Shutdown cancels all active requests and cleans up
func (rm *RequestManager) Shutdown() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for _, ctx := range rm.contexts {
		ctx.Cleanup()
	}

	rm.contexts = make(map[string]*RequestContext)

	rm.statsMu.Lock()
	rm.activeRequests = 0
	rm.statsMu.Unlock()
}

// ============================================================================
// Statistics Types
// ============================================================================

// RequestManagerStats contains request manager statistics
type RequestManagerStats struct {
	TotalRequests     int64
	ActiveRequests    int64
	CompletedRequests int64
	TimedOutRequests  int64
	MaxActive         int64
}
