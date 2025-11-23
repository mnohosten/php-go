package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionMethod_BasicInfo tests basic method information
func TestReflectionMethod_BasicInfo(t *testing.T) {
	classEntry := &types.ClassEntry{
		Name: "TestClass",
	}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		NumParams:  2,
	}

	rm := NewReflectionMethod(rc, "testMethod", methodDef)

	// Test GetName
	if name := rm.GetName(); name != "testMethod" {
		t.Errorf("GetName() = %s, want testMethod", name)
	}

	// Test GetDeclaringClass
	if declaringClass := rm.GetDeclaringClass(); declaringClass != rc {
		t.Error("GetDeclaringClass() returned wrong class")
	}

	// Test GetNumberOfParameters
	if numParams := rm.GetNumberOfParameters(); numParams != 2 {
		t.Errorf("GetNumberOfParameters() = %d, want 2", numParams)
	}
}

// TestReflectionMethod_Visibility tests visibility checks
func TestReflectionMethod_Visibility(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	tests := []struct {
		name        string
		visibility  types.PropertyVisibility
		isPublic    bool
		isProtected bool
		isPrivate   bool
	}{
		{
			name:        "public method",
			visibility:  types.VisibilityPublic,
			isPublic:    true,
			isProtected: false,
			isPrivate:   false,
		},
		{
			name:        "protected method",
			visibility:  types.VisibilityProtected,
			isPublic:    false,
			isProtected: true,
			isPrivate:   false,
		},
		{
			name:        "private method",
			visibility:  types.VisibilityPrivate,
			isPublic:    false,
			isProtected: false,
			isPrivate:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			methodDef := &types.MethodDef{
				Visibility: tt.visibility,
			}
			rm := NewReflectionMethod(rc, "method", methodDef)

			if isPublic := rm.IsPublic(); isPublic != tt.isPublic {
				t.Errorf("IsPublic() = %v, want %v", isPublic, tt.isPublic)
			}
			if isProtected := rm.IsProtected(); isProtected != tt.isProtected {
				t.Errorf("IsProtected() = %v, want %v", isProtected, tt.isProtected)
			}
			if isPrivate := rm.IsPrivate(); isPrivate != tt.isPrivate {
				t.Errorf("IsPrivate() = %v, want %v", isPrivate, tt.isPrivate)
			}
		})
	}
}

// TestReflectionMethod_Modifiers tests method modifiers
func TestReflectionMethod_Modifiers(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	tests := []struct {
		name          string
		methodDef     *types.MethodDef
		isStatic      bool
		isFinal       bool
		isAbstract    bool
		isConstructor bool
		isDestructor  bool
	}{
		{
			name: "regular method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
			},
			isStatic:      false,
			isFinal:       false,
			isAbstract:    false,
			isConstructor: false,
			isDestructor:  false,
		},
		{
			name: "static method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsStatic:   true,
			},
			isStatic: true,
		},
		{
			name: "final method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsFinal:    true,
			},
			isFinal: true,
		},
		{
			name: "abstract method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsAbstract: true,
			},
			isAbstract: true,
		},
		{
			name: "constructor",
			methodDef: &types.MethodDef{
				Visibility:    types.VisibilityPublic,
				IsConstructor: true,
			},
			isConstructor: true,
		},
		{
			name: "destructor",
			methodDef: &types.MethodDef{
				Visibility:   types.VisibilityPublic,
				IsDestructor: true,
			},
			isDestructor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewReflectionMethod(rc, "method", tt.methodDef)

			if isStatic := rm.IsStatic(); isStatic != tt.isStatic {
				t.Errorf("IsStatic() = %v, want %v", isStatic, tt.isStatic)
			}
			if isFinal := rm.IsFinal(); isFinal != tt.isFinal {
				t.Errorf("IsFinal() = %v, want %v", isFinal, tt.isFinal)
			}
			if isAbstract := rm.IsAbstract(); isAbstract != tt.isAbstract {
				t.Errorf("IsAbstract() = %v, want %v", isAbstract, tt.isAbstract)
			}
			if isConstructor := rm.IsConstructor(); isConstructor != tt.isConstructor {
				t.Errorf("IsConstructor() = %v, want %v", isConstructor, tt.isConstructor)
			}
			if isDestructor := rm.IsDestructor(); isDestructor != tt.isDestructor {
				t.Errorf("IsDestructor() = %v, want %v", isDestructor, tt.isDestructor)
			}
		})
	}
}

// TestReflectionMethod_Parameters tests parameter handling
func TestReflectionMethod_Parameters(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		NumParams:  3,
		Parameters: []*types.ParameterDef{
			{
				Name: "param1",
				Type: "string",
			},
			{
				Name:       "param2",
				Type:       "int",
				HasDefault: true,
				Default:    types.NewInt(42),
			},
			{
				Name:       "param3",
				Type:       "array",
				HasDefault: true,
				Default:    types.NewArray(types.NewArrayWithCapacity(0)),
			},
		},
	}

	rm := NewReflectionMethod(rc, "method", methodDef)

	// Test GetNumberOfParameters
	if numParams := rm.GetNumberOfParameters(); numParams != 3 {
		t.Errorf("GetNumberOfParameters() = %d, want 3", numParams)
	}

	// Test GetNumberOfRequiredParameters
	if requiredParams := rm.GetNumberOfRequiredParameters(); requiredParams != 1 {
		t.Errorf("GetNumberOfRequiredParameters() = %d, want 1", requiredParams)
	}

	// Test GetParameters
	params := rm.GetParameters()
	if len(params) != 3 {
		t.Fatalf("GetParameters() returned %d parameters, want 3", len(params))
	}

	// Verify first parameter
	if params[0].GetName() != "param1" {
		t.Errorf("params[0].GetName() = %s, want param1", params[0].GetName())
	}
	if params[0].IsOptional() {
		t.Error("params[0].IsOptional() = true, want false")
	}

	// Verify second parameter (optional)
	if params[1].GetName() != "param2" {
		t.Errorf("params[1].GetName() = %s, want param2", params[1].GetName())
	}
	if !params[1].IsOptional() {
		t.Error("params[1].IsOptional() = false, want true")
	}
}

