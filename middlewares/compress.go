package middlewares

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"
)

// Compress is responsible for compressing the payload with gzip or deflate and setting the proper
// headers when supported by the client.
func Compress() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		L:
			for _, encoding := range strings.Split(request.Header.Get("Accept-Encoding"), ",") {
				switch strings.TrimSpace(encoding) {
				case "gzip":
					writer.Header().Set("Content-Encoding", "gzip")
					writer.Header().Add("Vary", "Accept-Encoding")
					writer = &compressedResponseWriter{
						ResponseWriter: writer,
						compressWriter: gzip.NewWriter(writer),
						wroteHeader:    false,
					}
					break L
				case "deflate":
					writer.Header().Set("Content-Encoding", "deflate")
					writer.Header().Add("Vary", "Accept-Encoding")
					cw, _ := flate.NewWriter(writer, flate.DefaultCompression)
					writer = &compressedResponseWriter{
						ResponseWriter: writer,
						compressWriter: cw,
						wroteHeader:    false,
					}
					break L
				}
			}
			next.ServeHTTP(writer, request)
		})
	}
}

type flusher interface {
	Flush() error
}

// Private responseWriter intantiated by the gzip middleware.
// It encodes the payload with gzip and set the proper headers.
// It implements the following interfaces:
// http.ResponseWriter
// http.Flusher
// http.CloseNotifier
// http.Hijacker
type compressedResponseWriter struct {
	http.ResponseWriter
	compressWriter io.WriteCloser
	wroteHeader    bool
}

func (writer *compressedResponseWriter) Header() http.Header {
	return writer.ResponseWriter.Header()
}

// Set the right headers for compressed encoded responses.
func (writer *compressedResponseWriter) WriteHeader(code int) {
	writer.ResponseWriter.WriteHeader(code)
	writer.wroteHeader = true
}

// Make sure the local WriteHeader is called, and call the parent Flush.
// Provided in order to implement the http.Flusher interface.
func (writer *compressedResponseWriter) Flush() {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}
	writer.compressWriter.(flusher).Flush()
	writer.ResponseWriter.(http.Flusher).Flush()
}

// Call the parent CloseNotify.
// Provided in order to implement the http.CloseNotifier interface.
func (writer *compressedResponseWriter) CloseNotify() <-chan bool {
	return writer.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Provided in order to implement the http.Hijacker interface.
func (writer *compressedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return writer.ResponseWriter.(http.Hijacker).Hijack()
}

// Make sure the local WriteHeader is called, and encode the payload if necessary.
// Provided in order to implement the http.ResponseWriter interface.
func (writer *compressedResponseWriter) Write(b []byte) (int, error) {
	if writer.Header().Get("Content-Type") == "" {
		writer.Header().Set("Content-Type", http.DetectContentType(b))
	}
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}
	return writer.compressWriter.Write(b)
}
