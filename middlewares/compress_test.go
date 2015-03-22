package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Compress(t *testing.T) {
	w := httptest.NewRecorder()
	compressedRequest(w, "gzip")
	if w.HeaderMap.Get("Content-Encoding") != "gzip" {
		t.Fatalf("wrong content encoding, got %d want %d", w.HeaderMap.Get("Content-Encoding"), "gzip")
	}
	if w.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("wrong content type, got %s want %s", w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	}
	reader, _ := gzip.NewReader(bytes.NewReader(w.Body.Bytes()))
	buf := make([]byte, 10000)
	n, _ := reader.Read(buf)
	if string(buf[:n]) != buffer {
		fmt.Printf("%#v\n", string(buf))
		t.Fatalf("wrong buffer")
	}
}

func compressedRequest(writer *httptest.ResponseRecorder, compression string) {
	handler := Chain(Compress()).Then(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, buffer)
		writer.(http.Flusher).Flush()
	}))

	handler.ServeHTTP(writer, &http.Request{
		Method: "GET",
		Header: http.Header{
			"Accept-Encoding": []string{compression},
		},
	})
}

const buffer = `
     1  // Copyright 2010 The Go Authors. All rights reserved.
     2  // Use of this source code is governed by a BSD-style
     3  // license that can be found in the LICENSE file.
     4
     5  package gzip
     6
     7  import (
     8      "compress/flate"
     9      "errors"
    10      "fmt"
    11      "hash"
    12      "hash/crc32"
    13      "io"
    14  )
    15
    16  // These constants are copied from the flate package, so that code that imports
    17  // "compress/gzip" does not also have to import "compress/flate".
    18  const (
    19      NoCompression      = flate.NoCompression
    20      BestSpeed          = flate.BestSpeed
    21      BestCompression    = flate.BestCompression
    22      DefaultCompression = flate.DefaultCompression
    23  )
    24
    25  // A Writer is an io.WriteCloser.
    26  // Writes to a Writer are compressed and written to w.
    27  type Writer struct {
    28      Header
    29      w           io.Writer
    30      level       int
    31      wroteHeader bool
    32      compressor  *flate.Writer
    33      digest      hash.Hash32
    34      size        uint32
    35      closed      bool
    36      buf         [10]byte
    37      err         error
    38  }`
