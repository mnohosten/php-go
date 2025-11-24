package runtime

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ResourceLimitType represents the type of resource limit
type ResourceLimitType int

const (
	LimitTypeMemory          ResourceLimitType = iota // Memory usage limit
	LimitTypeExecutionTime                            // Execution time limit
	LimitTypeRecursionDepth                           // Maximum recursion depth
	LimitTypeInstructions                             // Maximum instruction count
	LimitTypeOutputSize                               // Maximum output buffer size
)

// String returns the string representation of a resource limit type
func (t ResourceLimitType) String() string {
	switch t {
	case LimitTypeMemory:
		return "memory"
	case LimitTypeExecutionTime:
		return "execution_time"
	case LimitTypeRecursionDepth:
		return "recursion_depth"
	case LimitTypeInstructions:
		return "instructions"
	case LimitTypeOutputSize:
		return "output_size"
	default:
		return "unknown"
	}
}

// ResourceLimitError represents an error when a resource limit is exceeded
type ResourceLimitError struct {
	Type    ResourceLimitType
	Limit   int64
	Current int64
	Message string
}

// Error implements the error interface
func (e *ResourceLimitError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s limit exceeded: current=%d, limit=%d", e.Type, e.Current, e.Limit)
}

// ResourceLimits holds all resource limit configurations
type ResourceLimits struct {
	// Memory limit in bytes (0 = unlimited)
	MemoryLimit int64

	// Execution time limit in nanoseconds (0 = unlimited)
	ExecutionTimeLimit int64

	// Maximum recursion depth (0 = unlimited, default = 1000)
	MaxRecursionDepth int64

	// Maximum instruction count (0 = unlimited)
	MaxInstructions int64

	// Maximum output buffer size in bytes (0 = unlimited)
	MaxOutputSize int64

	// Start time for execution time tracking
	startTime time.Time

	// Instruction counter
	instructionCount atomic.Int64

	// Current recursion depth
	recursionDepth atomic.Int64

	// Current output size
	outputSize atomic.Int64

	// Whether limits are enabled
	enabled bool

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// Mutex for thread safety
	mu sync.RWMutex

	// Callbacks for limit violations
	onMemoryLimitExceeded          func(current, limit int64)
	onExecutionTimeLimitExceeded   func(elapsed, limit time.Duration)
	onRecursionDepthLimitExceeded  func(current, limit int64)
	onInstructionLimitExceeded     func(current, limit int64)
	onOutputSizeLimitExceeded      func(current, limit int64)
}

// ResourceLimiter manages resource limits for PHP execution
type ResourceLimiter struct {
	limits *ResourceLimits
	mu     sync.RWMutex
}

// Global resource limiter instance
var globalLimiter *ResourceLimiter
var limiterOnce sync.Once

// GetResourceLimiter returns the global resource limiter instance
func GetResourceLimiter() *ResourceLimiter {
	limiterOnce.Do(func() {
		globalLimiter = NewResourceLimiter()
	})
	return globalLimiter
}

// NewResourceLimiter creates a new resource limiter with default limits
func NewResourceLimiter() *ResourceLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	return &ResourceLimiter{
		limits: &ResourceLimits{
			MemoryLimit:        0,    // Unlimited by default
			ExecutionTimeLimit: 0,    // Unlimited by default
			MaxRecursionDepth:  1000, // Default PHP limit
			MaxInstructions:    0,    // Unlimited by default
			MaxOutputSize:      0,    // Unlimited by default
			enabled:            true,
			ctx:                ctx,
			cancel:             cancel,
		},
	}
}

// SetMemoryLimit sets the memory limit in bytes
func (rl *ResourceLimiter) SetMemoryLimit(bytes int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.MemoryLimit = bytes
}

// GetMemoryLimit returns the current memory limit
func (rl *ResourceLimiter) GetMemoryLimit() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.MemoryLimit
}

// SetExecutionTimeLimit sets the execution time limit in seconds
func (rl *ResourceLimiter) SetExecutionTimeLimit(seconds int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.ExecutionTimeLimit = seconds * int64(time.Second)
}

