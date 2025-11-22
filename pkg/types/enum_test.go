package types

import "testing"

// ============================================================================
// Pure Enum Tests (Unbacked Enums)
// ============================================================================

func TestEnum_PureEnumDefinition(t *testing.T) {
	// Create a pure enum: enum Status { case Pending; case Active; case Inactive; }
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)
	enum.AddCase("Active", nil)
	enum.AddCase("Inactive", nil)

	if !enum.IsEnum {
		t.Error("Expected IsEnum to be true")
	}

	if enum.EnumBackingType != "" {
		t.Errorf("Expected no backing type, got %s", enum.EnumBackingType)
	}

	if len(enum.EnumCases) != 3 {
		t.Errorf("Expected 3 cases, got %d", len(enum.EnumCases))
	}

	// Check that cases exist
	if _, exists := enum.EnumCases["Pending"]; !exists {
		t.Error("Case 'Pending' should exist")
	}
	if _, exists := enum.EnumCases["Active"]; !exists {
		t.Error("Case 'Active' should exist")
	}
	if _, exists := enum.EnumCases["Inactive"]; !exists {
		t.Error("Case 'Inactive' should exist")
	}
}

func TestEnum_PureCaseAccess(t *testing.T) {
	enum := NewEnumEntry("Color", "")
	enum.AddCase("Red", nil)
	enum.AddCase("Green", nil)
	enum.AddCase("Blue", nil)

	// Access cases - for pure enums, the value is nil but the key should exist
	red, exists := enum.EnumCases["Red"]
	if !exists {
		t.Fatal("Red case should exist in EnumCases")
	}

	// Pure enum cases have nil as their value
	if red != nil {
		t.Error("Pure enum cases should have nil value")
	}

	// Pure enum cases should have the case name as their representation
	// In PHP: Color::Red->name returns "Red"
}

// ============================================================================
// Backed Enum Tests (Int and String)
// ============================================================================

func TestEnum_IntBackedEnum(t *testing.T) {
	// Create int-backed enum: enum Priority: int { case Low = 1; case Medium = 5; case High = 10; }
	enum := NewEnumEntry("Priority", "int")
	enum.AddCase("Low", NewInt(1))
	enum.AddCase("Medium", NewInt(5))
	enum.AddCase("High", NewInt(10))

	if enum.EnumBackingType != "int" {
		t.Errorf("Expected backing type 'int', got '%s'", enum.EnumBackingType)
	}

	// Check case values
	if low := enum.EnumCases["Low"]; low == nil || low.ToInt() != 1 {
		t.Error("Low case should have value 1")
	}

	if medium := enum.EnumCases["Medium"]; medium == nil || medium.ToInt() != 5 {
		t.Error("Medium case should have value 5")
	}

	if high := enum.EnumCases["High"]; high == nil || high.ToInt() != 10 {
		t.Error("High case should have value 10")
	}
}

func TestEnum_StringBackedEnum(t *testing.T) {
	// Create string-backed enum: enum Suit: string { case Hearts = 'H'; case Diamonds = 'D'; }
	enum := NewEnumEntry("Suit", "string")
	enum.AddCase("Hearts", NewString("H"))
	enum.AddCase("Diamonds", NewString("D"))
	enum.AddCase("Clubs", NewString("C"))
	enum.AddCase("Spades", NewString("S"))

	if enum.EnumBackingType != "string" {
		t.Errorf("Expected backing type 'string', got '%s'", enum.EnumBackingType)
	}

	// Check case values
	if hearts := enum.EnumCases["Hearts"]; hearts == nil || hearts.ToString() != "H" {
		t.Error("Hearts case should have value 'H'")
	}

	if diamonds := enum.EnumCases["Diamonds"]; diamonds == nil || diamonds.ToString() != "D" {
		t.Error("Diamonds case should have value 'D'")
	}
}

