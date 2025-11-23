package reflection

import (
	"fmt"
	"strings"

	"github.com/krizos/php-go/pkg/types"
)

// ReflectionClass represents a class reflection
// Provides methods to inspect and manipulate classes at runtime
type ReflectionClass struct {
	// The class being reflected
	class *types.ClassEntry

	// Name of the class
	name string
}

// NewReflectionClass creates a new ReflectionClass for the given class name
func NewReflectionClass(className string, classEntry *types.ClassEntry) *ReflectionClass {
	return &ReflectionClass{
		class: classEntry,
		name:  className,
	}
}

// NewReflectionClassFromObject creates a ReflectionClass from an object instance
func NewReflectionClassFromObject(obj *types.Object) *ReflectionClass {
	return &ReflectionClass{
		class: obj.ClassEntry,
		name:  obj.ClassName,
	}
}

// ============================================================================
// Basic Class Information
// ============================================================================

// GetName returns the class name
func (rc *ReflectionClass) GetName() string {
	if rc.class != nil {
		return rc.class.Name
	}
	return rc.name
}

// GetShortName returns the short name (without namespace)
func (rc *ReflectionClass) GetShortName() string {
	if rc.class != nil {
		return rc.class.ShortName
	}
	// Extract short name from full name
	parts := strings.Split(rc.name, "\\")
	return parts[len(parts)-1]
}

// GetNamespaceName returns the namespace name
func (rc *ReflectionClass) GetNamespaceName() string {
	if rc.class != nil {
		return rc.class.Namespace
	}
	// Extract namespace from full name
	parts := strings.Split(rc.name, "\\")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "\\")
	}
	return ""
}

// GetFileName returns the filename where the class was defined
func (rc *ReflectionClass) GetFileName() string {
	if rc.class != nil {
		return rc.class.FileName
	}
	return ""
}

// ============================================================================
// Class Modifiers
// ============================================================================

// IsFinal checks if the class is final
func (rc *ReflectionClass) IsFinal() bool {
	if rc.class != nil {
		return rc.class.IsFinal
	}
	return false
}

// IsAbstract checks if the class is abstract
func (rc *ReflectionClass) IsAbstract() bool {
	if rc.class != nil {
		return rc.class.IsAbstract
	}
	return false
}

// IsReadOnly checks if the class is readonly (PHP 8.2+)
func (rc *ReflectionClass) IsReadOnly() bool {
	if rc.class != nil {
		return rc.class.IsReadOnly
	}
	return false
}

// IsInterface checks if this is an interface
func (rc *ReflectionClass) IsInterface() bool {
	if rc.class != nil {
		return rc.class.IsInterface
	}
	return false
}

// IsTrait checks if this is a trait
func (rc *ReflectionClass) IsTrait() bool {
	if rc.class != nil {
		return rc.class.IsTrait
	}
	return false
}

// IsEnum checks if this is an enum (PHP 8.1+)
func (rc *ReflectionClass) IsEnum() bool {
	if rc.class != nil {
		return rc.class.IsEnum
	}
	return false
}

// IsInstantiable checks if the class can be instantiated
func (rc *ReflectionClass) IsInstantiable() bool {
	if rc.class == nil {
		return false
	}
	// Abstract classes and interfaces cannot be instantiated
	return !rc.class.IsAbstract && !rc.class.IsInterface && !rc.class.IsTrait
}

// ============================================================================
// Inheritance and Interfaces
// ============================================================================

// GetParentClass returns the parent class
func (rc *ReflectionClass) GetParentClass() *ReflectionClass {
	if rc.class != nil && rc.class.ParentClass != nil {
		return NewReflectionClass(rc.class.ParentClass.Name, rc.class.ParentClass)
	}
	return nil
}

// GetInterfaces returns all implemented interfaces
func (rc *ReflectionClass) GetInterfaces() []*ReflectionInterface {
	if rc.class == nil {
		return nil
	}

	interfaces := make([]*ReflectionInterface, 0, len(rc.class.Interfaces))
	for _, iface := range rc.class.Interfaces {
		interfaces = append(interfaces, NewReflectionInterface(iface.Name, iface))
	}
	return interfaces
}

// GetInterfaceNames returns names of all implemented interfaces
func (rc *ReflectionClass) GetInterfaceNames() []string {
	if rc.class == nil {
		return nil
	}

	names := make([]string, 0, len(rc.class.Interfaces))
	for _, iface := range rc.class.Interfaces {
		names = append(names, iface.Name)
	}
	return names
}

