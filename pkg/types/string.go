package types

import (
	"hash/fnv"
	"strings"
)

// String represents a PHP string value
// PHP strings are binary-safe and can contain null bytes
type String struct {
	val    string // The actual string data (Go strings are already UTF-8 and binary-safe)
	len    int    // Cached length in bytes
	hash   uint64 // Cached hash value (0 means not computed)
	interned bool // Whether this string is interned (deduplicated)
}

// NewPhpString creates a new PHP string
func NewPhpString(s string) *String {
	return &String{
		val:  s,
		len:  len(s),
		hash: 0, // Lazy hash computation
	}
}

// NewPhpStringFromBytes creates a new PHP string from bytes
// This is truly binary-safe
func NewPhpStringFromBytes(b []byte) *String {
	return &String{
		val:  string(b),
		len:  len(b),
		hash: 0,
	}
}

// ============================================================================
// String Properties
// ============================================================================

// Val returns the string value
func (s *String) Val() string {
	return s.val
}

// Len returns the length in bytes (not characters)
// This matches PHP's strlen() which returns byte length
func (s *String) Len() int {
	return s.len
}

// IsEmpty returns true if the string is empty
func (s *String) IsEmpty() bool {
	return s.len == 0
}

// ============================================================================
// String Operations
// ============================================================================

// Concat concatenates this string with another
// Returns a new String (strings are immutable)
func (s *String) Concat(other *String) *String {
	if s.len == 0 {
		return other
	}
	if other.len == 0 {
		return s
	}
	return NewPhpString(s.val + other.val)
}

// ConcatStr concatenates this string with a Go string
func (s *String) ConcatStr(str string) *String {
	if s.len == 0 {
		return NewPhpString(str)
	}
	if len(str) == 0 {
		return s
	}
	return NewPhpString(s.val + str)
}

// Substring returns a substring from start with given length
// Negative start counts from the end
// Negative length means "all except last N characters"
// This matches PHP's substr() behavior
func (s *String) Substring(start int, length int) *String {
	// Handle negative start (count from end)
	if start < 0 {
		start = s.len + start
		if start < 0 {
			start = 0
		}
	}

	// Start beyond string length
	if start >= s.len {
		return NewPhpString("")
	}

	// Handle negative length (all except last N)
	if length < 0 {
		length = s.len - start + length
		if length < 0 {
			return NewPhpString("")
		}
	}

	// Calculate end position
	end := start + length
	if end > s.len {
		end = s.len
	}

	return NewPhpString(s.val[start:end])
}

// SubstringFrom returns substring from start to end of string
func (s *String) SubstringFrom(start int) *String {
	return s.Substring(start, s.len-start)
}

// CharAt returns the byte at the given index
// Returns empty string if index out of bounds
func (s *String) CharAt(index int) byte {
	if index < 0 || index >= s.len {
		return 0
	}
	return s.val[index]
}

// IndexOf finds the first occurrence of substring
// Returns -1 if not found
func (s *String) IndexOf(substr string) int {
	return strings.Index(s.val, substr)
}

// LastIndexOf finds the last occurrence of substring
// Returns -1 if not found
func (s *String) LastIndexOf(substr string) int {
	return strings.LastIndex(s.val, substr)
}

// Contains checks if string contains substring
func (s *String) Contains(substr string) bool {
	return strings.Contains(s.val, substr)
}

// ToLower returns lowercase version
func (s *String) ToLower() *String {
	return NewPhpString(strings.ToLower(s.val))
}

// ToUpper returns uppercase version
func (s *String) ToUpper() *String {
	return NewPhpString(strings.ToUpper(s.val))
}

// Trim removes whitespace from both ends
func (s *String) Trim() *String {
	return NewPhpString(strings.TrimSpace(s.val))
}

// TrimLeft removes whitespace from start
func (s *String) TrimLeft() *String {
	return NewPhpString(strings.TrimLeft(s.val, " \t\n\r\x00\x0B"))
}

// TrimRight removes whitespace from end
func (s *String) TrimRight() *String {
	return NewPhpString(strings.TrimRight(s.val, " \t\n\r\x00\x0B"))
}

// Replace replaces all occurrences of old with new
func (s *String) Replace(old, new string) *String {
	return NewPhpString(strings.ReplaceAll(s.val, old, new))
}

// Split splits string by delimiter
func (s *String) Split(delim string) []*String {
	parts := strings.Split(s.val, delim)
	result := make([]*String, len(parts))
	for i, part := range parts {
		result[i] = NewPhpString(part)
	}
	return result
}

// ============================================================================
// Comparison
// ============================================================================

// Equals checks string equality
func (s *String) Equals(other *String) bool {
	if s.len != other.len {
		return false
	}
	return s.val == other.val
}

// Compare compares two strings lexicographically
// Returns: -1 if s < other, 0 if s == other, 1 if s > other
func (s *String) Compare(other *String) int {
	if s.val < other.val {
		return -1
	}
	if s.val > other.val {
		return 1
	}
	return 0
}

// ============================================================================
// Hashing
// ============================================================================

// Hash returns the hash value of the string
// Computes hash on first call, then caches it
func (s *String) Hash() uint64 {
	if s.hash == 0 && s.len > 0 {
		s.hash = s.computeHash()
	}
	return s.hash
}

// computeHash computes FNV-1a hash of the string
func (s *String) computeHash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.val))
	return h.Sum64()
}

// ============================================================================
// Binary Safety
// ============================================================================

// Bytes returns the raw bytes of the string
// This allows binary-safe operations
func (s *String) Bytes() []byte {
	return []byte(s.val)
}

// ContainsNull checks if string contains null bytes
func (s *String) ContainsNull() bool {
	return strings.Contains(s.val, "\x00")
}

// ============================================================================
// String Interning (for optimization)
// ============================================================================

// stringInternMap is a global map for string interning
// This is used to deduplicate common strings like property names
var stringInternMap = make(map[string]*String)

// Intern returns an interned version of the string
// Multiple calls with the same value will return the same String instance
// This saves memory for frequently used strings
func Intern(s string) *String {
	if interned, exists := stringInternMap[s]; exists {
		return interned
	}

	str := NewPhpString(s)
	str.interned = true
	str.hash = str.computeHash() // Compute hash eagerly for interned strings
	stringInternMap[s] = str
	return str
}

// IsInterned returns true if this string is interned
func (s *String) IsInterned() bool {
	return s.interned
}

// ============================================================================
// Conversion
// ============================================================================

// String implements fmt.Stringer
func (s *String) String() string {
	return s.val
}

// Copy creates a shallow copy of the string
// Since strings are immutable, this just returns the same instance
func (s *String) Copy() *String {
	// Strings are immutable, so we can return the same instance
	// But for consistency with the API, we create a new wrapper
	return &String{
		val:      s.val,
		len:      s.len,
		hash:     s.hash,
		interned: false, // Copies are not interned
	}
}
