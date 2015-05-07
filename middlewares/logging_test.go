package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	Start = time.Date(2015, 03, 22, 14, 47, 50, 2000, time.UTC)
	Stop  = time.Date(2015, 03, 22, 14, 47, 50, 976000, time.UTC)
)

func Test_Logging_ConstantValue(t *testing.T) {
	loggedRequest(t, StringConstant("const"), newRequest(t, "192.168.1.1", "GET", "http://server/"), "const\n")
}

func Test_Logging_Compose(t *testing.T) {
	loggedRequest(t, Compose(StringConstant("A"), StringConstant("B"), StringConstant("C")), newRequest(t, "192.168.1.1", "GET", "http://server/"), "ABC\n")
}

func Test_Logging_Enclose(t *testing.T) {
	loggedRequest(t, Enclose("(", ")", StringConstant("A")), newRequest(t, "192.168.1.1", "GET", "http://server/"), "(A)\n")
}

func Test_Logging_RemoteAddr(t *testing.T) {
	loggedRequest(t, RemoteAddr(), newRequest(t, "192.168.1.1:31954", "GET", "http://server/"), "192.168.1.1\n")
	loggedRequest(t, RemoteAddr(), newRequest(t, "192.168.1.1", "GET", "http://server/"), "192.168.1.1\n")
}

func Test_Logging_RemoteUser(t *testing.T) {
	loggedRequest(t, RemoteUser(), newRequest(t, "192.168.1.1", "GET", "http://user@server/"), "user\n")
	loggedRequest(t, RemoteUser(), newRequest(t, "192.168.1.1", "GET", "http://server/"), "-\n")
}

func Test_Logging_RequestTime(t *testing.T) {
	loggedRequest(t, RequestTime(""), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "22/Mar/2015:14:47:50 +0000\n")
	loggedRequest(t, RequestTime("2006-01-02 15:04:05.999999"), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "2015-03-22 14:47:50.000002\n")
}

func Test_Logging_RespondTime(t *testing.T) {
	loggedRequest(t, RespondTime(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "974\n")
}

func Test_Logging_RequestReferer(t *testing.T) {
	request := newRequest(t, "192.168.1.1", "GET", "http://server/")
	loggedRequest(t, RequestReferer(), request, "-\n")
	request.Header.Set("Referer", "url")
	loggedRequest(t, RequestReferer(), request, "url\n")
}

func Test_Logging_UserAgent(t *testing.T) {
	request := newRequest(t, "192.168.1.1", "GET", "http://server/")
	loggedRequest(t, RequestUserAgent(), request, "-\n")
	request.Header.Set("User-Agent", "firefox")
	loggedRequest(t, RequestUserAgent(), request, "firefox\n")
}

func Test_Logging_RequestMethod(t *testing.T) {
	loggedRequest(t, RequestMethod(), newRequest(t, "192.168.1.1", "GET", "http://server/"), "GET\n")
}

func Test_Logging_RequestURI(t *testing.T) {
	loggedRequest(t, RequestURI(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "/path\n")
}

func Test_Logging_RequestProto(t *testing.T) {
	loggedRequest(t, RequestProto(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "HTTP/1.1\n")
}

func Test_Logging_RequestInfo(t *testing.T) {
	loggedRequest(t, RequestInfo(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "GET /path HTTP/1.1\n")
}

func Test_Logging_ResponseStatus(t *testing.T) {
	loggedRequest(t, ResponseStatus(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "201\n")
}

func Test_Logging_BytesWritten(t *testing.T) {
	loggedRequest(t, BytesWritten(), newRequest(t, "192.168.1.1", "GET", "http://server/path"), "4\n")
}

func Test_Logging_ApacheCommonLog(t *testing.T) {
	loggedRequest(t,
		ApacheCommonLog(),
		newRequest(t, "192.168.1.1", "GET", "http://server/path"),
		"192.168.1.1 - - [22/Mar/2015:14:47:50 +0000] \"GET /path HTTP/1.1\" 201 4\n")
}

func Test_Logging_ApacheCombinedLog(t *testing.T) {
	request := newRequest(t, "192.168.1.1", "GET", "http://server/path")
	request.Header.Set("Referer", "url")
	request.Header.Set("User-Agent", "firefox")
	loggedRequest(t,
		ApacheCombinedLog(),
		request,
		"192.168.1.1 - - [22/Mar/2015:14:47:50 +0000] \"GET /path HTTP/1.1\" 201 4 \"url\" \"firefox\"\n")
}

func Benchmark_WriteLog(b *testing.B) {
	buffer := new(bytes.Buffer)
	request := newRequest(b, "192.168.1.1", "GET", "http://server/path")
	request.Header.Set("User-Agent", "firefox")
	recorder := &loggingResponseWriter{writer: httptest.NewRecorder()}
	recorder.WriteHeader(http.StatusOK)
	recorder.Write([]byte("body"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeLog(buffer, ApacheCombinedLog(), &Record{
			StartTime: Start,
			StopTime:  Stop,
			URL:       *request.URL,
			Request:   request,
			Writer:    recorder,
		})
	}
}

func newRequest(t testing.TB, remote string, method string, url string) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.RemoteAddr = remote
	return request
}

func loggedRequest(t *testing.T, logFunction LogFunction, request *http.Request, expectedLog string) {
	buffer := new(bytes.Buffer)
	logger := NewLogger(logFunction).SetWriter(buffer).SetTimer(fakeTimer(Start, Stop))
	handler := logger.Log(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte("body"))
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if expectedLog != buffer.String() {
		t.Errorf("Bad log: expected %#v, got %#v", expectedLog, buffer.String())
	}
}

func fakeTimer(values ...time.Time) func() time.Time {
	return func() time.Time {
		value := values[0]
		if len(values) > 1 {
			values = values[1:]
		}
		return value
	}
}
