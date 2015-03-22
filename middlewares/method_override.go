package middlewares

import (
	"net/http"
)

// OverrideMethod is a middleware which checks for the X-HTTP-Method-Override header
// or the _method form key, and overrides (if valid) request.Method with its value.
func OverrideMethod() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == "POST" {
				method := request.FormValue("_method")
				if method == "" {
					method = request.Header.Get("X-HTTP-Method-Override")
				}
				if method == "PUT" || method == "PATCH" || method == "DELETE" {
					request.Method = method
				}
			}
			next.ServeHTTP(writer, request)
		})
	}
}