// GetExecutionTimeLimit returns the current execution time limit in seconds
func (rl *ResourceLimiter) GetExecutionTimeLimit() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.ExecutionTimeLimit / int64(time.Second)
}

// SetMaxRecursionDepth sets the maximum recursion depth
func (rl *ResourceLimiter) SetMaxRecursionDepth(depth int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.MaxRecursionDepth = depth
}

// GetMaxRecursionDepth returns the current maximum recursion depth
func (rl *ResourceLimiter) GetMaxRecursionDepth() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.MaxRecursionDepth
}

// SetMaxInstructions sets the maximum instruction count
func (rl *ResourceLimiter) SetMaxInstructions(count int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.MaxInstructions = count
}

// GetMaxInstructions returns the current maximum instruction count
func (rl *ResourceLimiter) GetMaxInstructions() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.MaxInstructions
}

// SetMaxOutputSize sets the maximum output buffer size in bytes
func (rl *ResourceLimiter) SetMaxOutputSize(bytes int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.MaxOutputSize = bytes
}

// GetMaxOutputSize returns the current maximum output buffer size
func (rl *ResourceLimiter) GetMaxOutputSize() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.MaxOutputSize
}

// SetEnabled enables or disables limit checking
func (rl *ResourceLimiter) SetEnabled(enabled bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.enabled = enabled
}

// IsEnabled returns whether limit checking is enabled
func (rl *ResourceLimiter) IsEnabled() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limits.enabled
}

// Start starts tracking resource usage (resets counters and starts timer)
func (rl *ResourceLimiter) Start() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.startTime = time.Now()
	rl.limits.instructionCount.Store(0)
	rl.limits.recursionDepth.Store(0)
	rl.limits.outputSize.Store(0)
}

// Stop stops tracking and cancels any active checks
func (rl *ResourceLimiter) Stop() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.limits.cancel != nil {
		rl.limits.cancel()
	}
}

// Reset resets all counters and limits to default values
func (rl *ResourceLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Cancel existing context
	if rl.limits.cancel != nil {
		rl.limits.cancel()
	}

	// Create new context
	ctx, cancel := context.WithCancel(context.Background())

	// Reset to default values
	rl.limits.MemoryLimit = 0
	rl.limits.ExecutionTimeLimit = 0
	rl.limits.MaxRecursionDepth = 1000
	rl.limits.MaxInstructions = 0
	rl.limits.MaxOutputSize = 0
	rl.limits.startTime = time.Time{}
	rl.limits.instructionCount.Store(0)
	rl.limits.recursionDepth.Store(0)
	rl.limits.outputSize.Store(0)
	rl.limits.enabled = true
	rl.limits.ctx = ctx
	rl.limits.cancel = cancel
	rl.limits.onMemoryLimitExceeded = nil
	rl.limits.onExecutionTimeLimitExceeded = nil
	rl.limits.onRecursionDepthLimitExceeded = nil
	rl.limits.onInstructionLimitExceeded = nil
	rl.limits.onOutputSizeLimitExceeded = nil
}

// CheckMemoryLimit checks if memory usage exceeds the limit
func (rl *ResourceLimiter) CheckMemoryLimit() error {
	rl.mu.RLock()
	if !rl.limits.enabled || rl.limits.MemoryLimit <= 0 {
		rl.mu.RUnlock()
		return nil
	}
	limit := rl.limits.MemoryLimit
	callback := rl.limits.onMemoryLimitExceeded
	rl.mu.RUnlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	current := int64(m.Alloc)

	if current > limit {
		if callback != nil {
			callback(current, limit)
		}
		return &ResourceLimitError{
			Type:    LimitTypeMemory,
			Limit:   limit,
			Current: current,
			Message: fmt.Sprintf("memory limit of %d bytes exceeded (current: %d bytes)", limit, current),
		}
	}
	return nil
}

