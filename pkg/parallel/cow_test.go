package parallel

import (
	"sync"
	"testing"
)

// ============================================================================
// COWArray Tests
// ============================================================================

func TestNewCOWArray(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr := NewCOWArray(data)

	if arr == nil {
		t.Error("NewCOWArray returned nil")
	}
	if arr.Len() != 3 {
		t.Errorf("Length = %d, expected 3", arr.Len())
	}
	if arr.RefCount() != 1 {
		t.Errorf("RefCount = %d, expected 1", arr.RefCount())
	}
}

func TestCOWArrayGet(t *testing.T) {
	data := []interface{}{10, 20, 30}
	arr := NewCOWArray(data)

	if arr.Get(0).(int) != 10 {
		t.Errorf("Get(0) = %v, expected 10", arr.Get(0))
	}
	if arr.Get(1).(int) != 20 {
		t.Errorf("Get(1) = %v, expected 20", arr.Get(1))
	}
	if arr.Get(2).(int) != 30 {
		t.Errorf("Get(2) = %v, expected 30", arr.Get(2))
	}
}

func TestCOWArrayGetOutOfBounds(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr := NewCOWArray(data)

	if arr.Get(-1) != nil {
		t.Error("Get(-1) should return nil")
	}
	if arr.Get(10) != nil {
		t.Error("Get(10) should return nil")
	}
}

func TestCOWArraySetNoShare(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr := NewCOWArray(data)

	arr.Set(1, 99)

	if arr.Get(1).(int) != 99 {
		t.Errorf("Get(1) = %v, expected 99", arr.Get(1))
	}
	if arr.RefCount() != 1 {
		t.Error("RefCount should remain 1")
	}
}

func TestCOWArrayClone(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr1 := NewCOWArray(data)
	arr2 := arr1.Clone()

	if arr2.Len() != 3 {
		t.Errorf("Clone length = %d, expected 3", arr2.Len())
	}
	if arr1.RefCount() != 2 {
		t.Errorf("Original RefCount = %d, expected 2", arr1.RefCount())
	}
	if arr2.RefCount() != 1 {
		t.Errorf("Clone RefCount = %d, expected 1", arr2.RefCount())
	}

	// Should share data initially
	if arr1.Get(0).(int) != arr2.Get(0).(int) {
		t.Error("Clone should share data")
	}
}

func TestCOWArrayCopyOnWrite(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr1 := NewCOWArray(data)
	arr2 := arr1.Clone()

	// Modify arr1 - should trigger copy
	arr1.Set(1, 99)

	if arr1.Get(1).(int) != 99 {
		t.Errorf("arr1.Get(1) = %v, expected 99", arr1.Get(1))
	}
	if arr2.Get(1).(int) != 2 {
		t.Errorf("arr2.Get(1) = %v, expected 2 (unchanged)", arr2.Get(1))
	}

	// arr1 should have copied
	if arr1.RefCount() != 1 {
		t.Errorf("arr1 RefCount = %d, expected 1 after copy", arr1.RefCount())
	}
}

func TestCOWArrayAppend(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr := NewCOWArray(data)

	arr.Append(4)

	if arr.Len() != 4 {
		t.Errorf("Length = %d, expected 4", arr.Len())
	}
	if arr.Get(3).(int) != 4 {
		t.Errorf("Get(3) = %v, expected 4", arr.Get(3))
	}
}

func TestCOWArrayAppendCOW(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr1 := NewCOWArray(data)
	arr2 := arr1.Clone()

	arr1.Append(4)

	if arr1.Len() != 4 {
		t.Errorf("arr1 length = %d, expected 4", arr1.Len())
	}
	if arr2.Len() != 3 {
		t.Errorf("arr2 length = %d, expected 3", arr2.Len())
	}
}

func TestCOWArrayToSlice(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr := NewCOWArray(data)

	slice := arr.ToSlice()
	if len(slice) != 3 {
		t.Errorf("Slice length = %d, expected 3", len(slice))
	}

	// Modifying slice should not affect array
	slice[0] = 99
	if arr.Get(0).(int) != 1 {
		t.Error("Modifying slice should not affect array")
	}
}