func TestEnum_BackedEnumValidation(t *testing.T) {
	// Test that backed enums reject invalid backing types
	enum := NewEnumEntry("InvalidEnum", "bool")
	err := enum.Validate()
	if err == nil {
		t.Fatal("Expected error for invalid backing type 'bool'")
	}

	// Valid backing types
	enum2 := NewEnumEntry("ValidEnum", "int")
	if err := enum2.Validate(); err != nil {
		t.Errorf("int backing type should be valid: %v", err)
	}

	enum3 := NewEnumEntry("ValidEnum", "string")
	if err := enum3.Validate(); err != nil {
		t.Errorf("string backing type should be valid: %v", err)
	}
}

func TestEnum_BackedEnumValueTypeMismatch(t *testing.T) {
	// Test that int-backed enum rejects string values
	enum := NewEnumEntry("Priority", "int")
	enum.AddCase("Low", NewString("low")) // Wrong type

	err := enum.Validate()
	if err == nil {
		t.Fatal("Expected error for type mismatch in backed enum")
	}
}

// ============================================================================
// Enum Methods Tests
// ============================================================================

func TestEnum_EnumWithMethods(t *testing.T) {
	// Enums can have methods
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)
	enum.AddCase("Active", nil)

	// Add a method to the enum
	enum.Methods["isActive"] = &MethodDef{
		Name:       "isActive",
		Visibility: VisibilityPublic,
		NumParams:  0,
	}

	if _, exists := enum.Methods["isActive"]; !exists {
		t.Error("Enum should have isActive method")
	}
}

func TestEnum_EnumWithStaticMethods(t *testing.T) {
	// Enums can have static methods
	enum := NewEnumEntry("Color", "")
	enum.AddCase("Red", nil)

	enum.Methods["fromName"] = &MethodDef{
		Name:       "fromName",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		NumParams:  1,
	}

	method := enum.Methods["fromName"]
	if !method.IsStatic {
		t.Error("fromName should be a static method")
	}
}

// ============================================================================
// Enum Built-in Methods Tests
// ============================================================================

func TestEnum_CasesMethod(t *testing.T) {
	// cases() method returns all enum cases
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)
	enum.AddCase("Active", nil)
	enum.AddCase("Inactive", nil)

	// The cases() method should be automatically available
	cases := enum.GetCases()
	if len(cases) != 3 {
		t.Errorf("Expected 3 cases from cases(), got %d", len(cases))
	}
}

func TestEnum_FromMethod(t *testing.T) {
	// from() method for backed enums
	enum := NewEnumEntry("Priority", "int")
	enum.AddCase("Low", NewInt(1))
	enum.AddCase("Medium", NewInt(5))
	enum.AddCase("High", NewInt(10))

	// from(5) should return Medium case
	caseValue, err := enum.From(NewInt(5))
	if err != nil {
		t.Fatalf("from(5) should succeed: %v", err)
	}

	if caseValue != "Medium" {
		t.Errorf("from(5) should return 'Medium', got '%s'", caseValue)
	}

	// from(999) should error (value doesn't exist)
	_, err = enum.From(NewInt(999))
	if err == nil {
		t.Fatal("from(999) should fail for non-existent value")
	}
}

func TestEnum_TryFromMethod(t *testing.T) {
	// tryFrom() method for backed enums (returns null on failure)
	enum := NewEnumEntry("Suit", "string")
	enum.AddCase("Hearts", NewString("H"))
	enum.AddCase("Diamonds", NewString("D"))

	// tryFrom("H") should return Hearts
	caseValue, err := enum.TryFrom(NewString("H"))
	if err != nil {
		t.Fatalf("tryFrom('H') should succeed: %v", err)
	}

	if caseValue != "Hearts" {
		t.Errorf("tryFrom('H') should return 'Hearts', got '%s'", caseValue)
	}

	// tryFrom("X") should return null (not error)
	caseValue, err = enum.TryFrom(NewString("X"))
	if err != nil {
		t.Error("tryFrom should not error on non-existent value")
	}

	if caseValue != "" {
		t.Error("tryFrom should return empty string for non-existent value")
	}
}

