package types

import "testing"

// ============================================================================
// Integration Tests - Multiple OOP Features Combined
// ============================================================================

// Test class with inheritance, interfaces, and traits combined
func TestIntegration_CompleteClassStructure(t *testing.T) {
	// Create an interface
	logger := NewInterfaceEntry("Logger")
	logger.Methods["log"] = &MethodDef{
		Name:       "log",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "message", Type: "string"},
		},
		ReturnType: "void",
	}

	// Create a trait
	timestampable := NewTraitEntry("Timestampable")
	timestampable.Properties["createdAt"] = &PropertyDef{
		Name:       "createdAt",
		Visibility: VisibilityProtected,
		Type:       "int",
	}
	timestampable.Methods["getTimestamp"] = &MethodDef{
		Name:       "getTimestamp",
		Visibility: VisibilityPublic,
		ReturnType: "int",
	}

	// Create parent class
	baseEntity := NewClassEntry("BaseEntity")
	baseEntity.Properties["id"] = &PropertyDef{
		Name:       "id",
		Visibility: VisibilityProtected,
		Type:       "int",
	}
	baseEntity.Methods["getId"] = &MethodDef{
		Name:       "getId",
		Visibility: VisibilityPublic,
		ReturnType: "int",
	}

	// Create child class with everything
	user := NewClassEntry("User")
	user.ParentClass = baseEntity
	user.Interfaces = []*InterfaceEntry{logger}
	user.Traits = []*TraitEntry{timestampable}

	// Add User-specific properties
	user.Properties["username"] = &PropertyDef{
		Name:       "username",
		Visibility: VisibilityPrivate,
		Type:       "string",
	}
	user.Properties["email"] = &PropertyDef{
		Name:       "email",
		Visibility: VisibilityPrivate,
		Type:       "?string",
	}

	// Add magic methods
	user.MagicMethods["__construct"] = &MethodDef{
		Name:       "__construct",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "username", Type: "string"},
		},
	}
	user.MagicMethods["__toString"] = &MethodDef{
		Name:       "__toString",
		Visibility: VisibilityPublic,
		ReturnType: "string",
	}

	// Implement Logger interface
	user.Methods["log"] = &MethodDef{
		Name:       "log",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "message", Type: "string"},
		},
		ReturnType: "void",
	}

	// Apply inheritance, traits, and validate
	if err := user.InheritFrom(baseEntity); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	if err := user.ApplyTraits(); err != nil {
		t.Fatalf("Trait application failed: %v", err)
	}

	if !user.ImplementsInterface("Logger") {
		t.Fatal("User should implement Logger interface")
	}

	if err := user.ValidateMagicMethods(); err != nil {
		t.Fatalf("Magic method validation failed: %v", err)
	}

	// Verify inheritance worked
	if _, exists := user.Properties["id"]; !exists {
		t.Error("User should inherit 'id' property from BaseEntity")
	}

	if _, exists := user.Methods["getId"]; !exists {
		t.Error("User should inherit 'getId' method from BaseEntity")
	}

	// Verify trait application worked
	if _, exists := user.Properties["createdAt"]; !exists {
		t.Error("User should have 'createdAt' property from Timestampable trait")
	}

	if _, exists := user.Methods["getTimestamp"]; !exists {
		t.Error("User should have 'getTimestamp' method from Timestampable trait")
	}

	// Verify interface implementation
	if _, exists := user.Methods["log"]; !exists {
		t.Error("User should implement 'log' method from Logger interface")
	}

	// Verify magic methods
	if !user.HasMagicMethod("__construct") {
		t.Error("User should have __construct magic method")
	}

	if !user.HasMagicMethod("__toString") {
		t.Error("User should have __toString magic method")
	}

	// Verify User-specific properties
	if prop, exists := user.Properties["username"]; !exists || prop.Type != "string" {
		t.Error("User should have 'username' property with type 'string'")
	}

	if prop, exists := user.Properties["email"]; !exists || prop.Type != "?string" {
		t.Error("User should have 'email' property with type '?string'")
	}
}

