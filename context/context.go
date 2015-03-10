package context

import (
	"net/http"
)

type Context interface {
	// Set stores a value for a given key in a given request.
	Set(request *http.Request, key, value interface{})

	// Get returns a value stored for a given key in a given request.
	Get(request *http.Request, key interface{}) interface{}

	// GetOk returns stored value and presence state like multi-value return of map access.
	GetOk(request *http.Request, key interface{}) (interface{}, bool)

	// GetAll returns all stored values for the request as a map. Nil is returned for invalid requests.
	GetAll(request *http.Request, key interface{}) map[interface{}]interface{}

	// Delete removes a value stored for a given key in a given request.
	Delete(request *http.Request, key interface{})

	// Clear removes all values stored for a given request.
	Clear(request *http.Request)
}
