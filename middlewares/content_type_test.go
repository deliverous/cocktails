package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ContentType_WithNoSpecifiedContentType_ShouldAcceptAnything(t *testing.T) {
	handler := Chain(NewContentTypeChecker().Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json")
	acceptContentType(t, "POST", handler, "text/hml")
	acceptContentType(t, "POST", handler, "polop")
}

func Test_ContentType_WithSpecifiedContentType_ShouldAcceptThisContentTypeForAllMethod(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json")
	acceptContentType(t, "PUT", handler, "application/json")
	acceptContentType(t, "PATCH", handler, "application/json")
	acceptContentType(t, "GET", handler, "application/json")
	acceptContentType(t, "HEAD", handler, "application/json")
	acceptContentType(t, "DELETE", handler, "application/json")
}

func Test_ContentType_WithSpecifiedContentType_ShouldRejectOtherContentTypeOnlyForPostPutPatchMethods(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	rejectContentType(t, "POST", handler, "text/html")
	rejectContentType(t, "PUT", handler, "text/html")
	rejectContentType(t, "PATCH", handler, "text/html")

	acceptContentType(t, "GET", handler, "text/html")
	acceptContentType(t, "HEAD", handler, "text/html")
	acceptContentType(t, "DELETE", handler, "text/html")
}

func Test_ContentType_WithMultipleSpecifiedContentType_ShouldAcceptAll(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json", "text/xml").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json")
	acceptContentType(t, "POST", handler, "text/xml")
	rejectContentType(t, "POST", handler, "text/html")
}

func Test_ContentType_WithSpecifiedContentTypeAndCharset_ShouldAcceptThisContentTypeWithTheCharset(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-16").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-16")
	rejectContentType(t, "POST", handler, "application/json; charset=UTF-8")
}

func Test_ContentType_CharsetIsCaseInsensitive(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-16").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json; charset=utf-16")
	acceptContentType(t, "POST", handler, "application/json; charset=uTf-16")
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-16")
}

func Test_ContentType_CanSpecifyTheAssumedCharsetWhenMissingInContentType(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-8").SetAssumedCharset("UTF-16").Check).Then(nil)
	rejectContentType(t, "POST", handler, "application/json")
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-8")
}

func Test_ContentType_AssumedCharsetIsUTF8ByDefault(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-8").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json")
}

func Test_ContentType_DefaultCharsetIsUTF8(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-8")
	rejectContentType(t, "POST", handler, "application/json; charset=UTF-16")
}

func Test_ContentType_AcceptingEmptyCharset_ShouldAcceptAllCharsets(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset().Check).Then(nil)
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-8")
	acceptContentType(t, "POST", handler, "application/json; charset=UTF-16")
	acceptContentType(t, "POST", handler, "application/json; charset=any")
}

func Test_ContentType_SettingCustomErrorHandler_ShouldCallTheHandler(t *testing.T) {
	called := false
	errorHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		called = true
		writer.WriteHeader(http.StatusUnsupportedMediaType)
	})

	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetErrorHandler(errorHandler).Check).Then(nil)
	rejectContentType(t, "POST", handler, "text/html")
	expect(t, called, true)
}

func acceptContentType(t *testing.T, method string, handler http.Handler, contentType string) {
	if processRequestWithContent(t, method, handler, contentType).Code != http.StatusOK {
		t.Errorf("ContentType failure: sould have accepted %s", contentType)
	}
}

func rejectContentType(t *testing.T, method string, handler http.Handler, contentType string) {
	if processRequestWithContent(t, method, handler, contentType).Code != http.StatusUnsupportedMediaType {
		t.Errorf("ContentType failure: sould have rejected %s", contentType)
	}
}

func processRequestWithContent(t *testing.T, method string, handler http.Handler, contentType string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(method, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)
	handler.ServeHTTP(recorder, request)
	return recorder
}
