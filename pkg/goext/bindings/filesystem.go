package bindings

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// FilesystemExtension provides Go filesystem functions for PHP.
type FilesystemExtension struct {
	*goext.BaseExtension
}

// NewFilesystemExtension creates a new filesystem extension.
func NewFilesystemExtension() *FilesystemExtension {
	ext := &FilesystemExtension{
		BaseExtension: goext.NewBaseExtension("go_fs", "1.0.0"),
	}

	// Register functions
	ext.AddFunction("go_file_read", ext.fileRead)
	ext.AddFunction("go_file_write", ext.fileWrite)
	ext.AddFunction("go_file_append", ext.fileAppend)
	ext.AddFunction("go_file_exists", ext.fileExists)
	ext.AddFunction("go_file_delete", ext.fileDelete)
	ext.AddFunction("go_file_copy", ext.fileCopy)
	ext.AddFunction("go_file_move", ext.fileMove)
	ext.AddFunction("go_file_stat", ext.fileStat)
	ext.AddFunction("go_dir_create", ext.dirCreate)
	ext.AddFunction("go_dir_list", ext.dirList)
	ext.AddFunction("go_dir_remove", ext.dirRemove)

	// Register constants
	ext.AddConstant("FILE_PERM_0644", types.NewInt(0644))
	ext.AddConstant("FILE_PERM_0755", types.NewInt(0755))
	ext.AddConstant("FILE_PERM_0777", types.NewInt(0777))

	return ext
}

// fileRead implements go_file_read($path)
//
// Returns: file contents as string or false on error
func (e *FilesystemExtension) fileRead(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_file_read() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	content, err := os.ReadFile(path)
	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(string(content)), nil
}

// fileWrite implements go_file_write($path, $data, $perm = 0644)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) fileWrite(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("go_file_write() expects at least 2 arguments (path, data), got %d", len(args))
	}

	path := args[0].ToString()
	data := args[1].ToString()

	perm := os.FileMode(0644)
	if len(args) >= 3 {
		perm = os.FileMode(args[2].ToInt())
	}

	err := os.WriteFile(path, []byte(data), perm)
	return types.NewBool(err == nil), nil
}

// fileAppend implements go_file_append($path, $data, $perm = 0644)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) fileAppend(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("go_file_append() expects at least 2 arguments (path, data), got %d", len(args))
	}

	path := args[0].ToString()
	data := args[1].ToString()

	perm := os.FileMode(0644)
	if len(args) >= 3 {
		perm = os.FileMode(args[2].ToInt())
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return types.NewBool(false), nil
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return types.NewBool(err == nil), nil
}

// fileExists implements go_file_exists($path)
//
// Returns: true if file exists, false otherwise
func (e *FilesystemExtension) fileExists(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_file_exists() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	_, err := os.Stat(path)
	return types.NewBool(err == nil), nil
}

// fileDelete implements go_file_delete($path)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) fileDelete(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_file_delete() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	err := os.Remove(path)
	return types.NewBool(err == nil), nil
}

// fileCopy implements go_file_copy($src, $dst)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) fileCopy(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_file_copy() expects 2 arguments (src, dst), got %d", len(args))
	}

	src := args[0].ToString()
	dst := args[1].ToString()

	sourceFile, err := os.Open(src)
	if err != nil {
		return types.NewBool(false), nil
	}
	defer sourceFile.Close()

	// Create destination file with explicit permissions (0644)
	// Owner: read/write, Group/Others: read
	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return types.NewBool(false), nil
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return types.NewBool(err == nil), nil
}

// fileMove implements go_file_move($src, $dst)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) fileMove(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_file_move() expects 2 arguments (src, dst), got %d", len(args))
	}

	src := args[0].ToString()
	dst := args[1].ToString()

	err := os.Rename(src, dst)
	return types.NewBool(err == nil), nil
}

// fileStat implements go_file_stat($path)
//
// Returns: array with file info or false on error
//   - size: file size in bytes
//   - mode: file mode
//   - is_dir: true if directory
//   - mod_time: modification time (Unix timestamp)
func (e *FilesystemExtension) fileStat(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_file_stat() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	info, err := os.Stat(path)
	if err != nil {
		return types.NewBool(false), nil
	}

	result := types.NewEmptyArray()
	result.Set(types.NewString("size"), types.NewInt(info.Size()))
	result.Set(types.NewString("mode"), types.NewInt(int64(info.Mode())))
	result.Set(types.NewString("is_dir"), types.NewBool(info.IsDir()))
	result.Set(types.NewString("mod_time"), types.NewInt(info.ModTime().Unix()))

	return types.NewArray(result), nil
}

// dirCreate implements go_dir_create($path, $perm = 0755)
//
// Returns: true on success, false on error
func (e *FilesystemExtension) dirCreate(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_dir_create() expects at least 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	perm := os.FileMode(0755)
	if len(args) >= 2 {
		perm = os.FileMode(args[1].ToInt())
	}

	err := os.MkdirAll(path, perm)
	return types.NewBool(err == nil), nil
}

// dirList implements go_dir_list($path)
//
// Returns: array of filenames or false on error
func (e *FilesystemExtension) dirList(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_dir_list() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	entries, err := os.ReadDir(path)
	if err != nil {
		return types.NewBool(false), nil
	}

	result := types.NewEmptyArray()
	for _, entry := range entries {
		result.Append(types.NewString(entry.Name()))
	}

	return types.NewArray(result), nil
}

// dirRemove implements go_dir_remove($path)
//
// Returns: true on success, false on error
// Note: Removes directory and all contents (equivalent to rm -rf)
func (e *FilesystemExtension) dirRemove(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_dir_remove() expects 1 argument (path), got %d", len(args))
	}

	path := args[0].ToString()

	err := os.RemoveAll(path)
	return types.NewBool(err == nil), nil
}

// Helper: walk directory recursively
func (e *FilesystemExtension) walkDir(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})

	return files, err
}
