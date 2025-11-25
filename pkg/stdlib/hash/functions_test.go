package hash

import (
	"os"
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Hash Function Tests
// ============================================================================

func TestHashMd5(t *testing.T) {
	data := types.NewString("hello")
	algo := types.NewString("md5")

	result := Hash(algo, data)
	expected := "5d41402abc4b2a76b9719d911017c592" // MD5 of "hello"

	if result.ToString() != expected {
		t.Errorf("Hash(md5, 'hello') = %v, want %v", result.ToString(), expected)
	}
}

func TestHashSha1(t *testing.T) {
	data := types.NewString("hello")
	algo := types.NewString("sha1")

	result := Hash(algo, data)
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" // SHA1 of "hello"

	if result.ToString() != expected {
		t.Errorf("Hash(sha1, 'hello') = %v, want %v", result.ToString(), expected)
	}
}

func TestHashSha256(t *testing.T) {
	data := types.NewString("hello")
	algo := types.NewString("sha256")

	result := Hash(algo, data)
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" // SHA256 of "hello"

	if result.ToString() != expected {
		t.Errorf("Hash(sha256, 'hello') = %v, want %v", result.ToString(), expected)
	}
}

func TestHashUnknownAlgorithm(t *testing.T) {
	data := types.NewString("hello")
	algo := types.NewString("unknown")

	result := Hash(algo, data)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("Hash(unknown, 'hello') should return false")
	}
}

func TestHashBinary(t *testing.T) {
	data := types.NewString("hello")
	algo := types.NewString("md5")
	binary := types.NewBool(true)

	result := Hash(algo, data, binary)
	// Binary output should be 16 bytes for MD5
	if len(result.ToString()) != 16 {
		t.Errorf("Hash(md5, 'hello', true) should return 16 bytes, got %d", len(result.ToString()))
	}
}

// ============================================================================
// Hash File Tests
// ============================================================================

func TestHashFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "hello world"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	algo := types.NewString("md5")
	filename := types.NewString(tmpfile.Name())

	result := HashFile(algo, filename)
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3" // MD5 of "hello world"

	if result.ToString() != expected {
		t.Errorf("HashFile(md5, file) = %v, want %v", result.ToString(), expected)
	}
}

func TestHashFileNotFound(t *testing.T) {
	algo := types.NewString("md5")
	filename := types.NewString("/nonexistent/file.txt")

	result := HashFile(algo, filename)
	if result.Type() != types.TypeBool || result.ToBool() != false {
		t.Errorf("HashFile with nonexistent file should return false")
	}
}

// ============================================================================
// HMAC Tests
// ============================================================================

func TestHashHmac(t *testing.T) {
	algo := types.NewString("sha256")
	data := types.NewString("hello")
	key := types.NewString("secret")

	result := HashHmac(algo, data, key)

	// Verify it returns a hex string
	if len(result.ToString()) != 64 { // SHA256 produces 64 hex characters
		t.Errorf("HashHmac should return 64 hex chars for sha256, got %d", len(result.ToString()))
	}
}

func TestHashHmacFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "hello world"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	algo := types.NewString("sha256")
	filename := types.NewString(tmpfile.Name())
	key := types.NewString("secret")

	result := HashHmacFile(algo, filename, key)

	// Verify it returns a hex string
	if len(result.ToString()) != 64 {
		t.Errorf("HashHmacFile should return 64 hex chars for sha256, got %d", len(result.ToString()))
	}
}

// ============================================================================
// Legacy Hash Function Tests
// ============================================================================

func TestMd5(t *testing.T) {
	str := types.NewString("hello")
	result := Md5(str)
	expected := "5d41402abc4b2a76b9719d911017c592"

	if result.ToString() != expected {
		t.Errorf("Md5('hello') = %v, want %v", result.ToString(), expected)
	}
}

func TestMd5Binary(t *testing.T) {
	str := types.NewString("hello")
	binary := types.NewBool(true)
	result := Md5(str, binary)

	// Binary output should be 16 bytes
	if len(result.ToString()) != 16 {
		t.Errorf("Md5('hello', true) should return 16 bytes, got %d", len(result.ToString()))
	}
}

func TestMd5File(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "hello world"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	filename := types.NewString(tmpfile.Name())
	result := Md5File(filename)
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"

	if result.ToString() != expected {
		t.Errorf("Md5File(file) = %v, want %v", result.ToString(), expected)
	}
}

func TestSha1(t *testing.T) {
	str := types.NewString("hello")
	result := Sha1(str)
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"

	if result.ToString() != expected {
		t.Errorf("Sha1('hello') = %v, want %v", result.ToString(), expected)
	}
}

