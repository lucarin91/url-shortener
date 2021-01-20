package main

import (
	"fmt"
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
