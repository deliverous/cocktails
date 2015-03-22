package middlewares

import (
	"bufio"
	"bytes"
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

type LogFunction func(*bytes.Buffer, *Record)

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

func (record *Record) RemoteAddr() string {
	host, _, err := net.SplitHostPort(record.Request.RemoteAddr)
	if err != nil {
		host = record.Request.RemoteAddr
	}
	return host
}

func (record *Record) RemoteUser() string {
	username := "-"
	if record.URL.User != nil {
		if name := record.URL.User.Username(); name != "" {
			username = name
		}
	}
	return username
}

func (record *Record) RequestTime(format string) string {
	return record.StartTime.Format(format)
}

func (record *Record) RespondTime() string {
	if record.Duration == 0 {
		record.Duration = record.StopTime.Sub(record.StartTime)
	}
	return strconv.FormatInt(record.Duration.Nanoseconds()/time.Microsecond.Nanoseconds(), 10)
}

func (record *Record) RequestReferer() string {
	referer := record.Request.Referer()
	if referer == "" {
		referer = "-"
	}
	return referer
}

func (record *Record) RequestUserAgent() string {
	agent := record.Request.UserAgent()
	if agent == "" {
		agent = "-"
	}
	return agent
}

func (record *Record) Status() string {
	return strconv.Itoa(record.Writer.Status())
}

func (record *Record) BytesWritten() string {
	return strconv.Itoa(record.Writer.Size())
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
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	function(buffer, record)
	buffer.WriteByte('\n')
	writer.Write(buffer.Bytes())
}

func Compose(functions ...LogFunction) LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		for _, function := range functions {
			function(buffer, record)
		}
	}
}

func Enclose(begin string, end string, function LogFunction) LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(begin)
		function(buffer, record)
		buffer.WriteString(end)
	}
}

func StringConstant(value string) LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(value)
	}
}

func ByteConstant(value byte) LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteByte(value)
	}
}

func RemoteAddr() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RemoteAddr())
	}
}

func RemoteUser() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RemoteUser())
	}
}

func RequestTime(format string) LogFunction {
	if format == "" {
		format = "02/Jan/2006:15:04:05 -0700"
	}
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RequestTime(format))
	}
}

func RespondTime() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RespondTime())
	}
}

func RequestReferer() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RequestReferer())
	}
}

func RequestUserAgent() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.RequestUserAgent())
	}
}

func RequestMethod() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.Request.Method)
	}
}

func RequestURI() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.URL.RequestURI())
	}
}

func RequestProto() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.Request.Proto)
	}
}

func RequestInfo() LogFunction {
	return Compose(
		RequestMethod(),
		ByteConstant(' '),
		RequestURI(),
		ByteConstant(' '),
		RequestProto(),
	)
}

func ResponseStatus() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.Status())
	}
}

func BytesWritten() LogFunction {
	return func(buffer *bytes.Buffer, record *Record) {
		buffer.WriteString(record.BytesWritten())
	}
}

func ApacheCommonLog() LogFunction {
	return Compose(
		RemoteAddr(),
		StringConstant(" - "),
		RemoteUser(),
		Enclose(" [", "] ", RequestTime("")),
		Enclose("\"", "\" ", RequestInfo()),
		ResponseStatus(),
		StringConstant(" "),
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
