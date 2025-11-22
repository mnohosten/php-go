package types

import (
	"testing"
)

// ============================================================================
// ClassEntry Tests
// ============================================================================

func TestNewClassEntry(t *testing.T) {
	class := NewClassEntry("MyClass")

	if class.Name != "MyClass" {
		t.Errorf("Expected class name 'MyClass', got '%s'", class.Name)
	}

	if class.Constants == nil {
		t.Error("Expected Constants map to be initialized")
	}

	if class.Properties == nil {
		t.Error("Expected Properties map to be initialized")
	}

	if class.Methods == nil {
		t.Error("Expected Methods map to be initialized")
	}
}

func TestClassEntryGetMethod(t *testing.T) {
	// Create parent class
	parent := NewClassEntry("ParentClass")
	parentMethod := &MethodDef{
		Name:       "parentMethod",
		Visibility: VisibilityPublic,
	}
	parent.Methods["parentMethod"] = parentMethod

	// Create child class
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent
	childMethod := &MethodDef{
		Name:       "childMethod",
		Visibility: VisibilityPublic,
	}
	child.Methods["childMethod"] = childMethod

	// Test finding child's own method
	method, exists := child.GetMethod("childMethod")
	if !exists {
		t.Error("Expected to find childMethod in child class")
	}
	if method.Name != "childMethod" {
		t.Errorf("Expected method name 'childMethod', got '%s'", method.Name)
	}

	// Test finding inherited method
	method, exists = child.GetMethod("parentMethod")
	if !exists {
		t.Error("Expected to find parentMethod inherited from parent class")
	}
	if method.Name != "parentMethod" {
		t.Errorf("Expected method name 'parentMethod', got '%s'", method.Name)
	}

	// Test method not found
	_, exists = child.GetMethod("nonexistent")
	if exists {
		t.Error("Expected not to find nonexistent method")
	}
}

func TestClassEntryImplementsInterface(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("MyInterface")

	// Create class that implements interface
	class := NewClassEntry("MyClass")
	class.Interfaces = []*InterfaceEntry{iface}

	// Test direct implementation
	if !class.ImplementsInterface("MyInterface") {
		t.Error("Expected class to implement MyInterface")
	}

	// Test non-implemented interface
	if class.ImplementsInterface("OtherInterface") {
		t.Error("Expected class not to implement OtherInterface")
	}
}

func TestClassEntryImplementsInterfaceThroughParent(t *testing.T) {
	// Create interface
	iface := NewInterfaceEntry("MyInterface")

	// Create parent class that implements interface
	parent := NewClassEntry("ParentClass")
	parent.Interfaces = []*InterfaceEntry{iface}

	// Create child class
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent

	// Test inheritance of interface implementation
	if !child.ImplementsInterface("MyInterface") {
		t.Error("Expected child class to inherit interface implementation from parent")
	}
}

func TestClassEntryImplementsInterfaceHierarchy(t *testing.T) {
	// Create parent interface
	parentIface := NewInterfaceEntry("ParentInterface")

	// Create child interface that extends parent
	childIface := NewInterfaceEntry("ChildInterface")
	childIface.ParentInterfaces = []*InterfaceEntry{parentIface}

	// Create class that implements child interface
	class := NewClassEntry("MyClass")
	class.Interfaces = []*InterfaceEntry{childIface}

	// Test direct interface
	if !class.ImplementsInterface("ChildInterface") {
		t.Error("Expected class to implement ChildInterface")
	}

	// Test parent interface through hierarchy
	if !class.ImplementsInterface("ParentInterface") {
		t.Error("Expected class to implement ParentInterface through ChildInterface")
	}
}

// ============================================================================
// Object Creation and Management Tests
// ============================================================================

