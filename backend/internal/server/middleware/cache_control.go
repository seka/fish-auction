package middleware

import "net/http"

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
