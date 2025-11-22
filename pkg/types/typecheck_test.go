package types

import "testing"

// ============================================================================
// Property Type Hints Tests
// ============================================================================

func TestTypeCheck_PropertyTypeHint(t *testing.T) {
	class := NewClassEntry("Test")

	// Property with int type hint
	class.Properties["count"] = &PropertyDef{
		Name:       "count",
		Visibility: VisibilityPublic,
		Type:       "int",
	}

	// Property with string type hint
	class.Properties["name"] = &PropertyDef{
		Name:       "name",
		Visibility: VisibilityPublic,
		Type:       "string",
	}

	// Check that type hints are stored
	if class.Properties["count"].Type != "int" {
		t.Error("count property should have int type")
	}

	if class.Properties["name"].Type != "string" {
		t.Error("name property should have string type")
	}
}

func TestTypeCheck_NullablePropertyType(t *testing.T) {
	class := NewClassEntry("Test")

	// Nullable type: ?string
	class.Properties["nullableName"] = &PropertyDef{
		Name:       "nullableName",
		Visibility: VisibilityPublic,
		Type:       "?string",
	}

	// Check nullable type parsing
	typeInfo := ParseType(class.Properties["nullableName"].Type)
	if !typeInfo.IsNullable {
		t.Error("?string should be nullable")
	}

	if typeInfo.BaseType != "string" {
		t.Errorf("Base type should be string, got %s", typeInfo.BaseType)
	}
}

func TestTypeCheck_UnionPropertyType(t *testing.T) {
	// PHP 8.0+ union types: int|string
	class := NewClassEntry("Test")

	class.Properties["mixed"] = &PropertyDef{
		Name:       "mixed",
		Visibility: VisibilityPublic,
		Type:       "int|string",
	}

	typeInfo := ParseType(class.Properties["mixed"].Type)
	if !typeInfo.IsUnion {
		t.Error("int|string should be a union type")
	}

	if len(typeInfo.UnionTypes) != 2 {
		t.Errorf("Expected 2 union types, got %d", len(typeInfo.UnionTypes))
	}
}

// ============================================================================
// Parameter Type Checking Tests
// ============================================================================

func TestTypeCheck_ParameterTypes(t *testing.T) {
	method := &MethodDef{
		Name:       "setData",
		Visibility: VisibilityPublic,
		NumParams:  2,
		Parameters: []*ParameterDef{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "name",
				Type: "string",
			},
		},
	}

	if method.Parameters[0].Type != "int" {
		t.Error("First parameter should have int type")
	}

	if method.Parameters[1].Type != "string" {
		t.Error("Second parameter should have string type")
	}
}

func TestTypeCheck_ClassTypeParameter(t *testing.T) {
	method := &MethodDef{
		Name:       "process",
		Visibility: VisibilityPublic,
		NumParams:  1,
		Parameters: []*ParameterDef{
			{
				Name: "obj",
				Type: "MyClass",
			},
		},
	}

	// Check class type hint
	typeInfo := ParseType(method.Parameters[0].Type)
	if typeInfo.BaseType != "MyClass" {
		t.Errorf("Expected MyClass type, got %s", typeInfo.BaseType)
	}

	if !typeInfo.IsClass {
		t.Error("MyClass should be recognized as a class type")
	}
}

func TestTypeCheck_BuiltinTypes(t *testing.T) {
	// Test all built-in PHP types
	builtinTypes := []string{
		"int", "string", "float", "bool", "array",
		"object", "callable", "iterable", "mixed", "void", "never",
	}

	for _, typeName := range builtinTypes {
		typeInfo := ParseType(typeName)
		if !typeInfo.IsBuiltin {
			t.Errorf("%s should be recognized as built-in type", typeName)
		}

		if typeInfo.BaseType != typeName {
			t.Errorf("Base type should be %s, got %s", typeName, typeInfo.BaseType)
		}
	}
}

// ============================================================================
// Return Type Checking Tests
// ============================================================================

