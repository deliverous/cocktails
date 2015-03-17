package render

import (
	"net/http"
)

// Data built-in renderer.
type DataRender struct {
	ContentType string
	Charset     string
}

func (render *DataRender) SetContentType(value string) *DataRender {
	render.ContentType = value
	return render
}

func (render *DataRender) SetCharset(value string) *DataRender {
	render.Charset = value
	return render
}

// NewDataRender creates a new Data render with default values.
func NewDataRender() *DataRender {
	return &DataRender{
		ContentType: "application/json",
		Charset:     "UTF-8",
	}
}

// Render a data response.
func (render *DataRender) Render(writer http.ResponseWriter, status int, data []byte) error {
	writeHeader(writer, status, render.ContentType, render.Charset)
	writer.Write(data)
	return nil
}
