package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func expect(t *testing.T, value interface{}, expexted interface{}) {
	if value != expexted {
		t.Errorf("Expected %#v, got %#v.", expexted, value)
	}
}

func processRequest(t *testing.T, handler http.Handler) *httptest.ResponseRecorder {
	writer := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(writer, request)
	return writer
}
