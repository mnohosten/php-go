package types

import "fmt"

// ============================================================================
// Object Structure
// ============================================================================

// Object represents a PHP object instance
type Object struct {
	// Runtime object data
	ClassName  string                // Name of the class this object is an instance of
	ClassEntry *ClassEntry           // Reference to the class definition
	Properties map[string]*Property  // Instance properties (key: property name)
	ObjectID   uint64                // Unique object identifier for identity comparison

	// Object state
	IsDestroyed bool // Whether __destruct() has been called
}

// Property represents an object property with metadata
type Property struct {
	Value      *Value            // The property value
	Visibility PropertyVisibility // public, protected, private
	IsStatic   bool              // Static vs instance property
	Type       string            // Type declaration (e.g., "int", "string", "?MyClass")
	HasDefault bool              // Whether property has a default value
	Default    *Value            // Default value
	IsReadOnly bool              // Readonly property (PHP 8.1+)
	Hooks      *PropertyHooks    // Property hooks (PHP 8.4+)
}

// PropertyVisibility defines property access levels
type PropertyVisibility uint8

const (
	VisibilityPublic PropertyVisibility = iota
	VisibilityProtected
	VisibilityPrivate
)

// String returns the string representation of visibility
func (v PropertyVisibility) String() string {
	switch v {
	case VisibilityPublic:
		return "public"
	case VisibilityProtected:
		return "protected"
	case VisibilityPrivate:
		return "private"
	default:
		return "unknown"
	}
}

// PropertyHooks represents property get/set hooks (PHP 8.4+)
type PropertyHooks struct {
	Get *HookFunction // Getter hook: get { return ...; }
	Set *HookFunction // Setter hook: set { $field = $value; }
}

// HookFunction represents a property hook function
type HookFunction struct {
	Instructions []interface{} // Bytecode instructions for the hook
	NumLocals    int            // Number of local variables
}

// ============================================================================
// Class Entry (Class Definition/Metadata)
// ============================================================================

// ClassEntry represents a PHP class definition with full metadata
type ClassEntry struct {
	// Basic class information
	Name       string // Fully qualified class name
	ShortName  string // Short name (without namespace)
	Namespace  string // Namespace the class belongs to
	FileName   string // File where class was declared

	// Class modifiers
	IsFinal    bool // final class (cannot be extended)
	IsAbstract bool // abstract class (cannot be instantiated)
	IsReadOnly bool // readonly class (PHP 8.2+) - all properties are readonly

	// Inheritance and composition
	ParentClass  *ClassEntry            // Parent class (null if no parent)
	Interfaces   []*InterfaceEntry      // Implemented interfaces
	Traits       []*TraitEntry          // Used traits
	TraitAliases map[string]string      // Trait method aliases
	TraitPrecedence map[string]string   // Trait conflict resolution

	// Class constants
	Constants map[string]*ClassConstant // Class constants with visibility

	// Properties
	Properties       map[string]*PropertyDef // All properties (static + instance)
	StaticProperties map[string]*Value       // Static property values (shared across instances)
	DefaultProperties map[string]*Value      // Default values for instance properties

	// Methods
	Methods       map[string]*MethodDef // All methods (static + instance)
	Constructor   *MethodDef            // Constructor method (__construct)
	Destructor    *MethodDef            // Destructor method (__destruct)
	MagicMethods  map[string]*MethodDef // Magic methods (__get, __set, __call, etc.)

	// Type information
	IsInterface bool // Is this an interface?
	IsTrait     bool // Is this a trait?
	IsEnum      bool // Is this an enum? (PHP 8.1+)

	// Enum specific data
	EnumBackingType string           // Backing type for backed enums ("int" or "string")
	EnumCases       map[string]*Value // Enum cases (name => value)
}

// PropertyDef defines a class property with metadata
type PropertyDef struct {
	Name         string             // Property name
	Visibility   PropertyVisibility // public, protected, private
	IsStatic     bool               // Static vs instance property
	Type         string             // Type declaration
	HasDefault   bool               // Whether property has default value
	Default      *Value             // Default value
	IsReadOnly   bool               // readonly property (PHP 8.1+)
	Hooks        *PropertyHooks     // Property hooks (PHP 8.4+)
	DeclaringClass string           // Which class declared this property (for private props)
}

