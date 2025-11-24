# PHP-Go Configuration Guide

This guide covers all configuration options available in PHP-Go, from command-line arguments to runtime resource limits, logging, metrics, and extensions.

## Table of Contents

- [Command-Line Options](#command-line-options)
- [Resource Limits](#resource-limits)
- [Health Checks](#health-checks)
- [Logging Configuration](#logging-configuration)
- [Metrics Collection](#metrics-collection)
- [Extension System](#extension-system)
- [Plugin System](#plugin-system)
- [Runtime Configuration](#runtime-configuration)
- [Environment Variables](#environment-variables)
- [Best Practices](#best-practices)

## Command-Line Options

### Basic Usage

```bash
# Execute a PHP file
php-go <file>
php-go run <file>

# Show version
php-go --version
php-go -v

# Show help
php-go --help
php-go -h
```

### Development Commands

```bash
# Tokenize a PHP file (show lexer output)
php-go lex <file>
php-go lex --json <file>      # Output as JSON

# Parse a PHP file (show AST)
php-go parse <file>
php-go parse --json <file>    # Output as JSON

# Show parallelization demo
php-go demo
```

### Future Commands (Coming Soon)

```bash
# Interactive REPL mode
php-go -a

# Built-in web server
php-go -S host:port
```

## Resource Limits

PHP-Go provides comprehensive resource limiting to prevent runaway scripts and ensure system stability.

### Memory Limit

Control maximum memory usage:

```go
import "github.com/krizos/php-go/pkg/runtime"

// Set memory limit to 128MB
runtime.SetMemoryLimit(128 * 1024 * 1024)

// Get current limit
limit := runtime.GetResourceLimiter().GetMemoryLimit()

// Unlimited (0 = default)
runtime.SetMemoryLimit(0)
```

### Execution Time Limit

Control maximum script execution time:

```go
// Set execution time limit to 30 seconds
runtime.SetExecutionTimeLimit(30)

// Get current limit in seconds
limit := runtime.GetResourceLimiter().GetExecutionTimeLimit()

// Unlimited (0 = default)
runtime.SetExecutionTimeLimit(0)
```

### Recursion Depth Limit

Control maximum function call depth:

```go
// Set maximum recursion depth to 500
runtime.SetMaxRecursionDepth(500)

// Get current limit
limit := runtime.GetResourceLimiter().GetMaxRecursionDepth()

// Default is 1000 (same as PHP)
```

### Instruction Count Limit

Control maximum number of VM instructions:

```go
// Set maximum instruction count to 1 million
runtime.SetMaxInstructions(1000000)

// Get current limit
limit := runtime.GetResourceLimiter().GetMaxInstructions()

// Unlimited (0 = default)
```

### Output Buffer Size Limit

Control maximum output buffer size:

```go
// Set maximum output buffer size to 10MB
runtime.SetMaxOutputSize(10 * 1024 * 1024)

// Get current limit
limit := runtime.GetResourceLimiter().GetMaxOutputSize()

// Unlimited (0 = default)
```

### Resource Tracking

```go
// Start tracking resource usage
runtime.StartResourceTracking()

// Check limits during execution
if err := runtime.CheckResourceLimits(); err != nil {
    // Handle limit exceeded
    log.Printf("Resource limit exceeded: %v", err)
}

// Get current resource statistics
stats := runtime.GetResourceStats()
fmt.Printf("Memory: %v\n", stats["memory"])
fmt.Printf("Execution time: %v\n", stats["execution_time"])
fmt.Printf("Recursion depth: %v\n", stats["recursion_depth"])
fmt.Printf("Instructions: %v\n", stats["instructions"])
fmt.Printf("Output size: %v\n", stats["output_size"])

// Stop tracking
runtime.StopResourceTracking()
```

### Advanced Resource Configuration

```go
limiter := runtime.GetResourceLimiter()

// Enable/disable limit checking
limiter.SetEnabled(true)

// Set callbacks for limit violations
limiter.SetOnMemoryLimitExceeded(func(current, limit int64) {
    log.Printf("Memory limit exceeded: %d/%d bytes", current, limit)
})

limiter.SetOnExecutionTimeLimitExceeded(func(elapsed, limit time.Duration) {
    log.Printf("Execution time limit exceeded: %v/%v", elapsed, limit)
})

limiter.SetOnRecursionDepthLimitExceeded(func(current, limit int64) {
    log.Printf("Recursion depth limit exceeded: %d/%d", current, limit)
})

limiter.SetOnInstructionLimitExceeded(func(current, limit int64) {
    log.Printf("Instruction limit exceeded: %d/%d", current, limit)
})

limiter.SetOnOutputSizeLimitExceeded(func(current, limit int64) {
    log.Printf("Output size limit exceeded: %d/%d bytes", current, limit)
})

// Reset all limits to defaults
limiter.Reset()
```

## Health Checks

PHP-Go includes a health check system for monitoring application health in production.

### Built-in Health Checks

```go
import (
    "context"
    "github.com/krizos/php-go/pkg/runtime"
)

// Get health status
ctx := context.Background()
status := runtime.GetHealthStatus(ctx)
fmt.Printf("Health status: %s\n", status) // healthy, degraded, or unhealthy

// Get detailed health report
report := runtime.GetHealthReport(ctx)
fmt.Printf("Status: %s\n", report.Status)
fmt.Printf("Timestamp: %s\n", report.Timestamp)
for name, result := range report.Checks {
    fmt.Printf("  %s: %s - %s\n", name, result.Status, result.Message)
}
```

### Registering Custom Health Checks

```go
// Register a liveness check (is the app running?)
runtime.RegisterHealthCheck("my-liveness",
    runtime.HealthCheckTypeLiveness,
    func(ctx context.Context) (runtime.HealthStatus, string, error) {
        // Your check logic here
        return runtime.HealthStatusHealthy, "Service is alive", nil
    },
    5 * time.Second, // timeout
)

// Register a readiness check (can the app accept traffic?)
runtime.RegisterHealthCheck("database-ready",
    runtime.HealthCheckTypeReadiness,
    func(ctx context.Context) (runtime.HealthStatus, string, error) {
        // Check database connection
        if databaseConnected() {
            return runtime.HealthStatusHealthy, "Database connected", nil
        }
        return runtime.HealthStatusUnhealthy, "Database unavailable", nil
    },
    10 * time.Second,
)

// Register a startup check (has the app finished starting?)
runtime.RegisterHealthCheck("initialization",
    runtime.HealthCheckTypeStartup,
    func(ctx context.Context) (runtime.HealthStatus, string, error) {
        // Check if initialization is complete
        if isInitialized() {
            return runtime.HealthStatusHealthy, "Initialized", nil
        }
        return runtime.HealthStatusDegraded, "Still initializing", nil
    },
    30 * time.Second,
)
```

### Advanced Health Check Operations

```go
checker := runtime.GetGlobalHealthChecker()

// Run a specific health check
result, err := runtime.HealthCheckStatus(ctx, "database-ready")
if err != nil {
    log.Printf("Check failed: %v", err)
}
fmt.Printf("Status: %s, Message: %s\n", result.Status, result.Message)

// Run all checks of a specific type
readinessReport := checker.ReportByType(ctx, runtime.HealthCheckTypeReadiness)
fmt.Printf("Readiness status: %s\n", readinessReport.Status)

// Unregister a check
runtime.UnregisterHealthCheck("my-liveness")

// List all registered checks
checks := checker.List()
for _, name := range checks {
    fmt.Printf("Registered check: %s\n", name)
}
```

## Logging Configuration

PHP-Go provides a flexible logging system with multiple levels and output formats.

### Log Levels

```go
import "github.com/krizos/php-go/pkg/runtime"

// Available log levels (from most to least severe):
// - LogLevelFatal: Fatal errors that cause program termination
// - LogLevelError: Errors that prevent normal operation
// - LogLevelWarn:  Warning messages for potentially harmful situations
// - LogLevelInfo:  Informational messages about program execution
// - LogLevelDebug: Detailed debugging information
// - LogLevelTrace: Very detailed tracing information
```

### Basic Logging

```go
logger := runtime.NewLogger()

// Set log level
logger.SetLevel(runtime.LogLevelInfo)

// Log messages at different levels
logger.Fatal("Fatal error occurred")
logger.Error("An error occurred")
logger.Warn("Warning message")
logger.Info("Informational message")
logger.Debug("Debug information")
logger.Trace("Trace information")

// Log with fields (structured logging)
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
}).Info("User logged in")

// Log with component name
logger.WithComponent("auth").Info("Authentication successful")
```

### Log Formatters

#### Text Formatter

```go
// Create a text formatter
formatter := runtime.NewTextFormatter()
formatter.TimestampFormat = "2006-01-02 15:04:05"
formatter.IncludeLocation = true  // Include file:line
formatter.IncludeFields = true    // Include structured fields
formatter.ColorEnabled = false    // Enable/disable color output

logger.SetFormatter(formatter)
```

#### JSON Formatter

```go
// Create a JSON formatter for structured logging
jsonFormatter := runtime.NewJSONFormatter()
jsonFormatter.PrettyPrint = true      // Pretty-print JSON
jsonFormatter.IncludeLocation = true  // Include file:line
jsonFormatter.TimestampFormat = time.RFC3339

logger.SetFormatter(jsonFormatter)
```

### Log Outputs

```go
import "os"

// Log to stdout
logger.SetOutput(os.Stdout)

// Log to stderr
logger.SetOutput(os.Stderr)

// Log to file
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()
logger.SetOutput(file)

// Multiple outputs
logger.AddOutput(os.Stdout)
logger.AddOutput(file)
```

### Global Logger

```go
// Get global logger instance
logger := runtime.GetGlobalLogger()

// Configure global logger
logger.SetLevel(runtime.LogLevelDebug)
logger.SetFormatter(runtime.NewJSONFormatter())

// Use global logger convenience functions
runtime.LogInfo("Application started")
runtime.LogError("An error occurred")
runtime.LogDebug("Debug information")
```

## Metrics Collection

PHP-Go includes a metrics collection system for monitoring application performance.

### Metric Types

- **Counter**: Monotonically increasing value (e.g., request count)
- **Gauge**: Current value that can go up or down (e.g., active connections)
- **Histogram**: Distribution of values (e.g., request duration)
- **Timer**: Duration tracking (e.g., operation timing)

### Basic Metrics

```go
import "github.com/krizos/php-go/pkg/runtime"

collector := runtime.NewMetricsCollector()

// Counter - increment by value
collector.Counter("requests_total", 1, map[string]string{
    "method": "GET",
    "path": "/api/users",
})

// Inc - increment by 1 (convenience)
collector.Inc("http_requests", map[string]string{
    "status": "200",
})

// Gauge - set current value
collector.Gauge("active_connections", 42, nil)

// Histogram - record value distribution
collector.Histogram("request_duration_ms", 125.5, map[string]string{
    "endpoint": "/api/data",
})

// Timer - track operation duration
timer := collector.StartTimer("database_query", map[string]string{
    "query": "SELECT",
})
// ... perform operation ...
timer.Stop()
```

### Retrieving Metrics

```go
// Get all metrics
metrics := collector.GetAll()
for key, metric := range metrics {
    fmt.Printf("%s: %v\n", key, metric.Value)
}

// Get specific metric
metric, err := collector.Get("requests_total", map[string]string{
    "method": "GET",
})
if err == nil {
    fmt.Printf("Requests: %v\n", metric.Value)
}

// Get metrics by type
counters := collector.GetByType(runtime.MetricTypeCounter)
for _, metric := range counters {
    fmt.Printf("%s: %v\n", metric.Name, metric.Value)
}

// Get metrics by name prefix
apiMetrics := collector.GetByPrefix("api_")
```

### Global Metrics Collector

```go
// Get global metrics collector
collector := runtime.GetGlobalMetricsCollector()

// Use global convenience functions
runtime.IncrementCounter("requests", nil)
runtime.SetGauge("memory_usage_bytes", 1024*1024*50, nil)
runtime.RecordHistogram("response_time_ms", 45.2, nil)

// Reset all metrics
collector.Reset()

// Clear specific metric
collector.Clear("requests_total", nil)
```

## Extension System

PHP-Go supports loading Go-based extensions to add functionality.

### Loading Extensions

```go
import "github.com/krizos/php-go/pkg/goext"

// Get extension manager
mgr := goext.GetExtensionManager()

// Create an extension
ext := &goext.Extension{
    Name:        "my_extension",
    Version:     "1.0.0",
    Description: "My custom extension",
    Functions:   make(map[string]goext.PHPFunction),
    Classes:     make(map[string]*goext.PHPClass),
}

// Register functions
ext.RegisterFunction("my_func", func(args []types.Value) (types.Value, error) {
    // Your function implementation
    return types.NewString("Hello from extension"), nil
})

// Register extension
if err := mgr.Register(ext); err != nil {
    log.Fatalf("Failed to register extension: %v", err)
}

// List loaded extensions
extensions := mgr.List()
for _, name := range extensions {
    info := mgr.GetInfo(name)
    fmt.Printf("%s v%s: %s\n", info.Name, info.Version, info.Description)
}
```

### Extension Configuration

```go
// Enable/disable extension
mgr.Enable("my_extension")
mgr.Disable("my_extension")

// Check if extension is enabled
if mgr.IsEnabled("my_extension") {
    fmt.Println("Extension is enabled")
}

// Unregister extension
mgr.Unregister("my_extension")
```

## Plugin System

PHP-Go supports loading extensions from compiled Go plugins (.so files on Linux/macOS).

### Loading Plugins

```go
import "github.com/krizos/php-go/pkg/goext"

// Get plugin manager
pluginMgr := goext.GetPluginManager()

// Load a plugin
if err := pluginMgr.Load("/path/to/plugin.so"); err != nil {
    log.Fatalf("Failed to load plugin: %v", err)
}

// The plugin automatically registers its extension with the extension manager
```

### Plugin Structure

A plugin must export an `Extension` variable:

```go
// plugin.go
package main

import "github.com/krizos/php-go/pkg/goext"
import "github.com/krizos/php-go/pkg/types"

// Extension is the exported symbol that the plugin system looks for
var Extension = &goext.Extension{
    Name:        "my_plugin",
    Version:     "1.0.0",
    Description: "My custom plugin",
    Functions:   make(map[string]goext.PHPFunction),
}

func init() {
    Extension.RegisterFunction("plugin_hello", func(args []types.Value) (types.Value, error) {
        return types.NewString("Hello from plugin!"), nil
    })
}
```

Build the plugin:

```bash
go build -buildmode=plugin -o my_plugin.so plugin.go
```

### Plugin Management

```go
// List loaded plugins
plugins := pluginMgr.List()
for _, info := range plugins {
    fmt.Printf("Plugin: %s (loaded: %v)\n", info.Path, info.Loaded)
}

// Get plugin info
info, err := pluginMgr.GetInfo("/path/to/plugin.so")
if err == nil {
    fmt.Printf("Extension: %s\n", info.ExtensionName)
}

// Check if plugin is loaded
if pluginMgr.IsLoaded("/path/to/plugin.so") {
    fmt.Println("Plugin is loaded")
}
```

**Note**: Hot reloading is not supported due to Go plugin system limitations.

## Runtime Configuration

### Global Variables and Superglobals

```go
import "github.com/krizos/php-go/pkg/runtime"

// Set global variable
runtime.SetGlobal("myvar", types.NewString("value"))

// Get global variable
value := runtime.GetGlobal("myvar")

// Initialize superglobals ($_GET, $_POST, $_SERVER, etc.)
runtime.InitializeSuperglobals()

// Set $_SERVER values
runtime.SetServerVar("HTTP_HOST", "example.com")
runtime.SetServerVar("REQUEST_METHOD", "GET")

// Set $_GET values
runtime.SetGetVar("page", "1")
runtime.SetGetVar("limit", "10")

// Set $_POST values
runtime.SetPostVar("username", "john")
runtime.SetPostVar("password", "secret")
```

### Constants

```go
// Define constant
runtime.DefineConstant("MY_CONSTANT", types.NewString("value"))

// Check if constant is defined
if runtime.IsConstantDefined("MY_CONSTANT") {
    value := runtime.GetConstant("MY_CONSTANT")
}
```

### Error Handling

```go
// Set error handler
runtime.SetErrorHandler(func(errno int, errstr string, errfile string, errline int) {
    log.Printf("PHP Error [%d]: %s in %s:%d", errno, errstr, errfile, errline)
})

// Trigger error
runtime.TriggerError("Something went wrong", runtime.E_USER_WARNING)
```

### Output Buffering

```go
// Start output buffering
runtime.ObStart()

// Get current buffer contents
output := runtime.ObGetContents()

// Clean buffer (delete contents)
runtime.ObClean()

// End buffering and get contents
output = runtime.ObGetClean()

// End buffering and flush contents
runtime.ObEndFlush()
```

## Environment Variables

PHP-Go respects several environment variables for configuration:

### Logging

```bash
# Set log level
export PHPGO_LOG_LEVEL=debug

# Set log format (text or json)
export PHPGO_LOG_FORMAT=json

# Set log output file
export PHPGO_LOG_FILE=/var/log/php-go.log
```

### Resource Limits

```bash
# Memory limit in MB
export PHPGO_MEMORY_LIMIT=128

# Execution time limit in seconds
export PHPGO_EXECUTION_TIME_LIMIT=30

# Maximum recursion depth
export PHPGO_MAX_RECURSION_DEPTH=1000

# Maximum output buffer size in MB
export PHPGO_MAX_OUTPUT_SIZE=10
```

### Extensions

```bash
# Extension directory
export PHPGO_EXTENSION_DIR=/usr/lib/php-go/extensions

# Load extensions
export PHPGO_EXTENSIONS=json,hash,mysqli

# Plugin directory
export PHPGO_PLUGIN_DIR=/usr/lib/php-go/plugins
```

### Development

```bash
# Enable debug mode
export PHPGO_DEBUG=1

# Enable verbose output
export PHPGO_VERBOSE=1

# Enable profiling
export PHPGO_PROFILE=1
export PHPGO_PROFILE_OUTPUT=/tmp/php-go.prof
```

## Best Practices

### Production Configuration

```go
// Set conservative resource limits
runtime.SetMemoryLimit(256 * 1024 * 1024)      // 256MB
runtime.SetExecutionTimeLimit(30)               // 30 seconds
runtime.SetMaxRecursionDepth(500)               // 500 levels
runtime.SetMaxOutputSize(10 * 1024 * 1024)     // 10MB

// Use JSON logging for structured logs
logger := runtime.GetGlobalLogger()
logger.SetLevel(runtime.LogLevelInfo)
logger.SetFormatter(runtime.NewJSONFormatter())

// Register health checks
runtime.RegisterHealthCheck("database", runtime.HealthCheckTypeReadiness,
    checkDatabaseConnection, 5*time.Second)
runtime.RegisterHealthCheck("redis", runtime.HealthCheckTypeReadiness,
    checkRedisConnection, 3*time.Second)

// Enable metrics collection
collector := runtime.GetGlobalMetricsCollector()
// Metrics are automatically collected

// Start resource tracking
runtime.StartResourceTracking()
```

### Development Configuration

```go
// Generous resource limits for development
runtime.SetMemoryLimit(0)                      // Unlimited
runtime.SetExecutionTimeLimit(0)               // Unlimited
runtime.SetMaxRecursionDepth(10000)            // High limit

// Use text logging with colors
logger := runtime.GetGlobalLogger()
logger.SetLevel(runtime.LogLevelDebug)
formatter := runtime.NewTextFormatter()
formatter.ColorEnabled = true
formatter.IncludeLocation = true
logger.SetFormatter(formatter)
logger.SetOutput(os.Stdout)

// Enable all debugging
logger.Debug("Development mode enabled")
```

### Monitoring and Observability

```go
// Set up periodic health checks
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        ctx := context.Background()
        report := runtime.GetHealthReport(ctx)

        if report.Status != runtime.HealthStatusHealthy {
            logger.Warn("Health check degraded or unhealthy")
            for name, result := range report.Checks {
                if result.Status != runtime.HealthStatusHealthy {
                    logger.WithFields(map[string]interface{}{
                        "check": name,
                        "status": result.Status,
                    }).Warn(result.Message)
                }
            }
        }
    }
}()

// Collect and report metrics periodically
go func() {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        stats := runtime.GetResourceStats()
        collector := runtime.GetGlobalMetricsCollector()

        // Record resource usage as metrics
        if memory, ok := stats["memory"].(map[string]interface{}); ok {
            if current, ok := memory["current"].(uint64); ok {
                collector.Gauge("memory_usage_bytes", float64(current), nil)
            }
        }

        // Export metrics to monitoring system
        metrics := collector.GetAll()
        // Send to Prometheus, StatsD, etc.
    }
}()
```

### Error Handling

```go
// Set up error recovery
runtime.RegisterErrorHandler(func(err error) {
    logger.Error("Runtime error: " + err.Error())

    // Send to error tracking service (Sentry, etc.)
    // sentryClient.CaptureException(err)
})

// Set up panic recovery
defer func() {
    if r := recover(); r != nil {
        logger.Fatal(fmt.Sprintf("Panic recovered: %v", r))
        // Graceful shutdown
        runtime.StopResourceTracking()
    }
}()
```

### Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
)

// Set up signal handling for graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    logger.Info("Shutdown signal received")

    // Stop resource tracking
    runtime.StopResourceTracking()

    // Flush logs
    logger.Flush()

    // Close metrics collector
    collector.Close()

    // Unload extensions
    mgr := goext.GetExtensionManager()
    for _, name := range mgr.List() {
        mgr.Unregister(name)
    }

    os.Exit(0)
}()
```

## Configuration File Example

While PHP-Go doesn't currently support configuration files directly, you can implement your own:

```go
import (
    "encoding/json"
    "os"
)

type Config struct {
    ResourceLimits struct {
        MemoryMB        int `json:"memory_mb"`
        ExecutionTimeSec int `json:"execution_time_sec"`
        RecursionDepth   int `json:"recursion_depth"`
        OutputSizeMB     int `json:"output_size_mb"`
    } `json:"resource_limits"`

    Logging struct {
        Level  string `json:"level"`
        Format string `json:"format"`
        Output string `json:"output"`
    } `json:"logging"`

    Extensions struct {
        Directory string   `json:"directory"`
        Load      []string `json:"load"`
    } `json:"extensions"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func ApplyConfig(config *Config) {
    // Apply resource limits
    runtime.SetMemoryLimit(int64(config.ResourceLimits.MemoryMB * 1024 * 1024))
    runtime.SetExecutionTimeLimit(int64(config.ResourceLimits.ExecutionTimeSec))
    runtime.SetMaxRecursionDepth(int64(config.ResourceLimits.RecursionDepth))
    runtime.SetMaxOutputSize(int64(config.ResourceLimits.OutputSizeMB * 1024 * 1024))

    // Apply logging configuration
    logger := runtime.GetGlobalLogger()
    switch config.Logging.Level {
    case "debug":
        logger.SetLevel(runtime.LogLevelDebug)
    case "info":
        logger.SetLevel(runtime.LogLevelInfo)
    case "warn":
        logger.SetLevel(runtime.LogLevelWarn)
    case "error":
        logger.SetLevel(runtime.LogLevelError)
    }

    // Load extensions
    mgr := goext.GetExtensionManager()
    for _, ext := range config.Extensions.Load {
        // Load extension from directory
        pluginPath := filepath.Join(config.Extensions.Directory, ext+".so")
        goext.GetPluginManager().Load(pluginPath)
    }
}
```

Example `config.json`:

```json
{
  "resource_limits": {
    "memory_mb": 256,
    "execution_time_sec": 30,
    "recursion_depth": 1000,
    "output_size_mb": 10
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "/var/log/php-go.log"
  },
  "extensions": {
    "directory": "/usr/lib/php-go/extensions",
    "load": ["json", "hash", "mysqli"]
  }
}
```

## Next Steps

- Read the [Extension Development Guide](extension-development.md) to learn how to create custom extensions
- See the [Performance Tuning Guide](performance-tuning.md) for optimization tips
- Check the [API Reference](api-reference.md) for detailed API documentation
- Review the [Migration Guide](migration-guide.md) for migrating from standard PHP

For more information, visit the [PHP-Go documentation](../README.md).
