package main

import (
	"encoding/base64"
	"encoding/binary"
	"hash"
	"hash/crc32"
)

type Hasher interface {
	// Hash generates the hash value of v.
	Hash(v string) string
}

type MyHasher struct {
	hasher hash.Hash32
}

func NewMyHasher() MyHasher {
	return MyHasher{crc32.NewIEEE()}
}

func (h MyHasher) Hash(v string) string {
	h.hasher.Write([]byte(v))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, h.hasher.Sum32())
	h.hasher.Reset()
	str := base64.URLEncoding.EncodeToString(b)
	return str[:8]
}
