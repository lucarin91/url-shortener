package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Storage struct {
	URLPairs []URLPair `json:"url_pairs"`
}

type URLPair struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func LoadStorageFile(path string) (*Storage, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open storage file %q: %w", path, err)
	}
	var stg Storage
	err = json.Unmarshal(data, &stg)
	if err != nil {
		return nil, fmt.Errorf("cannot decode storage file %q: %w", path, err)
	}
	return &stg, nil
}
