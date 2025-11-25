# PHP-Go Stress Testing

This directory contains stress tests for PHP-Go, focusing on load testing, memory leak detection, concurrent request testing, and long-running process testing.

## Overview

Stress tests verify that PHP-Go can handle production workloads by testing:
- **Load Testing**: Sustained high request rates
- **Memory Leak Detection**: Memory stability over time
- **Concurrent Request Testing**: Multiple simultaneous requests
- **Long-Running Process Testing**: Extended execution periods

## Load Testing

### Running Load Tests

```bash
# Run all load tests (takes ~2-3 minutes)
go test -v ./tests/stress/

# Run with short mode (skips load tests)
go test -short ./tests/stress/

# Run specific test
go test -v ./tests/stress/ -run TestLoadSimpleScript

# Run with custom timeout
go test -v ./tests/stress/ -timeout 10m
```

### Load Test Configuration

Each load test uses a `LoadTestConfig` struct:

```go
type LoadTestConfig struct {
    Duration       time.Duration // How long to run the test
    ConcurrentVMs  int           // Number of concurrent VM instances
    RequestsPerVM  int           // Requests per VM (0 = unlimited)
    RampUpTime     time.Duration // Time to gradually increase load
    ThinkTime      time.Duration // Delay between requests per VM
    MaxMemoryMB    int           // Maximum allowed memory in MB
    MaxGoroutines  int           // Maximum allowed goroutines
}
```

### Available Load Tests

1. **TestLoadSimpleScript**
   - Tests basic arithmetic operations under load
   - 10 concurrent VMs, 10 seconds duration
   - Expected throughput: >100 req/sec
   - Memory limit: 500 MB

2. **TestLoadFunctionCalls**
   - Tests recursive function calls (Fibonacci)
   - 20 concurrent VMs, 10 seconds duration
   - Expected throughput: >50 req/sec
   - Tests function call overhead and stack management

3. **TestLoadArrayOperations**
   - Tests array creation, iteration, and manipulation
   - 15 concurrent VMs, 10 seconds duration
   - Tests memory allocation patterns for arrays

4. **TestLoadStringOperations**
   - Tests string concatenation and operations
   - 15 concurrent VMs, 10 seconds duration
   - Tests string memory management

5. **TestLoadHighConcurrency**
   - Tests with very high concurrency (100 VMs)
   - 15 seconds duration with 5-second ramp-up
   - Expected throughput: >500 req/sec
   - Tests goroutine management and synchronization

6. **TestLoadSustainedTraffic**
   - Tests sustained load over 30 seconds
   - 20 concurrent VMs with 1ms think time
   - Memory growth should be <50%
   - Detects memory leaks and resource exhaustion

### Load Test Results

Each test outputs comprehensive metrics:

```
=== Load Test Results: Simple Script ===
Duration: 10s
Total Requests: 15234
  Successful: 15234 (100.00%)
  Failed: 0 (0.00%)
Throughput: 1523.40 req/sec

Latency:
  Min: 245µs
  Avg: 654µs
  Max: 12.5ms
  P50: 612µs
  P95: 1.2ms
  P99: 2.8ms

Memory (MB):
  Initial: 5.23
  Final: 8.45
  Peak: 9.12
  Growth: 3.22 (61.57%)

Goroutines:
  Initial: 4
  Final: 14
  Peak: 24
```

### Understanding Metrics

- **Throughput**: Requests per second (higher is better)
- **Latency**: Response time distribution
  - **Avg**: Mean latency across all requests
  - **P50**: 50% of requests complete faster than this
  - **P95**: 95% of requests complete faster than this
  - **P99**: 99% of requests complete faster than this
- **Memory**: Memory consumption in MB
  - **Growth**: Increase from start to end (should be minimal)
  - **Peak**: Maximum memory during test
- **Goroutines**: Go routines active
  - Should return to near-initial levels after test

### Performance Baselines

Expected performance on modern hardware (Apple M4 Max or equivalent):

| Test Scenario          | Min Throughput | Max Latency (P99) | Max Memory Growth |
|------------------------|----------------|-------------------|-------------------|
| Simple Script          | >100 req/sec   | <5ms              | <100 MB           |
| Function Calls         | >50 req/sec    | <10ms             | <150 MB           |
| Array Operations       | >80 req/sec    | <8ms              | <200 MB           |
| String Operations      | >60 req/sec    | <10ms             | <250 MB           |
| High Concurrency       | >500 req/sec   | <20ms             | <500 MB           |
| Sustained Traffic      | >100 req/sec   | <10ms             | <50% growth       |

### Interpreting Failures

**High Failure Rate**
- Check error messages in test output
- May indicate VM bugs or resource exhaustion
- Review stack traces and error patterns

