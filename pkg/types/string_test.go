package types

import (
	"testing"
)

func TestNewPhpString(t *testing.T) {
	s := NewPhpString("hello")
	if s.Val() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", s.Val())
	}
	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}
}

func TestNewPhpStringFromBytes(t *testing.T) {
	bytes := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, 0x57, 0x6F, 0x72, 0x6C, 0x64}
	s := NewPhpStringFromBytes(bytes)

	if s.Len() != 11 {
		t.Errorf("Expected length 11, got %d", s.Len())
	}

	if !s.ContainsNull() {
		t.Error("Expected string to contain null byte")
	}
}

func TestStringConcat(t *testing.T) {
	s1 := NewPhpString("Hello")
	s2 := NewPhpString(" World")
	result := s1.Concat(s2)

	if result.Val() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.Val())
	}

	// Test empty string concatenation
	empty := NewPhpString("")
	result2 := s1.Concat(empty)
	if result2.Val() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result2.Val())
	}
}

func TestStringConcatStr(t *testing.T) {
	s := NewPhpString("Hello")
	result := s.ConcatStr(" World")

	if result.Val() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result.Val())
	}
}

func TestSubstring(t *testing.T) {
	s := NewPhpString("Hello World")

	// Positive start, positive length
	sub := s.Substring(0, 5)
	if sub.Val() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", sub.Val())
	}

	// Positive start, positive length (middle)
	sub = s.Substring(6, 5)
	if sub.Val() != "World" {
		t.Errorf("Expected 'World', got '%s'", sub.Val())
	}

	// Negative start (count from end)
	sub = s.Substring(-5, 5)
	if sub.Val() != "World" {
		t.Errorf("Expected 'World', got '%s'", sub.Val())
	}

	// Negative length (all except last N)
	sub = s.Substring(0, -6)
	if sub.Val() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", sub.Val())
	}

	// Out of bounds
	sub = s.Substring(100, 5)
	if sub.Val() != "" {
		t.Errorf("Expected empty string, got '%s'", sub.Val())
	}
}

func TestSubstringFrom(t *testing.T) {
	s := NewPhpString("Hello World")
	sub := s.SubstringFrom(6)

	if sub.Val() != "World" {
		t.Errorf("Expected 'World', got '%s'", sub.Val())
	}
}

func TestCharAt(t *testing.T) {
	s := NewPhpString("Hello")

	if s.CharAt(0) != 'H' {
		t.Errorf("Expected 'H', got '%c'", s.CharAt(0))
	}

	if s.CharAt(4) != 'o' {
		t.Errorf("Expected 'o', got '%c'", s.CharAt(4))
	}

	// Out of bounds
	if s.CharAt(10) != 0 {
		t.Errorf("Expected 0 for out of bounds, got %d", s.CharAt(10))
	}
}

func TestIndexOf(t *testing.T) {
	s := NewPhpString("Hello World")

	if s.IndexOf("World") != 6 {
		t.Errorf("Expected index 6, got %d", s.IndexOf("World"))
	}

	if s.IndexOf("xyz") != -1 {
		t.Errorf("Expected -1 for not found, got %d", s.IndexOf("xyz"))
	}
}

func TestLastIndexOf(t *testing.T) {
	s := NewPhpString("Hello World World")

	if s.LastIndexOf("World") != 12 {
		t.Errorf("Expected index 12, got %d", s.LastIndexOf("World"))
	}

	if s.LastIndexOf("xyz") != -1 {
		t.Errorf("Expected -1 for not found, got %d", s.LastIndexOf("xyz"))
	}
}

func TestContains(t *testing.T) {
	s := NewPhpString("Hello World")

	if !s.Contains("World") {
		t.Error("Expected string to contain 'World'")
	}

	if s.Contains("xyz") {
		t.Error("Expected string to not contain 'xyz'")
	}
}

func TestToLower(t *testing.T) {
	s := NewPhpString("Hello World")
	lower := s.ToLower()

	if lower.Val() != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", lower.Val())
	}
}

func TestToUpper(t *testing.T) {
	s := NewPhpString("Hello World")
	upper := s.ToUpper()

	if upper.Val() != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD', got '%s'", upper.Val())
	}
}

func TestTrim(t *testing.T) {
	s := NewPhpString("  Hello World  ")
	trimmed := s.Trim()

	if trimmed.Val() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", trimmed.Val())
	}
}

func TestTrimLeft(t *testing.T) {
	s := NewPhpString("  Hello World  ")
	trimmed := s.TrimLeft()

	if trimmed.Val() != "Hello World  " {
		t.Errorf("Expected 'Hello World  ', got '%s'", trimmed.Val())
	}
}

func TestTrimRight(t *testing.T) {
	s := NewPhpString("  Hello World  ")
	trimmed := s.TrimRight()

	if trimmed.Val() != "  Hello World" {
		t.Errorf("Expected '  Hello World', got '%s'", trimmed.Val())
	}
}

