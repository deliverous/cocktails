package render

import (
	"net/http"
)

type JSONPRender struct {
	ContentType string
	Charset     string
	Indent      bool
}

func (render *JSONPRender) SetContentType(value string) *JSONPRender {
	render.ContentType = value
	return render
}

func (render *JSONPRender) SetCharset(value string) *JSONPRender {
	render.Charset = value
	return render
}

func (render *JSONPRender) SetIndent(value bool) *JSONPRender {
	render.Indent = value
	return render
}

// NewJSONPRender creates a new JSONP render with default values.
func NewJSONPRender() *JSONPRender {
	return &JSONPRender{
		ContentType: "application/javascript",
		Charset:     "UTF-8",
		Indent:      false,
	}
}

// Render a JSONP response.
func (render *JSONPRender) Render(writer http.ResponseWriter, status int, callback string, v interface{}) error {
	result, err := marshallToJson(v, render.Indent)
	if err != nil {
		return err
	}

	writeHeader(writer, status, render.ContentType, render.Charset)
	writer.Write([]byte(callback + "("))
	writer.Write(result)
	writer.Write([]byte(");"))
	if render.Indent {
		writer.Write([]byte("\n"))
	}
	return nil
}
