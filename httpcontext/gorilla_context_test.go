package httpcontext

import (
	"net/http"
	"testing"
)

func Test_Gorilla_SettingAndGettingKey_ShouldBeOK(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	expect(t, context.Get(request, "key"), 1)
}

func Test_Gorilla_GettingUnknownKey_ShouldReturnsNil(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	expect(t, context.Get(request, "key"), nil)
}

func Test_Gorilla_GetOk_WithKnownKey_ShouldBeOK(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	value, ok := context.GetOk(request, "key")
	expect(t, value, 1)
	expect(t, ok, true)
}

func Test_Gorilla_GetOk_WithKnownKeyNil_ShouldReturnsNilAndTrue(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "key", nil)
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, true)
}

func Test_Gorilla_GetOk_WithUnknownKeyNil_ShouldReturnsNilAndFalse(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_Gorilla_GetAll(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "a", 1)
	context.Set(request, "b", 2)
	values := context.GetAll(request)
	expect(t, len(values), 2)
	expect(t, values["a"], 1)
	expect(t, values["b"], 2)
}

func Test_Gorilla_Delete_OnKnownKey_ShouldBeOK(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	context.Delete(request, "key")
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_Gorilla_Delete_OnUnknownKey_ShouldBeOK(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Delete(request, "key")
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_Gorilla_Clear(t *testing.T) {
	context := NewGorillaContext()
	request := createTestRequest()
	context.Set(request, "a", 1)
	context.Set(request, "b", 2)
	context.Clear(request)
	expect(t, len(context.GetAll(request)), 0)
}

func createTestRequest() *http.Request {
	request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	return request
}

func expect(t *testing.T, value interface{}, expexted interface{}) {
	if value != expexted {
		t.Errorf("Expected %#v, got %#v.", expexted, value)
	}
}