// MethodDef defines a class method with metadata
type MethodDef struct {
	Name           string             // Method name
	Visibility     PropertyVisibility // public, protected, private
	IsStatic       bool               // Static vs instance method
	IsFinal        bool               // final method (cannot be overridden)
	IsAbstract     bool               // abstract method (no implementation)
	Instructions   []interface{}      // Bytecode instructions
	NumLocals      int                // Number of local variables
	NumParams      int                // Number of parameters
	Parameters     []*ParameterDef    // Parameter definitions
	ReturnType     string             // Return type declaration
	ReturnByRef    bool               // Returns by reference
	IsConstructor  bool               // Is this __construct?
	IsDestructor   bool               // Is this __destruct?
	IsMagic        bool               // Is this a magic method?
	DeclaringClass string             // Which class declared this method
}

// ParameterDef defines a method parameter
type ParameterDef struct {
	Name         string  // Parameter name
	Type         string  // Type declaration
	HasDefault   bool    // Has default value
	Default      *Value  // Default value
	IsVariadic   bool    // Variadic parameter (...$args)
	PassedByRef  bool    // Passed by reference
	IsPromoted   bool    // Constructor promoted property (PHP 8.0+)
	Visibility   PropertyVisibility // Visibility if promoted
}

// ClassConstant represents a class constant with visibility
type ClassConstant struct {
	Name       string             // Constant name
	Value      *Value             // Constant value
	Visibility PropertyVisibility // public, protected, private (PHP 7.1+)
	IsFinal    bool               // final constant (PHP 8.1+) - cannot be overridden
}

// InterfaceEntry represents a PHP interface
type InterfaceEntry struct {
	Name          string                   // Interface name
	ParentInterfaces []*InterfaceEntry     // Extended interfaces
	Methods       map[string]*MethodDef    // Method signatures (all abstract)
	Constants     map[string]*ClassConstant // Interface constants
}

// TraitEntry represents a PHP trait
type TraitEntry struct {
	Name       string                 // Trait name
	Properties map[string]*PropertyDef // Trait properties
	Methods    map[string]*MethodDef   // Trait methods
	UsedTraits []*TraitEntry           // Traits used by this trait
}

// ============================================================================
// Object Creation and Management
// ============================================================================

// NewObjectFromClass creates a new object instance from a class entry
func NewObjectFromClass(classEntry *ClassEntry) *Object {
	obj := &Object{
		ClassName:  classEntry.Name,
		ClassEntry: classEntry,
		Properties: make(map[string]*Property),
		ObjectID:   nextObjectID(),
		IsDestroyed: false,
	}

	// Initialize instance properties with default values
	for name, propDef := range classEntry.Properties {
		if !propDef.IsStatic {
			// Copy the default value to avoid sharing between instances
			var value *Value
			if propDef.Default != nil {
				value = propDef.Default.Copy()
			}

			prop := &Property{
				Value:      value,
				Visibility: propDef.Visibility,
				IsStatic:   false,
				Type:       propDef.Type,
				HasDefault: propDef.HasDefault,
				Default:    propDef.Default, // Keep original for reference
				IsReadOnly: propDef.IsReadOnly,
				Hooks:      propDef.Hooks,
			}
			obj.Properties[name] = prop
		}
	}

	return obj
}

// NewObjectInstance creates a new object (legacy API for compatibility)
func NewObjectInstance(className string) *Object {
	return &Object{
		ClassName:  className,
		Properties: make(map[string]*Property),
		ObjectID:   nextObjectID(),
		IsDestroyed: false,
	}
}

// Global object ID counter for unique object identification
var objectIDCounter uint64 = 0

// nextObjectID generates a unique object ID
func nextObjectID() uint64 {
	objectIDCounter++
	return objectIDCounter
}

// NextObjectID generates a unique object ID (exported for use by VM)
func NextObjectID() uint64 {
	return nextObjectID()
}

// ============================================================================
// Property Access Methods
// ============================================================================

// GetProperty gets a property value with visibility checking
func (o *Object) GetProperty(name string, accessContext *ClassEntry) (*Value, bool) {
	prop, exists := o.Properties[name]
	if !exists {
		return nil, false
	}

	// Check visibility
	if !canAccessProperty(prop, accessContext, o.ClassEntry) {
		return nil, false
	}

	// Execute get hook if present
	if prop.Hooks != nil && prop.Hooks.Get != nil {
		// TODO: Execute get hook bytecode
		// For now, return the value directly
	}

	return prop.Value, true
}

