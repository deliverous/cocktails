package render

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

type JSONPGreeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func Test_JSONP_Render_ShouldSetTheStatus(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender(), http.StatusOK, nil)
	expect(t, recorder.Code, http.StatusOK)
}

func Test_JSONP_Render_WhenContentTypeAlreadySet_ShouldDoNothing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "bla")
	NewJSONPRender().Render(recorder, http.StatusOK, "claback", nil)
	expect(t, recorder.Header().Get("Content-Type"), "bla")
}

func Test_JSONP_Render_WhenNoContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender(), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/javascript; charset=UTF-8")
}

func Test_JSONP_Render_WithSpeficicContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender().SetContentType("my/json"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "my/json; charset=UTF-8")
}

func Test_JSONP_Render_WithSpeficicCharset_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender().SetCharset("UTF-16"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/javascript; charset=UTF-16")
}

func Test_JSONP_Render_WithEmptyCharset_ShouldUseUtf8(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender().SetCharset(""), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/javascript; charset=UTF-8")
}

func Test_JSONP_Render_WithNilData_ShouldRenderNull(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender(), http.StatusOK, nil)
	expect(t, recorder.Body.String(), `callback(null);`)
}

func Test_JSONP_Render_WithStruct_ShouldBeOK(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender(), http.StatusOK, JSONPGreeting{"hello", "world"})
	expect(t, recorder.Body.String(), `callback({"one":"hello","two":"world"});`)
}

func Test_JSONP_Render_WithIndent_ShouldBeOK(t *testing.T) {
	recorder, _ := renderJSONP(NewJSONPRender().SetIndent(true), http.StatusOK, JSONPGreeting{"hello", "world"})
	expect(t, recorder.Body.String(), "callback({\n  \"one\": \"hello\",\n  \"two\": \"world\"\n});\n")
}

func Test_JSONP_Render_WithInvalidData_ShouldFailWithoutWrittingAnyResponse(t *testing.T) {
	recorder, err := renderJSONP(NewJSONPRender(), http.StatusCreated, math.NaN)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), ``)
	refute(t, err, nil)
}

func renderJSONP(render *JSONPRender, status int, v interface{}) (*httptest.ResponseRecorder, error) {
	recorder := httptest.NewRecorder()
	err := render.Render(recorder, status, "callback", v)
	return recorder, err
}
