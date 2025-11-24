#!/bin/bash
# Memory Profile Analysis Script
# Analyzes memory profiles and generates reports

set -e

PROFILE_DIR="profiles"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

if [ ! -d "$PROFILE_DIR" ]; then
    echo -e "${RED}Error: profiles directory not found${NC}"
    echo -e "${YELLOW}Run ./profile_memory.sh first${NC}"
    exit 1
fi

# Find the most recent profile set
LATEST_TIMESTAMP=$(ls -1 "$PROFILE_DIR" | grep -o '[0-9]\{8\}_[0-9]\{6\}' | sort -u | tail -1)

if [ -z "$LATEST_TIMESTAMP" ]; then
    echo -e "${RED}Error: No profiles found${NC}"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Memory Profile Analysis${NC}"
echo -e "${BLUE}Timestamp: ${LATEST_TIMESTAMP}${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

REPORT_DIR="${PROFILE_DIR}/reports_${LATEST_TIMESTAMP}"
mkdir -p "$REPORT_DIR"

# Function to analyze a profile
analyze_profile() {
    local profile=$1
    local name=$(basename "$profile" .prof)
    local output="${REPORT_DIR}/${name}_analysis.txt"

    echo -e "${YELLOW}Analyzing: ${name}${NC}"

    {
        echo "============================================"
        echo "Memory Profile Analysis: ${name}"
        echo "============================================"
        echo ""

        echo "=== Top 20 Memory Allocations (by total allocated) ==="
        go tool pprof -text -alloc_space "$profile" 2>/dev/null | head -25
        echo ""

        echo "=== Top 20 Memory Allocations (by objects allocated) ==="
        go tool pprof -text -alloc_objects "$profile" 2>/dev/null | head -25
        echo ""

        echo "=== Top 20 In-Use Memory (by size) ==="
        go tool pprof -text -inuse_space "$profile" 2>/dev/null | head -25
        echo ""

        echo "=== Top 20 In-Use Memory (by objects) ==="
        go tool pprof -text -inuse_objects "$profile" 2>/dev/null | head -25
        echo ""

    } > "$output"

    echo -e "${GREEN}✓ Report saved: ${output}${NC}"
}

# Function to generate visual profiles
generate_visuals() {
    local profile=$1
    local name=$(basename "$profile" .prof)

    echo -e "${YELLOW}Generating visuals: ${name}${NC}"

    # Generate SVG flame graph for allocated space
    go tool pprof -svg -alloc_space "$profile" > "${REPORT_DIR}/${name}_alloc_space.svg" 2>/dev/null || true

    # Generate SVG flame graph for allocated objects
    go tool pprof -svg -alloc_objects "$profile" > "${REPORT_DIR}/${name}_alloc_objects.svg" 2>/dev/null || true

    # Generate PDF report (if supported)
    go tool pprof -pdf -alloc_space "$profile" > "${REPORT_DIR}/${name}_alloc_space.pdf" 2>/dev/null || true

    echo -e "${GREEN}✓ Visuals saved${NC}"
}

# Analyze all memory profiles
echo -e "${BLUE}=== Analyzing Memory Profiles ===${NC}"
for profile in "${PROFILE_DIR}"/*_mem_${LATEST_TIMESTAMP}.prof; do
    if [ -f "$profile" ]; then
        analyze_profile "$profile"
        generate_visuals "$profile"
    fi
done

echo ""

# Generate summary report
SUMMARY="${REPORT_DIR}/MEMORY_SUMMARY.md"
echo -e "${YELLOW}Generating summary report...${NC}"

{
    echo "# Memory Profile Summary Report"
    echo ""
    echo "**Timestamp**: ${LATEST_TIMESTAMP}"
    echo "**Platform**: $(go version)"
    echo ""
    echo "## Overview"
    echo ""
    echo "This report summarizes memory allocation patterns across all PHP-Go components."
    echo ""
    echo "## Benchmark Results"
    echo ""

    for bench_file in "${PROFILE_DIR}"/*_bench_${LATEST_TIMESTAMP}.txt; do
        if [ -f "$bench_file" ]; then
            name=$(basename "$bench_file" "_bench_${LATEST_TIMESTAMP}.txt")
            echo "### ${name}"
            echo ""
            echo '```'
            cat "$bench_file"
            echo '```'
            echo ""
        fi
    done

    echo "## Top Memory Allocators"
    echo ""
    echo "### VM/Runtime Components"
    echo ""

    # Parse and extract top allocators from VM profiles
    for profile in "${PROFILE_DIR}"/vm_*_mem_${LATEST_TIMESTAMP}.prof; do
        if [ -f "$profile" ]; then
            name=$(basename "$profile" "_mem_${LATEST_TIMESTAMP}.prof")
            echo "#### ${name}"
            echo ""
            echo '```'
            go tool pprof -text -alloc_space "$profile" 2>/dev/null | head -15
            echo '```'
            echo ""
        fi
    done

    echo "### Parser Components"
    echo ""
    for profile in "${PROFILE_DIR}"/parser_*_mem_${LATEST_TIMESTAMP}.prof; do
        if [ -f "$profile" ]; then
            name=$(basename "$profile" "_mem_${LATEST_TIMESTAMP}.prof")
            echo "#### ${name}"
            echo ""
            echo '```'
            go tool pprof -text -alloc_space "$profile" 2>/dev/null | head -15
            echo '```'
            echo ""
        fi
    done

    echo "## Analysis Files"
    echo ""
    echo "Detailed analysis files are available in: \`${REPORT_DIR}/\`"
    echo ""
    echo "### Text Reports"
    for file in "${REPORT_DIR}"/*_analysis.txt; do
        if [ -f "$file" ]; then
            echo "- \`$(basename "$file")\`"
        fi
    done
    echo ""

    echo "### Visual Reports (SVG/PDF)"
    for file in "${REPORT_DIR}"/*.svg "${REPORT_DIR}"/*.pdf; do
        if [ -f "$file" ]; then
            echo "- \`$(basename "$file")\`"
        fi
    done
    echo ""

    echo "## How to View Profiles Interactively"
    echo ""
    echo '```bash'
    echo "# Start interactive web UI"
    echo "go tool pprof -http=:8080 ${PROFILE_DIR}/<profile>.prof"
    echo ""
    echo "# View specific profile"
    echo "go tool pprof -http=:8080 ${PROFILE_DIR}/vm_simple_loop_mem_${LATEST_TIMESTAMP}.prof"
    echo '```'
    echo ""

    echo "## Next Steps"
    echo ""
    echo "1. Review top allocators in each component"
    echo "2. Identify optimization opportunities (see BOTTLENECK_ANALYSIS.md)"
    echo "3. Implement pooling for frequently allocated objects"
    echo "4. Re-profile after optimizations to measure improvements"
    echo ""

} > "$SUMMARY"

echo -e "${GREEN}✓ Summary report saved: ${SUMMARY}${NC}"
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Analysis Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Reports saved in: ${REPORT_DIR}/${NC}"
echo ""
echo -e "${YELLOW}View summary:${NC}"
echo -e "  ${GREEN}cat ${SUMMARY}${NC}"
echo ""
echo -e "${YELLOW}View detailed reports:${NC}"
echo -e "  ${GREEN}ls ${REPORT_DIR}/${NC}"
echo ""
echo -e "${YELLOW}Interactive profile viewer:${NC}"
echo -e "  ${GREEN}go tool pprof -http=:8080 ${PROFILE_DIR}/vm_simple_loop_mem_${LATEST_TIMESTAMP}.prof${NC}"
echo ""
