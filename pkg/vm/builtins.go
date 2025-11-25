package vm

import (
	arrayfuncs "github.com/krizos/php-go/pkg/stdlib/array"
	"github.com/krizos/php-go/pkg/stdlib/file"
	"github.com/krizos/php-go/pkg/stdlib/hash"
	jsonfuncs "github.com/krizos/php-go/pkg/stdlib/json"
	mathfuncs "github.com/krizos/php-go/pkg/stdlib/math"
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
	"strlen":       builtinStrlen,
	"substr":       builtinSubstr,
	"strpos":       builtinStrpos,
	"strrpos":      builtinStrrpos,
	"str_replace":  builtinStrReplace,
	"str_repeat":   builtinStrRepeat,
	"strtolower":   builtinStrtolower,
	"strtoupper":   builtinStrtoupper,
	"ucfirst":      builtinUcfirst,
	"lcfirst":      builtinLcfirst,
	"ucwords":      builtinUcwords,
	"trim":         builtinTrim,
	"ltrim":        builtinLtrim,
	"rtrim":        builtinRtrim,
	"explode":      builtinExplode,
	"implode":      builtinImplode,
	"join":         builtinJoin,
	"str_split":      builtinStrSplit,
	"str_shuffle":    builtinStrShuffle,
	"str_word_count": builtinStrWordCount,
	"substr_count":   builtinSubstrCount,
	"substr_replace": builtinSubstrReplace,

	// String search functions
	"strstr":  builtinStrstr,
	"strchr":  builtinStrchr,
	"stristr": builtinStristr,
	"strrchr": builtinStrrchr,
	"strpbrk": builtinStrpbrk,
	"strspn":  builtinStrspn,
	"strcspn": builtinStrcspn,

	// String formatting functions
	"sprintf":       builtinSprintf,
	"printf":        builtinPrintf,
	"vsprintf":      builtinVsprintf,
	"vprintf":       builtinVprintf,
	"number_format": builtinNumberFormat,
	"sscanf":        builtinSscanf,
	"str_getcsv":    builtinStrGetcsv,

	// HTML encoding functions
	"htmlspecialchars":        builtinHtmlspecialchars,
	"htmlentities":            builtinHtmlentities,
	"htmlspecialchars_decode": builtinHtmlspecialcharsDecode,
	"html_entity_decode":      builtinHtmlEntityDecode,
	"strip_tags":              builtinStripTags,
	"addslashes":              builtinAddslashes,
	"stripslashes":            builtinStripslashes,

	// Array functions
	"count":         builtinCount,
	"sizeof":        builtinSizeof,
	"array_keys":    builtinArrayKeys,
	"array_values":  builtinArrayValues,
	"array_push":    builtinArrayPush,
	"array_pop":     builtinArrayPop,
	"array_shift":   builtinArrayShift,
	"array_unshift": builtinArrayUnshift,
	"array_merge":   builtinArrayMerge,
	"in_array":      builtinInArray,
	"array_search":  builtinArraySearch,
	"array_slice":   builtinArraySlice,
	"array_reverse": builtinArrayReverse,
	"array_unique":  builtinArrayUnique,
	"array_combine": builtinArrayCombine,
	"array_flip":    builtinArrayFlip,
	"array_fill":    builtinArrayFill,
	"array_chunk":   builtinArrayChunk,
	"sort":          builtinSort,
	"rsort":         builtinRsort,
	"asort":         builtinAsort,
	"arsort":        builtinArsort,
	"ksort":         builtinKsort,
	"krsort":        builtinKrsort,
	"array_map":     builtinArrayMap,
	"array_filter":  builtinArrayFilter,
	"array_reduce":  builtinArrayReduce,
	"array_walk":           builtinArrayWalk,
	"array_walk_recursive": builtinArrayWalkRecursive,
	"array_diff":           builtinArrayDiff,
	"array_intersect":         builtinArrayIntersect,
	"array_pad":               builtinArrayPad,
	"array_sum":               builtinArraySum,
	"array_product":           builtinArrayProduct,
	"array_column":            builtinArrayColumn,
	"array_change_key_case":   builtinArrayChangeKeyCase,
	"array_replace":           builtinArrayReplace,
	"array_replace_recursive": builtinArrayReplaceRecursive,
	"current":                 builtinCurrent,
	"key":                     builtinKey,
	"reset":                   builtinReset,
	"end":                     builtinEnd,
	"next":                    builtinNext,
	"prev":                    builtinPrev,
	"array_key_exists":        builtinArrayKeyExists,
	"key_exists":              builtinKeyExists,
	"array_diff_key":          builtinArrayDiffKey,
	"array_intersect_key":     builtinArrayIntersectKey,
	"array_diff_assoc":        builtinArrayDiffAssoc,
	"array_intersect_assoc":   builtinArrayIntersectAssoc,
	"array_key_first":         builtinArrayKeyFirst,
	"array_key_last":          builtinArrayKeyLast,
	"array_count_values":      builtinArrayCountValues,

	// File functions
	"file_exists":        builtinFileExists,
	"file_get_contents":  builtinFileGetContents,
	"file_put_contents":  builtinFilePutContents,
	"file":               builtinFile,
	"readfile":           builtinReadfile,
	"fopen":              builtinFopen,
	"fclose":             builtinFclose,
	"fread":              builtinFread,
	"fwrite":             builtinFwrite,
	"fgets":              builtinFgets,
	"fgetc":              builtinFgetc,
	"is_file":            builtinIsFile,
	"is_dir":             builtinIsDir,
	"is_readable":        builtinIsReadable,
	"is_writable":        builtinIsWritable,
	"is_writeable":       builtinIsWritable, // Alias
	"filesize":           builtinFilesize,
	"filemtime":          builtinFilemtime,
	"filetype":           builtinFiletype,
	"mkdir":              builtinMkdir,
	"rmdir":              builtinRmdir,
	"unlink":             builtinUnlink,
	"rename":             builtinRename,
	"copy":               builtinCopy,
	"touch":              builtinTouch,
	"chmod":              builtinChmod,
	"scandir":            builtinScandir,
	"glob":               builtinGlob,
	"dirname":            builtinDirname,
	"basename":           builtinBasename,
	"pathinfo":           builtinPathinfo,
	"realpath":           builtinRealpath,

	// URL functions
	"urlencode":        builtinUrlencode,
	"urldecode":        builtinUrldecode,
	"rawurlencode":     builtinRawurlencode,
	"rawurldecode":     builtinRawurldecode,
	"base64_encode":    builtinBase64Encode,
	"base64_decode":    builtinBase64Decode,
	"parse_url":        builtinParseUrl,
	"http_build_query": builtinHttpBuildQuery,

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

	// JSON functions
	"json_encode":         builtinJsonEncode,
	"json_decode":         builtinJsonDecode,
	"json_last_error":     builtinJsonLastError,
	"json_last_error_msg": builtinJsonLastErrorMsg,

	// Type/Existence functions (function_exists registered in init to avoid cycle)
	"class_exists":  builtinClassExists,
	"method_exists": builtinMethodExists,

	// Serialization functions
	"serialize":   builtinSerialize,
	"unserialize": builtinUnserialize,

	// Math functions
	"abs":           builtinAbs,
	"ceil":          builtinCeil,
	"floor":         builtinFloor,
	"round":         builtinRound,
	"min":           builtinMin,
	"max":           builtinMax,
	"pow":           builtinPow,
	"sqrt":          builtinSqrt,
	"sin":           builtinSin,
	"cos":           builtinCos,
	"tan":           builtinTan,
	"asin":          builtinAsin,
	"acos":          builtinAcos,
	"atan":          builtinAtan,
	"atan2":         builtinAtan2,
	"deg2rad":       builtinDeg2rad,
	"rad2deg":       builtinRad2deg,
	"exp":           builtinExp,
	"log":           builtinLog,
	"log10":         builtinLog10,
	"log1p":         builtinLog1p,
	"expm1":         builtinExpm1,
	"rand":          builtinRand,
	"mt_rand":       builtinMtRand,
	"random_int":    builtinRandomInt,
	"getrandmax":    builtinGetRandMax,
	"mt_getrandmax": builtinMtGetRandMax,
	"pi":            builtinPi,
	"is_nan":        builtinIsNan,
	"is_infinite":   builtinIsInfinite,
	"is_finite":     builtinIsFinite,
	"hypot":         builtinHypot,
	"fmod":          builtinFmod,
	"intdiv":        builtinIntdiv,
	"fdiv":          builtinFdiv,
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

// builtinFilePutContents implements file_put_contents(string $filename, mixed $data, int $flags = 0): int|false
func builtinFilePutContents(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.FilePutContents(args[0], args[1], args[2:]...), nil
}

// builtinFile implements file(string $filename, int $flags = 0): array|false
func builtinFile(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.File(args[0], args[1:]...), nil
}

// builtinReadfile implements readfile(string $filename): int|false
func builtinReadfile(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Readfile(args[0]), nil
}

// builtinFopen implements fopen(string $filename, string $mode): resource|false
func builtinFopen(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Fopen(args[0], args[1]), nil
}

// builtinFclose implements fclose(resource $stream): bool
func builtinFclose(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Fclose(args[0]), nil
}

// builtinFread implements fread(resource $stream, int $length): string|false
func builtinFread(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Fread(args[0], args[1]), nil
}

// builtinFwrite implements fwrite(resource $stream, string $data, int $length = null): int|false
func builtinFwrite(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Fwrite(args[0], args[1], args[2:]...), nil
}

// builtinFgets implements fgets(resource $stream, int $length = null): string|false
func builtinFgets(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Fgets(args[0], args[1:]...), nil
}

// builtinFgetc implements fgetc(resource $stream): string|false
func builtinFgetc(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Fgetc(args[0]), nil
}

// builtinIsFile implements is_file(string $filename): bool
func builtinIsFile(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.IsFile(args[0]), nil
}

// builtinIsDir implements is_dir(string $filename): bool
func builtinIsDir(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.IsDir(args[0]), nil
}

// builtinIsReadable implements is_readable(string $filename): bool
func builtinIsReadable(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.IsReadable(args[0]), nil
}

// builtinIsWritable implements is_writable(string $filename): bool
func builtinIsWritable(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.IsWritable(args[0]), nil
}

// builtinFilesize implements filesize(string $filename): int|false
func builtinFilesize(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Filesize(args[0]), nil
}

// builtinFilemtime implements filemtime(string $filename): int|false
func builtinFilemtime(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Filemtime(args[0]), nil
}

// builtinFiletype implements filetype(string $filename): string|false
func builtinFiletype(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Filetype(args[0]), nil
}

// builtinMkdir implements mkdir(string $directory, int $permissions = 0777, bool $recursive = false): bool
func builtinMkdir(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Mkdir(args[0], args[1:]...), nil
}

// builtinRmdir implements rmdir(string $directory): bool
func builtinRmdir(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Rmdir(args[0]), nil
}

// builtinUnlink implements unlink(string $filename): bool
func builtinUnlink(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Unlink(args[0]), nil
}

// builtinRename implements rename(string $from, string $to): bool
func builtinRename(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Rename(args[0], args[1]), nil
}

// builtinCopy implements copy(string $from, string $to): bool
func builtinCopy(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Copy(args[0], args[1]), nil
}

// builtinTouch implements touch(string $filename, int $mtime = null, int $atime = null): bool
func builtinTouch(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Touch(args[0], args[1:]...), nil
}

// builtinChmod implements chmod(string $filename, int $permissions): bool
func builtinChmod(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return file.Chmod(args[0], args[1]), nil
}

// builtinScandir implements scandir(string $directory, int $sorting_order = 0): array|false
func builtinScandir(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Scandir(args[0], args[1:]...), nil
}

// builtinGlob implements glob(string $pattern, int $flags = 0): array|false
func builtinGlob(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Glob(args[0], args[1:]...), nil
}

// builtinDirname implements dirname(string $path, int $levels = 1): string
func builtinDirname(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return file.Dirname(args[0], args[1:]...), nil
}

// builtinBasename implements basename(string $path, string $suffix = ""): string
func builtinBasename(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return file.Basename(args[0], args[1:]...), nil
}

// builtinPathinfo implements pathinfo(string $path, int $flags = PATHINFO_ALL): array|string
func builtinPathinfo(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Pathinfo(args[0], args[1:]...), nil
}

// builtinRealpath implements realpath(string $path): string|false
func builtinRealpath(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return file.Realpath(args[0]), nil
}

// ============================================================================
// URL Functions
// ============================================================================

// builtinUrlencode implements urlencode(string $string): string
func builtinUrlencode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Urlencode(args[0]), nil
}

