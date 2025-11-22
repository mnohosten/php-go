package types

import (
	"testing"
)

// ============================================================================
// Resource Creation Tests
// ============================================================================

func TestNewResourceHandle(t *testing.T) {
	ResetResourceSystem() // Clean state

	res := NewResourceHandle("file", "test.txt")

	if res == nil {
		t.Fatal("NewResourceHandle returned nil")
	}

	if res.ID() != 1 {
		t.Errorf("Expected resource ID 1, got %d", res.ID())
	}

	if res.Type() != "file" {
		t.Errorf("Expected type 'file', got '%s'", res.Type())
	}

	if res.Data() != "test.txt" {
		t.Errorf("Expected data 'test.txt', got '%v'", res.Data())
	}

	if res.IsClosed() {
		t.Error("New resource should not be closed")
	}

	if !res.IsValid() {
		t.Error("New resource should be valid")
	}
}

func TestNewResourceHandleWithDestructor(t *testing.T) {
	ResetResourceSystem()

	destructorCalled := false
	destructor := func(data interface{}) {
		destructorCalled = true
	}

	res := NewResourceHandleWithDestructor("test", "data", destructor)

	if res == nil {
		t.Fatal("NewResourceHandleWithDestructor returned nil")
	}

	res.Close()

	if !destructorCalled {
		t.Error("Destructor should have been called on Close")
	}
}

func TestResourceIDIncrement(t *testing.T) {
	ResetResourceSystem()

	res1 := NewResourceHandle("type1", nil)
	res2 := NewResourceHandle("type2", nil)
	res3 := NewResourceHandle("type3", nil)

	if res1.ID() != 1 {
		t.Errorf("Expected res1 ID 1, got %d", res1.ID())
	}

	if res2.ID() != 2 {
		t.Errorf("Expected res2 ID 2, got %d", res2.ID())
	}

	if res3.ID() != 3 {
		t.Errorf("Expected res3 ID 3, got %d", res3.ID())
	}
}

// ============================================================================
// Resource Type Registry Tests
// ============================================================================

func TestRegisterResourceType(t *testing.T) {
	ResetResourceSystem()

	called := false
	destructor := func(data interface{}) {
		called = true
	}

	RegisterResourceType("mytype", destructor)

	rt, exists := GetResourceType("mytype")
	if !exists {
		t.Fatal("Resource type should be registered")
	}

	if rt.Name != "mytype" {
		t.Errorf("Expected name 'mytype', got '%s'", rt.Name)
	}

	// Test that registered destructor is used
	res := NewResourceHandle("mytype", "data")
	res.Close()

	if !called {
		t.Error("Registered destructor should be called")
	}
}

func TestGetResourceTypeNotFound(t *testing.T) {
	ResetResourceSystem()

	_, exists := GetResourceType("nonexistent")
	if exists {
		t.Error("Non-existent resource type should not be found")
	}
}

// ============================================================================
// Resource Properties Tests
// ============================================================================

func TestResourceID(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", nil)
	id := res.ID()

	if id == 0 {
		t.Error("Resource ID should not be 0")
	}

	// Nil resource
	var nilRes *Resource
	if nilRes.ID() != 0 {
		t.Error("Nil resource ID should be 0")
	}
}

func TestResourceType(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("database", nil)

	if res.Type() != "database" {
		t.Errorf("Expected type 'database', got '%s'", res.Type())
	}

	// Nil resource
	var nilRes *Resource
	if nilRes.Type() != "Unknown" {
		t.Error("Nil resource type should be 'Unknown'")
	}
}

func TestResourceData(t *testing.T) {
	ResetResourceSystem()

	type FileHandle struct {
		Name string
		FD   int
	}

	handle := &FileHandle{Name: "test.txt", FD: 42}
	res := NewResourceHandle("file", handle)

	data := res.Data()
	if data == nil {
		t.Fatal("Resource data should not be nil")
	}

	fh, ok := data.(*FileHandle)
	if !ok {
		t.Fatal("Data should be *FileHandle")
	}

	if fh.Name != "test.txt" {
		t.Errorf("Expected name 'test.txt', got '%s'", fh.Name)
	}
}

func TestResourceDataAfterClose(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	res.Close()

	if res.Data() != nil {
		t.Error("Closed resource data should be nil")
	}
}

func TestResourceIsClosed(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", nil)

	if res.IsClosed() {
		t.Error("New resource should not be closed")
	}

	res.Close()

	if !res.IsClosed() {
		t.Error("Closed resource should report as closed")
	}

	// Nil resource
	var nilRes *Resource
	if !nilRes.IsClosed() {
		t.Error("Nil resource should report as closed")
	}
}

