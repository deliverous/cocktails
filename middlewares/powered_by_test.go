package middlewares

import (
	"testing"
)

func Test_PoweredBy(t *testing.T) {
	tag := "cocktails"
	poweredBy := processRequest(t, Chain(PoweredBy(tag)).Then(nil)).HeaderMap.Get("X-Powered-By")
	if poweredBy != tag {
		t.Errorf("PoweredBy failure: expected '%s', got %#v", tag, poweredBy)
	}
}
