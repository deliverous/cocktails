package render

import (
	"testing"
)

func expect(t *testing.T, value interface{}, expexted interface{}) {
	if value != expexted {
		t.Errorf("Expected %#v, got %#v.", expexted, value)
	}
}

func refute(t *testing.T, value interface{}, expexted interface{}) {
	if value == expexted {
		t.Errorf("%#v not expected.", expexted, value)
	}
}
