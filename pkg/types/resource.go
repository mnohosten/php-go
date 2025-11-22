package types

// Resource represents a PHP resource (file handle, database connection, etc.)
// This will be fully implemented in Phase 4 and Phase 6.
// For now, we provide a minimal placeholder.
type Resource struct {
	ID   int
	Type string
	Data interface{}
}

var nextResourceID = 1

// NewResourceHandle creates a new resource
func NewResourceHandle(resourceType string, data interface{}) *Resource {
	id := nextResourceID
	nextResourceID++

	return &Resource{
		ID:   id,
		Type: resourceType,
		Data: data,
	}
}
