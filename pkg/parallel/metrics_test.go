package parallel

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// ============================================================================
// ParallelMetrics Tests
// ============================================================================

func TestNewParallelMetrics(t *testing.T) {
	pm := NewParallelMetrics()
	if pm == nil {
		t.Error("NewParallelMetrics returned nil")
	}

	stats := pm.GetStats()
	if stats.TotalTasks != 0 {
		t.Errorf("TotalTasks = %d, expected 0", stats.TotalTasks)
	}
	if stats.CurrentConcurrency != 0 {
		t.Errorf("CurrentConcurrency = %d, expected 0", stats.CurrentConcurrency)
	}
}

func TestMetricsRecordTaskStart(t *testing.T) {
	pm := NewParallelMetrics()

	concurrency := pm.RecordTaskStart("task1")
	if concurrency != 1 {
		t.Errorf("Concurrency = %d, expected 1", concurrency)
	}

	stats := pm.GetStats()
	if stats.TotalTasks != 1 {
		t.Errorf("TotalTasks = %d, expected 1", stats.TotalTasks)
	}
	if stats.CurrentConcurrency != 1 {
		t.Errorf("CurrentConcurrency = %d, expected 1", stats.CurrentConcurrency)
	}
	if stats.MaxConcurrency != 1 {
		t.Errorf("MaxConcurrency = %d, expected 1", stats.MaxConcurrency)
	}
}

func TestMetricsRecordTaskEnd(t *testing.T) {
	pm := NewParallelMetrics()

	startTime := time.Now()
	concurrency := pm.RecordTaskStart("task1")

	time.Sleep(10 * time.Millisecond)

	pm.RecordTaskEnd("task1", startTime, nil, false, concurrency)

	stats := pm.GetStats()
	if stats.CompletedTasks != 1 {
		t.Errorf("CompletedTasks = %d, expected 1", stats.CompletedTasks)
	}
	if stats.FailedTasks != 0 {
		t.Errorf("FailedTasks = %d, expected 0", stats.FailedTasks)
	}
	if stats.CurrentConcurrency != 0 {
		t.Errorf("CurrentConcurrency = %d, expected 0", stats.CurrentConcurrency)
	}
	if stats.AverageDuration < 10*time.Millisecond {
		t.Error("AverageDuration should be >= 10ms")
	}
}

func TestMetricsRecordTaskError(t *testing.T) {
	pm := NewParallelMetrics()

	startTime := time.Now()
	concurrency := pm.RecordTaskStart("task1")

	pm.RecordTaskEnd("task1", startTime, errors.New("test error"), false, concurrency)

	stats := pm.GetStats()
	if stats.CompletedTasks != 0 {
		t.Errorf("CompletedTasks = %d, expected 0", stats.CompletedTasks)
	}
	if stats.FailedTasks != 1 {
		t.Errorf("FailedTasks = %d, expected 1", stats.FailedTasks)
	}
}

func TestMetricsRecordTaskPanic(t *testing.T) {
	pm := NewParallelMetrics()

	startTime := time.Now()
	concurrency := pm.RecordTaskStart("task1")

	pm.RecordTaskEnd("task1", startTime, nil, true, concurrency)

	stats := pm.GetStats()
	if stats.PanickedTasks != 1 {
		t.Errorf("PanickedTasks = %d, expected 1", stats.PanickedTasks)
	}
	if stats.FailedTasks != 1 {
		t.Errorf("FailedTasks = %d, expected 1", stats.FailedTasks)
	}
}

func TestMetricsMultipleTasks(t *testing.T) {
	pm := NewParallelMetrics()

	for i := 0; i < 10; i++ {
		startTime := time.Now()
		concurrency := pm.RecordTaskStart(fmt.Sprintf("task%d", i))
		time.Sleep(1 * time.Millisecond)

		var err error
		if i%3 == 0 {
			err = errors.New("error")
		}

		pm.RecordTaskEnd(fmt.Sprintf("task%d", i), startTime, err, false, concurrency)
	}

	stats := pm.GetStats()
	if stats.TotalTasks != 10 {
		t.Errorf("TotalTasks = %d, expected 10", stats.TotalTasks)
	}
	if stats.CompletedTasks != 6 {
		t.Errorf("CompletedTasks = %d, expected 6", stats.CompletedTasks)
	}
	if stats.FailedTasks != 4 {
		t.Errorf("FailedTasks = %d, expected 4", stats.FailedTasks)
	}
}

