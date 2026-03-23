package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware compresses the response using gzip if the client supports it.
type GzipMiddleware struct{}

// NewGzipMiddleware creates a new GzipMiddleware instance.
func NewGzipMiddleware() *GzipMiddleware {
	return &GzipMiddleware{}
}

// Handle provides Handle related functionality.
func (m *GzipMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クライアントがGzip圧縮に対応している場合のみ適用。
		// JSON等のテキストデータの転送量を削減し、APIのレスポンス速度を向上（パフォーマンス最適化）。
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer func() { _ = gz.Close() }()

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
