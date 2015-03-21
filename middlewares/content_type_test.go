package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ContentType_WithNoSpecifiedContentType_ShouldAcceptAnything(t *testing.T) {
	handler := Chain(NewContentTypeChecker().Check).Then(nil)
	acceptContentType(t, handler, "application/json")
	acceptContentType(t, handler, "text/hml")
	acceptContentType(t, handler, "polop")
}

func Test_ContentType_WithSpecifiedContentType_ShouldAcceptThisContentType(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	acceptContentType(t, handler, "application/json")
}

func Test_ContentType_WithSpecifiedContentType_ShouldRejectOtherContentType(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	rejectContentType(t, handler, "text/html")
}

func Test_ContentType_WithMultipleSpecifiedContentType_ShouldAcceptAll(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json", "text/xml").Check).Then(nil)
	acceptContentType(t, handler, "application/json")
	acceptContentType(t, handler, "text/xml")
	rejectContentType(t, handler, "text/html")
}

func Test_ContentType_WithSpecifiedContentTypeAndCharset_ShouldAcceptThisContentTypeWithTheCharset(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-16").Check).Then(nil)
	acceptContentType(t, handler, "application/json; charset=UTF-16")
	rejectContentType(t, handler, "application/json; charset=UTF-8")
}

func Test_ContentType_CharsetIsCaseInsensitive(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-16").Check).Then(nil)
	acceptContentType(t, handler, "application/json; charset=utf-16")
	acceptContentType(t, handler, "application/json; charset=uTf-16")
	acceptContentType(t, handler, "application/json; charset=UTF-16")
}

func Test_ContentType_CanSpecifyTheAssumedCharsetWhenMissingInContentType(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-8").SetAssumedCharset("UTF-16").Check).Then(nil)
	rejectContentType(t, handler, "application/json")
	acceptContentType(t, handler, "application/json; charset=UTF-8")
}

func Test_ContentType_AssumedCharsetIsUTF8ByDefault(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset("UTF-8").Check).Then(nil)
	acceptContentType(t, handler, "application/json")
}

func Test_ContentType_DefaultCharsetIsUTF8(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").Check).Then(nil)
	acceptContentType(t, handler, "application/json; charset=UTF-8")
	rejectContentType(t, handler, "application/json; charset=UTF-16")
}

func Test_ContentType_AcceptingEmptyCharset_ShouldAcceptAllCharsets(t *testing.T) {
	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetAcceptedCharset().Check).Then(nil)
	acceptContentType(t, handler, "application/json; charset=UTF-8")
	acceptContentType(t, handler, "application/json; charset=UTF-16")
	acceptContentType(t, handler, "application/json; charset=any")
}

func Test_ContentType_SettingCustomErrorHandler_ShouldCallTheHandler(t *testing.T) {
	called := false
	errorHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		called = true
		writer.WriteHeader(http.StatusUnsupportedMediaType)
	})

	handler := Chain(NewContentTypeChecker().SetAcceptedContents("application/json").SetErrorHandler(errorHandler).Check).Then(nil)
	rejectContentType(t, handler, "text/html")
	expect(t, called, true)
}

func acceptContentType(t *testing.T, handler http.Handler, contentType string) {
	if processRequestWithContent(t, handler, contentType).Code != http.StatusOK {
		t.Errorf("ContentType failure: sould have accepted %s", contentType)
	}
}

func rejectContentType(t *testing.T, handler http.Handler, contentType string) {
	if processRequestWithContent(t, handler, contentType).Code != http.StatusUnsupportedMediaType {
		t.Errorf("ContentType failure: sould have rejected %s", contentType)
	}
}

func processRequestWithContent(t *testing.T, handler http.Handler, contentType string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)
	handler.ServeHTTP(recorder, request)
	return recorder
}
