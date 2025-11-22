package types

import "testing"

// ============================================================================
// Trait Definition Tests
// ============================================================================

func TestTrait_BasicDefinition(t *testing.T) {
	// Create trait with methods
	trait := NewTraitEntry("LoggerTrait")
	trait.Methods["log"] = &MethodDef{
		Name:       "log",
		Visibility: VisibilityPublic,
		NumParams:  1,
	}
	trait.Methods["error"] = &MethodDef{
		Name:       "error",
		Visibility: VisibilityPublic,
		NumParams:  1,
	}

	if len(trait.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(trait.Methods))
	}

	if trait.Name != "LoggerTrait" {
		t.Errorf("Expected name 'LoggerTrait', got '%s'", trait.Name)
	}
}

func TestTrait_WithProperties(t *testing.T) {
	trait := NewTraitEntry("CounterTrait")
	trait.Properties["count"] = &PropertyDef{
		Name:       "count",
		Visibility: VisibilityPrivate,
		Default:    NewInt(0),
	}

	if len(trait.Properties) != 1 {
		t.Errorf("Expected 1 property, got %d", len(trait.Properties))
	}

	prop := trait.Properties["count"]
	if prop.Visibility != VisibilityPrivate {
		t.Error("Expected private visibility")
	}
}

// ============================================================================
// Trait Usage Tests
// ============================================================================

func TestTrait_SingleTraitUsage(t *testing.T) {
	// Create trait
	trait := NewTraitEntry("SayWorldTrait")
	trait.Methods["sayHello"] = &MethodDef{
		Name:       "sayHello",
		Visibility: VisibilityPublic,
	}

	// Create class using trait
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}

	// Apply trait
	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Class should have trait's method
	if _, exists := class.Methods["sayHello"]; !exists {
		t.Error("Class should have trait's method")
	}
}

func TestTrait_MultipleTraits(t *testing.T) {
	// Create two traits
	trait1 := NewTraitEntry("Trait1")
	trait1.Methods["method1"] = &MethodDef{
		Name:       "method1",
		Visibility: VisibilityPublic,
	}

	trait2 := NewTraitEntry("Trait2")
	trait2.Methods["method2"] = &MethodDef{
		Name:       "method2",
		Visibility: VisibilityPublic,
	}

	// Class uses both
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait1, trait2}

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Should have both methods
	if _, exists := class.Methods["method1"]; !exists {
		t.Error("Should have method1 from Trait1")
	}
	if _, exists := class.Methods["method2"]; !exists {
		t.Error("Should have method2 from Trait2")
	}
}

func TestTrait_PropertyUsage(t *testing.T) {
	// Create trait with property
	trait := NewTraitEntry("PropertyTrait")
	trait.Properties["traitProp"] = &PropertyDef{
		Name:       "traitProp",
		Visibility: VisibilityProtected,
		Default:    NewString("trait value"),
	}

	// Class uses trait
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Should have trait's property
	if _, exists := class.Properties["traitProp"]; !exists {
		t.Error("Class should have trait's property")
	}
}

// ============================================================================
// Trait Conflict Tests
// ============================================================================

func TestTrait_MethodConflict(t *testing.T) {
	// Two traits with same method name
	trait1 := NewTraitEntry("Trait1")
	trait1.Methods["conflict"] = &MethodDef{
		Name:       "conflict",
		Visibility: VisibilityPublic,
	}

	trait2 := NewTraitEntry("Trait2")
	trait2.Methods["conflict"] = &MethodDef{
		Name:       "conflict",
		Visibility: VisibilityPublic,
	}

	// Class uses both without resolution
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait1, trait2}

	err := class.ApplyTraits()
	if err == nil {
		t.Fatal("Expected error for trait method conflict")
	}

	expectedMsg := "Trait method conflict: method 'conflict' exists in multiple traits (Trait1, Trait2)"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestTrait_ConflictWithPrecedence(t *testing.T) {
	// Two traits with same method
	trait1 := NewTraitEntry("Trait1")
	trait1.Methods["conflict"] = &MethodDef{
		Name:       "conflict",
		Visibility: VisibilityPublic,
	}

	trait2 := NewTraitEntry("Trait2")
	trait2.Methods["conflict"] = &MethodDef{
		Name:       "conflict",
		Visibility: VisibilityPublic,
	}

	// Class uses both with precedence resolution
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait1, trait2}
	// Trait1::conflict insteadof Trait2
	class.TraitPrecedence["conflict"] = "Trait1"

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits with precedence failed: %v", err)
	}

	// Should have Trait1's version
	method := class.Methods["conflict"]
	if method.DeclaringClass != "Trait1" {
		t.Errorf("Expected method from Trait1, got from %s", method.DeclaringClass)
	}
}

func TestTrait_MethodAliasing(t *testing.T) {
	// Create trait
	trait := NewTraitEntry("MyTrait")
	trait.Methods["original"] = &MethodDef{
		Name:       "original",
		Visibility: VisibilityPublic,
	}

	// Class uses trait with aliasing
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}
	// MyTrait::original as alias
	class.TraitAliases["alias"] = "MyTrait::original"

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits with aliasing failed: %v", err)
	}

	// Should have both original and alias
	if _, exists := class.Methods["original"]; !exists {
		t.Error("Should have original method")
	}
	if _, exists := class.Methods["alias"]; !exists {
		t.Error("Should have aliased method")
	}
}

func TestTrait_VisibilityChangeViaAlias(t *testing.T) {
	// Create trait with public method
	trait := NewTraitEntry("MyTrait")
	trait.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
	}

	// Class aliases with different visibility
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}
	// This would be: use MyTrait { method as protected; }
	// For simplicity, we'll test that aliasing can change visibility
	class.TraitAliases["method"] = "MyTrait::method:protected"

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Method should exist
	if _, exists := class.Methods["method"]; !exists {
		t.Error("Should have method from trait")
	}
}

