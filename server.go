package main

type Server struct {
	Hasher MyHasher
	Kvs    MyKVS
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
