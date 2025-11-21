# Phase 7: Parallelization & Multi-threading

## Overview

Phase 7 implements automatic parallelization - the key differentiator of PHP-Go. This leverages Go's goroutines to automatically parallelize safe PHP operations while maintaining correctness.

## Goals

1. **Safety Analysis**: Determine which code can be parallelized safely
2. **Request-Level Parallelism**: Handle multiple PHP requests concurrently
3. **Automatic Parallelization**: Detect and parallelize safe array operations
4. **Shared-Nothing Architecture**: Isolate request contexts
5. **Concurrency Primitives**: New APIs for explicit parallelism

## Strategy

### Level 1: Request-Level Parallelism (Like PHP-FPM)
- Each HTTP request runs in a separate goroutine
- No shared state between requests
- Easy to implement, big performance win

### Level 2: Automatic Array Operation Parallelism
- Detect pure functions in array_map, array_filter, etc.
- Automatically parallelize when safe
- Transparent to PHP code

### Level 3: Explicit Parallelism (New APIs)
- New PHP functions for explicit parallel execution
- Channels for communication
- Futures/Promises

## Components

### 1. Safety Analyzer

**File**: `pkg/parallel/analyzer.go`

```go
type SafetyAnalyzer struct {
    ast *ast.File
}

func (sa *SafetyAnalyzer) CanParallelize(fn *Function) *SafetyReport
func (sa *SafetyAnalyzer) AnalyzeFunction(fn *Function) *SafetyReport

type SafetyReport struct {
    Safe            bool
    Reasons         []string
    GlobalReads     []string
    GlobalWrites    []string
    StaticAccess    bool
    FileIO          bool
    NetworkIO       bool
    DatabaseAccess  bool
}
```

**Safe Patterns**:
- Pure functions (no side effects)
- Read-only access to immutable data
- No global variable modifications
- No static variable access
- No file/network I/O

**Unsafe Patterns**:
- Global variable modifications
- Static variable access
- File writes
- Database writes
- Shared resource access

### 2. Parallel Executor

**File**: `pkg/parallel/executor.go`

```go
type ParallelExecutor struct {
    vm       *VM
    pool     *WorkerPool
    maxWorkers int
}

func (pe *ParallelExecutor) ParallelMap(
    arr *Array,
    fn *Function,
) *Array

func (pe *ParallelExecutor) ParallelFilter(
    arr *Array,
    fn *Function,
) *Array
```

### 3. Worker Pool

**File**: `pkg/parallel/pool.go`

```go
type WorkerPool struct {
    workers   chan *Worker
    maxSize   int
}

type Worker struct {
    vm    *VM
    ctx   *Context
}

func (wp *WorkerPool) Submit(task *Task) *Future
func (wp *WorkerPool) Wait()
```

### 4. Request Context Isolation

**File**: `pkg/parallel/context.go`

```go
type RequestContext struct {
    id        string
    vm        *VM          // Isolated VM
    globals   *GlobalScope // Request-specific
    output    *Buffer
    session   *Session
    startTime time.Time
}

func NewRequestContext() *RequestContext
func (rc *RequestContext) Isolate() // Copy-on-write setup
```

## Implementation Tasks

### Task 7.1: Safety Analyzer
**Effort**: 16 hours

- [ ] AST analysis for side effects
- [ ] Global variable tracking
- [ ] Static variable detection
- [ ] I/O operation detection
- [ ] Pure function detection
- [ ] Safety report generation

### Task 7.2: Worker Pool
**Effort**: 10 hours

- [ ] Worker pool implementation
- [ ] Worker lifecycle management
- [ ] Task queue
- [ ] Load balancing
- [ ] Graceful shutdown

### Task 7.3: Request-Level Parallelism
**Effort**: 12 hours

- [ ] Request context isolation
- [ ] Goroutine per request
- [ ] Context cleanup
- [ ] Error handling per request
- [ ] Request timeout handling

### Task 7.4: Automatic Array Parallelization
**Effort**: 16 hours

- [ ] Parallel array_map()
- [ ] Parallel array_filter()
- [ ] Parallel array_reduce()
- [ ] Parallel array_walk()
- [ ] Automatic threshold detection (when to parallelize)
- [ ] Result aggregation

### Task 7.5: Explicit Parallelism APIs
**Effort**: 14 hours

- [ ] go_routine($callable) - Run in goroutine
- [ ] go_wait($futures) - Wait for completion
- [ ] go_channel() - Create channel
- [ ] go_send($channel, $value)
- [ ] go_recv($channel)
- [ ] go_parallel($callables) - Run in parallel

**New PHP APIs for explicit control!**

### Task 7.6: Synchronization Primitives
**Effort**: 10 hours

- [ ] Mutex wrapper
- [ ] RWMutex wrapper
- [ ] WaitGroup wrapper
- [ ] Atomic operations
- [ ] Lock management

### Task 7.7: Copy-on-Write Optimization
**Effort**: 12 hours

- [ ] COW for arrays
- [ ] COW for strings
- [ ] COW for object properties
- [ ] Minimize copying overhead

### Task 7.8: Performance Monitoring
**Effort**: 8 hours

- [ ] Parallelization metrics
- [ ] Worker pool stats
- [ ] Contention detection
- [ ] Performance profiling

### Task 7.9: Testing
**Effort**: 16 hours

- [ ] Safety analyzer tests
- [ ] Concurrent request tests
- [ ] Race condition tests
- [ ] Performance benchmarks
- [ ] Stress tests

## Estimated Timeline

**Total Effort**: ~115 hours (6 weeks)

## Success Criteria

- [ ] Multiple requests execute concurrently
- [ ] Array operations auto-parallelize when safe
- [ ] No race conditions
- [ ] Performance scales with cores
- [ ] Explicit APIs functional
- [ ] Safety analysis accurate

## Performance Goals

- Linear scaling for independent requests
- 2-4x speedup for parallel array operations
- Zero overhead when parallelism not beneficial

## Progress Tracking

- [ ] Task 7.1: Safety Analyzer (0%)
- [ ] Task 7.2: Worker Pool (0%)
- [ ] Task 7.3: Request-Level Parallelism (0%)
- [ ] Task 7.4: Automatic Array Parallelization (0%)
- [ ] Task 7.5: Explicit Parallelism APIs (0%)
- [ ] Task 7.6: Synchronization Primitives (0%)
- [ ] Task 7.7: Copy-on-Write Optimization (0%)
- [ ] Task 7.8: Performance Monitoring (0%)
- [ ] Task 7.9: Testing (0%)

**Overall Phase 7 Progress**: 0%

---

**Status**: Not Started
**Last Updated**: 2025-11-21
