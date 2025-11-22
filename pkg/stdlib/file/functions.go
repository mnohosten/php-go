package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// File Reading Functions
// ============================================================================

// FileGetContents reads entire file into a string
// file_get_contents(string $filename): string|false
func FileGetContents(filename *types.Value) *types.Value {
	path := filename.ToString()

	data, err := os.ReadFile(path)
	if err != nil {
		return types.NewBool(false)
	}

	return types.NewString(string(data))
}

// FilePutContents writes a string to a file
// file_put_contents(string $filename, mixed $data, int $flags = 0): int|false
func FilePutContents(filename *types.Value, data *types.Value, args ...*types.Value) *types.Value {
	path := filename.ToString()
	content := data.ToString()

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

	file, err := os.OpenFile(path, writeFlags, 0644)
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

	data, err := os.ReadFile(path)
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

	data, err := os.ReadFile(path)
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

	file, err := os.OpenFile(path, flags, 0644)
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

	n := int(length.ToInt())
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
		length := int(args[0].ToInt())
		if length < len(content) {
			content = content[:length]
		}
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
	_, err := os.Stat(path)
	return types.NewBool(err == nil)
}

// IsFile tells whether the filename is a regular file
// is_file(string $filename): bool
func IsFile(filename *types.Value) *types.Value {
	path := filename.ToString()
	info, err := os.Stat(path)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewBool(info.Mode().IsRegular())
}

// IsDir tells whether the filename is a directory
// is_dir(string $filename): bool
func IsDir(filename *types.Value) *types.Value {
	path := filename.ToString()
	info, err := os.Stat(path)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewBool(info.IsDir())
}

// IsReadable tells whether a file exists and is readable
// is_readable(string $filename): bool
func IsReadable(filename *types.Value) *types.Value {
	path := filename.ToString()
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
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

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		// File doesn't exist, check if directory is writable
		dir := filepath.Dir(path)
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
	info, err := os.Stat(path)
	if err != nil {
		return types.NewBool(false)
	}
	return types.NewInt(info.Size())
}

// Filetype gets file type
// filetype(string $filename): string|false
func Filetype(filename *types.Value) *types.Value {
	path := filename.ToString()
	info, err := os.Stat(path)
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

	permissions := os.FileMode(0777)
	if len(args) > 0 {
		permissions = os.FileMode(args[0].ToInt())
	}

	recursive := false
	if len(args) > 1 {
		recursive = args[1].ToBool()
	}

	var err error
	if recursive {
		err = os.MkdirAll(path, permissions)
	} else {
		err = os.Mkdir(path, permissions)
	}

	return types.NewBool(err == nil)
}

// Rmdir removes directory
// rmdir(string $directory): bool
func Rmdir(directory *types.Value) *types.Value {
	path := directory.ToString()
	err := os.Remove(path)
	return types.NewBool(err == nil)
}

// Scandir lists files and directories inside the specified path
// scandir(string $directory, int $sorting_order = SCANDIR_SORT_ASCENDING): array|false
func Scandir(directory *types.Value, args ...*types.Value) *types.Value {
	path := directory.ToString()

	entries, err := os.ReadDir(path)
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

	abs, err := filepath.Abs(pathStr)
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
	err := os.Remove(path)
	return types.NewBool(err == nil)
}

// Rename renames a file or directory
// rename(string $from, string $to): bool
func Rename(from *types.Value, to *types.Value) *types.Value {
	fromPath := from.ToString()
	toPath := to.ToString()
	err := os.Rename(fromPath, toPath)
	return types.NewBool(err == nil)
}

// Copy copies a file
// copy(string $from, string $to): bool
func Copy(from *types.Value, to *types.Value) *types.Value {
	fromPath := from.ToString()
	toPath := to.ToString()

	source, err := os.Open(fromPath)
	if err != nil {
		return types.NewBool(false)
	}
	defer source.Close()

	dest, err := os.Create(toPath)
	if err != nil {
		return types.NewBool(false)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return types.NewBool(err == nil)
}
