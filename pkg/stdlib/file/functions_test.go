package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// File Reading Tests
// ============================================================================

func TestFileGetContents(t *testing.T) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Hello, World!"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test reading
	result := FileGetContents(types.NewString(tmpfile.Name()))
	if result.ToString() != content {
		t.Errorf("FileGetContents() = %v, want %v", result.ToString(), content)
	}
}

func TestFileGetContentsNonexistent(t *testing.T) {
	result := FileGetContents(types.NewString("/nonexistent/file"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("FileGetContents(nonexistent) should return false")
	}
}

func TestFilePutContents(t *testing.T) {
	tmpfile := filepath.Join(os.TempDir(), "test_put_contents")
	defer os.Remove(tmpfile)

	content := "Test content"
	result := FilePutContents(types.NewString(tmpfile), types.NewString(content))

	if result.Type() != types.TypeInt {
		t.Errorf("FilePutContents should return int")
	}

	// Verify content
	data, err := os.ReadFile(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("File content = %v, want %v", string(data), content)
	}
}

func TestFilePutContentsAppend(t *testing.T) {
	tmpfile := filepath.Join(os.TempDir(), "test_append")
	defer os.Remove(tmpfile)

	// Write initial content
	FilePutContents(types.NewString(tmpfile), types.NewString("First\n"))

	// Append
	FilePutContents(types.NewString(tmpfile), types.NewString("Second\n"), types.NewInt(8)) // FILE_APPEND = 8

	data, err := os.ReadFile(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "First\nSecond\n"
	if string(data) != expected {
		t.Errorf("File content = %v, want %v", string(data), expected)
	}
}

func TestFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Line 1\nLine 2\nLine 3"
	tmpfile.WriteString(content)
	tmpfile.Close()

	result := File(types.NewString(tmpfile.Name()))
	if result.Type() != types.TypeArray {
		t.Errorf("File() should return array")
	}

	arr := result.ToArray()
	if arr.Len() != 3 {
		t.Errorf("File() should return 3 lines, got %d", arr.Len())
	}
}

func TestFileIgnoreNewlines(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("Line 1\nLine 2")
	tmpfile.Close()

	// FILE_IGNORE_NEW_LINES = 2
	result := File(types.NewString(tmpfile.Name()), types.NewInt(2))
	arr := result.ToArray()

	line0, _ := arr.Get(types.NewInt(0))
	if line0.ToString() == "Line 1\n" {
		t.Errorf("FILE_IGNORE_NEW_LINES should strip newlines")
	}
}

func TestReadfile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Test content for readfile"
	tmpfile.WriteString(content)
	tmpfile.Close()

	result := Readfile(types.NewString(tmpfile.Name()))
	if result.Type() != types.TypeInt {
		t.Errorf("Readfile() should return int (byte count)")
	}

	if result.ToInt() != int64(len(content)) {
		t.Errorf("Readfile() = %v, want %v bytes", result.ToInt(), len(content))
	}
}

func TestReadfileNonexistent(t *testing.T) {
	result := Readfile(types.NewString("/nonexistent/file"))
	if result.ToBool() != false {
		t.Errorf("Readfile(nonexistent) should return false")
	}
}

// ============================================================================
// File Handle Tests
// ============================================================================

func TestFopenFclose(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Open file
	handle := Fopen(types.NewString(tmpfile.Name()), types.NewString("r"))
	if handle.Type() != types.TypeResource {
		t.Errorf("Fopen() should return resource")
	}

	// Close file
	result := Fclose(handle)
	if result.ToBool() != true {
		t.Errorf("Fclose() should return true")
	}
}

func TestFopenInvalidMode(t *testing.T) {
	result := Fopen(types.NewString("/tmp/test"), types.NewString("invalid"))
	if result.ToBool() != false {
		t.Errorf("Fopen(invalid mode) should return false")
	}
}

func TestFopenModes(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	modes := []string{"r", "r+", "w", "w+", "a", "a+"}
	for _, mode := range modes {
		handle := Fopen(types.NewString(tmpfile.Name()), types.NewString(mode))
		if handle.Type() != types.TypeResource {
			t.Errorf("Fopen(mode=%q) should return resource", mode)
		}
		Fclose(handle)
	}
}

