package main

import (
	"testing"
)

func TestHash(t *testing.T) {
	h := NewMyHasher()

	v := h.Hash("about.sourcegraph.com/go/gophercon-2019-how-i-write-http-web-services-after-eight-years")
	if "SjknVAAA" != v {
		t.Fail()
	}

	v = h.Hash("golang.org/doc/effective_go.html")
	if "UAOuaAAA" != v {
		t.Fail()
	}

	v = h.Hash("golang.org")
	if "R7W7LQAA" != v {
		t.Fail()
	}
}
