# PHP-Go Hot Path Analysis

**Date**: November 24, 2025
**Platform**: Apple M4 Max (darwin/arm64, 14 cores)
**Analysis Phase**: Phase 10.7 - Optimization Pass
**Profiling Method**: CPU profiling with pprof

## Executive Summary

This document identifies the hottest execution paths (CPU-intensive code paths) in PHP-Go based on CPU profiling across all core components. This complements the existing bottleneck analysis (focused on allocations) with runtime execution data.

### Key Findings

1. **VM Execution Loop** - 9.03% of total CPU time (primary hot path)
2. **Value Creation (NewInt)** - 7.74% of CPU time + high allocation overhead
3. **Lexer NextToken** - 83.19% of lexer CPU time (expected, core loop)
4. **Parser Registration Overhead** - 12.44% of parser time (initialization cost)

---

## 1. VM Hot Paths (SimpleLoop Benchmark)

### CPU Profile Summary
**Total Samples**: 1.55s
**Primary Workload**: 10,000 loop iterations with integer arithmetic

### Top VM Functions by CPU Time

| Function | Flat% | Cum% | Role | Priority |
|----------|-------|------|------|----------|
| `VM.run()` | 0.65% | 9.03% | Main execution loop | CRITICAL |
| `VM.dispatch()` | 0% | 8.39% | Opcode dispatcher | CRITICAL |
| `types.NewInt()` | 1.29% | 7.74% | Integer value creation | CRITICAL |
| `VM.getOperandValue()` | 0.65% | 6.45% | Operand extraction | HIGH |
| `VM.GetConstant()` | 0% | 5.81% | Constant pool access | HIGH |
| `VM.opQMAssign()` | 0% | 5.16% | Assignment handler | HIGH |
| `VM.opAdd()` | 0% | 1.94% | Addition handler | MEDIUM |
| `VM.opFetch()` | 0% | 0.65% | Variable fetch | LOW |
| `VM.opJmp()` | 0% | 0.65% | Jump instruction | LOW |

### Analysis

#### 1.1 VM Execution Loop - CRITICAL HOT PATH
- **Location**: `pkg/vm/vm.go` - `VM.run()` method
- **CPU Time**: 9.03% cumulative
- **Issue**: Central execution loop with high call frequency
- **Root Cause**:
  - Opcode dispatch overhead (switch statement)
  - Stack manipulation on every instruction
  - No instruction caching or JIT
- **Impact**: Every instruction passes through this path ~10K+ times per benchmark
- **Optimization Opportunities**:
  1. **Direct threading** - Replace switch with computed goto (requires assembly or unsafe)
  2. **Instruction fusion** - Combine common instruction sequences
  3. **Stack caching** - Keep top stack values in registers/variables
  4. **Predicted dispatch** - Profile-guided opcode ordering in switch

#### 1.2 Value Creation Overhead - CRITICAL HOT PATH
- **Location**: `pkg/types/value.go` - `NewInt()` function
- **CPU Time**: 7.74% cumulative, 1.29% flat
- **Allocations**: 88% of total allocations (from memory profiling)
- **Issue**: Called ~20K times per 10K loop iterations (2x per iteration)
- **Root Cause**:
  - No value pooling
  - No integer caching (common values like 0, 1, 2)
  - Every arithmetic operation creates new Value
- **Impact**: Combined CPU + allocation bottleneck
- **Optimization Opportunities**:
  1. **sync.Pool for Values** - Reuse Value structs (CRITICAL - 80-90% allocation reduction)
  2. **Integer cache** - Pre-allocate common integers -128 to 1024 (HIGH - 60-70% NewInt reduction)
  3. **In-place operations** - Reuse result slot when safe (MEDIUM)
  4. **Escape analysis** - Stack allocate temporary values (RESEARCH)

#### 1.3 Operand Extraction - HIGH PRIORITY
- **Location**: `pkg/vm/vm.go` - `VM.getOperandValue()`
- **CPU Time**: 6.45% cumulative
- **Issue**: Called for every operand in every instruction
- **Root Cause**:
  - Multiple switch/if branches to decode operand type
  - Constant pool lookups
  - Stack indexing
