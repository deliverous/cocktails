package render

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JSONGreeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func Test_JSON_Render_ShouldSetTheStatus(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender(), http.StatusOK, nil)
	expect(t, recorder.Code, http.StatusOK)
}

func Test_JSON_Render_WhenContentTypeAlreadySet_ShouldDoNothing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "bla")
	NewJSONRender().Render(recorder, http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "bla")
}

func Test_JSON_Render_WhenNoContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender(), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-8")
}

func Test_JSON_Render_WithSpeficicContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender().SetContentType("my/json"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "my/json; charset=UTF-8")
}

func Test_JSON_Render_WithSpeficicCharset_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender().SetCharset("UTF-16"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-16")
}

func Test_JSON_Render_WithEmptyCharset_ShouldUseUtf8(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender().SetCharset(""), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-8")
}

func Test_JSON_Render_WithNilData_ShouldRenderNull(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender(), http.StatusOK, nil)
	expect(t, recorder.Body.String(), `null`)
}

func Test_JSON_Render_WithStruct_ShouldBeOK(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender(), http.StatusOK, JSONGreeting{"hello", "world"})
	expect(t, recorder.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_JSON_Render_WithIndent_ShouldBeOK(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender().SetIndent(true), http.StatusOK, JSONGreeting{"hello", "world"})
	expect(t, recorder.Body.String(), "{\n  \"one\": \"hello\",\n  \"two\": \"world\"\n}\n")
}

func Test_JSON_Render_WithPrefix_ShouldBeOK(t *testing.T) {
	recorder, _ := renderJSON(NewJSONRender().SetPrefix([]byte("prefix")), http.StatusOK, JSONGreeting{"hello", "world"})
	expect(t, recorder.Body.String(), `prefix{"one":"hello","two":"world"}`)
}

func Test_JSON_Render_WithInvalidData_ShouldFailWithoutWrittingAnyResponse(t *testing.T) {
	recorder, err := renderJSON(NewJSONRender(), http.StatusCreated, math.NaN)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), ``)
	refute(t, err, nil)
}

func renderJSON(render *JSONRender, status int, v interface{}) (*httptest.ResponseRecorder, error) {
	recorder := httptest.NewRecorder()
	err := render.Render(recorder, status, v)
	return recorder, err
}
