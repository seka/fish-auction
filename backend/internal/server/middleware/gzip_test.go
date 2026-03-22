package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
