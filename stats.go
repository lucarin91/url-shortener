package main

type Stats struct {
	ServerStats
	Handlers []Handler
}

type ServerStats struct {
	TotalURL uint64
	Redirects
}

type Redirects struct {
	Success uint64
	Failed  uint64
}

type Handler struct {
	name  string
	count uint64
}

type URLs int

const (
	Redirect   URLs = iota // 0
	Shorten                // 1
	Statistics             // 2
)
