package reflection

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ReflectionMethod represents a method reflection
type ReflectionMethod struct {
	// The class this method belongs to
	class *ReflectionClass

	// Method name
	name string

	// Method definition
	method *types.MethodDef
}

// NewReflectionMethod creates a new ReflectionMethod
func NewReflectionMethod(class *ReflectionClass, name string, method *types.MethodDef) *ReflectionMethod {
	return &ReflectionMethod{
		class:  class,
		name:   name,
		method: method,
	}
}

// GetName returns the method name
func (rm *ReflectionMethod) GetName() string {
	return rm.name
}

// GetDeclaringClass returns the class that declared this method
func (rm *ReflectionMethod) GetDeclaringClass() *ReflectionClass {
	return rm.class
}

// IsPublic checks if the method is public
func (rm *ReflectionMethod) IsPublic() bool {
	return rm.method.Visibility == types.VisibilityPublic
}

// IsProtected checks if the method is protected
func (rm *ReflectionMethod) IsProtected() bool {
	return rm.method.Visibility == types.VisibilityProtected
}

// IsPrivate checks if the method is private
func (rm *ReflectionMethod) IsPrivate() bool {
	return rm.method.Visibility == types.VisibilityPrivate
}

// IsStatic checks if the method is static
func (rm *ReflectionMethod) IsStatic() bool {
	return rm.method.IsStatic
}

// IsFinal checks if the method is final
func (rm *ReflectionMethod) IsFinal() bool {
	return rm.method.IsFinal
}

// IsAbstract checks if the method is abstract
func (rm *ReflectionMethod) IsAbstract() bool {
	return rm.method.IsAbstract
}

// IsConstructor checks if this is a constructor
func (rm *ReflectionMethod) IsConstructor() bool {
	return rm.method.IsConstructor
}

// IsDestructor checks if this is a destructor
func (rm *ReflectionMethod) IsDestructor() bool {
	return rm.method.IsDestructor
}

// GetNumberOfParameters returns the number of parameters
func (rm *ReflectionMethod) GetNumberOfParameters() int {
	return rm.method.NumParams
}

// GetNumberOfRequiredParameters returns the number of required parameters
func (rm *ReflectionMethod) GetNumberOfRequiredParameters() int {
	required := 0
	for _, param := range rm.method.Parameters {
		if !param.HasDefault {
			required++
		}
	}
	return required
}

// GetParameters returns all parameters
func (rm *ReflectionMethod) GetParameters() []*ReflectionParameter {
	params := make([]*ReflectionParameter, len(rm.method.Parameters))
	for i, param := range rm.method.Parameters {
		params[i] = NewReflectionParameter(rm, i, param)
	}
	return params
}

// GetReturnType returns the return type declaration
func (rm *ReflectionMethod) GetReturnType() string {
	return rm.method.ReturnType
}

// HasReturnType checks if the method has a return type
func (rm *ReflectionMethod) HasReturnType() bool {
	return rm.method.ReturnType != ""
}

// ReturnsReference checks if the method returns by reference
func (rm *ReflectionMethod) ReturnsReference() bool {
	return rm.method.ReturnByRef
}

// Invoke invokes the method on an object
// Note: In a full implementation, this would actually call the VM to execute the method
func (rm *ReflectionMethod) Invoke(obj *types.Object, args ...*types.Value) (*types.Value, error) {
	// This is a placeholder - full implementation would require VM integration
	return nil, fmt.Errorf("Method invocation not yet implemented")
}

// InvokeArgs invokes the method with arguments as an array
func (rm *ReflectionMethod) InvokeArgs(obj *types.Object, args []*types.Value) (*types.Value, error) {
	return rm.Invoke(obj, args...)
}

// SetAccessible makes the method accessible (bypasses visibility)
func (rm *ReflectionMethod) SetAccessible(accessible bool) {
	// In Go implementation, we can always access methods
	// This is mainly for API compatibility with PHP
}

// String returns a string representation of the method
func (rm *ReflectionMethod) String() string {
	visibility := "public"
	if rm.IsProtected() {
		visibility = "protected"
	} else if rm.IsPrivate() {
		visibility = "private"
	}

	modifiers := ""
	if rm.IsFinal() {
		modifiers += "final "
	}
	if rm.IsAbstract() {
		modifiers += "abstract "
	}
	if rm.IsStatic() {
		modifiers += "static "
	}

	returnType := ""
	if rm.HasReturnType() {
		returnType = ": " + rm.GetReturnType()
	}

	return fmt.Sprintf("Method [ %s %s%s%s() %s]",
		visibility, modifiers, rm.name, "", returnType)
}