func TestMetricsConcurrency(t *testing.T) {
	pm := NewParallelMetrics()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			startTime := time.Now()
			concurrency := pm.RecordTaskStart(fmt.Sprintf("task%d", id))
			time.Sleep(20 * time.Millisecond)
			pm.RecordTaskEnd(fmt.Sprintf("task%d", id), startTime, nil, false, concurrency)
		}(i)
	}

	// Give goroutines time to start
	time.Sleep(5 * time.Millisecond)

	stats := pm.GetStats()
	if stats.CurrentConcurrency == 0 {
		t.Error("CurrentConcurrency should be > 0")
	}

	wg.Wait()

	stats = pm.GetStats()
	if stats.MaxConcurrency < 2 {
		t.Errorf("MaxConcurrency = %d, expected >= 2", stats.MaxConcurrency)
	}
	if stats.CurrentConcurrency != 0 {
		t.Errorf("CurrentConcurrency = %d, expected 0", stats.CurrentConcurrency)
	}
}

func TestMetricsRecordContention(t *testing.T) {
	pm := NewParallelMetrics()

	pm.RecordLockWait()
	pm.RecordLockWait()
	pm.RecordChannelWait()

	stats := pm.GetStats()
	if stats.LockWaits != 2 {
		t.Errorf("LockWaits = %d, expected 2", stats.LockWaits)
	}
	if stats.ChannelWaits != 1 {
		t.Errorf("ChannelWaits = %d, expected 1", stats.ChannelWaits)
	}
}

func TestMetricsRecentTasks(t *testing.T) {
	pm := NewParallelMetrics()

	for i := 0; i < 5; i++ {
		startTime := time.Now()
		concurrency := pm.RecordTaskStart(fmt.Sprintf("task%d", i))
		pm.RecordTaskEnd(fmt.Sprintf("task%d", i), startTime, nil, false, concurrency)
	}

	recent := pm.GetRecentTasks(3)
	if len(recent) != 3 {
		t.Errorf("Recent tasks = %d, expected 3", len(recent))
	}

	// Should be last 3 tasks
	if recent[0].ID != "task2" {
		t.Errorf("First recent task = %s, expected task2", recent[0].ID)
	}
	if recent[2].ID != "task4" {
		t.Errorf("Last recent task = %s, expected task4", recent[2].ID)
	}
}

func TestMetricsRecentTasksLimit(t *testing.T) {
	pm := NewParallelMetrics()
	pm.maxRecentTasks = 10

	// Add more than limit
	for i := 0; i < 15; i++ {
		startTime := time.Now()
		concurrency := pm.RecordTaskStart(fmt.Sprintf("task%d", i))
		pm.RecordTaskEnd(fmt.Sprintf("task%d", i), startTime, nil, false, concurrency)
	}

	recent := pm.GetRecentTasks(0) // 0 = all
	if len(recent) != 10 {
		t.Errorf("Recent tasks = %d, expected 10", len(recent))
	}

	// Should be last 10 tasks (5-14)
	if recent[0].ID != "task5" {
		t.Errorf("First recent task = %s, expected task5", recent[0].ID)
	}
}

func TestMetricsReset(t *testing.T) {
	pm := NewParallelMetrics()

	startTime := time.Now()
	concurrency := pm.RecordTaskStart("task1")
	pm.RecordTaskEnd("task1", startTime, nil, false, concurrency)

	pm.Reset()

	stats := pm.GetStats()
	if stats.TotalTasks != 0 {
		t.Errorf("TotalTasks = %d, expected 0 after reset", stats.TotalTasks)
	}

	recent := pm.GetRecentTasks(0)
	if len(recent) != 0 {
		t.Errorf("Recent tasks = %d, expected 0 after reset", len(recent))
	}
}

func TestMetricsStatsString(t *testing.T) {
	stats := MetricsStats{
		TotalTasks:     100,
		CompletedTasks: 90,
		FailedTasks:    10,
	}

	str := stats.String()
	if str == "" {
		t.Error("Stats string should not be empty")
	}
}

// ============================================================================
// Profiler Tests
// ============================================================================

func TestNewProfiler(t *testing.T) {
	p := NewProfiler()
	if p == nil {
		t.Error("NewProfiler returned nil")
	}
	if !p.IsEnabled() {
		t.Error("Profiler should be enabled by default")
	}
}

func TestProfilerEnableDisable(t *testing.T) {
	p := NewProfiler()

	p.Disable()
	if p.IsEnabled() {
		t.Error("Profiler should be disabled")
	}

	p.Enable()
	if !p.IsEnabled() {
		t.Error("Profiler should be enabled")
	}
}