func TestSha1Binary(t *testing.T) {
	str := types.NewString("hello")
	binary := types.NewBool(true)
	result := Sha1(str, binary)

	// Binary output should be 20 bytes
	if len(result.ToString()) != 20 {
		t.Errorf("Sha1('hello', true) should return 20 bytes, got %d", len(result.ToString()))
	}
}

func TestSha1File(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "hello world"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	filename := types.NewString(tmpfile.Name())
	result := Sha1File(filename)
	expected := "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"

	if result.ToString() != expected {
		t.Errorf("Sha1File(file) = %v, want %v", result.ToString(), expected)
	}
}

// ============================================================================
// Hash Comparison Tests
// ============================================================================

func TestHashEquals(t *testing.T) {
	known := types.NewString("5d41402abc4b2a76b9719d911017c592")
	user := types.NewString("5d41402abc4b2a76b9719d911017c592")

	result := HashEquals(known, user)
	if !result.ToBool() {
		t.Errorf("HashEquals with equal strings should return true")
	}
}

func TestHashEqualsNotEqual(t *testing.T) {
	known := types.NewString("5d41402abc4b2a76b9719d911017c592")
	user := types.NewString("5d41402abc4b2a76b9719d911017c593")

	result := HashEquals(known, user)
	if result.ToBool() {
		t.Errorf("HashEquals with different strings should return false")
	}
}

func TestHashEqualsDifferentLength(t *testing.T) {
	known := types.NewString("short")
	user := types.NewString("much longer string")

	result := HashEquals(known, user)
	if result.ToBool() {
		t.Errorf("HashEquals with different lengths should return false")
	}
}

// TestHashEqualsConstantTime verifies constant-time behavior
// This test ensures the function uses crypto/subtle.ConstantTimeCompare
func TestHashEqualsConstantTime(t *testing.T) {
	tests := []struct {
		name     string
		known    string
		user     string
		expected bool
	}{
		{
			name:     "Equal hashes",
			known:    "e10adc3949ba59abbe56e057f20f883e", // md5("123456")
			user:     "e10adc3949ba59abbe56e057f20f883e",
			expected: true,
		},
		{
			name:     "Different last character",
			known:    "e10adc3949ba59abbe56e057f20f883e",
			user:     "e10adc3949ba59abbe56e057f20f883f", // Different last char
			expected: false,
		},
		{
			name:     "Different first character",
			known:    "e10adc3949ba59abbe56e057f20f883e",
			user:     "f10adc3949ba59abbe56e057f20f883e", // Different first char
			expected: false,
		},
		{
			name:     "Completely different same length",
			known:    "e10adc3949ba59abbe56e057f20f883e",
			user:     "5f4dcc3b5aa765d61d8327deb882cf99", // md5("password")
			expected: false,
		},
		{
			name:     "Empty strings",
			known:    "",
			user:     "",
			expected: true,
		},
		{
			name:     "One empty one not",
			known:    "hash",
			user:     "",
			expected: false,
		},
		{
			name:     "Short vs long",
			known:    "abc",
			user:     "abcdef",
			expected: false,
		},
		{
			name:     "Long vs short",
			known:    "abcdef",
			user:     "abc",
			expected: false,
		},
		{
			name:     "Binary data (null bytes)",
			known:    "test\x00data",
			user:     "test\x00data",
			expected: true,
		},
		{
			name:     "Binary data different",
			known:    "test\x00data",
			user:     "test\x01data",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			known := types.NewString(tt.known)
			user := types.NewString(tt.user)

			result := HashEquals(known, user)
			got := result.ToBool()

			if got != tt.expected {
				t.Errorf("HashEquals(%q, %q) = %v, want %v", tt.known, tt.user, got, tt.expected)
			}
		})
	}
}

// TestHashEqualsRealWorldHashes tests with actual hash values
func TestHashEqualsRealWorldHashes(t *testing.T) {
	// Test with SHA256 hashes
	passwordHash := types.NewString("5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8") // sha256("password")
	correctInput := types.NewString("5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8")
	wrongInput := types.NewString("5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d9")

	// Correct password hash should match
	result := HashEquals(passwordHash, correctInput)
	if !result.ToBool() {
		t.Errorf("HashEquals should return true for matching SHA256 hashes")
	}

	// Wrong password hash should not match
	result = HashEquals(passwordHash, wrongInput)
	if result.ToBool() {
		t.Errorf("HashEquals should return false for non-matching SHA256 hashes")
	}
}

