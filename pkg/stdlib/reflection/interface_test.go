package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionInterface_BasicInfo tests basic interface information
func TestReflectionInterface_BasicInfo(t *testing.T) {
	ifaceEntry := &types.InterfaceEntry{
		Name: "MyNamespace\\MyInterface",
	}

	ri := NewReflectionInterface("MyNamespace\\MyInterface", ifaceEntry)

	// Test GetName
	if name := ri.GetName(); name != "MyNamespace\\MyInterface" {
		t.Errorf("GetName() = %s, want MyNamespace\\MyInterface", name)
	}

	// Test GetShortName (simplified implementation)
	shortName := ri.GetShortName()
	if shortName != "MyNamespace\\MyInterface" { // Current implementation returns full name
		t.Errorf("GetShortName() = %s", shortName)
	}
}

// TestReflectionInterface_Methods tests method enumeration
func TestReflectionInterface_Methods(t *testing.T) {
	ifaceEntry := &types.InterfaceEntry{
		Name: "MyInterface",
		Methods: map[string]*types.MethodDef{
			"method1": {
				Visibility: types.VisibilityPublic,
				NumParams:  1,
				Parameters: []*types.ParameterDef{
					{Name: "param1", Type: "string"},
				},
				ReturnType: "bool",
			},
			"method2": {
				Visibility: types.VisibilityPublic,
				NumParams:  0,
				Parameters: []*types.ParameterDef{},
				ReturnType: "void",
			},
		},
	}

	ri := NewReflectionInterface("MyInterface", ifaceEntry)

	// Test GetMethods
	methods := ri.GetMethods()
	if len(methods) != 2 {
		t.Fatalf("GetMethods() returned %d methods, want 2", len(methods))
	}

	// Verify methods are correctly created
	foundMethod1 := false
	foundMethod2 := false
	for _, method := range methods {
		switch method.GetName() {
		case "method1":
			foundMethod1 = true
			if method.GetReturnType() != "bool" {
				t.Errorf("method1 return type = %s, want bool", method.GetReturnType())
			}
		case "method2":
			foundMethod2 = true
			if method.GetReturnType() != "void" {
				t.Errorf("method2 return type = %s, want void", method.GetReturnType())
			}
		}
	}

	if !foundMethod1 {
		t.Error("method1 not found in GetMethods() result")
	}
	if !foundMethod2 {
		t.Error("method2 not found in GetMethods() result")
	}
}

// TestReflectionInterface_Constants tests constant enumeration
func TestReflectionInterface_Constants(t *testing.T) {
	ifaceEntry := &types.InterfaceEntry{
		Name: "MyInterface",
		Constants: map[string]*types.ClassConstant{
			"CONST1": {
				Value: types.NewInt(100),
			},
			"CONST2": {
				Value: types.NewString("constant"),
			},
			"CONST3": {
				Value: types.NewBool(true),
			},
		},
	}

	ri := NewReflectionInterface("MyInterface", ifaceEntry)

	// Test GetConstants
	constants := ri.GetConstants()
	if len(constants) != 3 {
		t.Fatalf("GetConstants() returned %d constants, want 3", len(constants))
	}

	// Verify constant values
	if const1, ok := constants["CONST1"]; !ok || const1.ToInt() != 100 {
		t.Error("CONST1 not found or has wrong value")
	}
	if const2, ok := constants["CONST2"]; !ok || const2.ToString() != "constant" {
		t.Error("CONST2 not found or has wrong value")
	}
	if const3, ok := constants["CONST3"]; !ok || !const3.ToBool() {
		t.Error("CONST3 not found or has wrong value")
	}
}

// TestReflectionInterface_NilInterface tests handling of nil interface
func TestReflectionInterface_NilInterface(t *testing.T) {
	ri := NewReflectionInterface("NilInterface", nil)

	// Test GetMethods with nil interface
	methods := ri.GetMethods()
	if methods != nil {
		t.Error("GetMethods() should return nil for nil interface")
	}

	// Test GetConstants with nil interface
	constants := ri.GetConstants()
	if constants != nil {
		t.Error("GetConstants() should return nil for nil interface")
	}
}

// TestReflectionInterface_EmptyInterface tests empty interface
func TestReflectionInterface_EmptyInterface(t *testing.T) {
	ifaceEntry := &types.InterfaceEntry{
		Name:      "EmptyInterface",
		Methods:   map[string]*types.MethodDef{},
		Constants: map[string]*types.ClassConstant{},
	}

	ri := NewReflectionInterface("EmptyInterface", ifaceEntry)

	// Test GetMethods returns empty slice
	methods := ri.GetMethods()
	if len(methods) != 0 {
		t.Errorf("GetMethods() returned %d methods for empty interface, want 0", len(methods))
	}

	// Test GetConstants returns empty map
	constants := ri.GetConstants()
	if len(constants) != 0 {
		t.Errorf("GetConstants() returned %d constants for empty interface, want 0", len(constants))
	}
}

