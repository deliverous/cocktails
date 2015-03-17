package render

import (
	"encoding/xml"
	"net/http"
)

// XML built-in renderer.
type XMLRender struct {
	ContentType string
	Charset     string
	Indent      bool
	Prefix      []byte
}

func (render *XMLRender) SetContentType(value string) *XMLRender {
	render.ContentType = value
	return render
}

func (render *XMLRender) SetCharset(value string) *XMLRender {
	render.Charset = value
	return render
}

func (render *XMLRender) SetIndent(value bool) *XMLRender {
	render.Indent = value
	return render
}

func (render *XMLRender) SetPrefix(prefix []byte) *XMLRender {
	render.Prefix = prefix
	return render
}

// NewXMLRender creates a new XML render with default values.
func NewXMLRender() *XMLRender {
	return &XMLRender{
		ContentType: "text/xml",
		Charset:     "UTF-8",
		Indent:      false,
	}
}

// Render an XML response.
func (render *XMLRender) Render(writer http.ResponseWriter, status int, v interface{}) error {
	result, err := marshallToXml(v, render.Indent)
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

func marshallToXml(v interface{}, indent bool) ([]byte, error) {
	if indent {
		return xml.MarshalIndent(v, "", "  ")
	} else {
		return xml.Marshal(v)
	}
}
