package types

import "testing"

// ============================================================================
// __get / __set Tests (Dynamic Property Access)
// ============================================================================

func TestMagic_GetSet(t *testing.T) {
	// Create class with __get and __set
	class := NewClassEntry("DynamicProps")

	class.MagicMethods["__get"] = &MethodDef{
		Name:       "__get",
		Visibility: VisibilityPublic,
		NumParams:  1, // property name
		IsMagic:    true,
	}

	class.MagicMethods["__set"] = &MethodDef{
		Name:       "__set",
		Visibility: VisibilityPublic,
		NumParams:  2, // property name, value
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__get") {
		t.Error("Class should have __get magic method")
	}

	if !class.HasMagicMethod("__set") {
		t.Error("Class should have __set magic method")
	}
}

func TestMagic_GetMethod(t *testing.T) {
	class := NewClassEntry("Test")
	class.MagicMethods["__get"] = &MethodDef{
		Name:       "__get",
		Visibility: VisibilityPublic,
		NumParams:  1,
		IsMagic:    true,
	}

	method := class.GetMagicMethod("__get")
	if method == nil {
		t.Fatal("GetMagicMethod should return __get method")
	}

	if method.Name != "__get" {
		t.Errorf("Expected method name '__get', got '%s'", method.Name)
	}
}

// ============================================================================
// __isset / __unset Tests
// ============================================================================

func TestMagic_IssetUnset(t *testing.T) {
	class := NewClassEntry("Test")

	class.MagicMethods["__isset"] = &MethodDef{
		Name:       "__isset",
		Visibility: VisibilityPublic,
		NumParams:  1, // property name
		IsMagic:    true,
	}

	class.MagicMethods["__unset"] = &MethodDef{
		Name:       "__unset",
		Visibility: VisibilityPublic,
		NumParams:  1, // property name
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__isset") {
		t.Error("Class should have __isset magic method")
	}

	if !class.HasMagicMethod("__unset") {
		t.Error("Class should have __unset magic method")
	}
}

// ============================================================================
// __call / __callStatic Tests (Method Overloading)
// ============================================================================

func TestMagic_Call(t *testing.T) {
	class := NewClassEntry("Test")

	class.MagicMethods["__call"] = &MethodDef{
		Name:       "__call",
		Visibility: VisibilityPublic,
		NumParams:  2, // method name, arguments array
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__call") {
		t.Error("Class should have __call magic method")
	}

	method := class.GetMagicMethod("__call")
	if method.NumParams != 2 {
		t.Errorf("__call should have 2 parameters, got %d", method.NumParams)
	}
}

func TestMagic_CallStatic(t *testing.T) {
	class := NewClassEntry("Test")

	class.MagicMethods["__callStatic"] = &MethodDef{
		Name:       "__callStatic",
		Visibility: VisibilityPublic,
		IsStatic:   true,
		NumParams:  2, // method name, arguments array
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__callStatic") {
		t.Error("Class should have __callStatic magic method")
	}

	method := class.GetMagicMethod("__callStatic")
	if !method.IsStatic {
		t.Error("__callStatic should be static")
	}
}

// ============================================================================
// __toString Tests
// ============================================================================

func TestMagic_ToString(t *testing.T) {
	class := NewClassEntry("Stringable")

	class.MagicMethods["__toString"] = &MethodDef{
		Name:       "__toString",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
		ReturnType: "string",
	}

	if !class.HasMagicMethod("__toString") {
		t.Error("Class should have __toString magic method")
	}

	method := class.GetMagicMethod("__toString")
	if method.NumParams != 0 {
		t.Error("__toString should have no parameters")
	}
}

// ============================================================================
// __invoke Tests (Callable Objects)
// ============================================================================

func TestMagic_Invoke(t *testing.T) {
	class := NewClassEntry("Callable")

	class.MagicMethods["__invoke"] = &MethodDef{
		Name:       "__invoke",
		Visibility: VisibilityPublic,
		NumParams:  1, // Can have any number of parameters
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__invoke") {
		t.Error("Class should have __invoke magic method")
	}

	// Object with __invoke can be called as a function
}

// ============================================================================
// __clone Tests
// ============================================================================

func TestMagic_Clone(t *testing.T) {
	class := NewClassEntry("Cloneable")

	class.MagicMethods["__clone"] = &MethodDef{
		Name:       "__clone",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__clone") {
		t.Error("Class should have __clone magic method")
	}

	// __clone is called after object is cloned
	method := class.GetMagicMethod("__clone")
	if method.NumParams != 0 {
		t.Error("__clone should have no parameters")
	}
}

// ============================================================================
// __debugInfo Tests
// ============================================================================

func TestMagic_DebugInfo(t *testing.T) {
	class := NewClassEntry("Debuggable")

	class.MagicMethods["__debugInfo"] = &MethodDef{
		Name:       "__debugInfo",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
		ReturnType: "array",
	}

	if !class.HasMagicMethod("__debugInfo") {
		t.Error("Class should have __debugInfo magic method")
	}

	// __debugInfo returns array of properties to display in var_dump
}

// ============================================================================
// __serialize / __unserialize Tests (PHP 7.4+)
// ============================================================================

