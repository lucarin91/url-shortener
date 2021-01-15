package main

import "sync/atomic"

type Stats struct {
	ServerStats `json:"server_stats"`
}

func (s *Stats) Copy() *Stats {
	newS := &Stats{
		ServerStats{
			atomic.LoadUint64(&s.TotalURL),
			Redirects{
				atomic.LoadUint64(&s.Redirects.Success),
				atomic.LoadUint64(&s.Redirects.Failed),
			},
			make([]Handler, len(s.Handlers)),
		},
	}

	for i, v := range s.Handlers {
		newS.Handlers[i] = Handler{v.Name, atomic.LoadUint64(&v.Count)}
	}
	return newS
}

type ServerStats struct {
	TotalURL  uint64    `json:"total_url"`
	Redirects Redirects `json:"redirects"`
	Handlers  []Handler `json:"handlers"`
}

type Redirects struct {
	Success uint64 `json:"success"`
	Failed  uint64 `json:"failed"`
}

type Handler struct {
	Name  string `json:"name"`
	Count uint64 `json:"count"`
}

type URLs int

const (
	Redirect   URLs = iota // 0
	Shorten                // 1
	Statistics             // 2
)
