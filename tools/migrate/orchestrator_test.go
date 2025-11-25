package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMigrationOrchestrator(t *testing.T) {
	projectPath := "/tmp/test-project"
	projectName := "TestProject"

	o := NewMigrationOrchestrator(projectPath, projectName, false)

	if o.ProjectPath != projectPath {
		t.Errorf("Expected ProjectPath %s, got %s", projectPath, o.ProjectPath)
	}

	if o.ProjectName != projectName {
		t.Errorf("Expected ProjectName %s, got %s", projectName, o.ProjectName)
	}

	expectedOutputDir := filepath.Join(projectPath, ".php-go-migration")
	if o.OutputDir != expectedOutputDir {
		t.Errorf("Expected OutputDir %s, got %s", expectedOutputDir, o.OutputDir)
	}

	if len(o.Steps) == 0 {
		t.Error("Expected steps to be initialized")
	}

	if o.Logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestNewMigrationOrchestratorDefaultName(t *testing.T) {
	projectPath := "/tmp/my-project"
	o := NewMigrationOrchestrator(projectPath, "", false)

	expectedName := "my-project"
	if o.ProjectName != expectedName {
		t.Errorf("Expected ProjectName %s, got %s", expectedName, o.ProjectName)
	}
}

func TestInitializeSteps(t *testing.T) {
	o := NewMigrationOrchestrator("/tmp/test", "test", false)

	expectedSteps := []string{
		"create-output-dir",
		"run-analyzer",
		"convert-config",
		"create-checklist",
		"test-parse",
		"generate-report",
	}

	if len(o.Steps) != len(expectedSteps) {
		t.Fatalf("Expected %d steps, got %d", len(expectedSteps), len(o.Steps))
	}

	for i, expectedID := range expectedSteps {
		if o.Steps[i].ID != expectedID {
			t.Errorf("Step %d: expected ID %s, got %s", i, expectedID, o.Steps[i].ID)
		}

		if o.Steps[i].Status != StepPending {
			t.Errorf("Step %d: expected status %s, got %s", i, StepPending, o.Steps[i].Status)
		}

		if o.Steps[i].Action == nil {
			t.Errorf("Step %d: expected Action to be set", i)
		}
	}
}

func TestCreateOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	err := o.createOutputDir(o)
	if err != nil {
		t.Fatalf("createOutputDir failed: %v", err)
	}

	// Check that directory exists
	info, err := os.Stat(o.OutputDir)
	if err != nil {
		t.Fatalf("Output directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Output path is not a directory")
	}
}

func TestRunAnalyzer(t *testing.T) {
	// Create a temporary project with PHP files
	tmpDir := t.TempDir()

	// Create some test PHP files
	testFiles := map[string]string{
		"test1.php": "<?php echo 'hello'; ?>",
		"test2.php": "<?php $x = 1 + 2; ?>",
		"test3.php": "<?php function test() { return true; } ?>",
	}

	for name, content := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	o := NewMigrationOrchestrator(tmpDir, "test", false)
	if err := o.createOutputDir(o); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	err := o.runAnalyzer(o)
	if err != nil {
		t.Fatalf("runAnalyzer failed: %v", err)
	}

	// Check that results were stored
	analyzer, ok := o.Results["analyzer"].(map[string]interface{})
	if !ok {
		t.Fatal("analyzer results not stored")
	}

	if analyzer["total_files"] == nil {
		t.Error("total_files not set in results")
	}

	// Check that report file was created
	reportPath := filepath.Join(o.OutputDir, "compatibility-report.txt")
	if _, err := os.Stat(reportPath); err != nil {
		t.Errorf("Report file was not created: %v", err)
	}
}

func TestConvertConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test php.ini file
	phpIniContent := `[PHP]
error_reporting = 32767
display_errors = On
memory_limit = 256M
`
	phpIniPath := filepath.Join(tmpDir, "php.ini")
	if err := os.WriteFile(phpIniPath, []byte(phpIniContent), 0644); err != nil {
		t.Fatalf("Failed to create php.ini: %v", err)
	}

	o := NewMigrationOrchestrator(tmpDir, "test", false)
	if err := o.createOutputDir(o); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	err := o.convertConfig(o)
	if err != nil {
		t.Fatalf("convertConfig failed: %v", err)
	}

	// Check that config file was created
	if _, err := os.Stat(o.ConfigFile); err != nil {
		t.Errorf("Config file was not created: %v", err)
	}

	// Check that results were stored
	config, ok := o.Results["config"].(map[string]interface{})
	if !ok {
		t.Fatal("config results not stored")
	}

	if config["source"] == nil {
		t.Error("source not set in results")
	}

	if config["destination"] == nil {
		t.Error("destination not set in results")
	}
}