func TestProfilerSession(t *testing.T) {
	p := NewProfiler()

	session := p.Start("test")
	if session == nil {
		t.Error("Start should return session")
	}

	time.Sleep(10 * time.Millisecond)
	session.End()

	stats := p.GetProfile("test")
	if stats == nil {
		t.Error("Profile should exist")
	}
	if stats.CallCount != 1 {
		t.Errorf("CallCount = %d, expected 1", stats.CallCount)
	}
	if stats.TotalTime < 10*time.Millisecond {
		t.Error("TotalTime should be >= 10ms")
	}
}

func TestProfilerMultipleCalls(t *testing.T) {
	p := NewProfiler()

	for i := 0; i < 5; i++ {
		session := p.Start("test")
		time.Sleep(1 * time.Millisecond)
		session.End()
	}

	stats := p.GetProfile("test")
	if stats.CallCount != 5 {
		t.Errorf("CallCount = %d, expected 5", stats.CallCount)
	}
	if stats.AvgTime < 1*time.Millisecond {
		t.Error("AvgTime should be >= 1ms")
	}
}

func TestProfilerMultipleProfiles(t *testing.T) {
	p := NewProfiler()

	session1 := p.Start("profile1")
	time.Sleep(5 * time.Millisecond)
	session1.End()

	session2 := p.Start("profile2")
	time.Sleep(10 * time.Millisecond)
	session2.End()

	all := p.GetAllProfiles()
	if len(all) != 2 {
		t.Errorf("Profile count = %d, expected 2", len(all))
	}

	stats1 := p.GetProfile("profile1")
	stats2 := p.GetProfile("profile2")

	if stats1 == nil || stats2 == nil {
		t.Error("Both profiles should exist")
	}

	if stats2.TotalTime <= stats1.TotalTime {
		t.Error("profile2 should have longer total time")
	}
}

func TestProfilerDisabled(t *testing.T) {
	p := NewProfiler()
	p.Disable()

	session := p.Start("test")
	if session != nil {
		t.Error("Start should return nil when disabled")
	}

	// End should not panic on nil session
	session.End()

	stats := p.GetProfile("test")
	if stats != nil {
		t.Error("Profile should not exist when disabled")
	}
}

func TestProfilerReset(t *testing.T) {
	p := NewProfiler()

	session := p.Start("test")
	session.End()

	p.Reset()

	stats := p.GetProfile("test")
	if stats != nil {
		t.Error("Profile should not exist after reset")
	}

	all := p.GetAllProfiles()
	if len(all) != 0 {
		t.Errorf("Profile count = %d, expected 0 after reset", len(all))
	}
}

func TestProfilerConcurrent(t *testing.T) {
	p := NewProfiler()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				session := p.Start("concurrent")
				time.Sleep(1 * time.Millisecond)
				session.End()
			}
		}()
	}

	wg.Wait()

	stats := p.GetProfile("concurrent")
	if stats.CallCount != 100 {
		t.Errorf("CallCount = %d, expected 100", stats.CallCount)
	}
}

func TestProfileStatsString(t *testing.T) {
	stats := &ProfileStats{
		Name:      "test",
		CallCount: 10,
		TotalTime: 100 * time.Millisecond,
		AvgTime:   10 * time.Millisecond,
	}

	str := stats.String()
	if str == "" {
		t.Error("Profile stats string should not be empty")
	}
}

// ============================================================================
// ContentionDetector Tests
// ============================================================================

func TestNewContentionDetector(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)
	if cd == nil {
		t.Error("NewContentionDetector returned nil")
	}
	if !cd.IsEnabled() {
		t.Error("Detector should be enabled by default")
	}
	if cd.GetThreshold() != 1*time.Millisecond {
		t.Error("Threshold should be 1ms")
	}
}

func TestContentionDetectorEnableDisable(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)

	cd.Disable()
	if cd.IsEnabled() {
		t.Error("Detector should be disabled")
	}

	cd.Enable()
	if !cd.IsEnabled() {
		t.Error("Detector should be enabled")
	}
}

func TestContentionDetectorRecord(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)

	cd.RecordContention("resource1", 5*time.Millisecond, 2)

	contentions := cd.GetContentions(0)
	if len(contentions) != 1 {
		t.Errorf("Contentions = %d, expected 1", len(contentions))
	}

	if contentions[0].Resource != "resource1" {
		t.Errorf("Resource = %s, expected resource1", contentions[0].Resource)
	}
	if contentions[0].WaitTime != 5*time.Millisecond {
		t.Error("WaitTime should be 5ms")
	}
	if contentions[0].Goroutines != 2 {
		t.Error("Goroutines should be 2")
	}
}