func TestReplace(t *testing.T) {
	s := NewPhpString("Hello World")
	replaced := s.Replace("World", "Go")

	if replaced.Val() != "Hello Go" {
		t.Errorf("Expected 'Hello Go', got '%s'", replaced.Val())
	}

	// Multiple occurrences
	s2 := NewPhpString("foo bar foo")
	replaced2 := s2.Replace("foo", "baz")
	if replaced2.Val() != "baz bar baz" {
		t.Errorf("Expected 'baz bar baz', got '%s'", replaced2.Val())
	}
}

func TestSplit(t *testing.T) {
	s := NewPhpString("one,two,three")
	parts := s.Split(",")

	if len(parts) != 3 {
		t.Errorf("Expected 3 parts, got %d", len(parts))
	}

	if parts[0].Val() != "one" {
		t.Errorf("Expected 'one', got '%s'", parts[0].Val())
	}

	if parts[1].Val() != "two" {
		t.Errorf("Expected 'two', got '%s'", parts[1].Val())
	}

	if parts[2].Val() != "three" {
		t.Errorf("Expected 'three', got '%s'", parts[2].Val())
	}
}

func TestEquals(t *testing.T) {
	s1 := NewPhpString("Hello")
	s2 := NewPhpString("Hello")
	s3 := NewPhpString("World")

	if !s1.Equals(s2) {
		t.Error("Expected strings to be equal")
	}

	if s1.Equals(s3) {
		t.Error("Expected strings to not be equal")
	}
}

func TestCompare(t *testing.T) {
	s1 := NewPhpString("abc")
	s2 := NewPhpString("abc")
	s3 := NewPhpString("xyz")
	s4 := NewPhpString("aaa")

	if s1.Compare(s2) != 0 {
		t.Errorf("Expected 0, got %d", s1.Compare(s2))
	}

	if s1.Compare(s3) != -1 {
		t.Errorf("Expected -1, got %d", s1.Compare(s3))
	}

	if s1.Compare(s4) != 1 {
		t.Errorf("Expected 1, got %d", s1.Compare(s4))
	}
}

func TestHash(t *testing.T) {
	s := NewPhpString("Hello")

	hash1 := s.Hash()
	hash2 := s.Hash() // Should return cached hash

	if hash1 != hash2 {
		t.Error("Hash should be cached and consistent")
	}

	if hash1 == 0 {
		t.Error("Hash should not be zero for non-empty string")
	}

	// Different strings should have different hashes (usually)
	s2 := NewPhpString("World")
	hash3 := s2.Hash()
	if hash1 == hash3 {
		t.Error("Different strings should have different hashes")
	}
}

func TestBytes(t *testing.T) {
	s := NewPhpString("Hello")
	bytes := s.Bytes()

	expected := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}

	if len(bytes) != len(expected) {
		t.Errorf("Expected %d bytes, got %d", len(expected), len(bytes))
	}

	for i, b := range expected {
		if bytes[i] != b {
			t.Errorf("Byte %d: expected %x, got %x", i, b, bytes[i])
		}
	}
}

func TestContainsNull(t *testing.T) {
	s1 := NewPhpString("Hello")
	if s1.ContainsNull() {
		t.Error("String should not contain null bytes")
	}

	s2 := NewPhpString("Hello\x00World")
	if !s2.ContainsNull() {
		t.Error("String should contain null bytes")
	}
}

func TestIntern(t *testing.T) {
	s1 := Intern("common_string")
	s2 := Intern("common_string")

	// Should return the same instance
	if s1 != s2 {
		t.Error("Interned strings should be the same instance")
	}

	if !s1.IsInterned() {
		t.Error("String should be marked as interned")
	}

	// Different strings should be different instances
	s3 := Intern("different_string")
	if s1 == s3 {
		t.Error("Different interned strings should be different instances")
	}
}

func TestIsEmpty(t *testing.T) {
	s1 := NewPhpString("")
	if !s1.IsEmpty() {
		t.Error("Empty string should report as empty")
	}

	s2 := NewPhpString("Hello")
	if s2.IsEmpty() {
		t.Error("Non-empty string should not report as empty")
	}
}

func TestCopy(t *testing.T) {
	s := NewPhpString("Hello")
	copy := s.Copy()

	if copy.Val() != s.Val() {
		t.Error("Copy should have same value")
	}

	if copy.Len() != s.Len() {
		t.Error("Copy should have same length")
	}

	// Copy should not be interned even if original was
	s2 := Intern("test")
	copy2 := s2.Copy()
	if copy2.IsInterned() {
		t.Error("Copy of interned string should not be interned")
	}
}

func TestBinarySafety(t *testing.T) {
	// Test that strings can contain null bytes and other binary data
	bytes := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	s := NewPhpStringFromBytes(bytes)

	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}

	// Verify all bytes are preserved
	resultBytes := s.Bytes()
	for i, b := range bytes {
		if resultBytes[i] != b {
			t.Errorf("Byte %d: expected %x, got %x", i, b, resultBytes[i])
		}
	}
}
