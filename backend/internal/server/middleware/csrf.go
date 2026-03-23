package middleware

import (
	"net/http"
	"slices"
)

// CSRFMiddleware protects against Cross-Site Request Forgery.
// 現代的なCSRF対策として、ブラウザが提供する Origin と Fetch Metadata を検証します。
// 参考: 令和時代の API 実装のベースプラクティスと CSRF 対策 | blog.jxck.io
// https://blog.jxck.io/entries/2024-04-26/csrf.html
//
// 対策の仕組み:
// 1. SameSite Cookie (Lax/Strict): ブラウザがクロスドメインのリクエストに自動でクッキーを付与するのを抑制。
// 2. Origin Check: リクエストの発生元ドメインが許可されたものか検証。
// 3. Fetch Metadata Check: ブラウザが自動付与する Sec-Fetch-Site ヘッダーを検証し、cross-site を拒否。
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
			allowed := slices.Contains(m.allowedOrigins, origin)
			if !allowed {
				http.Error(w, "Forbidden (Origin mismatch)", http.StatusForbidden)
				return
			}
		}

		// ブラウザ以外のクライアント（curl コマンド、モバイルアプリ、古いツールなど）からリクエストを送る場合、Origin や Sec-Fetch-Site といったヘッダーをあえて付けずに送ることが可能。
		// API向け：Origin も Sec-Fetch-Site もない副作用リクエスト（直撃など）は一旦許可するが、より厳格にするなら Origin なしのPOSTを拒否する選択もあり得る。
		// ここでは実用性を考慮し、明示的な不一致がある場合のみ拒否する。
		next.ServeHTTP(w, r)
	})
}
