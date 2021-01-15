package main

import (
	"net/http"
	"sync/atomic"
)

func (s *Server) stat(i URLs, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&s.Stats.Handlers[i].Count, 1)
		h(w, r)
	}
}
