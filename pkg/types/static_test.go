package types

import "testing"

// ============================================================================
// Late Static Binding Basics
// ============================================================================

func TestStatic_CalledClassTracking(t *testing.T) {
	// Test that ClassEntry can track the called class for static:: resolution
	class := NewClassEntry("TestClass")

	// In a real scenario, the VM would track the called class
	// This test verifies the structure exists
	if class.Name != "TestClass" {
		t.Errorf("Expected class name TestClass, got %s", class.Name)
	}
}

func TestStatic_SelfVsStatic(t *testing.T) {
	// Create parent class with static method
	parent := NewClassEntry("Parent")
	parent.Methods["test"] = &MethodDef{
		Name:       "test",
		Visibility: VisibilityPublic,
		IsStatic:   true,
	}

	// Create child class extending parent
	child := NewClassEntry("Child")
	child.ParentClass = parent

	// Inherit methods
	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// The key difference:
	// - self:: would always resolve to "Parent"
	// - static:: would resolve to "Child" when called on Child

	// This is tested at runtime in the VM, but we verify the structure supports it
	if child.ParentClass != parent {
		t.Error("Child should have parent reference for static:: resolution")
	}
}

// ============================================================================
// Static Method Resolution
// ============================================================================

func TestStatic_StaticMethodInheritance(t *testing.T) {
	// Parent with static method
	parent := NewClassEntry("Parent")
	parent.Methods["whoAmI"] = &MethodDef{
		Name:       "whoAmI",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "string",
	}

	// Child inherits static method
	child := NewClassEntry("Child")
	child.ParentClass = parent

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Child should have the static method
	method, exists := child.Methods["whoAmI"]
	if !exists {
		t.Fatal("Child should inherit static method")
	}

	if !method.IsStatic {
		t.Error("Inherited method should be static")
	}
}

func TestStatic_OverrideStaticMethod(t *testing.T) {
	// Parent with static method
	parent := NewClassEntry("Parent")
	parent.Methods["getValue"] = &MethodDef{
		Name:           "getValue",
		Visibility:     VisibilityPublic,
		IsStatic:       true,
		DeclaringClass: "Parent",
	}

	// Child overrides static method
	child := NewClassEntry("Child")
	child.ParentClass = parent
	child.Methods["getValue"] = &MethodDef{
		Name:           "getValue",
		Visibility:     VisibilityPublic,
		IsStatic:       true,
		DeclaringClass: "Child",
	}

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Child's version should take precedence
	method := child.Methods["getValue"]
	if method.DeclaringClass != "Child" {
		t.Error("Child's static method should override parent's")
	}
}

// ============================================================================
// get_called_class() Support
// ============================================================================

func TestStatic_GetCalledClass(t *testing.T) {
	// The get_called_class() function returns the class name that was called
	// This is used for late static binding

	// Parent class
	parent := NewClassEntry("Parent")

	// Child class
	child := NewClassEntry("Child")
	child.ParentClass = parent

	// In runtime:
	// - Parent::method() would return "Parent"
	// - Child::method() would return "Child" (even if method is defined in Parent)

	// We verify the class structure supports this
	if parent.Name != "Parent" {
		t.Error("Parent name should be 'Parent'")
	}

	if child.Name != "Child" {
		t.Error("Child name should be 'Child'")
	}
}

// ============================================================================
// Static Property Access
// ============================================================================

func TestStatic_StaticPropertyAccess(t *testing.T) {
	// Static properties with late binding
	class := NewClassEntry("TestClass")

	class.Properties["counter"] = &PropertyDef{
		Name:       "counter",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		Type:       "int",
		Default:    NewInt(0),
	}

	// Verify static property
	prop := class.Properties["counter"]
	if !prop.IsStatic {
		t.Error("counter should be static property")
	}

	if prop.Default.ToInt() != 0 {
		t.Error("counter should have default value 0")
	}
}

func TestStatic_StaticPropertyInheritance(t *testing.T) {
	// Parent with static property
	parent := NewClassEntry("Parent")
	parent.Properties["shared"] = &PropertyDef{
		Name:       "shared",
		Visibility: VisibilityProtected,
		IsStatic:   true,
		Default:    NewString("parent"),
	}

	// Child inherits static property
	child := NewClassEntry("Child")
	child.ParentClass = parent

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Child should have access to inherited static property
	prop, exists := child.Properties["shared"]
	if !exists {
		t.Fatal("Child should inherit static property")
	}

	if !prop.IsStatic {
		t.Error("Inherited property should be static")
	}
}

