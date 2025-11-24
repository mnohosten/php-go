package runtime

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// RecoveryHandler is a function that handles recovered panics
type RecoveryHandler func(err interface{}, stack []byte) error

// RecoveryStrategy defines how to handle errors
type RecoveryStrategy int

const (
	// RecoveryStrategyNone does not recover from panics
	RecoveryStrategyNone RecoveryStrategy = iota
	// RecoveryStrategyLog logs the panic and continues
	RecoveryStrategyLog
	// RecoveryStrategyRestart attempts to restart the failed operation
	RecoveryStrategyRestart
	// RecoveryStrategyCustom uses a custom recovery handler
	RecoveryStrategyCustom
)

// RecoveryConfig configures error recovery behavior
type RecoveryConfig struct {
	// Strategy defines the recovery strategy
	Strategy RecoveryStrategy
	// MaxRetries defines maximum retry attempts for RecoveryStrategyRestart
	MaxRetries int
	// RetryDelay defines delay between retry attempts
	RetryDelay time.Duration
	// CustomHandler is used when Strategy is RecoveryStrategyCustom
	CustomHandler RecoveryHandler
	// OnPanic is called whenever a panic is recovered
	OnPanic func(err interface{}, stack []byte)
	// ShouldRecover determines if a panic should be recovered
	ShouldRecover func(err interface{}) bool
}

// DefaultRecoveryConfig returns a default recovery configuration
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		Strategy:   RecoveryStrategyLog,
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
		OnPanic: func(err interface{}, stack []byte) {
			GetGlobalLogger().WithFields(map[string]interface{}{
				"error": fmt.Sprintf("%v", err),
				"stack": string(stack),
			}).Error("Panic recovered")
		},
		ShouldRecover: func(err interface{}) bool {
			return true // Recover from all panics by default
		},
	}
}

// RecoveryManager manages error recovery
type RecoveryManager struct {
	mu            sync.RWMutex
	config        *RecoveryConfig
	panicCount    int
	lastPanic     time.Time
	panicHistory  []PanicRecord
	maxHistory    int
	errorHandlers []RecoveryErrorHandler
	exceptionHandler ExceptionHandler
}

// PanicRecord stores information about a recovered panic
type PanicRecord struct {
	Timestamp time.Time
	Error     interface{}
	Stack     []byte
	Recovered bool
}

// RecoveryErrorHandler is a custom error handler function
type RecoveryErrorHandler func(errType ErrorType, message string, file string, line int) bool

// ExceptionHandler is a custom exception handler function
type ExceptionHandler func(exception interface{}) bool

var (
	globalRecoveryManager     *RecoveryManager
	globalRecoveryManagerOnce sync.Once
)

// GetRecoveryManager returns the global recovery manager
func GetRecoveryManager() *RecoveryManager {
	globalRecoveryManagerOnce.Do(func() {
		globalRecoveryManager = NewRecoveryManager(DefaultRecoveryConfig())
	})
	return globalRecoveryManager
}

// NewRecoveryManager creates a new recovery manager
func NewRecoveryManager(config *RecoveryConfig) *RecoveryManager {
	if config == nil {
		config = DefaultRecoveryConfig()
	}
	return &RecoveryManager{
		config:       config,
		panicHistory: make([]PanicRecord, 0),
		maxHistory:   100,
		errorHandlers: make([]RecoveryErrorHandler, 0),
	}
}

// SetConfig updates the recovery configuration
func (rm *RecoveryManager) SetConfig(config *RecoveryConfig) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.config = config
}

// GetConfig returns the current recovery configuration
func (rm *RecoveryManager) GetConfig() *RecoveryConfig {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.config
}

// Recover wraps a function with panic recovery
func (rm *RecoveryManager) Recover(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			rm.recordPanic(r, stack, true)

			rm.mu.RLock()
			config := rm.config
			rm.mu.RUnlock()

			// Check if we should recover
			if config.ShouldRecover != nil && !config.ShouldRecover(r) {
				panic(r) // Re-panic if we shouldn't recover
			}

			// Call OnPanic callback
			if config.OnPanic != nil {
				config.OnPanic(r, stack)
			}

			// Apply recovery strategy
			switch config.Strategy {
			case RecoveryStrategyNone:
				panic(r) // Re-panic
			case RecoveryStrategyLog:
				err = fmt.Errorf("panic recovered: %v", r)
			case RecoveryStrategyCustom:
				if config.CustomHandler != nil {
					err = config.CustomHandler(r, stack)
				} else {
					err = fmt.Errorf("panic recovered: %v", r)
				}
			case RecoveryStrategyRestart:
				// For restart strategy, return a special error
				err = &RestartError{Cause: r, Stack: stack}
			default:
				err = fmt.Errorf("panic recovered: %v", r)
			}
		}
	}()

	return fn()
}

