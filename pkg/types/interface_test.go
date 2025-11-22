package types

import "testing"

// ============================================================================
// Interface Definition Tests
// ============================================================================

func TestInterface_BasicDefinition(t *testing.T) {
	// Create interface with methods
	iface := NewInterfaceEntry("TestInterface")
	iface.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
		IsAbstract: true, // Interface methods are always abstract
		NumParams:  0,
	}
	iface.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPublic,
		IsAbstract: true,
		NumParams:  1,
	}

	if len(iface.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(iface.Methods))
	}

	if iface.Name != "TestInterface" {
		t.Errorf("Expected name 'TestInterface', got '%s'", iface.Name)
	}
}

func TestInterface_WithConstants(t *testing.T) {
	iface := NewInterfaceEntry("ConfigInterface")
	iface.Constants["MAX_SIZE"] = &ClassConstant{
		Name:       "MAX_SIZE",
		Value:      NewInt(1000),
		Visibility: VisibilityPublic,
	}
	iface.Constants["DEFAULT_NAME"] = &ClassConstant{
		Name:       "DEFAULT_NAME",
		Value:      NewString("default"),
		Visibility: VisibilityPublic,
	}

	if len(iface.Constants) != 2 {
		t.Errorf("Expected 2 constants, got %d", len(iface.Constants))
	}

	maxSize := iface.Constants["MAX_SIZE"]
	if maxSize.Value.ToInt() != 1000 {
		t.Errorf("Expected MAX_SIZE to be 1000, got %d", maxSize.Value.ToInt())
	}
}

// ============================================================================
// Interface Inheritance Tests
// ============================================================================

func TestInterface_ExtendsSingleInterface(t *testing.T) {
	// Create parent interface
	parent := NewInterfaceEntry("ParentInterface")
	parent.Methods["parentMethod"] = &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Create child interface that extends parent
	child := NewInterfaceEntry("ChildInterface")
	child.ParentInterfaces = []*InterfaceEntry{parent}
	child.Methods["childMethod"] = &MethodDef{
		Name:       "childMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	if len(child.ParentInterfaces) != 1 {
		t.Errorf("Expected 1 parent interface, got %d", len(child.ParentInterfaces))
	}

	if child.ParentInterfaces[0].Name != "ParentInterface" {
		t.Error("Parent interface not set correctly")
	}
}

func TestInterface_ExtendsMultipleInterfaces(t *testing.T) {
	// Create two parent interfaces
	iface1 := NewInterfaceEntry("Interface1")
	iface1.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	iface2 := NewInterfaceEntry("Interface2")
	iface2.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Child extends both
	child := NewInterfaceEntry("ChildInterface")
	child.ParentInterfaces = []*InterfaceEntry{iface1, iface2}

	if len(child.ParentInterfaces) != 2 {
		t.Errorf("Expected 2 parent interfaces, got %d", len(child.ParentInterfaces))
	}
}

// ============================================================================
// Class Implementation Tests
// ============================================================================

func TestInterface_ClassImplementsSingle(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("TestInterface")
	iface.Methods["testMethod"] = &MethodDef{
		Name:       "testMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
		NumParams:  0,
	}

	// Create class that implements interface
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}
	class.Methods["testMethod"] = &MethodDef{
		Name:       "testMethod",
		Visibility: VisibilityPublic,
		IsAbstract: false, // Concrete implementation
		NumParams:  0,
	}

	// Validate implementation
	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Errorf("Interface validation failed: %v", err)
	}
}

func TestInterface_ClassImplementsMultiple(t *testing.T) {
	// Create two interfaces
	iface1 := NewInterfaceEntry("Interface1")
	iface1.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	iface2 := NewInterfaceEntry("Interface2")
	iface2.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Create class implementing both
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface1, iface2}
	class.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
	}
	class.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPublic,
	}

	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Errorf("Interface validation failed: %v", err)
	}

	if len(class.Interfaces) != 2 {
		t.Errorf("Expected 2 interfaces, got %d", len(class.Interfaces))
	}
}

func TestInterface_MissingMethodImplementation(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("TestInterface")
	iface.Methods["requiredMethod"] = &MethodDef{
		Name:       "requiredMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Create class without implementing required method
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}

	err := class.ValidateInterfaceImplementation()
	if err == nil {
		t.Fatal("Expected error for missing method implementation")
	}

	expectedMsg := "Class TestClass must implement method requiredMethod() from interface TestInterface"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInterface_WrongParameterCount(t *testing.T) {
	// Create interface with method expecting 2 parameters
	iface := NewInterfaceEntry("TestInterface")
	iface.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
		IsAbstract: true,
		NumParams:  2,
	}

	// Create class with method having wrong parameter count
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}
	class.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
		NumParams:  1, // Wrong count
	}

	err := class.ValidateInterfaceImplementation()
	if err == nil {
		t.Fatal("Expected error for wrong parameter count")
	}
}