- **Optimization Opportunities**:
  1. **Inline hot paths** - Ensure inlining of simple cases
  2. **Operand type prediction** - Most operands are constants or stack values
  3. **Fast path for common operands** - Optimize for CV (constant) and TMP (stack) types

#### 1.4 Constant Pool Access - HIGH PRIORITY
- **Location**: `pkg/vm/vm.go` - `VM.GetConstant()`
- **CPU Time**: 5.81% cumulative
- **Issue**: Array indexing with bounds checking
- **Root Cause**:
  - Go bounds checking on every access
  - No constant caching
  - Constants are copied, not referenced
- **Optimization Opportunities**:
  1. **Bounds check elimination** - Use unsafe or hints
  2. **Constant caching** - Cache frequently-used constants in frame
  3. **Constant inlining** - Embed small constants in instruction (like PHP opcache)

#### 1.5 Assignment Handler - HIGH PRIORITY
- **Location**: `pkg/vm/handlers_variables.go` - `VM.opQMAssign()`
- **CPU Time**: 5.16% cumulative
- **Issue**: Very common operation (every variable assignment)
- **Root Cause**:
  - Stack push operation
  - Value copying
  - Type checking
- **Optimization Opportunities**:
  1. **Eliminate redundant assignments** - Compiler optimization
  2. **Direct stack manipulation** - Avoid function call overhead
  3. **Copy-on-write** - Share values when safe

### VM Optimization Priority List

1. **CRITICAL** - Implement Value pooling (sync.Pool) → 80-90% allocation reduction, 30-40% CPU reduction
2. **CRITICAL** - Integer value cache (-128 to 1024) → 60-70% NewInt() elimination
3. **HIGH** - Optimize getOperandValue() → Inline and fast-path common cases
4. **HIGH** - Constant pool optimization → Bounds check elimination, caching
5. **MEDIUM** - Instruction fusion → Combine common sequences (FETCH + ADD, etc.)
6. **RESEARCH** - Direct threading → Replace switch with computed goto (requires unsafe/asm)

---

## 2. Lexer Hot Paths

### CPU Profile Summary
**Benchmark**: SimpleTokens
**Total Samples**: 1.19s
**Throughput**: 3.67M ops/sec

### Top Lexer Functions by CPU Time

| Function | Flat% | Cum% | Role | Priority |
|----------|-------|------|------|----------|
| `Lexer.NextToken()` | 29.41% | 83.19% | Main tokenization loop | CRITICAL |
| `Lexer.readChar()` | 18.49% | 18.49% | Character reading | CRITICAL |
| `Lexer.scanVariable()` | 5.88% | 10.08% | Variable scanning | MEDIUM |
| `Lexer.skipWhitespace()` | 5.04% | 10.08% | Whitespace skipping | MEDIUM |
| `Lexer.currentPosition()` | 4.20% | 4.20% | Position tracking | LOW |
| `Lexer.makeToken()` | 2.52% | 2.52% | Token creation | LOW |
| `Lexer.scanIdentifier()` | 1.68% | 5.04% | Identifier scanning | LOW |

### Analysis

#### 2.1 NextToken Loop - CRITICAL HOT PATH
- **Location**: `pkg/lexer/lexer.go` - `Lexer.NextToken()`
- **CPU Time**: 83.19% cumulative, 29.41% flat
- **Issue**: Central tokenization loop, called for every token
- **Root Cause**:
  - Large switch statement on current character
  - Multiple function calls (readChar, peekChar, etc.)
  - Position tracking overhead
- **Performance**: Already quite good (3.67M ops/sec, 313ns/op)
- **Optimization Opportunities**:
  1. **Character class table** - Replace switch with lookup table for character types
  2. **Inline hot paths** - Inline readChar, peekChar, makeToken (already done)
  3. **Reduce branching** - Profile-guided switch ordering
  4. **Token pooling** - Reuse Token structs (MEDIUM priority)

#### 2.2 Character Reading - CRITICAL HOT PATH
- **Location**: `pkg/lexer/lexer.go` - `Lexer.readChar()`
- **CPU Time**: 18.49% flat (already inlined)
- **Issue**: Called 1-10 times per token
- **Root Cause**:
  - String indexing with bounds checking
  - Position increment
  - Column tracking
