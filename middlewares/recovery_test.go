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
		processRequest(t, Chain(DefaultRecovery().Recover).Then(panicHandler))
	}()
	if didPanic {
		t.Error("Panic was propagated")
	}
}

func Test_WithRecovery_ShouldPrintTheStackToTheLogger(t *testing.T) {
	buffer := bytes.NewBufferString("")
	recovery := DefaultRecovery()
	recovery.Logger = log.New(buffer, "", 0)

	processRequest(t, Chain(recovery.Recover).Then(panicHandler))
	if !strings.Contains(buffer.String(), "here is a panic!") {
		t.Error("Stack was not printed into the logger")
	}
}

func Test_WithRecovery_WithPrintStack_ShouldPrintTheStackToTheResponseBody(t *testing.T) {
	recovery := DefaultRecovery()
	recovery.Logger = log.New(bytes.NewBufferString(""), "", 0)
	recovery.PrintStack = true
	recorder := processRequest(t, Chain(recovery.Recover).Then(panicHandler))
	if !strings.Contains(recorder.Body.String(), "here is a panic!") {
		t.Error("Stack was not printed into the response")
	}
}
