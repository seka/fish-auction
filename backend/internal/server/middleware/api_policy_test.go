package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCacheControlMiddleware(t *testing.T) {
	middleware := NewCacheControlMiddleware()
	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if val := rr.Header().Get("Cache-Control"); !strings.Contains(val, "no-store") {
		t.Errorf("Expected Cache-Control to contain no-store, got %q", val)
	}
}

func TestMaxBodyMiddleware(t *testing.T) {
	// Limit to 10 bytes
	middleware := NewMaxBodyMiddleware(10)
	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("Under limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345"))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	})

	t.Run("Over limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("1234567890123"))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected status RequestEntityTooLarge, got %v", rr.Code)
		}
	})
}

func TestGzipMiddleware(t *testing.T) {
	middleware := NewGzipMiddleware()
	handler := middleware.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! This is a long enough string to benefit from compression."))
	}))

	t.Run("No gzip support", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Header().Get("Content-Encoding") != "" {
			t.Error("Expected no Content-Encoding header")
		}
	})

	t.Run("Gzip support", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Header().Get("Content-Encoding") != "gzip" {
			t.Error("Expected Content-Encoding: gzip")
		}

		gz, err := gzip.NewReader(rr.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer gz.Close()

		body, _ := io.ReadAll(gz)
		if string(body) != "Hello, World! This is a long enough string to benefit from compression." {
			t.Errorf("Expected original body, got %q", string(body))
		}
	})
}