// TestReflectionMethod_ReturnType tests return type handling
func TestReflectionMethod_ReturnType(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Method with return type
	methodDef1 := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		ReturnType: "bool",
	}
	rm1 := NewReflectionMethod(rc, "methodWithReturn", methodDef1)

	if !rm1.HasReturnType() {
		t.Error("HasReturnType() = false for method with return type, want true")
	}
	if returnType := rm1.GetReturnType(); returnType != "bool" {
		t.Errorf("GetReturnType() = %s, want bool", returnType)
	}

	// Method without return type
	methodDef2 := &types.MethodDef{
		Visibility: types.VisibilityPublic,
		ReturnType: "",
	}
	rm2 := NewReflectionMethod(rc, "methodNoReturn", methodDef2)

	if rm2.HasReturnType() {
		t.Error("HasReturnType() = true for method without return type, want false")
	}
}

// TestReflectionMethod_ReturnsReference tests return by reference check
func TestReflectionMethod_ReturnsReference(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	// Method returning by value
	methodDef1 := &types.MethodDef{
		Visibility:    types.VisibilityPublic,
		ReturnByRef:   false,
	}
	rm1 := NewReflectionMethod(rc, "methodByValue", methodDef1)

	if rm1.ReturnsReference() {
		t.Error("ReturnsReference() = true for method returning by value, want false")
	}

	// Method returning by reference
	methodDef2 := &types.MethodDef{
		Visibility:    types.VisibilityPublic,
		ReturnByRef:   true,
	}
	rm2 := NewReflectionMethod(rc, "methodByRef", methodDef2)

	if !rm2.ReturnsReference() {
		t.Error("ReturnsReference() = false for method returning by reference, want true")
	}
}

// TestReflectionMethod_Invoke tests method invocation placeholder
func TestReflectionMethod_Invoke(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPublic,
	}
	rm := NewReflectionMethod(rc, "method", methodDef)

	obj := &types.Object{
		ClassName:  "TestClass",
		ClassEntry: classEntry,
		Properties: make(map[string]*types.Property),
	}

	// Test Invoke (should return error - not implemented)
	_, err := rm.Invoke(obj)
	if err == nil {
		t.Error("Invoke() should return error (not yet implemented)")
	}

	// Test InvokeArgs (should return error - not implemented)
	_, err = rm.InvokeArgs(obj, []*types.Value{})
	if err == nil {
		t.Error("InvokeArgs() should return error (not yet implemented)")
	}
}

// TestReflectionMethod_SetAccessible tests accessibility setter
func TestReflectionMethod_SetAccessible(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	methodDef := &types.MethodDef{
		Visibility: types.VisibilityPrivate,
	}
	rm := NewReflectionMethod(rc, "privateMethod", methodDef)

	// SetAccessible should not panic (it's a no-op for API compatibility)
	rm.SetAccessible(true)
	rm.SetAccessible(false)
}

// TestReflectionMethod_String tests string representation
func TestReflectionMethod_String(t *testing.T) {
	classEntry := &types.ClassEntry{Name: "TestClass"}
	rc := NewReflectionClass("TestClass", classEntry)

	tests := []struct {
		name      string
		methodDef *types.MethodDef
		methName  string
		contains  []string
	}{
		{
			name: "public method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
			},
			methName: "publicMethod",
			contains: []string{"public", "publicMethod"},
		},
		{
			name: "protected method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityProtected,
			},
			methName: "protectedMethod",
			contains: []string{"protected", "protectedMethod"},
		},
		{
			name: "private method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPrivate,
			},
			methName: "privateMethod",
			contains: []string{"private", "privateMethod"},
		},
		{
			name: "static method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsStatic:   true,
			},
			methName: "staticMethod",
			contains: []string{"static", "staticMethod"},
		},
		{
			name: "final method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsFinal:    true,
			},
			methName: "finalMethod",
			contains: []string{"final", "finalMethod"},
		},
		{
			name: "abstract method",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				IsAbstract: true,
			},
			methName: "abstractMethod",
			contains: []string{"abstract", "abstractMethod"},
		},
		{
			name: "method with return type",
			methodDef: &types.MethodDef{
				Visibility: types.VisibilityPublic,
				ReturnType: "string",
			},
			methName: "methodWithReturn",
			contains: []string{"methodWithReturn", "string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewReflectionMethod(rc, tt.methName, tt.methodDef)
			str := rm.String()

			for _, substr := range tt.contains {
				if !hasSubstring(str, substr) {
					t.Errorf("String() = %s, should contain %s", str, substr)
				}
			}
		})
	}
}