func TestConvertConfigNoPhpIni(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	if err := o.createOutputDir(o); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Should not fail when no php.ini is found
	err := o.convertConfig(o)
	if err != nil {
		t.Errorf("convertConfig should not fail when php.ini is not found: %v", err)
	}

	// Note: The test may find a system php.ini if PHP is installed
	// So we just check that the function doesn't error
}

func TestCreateChecklist(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	if err := o.createOutputDir(o); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	err := o.createChecklist(o)
	if err != nil {
		t.Fatalf("createChecklist failed: %v", err)
	}

	// Check that checklist file was created
	if _, err := os.Stat(o.ChecklistFile); err != nil {
		t.Errorf("Checklist file was not created: %v", err)
	}

	// Check that results were stored
	checklist, ok := o.Results["checklist"].(map[string]interface{})
	if !ok {
		t.Fatal("checklist results not stored")
	}

	if checklist["path"] == nil {
		t.Error("path not set in results")
	}

	if checklist["total_items"] == nil {
		t.Error("total_items not set in results")
	}
}

func TestCollectPhpFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure with PHP files
	subdirs := []string{"src", "tests", "vendor"}
	for _, dir := range subdirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Create PHP files
	files := map[string]string{
		"index.php":             "<?php echo 'test'; ?>",
		"src/app.php":           "<?php class App {} ?>",
		"tests/test.php":        "<?php test(); ?>",
		"vendor/autoload.php":   "<?php // autoload ?>",
		"README.md":             "# README",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	o := NewMigrationOrchestrator(tmpDir, "test", false)
	collected, err := o.collectPhpFiles()
	if err != nil {
		t.Fatalf("collectPhpFiles failed: %v", err)
	}

	// Should find at least index.php and src/app.php
	// (vendor should be excluded by default, but tests directory might be included
	// depending on exclude patterns)
	if len(collected) < 2 {
		t.Errorf("Expected at least 2 files, got %d", len(collected))
	}

	// Check that vendor files are excluded
	for _, file := range collected {
		if strings.Contains(file, "vendor") {
			t.Errorf("vendor files should be excluded: %s", file)
		}
	}

	// Check that non-PHP files are excluded
	for _, file := range collected {
		if !strings.HasSuffix(file, ".php") {
			t.Errorf("non-PHP files should be excluded: %s", file)
		}
	}
}

func TestOrchestratorGenerateReport(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	if err := o.createOutputDir(o); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Add some test results
	o.Results["analyzer"] = map[string]interface{}{
		"total_files":        100,
		"parse_success_rate": 95.5,
		"total_issues":       25,
	}

	o.Results["parse_test"] = map[string]interface{}{
		"total_files":   100,
		"success_count": 95,
		"fail_count":    5,
		"success_rate":  95.0,
	}

	err := o.generateReport(o)
	if err != nil {
		t.Fatalf("generateReport failed: %v", err)
	}

	// Check that markdown report was created
	reportPath := filepath.Join(o.OutputDir, "migration-report.md")
	if _, err := os.Stat(reportPath); err != nil {
		t.Errorf("Markdown report was not created: %v", err)
	}

	// Check that JSON report was created
	jsonPath := filepath.Join(o.OutputDir, "migration-report.json")
	if _, err := os.Stat(jsonPath); err != nil {
		t.Errorf("JSON report was not created: %v", err)
	}

	// Check report content
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report: %v", err)
	}

	reportStr := string(content)
	if !strings.Contains(reportStr, "PHP-Go Migration Report") {
		t.Error("Report should contain title")
	}

	if !strings.Contains(reportStr, "test") {
		t.Error("Report should contain project name")
	}

	if !strings.Contains(reportStr, "95.5") {
		t.Error("Report should contain parse success rate")
	}
}

