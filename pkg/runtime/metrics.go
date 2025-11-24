package runtime

import (
	"fmt"
	"sync"
	"time"
)

// MetricType represents the type of metric
type MetricType int

const (
	MetricTypeCounter   MetricType = iota // Counter - monotonically increasing value
	MetricTypeGauge                       // Gauge - current value that can go up or down
	MetricTypeHistogram                   // Histogram - distribution of values
	MetricTypeTimer                       // Timer - duration tracking
)

// String returns the string representation of a metric type
func (t MetricType) String() string {
	switch t {
	case MetricTypeCounter:
		return "counter"
	case MetricTypeGauge:
		return "gauge"
	case MetricTypeHistogram:
		return "histogram"
	case MetricTypeTimer:
		return "timer"
	default:
		return "unknown"
	}
}

// Metric represents a single metric
type Metric struct {
	Name      string
	Type      MetricType
	Value     float64
	Count     int64 // For histograms and timers
	Sum       float64
	Min       float64
	Max       float64
	Labels    map[string]string
	Timestamp time.Time
}

// MetricsCollector collects and manages metrics
type MetricsCollector struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
	}
}

// metricKey generates a unique key for a metric based on name and labels
func metricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}

	key := name
	for k, v := range labels {
		key += fmt.Sprintf(",%s=%s", k, v)
	}
	return key
}

// Counter increments a counter metric by the specified value
func (c *MetricsCollector) Counter(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	metric, exists := c.metrics[key]

	if !exists {
		metric = &Metric{
			Name:      name,
			Type:      MetricTypeCounter,
			Value:     0,
			Labels:    labels,
			Timestamp: time.Now(),
		}
		c.metrics[key] = metric
	}

	metric.Value += value
	metric.Timestamp = time.Now()
}

// Inc increments a counter by 1
func (c *MetricsCollector) Inc(name string, labels map[string]string) {
	c.Counter(name, 1, labels)
}

// Gauge sets a gauge metric to the specified value
func (c *MetricsCollector) Gauge(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	metric, exists := c.metrics[key]

	if !exists {
		metric = &Metric{
			Name:      name,
			Type:      MetricTypeGauge,
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
		c.metrics[key] = metric
	} else {
		metric.Value = value
		metric.Timestamp = time.Now()
	}
}

// GaugeInc increments a gauge by the specified value
func (c *MetricsCollector) GaugeInc(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	metric, exists := c.metrics[key]

	if !exists {
		metric = &Metric{
			Name:      name,
			Type:      MetricTypeGauge,
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
		c.metrics[key] = metric
	} else {
		metric.Value += value
		metric.Timestamp = time.Now()
	}
}

// GaugeDec decrements a gauge by the specified value
func (c *MetricsCollector) GaugeDec(name string, value float64, labels map[string]string) {
	c.GaugeInc(name, -value, labels)
}

// Histogram records a value in a histogram
func (c *MetricsCollector) Histogram(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	metric, exists := c.metrics[key]

	if !exists {
		metric = &Metric{
			Name:      name,
			Type:      MetricTypeHistogram,
			Count:     0,
			Sum:       0,
			Min:       value,
			Max:       value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
		c.metrics[key] = metric
	}

	metric.Count++
	metric.Sum += value
	if value < metric.Min {
		metric.Min = value
	}
	if value > metric.Max {
		metric.Max = value
	}
	metric.Value = metric.Sum / float64(metric.Count) // Average
	metric.Timestamp = time.Now()
}

// Timer records a duration in a timer
func (c *MetricsCollector) Timer(name string, duration time.Duration, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	value := duration.Seconds()
	metric, exists := c.metrics[key]

	if !exists {
		metric = &Metric{
			Name:      name,
			Type:      MetricTypeTimer,
			Count:     0,
			Sum:       0,
			Min:       value,
			Max:       value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
		c.metrics[key] = metric
	}

	metric.Count++
	metric.Sum += value
	if value < metric.Min {
		metric.Min = value
	}
	if value > metric.Max {
		metric.Max = value
	}
	metric.Value = metric.Sum / float64(metric.Count) // Average
	metric.Timestamp = time.Now()
}

// TimerFunc measures the duration of a function execution
func (c *MetricsCollector) TimerFunc(name string, labels map[string]string, fn func()) {
	start := time.Now()
	fn()
	c.Timer(name, time.Since(start), labels)
}

// Get retrieves a metric by name and labels
func (c *MetricsCollector) Get(name string, labels map[string]string) (*Metric, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := metricKey(name, labels)
	metric, exists := c.metrics[key]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	copy := *metric
	return &copy, true
}

// GetAll returns all metrics
func (c *MetricsCollector) GetAll() []*Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make([]*Metric, 0, len(c.metrics))
	for _, metric := range c.metrics {
		// Return copies to avoid race conditions
		copy := *metric
		metrics = append(metrics, &copy)
	}
	return metrics
}

// Reset resets a specific metric
func (c *MetricsCollector) Reset(name string, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := metricKey(name, labels)
	delete(c.metrics, key)
}

// ResetAll resets all metrics
func (c *MetricsCollector) ResetAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = make(map[string]*Metric)
}

// MetricsSnapshot represents a snapshot of all metrics at a point in time
type MetricsSnapshot struct {
	Timestamp time.Time
	Metrics   []*Metric
}

// Snapshot creates a snapshot of all current metrics
func (c *MetricsCollector) Snapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		Timestamp: time.Now(),
		Metrics:   c.GetAll(),
	}
}

