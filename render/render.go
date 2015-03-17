package render

import (
	"net/http"
)

func writeHeader(writer http.ResponseWriter, status int, contentType string, charset string) {
	if writer.Header().Get("Content-Type") != "" {
		return
	}
	if charset == "" {
		charset = "UTF-8"
	}
	writer.Header().Set("Content-Type", contentType+"; charset="+charset)
	writer.WriteHeader(status)
}
