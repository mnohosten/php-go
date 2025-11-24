#!/bin/bash
# Quick Memory Profile - Generates a sample profile for the most critical benchmark
# Use this for quick iteration during optimization work

set -e

PROFILE_DIR="profiles"
mkdir -p "$PROFILE_DIR"

echo "Running quick memory profile on SimpleLoop benchmark..."
echo ""

# Profile the most critical benchmark (SimpleLoop)
go test -run=^$ -bench=^BenchmarkSimpleLoop$ \
    -benchtime=1s \
    -benchmem \
    -memprofile="${PROFILE_DIR}/quick_mem.prof" \
    -cpuprofile="${PROFILE_DIR}/quick_cpu.prof" \
    -memprofilerate=1 \
    ./benchmarks

echo ""
echo "Profile saved to: ${PROFILE_DIR}/quick_mem.prof"
echo ""
echo "View with:"
echo "  go tool pprof -http=:8080 ${PROFILE_DIR}/quick_mem.prof"
echo ""
echo "Or generate a quick text report:"
echo "  go tool pprof -text -alloc_space ${PROFILE_DIR}/quick_mem.prof | head -30"
echo ""
