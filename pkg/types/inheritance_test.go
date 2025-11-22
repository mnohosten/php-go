package types

import "testing"

// ============================================================================
// Class Extension Tests
// ============================================================================

func TestInheritance_BasicExtension(t *testing.T) {
	// Create parent class
	parent := NewClassEntry("ParentClass")
	parent.Properties["parentProp"] = &PropertyDef{
		Name:       "parentProp",
		Visibility: VisibilityPublic,
		Default:    NewString("parent value"),
	}
	parent.Methods["parentMethod"] = &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
	}

	// Create child class
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent

	// Inherit from parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Verify child has parent's property
	if _, exists := child.Properties["parentProp"]; !exists {
		t.Error("Child should inherit parent's property")
	}

	// Verify child has parent's method
	if _, exists := child.Methods["parentMethod"]; !exists {
		t.Error("Child should inherit parent's method")
	}
}

func TestInheritance_PropertyInheritance(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Properties["publicProp"] = &PropertyDef{
		Name:       "publicProp",
		Visibility: VisibilityPublic,
		Default:    NewString("public"),
	}
	parent.Properties["protectedProp"] = &PropertyDef{
		Name:       "protectedProp",
		Visibility: VisibilityProtected,
		Default:    NewString("protected"),
	}
	parent.Properties["privateProp"] = &PropertyDef{
		Name:       "privateProp",
		Visibility: VisibilityPrivate,
		Default:    NewString("private"),
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Public and protected should be inherited
	if _, exists := child.Properties["publicProp"]; !exists {
		t.Error("Child should inherit public property")
	}
	if _, exists := child.Properties["protectedProp"]; !exists {
		t.Error("Child should inherit protected property")
	}

	// Private should NOT be inherited
	if _, exists := child.Properties["privateProp"]; exists {
		t.Error("Child should NOT inherit private property")
	}
}

func TestInheritance_MethodOverride(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
		IsFinal:    false,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent

	// Override the method
	child.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
		IsFinal:    false,
	}

	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Child's method should still be the override
	method := child.Methods["method"]
	if method.DeclaringClass != "" && method.DeclaringClass != "Child" {
		t.Errorf("Expected child's override, got method from %s", method.DeclaringClass)
	}
}

