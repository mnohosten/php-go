package runtime

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// Helper function to create a test object
func createTestObject(id uint64) *types.Object {
	return &types.Object{
		ClassName:  "TestClass",
		ObjectID:   id,
		Properties: make(map[string]*types.Property),
	}
}

// TestWeakReference_Creation tests weak reference creation
func TestWeakReference_Creation(t *testing.T) {
	obj := createTestObject(1)
	wr := NewWeakReference(obj)

	if wr == nil {
		t.Fatal("NewWeakReference returned nil")
	}

	if !wr.IsValid() {
		t.Error("New weak reference should be valid")
	}

	if wr.GetObjectID() != 1 {
		t.Errorf("GetObjectID() = %d, want 1", wr.GetObjectID())
	}
}

// TestWeakReference_CreationWithNil tests creating weak reference with nil
func TestWeakReference_CreationWithNil(t *testing.T) {
	wr := NewWeakReference(nil)

	if wr == nil {
		t.Fatal("NewWeakReference returned nil")
	}

	if wr.IsValid() {
		t.Error("Weak reference to nil should be invalid")
	}

	if wr.Get() != nil {
		t.Error("Get() should return nil for invalid reference")
	}
}

// TestWeakReference_Get tests retrieving the referenced object
func TestWeakReference_Get(t *testing.T) {
	obj := createTestObject(1)
	wr := NewWeakReference(obj)

	retrieved := wr.Get()
	if retrieved == nil {
		t.Fatal("Get() returned nil")
	}

	if retrieved.ObjectID != obj.ObjectID {
		t.Errorf("Retrieved object ID = %d, want %d", retrieved.ObjectID, obj.ObjectID)
	}

	if retrieved.ClassName != "TestClass" {
		t.Errorf("Retrieved object class = %s, want TestClass", retrieved.ClassName)
	}
}

// TestWeakReference_Invalidate tests invalidating a weak reference
func TestWeakReference_Invalidate(t *testing.T) {
	obj := createTestObject(1)
	wr := NewWeakReference(obj)

	if !wr.IsValid() {
		t.Error("Weak reference should be valid initially")
	}

	// Simulate garbage collection
	wr.Invalidate()

	if wr.IsValid() {
		t.Error("Weak reference should be invalid after Invalidate()")
	}

	if wr.Get() != nil {
		t.Error("Get() should return nil after invalidation")
	}
}

// TestWeakReference_ToValue tests converting to PHP value
func TestWeakReference_ToValue(t *testing.T) {
	obj := createTestObject(1)
	wr := NewWeakReference(obj)

	value := wr.ToValue()
	if value.Type() != types.TypeResource {
		t.Errorf("ToValue() type = %v, want TypeResource", value.Type())
	}
}

// TestWeakMap_Creation tests WeakMap creation
func TestWeakMap_Creation(t *testing.T) {
	wm := NewWeakMap()

	if wm == nil {
		t.Fatal("NewWeakMap returned nil")
	}

	if wm.Count() != 0 {
		t.Errorf("New WeakMap count = %d, want 0", wm.Count())
	}
}

// TestWeakMap_SetAndGet tests setting and getting values
func TestWeakMap_SetAndGet(t *testing.T) {
	wm := NewWeakMap()
	obj := createTestObject(1)
	value := types.NewString("test value")

	// Set value
	if err := wm.Set(obj, value); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	// Get value
	retrieved, exists := wm.Get(obj)
	if !exists {
		t.Fatal("Get() returned false for existing key")
	}

	if retrieved.ToString() != "test value" {
		t.Errorf("Get() value = %s, want test value", retrieved.ToString())
	}
}

// TestWeakMap_SetWithNilKey tests setting with nil key
func TestWeakMap_SetWithNilKey(t *testing.T) {
	wm := NewWeakMap()
	value := types.NewString("test")

	err := wm.Set(nil, value)
	if err == nil {
		t.Error("Set() with nil key should return error")
	}
}

// TestWeakMap_Has tests checking for key existence
func TestWeakMap_Has(t *testing.T) {
	wm := NewWeakMap()
	obj1 := createTestObject(1)
	obj2 := createTestObject(2)

	wm.Set(obj1, types.NewString("value"))

	if !wm.Has(obj1) {
		t.Error("Has() = false for existing key, want true")
	}

	if wm.Has(obj2) {
		t.Error("Has() = true for non-existing key, want false")
	}

	if wm.Has(nil) {
		t.Error("Has() = true for nil key, want false")
	}
}