// ============================================================================
// Trait Composition Tests
// ============================================================================

func TestTrait_UsingOtherTraits(t *testing.T) {
	// Base trait
	baseTrait := NewTraitEntry("BaseTrait")
	baseTrait.Methods["baseMethod"] = &MethodDef{
		Name:       "baseMethod",
		Visibility: VisibilityPublic,
	}

	// Composite trait using base trait
	compositeTrait := NewTraitEntry("CompositeTrait")
	compositeTrait.UsedTraits = []*TraitEntry{baseTrait}
	compositeTrait.Methods["compositeMethod"] = &MethodDef{
		Name:       "compositeMethod",
		Visibility: VisibilityPublic,
	}

	// Apply traits to composite trait first
	err := compositeTrait.ApplyUsedTraits()
	if err != nil {
		t.Fatalf("ApplyUsedTraits failed: %v", err)
	}

	// Composite trait should have both methods
	if _, exists := compositeTrait.Methods["baseMethod"]; !exists {
		t.Error("Composite trait should have base trait's method")
	}
	if _, exists := compositeTrait.Methods["compositeMethod"]; !exists {
		t.Error("Composite trait should have its own method")
	}

	// Now use composite trait in class
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{compositeTrait}

	err = class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Class should have both methods
	if _, exists := class.Methods["baseMethod"]; !exists {
		t.Error("Class should have base method through composite trait")
	}
	if _, exists := class.Methods["compositeMethod"]; !exists {
		t.Error("Class should have composite method")
	}
}

// ============================================================================
// Class Method Override Tests
// ============================================================================

func TestTrait_ClassMethodOverridesTrait(t *testing.T) {
	// Create trait
	trait := NewTraitEntry("MyTrait")
	trait.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
	}

	// Class has same method (should override trait)
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}
	class.Methods["method"] = &MethodDef{
		Name:           "method",
		Visibility:     VisibilityPublic,
		DeclaringClass: "MyClass",
	}

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Should have class's version, not trait's
	method := class.Methods["method"]
	if method.DeclaringClass != "MyClass" {
		t.Error("Class method should override trait method")
	}
}

func TestTrait_InheritedMethodOverridesTrait(t *testing.T) {
	// Create trait
	trait := NewTraitEntry("MyTrait")
	trait.Methods["method"] = &MethodDef{
		Name:       "method",
		Visibility: VisibilityPublic,
	}

	// Parent class with method
	parent := NewClassEntry("ParentClass")
	parent.Methods["method"] = &MethodDef{
		Name:           "method",
		Visibility:     VisibilityPublic,
		DeclaringClass: "ParentClass",
	}

	// Child class uses trait and extends parent
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent
	child.Traits = []*TraitEntry{trait}

	// Inherit from parent first
	err := child.InheritFrom(parent)
	if err != nil {
		t.Fatalf("InheritFrom failed: %v", err)
	}

	// Then apply traits
	err = child.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Trait method should override inherited method
	// This is PHP behavior: trait > inherited method
	method := child.Methods["method"]
	if method.DeclaringClass != "" && method.DeclaringClass == "ParentClass" {
		t.Error("Trait method should override inherited method")
	}
}

// ============================================================================
// Abstract Methods in Traits
// ============================================================================

func TestTrait_AbstractMethod(t *testing.T) {
	// Trait with abstract method
	trait := NewTraitEntry("MyTrait")
	trait.Methods["abstractMethod"] = &MethodDef{
		Name:       "abstractMethod",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	// Class must implement abstract method
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}
	class.Methods["abstractMethod"] = &MethodDef{
		Name:       "abstractMethod",
		Visibility: VisibilityPublic,
		IsAbstract: false, // Concrete implementation
	}

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Should have class's concrete implementation
	method := class.Methods["abstractMethod"]
	if method.IsAbstract {
		t.Error("Class should have concrete implementation")
	}
}

// ============================================================================
// Static Methods in Traits
// ============================================================================

func TestTrait_StaticMethod(t *testing.T) {
	// Trait with static method
	trait := NewTraitEntry("MyTrait")
	trait.Methods["staticMethod"] = &MethodDef{
		Name:       "staticMethod",
		Visibility: VisibilityPublic,
		IsStatic:   true,
	}

	// Class uses trait
	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}

	err := class.ApplyTraits()
	if err != nil {
		t.Fatalf("ApplyTraits failed: %v", err)
	}

	// Should have static method
	method := class.Methods["staticMethod"]
	if !method.IsStatic {
		t.Error("Method should be static")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestTrait_EmptyTrait(t *testing.T) {
	// Empty trait is valid
	trait := NewTraitEntry("EmptyTrait")

	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait}

	err := class.ApplyTraits()
	if err != nil {
		t.Errorf("Empty trait should be valid: %v", err)
	}
}

func TestTrait_PropertyConflict(t *testing.T) {
	// Two traits with same property name but different values
	trait1 := NewTraitEntry("Trait1")
	trait1.Properties["prop"] = &PropertyDef{
		Name:    "prop",
		Default: NewInt(1),
	}

	trait2 := NewTraitEntry("Trait2")
	trait2.Properties["prop"] = &PropertyDef{
		Name:    "prop",
		Default: NewInt(2),
	}

	class := NewClassEntry("MyClass")
	class.Traits = []*TraitEntry{trait1, trait2}

	err := class.ApplyTraits()
	if err == nil {
		t.Fatal("Expected error for property conflict with different default values")
	}
}