func TestInheritance_FinalMethodCannotBeOverridden(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Methods["finalMethod"] = &MethodDef{
		Name:       "finalMethod",
		Visibility: VisibilityPublic,
		IsFinal:    true,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent

	// Try to override final method
	child.Methods["finalMethod"] = &MethodDef{
		Name:       "finalMethod",
		Visibility: VisibilityPublic,
	}

	err := child.InheritFrom(parent)
	if err == nil {
		t.Fatal("Expected error when overriding final method")
	}

	expectedMsg := "Cannot override final method Parent::finalMethod() in Child"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInheritance_FinalClassCannotBeExtended(t *testing.T) {
	parent := NewClassEntry("FinalParent")
	parent.IsFinal = true

	child := NewClassEntry("Child")
	child.ParentClass = parent

	err := child.InheritFrom(parent)
	if err == nil {
		t.Fatal("Expected error when extending final class")
	}

	expectedMsg := "Class Child cannot extend final class FinalParent"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInheritance_VisibilityCannotBeReduced(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Methods["publicMethod"] = &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityPublic,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent

	// Try to reduce visibility from public to protected
	child.Methods["publicMethod"] = &MethodDef{
		Name:       "publicMethod",
		Visibility: VisibilityProtected,
	}

	err := child.InheritFrom(parent)
	if err == nil {
		t.Fatal("Expected error when reducing visibility")
	}
}

func TestInheritance_AbstractMethodMustBeImplemented(t *testing.T) {
	parent := NewClassEntry("AbstractParent")
	parent.IsAbstract = true
	parent.Methods["abstractMethod"] = &MethodDef{
		Name:       "abstractMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Concrete child without implementation
	child := NewClassEntry("ConcreteChild")
	child.ParentClass = parent
	child.IsAbstract = false

	err := child.InheritFrom(parent)
	if err != nil {
		// Inherit should succeed, but instantiation should fail
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Check if abstract methods are tracked
	if !child.HasAbstractMethods() {
		t.Error("Concrete child should track unimplemented abstract methods")
	}
}

func TestInheritance_MultiLevelInheritance(t *testing.T) {
	// Grandparent -> Parent -> Child
	grandparent := NewClassEntry("Grandparent")
	grandparent.Methods["grandparentMethod"] = &MethodDef{
		Name:       "grandparentMethod",
		Visibility: VisibilityPublic,
	}

	parent := NewClassEntry("Parent")
	parent.ParentClass = grandparent
	err := parent.InheritFrom(grandparent)
	if err != nil {
		t.Fatalf("Parent InheritFrom Grandparent failed: %v", err)
	}

	parent.Methods["parentMethod"] = &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent
	err = child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("Child InheritFrom Parent failed: %v", err)
	}

	// Child should have both grandparent and parent methods
	if _, exists := child.Methods["grandparentMethod"]; !exists {
		t.Error("Child should inherit grandparent's method through parent")
	}
	if _, exists := child.Methods["parentMethod"]; !exists {
		t.Error("Child should inherit parent's method")
	}
}

// ============================================================================
// Constructor Inheritance Tests
// ============================================================================

func TestInheritance_ConstructorNotInherited(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Constructor = &MethodDef{
		Name:         "__construct",
		Visibility:   VisibilityPublic,
		IsConstructor: true,
	}
	parent.Methods["__construct"] = parent.Constructor

	child := NewClassEntry("Child")
	child.ParentClass = parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Constructor should NOT be inherited
	if child.Constructor != nil {
		t.Error("Child should not inherit parent's constructor")
	}
}

// ============================================================================
// Method Lookup Tests
// ============================================================================

func TestInheritance_MethodLookup(t *testing.T) {
	parent := NewClassEntry("Parent")
	parent.Methods["parentMethod"] = &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
	}

	child := NewClassEntry("Child")
	child.ParentClass = parent
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	child.Methods["childMethod"] = &MethodDef{
		Name:       "childMethod",
		Visibility: VisibilityPublic,
	}

	// GetMethod should find both child's and parent's methods
	if method, exists := child.GetMethod("childMethod"); !exists {
		t.Error("Should find child's method")
	} else if method.Name != "childMethod" {
		t.Error("Wrong method returned")
	}

	if method, exists := child.GetMethod("parentMethod"); !exists {
		t.Error("Should find parent's method through inheritance")
	} else if method.Name != "parentMethod" {
		t.Error("Wrong method returned")
	}

	if _, exists := child.GetMethod("nonExistent"); exists {
		t.Error("Should not find non-existent method")
	}
}

// ============================================================================
// Helper Method Tests
// ============================================================================

func TestInheritance_IsSubclassOf(t *testing.T) {
	grandparent := NewClassEntry("Grandparent")
	parent := NewClassEntry("Parent")
	parent.ParentClass = grandparent

	child := NewClassEntry("Child")
	child.ParentClass = parent

	// child is subclass of parent
	if !isSubclassOf(child, parent) {
		t.Error("Child should be subclass of Parent")
	}

	// child is subclass of grandparent (transitive)
	if !isSubclassOf(child, grandparent) {
		t.Error("Child should be subclass of Grandparent")
	}

	// parent is subclass of grandparent
	if !isSubclassOf(parent, grandparent) {
		t.Error("Parent should be subclass of Grandparent")
	}

	// parent is NOT subclass of child
	if isSubclassOf(parent, child) {
		t.Error("Parent should NOT be subclass of Child")
	}

	// Not subclass of nil
	if isSubclassOf(child, nil) {
		t.Error("Should not be subclass of nil")
	}

	// Not subclass of self
	if isSubclassOf(child, child) {
		t.Error("Should not be subclass of self")
	}
}
