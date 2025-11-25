package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/util"
)

// ============================================================================
// Security - File Permissions
// ============================================================================

// File permission constants for secure file operations
// These constants make permission modes explicit and easier to audit
const (
	// FilePermissionPrivate creates files readable/writable only by owner (0600)
	// Use for sensitive data like credentials, tokens, or user-specific files
	FilePermissionPrivate os.FileMode = 0600

	// FilePermissionPublic creates files readable by all, writable by owner (0644)
	// Use for general application files, logs, and publicly readable data
	// This is the default for PHP file operations
	FilePermissionPublic os.FileMode = 0644

	// DirPermissionPrivate creates directories accessible only by owner (0700)
	// Use for user-specific directories containing sensitive data
	DirPermissionPrivate os.FileMode = 0700

	// DirPermissionPublic creates directories with standard Unix permissions (0755)
	// Owner: read/write/execute, Others: read/execute
	// This is the default for PHP directory operations
	DirPermissionPublic os.FileMode = 0755
)

// ============================================================================
// Security - Path Traversal Protection
// ============================================================================

// PathValidationConfig holds configuration for path validation
type PathValidationConfig struct {
	// BasePath restricts file operations to this directory and subdirectories
	// Empty string means no restriction (use with caution in production)
	BasePath string

	// AllowAbsolutePaths allows absolute paths outside BasePath
	// Only applies when BasePath is set
	AllowAbsolutePaths bool

	// AllowSymlinks allows following symbolic links within the base directory
	// When true: symlinks are allowed but must resolve to paths within BasePath
	// When false: any symlink usage is blocked
	AllowSymlinks bool

	// MaxFileSize limits the maximum size of files that can be read (in bytes)
	// 0 means no limit (use with caution in production)
	// Recommended: 10MB (10485760) for general use, 100MB (104857600) for larger files
	MaxFileSize int64

	// MaxWriteSize limits the maximum size of data that can be written (in bytes)
	// 0 means no limit (use with caution in production)
	// Recommended: 10MB (10485760) for general use
	MaxWriteSize int64
}

// DefaultPathValidationConfig returns a secure default configuration
// By default, restricts to current working directory
var DefaultPathValidationConfig = PathValidationConfig{
	BasePath:           "", // Empty = no restriction (for backward compatibility)
	AllowAbsolutePaths: false,
	AllowSymlinks:      false,
	MaxFileSize:        0, // 0 = no limit (for backward compatibility)
	MaxWriteSize:       0, // 0 = no limit (for backward compatibility)
}

// globalPathConfig is the global path validation configuration
// Can be modified by applications for stricter security
var globalPathConfig = DefaultPathValidationConfig

// SetPathValidationConfig sets the global path validation configuration
func SetPathValidationConfig(config PathValidationConfig) {
	globalPathConfig = config
}

// GetPathValidationConfig returns the current global path validation configuration
func GetPathValidationConfig() PathValidationConfig {
	return globalPathConfig
}