// Test enum implementing interface
func TestIntegration_EnumWithInterface(t *testing.T) {
	// Create an interface
	colorInterface := NewInterfaceEntry("Colorable")
	colorInterface.Methods["toHex"] = &MethodDef{
		Name:       "toHex",
		Visibility: VisibilityPublic,
		ReturnType: "string",
	}

	// Create enum implementing interface
	color := NewEnumEntry("Color", "string")
	color.AddCase("Red", NewString("#FF0000"))
	color.AddCase("Green", NewString("#00FF00"))
	color.AddCase("Blue", NewString("#0000FF"))

	color.Interfaces = []*InterfaceEntry{colorInterface}

	// Add interface method
	color.Methods["toHex"] = &MethodDef{
		Name:       "toHex",
		Visibility: VisibilityPublic,
		ReturnType: "string",
	}

	// Validate enum
	if err := color.Validate(); err != nil {
		t.Fatalf("Enum validation failed: %v", err)
	}

	// Verify interface implementation
	if !color.ImplementsInterface("Colorable") {
		t.Fatal("Enum should implement Colorable interface")
	}

	// Verify enum cases
	if len(color.EnumCases) != 3 {
		t.Errorf("Expected 3 cases, got %d", len(color.EnumCases))
	}

	// Verify backed values
	if red := color.EnumCases["Red"]; red == nil || red.ToString() != "#FF0000" {
		t.Error("Red case should have backing value '#FF0000'")
	}
}

// Test multi-level inheritance with late static binding
func TestIntegration_MultiLevelInheritanceWithStatic(t *testing.T) {
	// Grandparent class
	grandparent := NewClassEntry("Animal")
	grandparent.Methods["whoAmI"] = &MethodDef{
		Name:           "whoAmI",
		Visibility:     VisibilityPublic,
		IsStatic:       true,
		ReturnType:     "static",
		DeclaringClass: "Animal",
	}
	grandparent.Constants["TYPE"] = &ClassConstant{
		Name:       "TYPE",
		Value:      NewString("animal"),
		Visibility: VisibilityPublic,
	}

	// Parent class
	parent := NewClassEntry("Mammal")
	parent.ParentClass = grandparent
	parent.Constants["TYPE"] = &ClassConstant{
		Name:       "TYPE",
		Value:      NewString("mammal"),
		Visibility: VisibilityPublic,
	}

	if err := parent.InheritFrom(grandparent); err != nil {
		t.Fatalf("Parent inheritance failed: %v", err)
	}

	// Child class
	child := NewClassEntry("Dog")
	child.ParentClass = parent
	child.Constants["TYPE"] = &ClassConstant{
		Name:       "TYPE",
		Value:      NewString("dog"),
		Visibility: VisibilityPublic,
	}

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Child inheritance failed: %v", err)
	}

	// Verify inheritance chain
	if child.ParentClass != parent {
		t.Error("Dog's parent should be Mammal")
	}

	if parent.ParentClass != grandparent {
		t.Error("Mammal's parent should be Animal")
	}

	// Verify static method inherited
	method, exists := child.Methods["whoAmI"]
	if !exists {
		t.Fatal("Child should inherit static method 'whoAmI'")
	}

	if !method.IsStatic {
		t.Error("Inherited method should be static")
	}

	// Verify return type is 'static'
	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsStatic {
		t.Error("Method return type should be 'static'")
	}

	// Verify late static binding for constants
	// Dog::TYPE should be "dog"
	if constant := child.Constants["TYPE"]; constant.Value.ToString() != "dog" {
		t.Errorf("Dog::TYPE should be 'dog', got '%s'", constant.Value.ToString())
	}

	// self::TYPE in Animal would be "animal"
	if constant := grandparent.Constants["TYPE"]; constant.Value.ToString() != "animal" {
		t.Errorf("Animal::TYPE should be 'animal', got '%s'", constant.Value.ToString())
	}

	// static::TYPE when called on Dog would resolve to "dog"
	value, exists := child.GetStaticConstant("TYPE", true, child)
	if !exists || value.ToString() != "dog" {
		t.Error("static::TYPE on Dog should resolve to 'dog'")
	}
}

