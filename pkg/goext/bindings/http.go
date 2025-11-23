package bindings

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// HTTPExtension provides Go HTTP client bindings for PHP.
type HTTPExtension struct {
	*goext.BaseExtension
	client *http.Client
}

// NewHTTPExtension creates a new HTTP extension.
func NewHTTPExtension() *HTTPExtension {
	ext := &HTTPExtension{
		BaseExtension: goext.NewBaseExtension("go_http", "1.0.0"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Register functions
	ext.AddFunction("go_http_get", ext.httpGet)
	ext.AddFunction("go_http_post", ext.httpPost)
	ext.AddFunction("go_http_put", ext.httpPut)
	ext.AddFunction("go_http_delete", ext.httpDelete)
	ext.AddFunction("go_http_request", ext.httpRequest)

	// Register constants
	ext.AddConstant("HTTP_METHOD_GET", types.NewString("GET"))
	ext.AddConstant("HTTP_METHOD_POST", types.NewString("POST"))
	ext.AddConstant("HTTP_METHOD_PUT", types.NewString("PUT"))
	ext.AddConstant("HTTP_METHOD_DELETE", types.NewString("DELETE"))
	ext.AddConstant("HTTP_METHOD_PATCH", types.NewString("PATCH"))

	return ext
}

// httpGet implements go_http_get($url, $options = [])
//
// Options:
//   - headers: array of HTTP headers
//   - timeout: timeout in seconds (default: 30)
//
// Returns array:
//   - status: HTTP status code
//   - headers: response headers
//   - body: response body
func (e *HTTPExtension) httpGet(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_http_get() expects at least 1 argument (url), got %d", len(args))
	}

	url := args[0].ToString()

	// Parse options
	var options *types.Value
	if len(args) >= 2 {
		options = args[1]
	}

	return e.doRequest("GET", url, nil, options)
}

// httpPost implements go_http_post($url, $body, $options = [])
func (e *HTTPExtension) httpPost(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("go_http_post() expects at least 2 arguments (url, body), got %d", len(args))
	}

	url := args[0].ToString()
	body := args[1].ToString()

	var options *types.Value
	if len(args) >= 3 {
		options = args[2]
	}

	return e.doRequest("POST", url, []byte(body), options)
}

// httpPut implements go_http_put($url, $body, $options = [])
func (e *HTTPExtension) httpPut(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("go_http_put() expects at least 2 arguments (url, body), got %d", len(args))
	}

	url := args[0].ToString()
	body := args[1].ToString()

	var options *types.Value
	if len(args) >= 3 {
		options = args[2]
	}

	return e.doRequest("PUT", url, []byte(body), options)
}

// httpDelete implements go_http_delete($url, $options = [])
func (e *HTTPExtension) httpDelete(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("go_http_delete() expects at least 1 argument (url), got %d", len(args))
	}

	url := args[0].ToString()

	var options *types.Value
	if len(args) >= 2 {
		options = args[1]
	}

	return e.doRequest("DELETE", url, nil, options)
}

// httpRequest implements go_http_request($method, $url, $body = null, $options = [])
func (e *HTTPExtension) httpRequest(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("go_http_request() expects at least 2 arguments (method, url), got %d", len(args))
	}

	method := args[0].ToString()
	url := args[1].ToString()

	var body []byte
	if len(args) >= 3 && args[2].Type() != types.TypeNull {
		body = []byte(args[2].ToString())
	}

	var options *types.Value
	if len(args) >= 4 {
		options = args[3]
	}

	return e.doRequest(method, url, body, options)
}

// doRequest performs the actual HTTP request
func (e *HTTPExtension) doRequest(method, url string, body []byte, options *types.Value) (*types.Value, error) {
	// Create request
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Parse options
	client := e.client
	if options != nil && options.Type() == types.TypeArray {
		arr := options.ToArray()

		// Set headers
		if headersVal, ok := arr.Get(types.NewString("headers")); ok && headersVal.Type() == types.TypeArray {
			headers := headersVal.ToArray()
			headers.Each(func(key, val *types.Value) bool {
				req.Header.Set(key.ToString(), val.ToString())
				return true
			})
		}

		// Set timeout
		if timeoutVal, ok := arr.Get(types.NewString("timeout")); ok {
			timeout := time.Duration(timeoutVal.ToInt()) * time.Second
			client = &http.Client{Timeout: timeout}
		}
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Build response array
	result := types.NewEmptyArray()
	result.Set(types.NewString("status"), types.NewInt(int64(resp.StatusCode)))
	result.Set(types.NewString("body"), types.NewString(string(respBody)))

	// Add response headers
	headers := types.NewEmptyArray()
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers.Set(types.NewString(key), types.NewString(values[0]))
		}
	}
	result.Set(types.NewString("headers"), types.NewArray(headers))

	return types.NewArray(result), nil
}
