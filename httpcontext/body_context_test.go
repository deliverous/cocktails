package httpcontext

import (
	"testing"
)

func Test_BodyContext_SettingAndGettingKey_ShouldBeOK(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	expect(t, context.Get(request, "key"), 1)
}

func Test_BodyContext_GettingUnknownKey_ShouldReturnsNil(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	expect(t, context.Get(request, "key"), nil)
}

func Test_BodyContext_GetOk_WithKnownKey_ShouldBeOK(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	value, ok := context.GetOk(request, "key")
	expect(t, value, 1)
	expect(t, ok, true)
}

func Test_BodyContext_GetOk_WithKnownKeyNil_ShouldReturnsNilAndTrue(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "key", nil)
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, true)
}

func Test_BodyContext_GetOk_WithUnknownKeyNil_ShouldReturnsNilAndFalse(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_BodyContext_GetAll_WithKnownRequest(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "a", 1)
	context.Set(request, "b", 2)
	values := context.GetAll(request)
	expect(t, len(values), 2)
	expect(t, values["a"], 1)
	expect(t, values["b"], 2)
}

func Test_BodyContext_GetAll_WithUnknownRequest_ShouldREturnsAnEmptyMap(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	values := context.GetAll(request)
	expect(t, len(values), 0)
}

func Test_BodyContext_Delete_OnKnownKey_ShouldBeOK(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "key", 1)
	context.Delete(request, "key")
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_BodyContext_Delete_OnUnknownKey_ShouldBeOK(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Delete(request, "key")
	value, ok := context.GetOk(request, "key")
	expect(t, value, nil)
	expect(t, ok, false)
}

func Test_BodyContext_Clear(t *testing.T) {
	context := NewBodyContext()
	request := createTestRequest()
	context.Set(request, "a", 1)
	context.Set(request, "b", 2)
	context.Clear(request)
	expect(t, len(context.GetAll(request)), 0)
}
