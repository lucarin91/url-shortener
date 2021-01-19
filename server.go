package main

import (
	"fmt"
	"os"
	"os/signal"
)

type Server struct {
	Hasher Hasher
	Kvs    KVS
	Stats  Stats
}

func NewServer() *Server {
	s := &Server{
		NewMyHasher(),
		NewMyKVS(),
		Stats{
			ServerStats{Handlers: make([]Handler, 3)},
		},
	}
	s.Stats.Handlers[Redirect] = Handler{Name: "/"}
	s.Stats.Handlers[Shorten] = Handler{Name: "/shorten"}
	s.Stats.Handlers[Statistics] = Handler{Name: "/statistics"}
	return s
}

func (s *Server) SetStorage(stg *Storage) error {
	for _, v := range stg.URLPairs {
		err := s.Kvs.Store(v.Short, v.Long)
		if err != nil {
			return fmt.Errorf("cannot load url: %v", err)
		}
	}
	s.Stats.TotalURL = uint64(len(stg.URLPairs))
	return nil
}

func (s *Server) Storage() *Storage {
	return s.Kvs.Dump()
}

func (s *Server) ShutdownHandler(args *cmdArgs) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, os.Interrupt, os.Kill)
	go func() {
		<-termChan // Blocks here until interrupted
		err := s.Storage().SaveStorageFile(args.file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot save storage: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Save storage")
		os.Exit(0)
	}()
}
