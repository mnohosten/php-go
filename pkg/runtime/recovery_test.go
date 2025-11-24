package runtime

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDefaultRecoveryConfig(t *testing.T) {
	config := DefaultRecoveryConfig()
	if config == nil {
		t.Fatal("DefaultRecoveryConfig returned nil")
	}
	if config.Strategy != RecoveryStrategyLog {
		t.Errorf("expected strategy Log, got %v", config.Strategy)
	}
	if config.MaxRetries != 3 {
		t.Errorf("expected max retries 3, got %d", config.MaxRetries)
	}
	if config.RetryDelay != 100*time.Millisecond {
		t.Errorf("expected retry delay 100ms, got %v", config.RetryDelay)
	}
	if config.OnPanic == nil {
		t.Error("expected OnPanic callback to be set")
	}
	if config.ShouldRecover == nil {
		t.Error("expected ShouldRecover callback to be set")
	}
}

func TestNewRecoveryManager(t *testing.T) {
	rm := NewRecoveryManager(nil)
	if rm == nil {
		t.Fatal("NewRecoveryManager returned nil")
	}
	if rm.config == nil {
		t.Error("expected config to be initialized")
	}
	if rm.panicHistory == nil {
		t.Error("expected panic history to be initialized")
	}
	if rm.maxHistory != 100 {
		t.Errorf("expected max history 100, got %d", rm.maxHistory)
	}
}

func TestRecoveryManagerSetGetConfig(t *testing.T) {
	rm := NewRecoveryManager(nil)
	config := &RecoveryConfig{
		Strategy:   RecoveryStrategyRestart,
		MaxRetries: 5,
		RetryDelay: 200 * time.Millisecond,
	}
	rm.SetConfig(config)

	retrieved := rm.GetConfig()
	if retrieved.Strategy != RecoveryStrategyRestart {
		t.Errorf("expected strategy Restart, got %v", retrieved.Strategy)
	}
	if retrieved.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", retrieved.MaxRetries)
	}
}

func TestRecoverNoPanic(t *testing.T) {
	rm := NewRecoveryManager(nil)
	called := false
	err := rm.Recover(func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("function was not called")
	}
	if rm.GetPanicCount() != 0 {
		t.Errorf("expected panic count 0, got %d", rm.GetPanicCount())
	}
}

func TestRecoverWithPanicLog(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	err := rm.Recover(func() error {
		panic("test panic")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "test panic") {
		t.Errorf("expected error to contain 'test panic', got %v", err)
	}
	if rm.GetPanicCount() != 1 {
		t.Errorf("expected panic count 1, got %d", rm.GetPanicCount())
	}
}

func TestRecoverWithPanicNone(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy: RecoveryStrategyNone,
	})

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic to propagate")
		}
	}()

	_ = rm.Recover(func() error {
		panic("test panic")
	})
}

func TestRecoverWithCustomHandler(t *testing.T) {
	customCalled := false
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy: RecoveryStrategyCustom,
		CustomHandler: func(err interface{}, stack []byte) error {
			customCalled = true
			return errors.New("custom error")
		},
		ShouldRecover: func(err interface{}) bool { return true },
	})

	err := rm.Recover(func() error {
		panic("test panic")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !customCalled {
		t.Error("custom handler was not called")
	}
	if err.Error() != "custom error" {
		t.Errorf("expected 'custom error', got %v", err)
	}
}

func TestRecoverWithOnPanic(t *testing.T) {
	onPanicCalled := false
	var capturedErr interface{}
	var capturedStack []byte

	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy: RecoveryStrategyLog,
		OnPanic: func(err interface{}, stack []byte) {
			onPanicCalled = true
			capturedErr = err
			capturedStack = stack
		},
		ShouldRecover: func(err interface{}) bool { return true },
	})

	_ = rm.Recover(func() error {
		panic("test panic")
	})

	if !onPanicCalled {
		t.Error("OnPanic callback was not called")
	}
	if capturedErr == nil {
		t.Error("expected error to be captured")
	}
	if len(capturedStack) == 0 {
		t.Error("expected stack trace to be captured")
	}
}

func TestRecoverWithShouldRecover(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy: RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool {
			if s, ok := err.(string); ok {
				return s != "fatal"
			}
			return true
		},
	})

	// Test recoverable panic
	err := rm.Recover(func() error {
		panic("recoverable")
	})
	if err == nil {
		t.Error("expected error, got nil")
	}

	// Test non-recoverable panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic to propagate")
		} else if r != "fatal" {
			t.Errorf("expected 'fatal' panic, got %v", r)
		}
	}()
	_ = rm.Recover(func() error {
		panic("fatal")
	})
}

