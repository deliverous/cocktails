package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"testing"
)

var panicHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	panic("here is a panic!")
})

func testRecovery() *Recovery {
	return testRecoveryLoggingInto(bytes.NewBufferString(""))
}

func testRecoveryLoggingInto(buffer *bytes.Buffer) *Recovery {
	return NewRecovery(Logger(log.New(buffer, "", 0)))
}

func Test_WithoutRecovery_ShouldPanic(t *testing.T) {
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		processRequest(t, Chain().Then(panicHandler))
	}()
	if !didPanic {
		t.Error("Panic was not propagated")
	}
}

func Test_WithRecovery_ShouldNotPanic(t *testing.T) {
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		processRequest(t, Chain(testRecovery().Recover).Then(panicHandler))
	}()
	if didPanic {
		t.Error("Panic was propagated")
	}
}

func Test_WithRecovery_ShouldRespondsInternalServerError(t *testing.T) {
	recorder := processRequest(t, Chain(testRecovery().Recover).Then(panicHandler))
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Recovery failed to returns internal server error")
	}
}

func Test_WithRecovery_ShouldPrintTheStackToTheLogger(t *testing.T) {
	buffer := bytes.NewBufferString("")
	recovery := testRecoveryLoggingInto(buffer)
	recovery.PrintStack = false
	processRequest(t, Chain(recovery.Recover).Then(panicHandler))
	if !strings.Contains(buffer.String(), "here is a panic!") {
		t.Error("Stack was not printed into the logger")
	}
}

func Test_WithRecovery_WithPrintStack_ShouldPrintTheStackToTheResponseBody(t *testing.T) {
	recovery := testRecovery()
	recovery.PrintStack = true
	recorder := processRequest(t, Chain(recovery.Recover).Then(panicHandler))
	if !strings.Contains(recorder.Body.String(), "here is a panic!") {
		t.Error("Stack was not printed into the response")
	}
}
