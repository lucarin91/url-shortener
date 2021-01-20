package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const urlsPath = "urls.json"
const shortlyPath = "./shortly/shortly"
const shortlyURL = "http://localhost:8080"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n  Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	stg := &Storage{
		URLPairs: []URLPair{URLPair{Short: "R7W7LQAA", Long: "golang.org"}},
	}
	if err := stg.SaveStorageFile(urlsPath); err != nil {
		return fmt.Errorf("cannot init storage: %v", err)
	}

	fmt.Println("\n  Start server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, shortlyPath, "--load", urlsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	fmt.Println("  OK!")

	// Test server
	if err := CheckWrongRedirect("/H7DFLQAA"); err != nil {
		return err
	}
	if err := CheckGoodRedirect("/R7W7LQAA", "golang.org/"); err != nil {
		return err
	}
	url, err := CheckShortURL("/duckduckgo.com")
	if err != nil {
		return err
	}
	if err := CheckGoodRedirect(url, "duckduckgo.com/"); err != nil {
		return err
	}
	statsCheck := &Stats{ServerStats: ServerStats{
		TotalURL:  2,
		Redirects: Redirects{Success: 2, Failed: 1},
		Handlers: []Handler{
			Handler{Name: "/", Count: 3},
			Handler{Name: "/shorten", Count: 1},
			Handler{Name: "/statistics", Count: 1},
		},
	}}
	if err := CheckSatistics(statsCheck); err != nil {
		return err
	}

	fmt.Println("\n  Interrupt server")
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for child process: %v", err)
	}
	fmt.Println("  OK!")

	var finalStorage = &Storage{URLPairs: []URLPair{
		URLPair{Short: "R7W7LQAA", Long: "golang.org"},
		URLPair{Short: "tpxMrwAA", Long: "duckduckgo.com"},
	}}
	if err := CheckStorage(finalStorage); err != nil {
		return err
	}

	return nil
}
