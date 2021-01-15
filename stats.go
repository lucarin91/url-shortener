package main

type Stats struct {
	ServerStats
}

type ServerStats struct {
	TotalURL  uint64
	Redirects Redirects
	Handlers  []Handler
}

type Redirects struct {
	Success uint64
	Failed  uint64
}

type Handler struct {
	Name  string
	Count uint64
}

type URLs int

const (
	Redirect   URLs = iota // 0
	Shorten                // 1
	Statistics             // 2
)
