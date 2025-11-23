package reflection

import (
	"testing"

	"github.com/krizos/php-go/pkg/types"
)

// TestReflectionFunction_BasicInfo tests basic function information
func TestReflectionFunction_BasicInfo(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name:      "MyNamespace\\myFunction",
		FileName:  "/path/to/file.php",
		StartLine: 10,
		EndLine:   20,
		NumParams: 2,
	}

	rf := NewReflectionFunction("MyNamespace\\myFunction", funcDef)

	// Test GetName
	if name := rf.GetName(); name != "MyNamespace\\myFunction" {
		t.Errorf("GetName() = %s, want MyNamespace\\myFunction", name)
	}

	// Test GetShortName
	if shortName := rf.GetShortName(); shortName != "myFunction" {
		t.Errorf("GetShortName() = %s, want myFunction", shortName)
	}

	// Test GetNamespaceName
	if namespace := rf.GetNamespaceName(); namespace != "MyNamespace" {
		t.Errorf("GetNamespaceName() = %s, want MyNamespace", namespace)
	}

	// Test GetFileName
	if fileName := rf.GetFileName(); fileName != "/path/to/file.php" {
		t.Errorf("GetFileName() = %s, want /path/to/file.php", fileName)
	}

	// Test GetStartLine
	if startLine := rf.GetStartLine(); startLine != 10 {
		t.Errorf("GetStartLine() = %d, want 10", startLine)
	}

	// Test GetEndLine
	if endLine := rf.GetEndLine(); endLine != 20 {
		t.Errorf("GetEndLine() = %d, want 20", endLine)
	}
}

// TestReflectionFunction_NoNamespace tests function without namespace
func TestReflectionFunction_NoNamespace(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name: "simpleFunction",
	}

	rf := NewReflectionFunction("simpleFunction", funcDef)

	// Test GetShortName (should be same as full name)
	if shortName := rf.GetShortName(); shortName != "simpleFunction" {
		t.Errorf("GetShortName() = %s, want simpleFunction", shortName)
	}

	// Test GetNamespaceName (should be empty)
	if namespace := rf.GetNamespaceName(); namespace != "" {
		t.Errorf("GetNamespaceName() = %s, want empty string", namespace)
	}
}

// TestReflectionFunction_Characteristics tests function characteristics
func TestReflectionFunction_Characteristics(t *testing.T) {
	tests := []struct {
		name         string
		funcDef      *types.FunctionDef
		isInternal   bool
		isUserDefined bool
		isGenerator  bool
		isDeprecated bool
	}{
		{
			name: "user-defined function",
			funcDef: &types.FunctionDef{
				Name:       "userFunc",
				IsInternal: false,
			},
			isInternal:    false,
			isUserDefined: true,
			isGenerator:   false,
			isDeprecated:  false,
		},
		{
			name: "internal function",
			funcDef: &types.FunctionDef{
				Name:       "strlen",
				IsInternal: true,
			},
			isInternal:    true,
			isUserDefined: false,
			isGenerator:   false,
			isDeprecated:  false,
		},
		{
			name: "generator function",
			funcDef: &types.FunctionDef{
				Name:        "generatorFunc",
				IsGenerator: true,
			},
			isInternal:    false,
			isUserDefined: true,
			isGenerator:   true,
			isDeprecated:  false,
		},
		{
			name: "deprecated function",
			funcDef: &types.FunctionDef{
				Name:         "oldFunc",
				IsDeprecated: true,
			},
			isInternal:    false,
			isUserDefined: true,
			isGenerator:   false,
			isDeprecated:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := NewReflectionFunction(tt.funcDef.Name, tt.funcDef)

			if isInternal := rf.IsInternal(); isInternal != tt.isInternal {
				t.Errorf("IsInternal() = %v, want %v", isInternal, tt.isInternal)
			}

			if isUserDefined := rf.IsUserDefined(); isUserDefined != tt.isUserDefined {
				t.Errorf("IsUserDefined() = %v, want %v", isUserDefined, tt.isUserDefined)
			}

			if isGenerator := rf.IsGenerator(); isGenerator != tt.isGenerator {
				t.Errorf("IsGenerator() = %v, want %v", isGenerator, tt.isGenerator)
			}

			if isDeprecated := rf.IsDeprecated(); isDeprecated != tt.isDeprecated {
				t.Errorf("IsDeprecated() = %v, want %v", isDeprecated, tt.isDeprecated)
			}
		})
	}
}

