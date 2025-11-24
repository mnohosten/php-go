package main

import (
	"fmt"
	"time"

	"github.com/krizos/php-go/pkg/runtime"
)

func main() {
	// Create a custom metrics collector
	collector := runtime.NewMetricsCollector()

	// Or use the global collector (simpler)
	// runtime.Counter("requests", 1, nil)

	// Example 1: Counter - tracking total requests
	labels := map[string]string{"method": "GET", "endpoint": "/api/users"}
	collector.Counter("http_requests_total", 1, labels)
	collector.Counter("http_requests_total", 1, labels)
	collector.Counter("http_requests_total", 1, labels)

	// Example 2: Counter with different labels
	labels2 := map[string]string{"method": "POST", "endpoint": "/api/users"}
	collector.Counter("http_requests_total", 2, labels2)

	// Example 3: Gauge - tracking current active connections
	collector.Gauge("active_connections", 42, nil)
	collector.GaugeInc("active_connections", 5, nil)  // Add 5 connections
	collector.GaugeDec("active_connections", 3, nil)  // Remove 3 connections

	// Example 4: Histogram - tracking response sizes
	responseSizes := []float64{256, 512, 1024, 2048, 4096}
	for _, size := range responseSizes {
		collector.Histogram("response_size_bytes", size, map[string]string{"type": "json"})
	}

	// Example 5: Timer - tracking request durations
	durations := []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		75 * time.Millisecond,
	}
	for _, duration := range durations {
		collector.Timer("request_duration_seconds", duration, map[string]string{"handler": "user_list"})
	}

	// Example 6: TimerFunc - measure function execution time
	collector.TimerFunc("database_query_duration", map[string]string{"query": "select_users"}, func() {
		// Simulate database query
		time.Sleep(25 * time.Millisecond)
	})

	// Create a snapshot of all metrics
	snapshot := collector.Snapshot()

	// Output metrics in Prometheus format
	fmt.Println("=== Prometheus Format ===")
	fmt.Println(snapshot.FormatPrometheus())

	// Output metrics in JSON format
	fmt.Println("\n=== JSON Format ===")
	fmt.Println(snapshot.FormatJSON())

	// Example 7: Access individual metrics
	fmt.Println("\n=== Individual Metric Access ===")
	if metric, exists := collector.Get("active_connections", nil); exists {
		fmt.Printf("Active connections: %.0f\n", metric.Value)
	}

	// Example 8: Global metrics collector
	fmt.Println("\n=== Using Global Metrics ===")
	runtime.Counter("app_errors_total", 1, map[string]string{"severity": "high"})
	runtime.Gauge("memory_usage_bytes", 1024*1024*256, nil) // 256 MB
	runtime.Inc("cache_hits", map[string]string{"cache": "redis"})

	globalSnapshot := runtime.CreateMetricsSnapshot()
	fmt.Printf("Global metrics count: %d\n", len(globalSnapshot.Metrics))
}
