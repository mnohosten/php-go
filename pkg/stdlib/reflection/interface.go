package reflection

import (
	"github.com/krizos/php-go/pkg/types"
)

// ReflectionInterface represents an interface reflection
type ReflectionInterface struct {
	// Interface name
	name string

	// Interface entry
	iface *types.InterfaceEntry
}

// NewReflectionInterface creates a new ReflectionInterface
func NewReflectionInterface(name string, iface *types.InterfaceEntry) *ReflectionInterface {
	return &ReflectionInterface{
		name:  name,
		iface: iface,
	}
}

// GetName returns the interface name
func (ri *ReflectionInterface) GetName() string {
	return ri.name
}

// GetShortName returns the short name (without namespace)
func (ri *ReflectionInterface) GetShortName() string {
	// Extract short name from full name
	// Similar logic to ReflectionClass
	return ri.name // Simplified for now
}

// GetMethods returns all methods defined in the interface
func (ri *ReflectionInterface) GetMethods() []*ReflectionMethod {
	if ri.iface == nil {
		return nil
	}

	methods := make([]*ReflectionMethod, 0, len(ri.iface.Methods))
	for name, method := range ri.iface.Methods {
		// Create a ReflectionClass-like wrapper for the interface
		classWrapper := &ReflectionClass{
			name: ri.name,
		}
		methods = append(methods, NewReflectionMethod(classWrapper, name, method))
	}
	return methods
}

// GetConstants returns all constants defined in the interface
func (ri *ReflectionInterface) GetConstants() map[string]*types.Value {
	if ri.iface == nil {
		return nil
	}

	constants := make(map[string]*types.Value)
	for name, constant := range ri.iface.Constants {
		constants[name] = constant.Value
	}
	return constants
}

// String returns a string representation of the interface
func (ri *ReflectionInterface) String() string {
	return "Interface [ " + ri.name + " ]"
}
