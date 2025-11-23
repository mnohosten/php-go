package runtime

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestAttributeInstance_Creation tests attribute instance creation
func TestAttributeInstance_Creation(t *testing.T) {
	attr := NewAttributeInstance("MyAttribute", AttributeTargetClass)

	if attr.GetName() != "MyAttribute" {
		t.Errorf("GetName() = %s, want MyAttribute", attr.GetName())
	}

	if attr.Target != AttributeTargetClass {
		t.Errorf("Target = %v, want AttributeTargetClass", attr.Target)
	}

	if len(attr.GetArguments()) != 0 {
		t.Error("New attribute should have no arguments")
	}

	if len(attr.GetNamedArguments()) != 0 {
		t.Error("New attribute should have no named arguments")
	}
}

// TestAttributeInstance_Arguments tests adding arguments to attributes
func TestAttributeInstance_Arguments(t *testing.T) {
	attr := NewAttributeInstance("Route", AttributeTargetMethod)

	// Add positional arguments
	attr.AddArgument(types.NewString("/path"))
	attr.AddArgument(types.NewString("GET"))

	args := attr.GetArguments()
	if len(args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(args))
	}

	if args[0].ToString() != "/path" {
		t.Errorf("First argument = %s, want /path", args[0].ToString())
	}

	if args[1].ToString() != "GET" {
		t.Errorf("Second argument = %s, want GET", args[1].ToString())
	}
}

// TestAttributeInstance_NamedArguments tests named arguments
func TestAttributeInstance_NamedArguments(t *testing.T) {
	attr := NewAttributeInstance("Route", AttributeTargetMethod)

	// Add named arguments
	methods := types.NewArray(types.NewArrayWithCapacity(0))
	attr.AddNamedArgument("methods", methods)
	attr.AddNamedArgument("name", types.NewString("api.users"))

	named := attr.GetNamedArguments()
	if len(named) != 2 {
		t.Fatalf("Expected 2 named arguments, got %d", len(named))
	}

	if _, ok := named["methods"]; !ok {
		t.Error("Named argument 'methods' not found")
	}

	if name, ok := named["name"]; !ok {
		t.Error("Named argument 'name' not found")
	} else if name.ToString() != "api.users" {
		t.Errorf("Named argument 'name' = %s, want api.users", name.ToString())
	}
}

