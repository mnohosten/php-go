package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// IncludeManager manages file inclusion for require/include operations
type IncludeManager struct {
	mu sync.RWMutex

	// includePaths is the list of directories to search for included files
	// Similar to PHP's include_path configuration
	includePaths []string

	// includedFiles tracks files that have been included (for _once variants)
	// Maps canonicalized absolute path -> true
	includedFiles map[string]bool

	// currentScriptDir is the directory of the currently executing script
	currentScriptDir string

	// baseDir is the base directory for relative path resolution
	baseDir string
}

// NewIncludeManager creates a new include manager
func NewIncludeManager() *IncludeManager {
	cwd, _ := os.Getwd()
	return &IncludeManager{
		includePaths:  []string{cwd},
		includedFiles: make(map[string]bool),
		baseDir:       cwd,
	}
}

// SetIncludePaths sets the include path directories (like PHP's include_path)
func (im *IncludeManager) SetIncludePaths(paths []string) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.includePaths = paths
}

// AddIncludePath adds a directory to the include path
func (im *IncludeManager) AddIncludePath(path string) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.includePaths = append(im.includePaths, path)
}

// GetIncludePaths returns the current include paths
func (im *IncludeManager) GetIncludePaths() []string {
	im.mu.RLock()
	defer im.mu.RUnlock()
	paths := make([]string, len(im.includePaths))
	copy(paths, im.includePaths)
	return paths
}

// SetCurrentScriptDir sets the directory of the currently executing script
// This is used for resolving relative paths
func (im *IncludeManager) SetCurrentScriptDir(dir string) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.currentScriptDir = dir
}

// SetBaseDir sets the base directory for relative path resolution
func (im *IncludeManager) SetBaseDir(dir string) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.baseDir = dir
}

// ResolvePath resolves a file path for include/require operations
// It handles:
// - Absolute paths (returned as-is after canonicalization)
// - Relative paths starting with ./ or ../ (resolved relative to current script dir)
// - Other paths (searched in include paths)
func (im *IncludeManager) ResolvePath(path string) (string, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	// Handle absolute paths
	if filepath.IsAbs(path) {
		return im.canonicalizePath(path)
	}

	// Handle paths starting with ./ or ../
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") ||
	   strings.HasPrefix(path, ".\\") || strings.HasPrefix(path, "..\\") {
		// Resolve relative to current script directory
		baseDir := im.currentScriptDir
		if baseDir == "" {
			baseDir = im.baseDir
		}
		fullPath := filepath.Join(baseDir, path)
		return im.canonicalizePath(fullPath)
	}

	// Search in include paths
	for _, includePath := range im.includePaths {
		fullPath := filepath.Join(includePath, path)
		canonPath, err := im.canonicalizePath(fullPath)
		if err != nil {
			continue
		}

		// Check if file exists
		if _, err := os.Stat(canonPath); err == nil {
			return canonPath, nil
		}
	}

	// If not found in include paths, try relative to current script directory
	baseDir := im.currentScriptDir
	if baseDir == "" {
		baseDir = im.baseDir
	}
	fullPath := filepath.Join(baseDir, path)
	return im.canonicalizePath(fullPath)
}

// canonicalizePath converts a path to its canonical absolute form
// It resolves . and .. components and normalizes the path
func (im *IncludeManager) canonicalizePath(path string) (string, error) {
	// Clean the path (resolves . and .., removes redundant separators)
	cleanPath := filepath.Clean(path)

	// Convert to absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", err
	}

	// Evaluate symlinks to get the real path
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If the file doesn't exist yet, EvalSymlinks fails
		// In this case, return the absolute path without resolving symlinks
		// This matches PHP's behavior where the path is canonicalized even if file doesn't exist
		return absPath, nil
	}

	return realPath, nil
}

// MarkIncluded marks a file as included (for _once variants)
// Returns true if the file was already included, false otherwise
func (im *IncludeManager) MarkIncluded(path string) bool {
	im.mu.Lock()
	defer im.mu.Unlock()

	canonPath, err := im.canonicalizePath(path)
	if err != nil {
		// If we can't canonicalize, use the path as-is
		canonPath = path
	}

	if im.includedFiles[canonPath] {
		return true // Already included
	}

	im.includedFiles[canonPath] = true
	return false // Not previously included
}

// IsIncluded checks if a file has been included
func (im *IncludeManager) IsIncluded(path string) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()

	canonPath, err := im.canonicalizePath(path)
	if err != nil {
		canonPath = path
	}

	return im.includedFiles[canonPath]
}

// ClearIncluded clears the included files cache
// This is mainly useful for testing
func (im *IncludeManager) ClearIncluded() {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.includedFiles = make(map[string]bool)
}

// GetIncludedFiles returns a list of all included files
func (im *IncludeManager) GetIncludedFiles() []string {
	im.mu.RLock()
	defer im.mu.RUnlock()

	files := make([]string, 0, len(im.includedFiles))
	for file := range im.includedFiles {
		files = append(files, file)
	}
	return files
}

// Global include manager instance
var (
	globalIncludeManager     *IncludeManager
	globalIncludeManagerOnce sync.Once
)

// GetGlobalIncludeManager returns the singleton global include manager
func GetGlobalIncludeManager() *IncludeManager {
	globalIncludeManagerOnce.Do(func() {
		globalIncludeManager = NewIncludeManager()
	})
	return globalIncludeManager
}