// builtinUrldecode implements urldecode(string $string): string
func builtinUrldecode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Urldecode(args[0]), nil
}

// builtinRawurlencode implements rawurlencode(string $string): string
func builtinRawurlencode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Rawurlencode(args[0]), nil
}

// builtinRawurldecode implements rawurldecode(string $string): string
func builtinRawurldecode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Rawurldecode(args[0]), nil
}

// builtinBase64Encode implements base64_encode(string $string): string
func builtinBase64Encode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Base64Encode(args[0]), nil
}

// builtinBase64Decode implements base64_decode(string $string, bool $strict = false): string|false
func builtinBase64Decode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Base64Decode(args[0], args[1:]...), nil
}

// builtinParseUrl implements parse_url(string $url, int $component = -1): array|string|int|null|false
func builtinParseUrl(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return stringfuncs.ParseUrl(args[0], args[1:]...), nil
}

// builtinHttpBuildQuery implements http_build_query(array $data, string $numeric_prefix = "", string $arg_separator = null): string
func builtinHttpBuildQuery(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.HttpBuildQuery(args[0], args[1:]...), nil
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

// ============================================================================
// HTML Encoding Functions
// ============================================================================

// builtinHtmlspecialchars implements htmlspecialchars(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null, bool $double_encode = true): string
func builtinHtmlspecialchars(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Htmlspecialchars(args[0], args[1:]...), nil
}

// builtinHtmlentities implements htmlentities(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null, bool $double_encode = true): string
func builtinHtmlentities(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Htmlentities(args[0], args[1:]...), nil
}

// builtinHtmlspecialcharsDecode implements htmlspecialchars_decode(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401): string
func builtinHtmlspecialcharsDecode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.HtmlspecialcharsDecode(args[0], args[1:]...), nil
}