func TestBuildReport(t *testing.T) {
	o := NewMigrationOrchestrator("/tmp/test", "TestProject", false)

	// Add test results
	o.Results["analyzer"] = map[string]interface{}{
		"total_files":             100,
		"parse_success_rate":      95.5,
		"total_issues":            25,
		"unsupported_extensions":  3,
		"unsupported_functions":   10,
		"report_path":             "/tmp/test/report.txt",
	}

	report := o.buildReport()

	if !strings.Contains(report, "TestProject") {
		t.Error("Report should contain project name")
	}

	if !strings.Contains(report, "100") {
		t.Error("Report should contain total files count")
	}

	if !strings.Contains(report, "95.5") {
		t.Error("Report should contain parse success rate")
	}

	if !strings.Contains(report, "Next Steps") {
		t.Error("Report should contain next steps section")
	}

	if !strings.Contains(report, "Generated Files") {
		t.Error("Report should contain generated files section")
	}
}

func TestRunStepSuccess(t *testing.T) {
	o := NewMigrationOrchestrator("/tmp/test", "test", false)

	step := &MigrationStep{
		ID:          "test-step",
		Name:        "Test Step",
		Description: "A test step",
		Status:      StepPending,
		Action: func(o *MigrationOrchestrator) error {
			return nil
		},
	}

	err := o.runStep(step)
	if err != nil {
		t.Fatalf("runStep should not fail: %v", err)
	}

	if step.Status != StepCompleted {
		t.Errorf("Expected status %s, got %s", StepCompleted, step.Status)
	}

	if step.StartTime.IsZero() {
		t.Error("StartTime should be set")
	}

	if step.EndTime.IsZero() {
		t.Error("EndTime should be set")
	}
}

func TestRunStepFailure(t *testing.T) {
	o := NewMigrationOrchestrator("/tmp/test", "test", false)

	step := &MigrationStep{
		ID:          "test-step",
		Name:        "Test Step",
		Description: "A test step",
		Status:      StepPending,
		Action: func(o *MigrationOrchestrator) error {
			return &os.PathError{Op: "test", Path: "/test", Err: os.ErrNotExist}
		},
	}

	err := o.runStep(step)
	if err == nil {
		t.Fatal("runStep should fail")
	}

	if step.Status != StepFailed {
		t.Errorf("Expected status %s, got %s", StepFailed, step.Status)
	}

	if step.Error == nil {
		t.Error("Error should be set")
	}
}

func TestLogger(t *testing.T) {
	// Test non-verbose logger
	logger := NewLogger(false)
	if logger.Verbose {
		t.Error("Logger should not be verbose")
	}

	// Test verbose logger
	verboseLogger := NewLogger(true)
	if !verboseLogger.Verbose {
		t.Error("Logger should be verbose")
	}

	// Basic smoke test for logging functions (no assertions, just checking they don't panic)
	logger.Info("test info")
	logger.Debug("test debug") // Should not print
	logger.Warn("test warn")
	logger.Error("test error")

	verboseLogger.Debug("test debug") // Should print
}

func TestFindPhpIni(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	// Test with php.ini in project root
	phpIniPath := filepath.Join(tmpDir, "php.ini")
	if err := os.WriteFile(phpIniPath, []byte("[PHP]"), 0644); err != nil {
		t.Fatalf("Failed to create php.ini: %v", err)
	}

	found := o.findPhpIni()
	if found != phpIniPath {
		t.Errorf("Expected to find %s, got %s", phpIniPath, found)
	}
}

func TestFindPhpIniNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	o := NewMigrationOrchestrator(tmpDir, "test", false)

	found := o.findPhpIni()
	// Should return empty string or system php.ini
	// We can't test the exact value as it depends on the system
	if found != "" {
		// If found, it should be a valid path
		if _, err := os.Stat(found); err != nil {
			t.Errorf("Found path should exist: %s", found)
		}
	}
}

// Benchmark tests
func BenchmarkNewMigrationOrchestrator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewMigrationOrchestrator("/tmp/test", "test", false)
	}
}

func BenchmarkCollectPhpFiles(b *testing.B) {
	tmpDir := b.TempDir()

	// Create some test files
	for i := 0; i < 100; i++ {
		path := filepath.Join(tmpDir, "test"+string(rune(i))+".php")
		os.WriteFile(path, []byte("<?php ?>"), 0644)
	}

	o := NewMigrationOrchestrator(tmpDir, "test", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = o.collectPhpFiles()
	}
}
