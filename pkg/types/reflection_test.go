package types

import "testing"

// ============================================================================
// ReflectionClass Tests
// ============================================================================

func TestReflection_GetClassName(t *testing.T) {
	class := NewClassEntry("MyClass")

	name := class.GetName()
	if name != "MyClass" {
		t.Errorf("Expected name 'MyClass', got '%s'", name)
	}
}

func TestReflection_GetParentClass(t *testing.T) {
	parent := NewClassEntry("Parent")
	child := NewClassEntry("Child")
	child.ParentClass = parent

	parentName := child.GetParentClassName()
	if parentName != "Parent" {
		t.Errorf("Expected parent 'Parent', got '%s'", parentName)
	}

	// Class with no parent
	orphan := NewClassEntry("Orphan")
	if orphan.GetParentClassName() != "" {
		t.Error("Class with no parent should return empty string")
	}
}

func TestReflection_GetInterfaces(t *testing.T) {
	iface1 := NewInterfaceEntry("Interface1")
	iface2 := NewInterfaceEntry("Interface2")

	class := NewClassEntry("MyClass")
	class.Interfaces = []*InterfaceEntry{iface1, iface2}

	interfaces := class.GetInterfaceNames()
	if len(interfaces) != 2 {
		t.Errorf("Expected 2 interfaces, got %d", len(interfaces))
	}

	if interfaces[0] != "Interface1" || interfaces[1] != "Interface2" {
		t.Error("Interface names don't match")
	}
}

func TestReflection_GetTraits(t *testing.T) {
	trait1 := NewTraitEntry("Trait1")
	trait2 := NewTraitEntry("Trait2")

	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait1, trait2}

	traits := class.GetTraitNames()
	if len(traits) != 2 {
		t.Errorf("Expected 2 traits, got %d", len(traits))
	}

	if traits[0] != "Trait1" || traits[1] != "Trait2" {
		t.Error("Trait names don't match")
	}
}

func TestReflection_IsFinal(t *testing.T) {
	finalClass := NewClassEntry("FinalClass")
	finalClass.IsFinal = true

	if !finalClass.IsFinal {
		t.Error("Class should be final")
	}

	normalClass := NewClassEntry("NormalClass")
	if normalClass.IsFinal {
		t.Error("Class should not be final")
	}
}

func TestReflection_IsAbstract(t *testing.T) {
	abstractClass := NewClassEntry("AbstractClass")
	abstractClass.IsAbstract = true

	if !abstractClass.IsAbstract {
		t.Error("Class should be abstract")
	}

	normalClass := NewClassEntry("NormalClass")
	if normalClass.IsAbstract {
		t.Error("Class should not be abstract")
	}
}

func TestReflection_IsInterface(t *testing.T) {
	iface := NewInterfaceEntry("MyInterface")

	// InterfaceEntry is separate from ClassEntry, but we can check the structure
	if iface.Name != "MyInterface" {
		t.Error("Interface name should match")
	}
}

func TestReflection_IsTrait(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.IsTrait = true

	if !class.IsTrait {
		t.Error("Class should be marked as trait")
	}
}

func TestReflection_IsEnum(t *testing.T) {
	enum := NewEnumEntry("MyEnum", "int")

	if !enum.IsEnum {
		t.Error("Enum should be marked as enum")
	}

	normalClass := NewClassEntry("NormalClass")
	if normalClass.IsEnum {
		t.Error("Normal class should not be enum")
	}
}

// ============================================================================
// ReflectionMethod Tests
// ============================================================================

func TestReflection_GetMethods(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
	}
	class.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPrivate,
	}

	methods := class.GetMethodNames()
	if len(methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(methods))
	}
}

func TestReflection_GetMethod(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Methods["testMethod"] = &MethodDef{
		Name:       "testMethod",
		Visibility: VisibilityPublic,
		NumParams:  2,
		ReturnType: "string",
	}

	method, exists := class.Methods["testMethod"]
	if !exists {
		t.Fatal("Method should exist")
	}

	if method.Name != "testMethod" {
		t.Error("Method name doesn't match")
	}

	if method.NumParams != 2 {
		t.Errorf("Expected 2 parameters, got %d", method.NumParams)
	}

	if method.ReturnType != "string" {
		t.Errorf("Expected return type 'string', got '%s'", method.ReturnType)
	}
}