// TestAttributeTarget_String tests target string representation
func TestAttributeTarget_String(t *testing.T) {
	tests := []struct {
		target AttributeTarget
		want   string
	}{
		{AttributeTargetClass, "CLASS"},
		{AttributeTargetFunction, "FUNCTION"},
		{AttributeTargetMethod, "METHOD"},
		{AttributeTargetProperty, "PROPERTY"},
		{AttributeTargetConstant, "CONSTANT"},
		{AttributeTargetParameter, "PARAMETER"},
		{AttributeTargetAll, "ALL"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.target.String(); got != tt.want {
				t.Errorf("String() = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestAttributeDefinition_Creation tests attribute definition creation
func TestAttributeDefinition_Creation(t *testing.T) {
	def := NewAttributeDefinition("MyAttribute", AttributeTargetClass, false)

	if def.Name != "MyAttribute" {
		t.Errorf("Name = %s, want MyAttribute", def.Name)
	}

	if def.Targets != AttributeTargetClass {
		t.Errorf("Targets = %v, want AttributeTargetClass", def.Targets)
	}

	if def.Repeatable {
		t.Error("Repeatable should be false")
	}
}

// TestAttributeDefinition_ValidateTarget tests target validation
func TestAttributeDefinition_ValidateTarget(t *testing.T) {
	// Attribute that can only be used on methods
	methodOnlyDef := NewAttributeDefinition("MethodOnly", AttributeTargetMethod, false)

	// Should succeed for method
	if err := methodOnlyDef.ValidateTarget(AttributeTargetMethod); err != nil {
		t.Errorf("ValidateTarget(METHOD) failed: %v", err)
	}

	// Should fail for class
	if err := methodOnlyDef.ValidateTarget(AttributeTargetClass); err == nil {
		t.Error("ValidateTarget(CLASS) should fail for method-only attribute")
	}

	// Attribute that can be used anywhere
	allDef := NewAttributeDefinition("CanUseAnywhere", AttributeTargetAll, false)

	// Should succeed for any target
	targets := []AttributeTarget{
		AttributeTargetClass,
		AttributeTargetMethod,
		AttributeTargetProperty,
		AttributeTargetFunction,
		AttributeTargetParameter,
		AttributeTargetConstant,
	}

	for _, target := range targets {
		if err := allDef.ValidateTarget(target); err != nil {
			t.Errorf("ValidateTarget(%s) failed for ALL attribute: %v", target.String(), err)
		}
	}
}

// TestAttributeRegistry_Creation tests registry creation
func TestAttributeRegistry_Creation(t *testing.T) {
	registry := NewAttributeRegistry()

	if registry.definitions == nil {
		t.Error("Registry definitions map should be initialized")
	}

	// Check that built-in attributes are registered
	if _, ok := registry.Get(AttributeDeprecated); !ok {
		t.Error("Built-in Deprecated attribute should be registered")
	}
}

// TestAttributeRegistry_RegisterAndGet tests registering and retrieving definitions
func TestAttributeRegistry_RegisterAndGet(t *testing.T) {
	registry := NewAttributeRegistry()

	// Register custom attribute
	customDef := NewAttributeDefinition("CustomAttribute", AttributeTargetClass, true)
	registry.Register(customDef)

	// Retrieve it
	def, ok := registry.Get("CustomAttribute")
	if !ok {
		t.Fatal("Custom attribute should be retrievable")
	}

	if def.Name != "CustomAttribute" {
		t.Errorf("Retrieved definition name = %s, want CustomAttribute", def.Name)
	}

	if !def.Repeatable {
		t.Error("Custom attribute should be repeatable")
	}
}

// TestAttributeRegistry_Validate tests attribute validation
func TestAttributeRegistry_Validate(t *testing.T) {
	registry := NewAttributeRegistry()

	// Test with built-in attribute (Override - methods only)
	overrideAttr := NewAttributeInstance(AttributeOverride, AttributeTargetMethod)
	if err := registry.Validate(overrideAttr); err != nil {
		t.Errorf("Validate(Override on METHOD) failed: %v", err)
	}

	// Try to use Override on a class (should fail)
	overrideOnClass := NewAttributeInstance(AttributeOverride, AttributeTargetClass)
	if err := registry.Validate(overrideOnClass); err == nil {
		t.Error("Validate(Override on CLASS) should fail")
	}

	// Test with unknown attribute (should succeed - user-defined)
	unknownAttr := NewAttributeInstance("UnknownAttribute", AttributeTargetClass)
	if err := registry.Validate(unknownAttr); err != nil {
		t.Errorf("Validate(unknown attribute) failed: %v", err)
	}
}

// TestAttributeRegistry_BuiltinAttributes tests built-in attributes
func TestAttributeRegistry_BuiltinAttributes(t *testing.T) {
	registry := NewAttributeRegistry()

	builtins := []struct {
		name   string
		target AttributeTarget
	}{
		{AttributeDeprecated, AttributeTargetAll},
		{AttributeReturnTypeWillChange, AttributeTargetMethod},
		{AttributeAllowDynamicProperties, AttributeTargetClass},
		{AttributeSensitiveParameter, AttributeTargetParameter},
		{AttributeOverride, AttributeTargetMethod},
	}

	for _, builtin := range builtins {
		t.Run(builtin.name, func(t *testing.T) {
			_, ok := registry.Get(builtin.name)
			if !ok {
				t.Fatalf("Built-in attribute %s not found", builtin.name)
			}

			// Verify target
			instance := NewAttributeInstance(builtin.name, builtin.target)
			if err := registry.Validate(instance); err != nil {
				t.Errorf("Built-in attribute %s validation failed: %v", builtin.name, err)
			}
		})
	}
}

// TestAttributeStorage_Creation tests storage creation
func TestAttributeStorage_Creation(t *testing.T) {
	storage := NewAttributeStorage()

	if storage.ClassAttributes == nil {
		t.Error("ClassAttributes map should be initialized")
	}
	if storage.FunctionAttributes == nil {
		t.Error("FunctionAttributes map should be initialized")
	}
	if storage.MethodAttributes == nil {
		t.Error("MethodAttributes map should be initialized")
	}
	if storage.PropertyAttributes == nil {
		t.Error("PropertyAttributes map should be initialized")
	}
	if storage.ParameterAttributes == nil {
		t.Error("ParameterAttributes map should be initialized")
	}
	if storage.ConstantAttributes == nil {
		t.Error("ConstantAttributes map should be initialized")
	}
}

// TestAttributeStorage_ClassAttributes tests class attribute storage
func TestAttributeStorage_ClassAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr1 := NewAttributeInstance("Attribute1", AttributeTargetClass)
	attr2 := NewAttributeInstance("Attribute2", AttributeTargetClass)

	storage.AddClassAttribute("MyClass", attr1)
	storage.AddClassAttribute("MyClass", attr2)

	attrs := storage.GetClassAttributes("MyClass")
	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}

	if attrs[0].GetName() != "Attribute1" {
		t.Errorf("First attribute = %s, want Attribute1", attrs[0].GetName())
	}
	if attrs[1].GetName() != "Attribute2" {
		t.Errorf("Second attribute = %s, want Attribute2", attrs[1].GetName())
	}

	// Test non-existent class
	noAttrs := storage.GetClassAttributes("NonExistent")
	if noAttrs != nil {
		t.Error("Non-existent class should return nil")
	}
}

// TestAttributeStorage_FunctionAttributes tests function attribute storage
func TestAttributeStorage_FunctionAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr := NewAttributeInstance("Deprecated", AttributeTargetFunction)
	storage.AddFunctionAttribute("myFunction", attr)

	attrs := storage.GetFunctionAttributes("myFunction")
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].GetName() != "Deprecated" {
		t.Errorf("Attribute name = %s, want Deprecated", attrs[0].GetName())
	}
}

// TestAttributeStorage_MethodAttributes tests method attribute storage
func TestAttributeStorage_MethodAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr := NewAttributeInstance("Route", AttributeTargetMethod)
	attr.AddArgument(types.NewString("/users"))

	storage.AddMethodAttribute("UserController", "index", attr)

	attrs := storage.GetMethodAttributes("UserController", "index")
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].GetName() != "Route" {
		t.Errorf("Attribute name = %s, want Route", attrs[0].GetName())
	}

	args := attrs[0].GetArguments()
	if len(args) != 1 || args[0].ToString() != "/users" {
		t.Error("Attribute argument not preserved")
	}
}