// TestReflectionFunction_Variadic tests variadic function detection
func TestReflectionFunction_Variadic(t *testing.T) {
	// Non-variadic function
	funcDef1 := &types.FunctionDef{
		Name:      "regularFunc",
		NumParams: 2,
		Parameters: []*types.ParameterDef{
			{Name: "param1", Type: "string"},
			{Name: "param2", Type: "int"},
		},
	}
	rf1 := NewReflectionFunction("regularFunc", funcDef1)

	if rf1.IsVariadic() {
		t.Error("IsVariadic() = true for non-variadic function, want false")
	}

	// Variadic function
	funcDef2 := &types.FunctionDef{
		Name:      "variadicFunc",
		NumParams: 2,
		Parameters: []*types.ParameterDef{
			{Name: "param1", Type: "string"},
			{Name: "params", Type: "mixed", IsVariadic: true},
		},
	}
	rf2 := NewReflectionFunction("variadicFunc", funcDef2)

	if !rf2.IsVariadic() {
		t.Error("IsVariadic() = false for variadic function, want true")
	}
}

// TestReflectionFunction_Parameters tests parameter handling
func TestReflectionFunction_Parameters(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name:      "testFunc",
		NumParams: 3,
		Parameters: []*types.ParameterDef{
			{
				Name: "requiredParam",
				Type: "string",
			},
			{
				Name:       "optionalParam",
				Type:       "int",
				HasDefault: true,
				Default:    types.NewInt(42),
			},
			{
				Name:       "nullableParam",
				Type:       "?array",
				HasDefault: true,
				Default:    types.NewNull(),
			},
		},
	}

	rf := NewReflectionFunction("testFunc", funcDef)

	// Test GetNumberOfParameters
	if numParams := rf.GetNumberOfParameters(); numParams != 3 {
		t.Errorf("GetNumberOfParameters() = %d, want 3", numParams)
	}

	// Test GetNumberOfRequiredParameters
	if requiredParams := rf.GetNumberOfRequiredParameters(); requiredParams != 1 {
		t.Errorf("GetNumberOfRequiredParameters() = %d, want 1", requiredParams)
	}

	// Test GetParameters
	params := rf.GetParameters()
	if len(params) != 3 {
		t.Fatalf("GetParameters() returned %d parameters, want 3", len(params))
	}

	// Verify first parameter
	if params[0].GetName() != "requiredParam" {
		t.Errorf("params[0].GetName() = %s, want requiredParam", params[0].GetName())
	}
	if params[0].IsOptional() {
		t.Error("params[0].IsOptional() = true, want false")
	}

	// Verify second parameter (optional)
	if params[1].GetName() != "optionalParam" {
		t.Errorf("params[1].GetName() = %s, want optionalParam", params[1].GetName())
	}
	if !params[1].IsOptional() {
		t.Error("params[1].IsOptional() = false, want true")
	}

	// Verify third parameter (nullable)
	if params[2].GetName() != "nullableParam" {
		t.Errorf("params[2].GetName() = %s, want nullableParam", params[2].GetName())
	}
	if !params[2].AllowsNull() {
		t.Error("params[2].AllowsNull() = false, want true")
	}
}

// TestReflectionFunction_ReturnType tests return type handling
func TestReflectionFunction_ReturnType(t *testing.T) {
	// Function with return type
	funcDef1 := &types.FunctionDef{
		Name:       "funcWithReturn",
		ReturnType: "bool",
	}
	rf1 := NewReflectionFunction("funcWithReturn", funcDef1)

	if !rf1.HasReturnType() {
		t.Error("HasReturnType() = false for function with return type, want true")
	}
	if returnType := rf1.GetReturnType(); returnType != "bool" {
		t.Errorf("GetReturnType() = %s, want bool", returnType)
	}

	// Function without return type
	funcDef2 := &types.FunctionDef{
		Name:       "funcNoReturn",
		ReturnType: "",
	}
	rf2 := NewReflectionFunction("funcNoReturn", funcDef2)

	if rf2.HasReturnType() {
		t.Error("HasReturnType() = true for function without return type, want false")
	}
}

// TestReflectionFunction_ReturnsReference tests return by reference check
func TestReflectionFunction_ReturnsReference(t *testing.T) {
	// Function returning by value
	funcDef1 := &types.FunctionDef{
		Name:        "funcByValue",
		ReturnByRef: false,
	}
	rf1 := NewReflectionFunction("funcByValue", funcDef1)

	if rf1.ReturnsReference() {
		t.Error("ReturnsReference() = true for function returning by value, want false")
	}

	// Function returning by reference
	funcDef2 := &types.FunctionDef{
		Name:        "funcByRef",
		ReturnByRef: true,
	}
	rf2 := NewReflectionFunction("funcByRef", funcDef2)

	if !rf2.ReturnsReference() {
		t.Error("ReturnsReference() = false for function returning by reference, want true")
	}
}