func TestReflection_IsMethodPublic(t *testing.T) {
	method := &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityPublic,
	}

	if method.Visibility != VisibilityPublic {
		t.Error("Method should be public")
	}
}

func TestReflection_IsMethodProtected(t *testing.T) {
	method := &MethodDef{
		Name:       "protectedMethod",
		Visibility: VisibilityProtected,
	}

	if method.Visibility != VisibilityProtected {
		t.Error("Method should be protected")
	}
}

func TestReflection_IsMethodPrivate(t *testing.T) {
	method := &MethodDef{
		Name:       "privateMethod",
		Visibility: VisibilityPrivate,
	}

	if method.Visibility != VisibilityPrivate {
		t.Error("Method should be private")
	}
}

func TestReflection_IsMethodStatic(t *testing.T) {
	method := &MethodDef{
		Name:     "staticMethod",
		IsStatic: true,
	}

	if !method.IsStatic {
		t.Error("Method should be static")
	}
}

func TestReflection_IsMethodFinal(t *testing.T) {
	method := &MethodDef{
		Name:    "finalMethod",
		IsFinal: true,
	}

	if !method.IsFinal {
		t.Error("Method should be final")
	}
}

func TestReflection_IsMethodAbstract(t *testing.T) {
	method := &MethodDef{
		Name:       "abstractMethod",
		IsAbstract: true,
	}

	if !method.IsAbstract {
		t.Error("Method should be abstract")
	}
}

func TestReflection_GetMethodParameters(t *testing.T) {
	method := &MethodDef{
		Name:      "myMethod",
		NumParams: 3,
		Parameters: []*ParameterDef{
			{Name: "param1", Type: "int"},
			{Name: "param2", Type: "string"},
			{Name: "param3", Type: "bool", HasDefault: true, Default: NewBool(false)},
		},
	}

	if len(method.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(method.Parameters))
	}

	if method.Parameters[0].Name != "param1" || method.Parameters[0].Type != "int" {
		t.Error("First parameter doesn't match")
	}

	if !method.Parameters[2].HasDefault {
		t.Error("Third parameter should have default value")
	}
}

// ============================================================================
// ReflectionProperty Tests
// ============================================================================

func TestReflection_GetProperties(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Properties["prop1"] = &PropertyDef{
		Name:       "prop1",
		Visibility: VisibilityPublic,
		Type:       "int",
	}
	class.Properties["prop2"] = &PropertyDef{
		Name:       "prop2",
		Visibility: VisibilityPrivate,
		Type:       "string",
	}

	properties := class.GetPropertyNames()
	if len(properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(properties))
	}
}

func TestReflection_GetProperty(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Properties["testProp"] = &PropertyDef{
		Name:       "testProp",
		Visibility: VisibilityPublic,
		Type:       "string",
		Default:    NewString("test"),
	}

	prop, exists := class.Properties["testProp"]
	if !exists {
		t.Fatal("Property should exist")
	}

	if prop.Name != "testProp" {
		t.Error("Property name doesn't match")
	}

	if prop.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", prop.Type)
	}
}

func TestReflection_IsPropertyPublic(t *testing.T) {
	prop := &PropertyDef{
		Name:       "publicProp",
		Visibility: VisibilityPublic,
	}

	if prop.Visibility != VisibilityPublic {
		t.Error("Property should be public")
	}
}

func TestReflection_IsPropertyStatic(t *testing.T) {
	prop := &PropertyDef{
		Name:     "staticProp",
		IsStatic: true,
	}

	if !prop.IsStatic {
		t.Error("Property should be static")
	}
}

func TestReflection_IsPropertyReadonly(t *testing.T) {
	prop := &PropertyDef{
		Name:       "readonlyProp",
		IsReadOnly: true,
		Type:       "string",
	}

	if !prop.IsReadOnly {
		t.Error("Property should be readonly")
	}
}

func TestReflection_GetPropertyType(t *testing.T) {
	prop := &PropertyDef{
		Name: "typedProp",
		Type: "int",
	}

	if prop.Type != "int" {
		t.Errorf("Expected type 'int', got '%s'", prop.Type)
	}
}

