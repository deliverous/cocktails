package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_OverrideMethod_WhenPostWithHeaderOverride_ShouldOverrideOnlyPutPatchDelete(t *testing.T) {
	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "DELETE"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "PATCH"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "PUT"}, "PUT")

	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "GET"}, "POST")
	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "HEAD"}, "POST")
}

func Test_OverrideMethod_WhenPostWithFormOverride_ShouldOverrideOnlyPutPatchDelete(t *testing.T) {
	expectMethodIs(t, overrideDefinition{Method: "POST", FormOverride: "DELETE"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "POST", FormOverride: "PATCH"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "POST", FormOverride: "PUT"}, "PUT")

	expectMethodIs(t, overrideDefinition{Method: "POST", FormOverride: "GET"}, "POST")
	expectMethodIs(t, overrideDefinition{Method: "POST", FormOverride: "HEAD"}, "POST")
}

func Test_OverrideMethod_FormOverrideTakesPrecedenceOnHeaderOverride(t *testing.T) {
	expectMethodIs(t, overrideDefinition{Method: "POST", HeaderOverride: "PUT", FormOverride: "DELETE"}, "DELETE")
}

func Test_OverrideMethod_WhenNotPostWithHeaderOverride_ShouldNotOverride(t *testing.T) {
	// It isn't secure to override e.g a GET to a POST, so only POST requests are considered.
	// Likewise, the override method can only be a "write" method: PUT, PATCH or DELETE.

	expectMethodIs(t, overrideDefinition{Method: "DELETE", HeaderOverride: "GET"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "DELETE", HeaderOverride: "HEAD"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "DELETE", HeaderOverride: "PATCH"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "DELETE", HeaderOverride: "POST"}, "DELETE")
	expectMethodIs(t, overrideDefinition{Method: "DELETE", HeaderOverride: "PUT"}, "DELETE")

	expectMethodIs(t, overrideDefinition{Method: "GET", HeaderOverride: "DELETE"}, "GET")
	expectMethodIs(t, overrideDefinition{Method: "GET", HeaderOverride: "HEAD"}, "GET")
	expectMethodIs(t, overrideDefinition{Method: "GET", HeaderOverride: "PATCH"}, "GET")
	expectMethodIs(t, overrideDefinition{Method: "GET", HeaderOverride: "POST"}, "GET")
	expectMethodIs(t, overrideDefinition{Method: "GET", HeaderOverride: "PUT"}, "GET")

	expectMethodIs(t, overrideDefinition{Method: "HEAD", HeaderOverride: "DELETE"}, "HEAD")
	expectMethodIs(t, overrideDefinition{Method: "HEAD", HeaderOverride: "GET"}, "HEAD")
	expectMethodIs(t, overrideDefinition{Method: "HEAD", HeaderOverride: "PATCH"}, "HEAD")
	expectMethodIs(t, overrideDefinition{Method: "HEAD", HeaderOverride: "POST"}, "HEAD")
	expectMethodIs(t, overrideDefinition{Method: "HEAD", HeaderOverride: "PUT"}, "HEAD")

	expectMethodIs(t, overrideDefinition{Method: "PATCH", HeaderOverride: "DELETE"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "PATCH", HeaderOverride: "GET"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "PATCH", HeaderOverride: "HEAD"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "PATCH", HeaderOverride: "POST"}, "PATCH")
	expectMethodIs(t, overrideDefinition{Method: "PATCH", HeaderOverride: "PUT"}, "PATCH")

	expectMethodIs(t, overrideDefinition{Method: "PUT", HeaderOverride: "DELETE"}, "PUT")
	expectMethodIs(t, overrideDefinition{Method: "PUT", HeaderOverride: "GET"}, "PUT")
	expectMethodIs(t, overrideDefinition{Method: "PUT", HeaderOverride: "HEAD"}, "PUT")
	expectMethodIs(t, overrideDefinition{Method: "PUT", HeaderOverride: "PATCH"}, "PUT")
	expectMethodIs(t, overrideDefinition{Method: "PUT", HeaderOverride: "POST"}, "PUT")
}

type overrideDefinition struct {
	Method         string
	HeaderOverride string
	FormOverride   string
}

func expectMethodIs(t *testing.T, definition overrideDefinition, expectedMethod string) {
	var recordedMethod string
	handler := Chain(OverrideMethod()).Then(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		recordedMethod = request.Method
	}))

	request, err := http.NewRequest(definition.Method, "/", nil)
	if definition.HeaderOverride != "" {
		request.Header.Set("X-HTTP-Method-Override", definition.HeaderOverride)
	}
	if definition.FormOverride != "" {
		values := url.Values{}
		values.Add("_method", definition.FormOverride)
		request.Form = values
	}
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if recordedMethod != expectedMethod {
		t.Errorf("Bad method override: expected %#v, got %#v", expectedMethod, recordedMethod)
	}
}
