package httpcontext

import (
	"net/http"
	"sync"
)

// BigMap context stores values in a big map
type BigMapContext struct {
	mutex sync.RWMutex
	data  map[*http.Request]map[interface{}]interface{}
}

// NewBigMapContext creates a new BigMapContext
func NewBigMapContext() *BigMapContext {
	return &BigMapContext{
		data:  make(map[*http.Request]map[interface{}]interface{}),
	}
}

// Set stores a value for a given key in a given request.
func (context *BigMapContext) Set(request *http.Request, key, val interface{}) {
	context.mutex.Lock()
	if context.data[request] == nil {
		context.data[request] = make(map[interface{}]interface{})
	}
	context.data[request][key] = val
	context.mutex.Unlock()
}

// Get returns a value stored for a given key in a given request.
func (context *BigMapContext) Get(request *http.Request, key interface{}) interface{} {
	context.mutex.RLock()
	if ctx := context.data[request]; ctx != nil {
		value := ctx[key]
		context.mutex.RUnlock()
		return value
	}
	context.mutex.RUnlock()
	return nil
}

// GetOk returns stored value and presence state like multi-value return of map access.
func (context *BigMapContext) GetOk(request *http.Request, key interface{}) (interface{}, bool) {
	context.mutex.RLock()
	if _, ok := context.data[request]; ok {
		value, ok := context.data[request][key]
		context.mutex.RUnlock()
		return value, ok
	}
	context.mutex.RUnlock()
	return nil, false
}

// GetAll returns all stored values for the request as a map.
func (context *BigMapContext) GetAll(request *http.Request) map[interface{}]interface{} {
	context.mutex.RLock()
	if all, ok := context.data[request]; ok {
		result := make(map[interface{}]interface{}, len(all))
		for k, v := range all {
			result[k] = v
		}
		context.mutex.RUnlock()
		return result
	}
	context.mutex.RUnlock()
	return nil
}

// Delete removes a value stored for a given key in a given request.
func (context *BigMapContext) Delete(request *http.Request, key interface{}) {
	context.mutex.Lock()
	if context.data[request] != nil {
		delete(context.data[request], key)
	}
	context.mutex.Unlock()
}

// Clear removes all values stored for a given request.
func (context *BigMapContext) Clear(request *http.Request) {
	context.mutex.Lock()
	delete(context.data, request)
	context.mutex.Unlock()
}

// ClearBigMapContext is a middleware to cleanup request context at the end
func ClearBigMapContext(context *BigMapContext) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer context.Clear(request)
			handler.ServeHTTP(writer, request)
		})
	}
}

// This work is based on the gorilla context: https://github.com/gorilla/context
