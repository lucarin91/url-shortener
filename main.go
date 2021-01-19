package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type cmdArgs struct {
	addr string
	file string
}

func main() {
	// Read command arguments
	var args cmdArgs
	flag.StringVar(&args.addr, "addr", ":8080", "server address")
	flag.StringVar(&args.file, "load", "storage.json", "path to storage file")
	flag.Parse()

	if err := run(&args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args *cmdArgs) error {
	s := NewServer()

	// Load storage from file
	stg, err := LoadStorageFile(args.file)
	if err != nil {
		return err
	}
	// Set storage on the server
	err = s.SetStorage(stg)
	if err != nil {
		return err
	}
	fmt.Printf("[INFO] Load %d from storage file %q\n", s.Stats.TotalURL, args.file)

	// Configure handler
	http.HandleFunc("/", s.stat(Redirect, s.redirect))
	http.HandleFunc("/shorten/", s.stat(Shorten, s.shorten))
	http.HandleFunc("/statistics", s.stat(Statistics, s.statistics))

	fmt.Printf("Start on %q\n", args.addr)
	err = http.ListenAndServe(args.addr, nil)
	if err != nil {
		return fmt.Errorf("Http server: %v", err)
	}
	return nil
}
