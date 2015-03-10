package middlewares

import (
	"net/http"
	"testing"
)

func Test_Chain_NoMiddleware(t *testing.T) {
	ensureRequestContains(t,
		Chain().Then(taggerHandler("[APP]")),
		"[APP]")
}

func Test_Chain_ThenNil(t *testing.T) {
	ensureRequestContains(t,
		Chain(taggerMiddleware("[M1]")).Then(nil),
		"[M1]")
}

func Test_Chain_Then(t *testing.T) {
	ensureRequestContains(t,
		Chain(taggerMiddleware("[M1]")).Then(taggerHandler("[APP]")),
		"[M1][APP]")
}

func Test_Chain_Concat(t *testing.T) {
	first := Chain(taggerMiddleware("[M1]"))
	second := Chain(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	full := first.Concat(second)
	ensureRequestContains(t, first.Then(taggerHandler("[APP]")), "[M1][APP]")
	ensureRequestContains(t, second.Then(taggerHandler("[APP]")), "[M2][M3][APP]")
	ensureRequestContains(t, full.Then(taggerHandler("[APP]")), "[M1][M2][M3][APP]")
}

func Test_Chain_Extend(t *testing.T) {
	first := Chain(taggerMiddleware("[M1]"))
	full := first.Extend(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	ensureRequestContains(t, first.Then(taggerHandler("[APP]")), "[M1][APP]")
	ensureRequestContains(t, full.Then(taggerHandler("[APP]")), "[M1][M2][M3][APP]")
}

func ensureRequestContains(t *testing.T, handler http.Handler, expectedTrace string) {
	recorder, err := processRequest(handler)
	if err != nil {
		t.Fatal(err)
	}
	body := recorder.Body.String()
	if body != expectedTrace {
		t.Errorf("expected %#v, got %#v", expectedTrace, body)
	}
}

func taggerHandler(tag string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(tag))
	})
}

func taggerMiddleware(tag string) MiddlewareBuilder {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte(tag))
			next.ServeHTTP(writer, request)
		})
	}
}
