package stress

import (
	"bytes"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestLoadTesting tests handling of high request volumes
// Runs many PHP scripts in parallel to simulate high load
func TestLoadTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	const (
		numRequests    = 1000 // Number of concurrent requests
		numGoroutines  = 100  // Number of concurrent goroutines
		scriptPath     = "testdata/simple.php"
	)

	t.Logf("Load Testing: %d requests across %d goroutines", numRequests, numGoroutines)

	var (
		successCount uint64
		failureCount uint64
		wg           sync.WaitGroup
	)

	startTime := time.Now()
	requestsPerGoroutine := numRequests / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				cmd := exec.Command("../../php-go", "run", scriptPath)
				err := cmd.Run()
				if err != nil {
					atomic.AddUint64(&failureCount, 1)
				} else {
					atomic.AddUint64(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// Results
	success := atomic.LoadUint64(&successCount)
	failure := atomic.LoadUint64(&failureCount)
	total := success + failure
	successRate := float64(success) / float64(total) * 100
	requestsPerSecond := float64(total) / duration.Seconds()

	t.Logf("Load Test Results:")
	t.Logf("  Duration: %v", duration)
	t.Logf("  Total requests: %d", total)
	t.Logf("  Successful: %d (%.2f%%)", success, successRate)
	t.Logf("  Failed: %d", failure)
	t.Logf("  Requests/second: %.2f", requestsPerSecond)

	// Assert success rate is above threshold
	if successRate < 95.0 {
		t.Errorf("Success rate %.2f%% is below 95%% threshold", successRate)
	}
}

// TestMemoryLeakDetection runs long-running tests to identify memory leaks
// Monitors memory usage over time while executing PHP scripts
func TestMemoryLeakDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak detection in short mode")
	}

	const (
		duration   = 30 * time.Second // Run for 30 seconds
		sampleRate = 1 * time.Second  // Sample memory every second
		scriptPath = "testdata/memory_intensive.php"
	)

	t.Logf("Memory Leak Detection: Running for %v", duration)

	var memStats runtime.MemStats
	var samples []uint64

	ticker := time.NewTicker(sampleRate)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(duration)
		done <- true
	}()

	// Run PHP scripts continuously
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				cmd := exec.Command("../../php-go", "run", scriptPath)
				cmd.Run()
			}
		}
	}()

	// Sample memory usage
loop:
	for {
		select {
		case <-ticker.C:
			runtime.ReadMemStats(&memStats)
			samples = append(samples, memStats.Alloc)
			t.Logf("Memory: %d MB (Alloc), %d MB (Sys)",
				memStats.Alloc/1024/1024,
				memStats.Sys/1024/1024)

		case <-done:
			break loop
		}
	}

	// Analyze memory growth
	if len(samples) < 2 {
		t.Skip("Not enough samples to analyze memory growth")
	}

	firstHalf := samples[:len(samples)/2]
	secondHalf := samples[len(samples)/2:]

	avgFirst := average(firstHalf)
	avgSecond := average(secondHalf)
	growthPercent := float64(avgSecond-avgFirst) / float64(avgFirst) * 100

	t.Logf("Memory Leak Analysis:")
	t.Logf("  Samples collected: %d", len(samples))
	t.Logf("  Avg first half: %.2f MB", float64(avgFirst)/1024/1024)
	t.Logf("  Avg second half: %.2f MB", float64(avgSecond)/1024/1024)
	t.Logf("  Growth: %.2f%%", growthPercent)

	// Warn if memory grew significantly (potential leak)
	if growthPercent > 50.0 {
		t.Errorf("Memory grew by %.2f%%, possible memory leak detected", growthPercent)
	}
}

