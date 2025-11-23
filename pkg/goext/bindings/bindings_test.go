package bindings

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// JSON Extension Tests

func TestJSONExtension(t *testing.T) {
	ext := NewJSONExtension()

	if ext.Name() != "go_json" {
		t.Errorf("Expected name 'go_json', got '%s'", ext.Name())
	}

	funcs := ext.Functions()
	if _, ok := funcs["go_json_encode"]; !ok {
		t.Error("go_json_encode not registered")
	}
	if _, ok := funcs["go_json_decode"]; !ok {
		t.Error("go_json_decode not registered")
	}
}

func TestJSONEncode(t *testing.T) {
	ext := NewJSONExtension()
	testVM := vm.New()

	// Test simple string
	result, err := ext.jsonEncode(testVM, []*types.Value{
		types.NewString("hello"),
	})

	if err != nil {
		t.Fatalf("jsonEncode failed: %v", err)
	}

	if !strings.Contains(result.ToString(), "hello") {
		t.Errorf("Expected JSON to contain 'hello', got '%s'", result.ToString())
	}

	// Test array
	arr := types.NewEmptyArray()
	arr.Append(types.NewInt(1))
	arr.Append(types.NewInt(2))
	arr.Append(types.NewInt(3))

	result, err = ext.jsonEncode(testVM, []*types.Value{
		types.NewArray(arr),
	})

	if err != nil {
		t.Fatalf("jsonEncode array failed: %v", err)
	}

	jsonStr := result.ToString()
	if !strings.Contains(jsonStr, "1") || !strings.Contains(jsonStr, "2") {
		t.Errorf("Expected JSON array, got '%s'", jsonStr)
	}
}

func TestJSONDecode(t *testing.T) {
	ext := NewJSONExtension()
	testVM := vm.New()

	// Test simple JSON
	result, err := ext.jsonDecode(testVM, []*types.Value{
		types.NewString(`{"name":"John","age":30}`),
	})

	if err != nil {
		t.Fatalf("jsonDecode failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array result, got %s", result.TypeString())
	}

	arr := result.ToArray()
	nameVal, ok := arr.Get(types.NewString("name"))
	if !ok || nameVal.ToString() != "John" {
		t.Error("Failed to decode 'name' field correctly")
	}

	// Test invalid JSON
	result, err = ext.jsonDecode(testVM, []*types.Value{
		types.NewString(`{invalid json}`),
	})

	if err != nil {
		t.Fatalf("jsonDecode should not error: %v", err)
	}

	if result.Type() != types.TypeNull {
		t.Error("Expected null for invalid JSON")
	}
}