// builtinHtmlEntityDecode implements html_entity_decode(string $string, int $flags = ENT_QUOTES|ENT_SUBSTITUTE|ENT_HTML401, ?string $encoding = null): string
func builtinHtmlEntityDecode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.HtmlEntityDecode(args[0], args[1:]...), nil
}

// builtinStripTags implements strip_tags(string $string, array|string|null $allowed_tags = null): string
func builtinStripTags(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.StripTags(args[0], args[1:]...), nil
}

// builtinAddslashes implements addslashes(string $string): string
func builtinAddslashes(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Addslashes(args[0]), nil
}

// builtinStripslashes implements stripslashes(string $string): string
func builtinStripslashes(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Stripslashes(args[0]), nil
}

// ============================================================================
// Additional String Functions
// ============================================================================

// builtinSubstr implements substr(string $string, int $offset, ?int $length = null): string
func builtinSubstr(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewString(""), nil
	}
	return stringfuncs.Substr(args[0], args[1], args[2:]...), nil
}

// builtinStrpos implements strpos(string $haystack, string $needle, int $offset = 0): int|false
func builtinStrpos(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Strpos(args[0], args[1], args[2:]...), nil
}

// builtinStrrpos implements strrpos(string $haystack, string $needle, int $offset = 0): int|false
func builtinStrrpos(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Strrpos(args[0], args[1], args[2:]...), nil
}

