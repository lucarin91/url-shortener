package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Start...")
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := NewServer()
	http.HandleFunc("/", s.stat(Redirect, s.redirect))
	http.HandleFunc("/shorten/", s.stat(Shorten, s.shorten))
	http.HandleFunc("/statistics", s.stat(Statistics, s.statistics))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("Http server: %w", err)
	}
	return nil
}
