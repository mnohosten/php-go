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

// ValidateInterfaceImplementation validates that the class properly implements all its interfaces
func (ce *ClassEntry) ValidateInterfaceImplementation() error {
	for _, iface := range ce.Interfaces {
		if err := ce.validateSingleInterface(iface); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleInterface validates implementation of a single interface (including parent interfaces)
func (ce *ClassEntry) validateSingleInterface(iface *InterfaceEntry) error {
	// Check all methods in this interface
	for methodName, ifaceMethod := range iface.Methods {
		// Get the implementation from the class
		classMethod, exists := ce.GetMethod(methodName)
		if !exists {
			return fmt.Errorf("Class %s must implement method %s() from interface %s",
				ce.Name, methodName, iface.Name)
		}

		// Validate the implementation
		if err := validateInterfaceMethodImplementation(ifaceMethod, classMethod, iface.Name, ce.Name); err != nil {
			return err
		}
	}

	// Recursively check parent interfaces
	for _, parentIface := range iface.ParentInterfaces {
		if err := ce.validateSingleInterface(parentIface); err != nil {
			return err
		}
	}

	return nil
}

// validateInterfaceMethodImplementation checks if a class method properly implements an interface method
func validateInterfaceMethodImplementation(ifaceMethod, classMethod *MethodDef, ifaceName, className string) error {
	// Check parameter count
	if classMethod.NumParams != ifaceMethod.NumParams {
		return fmt.Errorf("Method %s::%s() must have %d parameter(s) to match interface %s",
			className, classMethod.Name, ifaceMethod.NumParams, ifaceName)
	}

	// Interface methods are always public, implementation must be public too
	if classMethod.Visibility != VisibilityPublic {
		return fmt.Errorf("Method %s::%s() must be public to implement interface %s",
			className, classMethod.Name, ifaceName)
	}

	// Method must not be abstract in concrete class (unless class is abstract)
	// This is checked elsewhere, not in interface validation

	return nil
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

	// Check if parent is an enum
	if parent.IsEnum {
		return fmt.Errorf("Class %s cannot extend enum %s", ce.Name, parent.Name)
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

	// Inherit magic methods (except constructor)
	for name, parentMagic := range parent.MagicMethods {
		// Skip constructor - constructors are not inherited
		if name == "__construct" {
			continue
		}

		// If child doesn't override, inherit it
		if _, exists := ce.MagicMethods[name]; !exists {
			inheritedMagic := &MethodDef{
				Name:           name,
				Visibility:     parentMagic.Visibility,
				IsStatic:       parentMagic.IsStatic,
				IsFinal:        parentMagic.IsFinal,
				IsAbstract:     parentMagic.IsAbstract,
				Instructions:   parentMagic.Instructions,
				NumLocals:      parentMagic.NumLocals,
				NumParams:      parentMagic.NumParams,
				Parameters:     parentMagic.Parameters,
				ReturnType:     parentMagic.ReturnType,
				ReturnByRef:    parentMagic.ReturnByRef,
				IsMagic:        true,
				DeclaringClass: parentMagic.DeclaringClass,
			}
			if inheritedMagic.DeclaringClass == "" {
				inheritedMagic.DeclaringClass = parent.Name
			}
			ce.MagicMethods[name] = inheritedMagic
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

// ============================================================================
// Trait Application
// ============================================================================

// ApplyTraits applies all traits to the class
// This implements PHP's trait composition rules:
// - Class methods override trait methods
// - Trait methods override inherited methods
// - Conflicts must be resolved with precedence or cause error
func (ce *ClassEntry) ApplyTraits() error {
	if len(ce.Traits) == 0 {
		return nil
	}

	// Track methods from traits to detect conflicts
	traitMethods := make(map[string][]*traitMethodSource)
	traitProperties := make(map[string][]*traitPropertySource)

	// Collect all methods and properties from all traits
	for _, trait := range ce.Traits {
		// First, ensure the trait itself has applied its used traits
		if err := trait.ApplyUsedTraits(); err != nil {
			return err
		}

		// Collect methods from this trait
		for methodName, method := range trait.Methods {
			source := &traitMethodSource{
				TraitName: trait.Name,
				Method:    method,
			}
			traitMethods[methodName] = append(traitMethods[methodName], source)
		}

		// Collect properties from this trait
		for propName, prop := range trait.Properties {
			source := &traitPropertySource{
				TraitName: trait.Name,
				Property:  prop,
			}
			traitProperties[propName] = append(traitProperties[propName], source)
		}
	}

	// Apply trait methods to class
	for methodName, sources := range traitMethods {
		// Skip if class already defines this method (class methods take precedence)
		if existingMethod, exists := ce.Methods[methodName]; exists {
			// Check if this is from the class itself, not inherited
			if existingMethod.DeclaringClass == ce.Name || existingMethod.DeclaringClass == "" {
				continue // Class method takes precedence
			}
		}

		// Check for conflicts
		if len(sources) > 1 {
			// Multiple traits define this method - need precedence resolution
			resolvedTrait, hasPrecedence := ce.TraitPrecedence[methodName]
			if !hasPrecedence {
				// No precedence defined - error
				traitNames := make([]string, len(sources))
				for i, src := range sources {
					traitNames[i] = src.TraitName
				}
				return fmt.Errorf("Trait method conflict: method '%s' exists in multiple traits (%s)",
					methodName, joinStrings(traitNames, ", "))
			}

			// Find the method from the resolved trait
			var selectedMethod *MethodDef
			for _, src := range sources {
				if src.TraitName == resolvedTrait {
					selectedMethod = src.Method
					break
				}
			}

			if selectedMethod == nil {
				return fmt.Errorf("Precedence resolution failed: trait %s does not define method %s",
					resolvedTrait, methodName)
			}

			// Copy the selected method to the class
			copiedMethod := copyMethodDef(selectedMethod)
			if copiedMethod.DeclaringClass == "" {
				copiedMethod.DeclaringClass = resolvedTrait
			}
			ce.Methods[methodName] = copiedMethod
		} else {
			// No conflict, just copy the method
			copiedMethod := copyMethodDef(sources[0].Method)
			if copiedMethod.DeclaringClass == "" {
				copiedMethod.DeclaringClass = sources[0].TraitName
			}
			ce.Methods[methodName] = copiedMethod
		}
	}

	// Apply trait properties to class
	for propName, sources := range traitProperties {
		// Skip if class already defines this property
		if _, exists := ce.Properties[propName]; exists {
			continue
		}

		// Check for conflicts
		if len(sources) > 1 {
			// Check if all sources have compatible definitions
			// In PHP, properties can conflict if they have different default values
			firstProp := sources[0].Property
			for i := 1; i < len(sources); i++ {
				if !arePropertiesCompatible(firstProp, sources[i].Property) {
					return fmt.Errorf("Trait property conflict: property '%s' has incompatible definitions in multiple traits",
						propName)
				}
			}
		}

		// Copy the property to the class
		ce.Properties[propName] = copyPropertyDef(sources[0].Property)
	}

	// Apply trait aliases
	for aliasName, aliasSpec := range ce.TraitAliases {
		// Parse alias spec: "TraitName::method" or "TraitName::method:visibility"
		var traitName, methodName, visibility string

		// Simple parsing - in real implementation this would be more robust
		parts := splitString(aliasSpec, "::")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid alias specification: %s", aliasSpec)
		}

		traitName = parts[0]
		methodParts := splitString(parts[1], ":")
		methodName = methodParts[0]
		if len(methodParts) > 1 {
			visibility = methodParts[1]
		}

		// Find the trait
		var sourceTrait *TraitEntry
		for _, trait := range ce.Traits {
			if trait.Name == traitName {
				sourceTrait = trait
				break
			}
		}

		if sourceTrait == nil {
			return fmt.Errorf("Alias error: trait %s not used by class", traitName)
		}

		// Find the method in the trait
		sourceMethod, exists := sourceTrait.Methods[methodName]
		if !exists {
			return fmt.Errorf("Alias error: method %s not found in trait %s", methodName, traitName)
		}

		// Create the alias
		aliasedMethod := copyMethodDef(sourceMethod)
		aliasedMethod.Name = aliasName

		// Set declaring class to the source trait
		if aliasedMethod.DeclaringClass == "" {
			aliasedMethod.DeclaringClass = traitName
		}

		// Change visibility if specified
		if visibility != "" {
			switch visibility {
			case "public":
				aliasedMethod.Visibility = VisibilityPublic
			case "protected":
				aliasedMethod.Visibility = VisibilityProtected
			case "private":
				aliasedMethod.Visibility = VisibilityPrivate
			}
		}

		ce.Methods[aliasName] = aliasedMethod
	}

	return nil
}

// ApplyUsedTraits applies traits used by this trait (trait composition)
func (te *TraitEntry) ApplyUsedTraits() error {
	if len(te.UsedTraits) == 0 {
		return nil
	}

	// Track methods from used traits to detect conflicts
	traitMethods := make(map[string][]*traitMethodSource)
	traitProperties := make(map[string][]*traitPropertySource)

	// Collect all methods and properties from used traits
	for _, usedTrait := range te.UsedTraits {
		// Recursively apply traits
		if err := usedTrait.ApplyUsedTraits(); err != nil {
			return err
		}

		// Collect methods
		for methodName, method := range usedTrait.Methods {
			source := &traitMethodSource{
				TraitName: usedTrait.Name,
				Method:    method,
			}
			traitMethods[methodName] = append(traitMethods[methodName], source)
		}

		// Collect properties
		for propName, prop := range usedTrait.Properties {
			source := &traitPropertySource{
				TraitName: usedTrait.Name,
				Property:  prop,
			}
			traitProperties[propName] = append(traitProperties[propName], source)
		}
	}

	// Apply methods to this trait
	for methodName, sources := range traitMethods {
		// Skip if this trait already defines this method (own methods take precedence)
		if _, exists := te.Methods[methodName]; exists {
			continue
		}

		// For trait composition, conflicts are an error (no precedence mechanism here)
		if len(sources) > 1 {
			traitNames := make([]string, len(sources))
			for i, src := range sources {
				traitNames[i] = src.TraitName
			}
			return fmt.Errorf("Trait method conflict in %s: method '%s' exists in multiple used traits (%s)",
				te.Name, methodName, joinStrings(traitNames, ", "))
		}

		// Copy the method
		te.Methods[methodName] = copyMethodDef(sources[0].Method)
	}

	// Apply properties to this trait
	for propName, sources := range traitProperties {
		// Skip if this trait already defines this property
		if _, exists := te.Properties[propName]; exists {
			continue
		}

		// Check for conflicts
		if len(sources) > 1 {
			firstProp := sources[0].Property
			for i := 1; i < len(sources); i++ {
				if !arePropertiesCompatible(firstProp, sources[i].Property) {
					return fmt.Errorf("Trait property conflict in %s: property '%s' has incompatible definitions",
						te.Name, propName)
				}
			}
		}

		// Copy the property
		te.Properties[propName] = copyPropertyDef(sources[0].Property)
	}

	return nil
}

// traitMethodSource tracks which trait a method came from
type traitMethodSource struct {
	TraitName string
	Method    *MethodDef
}

// traitPropertySource tracks which trait a property came from
type traitPropertySource struct {
	TraitName string
	Property  *PropertyDef
}

// copyMethodDef creates a copy of a method definition
func copyMethodDef(method *MethodDef) *MethodDef {
	return &MethodDef{
		Name:           method.Name,
		Visibility:     method.Visibility,
		IsStatic:       method.IsStatic,
		IsFinal:        method.IsFinal,
		IsAbstract:     method.IsAbstract,
		Instructions:   method.Instructions,
		NumLocals:      method.NumLocals,
		NumParams:      method.NumParams,
		Parameters:     method.Parameters,
		ReturnType:     method.ReturnType,
		ReturnByRef:    method.ReturnByRef,
		IsConstructor:  method.IsConstructor,
		IsDestructor:   method.IsDestructor,
		IsMagic:        method.IsMagic,
		DeclaringClass: method.DeclaringClass,
	}
}

// copyPropertyDef creates a copy of a property definition
func copyPropertyDef(prop *PropertyDef) *PropertyDef {
	return &PropertyDef{
		Name:           prop.Name,
		Visibility:     prop.Visibility,
		IsStatic:       prop.IsStatic,
		Type:           prop.Type,
		HasDefault:     prop.HasDefault,
		Default:        prop.Default,
		IsReadOnly:     prop.IsReadOnly,
		Hooks:          prop.Hooks,
		DeclaringClass: prop.DeclaringClass,
	}
}

// arePropertiesCompatible checks if two property definitions are compatible
// Properties conflict if they have different default values
func arePropertiesCompatible(prop1, prop2 *PropertyDef) bool {
	// Check if both have default values (either via HasDefault flag or non-nil Default)
	hasDef1 := prop1.HasDefault || prop1.Default != nil
	hasDef2 := prop2.HasDefault || prop2.Default != nil

	// If both have defaults, check if they're equal
	if hasDef1 && hasDef2 {
		if prop1.Default == nil && prop2.Default == nil {
			return true
		}
		if prop1.Default == nil || prop2.Default == nil {
			return false
		}
		// Check if default values are equal
		return prop1.Default.Equals(prop2.Default)
	}

	// If one has default and other doesn't, they're incompatible
	if hasDef1 != hasDef2 {
		return false
	}

	// Neither has defaults - compatible
	return true
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// Helper function to split strings
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	result := []string{}
	start := 0

	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])

	return result
}

// ============================================================================
// Enum Support (PHP 8.1+)
// ============================================================================

// NewEnumEntry creates a new enum entry
// backingType can be "" (pure enum), "int", or "string"
func NewEnumEntry(name string, backingType string) *ClassEntry {
	return &ClassEntry{
		Name:              name,
		IsEnum:            true,
		EnumBackingType:   backingType,
		EnumCases:         make(map[string]*Value),
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

// AddCase adds a case to the enum
// For pure enums, value should be nil
// For backed enums, value should be an int or string Value
func (ce *ClassEntry) AddCase(name string, value *Value) {
	if !ce.IsEnum {
		return
	}
	ce.EnumCases[name] = value
}

// Validate validates the enum definition
func (ce *ClassEntry) Validate() error {
	if !ce.IsEnum {
		return nil // Not an enum, skip validation
	}

	// Check backing type is valid
	if ce.EnumBackingType != "" && ce.EnumBackingType != "int" && ce.EnumBackingType != "string" {
		return fmt.Errorf("Enum backing type must be 'int' or 'string', got '%s'", ce.EnumBackingType)
	}

	// Enums cannot extend other classes
	if ce.ParentClass != nil {
		return fmt.Errorf("Enum %s cannot extend other classes", ce.Name)
	}

	// Enums cannot have instance properties (only constants and static properties allowed)
	for name, prop := range ce.Properties {
		if !prop.IsStatic {
			return fmt.Errorf("Enum %s cannot have instance property '%s'", ce.Name, name)
		}
	}

	// Note: We don't need to check for duplicate case names because EnumCases is a map,
	// and maps automatically prevent duplicates (adding the same key twice overwrites the value)

	// For backed enums, validate all cases have correct type
	if ce.EnumBackingType != "" {
		for caseName, caseValue := range ce.EnumCases {
			if caseValue == nil {
				return fmt.Errorf("Backed enum %s case '%s' must have a value", ce.Name, caseName)
			}

			// Check type matches backing type
			if ce.EnumBackingType == "int" {
				if caseValue.Type() != TypeInt {
					return fmt.Errorf("Backed enum %s case '%s' must have int value, got %v",
						ce.Name, caseName, caseValue.Type())
				}
			} else if ce.EnumBackingType == "string" {
				if caseValue.Type() != TypeString {
					return fmt.Errorf("Backed enum %s case '%s' must have string value, got %v",
						ce.Name, caseName, caseValue.Type())
				}
			}
		}
	} else {
		// Pure enums should not have backing values
		for caseName, caseValue := range ce.EnumCases {
			if caseValue != nil {
				return fmt.Errorf("Pure enum %s case '%s' should not have a backing value", ce.Name, caseName)
			}
		}
	}

	return nil
}

// GetCases returns all enum cases as a slice of case names
func (ce *ClassEntry) GetCases() []string {
	if !ce.IsEnum {
		return nil
	}

	cases := make([]string, 0, len(ce.EnumCases))
	for caseName := range ce.EnumCases {
		cases = append(cases, caseName)
	}
	return cases
}

// From returns the case name for a given backing value (backed enums only)
// Returns error if the value doesn't exist
func (ce *ClassEntry) From(value *Value) (string, error) {
	if !ce.IsEnum {
		return "", fmt.Errorf("From() can only be called on enums")
	}

	if ce.EnumBackingType == "" {
		return "", fmt.Errorf("From() is only available on backed enums")
	}

	// Find case with matching value
	for caseName, caseValue := range ce.EnumCases {
		if caseValue != nil && caseValue.Equals(value) {
			return caseName, nil
		}
	}

	// Value not found
	return "", fmt.Errorf("Value %v is not a valid backing value for enum %s", value, ce.Name)
}

// TryFrom returns the case name for a given backing value (backed enums only)
// Returns empty string if the value doesn't exist (instead of error)
func (ce *ClassEntry) TryFrom(value *Value) (string, error) {
	if !ce.IsEnum {
		return "", fmt.Errorf("TryFrom() can only be called on enums")
	}

	if ce.EnumBackingType == "" {
		return "", fmt.Errorf("TryFrom() is only available on backed enums")
	}

	// Find case with matching value
	for caseName, caseValue := range ce.EnumCases {
		if caseValue != nil && caseValue.Equals(value) {
			return caseName, nil
		}
	}

	// Value not found - return empty string (not error)
	return "", nil
}

// ============================================================================
// Type Checking and Validation
// ============================================================================

// TypeInfo represents parsed type information
type TypeInfo struct {
	BaseType    string   // The base type (int, string, ClassName, etc.)
	IsNullable  bool     // true if type starts with ?
	IsUnion     bool     // true if type contains |
	UnionTypes  []string // List of types in union (for int|string)
	IsBuiltin   bool     // true for built-in types (int, string, etc.)
	IsClass     bool     // true for class types
	IsSelf      bool     // true for 'self' type
	IsParent    bool     // true for 'parent' type
	IsStatic    bool     // true for 'static' type
}

// ParseType parses a type string and returns type information
func ParseType(typeStr string) *TypeInfo {
	if typeStr == "" {
		return &TypeInfo{}
	}

	info := &TypeInfo{}

	// Check for nullable type
	if len(typeStr) > 0 && typeStr[0] == '?' {
		info.IsNullable = true
		typeStr = typeStr[1:]
	}

	// Check for union type
	if containsString(typeStr, "|") {
		info.IsUnion = true
		info.UnionTypes = splitString(typeStr, "|")
		if len(info.UnionTypes) > 0 {
			info.BaseType = info.UnionTypes[0]
		}
	} else {
		info.BaseType = typeStr
	}

	// Check for built-in types
	builtinTypes := map[string]bool{
		"int":      true,
		"string":   true,
		"float":    true,
		"bool":     true,
		"array":    true,
		"object":   true,
		"callable": true,
		"iterable": true,
		"mixed":    true,
		"void":     true,
		"never":    true,
		"null":     true,
		"false":    true,
		"true":     true,
	}

	info.IsBuiltin = builtinTypes[info.BaseType]

	// Check for special types
	switch info.BaseType {
	case "self":
		info.IsSelf = true
	case "parent":
		info.IsParent = true
	case "static":
		info.IsStatic = true
	default:
		if !info.IsBuiltin {
			info.IsClass = true
		}
	}

	return info
}

// IsTypeCompatible checks if a value of valueType can be assigned to a variable of expectedType
func IsTypeCompatible(expectedType, valueType string) bool {
	if expectedType == valueType {
		return true
	}

	expectedInfo := ParseType(expectedType)
	valueInfo := ParseType(valueType)

	// mixed accepts any type
	if expectedInfo.BaseType == "mixed" {
		return true
	}

	// If expected type is nullable, it accepts the non-nullable version
	if expectedInfo.IsNullable && expectedInfo.BaseType == valueInfo.BaseType {
		return true
	}

	// If expected type is nullable, it also accepts null
	if expectedInfo.IsNullable && valueType == "null" {
		return true
	}

	// If expected is union type, check if value matches any of the union members
	if expectedInfo.IsUnion {
		for _, unionType := range expectedInfo.UnionTypes {
			if unionType == valueType {
				return true
			}
		}
	}

	// iterable accepts array
	if expectedInfo.BaseType == "iterable" && valueType == "array" {
		return true
	}

	// object accepts any class type
	if expectedInfo.BaseType == "object" && valueInfo.IsClass {
		return true
	}

	return false
}

// ValidatePropertyValue validates that a value matches a property's type
func ValidatePropertyValue(prop *PropertyDef, value *Value) error {
	if prop.Type == "" {
		return nil // No type constraint
	}

	typeInfo := ParseType(prop.Type)

	// Allow null for nullable types
	if value == nil || value.IsNull() {
		if typeInfo.IsNullable {
			return nil
		}
		return fmt.Errorf("Property %s cannot be null (type: %s)", prop.Name, prop.Type)
	}

	// Get value type
	valueTypeStr := getValueTypeString(value)

	// Check compatibility
	if !IsTypeCompatible(prop.Type, valueTypeStr) {
		return fmt.Errorf("Property %s expects type %s, got %s", prop.Name, prop.Type, valueTypeStr)
	}

	return nil
}

// getValueTypeString returns the type string for a Value
func getValueTypeString(v *Value) string {
	if v == nil {
		return "null"
	}

	switch v.Type() {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	case TypeArray:
		return "array"
	case TypeObject:
		// For objects, return the class name
		if obj, ok := v.data.(*Object); ok {
			return obj.ClassName
		}
		return "object"
	case TypeNull:
		return "null"
	default:
		return "mixed"
	}
}

// ValidateReadonlyProperty validates that a readonly property has a type hint
func ValidateReadonlyProperty(prop *PropertyDef) error {
	if !prop.IsReadOnly {
		return nil
	}

	if prop.Type == "" {
		return fmt.Errorf("Readonly property %s must have a type", prop.Name)
	}

	return nil
}

// ValidateReturnTypeCovariance checks if child return type is covariant with parent
// Return types are covariant: child can return a subtype of parent's return type
func ValidateReturnTypeCovariance(parentMethod, childMethod *MethodDef, parentType, childType string) error {
	if parentMethod.ReturnType == "" {
		return nil // No type constraint
	}

	if childMethod.ReturnType == "" {
		// Child has no return type but parent does - this is usually allowed
		return nil
	}

	// Same type is always OK
	if parentMethod.ReturnType == childMethod.ReturnType {
		return nil
	}

	parentInfo := ParseType(parentMethod.ReturnType)
	childInfo := ParseType(childMethod.ReturnType)

	// void return type must match exactly
	if parentInfo.BaseType == "void" && childInfo.BaseType != "void" {
		return fmt.Errorf("Return type must be void to match parent")
	}

	// Check if childType is a subclass of parentType (covariance)
	// For now, we'll accept any class type as potentially covariant
	// In a full implementation, we'd check the class hierarchy
	if parentInfo.IsClass && childInfo.IsClass {
		// Assume valid - actual check would require class hierarchy
		return nil
	}

	// Built-in types must match exactly (no covariance for primitives)
	if parentInfo.IsBuiltin && childInfo.IsBuiltin && parentInfo.BaseType != childInfo.BaseType {
		return fmt.Errorf("Return type %s is not covariant with parent return type %s",
			childMethod.ReturnType, parentMethod.ReturnType)
	}

	return nil
}

// ValidateParameterTypeContravariance checks if child parameter type is contravariant with parent
// Parameter types are contravariant: child can accept a supertype of parent's parameter type
func ValidateParameterTypeContravariance(parentParam, childParam *ParameterDef, parentType, childType string) error {
	if parentParam.Type == "" {
		return nil // No type constraint
	}

	if childParam.Type == "" {
		// Child has no type but parent does - this is usually not allowed
		return fmt.Errorf("Child parameter must have type %s to match parent", parentParam.Type)
	}

	// Same type is always OK
	if parentParam.Type == childParam.Type {
		return nil
	}

	parentInfo := ParseType(parentParam.Type)
	childInfo := ParseType(childParam.Type)

	// Check if childType is a superclass of parentType (contravariance)
	// For now, we'll accept any class type as potentially contravariant
	// In a full implementation, we'd check the class hierarchy
	if parentInfo.IsClass && childInfo.IsClass {
		// Assume valid - actual check would require class hierarchy
		return nil
	}

	// Built-in types must match exactly
	if parentInfo.IsBuiltin && childInfo.IsBuiltin && parentInfo.BaseType != childInfo.BaseType {
		return fmt.Errorf("Parameter type %s is not compatible with parent parameter type %s",
			childParam.Type, parentParam.Type)
	}

	return nil
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findInString(s, substr) >= 0
}

// findInString finds the index of substr in s, returns -1 if not found
func findInString(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ============================================================================
// Late Static Binding
// ============================================================================

// ResolveStaticClass resolves the class to use for static:: calls
// calledClass is the class that was actually called (from the VM frame)
// currentClass is the class where the method is defined
func ResolveStaticClass(calledClass, currentClass *ClassEntry) *ClassEntry {
	// static:: refers to the called class, not the defining class
	if calledClass != nil {
		return calledClass
	}
	// Fallback to current class if no called class info
	return currentClass
}

// ResolveSelfClass resolves the class to use for self:: calls
// self:: always refers to the class where the code is written
func ResolveSelfClass(currentClass *ClassEntry) *ClassEntry {
	return currentClass
}

// ResolveParentClass resolves the class to use for parent:: calls
// parent:: refers to the parent class of the current class
func ResolveParentClass(currentClass *ClassEntry) *ClassEntry {
	if currentClass != nil {
		return currentClass.ParentClass
	}
	return nil
}

// GetCalledClassName returns the name of the called class for get_called_class()
// This is used in late static binding contexts
func GetCalledClassName(calledClass *ClassEntry) string {
	if calledClass != nil {
		return calledClass.Name
	}
	return ""
}

// IsStaticContext checks if we're in a static method context
func (ce *ClassEntry) IsStaticContext(methodName string) bool {
	if method, exists := ce.Methods[methodName]; exists {
		return method.IsStatic
	}
	return false
}

// GetStaticProperty retrieves a static property value
// This is separate from instance property access
func (ce *ClassEntry) GetStaticProperty(name string) (*Value, bool) {
	if value, exists := ce.StaticProperties[name]; exists {
		return value, true
	}

	// Check parent class for inherited static properties
	if ce.ParentClass != nil {
		if prop, exists := ce.Properties[name]; exists {
			if prop.IsStatic && prop.Visibility != VisibilityPrivate {
				// Get from parent's static properties
				return ce.ParentClass.GetStaticProperty(name)
			}
		}
	}

	return nil, false
}

// SetStaticProperty sets a static property value
func (ce *ClassEntry) SetStaticProperty(name string, value *Value) bool {
	// Check if property is defined and static
	if prop, exists := ce.Properties[name]; exists {
		if prop.IsStatic {
			ce.StaticProperties[name] = value
			return true
		}
	}

	// Check parent class
	if ce.ParentClass != nil {
		if prop, exists := ce.Properties[name]; exists {
			if prop.IsStatic && prop.Visibility != VisibilityPrivate {
				return ce.ParentClass.SetStaticProperty(name, value)
			}
		}
	}

	return false
}

// GetStaticConstant retrieves a class constant
// static::CONST uses late binding, self::CONST uses early binding
func (ce *ClassEntry) GetStaticConstant(name string, useLateBind bool, calledClass *ClassEntry) (*Value, bool) {
	// For late binding (static::), check the called class first
	if useLateBind && calledClass != nil {
		if constant, exists := calledClass.Constants[name]; exists {
			return constant.Value, true
		}
	}

	// For self:: or when called class doesn't have it, use current class
	if constant, exists := ce.Constants[name]; exists {
		return constant.Value, true
	}

	// Check parent class
	if ce.ParentClass != nil {
		return ce.ParentClass.GetStaticConstant(name, false, nil)
	}

	return nil, false
}

// ============================================================================
// Reflection API
// ============================================================================

// GetName returns the class name
func (ce *ClassEntry) GetName() string {
	return ce.Name
}

// GetParentClassName returns the parent class name, or empty string if no parent
func (ce *ClassEntry) GetParentClassName() string {
	if ce.ParentClass != nil {
		return ce.ParentClass.Name
	}
	return ""
}

// GetInterfaceNames returns a list of interface names implemented by this class
func (ce *ClassEntry) GetInterfaceNames() []string {
	names := make([]string, len(ce.Interfaces))
	for i, iface := range ce.Interfaces {
		names[i] = iface.Name
	}
	return names
}

// GetTraitNames returns a list of trait names used by this class
func (ce *ClassEntry) GetTraitNames() []string {
	names := make([]string, len(ce.Traits))
	for i, trait := range ce.Traits {
		names[i] = trait.Name
	}
	return names
}

// GetMethodNames returns a list of all method names
func (ce *ClassEntry) GetMethodNames() []string {
	names := make([]string, 0, len(ce.Methods))
	for name := range ce.Methods {
		names = append(names, name)
	}
	return names
}

// GetPropertyNames returns a list of all property names
func (ce *ClassEntry) GetPropertyNames() []string {
	names := make([]string, 0, len(ce.Properties))
	for name := range ce.Properties {
		names = append(names, name)
	}
	return names
}

// GetConstantNames returns a list of all constant names
func (ce *ClassEntry) GetConstantNames() []string {
	names := make([]string, 0, len(ce.Constants))
	for name := range ce.Constants {
		names = append(names, name)
	}
	return names
}

// GetMethodsByVisibility returns methods filtered by visibility
func (ce *ClassEntry) GetMethodsByVisibility(visibility PropertyVisibility) []*MethodDef {
	methods := make([]*MethodDef, 0)
	for _, method := range ce.Methods {
		if method.Visibility == visibility {
			methods = append(methods, method)
		}
	}
	return methods
}

// GetPropertiesByVisibility returns properties filtered by visibility
func (ce *ClassEntry) GetPropertiesByVisibility(visibility PropertyVisibility) []*PropertyDef {
	properties := make([]*PropertyDef, 0)
	for _, prop := range ce.Properties {
		if prop.Visibility == visibility {
			properties = append(properties, prop)
		}
	}
	return properties
}

// GetStaticMethods returns all static methods
func (ce *ClassEntry) GetStaticMethods() []*MethodDef {
	methods := make([]*MethodDef, 0)
	for _, method := range ce.Methods {
		if method.IsStatic {
			methods = append(methods, method)
		}
	}
	return methods
}

// GetStaticProperties returns all static properties
func (ce *ClassEntry) GetStaticProperties() []*PropertyDef {
	properties := make([]*PropertyDef, 0)
	for _, prop := range ce.Properties {
		if prop.IsStatic {
			properties = append(properties, prop)
		}
	}
	return properties
}

// IsInstantiable returns true if the class can be instantiated
func (ce *ClassEntry) IsInstantiable() bool {
	// Abstract classes, interfaces, and traits cannot be instantiated
	return !ce.IsAbstract && !ce.IsInterface && !ce.IsTrait
}

// GetConstructor returns the constructor method, or nil if none
func (ce *ClassEntry) GetConstructor() *MethodDef {
	return ce.Constructor
}

// GetDestructor returns the destructor method, or nil if none
func (ce *ClassEntry) GetDestructor() *MethodDef {
	return ce.Destructor
}

// HasConstructor returns true if the class has a constructor
func (ce *ClassEntry) HasConstructor() bool {
	return ce.Constructor != nil
}

// HasDestructor returns true if the class has a destructor
func (ce *ClassEntry) HasDestructor() bool {
	return ce.Destructor != nil
}

// GetNamespaceName returns the namespace of the class
func (ce *ClassEntry) GetNamespaceName() string {
	return ce.Namespace
}

// GetShortName returns the short name (without namespace) of the class
func (ce *ClassEntry) GetShortName() string {
	if ce.ShortName != "" {
		return ce.ShortName
	}
	return ce.Name
}

// GetFileName returns the file where the class was defined
func (ce *ClassEntry) GetFileName() string {
	return ce.FileName
}

// GetModifiers returns a bitmask of class modifiers
func (ce *ClassEntry) GetModifiers() uint32 {
	var modifiers uint32 = 0
	if ce.IsFinal {
		modifiers |= 0x01 // Final
	}
	if ce.IsAbstract {
		modifiers |= 0x02 // Abstract
	}
	if ce.IsReadOnly {
		modifiers |= 0x04 // Readonly (PHP 8.2+)
	}
	return modifiers
}

// GetMethodModifiers returns a bitmask of method modifiers
func (m *MethodDef) GetModifiers() uint32 {
	var modifiers uint32 = 0
	if m.IsStatic {
		modifiers |= 0x01 // Static
	}
	if m.IsFinal {
		modifiers |= 0x02 // Final
	}
	if m.IsAbstract {
		modifiers |= 0x04 // Abstract
	}
	// Add visibility
	switch m.Visibility {
	case VisibilityPublic:
		modifiers |= 0x100 // Public
	case VisibilityProtected:
		modifiers |= 0x200 // Protected
	case VisibilityPrivate:
		modifiers |= 0x400 // Private
	}
	return modifiers
}

// GetPropertyModifiers returns a bitmask of property modifiers
func (p *PropertyDef) GetModifiers() uint32 {
	var modifiers uint32 = 0
	if p.IsStatic {
		modifiers |= 0x01 // Static
	}
	if p.IsReadOnly {
		modifiers |= 0x02 // Readonly
	}
	// Add visibility
	switch p.Visibility {
	case VisibilityPublic:
		modifiers |= 0x100 // Public
	case VisibilityProtected:
		modifiers |= 0x200 // Protected
	case VisibilityPrivate:
		modifiers |= 0x400 // Private
	}
	return modifiers
}

// ============================================================================
// Magic Methods
// ============================================================================

// HasMagicMethod checks if the class has a specific magic method
func (ce *ClassEntry) HasMagicMethod(name string) bool {
	// Check own magic methods
	if _, exists := ce.MagicMethods[name]; exists {
		return true
	}

	// Check if it's in Methods map (constructor/destructor)
	if name == "__construct" && ce.Constructor != nil {
		return true
	}
	if name == "__destruct" && ce.Destructor != nil {
		return true
	}

	// Check regular methods for magic methods
	if method, exists := ce.Methods[name]; exists && method.IsMagic {
		return true
	}

	// Check parent class
	if ce.ParentClass != nil {
		return ce.ParentClass.HasMagicMethod(name)
	}

	return false
}

// GetMagicMethod retrieves a magic method from the class hierarchy
func (ce *ClassEntry) GetMagicMethod(name string) *MethodDef {
	// Check own magic methods map
	if method, exists := ce.MagicMethods[name]; exists {
		return method
	}

	// Check constructor/destructor
	if name == "__construct" && ce.Constructor != nil {
		return ce.Constructor
	}
	if name == "__destruct" && ce.Destructor != nil {
		return ce.Destructor
	}

	// Check regular methods for magic methods
	if method, exists := ce.Methods[name]; exists && method.IsMagic {
		return method
	}

	// Check parent class
	if ce.ParentClass != nil {
		return ce.ParentClass.GetMagicMethod(name)
	}

	return nil
}

// ValidateMagicMethods validates magic method definitions
func (ce *ClassEntry) ValidateMagicMethods() error {
	// List of magic methods that must be public (except __construct/__destruct)
	publicOnlyMagic := map[string]bool{
		"__get":        true,
		"__set":        true,
		"__isset":      true,
		"__unset":      true,
		"__call":       true,
		"__callStatic": true,
		"__toString":   true,
		"__invoke":     true,
		"__clone":      true,
		"__debugInfo":  true,
		"__serialize":  true,
		"__unserialize": true,
		"__sleep":      true,
		"__wakeup":     true,
	}

	// Magic methods that must be static
	mustBeStatic := map[string]bool{
		"__callStatic": true,
		"__set_state":  true,
	}

	// Magic methods that must NOT be static
	mustBeInstance := map[string]bool{
		"__get":        true,
		"__set":        true,
		"__isset":      true,
		"__unset":      true,
		"__call":       true,
		"__toString":   true,
		"__invoke":     true,
		"__clone":      true,
		"__debugInfo":  true,
		"__serialize":  true,
		"__unserialize": true,
		"__sleep":      true,
		"__wakeup":     true,
	}

	// Validate magic methods in MagicMethods map
	for name, method := range ce.MagicMethods {
		// Check visibility
		if publicOnlyMagic[name] && method.Visibility != VisibilityPublic {
			return fmt.Errorf("Magic method %s must be public in class %s", name, ce.Name)
		}

		// Check static requirement
		if mustBeStatic[name] && !method.IsStatic {
			return fmt.Errorf("Magic method %s must be static in class %s", name, ce.Name)
		}

		// Check instance requirement
		if mustBeInstance[name] && method.IsStatic {
			return fmt.Errorf("Magic method %s cannot be static in class %s", name, ce.Name)
		}

		// Validate parameter counts for specific magic methods
		switch name {
		case "__get", "__isset", "__unset":
			if method.NumParams != 1 {
				return fmt.Errorf("Magic method %s must have exactly 1 parameter in class %s", name, ce.Name)
			}
		case "__set":
			if method.NumParams != 2 {
				return fmt.Errorf("Magic method %s must have exactly 2 parameters in class %s", name, ce.Name)
			}
		case "__call", "__callStatic":
			if method.NumParams != 2 {
				return fmt.Errorf("Magic method %s must have exactly 2 parameters in class %s", name, ce.Name)
			}
		case "__toString", "__clone", "__debugInfo", "__serialize", "__sleep", "__wakeup":
			if method.NumParams != 0 {
				return fmt.Errorf("Magic method %s must have no parameters in class %s", name, ce.Name)
			}
		case "__unserialize":
			if method.NumParams != 1 {
				return fmt.Errorf("Magic method %s must have exactly 1 parameter in class %s", name, ce.Name)
			}
		}
	}

	// Validate constructor/destructor if present
	if ce.Constructor != nil {
		// Constructor can have any visibility (private for singletons)
		// No specific validation needed
	}

	if ce.Destructor != nil && ce.Destructor.Visibility != VisibilityPublic {
		// Destructor should be public, but PHP allows other visibilities
		// We'll be lenient here
	}

	return nil
}