**Low Throughput**
- Compare against baselines above
- May indicate performance regression
- Run CPU profiling: `go test -cpuprofile=cpu.prof`

**High Memory Growth**
- >100% growth suggests memory leak
- Run memory profiling: `go test -memprofile=mem.prof`
- Check for missing Value.Release() calls

**Goroutine Leak**
- Final goroutines >> initial goroutines
- Indicates goroutines not being cleaned up
- Check for missing channel closes or WaitGroup.Done()

## Benchmarks

```bash
# Run load test benchmarks
go test -bench=BenchmarkLoadTest -benchmem ./tests/stress/

# Compare parallel vs sequential
go test -bench=BenchmarkLoadTest -benchmem -benchtime=10s ./tests/stress/
```

### Benchmark Results

Expected benchmark results (Apple M4 Max):

```
BenchmarkLoadTest-16              50000    24156 ns/op    8432 B/op    156 allocs/op
BenchmarkLoadTestParallel-16     500000     2847 ns/op    8432 B/op    156 allocs/op
```

The parallel benchmark should be ~8-10x faster, showing good scalability.

## Custom Load Tests

You can create custom load tests for specific scenarios:

```go
func TestLoadCustomScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }

    script := `<?php
    // Your PHP script here
    `

    config := LoadTestConfig{
        Duration:      20 * time.Second,
        ConcurrentVMs: 25,
        RequestsPerVM: 0, // Unlimited
        RampUpTime:    3 * time.Second,
        ThinkTime:     0,
        MaxMemoryMB:   500,
        MaxGoroutines: 1000,
    }

    result := RunLoadTest(t, "Custom Scenario", script, config)
    PrintLoadTestResult(t, "Custom Scenario", result)

    // Custom assertions
    if result.RequestsPerSecond < 100 {
        t.Errorf("Throughput too low: %.2f req/sec", result.RequestsPerSecond)
    }
}
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Stress Tests

on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly
  workflow_dispatch:      # Manual trigger

jobs:
  stress-test:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Run Load Tests
        run: go test -v -timeout 20m ./tests/stress/

      - name: Upload Results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: stress-test-results
          path: tests/stress/*.log
```

### Performance Regression Detection

Run load tests before and after changes:

```bash
# Before changes
go test -bench=BenchmarkLoadTest -benchmem ./tests/stress/ > before.txt

# Make changes...

# After changes
go test -bench=BenchmarkLoadTest -benchmem ./tests/stress/ > after.txt

# Compare
go install golang.org/x/perf/cmd/benchstat@latest
benchstat before.txt after.txt
```

## Troubleshooting

### Tests Timeout

Increase test timeout:
```bash
go test -v -timeout 30m ./tests/stress/
```

### Out of Memory

Reduce concurrent VMs or add more RAM:
```go
config := LoadTestConfig{
    ConcurrentVMs: 5,  // Reduced from 20
    // ...
}
```

### High CPU Usage

Add think time between requests:
```go
config := LoadTestConfig{
    ThinkTime: 10 * time.Millisecond,
    // ...
}
```

### Inconsistent Results

Run multiple times and average:
```bash
for i in {1..5}; do
    go test -v ./tests/stress/ -run TestLoadSimpleScript
done
```

## Future Enhancements

Planned additions to stress testing:

1. **Memory Leak Detection** (Phase 10.11)
   - Heap profiling over time
   - Automatic leak detection
   - Memory growth analysis

2. **Concurrent Request Testing** (Phase 10.11)
   - Request interleaving tests
   - Race condition detection
   - Shared state testing

3. **Long-Running Process Testing** (Phase 10.11)
   - Multi-hour execution tests
   - Resource exhaustion detection
   - Graceful degradation testing

4. **Advanced Metrics**
   - Response time percentiles (P90, P99.9)
   - Error rate over time
   - Resource utilization graphs
   - Throughput over time charts

5. **Distributed Load Testing**
   - Multiple load generators
   - Network latency simulation
   - Geographic distribution

## References

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Benchmarking](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Load Testing Best Practices](https://www.nginx.com/blog/load-testing-best-practices/)
- [PHP-Go Performance Tuning Guide](../../docs/user-guide/performance-tuning.md)

## Contributing

When adding new load tests:

1. Follow the naming convention: `TestLoad<Scenario>`
2. Use `testing.Short()` to skip in short mode
3. Document expected baselines in this README
4. Include meaningful assertions
5. Add comments explaining the test scenario
6. Consider edge cases (empty input, large data, etc.)

## Contact

For questions or issues with stress testing:
- Open an issue on GitHub
- Check the [User Guide](../../docs/user-guide/)
- Review [Performance Tuning Guide](../../docs/user-guide/performance-tuning.md)
