package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

func redirect(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateStats(state, Redirect)

		key := fmt.Sprint(r.URL)[1:]
		fmt.Printf("   (debug) key: '%v'\n", key)

		if len(key) != 8 {
			w.Write([]byte("error: invalid url"))
			return
		}

		url, e := state.Load(key)
		if e != nil {
			w.Write([]byte(fmt.Sprintf("error: %v", e)))
			atomic.AddUint64(&state.Stats.Redirects.Failed, 1)
			return
		}
		url = fmt.Sprintf("http://%v", url)
		fmt.Printf("   (debug) redirect to: '%v'\n", url)
		atomic.AddUint64(&state.Stats.Redirects.Success, 1)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}

func shorten(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateStats(state, Shorten)

		url := fmt.Sprint(r.URL)[9:]
		key := state.Hash(url)
		fmt.Printf("   (debug) url: '%v', key: '%v'\n", url, key)

		if len(url) == 0 {
			w.Write([]byte("error: invalid url"))
			return
		}

		e := state.Store(key, url)
		if e != nil {
			w.Write([]byte(fmt.Sprintf("error: %v", e)))
			return
		}

		sURL := fmt.Sprintf("http://%v/%v", r.Host, key)
		w.Write([]byte(sURL))
	}
}

func statistics(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateStats(state, Statistics)

		var b strings.Builder
		for _, v := range state.Handlers {
			fmt.Fprintf(&b, "  %v: %v\n", v.name, atomic.LoadUint64(&v.count))
		}
		hadlers := b.String()
		urls := len(state.MyKVS.storage)
		redirects := atomic.LoadUint64(&state.Stats.Redirects.Success)

		s := fmt.Sprintf("URLs: %v\nRedirect: %v\nHandler:\n%v", urls, redirects, hadlers)
		w.Write([]byte(s))
	}
}

func updateStats(state *State, i URLs) {
	atomic.AddUint64(&state.Stats.ServerStats.TotalURL, 1)
	atomic.AddUint64(&state.Stats.Handlers[i].count, 1)
}
