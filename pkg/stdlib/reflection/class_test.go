package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionClass_BasicInfo tests basic class information retrieval
func TestReflectionClass_BasicInfo(t *testing.T) {
	// Create a test class
	classEntry := &types.ClassEntry{
		Name:      "MyNamespace\\MyClass",
		ShortName: "MyClass",
		Namespace: "MyNamespace",
		FileName:  "/path/to/file.php",
	}

	rc := NewReflectionClass("MyNamespace\\MyClass", classEntry)

	// Test GetName
	if name := rc.GetName(); name != "MyNamespace\\MyClass" {
		t.Errorf("GetName() = %s, want MyNamespace\\MyClass", name)
	}

	// Test GetShortName
	if shortName := rc.GetShortName(); shortName != "MyClass" {
		t.Errorf("GetShortName() = %s, want MyClass", shortName)
	}

	// Test GetNamespaceName
	if namespace := rc.GetNamespaceName(); namespace != "MyNamespace" {
		t.Errorf("GetNamespaceName() = %s, want MyNamespace", namespace)
	}

	// Test GetFileName
	if fileName := rc.GetFileName(); fileName != "/path/to/file.php" {
		t.Errorf("GetFileName() = %s, want /path/to/file.php", fileName)
	}
}

// TestReflectionClass_Modifiers tests class modifier checks
func TestReflectionClass_Modifiers(t *testing.T) {
	tests := []struct {
		name       string
		class      *types.ClassEntry
		isFinal    bool
		isAbstract bool
		isReadOnly bool
		isInterface bool
		isTrait    bool
		isEnum     bool
		isInstantiable bool
	}{
		{
			name: "regular class",
			class: &types.ClassEntry{
				Name: "RegularClass",
			},
			isInstantiable: true,
		},
		{
			name: "final class",
			class: &types.ClassEntry{
				Name:    "FinalClass",
				IsFinal: true,
			},
			isFinal: true,
			isInstantiable: true,
		},
		{
			name: "abstract class",
			class: &types.ClassEntry{
				Name:       "AbstractClass",
				IsAbstract: true,
			},
			isAbstract: true,
			isInstantiable: false,
		},
		{
			name: "readonly class",
			class: &types.ClassEntry{
				Name:       "ReadOnlyClass",
				IsReadOnly: true,
			},
			isReadOnly: true,
			isInstantiable: true,
		},
		{
			name: "interface",
			class: &types.ClassEntry{
				Name:        "MyInterface",
				IsInterface: true,
			},
			isInterface: true,
			isInstantiable: false,
		},
		{
			name: "trait",
			class: &types.ClassEntry{
				Name:    "MyTrait",
				IsTrait: true,
			},
			isTrait: true,
			isInstantiable: false,
		},
		{
			name: "enum",
			class: &types.ClassEntry{
				Name:   "MyEnum",
				IsEnum: true,
			},
			isEnum: true,
			isInstantiable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := NewReflectionClass(tt.class.Name, tt.class)

			if isFinal := rc.IsFinal(); isFinal != tt.isFinal {
				t.Errorf("IsFinal() = %v, want %v", isFinal, tt.isFinal)
			}

			if isAbstract := rc.IsAbstract(); isAbstract != tt.isAbstract {
				t.Errorf("IsAbstract() = %v, want %v", isAbstract, tt.isAbstract)
			}

			if isReadOnly := rc.IsReadOnly(); isReadOnly != tt.isReadOnly {
				t.Errorf("IsReadOnly() = %v, want %v", isReadOnly, tt.isReadOnly)
			}

			if isInterface := rc.IsInterface(); isInterface != tt.isInterface {
				t.Errorf("IsInterface() = %v, want %v", isInterface, tt.isInterface)
			}

			if isTrait := rc.IsTrait(); isTrait != tt.isTrait {
				t.Errorf("IsTrait() = %v, want %v", isTrait, tt.isTrait)
			}

			if isEnum := rc.IsEnum(); isEnum != tt.isEnum {
				t.Errorf("IsEnum() = %v, want %v", isEnum, tt.isEnum)
			}

			if isInstantiable := rc.IsInstantiable(); isInstantiable != tt.isInstantiable {
				t.Errorf("IsInstantiable() = %v, want %v", isInstantiable, tt.isInstantiable)
			}
		})
	}
}

