package main

import (
	"errors"
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
		fmt.Fprintf(os.Stderr, "[ERR] %s\n", err)
		os.Exit(1)
	}
}

func run(args *cmdArgs) error {
	s := NewServer()

	// Load storage from file
	stg, err := LoadStorageFile(args.file)
	var e *os.PathError
	switch {
	// OK, load storage file on the server
	case err == nil:
		err = s.SetStorage(stg)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] Load %d from storage file %q\n", s.Stats.TotalURL, args.file)
	// File not exist, ignored
	case errors.As(err, &e) && os.IsNotExist(e):
		fmt.Printf("[WARN] Storage file %q not exist\n", args.file)
	// Other error, i.e., permission
	default:
		return err
	}

	// Handle graceful shoutdown
	s.ShutdownHandler(args)

	// Configure http handlers
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
