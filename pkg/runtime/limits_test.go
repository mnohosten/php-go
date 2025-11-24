package runtime

import (
	"sync"
	"testing"
	"time"
)

func TestResourceLimiter_NewResourceLimiter(t *testing.T) {
	limiter := NewResourceLimiter()
	if limiter == nil {
		t.Fatal("expected non-nil limiter")
	}
	if limiter.limits == nil {
		t.Fatal("expected non-nil limits")
	}
	if !limiter.limits.enabled {
		t.Error("expected limits to be enabled by default")
	}
	if limiter.limits.MaxRecursionDepth != 1000 {
		t.Errorf("expected default recursion depth of 1000, got %d", limiter.limits.MaxRecursionDepth)
	}
}

func TestResourceLimiter_GetResourceLimiter(t *testing.T) {
	limiter1 := GetResourceLimiter()
	limiter2 := GetResourceLimiter()
	if limiter1 != limiter2 {
		t.Error("expected GetResourceLimiter to return the same instance")
	}
}

func TestResourceLimiter_SetMemoryLimit(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test setting memory limit
	limiter.SetMemoryLimit(1024 * 1024) // 1MB
	if got := limiter.GetMemoryLimit(); got != 1024*1024 {
		t.Errorf("expected memory limit of 1048576, got %d", got)
	}

	// Test setting to 0 (unlimited)
	limiter.SetMemoryLimit(0)
	if got := limiter.GetMemoryLimit(); got != 0 {
		t.Errorf("expected memory limit of 0, got %d", got)
	}
}

func TestResourceLimiter_SetExecutionTimeLimit(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test setting execution time limit
	limiter.SetExecutionTimeLimit(30) // 30 seconds
	if got := limiter.GetExecutionTimeLimit(); got != 30 {
		t.Errorf("expected execution time limit of 30 seconds, got %d", got)
	}

	// Test setting to 0 (unlimited)
	limiter.SetExecutionTimeLimit(0)
	if got := limiter.GetExecutionTimeLimit(); got != 0 {
		t.Errorf("expected execution time limit of 0, got %d", got)
	}
}

func TestResourceLimiter_SetMaxRecursionDepth(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test setting recursion depth
	limiter.SetMaxRecursionDepth(500)
	if got := limiter.GetMaxRecursionDepth(); got != 500 {
		t.Errorf("expected max recursion depth of 500, got %d", got)
	}

	// Test setting to 0 (unlimited)
	limiter.SetMaxRecursionDepth(0)
	if got := limiter.GetMaxRecursionDepth(); got != 0 {
		t.Errorf("expected max recursion depth of 0, got %d", got)
	}
}

func TestResourceLimiter_SetMaxInstructions(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test setting instruction limit
	limiter.SetMaxInstructions(10000)
	if got := limiter.GetMaxInstructions(); got != 10000 {
		t.Errorf("expected max instructions of 10000, got %d", got)
	}

	// Test setting to 0 (unlimited)
	limiter.SetMaxInstructions(0)
	if got := limiter.GetMaxInstructions(); got != 0 {
		t.Errorf("expected max instructions of 0, got %d", got)
	}
}

func TestResourceLimiter_SetMaxOutputSize(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test setting output size limit
	limiter.SetMaxOutputSize(1024 * 1024) // 1MB
	if got := limiter.GetMaxOutputSize(); got != 1024*1024 {
		t.Errorf("expected max output size of 1048576, got %d", got)
	}

	// Test setting to 0 (unlimited)
	limiter.SetMaxOutputSize(0)
	if got := limiter.GetMaxOutputSize(); got != 0 {
		t.Errorf("expected max output size of 0, got %d", got)
	}
}

func TestResourceLimiter_SetEnabled(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test disabling limits
	limiter.SetEnabled(false)
	if limiter.IsEnabled() {
		t.Error("expected limits to be disabled")
	}

	// Test enabling limits
	limiter.SetEnabled(true)
	if !limiter.IsEnabled() {
		t.Error("expected limits to be enabled")
	}
}

func TestResourceLimiter_StartAndStop(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test starting resource tracking
	limiter.Start()
	if limiter.limits.startTime.IsZero() {
		t.Error("expected start time to be set")
	}

	startTime := limiter.limits.startTime

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Test that time is tracked
	if time.Since(startTime) < 10*time.Millisecond {
		t.Error("expected time to have elapsed")
	}

	// Test stopping
	limiter.Stop()
	// Note: Stop cancels the context, but doesn't reset the start time
}