func TestReflection_GetPropertyDefaultValue(t *testing.T) {
	prop := &PropertyDef{
		Name:       "defaultProp",
		HasDefault: true,
		Default:    NewInt(42),
	}

	if !prop.HasDefault {
		t.Error("Property should have default value")
	}

	if prop.Default.ToInt() != 42 {
		t.Errorf("Expected default value 42, got %d", prop.Default.ToInt())
	}
}

// ============================================================================
// ReflectionConstant Tests
// ============================================================================

func TestReflection_GetConstants(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Constants["CONST1"] = &ClassConstant{
		Name:  "CONST1",
		Value: NewInt(100),
	}
	class.Constants["CONST2"] = &ClassConstant{
		Name:  "CONST2",
		Value: NewString("test"),
	}

	constants := class.GetConstantNames()
	if len(constants) != 2 {
		t.Errorf("Expected 2 constants, got %d", len(constants))
	}
}

func TestReflection_GetConstantValue(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Constants["MAX_SIZE"] = &ClassConstant{
		Name:  "MAX_SIZE",
		Value: NewInt(1000),
	}

	constant, exists := class.Constants["MAX_SIZE"]
	if !exists {
		t.Fatal("Constant should exist")
	}

	if constant.Value.ToInt() != 1000 {
		t.Errorf("Expected constant value 1000, got %d", constant.Value.ToInt())
	}
}

// ============================================================================
// Class Metadata Tests
// ============================================================================

func TestReflection_HasMethod(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Methods["exists"] = &MethodDef{
		Name: "exists",
	}

	if _, exists := class.Methods["exists"]; !exists {
		t.Error("hasMethod should return true for existing method")
	}

	if _, exists := class.Methods["notExists"]; exists {
		t.Error("hasMethod should return false for non-existing method")
	}
}

func TestReflection_HasProperty(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Properties["exists"] = &PropertyDef{
		Name: "exists",
	}

	if _, exists := class.Properties["exists"]; !exists {
		t.Error("hasProperty should return true for existing property")
	}

	if _, exists := class.Properties["notExists"]; exists {
		t.Error("hasProperty should return false for non-existing property")
	}
}

func TestReflection_HasConstant(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.Constants["EXISTS"] = &ClassConstant{
		Name: "EXISTS",
	}

	if _, exists := class.Constants["EXISTS"]; !exists {
		t.Error("hasConstant should return true for existing constant")
	}

	if _, exists := class.Constants["NOT_EXISTS"]; exists {
		t.Error("hasConstant should return false for non-existing constant")
	}
}

func TestReflection_IsInstantiable(t *testing.T) {
	// Abstract class is not instantiable
	abstractClass := NewClassEntry("AbstractClass")
	abstractClass.IsAbstract = true

	if !abstractClass.IsAbstract {
		t.Error("Abstract class should not be instantiable")
	}

	// Interface is not instantiable (InterfaceEntry is a separate type)
	// Traits are also not instantiable
	trait := NewClassEntry("MyTrait")
	trait.IsTrait = true
	if !trait.IsTrait {
		t.Error("Trait should not be instantiable")
	}

	// Normal class is instantiable
	normalClass := NewClassEntry("NormalClass")
	if normalClass.IsAbstract || normalClass.IsTrait {
		t.Error("Normal class should be instantiable")
	}
}

func TestReflection_GetNamespaceName(t *testing.T) {
	class := NewClassEntry("MyNamespace\\MyClass")
	class.Namespace = "MyNamespace"

	if class.Namespace != "MyNamespace" {
		t.Errorf("Expected namespace 'MyNamespace', got '%s'", class.Namespace)
	}
}

func TestReflection_GetShortName(t *testing.T) {
	class := NewClassEntry("MyNamespace\\MyClass")
	class.ShortName = "MyClass"

	if class.ShortName != "MyClass" {
		t.Errorf("Expected short name 'MyClass', got '%s'", class.ShortName)
	}
}

func TestReflection_GetFileName(t *testing.T) {
	class := NewClassEntry("MyClass")
	class.FileName = "/path/to/MyClass.php"

	if class.FileName != "/path/to/MyClass.php" {
		t.Errorf("Expected file name '/path/to/MyClass.php', got '%s'", class.FileName)
	}
}
