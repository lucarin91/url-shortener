package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprint(r.URL)[1:]
	log.Printf("(debug) [redirect] key: '%v'\n", key)

	if len(key) != 8 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: invalid url")
		return
	}

	url, e := s.Kvs.Load(key)
	if e != nil {
		atomic.AddUint64(&s.Stats.Redirects.Failed, 1)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "error: %v", e)
		return
	}
	atomic.AddUint64(&s.Stats.Redirects.Success, 1)

	url = fmt.Sprintf("http://%v", url)
	log.Printf("(debug) [redirect] redirect to: '%v'\n", url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (s *Server) shorten(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprint(r.URL)[9:]
	key := s.Hasher.Hash(url)
	log.Printf("(debug) [shorten] url: '%v', key: '%v'\n", url, key)

	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: invalid url")
		return
	}

	e := s.Kvs.Store(key, url)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %v", e)
		return
	}
	atomic.AddUint64(&s.Stats.TotalURL, 1)

	fmt.Fprintf(w, "http://%v/%v", r.Host, key)
}

func (s *Server) statistics(w http.ResponseWriter, r *http.Request) {
	if "json" == r.URL.Query().Get("format") {
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(s.Stats.Copy())
	} else {
		fmt.Fprintf(w, "URLs: %v\n", atomic.LoadUint64(&s.Stats.TotalURL))
		fmt.Fprintf(w, "Redirect: %v\n", atomic.LoadUint64(&s.Stats.Redirects.Success))
		fmt.Fprintf(w, "Handler:\n")
		for _, v := range s.Stats.Handlers {
			fmt.Fprintf(w, "  %v: %v\n", v.Name, atomic.LoadUint64(&v.Count))
		}
	}
}