// TestWeakMap_Delete tests deleting entries
func TestWeakMap_Delete(t *testing.T) {
	wm := NewWeakMap()
	obj := createTestObject(1)

	wm.Set(obj, types.NewString("value"))

	if !wm.Has(obj) {
		t.Fatal("Key should exist before deletion")
	}

	deleted := wm.Delete(obj)
	if !deleted {
		t.Error("Delete() = false, want true")
	}

	if wm.Has(obj) {
		t.Error("Key should not exist after deletion")
	}

	// Delete non-existent key
	deleted = wm.Delete(obj)
	if deleted {
		t.Error("Delete() = true for non-existent key, want false")
	}
}

// TestWeakMap_Count tests counting entries
func TestWeakMap_Count(t *testing.T) {
	wm := NewWeakMap()

	if wm.Count() != 0 {
		t.Errorf("Initial count = %d, want 0", wm.Count())
	}

	obj1 := createTestObject(1)
	obj2 := createTestObject(2)
	obj3 := createTestObject(3)

	wm.Set(obj1, types.NewString("value1"))
	wm.Set(obj2, types.NewString("value2"))
	wm.Set(obj3, types.NewString("value3"))

	if wm.Count() != 3 {
		t.Errorf("Count after 3 sets = %d, want 3", wm.Count())
	}

	wm.Delete(obj2)

	if wm.Count() != 2 {
		t.Errorf("Count after deletion = %d, want 2", wm.Count())
	}
}

// TestWeakMap_MultipleValues tests storing multiple values
func TestWeakMap_MultipleValues(t *testing.T) {
	wm := NewWeakMap()

	objects := []*types.Object{
		createTestObject(1),
		createTestObject(2),
		createTestObject(3),
	}

	values := []*types.Value{
		types.NewString("value1"),
		types.NewInt(42),
		types.NewBool(true),
	}

	// Set all values
	for i, obj := range objects {
		if err := wm.Set(obj, values[i]); err != nil {
			t.Fatalf("Set(%d) error: %v", i, err)
		}
	}

	// Verify all values
	for i, obj := range objects {
		retrieved, exists := wm.Get(obj)
		if !exists {
			t.Errorf("Get(%d) returned false", i)
			continue
		}

		if retrieved.Type() != values[i].Type() {
			t.Errorf("Get(%d) type = %v, want %v", i, retrieved.Type(), values[i].Type())
		}
	}
}

// TestWeakMap_Cleanup tests cleanup of invalid entries
func TestWeakMap_Cleanup(t *testing.T) {
	wm := NewWeakMap()

	obj1 := createTestObject(1)
	obj2 := createTestObject(2)

	wm.Set(obj1, types.NewString("value1"))
	wm.Set(obj2, types.NewString("value2"))

	// Invalidate one reference
	if ref, exists := wm.references[obj1.ObjectID]; exists {
		ref.Invalidate()
	}

	removed := wm.Cleanup()
	if removed != 1 {
		t.Errorf("Cleanup() removed = %d, want 1", removed)
	}

	if wm.Count() != 1 {
		t.Errorf("Count after cleanup = %d, want 1", wm.Count())
	}
}

// TestWeakMap_ToValue tests converting to PHP value
func TestWeakMap_ToValue(t *testing.T) {
	wm := NewWeakMap()

	value := wm.ToValue()
	if value.Type() != types.TypeResource {
		t.Errorf("ToValue() type = %v, want TypeResource", value.Type())
	}
}

// TestWeakReferenceRegistry_Creation tests registry creation
func TestWeakReferenceRegistry_Creation(t *testing.T) {
	registry := NewWeakReferenceRegistry()

	if registry == nil {
		t.Fatal("NewWeakReferenceRegistry returned nil")
	}

	numObjects, numRefs, numMaps := registry.GetStats()
	if numObjects != 0 || numRefs != 0 || numMaps != 0 {
		t.Error("New registry should have zero stats")
	}
}

// TestWeakReferenceRegistry_RegisterWeakReference tests registering weak references
func TestWeakReferenceRegistry_RegisterWeakReference(t *testing.T) {
	registry := NewWeakReferenceRegistry()
	obj := createTestObject(1)
	wr := NewWeakReference(obj)

	registry.RegisterWeakReference(obj.ObjectID, wr)

	numObjects, numRefs, _ := registry.GetStats()
	if numObjects != 1 {
		t.Errorf("numObjects = %d, want 1", numObjects)
	}
	if numRefs != 1 {
		t.Errorf("numRefs = %d, want 1", numRefs)
	}
}

// TestWeakReferenceRegistry_RegisterWeakMap tests registering weak maps
func TestWeakReferenceRegistry_RegisterWeakMap(t *testing.T) {
	registry := NewWeakReferenceRegistry()
	obj := createTestObject(1)
	wm := NewWeakMap()

	registry.RegisterWeakMap(obj.ObjectID, wm)

	_, _, numMaps := registry.GetStats()
	if numMaps != 1 {
		t.Errorf("numMaps = %d, want 1", numMaps)
	}
}

