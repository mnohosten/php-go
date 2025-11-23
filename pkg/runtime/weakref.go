package runtime

import (
	"fmt"
	"sync"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// WeakReference (PHP 7.4+)
// ============================================================================

// WeakReference allows holding a reference to an object without preventing
// it from being garbage collected
type WeakReference struct {
	// The weakly-held object (may be nil if collected)
	object *types.Object

	// Unique ID of the referenced object for tracking
	objectID uint64

	// Whether the reference is still valid
	valid bool

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// NewWeakReference creates a new weak reference to an object
func NewWeakReference(obj *types.Object) *WeakReference {
	if obj == nil {
		return &WeakReference{
			object:   nil,
			objectID: 0,
			valid:    false,
		}
	}

	return &WeakReference{
		object:   obj,
		objectID: obj.ObjectID,
		valid:    true,
	}
}

// Get returns the weakly-referenced object, or nil if it has been collected
func (wr *WeakReference) Get() *types.Object {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	if !wr.valid {
		return nil
	}

	return wr.object
}

// IsValid checks if the weak reference is still valid
// Returns false if the object has been garbage collected
func (wr *WeakReference) IsValid() bool {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	return wr.valid
}

// Invalidate marks the reference as invalid (called by GC)
func (wr *WeakReference) Invalidate() {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.object = nil
	wr.valid = false
}

// GetObjectID returns the ID of the referenced object
func (wr *WeakReference) GetObjectID() uint64 {
	wr.mu.RLock()
	defer wr.mu.RUnlock()

	return wr.objectID
}

// ToValue converts the WeakReference to a PHP value (as a resource)
func (wr *WeakReference) ToValue() *types.Value {
	resource := types.NewResourceHandle("WeakReference", wr)
	return types.NewResource(resource)
}

// ============================================================================
// WeakMap (PHP 8.0+)
// ============================================================================

// WeakMap is a map that holds weak references to objects as keys
// When an object key is garbage collected, its entry is automatically removed
type WeakMap struct {
	// Map from object ID to value
	entries map[uint64]*weakMapEntry

	// Weak references to track object lifetime
	references map[uint64]*WeakReference

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// weakMapEntry represents an entry in the WeakMap
type weakMapEntry struct {
	key   *types.Object // The object key (weak)
	value *types.Value  // The associated value (strong)
}

// NewWeakMap creates a new WeakMap
func NewWeakMap() *WeakMap {
	return &WeakMap{
		entries:    make(map[uint64]*weakMapEntry),
		references: make(map[uint64]*WeakReference),
	}
}

// Set associates a value with an object key
func (wm *WeakMap) Set(key *types.Object, value *types.Value) error {
	if key == nil {
		return fmt.Errorf("WeakMap key must be an object")
	}

	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Create weak reference if not exists
	if _, exists := wm.references[key.ObjectID]; !exists {
		wm.references[key.ObjectID] = NewWeakReference(key)
	}

	// Store entry
	wm.entries[key.ObjectID] = &weakMapEntry{
		key:   key,
		value: value,
	}

	return nil
}

// Get retrieves the value associated with an object key
func (wm *WeakMap) Get(key *types.Object) (*types.Value, bool) {
	if key == nil {
		return nil, false
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()

	entry, exists := wm.entries[key.ObjectID]
	if !exists {
		return nil, false
	}

	// Check if reference is still valid
	ref, refExists := wm.references[key.ObjectID]
	if !refExists || !ref.IsValid() {
		return nil, false
	}

	return entry.value, true
}

// Has checks if a key exists in the WeakMap
func (wm *WeakMap) Has(key *types.Object) bool {
	if key == nil {
		return false
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()

	_, exists := wm.entries[key.ObjectID]
	if !exists {
		return false
	}

	// Check if reference is still valid
	ref, refExists := wm.references[key.ObjectID]
	if !refExists || !ref.IsValid() {
		return false
	}

	return true
}

// Delete removes an entry from the WeakMap
func (wm *WeakMap) Delete(key *types.Object) bool {
	if key == nil {
		return false
	}

	wm.mu.Lock()
	defer wm.mu.Unlock()

	_, exists := wm.entries[key.ObjectID]
	if !exists {
		return false
	}

	delete(wm.entries, key.ObjectID)
	delete(wm.references, key.ObjectID)

	return true
}

// Count returns the number of entries in the WeakMap
func (wm *WeakMap) Count() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	// Count only valid entries
	count := 0
	for objectID, ref := range wm.references {
		if ref.IsValid() {
			if _, exists := wm.entries[objectID]; exists {
				count++
			}
		}
	}

	return count
}

// Cleanup removes entries whose keys have been garbage collected
func (wm *WeakMap) Cleanup() int {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	removed := 0
	for objectID, ref := range wm.references {
		if !ref.IsValid() {
			delete(wm.entries, objectID)
			delete(wm.references, objectID)
			removed++
		}
	}

	return removed
}

// ToValue converts the WeakMap to a PHP value (as a resource)
func (wm *WeakMap) ToValue() *types.Value {
	resource := types.NewResourceHandle("WeakMap", wm)
	return types.NewResource(resource)
}

// ============================================================================
// WeakReference Registry
// ============================================================================

// WeakReferenceRegistry tracks all weak references for garbage collection
type WeakReferenceRegistry struct {
	// Map from object ID to list of weak references
	references map[uint64][]*WeakReference

	// Map from object ID to list of WeakMaps containing the object
	weakMaps map[uint64][]*WeakMap

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// NewWeakReferenceRegistry creates a new registry
func NewWeakReferenceRegistry() *WeakReferenceRegistry {
	return &WeakReferenceRegistry{
		references: make(map[uint64][]*WeakReference),
		weakMaps:   make(map[uint64][]*WeakMap),
	}
}

// RegisterWeakReference registers a weak reference for an object
func (wrr *WeakReferenceRegistry) RegisterWeakReference(objectID uint64, ref *WeakReference) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	wrr.references[objectID] = append(wrr.references[objectID], ref)
}

// RegisterWeakMap registers a WeakMap that contains an object
func (wrr *WeakReferenceRegistry) RegisterWeakMap(objectID uint64, wm *WeakMap) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	wrr.weakMaps[objectID] = append(wrr.weakMaps[objectID], wm)
}

// OnObjectDestroyed notifies the registry that an object has been destroyed
// This invalidates all weak references and removes WeakMap entries
func (wrr *WeakReferenceRegistry) OnObjectDestroyed(objectID uint64) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// Invalidate all weak references
	if refs, exists := wrr.references[objectID]; exists {
		for _, ref := range refs {
			ref.Invalidate()
		}
		delete(wrr.references, objectID)
	}

	// Cleanup WeakMaps
	if maps, exists := wrr.weakMaps[objectID]; exists {
		for _, wm := range maps {
			wm.mu.Lock()
			delete(wm.entries, objectID)
			if ref, ok := wm.references[objectID]; ok {
				ref.Invalidate()
				delete(wm.references, objectID)
			}
			wm.mu.Unlock()
		}
		delete(wrr.weakMaps, objectID)
	}
}

// Cleanup performs cleanup of all registered weak structures
func (wrr *WeakReferenceRegistry) Cleanup() (refsInvalidated int, entriesRemoved int) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// Cleanup invalid weak references
	for objectID, refs := range wrr.references {
		validRefs := make([]*WeakReference, 0)
		for _, ref := range refs {
			if ref.IsValid() {
				validRefs = append(validRefs, ref)
			} else {
				refsInvalidated++
			}
		}
		if len(validRefs) == 0 {
			delete(wrr.references, objectID)
		} else {
			wrr.references[objectID] = validRefs
		}
	}

	// Cleanup WeakMaps
	for objectID, maps := range wrr.weakMaps {
		validMaps := make([]*WeakMap, 0)
		for _, wm := range maps {
			removed := wm.Cleanup()
			entriesRemoved += removed
			// Keep the WeakMap in registry if it still has entries
			if wm.Count() > 0 {
				validMaps = append(validMaps, wm)
			}
		}
		if len(validMaps) == 0 {
			delete(wrr.weakMaps, objectID)
		} else {
			wrr.weakMaps[objectID] = validMaps
		}
	}

	return refsInvalidated, entriesRemoved
}

// GetStats returns statistics about weak references
func (wrr *WeakReferenceRegistry) GetStats() (numObjects int, numReferences int, numWeakMaps int) {
	wrr.mu.RLock()
	defer wrr.mu.RUnlock()

	numObjects = len(wrr.references)
	for _, refs := range wrr.references {
		numReferences += len(refs)
	}
	for _, maps := range wrr.weakMaps {
		numWeakMaps += len(maps)
	}

	return numObjects, numReferences, numWeakMaps
}

// ============================================================================
// Helper Functions
// ============================================================================

// CreateWeakReference is a helper to create a WeakReference from a PHP value
func CreateWeakReference(value *types.Value) (*WeakReference, error) {
	if value.Type() != types.TypeObject {
		return nil, fmt.Errorf("WeakReference can only reference objects, got %v", value.Type())
	}

	obj := value.ToObject()
	if obj == nil {
		return nil, fmt.Errorf("Invalid object value")
	}

	return NewWeakReference(obj), nil
}

// CreateWeakMap is a helper to create a new WeakMap
func CreateWeakMap() *WeakMap {
	return NewWeakMap()
}

// WeakReferenceFromValue extracts a WeakReference from a resource value
func WeakReferenceFromValue(value *types.Value) (*WeakReference, error) {
	if value.Type() != types.TypeResource {
		return nil, fmt.Errorf("Expected WeakReference resource, got %v", value.Type())
	}

	resource := value.ToResource()
	if resource == nil {
		return nil, fmt.Errorf("Invalid resource value")
	}

	wr, ok := resource.Data().(*WeakReference)
	if !ok {
		return nil, fmt.Errorf("Invalid WeakReference resource")
	}

	return wr, nil
}

// WeakMapFromValue extracts a WeakMap from a resource value
func WeakMapFromValue(value *types.Value) (*WeakMap, error) {
	if value.Type() != types.TypeResource {
		return nil, fmt.Errorf("Expected WeakMap resource, got %v", value.Type())
	}

	resource := value.ToResource()
	if resource == nil {
		return nil, fmt.Errorf("Invalid resource value")
	}

	wm, ok := resource.Data().(*WeakMap)
	if !ok {
		return nil, fmt.Errorf("Invalid WeakMap resource")
	}

	return wm, nil
}
