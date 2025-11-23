package parallel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Performance Monitoring - Metrics and profiling for parallel execution
// ============================================================================

// ParallelMetrics tracks detailed metrics for parallel execution
type ParallelMetrics struct {
	// Execution counts
	totalTasks       atomic.Int64
	completedTasks   atomic.Int64
	failedTasks      atomic.Int64
	panickedTasks    atomic.Int64

	// Timing metrics
	totalDuration    atomic.Int64 // nanoseconds
	minDuration      atomic.Int64 // nanoseconds
	maxDuration      atomic.Int64 // nanoseconds

	// Concurrency metrics
	currentConcurrency atomic.Int32
	maxConcurrency     atomic.Int32

	// Contention metrics
	lockWaits        atomic.Int64
	channelWaits     atomic.Int64

	// Task history (limited size)
	recentTasks      []TaskMetric
	recentTasksMu    sync.RWMutex
	maxRecentTasks   int

	startTime        time.Time
	mu               sync.RWMutex
}

// TaskMetric represents metrics for a single task
type TaskMetric struct {
	ID           string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Success      bool
	Error        string
	Panicked     bool
	Concurrency  int32 // Concurrent tasks when this started
}

// NewParallelMetrics creates a new metrics tracker
func NewParallelMetrics() *ParallelMetrics {
	pm := &ParallelMetrics{
		startTime:      time.Now(),
		maxRecentTasks: 100, // Keep last 100 tasks
		recentTasks:    make([]TaskMetric, 0, 100),
	}
	pm.minDuration.Store(int64(^uint64(0) >> 1)) // Max int64
	return pm
}

// RecordTaskStart records the start of a task
func (pm *ParallelMetrics) RecordTaskStart(id string) int32 {
	pm.totalTasks.Add(1)
	concurrency := pm.currentConcurrency.Add(1)

	// Update max concurrency
	for {
		max := pm.maxConcurrency.Load()
		if concurrency <= max || pm.maxConcurrency.CompareAndSwap(max, concurrency) {
			break
		}
	}

	return concurrency
}

// RecordTaskEnd records the completion of a task
func (pm *ParallelMetrics) RecordTaskEnd(id string, startTime time.Time, err error, panicked bool, concurrency int32) {
	pm.currentConcurrency.Add(-1)

	duration := time.Since(startTime)
	durationNanos := duration.Nanoseconds()

	pm.totalDuration.Add(durationNanos)

	// Update min duration
	for {
		min := pm.minDuration.Load()
		if durationNanos >= min || pm.minDuration.CompareAndSwap(min, durationNanos) {
			break
		}
	}

	// Update max duration
	for {
		max := pm.maxDuration.Load()
		if durationNanos <= max || pm.maxDuration.CompareAndSwap(max, durationNanos) {
			break
		}
	}

	if panicked {
		pm.panickedTasks.Add(1)
		pm.failedTasks.Add(1)
	} else if err != nil {
		pm.failedTasks.Add(1)
	} else {
		pm.completedTasks.Add(1)
	}

	// Record task metric
	metric := TaskMetric{
		ID:          id,
		StartTime:   startTime,
		EndTime:     time.Now(),
		Duration:    duration,
		Success:     err == nil && !panicked,
		Panicked:    panicked,
		Concurrency: concurrency,
	}
	if err != nil {
		metric.Error = err.Error()
	}

	pm.addRecentTask(metric)
}

// RecordLockWait records a lock wait
func (pm *ParallelMetrics) RecordLockWait() {
	pm.lockWaits.Add(1)
}

// RecordChannelWait records a channel wait
func (pm *ParallelMetrics) RecordChannelWait() {
	pm.channelWaits.Add(1)
}

// addRecentTask adds a task to recent history
func (pm *ParallelMetrics) addRecentTask(metric TaskMetric) {
	pm.recentTasksMu.Lock()
	defer pm.recentTasksMu.Unlock()

	if len(pm.recentTasks) >= pm.maxRecentTasks {
		// Remove oldest
		pm.recentTasks = pm.recentTasks[1:]
	}

	pm.recentTasks = append(pm.recentTasks, metric)
}

