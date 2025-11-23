package parallel

import (
	"sync"
	"sync/atomic"
)

// ============================================================================
// Copy-on-Write Optimization - Minimize copying overhead in parallel code
// ============================================================================

// COW (Copy-on-Write) allows multiple readers to share data without copying
// until a write occurs. This is crucial for efficient parallel execution.

// ============================================================================
// COW Array - Copy-on-Write array wrapper
// ============================================================================

// COWArray provides copy-on-write semantics for arrays
type COWArray struct {
	data     []interface{}
	refCount atomic.Int32
	mu       sync.RWMutex
}

// NewCOWArray creates a new COW array from existing data
func NewCOWArray(data []interface{}) *COWArray {
	arr := &COWArray{
		data: data,
	}
	arr.refCount.Store(1)
	return arr
}

// Clone creates a new reference to the same underlying data
func (ca *COWArray) Clone() *COWArray {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	clone := &COWArray{
		data: ca.data, // Share the same slice
	}
	clone.refCount.Store(1)
	ca.refCount.Add(1)

	return clone
}

// Get returns element at index (read-only, no copy)
func (ca *COWArray) Get(index int) interface{} {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if index < 0 || index >= len(ca.data) {
		return nil
	}

	return ca.data[index]
}

// Set sets element at index (triggers copy if shared)
func (ca *COWArray) Set(index int, value interface{}) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if index < 0 || index >= len(ca.data) {
		return
	}

	// Copy if shared (refCount > 1)
	if ca.refCount.Load() > 1 {
		ca.copyData()
	}

	ca.data[index] = value
}

// Append appends element (triggers copy if shared)
func (ca *COWArray) Append(value interface{}) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	// Copy if shared
	if ca.refCount.Load() > 1 {
		ca.copyData()
	}

	ca.data = append(ca.data, value)
}

// Len returns the length
func (ca *COWArray) Len() int {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	return len(ca.data)
}

// ToSlice returns a copy of the underlying slice
func (ca *COWArray) ToSlice() []interface{} {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	result := make([]interface{}, len(ca.data))
	copy(result, ca.data)
	return result
}

// IsShared returns true if the array is shared
func (ca *COWArray) IsShared() bool {
	return ca.refCount.Load() > 1
}

// RefCount returns the current reference count
func (ca *COWArray) RefCount() int32 {
	return ca.refCount.Load()
}

// copyData creates a copy of the underlying data
// Must be called with write lock held
func (ca *COWArray) copyData() {
	newData := make([]interface{}, len(ca.data))
	copy(newData, ca.data)
	ca.data = newData
	ca.refCount.Store(1)
}

// ============================================================================
// COW String - Copy-on-Write string wrapper
// ============================================================================

// COWString provides copy-on-write semantics for strings
// In Go, strings are immutable, but we can optimize buffer operations
type COWString struct {
	data     []byte // Underlying byte buffer
	refCount atomic.Int32
	mu       sync.RWMutex
}

// NewCOWString creates a new COW string
func NewCOWString(s string) *COWString {
	str := &COWString{
		data: []byte(s),
	}
	str.refCount.Store(1)
	return str
}

// NewCOWStringFromBytes creates a COW string from bytes
func NewCOWStringFromBytes(data []byte) *COWString {
	str := &COWString{
		data: data,
	}
	str.refCount.Store(1)
	return str
}

// Clone creates a new reference to the same underlying data
func (cs *COWString) Clone() *COWString {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	clone := &COWString{
		data: cs.data, // Share the same buffer
	}
	clone.refCount.Store(1)
	cs.refCount.Add(1)

	return clone
}

// String returns the string value
func (cs *COWString) String() string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return string(cs.data)
}

// Bytes returns a copy of the underlying bytes
func (cs *COWString) Bytes() []byte {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	result := make([]byte, len(cs.data))
	copy(result, cs.data)
	return result
}

// Append appends string (triggers copy if shared)
func (cs *COWString) Append(s string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Copy if shared
	if cs.refCount.Load() > 1 {
		cs.copyData()
	}

	cs.data = append(cs.data, []byte(s)...)
}

// Set replaces the string value (triggers copy if shared)
func (cs *COWString) Set(s string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Copy if shared
	if cs.refCount.Load() > 1 {
		cs.copyData()
	}

	cs.data = []byte(s)
}

// Len returns the length
func (cs *COWString) Len() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return len(cs.data)
}

// IsShared returns true if the string is shared
func (cs *COWString) IsShared() bool {
	return cs.refCount.Load() > 1
}

// RefCount returns the current reference count
func (cs *COWString) RefCount() int32 {
	return cs.refCount.Load()
}

// copyData creates a copy of the underlying data
// Must be called with write lock held
func (cs *COWString) copyData() {
	newData := make([]byte, len(cs.data))
	copy(newData, cs.data)
	cs.data = newData
	cs.refCount.Store(1)
}

// ============================================================================
// COW Map - Copy-on-Write map for object properties
// ============================================================================

// COWMap provides copy-on-write semantics for maps
type COWMap struct {
	data     map[string]interface{}
	refCount atomic.Int32
	mu       sync.RWMutex
}

// NewCOWMap creates a new COW map
func NewCOWMap() *COWMap {
	m := &COWMap{
		data: make(map[string]interface{}),
	}
	m.refCount.Store(1)
	return m
}

