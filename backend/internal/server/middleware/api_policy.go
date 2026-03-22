package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// CacheControlMiddleware adds Cache-Control headers to prevent caching of API responses.
type CacheControlMiddleware struct{}

func NewCacheControlMiddleware() *CacheControlMiddleware {
	return &CacheControlMiddleware{}
}

func (m *CacheControlMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// MaxBodyMiddleware limits the request body size to prevent DoS attacks.
type MaxBodyMiddleware struct {
	limit int64
}

func NewMaxBodyMiddleware(limit int64) *MaxBodyMiddleware {
	return &MaxBodyMiddleware{
		limit: limit,
	}
}

func (m *MaxBodyMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, m.limit)
		next.ServeHTTP(w, r)
	})
}

// GzipMiddleware compresses the response using gzip if the client supports it.
type GzipMiddleware struct{}

func NewGzipMiddleware() *GzipMiddleware {
	return &GzipMiddleware{}
}

func (m *GzipMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

// gzipResponseWriter wraps http.ResponseWriter to support gzip compression.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