// CheckExecutionTimeLimit checks if execution time exceeds the limit
func (rl *ResourceLimiter) CheckExecutionTimeLimit() error {
	rl.mu.RLock()
	if !rl.limits.enabled || rl.limits.ExecutionTimeLimit <= 0 {
		rl.mu.RUnlock()
		return nil
	}
	limit := rl.limits.ExecutionTimeLimit
	startTime := rl.limits.startTime
	callback := rl.limits.onExecutionTimeLimitExceeded
	rl.mu.RUnlock()

	if startTime.IsZero() {
		return nil
	}

	elapsed := time.Since(startTime)
	if elapsed > time.Duration(limit) {
		if callback != nil {
			callback(elapsed, time.Duration(limit))
		}
		return &ResourceLimitError{
			Type:    LimitTypeExecutionTime,
			Limit:   limit,
			Current: int64(elapsed),
			Message: fmt.Sprintf("execution time limit of %v exceeded (current: %v)", time.Duration(limit), elapsed),
		}
	}
	return nil
}

// IncrementInstructionCount increments the instruction counter and checks the limit
func (rl *ResourceLimiter) IncrementInstructionCount() error {
	rl.mu.RLock()
	if !rl.limits.enabled || rl.limits.MaxInstructions <= 0 {
		rl.mu.RUnlock()
		return nil
	}
	limit := rl.limits.MaxInstructions
	callback := rl.limits.onInstructionLimitExceeded
	rl.mu.RUnlock()

	current := rl.limits.instructionCount.Add(1)
	if current > limit {
		if callback != nil {
			callback(current, limit)
		}
		return &ResourceLimitError{
			Type:    LimitTypeInstructions,
			Limit:   limit,
			Current: current,
			Message: fmt.Sprintf("instruction limit of %d exceeded (current: %d)", limit, current),
		}
	}
	return nil
}

// GetInstructionCount returns the current instruction count
func (rl *ResourceLimiter) GetInstructionCount() int64 {
	return rl.limits.instructionCount.Load()
}

// IncrementRecursionDepth increments the recursion depth counter and checks the limit
func (rl *ResourceLimiter) IncrementRecursionDepth() error {
	rl.mu.RLock()
	if !rl.limits.enabled || rl.limits.MaxRecursionDepth <= 0 {
		rl.mu.RUnlock()
		return nil
	}
	limit := rl.limits.MaxRecursionDepth
	callback := rl.limits.onRecursionDepthLimitExceeded
	rl.mu.RUnlock()

	current := rl.limits.recursionDepth.Add(1)
	if current > limit {
		if callback != nil {
			callback(current, limit)
		}
		return &ResourceLimitError{
			Type:    LimitTypeRecursionDepth,
			Limit:   limit,
			Current: current,
			Message: fmt.Sprintf("maximum recursion depth of %d exceeded", limit),
		}
	}
	return nil
}

// DecrementRecursionDepth decrements the recursion depth counter
func (rl *ResourceLimiter) DecrementRecursionDepth() {
	rl.limits.recursionDepth.Add(-1)
}

// GetRecursionDepth returns the current recursion depth
func (rl *ResourceLimiter) GetRecursionDepth() int64 {
	return rl.limits.recursionDepth.Load()
}

// AddOutputSize adds to the output size counter and checks the limit
func (rl *ResourceLimiter) AddOutputSize(bytes int64) error {
	rl.mu.RLock()
	if !rl.limits.enabled || rl.limits.MaxOutputSize <= 0 {
		rl.mu.RUnlock()
		return nil
	}
	limit := rl.limits.MaxOutputSize
	callback := rl.limits.onOutputSizeLimitExceeded
	rl.mu.RUnlock()

	current := rl.limits.outputSize.Add(bytes)
	if current > limit {
		if callback != nil {
			callback(current, limit)
		}
		return &ResourceLimitError{
			Type:    LimitTypeOutputSize,
			Limit:   limit,
			Current: current,
			Message: fmt.Sprintf("output size limit of %d bytes exceeded (current: %d bytes)", limit, current),
		}
	}
	return nil
}

// GetOutputSize returns the current output size
func (rl *ResourceLimiter) GetOutputSize() int64 {
	return rl.limits.outputSize.Load()
}

// CheckAllLimits checks all configured limits and returns the first error encountered
func (rl *ResourceLimiter) CheckAllLimits() error {
	// Check memory limit
	if err := rl.CheckMemoryLimit(); err != nil {
		return err
	}

	// Check execution time limit
	if err := rl.CheckExecutionTimeLimit(); err != nil {
		return err
	}

	return nil
}

