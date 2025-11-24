package vm

import (
	"github.com/krizos/php-go/pkg/stdlib/file"
	"github.com/krizos/php-go/pkg/stdlib/hash"
	stringfuncs "github.com/krizos/php-go/pkg/stdlib/string"
	varfuncs "github.com/krizos/php-go/pkg/stdlib/var"
	"github.com/krizos/php-go/pkg/types"
)

// BuiltinFunction represents a native Go function that implements a PHP built-in
type BuiltinFunction func(args []*types.Value) (*types.Value, error)

// builtinFunctions maps function names to their implementations
var builtinFunctions = map[string]BuiltinFunction{
	// Variable debugging functions
	"var_dump": builtinVarDump,

	// String functions
	"strlen":     builtinStrlen,
	"str_repeat": builtinStrRepeat,

	// File functions
	"file_exists":       builtinFileExists,
	"file_get_contents": builtinFileGetContents,

	// Hash functions
	"hash":            builtinHash,
	"hash_file":       builtinHashFile,
	"hash_hmac":       builtinHashHmac,
	"hash_hmac_file":  builtinHashHmacFile,
	"hash_equals":     builtinHashEquals,
	"hash_algos":      builtinHashAlgos,
	"hash_hmac_algos": builtinHashHmacAlgos,
	"hash_init":       builtinHashInit,
	"hash_update":     builtinHashUpdate,
	"hash_final":      builtinHashFinal,
	"hash_copy":       builtinHashCopy,
	"md5":             builtinMd5,
	"md5_file":        builtinMd5File,
	"sha1":            builtinSha1,
	"sha1_file":       builtinSha1File,
	"crc32":           builtinCrc32,
}

// IsBuiltin checks if a function name is a built-in function
func IsBuiltin(name string) bool {
	_, exists := builtinFunctions[name]
	return exists
}

// CallBuiltin calls a built-in function with the given arguments
func CallBuiltin(name string, args []*types.Value) (*types.Value, error) {
	fn, exists := builtinFunctions[name]
	if !exists {
		return nil, nil // Not a built-in
	}
	return fn(args)
}

// ============================================================================
// Variable Debugging Functions
// ============================================================================

// builtinVarDump implements var_dump(mixed ...$vars): void
func builtinVarDump(args []*types.Value) (*types.Value, error) {
	return varfuncs.VarDump(args...), nil
}

// ============================================================================
// String Functions
// ============================================================================

// builtinStrlen implements strlen(string $string): int
func builtinStrlen(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		// PHP strlen() requires exactly 1 argument
		return types.NewInt(0), nil
	}

	str := args[0].ToString()
	return types.NewInt(int64(len(str))), nil
}

// builtinStrRepeat implements str_repeat(string $string, int $times): string
func builtinStrRepeat(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewString(""), nil
	}
	return stringfuncs.StrRepeat(args[0], args[1]), nil
}

// ============================================================================
// File Functions
// ============================================================================

// builtinFileExists implements file_exists(string $filename): bool
func builtinFileExists(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.FileExists(args[0]), nil
}

// builtinFileGetContents implements file_get_contents(string $filename): string|false
func builtinFileGetContents(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.FileGetContents(args[0]), nil
}

// ============================================================================
// Hash Functions
// ============================================================================

// builtinHash implements hash(string $algo, string $data, bool $binary = false): string
func builtinHash(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return hash.Hash(args[0], args[1], args[2:]...), nil
}

// builtinHashFile implements hash_file(string $algo, string $filename, bool $binary = false): string|false
func builtinHashFile(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return hash.HashFile(args[0], args[1], args[2:]...), nil
}

// builtinHashHmac implements hash_hmac(string $algo, string $data, string $key, bool $binary = false): string
func builtinHashHmac(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewBool(false), nil
	}
	return hash.HashHmac(args[0], args[1], args[2], args[3:]...), nil
}

// builtinHashHmacFile implements hash_hmac_file(string $algo, string $filename, string $key, bool $binary = false): string|false
func builtinHashHmacFile(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewBool(false), nil
	}
	return hash.HashHmacFile(args[0], args[1], args[2], args[3:]...), nil
}

// builtinHashEquals implements hash_equals(string $known_string, string $user_string): bool
func builtinHashEquals(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return hash.HashEquals(args[0], args[1]), nil
}

// builtinHashAlgos implements hash_algos(): array
func builtinHashAlgos(args []*types.Value) (*types.Value, error) {
	return hash.HashAlgos(), nil
}

// builtinHashHmacAlgos implements hash_hmac_algos(): array
func builtinHashHmacAlgos(args []*types.Value) (*types.Value, error) {
	return hash.HashHmacAlgos(), nil
}

// builtinMd5 implements md5(string $str, bool $binary = false): string
func builtinMd5(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return hash.Md5(args[0], args[1:]...), nil
}

// builtinMd5File implements md5_file(string $filename, bool $binary = false): string|false
func builtinMd5File(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return hash.Md5File(args[0], args[1:]...), nil
}

// builtinSha1 implements sha1(string $str, bool $binary = false): string
func builtinSha1(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return hash.Sha1(args[0], args[1:]...), nil
}

// builtinSha1File implements sha1_file(string $filename, bool $binary = false): string|false
func builtinSha1File(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return hash.Sha1File(args[0], args[1:]...), nil
}

// builtinCrc32 implements crc32(string $str): int
func builtinCrc32(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return hash.Crc32(args[0]), nil
}

// builtinHashInit implements hash_init(string $algo, int $options = 0, string $key = ""): HashContext
func builtinHashInit(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return hash.HashInit(args[0], args[1:]...), nil
}

// builtinHashUpdate implements hash_update(HashContext $context, string $data): bool
func builtinHashUpdate(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return hash.HashUpdate(args[0], args[1]), nil
}

// builtinHashFinal implements hash_final(HashContext $context, bool $binary = false): string
func builtinHashFinal(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return hash.HashFinal(args[0], args[1:]...), nil
}

// builtinHashCopy implements hash_copy(HashContext $context): HashContext
func builtinHashCopy(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return hash.HashCopy(args[0]), nil
}
