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
	if err != nil {
		return nil, fmt.Errorf("cannot open storage file %q: %v", path, err)
	}
	var stg Storage
	e := json.NewDecoder(f)
	err = e.Decode(&stg)
	if err != nil {
		return nil, fmt.Errorf("cannot decode storage file %q: %v", path, err)
	}
	return &stg, nil
}