// validatePath validates a file path for security issues
// Returns the cleaned absolute path or an error
func validatePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path not allowed")
	}

	// Clean the path to remove . and .. components
	cleaned := filepath.Clean(path)

	// Check for null bytes (security issue)
	if strings.Contains(cleaned, "\x00") {
		return "", errors.New("null byte in path not allowed")
	}

	// If no base path restriction, just return cleaned path
	if globalPathConfig.BasePath == "" {
		return cleaned, nil
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Get absolute base path
	absBase, err := filepath.Abs(globalPathConfig.BasePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base path: %w", err)
	}

	// Check if path is within base directory
	relPath, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	// Check for path traversal attempts
	if strings.HasPrefix(relPath, "..") {
		if !globalPathConfig.AllowAbsolutePaths {
			return "", fmt.Errorf("path traversal detected: %s escapes base directory %s", path, absBase)
		}
	}

	// Evaluate symlinks for both paths to ensure consistent comparison
	// This handles cases where system directories themselves are symlinks (e.g., /var -> /private/var)
	evalPath, err1 := filepath.EvalSymlinks(absPath)
	evalBase, err2 := filepath.EvalSymlinks(absBase)

	// Check symlink behavior based on configuration
	if err1 == nil && err2 == nil {
		// Always check that the resolved path stays within base directory
		// (applies whether symlinks are allowed or not)
		evalRel, err := filepath.Rel(evalBase, evalPath)
		if err != nil || strings.HasPrefix(evalRel, "..") {
			return "", fmt.Errorf("path escapes base directory: %s -> %s", path, evalPath)
		}

		// If symlinks are not allowed, check if the resolved path differs from original
		// Note: We only check this AFTER confirming the resolved path is within base
		// This prevents false positives from system-level symlinks (e.g., /var -> /private/var)
		if !globalPathConfig.AllowSymlinks {
			// Compare relative paths to detect user-introduced symlinks
			// Get relative path from base for both original and resolved
			origRel, err1 := filepath.Rel(absBase, absPath)
			evalRelPath, err2 := filepath.Rel(evalBase, evalPath)

			if err1 == nil && err2 == nil && origRel != evalRelPath {
				// The relative paths differ, meaning a symlink was followed within the base directory
				return "", fmt.Errorf("symlink detected and not allowed: %s -> %s", path, evalPath)
			}
		}
	}

	// Return the cleaned path (or absolute path if base path was configured)
	// This ensures file operations work correctly with both relative and absolute paths
	if globalPathConfig.BasePath != "" {
		return absPath, nil
	}
	return cleaned, nil
}

// validateFileSize checks if a file's size is within the configured limit
// Returns an error if the file is too large
func validateFileSize(path string) error {
	// If no limit is configured, allow any size
	if globalPathConfig.MaxFileSize <= 0 {
		return nil
	}

	// Get file info to check size
	info, err := os.Stat(path)
	if err != nil {
		// If file doesn't exist or can't be stat'd, let the actual read operation handle it
		return nil
	}

	// Check if file size exceeds limit
	if info.Size() > globalPathConfig.MaxFileSize {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", info.Size(), globalPathConfig.MaxFileSize)
	}

	return nil
}

// validateWriteSize checks if the data size is within the configured write limit
// Returns an error if the data is too large
func validateWriteSize(dataSize int64) error {
	// If no limit is configured, allow any size
	if globalPathConfig.MaxWriteSize <= 0 {
		return nil
	}

	// Check if write size exceeds limit
	if dataSize > globalPathConfig.MaxWriteSize {
		return fmt.Errorf("write size (%d bytes) exceeds maximum allowed size (%d bytes)", dataSize, globalPathConfig.MaxWriteSize)
	}

	return nil
}

// ============================================================================
// File Reading Functions
// ============================================================================