func TestInterface_InsufficientVisibility(t *testing.T) {
	// Create interface (all interface methods are public)
	iface := NewInterfaceEntry("TestInterface")
	iface.Methods["publicMethod"] = &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Create class with protected method (less visibility)
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}
	class.Methods["publicMethod"] = &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityProtected, // Should be public
	}

	err := class.ValidateInterfaceImplementation()
	if err == nil {
		t.Fatal("Expected error for reduced visibility")
	}
}

// ============================================================================
// Interface Checking Tests
// ============================================================================

func TestInterface_ImplementsInterface(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("TestInterface")

	// Create class implementing interface
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}

	// Check if class implements interface
	if !class.ImplementsInterface("TestInterface") {
		t.Error("Class should implement TestInterface")
	}

	if class.ImplementsInterface("NonExistentInterface") {
		t.Error("Class should not implement non-existent interface")
	}
}

func TestInterface_ImplementsInterfaceViaParent(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("TestInterface")

	// Create parent class implementing interface
	parent := NewClassEntry("ParentClass")
	parent.Interfaces = []*InterfaceEntry{iface}

	// Create child class
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent

	// Child should implement interface through parent
	if !child.ImplementsInterface("TestInterface") {
		t.Error("Child should implement interface through parent")
	}
}

func TestInterface_ImplementsExtendedInterface(t *testing.T) {
	// Create parent interface
	parentIface := NewInterfaceEntry("ParentInterface")

	// Create child interface extending parent
	childIface := NewInterfaceEntry("ChildInterface")
	childIface.ParentInterfaces = []*InterfaceEntry{parentIface}

	// Create class implementing child interface
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{childIface}

	// Should implement both child and parent interfaces
	if !class.ImplementsInterface("ChildInterface") {
		t.Error("Class should implement ChildInterface")
	}

	if !class.ImplementsInterface("ParentInterface") {
		t.Error("Class should implement ParentInterface (through ChildInterface)")
	}
}

// ============================================================================
// Interface Constant Access Tests
// ============================================================================

func TestInterface_AccessConstants(t *testing.T) {
	iface := NewInterfaceEntry("ConfigInterface")
	iface.Constants["VERSION"] = &ClassConstant{
		Name:       "VERSION",
		Value:      NewString("1.0.0"),
		Visibility: VisibilityPublic,
	}

	class := NewClassEntry("Config")
	class.Interfaces = []*InterfaceEntry{iface}

	// Class should be able to access interface constants
	// This is tested through the instanceof and interface implementation
	if len(iface.Constants) != 1 {
		t.Error("Interface should have 1 constant")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestInterface_EmptyInterface(t *testing.T) {
	// Empty interface is valid in PHP
	iface := NewInterfaceEntry("EmptyInterface")

	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{iface}

	err := class.ValidateInterfaceImplementation()
	if err != nil {
		t.Errorf("Empty interface should be valid: %v", err)
	}
}

func TestInterface_MultiLevelInheritance(t *testing.T) {
	// Grandparent interface
	grandparent := NewInterfaceEntry("GrandparentInterface")
	grandparent.Methods["grandparentMethod"] = &MethodDef{
		Name:       "grandparentMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Parent interface extends grandparent
	parent := NewInterfaceEntry("ParentInterface")
	parent.ParentInterfaces = []*InterfaceEntry{grandparent}
	parent.Methods["parentMethod"] = &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Child interface extends parent
	child := NewInterfaceEntry("ChildInterface")
	child.ParentInterfaces = []*InterfaceEntry{parent}

	// Class implements child interface
	class := NewClassEntry("TestClass")
	class.Interfaces = []*InterfaceEntry{child}

	// Should implement all three interfaces
	if !class.ImplementsInterface("ChildInterface") {
		t.Error("Should implement ChildInterface")
	}
	if !class.ImplementsInterface("ParentInterface") {
		t.Error("Should implement ParentInterface")
	}
	if !class.ImplementsInterface("GrandparentInterface") {
		t.Error("Should implement GrandparentInterface")
	}
}
