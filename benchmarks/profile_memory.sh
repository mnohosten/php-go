#!/bin/bash
# Memory Profiling Script for PHP-Go
# Generates comprehensive memory profiles using Go's pprof tool

set -e

PROFILE_DIR="profiles"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}PHP-Go Memory Profiling Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Create profiles directory
mkdir -p "$PROFILE_DIR"

# Function to run a benchmark with memory profiling
profile_benchmark() {
    local name=$1
    local package=$2
    local bench_name=$3
    local description=$4

    echo -e "${YELLOW}Profiling: ${description}${NC}"

    local mem_profile="${PROFILE_DIR}/${name}_mem_${TIMESTAMP}.prof"
    local cpu_profile="${PROFILE_DIR}/${name}_cpu_${TIMESTAMP}.prof"
    local alloc_profile="${PROFILE_DIR}/${name}_alloc_${TIMESTAMP}.prof"

    # Run benchmark with memory and CPU profiling
    go test -bench="^${bench_name}$" \
        -benchtime=3s \
        -benchmem \
        -memprofile="$mem_profile" \
        -cpuprofile="$cpu_profile" \
        -memprofilerate=1 \
        "$package" > "${PROFILE_DIR}/${name}_bench_${TIMESTAMP}.txt"

    # Generate allocation profile (tracks all allocations)
    go test -bench="^${bench_name}$" \
        -benchtime=3s \
        -benchmem \
        -memprofile="$alloc_profile" \
        -memprofilerate=1 \
        "$package" > /dev/null 2>&1

    echo -e "${GREEN}✓ Profiles saved:${NC}"
    echo -e "  Memory: $mem_profile"
    echo -e "  CPU: $cpu_profile"
    echo -e "  Alloc: $alloc_profile"
    echo ""
}

# Profile VM/Runtime (macro benchmarks)
echo -e "${BLUE}=== VM/Runtime Profiling ===${NC}"
profile_benchmark "vm_simple_loop" "./benchmarks" "BenchmarkSimpleLoop" "Simple Loop (10K iterations)"
profile_benchmark "vm_function_calls" "./benchmarks" "BenchmarkFunctionCalls" "Function Calls (5K calls)"
profile_benchmark "vm_recursion" "./benchmarks" "BenchmarkRecursion" "Recursion (Fibonacci 15)"
profile_benchmark "vm_string_concat" "./benchmarks" "BenchmarkStringConcatenation" "String Concatenation"

# Profile Parser
echo -e "${BLUE}=== Parser Profiling ===${NC}"
profile_benchmark "parser_simple_expr" "./pkg/parser" "BenchmarkSimpleExpression" "Parser - Simple Expression"
profile_benchmark "parser_complex_expr" "./pkg/parser" "BenchmarkComplexExpression" "Parser - Complex Expression"
profile_benchmark "parser_complex_class" "./pkg/parser" "BenchmarkComplexClass" "Parser - Complex Class"
profile_benchmark "parser_large_file" "./pkg/parser" "BenchmarkLargeFile" "Parser - Large File"

# Profile Compiler
echo -e "${BLUE}=== Compiler Profiling ===${NC}"
profile_benchmark "compiler_simple_expr" "./pkg/compiler" "BenchmarkSimpleExpression" "Compiler - Simple Expression"
profile_benchmark "compiler_complex_class" "./pkg/compiler" "BenchmarkComplexClass" "Compiler - Complex Class"
profile_benchmark "compiler_large_file" "./pkg/compiler" "BenchmarkLargeFile" "Compiler - Large File"

# Profile Lexer
echo -e "${BLUE}=== Lexer Profiling ===${NC}"
profile_benchmark "lexer_simple" "./pkg/lexer" "BenchmarkSimpleTokens" "Lexer - Simple Tokens"
profile_benchmark "lexer_complex" "./pkg/lexer" "BenchmarkComplexExpression" "Lexer - Complex Expression"
profile_benchmark "lexer_large_file" "./pkg/lexer" "BenchmarkLargeFile" "Lexer - Large File"

# Profile Type System
echo -e "${BLUE}=== Type System Profiling ===${NC}"
profile_benchmark "types_conversions" "./pkg/types" "BenchmarkIntToString" "Type System - Conversions"
profile_benchmark "types_juggling" "./pkg/types" "BenchmarkTypeJugglingIntString" "Type System - Type Juggling"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Memory Profiling Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Profile files saved in: ${PROFILE_DIR}/${NC}"
echo ""
echo -e "${YELLOW}View profiles with:${NC}"
echo -e "  ${GREEN}go tool pprof -http=:8080 ${PROFILE_DIR}/<profile>.prof${NC}"
echo ""
echo -e "${YELLOW}Generate reports:${NC}"
echo -e "  ${GREEN}go tool pprof -text ${PROFILE_DIR}/<profile>.prof${NC}"
echo -e "  ${GREEN}go tool pprof -top ${PROFILE_DIR}/<profile>.prof${NC}"
echo -e "  ${GREEN}go tool pprof -list=<function> ${PROFILE_DIR}/<profile>.prof${NC}"
echo ""
echo -e "${YELLOW}Compare profiles:${NC}"
echo -e "  ${GREEN}go tool pprof -base=before.prof after.prof${NC}"
echo ""
echo -e "${YELLOW}Generate flame graphs:${NC}"
echo -e "  ${GREEN}go tool pprof -svg ${PROFILE_DIR}/<profile>.prof > graph.svg${NC}"
echo ""