// GetStats returns current statistics
func (pm *ParallelMetrics) GetStats() MetricsStats {
	total := pm.totalTasks.Load()
	completed := pm.completedTasks.Load()
	totalDuration := pm.totalDuration.Load()

	var avgDuration time.Duration
	if completed > 0 {
		avgDuration = time.Duration(totalDuration / completed)
	}

	return MetricsStats{
		TotalTasks:         total,
		CompletedTasks:     completed,
		FailedTasks:        pm.failedTasks.Load(),
		PanickedTasks:      pm.panickedTasks.Load(),
		CurrentConcurrency: pm.currentConcurrency.Load(),
		MaxConcurrency:     pm.maxConcurrency.Load(),
		AverageDuration:    avgDuration,
		MinDuration:        time.Duration(pm.minDuration.Load()),
		MaxDuration:        time.Duration(pm.maxDuration.Load()),
		LockWaits:          pm.lockWaits.Load(),
		ChannelWaits:       pm.channelWaits.Load(),
		Uptime:             time.Since(pm.startTime),
	}
}

// GetRecentTasks returns recent task history
func (pm *ParallelMetrics) GetRecentTasks(limit int) []TaskMetric {
	pm.recentTasksMu.RLock()
	defer pm.recentTasksMu.RUnlock()

	if limit <= 0 || limit > len(pm.recentTasks) {
		limit = len(pm.recentTasks)
	}

	// Return last N tasks
	start := len(pm.recentTasks) - limit
	result := make([]TaskMetric, limit)
	copy(result, pm.recentTasks[start:])
	return result
}

// Reset resets all metrics
func (pm *ParallelMetrics) Reset() {
	pm.totalTasks.Store(0)
	pm.completedTasks.Store(0)
	pm.failedTasks.Store(0)
	pm.panickedTasks.Store(0)
	pm.totalDuration.Store(0)
	pm.minDuration.Store(int64(^uint64(0) >> 1))
	pm.maxDuration.Store(0)
	pm.currentConcurrency.Store(0)
	pm.maxConcurrency.Store(0)
	pm.lockWaits.Store(0)
	pm.channelWaits.Store(0)

	pm.recentTasksMu.Lock()
	pm.recentTasks = make([]TaskMetric, 0, pm.maxRecentTasks)
	pm.recentTasksMu.Unlock()

	pm.mu.Lock()
	pm.startTime = time.Now()
	pm.mu.Unlock()
}

// MetricsStats contains snapshot of metrics
type MetricsStats struct {
	TotalTasks         int64
	CompletedTasks     int64
	FailedTasks        int64
	PanickedTasks      int64
	CurrentConcurrency int32
	MaxConcurrency     int32
	AverageDuration    time.Duration
	MinDuration        time.Duration
	MaxDuration        time.Duration
	LockWaits          int64
	ChannelWaits       int64
	Uptime             time.Duration
}

// String returns a formatted string of stats
func (ms MetricsStats) String() string {
	successRate := float64(0)
	if ms.TotalTasks > 0 {
		successRate = float64(ms.CompletedTasks) / float64(ms.TotalTasks) * 100
	}

	return fmt.Sprintf(`Parallel Metrics:
  Total Tasks:     %d
  Completed:       %d (%.1f%%)
  Failed:          %d
  Panicked:        %d
  Concurrency:     %d (current), %d (max)
  Duration:        %v (avg), %v (min), %v (max)
  Contention:      %d lock waits, %d channel waits
  Uptime:          %v`,
		ms.TotalTasks, ms.CompletedTasks, successRate,
		ms.FailedTasks, ms.PanickedTasks,
		ms.CurrentConcurrency, ms.MaxConcurrency,
		ms.AverageDuration, ms.MinDuration, ms.MaxDuration,
		ms.LockWaits, ms.ChannelWaits,
		ms.Uptime)
}

// ============================================================================
// Profiler - Performance profiling for parallel execution
// ============================================================================

// Profiler tracks performance profiles
type Profiler struct {
	enabled    atomic.Bool
	profiles   map[string]*Profile
	mu         sync.RWMutex
	startTime  time.Time
}

// Profile represents a performance profile
type Profile struct {
	Name         string
	CallCount    atomic.Int64
	TotalTime    atomic.Int64 // nanoseconds
	MinTime      atomic.Int64 // nanoseconds
	MaxTime      atomic.Int64 // nanoseconds
	mu           sync.RWMutex
}