// TestWeakReferenceRegistry_OnObjectDestroyed tests object destruction notification
func TestWeakReferenceRegistry_OnObjectDestroyed(t *testing.T) {
	registry := NewWeakReferenceRegistry()
	obj := createTestObject(1)

	// Create weak reference and register
	wr := NewWeakReference(obj)
	registry.RegisterWeakReference(obj.ObjectID, wr)

	if !wr.IsValid() {
		t.Fatal("Weak reference should be valid initially")
	}

	// Notify object destroyed
	registry.OnObjectDestroyed(obj.ObjectID)

	if wr.IsValid() {
		t.Error("Weak reference should be invalid after object destroyed")
	}

	numObjects, numRefs, _ := registry.GetStats()
	if numObjects != 0 {
		t.Errorf("numObjects after destruction = %d, want 0", numObjects)
	}
	if numRefs != 0 {
		t.Errorf("numRefs after destruction = %d, want 0", numRefs)
	}
}

// TestWeakReferenceRegistry_OnObjectDestroyedWithWeakMap tests WeakMap cleanup on destruction
func TestWeakReferenceRegistry_OnObjectDestroyedWithWeakMap(t *testing.T) {
	registry := NewWeakReferenceRegistry()
	obj := createTestObject(1)
	wm := NewWeakMap()

	// Add entry to WeakMap
	wm.Set(obj, types.NewString("value"))
	registry.RegisterWeakMap(obj.ObjectID, wm)

	if !wm.Has(obj) {
		t.Fatal("WeakMap should contain object initially")
	}

	// Notify object destroyed
	registry.OnObjectDestroyed(obj.ObjectID)

	if wm.Has(obj) {
		t.Error("WeakMap should not contain object after destruction")
	}
}

// TestWeakReferenceRegistry_Cleanup tests registry cleanup
func TestWeakReferenceRegistry_Cleanup(t *testing.T) {
	registry := NewWeakReferenceRegistry()

	obj1 := createTestObject(1)
	obj2 := createTestObject(2)

	wr1 := NewWeakReference(obj1)
	wr2 := NewWeakReference(obj2)

	registry.RegisterWeakReference(obj1.ObjectID, wr1)
	registry.RegisterWeakReference(obj2.ObjectID, wr2)

	// Invalidate one reference manually
	wr1.Invalidate()

	refsInvalidated, _ := registry.Cleanup()
	if refsInvalidated != 1 {
		t.Errorf("Cleanup() refs invalidated = %d, want 1", refsInvalidated)
	}

	numObjects, numRefs, _ := registry.GetStats()
	if numObjects != 1 {
		t.Errorf("numObjects after cleanup = %d, want 1", numObjects)
	}
	if numRefs != 1 {
		t.Errorf("numRefs after cleanup = %d, want 1", numRefs)
	}
}

// TestWeakReferenceRegistry_MultipleReferences tests multiple references to same object
func TestWeakReferenceRegistry_MultipleReferences(t *testing.T) {
	registry := NewWeakReferenceRegistry()
	obj := createTestObject(1)

	wr1 := NewWeakReference(obj)
	wr2 := NewWeakReference(obj)
	wr3 := NewWeakReference(obj)

	registry.RegisterWeakReference(obj.ObjectID, wr1)
	registry.RegisterWeakReference(obj.ObjectID, wr2)
	registry.RegisterWeakReference(obj.ObjectID, wr3)

	numObjects, numRefs, _ := registry.GetStats()
	if numObjects != 1 {
		t.Errorf("numObjects = %d, want 1", numObjects)
	}
	if numRefs != 3 {
		t.Errorf("numRefs = %d, want 3", numRefs)
	}

	// Destroy object - all references should be invalidated
	registry.OnObjectDestroyed(obj.ObjectID)

	if wr1.IsValid() || wr2.IsValid() || wr3.IsValid() {
		t.Error("All weak references should be invalid after object destroyed")
	}
}

// TestCreateWeakReference tests helper function
func TestCreateWeakReference(t *testing.T) {
	obj := createTestObject(1)
	value := types.NewObject(obj)

	wr, err := CreateWeakReference(value)
	if err != nil {
		t.Fatalf("CreateWeakReference() error: %v", err)
	}

	if wr == nil {
		t.Fatal("CreateWeakReference returned nil")
	}

	if !wr.IsValid() {
		t.Error("Created weak reference should be valid")
	}

	// Test with non-object value
	_, err = CreateWeakReference(types.NewString("not an object"))
	if err == nil {
		t.Error("CreateWeakReference() with non-object should return error")
	}
}

// TestCreateWeakMap tests helper function
func TestCreateWeakMap(t *testing.T) {
	wm := CreateWeakMap()

	if wm == nil {
		t.Fatal("CreateWeakMap returned nil")
	}

	if wm.Count() != 0 {
		t.Error("Created WeakMap should be empty")
	}
}

