package middlewares

import (
	"net/http"
)

var (
	emptyHandler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
)

// TODO: find a better name. this is a decorator pattern

// MiddlewareBuilder is the prototype of all middleware builder.
type MiddlewareBuilder func(http.Handler) http.Handler

// MiddlewareChain if an immutable list of MiddlewareBuilder that represents a chain of middleware.
type MiddlewareChain struct {
	builders []MiddlewareBuilder
}

// Chain creates a new MiddlewareChain by only storing the builders.
func Chain(builders ...MiddlewareBuilder) MiddlewareChain {
	middlewares := MiddlewareChain{}
	middlewares.builders = append(middlewares.builders, builders...)
	return middlewares
}

// Then chains all middlewares into the result handler.
// Then() treats nil as an empty handler.
//
// A chain can be safely reused by calling Then() several times.
func (middlewares MiddlewareChain) Then(handler http.Handler) http.Handler {
	if handler == nil {
		handler = emptyHandler
	}
	for i := len(middlewares.builders) - 1; i >= 0; i-- {
		handler = middlewares.builders[i](handler)
	}
	return handler
}

// Extend creates a new MiddlewareChain by appending the builders of the given chain at the end of current chain.
//    Chain(a, b, c).Concat(Chain(d, e))
// is equivalent to:
//    Chain(a, b, c, d, e)
func (middlewares MiddlewareChain) Concat(another MiddlewareChain) MiddlewareChain {
	return middlewares.Extend(another.builders...)
}

// Extend creates a new MiddlewareChain by appending the given builders at the end of current chain.
//    Chain(a, b, c).Extend(d, e)
// is equivalent to:
//    Chain(a, b, c, d, e)
func (middlewares MiddlewareChain) Extend(builders ...MiddlewareBuilder) MiddlewareChain {
	newBuilders := make([]MiddlewareBuilder, len(middlewares.builders)+len(builders))
	copy(newBuilders, middlewares.builders)
	copy(newBuilders[len(middlewares.builders):], builders)
	return MiddlewareChain{builders: newBuilders}
}
