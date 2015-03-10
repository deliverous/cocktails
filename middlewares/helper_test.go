package middlewares

import (
	"net/http"
	"net/http/httptest"
)

func processRequest(handler http.Handler) (*httptest.ResponseRecorder, error) {
	writer := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	handler.ServeHTTP(writer, request)
	return writer, nil
}