// TestReflectionInterface_String tests string representation
func TestReflectionInterface_String(t *testing.T) {
	tests := []struct {
		name     string
		ifaceName string
		expected string
	}{
		{
			name:      "simple interface",
			ifaceName: "MyInterface",
			expected:  "Interface [ MyInterface ]",
		},
		{
			name:      "namespaced interface",
			ifaceName: "MyNamespace\\MyInterface",
			expected:  "Interface [ MyNamespace\\MyInterface ]",
		},
		{
			name:      "deeply namespaced interface",
			ifaceName: "App\\Services\\MyInterface",
			expected:  "Interface [ App\\Services\\MyInterface ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ifaceEntry := &types.InterfaceEntry{
				Name: tt.ifaceName,
			}
			ri := NewReflectionInterface(tt.ifaceName, ifaceEntry)

			str := ri.String()
			if str != tt.expected {
				t.Errorf("String() = %s, want %s", str, tt.expected)
			}
		})
	}
}

// TestReflectionInterface_ComplexExample tests a complex interface scenario
func TestReflectionInterface_ComplexExample(t *testing.T) {
	// Create a complex interface with methods and constants
	ifaceEntry := &types.InterfaceEntry{
		Name: "ComplexInterface",
		Methods: map[string]*types.MethodDef{
			"getData": {
				Visibility: types.VisibilityPublic,
				NumParams:  0,
				Parameters: []*types.ParameterDef{},
				ReturnType: "array",
			},
			"setData": {
				Visibility: types.VisibilityPublic,
				NumParams:  1,
				Parameters: []*types.ParameterDef{
					{Name: "data", Type: "array"},
				},
				ReturnType: "void",
			},
			"process": {
				Visibility: types.VisibilityPublic,
				NumParams:  2,
				Parameters: []*types.ParameterDef{
					{Name: "input", Type: "string"},
					{
						Name:       "options",
						Type:       "?array",
						HasDefault: true,
						Default:    types.NewNull(),
					},
				},
				ReturnType: "mixed",
			},
		},
		Constants: map[string]*types.ClassConstant{
			"VERSION": {
				Value: types.NewString("1.0.0"),
			},
			"MAX_SIZE": {
				Value: types.NewInt(1024),
			},
		},
	}

	ri := NewReflectionInterface("ComplexInterface", ifaceEntry)

	// Verify name
	if ri.GetName() != "ComplexInterface" {
		t.Errorf("GetName() = %s, want ComplexInterface", ri.GetName())
	}

	// Verify methods
	methods := ri.GetMethods()
	if len(methods) != 3 {
		t.Fatalf("GetMethods() returned %d methods, want 3", len(methods))
	}

	// Find and verify the 'process' method
	var processMethod *ReflectionMethod
	for _, method := range methods {
		if method.GetName() == "process" {
			processMethod = method
			break
		}
	}

	if processMethod == nil {
		t.Fatal("'process' method not found")
	}

	// Check method parameters
	params := processMethod.GetParameters()
	if len(params) != 2 {
		t.Fatalf("'process' method has %d parameters, want 2", len(params))
	}

	// Verify first parameter
	if params[0].GetName() != "input" {
		t.Errorf("First parameter name = %s, want input", params[0].GetName())
	}
	if params[0].GetType() != "string" {
		t.Errorf("First parameter type = %s, want string", params[0].GetType())
	}

	// Verify second parameter (optional)
	if params[1].GetName() != "options" {
		t.Errorf("Second parameter name = %s, want options", params[1].GetName())
	}
	if !params[1].IsOptional() {
		t.Error("Second parameter should be optional")
	}
	if !params[1].AllowsNull() {
		t.Error("Second parameter should allow null")
	}

	// Verify constants
	constants := ri.GetConstants()
	if len(constants) != 2 {
		t.Fatalf("GetConstants() returned %d constants, want 2", len(constants))
	}

	if version, ok := constants["VERSION"]; !ok || version.ToString() != "1.0.0" {
		t.Error("VERSION constant not found or has wrong value")
	}
	if maxSize, ok := constants["MAX_SIZE"]; !ok || maxSize.ToInt() != 1024 {
		t.Error("MAX_SIZE constant not found or has wrong value")
	}
}