func TestResourceLimiter_Reset(t *testing.T) {
	limiter := NewResourceLimiter()

	// Set some values
	limiter.SetMemoryLimit(1024)
	limiter.SetExecutionTimeLimit(30)
	limiter.SetMaxRecursionDepth(100)
	limiter.SetMaxInstructions(5000)
	limiter.SetMaxOutputSize(2048)
	limiter.Start()
	limiter.limits.instructionCount.Add(100)
	limiter.limits.recursionDepth.Add(5)
	limiter.limits.outputSize.Add(512)

	// Reset
	limiter.Reset()

	// Check all values are reset to defaults
	if limiter.GetMemoryLimit() != 0 {
		t.Errorf("expected memory limit to be reset to 0, got %d", limiter.GetMemoryLimit())
	}
	if limiter.GetExecutionTimeLimit() != 0 {
		t.Errorf("expected execution time limit to be reset to 0, got %d", limiter.GetExecutionTimeLimit())
	}
	if limiter.GetMaxRecursionDepth() != 1000 {
		t.Errorf("expected max recursion depth to be reset to 1000, got %d", limiter.GetMaxRecursionDepth())
	}
	if limiter.GetMaxInstructions() != 0 {
		t.Errorf("expected max instructions to be reset to 0, got %d", limiter.GetMaxInstructions())
	}
	if limiter.GetMaxOutputSize() != 0 {
		t.Errorf("expected max output size to be reset to 0, got %d", limiter.GetMaxOutputSize())
	}
	if limiter.GetInstructionCount() != 0 {
		t.Errorf("expected instruction count to be reset to 0, got %d", limiter.GetInstructionCount())
	}
	if limiter.GetRecursionDepth() != 0 {
		t.Errorf("expected recursion depth to be reset to 0, got %d", limiter.GetRecursionDepth())
	}
	if limiter.GetOutputSize() != 0 {
		t.Errorf("expected output size to be reset to 0, got %d", limiter.GetOutputSize())
	}
	if !limiter.IsEnabled() {
		t.Error("expected limits to be enabled after reset")
	}
}

func TestResourceLimiter_CheckMemoryLimit_NoLimit(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(0) // Unlimited

	// Should not error
	if err := limiter.CheckMemoryLimit(); err != nil {
		t.Errorf("expected no error for unlimited memory, got %v", err)
	}
}

func TestResourceLimiter_CheckMemoryLimit_WithLimit(t *testing.T) {
	limiter := NewResourceLimiter()

	// Set an impossibly low limit that will always fail
	limiter.SetMemoryLimit(1) // 1 byte

	// Should error
	err := limiter.CheckMemoryLimit()
	if err == nil {
		t.Error("expected error for exceeded memory limit")
	}

	// Check error type
	if _, ok := err.(*ResourceLimitError); !ok {
		t.Errorf("expected ResourceLimitError, got %T", err)
	}

	// Check error details
	limitErr := err.(*ResourceLimitError)
	if limitErr.Type != LimitTypeMemory {
		t.Errorf("expected error type LimitTypeMemory, got %v", limitErr.Type)
	}
}

func TestResourceLimiter_CheckMemoryLimit_Disabled(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(1) // Set low limit
	limiter.SetEnabled(false)   // But disable checking

	// Should not error because checking is disabled
	if err := limiter.CheckMemoryLimit(); err != nil {
		t.Errorf("expected no error when limits are disabled, got %v", err)
	}
}

func TestResourceLimiter_CheckExecutionTimeLimit_NoLimit(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetExecutionTimeLimit(0) // Unlimited
	limiter.Start()

	// Should not error
	if err := limiter.CheckExecutionTimeLimit(); err != nil {
		t.Errorf("expected no error for unlimited execution time, got %v", err)
	}
}