func TestNewObjectFromClass(t *testing.T) {
	// Create class with properties
	class := NewClassEntry("MyClass")
	class.Properties["name"] = &PropertyDef{
		Name:       "name",
		Visibility: VisibilityPublic,
		HasDefault: true,
		Default:    NewString("default"),
	}
	class.Properties["count"] = &PropertyDef{
		Name:       "count",
		Visibility: VisibilityPrivate,
		HasDefault: true,
		Default:    NewInt(0),
	}

	// Create object
	obj := NewObjectFromClass(class)

	if obj.ClassName != "MyClass" {
		t.Errorf("Expected class name 'MyClass', got '%s'", obj.ClassName)
	}

	if obj.ClassEntry != class {
		t.Error("Expected ClassEntry to be set correctly")
	}

	if obj.ObjectID == 0 {
		t.Error("Expected ObjectID to be non-zero")
	}

	if obj.IsDestroyed {
		t.Error("Expected IsDestroyed to be false")
	}

	// Check properties were initialized with defaults
	if len(obj.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(obj.Properties))
	}

	nameProp, exists := obj.Properties["name"]
	if !exists {
		t.Error("Expected 'name' property to exist")
	}
	if nameProp.Value.ToString() != "default" {
		t.Errorf("Expected name default value 'default', got '%s'", nameProp.Value.ToString())
	}
}

func TestNewObjectInstance(t *testing.T) {
	obj := NewObjectInstance("TestClass")

	if obj.ClassName != "TestClass" {
		t.Errorf("Expected class name 'TestClass', got '%s'", obj.ClassName)
	}

	if obj.Properties == nil {
		t.Error("Expected Properties map to be initialized")
	}

	if obj.ObjectID == 0 {
		t.Error("Expected ObjectID to be non-zero")
	}
}

func TestObjectIDUniqueness(t *testing.T) {
	obj1 := NewObjectInstance("Class1")
	obj2 := NewObjectInstance("Class2")

	if obj1.ObjectID == obj2.ObjectID {
		t.Error("Expected ObjectIDs to be unique")
	}

	if obj2.ObjectID != obj1.ObjectID+1 {
		t.Error("Expected ObjectIDs to be sequential")
	}
}

// ============================================================================
// Property Access Tests
// ============================================================================

func TestObjectGetPropertyPublic(t *testing.T) {
	// Create class and object
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)

	// Add public property
	obj.Properties["publicProp"] = &Property{
		Value:      NewString("test"),
		Visibility: VisibilityPublic,
	}

	// Access from any context (nil = no context)
	value, exists := obj.GetProperty("publicProp", nil)
	if !exists {
		t.Error("Expected to access public property")
	}
	if value.ToString() != "test" {
		t.Errorf("Expected value 'test', got '%s'", value.ToString())
	}
}

func TestObjectGetPropertyProtected(t *testing.T) {
	// Create class hierarchy
	parent := NewClassEntry("ParentClass")
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent

	obj := NewObjectFromClass(parent)
	obj.Properties["protectedProp"] = &Property{
		Value:      NewString("protected"),
		Visibility: VisibilityProtected,
	}

	// Access from same class - should succeed
	value, exists := obj.GetProperty("protectedProp", parent)
	if !exists {
		t.Error("Expected to access protected property from same class")
	}
	if value.ToString() != "protected" {
		t.Errorf("Expected value 'protected', got '%s'", value.ToString())
	}

	// Access from child class - should succeed
	value, exists = obj.GetProperty("protectedProp", child)
	if !exists {
		t.Error("Expected to access protected property from child class")
	}

	// Access from unrelated class - should fail
	unrelated := NewClassEntry("UnrelatedClass")
	_, exists = obj.GetProperty("protectedProp", unrelated)
	if exists {
		t.Error("Expected not to access protected property from unrelated class")
	}

	// Access from no context - should fail
	_, exists = obj.GetProperty("protectedProp", nil)
	if exists {
		t.Error("Expected not to access protected property without context")
	}
}

func TestObjectGetPropertyPrivate(t *testing.T) {
	// Create class
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)
	obj.Properties["privateProp"] = &Property{
		Value:      NewString("private"),
		Visibility: VisibilityPrivate,
	}

	// Access from same class - should succeed
	value, exists := obj.GetProperty("privateProp", class)
	if !exists {
		t.Error("Expected to access private property from same class")
	}
	if value.ToString() != "private" {
		t.Errorf("Expected value 'private', got '%s'", value.ToString())
	}

	// Access from child class - should fail
	child := NewClassEntry("ChildClass")
	child.ParentClass = class
	_, exists = obj.GetProperty("privateProp", child)
	if exists {
		t.Error("Expected not to access private property from child class")
	}

	// Access from no context - should fail
	_, exists = obj.GetProperty("privateProp", nil)
	if exists {
		t.Error("Expected not to access private property without context")
	}
}