func TestCOWArrayIsShared(t *testing.T) {
	data := []interface{}{1, 2, 3}
	arr1 := NewCOWArray(data)

	if arr1.IsShared() {
		t.Error("Array should not be shared initially")
	}

	arr2 := arr1.Clone()
	if !arr1.IsShared() {
		t.Error("Array should be shared after clone")
	}

	arr1.Set(0, 99) // Triggers copy
	if arr1.IsShared() {
		t.Error("Array should not be shared after copy")
	}
	if arr2.IsShared() {
		t.Error("arr2 should not be shared")
	}
}

// ============================================================================
// COWString Tests
// ============================================================================

func TestNewCOWString(t *testing.T) {
	str := NewCOWString("hello")

	if str == nil {
		t.Error("NewCOWString returned nil")
	}
	if str.String() != "hello" {
		t.Errorf("String() = %s, expected 'hello'", str.String())
	}
	if str.RefCount() != 1 {
		t.Errorf("RefCount = %d, expected 1", str.RefCount())
	}
}

func TestNewCOWStringFromBytes(t *testing.T) {
	data := []byte("world")
	str := NewCOWStringFromBytes(data)

	if str.String() != "world" {
		t.Errorf("String() = %s, expected 'world'", str.String())
	}
}

func TestCOWStringClone(t *testing.T) {
	str1 := NewCOWString("hello")
	str2 := str1.Clone()

	if str2.String() != "hello" {
		t.Errorf("Clone String() = %s, expected 'hello'", str2.String())
	}
	if str1.RefCount() != 2 {
		t.Errorf("Original RefCount = %d, expected 2", str1.RefCount())
	}
}

func TestCOWStringAppend(t *testing.T) {
	str := NewCOWString("hello")
	str.Append(" world")

	if str.String() != "hello world" {
		t.Errorf("String() = %s, expected 'hello world'", str.String())
	}
}

func TestCOWStringAppendCOW(t *testing.T) {
	str1 := NewCOWString("hello")
	str2 := str1.Clone()

	str1.Append(" world")

	if str1.String() != "hello world" {
		t.Errorf("str1 = %s, expected 'hello world'", str1.String())
	}
	if str2.String() != "hello" {
		t.Errorf("str2 = %s, expected 'hello'", str2.String())
	}
	if str1.RefCount() != 1 {
		t.Error("str1 should have copied")
	}
}

func TestCOWStringSet(t *testing.T) {
	str := NewCOWString("hello")
	str.Set("goodbye")

	if str.String() != "goodbye" {
		t.Errorf("String() = %s, expected 'goodbye'", str.String())
	}
}

func TestCOWStringSetCOW(t *testing.T) {
	str1 := NewCOWString("hello")
	str2 := str1.Clone()

	str1.Set("goodbye")

	if str1.String() != "goodbye" {
		t.Errorf("str1 = %s, expected 'goodbye'", str1.String())
	}
	if str2.String() != "hello" {
		t.Errorf("str2 = %s, expected 'hello'", str2.String())
	}
}

func TestCOWStringLen(t *testing.T) {
	str := NewCOWString("hello")

	if str.Len() != 5 {
		t.Errorf("Len() = %d, expected 5", str.Len())
	}
}

func TestCOWStringBytes(t *testing.T) {
	str := NewCOWString("hello")
	bytes := str.Bytes()

	if string(bytes) != "hello" {
		t.Errorf("Bytes() = %s, expected 'hello'", string(bytes))
	}

	// Modifying bytes should not affect string
	bytes[0] = 'X'
	if str.String() != "hello" {
		t.Error("Modifying bytes should not affect string")
	}
}

func TestCOWStringIsShared(t *testing.T) {
	str1 := NewCOWString("hello")

	if str1.IsShared() {
		t.Error("String should not be shared initially")
	}

	_ = str1.Clone()
	if !str1.IsShared() {
		t.Error("String should be shared after clone")
	}

	str1.Set("goodbye") // Triggers copy
	if str1.IsShared() {
		t.Error("String should not be shared after copy")
	}
}

// ============================================================================
// COWMap Tests
// ============================================================================

func TestNewCOWMap(t *testing.T) {
	m := NewCOWMap()

	if m == nil {
		t.Error("NewCOWMap returned nil")
	}
	if m.Len() != 0 {
		t.Errorf("Len() = %d, expected 0", m.Len())
	}
	if m.RefCount() != 1 {
		t.Errorf("RefCount = %d, expected 1", m.RefCount())
	}
}