// TestAttributeStorage_PropertyAttributes tests property attribute storage
func TestAttributeStorage_PropertyAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr := NewAttributeInstance("Column", AttributeTargetProperty)
	storage.AddPropertyAttribute("User", "email", attr)

	attrs := storage.GetPropertyAttributes("User", "email")
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].GetName() != "Column" {
		t.Errorf("Attribute name = %s, want Column", attrs[0].GetName())
	}
}

// TestAttributeStorage_ParameterAttributes tests parameter attribute storage
func TestAttributeStorage_ParameterAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr := NewAttributeInstance(AttributeSensitiveParameter, AttributeTargetParameter)
	storage.AddParameterAttribute("User", "login", "password", attr)

	attrs := storage.GetParameterAttributes("User", "login", "password")
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].GetName() != AttributeSensitiveParameter {
		t.Errorf("Attribute name = %s, want %s", attrs[0].GetName(), AttributeSensitiveParameter)
	}
}

// TestAttributeStorage_ConstantAttributes tests constant attribute storage
func TestAttributeStorage_ConstantAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	attr := NewAttributeInstance("Description", AttributeTargetConstant)
	attr.AddArgument(types.NewString("API version constant"))

	storage.AddConstantAttribute("Api", "VERSION", attr)

	attrs := storage.GetConstantAttributes("Api", "VERSION")
	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].GetName() != "Description" {
		t.Errorf("Attribute name = %s, want Description", attrs[0].GetName())
	}
}

