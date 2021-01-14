package main

import (
	"fmt"
	"net/http"
)

type handleFun func(http.ResponseWriter, *http.Request)

func root(kv KVS) handleFun {
	return func(w http.ResponseWriter, r *http.Request) {
		key := fmt.Sprint(r.URL)[1:]
		fmt.Printf("   (debug) key: '%v'\n", key)

		if len(key) != 8 {
			w.Write([]byte("error: invalid url"))
			return
		}

		url, e := kv.Load(key)
		if e != nil {
			w.Write([]byte(fmt.Sprintf("error: %v", e)))
			return
		}
		url = fmt.Sprintf("http://%v", url)
		fmt.Printf("   (debug) redirect to: '%v'\n", url)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}

func shorten(kv KVS, hasher Hasher) handleFun {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprint(r.URL)[9:]
		key := hasher.Hash(url)
		fmt.Printf("   (debug) url: '%v', key: '%v'\n", url, key)

		if len(url) == 0 {
			w.Write([]byte("error: invalid url"))
			return
		}

		e := kv.Store(key, url)
		if e != nil {
			w.Write([]byte(fmt.Sprintf("error: %v", e)))
			return
		}

		sURL := fmt.Sprintf("http://%v/%v", r.Host, key)
		w.Write([]byte(sURL))
	}
}

func stats() handleFun {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TBD"))
	}
}