func TestFcloseNonResource(t *testing.T) {
	result := Fclose(types.NewString("not a resource"))
	if result.ToBool() != false {
		t.Errorf("Fclose(non-resource) should return false")
	}
}

func TestFreadNonResource(t *testing.T) {
	result := Fread(types.NewString("not a resource"), types.NewInt(10))
	if result.ToBool() != false {
		t.Errorf("Fread(non-resource) should return false")
	}
}

func TestFwriteWithLength(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	handle := Fopen(types.NewString(tmpfile.Name()), types.NewString("w"))

	// Write only first 5 characters
	result := Fwrite(handle, types.NewString("Hello, World!"), types.NewInt(5))
	if result.ToInt() != 5 {
		t.Errorf("Fwrite with length should write 5 bytes, wrote %d", result.ToInt())
	}

	Fclose(handle)
}

func TestFreadFwrite(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Open for writing
	handle := Fopen(types.NewString(tmpfile.Name()), types.NewString("w"))

	// Write
	content := "Hello, File!"
	written := Fwrite(handle, types.NewString(content))
	if written.ToInt() != int64(len(content)) {
		t.Errorf("Fwrite() = %v, want %v", written.ToInt(), len(content))
	}

	Fclose(handle)

	// Open for reading
	handle = Fopen(types.NewString(tmpfile.Name()), types.NewString("r"))

	// Read
	read := Fread(handle, types.NewInt(13))
	if read.ToString() != content {
		t.Errorf("Fread() = %v, want %v", read.ToString(), content)
	}

	Fclose(handle)
}

func TestFgets(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("Line 1\nLine 2\nLine 3")
	tmpfile.Close()

	handle := Fopen(types.NewString(tmpfile.Name()), types.NewString("r"))

	line1 := Fgets(handle)
	if line1.ToString() != "Line 1\n" {
		t.Errorf("Fgets() = %v, want 'Line 1\\n'", line1.ToString())
	}

	line2 := Fgets(handle)
	if line2.ToString() != "Line 2\n" {
		t.Errorf("Fgets() = %v, want 'Line 2\\n'", line2.ToString())
	}

	Fclose(handle)
}

func TestFgetc(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("ABC")
	tmpfile.Close()

	handle := Fopen(types.NewString(tmpfile.Name()), types.NewString("r"))

	char1 := Fgetc(handle)
	if char1.ToString() != "A" {
		t.Errorf("Fgetc() = %v, want 'A'", char1.ToString())
	}

	char2 := Fgetc(handle)
	if char2.ToString() != "B" {
		t.Errorf("Fgetc() = %v, want 'B'", char2.ToString())
	}

	Fclose(handle)
}

// ============================================================================
// File Information Tests
// ============================================================================

func TestFileExists(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Existing file
	result := FileExists(types.NewString(tmpfile.Name()))
	if result.ToBool() != true {
		t.Errorf("FileExists(existing) should return true")
	}

	// Non-existing file
	result = FileExists(types.NewString("/nonexistent"))
	if result.ToBool() != false {
		t.Errorf("FileExists(nonexistent) should return false")
	}
}

func TestIsFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Regular file
	result := IsFile(types.NewString(tmpfile.Name()))
	if result.ToBool() != true {
		t.Errorf("IsFile(file) should return true")
	}

	// Directory
	tmpdir := os.TempDir()
	result = IsFile(types.NewString(tmpdir))
	if result.ToBool() != false {
		t.Errorf("IsFile(directory) should return false")
	}
}

func TestIsDir(t *testing.T) {
	tmpdir := os.TempDir()

	result := IsDir(types.NewString(tmpdir))
	if result.ToBool() != true {
		t.Errorf("IsDir(directory) should return true")
	}

	// File
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	result = IsDir(types.NewString(tmpfile.Name()))
	if result.ToBool() != false {
		t.Errorf("IsDir(file) should return false")
	}
}

func TestIsReadable(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	result := IsReadable(types.NewString(tmpfile.Name()))
	if result.ToBool() != true {
		t.Errorf("IsReadable(readable file) should return true")
	}
}

func TestIsWritable(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	result := IsWritable(types.NewString(tmpfile.Name()))
	if result.ToBool() != true {
		t.Errorf("IsWritable(writable file) should return true")
	}
}

