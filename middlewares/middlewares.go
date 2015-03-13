package middlewares

import (
	"net/http"
)

var (
	emptyHandler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
)

// Middleware is the prototype of all middleware.
type Middleware func(http.Handler) http.Handler

// MiddlewareChain if an immutable list of Middleware that represents a chain of middleware.
type MiddlewareChain struct {
	middlewares []Middleware
}

// Chain creates a new MiddlewareChain by only storing the middlewares.
func Chain(middlewares ...Middleware) MiddlewareChain {
	return MiddlewareChain{middlewares: middlewares}
}

// Then chains all middlewares into the result handler.
// Then() treats nil as an empty handler.
//
// A chain can be safely reused by calling Then() several times.
func (chain MiddlewareChain) Then(handler http.Handler) http.Handler {
	if handler == nil {
		handler = emptyHandler
	}
	for i := len(chain.middlewares) - 1; i >= 0; i-- {
		handler = chain.middlewares[i](handler)
	}
	return handler
}

// Extend creates a new MiddlewareChain by appending the middlewares of the given chain at the end of current chain.
//    Chain(a, b, c).Concat(Chain(d, e))
// returns a new chain equivalent to:
//    Chain(a, b, c, d, e)
func (chain MiddlewareChain) Concat(another MiddlewareChain) MiddlewareChain {
	return chain.Extend(another.middlewares...)
}

// Extend creates a new MiddlewareChain by appending the given middlewares at the end of current chain.
//    Chain(a, b, c).Extend(d, e)
// returns a new chain equivalent to:
//    Chain(a, b, c, d, e)
func (chain MiddlewareChain) Extend(middlewares ...Middleware) MiddlewareChain {
	newMiddlewares := make([]Middleware, len(chain.middlewares)+len(middlewares))
	copy(newMiddlewares, chain.middlewares)
	copy(newMiddlewares[len(chain.middlewares):], middlewares)
	return MiddlewareChain{middlewares: newMiddlewares}
}
