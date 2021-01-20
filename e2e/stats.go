package main

type Stats struct {
	ServerStats `json:"server_stats"`
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
