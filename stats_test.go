package main

import (
	"reflect"
	"testing"
)

func TestStatsCopy(t *testing.T) {
	s := &Stats{
		ServerStats{
			TotalURL: 1,
			Redirects: Redirects{
				Success: 2,
				Failed:  3,
			},
			Handlers: []Handler{
				Handler{Name: "qwe", Count: 4},
				Handler{Name: "qwee", Count: 5},
				Handler{Name: "qweee", Count: 6},
			},
		},
	}
	newS := s.Copy()

	if !reflect.DeepEqual(s, newS) {
		t.Errorf("got %+v, want %+v", newS, s)
	}
}