// TestReflectionFunction_DocComment tests documentation comment retrieval
func TestReflectionFunction_DocComment(t *testing.T) {
	docComment := "/** This is a test function */"

	funcDef := &types.FunctionDef{
		Name:       "documentedFunc",
		DocComment: docComment,
	}

	rf := NewReflectionFunction("documentedFunc", funcDef)

	if doc := rf.GetDocComment(); doc != docComment {
		t.Errorf("GetDocComment() = %s, want %s", doc, docComment)
	}
}

// TestReflectionFunction_Invoke tests function invocation placeholder
func TestReflectionFunction_Invoke(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name: "testFunc",
	}

	rf := NewReflectionFunction("testFunc", funcDef)

	// Test Invoke (should return error - not implemented)
	_, err := rf.Invoke()
	if err == nil {
		t.Error("Invoke() should return error (not yet implemented)")
	}

	// Test InvokeArgs (should return error - not implemented)
	_, err = rf.InvokeArgs([]*types.Value{})
	if err == nil {
		t.Error("InvokeArgs() should return error (not yet implemented)")
	}

	// Test GetClosure (should return error - not implemented)
	_, err = rf.GetClosure()
	if err == nil {
		t.Error("GetClosure() should return error (not yet implemented)")
	}
}

// TestReflectionFunction_String tests string representation
func TestReflectionFunction_String(t *testing.T) {
	tests := []struct {
		name     string
		funcDef  *types.FunctionDef
		contains []string
	}{
		{
			name: "user-defined function",
			funcDef: &types.FunctionDef{
				Name:      "myFunc",
				FileName:  "/path/to/file.php",
				StartLine: 10,
				EndLine:   20,
			},
			contains: []string{"myFunc", "/path/to/file.php", "10-20"},
		},
		{
			name: "internal function",
			funcDef: &types.FunctionDef{
				Name:       "strlen",
				IsInternal: true,
			},
			contains: []string{"strlen", "internal"},
		},
		{
			name: "function with return type",
			funcDef: &types.FunctionDef{
				Name:       "funcWithReturn",
				ReturnType: "string",
			},
			contains: []string{"funcWithReturn", "string"},
		},
		{
			name: "function returning by reference",
			funcDef: &types.FunctionDef{
				Name:        "funcByRef",
				ReturnByRef: true,
			},
			contains: []string{"funcByRef", "&"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := NewReflectionFunction(tt.funcDef.Name, tt.funcDef)
			str := rf.String()

			for _, substr := range tt.contains {
				if !hasSubstring(str, substr) {
					t.Errorf("String() = %s, should contain %s", str, substr)
				}
			}
		})
	}
}

// TestReflectionFunction_IsClosure tests closure detection
func TestReflectionFunction_IsClosure(t *testing.T) {
	// Regular function
	funcDef1 := &types.FunctionDef{
		Name: "regularFunc",
	}
	rf1 := NewReflectionFunction("regularFunc", funcDef1)

	if rf1.IsClosure() {
		t.Error("IsClosure() = true for regular function, want false")
	}

	// Closure (name starts with {)
	funcDef2 := &types.FunctionDef{
		Name: "{closure}",
	}
	rf2 := NewReflectionFunction("{closure}", funcDef2)

	if !rf2.IsClosure() {
		t.Error("IsClosure() = false for closure, want true")
	}

	// Anonymous function (empty name)
	funcDef3 := &types.FunctionDef{
		Name: "",
	}
	rf3 := NewReflectionFunction("", funcDef3)

	if !rf3.IsClosure() {
		t.Error("IsClosure() = false for anonymous function, want true")
	}
}

// TestReflectionFunction_Extension tests extension information
func TestReflectionFunction_Extension(t *testing.T) {
	// User-defined function (no extension)
	funcDef1 := &types.FunctionDef{
		Name:       "userFunc",
		IsInternal: false,
	}
	rf1 := NewReflectionFunction("userFunc", funcDef1)

	if ext := rf1.GetExtension(); ext != "" {
		t.Errorf("GetExtension() = %s for user function, want empty string", ext)
	}
	if extName := rf1.GetExtensionName(); extName != "" {
		t.Errorf("GetExtensionName() = %s for user function, want empty string", extName)
	}

	// Internal function
	funcDef2 := &types.FunctionDef{
		Name:       "strlen",
		IsInternal: true,
	}
	rf2 := NewReflectionFunction("strlen", funcDef2)

	if ext := rf2.GetExtension(); ext != "Core" {
		t.Errorf("GetExtension() = %s for internal function, want Core", ext)
	}
	if extName := rf2.GetExtensionName(); extName != "Core" {
		t.Errorf("GetExtensionName() = %s for internal function, want Core", extName)
	}
}