// ============================================================================
// Static Return Type Resolution
// ============================================================================

func TestStatic_StaticReturnType(t *testing.T) {
	// Method returning 'static' type (late static binding for return)
	method := &MethodDef{
		Name:       "create",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "static",
	}

	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsStatic {
		t.Error("'static' should be recognized as static return type")
	}

	// When Parent::create() returns static, it returns Parent
	// When Child::create() returns static, it returns Child
	// This is resolved at runtime
}

func TestStatic_StaticVsSelfReturnType(t *testing.T) {
	// self returns the defining class
	// static returns the called class

	selfMethod := &MethodDef{
		Name:       "getInstance",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "self",
	}

	staticMethod := &MethodDef{
		Name:       "make",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "static",
	}

	selfType := ParseType(selfMethod.ReturnType)
	staticType := ParseType(staticMethod.ReturnType)

	if !selfType.IsSelf {
		t.Error("'self' should be recognized")
	}

	if !staticType.IsStatic {
		t.Error("'static' should be recognized")
	}

	// self and static are different
	if selfType.IsStatic || staticType.IsSelf {
		t.Error("'self' and 'static' should be distinct types")
	}
}

// ============================================================================
// Forward Static Call
// ============================================================================

func TestStatic_ForwardStaticCall(t *testing.T) {
	// PHP's forward_static_call() preserves the called class through the call chain
	// This is already handled by the VM's calledClass tracking

	parent := NewClassEntry("Parent")
	parent.Methods["start"] = &MethodDef{
		Name:       "start",
		Visibility: VisibilityPublic,
		IsStatic:   true,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// When Child::start() calls another static method,
	// that method should see Child as the called class
	// This is tracked by the VM frame's calledClass field
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestStatic_StaticInTrait(t *testing.T) {
	// Traits can use static:: too
	trait := NewTraitEntry("MyTrait")
	trait.Methods["traitMethod"] = &MethodDef{
		Name:       "traitMethod",
		Visibility: VisibilityPublic,
		IsStatic:   true,
	}

	// When a class uses this trait, static:: refers to the using class
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}

	if err := class.ApplyTraits(); err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Class should have the static method from trait
	method, exists := class.Methods["traitMethod"]
	if !exists {
		t.Fatal("Class should have trait's static method")
	}

	if !method.IsStatic {
		t.Error("Method from trait should be static")
	}
}

func TestStatic_MultiLevelInheritance(t *testing.T) {
	// Grandparent -> Parent -> Child
	// static:: should resolve to the actual called class, not intermediate

	grandparent := NewClassEntry("Grandparent")
	grandparent.Methods["staticMethod"] = &MethodDef{
		Name:           "staticMethod",
		Visibility:     VisibilityPublic,
		IsStatic:       true,
		DeclaringClass: "Grandparent",
	}

	parent := NewClassEntry("Parent")
	parent.ParentClass = grandparent
	if err := parent.InheritFrom(grandparent); err != nil {
		t.Fatalf("Parent inheritance failed: %v", err)
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent
	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Child inheritance failed: %v", err)
	}

	// Child::staticMethod() should resolve static:: to "Child"
	method, exists := child.Methods["staticMethod"]
	if !exists {
		t.Fatal("Child should inherit static method")
	}

	if !method.IsStatic {
		t.Error("Inherited method should be static")
	}

	// The declaring class is still Grandparent, but static:: at runtime
	// would resolve to Child
	if method.DeclaringClass != "Grandparent" {
		t.Error("Method should show original declaring class")
	}
}

// ============================================================================
// Static Constant Access
// ============================================================================

func TestStatic_StaticConstantWithInheritance(t *testing.T) {
	// Constants can be accessed with static:: too
	parent := NewClassEntry("Parent")
	parent.Constants["VALUE"] = &ClassConstant{
		Name:       "VALUE",
		Value:      NewInt(100),
		Visibility: VisibilityPublic,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent
	child.Constants["VALUE"] = &ClassConstant{
		Name:       "VALUE",
		Value:      NewInt(200),
		Visibility: VisibilityPublic,
	}

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Child::VALUE should be 200
	constant := child.Constants["VALUE"]
	if constant.Value.ToInt() != 200 {
		t.Errorf("Child::VALUE should be 200, got %d", constant.Value.ToInt())
	}

	// self::VALUE in parent would be 100
	// static::VALUE when called on Child would be 200
}
