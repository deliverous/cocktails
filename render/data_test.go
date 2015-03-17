package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Data_Render_ShouldSetTheStatus(t *testing.T) {
	recorder, _ := renderData(NewDataRender(), http.StatusOK, nil)
	expect(t, recorder.Code, http.StatusOK)
}

func Test_Data_Render_WhenContentTypeAlreadySet_ShouldDoNothing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "bla")
	NewDataRender().Render(recorder, http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "bla")
}

func Test_Data_Render_WhenNoContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderData(NewDataRender(), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-8")
}

func Test_Data_Render_WithSpeficicContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderData(NewDataRender().SetContentType("my/json"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "my/json; charset=UTF-8")
}

func Test_Data_Render_WithSpeficicCharset_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderData(NewDataRender().SetCharset("UTF-16"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-16")
}

func Test_Data_Render_WithEmptyCharset_ShouldUseUtf8(t *testing.T) {
	recorder, _ := renderData(NewDataRender().SetCharset(""), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/json; charset=UTF-8")
}

func Test_Data_Render_WithData_ShouldBeOK(t *testing.T) {
	recorder, _ := renderData(NewDataRender(), http.StatusOK, []byte("Hello world"))
	expect(t, recorder.Body.String(), `Hello world`)
}

func renderData(render *DataRender, status int, data []byte) (*httptest.ResponseRecorder, error) {
	recorder := httptest.NewRecorder()
	err := render.Render(recorder, status, data)
	return recorder, err
}
