package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Hasher MyHasher
	Kvs    MyKVS
	Stats  Stats
}

func main() {
	server := &Server{
		NewMyHasher(),
		NewMyKVS(),
		Stats{
			ServerStats{Handlers: make([]Handler, 3)},
		},
	}
	server.Stats.Handlers[Redirect] = Handler{"/", 0}
	server.Stats.Handlers[Shorten] = Handler{"/shorten", 0}
	server.Stats.Handlers[Statistics] = Handler{"/statistics", 0}

	http.HandleFunc("/", server.redirect)
	http.HandleFunc("/shorten/", server.shorten)
	http.HandleFunc("/statistics", server.statistics)

	fmt.Println("Start...")
	http.ListenAndServe(":8080", nil)
}
