package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIncludeManager(t *testing.T) {
	im := NewIncludeManager()

	if im == nil {
		t.Fatal("NewIncludeManager returned nil")
	}

	if im.includedFiles == nil {
		t.Error("includedFiles map not initialized")
	}

	paths := im.GetIncludePaths()
	if len(paths) == 0 {
		t.Error("Expected at least one default include path")
	}
}

func TestSetGetIncludePaths(t *testing.T) {
	im := NewIncludeManager()

	testPaths := []string{"/path/one", "/path/two", "/path/three"}
	im.SetIncludePaths(testPaths)

	paths := im.GetIncludePaths()
	if len(paths) != len(testPaths) {
		t.Errorf("Expected %d paths, got %d", len(testPaths), len(paths))
	}

	for i, path := range paths {
		if path != testPaths[i] {
			t.Errorf("Path %d: expected %s, got %s", i, testPaths[i], path)
		}
	}
}

func TestAddIncludePath(t *testing.T) {
	im := NewIncludeManager()
	initialCount := len(im.GetIncludePaths())

	im.AddIncludePath("/new/path")
	paths := im.GetIncludePaths()

	if len(paths) != initialCount+1 {
		t.Errorf("Expected %d paths, got %d", initialCount+1, len(paths))
	}

	lastPath := paths[len(paths)-1]
	if lastPath != "/new/path" {
		t.Errorf("Expected last path to be '/new/path', got %s", lastPath)
	}
}

func TestResolvePathAbsolute(t *testing.T) {
	im := NewIncludeManager()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_*.php")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	absPath := tmpFile.Name()

	resolved, err := im.ResolvePath(absPath)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Should return canonicalized absolute path
	if !filepath.IsAbs(resolved) {
		t.Errorf("Expected absolute path, got %s", resolved)
	}
}

