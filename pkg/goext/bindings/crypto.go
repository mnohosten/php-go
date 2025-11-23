package bindings

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// CryptoExtension provides Go crypto functions for PHP.
type CryptoExtension struct {
	*goext.BaseExtension
}

// NewCryptoExtension creates a new crypto extension.
func NewCryptoExtension() *CryptoExtension {
	ext := &CryptoExtension{
		BaseExtension: goext.NewBaseExtension("go_crypto", "1.0.0"),
	}

	// Register hash functions
	ext.AddFunction("go_hash_md5", ext.hashMD5)
	ext.AddFunction("go_hash_sha1", ext.hashSHA1)
	ext.AddFunction("go_hash_sha256", ext.hashSHA256)
	ext.AddFunction("go_hash_sha512", ext.hashSHA512)

	// Register encoding functions
	ext.AddFunction("go_base64_encode", ext.base64Encode)
	ext.AddFunction("go_base64_decode", ext.base64Decode)
	ext.AddFunction("go_hex_encode", ext.hexEncode)
	ext.AddFunction("go_hex_decode", ext.hexDecode)

	// Register encryption functions
	ext.AddFunction("go_aes_encrypt", ext.aesEncrypt)
	ext.AddFunction("go_aes_decrypt", ext.aesDecrypt)
	ext.AddFunction("go_random_bytes", ext.randomBytes)

	// Register constants
	ext.AddConstant("CRYPTO_AES_BLOCK_SIZE", types.NewInt(aes.BlockSize))

	return ext
}

// hashMD5 implements go_hash_md5($data)
//
// Returns: MD5 hash as hex string
func (e *CryptoExtension) hashMD5(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hash_md5() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	hash := md5.Sum([]byte(data))

	return types.NewString(hex.EncodeToString(hash[:])), nil
}

// hashSHA1 implements go_hash_sha1($data)
//
// Returns: SHA1 hash as hex string
func (e *CryptoExtension) hashSHA1(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hash_sha1() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	hash := sha1.Sum([]byte(data))

	return types.NewString(hex.EncodeToString(hash[:])), nil
}

// hashSHA256 implements go_hash_sha256($data)
//
// Returns: SHA256 hash as hex string
func (e *CryptoExtension) hashSHA256(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hash_sha256() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	hash := sha256.Sum256([]byte(data))

	return types.NewString(hex.EncodeToString(hash[:])), nil
}

// hashSHA512 implements go_hash_sha512($data)
//
// Returns: SHA512 hash as hex string
func (e *CryptoExtension) hashSHA512(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hash_sha512() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	hash := sha512.Sum512([]byte(data))

	return types.NewString(hex.EncodeToString(hash[:])), nil
}

// base64Encode implements go_base64_encode($data)
//
// Returns: base64 encoded string
func (e *CryptoExtension) base64Encode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_base64_encode() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return types.NewString(encoded), nil
}

// base64Decode implements go_base64_decode($data)
//
// Returns: decoded string or false on error
func (e *CryptoExtension) base64Decode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_base64_decode() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(string(decoded)), nil
}

// hexEncode implements go_hex_encode($data)
//
// Returns: hex encoded string
func (e *CryptoExtension) hexEncode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hex_encode() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	encoded := hex.EncodeToString([]byte(data))

	return types.NewString(encoded), nil
}

// hexDecode implements go_hex_decode($data)
//
// Returns: decoded string or false on error
func (e *CryptoExtension) hexDecode(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_hex_decode() expects 1 argument (data), got %d", len(args))
	}

	data := args[0].ToString()
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(string(decoded)), nil
}

// aesEncrypt implements go_aes_encrypt($data, $key)
//
// Encrypts data using AES-256-GCM
// Key must be 32 bytes (256 bits)
//
// Returns: encrypted data (base64 encoded) or false on error
func (e *CryptoExtension) aesEncrypt(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_aes_encrypt() expects 2 arguments (data, key), got %d", len(args))
	}

	data := []byte(args[0].ToString())
	key := []byte(args[1].ToString())

	// Key must be 16, 24, or 32 bytes
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return types.NewBool(false), nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return types.NewBool(false), nil
	}

	// Use GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return types.NewBool(false), nil
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return types.NewBool(false), nil
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Return base64 encoded
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return types.NewString(encoded), nil
}

// aesDecrypt implements go_aes_decrypt($encrypted, $key)
//
// Decrypts data using AES-256-GCM
// Encrypted data must be base64 encoded
//
// Returns: decrypted data or false on error
func (e *CryptoExtension) aesDecrypt(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_aes_decrypt() expects 2 arguments (encrypted, key), got %d", len(args))
	}

	encryptedStr := args[0].ToString()
	key := []byte(args[1].ToString())

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return types.NewBool(false), nil
	}

	// Key must be 16, 24, or 32 bytes
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return types.NewBool(false), nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return types.NewBool(false), nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return types.NewBool(false), nil
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return types.NewBool(false), nil
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(string(plaintext)), nil
}

// randomBytes implements go_random_bytes($length)
//
// Returns: cryptographically secure random bytes as hex string
func (e *CryptoExtension) randomBytes(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_random_bytes() expects 1 argument (length), got %d", len(args))
	}

	length := int(args[0].ToInt())
	if length <= 0 {
		return types.NewBool(false), nil
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return types.NewBool(false), nil
	}

	return types.NewString(hex.EncodeToString(bytes)), nil
}