// builtinStrReplace implements str_replace(mixed $search, mixed $replace, mixed $subject): mixed
func builtinStrReplace(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewString(""), nil
	}
	return stringfuncs.StrReplace(args[0], args[1], args[2]), nil
}

// builtinStrtolower implements strtolower(string $string): string
func builtinStrtolower(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Strtolower(args[0]), nil
}

// builtinStrtoupper implements strtoupper(string $string): string
func builtinStrtoupper(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Strtoupper(args[0]), nil
}

// builtinUcfirst implements ucfirst(string $string): string
func builtinUcfirst(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Ucfirst(args[0]), nil
}

// builtinLcfirst implements lcfirst(string $string): string
func builtinLcfirst(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Lcfirst(args[0]), nil
}

// builtinUcwords implements ucwords(string $string, string $separators = " \t\r\n\f\v"): string
func builtinUcwords(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Ucwords(args[0]), nil
}

// builtinTrim implements trim(string $string, string $characters = " \t\n\r\0\x0B"): string
func builtinTrim(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Trim(args[0], args[1:]...), nil
}

// builtinLtrim implements ltrim(string $string, string $characters = " \t\n\r\0\x0B"): string
func builtinLtrim(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Ltrim(args[0], args[1:]...), nil
}

// builtinRtrim implements rtrim(string $string, string $characters = " \t\n\r\0\x0B"): string
func builtinRtrim(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Rtrim(args[0], args[1:]...), nil
}

// builtinExplode implements explode(string $delimiter, string $string, int $limit = PHP_INT_MAX): array
func builtinExplode(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Explode(args[0], args[1], args[2:]...), nil
}

// builtinImplode implements implode(string $separator, array $array): string
func builtinImplode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Implode(args[0], args[1:]...), nil
}

// builtinJoin implements join(string $separator, array $array): string (alias of implode)
func builtinJoin(args []*types.Value) (*types.Value, error) {
	return builtinImplode(args)
}

// builtinStrSplit implements str_split(string $string, int $length = 1): array|false
func builtinStrSplit(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return stringfuncs.StrSplit(args[0], args[1:]...), nil
}

// builtinStrShuffle implements str_shuffle(string $string): string
func builtinStrShuffle(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.StrShuffle(args[0]), nil
}

// builtinStrWordCount implements str_word_count(string $string, int $format = 0, ?string $characters = null): array|int
func builtinStrWordCount(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return stringfuncs.StrWordCount(args[0], args[1:]...), nil
}

// builtinSubstrCount implements substr_count(string $haystack, string $needle, int $offset = 0, ?int $length = null): int
func builtinSubstrCount(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(0), nil
	}
	return stringfuncs.SubstrCount(args[0], args[1], args[2:]...), nil
}

// builtinSubstrReplace implements substr_replace(array|string $string, array|string $replace, array|int $offset, array|int|null $length = null): string|array
func builtinSubstrReplace(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewString(""), nil
	}
	return stringfuncs.SubstrReplace(args[0], args[1], args[2], args[3:]...), nil
}

// ============================================================================
// String Search Functions
// ============================================================================

// builtinStrstr implements strstr(string $haystack, mixed $needle, bool $before_needle = false): string|false
func builtinStrstr(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Strstr(args[0], args[1], args[2:]...), nil
}

// builtinStrchr implements strchr (alias of strstr)
func builtinStrchr(args []*types.Value) (*types.Value, error) {
	return builtinStrstr(args)
}