func TestRecoverWithRetry(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyRestart,
		MaxRetries:    3,
		RetryDelay:    10 * time.Millisecond,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	attempts := 0
	err := rm.RecoverWithRetry(func() error {
		attempts++
		if attempts < 3 {
			panic("retry me")
		}
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRecoverWithRetryMaxExceeded(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyRestart,
		MaxRetries:    2,
		RetryDelay:    10 * time.Millisecond,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	attempts := 0
	err := rm.RecoverWithRetry(func() error {
		attempts++
		panic("always fail")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "max retries") {
		t.Errorf("expected 'max retries' error, got %v", err)
	}
	if attempts != 3 { // Initial + 2 retries
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRecoverWithRetryNonRestartError(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		MaxRetries:    3,
		RetryDelay:    10 * time.Millisecond,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	attempts := 0
	err := rm.RecoverWithRetry(func() error {
		attempts++
		panic("fail once")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 1 { // Should not retry for non-restart errors
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestPanicHistory(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	// Generate some panics
	for i := 0; i < 5; i++ {
		_ = rm.Recover(func() error {
			panic(fmt.Sprintf("panic %d", i))
		})
	}

	if rm.GetPanicCount() != 5 {
		t.Errorf("expected panic count 5, got %d", rm.GetPanicCount())
	}

	history := rm.GetPanicHistory()
	if len(history) != 5 {
		t.Errorf("expected history length 5, got %d", len(history))
	}

	for i, record := range history {
		if !record.Recovered {
			t.Errorf("record %d should be marked as recovered", i)
		}
		if record.Error == nil {
			t.Errorf("record %d should have an error", i)
		}
		if len(record.Stack) == 0 {
			t.Errorf("record %d should have a stack trace", i)
		}
		if record.Timestamp.IsZero() {
			t.Errorf("record %d should have a timestamp", i)
		}
	}

	lastPanic := rm.GetLastPanic()
	if lastPanic.IsZero() {
		t.Error("expected last panic time to be set")
	}
}

func TestPanicHistoryLimit(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})
	rm.maxHistory = 10

	// Generate more panics than the limit
	for i := 0; i < 15; i++ {
		_ = rm.Recover(func() error {
			panic(fmt.Sprintf("panic %d", i))
		})
	}

	history := rm.GetPanicHistory()
	if len(history) != 10 {
		t.Errorf("expected history length 10 (limited), got %d", len(history))
	}
	if rm.GetPanicCount() != 15 {
		t.Errorf("expected panic count 15, got %d", rm.GetPanicCount())
	}
}

func TestRecoveryReset(t *testing.T) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	// Generate some panics
	_ = rm.Recover(func() error {
		panic("test")
	})

	if rm.GetPanicCount() == 0 {
		t.Fatal("expected non-zero panic count before reset")
	}

	rm.Reset()

	if rm.GetPanicCount() != 0 {
		t.Errorf("expected panic count 0 after reset, got %d", rm.GetPanicCount())
	}
	if !rm.GetLastPanic().IsZero() {
		t.Error("expected last panic time to be zero after reset")
	}
	if len(rm.GetPanicHistory()) != 0 {
		t.Errorf("expected empty history after reset, got %d", len(rm.GetPanicHistory()))
	}
}

func TestErrorHandlers(t *testing.T) {
	rm := NewRecoveryManager(nil)

	called1 := false
	called2 := false

	var handler1 RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		called1 = true
		return false // Don't handle
	}

	var handler2 RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		called2 = true
		return true // Handle it
	}

	rm.PushErrorHandler(handler1)
	rm.PushErrorHandler(handler2)

	handled := rm.HandleError(E_WARNING, "test warning", "test.php", 10)

	if !handled {
		t.Error("expected error to be handled")
	}
	if !called2 {
		t.Error("handler2 should be called")
	}
	if called1 {
		t.Error("handler1 should not be called (handler2 returned true)")
	}
}

func TestErrorHandlersPop(t *testing.T) {
	rm := NewRecoveryManager(nil)

	called := false
	var handler RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		called = true
		return true
	}

	rm.PushErrorHandler(handler)
	rm.PopErrorHandler()

	handled := rm.HandleError(E_WARNING, "test warning", "test.php", 10)

	if handled {
		t.Error("expected error not to be handled after pop")
	}
	if called {
		t.Error("handler should not be called after pop")
	}
}

func TestErrorHandlersOrder(t *testing.T) {
	rm := NewRecoveryManager(nil)

	var order []int

	var handler1 RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		order = append(order, 1)
		return false
	}

	var handler2 RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		order = append(order, 2)
		return false
	}

	var handler3 RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		order = append(order, 3)
		return false
	}

	rm.PushErrorHandler(handler1)
	rm.PushErrorHandler(handler2)
	rm.PushErrorHandler(handler3)

	rm.HandleError(E_WARNING, "test", "test.php", 10)

	// Handlers should be called in reverse order (most recent first)
	if len(order) != 3 {
		t.Fatalf("expected 3 handlers called, got %d", len(order))
	}
	if order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Errorf("expected order [3, 2, 1], got %v", order)
	}
}

func TestExceptionHandler(t *testing.T) {
	rm := NewRecoveryManager(nil)

	called := false
	var capturedEx interface{}

	handler := func(exception interface{}) bool {
		called = true
		capturedEx = exception
		return true
	}

	rm.SetExceptionHandler(handler)

	handled := rm.HandleException("test exception")

	if !handled {
		t.Error("expected exception to be handled")
	}
	if !called {
		t.Error("exception handler should be called")
	}
	if capturedEx != "test exception" {
		t.Errorf("expected exception 'test exception', got %v", capturedEx)
	}
}

