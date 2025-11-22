package runtime

import (
	"fmt"
	"os"
	"time"

	"github.com/krizos/php-go/pkg/types"
)

// Runtime manages PHP runtime state including superglobals, constants, and configuration
type Runtime struct {
	// Superglobals
	GET     *types.Value // $_GET
	POST    *types.Value // $_POST
	REQUEST *types.Value // $_REQUEST
	SERVER  *types.Value // $_SERVER
	ENV     *types.Value // $_ENV
	COOKIE  *types.Value // $_COOKIE
	FILES   *types.Value // $_FILES
	SESSION *types.Value // $_SESSION
	GLOBALS *types.Value // $GLOBALS

	// Constants (name -> value)
	constants map[string]*types.Value

	// Error handling
	errorReporting int
	errorHandler   ErrorHandler

	// Output buffering
	outputBuffers []*OutputBuffer
	currentBuffer *OutputBuffer

	// Execution context
	scriptPath string
	scriptDir  string
	startTime  time.Time
}

// ErrorHandler is a function that handles errors
type ErrorHandler func(errorType ErrorType, message string, file string, line int)

// New creates a new runtime instance
func New() *Runtime {
	rt := &Runtime{
		GET:            types.NewArray(types.NewEmptyArray()),
		POST:           types.NewArray(types.NewEmptyArray()),
		REQUEST:        types.NewArray(types.NewEmptyArray()),
		SERVER:         types.NewArray(types.NewEmptyArray()),
		ENV:            types.NewArray(types.NewEmptyArray()),
		COOKIE:         types.NewArray(types.NewEmptyArray()),
		FILES:          types.NewArray(types.NewEmptyArray()),
		SESSION:        types.NewArray(types.NewEmptyArray()),
		GLOBALS:        types.NewArray(types.NewEmptyArray()),
		constants:      make(map[string]*types.Value),
		errorReporting: int(E_ALL),
		outputBuffers:  make([]*OutputBuffer, 0),
		startTime:      time.Now(),
	}

	// Initialize built-in constants
	rt.initBuiltinConstants()

	// Initialize $_SERVER
	rt.initServerSuperglobal()

	return rt
}

// ============================================================================
// Constants Management
// ============================================================================

// DefineConstant defines a constant
func (rt *Runtime) DefineConstant(name string, value *types.Value) error {
	if _, exists := rt.constants[name]; exists {
		return fmt.Errorf("constant '%s' already defined", name)
	}

	rt.constants[name] = value
	return nil
}

// GetConstant retrieves a constant
func (rt *Runtime) GetConstant(name string) (*types.Value, bool) {
	val, ok := rt.constants[name]
	return val, ok
}

// ConstantExists checks if a constant exists
func (rt *Runtime) ConstantExists(name string) bool {
	_, exists := rt.constants[name]
	return exists
}

// initBuiltinConstants initializes PHP built-in constants
func (rt *Runtime) initBuiltinConstants() {
	// PHP version constants
	rt.constants["PHP_VERSION"] = types.NewString("8.4.0-dev")
	rt.constants["PHP_MAJOR_VERSION"] = types.NewInt(8)
	rt.constants["PHP_MINOR_VERSION"] = types.NewInt(4)
	rt.constants["PHP_RELEASE_VERSION"] = types.NewInt(0)

	// Boolean constants
	rt.constants["TRUE"] = types.NewBool(true)
	rt.constants["FALSE"] = types.NewBool(false)
	rt.constants["NULL"] = types.NewNull()

	// Path constants (will be updated when script runs)
	rt.constants["PHP_EOL"] = types.NewString("\n")
	rt.constants["DIRECTORY_SEPARATOR"] = types.NewString(string(os.PathSeparator))

	// Math constants
	rt.constants["PHP_INT_MAX"] = types.NewInt(9223372036854775807)
	rt.constants["PHP_INT_MIN"] = types.NewInt(-9223372036854775808)
	rt.constants["PHP_FLOAT_MAX"] = types.NewFloat(1.7976931348623157e+308)
	rt.constants["PHP_FLOAT_MIN"] = types.NewFloat(2.2250738585072014e-308)
}

// ============================================================================
// Superglobals Management
// ============================================================================

// initServerSuperglobal initializes $_SERVER with default values
func (rt *Runtime) initServerSuperglobal() {
	server := rt.SERVER.ToArray()

	// Script information
	if rt.scriptPath != "" {
		server.Set(types.NewString("SCRIPT_FILENAME"), types.NewString(rt.scriptPath))
		server.Set(types.NewString("SCRIPT_NAME"), types.NewString(rt.scriptPath))
	}

	// Server information
	hostname, _ := os.Hostname()
	server.Set(types.NewString("SERVER_NAME"), types.NewString(hostname))
	server.Set(types.NewString("SERVER_SOFTWARE"), types.NewString("PHP-Go/1.0"))

	// Request information
	server.Set(types.NewString("REQUEST_METHOD"), types.NewString("CLI"))
	server.Set(types.NewString("REQUEST_TIME"), types.NewInt(rt.startTime.Unix()))
	server.Set(types.NewString("REQUEST_TIME_FLOAT"), types.NewFloat(float64(rt.startTime.UnixNano())/1e9))

	// Environment
	for _, env := range os.Environ() {
		// Parse KEY=VALUE
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				key := env[:i]
				value := env[i+1:]
				server.Set(types.NewString(key), types.NewString(value))
				break
			}
		}
	}
}