func TestNewCOWMapFromMap(t *testing.T) {
	data := map[string]interface{}{
		"a": 1,
		"b": 2,
	}
	m := NewCOWMapFromMap(data)

	if m.Len() != 2 {
		t.Errorf("Len() = %d, expected 2", m.Len())
	}
}

func TestCOWMapGetSet(t *testing.T) {
	m := NewCOWMap()
	m.Set("key1", "value1")

	val, exists := m.Get("key1")
	if !exists {
		t.Error("Key should exist")
	}
	if val.(string) != "value1" {
		t.Errorf("Get('key1') = %v, expected 'value1'", val)
	}
}

func TestCOWMapClone(t *testing.T) {
	m1 := NewCOWMap()
	m1.Set("key1", "value1")

	m2 := m1.Clone()

	val, exists := m2.Get("key1")
	if !exists {
		t.Error("Cloned map should have key")
	}
	if val.(string) != "value1" {
		t.Error("Cloned map should have same value")
	}

	if m1.RefCount() != 2 {
		t.Errorf("Original RefCount = %d, expected 2", m1.RefCount())
	}
}

func TestCOWMapSetCOW(t *testing.T) {
	m1 := NewCOWMap()
	m1.Set("key1", "value1")

	m2 := m1.Clone()

	m1.Set("key1", "modified")

	val1, _ := m1.Get("key1")
	if val1.(string) != "modified" {
		t.Errorf("m1 value = %v, expected 'modified'", val1)
	}

	val2, _ := m2.Get("key1")
	if val2.(string) != "value1" {
		t.Errorf("m2 value = %v, expected 'value1'", val2)
	}

	if m1.RefCount() != 1 {
		t.Error("m1 should have copied")
	}
}

func TestCOWMapDelete(t *testing.T) {
	m := NewCOWMap()
	m.Set("key1", "value1")
	m.Set("key2", "value2")

	m.Delete("key1")

	if m.Has("key1") {
		t.Error("key1 should be deleted")
	}
	if !m.Has("key2") {
		t.Error("key2 should still exist")
	}
}

func TestCOWMapDeleteCOW(t *testing.T) {
	m1 := NewCOWMap()
	m1.Set("key1", "value1")

	m2 := m1.Clone()

	m1.Delete("key1")

	if m1.Has("key1") {
		t.Error("m1 should not have key1")
	}
	if !m2.Has("key1") {
		t.Error("m2 should still have key1")
	}
}

func TestCOWMapHas(t *testing.T) {
	m := NewCOWMap()
	m.Set("key1", "value1")

	if !m.Has("key1") {
		t.Error("Should have key1")
	}
	if m.Has("key2") {
		t.Error("Should not have key2")
	}
}

func TestCOWMapKeys(t *testing.T) {
	m := NewCOWMap()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys length = %d, expected 3", len(keys))
	}

	// Check all keys present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	if !keyMap["a"] || !keyMap["b"] || !keyMap["c"] {
		t.Error("Not all keys returned")
	}
}

func TestCOWMapToMap(t *testing.T) {
	m := NewCOWMap()
	m.Set("a", 1)
	m.Set("b", 2)

	data := m.ToMap()
	if len(data) != 2 {
		t.Errorf("Map length = %d, expected 2", len(data))
	}

	// Modifying returned map should not affect COWMap
	data["c"] = 3
	if m.Has("c") {
		t.Error("Modifying returned map should not affect COWMap")
	}
}

func TestCOWMapIsShared(t *testing.T) {
	m1 := NewCOWMap()
	m1.Set("key", "value")

	if m1.IsShared() {
		t.Error("Map should not be shared initially")
	}

	_ = m1.Clone()
	if !m1.IsShared() {
		t.Error("Map should be shared after clone")
	}

	m1.Set("key", "modified") // Triggers copy
	if m1.IsShared() {
		t.Error("Map should not be shared after copy")
	}
}

// ============================================================================
// COWManager Tests
// ============================================================================

func TestNewCOWManager(t *testing.T) {
	mgr := NewCOWManager()

	if mgr == nil {
		t.Error("NewCOWManager returned nil")
	}

	stats := mgr.GetStats()
	if stats.TotalShares != 0 {
		t.Error("TotalShares should be 0")
	}
	if stats.TotalCopies != 0 {
		t.Error("TotalCopies should be 0")
	}
}