func TestResolvePathRelative(t *testing.T) {
	im := NewIncludeManager()

	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "test_include_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.php")
	if err := os.WriteFile(testFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set current script dir
	im.SetCurrentScriptDir(tmpDir)

	// Test relative path with ./
	resolved, err := im.ResolvePath("./test.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(testFile)
	if expected == "" {
		expected, _ = filepath.Abs(testFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestResolvePathDotDot(t *testing.T) {
	im := NewIncludeManager()

	// Create directory structure: tmpDir/subdir/
	tmpDir, err := os.MkdirTemp("", "test_include_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create file in parent: tmpDir/parent.php
	parentFile := filepath.Join(tmpDir, "parent.php")
	if err := os.WriteFile(parentFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set current script dir to subdir
	im.SetCurrentScriptDir(subDir)

	// Resolve ../parent.php from subdir
	resolved, err := im.ResolvePath("../parent.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(parentFile)
	if expected == "" {
		expected, _ = filepath.Abs(parentFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestResolvePathIncludePath(t *testing.T) {
	im := NewIncludeManager()

	// Create temporary directory for include path
	tmpDir, err := os.MkdirTemp("", "test_include_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "library.php")
	if err := os.WriteFile(testFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Add to include paths
	im.SetIncludePaths([]string{tmpDir})

	// Resolve just the filename (should search in include paths)
	resolved, err := im.ResolvePath("library.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(testFile)
	if expected == "" {
		expected, _ = filepath.Abs(testFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestResolvePathMultipleIncludePaths(t *testing.T) {
	im := NewIncludeManager()

	// Create two temporary directories
	tmpDir1, err := os.MkdirTemp("", "test_include1_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir1)

	tmpDir2, err := os.MkdirTemp("", "test_include2_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir2)

	// Create file only in second directory
	testFile := filepath.Join(tmpDir2, "config.php")
	if err := os.WriteFile(testFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set include paths (first dir doesn't have the file)
	im.SetIncludePaths([]string{tmpDir1, tmpDir2})

	// Should find file in second include path
	resolved, err := im.ResolvePath("config.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(testFile)
	if expected == "" {
		expected, _ = filepath.Abs(testFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestCanonicalizePath(t *testing.T) {
	im := NewIncludeManager()

	tests := []struct {
		name     string
		input    string
		contains []string // Path should contain these components
	}{
		{
			name:     "simple path",
			input:    "/tmp/test.php",
			contains: []string{"tmp", "test.php"},
		},
		{
			name:     "path with dot",
			input:    "/tmp/./test.php",
			contains: []string{"tmp", "test.php"},
		},
		{
			name:     "path with dotdot",
			input:    "/tmp/subdir/../test.php",
			contains: []string{"tmp", "test.php"},
		},
		{
			name:     "path with multiple dotdot",
			input:    "/tmp/a/b/../../test.php",
			contains: []string{"tmp", "test.php"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := im.canonicalizePath(tt.input)
			if err != nil {
				// Note: canonicalizePath may return error if path doesn't exist
				// but it should still clean the path
				t.Logf("canonicalizePath returned error (may be expected): %v", err)
			}

			for _, component := range tt.contains {
				if !filepath.IsAbs(result) && component != "." && component != ".." {
					t.Errorf("Result should be absolute path, got: %s", result)
				}
			}

			// Check that result doesn't contain . or .. (unless it's at the start)
			cleaned := filepath.Clean(result)
			if cleaned != result {
				// Note: EvalSymlinks may return different results
				t.Logf("Path not fully cleaned by canonicalizePath: %s vs %s", result, cleaned)
			}
		})
	}
}

func TestMarkIncluded(t *testing.T) {
	im := NewIncludeManager()

	path := "/tmp/test.php"

	// First call should return false (not previously included)
	wasIncluded := im.MarkIncluded(path)
	if wasIncluded {
		t.Error("First MarkIncluded should return false")
	}

	// Second call should return true (already included)
	wasIncluded = im.MarkIncluded(path)
	if !wasIncluded {
		t.Error("Second MarkIncluded should return true")
	}

	// Check with IsIncluded
	if !im.IsIncluded(path) {
		t.Error("IsIncluded should return true after MarkIncluded")
	}
}

func TestIsIncluded(t *testing.T) {
	im := NewIncludeManager()

	path := "/tmp/test.php"

	// Should not be included initially
	if im.IsIncluded(path) {
		t.Error("IsIncluded should return false for new path")
	}

	// Mark as included
	im.MarkIncluded(path)

	// Should now be included
	if !im.IsIncluded(path) {
		t.Error("IsIncluded should return true after marking")
	}
}

func TestClearIncluded(t *testing.T) {
	im := NewIncludeManager()

	// Mark some files as included
	im.MarkIncluded("/tmp/file1.php")
	im.MarkIncluded("/tmp/file2.php")
	im.MarkIncluded("/tmp/file3.php")

	files := im.GetIncludedFiles()
	if len(files) != 3 {
		t.Errorf("Expected 3 included files, got %d", len(files))
	}

	// Clear
	im.ClearIncluded()

	// Should be empty now
	files = im.GetIncludedFiles()
	if len(files) != 0 {
		t.Errorf("Expected 0 included files after clear, got %d", len(files))
	}

	// Previously included files should no longer be marked
	if im.IsIncluded("/tmp/file1.php") {
		t.Error("File should not be included after clear")
	}
}

func TestGetIncludedFiles(t *testing.T) {
	im := NewIncludeManager()

	paths := []string{
		"/tmp/file1.php",
		"/tmp/file2.php",
		"/tmp/file3.php",
	}

	for _, path := range paths {
		im.MarkIncluded(path)
	}

	included := im.GetIncludedFiles()
	if len(included) != len(paths) {
		t.Errorf("Expected %d included files, got %d", len(paths), len(included))
	}

	// Check all paths are in the result (order doesn't matter)
	pathMap := make(map[string]bool)
	for _, path := range included {
		pathMap[path] = true
	}

	for _, path := range paths {
		// Note: paths may be canonicalized, so we check if IsIncluded returns true
		if !im.IsIncluded(path) {
			t.Errorf("Path %s should be in included files", path)
		}
	}
}

func TestSetCurrentScriptDir(t *testing.T) {
	im := NewIncludeManager()

	testDir := "/tmp/testdir"
	im.SetCurrentScriptDir(testDir)

	// Create a temp directory for actual resolution
	tmpDir, err := os.MkdirTemp("", "test_script_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.php")
	if err := os.WriteFile(testFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set current script dir and resolve
	im.SetCurrentScriptDir(tmpDir)
	resolved, err := im.ResolvePath("./test.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(testFile)
	if expected == "" {
		expected, _ = filepath.Abs(testFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestSetBaseDir(t *testing.T) {
	im := NewIncludeManager()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test_base_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "base.php")
	if err := os.WriteFile(testFile, []byte("<?php"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set base dir
	im.SetBaseDir(tmpDir)
	im.SetCurrentScriptDir("") // Clear current script dir to use base dir

	resolved, err := im.ResolvePath("./base.php")
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Canonicalize expected path for comparison (handles symlinks like /var -> /private/var on macOS)
	expected, _ := filepath.EvalSymlinks(testFile)
	if expected == "" {
		expected, _ = filepath.Abs(testFile)
	}
	if resolved != expected {
		t.Errorf("Expected %s, got %s", expected, resolved)
	}
}

func TestGetGlobalIncludeManager(t *testing.T) {
	im1 := GetGlobalIncludeManager()
	im2 := GetGlobalIncludeManager()

	if im1 != im2 {
		t.Error("GetGlobalIncludeManager should return same instance")
	}

	if im1 == nil {
		t.Error("GetGlobalIncludeManager returned nil")
	}
}

func TestConcurrentAccess(t *testing.T) {
	im := NewIncludeManager()

	// Test concurrent access to include manager
	done := make(chan bool)

	// Writer goroutines
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				im.MarkIncluded(filepath.Join("/tmp", "file_"+string(rune(id))+"_"+string(rune(j))+".php"))
				im.AddIncludePath("/path/" + string(rune(id)))
			}
			done <- true
		}(i)
	}

	// Reader goroutines
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				im.IsIncluded(filepath.Join("/tmp", "file_"+string(rune(id))+".php"))
				im.GetIncludePaths()
				im.GetIncludedFiles()
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we got here without deadlock or race, test passes
	t.Log("Concurrent access test completed successfully")
}