// FileGetContents reads entire file into a string
// file_get_contents(string $filename): string|false
func FileGetContents(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	// Check file size limit
	if err := validateFileSize(validPath); err != nil {
		// File too large - return false
		return types.NewBool(false)
	}

	data, err := os.ReadFile(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(string(data))
}

// FilePutContents writes a string to a file
// file_put_contents(string $filename, mixed $data, int $flags = 0): int|false
func FilePutContents(filename *types.Value, data *types.Value, args ...*types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	content := data.ToString()

	// Check write size limit
	if err := validateWriteSize(int64(len(content))); err != nil {
		// Data too large - return false
		return types.NewBool(false)
	}

	flags := 0
	if len(args) > 0 {
		flags = int(args[0].ToInt())
	}

	// FILE_APPEND flag (8)
	var writeFlags int
	if flags&8 != 0 {
		writeFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		writeFlags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	file, err := os.OpenFile(validPath, writeFlags, FilePermissionPublic)
	if err != nil {
		return types.NewBool(false)
	}
	defer file.Close()

	n, err := file.WriteString(content)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewInt(int64(n))
}

// File reads entire file into an array
// file(string $filename, int $flags = 0): array|false
func File(filename *types.Value, args ...*types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	// Check file size limit
	if err := validateFileSize(validPath); err != nil {
		// File too large - return false
		return types.NewBool(false)
	}

	data, err := os.ReadFile(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	flags := 0
	if len(args) > 0 {
		flags = int(args[0].ToInt())
	}

	// Split into lines
	content := string(data)
	lines := strings.Split(content, "\n")

	// FILE_IGNORE_NEW_LINES flag (2)
	skipNewlines := flags&2 != 0

	// FILE_SKIP_EMPTY_LINES flag (4)
	skipEmpty := flags&4 != 0

	arr := types.NewEmptyArray()
	for _, line := range lines {
		if skipEmpty && line == "" {
			continue
		}

		if skipNewlines {
			arr.Append(types.NewString(line))
		} else {
			arr.Append(types.NewString(line + "\n"))
		}
	}

	return types.NewArray(arr)
}

// Readfile outputs a file
// readfile(string $filename): int|false
func Readfile(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	// Check file size limit
	if err := validateFileSize(validPath); err != nil {
		// File too large - return false
		return types.NewBool(false)
	}

	data, err := os.ReadFile(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	// In real implementation, this would output to stdout
	// For now, just return the byte count
	return types.NewInt(int64(len(data)))
}

// ============================================================================
// File Handle Functions
// ============================================================================

// Fopen opens a file or URL
// fopen(string $filename, string $mode): resource|false
func Fopen(filename *types.Value, mode *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	modeStr := mode.ToString()

	var flags int
	switch modeStr {
	case "r":
		flags = os.O_RDONLY
	case "r+":
		flags = os.O_RDWR
	case "w":
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "w+":
		flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a":
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "a+":
		flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	case "x":
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	case "x+":
		flags = os.O_RDWR | os.O_CREATE | os.O_EXCL
	default:
		return types.NewBool(false)
	}

	file, err := os.OpenFile(validPath, flags, FilePermissionPublic)
	if err != nil {
		return types.NewBool(false)
	}

	resource := types.NewResourceHandle("file", file)
	return types.NewResource(resource)
}

// Fclose closes an open file pointer
// fclose(resource $stream): bool
func Fclose(stream *types.Value) *types.Value {
	if stream.Type() != types.TypeResource {
		return types.NewBool(false)
	}

	res := stream.ToResource()
	if res.Type() != "file" {
		return types.NewBool(false)
	}

	if file, ok := res.Data().(*os.File); ok {
		err := file.Close()
		return types.NewBool(err == nil)
	}

	return types.NewBool(false)
}

// Fread reads from file pointer
// fread(resource $stream, int $length): string|false
func Fread(stream *types.Value, length *types.Value) *types.Value {
	if stream.Type() != types.TypeResource {
		return types.NewBool(false)
	}

	res := stream.ToResource()
	if res.Type() != "file" {
		return types.NewBool(false)
	}

	file, ok := res.Data().(*os.File)
	if !ok {
		return types.NewBool(false)
	}

	// Safely convert int64 to int with overflow checking
	lengthInt64 := length.ToInt()
	n, err := util.SafeConvertToInt(lengthInt64, "read_length")
	if err != nil {
		return types.NewBool(false)
	}

	// Check if read size exceeds limit
	if err := validateWriteSize(int64(n)); err != nil {
		// Using validateWriteSize as it checks size limit for memory allocation
		return types.NewBool(false)
	}

	buf := make([]byte, n)

	bytesRead, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return types.NewBool(false)
	}

	return types.NewString(string(buf[:bytesRead]))
}

// Fwrite writes to file pointer
// fwrite(resource $stream, string $data, int $length = null): int|false
func Fwrite(stream *types.Value, data *types.Value, args ...*types.Value) *types.Value {
	if stream.Type() != types.TypeResource {
		return types.NewBool(false)
	}

	res := stream.ToResource()
	if res.Type() != "file" {
		return types.NewBool(false)
	}

	file, ok := res.Data().(*os.File)
	if !ok {
		return types.NewBool(false)
	}

	content := data.ToString()

	// Optional length parameter
	if len(args) > 0 {
		lengthInt64 := args[0].ToInt()
		lengthInt, err := util.SafeConvertToInt(lengthInt64, "write_length")
		if err != nil {
			return types.NewBool(false)
		}
		if lengthInt < len(content) {
			content = content[:lengthInt]
		}
	}

	// Check write size limit
	if err := validateWriteSize(int64(len(content))); err != nil {
		// Data too large - return false
		return types.NewBool(false)
	}

	n, err := file.WriteString(content)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewInt(int64(n))
}

// Fgets reads line from file pointer
// fgets(resource $stream, int $length = null): string|false
func Fgets(stream *types.Value, args ...*types.Value) *types.Value {
	if stream.Type() != types.TypeResource {
		return types.NewBool(false)
	}

	res := stream.ToResource()
	if res.Type() != "file" {
		return types.NewBool(false)
	}

	file, ok := res.Data().(*os.File)
	if !ok {
		return types.NewBool(false)
	}

	// Read one byte at a time until newline
	var line strings.Builder
	buf := make([]byte, 1)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF && line.Len() > 0 {
				break
			}
			return types.NewBool(false)
		}

		if n > 0 {
			line.WriteByte(buf[0])
			if buf[0] == '\n' {
				break
			}
		}
	}

	return types.NewString(line.String())
}

// Fgetc reads character from file pointer
// fgetc(resource $stream): string|false
func Fgetc(stream *types.Value) *types.Value {
	if stream.Type() != types.TypeResource {
		return types.NewBool(false)
	}

	res := stream.ToResource()
	if res.Type() != "file" {
		return types.NewBool(false)
	}

	file, ok := res.Data().(*os.File)
	if !ok {
		return types.NewBool(false)
	}

	buf := make([]byte, 1)
	n, err := file.Read(buf)
	if err != nil || n == 0 {
		return types.NewBool(false)
	}

	return types.NewString(string(buf[0]))
}

// ============================================================================
// File Information Functions
// ============================================================================

// FileExists checks whether a file or directory exists
// file_exists(string $filename): bool
func FileExists(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	_, err = os.Stat(validPath)
	return types.NewBool(err == nil)
}

// IsFile tells whether the filename is a regular file
// is_file(string $filename): bool
func IsFile(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	info, err := os.Stat(validPath)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewBool(info.Mode().IsRegular())
}

// IsDir tells whether the filename is a directory
// is_dir(string $filename): bool
func IsDir(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	info, err := os.Stat(validPath)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewBool(info.IsDir())
}

// IsReadable tells whether a file exists and is readable
// is_readable(string $filename): bool
func IsReadable(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	file, err := os.OpenFile(validPath, os.O_RDONLY, 0)
	if err != nil {
		return types.NewBool(false)
	}
	file.Close()
	return types.NewBool(true)
}

// IsWritable tells whether the filename is writable
// is_writable(string $filename): bool
func IsWritable(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	// Check if file exists
	info, err := os.Stat(validPath)
	if err != nil {
		// File doesn't exist, check if directory is writable
		dir := filepath.Dir(validPath)
		dirInfo, err := os.Stat(dir)
		if err != nil {
			return types.NewBool(false)
		}
		return types.NewBool(dirInfo.Mode().Perm()&0200 != 0)
	}

	// File exists, check if writable
	return types.NewBool(info.Mode().Perm()&0200 != 0)
}

// Filesize gets file size
// filesize(string $filename): int|false
func Filesize(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	info, err := os.Stat(validPath)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewInt(info.Size())
}

// Filetype gets file type
// filetype(string $filename): string|false
func Filetype(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	info, err := os.Stat(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	mode := info.Mode()
	switch {
	case mode.IsRegular():
		return types.NewString("file")
	case mode.IsDir():
		return types.NewString("dir")
	case mode&os.ModeSymlink != 0:
		return types.NewString("link")
	case mode&os.ModeNamedPipe != 0:
		return types.NewString("fifo")
	case mode&os.ModeCharDevice != 0:
		return types.NewString("char")
	case mode&os.ModeDevice != 0:
		return types.NewString("block")
	default:
		return types.NewString("unknown")
	}
}

// ============================================================================
// Directory Functions
// ============================================================================

// Mkdir makes directory
// mkdir(string $directory, int $permissions = 0777, bool $recursive = false): bool
func Mkdir(directory *types.Value, args ...*types.Value) *types.Value {
	path := directory.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	permissions := os.FileMode(0777)
	if len(args) > 0 {
		permissions = os.FileMode(args[0].ToInt())
	}

	recursive := false
	if len(args) > 1 {
		recursive = args[1].ToBool()
	}

	if recursive {
		err = os.MkdirAll(validPath, permissions)
	} else {
		err = os.Mkdir(validPath, permissions)
	}

	return types.NewBool(err == nil)
}

// Rmdir removes directory
// rmdir(string $directory): bool
func Rmdir(directory *types.Value) *types.Value {
	path := directory.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	err = os.Remove(validPath)
	return types.NewBool(err == nil)
}

// Scandir lists files and directories inside the specified path
// scandir(string $directory, int $sorting_order = SCANDIR_SORT_ASCENDING): array|false
func Scandir(directory *types.Value, args ...*types.Value) *types.Value {
	path := directory.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	arr := types.NewEmptyArray()

	// Add . and ..
	arr.Append(types.NewString("."))
	arr.Append(types.NewString(".."))

	// Add all entries
	for _, entry := range entries {
		arr.Append(types.NewString(entry.Name()))
	}

	// TODO: Handle sorting_order parameter (0 = ascending, 1 = descending)

	return types.NewArray(arr)
}

// Glob finds pathnames matching a pattern
// glob(string $pattern, int $flags = 0): array|false
func Glob(pattern *types.Value, args ...*types.Value) *types.Value {
	patternStr := pattern.ToString()

	matches, err := filepath.Glob(patternStr)
	if err != nil {
		return types.NewBool(false)
	}

	arr := types.NewEmptyArray()
	for _, match := range matches {
		arr.Append(types.NewString(match))
	}

	return types.NewArray(arr)
}

// ============================================================================
// Path Functions
// ============================================================================

// Dirname returns a parent directory's path
// dirname(string $path, int $levels = 1): string
func Dirname(path *types.Value, args ...*types.Value) *types.Value {
	pathStr := path.ToString()

	levels := 1
	if len(args) > 0 {
		levels = int(args[0].ToInt())
	}

	result := pathStr
	for i := 0; i < levels; i++ {
		result = filepath.Dir(result)
	}

	return types.NewString(result)
}

// Basename returns trailing name component of path
// basename(string $path, string $suffix = ""): string
func Basename(path *types.Value, args ...*types.Value) *types.Value {
	pathStr := path.ToString()
	base := filepath.Base(pathStr)

	if len(args) > 0 {
		suffix := args[0].ToString()
		if strings.HasSuffix(base, suffix) {
			base = base[:len(base)-len(suffix)]
		}
	}

	return types.NewString(base)
}

// Pathinfo returns information about a file path
// pathinfo(string $path, int $flags = PATHINFO_ALL): array|string
func Pathinfo(path *types.Value, args ...*types.Value) *types.Value {
	pathStr := path.ToString()

	dir := filepath.Dir(pathStr)
	base := filepath.Base(pathStr)
	ext := filepath.Ext(pathStr)
	filename := base
	if ext != "" {
		filename = base[:len(base)-len(ext)]
		ext = ext[1:] // Remove leading dot
	}

	// PATHINFO_ALL (default)
	if len(args) == 0 {
		arr := types.NewEmptyArray()
		arr.Set(types.NewString("dirname"), types.NewString(dir))
		arr.Set(types.NewString("basename"), types.NewString(base))
		arr.Set(types.NewString("extension"), types.NewString(ext))
		arr.Set(types.NewString("filename"), types.NewString(filename))
		return types.NewArray(arr)
	}

	// Specific component
	flags := int(args[0].ToInt())
	switch flags {
	case 1: // PATHINFO_DIRNAME
		return types.NewString(dir)
	case 2: // PATHINFO_BASENAME
		return types.NewString(base)
	case 4: // PATHINFO_EXTENSION
		return types.NewString(ext)
	case 8: // PATHINFO_FILENAME
		return types.NewString(filename)
	default:
		return types.NewBool(false)
	}
}

// Realpath returns canonicalized absolute pathname
// realpath(string $path): string|false
func Realpath(path *types.Value) *types.Value {
	pathStr := path.ToString()

	// Validate path for security
	validPath, err := validatePath(pathStr)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	abs, err := filepath.Abs(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	// Resolve symlinks
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(real)
}

// Unlink deletes a file
// unlink(string $filename): bool
func Unlink(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	err = os.Remove(validPath)
	return types.NewBool(err == nil)
}

// Rename renames a file or directory
// rename(string $from, string $to): bool
func Rename(from *types.Value, to *types.Value) *types.Value {
	fromPath := from.ToString()
	toPath := to.ToString()

	// Validate both paths for security
	validFrom, err := validatePath(fromPath)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	validTo, err := validatePath(toPath)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	err = os.Rename(validFrom, validTo)
	return types.NewBool(err == nil)
}

// Copy copies a file
// copy(string $from, string $to): bool
func Copy(from *types.Value, to *types.Value) *types.Value {
	fromPath := from.ToString()
	toPath := to.ToString()

	// Validate both paths for security
	validFrom, err := validatePath(fromPath)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	validTo, err := validatePath(toPath)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	source, err := os.Open(validFrom)
	if err != nil {
		return types.NewBool(false)
	}
	defer source.Close()

	dest, err := os.Create(validTo)
	if err != nil {
		return types.NewBool(false)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return types.NewBool(err == nil)
}

// Filemtime gets file modification time
// filemtime(string $filename): int|false
func Filemtime(filename *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	info, err := os.Stat(validPath)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewInt(info.ModTime().Unix())
}

// Touch sets access and modification time of file
// touch(string $filename, int $mtime = null, int $atime = null): bool
func Touch(filename *types.Value, args ...*types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	// Check if file exists, create if not
	_, err = os.Stat(validPath)
	if os.IsNotExist(err) {
		file, err := os.Create(validPath)
		if err != nil {
			return types.NewBool(false)
		}
		file.Close()
	}

	// Get current time for defaults
	now := time.Now()
	mtime := now
	atime := now

	// Optional mtime parameter
	if len(args) > 0 && args[0] != nil && args[0].Type() != types.TypeNull {
		mtime = time.Unix(args[0].ToInt(), 0)
	}

	// Optional atime parameter
	if len(args) > 1 && args[1] != nil && args[1].Type() != types.TypeNull {
		atime = time.Unix(args[1].ToInt(), 0)
	}

	err = os.Chtimes(validPath, atime, mtime)
	return types.NewBool(err == nil)
}

// Chmod changes file mode
// chmod(string $filename, int $permissions): bool
func Chmod(filename *types.Value, permissions *types.Value) *types.Value {
	path := filename.ToString()

	// Validate path for security
	validPath, err := validatePath(path)
	if err != nil {
		// Security violation - return false
		return types.NewBool(false)
	}

	mode := os.FileMode(permissions.ToInt())
	err = os.Chmod(validPath, mode)
	return types.NewBool(err == nil)
}