// NewProfiler creates a new profiler
func NewProfiler() *Profiler {
	p := &Profiler{
		profiles:  make(map[string]*Profile),
		startTime: time.Now(),
	}
	p.enabled.Store(true)
	return p
}

// Enable enables profiling
func (p *Profiler) Enable() {
	p.enabled.Store(true)
}

// Disable disables profiling
func (p *Profiler) Disable() {
	p.enabled.Store(false)
}

// IsEnabled returns true if profiling is enabled
func (p *Profiler) IsEnabled() bool {
	return p.enabled.Load()
}

// Start starts profiling a section
func (p *Profiler) Start(name string) *ProfilerSession {
	if !p.enabled.Load() {
		return nil
	}

	profile := p.getProfile(name)
	return &ProfilerSession{
		profile:   profile,
		startTime: time.Now(),
	}
}

// getProfile gets or creates a profile
func (p *Profiler) getProfile(name string) *Profile {
	p.mu.RLock()
	profile, exists := p.profiles[name]
	p.mu.RUnlock()

	if exists {
		return profile
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check
	profile, exists = p.profiles[name]
	if exists {
		return profile
	}

	profile = &Profile{
		Name: name,
	}
	profile.MinTime.Store(int64(^uint64(0) >> 1)) // Max int64

	p.profiles[name] = profile
	return profile
}

// GetProfile returns a specific profile
func (p *Profiler) GetProfile(name string) *ProfileStats {
	p.mu.RLock()
	profile, exists := p.profiles[name]
	p.mu.RUnlock()

	if !exists {
		return nil
	}

	return profile.GetStats()
}

// GetAllProfiles returns all profiles
func (p *Profiler) GetAllProfiles() map[string]*ProfileStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*ProfileStats, len(p.profiles))
	for name, profile := range p.profiles {
		result[name] = profile.GetStats()
	}

	return result
}

// Reset resets all profiles
func (p *Profiler) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.profiles = make(map[string]*Profile)
	p.startTime = time.Now()
}

// ProfilerSession represents an active profiling session
type ProfilerSession struct {
	profile   *Profile
	startTime time.Time
}

// End ends the profiling session
func (ps *ProfilerSession) End() {
	if ps == nil || ps.profile == nil {
		return
	}

	duration := time.Since(ps.startTime)
	durationNanos := duration.Nanoseconds()

	ps.profile.CallCount.Add(1)
	ps.profile.TotalTime.Add(durationNanos)

	// Update min
	for {
		min := ps.profile.MinTime.Load()
		if durationNanos >= min || ps.profile.MinTime.CompareAndSwap(min, durationNanos) {
			break
		}
	}

	// Update max
	for {
		max := ps.profile.MaxTime.Load()
		if durationNanos <= max || ps.profile.MaxTime.CompareAndSwap(max, durationNanos) {
			break
		}
	}
}

// GetStats returns statistics for this profile
func (p *Profile) GetStats() *ProfileStats {
	callCount := p.CallCount.Load()
	totalTime := p.TotalTime.Load()

	var avgTime time.Duration
	if callCount > 0 {
		avgTime = time.Duration(totalTime / callCount)
	}

	return &ProfileStats{
		Name:      p.Name,
		CallCount: callCount,
		TotalTime: time.Duration(totalTime),
		AvgTime:   avgTime,
		MinTime:   time.Duration(p.MinTime.Load()),
		MaxTime:   time.Duration(p.MaxTime.Load()),
	}
}

// ProfileStats contains profile statistics
type ProfileStats struct {
	Name      string
	CallCount int64
	TotalTime time.Duration
	AvgTime   time.Duration
	MinTime   time.Duration
	MaxTime   time.Duration
}

// String returns formatted profile stats
func (ps *ProfileStats) String() string {
	return fmt.Sprintf(`Profile "%s":
  Calls:     %d
  Total:     %v
  Average:   %v
  Min:       %v
  Max:       %v`,
		ps.Name, ps.CallCount,
		ps.TotalTime, ps.AvgTime,
		ps.MinTime, ps.MaxTime)
}

// ============================================================================
// ContentionDetector - Detects and reports contention
// ============================================================================