func TestContentionDetectorThreshold(t *testing.T) {
	cd := NewContentionDetector(10 * time.Millisecond)

	// Below threshold - should not record
	cd.RecordContention("resource1", 5*time.Millisecond, 1)

	if cd.GetContentionCount() != 0 {
		t.Error("Should not record below threshold")
	}

	// Above threshold - should record
	cd.RecordContention("resource2", 15*time.Millisecond, 1)

	if cd.GetContentionCount() != 1 {
		t.Error("Should record above threshold")
	}
}

func TestContentionDetectorLimit(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)
	cd.maxEvents = 10

	// Add more than limit
	for i := 0; i < 15; i++ {
		cd.RecordContention(fmt.Sprintf("resource%d", i), 5*time.Millisecond, 1)
	}

	if cd.GetContentionCount() != 10 {
		t.Errorf("Contention count = %d, expected 10", cd.GetContentionCount())
	}

	contentions := cd.GetContentions(0)
	// Should have last 10 (5-14)
	if contentions[0].Resource != "resource5" {
		t.Errorf("First resource = %s, expected resource5", contentions[0].Resource)
	}
}

func TestContentionDetectorGetLimit(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)

	for i := 0; i < 5; i++ {
		cd.RecordContention(fmt.Sprintf("resource%d", i), 5*time.Millisecond, 1)
	}

	contentions := cd.GetContentions(3)
	if len(contentions) != 3 {
		t.Errorf("Contentions = %d, expected 3", len(contentions))
	}

	// Should be last 3
	if contentions[0].Resource != "resource2" {
		t.Errorf("First resource = %s, expected resource2", contentions[0].Resource)
	}
}

func TestContentionDetectorClear(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)

	cd.RecordContention("resource1", 5*time.Millisecond, 1)
	cd.Clear()

	if cd.GetContentionCount() != 0 {
		t.Error("Count should be 0 after clear")
	}
}

func TestContentionDetectorSetThreshold(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)

	cd.SetThreshold(10 * time.Millisecond)

	if cd.GetThreshold() != 10*time.Millisecond {
		t.Error("Threshold should be 10ms")
	}
}

func TestContentionDetectorDisabled(t *testing.T) {
	cd := NewContentionDetector(1 * time.Millisecond)
	cd.Disable()

	cd.RecordContention("resource1", 5*time.Millisecond, 1)

	if cd.GetContentionCount() != 0 {
		t.Error("Should not record when disabled")
	}
}

// ============================================================================
// Global Metrics Tests
// ============================================================================

func TestGetGlobalMetrics(t *testing.T) {
	gm := GetGlobalMetrics()
	if gm == nil {
		t.Error("GetGlobalMetrics returned nil")
	}

	// Should return same instance
	gm2 := GetGlobalMetrics()
	if gm != gm2 {
		t.Error("Should return same instance")
	}
}

func TestGetGlobalProfiler(t *testing.T) {
	gp := GetGlobalProfiler()
	if gp == nil {
		t.Error("GetGlobalProfiler returned nil")
	}

	// Should return same instance
	gp2 := GetGlobalProfiler()
	if gp != gp2 {
		t.Error("Should return same instance")
	}
}

func TestGetGlobalContentionDetector(t *testing.T) {
	gcd := GetGlobalContentionDetector()
	if gcd == nil {
		t.Error("GetGlobalContentionDetector returned nil")
	}

	// Should return same instance
	gcd2 := GetGlobalContentionDetector()
	if gcd != gcd2 {
		t.Error("Should return same instance")
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestMetricsIntegration(t *testing.T) {
	pm := NewParallelMetrics()
	p := NewProfiler()

	// Simulate parallel workload
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Profile the task
			session := p.Start("task-execution")

			// Record task start
			startTime := time.Now()
			concurrency := pm.RecordTaskStart(fmt.Sprintf("task%d", id))

			// Simulate work
			time.Sleep(time.Duration(id) * time.Millisecond)

			// Record task end
			var err error
			if id%4 == 0 {
				err = errors.New("simulated error")
			}
			pm.RecordTaskEnd(fmt.Sprintf("task%d", id), startTime, err, false, concurrency)

			session.End()
		}(i)
	}

	wg.Wait()

	// Check metrics
	stats := pm.GetStats()
	if stats.TotalTasks != 10 {
		t.Errorf("TotalTasks = %d, expected 10", stats.TotalTasks)
	}
	if stats.CompletedTasks+stats.FailedTasks != 10 {
		t.Error("Completed + Failed should equal Total")
	}

	// Check profiler
	profileStats := p.GetProfile("task-execution")
	if profileStats.CallCount != 10 {
		t.Errorf("Profile CallCount = %d, expected 10", profileStats.CallCount)
	}
}