// builtinStristr implements stristr(string $haystack, mixed $needle, bool $before_needle = false): string|false
func builtinStristr(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Stristr(args[0], args[1], args[2:]...), nil
}

// builtinStrrchr implements strrchr(string $haystack, mixed $needle): string|false
func builtinStrrchr(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Strrchr(args[0], args[1]), nil
}

// builtinStrpbrk implements strpbrk(string $string, string $characters): string|false
func builtinStrpbrk(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return stringfuncs.Strpbrk(args[0], args[1]), nil
}

// builtinStrspn implements strspn(string $string, string $characters, int $offset = 0, ?int $length = null): int
func builtinStrspn(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(0), nil
	}
	return stringfuncs.Strspn(args[0], args[1], args[2:]...), nil
}

// builtinStrcspn implements strcspn(string $string, string $characters, int $offset = 0, ?int $length = null): int
func builtinStrcspn(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(0), nil
	}
	return stringfuncs.Strcspn(args[0], args[1], args[2:]...), nil
}

// ============================================================================
// String Formatting Functions
// ============================================================================

// builtinSprintf implements sprintf(string $format, mixed ...$values): string
func builtinSprintf(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString(""), nil
	}
	return stringfuncs.Sprintf(args[0], args[1:]...), nil
}

// builtinPrintf implements printf(string $format, mixed ...$values): int
func builtinPrintf(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return stringfuncs.Printf(args[0], args[1:]...), nil
}

// builtinVsprintf implements vsprintf(string $format, array $values): string
func builtinVsprintf(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewString(""), nil
	}
	return stringfuncs.Vsprintf(args[0], args[1]), nil
}

// builtinVprintf implements vprintf(string $format, array $values): int
func builtinVprintf(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(0), nil
	}
	return stringfuncs.Vprintf(args[0], args[1]), nil
}

// builtinNumberFormat implements number_format(float $num, int $decimals = 0, ?string $decimal_separator = ".", ?string $thousands_separator = ","): string
func builtinNumberFormat(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString("0"), nil
	}
	return stringfuncs.NumberFormat(args[0], args[1:]...), nil
}

// builtinSscanf implements sscanf(string $string, string $format, mixed &...$vars): array|int|null
func builtinSscanf(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewNull(), nil
	}
	return stringfuncs.Sscanf(args[0], args[1], args[2:]...), nil
}

// builtinStrGetcsv implements str_getcsv(string $string, string $separator = ",", string $enclosure = "\"", string $escape = "\\"): array
func builtinStrGetcsv(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return stringfuncs.StrGetcsv(args[0], args[1:]...), nil
}

// ============================================================================
// Array Functions
// ============================================================================

// builtinCount implements count(Countable|array $value, int $mode = COUNT_NORMAL): int
func builtinCount(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return arrayfuncs.Count(args[0]), nil
}

// builtinSizeof implements sizeof(Countable|array $value, int $mode = COUNT_NORMAL): int (alias of count)
func builtinSizeof(args []*types.Value) (*types.Value, error) {
	return builtinCount(args)
}

// builtinArrayKeys implements array_keys(array $array, mixed $filter_value, bool $strict = false): array
func builtinArrayKeys(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayKeys(args[0]), nil
}

// builtinArrayValues implements array_values(array $array): array
func builtinArrayValues(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayValues(args[0]), nil
}

// builtinArrayPush implements array_push(array &$array, mixed ...$values): int
func builtinArrayPush(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return arrayfuncs.ArrayPush(args[0], args[1:]...), nil
}

// builtinArrayPop implements array_pop(array &$array): mixed
func builtinArrayPop(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return arrayfuncs.ArrayPop(args[0]), nil
}

// builtinArrayShift implements array_shift(array &$array): mixed
func builtinArrayShift(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return arrayfuncs.ArrayShift(args[0]), nil
}

// builtinArrayUnshift implements array_unshift(array &$array, mixed ...$values): int
func builtinArrayUnshift(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return arrayfuncs.ArrayUnshift(args[0], args[1:]...), nil
}

// builtinArrayMerge implements array_merge(array ...$arrays): array
func builtinArrayMerge(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayMerge(args...), nil
}

// builtinInArray implements in_array(mixed $needle, array $haystack, bool $strict = false): bool
func builtinInArray(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.InArray(args[0], args[1], args[2:]...), nil
}

// builtinArraySearch implements array_search(mixed $needle, array $haystack, bool $strict = false): int|string|false
func builtinArraySearch(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.ArraySearch(args[0], args[1], args[2:]...), nil
}

// builtinArraySlice implements array_slice(array $array, int $offset, ?int $length = null, bool $preserve_keys = false): array
func builtinArraySlice(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArraySlice(args[0], args[1], args[2:]...), nil
}

// builtinArrayReverse implements array_reverse(array $array, bool $preserve_keys = false): array
func builtinArrayReverse(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayReverse(args[0], args[1:]...), nil
}