// TestWeakReferenceFromValue tests extracting WeakReference from value
func TestWeakReferenceFromValue(t *testing.T) {
	obj := createTestObject(1)
	wr := NewWeakReference(obj)
	value := wr.ToValue()

	extracted, err := WeakReferenceFromValue(value)
	if err != nil {
		t.Fatalf("WeakReferenceFromValue() error: %v", err)
	}

	if extracted != wr {
		t.Error("Extracted WeakReference should be the same instance")
	}

	// Test with non-resource value
	_, err = WeakReferenceFromValue(types.NewString("not a resource"))
	if err == nil {
		t.Error("WeakReferenceFromValue() with non-resource should return error")
	}
}

// TestWeakMapFromValue tests extracting WeakMap from value
func TestWeakMapFromValue(t *testing.T) {
	wm := NewWeakMap()
	value := wm.ToValue()

	extracted, err := WeakMapFromValue(value)
	if err != nil {
		t.Fatalf("WeakMapFromValue() error: %v", err)
	}

	if extracted != wm {
		t.Error("Extracted WeakMap should be the same instance")
	}

	// Test with non-resource value
	_, err = WeakMapFromValue(types.NewInt(42))
	if err == nil {
		t.Error("WeakMapFromValue() with non-resource should return error")
	}
}

// TestWeakMap_OverwriteValue tests overwriting a value for existing key
func TestWeakMap_OverwriteValue(t *testing.T) {
	wm := NewWeakMap()
	obj := createTestObject(1)

	// Set initial value
	wm.Set(obj, types.NewString("first"))

	// Overwrite with new value
	wm.Set(obj, types.NewString("second"))

	retrieved, exists := wm.Get(obj)
	if !exists {
		t.Fatal("Key should exist")
	}

	if retrieved.ToString() != "second" {
		t.Errorf("Get() = %s, want second", retrieved.ToString())
	}

	// Count should still be 1
	if wm.Count() != 1 {
		t.Errorf("Count() = %d, want 1", wm.Count())
	}
}

// TestWeakReferenceRegistry_ComplexScenario tests a complex real-world scenario
func TestWeakReferenceRegistry_ComplexScenario(t *testing.T) {
	registry := NewWeakReferenceRegistry()

	// Create multiple objects
	obj1 := createTestObject(1)
	obj2 := createTestObject(2)
	obj3 := createTestObject(3)

	// Create weak references
	wr1 := NewWeakReference(obj1)
	wr2a := NewWeakReference(obj2)
	wr2b := NewWeakReference(obj2) // Multiple refs to same object
	wr3 := NewWeakReference(obj3)

	registry.RegisterWeakReference(obj1.ObjectID, wr1)
	registry.RegisterWeakReference(obj2.ObjectID, wr2a)
	registry.RegisterWeakReference(obj2.ObjectID, wr2b)
	registry.RegisterWeakReference(obj3.ObjectID, wr3)

	// Create WeakMaps
	wm1 := NewWeakMap()
	wm2 := NewWeakMap()

	wm1.Set(obj1, types.NewString("value1"))
	wm1.Set(obj2, types.NewString("value2"))
	wm2.Set(obj2, types.NewInt(42))
	wm2.Set(obj3, types.NewBool(true))

	registry.RegisterWeakMap(obj1.ObjectID, wm1)
	registry.RegisterWeakMap(obj2.ObjectID, wm1)
	registry.RegisterWeakMap(obj2.ObjectID, wm2)
	registry.RegisterWeakMap(obj3.ObjectID, wm2)

	// Verify initial state
	numObjects, numRefs, numMaps := registry.GetStats()
	if numObjects != 3 {
		t.Errorf("Initial numObjects = %d, want 3", numObjects)
	}
	if numRefs != 4 {
		t.Errorf("Initial numRefs = %d, want 4", numRefs)
	}
	if numMaps != 4 {
		t.Errorf("Initial numMaps = %d, want 4", numMaps)
	}

	// Destroy obj2
	registry.OnObjectDestroyed(obj2.ObjectID)

	// Verify obj2 references are invalid
	if wr2a.IsValid() || wr2b.IsValid() {
		t.Error("obj2 weak references should be invalid")
	}

	// Verify obj2 removed from WeakMaps
	if wm1.Has(obj2) {
		t.Error("wm1 should not contain obj2")
	}
	if wm2.Has(obj2) {
		t.Error("wm2 should not contain obj2")
	}

	// Verify obj1 and obj3 still valid
	if !wr1.IsValid() || !wr3.IsValid() {
		t.Error("obj1 and obj3 weak references should still be valid")
	}
	if !wm1.Has(obj1) || !wm2.Has(obj3) {
		t.Error("obj1 and obj3 should still be in WeakMaps")
	}
}
