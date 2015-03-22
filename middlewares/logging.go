package middlewares

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type LoggingResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
}

type loggingResponseWriter struct {
	writer http.ResponseWriter
	status int
	size   int
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.writer.Header()
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	if l.status == 0 {
		// Sets the status to StatusOK if status was not set previously
		l.WriteHeader(http.StatusOK)
	}
	size, err := l.writer.Write(b)
	l.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(s int) {
	l.writer.WriteHeader(s)
	l.status = s
}

func (l *loggingResponseWriter) Status() int {
	return l.status
}

func (l *loggingResponseWriter) Size() int {
	return l.size
}

type loggingHijacker struct {
	loggingResponseWriter
}

func (l *loggingHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := l.loggingResponseWriter.writer.(http.Hijacker)
	conn, rw, err := h.Hijack()
	if err == nil && l.loggingResponseWriter.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and WriteHeader has not been called yet
		l.loggingResponseWriter.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}

type LogFunction func([]byte, *Record) []byte

type Logger struct {
	Timer       func() time.Time
	Writer      io.Writer
	LogFunction LogFunction
}

func NewLogger(function LogFunction) *Logger {
	return &Logger{
		Timer:       time.Now,
		Writer:      os.Stdout,
		LogFunction: function,
	}
}

func (logger *Logger) SetWriter(writer io.Writer) *Logger {
	logger.Writer = writer
	return logger
}

func (logger *Logger) SetTimer(function func() time.Time) *Logger {
	logger.Timer = function
	return logger
}

type Record struct {
	StartTime time.Time
	StopTime  time.Time
	Duration  time.Duration
	URL       url.URL
	Request   *http.Request
	Writer    LoggingResponseWriter
}

func (logger *Logger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := logger.Timer()
		loggingWriter := loggingResponseWriter{writer: writer}
		var subWriter LoggingResponseWriter
		if _, ok := writer.(http.Hijacker); ok {
			subWriter = &loggingHijacker{loggingResponseWriter: loggingWriter}
		} else {
			subWriter = &loggingWriter
		}
		url := *request.URL
		next.ServeHTTP(subWriter, request)
		stop := logger.Timer()
		writeLog(logger.Writer, logger.LogFunction, &Record{
			StartTime: start,
			StopTime:  stop,
			URL:       url,
			Request:   request,
			Writer:    &loggingWriter,
		})
	})
}

func writeLog(writer io.Writer, function LogFunction, record *Record) {
	buffer := make([]byte, 0, 1024)
	buffer = function(buffer, record)
	buffer = append(buffer, '\n')
	writer.Write(buffer)
}

func Compose(functions ...LogFunction) LogFunction {
	return func(buffer []byte, record *Record) []byte {
		for _, function := range functions {
			buffer = function(buffer, record)
		}
		return buffer
	}
}

func Enclose(begin string, end string, function LogFunction) LogFunction {
	return func(buffer []byte, record *Record) []byte {
		buffer = append(buffer, begin...)
		buffer = function(buffer, record)
		return append(buffer, end...)
	}
}

func Constant(value string) LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, value...)
	}
}

func RemoteAddr() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		host, _, err := net.SplitHostPort(record.Request.RemoteAddr)
		if err != nil {
			host = record.Request.RemoteAddr
		}
		return append(buffer, host...)
	}
}

func RemoteUser() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		username := "-"
		if record.URL.User != nil {
			if name := record.URL.User.Username(); name != "" {
				username = name
			}
		}
		return append(buffer, username...)
	}
}

func RequestTime(format string) LogFunction {
	if format == "" {
		format = "02/Jan/2006:15:04:05 -0700"
	}
	return func(buffer []byte, record *Record) []byte {
		record.StartTime.Format(format)
		return append(buffer, record.StartTime.Format(format)...)
	}
}

func RespondTime() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		if record.Duration == 0 {
			record.Duration = record.StopTime.Sub(record.StartTime)
		}
		return append(buffer, fmt.Sprintf("%d", record.Duration/time.Microsecond)...)
	}
}

func RequestReferer() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		referer := record.Request.Referer()
		if referer == "" {
			referer = "-"
		}
		return append(buffer, referer...)
	}
}

func RequestUserAgent() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		agent := record.Request.UserAgent()
		if agent == "" {
			agent = "-"
		}
		return append(buffer, agent...)
	}
}

func RequestMethod() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, record.Request.Method...)
	}
}

func RequestURI() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, record.URL.RequestURI()...)
	}
}

func RequestProto() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, record.Request.Proto...)
	}
}

func RequestInfo() LogFunction {
	return Compose(RequestMethod(), Constant(" "), RequestURI(), Constant(" "), RequestProto())
}

func ResponseStatus() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, strconv.Itoa(record.Writer.Status())...)
	}
}

func BytesWritten() LogFunction {
	return func(buffer []byte, record *Record) []byte {
		return append(buffer, strconv.Itoa(record.Writer.Size())...)
	}
}

func ApacheCommonLog() LogFunction {
	return Compose(
		RemoteAddr(),
		Constant(" - "),
		RemoteUser(),
		Enclose(" [", "] ", RequestTime("")),
		Enclose("\"", "\" ", RequestInfo()),
		ResponseStatus(),
		Constant(" "),
		BytesWritten(),
	)
}

func ApacheCombinedLog() LogFunction {
	return Compose(
		ApacheCommonLog(),
		Enclose(" \"", "\" ", RequestReferer()),
		Enclose("\"", "\"", RequestUserAgent()),
	)
}
