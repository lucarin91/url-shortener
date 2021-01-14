package main

import (
	"fmt"
	"sync"
)

type KVS interface {
	// Store stores a key value pair.
	Store(k, v string) error

	// Load returns the value associated with a given key.
	Load(k string) (string, error)
}

type MyKVS struct {
	storage map[string]string
	sync.RWMutex
}

func NewMyKVS() MyKVS {
	return MyKVS{storage: make(map[string]string)}
}

func (kvs MyKVS) Store(k, v string) error {
	kvs.Lock()
	defer kvs.Unlock()
	_, ok := kvs.storage[k]
	if ok {
		return fmt.Errorf("already exist")
	}
	kvs.storage[k] = v
	return nil
}

func (kvs MyKVS) Load(k string) (string, error) {
	kvs.RLock()
	defer kvs.RUnlock()
	v, ok := kvs.storage[k]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return v, nil
}