func TestResourceLimiter_CheckExecutionTimeLimit_WithLimit(t *testing.T) {
	limiter := NewResourceLimiter()

	// Set a very short limit
	limiter.limits.ExecutionTimeLimit = 1 // 1 nanosecond
	limiter.Start()

	// Wait a bit to exceed the limit
	time.Sleep(10 * time.Millisecond)

	// Should error
	err := limiter.CheckExecutionTimeLimit()
	if err == nil {
		t.Error("expected error for exceeded execution time limit")
	}

	// Check error type
	if _, ok := err.(*ResourceLimitError); !ok {
		t.Errorf("expected ResourceLimitError, got %T", err)
	}

	// Check error details
	limitErr := err.(*ResourceLimitError)
	if limitErr.Type != LimitTypeExecutionTime {
		t.Errorf("expected error type LimitTypeExecutionTime, got %v", limitErr.Type)
	}
}

func TestResourceLimiter_CheckExecutionTimeLimit_NotStarted(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetExecutionTimeLimit(1) // Set short limit

	// Should not error because tracking hasn't started
	if err := limiter.CheckExecutionTimeLimit(); err != nil {
		t.Errorf("expected no error when tracking not started, got %v", err)
	}
}

func TestResourceLimiter_IncrementInstructionCount(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMaxInstructions(5)

	// Increment within limit
	for i := 0; i < 5; i++ {
		if err := limiter.IncrementInstructionCount(); err != nil {
			t.Errorf("expected no error for instruction %d, got %v", i+1, err)
		}
	}

	// Check count
	if got := limiter.GetInstructionCount(); got != 5 {
		t.Errorf("expected instruction count of 5, got %d", got)
	}

	// Exceed limit
	err := limiter.IncrementInstructionCount()
	if err == nil {
		t.Error("expected error for exceeded instruction limit")
	}

	// Check error type
	limitErr, ok := err.(*ResourceLimitError)
	if !ok {
		t.Errorf("expected ResourceLimitError, got %T", err)
	} else {
		if limitErr.Type != LimitTypeInstructions {
			t.Errorf("expected error type LimitTypeInstructions, got %v", limitErr.Type)
		}
		if limitErr.Current != 6 {
			t.Errorf("expected current count of 6, got %d", limitErr.Current)
		}
		if limitErr.Limit != 5 {
			t.Errorf("expected limit of 5, got %d", limitErr.Limit)
		}
	}
}

func TestResourceLimiter_IncrementInstructionCount_NoLimit(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMaxInstructions(0) // Unlimited

	// Should be able to increment many times without error
	for i := 0; i < 1000; i++ {
		if err := limiter.IncrementInstructionCount(); err != nil {
			t.Errorf("expected no error for unlimited instructions, got %v", err)
		}
	}

	// Note: When limit is 0 (unlimited), counter is not incremented (optimization)
	// So we don't test the counter value here, just that there's no error
}

func TestResourceLimiter_RecursionDepth(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMaxRecursionDepth(3)

	// Increment within limit
	for i := 0; i < 3; i++ {
		if err := limiter.IncrementRecursionDepth(); err != nil {
			t.Errorf("expected no error for recursion depth %d, got %v", i+1, err)
		}
	}

	// Check depth
	if got := limiter.GetRecursionDepth(); got != 3 {
		t.Errorf("expected recursion depth of 3, got %d", got)
	}

	// Exceed limit
	err := limiter.IncrementRecursionDepth()
	if err == nil {
		t.Error("expected error for exceeded recursion depth")
	}

	// Check error type
	limitErr, ok := err.(*ResourceLimitError)
	if !ok {
		t.Errorf("expected ResourceLimitError, got %T", err)
	} else {
		if limitErr.Type != LimitTypeRecursionDepth {
			t.Errorf("expected error type LimitTypeRecursionDepth, got %v", limitErr.Type)
		}
	}

	// Test decrement
	limiter.DecrementRecursionDepth()
	if got := limiter.GetRecursionDepth(); got != 3 {
		t.Errorf("expected recursion depth of 3 after decrement, got %d", got)
	}

	// Decrement to 0
	limiter.DecrementRecursionDepth()
	limiter.DecrementRecursionDepth()
	limiter.DecrementRecursionDepth()
	if got := limiter.GetRecursionDepth(); got != 0 {
		t.Errorf("expected recursion depth of 0 after all decrements, got %d", got)
	}
}

