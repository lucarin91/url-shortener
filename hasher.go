package main

import (
	"crypto/sha1"
	"encoding/base64"
	"hash"
)

type Hasher interface {
	// Hash generates the hash value of v.
	Hash(v string) string
}

type MyHasher struct {
	hasher hash.Hash
}

func NewMyHasher() MyHasher {
	return MyHasher{sha1.New()}
}

func (h MyHasher) Hash(v string) string {
	shaHash := h.hasher.Sum([]byte(v))
	str := base64.URLEncoding.EncodeToString(shaHash)
	return str[:8]
}
