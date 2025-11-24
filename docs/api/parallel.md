# Parallel API Reference

Package: `github.com/krizos/php-go/pkg/parallel`

## Overview

The parallel package will implement automatic parallelization for PHP code, allowing safe concurrent execution of PHP scripts. This package is planned for Phase 7 and is currently a placeholder.

**Status**: Not yet implemented (Phase 7 - Planned)

**Implementation Timeline**: Phase 7 of the project roadmap

**Goal**: Automatic parallelization of PHP code with safety analysis and worker pools

## Planned Features

### 1. Safety Analysis
Analyze PHP code to determine if it's safe to parallelize:
- Detect shared state access
- Identify data dependencies
- Check for side effects
- Analyze I/O operations
- Detect potential race conditions

### 2. Parallel Execution Contexts
```go
type ParallelContext struct {
    // Worker pool for executing tasks
    // Copy-on-write data structures
    // Isolated state per goroutine
}
```

### 3. Worker Pools
```go
type WorkerPool struct {
    // Configurable number of workers
    // Task queue
    // Result collection
}
```

### 4. Copy-on-Write (COW)
Implement COW semantics for arrays and strings to enable safe sharing:
```go
type CopyOnWrite struct {
    // Reference counting
    // Lazy copying
    // Thread-safe operations
}
```

### 5. Parallel Array Operations
```go
// Parallel map
func ParallelMap(arr *types.Array, fn Callable) *types.Array

// Parallel filter
func ParallelFilter(arr *types.Array, fn Callable) *types.Array

// Parallel reduce
func ParallelReduce(arr *types.Array, fn Callable, initial *types.Value) *types.Value
```

### 6. Synchronization Primitives
```go
type Mutex struct { /* ... */ }
type RWMutex struct { /* ... */ }
type Channel struct { /* ... */ }
type WaitGroup struct { /* ... */ }
```

### 7. Concurrent Data Structures
```go
type ConcurrentArray struct {
    // Thread-safe array operations
}

type ConcurrentMap struct {
    // Thread-safe map operations
}
```

## Planned API

### Parallel Execution

```go
// Execute function in parallel with given inputs
func Parallel(fn Callable, inputs []*types.Value, options *ParallelOptions) ([]*types.Value, error)

// Parallel options
type ParallelOptions struct {
    MaxWorkers int
    Timeout    time.Duration
    Strategy   ParallelStrategy
}

type ParallelStrategy int

const (
    StrategyAuto ParallelStrategy = iota
    StrategyDataParallel
    StrategyTaskParallel
)
```

### Safety Analysis

```go
// Analyze if code is safe to parallelize
func AnalyzeSafety(code *ast.Program) (*SafetyReport, error)

type SafetyReport struct {
    IsSafe        bool
    Warnings      []string
    SharedState   []string
    Dependencies  []Dependency
    Recommendations []string
}

type Dependency struct {
    Type  DependencyType
    From  ast.Node
    To    ast.Node
}
```

### Worker Pool

```go
// Create worker pool
func NewWorkerPool(size int) *WorkerPool

// Submit task to pool
func (wp *WorkerPool) Submit(task Task) Future

// Task interface
type Task interface {
    Execute() (*types.Value, error)
}

// Future for async results
type Future interface {
    Get() (*types.Value, error)
    IsDone() bool
    Cancel() error
}
```

### Parallel Collections

```go
// Parallel array
type ParallelArray struct {
    data *types.Array
    lock sync.RWMutex
}

func NewParallelArray() *ParallelArray
func (pa *ParallelArray) Get(key *types.Value) (*types.Value, error)
func (pa *ParallelArray) Set(key, value *types.Value) error
func (pa *ParallelArray) Map(fn Callable) (*ParallelArray, error)
func (pa *ParallelArray) Filter(fn Callable) (*ParallelArray, error)
func (pa *ParallelArray) Reduce(fn Callable, initial *types.Value) (*types.Value, error)
```

### Copy-on-Write

```go
// COW array
type COWArray struct {
    // Shared read access
    // Copy on first write
    // Reference counting
}

func NewCOWArray(arr *types.Array) *COWArray
func (cw *COWArray) Read() *types.Array
func (cw *COWArray) Write() *types.Array
func (cw *COWArray) IsShared() bool
```

### Synchronization

