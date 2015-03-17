package render

import (
	"encoding/json"
	"net/http"
)

type JSONRender struct {
	ContentType string
	Charset     string
	Indent      bool
	Prefix      []byte
}

func (render *JSONRender) SetContentType(value string) *JSONRender {
	render.ContentType = value
	return render
}

func (render *JSONRender) SetCharset(value string) *JSONRender {
	render.Charset = value
	return render
}

func (render *JSONRender) SetIndent(value bool) *JSONRender {
	render.Indent = value
	return render
}

func (render *JSONRender) SetPrefix(prefix []byte) *JSONRender {
	render.Prefix = prefix
	return render
}

// NewJSONRender creates a new JSON render with default values.
func NewJSONRender() *JSONRender {
	return &JSONRender{
		ContentType: "application/json",
		Charset:     "UTF-8",
		Indent:      false,
	}
}

// Render a JSON response.
func (render *JSONRender) Render(writer http.ResponseWriter, status int, v interface{}) error {
	result, err := marshallToJson(v, render.Indent)
	if err != nil {
		return err
	}

	writeHeader(writer, status, render.ContentType, render.Charset)
	if len(render.Prefix) > 0 {
		writer.Write(render.Prefix)
	}
	writer.Write(result)
	if render.Indent {
		writer.Write([]byte("\n"))
	}
	return nil
}

func marshallToJson(v interface{}, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(v, "", "  ")
	} else {
		return json.Marshal(v)
	}
}
