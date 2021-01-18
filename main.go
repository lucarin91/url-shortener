package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	addr string
	file string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "server address")
	flag.StringVar(&file, "load", "storage.json", "path to storage file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := NewServer()

	_, err := os.Stat(file)
	if err == nil {
		stg, err := LoadStorageFile(file)
		if err != nil {
			return err
		}
		s.LoadStorage(stg)
		fmt.Printf("[INFO] Load %d from storage file %q\n", s.Stats.TotalURL, file)
	} else {
		fmt.Printf("[WARN] Storage file %q not found\n", file)
	}

	http.HandleFunc("/", s.stat(Redirect, s.redirect))
	http.HandleFunc("/shorten/", s.stat(Shorten, s.shorten))
	http.HandleFunc("/statistics", s.stat(Statistics, s.statistics))

	fmt.Printf("Start on %q\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("Http server: %w", err)
	}
	return nil
}
