package main

import (
	"testing"
)

var flagtests = []struct {
	in  string
	out string
}{
	{"about.sourcegraph.com/go/gophercon-2019-how-i-write-http-web-services-after-eight-years", "SjknVAAA"},
	{"golang.org/doc/effective_go.html", "UAOuaAAA"},
	{"golang.org", "R7W7LQAA"},
}

func TestHash(t *testing.T) {
	h := NewMyHasher()
	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {
			v := h.Hash(tt.in)
			if v != tt.out {
				t.Errorf("got %q, want %q", v, tt.out)
			}
		})
	}
}