func TestObjectGetPropertyNonexistent(t *testing.T) {
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)

	_, exists := obj.GetProperty("nonexistent", nil)
	if exists {
		t.Error("Expected not to find nonexistent property")
	}
}

func TestObjectSetPropertyPublic(t *testing.T) {
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)

	// Set public property from any context
	success := obj.SetProperty("publicProp", NewString("value"), nil)
	if !success {
		t.Error("Expected to set public property")
	}

	// Verify value was set
	value, _ := obj.GetProperty("publicProp", nil)
	if value.ToString() != "value" {
		t.Errorf("Expected value 'value', got '%s'", value.ToString())
	}
}

func TestObjectSetPropertyProtected(t *testing.T) {
	// Create class hierarchy
	parent := NewClassEntry("ParentClass")
	child := NewClassEntry("ChildClass")
	child.ParentClass = parent

	obj := NewObjectFromClass(parent)
	obj.Properties["protectedProp"] = &Property{
		Value:      NewString("original"),
		Visibility: VisibilityProtected,
	}

	// Set from same class - should succeed
	success := obj.SetProperty("protectedProp", NewString("modified"), parent)
	if !success {
		t.Error("Expected to set protected property from same class")
	}

	// Set from child class - should succeed
	success = obj.SetProperty("protectedProp", NewString("by-child"), child)
	if !success {
		t.Error("Expected to set protected property from child class")
	}

	// Set from unrelated class - should fail
	unrelated := NewClassEntry("UnrelatedClass")
	success = obj.SetProperty("protectedProp", NewString("by-unrelated"), unrelated)
	if success {
		t.Error("Expected not to set protected property from unrelated class")
	}
}

func TestObjectSetPropertyPrivate(t *testing.T) {
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)
	obj.Properties["privateProp"] = &Property{
		Value:      NewString("original"),
		Visibility: VisibilityPrivate,
	}

	// Set from same class - should succeed
	success := obj.SetProperty("privateProp", NewString("modified"), class)
	if !success {
		t.Error("Expected to set private property from same class")
	}

	// Set from child class - should fail
	child := NewClassEntry("ChildClass")
	child.ParentClass = class
	success = obj.SetProperty("privateProp", NewString("by-child"), child)
	if success {
		t.Error("Expected not to set private property from child class")
	}
}

func TestObjectSetPropertyReadonly(t *testing.T) {
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)
	obj.Properties["readonlyProp"] = &Property{
		Value:      NewString("initial"),
		Visibility: VisibilityPublic,
		IsReadOnly: true,
	}

	// Try to modify readonly property - should fail
	success := obj.SetProperty("readonlyProp", NewString("modified"), nil)
	if success {
		t.Error("Expected not to modify readonly property")
	}

	// Value should remain unchanged
	value, _ := obj.GetProperty("readonlyProp", nil)
	if value.ToString() != "initial" {
		t.Errorf("Expected value 'initial', got '%s'", value.ToString())
	}
}

func TestObjectSetPropertyReadonlyUninitialized(t *testing.T) {
	class := NewClassEntry("MyClass")
	obj := NewObjectFromClass(class)
	obj.Properties["readonlyProp"] = &Property{
		Value:      nil, // Uninitialized
		Visibility: VisibilityPublic,
		IsReadOnly: true,
	}

	// Setting uninitialized readonly property - should succeed (during construction)
	success := obj.SetProperty("readonlyProp", NewString("value"), nil)
	if !success {
		t.Error("Expected to set uninitialized readonly property")
	}

	// Try to modify again - should fail
	success = obj.SetProperty("readonlyProp", NewString("modified"), nil)
	if success {
		t.Error("Expected not to modify readonly property after initialization")
	}
}

// ============================================================================
// PropertyVisibility Tests
// ============================================================================

func TestPropertyVisibilityString(t *testing.T) {
	tests := []struct {
		visibility PropertyVisibility
		expected   string
	}{
		{VisibilityPublic, "public"},
		{VisibilityProtected, "protected"},
		{VisibilityPrivate, "private"},
	}

	for _, test := range tests {
		result := test.visibility.String()
		if result != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, result)
		}
	}
}