- **Optimization Opportunities**:
  1. **Byte slice instead of string** - Avoid UTF-8 overhead for ASCII
  2. **Batch character reading** - Read multiple chars when safe
  3. **Simplified position tracking** - Only track when needed

#### 2.3 Variable Scanning - MEDIUM PRIORITY
- **Location**: `pkg/lexer/lexer.go` - `Lexer.scanVariable()`
- **CPU Time**: 10.08% cumulative, 5.88% flat
- **Issue**: Variables are extremely common in PHP
- **Root Cause**:
  - Character-by-character scanning
  - String building
- **Optimization Opportunities**:
  1. **Fast path for simple vars** - `$a`, `$i`, etc. without loop
  2. **Buffer pooling** - Reuse string builders

### Lexer Optimization Priority List

1. **MEDIUM** - Token pooling → 30-40% allocation reduction
2. **MEDIUM** - Character class table → Reduce switch overhead
3. **LOW** - Byte slice input → Faster character access for ASCII
4. **LOW** - Variable scanning fast path → Optimize common case

**Note**: Lexer is already highly optimized (3.67M ops/sec). Further optimization has diminishing returns.

---

## 3. Parser Hot Paths

### CPU Profile Summary
**Benchmark**: SimpleExpression
**Total Samples**: 2.25s
**Throughput**: 307K ops/sec

### Top Parser Functions by CPU Time

| Function | Flat% | Cum% | Role | Priority |
|----------|-------|------|------|----------|
| `Parser.registerExpressionParsers()` | 0% | 9.78% | Parser initialization | HIGH |
| `Parser.New()` | 0% | 10.22% | Parser creation | HIGH |
| `Parser.ParseProgram()` | 0% | 2.22% | Top-level parsing | LOW |
| `Parser.parseExpression()` | 0% | 1.33% | Expression parsing | LOW |
| `Parser.parseExpressionStatement()` | 0% | 1.78% | Statement parsing | LOW |

### Analysis

#### 3.1 Parser Initialization Overhead - HIGH PRIORITY
- **Location**: `pkg/parser/parser.go` - `registerExpressionParsers()`
- **CPU Time**: 9.78% cumulative (INITIALIZATION OVERHEAD)
- **Issue**: Called on every `parser.New()` invocation
- **Root Cause**:
  - Registering 50+ expression parser functions in maps
  - Map allocations for prefix/infix parsers
  - Function value creation
- **Impact**: Benchmark includes parser creation time, inflating these numbers
- **Optimization Opportunities**:
  1. **Global parser registration** - Register parsers once at init() time
  2. **Parser pooling** - Reuse Parser instances (sync.Pool)
  3. **Static dispatch table** - Use arrays instead of maps for common tokens

#### 3.2 Actual Parsing Performance
- **Combined parsing time**: ~1.78% cumulative for actual parsing work
- **Performance**: 307K ops/sec at 3.8μs/op (reasonable for full AST construction)
- **Analysis**: Parser initialization dominates benchmark; actual parsing is efficient

### Parser Optimization Priority List

1. **HIGH** - Parser pooling (sync.Pool) → Eliminate initialization overhead
2. **MEDIUM** - Global expression parser registration → Reduce New() cost
3. **MEDIUM** - AST node pooling → Reduce allocation overhead (from bottleneck analysis)
4. **LOW** - Static dispatch tables → Replace maps with arrays for common tokens

---

## 4. Compiler Hot Paths

### CPU Profile Summary
**Benchmark**: SimpleExpression
**Total Samples**: 2.30s
**Throughput**: 1.34M ops/sec

### Top Compiler Functions by CPU Time

| Function | Flat% | Cum% | Role | Priority |
|----------|-------|------|------|----------|
| `Compiler.EmitWithLine()` | 0.87% | 1.74% | Instruction emission | MEDIUM |
| `Compiler.Compile()` | 0.43% | 5.22% | Main compilation | LOW |
| `Compiler.DefineVariable()` | 0% | 1.30% | Symbol table ops | LOW |
| `Compiler.AddConstant()` | 0% | 0.87% | Constant addition | LOW |

### Analysis