func TestResourceLimiter_AddOutputSize(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMaxOutputSize(100)

	// Add within limit
	if err := limiter.AddOutputSize(50); err != nil {
		t.Errorf("expected no error for output size within limit, got %v", err)
	}

	if got := limiter.GetOutputSize(); got != 50 {
		t.Errorf("expected output size of 50, got %d", got)
	}

	// Add more within limit
	if err := limiter.AddOutputSize(30); err != nil {
		t.Errorf("expected no error for output size within limit, got %v", err)
	}

	if got := limiter.GetOutputSize(); got != 80 {
		t.Errorf("expected output size of 80, got %d", got)
	}

	// Exceed limit
	err := limiter.AddOutputSize(50)
	if err == nil {
		t.Error("expected error for exceeded output size limit")
	}

	// Check error type
	limitErr, ok := err.(*ResourceLimitError)
	if !ok {
		t.Errorf("expected ResourceLimitError, got %T", err)
	} else {
		if limitErr.Type != LimitTypeOutputSize {
			t.Errorf("expected error type LimitTypeOutputSize, got %v", limitErr.Type)
		}
		if limitErr.Current != 130 {
			t.Errorf("expected current size of 130, got %d", limitErr.Current)
		}
		if limitErr.Limit != 100 {
			t.Errorf("expected limit of 100, got %d", limitErr.Limit)
		}
	}
}

func TestResourceLimiter_CheckAllLimits(t *testing.T) {
	limiter := NewResourceLimiter()

	// Set limits that won't be exceeded
	limiter.SetMemoryLimit(1024 * 1024 * 1024) // 1GB (shouldn't exceed)
	limiter.SetExecutionTimeLimit(10)          // 10 seconds
	limiter.Start()

	// Should not error
	if err := limiter.CheckAllLimits(); err != nil {
		t.Errorf("expected no error for non-exceeded limits, got %v", err)
	}

	// Now set a limit that will be exceeded
	limiter.SetMemoryLimit(1) // 1 byte (will definitely exceed)

	// Should error
	err := limiter.CheckAllLimits()
	if err == nil {
		t.Error("expected error when a limit is exceeded")
	}
}

func TestResourceLimiter_GetStats(t *testing.T) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(1024 * 1024)
	limiter.SetExecutionTimeLimit(30)
	limiter.SetMaxRecursionDepth(100)
	limiter.SetMaxInstructions(10000)
	limiter.SetMaxOutputSize(2048)
	limiter.Start()

	// Add some usage
	limiter.limits.instructionCount.Add(500)
	limiter.limits.recursionDepth.Add(10)
	limiter.limits.outputSize.Add(256)

	// Get stats
	stats := limiter.GetStats()

	// Check structure
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}

	// Check memory stats
	if memory, ok := stats["memory"].(map[string]interface{}); ok {
		if limit := memory["limit"].(int64); limit != 1024*1024 {
			t.Errorf("expected memory limit of 1048576, got %d", limit)
		}
		if enabled := memory["enabled"].(bool); !enabled {
			t.Error("expected memory limit to be enabled")
		}
	} else {
		t.Error("expected memory stats in result")
	}

	// Check execution time stats
	if execTime, ok := stats["execution_time"].(map[string]interface{}); ok {
		if limit := execTime["limit"].(float64); limit != 30.0 {
			t.Errorf("expected execution time limit of 30, got %f", limit)
		}
		if enabled := execTime["enabled"].(bool); !enabled {
			t.Error("expected execution time limit to be enabled")
		}
	} else {
		t.Error("expected execution_time stats in result")
	}

	// Check recursion depth stats
	if recursion, ok := stats["recursion_depth"].(map[string]interface{}); ok {
		if current := recursion["current"].(int64); current != 10 {
			t.Errorf("expected current recursion depth of 10, got %d", current)
		}
		if limit := recursion["limit"].(int64); limit != 100 {
			t.Errorf("expected recursion depth limit of 100, got %d", limit)
		}
	} else {
		t.Error("expected recursion_depth stats in result")
	}

	// Check instructions stats
	if instructions, ok := stats["instructions"].(map[string]interface{}); ok {
		if current := instructions["current"].(int64); current != 500 {
			t.Errorf("expected current instruction count of 500, got %d", current)
		}
		if limit := instructions["limit"].(int64); limit != 10000 {
			t.Errorf("expected instruction limit of 10000, got %d", limit)
		}
	} else {
		t.Error("expected instructions stats in result")
	}

	// Check output size stats
	if output, ok := stats["output_size"].(map[string]interface{}); ok {
		if current := output["current"].(int64); current != 256 {
			t.Errorf("expected current output size of 256, got %d", current)
		}
		if limit := output["limit"].(int64); limit != 2048 {
			t.Errorf("expected output size limit of 2048, got %d", limit)
		}
	} else {
		t.Error("expected output_size stats in result")
	}
}

