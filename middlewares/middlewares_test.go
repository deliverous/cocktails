package middlewares

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Chain_NoMiddleware(t *testing.T) {
	ensureRequestContains(t,
		Chain().Then(taggerHandler("[APP]")),
		"[APP]")
}

func Test_Chain_Then_WithNil_ShouldUseVoidHandler(t *testing.T) {
	ensureRequestContains(t,
		Chain(taggerMiddleware("[M1]")).Then(nil),
		"[M1]")
}

func Test_Chain_Then_WithHandler_ShouldCallTheHandler(t *testing.T) {
	ensureRequestContains(t,
		Chain(taggerMiddleware("[M1]")).Then(taggerHandler("[APP]")),
		"[M1][APP]")
}

func Test_Chain_Concat_ShouldConcatTheTwoChains(t *testing.T) {
	full := Chain(taggerMiddleware("[M1]")).Concat(Chain(taggerMiddleware("[M2]"), taggerMiddleware("[M3]")))
	ensureRequestContains(t, full.Then(taggerHandler("[APP]")), "[M1][M2][M3][APP]")
}

func Test_Chain_Concat_ShouldNotModifyTheReceiver(t *testing.T) {
	first := Chain(taggerMiddleware("[M1]"))
	first.Concat(Chain(taggerMiddleware("[M2]"), taggerMiddleware("[M3]")))
	ensureRequestContains(t, first.Then(taggerHandler("[APP]")), "[M1][APP]")
}

func Test_Chain_Concat_ShouldNotModifyTheArgument(t *testing.T) {
	second := Chain(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	Chain(taggerMiddleware("[M1]")).Concat(second)
	ensureRequestContains(t, second.Then(taggerHandler("[APP]")), "[M2][M3][APP]")
}

func Test_Chain_Extend_ShouldExtendsTheChainWithGivenMiddleware(t *testing.T) {
	full := Chain(taggerMiddleware("[M1]")).Extend(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	ensureRequestContains(t, full.Then(taggerHandler("[APP]")), "[M1][M2][M3][APP]")
}

func Test_Chain_Extend_ShouldNotModifyTheReceiver(t *testing.T) {
	first := Chain(taggerMiddleware("[M1]"))
	first.Extend(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	ensureRequestContains(t, first.Then(taggerHandler("[APP]")), "[M1][APP]")
}

func Test_Chain_ExtendShouldBeImmutable(t *testing.T) {
	first := Chain(taggerMiddleware("[M1]"))
	full := first.Extend(taggerMiddleware("[M2]"), taggerMiddleware("[M3]"))
	if repr(&first.builders[0]) == repr(&full.builders[0]) {
		t.Errorf("Extand failed to be immutable")
	}
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

func repr(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}
