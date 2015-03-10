package context

import (
	"net/http"
	"sync"
	"time"
)

// Gorilla context store values in a big map
type GorillaContext struct {
	mutex sync.RWMutex
	data  map[*http.Request]map[interface{}]interface{}
	datat map[*http.Request]int64
}

// NewGorillaContext creates a new GorillaContext
func NewGorillaContext() *GorillaContext {
	return &GorillaContext{
		data:  make(map[*http.Request]map[interface{}]interface{}),
		datat: make(map[*http.Request]int64),
	}
}

// Set stores a value for a given key in a given request.
func (context *GorillaContext) Set(request *http.Request, key, val interface{}) {
	context.mutex.Lock()
	if context.data[request] == nil {
		context.data[request] = make(map[interface{}]interface{})
		context.datat[request] = time.Now().Unix()
	}
	context.data[request][key] = val
	context.mutex.Unlock()
}

// Get returns a value stored for a given key in a given request.
func (context *GorillaContext) Get(request *http.Request, key interface{}) interface{} {
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
func (context *GorillaContext) GetOk(request *http.Request, key interface{}) (interface{}, bool) {
	context.mutex.RLock()
	if _, ok := context.data[request]; ok {
		value, ok := context.data[request][key]
		context.mutex.RUnlock()
		return value, ok
	}
	context.mutex.RUnlock()
	return nil, false
}

// GetAll returns all stored values for the request as a map. Nil is returned for invalid requests.
func (context *GorillaContext) GetAll(request *http.Request) map[interface{}]interface{} {
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
func (context *GorillaContext) Delete(request *http.Request, key interface{}) {
	context.mutex.Lock()
	if context.data[request] != nil {
		delete(context.data[request], key)
	}
	context.mutex.Unlock()
}

// Clear removes all values stored for a given request.
func (context *GorillaContext) Clear(request *http.Request) {
	context.mutex.Lock()
	delete(context.data, request)
	delete(context.datat, request)
	context.mutex.Unlock()
}

// ClearGorillaContext is a middleware to cleanup request context at the end
func ClearGorillaContext(context *GorillaContext) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer context.Clear(request)
			handler.ServeHTTP(writer, request)
		})
	}
}

// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