// SetProperty sets a property value with visibility and readonly checking
func (o *Object) SetProperty(name string, value *Value, accessContext *ClassEntry) bool {
	prop, exists := o.Properties[name]
	if !exists {
		// Dynamic property creation (if allowed)
		prop = &Property{
			Value:      value,
			Visibility: VisibilityPublic,
			IsStatic:   false,
		}
		o.Properties[name] = prop
		return true
	}

	// Check visibility
	if !canAccessProperty(prop, accessContext, o.ClassEntry) {
		return false
	}

	// Check readonly
	if prop.IsReadOnly && prop.Value != nil {
		// Readonly properties can only be set once during construction
		return false
	}

	// Execute set hook if present
	if prop.Hooks != nil && prop.Hooks.Set != nil {
		// TODO: Execute set hook bytecode
		// For now, set the value directly
	}

	prop.Value = value
	return true
}

// canAccessProperty checks if a property can be accessed from a given context
func canAccessProperty(prop *Property, accessContext *ClassEntry, ownerClass *ClassEntry) bool {
	switch prop.Visibility {
	case VisibilityPublic:
		return true
	case VisibilityProtected:
		// Accessible from same class or subclasses
		return accessContext != nil && (accessContext == ownerClass || isSubclassOf(accessContext, ownerClass))
	case VisibilityPrivate:
		// Only accessible from the declaring class
		return accessContext != nil && accessContext == ownerClass
	default:
		return false
	}
}

// isSubclassOf checks if childClass is a subclass of parentClass
func isSubclassOf(childClass, parentClass *ClassEntry) bool {
	if childClass == nil || parentClass == nil {
		return false
	}

	current := childClass.ParentClass
	for current != nil {
		if current == parentClass {
			return true
		}
		current = current.ParentClass
	}
	return false
}

// ============================================================================
// Method Access
// ============================================================================

// GetMethod retrieves a method from the class hierarchy
func (ce *ClassEntry) GetMethod(name string) (*MethodDef, bool) {
	// Check current class
	if method, exists := ce.Methods[name]; exists {
		return method, true
	}

	// Check parent class
	if ce.ParentClass != nil {
		return ce.ParentClass.GetMethod(name)
	}

	return nil, false
}

// ImplementsInterface checks if the class implements an interface
func (ce *ClassEntry) ImplementsInterface(interfaceName string) bool {
	for _, iface := range ce.Interfaces {
		if iface.Name == interfaceName {
			return true
		}
		// Check parent interfaces
		if ifaceImplements(iface, interfaceName) {
			return true
		}
	}

	// Check parent class
	if ce.ParentClass != nil {
		return ce.ParentClass.ImplementsInterface(interfaceName)
	}

	return false
}

// ifaceImplements checks if an interface extends another interface
func ifaceImplements(iface *InterfaceEntry, interfaceName string) bool {
	for _, parent := range iface.ParentInterfaces {
		if parent.Name == interfaceName {
			return true
		}
		if ifaceImplements(parent, interfaceName) {
			return true
		}
	}
	return false
}

// ============================================================================
// Inheritance
// ============================================================================

// InheritFrom performs inheritance from a parent class
// This copies properties and methods from parent to child, performing all necessary checks
func (ce *ClassEntry) InheritFrom(parent *ClassEntry) error {
	if parent == nil {
		return nil
	}

	// Check if parent is final
	if parent.IsFinal {
		return fmt.Errorf("Class %s cannot extend final class %s", ce.Name, parent.Name)
	}

	// Set parent reference
	ce.ParentClass = parent

	// Inherit properties (skip private properties)
	for name, parentProp := range parent.Properties {
		// Skip private properties
		if parentProp.Visibility == VisibilityPrivate {
			continue
		}

		// If child doesn't override this property, inherit it
		if _, exists := ce.Properties[name]; !exists {
			// Create a copy of the property definition
			inheritedProp := &PropertyDef{
				Name:           name,
				Visibility:     parentProp.Visibility,
				IsStatic:       parentProp.IsStatic,
				Type:           parentProp.Type,
				HasDefault:     parentProp.HasDefault,
				Default:        parentProp.Default,
				IsReadOnly:     parentProp.IsReadOnly,
				Hooks:          parentProp.Hooks,
				DeclaringClass: parentProp.DeclaringClass,
			}
			if inheritedProp.DeclaringClass == "" {
				inheritedProp.DeclaringClass = parent.Name
			}
			ce.Properties[name] = inheritedProp
		}
	}

	// Inherit methods (skip private methods and constructors)
	for name, parentMethod := range parent.Methods {
		// Skip private methods
		if parentMethod.Visibility == VisibilityPrivate {
			continue
		}

		// Skip constructors - they are not inherited
		if parentMethod.IsConstructor {
			continue
		}

		// Check if child overrides this method
		if childMethod, exists := ce.Methods[name]; exists {
			// Validate the override
			if err := validateMethodOverride(parentMethod, childMethod, parent.Name, ce.Name); err != nil {
				return err
			}
		} else {
			// Inherit the method
			inheritedMethod := &MethodDef{
				Name:           name,
				Visibility:     parentMethod.Visibility,
				IsStatic:       parentMethod.IsStatic,
				IsFinal:        parentMethod.IsFinal,
				IsAbstract:     parentMethod.IsAbstract,
				Instructions:   parentMethod.Instructions,
				NumLocals:      parentMethod.NumLocals,
				NumParams:      parentMethod.NumParams,
				Parameters:     parentMethod.Parameters,
				ReturnType:     parentMethod.ReturnType,
				ReturnByRef:    parentMethod.ReturnByRef,
				IsMagic:        parentMethod.IsMagic,
				DeclaringClass: parentMethod.DeclaringClass,
			}
			if inheritedMethod.DeclaringClass == "" {
				inheritedMethod.DeclaringClass = parent.Name
			}
			ce.Methods[name] = inheritedMethod
		}
	}

	// Inherit constants
	for name, parentConst := range parent.Constants {
		// Skip private constants
		if parentConst.Visibility == VisibilityPrivate {
			continue
		}

		// If child doesn't override, inherit it
		if _, exists := ce.Constants[name]; !exists {
			ce.Constants[name] = parentConst
		}
	}

	return nil
}