func TestCOWManagerRecordShare(t *testing.T) {
	mgr := NewCOWManager()

	mgr.RecordShare(1024)
	mgr.RecordShare(2048)

	stats := mgr.GetStats()
	if stats.TotalShares != 2 {
		t.Errorf("TotalShares = %d, expected 2", stats.TotalShares)
	}
	if stats.BytesSaved != 3072 {
		t.Errorf("BytesSaved = %d, expected 3072", stats.BytesSaved)
	}
}

func TestCOWManagerRecordCopy(t *testing.T) {
	mgr := NewCOWManager()

	mgr.RecordCopy(512)
	mgr.RecordCopy(1024)

	stats := mgr.GetStats()
	if stats.TotalCopies != 2 {
		t.Errorf("TotalCopies = %d, expected 2", stats.TotalCopies)
	}
	if stats.BytesAlloced != 1536 {
		t.Errorf("BytesAlloced = %d, expected 1536", stats.BytesAlloced)
	}
}

func TestCOWManagerReset(t *testing.T) {
	mgr := NewCOWManager()

	mgr.RecordShare(1024)
	mgr.RecordCopy(512)

	mgr.Reset()

	stats := mgr.GetStats()
	if stats.TotalShares != 0 || stats.TotalCopies != 0 {
		t.Error("Stats should be reset")
	}
}

func TestCOWStatsEfficiency(t *testing.T) {
	stats := COWStats{
		BytesSaved:   8000,
		BytesAlloced: 2000,
	}

	efficiency := stats.Efficiency()
	expected := 0.8
	if efficiency != expected {
		t.Errorf("Efficiency = %f, expected %f", efficiency, expected)
	}
}

func TestCOWStatsEfficiencyZero(t *testing.T) {
	stats := COWStats{}

	efficiency := stats.Efficiency()
	if efficiency != 0.0 {
		t.Errorf("Efficiency = %f, expected 0.0", efficiency)
	}
}

func TestGetGlobalCOWManager(t *testing.T) {
	mgr1 := GetGlobalCOWManager()
	if mgr1 == nil {
		t.Error("GetGlobalCOWManager returned nil")
	}

	// Should return same instance
	mgr2 := GetGlobalCOWManager()
	if mgr1 != mgr2 {
		t.Error("Should return same instance")
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestShouldUseCOW(t *testing.T) {
	tests := []struct {
		dataSize   int
		shareCount int
		expected   bool
	}{
		{100, 2, false},    // Too small
		{2000, 1, false},   // Not shared
		{2000, 2, true},    // Large and shared
		{1024, 2, false},   // Exactly at threshold
		{1025, 2, true},    // Just above threshold
	}

	for _, tt := range tests {
		result := ShouldUseCOW(tt.dataSize, tt.shareCount)
		if result != tt.expected {
			t.Errorf("ShouldUseCOW(%d, %d) = %v, expected %v",
				tt.dataSize, tt.shareCount, result, tt.expected)
		}
	}
}

func TestEstimateArraySize(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	size := EstimateArraySize(arr)

	expected := 5 * 8 // 5 elements * 8 bytes
	if size != expected {
		t.Errorf("EstimateArraySize = %d, expected %d", size, expected)
	}
}

func TestEstimateMapSize(t *testing.T) {
	m := map[string]interface{}{
		"ab":  1,
		"xyz": 2,
	}
	size := EstimateMapSize(m)

	// "ab" (2) + "xyz" (3) = 5 bytes for keys + 2*8 for pointers = 21
	expected := 5 + 16
	if size != expected {
		t.Errorf("EstimateMapSize = %d, expected %d", size, expected)
	}
}

// ============================================================================
// Concurrent Tests
// ============================================================================

func TestCOWArrayConcurrent(t *testing.T) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	arr := NewCOWArray(data)
	var wg sync.WaitGroup

	// Multiple readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				arr.Get(j)
			}
		}()
	}

	// One writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			arr.Set(j, j*2)
		}
	}()

	wg.Wait()
}

func TestCOWMapConcurrent(t *testing.T) {
	m := NewCOWMap()

	var wg sync.WaitGroup

	// Multiple writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.Set(string(rune('a'+id)), j)
			}
		}(i)
	}

	// Multiple readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				m.Get("a")
			}
		}()
	}

	wg.Wait()
}