// TestReflectionFunction_IsDisabled tests disabled function check
func TestReflectionFunction_IsDisabled(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name: "testFunc",
	}

	rf := NewReflectionFunction("testFunc", funcDef)

	// Currently always returns false (placeholder)
	if rf.IsDisabled() {
		t.Error("IsDisabled() = true, want false")
	}
}

// TestReflectionFunction_ComplexExample tests a complex function scenario
func TestReflectionFunction_ComplexExample(t *testing.T) {
	funcDef := &types.FunctionDef{
		Name:        "MyNamespace\\complexFunction",
		FileName:    "/path/to/file.php",
		StartLine:   100,
		EndLine:     150,
		NumParams:   4,
		ReturnType:  "?array",
		ReturnByRef: true,
		IsGenerator: true,
		DocComment:  "/** Complex generator function */",
		Parameters: []*types.ParameterDef{
			{
				Name: "required",
				Type: "string",
			},
			{
				Name:       "optional",
				Type:       "int",
				HasDefault: true,
				Default:    types.NewInt(10),
			},
			{
				Name:        "byRef",
				Type:        "array",
				PassedByRef: true,
			},
			{
				Name:       "variadic",
				Type:       "mixed",
				IsVariadic: true,
			},
		},
	}

	rf := NewReflectionFunction("MyNamespace\\complexFunction", funcDef)

	// Verify basic info
	if rf.GetName() != "MyNamespace\\complexFunction" {
		t.Errorf("GetName() = %s, want MyNamespace\\complexFunction", rf.GetName())
	}
	if rf.GetShortName() != "complexFunction" {
		t.Errorf("GetShortName() = %s, want complexFunction", rf.GetShortName())
	}
	if rf.GetNamespaceName() != "MyNamespace" {
		t.Errorf("GetNamespaceName() = %s, want MyNamespace", rf.GetNamespaceName())
	}

	// Verify characteristics
	if !rf.IsGenerator() {
		t.Error("IsGenerator() = false, want true")
	}
	if !rf.IsVariadic() {
		t.Error("IsVariadic() = false, want true")
	}
	if !rf.ReturnsReference() {
		t.Error("ReturnsReference() = false, want true")
	}

	// Verify parameters
	if rf.GetNumberOfParameters() != 4 {
		t.Errorf("GetNumberOfParameters() = %d, want 4", rf.GetNumberOfParameters())
	}
	// Required parameters: required (no default), byRef (no default), variadic (no default) = 3
	// Optional parameters: optional (has default) = 1
	if rf.GetNumberOfRequiredParameters() != 3 {
		t.Errorf("GetNumberOfRequiredParameters() = %d, want 3", rf.GetNumberOfRequiredParameters())
	}

	params := rf.GetParameters()
	if len(params) != 4 {
		t.Fatalf("GetParameters() returned %d parameters, want 4", len(params))
	}

	// Check by-reference parameter
	if !params[2].IsPassedByReference() {
		t.Error("Third parameter should be passed by reference")
	}

	// Check variadic parameter
	if !params[3].IsVariadic() {
		t.Error("Fourth parameter should be variadic")
	}

	// Verify return type
	if rf.GetReturnType() != "?array" {
		t.Errorf("GetReturnType() = %s, want ?array", rf.GetReturnType())
	}

	// Verify documentation
	if rf.GetDocComment() != "/** Complex generator function */" {
		t.Errorf("GetDocComment() incorrect")
	}

	// Verify file info
	if rf.GetFileName() != "/path/to/file.php" {
		t.Errorf("GetFileName() = %s, want /path/to/file.php", rf.GetFileName())
	}
	if rf.GetStartLine() != 100 {
		t.Errorf("GetStartLine() = %d, want 100", rf.GetStartLine())
	}
	if rf.GetEndLine() != 150 {
		t.Errorf("GetEndLine() = %d, want 150", rf.GetEndLine())
	}
}

// TestReflectionFunction_NilFunctionDef tests handling of nil function definition
func TestReflectionFunction_NilFunctionDef(t *testing.T) {
	rf := NewReflectionFunction("nilFunc", nil)

	// Should not panic and should return default values
	if rf.GetFileName() != "" {
		t.Error("GetFileName() should return empty string for nil function")
	}
	if rf.GetStartLine() != 0 {
		t.Error("GetStartLine() should return 0 for nil function")
	}
	if rf.GetNumberOfParameters() != 0 {
		t.Error("GetNumberOfParameters() should return 0 for nil function")
	}
	if rf.GetParameters() != nil {
		t.Error("GetParameters() should return nil for nil function")
	}
}
