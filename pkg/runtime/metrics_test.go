package runtime

import (
	"strings"
	"testing"
	"time"
)

// TestMetricType tests metric type string representation
func TestMetricType(t *testing.T) {
	tests := []struct {
		metricType MetricType
		expected   string
	}{
		{MetricTypeCounter, "counter"},
		{MetricTypeGauge, "gauge"},
		{MetricTypeHistogram, "histogram"},
		{MetricTypeTimer, "timer"},
		{MetricType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.metricType.String(); got != tt.expected {
				t.Errorf("MetricType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestMetricKey tests metric key generation
func TestMetricKey(t *testing.T) {
	tests := []struct {
		name     string
		metricName string
		labels   map[string]string
		contains []string
	}{
		{
			name:       "no labels",
			metricName: "test_metric",
			labels:     nil,
			contains:   []string{"test_metric"},
		},
		{
			name:       "with labels",
			metricName: "test_metric",
			labels:     map[string]string{"env": "prod", "host": "server1"},
			contains:   []string{"test_metric", "env=prod", "host=server1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := metricKey(tt.metricName, tt.labels)
			for _, substr := range tt.contains {
				if !strings.Contains(key, substr) {
					t.Errorf("metricKey() = %v, should contain %v", key, substr)
				}
			}
		})
	}
}

// TestCounter tests counter metric
func TestCounter(t *testing.T) {
	collector := NewMetricsCollector()

	// Test basic counter increment
	collector.Counter("requests", 1, nil)
	collector.Counter("requests", 2, nil)
	collector.Counter("requests", 3, nil)

	metric, exists := collector.Get("requests", nil)
	if !exists {
		t.Fatal("Counter metric not found")
	}

	if metric.Type != MetricTypeCounter {
		t.Errorf("Expected counter type, got %v", metric.Type)
	}

	if metric.Value != 6 {
		t.Errorf("Expected counter value 6, got %v", metric.Value)
	}

	if metric.Name != "requests" {
		t.Errorf("Expected metric name 'requests', got %v", metric.Name)
	}
}

// TestInc tests counter increment by 1
func TestInc(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Inc("events", nil)
	collector.Inc("events", nil)
	collector.Inc("events", nil)

	metric, exists := collector.Get("events", nil)
	if !exists {
		t.Fatal("Counter metric not found")
	}

	if metric.Value != 3 {
		t.Errorf("Expected counter value 3, got %v", metric.Value)
	}
}

// TestCounterWithLabels tests counter with labels
func TestCounterWithLabels(t *testing.T) {
	collector := NewMetricsCollector()

	labels1 := map[string]string{"status": "200"}
	labels2 := map[string]string{"status": "404"}

	collector.Counter("http_requests", 5, labels1)
	collector.Counter("http_requests", 2, labels2)

	metric1, exists := collector.Get("http_requests", labels1)
	if !exists {
		t.Fatal("Counter metric with labels1 not found")
	}

	metric2, exists := collector.Get("http_requests", labels2)
	if !exists {
		t.Fatal("Counter metric with labels2 not found")
	}

	if metric1.Value != 5 {
		t.Errorf("Expected metric1 value 5, got %v", metric1.Value)
	}

	if metric2.Value != 2 {
		t.Errorf("Expected metric2 value 2, got %v", metric2.Value)
	}
}

// TestGauge tests gauge metric
func TestGauge(t *testing.T) {
	collector := NewMetricsCollector()

	// Test gauge set
	collector.Gauge("temperature", 20.5, nil)
	metric, _ := collector.Get("temperature", nil)
	if metric.Value != 20.5 {
		t.Errorf("Expected gauge value 20.5, got %v", metric.Value)
	}

	// Test gauge update
	collector.Gauge("temperature", 25.3, nil)
	metric, _ = collector.Get("temperature", nil)
	if metric.Value != 25.3 {
		t.Errorf("Expected gauge value 25.3, got %v", metric.Value)
	}

	if metric.Type != MetricTypeGauge {
		t.Errorf("Expected gauge type, got %v", metric.Type)
	}
}

// TestGaugeInc tests gauge increment
func TestGaugeInc(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Gauge("memory", 100, nil)
	collector.GaugeInc("memory", 50, nil)
	collector.GaugeInc("memory", 25, nil)

	metric, exists := collector.Get("memory", nil)
	if !exists {
		t.Fatal("Gauge metric not found")
	}

	if metric.Value != 175 {
		t.Errorf("Expected gauge value 175, got %v", metric.Value)
	}
}

// TestGaugeDec tests gauge decrement
func TestGaugeDec(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Gauge("connections", 100, nil)
	collector.GaugeDec("connections", 30, nil)
	collector.GaugeDec("connections", 20, nil)

	metric, exists := collector.Get("connections", nil)
	if !exists {
		t.Fatal("Gauge metric not found")
	}

	if metric.Value != 50 {
		t.Errorf("Expected gauge value 50, got %v", metric.Value)
	}
}

// TestHistogram tests histogram metric
func TestHistogram(t *testing.T) {
	collector := NewMetricsCollector()

	values := []float64{10, 20, 30, 40, 50}
	for _, v := range values {
		collector.Histogram("response_size", v, nil)
	}

	metric, exists := collector.Get("response_size", nil)
	if !exists {
		t.Fatal("Histogram metric not found")
	}

	if metric.Type != MetricTypeHistogram {
		t.Errorf("Expected histogram type, got %v", metric.Type)
	}

	if metric.Count != 5 {
		t.Errorf("Expected count 5, got %v", metric.Count)
	}

	if metric.Sum != 150 {
		t.Errorf("Expected sum 150, got %v", metric.Sum)
	}

	if metric.Min != 10 {
		t.Errorf("Expected min 10, got %v", metric.Min)
	}

	if metric.Max != 50 {
		t.Errorf("Expected max 50, got %v", metric.Max)
	}

	if metric.Value != 30 { // Average
		t.Errorf("Expected average 30, got %v", metric.Value)
	}
}

// TestTimer tests timer metric
func TestTimer(t *testing.T) {
	collector := NewMetricsCollector()

	durations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
	}

	for _, d := range durations {
		collector.Timer("request_duration", d, nil)
	}

	metric, exists := collector.Get("request_duration", nil)
	if !exists {
		t.Fatal("Timer metric not found")
	}

	if metric.Type != MetricTypeTimer {
		t.Errorf("Expected timer type, got %v", metric.Type)
	}

	if metric.Count != 3 {
		t.Errorf("Expected count 3, got %v", metric.Count)
	}

	// Sum should be 0.6 seconds (600ms)
	expectedSum := 0.6
	if metric.Sum < expectedSum-0.01 || metric.Sum > expectedSum+0.01 {
		t.Errorf("Expected sum around %v, got %v", expectedSum, metric.Sum)
	}

	// Average should be 0.2 seconds (200ms)
	expectedAvg := 0.2
	if metric.Value < expectedAvg-0.01 || metric.Value > expectedAvg+0.01 {
		t.Errorf("Expected average around %v, got %v", expectedAvg, metric.Value)
	}
}

// TestTimerFunc tests timer function wrapper
func TestTimerFunc(t *testing.T) {
	collector := NewMetricsCollector()

	// Measure a simple function
	collector.TimerFunc("function_duration", nil, func() {
		time.Sleep(50 * time.Millisecond)
	})

	metric, exists := collector.Get("function_duration", nil)
	if !exists {
		t.Fatal("Timer metric not found")
	}

	if metric.Count != 1 {
		t.Errorf("Expected count 1, got %v", metric.Count)
	}

	// Should be around 0.05 seconds (50ms)
	if metric.Value < 0.04 || metric.Value > 0.1 {
		t.Errorf("Expected value around 0.05, got %v", metric.Value)
	}
}

// TestGetAll tests retrieving all metrics
func TestGetAll(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Counter("metric1", 1, nil)
	collector.Gauge("metric2", 2, nil)
	collector.Histogram("metric3", 3, nil)

	metrics := collector.GetAll()

	if len(metrics) != 3 {
		t.Errorf("Expected 3 metrics, got %v", len(metrics))
	}

	// Check that all metrics are present
	names := make(map[string]bool)
	for _, m := range metrics {
		names[m.Name] = true
	}

	if !names["metric1"] || !names["metric2"] || !names["metric3"] {
		t.Error("Not all expected metrics found")
	}
}

// TestReset tests resetting a specific metric
func TestReset(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Counter("reset_test", 10, nil)

	_, exists := collector.Get("reset_test", nil)
	if !exists {
		t.Fatal("Metric should exist before reset")
	}

	collector.Reset("reset_test", nil)

	_, exists = collector.Get("reset_test", nil)
	if exists {
		t.Error("Metric should not exist after reset")
	}
}

// TestResetAll tests resetting all metrics
func TestResetAll(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Counter("metric1", 1, nil)
	collector.Gauge("metric2", 2, nil)
	collector.Histogram("metric3", 3, nil)

	if len(collector.GetAll()) != 3 {
		t.Fatal("Expected 3 metrics before reset")
	}

	collector.ResetAll()

	if len(collector.GetAll()) != 0 {
		t.Error("Expected 0 metrics after reset")
	}
}

// TestSnapshot tests creating a metrics snapshot
func TestSnapshot(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Counter("requests", 100, nil)
	collector.Gauge("memory", 50, nil)

	snapshot := collector.Snapshot()

	if snapshot == nil {
		t.Fatal("Snapshot should not be nil")
	}

	if len(snapshot.Metrics) != 2 {
		t.Errorf("Expected 2 metrics in snapshot, got %v", len(snapshot.Metrics))
	}

	if snapshot.Timestamp.IsZero() {
		t.Error("Snapshot timestamp should not be zero")
	}
}

// TestFormatPrometheus tests Prometheus format output
func TestFormatPrometheus(t *testing.T) {
	collector := NewMetricsCollector()

	labels := map[string]string{"env": "test"}
	collector.Counter("http_requests_total", 42, labels)
	collector.Gauge("cpu_usage", 75.5, nil)
	collector.Histogram("request_size", 100, nil)
	collector.Histogram("request_size", 200, nil)

	snapshot := collector.Snapshot()
	output := snapshot.FormatPrometheus()

	// Check for expected content
	expectedStrings := []string{
		"# TYPE http_requests_total counter",
		`http_requests_total{env="test"}`,
		"# TYPE cpu_usage gauge",
		"cpu_usage",
		"75.5",
		"# TYPE request_size histogram",
		"request_size_count",
		"request_size_sum",
		"request_size_min",
		"request_size_max",
		"request_size_avg",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Prometheus output missing expected string: %v", expected)
		}
	}
}

// TestFormatJSON tests JSON format output
func TestFormatJSON(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Counter("requests", 10, nil)
	collector.Gauge("temperature", 25.5, map[string]string{"sensor": "1"})

	snapshot := collector.Snapshot()
	output := snapshot.FormatJSON()

	// Check for expected JSON structure
	expectedStrings := []string{
		`"timestamp"`,
		`"metrics"`,
		`"name"`,
		`"type"`,
		`"value"`,
		`"requests"`,
		`"counter"`,
		`"temperature"`,
		`"gauge"`,
		`"labels"`,
		`"sensor"`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("JSON output missing expected string: %v", expected)
		}
	}
}

// TestConcurrency tests concurrent metric updates
func TestConcurrency(t *testing.T) {
	collector := NewMetricsCollector()

	done := make(chan bool)
	iterations := 1000

	// Start multiple goroutines updating the same counter
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < iterations; j++ {
				collector.Inc("concurrent_counter", nil)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	metric, exists := collector.Get("concurrent_counter", nil)
	if !exists {
		t.Fatal("Concurrent counter metric not found")
	}

	expected := float64(10 * iterations)
	if metric.Value != expected {
		t.Errorf("Expected counter value %v, got %v", expected, metric.Value)
	}
}

// TestGlobalMetrics tests global metrics functions
func TestGlobalMetrics(t *testing.T) {
	// Reset global metrics
	SetGlobalMetrics(NewMetricsCollector())

	Counter("global_counter", 5, nil)
	Inc("global_inc", nil)
	Gauge("global_gauge", 100, nil)
	GaugeInc("global_gauge_inc", 10, nil)
	GaugeDec("global_gauge_dec", 5, map[string]string{"test": "true"})
	Histogram("global_histogram", 50, nil)
	Timer("global_timer", 100*time.Millisecond, nil)

	metrics := GetGlobalMetrics().GetAll()

	if len(metrics) != 7 {
		t.Errorf("Expected 7 global metrics, got %v", len(metrics))
	}
}

// TestGlobalTimerFunc tests global timer function wrapper
func TestGlobalTimerFunc(t *testing.T) {
	// Reset global metrics
	SetGlobalMetrics(NewMetricsCollector())

	TimerFunc("global_timer_func", nil, func() {
		time.Sleep(10 * time.Millisecond)
	})

	metric, exists := GetGlobalMetrics().Get("global_timer_func", nil)
	if !exists {
		t.Fatal("Global timer func metric not found")
	}

	if metric.Count != 1 {
		t.Errorf("Expected count 1, got %v", metric.Count)
	}
}

// TestMetricTimestamps tests that timestamps are updated correctly
func TestMetricTimestamps(t *testing.T) {
	collector := NewMetricsCollector()

	before := time.Now()
	time.Sleep(10 * time.Millisecond)

	collector.Counter("timestamp_test", 1, nil)

	metric, exists := collector.Get("timestamp_test", nil)
	if !exists {
		t.Fatal("Metric not found")
	}

	if metric.Timestamp.Before(before) {
		t.Error("Metric timestamp should be after start time")
	}

	if metric.Timestamp.After(time.Now()) {
		t.Error("Metric timestamp should not be in the future")
	}
}

// BenchmarkCounter benchmarks counter operations
func BenchmarkCounter(b *testing.B) {
	collector := NewMetricsCollector()
	labels := map[string]string{"benchmark": "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Counter("bench_counter", 1, labels)
	}
}

// BenchmarkGauge benchmarks gauge operations
func BenchmarkGauge(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Gauge("bench_gauge", float64(i), nil)
	}
}

// BenchmarkHistogram benchmarks histogram operations
func BenchmarkHistogram(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Histogram("bench_histogram", float64(i), nil)
	}
}

// BenchmarkTimer benchmarks timer operations
func BenchmarkTimer(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Timer("bench_timer", time.Duration(i)*time.Millisecond, nil)
	}
}

