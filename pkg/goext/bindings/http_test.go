package bindings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

func TestHTTPExtension(t *testing.T) {
	ext := NewHTTPExtension()

	if ext.Name() != "go_http" {
		t.Errorf("Expected name 'go_http', got '%s'", ext.Name())
	}

	if ext.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", ext.Version())
	}

	// Check functions registered
	funcs := ext.Functions()
	expectedFuncs := []string{
		"go_http_get",
		"go_http_post",
		"go_http_put",
		"go_http_delete",
		"go_http_request",
	}

	for _, name := range expectedFuncs {
		if _, ok := funcs[name]; !ok {
			t.Errorf("Function '%s' not registered", name)
		}
	}

	// Check constants
	constants := ext.Constants()
	if val, ok := constants["HTTP_METHOD_GET"]; !ok || val.ToString() != "GET" {
		t.Error("HTTP_METHOD_GET constant not registered correctly")
	}
}

func TestHTTPGet(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	// Test simple GET
	result, err := ext.httpGet(testVM, []*types.Value{
		types.NewString(server.URL),
	})

	if err != nil {
		t.Fatalf("httpGet failed: %v", err)
	}

	if result.Type() != types.TypeArray {
		t.Fatalf("Expected array result, got %s", result.TypeString())
	}

	arr := result.ToArray()

	// Check status
	statusVal, ok := arr.Get(types.NewString("status"))
	if !ok {
		t.Fatal("Status not in response")
	}
	if statusVal.ToInt() != 200 {
		t.Errorf("Expected status 200, got %d", statusVal.ToInt())
	}

	// Check body
	bodyVal, ok := arr.Get(types.NewString("body"))
	if !ok {
		t.Fatal("Body not in response")
	}
	if bodyVal.ToString() != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", bodyVal.ToString())
	}
}

func TestHTTPGetWithHeaders(t *testing.T) {
	// Create test server that checks headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Error("Custom header not received")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	// Build options with headers
	options := types.NewEmptyArray()
	headers := types.NewEmptyArray()
	headers.Set(types.NewString("X-Custom-Header"), types.NewString("test-value"))
	options.Set(types.NewString("headers"), types.NewArray(headers))

	result, err := ext.httpGet(testVM, []*types.Value{
		types.NewString(server.URL),
		types.NewArray(options),
	})

	if err != nil {
		t.Fatalf("httpGet with headers failed: %v", err)
	}

	arr := result.ToArray()
	statusVal, _ := arr.Get(types.NewString("status"))
	if statusVal.ToInt() != 200 {
		t.Errorf("Expected status 200, got %d", statusVal.ToInt())
	}
}

func TestHTTPPost(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	result, err := ext.httpPost(testVM, []*types.Value{
		types.NewString(server.URL),
		types.NewString("test data"),
	})

	if err != nil {
		t.Fatalf("httpPost failed: %v", err)
	}

	arr := result.ToArray()
	statusVal, _ := arr.Get(types.NewString("status"))
	if statusVal.ToInt() != 201 {
		t.Errorf("Expected status 201, got %d", statusVal.ToInt())
	}
}

func TestHTTPPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	result, err := ext.httpPut(testVM, []*types.Value{
		types.NewString(server.URL),
		types.NewString("update data"),
	})

	if err != nil {
		t.Fatalf("httpPut failed: %v", err)
	}

	arr := result.ToArray()
	statusVal, _ := arr.Get(types.NewString("status"))
	if statusVal.ToInt() != 200 {
		t.Errorf("Expected status 200, got %d", statusVal.ToInt())
	}
}

func TestHTTPDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	result, err := ext.httpDelete(testVM, []*types.Value{
		types.NewString(server.URL),
	})

	if err != nil {
		t.Fatalf("httpDelete failed: %v", err)
	}

	arr := result.ToArray()
	statusVal, _ := arr.Get(types.NewString("status"))
	if statusVal.ToInt() != 204 {
		t.Errorf("Expected status 204, got %d", statusVal.ToInt())
	}
}

func TestHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ext := NewHTTPExtension()
	testVM := vm.New()

	result, err := ext.httpRequest(testVM, []*types.Value{
		types.NewString("PATCH"),
		types.NewString(server.URL),
		types.NewString("patch data"),
	})

	if err != nil {
		t.Fatalf("httpRequest failed: %v", err)
	}

	arr := result.ToArray()
	statusVal, _ := arr.Get(types.NewString("status"))
	if statusVal.ToInt() != 200 {
		t.Errorf("Expected status 200, got %d", statusVal.ToInt())
	}
}

func TestHTTPGetInvalidArgs(t *testing.T) {
	ext := NewHTTPExtension()
	testVM := vm.New()

	_, err := ext.httpGet(testVM, []*types.Value{})
	if err == nil {
		t.Error("Expected error for missing URL argument")
	}
}

func TestHTTPPostInvalidArgs(t *testing.T) {
	ext := NewHTTPExtension()
	testVM := vm.New()

	_, err := ext.httpPost(testVM, []*types.Value{
		types.NewString("http://example.com"),
	})
	if err == nil {
		t.Error("Expected error for missing body argument")
	}
}