// ImplementsInterface checks if the class implements a specific interface
func (rc *ReflectionClass) ImplementsInterface(interfaceName string) bool {
	if rc.class == nil {
		return false
	}

	for _, iface := range rc.class.Interfaces {
		if iface.Name == interfaceName {
			return true
		}
	}
	return false
}

// GetTraits returns all used traits
func (rc *ReflectionClass) GetTraits() []*ReflectionClass {
	if rc.class == nil {
		return nil
	}

	traits := make([]*ReflectionClass, 0, len(rc.class.Traits))
	for _, trait := range rc.class.Traits {
		// Convert TraitEntry to ClassEntry-like structure
		// Note: TraitEntry should have similar structure to ClassEntry
		traitClass := &types.ClassEntry{
			Name:    trait.Name,
			IsTrait: true,
		}
		traits = append(traits, NewReflectionClass(trait.Name, traitClass))
	}
	return traits
}

// GetTraitNames returns names of all used traits
func (rc *ReflectionClass) GetTraitNames() []string {
	if rc.class == nil {
		return nil
	}

	names := make([]string, 0, len(rc.class.Traits))
	for _, trait := range rc.class.Traits {
		names = append(names, trait.Name)
	}
	return names
}

// IsSubclassOf checks if this class is a subclass of another
func (rc *ReflectionClass) IsSubclassOf(className string) bool {
	if rc.class == nil {
		return false
	}

	// Walk up the parent chain
	parent := rc.class.ParentClass
	for parent != nil {
		if parent.Name == className {
			return true
		}
		parent = parent.ParentClass
	}
	return false
}

// ============================================================================
// Properties
// ============================================================================

// GetProperties returns all properties
func (rc *ReflectionClass) GetProperties(filter int) []*ReflectionProperty {
	if rc.class == nil {
		return nil
	}

	properties := make([]*ReflectionProperty, 0)
	for name, prop := range rc.class.Properties {
		// Apply filter
		if filter != 0 && !matchesPropertyFilter(prop, filter) {
			continue
		}

		properties = append(properties, NewReflectionProperty(rc, name, prop))
	}
	return properties
}

// GetProperty returns a specific property by name
func (rc *ReflectionClass) GetProperty(name string) (*ReflectionProperty, error) {
	if rc.class == nil {
		return nil, fmt.Errorf("Class not loaded")
	}

	if prop, ok := rc.class.Properties[name]; ok {
		return NewReflectionProperty(rc, name, prop), nil
	}
	return nil, fmt.Errorf("Property %s does not exist", name)
}

// HasProperty checks if a property exists
func (rc *ReflectionClass) HasProperty(name string) bool {
	if rc.class == nil {
		return false
	}
	_, ok := rc.class.Properties[name]
	return ok
}

// GetStaticProperties returns all static properties
func (rc *ReflectionClass) GetStaticProperties() map[string]*types.Value {
	if rc.class == nil {
		return nil
	}
	return rc.class.StaticProperties
}

// GetDefaultProperties returns default values for instance properties
func (rc *ReflectionClass) GetDefaultProperties() map[string]*types.Value {
	if rc.class == nil {
		return nil
	}
	return rc.class.DefaultProperties
}

// ============================================================================
// Methods
// ============================================================================

// GetMethods returns all methods
func (rc *ReflectionClass) GetMethods(filter int) []*ReflectionMethod {
	if rc.class == nil {
		return nil
	}

	methods := make([]*ReflectionMethod, 0)
	for name, method := range rc.class.Methods {
		// Apply filter
		if filter != 0 && !matchesMethodFilter(method, filter) {
			continue
		}

		methods = append(methods, NewReflectionMethod(rc, name, method))
	}
	return methods
}

// GetMethod returns a specific method by name
func (rc *ReflectionClass) GetMethod(name string) (*ReflectionMethod, error) {
	if rc.class == nil {
		return nil, fmt.Errorf("Class not loaded")
	}

	if method, ok := rc.class.Methods[name]; ok {
		return NewReflectionMethod(rc, name, method), nil
	}
	return nil, fmt.Errorf("Method %s does not exist", name)
}

// HasMethod checks if a method exists
func (rc *ReflectionClass) HasMethod(name string) bool {
	if rc.class == nil {
		return false
	}
	_, ok := rc.class.Methods[name]
	return ok
}

// GetConstructor returns the constructor method
func (rc *ReflectionClass) GetConstructor() *ReflectionMethod {
	if rc.class == nil || rc.class.Constructor == nil {
		return nil
	}
	return NewReflectionMethod(rc, "__construct", rc.class.Constructor)
}

// HasConstructor checks if the class has a constructor
func (rc *ReflectionClass) HasConstructor() bool {
	return rc.class != nil && rc.class.Constructor != nil
}