// TestConcurrentRequestTesting runs multi-threaded stress tests
// Tests for race conditions and concurrent access issues
func TestConcurrentRequestTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent request test in short mode")
	}

	const (
		numConcurrent = 50                               // Number of concurrent workers
		duration      = 10 * time.Second                  // Run for 10 seconds
		scriptPath    = "testdata/concurrent_access.php"
	)

	t.Logf("Concurrent Request Testing: %d concurrent workers for %v", numConcurrent, duration)

	var (
		requestCount uint64
		errorCount   uint64
		wg           sync.WaitGroup
	)

	startTime := time.Now()
	done := make(chan bool)

	// Start timer
	go func() {
		time.Sleep(duration)
		done <- true
	}()

	// Start concurrent workers
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					cmd := exec.Command("../../php-go", "run", scriptPath)
					var stderr bytes.Buffer
					cmd.Stderr = &stderr

					err := cmd.Run()
					atomic.AddUint64(&requestCount, 1)

					if err != nil {
						atomic.AddUint64(&errorCount, 1)
						// Check for race condition errors
						if bytes.Contains(stderr.Bytes(), []byte("race")) ||
							bytes.Contains(stderr.Bytes(), []byte("concurrent")) {
							t.Errorf("Race condition detected in worker %d: %v", workerID, stderr.String())
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	requests := atomic.LoadUint64(&requestCount)
	errors := atomic.LoadUint64(&errorCount)
	errorRate := float64(errors) / float64(requests) * 100
	throughput := float64(requests) / elapsed.Seconds()

	t.Logf("Concurrent Request Results:")
	t.Logf("  Duration: %v", elapsed)
	t.Logf("  Total requests: %d", requests)
	t.Logf("  Errors: %d (%.2f%%)", errors, errorRate)
	t.Logf("  Throughput: %.2f requests/second", throughput)

	// Assert error rate is low
	if errorRate > 1.0 {
		t.Errorf("Error rate %.2f%% exceeds 1%% threshold", errorRate)
	}
}

// TestLongRunningProcess tests stability over extended periods
// Ensures the interpreter doesn't crash or degrade over time
func TestLongRunningProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running process test in short mode")
	}

	const (
		duration   = 5 * time.Minute // Run for 5 minutes
		scriptPath = "testdata/long_running.php"
	)

	t.Logf("Long-Running Process Test: Running for %v", duration)

	cmd := exec.Command("../../php-go", "run", scriptPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the process
	startTime := time.Now()
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}

	// Wait with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(duration):
		// Process is still running after duration - this is good
		cmd.Process.Kill()
		elapsed := time.Since(startTime)

		t.Logf("Long-Running Process Results:")
		t.Logf("  Duration: %v", elapsed)
		t.Logf("  Status: Process ran successfully for full duration")
		t.Logf("  Stdout size: %d bytes", stdout.Len())
		t.Logf("  Stderr size: %d bytes", stderr.Len())

		// Check if there were any errors in stderr
		if stderr.Len() > 0 {
			t.Logf("  Stderr content (first 500 bytes): %s", stderr.Bytes()[:min(500, stderr.Len())])
		}

	case err := <-done:
		// Process ended prematurely
		elapsed := time.Since(startTime)

		if err != nil {
			t.Errorf("Process ended prematurely after %v with error: %v", elapsed, err)
			t.Logf("  Stderr: %s", stderr.String())
		} else {
			t.Logf("Process completed normally after %v", elapsed)
		}
	}
}

// TestStressTestSuite runs a comprehensive stress test combining all scenarios
func TestStressTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive stress test suite in short mode")
	}

	t.Run("Load", func(t *testing.T) {
		// Reduced load test
		runReducedLoadTest(t)
	})

	t.Run("Memory", func(t *testing.T) {
		// Quick memory check
		runQuickMemoryCheck(t)
	})

	t.Run("Concurrent", func(t *testing.T) {
		// Concurrent access test
		runQuickConcurrentTest(t)
	})
}

// Helper functions

func average(values []uint64) uint64 {
	if len(values) == 0 {
		return 0
	}
	var sum uint64
	for _, v := range values {
		sum += v
	}
	return sum / uint64(len(values))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runReducedLoadTest(t *testing.T) {
	const (
		numRequests   = 100
		numGoroutines = 10
		scriptPath    = "testdata/simple.php"
	)

	var successCount uint64
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numRequests/numGoroutines; j++ {
				cmd := exec.Command("../../php-go", "run", scriptPath)
				if cmd.Run() == nil {
					atomic.AddUint64(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	t.Logf("Reduced load test: %d/%d successful", successCount, numRequests)
}

func runQuickMemoryCheck(t *testing.T) {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)
	before := memStats.Alloc

	// Run some PHP scripts
	for i := 0; i < 10; i++ {
		cmd := exec.Command("../../php-go", "run", "testdata/memory_intensive.php")
		cmd.Run()
	}

	runtime.GC() // Force garbage collection
	time.Sleep(100 * time.Millisecond)

	runtime.ReadMemStats(&memStats)
	after := memStats.Alloc

	growth := int64(after) - int64(before)
	t.Logf("Memory check: Before=%d MB, After=%d MB, Growth=%d MB",
		before/1024/1024, after/1024/1024, growth/1024/1024)
}

func runQuickConcurrentTest(t *testing.T) {
	const numConcurrent = 20
	var wg sync.WaitGroup
	var errorCount uint64

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.Command("../../php-go", "run", "testdata/concurrent_access.php")
			if cmd.Run() != nil {
				atomic.AddUint64(&errorCount, 1)
			}
		}()
	}

	wg.Wait()
	t.Logf("Quick concurrent test: %d errors", errorCount)
}