func TestIsWritableNonexistentInValidDir(t *testing.T) {
	// File doesn't exist, but dir is writable
	tmpdir := os.TempDir()
	nonexistentFile := filepath.Join(tmpdir, "nonexistent_file")

	result := IsWritable(types.NewString(nonexistentFile))
	if result.ToBool() != true {
		t.Errorf("IsWritable(nonexistent in writable dir) should return true")
	}
}

func TestFilesize(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Hello!"
	tmpfile.WriteString(content)
	tmpfile.Close()

	result := Filesize(types.NewString(tmpfile.Name()))
	if result.ToInt() != int64(len(content)) {
		t.Errorf("Filesize() = %v, want %v", result.ToInt(), len(content))
	}
}

func TestFiletype(t *testing.T) {
	// File
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	result := Filetype(types.NewString(tmpfile.Name()))
	if result.ToString() != "file" {
		t.Errorf("Filetype(file) = %v, want 'file'", result.ToString())
	}

	// Directory
	tmpdir := os.TempDir()
	result = Filetype(types.NewString(tmpdir))
	if result.ToString() != "dir" {
		t.Errorf("Filetype(dir) = %v, want 'dir'", result.ToString())
	}
}

// ============================================================================
// Directory Tests
// ============================================================================

func TestMkdirRmdir(t *testing.T) {
	tmpdir := filepath.Join(os.TempDir(), "test_mkdir")
	defer os.RemoveAll(tmpdir)

	// Create directory
	result := Mkdir(types.NewString(tmpdir))
	if result.ToBool() != true {
		t.Errorf("Mkdir() should return true")
	}

	// Verify it exists
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		t.Errorf("Directory was not created")
	}

	// Remove directory
	result = Rmdir(types.NewString(tmpdir))
	if result.ToBool() != true {
		t.Errorf("Rmdir() should return true")
	}
}

func TestMkdirRecursive(t *testing.T) {
	tmpdir := filepath.Join(os.TempDir(), "test_mkdir_recursive", "subdir", "subsubdir")
	defer os.RemoveAll(filepath.Join(os.TempDir(), "test_mkdir_recursive"))

	// Create recursively
	result := Mkdir(types.NewString(tmpdir), types.NewInt(0755), types.NewBool(true))
	if result.ToBool() != true {
		t.Errorf("Mkdir(recursive) should return true")
	}

	// Verify it exists
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		t.Errorf("Recursive directory was not created")
	}
}

func TestScandir(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test_scandir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create some files
	os.WriteFile(filepath.Join(tmpdir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpdir, "file2.txt"), []byte("test"), 0644)

	result := Scandir(types.NewString(tmpdir))
	if result.Type() != types.TypeArray {
		t.Errorf("Scandir() should return array")
	}

	arr := result.ToArray()
	if arr.Len() < 4 { // ., .., file1.txt, file2.txt
		t.Errorf("Scandir() should return at least 4 entries, got %d", arr.Len())
	}
}

func TestGlob(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test_glob")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create some files
	os.WriteFile(filepath.Join(tmpdir, "test1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpdir, "test2.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpdir, "other.log"), []byte(""), 0644)

	pattern := filepath.Join(tmpdir, "*.txt")
	result := Glob(types.NewString(pattern))

	if result.Type() != types.TypeArray {
		t.Errorf("Glob() should return array")
	}

	arr := result.ToArray()
	if arr.Len() != 2 {
		t.Errorf("Glob(*.txt) should return 2 files, got %d", arr.Len())
	}
}

// ============================================================================
// Path Tests
// ============================================================================

func TestDirname(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.txt", "/path/to"},
		{"/path/to/", "/path/to"}, // Go's filepath.Dir doesn't strip trailing slash
		{"file.txt", "."},
	}

	for _, tt := range tests {
		result := Dirname(types.NewString(tt.path))
		if result.ToString() != tt.expected {
			t.Errorf("Dirname(%q) = %v, want %v", tt.path, result.ToString(), tt.expected)
		}
	}
}

