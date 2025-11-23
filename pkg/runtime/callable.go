package runtime

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// First-Class Callable Support (PHP 8.1+)
// ============================================================================

// CallableType represents the type of callable
type CallableType int

const (
	CallableTypeFunction CallableType = iota
	CallableTypeMethod
	CallableTypeStaticMethod
	CallableTypeClosure
)

// Callable represents a first-class callable reference
// Created using the ... syntax: strlen(...), $obj->method(...), Class::method(...)
type Callable struct {
	Type CallableType

	// For functions
	FunctionName string

	// For methods
	Object     *types.Object // nil for static methods
	MethodName string
	ClassName  string // For static methods

	// For closures
	Closure *types.Value
}

// NewFunctionCallable creates a callable reference to a function
// Example: strlen(...)
func NewFunctionCallable(functionName string) *Callable {
	return &Callable{
		Type:         CallableTypeFunction,
		FunctionName: functionName,
	}
}

// NewMethodCallable creates a callable reference to an instance method
// Example: $obj->method(...)
func NewMethodCallable(object *types.Object, methodName string) *Callable {
	return &Callable{
		Type:       CallableTypeMethod,
		Object:     object,
		MethodName: methodName,
	}
}

// NewStaticMethodCallable creates a callable reference to a static method
// Example: Class::method(...)
func NewStaticMethodCallable(className string, methodName string) *Callable {
	return &Callable{
		Type:       CallableTypeStaticMethod,
		ClassName:  className,
		MethodName: methodName,
	}
}

// NewClosureCallable creates a callable reference from a closure
func NewClosureCallable(closure *types.Value) *Callable {
	return &Callable{
		Type:    CallableTypeClosure,
		Closure: closure,
	}
}

// Invoke calls the callable with the given arguments
func (c *Callable) Invoke(args []*types.Value) (*types.Value, error) {
	if c == nil {
		return nil, fmt.Errorf("Cannot invoke null callable")
	}

	switch c.Type {
	case CallableTypeFunction:
		// Call function by name
		// This would integrate with the VM's function call mechanism
		return nil, fmt.Errorf("Function call not yet implemented in callable")

	case CallableTypeMethod:
		// Call instance method
		if c.Object == nil {
			return nil, fmt.Errorf("Cannot call method on null object")
		}
		// This would integrate with the VM's method call mechanism
		return nil, fmt.Errorf("Method call not yet implemented in callable")

	case CallableTypeStaticMethod:
		// Call static method
		// This would integrate with the VM's static method call mechanism
		return nil, fmt.Errorf("Static method call not yet implemented in callable")

	case CallableTypeClosure:
		// Call closure
		if c.Closure == nil {
			return nil, fmt.Errorf("Cannot invoke null closure")
		}
		// This would integrate with the VM's closure call mechanism
		return nil, fmt.Errorf("Closure call not yet implemented in callable")

	default:
		return nil, fmt.Errorf("Unknown callable type")
	}
}

// GetName returns a string representation of the callable
func (c *Callable) GetName() string {
	if c == nil {
		return ""
	}

	switch c.Type {
	case CallableTypeFunction:
		return c.FunctionName

	case CallableTypeMethod:
		if c.Object != nil {
			return c.Object.ClassName + "->" + c.MethodName
		}
		return c.MethodName

	case CallableTypeStaticMethod:
		return c.ClassName + "::" + c.MethodName

	case CallableTypeClosure:
		return "{closure}"

	default:
		return ""
	}
}

// ToValue converts the callable to a PHP value (as a resource or closure)
func (c *Callable) ToValue() *types.Value {
	if c == nil {
		return types.NewNull()
	}

	// For closures, return the closure value directly
	if c.Type == CallableTypeClosure && c.Closure != nil {
		return c.Closure
	}

	// For other callables, wrap as a resource
	resource := types.NewResourceHandle("Callable", c)
	return types.NewResource(resource)
}

// IsCallable checks if a value is callable
// This implements PHP's is_callable() function
func IsCallable(value *types.Value) bool {
	if value == nil {
		return false
	}

	switch value.Type() {
	case types.TypeString:
		// Function name as string
		return true // Would need to check if function exists

	case types.TypeArray:
		// Array callable: ['Class', 'method'] or [$obj, 'method']
		arr := value.ToArray()
		if arr.Len() == 2 {
			return true // Would need to validate the array structure
		}
		return false

	case types.TypeObject:
		// Object with __invoke method
		obj := value.ToObject()
		if obj != nil && obj.ClassEntry != nil {
			_, hasInvoke := obj.ClassEntry.Methods["__invoke"]
			return hasInvoke
		}
		return false

	case types.TypeResource:
		// Check if it's a Callable resource
		resource := value.ToResource()
		if resource != nil {
			_, ok := resource.Data().(*Callable)
			return ok
		}
		return false

	default:
		return false
	}
}

// CallableFromValue extracts a Callable from a resource value
func CallableFromValue(value *types.Value) (*Callable, error) {
	if value.Type() != types.TypeResource {
		return nil, fmt.Errorf("Expected Callable resource, got %v", value.Type())
	}

	resource := value.ToResource()
	if resource == nil {
		return nil, fmt.Errorf("Invalid resource value")
	}

	callable, ok := resource.Data().(*Callable)
	if !ok {
		return nil, fmt.Errorf("Invalid Callable resource")
	}

	return callable, nil
}

// ============================================================================
// Callable Validation
// ============================================================================

// ValidateCallable validates that a callable reference is valid
// This checks that the referenced function/method exists
func ValidateCallable(callable *Callable) error {
	if callable == nil {
		return fmt.Errorf("Callable is nil")
	}

	switch callable.Type {
	case CallableTypeFunction:
		if callable.FunctionName == "" {
			return fmt.Errorf("Function name is empty")
		}
		// Would check if function exists in function table
		return nil

	case CallableTypeMethod:
		if callable.Object == nil {
			return fmt.Errorf("Object is nil for method callable")
		}
		if callable.MethodName == "" {
			return fmt.Errorf("Method name is empty")
		}
		// Would check if method exists on object
		return nil

	case CallableTypeStaticMethod:
		if callable.ClassName == "" {
			return fmt.Errorf("Class name is empty for static method callable")
		}
		if callable.MethodName == "" {
			return fmt.Errorf("Method name is empty")
		}
		// Would check if static method exists on class
		return nil

	case CallableTypeClosure:
		if callable.Closure == nil {
			return fmt.Errorf("Closure is nil")
		}
		return nil

	default:
		return fmt.Errorf("Unknown callable type: %d", callable.Type)
	}
}