func TestEnum_PureEnumNoFromMethod(t *testing.T) {
	// Pure enums should not have from() or tryFrom()
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)

	// Calling from() on pure enum should error
	_, err := enum.From(NewInt(1))
	if err == nil {
		t.Fatal("from() should not be available on pure enums")
	}
}

// ============================================================================
// Enum Interfaces Tests
// ============================================================================

func TestEnum_ImplementsInterfaces(t *testing.T) {
	// Enums can implement interfaces
	iface := NewInterfaceEntry("Renderable")
	iface.Methods["render"] = &MethodDef{
		Name:       "render",
		Visibility: VisibilityPublic,
		IsAbstract: true,
	}

	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)
	enum.Interfaces = []*InterfaceEntry{iface}
	enum.Methods["render"] = &MethodDef{
		Name:       "render",
		Visibility: VisibilityPublic,
	}

	// Validate interface implementation
	err := enum.ValidateInterfaceImplementation()
	if err != nil {
		t.Errorf("Enum should be able to implement interfaces: %v", err)
	}
}

// ============================================================================
// Enum Restrictions Tests
// ============================================================================

func TestEnum_CannotExtend(t *testing.T) {
	// Enums cannot extend other classes
	enum := NewEnumEntry("Status", "")
	parentClass := NewClassEntry("BaseClass")

	enum.ParentClass = parentClass

	err := enum.Validate()
	if err == nil {
		t.Fatal("Enums should not be able to extend classes")
	}
}

func TestEnum_CannotBeExtended(t *testing.T) {
	// Classes cannot extend enums
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)

	childClass := NewClassEntry("MyClass")
	childClass.ParentClass = enum

	err := childClass.InheritFrom(enum)
	if err == nil {
		t.Fatal("Classes should not be able to extend enums")
	}
}

func TestEnum_CannotHaveProperties(t *testing.T) {
	// Enums cannot have properties (except readonly promoted constructor properties in PHP 8.2+)
	// For now, we'll test that standard properties are not allowed
	enum := NewEnumEntry("Status", "")
	enum.Properties["invalidProp"] = &PropertyDef{
		Name:       "invalidProp",
		Visibility: VisibilityPublic,
	}

	err := enum.Validate()
	if err == nil {
		t.Fatal("Enums should not be able to have non-const properties")
	}
}

// ============================================================================
// Enum Edge Cases
// ============================================================================

func TestEnum_EmptyEnum(t *testing.T) {
	// Enum with no cases (technically valid but unusual)
	enum := NewEnumEntry("EmptyEnum", "")

	if len(enum.EnumCases) != 0 {
		t.Error("Empty enum should have 0 cases")
	}

	// Should still be valid
	if err := enum.Validate(); err != nil {
		t.Errorf("Empty enum should be valid: %v", err)
	}
}

func TestEnum_DuplicateCases(t *testing.T) {
	// Note: Since EnumCases is a map, adding the same case name twice
	// will just overwrite the first one, so there won't be duplicates.
	// This test verifies that duplicate additions don't cause issues.
	enum := NewEnumEntry("Status", "")
	enum.AddCase("Pending", nil)
	enum.AddCase("Pending", nil) // This overwrites the first one

	// Should still be valid (only one "Pending" case in the map)
	err := enum.Validate()
	if err != nil {
		t.Errorf("Enum should be valid: %v", err)
	}

	// Should have only 1 case
	if len(enum.EnumCases) != 1 {
		t.Errorf("Expected 1 case after duplicate addition, got %d", len(enum.EnumCases))
	}
}

func TestEnum_DuplicateBackingValues(t *testing.T) {
	// In backed enums, duplicate values are allowed (unlike case names)
	enum := NewEnumEntry("Status", "int")
	enum.AddCase("Active", NewInt(1))
	enum.AddCase("Running", NewInt(1)) // Same value, different name - allowed

	err := enum.Validate()
	if err != nil {
		t.Errorf("Duplicate backing values should be allowed: %v", err)
	}
}