// Test abstract class with traits and interfaces
func TestIntegration_AbstractClassWithTraitsAndInterfaces(t *testing.T) {
	// Create interface
	repository := NewInterfaceEntry("Repository")
	repository.Methods["save"] = &MethodDef{
		Name:       "save",
		Visibility: VisibilityPublic,
		ReturnType: "bool",
	}
	repository.Methods["find"] = &MethodDef{
		Name:       "find",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "id", Type: "int"},
		},
		ReturnType: "?object",
	}

	// Create trait
	cacheable := NewTraitEntry("Cacheable")
	cacheable.Methods["clearCache"] = &MethodDef{
		Name:       "clearCache",
		Visibility: VisibilityProtected,
		ReturnType: "void",
	}

	// Create abstract class
	abstractRepo := NewClassEntry("AbstractRepository")
	abstractRepo.IsAbstract = true
	abstractRepo.Interfaces = []*InterfaceEntry{repository}
	abstractRepo.Traits = []*TraitEntry{cacheable}

	// Abstract method (no implementation)
	abstractRepo.Methods["save"] = &MethodDef{
		Name:       "save",
		Visibility: VisibilityPublic,
		IsAbstract: true,
		ReturnType: "bool",
	}

	// Concrete method
	abstractRepo.Methods["find"] = &MethodDef{
		Name:       "find",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "id", Type: "int"},
		},
		ReturnType: "?object",
	}

	// Apply traits
	if err := abstractRepo.ApplyTraits(); err != nil {
		t.Fatalf("Trait application failed: %v", err)
	}

	// Verify abstract class properties
	if !abstractRepo.IsAbstract {
		t.Error("AbstractRepository should be abstract")
	}

	// Verify trait methods
	if _, exists := abstractRepo.Methods["clearCache"]; !exists {
		t.Error("AbstractRepository should have clearCache method from trait")
	}

	// Verify abstract method
	saveMethod, exists := abstractRepo.Methods["save"]
	if !exists {
		t.Fatal("AbstractRepository should have save method")
	}
	if !saveMethod.IsAbstract {
		t.Error("save method should be abstract")
	}

	// Create concrete implementation
	userRepo := NewClassEntry("UserRepository")
	userRepo.ParentClass = abstractRepo

	// Implement abstract method
	userRepo.Methods["save"] = &MethodDef{
		Name:       "save",
		Visibility: VisibilityPublic,
		ReturnType: "bool",
	}

	if err := userRepo.InheritFrom(abstractRepo); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Verify concrete class has all methods
	if _, exists := userRepo.Methods["save"]; !exists {
		t.Error("UserRepository should have save method")
	}

	if _, exists := userRepo.Methods["find"]; !exists {
		t.Error("UserRepository should inherit find method")
	}

	if _, exists := userRepo.Methods["clearCache"]; !exists {
		t.Error("UserRepository should inherit clearCache from trait")
	}

	// Verify interface implementation
	if !userRepo.ImplementsInterface("Repository") {
		t.Fatal("UserRepository should implement Repository interface")
	}
}

// Test readonly class with typed properties
func TestIntegration_ReadonlyClassWithTypedProperties(t *testing.T) {
	// Create readonly class (PHP 8.2+)
	valueObject := NewClassEntry("Point")
	valueObject.IsReadOnly = true

	// Readonly properties must have types
	valueObject.Properties["x"] = &PropertyDef{
		Name:       "x",
		Visibility: VisibilityPublic,
		Type:       "float",
		IsReadOnly: true,
	}

	valueObject.Properties["y"] = &PropertyDef{
		Name:       "y",
		Visibility: VisibilityPublic,
		Type:       "float",
		IsReadOnly: true,
	}

	// Constructor to initialize readonly properties
	valueObject.MagicMethods["__construct"] = &MethodDef{
		Name:       "__construct",
		Visibility: VisibilityPublic,
		NumParams:  2,
		Parameters: []*ParameterDef{
			{Name: "x", Type: "float"},
			{Name: "y", Type: "float"},
		},
	}

	// Method with return type
	valueObject.Methods["distanceFromOrigin"] = &MethodDef{
		Name:       "distanceFromOrigin",
		Visibility: VisibilityPublic,
		ReturnType: "float",
	}

	// Verify readonly class
	if !valueObject.IsReadOnly {
		t.Error("Point should be readonly class")
	}

	// Verify all properties are readonly
	for _, prop := range valueObject.Properties {
		if !prop.IsReadOnly {
			t.Errorf("Property %s should be readonly", prop.Name)
		}
	}

	// Validate readonly properties have types
	for _, prop := range valueObject.Properties {
		if err := ValidateReadonlyProperty(prop); err != nil {
			t.Errorf("Readonly property validation failed: %v", err)
		}
	}

	// Verify type checking
	xProp := valueObject.Properties["x"]
	if err := ValidatePropertyValue(xProp, NewFloat(3.14)); err != nil {
		t.Errorf("float value should be valid for float property: %v", err)
	}

	if err := ValidatePropertyValue(xProp, NewString("invalid")); err == nil {
		t.Error("string value should be invalid for float property")
	}
}