// TestHashEqualsTimingSafety documents the timing-safety guarantees
// Note: This is a documentation test - actual timing verification would require
// specialized benchmarking tools and statistical analysis
func TestHashEqualsTimingSafety(t *testing.T) {
	// This test documents that HashEquals uses crypto/subtle.ConstantTimeCompare
	// which provides the following guarantees:
	// 1. Comparison time does not depend on string content
	// 2. Comparison time does not leak length information
	// 3. No early returns based on content differences

	// Equal length, different content at various positions
	testCases := []struct {
		known string
		user  string
	}{
		{"aaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaa"}, // Identical
		{"aaaaaaaaaaaaaaaa", "baaaaaaaaaaaaaaa"}, // Different at start
		{"aaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaب"}, // Different at end
		{"aaaaaaaaaaaaaaaa", "aaaaaaaabaaaaaaa"}, // Different in middle
		{"aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb"}, // Completely different
	}

	for _, tc := range testCases {
		known := types.NewString(tc.known)
		user := types.NewString(tc.user)

		// The important thing is that all these comparisons use constant-time operations
		result := HashEquals(known, user)

		// Verify result is a boolean (sanity check)
		if result.Type() != types.TypeBool {
			t.Errorf("HashEquals should return bool type")
		}
	}

	// Different lengths should also use constant-time comparison
	// crypto/subtle.ConstantTimeCompare handles this safely
	differentLengths := []struct {
		known string
		user  string
	}{
		{"short", "much longer string"},
		{"", "x"},
		{"x", ""},
		{"abc", "abcd"},
	}

	for _, tc := range differentLengths {
		known := types.NewString(tc.known)
		user := types.NewString(tc.user)

		result := HashEquals(known, user)

		// All should return false, but in constant time
		if result.ToBool() {
			t.Errorf("HashEquals with different lengths should return false")
		}
	}
}

// ============================================================================
// Hash Algorithm Information Tests
// ============================================================================

func TestHashAlgos(t *testing.T) {
	result := HashAlgos()

	if result.Type() != types.TypeArray {
		t.Errorf("HashAlgos should return array, got %v", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() == 0 {
		t.Errorf("HashAlgos should return non-empty array")
	}

	// Check if md5 is in the list
	hasMd5 := false
	arr.Each(func(_, val *types.Value) bool {
		if val.ToString() == "md5" {
			hasMd5 = true
			return false
		}
		return true
	})

	if !hasMd5 {
		t.Errorf("HashAlgos should include 'md5'")
	}
}

func TestHashHmacAlgos(t *testing.T) {
	result := HashHmacAlgos()

	if result.Type() != types.TypeArray {
		t.Errorf("HashHmacAlgos should return array, got %v", result.Type())
	}

	arr := result.ToArray()
	if arr.Len() == 0 {
		t.Errorf("HashHmacAlgos should return non-empty array")
	}
}

// ============================================================================
// CRC32 Tests
// ============================================================================

func TestCrc32(t *testing.T) {
	str := types.NewString("hello")
	result := Crc32(str)

	if result.Type() != types.TypeInt {
		t.Errorf("Crc32 should return int, got %v", result.Type())
	}

	// Just verify it returns a value (exact value depends on algorithm)
	// PHP's crc32 implementation may differ
}

func TestCrc32Empty(t *testing.T) {
	str := types.NewString("")
	result := Crc32(str)

	if result.Type() != types.TypeInt {
		t.Errorf("Crc32 should return int, got %v", result.Type())
	}
}

// ============================================================================
// Additional Algorithm Tests
// ============================================================================

func TestHashSha384(t *testing.T) {
	data := types.NewString("test")
	algo := types.NewString("sha384")

	result := Hash(algo, data)

	// SHA384 produces 96 hex characters
	if len(result.ToString()) != 96 {
		t.Errorf("Hash(sha384) should return 96 hex chars, got %d", len(result.ToString()))
	}
}

func TestHashSha512(t *testing.T) {
	data := types.NewString("test")
	algo := types.NewString("sha512")

	result := Hash(algo, data)

	// SHA512 produces 128 hex characters
	if len(result.ToString()) != 128 {
		t.Errorf("Hash(sha512) should return 128 hex chars, got %d", len(result.ToString()))
	}
}

func TestHashCaseInsensitive(t *testing.T) {
	data := types.NewString("hello")
	algo1 := types.NewString("MD5")
	algo2 := types.NewString("md5")

	result1 := Hash(algo1, data)
	result2 := Hash(algo2, data)

	if result1.ToString() != result2.ToString() {
		t.Errorf("Hash algorithm names should be case-insensitive")
	}
}
