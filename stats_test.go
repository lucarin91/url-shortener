package main

import (
	"testing"
)

func TestStatsCopy(t *testing.T) {
	s := Stats{
		ServerStats{
			1,
			Redirects{
				2,
				3,
			},
			[]Handler{Handler{"qwe", 4}, Handler{"qwee", 5}, Handler{"qweee", 6}},
		},
	}
	newS := s.Copy()

	if s.TotalURL != newS.TotalURL ||
		s.Redirects.Failed != newS.Redirects.Failed ||
		s.Redirects.Success != newS.Redirects.Success {
		t.Fail()
	}

	for i, v := range s.Handlers {
		if newS.Handlers[i].Name != v.Name ||
			newS.Handlers[i].Count != v.Count {
			t.Fail()
		}
	}
}
