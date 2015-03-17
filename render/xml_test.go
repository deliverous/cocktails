package render

import (
	"encoding/xml"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

type XMLGreeting struct {
	XMLName xml.Name `xml:"greeting"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func Test_XML_Render_ShouldSetTheStatus(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender(), http.StatusOK, nil)
	expect(t, recorder.Code, http.StatusOK)
}

func Test_XML_Render_WhenContentTypeAlreadySet_ShouldDoNothing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "bla")
	NewXMLRender().Render(recorder, http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "bla")
}

func Test_XML_Render_WhenNoContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender(), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/xml; charset=UTF-8")
}

func Test_XML_Render_WithSpeficicContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender().SetContentType("my/xml"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "my/xml; charset=UTF-8")
}

func Test_XML_Render_WithSpeficicCharset_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender().SetCharset("UTF-16"), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/xml; charset=UTF-16")
}

func Test_XML_Render_WithEmptyCharset_ShouldUseUtf8(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender().SetCharset(""), http.StatusOK, nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/xml; charset=UTF-8")
}

func Test_XML_Render_WithNilData_ShouldRenderEmpty(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender(), http.StatusOK, nil)
	expect(t, recorder.Body.String(), ``)
}

func Test_XML_Render_WithStruct_ShouldBeOK(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender(), http.StatusOK, XMLGreeting{One: "hello", Two: "world"})
	expect(t, recorder.Body.String(), `<greeting one="hello" two="world"></greeting>`)
}

func Test_XML_Render_WithIndent_ShouldBeOK(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender().SetIndent(true), http.StatusOK, XMLGreeting{One: "hello", Two: "world"})
	expect(t, recorder.Body.String(), "<greeting one=\"hello\" two=\"world\"></greeting>\n")
}

func Test_XML_Render_WithPrefix_ShouldBeOK(t *testing.T) {
	recorder, _ := renderXML(NewXMLRender().SetPrefix([]byte("prefix")), http.StatusOK, XMLGreeting{One: "hello", Two: "world"})
	expect(t, recorder.Body.String(), `prefix<greeting one="hello" two="world"></greeting>`)
}

func Test_XML_Render_WithInvalidData_ShouldFailWithoutWrittingAnyResponse(t *testing.T) {
	recorder, err := renderXML(NewXMLRender(), http.StatusCreated, math.NaN)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), ``)
	refute(t, err, nil)
}

func renderXML(render *XMLRender, status int, v interface{}) (*httptest.ResponseRecorder, error) {
	recorder := httptest.NewRecorder()
	err := render.Render(recorder, status, v)
	return recorder, err
}