func TestTypeCheck_ReturnType(t *testing.T) {
	method := &MethodDef{
		Name:       "getId",
		Visibility: VisibilityPublic,
		ReturnType: "int",
	}

	if method.ReturnType != "int" {
		t.Error("Return type should be int")
	}
}

func TestTypeCheck_VoidReturnType(t *testing.T) {
	method := &MethodDef{
		Name:       "doSomething",
		Visibility: VisibilityPublic,
		ReturnType: "void",
	}

	typeInfo := ParseType(method.ReturnType)
	if typeInfo.BaseType != "void" {
		t.Error("Return type should be void")
	}
}

func TestTypeCheck_NullableReturnType(t *testing.T) {
	method := &MethodDef{
		Name:       "findUser",
		Visibility: VisibilityPublic,
		ReturnType: "?User",
	}

	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsNullable {
		t.Error("?User should be nullable")
	}

	if typeInfo.BaseType != "User" {
		t.Errorf("Base type should be User, got %s", typeInfo.BaseType)
	}
}

// ============================================================================
// Type Variance Tests (Covariance/Contravariance)
// ============================================================================

func TestTypeCheck_ReturnTypeCovariance(t *testing.T) {
	// Parent method returns Animal
	parentMethod := &MethodDef{
		Name:       "getAnimal",
		Visibility: VisibilityPublic,
		ReturnType: "Animal",
	}

	// Child method returns Dog (subclass of Animal) - covariant
	childMethod := &MethodDef{
		Name:       "getAnimal",
		Visibility: VisibilityPublic,
		ReturnType: "Dog",
	}

	// This should be valid (covariance)
	err := ValidateReturnTypeCovariance(parentMethod, childMethod, "Animal", "Dog")
	if err != nil {
		t.Errorf("Return type covariance should be valid: %v", err)
	}
}

func TestTypeCheck_ReturnTypeInvariance(t *testing.T) {
	// Test with built-in types which must match exactly
	parentMethod := &MethodDef{
		Name:       "getValue",
		Visibility: VisibilityPublic,
		ReturnType: "int",
	}

	// Child returns string - NOT allowed
	childMethod := &MethodDef{
		Name:       "getValue",
		Visibility: VisibilityPublic,
		ReturnType: "string",
	}

	// This should be invalid (built-in types must match exactly)
	err := ValidateReturnTypeCovariance(parentMethod, childMethod, "int", "string")
	if err == nil {
		t.Fatal("Changing return type from int to string should be invalid")
	}
}

func TestTypeCheck_ParameterTypeContravariance(t *testing.T) {
	// Parent method accepts Dog parameter
	parentParam := &ParameterDef{
		Name: "dog",
		Type: "Dog",
	}

	// Child method accepts Animal parameter (superclass) - contravariant
	childParam := &ParameterDef{
		Name: "dog",
		Type: "Animal",
	}

	// This should be valid (contravariance for parameters)
	err := ValidateParameterTypeContravariance(parentParam, childParam, "Dog", "Animal")
	if err != nil {
		t.Errorf("Parameter type contravariance should be valid: %v", err)
	}
}

// ============================================================================
// Type Compatibility Tests
// ============================================================================

func TestTypeCheck_TypeCompatibility(t *testing.T) {
	// int is compatible with int
	if !IsTypeCompatible("int", "int") {
		t.Error("int should be compatible with int")
	}

	// int is NOT compatible with string
	if IsTypeCompatible("int", "string") {
		t.Error("int should not be compatible with string")
	}
}

func TestTypeCheck_NullableCompatibility(t *testing.T) {
	// ?string is compatible with string (nullable accepts non-nullable)
	if !IsTypeCompatible("?string", "string") {
		t.Error("?string should accept string value")
	}

	// string is NOT compatible with ?string (non-nullable doesn't accept null)
	// This is checked at runtime, not compile time
}

