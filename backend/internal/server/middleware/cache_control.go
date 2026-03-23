package middleware

import "net/http"

// CacheControlMiddleware adds Cache-Control headers to prevent caching of API responses.
type CacheControlMiddleware struct{}

// NewCacheControlMiddleware creates a new CacheControlMiddleware instance.
func NewCacheControlMiddleware() *CacheControlMiddleware {
	return &CacheControlMiddleware{}
}

// Handle provides Handle related functionality.
func (m *CacheControlMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no-store: キャッシュを保存させない。個人情報や最新性が重要なAPIレスポンスの保護。
		// no-cache, must-revalidate: キャッシュを使用する前に必ずサーバーに再検証を求める。
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		// 古いブラウザ（HTTP/1.0）向けの互換性維持
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