func TestResourceIsValid(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", nil)

	if !res.IsValid() {
		t.Error("New resource should be valid")
	}

	res.Close()

	if res.IsValid() {
		t.Error("Closed resource should not be valid")
	}

	// Nil resource
	var nilRes *Resource
	if nilRes.IsValid() {
		t.Error("Nil resource should not be valid")
	}
}

// ============================================================================
// Resource Management Tests
// ============================================================================

func TestResourceClose(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")

	err := res.Close()
	if err != nil {
		t.Errorf("Close should not return error: %v", err)
	}

	if !res.IsClosed() {
		t.Error("Resource should be closed")
	}

	// Closing again should return error
	err = res.Close()
	if err == nil {
		t.Error("Closing already closed resource should return error")
	}
}

func TestResourceCloseNil(t *testing.T) {
	var nilRes *Resource
	err := nilRes.Close()
	if err == nil {
		t.Error("Closing nil resource should return error")
	}
}

func TestResourceCloseCallsDestructor(t *testing.T) {
	ResetResourceSystem()

	destructorCalled := false
	var destructorData interface{}

	destructor := func(data interface{}) {
		destructorCalled = true
		destructorData = data
	}

	res := NewResourceHandleWithDestructor("test", "mydata", destructor)
	res.Close()

	if !destructorCalled {
		t.Error("Destructor should be called on Close")
	}

	if destructorData != "mydata" {
		t.Error("Destructor should receive resource data")
	}
}

func TestResourceSetData(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "initial")

	err := res.SetData("updated")
	if err != nil {
		t.Errorf("SetData should not return error: %v", err)
	}

	if res.Data() != "updated" {
		t.Error("Data should be updated")
	}
}

func TestResourceSetDataClosed(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	res.Close()

	err := res.SetData("new")
	if err == nil {
		t.Error("SetData on closed resource should return error")
	}
}

func TestResourceSetDataNil(t *testing.T) {
	var nilRes *Resource
	err := nilRes.SetData("data")
	if err == nil {
		t.Error("SetData on nil resource should return error")
	}
}

// ============================================================================
// Global Resource Management Tests
// ============================================================================

func TestGetResourceByID(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	id := res.ID()

	retrieved, exists := GetResourceByID(id)
	if !exists {
		t.Error("Resource should exist")
	}

	if retrieved != res {
		t.Error("Retrieved resource should be the same instance")
	}
}

func TestGetResourceByIDNotFound(t *testing.T) {
	ResetResourceSystem()

	_, exists := GetResourceByID(999)
	if exists {
		t.Error("Non-existent resource should not be found")
	}
}

func TestGetResourceByIDAfterClose(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	id := res.ID()

	res.Close()

	_, exists := GetResourceByID(id)
	if exists {
		t.Error("Closed resource should not be in active resources")
	}
}

func TestGetActiveResourceCount(t *testing.T) {
	ResetResourceSystem()

	if GetActiveResourceCount() != 0 {
		t.Error("Initial count should be 0")
	}

	res1 := NewResourceHandle("file", nil)
	if GetActiveResourceCount() != 1 {
		t.Errorf("Expected count 1, got %d", GetActiveResourceCount())
	}

	res2 := NewResourceHandle("db", nil)
	if GetActiveResourceCount() != 2 {
		t.Errorf("Expected count 2, got %d", GetActiveResourceCount())
	}

	res1.Close()
	if GetActiveResourceCount() != 1 {
		t.Errorf("Expected count 1 after close, got %d", GetActiveResourceCount())
	}

	res2.Close()
	if GetActiveResourceCount() != 0 {
		t.Errorf("Expected count 0 after all closed, got %d", GetActiveResourceCount())
	}
}

func TestCloseAllResources(t *testing.T) {
	ResetResourceSystem()

	res1 := NewResourceHandle("file", nil)
	res2 := NewResourceHandle("db", nil)
	res3 := NewResourceHandle("socket", nil)

	if GetActiveResourceCount() != 3 {
		t.Error("Expected 3 active resources")
	}

	CloseAllResources()

	if GetActiveResourceCount() != 0 {
		t.Error("All resources should be closed")
	}

	if !res1.IsClosed() || !res2.IsClosed() || !res3.IsClosed() {
		t.Error("All resources should be marked as closed")
	}
}

func TestGetActiveResourceIDs(t *testing.T) {
	ResetResourceSystem()

	res1 := NewResourceHandle("file", nil)
	res2 := NewResourceHandle("db", nil)

	ids := GetActiveResourceIDs()

	if len(ids) != 2 {
		t.Errorf("Expected 2 IDs, got %d", len(ids))
	}

	// Check IDs are present
	foundRes1 := false
	foundRes2 := false
	for _, id := range ids {
		if id == res1.ID() {
			foundRes1 = true
		}
		if id == res2.ID() {
			foundRes2 = true
		}
	}

	if !foundRes1 || !foundRes2 {
		t.Error("Active resource IDs should include all created resources")
	}
}

