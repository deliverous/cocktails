package middlewares

import (
	"mime"
	"net/http"
	"strings"
)

type ContentTypeChecker struct {
	AcceptedContents []string
	AcceptedCharsets []string
	AssumedCharset   string
	ErrorHandler     http.Handler
}

func NewContentTypeChecker() *ContentTypeChecker {
	return &ContentTypeChecker{
		AssumedCharset:   "UTF-8",
		AcceptedCharsets: []string{"UTF-8"},
		ErrorHandler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
		}),
	}
}

func (m *ContentTypeChecker) SetAcceptedContents(contents ...string) *ContentTypeChecker {
	m.AcceptedContents = contents
	return m
}

func (m *ContentTypeChecker) SetAcceptedCharset(charsets ...string) *ContentTypeChecker {
	m.AcceptedCharsets = make([]string, len(charsets))
	for i, charset := range charsets {
		m.AcceptedCharsets[i] = strings.ToUpper(charset)
	}
	return m
}

func (m *ContentTypeChecker) SetAssumedCharset(charset string) *ContentTypeChecker {
	m.AssumedCharset = charset
	return m
}

func (m *ContentTypeChecker) SetErrorHandler(handler http.Handler) *ContentTypeChecker {
	m.ErrorHandler = handler
	return m
}

func (m *ContentTypeChecker) acceptContent(value string) bool {
	if len(m.AcceptedContents) == 0 {
		return true
	}
	for _, content := range m.AcceptedContents {
		if value == content {
			return true
		}
	}
	return false
}

func (m *ContentTypeChecker) acceptCharset(value string) bool {
	if len(m.AcceptedCharsets) == 0 {
		return true
	}
	value = strings.ToUpper(value)
	for _, charset := range m.AcceptedCharsets {
		if value == charset {
			return true
		}
	}
	return false
}

func (m *ContentTypeChecker) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !(request.Method == "PUT" || request.Method == "POST" || request.Method == "PATCH") {
			next.ServeHTTP(writer, request)
			return
		}

		mediatype, params, _ := mime.ParseMediaType(request.Header.Get("Content-Type"))
		charset, ok := params["charset"]
		if !ok {
			charset = m.AssumedCharset
		}
		if m.acceptContent(mediatype) && m.acceptCharset(charset) {
			next.ServeHTTP(writer, request)
		} else {
			m.ErrorHandler.ServeHTTP(writer, request)
		}
	})
}