// validateMethodOverride checks if a method override is valid
func validateMethodOverride(parentMethod, childMethod *MethodDef, parentClassName, childClassName string) error {
	// Cannot override final methods
	if parentMethod.IsFinal {
		return fmt.Errorf("Cannot override final method %s::%s() in %s", parentClassName, parentMethod.Name, childClassName)
	}

	// Cannot reduce visibility
	if !isVisibilityCompatible(parentMethod.Visibility, childMethod.Visibility) {
		return fmt.Errorf("Access level to %s::%s() must be %s (as in class %s) or weaker",
			childClassName, childMethod.Name, parentMethod.Visibility, parentClassName)
	}

	// Abstract methods can be implemented by concrete methods
	if parentMethod.IsAbstract && !childMethod.IsAbstract {
		// This is valid - implementing an abstract method
		return nil
	}

	return nil
}

// isVisibilityCompatible checks if child visibility is compatible with parent
// public >= protected >= private
func isVisibilityCompatible(parent, child PropertyVisibility) bool {
	// Same visibility is always OK
	if parent == child {
		return true
	}

	// Parent is protected: child can be protected or public
	if parent == VisibilityProtected && child == VisibilityPublic {
		return true
	}

	// Parent is private: child can be anything (but private methods aren't inherited anyway)
	if parent == VisibilityPrivate {
		return true
	}

	// All other cases are reducing visibility, which is not allowed
	return false
}

// HasAbstractMethods returns true if the class has any unimplemented abstract methods
func (ce *ClassEntry) HasAbstractMethods() bool {
	// Check own methods
	for _, method := range ce.Methods {
		if method.IsAbstract {
			return true
		}
	}

	// Check parent
	if ce.ParentClass != nil {
		return ce.ParentClass.HasAbstractMethods()
	}

	return false
}

// ============================================================================
// Helper Functions
// ============================================================================

// NewClassEntry creates a new class entry with default values
func NewClassEntry(name string) *ClassEntry {
	return &ClassEntry{
		Name:              name,
		Constants:         make(map[string]*ClassConstant),
		Properties:        make(map[string]*PropertyDef),
		StaticProperties:  make(map[string]*Value),
		DefaultProperties: make(map[string]*Value),
		Methods:           make(map[string]*MethodDef),
		MagicMethods:      make(map[string]*MethodDef),
		Interfaces:        make([]*InterfaceEntry, 0),
		Traits:            make([]*TraitEntry, 0),
		TraitAliases:      make(map[string]string),
		TraitPrecedence:   make(map[string]string),
	}
}

// NewInterfaceEntry creates a new interface entry
func NewInterfaceEntry(name string) *InterfaceEntry {
	return &InterfaceEntry{
		Name:             name,
		ParentInterfaces: make([]*InterfaceEntry, 0),
		Methods:          make(map[string]*MethodDef),
		Constants:        make(map[string]*ClassConstant),
	}
}

// NewTraitEntry creates a new trait entry
func NewTraitEntry(name string) *TraitEntry {
	return &TraitEntry{
		Name:       name,
		Properties: make(map[string]*PropertyDef),
		Methods:    make(map[string]*MethodDef),
		UsedTraits: make([]*TraitEntry, 0),
	}
}
