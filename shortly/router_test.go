package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRedirect(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		s := NewServer()

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		s.redirect(w, req)

		resp := w.Result()

		if http.StatusBadRequest != resp.StatusCode {
			t.Errorf("got %q, want %q", resp.StatusCode, http.StatusBadRequest)
		}

	})

	t.Run("not found", func(t *testing.T) {
		s := NewServer()

		req := httptest.NewRequest("GET", "/12345678", nil)
		w := httptest.NewRecorder()
		s.redirect(w, req)

		resp := w.Result()

		if http.StatusNotFound != resp.StatusCode {
			t.Errorf("got %q, want %q", resp.StatusCode, http.StatusNotFound)
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
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusMovedPermanently)
		}
		if "http://value1" != resp.Header.Get("Location") {
			t.Errorf("got %q, want %q", resp.Header.Get("Location"), "http://value1")
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

		if http.StatusBadRequest != resp.StatusCode {
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusBadRequest)
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

		if http.StatusBadRequest != resp.StatusCode {
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusBadRequest)
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
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusOK)
		}
		if "http://example.com/keykey12" != string(body) {
			t.Errorf("got %q, want %q", string(body), "http://example.com/keykey12")
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
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusOK)
		}
		r := "URLs: 1\n" +
			"Redirect: 3\n" +
			"Handler:\n" +
			"  /: 4\n" +
			"  /shorten: 5\n" +
			"  /statistics: 6\n"
		if r != string(body) {
			t.Errorf("got %v, want %v", string(body), r)
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
			t.Errorf("got %v, want %v", resp.StatusCode, http.StatusOK)
		}

		res := Stats{}
		json.Unmarshal(body, &res)
		if !reflect.DeepEqual(&s.Stats, &res) {
			t.Errorf("got %+v, want %+v", res, s.Stats)
		}
	})
}
