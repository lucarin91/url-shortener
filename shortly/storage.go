package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage struct {
	URLPairs []URLPair `json:"url_pairs"`
}

type URLPair struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func LoadStorageFile(path string) (*Storage, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}
	var stg Storage
	e := json.NewDecoder(f)
	err = e.Decode(&stg)
	if err != nil {
		return nil, fmt.Errorf("decode storage: %w", err)
	}
	return &stg, nil
}

func (stg *Storage) SaveStorageFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0664)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("open storage: %v", err)
	}
	e := json.NewEncoder(f)
	e.Encode(stg)
	if err != nil {
		return fmt.Errorf("encode storage: %v", err)
	}
	return nil
}