// RecoverWithRetry wraps a function with panic recovery and retry logic
func (rm *RecoveryManager) RecoverWithRetry(fn func() error) error {
	rm.mu.RLock()
	maxRetries := rm.config.MaxRetries
	retryDelay := rm.config.RetryDelay
	rm.mu.RUnlock()

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		lastErr = rm.Recover(fn)
		if lastErr == nil {
			return nil // Success
		}

		// If it's not a restart error, don't retry
		if _, ok := lastErr.(*RestartError); !ok {
			return lastErr
		}
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}

// recordPanic records a panic in the history
func (rm *RecoveryManager) recordPanic(err interface{}, stack []byte, recovered bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.panicCount++
	rm.lastPanic = time.Now()

	record := PanicRecord{
		Timestamp: time.Now(),
		Error:     err,
		Stack:     stack,
		Recovered: recovered,
	}

	rm.panicHistory = append(rm.panicHistory, record)
	if len(rm.panicHistory) > rm.maxHistory {
		rm.panicHistory = rm.panicHistory[1:]
	}
}

// GetPanicCount returns the total number of panics recovered
func (rm *RecoveryManager) GetPanicCount() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.panicCount
}

// GetLastPanic returns the time of the last panic
func (rm *RecoveryManager) GetLastPanic() time.Time {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.lastPanic
}

// GetPanicHistory returns the panic history
func (rm *RecoveryManager) GetPanicHistory() []PanicRecord {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	history := make([]PanicRecord, len(rm.panicHistory))
	copy(history, rm.panicHistory)
	return history
}

// Reset resets the panic statistics
func (rm *RecoveryManager) Reset() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.panicCount = 0
	rm.lastPanic = time.Time{}
	rm.panicHistory = make([]PanicRecord, 0)
}

// PushErrorHandler adds a custom error handler to the stack
func (rm *RecoveryManager) PushErrorHandler(handler RecoveryErrorHandler) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.errorHandlers = append(rm.errorHandlers, handler)
}

// PopErrorHandler removes the most recent error handler
func (rm *RecoveryManager) PopErrorHandler() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if len(rm.errorHandlers) > 0 {
		rm.errorHandlers = rm.errorHandlers[:len(rm.errorHandlers)-1]
	}
}

// HandleError calls the registered error handlers
func (rm *RecoveryManager) HandleError(errType ErrorType, message string, file string, line int) bool {
	rm.mu.RLock()
	handlers := make([]RecoveryErrorHandler, len(rm.errorHandlers))
	copy(handlers, rm.errorHandlers)
	rm.mu.RUnlock()

	// Call handlers in reverse order (most recent first)
	for i := len(handlers) - 1; i >= 0; i-- {
		if handlers[i](errType, message, file, line) {
			return true // Handler handled the error
		}
	}

	return false // No handler handled the error
}

// SetExceptionHandler sets the exception handler
func (rm *RecoveryManager) SetExceptionHandler(handler ExceptionHandler) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.exceptionHandler = handler
}

// HandleException calls the registered exception handler
func (rm *RecoveryManager) HandleException(exception interface{}) bool {
	rm.mu.RLock()
	handler := rm.exceptionHandler
	rm.mu.RUnlock()

	if handler != nil {
		return handler(exception)
	}
	return false
}

// RestartError indicates an operation should be restarted
type RestartError struct {
	Cause interface{}
	Stack []byte
}

func (e *RestartError) Error() string {
	return fmt.Sprintf("restart error: %v", e.Cause)
}

// Package-level functions for convenience

// SetRecoveryConfig sets the global recovery configuration
func SetRecoveryConfig(config *RecoveryConfig) {
	GetRecoveryManager().SetConfig(config)
}

// Recover wraps a function with panic recovery using the global manager
func Recover(fn func() error) error {
	return GetRecoveryManager().Recover(fn)
}

// RecoverWithRetry wraps a function with panic recovery and retry using the global manager
func RecoverWithRetry(fn func() error) error {
	return GetRecoveryManager().RecoverWithRetry(fn)
}

// GetPanicCount returns the global panic count
func GetPanicCount() int {
	return GetRecoveryManager().GetPanicCount()
}

// GetPanicHistory returns the global panic history
func GetPanicHistory() []PanicRecord {
	return GetRecoveryManager().GetPanicHistory()
}

// ResetRecoveryStats resets the global recovery statistics
func ResetRecoveryStats() {
	GetRecoveryManager().Reset()
}

// PushErrorHandler adds an error handler to the global manager
func PushErrorHandler(handler RecoveryErrorHandler) {
	GetRecoveryManager().PushErrorHandler(handler)
}

// PopErrorHandler removes the most recent error handler from the global manager
func PopErrorHandler() {
	GetRecoveryManager().PopErrorHandler()
}

// HandleError handles an error using the global manager
func HandleError(errType ErrorType, message string, file string, line int) bool {
	return GetRecoveryManager().HandleError(errType, message, file, line)
}

// SetExceptionHandler sets the global exception handler
func SetExceptionHandler(handler ExceptionHandler) {
	GetRecoveryManager().SetExceptionHandler(handler)
}

// HandleException handles an exception using the global manager
func HandleException(exception interface{}) bool {
	return GetRecoveryManager().HandleException(exception)
}