func TestMagic_Serialize(t *testing.T) {
	class := NewClassEntry("Serializable")

	class.MagicMethods["__serialize"] = &MethodDef{
		Name:       "__serialize",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
		ReturnType: "array",
	}

	if !class.HasMagicMethod("__serialize") {
		t.Error("Class should have __serialize magic method")
	}

	// __serialize returns array representation for serialization
}

func TestMagic_Unserialize(t *testing.T) {
	class := NewClassEntry("Serializable")

	class.MagicMethods["__unserialize"] = &MethodDef{
		Name:       "__unserialize",
		Visibility: VisibilityPublic,
		NumParams:  1, // array data
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__unserialize") {
		t.Error("Class should have __unserialize magic method")
	}

	// __unserialize restores object from array
}

// ============================================================================
// __sleep / __wakeup Tests (Legacy Serialization)
// ============================================================================

func TestMagic_Sleep(t *testing.T) {
	class := NewClassEntry("Sleeper")

	class.MagicMethods["__sleep"] = &MethodDef{
		Name:       "__sleep",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
		ReturnType: "array",
	}

	if !class.HasMagicMethod("__sleep") {
		t.Error("Class should have __sleep magic method")
	}

	// __sleep returns array of property names to serialize
}

func TestMagic_Wakeup(t *testing.T) {
	class := NewClassEntry("Sleeper")

	class.MagicMethods["__wakeup"] = &MethodDef{
		Name:       "__wakeup",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
	}

	if !class.HasMagicMethod("__wakeup") {
		t.Error("Class should have __wakeup magic method")
	}

	// __wakeup is called after unserialization
}

// ============================================================================
// Magic Method Visibility Tests
// ============================================================================

func TestMagic_VisibilityEnforcement(t *testing.T) {
	// Most magic methods must be public
	class := NewClassEntry("Test")

	// __construct can be private (for singletons)
	class.Constructor = &MethodDef{
		Name:         "__construct",
		Visibility:   VisibilityPrivate,
		IsConstructor: true,
		IsMagic:      true,
	}

	// This is valid for singleton pattern
	if err := class.ValidateMagicMethods(); err != nil {
		t.Errorf("Private __construct should be allowed: %v", err)
	}
}

func TestMagic_NonPublicMagicMethod(t *testing.T) {
	// __toString must be public
	class := NewClassEntry("Test")

	class.MagicMethods["__toString"] = &MethodDef{
		Name:       "__toString",
		Visibility: VisibilityPrivate, // Invalid
		IsMagic:    true,
	}

	err := class.ValidateMagicMethods()
	if err == nil {
		t.Fatal("__toString must be public")
	}
}

// ============================================================================
// Magic Method Inheritance Tests
// ============================================================================

func TestMagic_InheritedMagicMethods(t *testing.T) {
	// Parent class with __get
	parent := NewClassEntry("Parent")
	parent.MagicMethods["__get"] = &MethodDef{
		Name:       "__get",
		Visibility: VisibilityPublic,
		NumParams:  1,
		IsMagic:    true,
	}

	// Child class inherits __get
	child := NewClassEntry("Child")
	child.ParentClass = parent

	// After inheritance, child should have access to __get
	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Magic methods should be inherited
	if !child.HasMagicMethod("__get") {
		t.Error("Child should inherit __get from parent")
	}
}

func TestMagic_OverrideMagicMethod(t *testing.T) {
	// Parent class with __toString
	parent := NewClassEntry("Parent")
	parent.MagicMethods["__toString"] = &MethodDef{
		Name:       "__toString",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
	}

	// Child overrides __toString
	child := NewClassEntry("Child")
	child.ParentClass = parent
	child.MagicMethods["__toString"] = &MethodDef{
		Name:       "__toString",
		Visibility: VisibilityPublic,
		NumParams:  0,
		IsMagic:    true,
	}

	if err := child.InheritFrom(parent); err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	// Child's version should take precedence
	method := child.GetMagicMethod("__toString")
	if method == nil {
		t.Fatal("Child should have __toString")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestMagic_StaticCallOnInstanceMethod(t *testing.T) {
	// __call is for instance methods, __callStatic is for static
	// They should be separate
	class := NewClassEntry("Test")

	class.MagicMethods["__call"] = &MethodDef{
		Name:       "__call",
		Visibility: VisibilityPublic,
		IsStatic:   false, // Must be instance method
		NumParams:  2,
		IsMagic:    true,
	}

	err := class.ValidateMagicMethods()
	if err != nil {
		t.Errorf("__call as instance method should be valid: %v", err)
	}

	// __call cannot be static
	class.MagicMethods["__call"].IsStatic = true
	err = class.ValidateMagicMethods()
	if err == nil {
		t.Fatal("__call cannot be static")
	}
}

func TestMagic_CallStaticOnNonStatic(t *testing.T) {
	class := NewClassEntry("Test")

	class.MagicMethods["__callStatic"] = &MethodDef{
		Name:       "__callStatic",
		Visibility: VisibilityPublic,
		IsStatic:   false, // Must be static
		NumParams:  2,
		IsMagic:    true,
	}

	err := class.ValidateMagicMethods()
	if err == nil {
		t.Fatal("__callStatic must be static")
	}
}