func TestGetResourcesByType(t *testing.T) {
	ResetResourceSystem()

	file1 := NewResourceHandle("file", "file1.txt")
	file2 := NewResourceHandle("file", "file2.txt")
	db := NewResourceHandle("database", "mydb")

	fileResources := GetResourcesByType("file")
	if len(fileResources) != 2 {
		t.Errorf("Expected 2 file resources, got %d", len(fileResources))
	}

	dbResources := GetResourcesByType("database")
	if len(dbResources) != 1 {
		t.Errorf("Expected 1 database resource, got %d", len(dbResources))
	}

	// Check that we got the right resources
	foundFile1 := false
	foundFile2 := false
	for _, res := range fileResources {
		if res == file1 {
			foundFile1 = true
		}
		if res == file2 {
			foundFile2 = true
		}
	}

	if !foundFile1 || !foundFile2 {
		t.Error("File resources should contain both file handles")
	}

	if dbResources[0] != db {
		t.Error("Database resources should contain the db handle")
	}

	// Non-existent type
	nonexistent := GetResourcesByType("nonexistent")
	if len(nonexistent) != 0 {
		t.Error("Non-existent type should return empty slice")
	}
}

// ============================================================================
// String Representation Tests
// ============================================================================

func TestResourceString(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	str := res.String()

	if str == "" {
		t.Error("String representation should not be empty")
	}

	// Should contain type
	if len(str) < 4 {
		t.Errorf("String representation seems too short: '%s'", str)
	}
}

func TestResourceStringClosed(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", "data")
	res.Close()

	str := res.String()
	// Closed resources should indicate they're closed in string representation
	if len(str) == 0 {
		t.Error("Closed resource string should not be empty")
	}
}

func TestResourceStringNil(t *testing.T) {
	var nilRes *Resource
	str := nilRes.String()

	if str != "Resource(nil)" {
		t.Errorf("Expected 'Resource(nil)', got '%s'", str)
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestResourceConcurrentAccess(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("file", 0)

	done := make(chan bool)

	// Multiple goroutines accessing resource
	for i := 0; i < 10; i++ {
		go func(val int) {
			res.SetData(val)
			_ = res.Data()
			_ = res.IsClosed()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestResourceConcurrentCreate(t *testing.T) {
	ResetResourceSystem()

	done := make(chan *Resource)

	// Create multiple resources concurrently
	for i := 0; i < 10; i++ {
		go func() {
			res := NewResourceHandle("test", nil)
			done <- res
		}()
	}

	// Collect resources
	resources := make([]*Resource, 10)
	for i := 0; i < 10; i++ {
		resources[i] = <-done
	}

	// All should have unique IDs
	ids := make(map[int]bool)
	for _, res := range resources {
		id := res.ID()
		if ids[id] {
			t.Errorf("Duplicate resource ID: %d", id)
		}
		ids[id] = true
	}
}

// ============================================================================
// Reset Tests
// ============================================================================

func TestResetResourceSystem(t *testing.T) {
	ResetResourceSystem()

	// Create some resources
	res1 := NewResourceHandle("file", nil)
	res2 := NewResourceHandle("db", nil)
	RegisterResourceType("mytype", nil)

	// Reset
	ResetResourceSystem()

	// Check everything is reset
	if GetActiveResourceCount() != 0 {
		t.Error("Active resource count should be 0 after reset")
	}

	if res1.IsValid() || res2.IsValid() {
		t.Error("Resources should be closed after reset")
	}

	if _, exists := GetResourceType("mytype"); exists {
		t.Error("Resource type registry should be cleared")
	}

	// Next resource should have ID 1
	newRes := NewResourceHandle("test", nil)
	if newRes.ID() != 1 {
		t.Errorf("First resource after reset should have ID 1, got %d", newRes.ID())
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestResourceWithNilData(t *testing.T) {
	ResetResourceSystem()

	res := NewResourceHandle("test", nil)

	if res.Data() != nil {
		t.Error("Resource with nil data should return nil")
	}
}

func TestResourceWithComplexData(t *testing.T) {
	ResetResourceSystem()

	type ComplexData struct {
		Field1 string
		Field2 int
		Field3 []string
	}

	data := &ComplexData{
		Field1: "test",
		Field2: 42,
		Field3: []string{"a", "b", "c"},
	}

	res := NewResourceHandle("complex", data)

	retrieved := res.Data().(*ComplexData)
	if retrieved.Field1 != "test" || retrieved.Field2 != 42 {
		t.Error("Complex data not preserved correctly")
	}
}
