package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

func CheckWrongRedirect(url string) error {
	r, err := logGet(shortlyURL + url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusNotFound {
		return fmt.Errorf("code %v, want %v", r.StatusCode, http.StatusBadRequest)
	}
	fmt.Println("  OK!")
	return nil
}

func CheckGoodRedirect(short, long string) error {
	r, err := logGet(shortlyURL + short)
	if err != nil {
		return err
	}
	if http.StatusOK != r.StatusCode {
		return fmt.Errorf("code %v, want %v", r.StatusCode, http.StatusOK)
	}
	u, _ := url.Parse("http://" + long)
	if u.Host != r.Request.URL.Host || u.Path != r.Request.URL.Path {
		return fmt.Errorf("redirect to %q, want %q", r.Request.URL, long)
	}
	fmt.Println("  OK!")
	return nil
}

func CheckShortURL(url string) (string, error) {
	r, err := logGet(shortlyURL + "/shorten" + url)
	if err != nil {
		return "", err
	}
	if http.StatusOK != r.StatusCode {
		return "", fmt.Errorf("code %v, want %v", r.StatusCode, http.StatusOK)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body.Close()
	fmt.Println("  OK!")
	return string(body)[len(shortlyURL):], nil
}

func CheckSatistics(statsCheck *Stats) error {
	r, err := logGet(shortlyURL + "/statistics?format=json")
	if err != nil {
		return err
	}
	if http.StatusOK != r.StatusCode {
		return fmt.Errorf("code %v, want %v", r.StatusCode, http.StatusOK)
	}
	e := json.NewDecoder(r.Body)
	var stats Stats
	if err := e.Decode(&stats); err != nil {
		return err
	}
	r.Body.Close()

	if !reflect.DeepEqual(statsCheck, &stats) {
		return fmt.Errorf("statistics %+v,\n want %+v", statsCheck, stats)
	}
	fmt.Println("  OK!")
	return nil
}

func CheckStorage(finalStorage *Storage) error {
	fmt.Println("\n  Check storage")
	stg, err := LoadStorageFile(urlsPath)
	if err != nil {
		return err
	}

	if !stg.Equal(finalStorage) {
		return fmt.Errorf("storage %+v, want %+v", stg, finalStorage)
	}
	fmt.Println("  OK!")
	return nil
}

func logGet(url string) (*http.Response, error) {
	fmt.Printf("\n  GET: %v\n", url)
	return http.Get(url)
}