func TestResourceLimiter_Callbacks(t *testing.T) {
	limiter := NewResourceLimiter()

	// Test memory limit callback
	var memoryCallbackCalled bool
	limiter.SetOnMemoryLimitExceeded(func(current, limit int64) {
		memoryCallbackCalled = true
		if current <= 0 {
			t.Error("expected positive current memory value")
		}
		if limit != 1 {
			t.Errorf("expected limit of 1, got %d", limit)
		}
	})
	limiter.SetMemoryLimit(1)
	limiter.CheckMemoryLimit()
	if !memoryCallbackCalled {
		t.Error("expected memory limit callback to be called")
	}

	// Test execution time callback
	var execTimeCallbackCalled bool
	limiter.SetOnExecutionTimeLimitExceeded(func(elapsed, limit time.Duration) {
		execTimeCallbackCalled = true
		if elapsed <= 0 {
			t.Error("expected positive elapsed time")
		}
		if limit != 1 {
			t.Errorf("expected limit of 1 nanosecond, got %v", limit)
		}
	})
	limiter.limits.ExecutionTimeLimit = 1
	limiter.Start()
	time.Sleep(10 * time.Millisecond)
	limiter.CheckExecutionTimeLimit()
	if !execTimeCallbackCalled {
		t.Error("expected execution time callback to be called")
	}

	// Test recursion depth callback
	limiter2 := NewResourceLimiter()
	var recursionCallbackCalled bool
	limiter2.SetOnRecursionDepthLimitExceeded(func(current, limit int64) {
		recursionCallbackCalled = true
		if current != 2 {
			t.Errorf("expected current depth of 2, got %d", current)
		}
		if limit != 1 {
			t.Errorf("expected limit of 1, got %d", limit)
		}
	})
	limiter2.SetMaxRecursionDepth(1)
	limiter2.IncrementRecursionDepth()
	limiter2.IncrementRecursionDepth()
	if !recursionCallbackCalled {
		t.Error("expected recursion depth callback to be called")
	}

	// Test instruction limit callback
	limiter3 := NewResourceLimiter()
	var instructionCallbackCalled bool
	limiter3.SetOnInstructionLimitExceeded(func(current, limit int64) {
		instructionCallbackCalled = true
		if current != 6 {
			t.Errorf("expected current count of 6, got %d", current)
		}
		if limit != 5 {
			t.Errorf("expected limit of 5, got %d", limit)
		}
	})
	limiter3.SetMaxInstructions(5)
	for i := 0; i < 6; i++ {
		limiter3.IncrementInstructionCount()
	}
	if !instructionCallbackCalled {
		t.Error("expected instruction limit callback to be called")
	}

	// Test output size callback
	limiter4 := NewResourceLimiter()
	var outputCallbackCalled bool
	limiter4.SetOnOutputSizeLimitExceeded(func(current, limit int64) {
		outputCallbackCalled = true
		if current != 150 {
			t.Errorf("expected current size of 150, got %d", current)
		}
		if limit != 100 {
			t.Errorf("expected limit of 100, got %d", limit)
		}
	})
	limiter4.SetMaxOutputSize(100)
	limiter4.AddOutputSize(150)
	if !outputCallbackCalled {
		t.Error("expected output size callback to be called")
	}
}

