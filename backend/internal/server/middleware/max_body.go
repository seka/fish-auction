package middleware

import (
	"net/http"
)

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
