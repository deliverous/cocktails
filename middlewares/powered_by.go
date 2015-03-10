package middlewares

import (
	"net/http"
)

// PoweredBy is a middleware to add the X-Powered-By header
func PoweredBy(tag string) MiddlewareBuilder {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("X-Powered-By", tag)
			next.ServeHTTP(writer, request)
		})
	}
}