// SetScriptPath sets the current script path
func (rt *Runtime) SetScriptPath(path string) {
	rt.scriptPath = path
	rt.initServerSuperglobal()
}

// GetSuperglobal retrieves a superglobal by name
func (rt *Runtime) GetSuperglobal(name string) (*types.Value, bool) {
	switch name {
	case "_GET":
		return rt.GET, true
	case "_POST":
		return rt.POST, true
	case "_REQUEST":
		return rt.REQUEST, true
	case "_SERVER":
		return rt.SERVER, true
	case "_ENV":
		return rt.ENV, true
	case "_COOKIE":
		return rt.COOKIE, true
	case "_FILES":
		return rt.FILES, true
	case "_SESSION":
		return rt.SESSION, true
	case "GLOBALS":
		return rt.GLOBALS, true
	default:
		return nil, false
	}
}

// ============================================================================
// Error Handling
// ============================================================================

// SetErrorReporting sets the error reporting level
func (rt *Runtime) SetErrorReporting(level int) {
	rt.errorReporting = level
}

// GetErrorReporting gets the error reporting level
func (rt *Runtime) GetErrorReporting() int {
	return rt.errorReporting
}

// SetErrorHandler sets a custom error handler
func (rt *Runtime) SetErrorHandler(handler ErrorHandler) {
	rt.errorHandler = handler
}

// TriggerError triggers an error
func (rt *Runtime) TriggerError(errorType ErrorType, message string, file string, line int) {
	// Check if this error type should be reported
	if rt.errorReporting&int(errorType) == 0 {
		return
	}

	// Call custom error handler if set
	if rt.errorHandler != nil {
		rt.errorHandler(errorType, message, file, line)
		return
	}

	// Default error handling - print to stderr
	fmt.Fprintf(os.Stderr, "%s: %s in %s on line %d\n",
		errorType.String(), message, file, line)
}

// ============================================================================
// Output Buffering
// ============================================================================

// StartOutputBuffering starts a new output buffer
func (rt *Runtime) StartOutputBuffering() {
	buffer := NewOutputBuffer()
	rt.outputBuffers = append(rt.outputBuffers, buffer)
	rt.currentBuffer = buffer
}

// EndOutputBuffering ends the current output buffer and returns its contents
func (rt *Runtime) EndOutputBuffering() string {
	if len(rt.outputBuffers) == 0 {
		return ""
	}

	// Get current buffer
	buffer := rt.outputBuffers[len(rt.outputBuffers)-1]
	contents := buffer.GetContents()

	// Pop buffer
	rt.outputBuffers = rt.outputBuffers[:len(rt.outputBuffers)-1]

	// Update current buffer
	if len(rt.outputBuffers) > 0 {
		rt.currentBuffer = rt.outputBuffers[len(rt.outputBuffers)-1]
	} else {
		rt.currentBuffer = nil
	}

	return contents
}

// CleanOutputBuffer ends the current buffer and discards its contents
func (rt *Runtime) CleanOutputBuffer() {
	if len(rt.outputBuffers) == 0 {
		return
	}

	// Pop buffer without returning contents
	rt.outputBuffers = rt.outputBuffers[:len(rt.outputBuffers)-1]

	// Update current buffer
	if len(rt.outputBuffers) > 0 {
		rt.currentBuffer = rt.outputBuffers[len(rt.outputBuffers)-1]
	} else {
		rt.currentBuffer = nil
	}
}

// GetOutputBufferContents gets the contents of the current output buffer without ending it
func (rt *Runtime) GetOutputBufferContents() string {
	if rt.currentBuffer == nil {
		return ""
	}

	return rt.currentBuffer.GetContents()
}

// FlushOutputBuffer flushes the current output buffer
func (rt *Runtime) FlushOutputBuffer() string {
	if rt.currentBuffer == nil {
		return ""
	}

	contents := rt.currentBuffer.GetContents()
	rt.currentBuffer.Clear()
	return contents
}

// Write writes to the current output buffer (or stdout if no buffer active)
func (rt *Runtime) Write(data string) {
	if rt.currentBuffer != nil {
		rt.currentBuffer.Write(data)
	} else {
		fmt.Print(data)
	}
}

// GetOutputBufferLevel returns the nesting level of output buffers
func (rt *Runtime) GetOutputBufferLevel() int {
	return len(rt.outputBuffers)
}
