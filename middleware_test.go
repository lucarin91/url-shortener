package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStat(t *testing.T) {
	s := NewServer()
	s.Stats.Handlers[Redirect].Count = 1
	s.Stats.Handlers[Shorten].Count = 2
	s.Stats.Handlers[Statistics].Count = 3

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.stat(Redirect, func(w http.ResponseWriter, r *http.Request) {})(w, req)
	s.stat(Shorten, func(w http.ResponseWriter, r *http.Request) {})(w, req)
	s.stat(Statistics, func(w http.ResponseWriter, r *http.Request) {})(w, req)

	if 2 != s.Stats.Handlers[Redirect].Count {
		t.Errorf("got %v, want %v", s.Stats.Handlers[Redirect].Count, 2)
	}

	if 3 != s.Stats.Handlers[Shorten].Count {
		t.Errorf("got %v, want %v", s.Stats.Handlers[Shorten].Count, 3)
	}

	if 4 != s.Stats.Handlers[Statistics].Count {
		t.Errorf("got %v, want %v", s.Stats.Handlers[Statistics].Count, 4)
	}
}