// BenchmarkSnapshot benchmarks snapshot creation
func BenchmarkSnapshot(b *testing.B) {
	collector := NewMetricsCollector()

	// Add some metrics
	for i := 0; i < 100; i++ {
		collector.Counter("counter", 1, map[string]string{"id": string(rune(i))})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collector.Snapshot()
	}
}

// BenchmarkFormatPrometheus benchmarks Prometheus format generation
func BenchmarkFormatPrometheus(b *testing.B) {
	collector := NewMetricsCollector()

	// Add various metrics
	collector.Counter("requests", 100, map[string]string{"method": "GET"})
	collector.Gauge("memory", 1024, nil)
	collector.Histogram("size", 500, nil)
	collector.Timer("duration", 100*time.Millisecond, nil)

	snapshot := collector.Snapshot()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = snapshot.FormatPrometheus()
	}
}

// BenchmarkFormatJSON benchmarks JSON format generation
func BenchmarkFormatJSON(b *testing.B) {
	collector := NewMetricsCollector()

	// Add various metrics
	collector.Counter("requests", 100, map[string]string{"method": "GET"})
	collector.Gauge("memory", 1024, nil)
	collector.Histogram("size", 500, nil)
	collector.Timer("duration", 100*time.Millisecond, nil)

	snapshot := collector.Snapshot()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = snapshot.FormatJSON()
	}
}
