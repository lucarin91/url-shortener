package main

import (
	"errors"
	"sync"
)

type KVS interface {
	// Store stores a key value pair.
	Store(k, v string) error

	// Load returns the value associated with a given key.
	Load(k string) (string, error)

	// Dump the internal storage
	Dump() *Storage
}

type MyKVS struct {
	storage map[string]string
	sync.RWMutex
}

func NewMyKVS() MyKVS {
	return MyKVS{storage: make(map[string]string)}
}

var errAlreadyExists = errors.New("can't add entry: already exists")
var errNotFound = errors.New("can't load entry: not found")

func (kvs MyKVS) Store(k, v string) error {
	kvs.Lock()
	defer kvs.Unlock()
	_, ok := kvs.storage[k]
	if ok {
		return errAlreadyExists
	}
	kvs.storage[k] = v
	return nil
}

func (kvs MyKVS) Load(k string) (string, error) {
	kvs.RLock()
	defer kvs.RUnlock()
	v, ok := kvs.storage[k]
	if !ok {
		return "", errNotFound
	}
	return v, nil
}

func (kvs MyKVS) Dump() *Storage {
	kvs.Lock()
	defer kvs.Unlock()
	stg := Storage{URLPairs: make([]URLPair, 0, len(kvs.storage))}
	for key, value := range kvs.storage {
		stg.URLPairs = append(stg.URLPairs, URLPair{Short: key, Long: value})
	}
	return &stg
}