```go
// Mutex for PHP
type PHPMutex struct {
    mu sync.Mutex
}

func NewMutex() *PHPMutex
func (m *PHPMutex) Lock()
func (m *PHPMutex) Unlock()
func (m *PHPMutex) TryLock() bool

// Channels for PHP
type PHPChannel struct {
    ch chan *types.Value
}

func NewChannel(capacity int) *PHPChannel
func (c *PHPChannel) Send(value *types.Value) error
func (c *PHPChannel) Receive() (*types.Value, error)
func (c *PHPChannel) Close()
```

## Usage Example (Future)

Once implemented, parallel execution will work like this:

```go
package main

import (
    "github.com/krizos/php-go/pkg/parallel"
    "github.com/krizos/php-go/pkg/types"
)

func main() {
    // Create array of values to process
    arr := types.NewEmptyArray()
    for i := 0; i < 1000; i++ {
        arr.Append(types.NewInt(int64(i)))
    }

    // Process in parallel
    results, err := parallel.ParallelMap(arr, func(val *types.Value) *types.Value {
        // This closure will run in parallel across workers
        n := val.ToInt()
        return types.NewInt(n * n)
    })

    if err != nil {
        panic(err)
    }

    // results now contains squares of all numbers
}
```

## PHP Integration (Future)

Parallel execution in PHP code:

```php
<?php
// Parallel array mapping
$numbers = range(1, 1000);
$squares = parallel_map($numbers, function($n) {
    return $n * $n;
});

// Parallel foreach
parallel_foreach($numbers as $n) {
    // Each iteration runs in parallel
    process($n);
}

// Manual worker pool
$pool = new WorkerPool(8);
foreach ($tasks as $task) {
    $future = $pool->submit($task);
    $futures[] = $future;
}

// Collect results
$results = array_map(fn($f) => $f->get(), $futures);
```

## Implementation Plan (Phase 7)

### Task 7.1: Safety Analyzer (40 hours)
Implement static analysis to detect parallelization safety.

### Task 7.2: Worker Pools (30 hours)
Implement worker pool with task scheduling.

### Task 7.3: Copy-on-Write (35 hours)
Implement COW for arrays and strings.

### Task 7.4: Parallel Array Operations (30 hours)
Implement parallel map, filter, reduce.

### Task 7.5: Synchronization Primitives (25 hours)
Implement mutexes, channels, wait groups.

### Task 7.6: Integration with VM (40 hours)
Integrate parallel execution with VM.

### Task 7.7: Testing and Benchmarking (50 hours)
Comprehensive testing and performance optimization.

**Total Estimated**: 250 hours for Phase 7

## Design Considerations

1. **Safety First**: Only parallelize when safe
2. **Backward Compatibility**: Single-threaded PHP code works unchanged
3. **Performance**: Overhead should be minimal for serial code
4. **Debugging**: Clear error messages for race conditions
5. **Resource Management**: Automatic cleanup of workers and locks

## Performance Goals

- **Speedup**: 2-8x for parallelizable workloads on multi-core CPUs
- **Overhead**: <5% for serial execution
- **Scalability**: Linear scaling up to 8-16 cores

## Challenges

1. **PHP's Mutable State**: PHP code often has shared mutable state
2. **Global Variables**: Need isolation or synchronization
3. **I/O Operations**: File handles, database connections need special handling
4. **Error Handling**: Exceptions in parallel code are complex
5. **Debugging**: Race conditions are hard to debug

## Related Work

- PHP's pthreads extension (deprecated in PHP 8)
- PHP Parallel extension (pecl/parallel)
- ReactPHP for async I/O
- Swoole for concurrent PHP

## References

- Phase 7 Documentation: `/Users/krizos/code/mnohosten/php-go/docs/phases/07-parallelization/`
- Go concurrency patterns
- Copy-on-write semantics
- Work stealing algorithms

## Safety Analysis Rules

### Safe to Parallelize
- Pure functions (no side effects)
- Read-only data access
- Local variables only
- Immutable data structures

### Unsafe to Parallelize
- Shared mutable state
- Global variable writes
- File I/O without locking
- Database operations
- Resource creation/destruction

## Future Enhancements (Beyond Phase 7)

- GPU acceleration for numeric operations
- Distributed execution across machines
- Automatic parallelization hints from profiler
- Parallel debugging tools
- Integration with async/await (Phase 9)