// TestReflectionClass_Inheritance tests inheritance and interface checks
func TestReflectionClass_Inheritance(t *testing.T) {
	// Create parent class
	parentClass := &types.ClassEntry{
		Name: "ParentClass",
	}

	// Create grandparent class
	grandparentClass := &types.ClassEntry{
		Name: "GrandparentClass",
	}
	parentClass.ParentClass = grandparentClass

	// Create interfaces
	iface1 := &types.InterfaceEntry{
		Name: "Interface1",
	}
	iface2 := &types.InterfaceEntry{
		Name: "Interface2",
	}

	// Create child class
	childClass := &types.ClassEntry{
		Name:        "ChildClass",
		ParentClass: parentClass,
		Interfaces:  []*types.InterfaceEntry{iface1, iface2},
	}

	rc := NewReflectionClass("ChildClass", childClass)

	// Test GetParentClass
	parent := rc.GetParentClass()
	if parent == nil {
		t.Fatal("GetParentClass() returned nil")
	}
	if parent.GetName() != "ParentClass" {
		t.Errorf("GetParentClass().GetName() = %s, want ParentClass", parent.GetName())
	}

	// Test IsSubclassOf
	if !rc.IsSubclassOf("ParentClass") {
		t.Error("IsSubclassOf(ParentClass) = false, want true")
	}
	if !rc.IsSubclassOf("GrandparentClass") {
		t.Error("IsSubclassOf(GrandparentClass) = false, want true")
	}
	if rc.IsSubclassOf("UnrelatedClass") {
		t.Error("IsSubclassOf(UnrelatedClass) = true, want false")
	}

	// Test GetInterfaces
	interfaces := rc.GetInterfaces()
	if len(interfaces) != 2 {
		t.Fatalf("GetInterfaces() returned %d interfaces, want 2", len(interfaces))
	}

	// Test GetInterfaceNames
	ifaceNames := rc.GetInterfaceNames()
	if len(ifaceNames) != 2 {
		t.Fatalf("GetInterfaceNames() returned %d names, want 2", len(ifaceNames))
	}

	// Test ImplementsInterface
	if !rc.ImplementsInterface("Interface1") {
		t.Error("ImplementsInterface(Interface1) = false, want true")
	}
	if !rc.ImplementsInterface("Interface2") {
		t.Error("ImplementsInterface(Interface2) = false, want true")
	}
	if rc.ImplementsInterface("Interface3") {
		t.Error("ImplementsInterface(Interface3) = true, want false")
	}
}

// TestReflectionClass_Properties tests property enumeration and access
func TestReflectionClass_Properties(t *testing.T) {
	// Create class with properties
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"publicProp": {
				Visibility: types.VisibilityPublic,
				Type:       "string",
				HasDefault: true,
				Default:    types.NewString("default"),
			},
			"protectedProp": {
				Visibility: types.VisibilityProtected,
				Type:       "int",
			},
			"privateProp": {
				Visibility: types.VisibilityPrivate,
				Type:       "float",
			},
			"staticProp": {
				Visibility: types.VisibilityPublic,
				Type:       "bool",
				IsStatic:   true,
			},
		},
		StaticProperties: map[string]*types.Value{
			"staticProp": types.NewBool(true),
		},
		DefaultProperties: map[string]*types.Value{
			"publicProp": types.NewString("default"),
		},
	}

	rc := NewReflectionClass("TestClass", classEntry)

	// Test GetProperties with no filter
	allProps := rc.GetProperties(0)
	if len(allProps) != 4 {
		t.Errorf("GetProperties(0) returned %d properties, want 4", len(allProps))
	}

	// Test GetProperties with IS_PUBLIC filter
	publicProps := rc.GetProperties(IS_PUBLIC)
	if len(publicProps) != 2 { // publicProp and staticProp
		t.Errorf("GetProperties(IS_PUBLIC) returned %d properties, want 2", len(publicProps))
	}

	// Test GetProperties with IS_STATIC filter
	staticProps := rc.GetProperties(IS_STATIC)
	if len(staticProps) != 1 {
		t.Errorf("GetProperties(IS_STATIC) returned %d properties, want 1", len(staticProps))
	}

	// Test GetProperty
	prop, err := rc.GetProperty("publicProp")
	if err != nil {
		t.Fatalf("GetProperty(publicProp) error: %v", err)
	}
	if prop.GetName() != "publicProp" {
		t.Errorf("GetProperty(publicProp).GetName() = %s, want publicProp", prop.GetName())
	}

	// Test GetProperty with non-existent property
	_, err = rc.GetProperty("nonExistent")
	if err == nil {
		t.Error("GetProperty(nonExistent) should return error")
	}

	// Test HasProperty
	if !rc.HasProperty("publicProp") {
		t.Error("HasProperty(publicProp) = false, want true")
	}
	if rc.HasProperty("nonExistent") {
		t.Error("HasProperty(nonExistent) = true, want false")
	}

	// Test GetStaticProperties
	staticPropsMap := rc.GetStaticProperties()
	if len(staticPropsMap) != 1 {
		t.Errorf("GetStaticProperties() returned %d properties, want 1", len(staticPropsMap))
	}

	// Test GetDefaultProperties
	defaultPropsMap := rc.GetDefaultProperties()
	if len(defaultPropsMap) != 1 {
		t.Errorf("GetDefaultProperties() returned %d properties, want 1", len(defaultPropsMap))
	}
}