// TestAttributeStorage_MultipleAttributes tests storing multiple attributes
func TestAttributeStorage_MultipleAttributes(t *testing.T) {
	storage := NewAttributeStorage()

	// Add multiple attributes to the same class
	attr1 := NewAttributeInstance("Attr1", AttributeTargetClass)
	attr2 := NewAttributeInstance("Attr2", AttributeTargetClass)
	attr3 := NewAttributeInstance("Attr3", AttributeTargetClass)

	storage.AddClassAttribute("TestClass", attr1)
	storage.AddClassAttribute("TestClass", attr2)
	storage.AddClassAttribute("TestClass", attr3)

	attrs := storage.GetClassAttributes("TestClass")
	if len(attrs) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(attrs))
	}

	// Verify order is preserved
	if attrs[0].GetName() != "Attr1" || attrs[1].GetName() != "Attr2" || attrs[2].GetName() != "Attr3" {
		t.Error("Attribute order not preserved")
	}
}

// TestAttributeStorage_ComplexScenario tests a complex real-world scenario
func TestAttributeStorage_ComplexScenario(t *testing.T) {
	storage := NewAttributeStorage()
	registry := NewAttributeRegistry()

	// Class with AllowDynamicProperties
	classAttr := NewAttributeInstance(AttributeAllowDynamicProperties, AttributeTargetClass)
	if err := registry.Validate(classAttr); err != nil {
		t.Fatalf("Class attribute validation failed: %v", err)
	}
	storage.AddClassAttribute("DynamicClass", classAttr)

	// Method with Route attribute
	routeAttr := NewAttributeInstance("Route", AttributeTargetMethod)
	routeAttr.AddArgument(types.NewString("/api/users"))
	routeAttr.AddNamedArgument("methods", types.NewArray(types.NewArrayWithCapacity(0)))
	storage.AddMethodAttribute("UserController", "list", routeAttr)

	// Parameter with SensitiveParameter
	sensitiveAttr := NewAttributeInstance(AttributeSensitiveParameter, AttributeTargetParameter)
	if err := registry.Validate(sensitiveAttr); err != nil {
		t.Fatalf("Parameter attribute validation failed: %v", err)
	}
	storage.AddParameterAttribute("Auth", "login", "password", sensitiveAttr)

	// Property with custom attribute
	columnAttr := NewAttributeInstance("Column", AttributeTargetProperty)
	columnAttr.AddNamedArgument("name", types.NewString("user_email"))
	columnAttr.AddNamedArgument("nullable", types.NewBool(false))
	storage.AddPropertyAttribute("User", "email", columnAttr)

	// Verify all stored correctly
	if classAttrs := storage.GetClassAttributes("DynamicClass"); len(classAttrs) != 1 {
		t.Error("Class attribute not stored")
	}

	if methodAttrs := storage.GetMethodAttributes("UserController", "list"); len(methodAttrs) != 1 {
		t.Error("Method attribute not stored")
	} else {
		if len(methodAttrs[0].GetArguments()) != 1 {
			t.Error("Method attribute arguments not stored")
		}
		if len(methodAttrs[0].GetNamedArguments()) != 1 {
			t.Error("Method attribute named arguments not stored")
		}
	}

	if paramAttrs := storage.GetParameterAttributes("Auth", "login", "password"); len(paramAttrs) != 1 {
		t.Error("Parameter attribute not stored")
	}

	if propAttrs := storage.GetPropertyAttributes("User", "email"); len(propAttrs) != 1 {
		t.Error("Property attribute not stored")
	} else {
		named := propAttrs[0].GetNamedArguments()
		if len(named) != 2 {
			t.Errorf("Property attribute should have 2 named arguments, got %d", len(named))
		}
	}
}
