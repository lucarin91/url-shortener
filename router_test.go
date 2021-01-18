package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		s := NewServer()

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		s.redirect(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusBadRequest != resp.StatusCode {
			t.Fail()
		}
		if "error: invalid url" != string(body) {
			t.Fail()
		}
	})

	t.Run("not found", func(t *testing.T) {
		s := NewServer()

		req := httptest.NewRequest("GET", "/12345678", nil)
		w := httptest.NewRecorder()
		s.redirect(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusNotFound != resp.StatusCode {
			t.Fail()
		}
		if "error: not found" != string(body) {
			t.Fail()
		}
	})

	t.Run("redirect", func(t *testing.T) {
		s := NewServer()
		s.Kvs.Store("keykey12", "value1")

		req := httptest.NewRequest("GET", "/keykey12", nil)
		w := httptest.NewRecorder()
		s.redirect(w, req)

		resp := w.Result()
		if http.StatusMovedPermanently != resp.StatusCode {
			t.Fail()
		}
		if "http://value1" != resp.Header.Get("Location") {
			t.Fail()
		}
	})

}

type MockedHasher struct{}

func (h MockedHasher) Hash(s string) string {
	return s
}

func TestShorten(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		s := NewServer()

		req := httptest.NewRequest("GET", "/shorten/", nil)
		w := httptest.NewRecorder()
		s.shorten(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusBadRequest != resp.StatusCode {
			t.Fail()
		}
		if "error: invalid url" != string(body) {
			t.Fail()
		}
	})

	t.Run("already exist", func(t *testing.T) {
		s := NewServer()
		s.Kvs.Store("keykey12", "value1")
		s.Hasher = MockedHasher{}

		req := httptest.NewRequest("GET", "/shorten/keykey12", nil)
		w := httptest.NewRecorder()
		s.shorten(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusBadRequest != resp.StatusCode {
			t.Fail()
		}
		if "error: already exist" != string(body) {
			t.Fail()
		}
	})

	t.Run("shorten", func(t *testing.T) {
		s := NewServer()
		s.Hasher = MockedHasher{}

		req := httptest.NewRequest("GET", "/shorten/keykey12", nil)
		w := httptest.NewRecorder()
		s.shorten(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusOK != resp.StatusCode {
			t.Fail()
		}
		if "http://example.com/keykey12" != string(body) {
			t.Fail()
		}
	})

}

func TestStatistics(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		s := NewServer()
		s.Stats.TotalURL = 1
		s.Stats.Redirects.Failed = 2
		s.Stats.Redirects.Success = 3
		s.Stats.Handlers[Redirect].Count = 4
		s.Stats.Handlers[Shorten].Count = 5
		s.Stats.Handlers[Statistics].Count = 6

		req := httptest.NewRequest("GET", "/statistics", nil)
		w := httptest.NewRecorder()
		s.statistics(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusOK != resp.StatusCode {
			t.Fail()
		}
		r := "URLs: 1\n" +
			"Redirect: 3\n" +
			"Handler:\n" +
			"  /: 4\n" +
			"  /shorten: 5\n" +
			"  /statistics: 6\n"
		if r != string(body) {
			t.Fail()
		}
	})

	t.Run("json", func(t *testing.T) {
		s := NewServer()
		s.Stats.TotalURL = 1
		s.Stats.Redirects.Failed = 2
		s.Stats.Redirects.Success = 3
		s.Stats.Handlers[Redirect].Count = 4
		s.Stats.Handlers[Shorten].Count = 5
		s.Stats.Handlers[Statistics].Count = 6

		req := httptest.NewRequest("GET", "/statistics?format=json", nil)
		w := httptest.NewRecorder()
		s.statistics(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if http.StatusOK != resp.StatusCode {
			t.Fail()
		}
		r, _ := json.MarshalIndent(s.Stats, "", "  ")
		if string(r)+"\n" != string(body) {
			t.Fail()
		}
	})
}