// builtinArrayUnique implements array_unique(array $array, int $flags = SORT_STRING): array
func builtinArrayUnique(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayUnique(args[0]), nil
}

// builtinArrayCombine implements array_combine(array $keys, array $values): array|false
func builtinArrayCombine(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.ArrayCombine(args[0], args[1]), nil
}

// builtinArrayFlip implements array_flip(array $array): array
func builtinArrayFlip(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayFlip(args[0]), nil
}

// builtinArrayFill implements array_fill(int $start_index, int $count, mixed $value): array
func builtinArrayFill(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayFill(args[0], args[1], args[2]), nil
}

// builtinArrayChunk implements array_chunk(array $array, int $length, bool $preserve_keys = false): array
func builtinArrayChunk(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayChunk(args[0], args[1], args[2:]...), nil
}

// builtinSort implements sort(array &$array, int $flags = SORT_REGULAR): true
func builtinSort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Sort(args[0], args[1:]...), nil
}

// builtinRsort implements rsort(array &$array, int $flags = SORT_REGULAR): true
func builtinRsort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Rsort(args[0], args[1:]...), nil
}

// builtinAsort implements asort(array &$array, int $flags = SORT_REGULAR): true
func builtinAsort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Asort(args[0], args[1:]...), nil
}

// builtinArsort implements arsort(array &$array, int $flags = SORT_REGULAR): true
func builtinArsort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Arsort(args[0], args[1:]...), nil
}

// builtinKsort implements ksort(array &$array, int $flags = SORT_REGULAR): true
func builtinKsort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Ksort(args[0], args[1:]...), nil
}

// builtinKrsort implements krsort(array &$array, int $flags = SORT_REGULAR): true
func builtinKrsort(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return arrayfuncs.Krsort(args[0], args[1:]...), nil
}

// builtinArrayMap implements array_map(?callable $callback, array $array, array ...$arrays): array
func builtinArrayMap(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayMap(args[0], args[1:]...), nil
}

// builtinArrayFilter implements array_filter(array $array, ?callable $callback = null, int $mode = 0): array
func builtinArrayFilter(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayFilter(args[0], args[1:]...), nil
}

// builtinArrayReduce implements array_reduce(array $array, callable $callback, mixed $initial = null): mixed
func builtinArrayReduce(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewNull(), nil
	}
	return arrayfuncs.ArrayReduce(args[0], args[1], args[2:]...), nil
}

// builtinArrayWalk implements array_walk(array &$array, callable $callback, mixed $arg = null): bool
func builtinArrayWalk(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.ArrayWalk(args[0], args[1], args[2:]...), nil
}

// builtinArrayWalkRecursive implements array_walk_recursive(array &$array, callable $callback, mixed $arg = null): bool
func builtinArrayWalkRecursive(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.ArrayWalkRecursive(args[0], args[1], args[2:]...), nil
}

// builtinArrayDiff implements array_diff(array $array, array ...$arrays): array
func builtinArrayDiff(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayDiff(args...), nil
}

// builtinArrayIntersect implements array_intersect(array $array, array ...$arrays): array
func builtinArrayIntersect(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayIntersect(args...), nil
}

// builtinArrayPad implements array_pad(array $array, int $length, mixed $value): array
func builtinArrayPad(args []*types.Value) (*types.Value, error) {
	if len(args) < 3 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayPad(args[0], args[1], args[2]), nil
}

// builtinArraySum implements array_sum(array $array): int|float
func builtinArraySum(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return arrayfuncs.ArraySum(args[0]), nil
}

// builtinArrayProduct implements array_product(array $array): int|float
func builtinArrayProduct(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return arrayfuncs.ArrayProduct(args[0]), nil
}

// builtinArrayColumn implements array_column(array $array, int|string|null $column_key, int|string|null $index_key = null): array
func builtinArrayColumn(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	if len(args) >= 3 {
		return arrayfuncs.ArrayColumn(args[0], args[1], args[2]), nil
	}
	return arrayfuncs.ArrayColumn(args[0], args[1]), nil
}

// builtinArrayChangeKeyCase implements array_change_key_case(array $array, int $case = CASE_LOWER): array
func builtinArrayChangeKeyCase(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	if len(args) >= 2 {
		return arrayfuncs.ArrayChangeKeyCase(args[0], args[1]), nil
	}
	return arrayfuncs.ArrayChangeKeyCase(args[0]), nil
}

// builtinArrayReplace implements array_replace(array $array, array ...$replacements): array
func builtinArrayReplace(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayReplace(args...), nil
}

// builtinArrayReplaceRecursive implements array_replace_recursive(array $array, array ...$replacements): array
func builtinArrayReplaceRecursive(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayReplaceRecursive(args...), nil
}