// GetStats returns current resource usage statistics
func (rl *ResourceLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var elapsed time.Duration
	if !rl.limits.startTime.IsZero() {
		elapsed = time.Since(rl.limits.startTime)
	}

	return map[string]interface{}{
		"memory": map[string]interface{}{
			"current": m.Alloc,
			"limit":   rl.limits.MemoryLimit,
			"enabled": rl.limits.MemoryLimit > 0,
		},
		"execution_time": map[string]interface{}{
			"elapsed": elapsed.Seconds(),
			"limit":   time.Duration(rl.limits.ExecutionTimeLimit).Seconds(),
			"enabled": rl.limits.ExecutionTimeLimit > 0,
		},
		"recursion_depth": map[string]interface{}{
			"current": rl.limits.recursionDepth.Load(),
			"limit":   rl.limits.MaxRecursionDepth,
			"enabled": rl.limits.MaxRecursionDepth > 0,
		},
		"instructions": map[string]interface{}{
			"current": rl.limits.instructionCount.Load(),
			"limit":   rl.limits.MaxInstructions,
			"enabled": rl.limits.MaxInstructions > 0,
		},
		"output_size": map[string]interface{}{
			"current": rl.limits.outputSize.Load(),
			"limit":   rl.limits.MaxOutputSize,
			"enabled": rl.limits.MaxOutputSize > 0,
		},
	}
}

// SetOnMemoryLimitExceeded sets a callback for when memory limit is exceeded
func (rl *ResourceLimiter) SetOnMemoryLimitExceeded(callback func(current, limit int64)) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.onMemoryLimitExceeded = callback
}

// SetOnExecutionTimeLimitExceeded sets a callback for when execution time limit is exceeded
func (rl *ResourceLimiter) SetOnExecutionTimeLimitExceeded(callback func(elapsed, limit time.Duration)) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.onExecutionTimeLimitExceeded = callback
}

// SetOnRecursionDepthLimitExceeded sets a callback for when recursion depth limit is exceeded
func (rl *ResourceLimiter) SetOnRecursionDepthLimitExceeded(callback func(current, limit int64)) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.onRecursionDepthLimitExceeded = callback
}

// SetOnInstructionLimitExceeded sets a callback for when instruction limit is exceeded
func (rl *ResourceLimiter) SetOnInstructionLimitExceeded(callback func(current, limit int64)) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.onInstructionLimitExceeded = callback
}

// SetOnOutputSizeLimitExceeded sets a callback for when output size limit is exceeded
func (rl *ResourceLimiter) SetOnOutputSizeLimitExceeded(callback func(current, limit int64)) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits.onOutputSizeLimitExceeded = callback
}

// Package-level convenience functions using the global limiter

// SetMemoryLimit sets the memory limit in bytes
func SetMemoryLimit(bytes int64) {
	GetResourceLimiter().SetMemoryLimit(bytes)
}

// SetExecutionTimeLimit sets the execution time limit in seconds
func SetExecutionTimeLimit(seconds int64) {
	GetResourceLimiter().SetExecutionTimeLimit(seconds)
}

// SetMaxRecursionDepth sets the maximum recursion depth
func SetMaxRecursionDepth(depth int64) {
	GetResourceLimiter().SetMaxRecursionDepth(depth)
}

// SetMaxInstructions sets the maximum instruction count
func SetMaxInstructions(count int64) {
	GetResourceLimiter().SetMaxInstructions(count)
}

// SetMaxOutputSize sets the maximum output buffer size in bytes
func SetMaxOutputSize(bytes int64) {
	GetResourceLimiter().SetMaxOutputSize(bytes)
}

// StartResourceTracking starts tracking resource usage
func StartResourceTracking() {
	GetResourceLimiter().Start()
}

// StopResourceTracking stops tracking resource usage
func StopResourceTracking() {
	GetResourceLimiter().Stop()
}

// CheckResourceLimits checks all configured resource limits
func CheckResourceLimits() error {
	return GetResourceLimiter().CheckAllLimits()
}

// GetResourceStats returns current resource usage statistics
func GetResourceStats() map[string]interface{} {
	return GetResourceLimiter().GetStats()
}
