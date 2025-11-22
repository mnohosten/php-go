package types

// Object represents a PHP object instance
// This will be fully implemented in Phase 5.
// For now, we provide a minimal placeholder.
type Object struct {
	ClassName  string
	Properties map[string]*Value
}

// NewObject creates a new object
func NewObjectInstance(className string) *Object {
	return &Object{
		ClassName:  className,
		Properties: make(map[string]*Value),
	}
}