// builtinCurrent implements current(array $array): mixed
func builtinCurrent(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.Current(args[0]), nil
}

// builtinKey implements key(array $array): int|string|null
func builtinKey(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return arrayfuncs.Key(args[0]), nil
}

// builtinReset implements reset(array &$array): mixed
func builtinReset(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.Reset(args[0]), nil
}

// builtinEnd implements end(array &$array): mixed
func builtinEnd(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.End(args[0]), nil
}

// builtinNext implements next(array &$array): mixed
func builtinNext(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.Next(args[0]), nil
}

// builtinPrev implements prev(array &$array): mixed
func builtinPrev(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.Prev(args[0]), nil
}

// builtinArrayKeyExists implements array_key_exists(int|string $key, array $array): bool
func builtinArrayKeyExists(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.ArrayKeyExists(args[0], args[1]), nil
}

// builtinKeyExists implements key_exists(int|string $key, array $array): bool
func builtinKeyExists(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return arrayfuncs.KeyExists(args[0], args[1]), nil
}

// builtinArrayDiffKey implements array_diff_key(array $array, array ...$arrays): array
func builtinArrayDiffKey(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayDiffKey(args...), nil
}

// builtinArrayIntersectKey implements array_intersect_key(array $array, array ...$arrays): array
func builtinArrayIntersectKey(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayIntersectKey(args...), nil
}

// builtinArrayDiffAssoc implements array_diff_assoc(array $array, array ...$arrays): array
func builtinArrayDiffAssoc(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayDiffAssoc(args...), nil
}

// builtinArrayIntersectAssoc implements array_intersect_assoc(array $array, array ...$arrays): array
func builtinArrayIntersectAssoc(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayIntersectAssoc(args...), nil
}

// builtinArrayKeyFirst implements array_key_first(array $array): int|string|null
func builtinArrayKeyFirst(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return arrayfuncs.ArrayKeyFirst(args[0]), nil
}

// builtinArrayKeyLast implements array_key_last(array $array): int|string|null
func builtinArrayKeyLast(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return arrayfuncs.ArrayKeyLast(args[0]), nil
}

// builtinArrayCountValues implements array_count_values(array $array): array
func builtinArrayCountValues(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewArray(types.NewEmptyArray()), nil
	}
	return arrayfuncs.ArrayCountValues(args[0]), nil
}

// ============================================================================
// JSON Functions
// ============================================================================

// builtinJsonEncode implements json_encode(mixed $value, int $flags = 0, int $depth = 512): string|false
func builtinJsonEncode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return jsonfuncs.JsonEncode(args[0], args[1:]...), nil
}

// builtinJsonDecode implements json_decode(string $json, bool $assoc = false, int $depth = 512, int $flags = 0): mixed
func builtinJsonDecode(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return jsonfuncs.JsonDecode(args[0], args[1:]...), nil
}

// builtinJsonLastError implements json_last_error(): int
func builtinJsonLastError(args []*types.Value) (*types.Value, error) {
	return jsonfuncs.JsonLastError(), nil
}

// builtinJsonLastErrorMsg implements json_last_error_msg(): string
func builtinJsonLastErrorMsg(args []*types.Value) (*types.Value, error) {
	return jsonfuncs.JsonLastErrorMsg(), nil
}

// ============================================================================
// Type/Existence Functions
// ============================================================================

// init registers function_exists after builtinFunctions is initialized to avoid cycle
func init() {
	builtinFunctions["function_exists"] = builtinFunctionExistsImpl
}

// builtinFunctionExistsImpl implements function_exists(string $function_name): bool
func builtinFunctionExistsImpl(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	funcName := args[0].ToString()
	// Check if it's a builtin function
	_, exists := builtinFunctions[funcName]
	return types.NewBool(exists), nil
}

// builtinClassExists implements class_exists(string $class, bool $autoload = true): bool
func builtinClassExists(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	// For now, use the stub - full implementation needs VM class registry access
	return varfuncs.ClassExistsStub(args[0], args[1:]...), nil
}

// builtinMethodExists implements method_exists(object|string $object_or_class, string $method): bool
func builtinMethodExists(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewBool(false), nil
	}
	return varfuncs.MethodExistsStub(args[0], args[1]), nil
}

// ============================================================================
// Serialization Functions
// ============================================================================

// builtinSerialize implements serialize(mixed $value): string
func builtinSerialize(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewString("N;"), nil
	}
	return varfuncs.Serialize(args[0]), nil
}

// builtinUnserialize implements unserialize(string $data, array $options = []): mixed
func builtinUnserialize(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return varfuncs.Unserialize(args[0], args[1:]...), nil
}

// ============================================================================
// Math Functions
// ============================================================================

// builtinAbs implements abs(int|float $num): int|float
func builtinAbs(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewInt(0), nil
	}
	return mathfuncs.Abs(args[0]), nil
}