// TestReflectionClass_Methods tests method enumeration and access
func TestReflectionClass_Methods(t *testing.T) {
	// Create constructor
	constructor := &types.MethodDef{
		Visibility:    types.VisibilityPublic,
		IsConstructor: true,
		NumParams:     0,
		Parameters:    []*types.ParameterDef{},
	}

	// Create destructor
	destructor := &types.MethodDef{
		Visibility:   types.VisibilityPublic,
		IsDestructor: true,
	}

	// Create class with methods
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Methods: map[string]*types.MethodDef{
			"__construct": constructor,
			"__destruct":  destructor,
			"publicMethod": {
				Visibility: types.VisibilityPublic,
				NumParams:  2,
				Parameters: []*types.ParameterDef{
					{Name: "param1", Type: "string"},
					{Name: "param2", Type: "int"},
				},
				ReturnType: "bool",
			},
			"protectedMethod": {
				Visibility: types.VisibilityProtected,
			},
			"privateMethod": {
				Visibility: types.VisibilityPrivate,
			},
			"staticMethod": {
				Visibility: types.VisibilityPublic,
				IsStatic:   true,
			},
			"finalMethod": {
				Visibility: types.VisibilityPublic,
				IsFinal:    true,
			},
			"abstractMethod": {
				Visibility: types.VisibilityPublic,
				IsAbstract: true,
			},
		},
		Constructor: constructor,
		Destructor:  destructor,
	}

	rc := NewReflectionClass("TestClass", classEntry)

	// Test GetMethods with no filter
	allMethods := rc.GetMethods(0)
	if len(allMethods) != 8 {
		t.Errorf("GetMethods(0) returned %d methods, want 8", len(allMethods))
	}

	// Test GetMethods with IS_PUBLIC filter
	publicMethods := rc.GetMethods(IS_PUBLIC)
	if len(publicMethods) != 6 { // __construct, __destruct, publicMethod, staticMethod, finalMethod, abstractMethod
		t.Errorf("GetMethods(IS_PUBLIC) returned %d methods, want 6", len(publicMethods))
	}

	// Test GetMethods with IS_STATIC filter
	staticMethods := rc.GetMethods(IS_STATIC)
	if len(staticMethods) != 1 {
		t.Errorf("GetMethods(IS_STATIC) returned %d methods, want 1", len(staticMethods))
	}

	// Test GetMethod
	method, err := rc.GetMethod("publicMethod")
	if err != nil {
		t.Fatalf("GetMethod(publicMethod) error: %v", err)
	}
	if method.GetName() != "publicMethod" {
		t.Errorf("GetMethod(publicMethod).GetName() = %s, want publicMethod", method.GetName())
	}

	// Test GetMethod with non-existent method
	_, err = rc.GetMethod("nonExistent")
	if err == nil {
		t.Error("GetMethod(nonExistent) should return error")
	}

	// Test HasMethod
	if !rc.HasMethod("publicMethod") {
		t.Error("HasMethod(publicMethod) = false, want true")
	}
	if rc.HasMethod("nonExistent") {
		t.Error("HasMethod(nonExistent) = true, want false")
	}

	// Test GetConstructor
	ctor := rc.GetConstructor()
	if ctor == nil {
		t.Error("GetConstructor() returned nil")
	}

	// Test HasConstructor
	if !rc.HasConstructor() {
		t.Error("HasConstructor() = false, want true")
	}

	// Test GetDestructor
	dtor := rc.GetDestructor()
	if dtor == nil {
		t.Error("GetDestructor() returned nil")
	}
}

// TestReflectionClass_Constants tests constant access
func TestReflectionClass_Constants(t *testing.T) {
	// Create class with constants
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Constants: map[string]*types.ClassConstant{
			"CONST1": {
				Value: types.NewInt(42),
			},
			"CONST2": {
				Value: types.NewString("hello"),
			},
		},
	}

	rc := NewReflectionClass("TestClass", classEntry)

	// Test GetConstants
	constants := rc.GetConstants()
	if len(constants) != 2 {
		t.Errorf("GetConstants() returned %d constants, want 2", len(constants))
	}

	// Test GetConstant
	const1, err := rc.GetConstant("CONST1")
	if err != nil {
		t.Fatalf("GetConstant(CONST1) error: %v", err)
	}
	if const1.ToInt() != 42 {
		t.Errorf("GetConstant(CONST1).ToInt() = %d, want 42", const1.ToInt())
	}

	// Test GetConstant with non-existent constant
	_, err = rc.GetConstant("CONST3")
	if err == nil {
		t.Error("GetConstant(CONST3) should return error")
	}

	// Test HasConstant
	if !rc.HasConstant("CONST1") {
		t.Error("HasConstant(CONST1) = false, want true")
	}
	if rc.HasConstant("CONST3") {
		t.Error("HasConstant(CONST3) = true, want false")
	}
}

