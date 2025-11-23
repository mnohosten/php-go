package runtime

import (
	"fmt"

	"github.com/krizos/php-go/pkg/types"
)

// ============================================================================
// Attribute System (PHP 8.0+)
// ============================================================================

// AttributeInstance represents an instantiated PHP attribute
// Attributes are metadata that can be attached to classes, methods, properties, etc.
type AttributeInstance struct {
	Name      string            // Fully qualified attribute name
	Arguments []*types.Value    // Positional arguments
	Named     map[string]*types.Value // Named arguments
	Target    AttributeTarget   // Where the attribute is applied
}

// AttributeTarget specifies what an attribute can be applied to
type AttributeTarget uint16

const (
	AttributeTargetClass      AttributeTarget = 1 << 0  // Classes
	AttributeTargetFunction   AttributeTarget = 1 << 1  // Functions
	AttributeTargetMethod     AttributeTarget = 1 << 2  // Methods
	AttributeTargetProperty   AttributeTarget = 1 << 3  // Properties
	AttributeTargetConstant   AttributeTarget = 1 << 4  // Class constants
	AttributeTargetParameter  AttributeTarget = 1 << 5  // Function/method parameters
	AttributeTargetAll        AttributeTarget = 0xFFFF  // All targets
)

// String returns a string representation of the target
func (t AttributeTarget) String() string {
	switch t {
	case AttributeTargetClass:
		return "CLASS"
	case AttributeTargetFunction:
		return "FUNCTION"
	case AttributeTargetMethod:
		return "METHOD"
	case AttributeTargetProperty:
		return "PROPERTY"
	case AttributeTargetConstant:
		return "CONSTANT"
	case AttributeTargetParameter:
		return "PARAMETER"
	case AttributeTargetAll:
		return "ALL"
	default:
		return fmt.Sprintf("TARGET(%d)", t)
	}
}

// NewAttributeInstance creates a new attribute instance
func NewAttributeInstance(name string, target AttributeTarget) *AttributeInstance {
	return &AttributeInstance{
		Name:      name,
		Arguments: make([]*types.Value, 0),
		Named:     make(map[string]*types.Value),
		Target:    target,
	}
}

// AddArgument adds a positional argument
func (ai *AttributeInstance) AddArgument(value *types.Value) {
	ai.Arguments = append(ai.Arguments, value)
}

// AddNamedArgument adds a named argument
func (ai *AttributeInstance) AddNamedArgument(name string, value *types.Value) {
	ai.Named[name] = value
}

// GetName returns the attribute name
func (ai *AttributeInstance) GetName() string {
	return ai.Name
}

// GetArguments returns all positional arguments
func (ai *AttributeInstance) GetArguments() []*types.Value {
	return ai.Arguments
}

// GetNamedArguments returns all named arguments
func (ai *AttributeInstance) GetNamedArguments() map[string]*types.Value {
	return ai.Named
}

// ============================================================================
// Built-in Attributes
// ============================================================================

// Built-in attribute names
const (
	// Deprecated marks a feature as deprecated
	AttributeDeprecated = "Deprecated"

	// ReturnTypeWillChange marks a method where return type might change
	AttributeReturnTypeWillChange = "#[ReturnTypeWillChange]"

	// AllowDynamicProperties allows dynamic properties on a class (PHP 8.2+)
	AttributeAllowDynamicProperties = "#[AllowDynamicProperties]"

	// SensitiveParameter marks a parameter as sensitive (masks in stack traces)
	AttributeSensitiveParameter = "#[SensitiveParameter]"

	// Override explicitly marks a method as overriding parent method (PHP 8.3+)
	AttributeOverride = "#[Override]"
)

// AttributeDefinition defines metadata for an attribute class
type AttributeDefinition struct {
	Name       string          // Attribute name
	Targets    AttributeTarget // Allowed targets
	Repeatable bool            // Can be used multiple times on same target
	Flags      int             // Additional flags
}

// NewAttributeDefinition creates a new attribute definition
func NewAttributeDefinition(name string, targets AttributeTarget, repeatable bool) *AttributeDefinition {
	return &AttributeDefinition{
		Name:       name,
		Targets:    targets,
		Repeatable: repeatable,
		Flags:      0,
	}
}

// ValidateTarget checks if the attribute can be applied to the given target
func (ad *AttributeDefinition) ValidateTarget(target AttributeTarget) error {
	if ad.Targets&target == 0 && ad.Targets != AttributeTargetAll {
		return fmt.Errorf("Attribute %s cannot be applied to %s", ad.Name, target.String())
	}
	return nil
}

// ============================================================================
// Attribute Registry
// ============================================================================

// AttributeRegistry manages attribute definitions
type AttributeRegistry struct {
	definitions map[string]*AttributeDefinition
}

// NewAttributeRegistry creates a new attribute registry
func NewAttributeRegistry() *AttributeRegistry {
	registry := &AttributeRegistry{
		definitions: make(map[string]*AttributeDefinition),
	}

	// Register built-in attributes
	registry.RegisterBuiltinAttributes()

	return registry
}

