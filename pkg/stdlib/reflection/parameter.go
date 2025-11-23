package reflection

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ReflectionParameter represents a parameter reflection
type ReflectionParameter struct {
	// The method this parameter belongs to
	method *ReflectionMethod

	// Parameter position (0-based)
	position int

	// Parameter definition
	parameter *types.ParameterDef
}

// NewReflectionParameter creates a new ReflectionParameter
func NewReflectionParameter(method *ReflectionMethod, position int, parameter *types.ParameterDef) *ReflectionParameter {
	return &ReflectionParameter{
		method:    method,
		position:  position,
		parameter: parameter,
	}
}

// GetName returns the parameter name
func (rp *ReflectionParameter) GetName() string {
	return rp.parameter.Name
}

// GetPosition returns the parameter position
func (rp *ReflectionParameter) GetPosition() int {
	return rp.position
}

// GetType returns the parameter type declaration
func (rp *ReflectionParameter) GetType() string {
	return rp.parameter.Type
}

// HasType checks if the parameter has a type declaration
func (rp *ReflectionParameter) HasType() bool {
	return rp.parameter.Type != ""
}

// IsOptional checks if the parameter is optional (has default value)
func (rp *ReflectionParameter) IsOptional() bool {
	return rp.parameter.HasDefault
}

// GetDefaultValue returns the default value
func (rp *ReflectionParameter) GetDefaultValue() *types.Value {
	if rp.parameter.HasDefault {
		return rp.parameter.Default
	}
	return types.NewNull()
}

// IsPassedByReference checks if the parameter is passed by reference
func (rp *ReflectionParameter) IsPassedByReference() bool {
	return rp.parameter.PassedByRef
}

// IsVariadic checks if the parameter is variadic
func (rp *ReflectionParameter) IsVariadic() bool {
	return rp.parameter.IsVariadic
}

// AllowsNull checks if the parameter allows null
func (rp *ReflectionParameter) AllowsNull() bool {
	// If no type, allows null
	if !rp.HasType() {
		return true
	}
	// Check if type starts with ? (nullable)
	return len(rp.parameter.Type) > 0 && rp.parameter.Type[0] == '?'
}

// GetDeclaringFunction returns the method this parameter belongs to
func (rp *ReflectionParameter) GetDeclaringFunction() *ReflectionMethod {
	return rp.method
}

// String returns a string representation of the parameter
func (rp *ReflectionParameter) String() string {
	typeStr := ""
	if rp.HasType() {
		typeStr = rp.GetType() + " "
	}

	byRef := ""
	if rp.IsPassedByReference() {
		byRef = "&"
	}

	variadic := ""
	if rp.IsVariadic() {
		variadic = "..."
	}

	defaultVal := ""
	if rp.IsOptional() {
		defaultVal = fmt.Sprintf(" = %v", rp.GetDefaultValue())
	}

	return fmt.Sprintf("Parameter #%d [ %s%s%s$%s%s ]",
		rp.position, typeStr, byRef, variadic, rp.GetName(), defaultVal)
}
