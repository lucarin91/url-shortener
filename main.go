package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

	stg, err := LoadStorageFile(args.file)
	var e *os.PathError
	switch {
	case err == nil:
		// OK, load storage file on the server
		err = s.SetStorage(stg)
		if err != nil {
			return err
		}
		fmt.Printf("[INFO] Load %d from storage file %q\n", s.Stats.TotalURL, args.file)
	case errors.As(err, &e) && os.IsNotExist(e):
		// File not exist, ignored
		fmt.Printf("[WARN] Storage file %q not exist\n", args.file)
	default:
		// Other error, i.e., permission
		return err
	}

	signalHandler(func() {
		err := s.Storage().SaveStorageFile(args.file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot save storage: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Save storage")
		os.Exit(0)
	})

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

func signalHandler(cb func()) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, os.Interrupt, os.Kill)
	go func() {
		<-termChan // Blocks here until interrupted
		cb()
	}()
}
