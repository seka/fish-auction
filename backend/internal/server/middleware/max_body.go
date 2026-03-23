package middleware

import (
	"net/http"
)

// MaxBodyMiddleware limits the request body size to prevent DoS attacks.
type MaxBodyMiddleware struct {
	limit int64
}

// NewMaxBodyMiddleware creates a new MaxBodyMiddleware instance.
func NewMaxBodyMiddleware(limit int64) *MaxBodyMiddleware {
	return &MaxBodyMiddleware{
		limit: limit,
	}
}

// Handle provides Handle related functionality.
func (m *MaxBodyMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// http.MaxBytesReader: 指定サイズを超えた読み込みが発生した時点でエラーを返し接続を切る。
		// サーバーのメモリ消費を抑え、巨大なリクエストによるDoS攻撃（Resource Exhaustion）を防止。
		r.Body = http.MaxBytesReader(w, r.Body, m.limit)
		next.ServeHTTP(w, r)
	})
}