func TestTypeCheck_MixedType(t *testing.T) {
	// mixed accepts any type
	types := []string{"int", "string", "float", "bool", "array", "object", "MyClass"}

	for _, typ := range types {
		if !IsTypeCompatible("mixed", typ) {
			t.Errorf("mixed should accept %s", typ)
		}
	}
}

func TestTypeCheck_UnionTypeCompatibility(t *testing.T) {
	// int|string accepts int
	if !IsTypeCompatible("int|string", "int") {
		t.Error("int|string should accept int")
	}

	// int|string accepts string
	if !IsTypeCompatible("int|string", "string") {
		t.Error("int|string should accept string")
	}

	// int|string does NOT accept float
	if IsTypeCompatible("int|string", "float") {
		t.Error("int|string should not accept float")
	}
}

// ============================================================================
// Type Validation Tests
// ============================================================================

func TestTypeCheck_ValidatePropertyType(t *testing.T) {
	prop := &PropertyDef{
		Name:    "age",
		Type:    "int",
		Default: NewInt(25),
	}

	// Valid: int value for int property
	err := ValidatePropertyValue(prop, NewInt(30))
	if err != nil {
		t.Errorf("int value should be valid for int property: %v", err)
	}

	// Invalid: string value for int property
	err = ValidatePropertyValue(prop, NewString("abc"))
	if err == nil {
		t.Fatal("string value should be invalid for int property")
	}
}

func TestTypeCheck_ValidateNullableProperty(t *testing.T) {
	prop := &PropertyDef{
		Name:    "name",
		Type:    "?string",
		Default: nil,
	}

	// Valid: null for nullable property
	err := ValidatePropertyValue(prop, nil)
	if err != nil {
		t.Errorf("null should be valid for ?string property: %v", err)
	}

	// Valid: string for nullable property
	err = ValidatePropertyValue(prop, NewString("John"))
	if err != nil {
		t.Errorf("string should be valid for ?string property: %v", err)
	}
}

func TestTypeCheck_ReadonlyPropertyType(t *testing.T) {
	// Readonly properties must have type (PHP 8.1+)
	prop := &PropertyDef{
		Name:       "id",
		IsReadOnly: true,
		Type:       "int",
	}

	err := ValidateReadonlyProperty(prop)
	if err != nil {
		t.Errorf("Readonly property with type should be valid: %v", err)
	}

	// Readonly without type is invalid
	prop2 := &PropertyDef{
		Name:       "code",
		IsReadOnly: true,
		Type:       "",
	}

	err = ValidateReadonlyProperty(prop2)
	if err == nil {
		t.Fatal("Readonly property must have type hint")
	}
}

// ============================================================================
// Array and Iterable Type Tests
// ============================================================================

func TestTypeCheck_ArrayType(t *testing.T) {
	typeInfo := ParseType("array")

	if !typeInfo.IsBuiltin {
		t.Error("array should be built-in type")
	}

	if typeInfo.BaseType != "array" {
		t.Error("Base type should be array")
	}
}

func TestTypeCheck_IterableType(t *testing.T) {
	// iterable accepts array and Traversable objects
	if !IsTypeCompatible("iterable", "array") {
		t.Error("iterable should accept array")
	}
}

// ============================================================================
// Self, Parent, Static Type Tests
// ============================================================================

func TestTypeCheck_SelfType(t *testing.T) {
	method := &MethodDef{
		Name:       "getInstance",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "self",
	}

	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsSelf {
		t.Error("'self' should be recognized as self type")
	}
}

func TestTypeCheck_ParentType(t *testing.T) {
	method := &MethodDef{
		Name:       "callParent",
		Visibility: VisibilityPublic,
		ReturnType: "parent",
	}

	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsParent {
		t.Error("'parent' should be recognized as parent type")
	}
}

func TestTypeCheck_StaticType(t *testing.T) {
	// 'static' return type (late static binding)
	method := &MethodDef{
		Name:       "create",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		ReturnType: "static",
	}

	typeInfo := ParseType(method.ReturnType)
	if !typeInfo.IsStatic {
		t.Error("'static' should be recognized as static type")
	}
}
