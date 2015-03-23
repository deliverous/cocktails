package middlewares

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Compress_WhenGzipIsAccepted_ShouldUseGzip(t *testing.T) {
	ensureCompression(t, []string{"gzip"}, "gzip", gunzip)
}

func Test_Compress_WhenDeflateIsAccepted_ShouldUseDeflate(t *testing.T) {
	ensureCompression(t, []string{"deflate"}, "deflate", unflate)
}

func Test_Compress_WhenMultipleEncodingAccepted_ShouldUseFirst(t *testing.T) {
	ensureCompression(t, []string{"gzip", "deflate"}, "gzip", gunzip)
	ensureCompression(t, []string{"deflate", "gzip"}, "deflate", unflate)
}

func Test_Compress_WhenNoEncodingSpecified_ShouldNotCompress(t *testing.T) {
	ensureCompression(t, []string{}, "", raw)
}

func Test_Compress_WhenUnknownEncodingSpecified_ShouldNotCompress(t *testing.T) {
	ensureCompression(t, []string{"bzip"}, "", raw)
}

func Test_Compress_WhenCompression_ShouldDetectContentType(t *testing.T) {
	recorder := httptest.NewRecorder()
	compressedRequest(recorder, []string{"gzip"})
	if recorder.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("Wrong content type, expected %#v, got %#v", "text/plain; charset=utf-8", recorder.HeaderMap.Get("Content-Type"))
	}
}

func ensureCompression(t *testing.T, acceptedEncoding []string, expectedEncoding string, decoder func([]byte) string) {
	recorder := httptest.NewRecorder()
	compressedRequest(recorder, acceptedEncoding)
	if recorder.HeaderMap.Get("Content-Encoding") != expectedEncoding {
		t.Errorf("Wrong content encoding, expected %#v, got %#v", expectedEncoding, recorder.HeaderMap.Get("Content-Encoding"))
	}
	data := decoder(recorder.Body.Bytes())
	if data != bodyContent {
		t.Errorf("Wrong content, expected %d chars, got %d chars", len(bodyContent), len(data))
	}
}

func compressedRequest(recorder *httptest.ResponseRecorder, acceptedEncoding []string) {
	handler := Chain(Compress()).Then(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, bodyContent)
		writer.(http.Flusher).Flush()
	}))
	handler.ServeHTTP(recorder, &http.Request{
		Method: "GET",
		Header: http.Header{
			"Accept-Encoding": acceptedEncoding,
		},
	})
}

func gunzip(data []byte) string {
	buffer := new(bytes.Buffer)
	reader, _ := gzip.NewReader(bytes.NewReader(data))
	buffer.ReadFrom(reader)
	return buffer.String()
}

func unflate(data []byte) string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(flate.NewReader(bytes.NewReader(data)))
	return buffer.String()
}

func raw(data []byte) string {
	return string(data)
}

const bodyContent = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed
in massa quis nisl vestibulum vehicula. Aenean a augue eu metus congue
placerat. Class aptent taciti sociosqu ad litora torquent per conubia nostra,
per inceptos himenaeos. Vivamus eleifend sodales facilisis. Mauris tempor
tempor condimentum. Curabitur ac lorem nec tortor imperdiet varius sodales
mollis nisi. Aliquam ut nisl ut dui blandit faucibus. Phasellus nec mauris
turpis. Suspendisse dapibus auctor est eget vulputate. Suspendisse maximus,
enim a vehicula rutrum, quam diam vehicula purus, ut fermentum leo augue sed
sapien.

Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac
turpis egestas. Fusce nec aliquam augue, nec vulputate ligula. Curabitur eget
tellus ipsum. Maecenas mattis lorem nec magna luctus, et rhoncus orci lacinia.
Vivamus lobortis ipsum at massa maximus pulvinar. Vivamus non ante et enim
congue aliquet a in arcu. Nunc sollicitudin purus leo, sed varius leo semper
quis. Mauris nulla sem, auctor volutpat condimentum a, pharetra posuere
lectus. Integer tempor lectus sit amet hendrerit pretium. Ut elementum mollis
tincidunt.

Etiam pellentesque sagittis turpis, non rutrum velit dapibus ultrices.
Praesent lobortis elit eget accumsan fermentum. Morbi rutrum imperdiet risus
sed efficitur. Morbi cursus sapien ac neque facilisis efficitur. Nulla aliquet
libero et quam auctor, blandit convallis lorem pellentesque. Curabitur
convallis finibus elit ut suscipit. Aenean sem augue, tempor sed odio quis,
pulvinar tincidunt turpis. Aliquam erat volutpat. Donec libero ipsum, posuere
pretium euismod non, lobortis ac leo. Nullam suscipit massa eu tristique
euismod. Praesent congue dictum egestas. Sed ac facilisis enim, a venenatis
purus. Vivamus pharetra venenatis odio, sit amet vehicula ex lacinia non.
Mauris congue risus diam, ut commodo leo lacinia id. Aenean consectetur, est
sit amet tristique ultricies, arcu orci rutrum velit, et auctor lectus mauris
eget risus.

Integer ornare nunc orci, a placerat diam convallis sit amet. Aenean augue
nisl, imperdiet vel faucibus eget, elementum id nunc. Aliquam facilisis lorem
ante, non elementum erat molestie at. Nulla dapibus pulvinar luctus. Curabitur
dignissim pretium eros, id vestibulum leo vulputate non. Phasellus lacinia,
arcu sed lobortis gravida, arcu velit viverra tellus, ullamcorper commodo quam
tellus et orci. Vivamus eget ultrices felis. Lorem ipsum dolor sit amet,
consectetur adipiscing elit. Duis feugiat volutpat turpis, nec sodales dui
fringilla in.

Integer id leo non justo finibus dictum id in dui. Vivamus sed finibus nisi.
Sed in diam in mauris dignissim molestie. Aliquam erat volutpat. Phasellus et
sollicitudin velit, vitae pharetra quam. Nunc vulputate varius massa, vitae
mattis turpis rutrum eget. Fusce velit justo, pharetra sit amet placerat at,
ultrices et magna. Nunc varius justo sit amet iaculis egestas. Integer in
lacus sit amet sem finibus mattis sit amet efficitur ante. Aliquam interdum
metus nec facilisis condimentum.

Proin eleifend non justo in varius. In iaculis ex lacus, vel volutpat lorem
ultricies id. Mauris interdum nulla rhoncus posuere tempus. Nulla nec eros sit
amet nibh blandit mattis. Integer elementum sem nec ligula posuere imperdiet.
Sed sagittis ipsum mi, sed lobortis urna congue ut. Phasellus eget sodales
velit. Fusce pellentesque ipsum non nibh luctus lobortis. Interdum et
malesuada fames ac ante ipsum primis in faucibus. Nam convallis massa vel
pretium euismod. Nulla aliquet faucibus massa, non finibus felis.

Ut id consequat tellus. Sed sagittis arcu in urna egestas, et gravida metus
tincidunt. Sed viverra lorem at libero vulputate tincidunt eget imperdiet
urna. Aenean semper, nunc tincidunt convallis sagittis, elit sapien vestibulum
dolor, sed maximus sapien arcu at risus. Aenean vitae diam sed nisl tincidunt
molestie. Donec sit amet pellentesque leo. Mauris semper magna vitae nunc
porttitor, ac pretium ipsum tincidunt. Nulla blandit ipsum quis magna varius,
nec rhoncus leo consequat. Integer consectetur a nisl et cursus. Suspendisse
ac dignissim cras amet.`