// NewCOWMapFromMap creates a COW map from existing map
func NewCOWMapFromMap(data map[string]interface{}) *COWMap {
	m := &COWMap{
		data: data,
	}
	m.refCount.Store(1)
	return m
}

// Clone creates a new reference to the same underlying data
func (cm *COWMap) Clone() *COWMap {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clone := &COWMap{
		data: cm.data, // Share the same map
	}
	clone.refCount.Store(1)
	cm.refCount.Add(1)

	return clone
}

// Get returns value for key (read-only, no copy)
func (cm *COWMap) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	value, exists := cm.data[key]
	return value, exists
}

// Set sets value for key (triggers copy if shared)
func (cm *COWMap) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Copy if shared
	if cm.refCount.Load() > 1 {
		cm.copyData()
	}

	cm.data[key] = value
}

// Delete deletes key (triggers copy if shared)
func (cm *COWMap) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Copy if shared
	if cm.refCount.Load() > 1 {
		cm.copyData()
	}

	delete(cm.data, key)
}

// Has checks if key exists
func (cm *COWMap) Has(key string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	_, exists := cm.data[key]
	return exists
}

// Len returns the number of keys
func (cm *COWMap) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.data)
}

// Keys returns all keys
func (cm *COWMap) Keys() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	keys := make([]string, 0, len(cm.data))
	for k := range cm.data {
		keys = append(keys, k)
	}
	return keys
}

// ToMap returns a copy of the underlying map
func (cm *COWMap) ToMap() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]interface{}, len(cm.data))
	for k, v := range cm.data {
		result[k] = v
	}
	return result
}

// IsShared returns true if the map is shared
func (cm *COWMap) IsShared() bool {
	return cm.refCount.Load() > 1
}

// RefCount returns the current reference count
func (cm *COWMap) RefCount() int32 {
	return cm.refCount.Load()
}

// copyData creates a copy of the underlying data
// Must be called with write lock held
func (cm *COWMap) copyData() {
	newData := make(map[string]interface{}, len(cm.data))
	for k, v := range cm.data {
		newData[k] = v
	}
	cm.data = newData
	cm.refCount.Store(1)
}

// ============================================================================
// COW Manager - Manages COW optimization metrics
// ============================================================================

// COWManager tracks copy-on-write optimization metrics
type COWManager struct {
	totalShares   atomic.Int64
	totalCopies   atomic.Int64
	bytesSaved    atomic.Int64
	bytesAlloced  atomic.Int64
}

// NewCOWManager creates a new COW manager
func NewCOWManager() *COWManager {
	return &COWManager{}
}

// RecordShare records a share operation
func (m *COWManager) RecordShare(bytesSaved int64) {
	m.totalShares.Add(1)
	m.bytesSaved.Add(bytesSaved)
}

// RecordCopy records a copy operation
func (m *COWManager) RecordCopy(bytesAlloced int64) {
	m.totalCopies.Add(1)
	m.bytesAlloced.Add(bytesAlloced)
}

// GetStats returns COW statistics
func (m *COWManager) GetStats() COWStats {
	return COWStats{
		TotalShares:  m.totalShares.Load(),
		TotalCopies:  m.totalCopies.Load(),
		BytesSaved:   m.bytesSaved.Load(),
		BytesAlloced: m.bytesAlloced.Load(),
	}
}

// Reset resets all metrics
func (m *COWManager) Reset() {
	m.totalShares.Store(0)
	m.totalCopies.Store(0)
	m.bytesSaved.Store(0)
	m.bytesAlloced.Store(0)
}

// COWStats contains COW statistics
type COWStats struct {
	TotalShares  int64
	TotalCopies  int64
	BytesSaved   int64
	BytesAlloced int64
}

// Efficiency returns the efficiency ratio (0.0 - 1.0)
func (s COWStats) Efficiency() float64 {
	total := s.BytesSaved + s.BytesAlloced
	if total == 0 {
		return 0.0
	}
	return float64(s.BytesSaved) / float64(total)
}

// ============================================================================
// Global COW Manager
// ============================================================================

var (
	globalCOWManager *COWManager
	cowOnce          sync.Once
)

// GetGlobalCOWManager returns the global COW manager
func GetGlobalCOWManager() *COWManager {
	cowOnce.Do(func() {
		globalCOWManager = NewCOWManager()
	})
	return globalCOWManager
}

// ============================================================================
// Helper Functions
// ============================================================================

// ShouldUseCOW determines if COW optimization should be used
func ShouldUseCOW(dataSize int, shareCount int) bool {
	// Use COW if:
	// 1. Data is large enough (> 1KB)
	// 2. Will be shared multiple times
	return dataSize > 1024 && shareCount > 1
}

// EstimateArraySize estimates memory size of an array
func EstimateArraySize(arr []interface{}) int {
	// Rough estimate: 8 bytes per pointer + overhead
	return len(arr) * 8
}

// EstimateMapSize estimates memory size of a map
func EstimateMapSize(m map[string]interface{}) int {
	// Rough estimate: key + pointer overhead
	size := 0
	for k := range m {
		size += len(k) + 8 // Key size + pointer
	}
	return size
}