func TestJSONValidate(t *testing.T) {
	ext := NewJSONExtension()
	testVM := vm.New()

	// Valid JSON
	result, err := ext.jsonValidate(testVM, []*types.Value{
		types.NewString(`{"valid":true}`),
	})

	if err != nil {
		t.Fatalf("jsonValidate failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("Expected valid JSON to return true")
	}

	// Invalid JSON
	result, err = ext.jsonValidate(testVM, []*types.Value{
		types.NewString(`{invalid}`),
	})

	if err != nil {
		t.Fatalf("jsonValidate failed: %v", err)
	}

	if result.ToBool() {
		t.Error("Expected invalid JSON to return false")
	}
}

// Time Extension Tests

func TestTimeExtension(t *testing.T) {
	ext := NewTimeExtension()

	if ext.Name() != "go_time" {
		t.Errorf("Expected name 'go_time', got '%s'", ext.Name())
	}

	// Check constants
	constants := ext.Constants()
	if val, ok := constants["TIME_RFC3339"]; !ok || val.ToString() != time.RFC3339 {
		t.Error("TIME_RFC3339 constant not registered correctly")
	}
}

func TestTimeNow(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	result, err := ext.timeNow(testVM, []*types.Value{})

	if err != nil {
		t.Fatalf("timeNow failed: %v", err)
	}

	timestamp := result.ToInt()
	now := time.Now().Unix()

	// Should be within 1 second
	if timestamp < now-1 || timestamp > now+1 {
		t.Errorf("Timestamp %d not close to current time %d", timestamp, now)
	}
}

func TestTimeUnix(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	// Test known timestamp: 2023-01-01 00:00:00 UTC
	timestamp := int64(1672531200)

	result, err := ext.timeUnix(testVM, []*types.Value{
		types.NewInt(timestamp),
	})

	if err != nil {
		t.Fatalf("timeUnix failed: %v", err)
	}

	arr := result.ToArray()
	yearVal, _ := arr.Get(types.NewString("year"))
	if yearVal.ToInt() != 2023 {
		t.Errorf("Expected year 2023, got %d", yearVal.ToInt())
	}
}

func TestTimeFormat(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	timestamp := int64(1672531200) // 2023-01-01 00:00:00 UTC

	result, err := ext.timeFormat(testVM, []*types.Value{
		types.NewInt(timestamp),
		types.NewString(time.RFC3339),
	})

	if err != nil {
		t.Fatalf("timeFormat failed: %v", err)
	}

	formatted := result.ToString()
	if !strings.Contains(formatted, "2023") {
		t.Errorf("Expected formatted time to contain '2023', got '%s'", formatted)
	}
}

func TestTimeParse(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	result, err := ext.timeParse(testVM, []*types.Value{
		types.NewString(time.RFC3339),
		types.NewString("2023-01-01T00:00:00Z"),
	})

	if err != nil {
		t.Fatalf("timeParse failed: %v", err)
	}

	if result.ToInt() != 1672531200 {
		t.Errorf("Expected timestamp 1672531200, got %d", result.ToInt())
	}

	// Test invalid format
	result, err = ext.timeParse(testVM, []*types.Value{
		types.NewString(time.RFC3339),
		types.NewString("invalid"),
	})

	if err != nil {
		t.Fatalf("timeParse should not error: %v", err)
	}

	if result.Type() != types.TypeBool || result.ToBool() {
		t.Error("Expected false for invalid time string")
	}
}

func TestTimeAdd(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	timestamp := int64(1000)
	duration := int64(500)

	result, err := ext.timeAdd(testVM, []*types.Value{
		types.NewInt(timestamp),
		types.NewInt(duration),
	})

	if err != nil {
		t.Fatalf("timeAdd failed: %v", err)
	}

	if result.ToInt() != 1500 {
		t.Errorf("Expected 1500, got %d", result.ToInt())
	}
}

func TestTimeDiff(t *testing.T) {
	ext := NewTimeExtension()
	testVM := vm.New()

	result, err := ext.timeDiff(testVM, []*types.Value{
		types.NewInt(1000),
		types.NewInt(1500),
	})

	if err != nil {
		t.Fatalf("timeDiff failed: %v", err)
	}

	if result.ToInt() != 500 {
		t.Errorf("Expected 500, got %d", result.ToInt())
	}
}

// Filesystem Extension Tests

func TestFilesystemExtension(t *testing.T) {
	ext := NewFilesystemExtension()

	if ext.Name() != "go_fs" {
		t.Errorf("Expected name 'go_fs', got '%s'", ext.Name())
	}
}

func TestFileReadWrite(t *testing.T) {
	ext := NewFilesystemExtension()
	testVM := vm.New()

	// Create temp file
	tmpfile := "/tmp/php-go-test-file.txt"
	defer os.Remove(tmpfile)

	// Write file
	result, err := ext.fileWrite(testVM, []*types.Value{
		types.NewString(tmpfile),
		types.NewString("Hello, World!"),
	})

	if err != nil {
		t.Fatalf("fileWrite failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("fileWrite returned false")
	}

	// Read file
	result, err = ext.fileRead(testVM, []*types.Value{
		types.NewString(tmpfile),
	})

	if err != nil {
		t.Fatalf("fileRead failed: %v", err)
	}

	if result.ToString() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result.ToString())
	}
}

