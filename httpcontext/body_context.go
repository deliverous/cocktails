package httpcontext

import (
	"io"
	"net/http"
)

// BodyContext stores value into the request.
// It currently accomplishes this by replacing the http.Request’s Body with
// a ContextReadCloser, which wraps the original io.ReadCloser.
type BodyContext struct {
}

// NewBodyContext creates a new BodyContext
func NewBodyContext() *BodyContext {
	return &BodyContext{}
}

// Set stores a value for a given key in a given request.
func (context *BodyContext) Set(request *http.Request, key interface{}, value interface{}) {
	getContextReadCloser(request).Context()[key] = value
}

// Get returns a value stored for a given key in a given request.
func (context *BodyContext) Get(request *http.Request, key interface{}) interface{} {
	return getContextReadCloser(request).Context()[key]
}

// GetOk returns stored value and presence state like multi-value return of map access.
func (context *BodyContext) GetOk(request *http.Request, key interface{}) (interface{}, bool) {
	value, ok := getContextReadCloser(request).Context()[key]
	return value, ok
}

// GetAll returns all stored values for the request as a map.
func (context *BodyContext) GetAll(request *http.Request) map[interface{}]interface{} {
	return getContextReadCloser(request).Context()
}

// Delete removes a value stored for a given key in a given request.
func (context *BodyContext) Delete(request *http.Request, key interface{}) {
	delete(getContextReadCloser(request).Context(), key)
}

// Clear removes all values stored for a given request.
func (context *BodyContext) Clear(request *http.Request) {
	getContextReadCloser(request).ClearContext()
}

// ContextReadCloser augments the io.ReadCloser interface with a Context() method.
type ContextReadCloser interface {
	io.ReadCloser
	Context() map[interface{}]interface{}
	ClearContext()
}

type contextReadCloser struct {
	io.ReadCloser
	context map[interface{}]interface{}
}

func (crc *contextReadCloser) Context() map[interface{}]interface{} {
	return crc.context
}

func (crc *contextReadCloser) ClearContext() {
	crc.context = make(map[interface{}]interface{})
}

func getContextReadCloser(request *http.Request) ContextReadCloser {
	crc, ok := request.Body.(ContextReadCloser)
	if !ok {
		crc = &contextReadCloser{
			ReadCloser: request.Body,
			context:    make(map[interface{}]interface{}),
		}
		request.Body = crc
	}
	return crc
}
