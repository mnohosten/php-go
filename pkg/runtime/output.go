package runtime

import "bytes"

// OutputBuffer represents an output buffer
type OutputBuffer struct {
	buffer bytes.Buffer
}

// NewOutputBuffer creates a new output buffer
func NewOutputBuffer() *OutputBuffer {
	return &OutputBuffer{}
}

// Write writes data to the buffer
func (ob *OutputBuffer) Write(data string) {
	ob.buffer.WriteString(data)
}

// GetContents returns the contents of the buffer
func (ob *OutputBuffer) GetContents() string {
	return ob.buffer.String()
}

// Clear clears the buffer
func (ob *OutputBuffer) Clear() {
	ob.buffer.Reset()
}

// Len returns the length of the buffer
func (ob *OutputBuffer) Len() int {
	return ob.buffer.Len()
}