func TestFileExists(t *testing.T) {
	ext := NewFilesystemExtension()
	testVM := vm.New()

	// Non-existent file
	result, err := ext.fileExists(testVM, []*types.Value{
		types.NewString("/tmp/nonexistent-file-12345.txt"),
	})

	if err != nil {
		t.Fatalf("fileExists failed: %v", err)
	}

	if result.ToBool() {
		t.Error("Expected false for non-existent file")
	}

	// Create temp file
	tmpfile := "/tmp/php-go-test-exists.txt"
	os.WriteFile(tmpfile, []byte("test"), 0644)
	defer os.Remove(tmpfile)

	result, err = ext.fileExists(testVM, []*types.Value{
		types.NewString(tmpfile),
	})

	if err != nil {
		t.Fatalf("fileExists failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("Expected true for existing file")
	}
}

func TestFileCopy(t *testing.T) {
	ext := NewFilesystemExtension()
	testVM := vm.New()

	// Create source file
	src := "/tmp/php-go-test-src.txt"
	dst := "/tmp/php-go-test-dst.txt"
	os.WriteFile(src, []byte("copy test"), 0644)
	defer os.Remove(src)
	defer os.Remove(dst)

	result, err := ext.fileCopy(testVM, []*types.Value{
		types.NewString(src),
		types.NewString(dst),
	})

	if err != nil {
		t.Fatalf("fileCopy failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("fileCopy returned false")
	}

	// Verify copy
	content, _ := os.ReadFile(dst)
	if string(content) != "copy test" {
		t.Errorf("Expected 'copy test', got '%s'", string(content))
	}
}

func TestDirOperations(t *testing.T) {
	ext := NewFilesystemExtension()
	testVM := vm.New()

	tmpdir := "/tmp/php-go-test-dir"
	defer os.RemoveAll(tmpdir)

	// Create directory
	result, err := ext.dirCreate(testVM, []*types.Value{
		types.NewString(tmpdir),
	})

	if err != nil {
		t.Fatalf("dirCreate failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("dirCreate returned false")
	}

	// List directory
	result, err = ext.dirList(testVM, []*types.Value{
		types.NewString(tmpdir),
	})

	if err != nil {
		t.Fatalf("dirList failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Error("Expected array from dirList")
	}

	// Remove directory
	result, err = ext.dirRemove(testVM, []*types.Value{
		types.NewString(tmpdir),
	})

	if err != nil {
		t.Fatalf("dirRemove failed: %v", err)
	}

	if !result.ToBool() {
		t.Error("dirRemove returned false")
	}
}

// Crypto Extension Tests

func TestCryptoExtension(t *testing.T) {
	ext := NewCryptoExtension()

	if ext.Name() != "go_crypto" {
		t.Errorf("Expected name 'go_crypto', got '%s'", ext.Name())
	}
}

func TestHashFunctions(t *testing.T) {
	ext := NewCryptoExtension()
	testVM := vm.New()

	data := types.NewString("hello world")

	// Test MD5
	result, err := ext.hashMD5(testVM, []*types.Value{data})
	if err != nil {
		t.Fatalf("hashMD5 failed: %v", err)
	}
	if result.ToString() != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Errorf("Unexpected MD5 hash: %s", result.ToString())
	}

	// Test SHA1
	result, err = ext.hashSHA1(testVM, []*types.Value{data})
	if err != nil {
		t.Fatalf("hashSHA1 failed: %v", err)
	}
	if len(result.ToString()) != 40 {
		t.Error("SHA1 hash should be 40 characters")
	}

	// Test SHA256
	result, err = ext.hashSHA256(testVM, []*types.Value{data})
	if err != nil {
		t.Fatalf("hashSHA256 failed: %v", err)
	}
	if len(result.ToString()) != 64 {
		t.Error("SHA256 hash should be 64 characters")
	}

	// Test SHA512
	result, err = ext.hashSHA512(testVM, []*types.Value{data})
	if err != nil {
		t.Fatalf("hashSHA512 failed: %v", err)
	}
	if len(result.ToString()) != 128 {
		t.Error("SHA512 hash should be 128 characters")
	}
}

func TestBase64Encoding(t *testing.T) {
	ext := NewCryptoExtension()
	testVM := vm.New()

	original := "Hello, World!"

	// Encode
	encoded, err := ext.base64Encode(testVM, []*types.Value{
		types.NewString(original),
	})

	if err != nil {
		t.Fatalf("base64Encode failed: %v", err)
	}

	// Decode
	decoded, err := ext.base64Decode(testVM, []*types.Value{
		encoded,
	})

	if err != nil {
		t.Fatalf("base64Decode failed: %v", err)
	}

	if decoded.ToString() != original {
		t.Errorf("Expected '%s', got '%s'", original, decoded.ToString())
	}
}

func TestHexEncoding(t *testing.T) {
	ext := NewCryptoExtension()
	testVM := vm.New()

	original := "test"

	// Encode
	encoded, err := ext.hexEncode(testVM, []*types.Value{
		types.NewString(original),
	})

	if err != nil {
		t.Fatalf("hexEncode failed: %v", err)
	}

	if encoded.ToString() != "74657374" {
		t.Errorf("Expected '74657374', got '%s'", encoded.ToString())
	}

	// Decode
	decoded, err := ext.hexDecode(testVM, []*types.Value{
		encoded,
	})

	if err != nil {
		t.Fatalf("hexDecode failed: %v", err)
	}

	if decoded.ToString() != original {
		t.Errorf("Expected '%s', got '%s'", original, decoded.ToString())
	}
}

func TestAESEncryption(t *testing.T) {
	ext := NewCryptoExtension()
	testVM := vm.New()

	data := "Secret message"
	key := "12345678901234567890123456789012" // 32 bytes

	// Encrypt
	encrypted, err := ext.aesEncrypt(testVM, []*types.Value{
		types.NewString(data),
		types.NewString(key),
	})

	if err != nil {
		t.Fatalf("aesEncrypt failed: %v", err)
	}

	if encrypted.Type() == types.TypeBool {
		t.Fatal("Encryption failed")
	}

	// Decrypt
	decrypted, err := ext.aesDecrypt(testVM, []*types.Value{
		encrypted,
		types.NewString(key),
	})

	if err != nil {
		t.Fatalf("aesDecrypt failed: %v", err)
	}

	if decrypted.ToString() != data {
		t.Errorf("Expected '%s', got '%s'", data, decrypted.ToString())
	}
}

func TestRandomBytes(t *testing.T) {
	ext := NewCryptoExtension()
	testVM := vm.New()

	result, err := ext.randomBytes(testVM, []*types.Value{
		types.NewInt(16),
	})

	if err != nil {
		t.Fatalf("randomBytes failed: %v", err)
	}

	// Should return 32 hex characters (16 bytes * 2)
	if len(result.ToString()) != 32 {
		t.Errorf("Expected 32 characters, got %d", len(result.ToString()))
	}

	// Test two calls return different values
	result2, _ := ext.randomBytes(testVM, []*types.Value{
		types.NewInt(16),
	})

	if result.ToString() == result2.ToString() {
		t.Error("Random bytes should return different values")
	}
}