// ContentionDetector detects lock and resource contention
type ContentionDetector struct {
	enabled       atomic.Bool
	threshold     time.Duration // Threshold for reporting contention
	contentions   []ContentionEvent
	contentionsMu sync.RWMutex
	maxEvents     int
}

// ContentionEvent represents a contention event
type ContentionEvent struct {
	Timestamp  time.Time
	Resource   string
	WaitTime   time.Duration
	Goroutines int
}

// NewContentionDetector creates a new contention detector
func NewContentionDetector(threshold time.Duration) *ContentionDetector {
	cd := &ContentionDetector{
		threshold:   threshold,
		maxEvents:   1000,
		contentions: make([]ContentionEvent, 0, 1000),
	}
	cd.enabled.Store(true)
	return cd
}

// Enable enables contention detection
func (cd *ContentionDetector) Enable() {
	cd.enabled.Store(true)
}

// Disable disables contention detection
func (cd *ContentionDetector) Disable() {
	cd.enabled.Store(false)
}

// IsEnabled returns true if detection is enabled
func (cd *ContentionDetector) IsEnabled() bool {
	return cd.enabled.Load()
}

// RecordContention records a contention event
func (cd *ContentionDetector) RecordContention(resource string, waitTime time.Duration, goroutines int) {
	if !cd.enabled.Load() {
		return
	}

	if waitTime < cd.threshold {
		return
	}

	event := ContentionEvent{
		Timestamp:  time.Now(),
		Resource:   resource,
		WaitTime:   waitTime,
		Goroutines: goroutines,
	}

	cd.contentionsMu.Lock()
	defer cd.contentionsMu.Unlock()

	if len(cd.contentions) >= cd.maxEvents {
		// Remove oldest
		cd.contentions = cd.contentions[1:]
	}

	cd.contentions = append(cd.contentions, event)
}

// GetContentions returns recent contention events
func (cd *ContentionDetector) GetContentions(limit int) []ContentionEvent {
	cd.contentionsMu.RLock()
	defer cd.contentionsMu.RUnlock()

	if limit <= 0 || limit > len(cd.contentions) {
		limit = len(cd.contentions)
	}

	// Return last N events
	start := len(cd.contentions) - limit
	result := make([]ContentionEvent, limit)
	copy(result, cd.contentions[start:])
	return result
}

// GetContentionCount returns the total number of contention events
func (cd *ContentionDetector) GetContentionCount() int {
	cd.contentionsMu.RLock()
	defer cd.contentionsMu.RUnlock()
	return len(cd.contentions)
}

// Clear clears all contention events
func (cd *ContentionDetector) Clear() {
	cd.contentionsMu.Lock()
	defer cd.contentionsMu.Unlock()
	cd.contentions = make([]ContentionEvent, 0, cd.maxEvents)
}

// SetThreshold sets the contention threshold
func (cd *ContentionDetector) SetThreshold(threshold time.Duration) {
	cd.threshold = threshold
}

// GetThreshold returns the current threshold
func (cd *ContentionDetector) GetThreshold() time.Duration {
	return cd.threshold
}

// ============================================================================
// GlobalMetrics - Global metrics instance
// ============================================================================

var (
	globalMetrics  *ParallelMetrics
	globalProfiler *Profiler
	globalDetector *ContentionDetector
	metricsOnce    sync.Once
)

// GetGlobalMetrics returns the global metrics instance
func GetGlobalMetrics() *ParallelMetrics {
	metricsOnce.Do(func() {
		globalMetrics = NewParallelMetrics()
		globalProfiler = NewProfiler()
		globalDetector = NewContentionDetector(1 * time.Millisecond)
	})
	return globalMetrics
}

// GetGlobalProfiler returns the global profiler instance
func GetGlobalProfiler() *Profiler {
	metricsOnce.Do(func() {
		globalMetrics = NewParallelMetrics()
		globalProfiler = NewProfiler()
		globalDetector = NewContentionDetector(1 * time.Millisecond)
	})
	return globalProfiler
}

// GetGlobalContentionDetector returns the global contention detector
func GetGlobalContentionDetector() *ContentionDetector {
	metricsOnce.Do(func() {
		globalMetrics = NewParallelMetrics()
		globalProfiler = NewProfiler()
		globalDetector = NewContentionDetector(1 * time.Millisecond)
	})
	return globalDetector
}
