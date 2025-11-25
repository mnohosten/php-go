package migrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MigrationStep represents a single step in the automated migration process
type MigrationStep struct {
	ID          string
	Name        string
	Description string
	Action      func(*MigrationOrchestrator) error
	Required    bool
	Status      StepStatus
	Error       error
	StartTime   time.Time
	EndTime     time.Time
}

// StepStatus represents the status of a migration step
type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepSkipped   StepStatus = "skipped"
)

// MigrationOrchestrator orchestrates the automated migration process
type MigrationOrchestrator struct {
	ProjectPath     string
	ProjectName     string
	OutputDir       string
	ConfigFile      string
	ChecklistFile   string
	ExcludePatterns []string
	Steps           []*MigrationStep
	Results         map[string]interface{}
	Logger          *Logger
}

// Logger provides simple logging for the orchestrator
type Logger struct {
	Verbose bool
}

// NewLogger creates a new logger
func NewLogger(verbose bool) *Logger {
	return &Logger{Verbose: verbose}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Verbose {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

// NewMigrationOrchestrator creates a new migration orchestrator
func NewMigrationOrchestrator(projectPath, projectName string, verbose bool) *MigrationOrchestrator {
	if projectName == "" {
		projectName = filepath.Base(projectPath)
	}

	outputDir := filepath.Join(projectPath, ".php-go-migration")

	o := &MigrationOrchestrator{
		ProjectPath:     projectPath,
		ProjectName:     projectName,
		OutputDir:       outputDir,
		ConfigFile:      filepath.Join(outputDir, "php-go.yaml"),
		ChecklistFile:   filepath.Join(outputDir, "migration-checklist.json"),
		ExcludePatterns: []string{"vendor", "node_modules", ".git", "cache", "storage/logs"},
		Results:         make(map[string]interface{}),
		Logger:          NewLogger(verbose),
	}

	o.initializeSteps()
	return o
}

// initializeSteps sets up the migration steps
func (o *MigrationOrchestrator) initializeSteps() {
	o.Steps = []*MigrationStep{
		{
			ID:          "create-output-dir",
			Name:        "Create Output Directory",
			Description: "Create directory for migration artifacts",
			Action:      o.createOutputDir,
			Required:    true,
			Status:      StepPending,
		},
		{
			ID:          "run-analyzer",
			Name:        "Run Compatibility Analyzer",
			Description: "Analyze codebase for compatibility issues",
			Action:      o.runAnalyzer,
			Required:    true,
			Status:      StepPending,
		},
		{
			ID:          "convert-config",
			Name:        "Convert Configuration",
			Description: "Convert php.ini to PHP-Go YAML format",
			Action:      o.convertConfig,
			Required:    false,
			Status:      StepPending,
		},
		{
			ID:          "create-checklist",
			Name:        "Create Migration Checklist",
			Description: "Initialize migration tracking checklist",
			Action:      o.createChecklist,
			Required:    true,
			Status:      StepPending,
		},
		{
			ID:          "test-parse",
			Name:        "Test Parse Success Rate",
			Description: "Parse all PHP files and report success rate",
			Action:      o.testParse,
			Required:    true,
			Status:      StepPending,
		},
		{
			ID:          "generate-report",
			Name:        "Generate Migration Report",
			Description: "Generate comprehensive migration report",
			Action:      o.generateReport,
			Required:    true,
			Status:      StepPending,
		},
	}
}

// Run executes all migration steps
func (o *MigrationOrchestrator) Run() error {
	o.Logger.Info("Starting automated migration for project: %s", o.ProjectName)
	o.Logger.Info("Project path: %s", o.ProjectPath)
	o.Logger.Info("Output directory: %s", o.OutputDir)
	fmt.Println()

	for _, step := range o.Steps {
		if err := o.runStep(step); err != nil {
			if step.Required {
				o.Logger.Error("Required step '%s' failed: %v", step.Name, err)
				return fmt.Errorf("migration failed at step '%s': %w", step.Name, err)
			}
			o.Logger.Warn("Optional step '%s' failed: %v", step.Name, err)
			step.Status = StepSkipped
		}
	}

	o.Logger.Info("Migration process completed successfully!")
	o.printSummary()
	return nil
}

// runStep executes a single migration step
func (o *MigrationOrchestrator) runStep(step *MigrationStep) error {
	o.Logger.Info("Running step: %s", step.Name)
	o.Logger.Debug("  Description: %s", step.Description)

	step.Status = StepRunning
	step.StartTime = time.Now()

	err := step.Action(o)
	step.EndTime = time.Now()
	duration := step.EndTime.Sub(step.StartTime)

	if err != nil {
		step.Status = StepFailed
		step.Error = err
		o.Logger.Error("  Step failed after %v: %v", duration, err)
		return err
	}

	step.Status = StepCompleted
	o.Logger.Info("  Step completed in %v", duration)
	fmt.Println()
	return nil
}

// createOutputDir creates the output directory for migration artifacts
func (o *MigrationOrchestrator) createOutputDir(orch *MigrationOrchestrator) error {
	if err := os.MkdirAll(o.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	o.Logger.Debug("Created output directory: %s", o.OutputDir)
	return nil
}

// runAnalyzer runs the compatibility analyzer
func (o *MigrationOrchestrator) runAnalyzer(orch *MigrationOrchestrator) error {
	o.Logger.Info("  Analyzing project for compatibility issues...")

	analyzer := NewCompatibilityAnalyzer(o.ProjectPath)
	analyzer.SetExcludePatterns(o.ExcludePatterns)
	err := analyzer.Analyze()
	if err != nil {
		return fmt.Errorf("analyzer failed: %w", err)
	}

	report := analyzer.GenerateReport()

	// Save report to file
	reportPath := filepath.Join(o.OutputDir, "compatibility-report.txt")
	if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// Calculate parse success rate
	successRate := 0.0
	if len(analyzer.Files) > 0 {
		successCount := 0
		for _, f := range analyzer.Files {
			if f.ParseSuccess {
				successCount++
			}
		}
		successRate = float64(successCount) / float64(len(analyzer.Files)) * 100
	}

	// Store results
	o.Results["analyzer"] = map[string]interface{}{
		"total_files":             len(analyzer.Files),
		"successfully_parsed":     len(analyzer.Files) - len(analyzer.ParseErrors),
		"parse_failures":          len(analyzer.ParseErrors),
		"total_issues":            len(analyzer.Issues),
		"unsupported_extensions":  len(analyzer.UnsupportedExts),
		"unsupported_functions":   len(analyzer.UnsupportedFuncs),
		"parse_success_rate":      successRate,
		"report_path":             reportPath,
	}

	o.Logger.Info("  Analysis complete:")
	o.Logger.Info("    Total files: %d", len(analyzer.Files))
	o.Logger.Info("    Parse success rate: %.2f%%", successRate)
	o.Logger.Info("    Issues found: %d", len(analyzer.Issues))
	o.Logger.Info("    Report saved to: %s", reportPath)

	return nil
}

// convertConfig converts php.ini to PHP-Go YAML format
func (o *MigrationOrchestrator) convertConfig(orch *MigrationOrchestrator) error {
	// Find php.ini file
	phpIniPath := o.findPhpIni()
	if phpIniPath == "" {
		o.Logger.Warn("  No php.ini found, skipping configuration conversion")
		return nil
	}

	o.Logger.Info("  Converting configuration from: %s", phpIniPath)

	config, warnings, err := ConvertFile(phpIniPath)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Log warnings if any
	if len(warnings) > 0 {
		o.Logger.Warn("  Configuration conversion warnings:")
		for _, w := range warnings {
			o.Logger.Warn("    - %s", w)
		}
	}

	// Convert config to YAML
	yamlContent := config.ToYAML()

	// Save to file
	if err := os.WriteFile(o.ConfigFile, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	o.Results["config"] = map[string]interface{}{
		"source":      phpIniPath,
		"destination": o.ConfigFile,
	}

	o.Logger.Info("  Configuration converted and saved to: %s", o.ConfigFile)
	return nil
}

// findPhpIni searches for php.ini in common locations
func (o *MigrationOrchestrator) findPhpIni() string {
	// Check project directory
	candidates := []string{
		filepath.Join(o.ProjectPath, "php.ini"),
		filepath.Join(o.ProjectPath, "config", "php.ini"),
		filepath.Join(o.ProjectPath, ".php.ini"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try to get system php.ini location using php command
	cmd := exec.Command("php", "--ini")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Loaded Configuration File:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					path := strings.TrimSpace(parts[1])
					if path != "" && path != "(none)" {
						return path
					}
				}
			}
		}
	}

	return ""
}

// createChecklist creates a migration checklist
func (o *MigrationOrchestrator) createChecklist(orch *MigrationOrchestrator) error {
	o.Logger.Info("  Creating migration checklist...")

	checklist := NewMigrationChecklist(o.ProjectName)
	if err := checklist.Save(o.ChecklistFile); err != nil {
		return fmt.Errorf("failed to save checklist: %w", err)
	}

	o.Results["checklist"] = map[string]interface{}{
		"path":        o.ChecklistFile,
		"total_items": len(checklist.Items),
	}

	o.Logger.Info("  Checklist created with %d items", len(checklist.Items))
	o.Logger.Info("  Saved to: %s", o.ChecklistFile)
	return nil
}

// testParse tests the parse success rate
func (o *MigrationOrchestrator) testParse(orch *MigrationOrchestrator) error {
	o.Logger.Info("  Testing parse success rate...")

	// This will use the php-go binary if available
	phpGoBin := o.findPhpGoBinary()
	if phpGoBin == "" {
		o.Logger.Warn("  php-go binary not found, skipping parse test")
		o.Logger.Warn("  Run 'go build -o php-go ./cmd/php-go' to enable parse testing")
		return nil
	}

	// Collect all PHP files
	files, err := o.collectPhpFiles()
	if err != nil {
		return fmt.Errorf("failed to collect PHP files: %w", err)
	}

	o.Logger.Info("  Found %d PHP files to test", len(files))

	// Test parsing with php-go
	successCount := 0
	failCount := 0
	var failures []string

	for i, file := range files {
		if i > 0 && i%100 == 0 {
			o.Logger.Debug("  Progress: %d/%d files tested", i, len(files))
		}

		cmd := exec.Command(phpGoBin, "parse", file)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			failCount++
			failures = append(failures, file)
		} else {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(len(files)) * 100

	// Save failures to file
	if len(failures) > 0 {
		failuresPath := filepath.Join(o.OutputDir, "parse-failures.txt")
		content := strings.Join(failures, "\n")
		if err := os.WriteFile(failuresPath, []byte(content), 0644); err != nil {
			o.Logger.Warn("  Failed to save parse failures: %v", err)
		} else {
			o.Logger.Info("  Parse failures saved to: %s", failuresPath)
		}
	}

	o.Results["parse_test"] = map[string]interface{}{
		"total_files":    len(files),
		"success_count":  successCount,
		"fail_count":     failCount,
		"success_rate":   successRate,
	}

	o.Logger.Info("  Parse test complete:")
	o.Logger.Info("    Success: %d/%d (%.2f%%)", successCount, len(files), successRate)
	o.Logger.Info("    Failures: %d", failCount)

	return nil
}

// collectPhpFiles collects all PHP files in the project
func (o *MigrationOrchestrator) collectPhpFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(o.ProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded patterns
		for _, pattern := range o.ExcludePatterns {
			if strings.Contains(path, pattern) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Collect .php files
		if !info.IsDir() && strings.HasSuffix(path, ".php") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// findPhpGoBinary finds the php-go binary
func (o *MigrationOrchestrator) findPhpGoBinary() string {
	// Check in project root
	candidates := []string{
		filepath.Join(o.ProjectPath, "php-go"),
		filepath.Join(o.ProjectPath, "..", "php-go"),
		"./php-go",
		"php-go",
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	// Try PATH
	path, err := exec.LookPath("php-go")
	if err == nil {
		return path
	}

	return ""
}

// generateReport generates a comprehensive migration report
func (o *MigrationOrchestrator) generateReport(orch *MigrationOrchestrator) error {
	o.Logger.Info("  Generating migration report...")

	report := o.buildReport()
	reportPath := filepath.Join(o.OutputDir, "migration-report.md")

	if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// Also save as JSON for programmatic access
	jsonReport, err := json.MarshalIndent(o.Results, "", "  ")
	if err != nil {
		o.Logger.Warn("  Failed to generate JSON report: %v", err)
	} else {
		jsonPath := filepath.Join(o.OutputDir, "migration-report.json")
		if err := os.WriteFile(jsonPath, jsonReport, 0644); err != nil {
			o.Logger.Warn("  Failed to save JSON report: %v", err)
		} else {
			o.Logger.Info("  JSON report saved to: %s", jsonPath)
		}
	}

	o.Logger.Info("  Migration report saved to: %s", reportPath)
	return nil
}

// buildReport builds the migration report content
func (o *MigrationOrchestrator) buildReport() string {
	var b strings.Builder

	b.WriteString("# PHP-Go Migration Report\n\n")
	b.WriteString(fmt.Sprintf("**Project**: %s\n", o.ProjectName))
	b.WriteString(fmt.Sprintf("**Date**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	b.WriteString("## Summary\n\n")

	// Analyzer results
	if analyzer, ok := o.Results["analyzer"].(map[string]interface{}); ok {
		b.WriteString("### Compatibility Analysis\n\n")
		b.WriteString(fmt.Sprintf("- **Total Files**: %v\n", analyzer["total_files"]))
		b.WriteString(fmt.Sprintf("- **Parse Success Rate**: %.2f%%\n", analyzer["parse_success_rate"]))
		b.WriteString(fmt.Sprintf("- **Issues Found**: %v\n", analyzer["total_issues"]))
		b.WriteString(fmt.Sprintf("- **Unsupported Extensions**: %v\n", analyzer["unsupported_extensions"]))
		b.WriteString(fmt.Sprintf("- **Unsupported Functions**: %v\n", analyzer["unsupported_functions"]))
		b.WriteString(fmt.Sprintf("- **Report**: %v\n\n", analyzer["report_path"]))
	}

	// Parse test results
	if parseTest, ok := o.Results["parse_test"].(map[string]interface{}); ok {
		b.WriteString("### Parse Test Results\n\n")
		b.WriteString(fmt.Sprintf("- **Total Files Tested**: %v\n", parseTest["total_files"]))
		b.WriteString(fmt.Sprintf("- **Success Rate**: %.2f%%\n", parseTest["success_rate"]))
		b.WriteString(fmt.Sprintf("- **Successful**: %v\n", parseTest["success_count"]))
		b.WriteString(fmt.Sprintf("- **Failed**: %v\n\n", parseTest["fail_count"]))
	}

	// Configuration
	if config, ok := o.Results["config"].(map[string]interface{}); ok {
		b.WriteString("### Configuration\n\n")
		b.WriteString(fmt.Sprintf("- **Source**: %v\n", config["source"]))
		b.WriteString(fmt.Sprintf("- **Destination**: %v\n\n", config["destination"]))
	}

	// Checklist
	if checklist, ok := o.Results["checklist"].(map[string]interface{}); ok {
		b.WriteString("### Migration Checklist\n\n")
		b.WriteString(fmt.Sprintf("- **Total Tasks**: %v\n", checklist["total_items"]))
		b.WriteString(fmt.Sprintf("- **Checklist File**: %v\n\n", checklist["path"]))
	}

	b.WriteString("## Next Steps\n\n")
	b.WriteString("1. Review the compatibility report to understand migration challenges\n")
	b.WriteString("2. Check the parse test results to identify files that need fixes\n")
	b.WriteString("3. Review the generated configuration file and adjust as needed\n")
	b.WriteString("4. Use the migration checklist to track progress:\n")
	b.WriteString(fmt.Sprintf("   ```\n   php-go-migrate status -file %s\n   ```\n", o.ChecklistFile))
	b.WriteString("5. Address critical issues identified in the compatibility report\n")
	b.WriteString("6. Run functional tests to verify application behavior\n")
	b.WriteString("7. Consult the migration guide: `docs/user-guide/migration.md`\n\n")

	b.WriteString("## Generated Files\n\n")
	b.WriteString(fmt.Sprintf("- Compatibility Report: `%s`\n", filepath.Join(o.OutputDir, "compatibility-report.txt")))
	b.WriteString(fmt.Sprintf("- Migration Checklist: `%s`\n", o.ChecklistFile))
	if _, ok := o.Results["config"]; ok {
		b.WriteString(fmt.Sprintf("- Configuration: `%s`\n", o.ConfigFile))
	}
	if _, ok := o.Results["parse_test"]; ok {
		b.WriteString(fmt.Sprintf("- Parse Failures: `%s`\n", filepath.Join(o.OutputDir, "parse-failures.txt")))
	}
	b.WriteString(fmt.Sprintf("- JSON Report: `%s`\n", filepath.Join(o.OutputDir, "migration-report.json")))

	return b.String()
}

// printSummary prints a summary of the migration process
func (o *MigrationOrchestrator) printSummary() {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("           MIGRATION PROCESS COMPLETE              ")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Print step summary
	completedCount := 0
	failedCount := 0
	skippedCount := 0

	for _, step := range o.Steps {
		switch step.Status {
		case StepCompleted:
			completedCount++
		case StepFailed:
			failedCount++
		case StepSkipped:
			skippedCount++
		}
	}

	fmt.Printf("Steps Completed: %d/%d\n", completedCount, len(o.Steps))
	if failedCount > 0 {
		fmt.Printf("Steps Failed: %d\n", failedCount)
	}
	if skippedCount > 0 {
		fmt.Printf("Steps Skipped: %d\n", skippedCount)
	}
	fmt.Println()

	// Print key results
	if analyzer, ok := o.Results["analyzer"].(map[string]interface{}); ok {
		fmt.Printf("Parse Success Rate: %.2f%%\n", analyzer["parse_success_rate"])
		fmt.Printf("Issues Found: %v\n", analyzer["total_issues"])
	}
	fmt.Println()

	fmt.Printf("All migration artifacts saved to: %s\n", o.OutputDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Review the migration report: %s\n", filepath.Join(o.OutputDir, "migration-report.md"))
	fmt.Printf("  2. Check migration progress: php-go-migrate status -file %s\n", o.ChecklistFile)
	fmt.Println("  3. Consult the migration guide: docs/user-guide/migration.md")
	fmt.Println()
}
