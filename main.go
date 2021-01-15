package main

import (
	"fmt"
	"net/http"
)

type State struct {
	MyHasher
	MyKVS
	Stats
}

func main() {
	state := &State{NewMyHasher(),
		NewMyKVS(),
		Stats{
			ServerStats{},
			make([]Handler, 3),
		},
	}
	state.Stats.Handlers[Redirect] = Handler{"/", 0}
	state.Stats.Handlers[Shorten] = Handler{"/shorten", 0}
	state.Stats.Handlers[Statistics] = Handler{"/statistics", 0}

	http.HandleFunc("/", redirect(state))
	http.HandleFunc("/shorten/", shorten(state))
	http.HandleFunc("/statistics", statistics(state))

	fmt.Println("Start...")
	http.ListenAndServe(":8080", nil)
}
