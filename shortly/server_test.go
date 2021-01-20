package main

import (
	"testing"
)

func TestSetStorage(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		stg := &Storage{
			URLPairs: []URLPair{
				URLPair{Short: "short1", Long: "long1"},
				URLPair{Short: "short2", Long: "long2"},
				URLPair{Short: "short3", Long: "long3"},
				URLPair{Short: "short4", Long: "long4"},
			},
		}
		s := NewServer()
		e := s.SetStorage(stg)
		if e != nil {
			t.Errorf("got %q, want %q", e, "nil")
		}
		for _, v := range stg.URLPairs {
			_, e := s.Kvs.Load(v.Short)
			if e != nil {
				t.Errorf("key %q not found", v.Short)
			}
		}
		if 4 != s.Stats.TotalURL || len(stg.URLPairs) != int(s.Stats.TotalURL) {
			t.Errorf("got %v, want %v", len(stg.URLPairs), int(s.Stats.TotalURL))
		}
	})

	t.Run("already exists", func(t *testing.T) {
		stg := &Storage{
			URLPairs: []URLPair{
				URLPair{Short: "short", Long: "long1"},
				URLPair{Short: "short", Long: "long2"},
			},
		}
		s := NewServer()
		e := s.SetStorage(stg)
		if e == nil {
			t.Errorf("got %v, want %q", e, "error")
		}
	})
}