// RegisterBuiltinAttributes registers PHP's built-in attributes
func (ar *AttributeRegistry) RegisterBuiltinAttributes() {
	// Deprecated - can be applied to anything
	ar.Register(NewAttributeDefinition(
		AttributeDeprecated,
		AttributeTargetAll,
		false,
	))

	// ReturnTypeWillChange - methods only
	ar.Register(NewAttributeDefinition(
		AttributeReturnTypeWillChange,
		AttributeTargetMethod,
		false,
	))

	// AllowDynamicProperties - classes only
	ar.Register(NewAttributeDefinition(
		AttributeAllowDynamicProperties,
		AttributeTargetClass,
		false,
	))

	// SensitiveParameter - parameters only
	ar.Register(NewAttributeDefinition(
		AttributeSensitiveParameter,
		AttributeTargetParameter,
		false,
	))

	// Override - methods only
	ar.Register(NewAttributeDefinition(
		AttributeOverride,
		AttributeTargetMethod,
		false,
	))
}

// Register registers an attribute definition
func (ar *AttributeRegistry) Register(def *AttributeDefinition) {
	ar.definitions[def.Name] = def
}

// Get retrieves an attribute definition
func (ar *AttributeRegistry) Get(name string) (*AttributeDefinition, bool) {
	def, ok := ar.definitions[name]
	return def, ok
}

// Validate validates an attribute instance against its definition
func (ar *AttributeRegistry) Validate(instance *AttributeInstance) error {
	def, ok := ar.Get(instance.Name)
	if !ok {
		// Unknown attributes are allowed (user-defined)
		return nil
	}

	// Validate target
	return def.ValidateTarget(instance.Target)
}

// ============================================================================
// Attribute Storage
// ============================================================================

// AttributeStorage stores attributes for different entities
type AttributeStorage struct {
	// Class attributes
	ClassAttributes map[string][]*AttributeInstance

	// Function attributes
	FunctionAttributes map[string][]*AttributeInstance

	// Method attributes (className::methodName)
	MethodAttributes map[string][]*AttributeInstance

	// Property attributes (className::propertyName)
	PropertyAttributes map[string][]*AttributeInstance

	// Parameter attributes (className::methodName::parameterName or functionName::parameterName)
	ParameterAttributes map[string][]*AttributeInstance

	// Constant attributes (className::constantName)
	ConstantAttributes map[string][]*AttributeInstance
}

// NewAttributeStorage creates a new attribute storage
func NewAttributeStorage() *AttributeStorage {
	return &AttributeStorage{
		ClassAttributes:     make(map[string][]*AttributeInstance),
		FunctionAttributes:  make(map[string][]*AttributeInstance),
		MethodAttributes:    make(map[string][]*AttributeInstance),
		PropertyAttributes:  make(map[string][]*AttributeInstance),
		ParameterAttributes: make(map[string][]*AttributeInstance),
		ConstantAttributes:  make(map[string][]*AttributeInstance),
	}
}

// AddClassAttribute adds an attribute to a class
func (as *AttributeStorage) AddClassAttribute(className string, attr *AttributeInstance) {
	as.ClassAttributes[className] = append(as.ClassAttributes[className], attr)
}

// GetClassAttributes returns all attributes for a class
func (as *AttributeStorage) GetClassAttributes(className string) []*AttributeInstance {
	return as.ClassAttributes[className]
}

// AddFunctionAttribute adds an attribute to a function
func (as *AttributeStorage) AddFunctionAttribute(functionName string, attr *AttributeInstance) {
	as.FunctionAttributes[functionName] = append(as.FunctionAttributes[functionName], attr)
}

// GetFunctionAttributes returns all attributes for a function
func (as *AttributeStorage) GetFunctionAttributes(functionName string) []*AttributeInstance {
	return as.FunctionAttributes[functionName]
}

// AddMethodAttribute adds an attribute to a method
func (as *AttributeStorage) AddMethodAttribute(className, methodName string, attr *AttributeInstance) {
	key := className + "::" + methodName
	as.MethodAttributes[key] = append(as.MethodAttributes[key], attr)
}

// GetMethodAttributes returns all attributes for a method
func (as *AttributeStorage) GetMethodAttributes(className, methodName string) []*AttributeInstance {
	key := className + "::" + methodName
	return as.MethodAttributes[key]
}

// AddPropertyAttribute adds an attribute to a property
func (as *AttributeStorage) AddPropertyAttribute(className, propertyName string, attr *AttributeInstance) {
	key := className + "::" + propertyName
	as.PropertyAttributes[key] = append(as.PropertyAttributes[key], attr)
}

// GetPropertyAttributes returns all attributes for a property
func (as *AttributeStorage) GetPropertyAttributes(className, propertyName string) []*AttributeInstance {
	key := className + "::" + propertyName
	return as.PropertyAttributes[key]
}

// AddParameterAttribute adds an attribute to a parameter
func (as *AttributeStorage) AddParameterAttribute(scope, name, paramName string, attr *AttributeInstance) {
	key := scope + "::" + name + "::" + paramName
	as.ParameterAttributes[key] = append(as.ParameterAttributes[key], attr)
}

// GetParameterAttributes returns all attributes for a parameter
func (as *AttributeStorage) GetParameterAttributes(scope, name, paramName string) []*AttributeInstance {
	key := scope + "::" + name + "::" + paramName
	return as.ParameterAttributes[key]
}

// AddConstantAttribute adds an attribute to a constant
func (as *AttributeStorage) AddConstantAttribute(className, constantName string, attr *AttributeInstance) {
	key := className + "::" + constantName
	as.ConstantAttributes[key] = append(as.ConstantAttributes[key], attr)
}

// GetConstantAttributes returns all attributes for a constant
func (as *AttributeStorage) GetConstantAttributes(className, constantName string) []*AttributeInstance {
	key := className + "::" + constantName
	return as.ConstantAttributes[key]
}
