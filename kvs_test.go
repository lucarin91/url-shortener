package main

import (
	"testing"
)

func TestKVS(t *testing.T) {

	t.Run("not found", func(t *testing.T) {
		kvs := NewMyKVS()

		_, err := kvs.Load("SjknVAAA")
		if errNotFound != err {
			t.Errorf("got %q, want %q", err, errNotFound)
		}
	})

	t.Run("store", func(t *testing.T) {
		kvs := NewMyKVS()

		err := kvs.Store("SjknVAAA", "golang.org")
		if nil != err {
			t.Errorf("got %q, want %q", err, "nil")
		}
	})

	t.Run("already exist", func(t *testing.T) {
		kvs := NewMyKVS()

		kvs.Store("SjknVAAA", "golang.org")
		err := kvs.Store("SjknVAAA", "golang.org")
		if errAlreadyExists != err {
			t.Errorf("got %q, want %q", err, errAlreadyExists)
		}
	})

	t.Run("load", func(t *testing.T) {
		kvs := NewMyKVS()

		err := kvs.Store("SjknVAAA", "golang.org")
		v, err := kvs.Load("SjknVAAA")
		if nil != err {
			t.Errorf("got %q, want %q", err, "nil")
		}
		if "golang.org" != v {
			t.Errorf("got %q, want %q", v, "golang.org")
		}
	})
}
