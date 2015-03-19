package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Template_Render_ShouldSetTheStatus(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender(), http.StatusOK, "main", nil)
	expect(t, recorder.Code, http.StatusOK)
}

func Test_Template_Render_WhenContentTypeAlreadySet_ShouldDoNothing(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "bla")
	staticTemplateRender().Render(recorder, http.StatusOK, "main", nil)
	expect(t, recorder.Header().Get("Content-Type"), "bla")
}

func Test_Template_Render_WhenNoContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender(), http.StatusOK, "main", nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/html; charset=UTF-8")
}

func Test_Template_Render_WithSpeficicContentType_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender().SetContentType("application/xhtml+xml"), http.StatusOK, "main", nil)
	expect(t, recorder.Header().Get("Content-Type"), "application/xhtml+xml; charset=UTF-8")
}

func Test_Template_Render_WithSpeficicCharset_ShouldSetTheContentType(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender().SetCharset("UTF-16"), http.StatusOK, "main", nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/html; charset=UTF-16")
}

func Test_Template_Render_WithEmptyCharset_ShouldUseUtf8(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender().SetCharset(""), http.StatusOK, "main", nil)
	expect(t, recorder.Header().Get("Content-Type"), "text/html; charset=UTF-8")
}

func Test_Template_Render_WithoutBinding_ShouldBeOK(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender(), http.StatusOK, "main", nil)
	expect(t, recorder.Body.String(), `<h1>Main</h1>`)
}

func Test_Template_Render_WithBinding_ShouldBeOK(t *testing.T) {
	recorder, _ := renderTemplate(staticTemplateRender(), http.StatusOK, "hello", "World")
	expect(t, recorder.Body.String(), `<h1>Hello World</h1>`)
}

func Test_Template_Render_WithSpecificDelimiter_ShouldBeOK(t *testing.T) {
	factory := NewStaticTemplateFactory("static").SetSubTemplates(func() (string, string) { return "bracket", "<h1>Hello [[.]]</h1>" })
	recorder, _ := renderTemplate(NewTemplateRender().SetFactory(factory), http.StatusCreated, "bracket", "World")
	expect(t, recorder.Body.String(), `<h1>Hello [[.]]</h1>`)
	factory.SetDelimiters("[[", "]]")
	recorder, _ = renderTemplate(NewTemplateRender().SetFactory(factory), http.StatusCreated, "bracket", "World")
	expect(t, recorder.Body.String(), `<h1>Hello World</h1>`)
}

func Test_Template_Render_WithSpecificFunctions_ShouldBeOK(t *testing.T) {
	factory := NewStaticTemplateFactory("static").SetFunctions(template.FuncMap{
		"title": strings.Title,
	}).SetSubTemplates(func() (string, string) { return "main", "<h1>{{. | title}}</h1>" })
	recorder, _ := renderTemplate(NewTemplateRender().SetFactory(factory), http.StatusCreated, "main", "hello world")
	expect(t, recorder.Body.String(), `<h1>Hello World</h1>`)
}

func Test_Template_Render_WithInvalidName_ShouldFail(t *testing.T) {
	recorder, err := renderTemplate(staticTemplateRender(), http.StatusCreated, "unknown", nil)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), ``)
	refute(t, err, nil)
}

func Test_Template_Render_InvalidTemplate_ShouldNotCompile(t *testing.T) {
	err := NewTemplateRender().SetFactory(NewStaticTemplateFactory("static").SetSubTemplates(invalid_template)).CompileTemplates()
	refute(t, err, nil)
}

func Test_Template_Render_InvalidTemplate_ShouldNotRender(t *testing.T) {
	render := NewTemplateRender().SetFactory(NewStaticTemplateFactory("static").SetSubTemplates(invalid_template))
	recorder, err := renderTemplate(render, http.StatusCreated, "invalid", nil)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), ``)
	refute(t, err, nil)
}

func staticTemplateRender() *TemplateRender {
	return NewTemplateRender().SetFactory(NewStaticTemplateFactory("static").SetSubTemplates(main_template, hello_template))
}

func main_template() (string, string) {
	return "main", `<h1>Main</h1>`
}

func hello_template() (string, string) {
	return "hello", `<h1>Hello {{.}}</h1>`
}

func invalid_template() (string, string) {
	return "invalid", "<h1>Invalid template</h1>{{.}"
}

func renderTemplate(render *TemplateRender, status int, name string, v interface{}) (*httptest.ResponseRecorder, error) {
	recorder := httptest.NewRecorder()
	err := render.Render(recorder, status, name, v)
	return recorder, err
}
