package main

import (
	"reflect"
	"testing"
)

func TestKVSLoad(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
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

	t.Run("not found", func(t *testing.T) {
		kvs := NewMyKVS()

		_, err := kvs.Load("SjknVAAA")
		if errNotFound != err {
			t.Errorf("got %q, want %q", err, errNotFound)
		}
	})
}

func TestKVSStore(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
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
}

func TestKVSDump(t *testing.T) {
	kvs := NewMyKVS()
	kvs.storage["key1"] = "url1"
	kvs.storage["key2"] = "url2"
	kvs.storage["key3"] = "url3"
	kvs.storage["key4"] = "url4"
	kvs.storage["key5"] = "url5"

	stg := kvs.Dump()
	check := make(map[string]string)
	for _, v := range stg.URLPairs {
		check[v.Short] = v.Long
	}

	if !reflect.DeepEqual(kvs.storage, check) {
		t.Errorf("got %q, want %q", kvs.storage, check)
	}
}