// builtinCeil implements ceil(int|float $num): float
func builtinCeil(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Ceil(args[0]), nil
}

// builtinFloor implements floor(int|float $num): float
func builtinFloor(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Floor(args[0]), nil
}

// builtinRound implements round(int|float $num, int $precision = 0): float
func builtinRound(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Round(args[0], args[1:]...), nil
}

// builtinMin implements min(mixed ...$values): mixed
func builtinMin(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return mathfuncs.Min(args...), nil
}

// builtinMax implements max(mixed ...$values): mixed
func builtinMax(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewNull(), nil
	}
	return mathfuncs.Max(args...), nil
}

// builtinPow implements pow(mixed $base, mixed $exp): int|float
func builtinPow(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(1), nil
	}
	return mathfuncs.Pow(args[0], args[1]), nil
}

// builtinSqrt implements sqrt(float $num): float
func builtinSqrt(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Sqrt(args[0]), nil
}

// builtinSin implements sin(float $num): float
func builtinSin(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Sin(args[0]), nil
}

// builtinCos implements cos(float $num): float
func builtinCos(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(1), nil
	}
	return mathfuncs.Cos(args[0]), nil
}

// builtinTan implements tan(float $num): float
func builtinTan(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Tan(args[0]), nil
}

// builtinAsin implements asin(float $num): float
func builtinAsin(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Asin(args[0]), nil
}

// builtinAcos implements acos(float $num): float
func builtinAcos(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Acos(args[0]), nil
}

// builtinAtan implements atan(float $num): float
func builtinAtan(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Atan(args[0]), nil
}

// builtinAtan2 implements atan2(float $y, float $x): float
func builtinAtan2(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Atan2(args[0], args[1]), nil
}

// builtinDeg2rad implements deg2rad(float $num): float
func builtinDeg2rad(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Deg2rad(args[0]), nil
}

// builtinRad2deg implements rad2deg(float $num): float
func builtinRad2deg(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Rad2deg(args[0]), nil
}

// builtinExp implements exp(float $num): float
func builtinExp(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(1), nil
	}
	return mathfuncs.Exp(args[0]), nil
}

// builtinLog implements log(float $num, float $base = M_E): float
func builtinLog(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Log(args[0], args[1:]...), nil
}

// builtinLog10 implements log10(float $num): float
func builtinLog10(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Log10(args[0]), nil
}

// builtinLog1p implements log1p(float $num): float
func builtinLog1p(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Log1p(args[0]), nil
}

// builtinExpm1 implements expm1(float $num): float
func builtinExpm1(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Expm1(args[0]), nil
}

// builtinRand implements rand(int $min = 0, int $max = getrandmax()): int
func builtinRand(args []*types.Value) (*types.Value, error) {
	return mathfuncs.Rand(args...), nil
}

// builtinMtRand implements mt_rand(int $min = 0, int $max = mt_getrandmax()): int
func builtinMtRand(args []*types.Value) (*types.Value, error) {
	return mathfuncs.MtRand(args...), nil
}

// builtinRandomInt implements random_int(int $min, int $max): int
func builtinRandomInt(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewInt(0), nil
	}
	return mathfuncs.RandomInt(args[0], args[1]), nil
}

// builtinGetRandMax implements getrandmax(): int
func builtinGetRandMax(args []*types.Value) (*types.Value, error) {
	return mathfuncs.GetRandMax(), nil
}

// builtinMtGetRandMax implements mt_getrandmax(): int
func builtinMtGetRandMax(args []*types.Value) (*types.Value, error) {
	return mathfuncs.MtGetRandMax(), nil
}

// builtinPi implements pi(): float
func builtinPi(args []*types.Value) (*types.Value, error) {
	return mathfuncs.Pi(), nil
}

// builtinIsNan implements is_nan(float $num): bool
func builtinIsNan(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return mathfuncs.IsNan(args[0]), nil
}

// builtinIsInfinite implements is_infinite(float $num): bool
func builtinIsInfinite(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(false), nil
	}
	return mathfuncs.IsInfinite(args[0]), nil
}

// builtinIsFinite implements is_finite(float $num): bool
func builtinIsFinite(args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return types.NewBool(true), nil
	}
	return mathfuncs.IsFinite(args[0]), nil
}

// builtinHypot implements hypot(float $x, float $y): float
func builtinHypot(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Hypot(args[0], args[1]), nil
}

// builtinFmod implements fmod(float $x, float $y): float
func builtinFmod(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Fmod(args[0], args[1]), nil
}

// builtinIntdiv implements intdiv(int $num1, int $num2): int
func builtinIntdiv(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewNull(), nil
	}
	return mathfuncs.Intdiv(args[0], args[1]), nil
}

// builtinFdiv implements fdiv(float $num1, float $num2): float
func builtinFdiv(args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return types.NewFloat(0), nil
	}
	return mathfuncs.Fdiv(args[0], args[1]), nil
}