// FormatPrometheus formats metrics in Prometheus text format
func (s *MetricsSnapshot) FormatPrometheus() string {
	output := ""

	for _, metric := range s.Metrics {
		// Add type hint
		output += fmt.Sprintf("# TYPE %s %s\n", metric.Name, metric.Type.String())

		// Format labels
		labels := ""
		if len(metric.Labels) > 0 {
			labels = "{"
			first := true
			for k, v := range metric.Labels {
				if !first {
					labels += ","
				}
				labels += fmt.Sprintf(`%s="%s"`, k, v)
				first = false
			}
			labels += "}"
		}

		// Add metric value(s)
		switch metric.Type {
		case MetricTypeCounter, MetricTypeGauge:
			output += fmt.Sprintf("%s%s %g %d\n", metric.Name, labels, metric.Value, metric.Timestamp.UnixMilli())

		case MetricTypeHistogram, MetricTypeTimer:
			// For histograms and timers, output count, sum, min, max, and avg
			output += fmt.Sprintf("%s_count%s %d %d\n", metric.Name, labels, metric.Count, metric.Timestamp.UnixMilli())
			output += fmt.Sprintf("%s_sum%s %g %d\n", metric.Name, labels, metric.Sum, metric.Timestamp.UnixMilli())
			output += fmt.Sprintf("%s_min%s %g %d\n", metric.Name, labels, metric.Min, metric.Timestamp.UnixMilli())
			output += fmt.Sprintf("%s_max%s %g %d\n", metric.Name, labels, metric.Max, metric.Timestamp.UnixMilli())
			output += fmt.Sprintf("%s_avg%s %g %d\n", metric.Name, labels, metric.Value, metric.Timestamp.UnixMilli())
		}

		output += "\n"
	}

	return output
}

// FormatJSON formats metrics in JSON format
func (s *MetricsSnapshot) FormatJSON() string {
	output := "{\n"
	output += fmt.Sprintf(`  "timestamp": "%s",`, s.Timestamp.Format(time.RFC3339)) + "\n"
	output += `  "metrics": [` + "\n"

	for i, metric := range s.Metrics {
		if i > 0 {
			output += ",\n"
		}

		output += "    {\n"
		output += fmt.Sprintf(`      "name": "%s",`, metric.Name) + "\n"
		output += fmt.Sprintf(`      "type": "%s",`, metric.Type.String()) + "\n"
		output += fmt.Sprintf(`      "value": %g,`, metric.Value) + "\n"

		if metric.Type == MetricTypeHistogram || metric.Type == MetricTypeTimer {
			output += fmt.Sprintf(`      "count": %d,`, metric.Count) + "\n"
			output += fmt.Sprintf(`      "sum": %g,`, metric.Sum) + "\n"
			output += fmt.Sprintf(`      "min": %g,`, metric.Min) + "\n"
			output += fmt.Sprintf(`      "max": %g,`, metric.Max) + "\n"
		}

		if len(metric.Labels) > 0 {
			output += `      "labels": {` + "\n"
			first := true
			for k, v := range metric.Labels {
				if !first {
					output += ",\n"
				}
				output += fmt.Sprintf(`        "%s": "%s"`, k, v)
				first = false
			}
			output += "\n      },\n"
		}

		output += fmt.Sprintf(`      "timestamp": "%s"`, metric.Timestamp.Format(time.RFC3339)) + "\n"
		output += "    }"
	}

	output += "\n  ]\n"
	output += "}\n"

	return output
}

// Global metrics collector
var globalMetrics = NewMetricsCollector()

// GetGlobalMetrics returns the global metrics collector
func GetGlobalMetrics() *MetricsCollector {
	return globalMetrics
}

// SetGlobalMetrics sets the global metrics collector
func SetGlobalMetrics(collector *MetricsCollector) {
	globalMetrics = collector
}

// Package-level convenience functions that use the global collector

// Counter increments a counter using the global collector
func Counter(name string, value float64, labels map[string]string) {
	globalMetrics.Counter(name, value, labels)
}

// Inc increments a counter by 1 using the global collector
func Inc(name string, labels map[string]string) {
	globalMetrics.Inc(name, labels)
}

// Gauge sets a gauge using the global collector
func Gauge(name string, value float64, labels map[string]string) {
	globalMetrics.Gauge(name, value, labels)
}

// GaugeInc increments a gauge using the global collector
func GaugeInc(name string, value float64, labels map[string]string) {
	globalMetrics.GaugeInc(name, value, labels)
}

// GaugeDec decrements a gauge using the global collector
func GaugeDec(name string, value float64, labels map[string]string) {
	globalMetrics.GaugeDec(name, value, labels)
}

// Histogram records a histogram value using the global collector
func Histogram(name string, value float64, labels map[string]string) {
	globalMetrics.Histogram(name, value, labels)
}

// Timer records a timer duration using the global collector
func Timer(name string, duration time.Duration, labels map[string]string) {
	globalMetrics.Timer(name, duration, labels)
}

// TimerFunc measures function duration using the global collector
func TimerFunc(name string, labels map[string]string, fn func()) {
	globalMetrics.TimerFunc(name, labels, fn)
}

// CreateMetricsSnapshot creates a snapshot using the global collector
func CreateMetricsSnapshot() *MetricsSnapshot {
	return globalMetrics.Snapshot()
}
