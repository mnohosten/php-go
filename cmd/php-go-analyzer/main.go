package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/krizos/php-go/tools/migrate"
)

func main() {
	// Parse command-line flags
	var (
		rootPath    = flag.String("path", ".", "Root path to analyze")
		exclude     = flag.String("exclude", "vendor,node_modules,.git", "Comma-separated patterns to exclude")
		outputFile  = flag.String("output", "", "Output file for report (default: stdout)")
		showHelp    = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	// Create analyzer
	analyzer := migrate.NewCompatibilityAnalyzer(*rootPath)

	// Set exclude patterns
	if *exclude != "" {
		patterns := strings.Split(*exclude, ",")
		for i, p := range patterns {
			patterns[i] = strings.TrimSpace(p)
		}
		analyzer.SetExcludePatterns(patterns)
	}

	// Run analysis
	fmt.Fprintf(os.Stderr, "Analyzing PHP files in: %s\n", *rootPath)
	fmt.Fprintf(os.Stderr, "Excluding patterns: %s\n", *exclude)
	fmt.Fprintln(os.Stderr, "")

	err := analyzer.Analyze()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during analysis: %v\n", err)
		os.Exit(1)
	}

	// Generate report
	report := analyzer.GenerateReport()

	// Output report
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(report), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing report to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Report written to: %s\n", *outputFile)
	} else {
		fmt.Print(report)
	}
}

func printHelp() {
	fmt.Println("PHP-Go Compatibility Analyzer")
	fmt.Println("")
	fmt.Println("Analyzes PHP applications for compatibility with PHP-Go, identifying")
	fmt.Println("unsupported extensions, functions, and language features.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  php-go-analyzer [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -path string")
	fmt.Println("        Root path to analyze (default: current directory)")
	fmt.Println("  -exclude string")
	fmt.Println("        Comma-separated patterns to exclude (default: vendor,node_modules,.git)")
	fmt.Println("  -output string")
	fmt.Println("        Output file for report (default: stdout)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Analyze current directory")
	fmt.Println("  php-go-analyzer")
	fmt.Println("")
	fmt.Println("  # Analyze specific directory")
	fmt.Println("  php-go-analyzer -path /path/to/project")
	fmt.Println("")
	fmt.Println("  # Save report to file")
	fmt.Println("  php-go-analyzer -path /path/to/project -output report.txt")
	fmt.Println("")
	fmt.Println("  # Exclude additional patterns")
	fmt.Println("  php-go-analyzer -exclude vendor,tests,cache")
	fmt.Println("")
}
