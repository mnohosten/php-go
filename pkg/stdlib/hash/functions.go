package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Hash Algorithms Map
// ============================================================================

// getHashAlgorithm returns the hash.Hash for the given algorithm name
func getHashAlgorithm(algo string) hash.Hash {
	algo = strings.ToLower(algo)

	switch algo {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha224":
		return sha256.New224()
	case "sha384":
		return sha512.New384()
	case "sha512":
		return sha512.New()
	case "sha512/224":
		return sha512.New512_224()
	case "sha512/256":
		return sha512.New512_256()
	default:
		return nil
	}
}

// ============================================================================
// Hash Functions
// ============================================================================

// Hash generates a hash value (message digest)
// hash(string $algo, string $data, bool $binary = false): string
func Hash(algo, data *types.Value, binary ...*types.Value) *types.Value {
	h := getHashAlgorithm(algo.ToString())
	if h == nil {
		// Unknown algorithm
		return types.NewBool(false)
	}

	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	// Compute hash
	h.Write([]byte(data.ToString()))
	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// HashFile generates a hash value for a file
// hash_file(string $algo, string $filename, bool $binary = false): string|false
func HashFile(algo, filename *types.Value, binary ...*types.Value) *types.Value {
	h := getHashAlgorithm(algo.ToString())
	if h == nil {
		return types.NewBool(false)
	}

	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	// Read file
	file, err := os.Open(filename.ToString())
	if err != nil {
		return types.NewBool(false)
	}
	defer file.Close()

	// Compute hash
	if _, err := io.Copy(h, file); err != nil {
		return types.NewBool(false)
	}

	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// HashHmac generates a keyed hash value using HMAC
// hash_hmac(string $algo, string $data, string $key, bool $binary = false): string
func HashHmac(algo, data, key *types.Value, binary ...*types.Value) *types.Value {
	h := getHashAlgorithm(algo.ToString())
	if h == nil {
		return types.NewBool(false)
	}

	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	// Create HMAC
	mac := hmac.New(func() hash.Hash {
		return getHashAlgorithm(algo.ToString())
	}, []byte(key.ToString()))

	mac.Write([]byte(data.ToString()))
	hashBytes := mac.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// HashHmacFile generates a keyed hash value for a file using HMAC
// hash_hmac_file(string $algo, string $filename, string $key, bool $binary = false): string|false
func HashHmacFile(algo, filename, key *types.Value, binary ...*types.Value) *types.Value {
	h := getHashAlgorithm(algo.ToString())
	if h == nil {
		return types.NewBool(false)
	}

	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	// Read file
	file, err := os.Open(filename.ToString())
	if err != nil {
		return types.NewBool(false)
	}
	defer file.Close()

	// Create HMAC
	mac := hmac.New(func() hash.Hash {
		return getHashAlgorithm(algo.ToString())
	}, []byte(key.ToString()))

	if _, err := io.Copy(mac, file); err != nil {
		return types.NewBool(false)
	}

	hashBytes := mac.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// ============================================================================
// Legacy Hash Functions (md5, sha1)
// ============================================================================

// Md5 calculates the MD5 hash of a string
// md5(string $str, bool $binary = false): string
func Md5(str *types.Value, binary ...*types.Value) *types.Value {
	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	h := md5.New()
	h.Write([]byte(str.ToString()))
	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// Md5File calculates the MD5 hash of a file
// md5_file(string $filename, bool $binary = false): string|false
func Md5File(filename *types.Value, binary ...*types.Value) *types.Value {
	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	file, err := os.Open(filename.ToString())
	if err != nil {
		return types.NewBool(false)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return types.NewBool(false)
	}

	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// Sha1 calculates the SHA1 hash of a string
// sha1(string $str, bool $binary = false): string
func Sha1(str *types.Value, binary ...*types.Value) *types.Value {
	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	h := sha1.New()
	h.Write([]byte(str.ToString()))
	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// Sha1File calculates the SHA1 hash of a file
// sha1_file(string $filename, bool $binary = false): string|false
func Sha1File(filename *types.Value, binary ...*types.Value) *types.Value {
	rawBinary := false
	if len(binary) > 0 && binary[0] != nil {
		rawBinary = binary[0].ToBool()
	}

	file, err := os.Open(filename.ToString())
	if err != nil {
		return types.NewBool(false)
	}
	defer file.Close()

	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		return types.NewBool(false)
	}

	hashBytes := h.Sum(nil)

	if rawBinary {
		return types.NewString(string(hashBytes))
	}

	return types.NewString(hex.EncodeToString(hashBytes))
}

// ============================================================================
// Hash Comparison
// ============================================================================

// HashEquals performs a timing-safe string comparison
// hash_equals(string $known_string, string $user_string): bool
func HashEquals(knownString, userString *types.Value) *types.Value {
	known := knownString.ToString()
	user := userString.ToString()

	// Different lengths - definitely not equal
	if len(known) != len(user) {
		return types.NewBool(false)
	}

	// Timing-safe comparison
	result := 0
	for i := 0; i < len(known); i++ {
		result |= int(known[i]) ^ int(user[i])
	}

	return types.NewBool(result == 0)
}

// ============================================================================
// Hash Algorithm Information
// ============================================================================

// HashAlgos returns a list of registered hashing algorithms
// hash_algos(): array
func HashAlgos() *types.Value {
	algos := []string{
		"md5",
		"sha1",
		"sha224",
		"sha256",
		"sha384",
		"sha512",
		"sha512/224",
		"sha512/256",
	}

	arr := types.NewEmptyArray()
	for _, algo := range algos {
		arr.Append(types.NewString(algo))
	}

	return types.NewArray(arr)
}

// HashHmacAlgos returns a list of registered hashing algorithms suitable for hash_hmac
// hash_hmac_algos(): array
func HashHmacAlgos() *types.Value {
	// Same as hash_algos for our implementation
	return HashAlgos()
}

// ============================================================================
// Additional Hash Functions
// ============================================================================

// Crc32 calculates the CRC32 polynomial of a string
// crc32(string $str): int
func Crc32(str *types.Value) *types.Value {
	// Note: PHP's crc32 uses a different polynomial than Go's hash/crc32
	// This is a simplified implementation
	// For full compatibility, we'd need to match PHP's exact algorithm
	data := []byte(str.ToString())

	// Simple CRC32 implementation (not the same as PHP's)
	// In production, would need exact PHP compatibility
	var crc uint32 = 0xFFFFFFFF
	for _, b := range data {
		crc ^= uint32(b)
		for j := 0; j < 8; j++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xEDB88320
			} else {
				crc >>= 1
			}
		}
	}

	return types.NewInt(int64(^crc))
}

// HashPbkdf2 generates a PBKDF2 key derivation of a password
// hash_pbkdf2(string $algo, string $password, string $salt, int $iterations, int $length = 0, bool $binary = false): string
func HashPbkdf2(algo, password, salt, iterations *types.Value, length ...*types.Value) *types.Value {
	// This would require implementing PBKDF2
	// For now, return a placeholder
	// TODO: Implement full PBKDF2 support
	return types.NewBool(false)
}