// TestReflectionClass_NewInstance tests object instantiation
func TestReflectionClass_NewInstance(t *testing.T) {
	// Create instantiable class
	classEntry := &types.ClassEntry{
		Name: "TestClass",
		Properties: map[string]*types.PropertyDef{
			"prop1": {
				Visibility: types.VisibilityPublic,
			},
		},
		DefaultProperties: map[string]*types.Value{
			"prop1": types.NewString("default"),
		},
	}

	rc := NewReflectionClass("TestClass", classEntry)

	// Test NewInstance
	obj, err := rc.NewInstance()
	if err != nil {
		t.Fatalf("NewInstance() error: %v", err)
	}
	if obj == nil {
		t.Fatal("NewInstance() returned nil")
	}
	if obj.ClassName != "TestClass" {
		t.Errorf("NewInstance().ClassName = %s, want TestClass", obj.ClassName)
	}

	// Test that default properties are initialized
	if len(obj.Properties) != 1 {
		t.Errorf("NewInstance() object has %d properties, want 1", len(obj.Properties))
	}

	// Test abstract class cannot be instantiated
	abstractClass := &types.ClassEntry{
		Name:       "AbstractClass",
		IsAbstract: true,
	}
	rcAbstract := NewReflectionClass("AbstractClass", abstractClass)
	_, err = rcAbstract.NewInstance()
	if err == nil {
		t.Error("NewInstance() on abstract class should return error")
	}
}

// TestReflectionClass_String tests string representation
func TestReflectionClass_String(t *testing.T) {
	tests := []struct {
		name     string
		class    *types.ClassEntry
		expected string
	}{
		{
			name: "regular class",
			class: &types.ClassEntry{
				Name: "MyClass",
			},
			expected: "Class [ MyClass ]",
		},
		{
			name: "final class",
			class: &types.ClassEntry{
				Name:    "FinalClass",
				IsFinal: true,
			},
			expected: "Class [ final FinalClass ]",
		},
		{
			name: "abstract class",
			class: &types.ClassEntry{
				Name:       "AbstractClass",
				IsAbstract: true,
			},
			expected: "Class [ abstract AbstractClass ]",
		},
		{
			name: "interface",
			class: &types.ClassEntry{
				Name:        "MyInterface",
				IsInterface: true,
			},
			expected: "Interface [ MyInterface ]",
		},
		{
			name: "trait",
			class: &types.ClassEntry{
				Name:    "MyTrait",
				IsTrait: true,
			},
			expected: "Trait [ MyTrait ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := NewReflectionClass(tt.class.Name, tt.class)
			str := rc.String()
			if str != tt.expected {
				t.Errorf("String() = %s, want %s", str, tt.expected)
			}
		})
	}
}

// TestReflectionClass_Traits tests trait retrieval
func TestReflectionClass_Traits(t *testing.T) {
	// Create traits
	trait1 := &types.TraitEntry{
		Name: "Trait1",
	}
	trait2 := &types.TraitEntry{
		Name: "Trait2",
	}

	// Create class using traits
	classEntry := &types.ClassEntry{
		Name:   "TestClass",
		Traits: []*types.TraitEntry{trait1, trait2},
	}

	rc := NewReflectionClass("TestClass", classEntry)

	// Test GetTraits
	traits := rc.GetTraits()
	if len(traits) != 2 {
		t.Fatalf("GetTraits() returned %d traits, want 2", len(traits))
	}

	// Test GetTraitNames
	traitNames := rc.GetTraitNames()
	if len(traitNames) != 2 {
		t.Fatalf("GetTraitNames() returned %d names, want 2", len(traitNames))
	}
	if traitNames[0] != "Trait1" || traitNames[1] != "Trait2" {
		t.Errorf("GetTraitNames() = %v, want [Trait1 Trait2]", traitNames)
	}
}

// TestReflectionClass_FromObject tests creating reflection from object instance
func TestReflectionClass_FromObject(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "MyClass",
	}

	obj := &types.Object{
		ClassName:  "MyClass",
		ClassEntry: classEntry,
		Properties: make(map[string]*types.Property),
	}

	rc := NewReflectionClassFromObject(obj)

	if rc.GetName() != "MyClass" {
		t.Errorf("GetName() = %s, want MyClass", rc.GetName())
	}
}
