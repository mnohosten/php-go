package compiler

import (
	"fmt"

	"github.com/krizos/php-go/pkg/ast"
	"github.com/krizos/php-go/pkg/types"
)

// NamedArgumentResolver handles resolution of named arguments to positional order
type NamedArgumentResolver struct {
	// Cache of parameter names for functions/methods
	// Key: function/method name
	// Value: slice of parameter names in order
	parameterCache map[string][]string
}

// NewNamedArgumentResolver creates a new named argument resolver
func NewNamedArgumentResolver() *NamedArgumentResolver {
	return &NamedArgumentResolver{
		parameterCache: make(map[string][]string),
	}
}

// RegisterFunction registers a function's parameter names
func (nar *NamedArgumentResolver) RegisterFunction(name string, params []string) {
	nar.parameterCache[name] = params
}

// RegisterMethod registers a method's parameter names
func (nar *NamedArgumentResolver) RegisterMethod(className string, methodName string, params []string) {
	key := className + "::" + methodName
	nar.parameterCache[key] = params
}

// ResolveArguments resolves function call arguments, reordering named arguments
// to match the positional parameter order.
//
// Rules:
// 1. Positional arguments must come before named arguments
// 2. Named arguments can be in any order
// 3. All required parameters must be provided
// 4. Named arguments override positional for the same parameter
//
// Returns:
// - resolvedArgs: Arguments in correct positional order
// - error: if validation fails
func (nar *NamedArgumentResolver) ResolveArguments(
	functionName string,
	arguments []*ast.Argument,
	params []*types.ParameterDef,
) ([]*ast.Argument, error) {

	if len(arguments) == 0 {
		return arguments, nil
	}

	// Check if there are any named arguments
	hasNamed := false
	firstNamedIndex := -1
	for i, arg := range arguments {
		if arg.Name != "" {
			hasNamed = true
			if firstNamedIndex == -1 {
				firstNamedIndex = i
			}
		} else if hasNamed {
			// Positional argument after named argument
			return nil, fmt.Errorf("Cannot use positional argument after named argument in %s()", functionName)
		}
	}

	// If no named arguments, return as-is
	if !hasNamed {
		return arguments, nil
	}

	// Create a map to track which parameters have been filled
	paramMap := make(map[string]*ast.Argument)
	filledPositions := make(map[int]bool)

	// First, place positional arguments (before first named argument)
	for i := 0; i < firstNamedIndex; i++ {
		if i >= len(params) {
			return nil, fmt.Errorf("Too many positional arguments for %s()", functionName)
		}
		paramMap[params[i].Name] = arguments[i]
		filledPositions[i] = true
	}

	// Then, place named arguments
	for i := firstNamedIndex; i < len(arguments); i++ {
		arg := arguments[i]
		if arg.Name == "" {
			continue // Should not happen due to earlier validation
		}

		// Find the parameter index for this name
		paramIndex := -1
		for j, param := range params {
			if param.Name == arg.Name {
				paramIndex = j
				break
			}
		}

		if paramIndex == -1 {
			return nil, fmt.Errorf("Unknown named parameter: %s for %s()", arg.Name, functionName)
		}

		// Check if this position was already filled by a positional argument
		if filledPositions[paramIndex] {
			return nil, fmt.Errorf("Named parameter %s is already filled by positional argument in %s()", arg.Name, functionName)
		}

		paramMap[params[paramIndex].Name] = arg
		filledPositions[paramIndex] = true
	}

	// Build the resolved argument list in positional order
	resolvedArgs := make([]*ast.Argument, len(params))
	for i, param := range params {
		if arg, exists := paramMap[param.Name]; exists {
			resolvedArgs[i] = arg
		} else if param.HasDefault {
			// Parameter has default value, leave nil (will be handled at runtime)
			resolvedArgs[i] = nil
		} else if param.IsVariadic {
			// Variadic parameter can be empty
			resolvedArgs[i] = nil
		} else {
			// Required parameter not provided
			return nil, fmt.Errorf("Missing required parameter: $%s for %s()", param.Name, functionName)
		}
	}

	// Remove trailing nils for parameters with defaults
	actualLength := len(resolvedArgs)
	for i := len(resolvedArgs) - 1; i >= 0; i-- {
		if resolvedArgs[i] != nil {
			break
		}
		actualLength = i
	}
	resolvedArgs = resolvedArgs[:actualLength]

	return resolvedArgs, nil
}

// ValidateArguments performs validation on arguments before compilation
func (nar *NamedArgumentResolver) ValidateArguments(
	functionName string,
	arguments []*ast.Argument,
) error {

	hasNamed := false
	for _, arg := range arguments {
		if arg.Name != "" {
			hasNamed = true
		} else if hasNamed {
			return fmt.Errorf("Cannot use positional argument after named argument in %s()", functionName)
		}
	}

	return nil
}

// GetParameterNames extracts parameter names from parameter definitions
func GetParameterNames(params []*types.ParameterDef) []string {
	names := make([]string, len(params))
	for i, param := range params {
		names[i] = param.Name
	}
	return names
}

// HasNamedArguments checks if any arguments are named
func HasNamedArguments(arguments []*ast.Argument) bool {
	for _, arg := range arguments {
		if arg.Name != "" {
			return true
		}
	}
	return false
}

// CountPositionalArguments counts the number of positional arguments
// (before the first named argument)
func CountPositionalArguments(arguments []*ast.Argument) int {
	for i, arg := range arguments {
		if arg.Name != "" {
			return i
		}
	}
	return len(arguments)
}
