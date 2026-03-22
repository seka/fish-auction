package middleware

import (
	"net/http"
)

// CSRFMiddleware protects against Cross-Site Request Forgery.
// It leverages Origin and Fetch Metadata headers for modern browser-based protection.
type CSRFMiddleware struct {
	allowedOrigins []string
}

// NewCSRFMiddleware creates a new CSRFMiddleware instance.
func NewCSRFMiddleware(allowedOrigins []string) *CSRFMiddleware {
	return &CSRFMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Handle wraps an http.Handler with CSRF protection logic.
func (m *CSRFMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET, HEAD, OPTIONS, TRACE は副作用がない「安全な」メソッドとされるためスキップ
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// 1. Fetch Metadata (Sec-Fetch-Site) の検証
		// ブラウザがサポートしている場合、same-origin または same-site でない限り拒否
		fetchSite := r.Header.Get("Sec-Fetch-Site")
		if fetchSite != "" {
			if fetchSite != "same-origin" && fetchSite != "same-site" {
				// none はURLの直接入力などのため、APIベースの副作用を許可しない
				http.Error(w, "Forbidden (Sec-Fetch-Site mismatch)", http.StatusForbidden)
				return
			}
		}

		// 2. Origin ヘッダーの検証
		// Sec-Fetch-Site がない場合や古いブラウザ向けの多層防御
		origin := r.Header.Get("Origin")
		if origin != "" {
			allowed := false
			for _, ao := range m.allowedOrigins {
				if origin == ao {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden (Origin mismatch)", http.StatusForbidden)
				return
			}
		}

		// API向け：Origin も Sec-Fetch-Site もない副作用リクエスト（直撃など）をどう扱うか？
		// 令和時代のAPIとしては、ブラウザからの正当なリクエストにはいずれかが付与されているはずなので、
		// 両方ない場合も一旦許可するが、より厳格にするなら Origin なしのPOSTを拒否する選択もあり得る。
		// ここでは実用性を考慮し、明示的な不一致がある場合のみ拒否する。

		next.ServeHTTP(w, r)
	})
}
