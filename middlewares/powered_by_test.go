package middlewares

import (
	"testing"
)

func Test_PoweredBy(t *testing.T) {
	tag := "cocktails"
	recorder, err := processRequest(Chain(PoweredBy(tag)).Then(nil))
	if err != nil {
		t.Fatal(err)
	}

	poweredBy := recorder.HeaderMap.Get("X-Powered-By")
	if poweredBy != tag {
		t.Errorf("PoweredBy failure: expected '%s', got %#v", tag, poweredBy)
	}
}