#### 4.1 Instruction Emission - MEDIUM PRIORITY
- **Location**: `pkg/compiler/compiler.go` - `EmitWithLine()`
- **CPU Time**: 1.74% cumulative, 0.87% flat
- **Issue**: Called for every emitted instruction
- **Root Cause**:
  - Instruction struct creation
  - Slice append (dynamic growth)
  - Line number tracking
- **Optimization Opportunities**:
  1. **Pre-size instruction buffer** - Estimate size from AST node count
  2. **Batch emission** - Emit multiple instructions at once for common patterns
  3. **Simplified line tracking** - Only track when debugging enabled

#### 4.2 Overall Compiler Performance
- **Performance**: 1.34M ops/sec at 891ns/op (EXCELLENT)
- **CPU Time**: Very low (5.22% total for main compilation)
- **Analysis**: Compiler is highly efficient; optimization has low priority

### Compiler Optimization Priority List

1. **MEDIUM** - Pre-size instruction buffers → Reduce allocation overhead
2. **LOW** - Batch instruction emission → Reduce function call overhead
3. **LOW** - Constant table optimization → String interning

**Note**: Compiler is not a bottleneck. Focus on VM and runtime optimizations first.

---

## 5. Overall Hot Path Summary

### Critical Hot Paths (Ordered by Impact)

1. **VM Execution Loop** (`VM.run()`, `VM.dispatch()`)
   - **Impact**: 9.03% CPU time, core of interpreter
   - **Fix**: Value pooling, instruction fusion, direct threading
   - **Effort**: 8-12 hours (value pooling), 20-40 hours (direct threading research)

2. **Value Creation** (`types.NewInt()`)
   - **Impact**: 7.74% CPU + 88% of allocations
   - **Fix**: sync.Pool, integer cache, in-place operations
   - **Effort**: 6-8 hours

3. **Operand Handling** (`VM.getOperandValue()`, `VM.GetConstant()`)
   - **Impact**: 12.26% combined CPU time
   - **Fix**: Inlining, fast paths, bounds check elimination
   - **Effort**: 4-6 hours

4. **Lexer Tokenization** (`Lexer.NextToken()`, `Lexer.readChar()`)
   - **Impact**: 83.19% of lexer time (but lexer is fast)
   - **Fix**: Token pooling, character class table
   - **Effort**: 4-6 hours

5. **Parser Initialization** (`Parser.New()`, `registerExpressionParsers()`)
   - **Impact**: 10.22% of parser time (initialization overhead)
   - **Fix**: Parser pooling, global registration
   - **Effort**: 3-4 hours

### Optimization Roadmap (Aligned with Bottleneck Analysis)

#### Phase 1: VM Critical Path (12-16 hours) - HIGHEST ROI
1. Implement Value pooling (sync.Pool) - 6-8h → **80-90% allocation reduction, 30-40% CPU reduction**
2. Implement integer cache (-128 to 1024) - 2-3h → **60-70% NewInt() reduction**
3. Optimize getOperandValue() (inlining, fast paths) - 2-3h → **20-30% operand overhead reduction**
4. Constant pool optimization (bounds checks, caching) - 2-3h → **15-20% constant access speedup**

**Expected Result**: **5-10x faster VM execution**, ~200K allocs → <20K allocs for SimpleLoop

#### Phase 2: Parser/Lexer Pooling (6-8 hours) - MEDIUM ROI
1. Implement Parser pooling (sync.Pool) - 2-3h → **Eliminate 10% initialization overhead**
2. Implement Token pooling (sync.Pool) - 2-3h → **30-40% lexer allocation reduction**
3. AST node pooling (sync.Pool) - 2-3h → **40-50% parser allocation reduction**

**Expected Result**: **2-3x faster parsing**, better startup time

#### Phase 3: Instruction-Level Optimization (12-20 hours) - ADVANCED
1. Instruction fusion - 6-8h → **Combine common opcode sequences**
2. Stack caching - 4-6h → **Keep top values in local variables**
3. Direct threading research - 20-40h → **Replace switch with computed goto** (RESEARCH)

**Expected Result**: **2-5x additional VM speedup** (highly dependent on implementation)

---

## 6. Measurement Plan

### Before/After Profiling

For each optimization phase:

```bash
# Baseline profile
go test -run=^$ -bench=BenchmarkSimpleLoop -benchtime=3s \
  -cpuprofile=before_cpu.prof -memprofile=before_mem.prof \
  -benchmem ./benchmarks/ > before.txt

# After optimization
go test -run=^$ -bench=BenchmarkSimpleLoop -benchtime=3s \
  -cpuprofile=after_cpu.prof -memprofile=after_mem.prof \
  -benchmem ./benchmarks/ > after.txt

# Compare
benchstat before.txt after.txt

# Profile comparison
go tool pprof -base before_cpu.prof after_cpu.prof
go tool pprof -base before_mem.prof after_mem.prof
```

### Success Metrics

**Phase 1 Success Criteria**:
- SimpleLoop CPU time: <1.0ms (from 4.1ms) → **4x improvement**
- NewInt CPU%: <2% (from 7.74%) → **3.8x reduction**
- Total allocations: <20K (from 199K) → **10x reduction**
- No test regressions

**Phase 2 Success Criteria**:
- Parser creation overhead: <5% (from 10.22%) → **2x reduction**
- Lexer allocations: <100 (from ~150) → **30% reduction**
- Parser allocations: <80 (from 137) → **40% reduction**

**Phase 3 Success Criteria**:
- VM dispatch overhead: <5% (from 8.39%) → **40% reduction**
- Overall CPU time: <500μs for SimpleLoop → **8x improvement from baseline**

---

## 7. Profiling Tools Reference

### CPU Profiling

```bash
# Profile specific component
go test -run=^$ -bench=BenchmarkName -cpuprofile=cpu.prof ./pkg/component/

# View top functions
go tool pprof -text -cum cpu.prof | head -50

# Interactive view
go tool pprof -http=:8080 cpu.prof

# Compare profiles
go tool pprof -base=before.prof after.prof
```

### Memory Profiling

```bash
# Profile allocations
go test -run=^$ -bench=BenchmarkName -memprofile=mem.prof -benchmem ./pkg/component/

# View allocation hot spots
go tool pprof -text -alloc_space mem.prof | head -50
go tool pprof -text -alloc_objects mem.prof | head -50

# Interactive view
go tool pprof -http=:8080 mem.prof
```

### Combined Analysis

```bash
# Run full profiling suite
./benchmarks/profile_memory.sh

# Analyze all profiles
./benchmarks/analyze_memory.sh

# Quick profile for iteration
./benchmarks/quick_profile.sh
```

---

## 8. Next Actions

### Immediate (This Week)
1. **Start Value pooling implementation** - `pkg/types/value.go`
   - Add sync.Pool for Value structs
   - Implement Get/Put methods
   - Update VM handlers to use pool
   - Measure impact with benchmarks

### Short Term (Next 2 Weeks)
1. **Integer cache implementation** - `pkg/types/value.go`
   - Pre-allocate integers -128 to 1024
   - Update NewInt() to return cached values
   - Measure NewInt() CPU reduction

2. **Operand access optimization** - `pkg/vm/vm.go`
   - Add fast paths for CV and TMP operands
   - Inline critical paths
   - Bounds check elimination

### Medium Term (Next Month)
1. **Parser/Lexer pooling** - `pkg/parser/`, `pkg/lexer/`
2. **AST node pooling** - `pkg/ast/`
3. **Instruction fusion research** - `pkg/vm/`

---

## 9. Conclusion

CPU profiling confirms the findings from allocation profiling:

1. **VM is the critical hot path** - Value creation and operand handling dominate
2. **Pooling is the #1 optimization** - Eliminates both CPU and allocation overhead
3. **Lexer/Parser are efficient** - Not primary bottlenecks
4. **Compiler is excellent** - Very fast, low overhead

The hot path analysis provides specific function-level targets for optimization:
- `types.NewInt()` - 7.74% CPU, 88% allocations → **CRITICAL TARGET**
- `VM.getOperandValue()` - 6.45% CPU → **HIGH PRIORITY**
- `VM.dispatch()` - 8.39% CPU → **MEDIUM PRIORITY** (harder to optimize)

**Recommended first step**: Implement Value pooling for `NewInt()`, `NewBool()`, `NewFloat()` as it addresses both the top CPU hot path and the top allocation bottleneck.

**Expected overall improvement after all optimizations**: **5-10x faster execution** with **10-20x fewer allocations**.
