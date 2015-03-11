package httpcontext

import (
	"net/http"
	"testing"
)

func createTestRequest() *http.Request {
	request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	return request
}

func expect(t *testing.T, value interface{}, expexted interface{}) {
	if value != expexted {
		t.Errorf("Expected %#v, got %#v.", expexted, value)
	}
}