func TestDirnameLevels(t *testing.T) {
	path := "/path/to/file.txt"
	result := Dirname(types.NewString(path), types.NewInt(2))

	expected := "/path"
	if result.ToString() != expected {
		t.Errorf("Dirname(%q, 2) = %v, want %v", path, result.ToString(), expected)
	}
}

func TestBasename(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.txt", "file.txt"},
		{"/path/to/", "to"},
		{"file.txt", "file.txt"},
	}

	for _, tt := range tests {
		result := Basename(types.NewString(tt.path))
		if result.ToString() != tt.expected {
			t.Errorf("Basename(%q) = %v, want %v", tt.path, result.ToString(), tt.expected)
		}
	}
}

func TestBasenameSuffix(t *testing.T) {
	result := Basename(types.NewString("/path/to/file.txt"), types.NewString(".txt"))

	expected := "file"
	if result.ToString() != expected {
		t.Errorf("Basename with suffix = %v, want %v", result.ToString(), expected)
	}
}

func TestPathinfo(t *testing.T) {
	result := Pathinfo(types.NewString("/path/to/file.txt"))

	if result.Type() != types.TypeArray {
		t.Errorf("Pathinfo() should return array")
	}

	arr := result.ToArray()

	dirname, _ := arr.Get(types.NewString("dirname"))
	if dirname.ToString() != "/path/to" {
		t.Errorf("dirname = %v, want '/path/to'", dirname.ToString())
	}

	basename, _ := arr.Get(types.NewString("basename"))
	if basename.ToString() != "file.txt" {
		t.Errorf("basename = %v, want 'file.txt'", basename.ToString())
	}

	extension, _ := arr.Get(types.NewString("extension"))
	if extension.ToString() != "txt" {
		t.Errorf("extension = %v, want 'txt'", extension.ToString())
	}

	filename, _ := arr.Get(types.NewString("filename"))
	if filename.ToString() != "file" {
		t.Errorf("filename = %v, want 'file'", filename.ToString())
	}
}

func TestPathinfoComponent(t *testing.T) {
	// PATHINFO_EXTENSION = 4
	result := Pathinfo(types.NewString("/path/to/file.txt"), types.NewInt(4))

	if result.ToString() != "txt" {
		t.Errorf("Pathinfo(PATHINFO_EXTENSION) = %v, want 'txt'", result.ToString())
	}
}

func TestRealpath(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	result := Realpath(types.NewString(tmpfile.Name()))
	if result.Type() != types.TypeString {
		t.Errorf("Realpath() should return string")
	}

	// Should be absolute path
	path := result.ToString()
	if !filepath.IsAbs(path) {
		t.Errorf("Realpath() should return absolute path, got %v", path)
	}
}

// ============================================================================
// File Operation Tests
// ============================================================================

func TestUnlink(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	result := Unlink(types.NewString(tmpfile.Name()))
	if result.ToBool() != true {
		t.Errorf("Unlink() should return true")
	}

	// Verify it's deleted
	if _, err := os.Stat(tmpfile.Name()); !os.IsNotExist(err) {
		t.Errorf("File should be deleted")
	}
}

func TestRename(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	oldName := tmpfile.Name()
	newName := tmpfile.Name() + ".renamed"

	defer os.Remove(oldName)
	defer os.Remove(newName)

	result := Rename(types.NewString(oldName), types.NewString(newName))
	if result.ToBool() != true {
		t.Errorf("Rename() should return true")
	}

	// Verify new file exists
	if _, err := os.Stat(newName); os.IsNotExist(err) {
		t.Errorf("Renamed file should exist")
	}

	// Verify old file doesn't exist
	if _, err := os.Stat(oldName); !os.IsNotExist(err) {
		t.Errorf("Old file should not exist")
	}
}

func TestCopy(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Test content"
	tmpfile.WriteString(content)
	tmpfile.Close()

	destPath := tmpfile.Name() + ".copy"
	defer os.Remove(destPath)

	result := Copy(types.NewString(tmpfile.Name()), types.NewString(destPath))
	if result.ToBool() != true {
		t.Errorf("Copy() should return true")
	}

	// Verify copied file content
	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("Copied file content = %v, want %v", string(data), content)
	}
}

// ============================================================================
// File Permission Tests
// ============================================================================

func TestFilePermissionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant os.FileMode
		expected os.FileMode
	}{
		{"FilePermissionPrivate", FilePermissionPrivate, 0600},
		{"FilePermissionPublic", FilePermissionPublic, 0644},
		{"DirPermissionPrivate", DirPermissionPrivate, 0700},
		{"DirPermissionPublic", DirPermissionPublic, 0755},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %o, got %o", tt.expected, tt.constant)
			}
		})
	}
}

// ============================================================================
// Integer Overflow Protection Tests (Phase 2)
// ============================================================================

func TestFreadIntegerOverflowProtection(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "test content"
	tmpfile.WriteString(content)
	tmpfile.Close()

	fp := Fopen(types.NewString(tmpfile.Name()), types.NewString("r"))
	defer Fclose(fp)

	// Normal read should work
	result := Fread(fp, types.NewInt(int64(len(content))))
	if result.Type() != types.TypeString {
		t.Error("Expected normal read to succeed")
	}
}

func TestFwriteIntegerOverflowProtection(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	fp := Fopen(types.NewString(tmpfile.Name()), types.NewString("w"))
	defer Fclose(fp)

	// Normal write should work
	content := "test content"
	result := Fwrite(fp, types.NewString(content), types.NewInt(5))
	if result.ToInt() != 5 {
		t.Error("Expected normal write to succeed")
	}
}

// ============================================================================
// Filemtime, Touch, Chmod Tests
// ============================================================================

func TestFilemtime(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test_filemtime")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Get modification time
	result := Filemtime(types.NewString(tmpfile.Name()))
	if result.Type() != types.TypeInt {
		t.Errorf("Expected int type, got %v", result.Type())
	}

	// Modification time should be recent (within last minute)
	mtime := result.ToInt()
	now := time.Now().Unix()
	if now-mtime > 60 {
		t.Errorf("Expected recent modification time, got %d (now: %d)", mtime, now)
	}
}

func TestFilemtimeNonExistent(t *testing.T) {
	result := Filemtime(types.NewString("/nonexistent/file"))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for nonexistent file")
	}
}

func TestTouch(t *testing.T) {
	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "test_touch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Touch creates a new file if it doesn't exist
	newFile := tmpdir + "/newfile.txt"
	result := Touch(types.NewString(newFile))
	if !result.ToBool() {
		t.Error("Expected touch to succeed for new file")
	}

	// File should exist now
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("Expected file to be created by touch")
	}
}

func TestTouchExistingFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test_touch_existing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Get original mtime
	info1, _ := os.Stat(tmpfile.Name())
	origMtime := info1.ModTime()

	// Wait a bit and touch
	time.Sleep(10 * time.Millisecond)

	result := Touch(types.NewString(tmpfile.Name()))
	if !result.ToBool() {
		t.Error("Expected touch to succeed for existing file")
	}

	// mtime should be updated
	info2, _ := os.Stat(tmpfile.Name())
	newMtime := info2.ModTime()

	if !newMtime.After(origMtime) && !newMtime.Equal(origMtime) {
		t.Error("Expected modification time to be updated")
	}
}

func TestTouchWithMtime(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test_touch_mtime")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Set a specific mtime (1 day ago)
	specificTime := time.Now().Add(-24 * time.Hour).Unix()
	result := Touch(types.NewString(tmpfile.Name()), types.NewInt(specificTime))
	if !result.ToBool() {
		t.Error("Expected touch with mtime to succeed")
	}

	// Verify mtime was set
	info, _ := os.Stat(tmpfile.Name())
	if info.ModTime().Unix() != specificTime {
		t.Errorf("Expected mtime %d, got %d", specificTime, info.ModTime().Unix())
	}
}

func TestChmod(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test_chmod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Change permissions to 0600
	result := Chmod(types.NewString(tmpfile.Name()), types.NewInt(0600))
	if !result.ToBool() {
		t.Error("Expected chmod to succeed")
	}

	// Verify permissions
	info, _ := os.Stat(tmpfile.Name())
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("Expected permissions 0600, got %o", perm)
	}
}

func TestChmodNonExistent(t *testing.T) {
	result := Chmod(types.NewString("/nonexistent/file"), types.NewInt(0644))
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Error("Expected false for nonexistent file")
	}
}