// Test trait with precedence and aliasing
func TestIntegration_TraitPrecedenceAndAliasing(t *testing.T) {
	// Create two traits with conflicting methods
	trait1 := NewTraitEntry("Logger")
	trait1.Methods["log"] = &MethodDef{
		Name:       "log",
		Visibility: VisibilityPublic,
	}

	trait2 := NewTraitEntry("FileLogger")
	trait2.Methods["log"] = &MethodDef{
		Name:       "log",
		Visibility: VisibilityPublic,
	}

	// Create class using both traits
	app := NewClassEntry("Application")
	app.Traits = []*TraitEntry{trait1, trait2}

	// Resolve conflict with precedence (Logger::log insteadof FileLogger)
	app.TraitPrecedence = make(map[string]string)
	app.TraitPrecedence["log"] = "Logger"

	// Alias FileLogger's log method with protected visibility
	app.TraitAliases = make(map[string]string)
	app.TraitAliases["fileLog"] = "FileLogger::log:protected"

	// Apply traits
	if err := app.ApplyTraits(); err != nil {
		t.Fatalf("Trait application failed: %v", err)
	}

	// Verify precedence worked
	logMethod, exists := app.Methods["log"]
	if !exists {
		t.Fatal("Application should have log method")
	}

	// log should come from Logger trait
	if logMethod.DeclaringClass != "Logger" {
		t.Errorf("log method should be from Logger trait, got %s", logMethod.DeclaringClass)
	}

	// Verify alias worked
	fileLogMethod, exists := app.Methods["fileLog"]
	if !exists {
		t.Fatal("Application should have fileLog aliased method")
	}

	// Alias should have new visibility
	if fileLogMethod.Visibility != VisibilityProtected {
		t.Error("fileLog should have protected visibility")
	}

	// Original name should come from FileLogger
	if fileLogMethod.DeclaringClass != "FileLogger" {
		t.Errorf("fileLog should be from FileLogger trait, got %s", fileLogMethod.DeclaringClass)
	}
}

// Test reflection on complex class
func TestIntegration_ReflectionOnComplexClass(t *testing.T) {
	// Create a complex class
	class := NewClassEntry("ComplexClass")
	class.IsFinal = true

	// Add properties with various modifiers
	class.Properties["publicProp"] = &PropertyDef{
		Name:       "publicProp",
		Visibility: VisibilityPublic,
		Type:       "int",
		Default:    NewInt(10),
	}

	class.Properties["protectedStatic"] = &PropertyDef{
		Name:       "protectedStatic",
		Visibility: VisibilityProtected,
		IsStatic:   true,
		Type:       "string",
		Default:    NewString("test"),
	}

	class.Properties["privateReadonly"] = &PropertyDef{
		Name:       "privateReadonly",
		Visibility: VisibilityPrivate,
		IsReadOnly: true,
		Type:       "float",
	}

	// Add methods with various modifiers
	class.Methods["publicMethod"] = &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityPublic,
		ReturnType: "void",
	}

	class.Methods["protectedStatic"] = &MethodDef{
		Name:       "protectedStatic",
		Visibility: VisibilityProtected,
		IsStatic:   true,
		ReturnType: "int",
	}

	class.Methods["privateFinal"] = &MethodDef{
		Name:       "privateFinal",
		Visibility: VisibilityPrivate,
		IsFinal:    true,
		ReturnType: "bool",
	}

	// Add constants
	class.Constants["PUBLIC_CONST"] = &ClassConstant{
		Name:       "PUBLIC_CONST",
		Value:      NewInt(100),
		Visibility: VisibilityPublic,
	}

	class.Constants["PRIVATE_CONST"] = &ClassConstant{
		Name:       "PRIVATE_CONST",
		Value:      NewString("secret"),
		Visibility: VisibilityPrivate,
	}

	// Test reflection methods
	if class.GetName() != "ComplexClass" {
		t.Error("GetName() should return 'ComplexClass'")
	}

	// Get all property names
	propNames := class.GetPropertyNames()
	if len(propNames) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(propNames))
	}

	// Get properties by visibility
	publicProps := class.GetPropertiesByVisibility(VisibilityPublic)
	if len(publicProps) != 1 {
		t.Errorf("Expected 1 public property, got %d", len(publicProps))
	}

	protectedProps := class.GetPropertiesByVisibility(VisibilityProtected)
	if len(protectedProps) != 1 {
		t.Errorf("Expected 1 protected property, got %d", len(protectedProps))
	}

	// Get all method names
	methodNames := class.GetMethodNames()
	if len(methodNames) != 3 {
		t.Errorf("Expected 3 methods, got %d", len(methodNames))
	}

	// Get methods by visibility
	publicMethods := class.GetMethodsByVisibility(VisibilityPublic)
	if len(publicMethods) != 1 {
		t.Errorf("Expected 1 public method, got %d", len(publicMethods))
	}

	// Get all constant names
	constNames := class.GetConstantNames()
	if len(constNames) != 2 {
		t.Errorf("Expected 2 constants, got %d", len(constNames))
	}

	// Get class modifiers
	modifiers := class.GetModifiers()
	if modifiers&0x01 == 0 {
		t.Error("Class should have final modifier bit set")
	}

	// Get method modifiers
	staticMethod := class.Methods["protectedStatic"]
	methodMod := staticMethod.GetModifiers()
	if methodMod&0x01 == 0 {
		t.Error("Method should have static modifier bit set")
	}
	if methodMod&0x200 == 0 {
		t.Error("Method should have protected visibility bit set")
	}

	// Get property modifiers
	staticProp := class.Properties["protectedStatic"]
	propMod := staticProp.GetModifiers()
	if propMod&0x01 == 0 {
		t.Error("Property should have static modifier bit set")
	}

	readonlyProp := class.Properties["privateReadonly"]
	readonlyMod := readonlyProp.GetModifiers()
	if readonlyMod&0x02 == 0 {
		t.Error("Property should have readonly modifier bit set")
	}
}