func TestExceptionHandlerNil(t *testing.T) {
	rm := NewRecoveryManager(nil)

	handled := rm.HandleException("test exception")

	if handled {
		t.Error("expected exception not to be handled (no handler set)")
	}
}

func TestRestartError(t *testing.T) {
	err := &RestartError{
		Cause: "test cause",
		Stack: []byte("test stack"),
	}

	if !strings.Contains(err.Error(), "test cause") {
		t.Errorf("expected error to contain 'test cause', got %v", err.Error())
	}
}

func TestGlobalRecoveryManager(t *testing.T) {
	// Reset global manager for clean test
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	rm := GetRecoveryManager()
	if rm == nil {
		t.Fatal("GetRecoveryManager returned nil")
	}

	rm2 := GetRecoveryManager()
	if rm != rm2 {
		t.Error("GetRecoveryManager should return same instance")
	}
}

func TestPackageLevelRecover(t *testing.T) {
	// Reset global manager
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	SetRecoveryConfig(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	err := Recover(func() error {
		panic("test panic")
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if GetPanicCount() != 1 {
		t.Errorf("expected panic count 1, got %d", GetPanicCount())
	}
}

func TestPackageLevelRecoverWithRetry(t *testing.T) {
	// Reset global manager
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	SetRecoveryConfig(&RecoveryConfig{
		Strategy:      RecoveryStrategyRestart,
		MaxRetries:    2,
		RetryDelay:    10 * time.Millisecond,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	attempts := 0
	err := RecoverWithRetry(func() error {
		attempts++
		if attempts < 2 {
			panic("retry")
		}
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestPackageLevelErrorHandlers(t *testing.T) {
	// Reset global manager
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	called := false
	var handler RecoveryErrorHandler = func(errType ErrorType, message string, file string, line int) bool {
		called = true
		return true
	}

	PushErrorHandler(handler)
	handled := HandleError(E_WARNING, "test", "test.php", 10)

	if !handled {
		t.Error("expected error to be handled")
	}
	if !called {
		t.Error("handler should be called")
	}

	PopErrorHandler()
	called = false
	handled = HandleError(E_WARNING, "test", "test.php", 10)

	if handled {
		t.Error("expected error not to be handled after pop")
	}
	if called {
		t.Error("handler should not be called after pop")
	}
}

func TestPackageLevelExceptionHandler(t *testing.T) {
	// Reset global manager
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	called := false
	handler := func(exception interface{}) bool {
		called = true
		return true
	}

	SetExceptionHandler(handler)
	handled := HandleException("test exception")

	if !handled {
		t.Error("expected exception to be handled")
	}
	if !called {
		t.Error("handler should be called")
	}
}

func TestPackageLevelReset(t *testing.T) {
	// Reset global manager
	globalRecoveryManager = nil
	globalRecoveryManagerOnce = sync.Once{}

	SetRecoveryConfig(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
	})

	_ = Recover(func() error {
		panic("test")
	})

	if GetPanicCount() == 0 {
		t.Fatal("expected non-zero panic count before reset")
	}

	history := GetPanicHistory()
	if len(history) == 0 {
		t.Fatal("expected non-empty panic history before reset")
	}

	ResetRecoveryStats()

	if GetPanicCount() != 0 {
		t.Errorf("expected panic count 0 after reset, got %d", GetPanicCount())
	}

	history = GetPanicHistory()
	if len(history) != 0 {
		t.Errorf("expected empty panic history after reset, got %d", len(history))
	}
}

// Benchmarks

func BenchmarkRecoverNoPanic(b *testing.B) {
	rm := NewRecoveryManager(DefaultRecoveryConfig())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rm.Recover(func() error {
			return nil
		})
	}
}

func BenchmarkRecoverWithPanic(b *testing.B) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyLog,
		ShouldRecover: func(err interface{}) bool { return true },
		OnPanic:       func(err interface{}, stack []byte) {}, // No-op
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rm.Recover(func() error {
			panic("test")
		})
	}
}

func BenchmarkRecoverWithRetry(b *testing.B) {
	rm := NewRecoveryManager(&RecoveryConfig{
		Strategy:      RecoveryStrategyRestart,
		MaxRetries:    3,
		RetryDelay:    time.Microsecond,
		ShouldRecover: func(err interface{}) bool { return true },
		OnPanic:       func(err interface{}, stack []byte) {}, // No-op
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempts := 0
		_ = rm.RecoverWithRetry(func() error {
			attempts++
			if attempts < 2 {
				panic("retry")
			}
			return nil
		})
	}
}

func BenchmarkErrorHandler(b *testing.B) {
	rm := NewRecoveryManager(nil)
	rm.PushErrorHandler(func(errType ErrorType, message string, file string, line int) bool {
		return true
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.HandleError(E_WARNING, "test", "test.php", 10)
	}
}

func BenchmarkExceptionHandler(b *testing.B) {
	rm := NewRecoveryManager(nil)
	rm.SetExceptionHandler(func(exception interface{}) bool {
		return true
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.HandleException("test exception")
	}
}