// ============================================================================
// Interface and Trait Tests
// ============================================================================

func TestNewInterfaceEntry(t *testing.T) {
	iface := NewInterfaceEntry("MyInterface")

	if iface.Name != "MyInterface" {
		t.Errorf("Expected interface name 'MyInterface', got '%s'", iface.Name)
	}

	if iface.Methods == nil {
		t.Error("Expected Methods map to be initialized")
	}

	if iface.Constants == nil {
		t.Error("Expected Constants map to be initialized")
	}

	if iface.ParentInterfaces == nil {
		t.Error("Expected ParentInterfaces slice to be initialized")
	}
}

func TestNewTraitEntry(t *testing.T) {
	trait := NewTraitEntry("MyTrait")

	if trait.Name != "MyTrait" {
		t.Errorf("Expected trait name 'MyTrait', got '%s'", trait.Name)
	}

	if trait.Properties == nil {
		t.Error("Expected Properties map to be initialized")
	}

	if trait.Methods == nil {
		t.Error("Expected Methods map to be initialized")
	}

	if trait.UsedTraits == nil {
		t.Error("Expected UsedTraits slice to be initialized")
	}
}

// ============================================================================
// Inheritance Tests
// ============================================================================

func TestIsSubclassOf(t *testing.T) {
	// Create class hierarchy: GrandParent -> Parent -> Child
	grandParent := NewClassEntry("GrandParent")
	parent := NewClassEntry("Parent")
	parent.ParentClass = grandParent
	child := NewClassEntry("Child")
	child.ParentClass = parent

	// Test direct parent
	if !isSubclassOf(child, parent) {
		t.Error("Expected Child to be subclass of Parent")
	}

	// Test grandparent
	if !isSubclassOf(child, grandParent) {
		t.Error("Expected Child to be subclass of GrandParent")
	}

	// Test same class (should be false)
	if isSubclassOf(child, child) {
		t.Error("Expected class not to be subclass of itself")
	}

	// Test unrelated class
	unrelated := NewClassEntry("Unrelated")
	if isSubclassOf(child, unrelated) {
		t.Error("Expected Child not to be subclass of Unrelated")
	}

	// Test reverse hierarchy
	if isSubclassOf(parent, child) {
		t.Error("Expected Parent not to be subclass of Child")
	}

	// Test nil classes
	if isSubclassOf(nil, parent) {
		t.Error("Expected nil not to be subclass")
	}
	if isSubclassOf(child, nil) {
		t.Error("Expected nil not to be parent")
	}
}

// ============================================================================
// Property Access Helper Tests
// ============================================================================

func TestCanAccessProperty(t *testing.T) {
	// Create class hierarchy
	parent := NewClassEntry("Parent")
	child := NewClassEntry("Child")
	child.ParentClass = parent
	unrelated := NewClassEntry("Unrelated")

	tests := []struct {
		name           string
		visibility     PropertyVisibility
		accessContext  *ClassEntry
		ownerClass     *ClassEntry
		expectedAccess bool
	}{
		// Public access
		{"public from same class", VisibilityPublic, parent, parent, true},
		{"public from child", VisibilityPublic, child, parent, true},
		{"public from unrelated", VisibilityPublic, unrelated, parent, true},
		{"public from nil", VisibilityPublic, nil, parent, true},

		// Protected access
		{"protected from same class", VisibilityProtected, parent, parent, true},
		{"protected from child", VisibilityProtected, child, parent, true},
		{"protected from unrelated", VisibilityProtected, unrelated, parent, false},
		{"protected from nil", VisibilityProtected, nil, parent, false},

		// Private access
		{"private from same class", VisibilityPrivate, parent, parent, true},
		{"private from child", VisibilityPrivate, child, parent, false},
		{"private from unrelated", VisibilityPrivate, unrelated, parent, false},
		{"private from nil", VisibilityPrivate, nil, parent, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prop := &Property{
				Visibility: test.visibility,
			}
			result := canAccessProperty(prop, test.accessContext, test.ownerClass)
			if result != test.expectedAccess {
				t.Errorf("Expected access=%v, got %v", test.expectedAccess, result)
			}
		})
	}
}