func TestResourceLimitError_Error(t *testing.T) {
	// Test with custom message
	err := &ResourceLimitError{
		Type:    LimitTypeMemory,
		Limit:   1024,
		Current: 2048,
		Message: "custom error message",
	}
	if got := err.Error(); got != "custom error message" {
		t.Errorf("expected custom message, got %s", got)
	}

	// Test without custom message
	err = &ResourceLimitError{
		Type:    LimitTypeExecutionTime,
		Limit:   5000,
		Current: 10000,
	}
	expected := "execution_time limit exceeded: current=10000, limit=5000"
	if got := err.Error(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestResourceLimitType_String(t *testing.T) {
	tests := []struct {
		limitType ResourceLimitType
		expected  string
	}{
		{LimitTypeMemory, "memory"},
		{LimitTypeExecutionTime, "execution_time"},
		{LimitTypeRecursionDepth, "recursion_depth"},
		{LimitTypeInstructions, "instructions"},
		{LimitTypeOutputSize, "output_size"},
		{ResourceLimitType(999), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.limitType.String(); got != tt.expected {
			t.Errorf("expected %s for type %d, got %s", tt.expected, tt.limitType, got)
		}
	}
}

func TestResourceLimiter_PackageLevelFunctions(t *testing.T) {
	// Reset global limiter before testing
	globalLimiter = nil
	limiterOnce = sync.Once{}

	// Test package-level setters and getters
	SetMemoryLimit(2048)
	SetExecutionTimeLimit(60)
	SetMaxRecursionDepth(500)
	SetMaxInstructions(50000)
	SetMaxOutputSize(4096)

	limiter := GetResourceLimiter()
	if got := limiter.GetMemoryLimit(); got != 2048 {
		t.Errorf("expected memory limit of 2048, got %d", got)
	}
	if got := limiter.GetExecutionTimeLimit(); got != 60 {
		t.Errorf("expected execution time limit of 60, got %d", got)
	}
	if got := limiter.GetMaxRecursionDepth(); got != 500 {
		t.Errorf("expected max recursion depth of 500, got %d", got)
	}
	if got := limiter.GetMaxInstructions(); got != 50000 {
		t.Errorf("expected max instructions of 50000, got %d", got)
	}
	if got := limiter.GetMaxOutputSize(); got != 4096 {
		t.Errorf("expected max output size of 4096, got %d", got)
	}

	// Test package-level Start/Stop
	StartResourceTracking()
	if limiter.limits.startTime.IsZero() {
		t.Error("expected start time to be set")
	}

	StopResourceTracking()

	// Test package-level CheckResourceLimits
	limiter.Reset()
	limiter.SetMemoryLimit(1024 * 1024 * 1024) // High limit
	if err := CheckResourceLimits(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test package-level GetResourceStats
	stats := GetResourceStats()
	if stats == nil {
		t.Error("expected non-nil stats")
	}
}

// Benchmarks

func BenchmarkResourceLimiter_IncrementInstructionCount(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMaxInstructions(0) // Unlimited to avoid errors

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.IncrementInstructionCount()
	}
}

func BenchmarkResourceLimiter_IncrementRecursionDepth(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMaxRecursionDepth(0) // Unlimited to avoid errors

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.IncrementRecursionDepth()
	}
}

func BenchmarkResourceLimiter_DecrementRecursionDepth(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMaxRecursionDepth(0) // Unlimited to avoid errors

	// Pre-increment to avoid going negative
	for i := 0; i < b.N; i++ {
		limiter.IncrementRecursionDepth()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.DecrementRecursionDepth()
	}
}

func BenchmarkResourceLimiter_AddOutputSize(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMaxOutputSize(0) // Unlimited to avoid errors

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.AddOutputSize(100)
	}
}

func BenchmarkResourceLimiter_CheckMemoryLimit(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(1024 * 1024 * 1024) // High limit to avoid failures

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.CheckMemoryLimit()
	}
}

func BenchmarkResourceLimiter_CheckExecutionTimeLimit(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetExecutionTimeLimit(3600) // 1 hour to avoid failures
	limiter.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.CheckExecutionTimeLimit()
	}
}

func BenchmarkResourceLimiter_CheckAllLimits(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(1024 * 1024 * 1024) // High limits
	limiter.SetExecutionTimeLimit(3600)
	limiter.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.CheckAllLimits()
	}
}

func BenchmarkResourceLimiter_GetStats(b *testing.B) {
	limiter := NewResourceLimiter()
	limiter.SetMemoryLimit(1024 * 1024)
	limiter.SetExecutionTimeLimit(30)
	limiter.SetMaxRecursionDepth(100)
	limiter.SetMaxInstructions(10000)
	limiter.SetMaxOutputSize(2048)
	limiter.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.GetStats()
	}
}
