package main

import (
	"fmt"
	"net/http"
)

func main() {
	hasher := NewMyHasher()
	kv := NewMyKVS()

	http.HandleFunc("/", root(kv))
	http.HandleFunc("/shorten/", shorten(kv, hasher))
	http.HandleFunc("/stats", stats())

	fmt.Println("Start...")
	http.ListenAndServe(":8080", nil)
}
