package main

import (
	"testing"
)

func TestKVS(t *testing.T) {

	t.Run("not found", func(t *testing.T) {
		kvs := NewMyKVS()

		_, err := kvs.Load("SjknVAAA")
		if nil == err {
			t.Fail()
		}
		if "not found" != err.Error() {
			t.Fail()
		}
	})

	t.Run("store", func(t *testing.T) {
		kvs := NewMyKVS()

		err := kvs.Store("SjknVAAA", "golang.org")
		if nil != err {
			t.Fail()
		}
	})

	t.Run("already exist", func(t *testing.T) {
		kvs := NewMyKVS()

		kvs.Store("SjknVAAA", "golang.org")
		err := kvs.Store("SjknVAAA", "golang.org")
		if nil == err {
			t.Fail()
		}
		if "already exist" != err.Error() {
			t.Fail()
		}
	})

	t.Run("load", func(t *testing.T) {
		kvs := NewMyKVS()

		err := kvs.Store("SjknVAAA", "golang.org")
		v, err := kvs.Load("SjknVAAA")
		if nil != err {
			t.Fail()
		}
		if "golang.org" != v {
			t.Fail()
		}
	})
}