// Test type compatibility across inheritance
func TestIntegration_TypeCompatibilityInheritance(t *testing.T) {
	// Create parent class
	animal := NewClassEntry("Animal")
	animal.Methods["makeSound"] = &MethodDef{
		Name:       "makeSound",
		Visibility: VisibilityPublic,
		ReturnType: "string",
	}

	// Create child class with covariant return type
	dog := NewClassEntry("Dog")
	dog.ParentClass = animal

	// In PHP 7.4+, child can have more specific return type
	// But for built-in types, must match exactly
	dog.Methods["makeSound"] = &MethodDef{
		Name:       "makeSound",
		Visibility: VisibilityPublic,
		ReturnType: "string", // Must match parent
	}

	if err := dog.InheritFrom(animal); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Test method with object return type
	animal.Methods["getParent"] = &MethodDef{
		Name:       "getParent",
		Visibility: VisibilityPublic,
		ReturnType: "?Animal",
	}

	dog.Methods["getParent"] = &MethodDef{
		Name:       "getParent",
		Visibility: VisibilityPublic,
		ReturnType: "?Dog", // Covariant - more specific
	}

	// Verify type parsing
	animalReturnType := ParseType(animal.Methods["getParent"].ReturnType)
	if !animalReturnType.IsNullable {
		t.Error("?Animal should be nullable")
	}
	if animalReturnType.BaseType != "Animal" {
		t.Error("Base type should be Animal")
	}

	dogReturnType := ParseType(dog.Methods["getParent"].ReturnType)
	if !dogReturnType.IsNullable {
		t.Error("?Dog should be nullable")
	}
	if dogReturnType.BaseType != "Dog" {
		t.Error("Base type should be Dog")
	}
}

// Test union types and mixed type
func TestIntegration_UnionAndMixedTypes(t *testing.T) {
	class := NewClassEntry("FlexibleClass")

	// Property with union type
	class.Properties["flexible"] = &PropertyDef{
		Name:       "flexible",
		Visibility: VisibilityPublic,
		Type:       "int|string|null",
	}

	// Method with union parameter
	class.Methods["process"] = &MethodDef{
		Name:       "process",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{Name: "value", Type: "int|float"},
		},
		ReturnType: "string|bool",
	}

	// Property with mixed type
	class.Properties["anything"] = &PropertyDef{
		Name:       "anything",
		Visibility: VisibilityPublic,
		Type:       "mixed",
	}

	// Parse union type
	flexibleType := ParseType(class.Properties["flexible"].Type)
	if !flexibleType.IsUnion {
		t.Error("int|string|null should be a union type")
	}
	if len(flexibleType.UnionTypes) != 3 {
		t.Errorf("Expected 3 union types, got %d", len(flexibleType.UnionTypes))
	}

	// Parse mixed type
	mixedType := ParseType(class.Properties["anything"].Type)
	if !mixedType.IsBuiltin {
		t.Error("mixed should be a builtin type")
	}
	if mixedType.BaseType != "mixed" {
		t.Error("Base type should be mixed")
	}

	// Test type compatibility
	if !IsTypeCompatible("int|string", "int") {
		t.Error("int|string should accept int")
	}

	if !IsTypeCompatible("int|string", "string") {
		t.Error("int|string should accept string")
	}

	if IsTypeCompatible("int|string", "float") {
		t.Error("int|string should not accept float")
	}

	if !IsTypeCompatible("mixed", "int") {
		t.Error("mixed should accept int")
	}

	if !IsTypeCompatible("mixed", "object") {
		t.Error("mixed should accept object")
	}
}
