package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownHandler is a function that is called during graceful shutdown
// It receives a context that will be cancelled when the shutdown timeout is reached
type ShutdownHandler func(ctx context.Context) error

// ShutdownManager manages graceful shutdown of the application
type ShutdownManager struct {
	mu sync.RWMutex

	// Shutdown handlers registered by priority
	handlers []shutdownHandlerEntry

	// Signal channel for OS signals
	signalChan chan os.Signal

	// Shutdown timeout (default 30 seconds)
	timeout time.Duration

	// Whether shutdown has been initiated
	shutdownInitiated bool

	// Channels for coordination
	shutdownChan chan struct{} // Closed when shutdown is initiated
	doneChan     chan struct{} // Closed when shutdown is complete
}

// shutdownHandlerEntry represents a handler with its priority
type shutdownHandlerEntry struct {
	name     string
	handler  ShutdownHandler
	priority int // Higher priority runs first
}

// Priority constants for shutdown handlers
const (
	ShutdownPriorityHighest = 100 // User-critical cleanup (close DB connections, flush logs)
	ShutdownPriorityHigh    = 75  // Important cleanup (stop accepting requests)
	ShutdownPriorityNormal  = 50  // Standard cleanup (flush metrics, save state)
	ShutdownPriorityLow     = 25  // Best-effort cleanup (optional tasks)
	ShutdownPriorityLowest  = 0   // Final cleanup (close files, cleanup temp resources)
)

// Global shutdown manager
var globalShutdownManager *ShutdownManager
var shutdownManagerOnce sync.Once

// GetShutdownManager returns the global shutdown manager
func GetShutdownManager() *ShutdownManager {
	shutdownManagerOnce.Do(func() {
		globalShutdownManager = NewShutdownManager(30 * time.Second)
	})
	return globalShutdownManager
}

// NewShutdownManager creates a new shutdown manager with the specified timeout
func NewShutdownManager(timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		handlers:     make([]shutdownHandlerEntry, 0),
		signalChan:   make(chan os.Signal, 1),
		timeout:      timeout,
		shutdownChan: make(chan struct{}),
		doneChan:     make(chan struct{}),
	}
}

// RegisterHandler registers a shutdown handler with the specified priority
// Higher priority handlers are called first during shutdown
func (sm *ShutdownManager) RegisterHandler(name string, handler ShutdownHandler, priority int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	entry := shutdownHandlerEntry{
		name:     name,
		handler:  handler,
		priority: priority,
	}

	// Insert in priority order (highest first)
	inserted := false
	for i, h := range sm.handlers {
		if priority > h.priority {
			// Insert at position i
			sm.handlers = append(sm.handlers[:i], append([]shutdownHandlerEntry{entry}, sm.handlers[i:]...)...)
			inserted = true
			break
		}
	}

	if !inserted {
		sm.handlers = append(sm.handlers, entry)
	}
}

// ListenForSignals starts listening for OS signals (SIGTERM, SIGINT)
// When a signal is received, it initiates graceful shutdown
func (sm *ShutdownManager) ListenForSignals() {
	signal.Notify(sm.signalChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sm.signalChan
		Infof("Received signal: %v, initiating graceful shutdown", sig)
		sm.Shutdown()
	}()
}

// Shutdown initiates graceful shutdown of the application
// It calls all registered handlers in priority order with a timeout
func (sm *ShutdownManager) Shutdown() error {
	sm.mu.Lock()
	if sm.shutdownInitiated {
		sm.mu.Unlock()
		return fmt.Errorf("shutdown already initiated")
	}
	sm.shutdownInitiated = true
	close(sm.shutdownChan)
	handlers := make([]shutdownHandlerEntry, len(sm.handlers))
	copy(handlers, sm.handlers)
	sm.mu.Unlock()

	Infof("Starting graceful shutdown (timeout: %v)", sm.timeout)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	// Call each handler in priority order
	var errors []error
	for _, entry := range handlers {
		Debugf("Calling shutdown handler: %s (priority: %d)", entry.name, entry.priority)

		if err := sm.callHandler(ctx, entry); err != nil {
			Errorf("Shutdown handler '%s' failed: %v", entry.name, err)
			errors = append(errors, fmt.Errorf("%s: %w", entry.name, err))
		} else {
			Debugf("Shutdown handler '%s' completed successfully", entry.name)
		}

		// Check if context was cancelled (timeout)
		if ctx.Err() != nil {
			Warn("Shutdown timeout reached, remaining handlers will be skipped")
			break
		}
	}

	// Mark shutdown as complete
	close(sm.doneChan)

	if len(errors) > 0 {
		Infof("Graceful shutdown completed with %d error(s)", len(errors))
		return fmt.Errorf("shutdown completed with errors: %v", errors)
	}

	Info("Graceful shutdown completed successfully")
	return nil
}

// callHandler calls a shutdown handler with a timeout
func (sm *ShutdownManager) callHandler(ctx context.Context, entry shutdownHandlerEntry) error {
	// Create a channel to receive the result
	errChan := make(chan error, 1)

	// Run handler in goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("handler panicked: %v", r)
			}
		}()
		errChan <- entry.handler(ctx)
	}()

	// Wait for handler to complete or context to be cancelled
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("handler timed out")
	}
}

// WaitForShutdown blocks until shutdown is initiated
func (sm *ShutdownManager) WaitForShutdown() {
	<-sm.shutdownChan
}

// WaitForShutdownComplete blocks until shutdown is complete
func (sm *ShutdownManager) WaitForShutdownComplete() {
	<-sm.doneChan
}

// IsShuttingDown returns true if shutdown has been initiated
func (sm *ShutdownManager) IsShuttingDown() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.shutdownInitiated
}

// SetTimeout sets the shutdown timeout
func (sm *ShutdownManager) SetTimeout(timeout time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.timeout = timeout
}

// GetTimeout returns the current shutdown timeout
func (sm *ShutdownManager) GetTimeout() time.Duration {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.timeout
}

// Package-level convenience functions

// RegisterShutdownHandler registers a shutdown handler with the global shutdown manager
func RegisterShutdownHandler(name string, handler ShutdownHandler, priority int) {
	GetShutdownManager().RegisterHandler(name, handler, priority)
}

// ListenForShutdownSignals starts listening for OS signals with the global shutdown manager
func ListenForShutdownSignals() {
	GetShutdownManager().ListenForSignals()
}

// Shutdown initiates graceful shutdown with the global shutdown manager
func Shutdown() error {
	return GetShutdownManager().Shutdown()
}

// WaitForShutdown blocks until shutdown is initiated with the global shutdown manager
func WaitForShutdown() {
	GetShutdownManager().WaitForShutdown()
}

// WaitForShutdownComplete blocks until shutdown is complete with the global shutdown manager
func WaitForShutdownComplete() {
	GetShutdownManager().WaitForShutdownComplete()
}

// IsShuttingDown returns true if shutdown has been initiated with the global shutdown manager
func IsShuttingDown() bool {
	return GetShutdownManager().IsShuttingDown()
}

// SetShutdownTimeout sets the shutdown timeout for the global shutdown manager
func SetShutdownTimeout(timeout time.Duration) {
	GetShutdownManager().SetTimeout(timeout)
}

// GetShutdownTimeout returns the current shutdown timeout from the global shutdown manager
func GetShutdownTimeout() time.Duration {
	return GetShutdownManager().GetTimeout()
}