// GetDestructor returns the destructor method
func (rc *ReflectionClass) GetDestructor() *ReflectionMethod {
	if rc.class == nil || rc.class.Destructor == nil {
		return nil
	}
	return NewReflectionMethod(rc, "__destruct", rc.class.Destructor)
}

// ============================================================================
// Constants
// ============================================================================

// GetConstants returns all class constants
func (rc *ReflectionClass) GetConstants() map[string]*types.Value {
	if rc.class == nil {
		return nil
	}

	constants := make(map[string]*types.Value)
	for name, constant := range rc.class.Constants {
		constants[name] = constant.Value
	}
	return constants
}

// GetConstant returns a specific constant value
func (rc *ReflectionClass) GetConstant(name string) (*types.Value, error) {
	if rc.class == nil {
		return nil, fmt.Errorf("Class not loaded")
	}

	if constant, ok := rc.class.Constants[name]; ok {
		return constant.Value, nil
	}
	return nil, fmt.Errorf("Constant %s does not exist", name)
}

// HasConstant checks if a constant exists
func (rc *ReflectionClass) HasConstant(name string) bool {
	if rc.class == nil {
		return false
	}
	_, ok := rc.class.Constants[name]
	return ok
}

// ============================================================================
// Instantiation
// ============================================================================

// NewInstance creates a new instance of the class
func (rc *ReflectionClass) NewInstance(args ...*types.Value) (*types.Object, error) {
	if !rc.IsInstantiable() {
		return nil, fmt.Errorf("Cannot instantiate %s", rc.GetName())
	}

	// Create object
	obj := &types.Object{
		ClassName:  rc.class.Name,
		ClassEntry: rc.class,
		Properties: make(map[string]*types.Property),
	}

	// Initialize default properties
	for name, defaultVal := range rc.class.DefaultProperties {
		if propDef, ok := rc.class.Properties[name]; ok {
			obj.Properties[name] = &types.Property{
				Value:      defaultVal,
				Visibility: propDef.Visibility,
				IsStatic:   false,
			}
		}
	}

	// TODO: Call constructor with args

	return obj, nil
}

// ============================================================================
// String Representation
// ============================================================================

// String returns a string representation of the class
func (rc *ReflectionClass) String() string {
	var sb strings.Builder

	if rc.IsInterface() {
		sb.WriteString("Interface [ ")
	} else if rc.IsTrait() {
		sb.WriteString("Trait [ ")
	} else {
		sb.WriteString("Class [ ")
	}

	if rc.IsFinal() {
		sb.WriteString("final ")
	}
	if rc.IsAbstract() {
		sb.WriteString("abstract ")
	}
	if rc.IsReadOnly() {
		sb.WriteString("readonly ")
	}

	sb.WriteString(rc.GetName())
	sb.WriteString(" ]")

	return sb.String()
}

// ============================================================================
// Helper Functions
// ============================================================================

// Reflection filter constants
const (
	IS_PUBLIC    = 1 << 0
	IS_PROTECTED = 1 << 1
	IS_PRIVATE   = 1 << 2
	IS_STATIC    = 1 << 3
	IS_FINAL     = 1 << 4
	IS_ABSTRACT  = 1 << 5
)

// matchesPropertyFilter checks if a property matches the filter
func matchesPropertyFilter(prop *types.PropertyDef, filter int) bool {
	match := true

	if filter&IS_PUBLIC != 0 {
		match = match && prop.Visibility == types.VisibilityPublic
	}
	if filter&IS_PROTECTED != 0 {
		match = match && prop.Visibility == types.VisibilityProtected
	}
	if filter&IS_PRIVATE != 0 {
		match = match && prop.Visibility == types.VisibilityPrivate
	}
	if filter&IS_STATIC != 0 {
		match = match && prop.IsStatic
	}

	return match
}

// matchesMethodFilter checks if a method matches the filter
func matchesMethodFilter(method *types.MethodDef, filter int) bool {
	match := true

	if filter&IS_PUBLIC != 0 {
		match = match && method.Visibility == types.VisibilityPublic
	}
	if filter&IS_PROTECTED != 0 {
		match = match && method.Visibility == types.VisibilityProtected
	}
	if filter&IS_PRIVATE != 0 {
		match = match && method.Visibility == types.VisibilityPrivate
	}
	if filter&IS_STATIC != 0 {
		match = match && method.IsStatic
	}
	if filter&IS_FINAL != 0 {
		match = match && method.IsFinal
	}
	if filter&IS_ABSTRACT != 0 {
		match = match && method.IsAbstract
	}

	return match
}
