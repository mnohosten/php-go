package types

import (
	"fmt"
	"sync"
)

// Resource represents a PHP resource (file handle, database connection, etc.)
// Resources are special variable types that hold references to external resources
type Resource struct {
	id         int
	typ        string
	data       interface{}
	destructor func(interface{})
	closed     bool
	mutex      sync.RWMutex
}

// Global resource tracking
var (
	nextResourceID   = 1
	resourceIDMutex  sync.Mutex
	activeResources  = make(map[int]*Resource)
	resourcesMutex   sync.RWMutex
	resourceRegistry = make(map[string]*ResourceType)
	registryMutex    sync.RWMutex
)

// ResourceType defines a type of resource that can be created
type ResourceType struct {
	Name       string
	Destructor func(interface{})
}

// ============================================================================
// Resource Type Registry
// ============================================================================

// RegisterResourceType registers a new resource type
func RegisterResourceType(name string, destructor func(interface{})) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	resourceRegistry[name] = &ResourceType{
		Name:       name,
		Destructor: destructor,
	}
}

// GetResourceType retrieves a registered resource type
func GetResourceType(name string) (*ResourceType, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	rt, exists := resourceRegistry[name]
	return rt, exists
}

// ============================================================================
// Resource Creation
// ============================================================================

// NewResourceHandle creates a new resource
func NewResourceHandle(resourceType string, data interface{}) *Resource {
	resourceIDMutex.Lock()
	id := nextResourceID
	nextResourceID++
	resourceIDMutex.Unlock()

	// Get destructor from registry if available
	var destructor func(interface{})
	if rt, exists := GetResourceType(resourceType); exists {
		destructor = rt.Destructor
	}

	resource := &Resource{
		id:         id,
		typ:        resourceType,
		data:       data,
		destructor: destructor,
		closed:     false,
	}

	// Track the resource
	resourcesMutex.Lock()
	activeResources[id] = resource
	resourcesMutex.Unlock()

	return resource
}

// NewResourceHandleWithDestructor creates a resource with a custom destructor
func NewResourceHandleWithDestructor(resourceType string, data interface{}, destructor func(interface{})) *Resource {
	resourceIDMutex.Lock()
	id := nextResourceID
	nextResourceID++
	resourceIDMutex.Unlock()

	resource := &Resource{
		id:         id,
		typ:        resourceType,
		data:       data,
		destructor: destructor,
		closed:     false,
	}

	resourcesMutex.Lock()
	activeResources[id] = resource
	resourcesMutex.Unlock()

	return resource
}

// ============================================================================
// Resource Properties
// ============================================================================

// ID returns the resource ID
func (r *Resource) ID() int {
	if r == nil {
		return 0
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.id
}

// Type returns the resource type
func (r *Resource) Type() string {
	if r == nil {
		return "Unknown"
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.typ
}

// Data returns the underlying data
func (r *Resource) Data() interface{} {
	if r == nil {
		return nil
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if r.closed {
		return nil
	}
	return r.data
}

// IsClosed returns true if the resource has been closed
func (r *Resource) IsClosed() bool {
	if r == nil {
		return true
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.closed
}

// IsValid returns true if the resource is valid and not closed
func (r *Resource) IsValid() bool {
	return r != nil && !r.IsClosed()
}

// ============================================================================
// Resource Management
// ============================================================================

// Close closes the resource and calls its destructor if set
func (r *Resource) Close() error {
	if r == nil {
		return fmt.Errorf("cannot close nil resource")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return fmt.Errorf("resource already closed")
	}

	// Call destructor if set
	if r.destructor != nil && r.data != nil {
		r.destructor(r.data)
	}

	r.closed = true
	r.data = nil

	// Remove from active resources
	resourcesMutex.Lock()
	delete(activeResources, r.id)
	resourcesMutex.Unlock()

	return nil
}

// SetData updates the resource data (must be valid/open)
func (r *Resource) SetData(data interface{}) error {
	if r == nil {
		return fmt.Errorf("cannot set data on nil resource")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return fmt.Errorf("resource is closed")
	}

	r.data = data
	return nil
}

// ============================================================================
// Global Resource Management
// ============================================================================

// GetResourceByID retrieves a resource by its ID
func GetResourceByID(id int) (*Resource, bool) {
	resourcesMutex.RLock()
	defer resourcesMutex.RUnlock()

	res, exists := activeResources[id]
	return res, exists
}

// GetActiveResourceCount returns the number of active resources
func GetActiveResourceCount() int {
	resourcesMutex.RLock()
	defer resourcesMutex.RUnlock()
	return len(activeResources)
}

// CloseAllResources closes all active resources
// This should be called on shutdown
func CloseAllResources() {
	resourcesMutex.Lock()
	resources := make([]*Resource, 0, len(activeResources))
	for _, res := range activeResources {
		resources = append(resources, res)
	}
	resourcesMutex.Unlock()

	for _, res := range resources {
		res.Close()
	}
}

// GetActiveResourceIDs returns a slice of all active resource IDs
func GetActiveResourceIDs() []int {
	resourcesMutex.RLock()
	defer resourcesMutex.RUnlock()

	ids := make([]int, 0, len(activeResources))
	for id := range activeResources {
		ids = append(ids, id)
	}
	return ids
}

// GetResourcesByType returns all resources of a specific type
func GetResourcesByType(resourceType string) []*Resource {
	resourcesMutex.RLock()
	defer resourcesMutex.RUnlock()

	results := make([]*Resource, 0)
	for _, res := range activeResources {
		if res.Type() == resourceType {
			results = append(results, res)
		}
	}
	return results
}

// ============================================================================
// String Representation
// ============================================================================

// String returns a string representation of the resource
func (r *Resource) String() string {
	if r == nil {
		return "Resource(nil)"
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.closed {
		return fmt.Sprintf("Resource(#%d, type=%s, closed)", r.id, r.typ)
	}

	return fmt.Sprintf("Resource(#%d, type=%s)", r.id, r.typ)
}

// ============================================================================
// Utility Functions
// ============================================================================

// ResetResourceSystem resets the resource system (for testing)
func ResetResourceSystem() {
	resourceIDMutex.Lock()
	nextResourceID = 1
	resourceIDMutex.Unlock()

	// Close all resources first
	CloseAllResources()

	resourcesMutex.Lock()
	activeResources = make(map[int]*Resource)
	resourcesMutex.Unlock()

	registryMutex.Lock()
	resourceRegistry = make(map[string]*ResourceType)
	registryMutex.Unlock()
}
