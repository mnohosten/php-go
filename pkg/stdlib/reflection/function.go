package reflection

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ReflectionFunction represents a function reflection
// Provides methods to inspect and manipulate functions at runtime
type ReflectionFunction struct {
	// The function being reflected
	function *types.FunctionDef

	// Function name
	name string
}

// NewReflectionFunction creates a new ReflectionFunction
func NewReflectionFunction(name string, function *types.FunctionDef) *ReflectionFunction {
	return &ReflectionFunction{
		function: function,
		name:     name,
	}
}

// ============================================================================
// Basic Function Information
// ============================================================================

// GetName returns the function name
func (rf *ReflectionFunction) GetName() string {
	if rf.function != nil {
		return rf.function.Name
	}
	return rf.name
}

// GetShortName returns the short name (without namespace)
func (rf *ReflectionFunction) GetShortName() string {
	name := rf.GetName()
	// Extract short name from full name
	// In PHP, functions don't have namespaces in the same way as classes
	// but they can be namespaced: MyNamespace\myFunction
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '\\' {
			return name[i+1:]
		}
	}
	return name
}

// GetNamespaceName returns the namespace name
func (rf *ReflectionFunction) GetNamespaceName() string {
	name := rf.GetName()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '\\' {
			return name[:i]
		}
	}
	return ""
}

// GetFileName returns the filename where the function was defined
func (rf *ReflectionFunction) GetFileName() string {
	if rf.function != nil {
		return rf.function.FileName
	}
	return ""
}

// GetStartLine returns the starting line number
func (rf *ReflectionFunction) GetStartLine() int {
	if rf.function != nil {
		return rf.function.StartLine
	}
	return 0
}

// GetEndLine returns the ending line number
func (rf *ReflectionFunction) GetEndLine() int {
	if rf.function != nil {
		return rf.function.EndLine
	}
	return 0
}

// ============================================================================
// Function Characteristics
// ============================================================================

// IsInternal checks if this is a built-in function
func (rf *ReflectionFunction) IsInternal() bool {
	if rf.function != nil {
		return rf.function.IsInternal
	}
	return false
}

// IsUserDefined checks if this is a user-defined function
func (rf *ReflectionFunction) IsUserDefined() bool {
	return !rf.IsInternal()
}

// IsGenerator checks if this is a generator function
func (rf *ReflectionFunction) IsGenerator() bool {
	if rf.function != nil {
		return rf.function.IsGenerator
	}
	return false
}

// IsDeprecated checks if this function is deprecated
func (rf *ReflectionFunction) IsDeprecated() bool {
	if rf.function != nil {
		return rf.function.IsDeprecated
	}
	return false
}

// IsVariadic checks if this function accepts variable number of arguments
func (rf *ReflectionFunction) IsVariadic() bool {
	if rf.function == nil {
		return false
	}
	// Check if last parameter is variadic
	if len(rf.function.Parameters) > 0 {
		lastParam := rf.function.Parameters[len(rf.function.Parameters)-1]
		return lastParam.IsVariadic
	}
	return false
}

// ============================================================================
// Parameters
// ============================================================================

// GetNumberOfParameters returns the total number of parameters
func (rf *ReflectionFunction) GetNumberOfParameters() int {
	if rf.function != nil {
		return rf.function.NumParams
	}
	return 0
}

// GetNumberOfRequiredParameters returns the number of required parameters
func (rf *ReflectionFunction) GetNumberOfRequiredParameters() int {
	if rf.function == nil {
		return 0
	}

	required := 0
	for _, param := range rf.function.Parameters {
		if !param.HasDefault {
			required++
		}
	}
	return required
}

// GetParameters returns all parameters
func (rf *ReflectionFunction) GetParameters() []*ReflectionParameter {
	if rf.function == nil {
		return nil
	}

	params := make([]*ReflectionParameter, len(rf.function.Parameters))
	for i, param := range rf.function.Parameters {
		// Create a minimal ReflectionMethod wrapper for the function
		methodWrapper := &ReflectionMethod{
			name: rf.name,
		}
		params[i] = NewReflectionParameter(methodWrapper, i, param)
	}
	return params
}

// ============================================================================
// Return Type
// ============================================================================

// GetReturnType returns the return type declaration
func (rf *ReflectionFunction) GetReturnType() string {
	if rf.function != nil {
		return rf.function.ReturnType
	}
	return ""
}

// HasReturnType checks if the function has a return type
func (rf *ReflectionFunction) HasReturnType() bool {
	return rf.GetReturnType() != ""
}

// ReturnsReference checks if the function returns by reference
func (rf *ReflectionFunction) ReturnsReference() bool {
	if rf.function != nil {
		return rf.function.ReturnByRef
	}
	return false
}

// ============================================================================
// Documentation
// ============================================================================

// GetDocComment returns the documentation comment
func (rf *ReflectionFunction) GetDocComment() string {
	if rf.function != nil {
		return rf.function.DocComment
	}
	return ""
}

// ============================================================================
// Invocation
// ============================================================================

// Invoke invokes the function with given arguments
// Note: In a full implementation, this would actually call the VM to execute the function
func (rf *ReflectionFunction) Invoke(args ...*types.Value) (*types.Value, error) {
	// This is a placeholder - full implementation would require VM integration
	return nil, fmt.Errorf("Function invocation not yet implemented")
}

// InvokeArgs invokes the function with arguments as an array
func (rf *ReflectionFunction) InvokeArgs(args []*types.Value) (*types.Value, error) {
	return rf.Invoke(args...)
}

// GetClosure returns a closure for the function
// Note: This creates a callable closure from the function
func (rf *ReflectionFunction) GetClosure() (*types.Value, error) {
	// This is a placeholder - full implementation would require VM integration
	return nil, fmt.Errorf("GetClosure not yet implemented")
}

// ============================================================================
// String Representation
// ============================================================================

// String returns a string representation of the function
func (rf *ReflectionFunction) String() string {
	returnType := ""
	if rf.HasReturnType() {
		returnType = ": " + rf.GetReturnType()
	}

	byRef := ""
	if rf.ReturnsReference() {
		byRef = "&"
	}

	internal := ""
	if rf.IsInternal() {
		internal = "<internal>"
	} else if rf.function != nil && rf.function.FileName != "" {
		internal = fmt.Sprintf("<%s:%d-%d>", rf.function.FileName, rf.function.StartLine, rf.function.EndLine)
	}

	return fmt.Sprintf("Function [ %s function %s%s() %s]",
		internal, byRef, rf.GetName(), returnType)
}

// ============================================================================
// Additional Utility Methods
// ============================================================================

// IsClosure checks if this is a closure/anonymous function
func (rf *ReflectionFunction) IsClosure() bool {
	// Closures typically have names like {closure} or empty names
	name := rf.GetName()
	return name == "" || (len(name) > 0 && name[0] == '{')
}

// IsDisabled checks if the function is disabled (via disable_functions)
// Note: This is a placeholder for PHP's disable_functions ini setting
func (rf *ReflectionFunction) IsDisabled() bool {
	return false
}

// GetExtension returns the extension that defined this function
// Note: Only relevant for internal functions
func (rf *ReflectionFunction) GetExtension() string {
	if rf.IsInternal() {
		// This could be expanded to track which extension provides which function
		return "Core"
	}
	return ""
}

// GetExtensionName returns the name of the extension that defined this function
func (rf *ReflectionFunction) GetExtensionName() string {
	return rf.GetExtension()
}
