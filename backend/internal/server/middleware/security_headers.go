package middleware

import "net/http"

type SecurityHeadersMiddleware struct{}

func NewSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

func (m *SecurityHeadersMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リソースの読み込み制限（APIのため一切の読み込みとiframe埋め込みを禁止）
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		// HTTPS接続を強制（約20年の長期記憶）
		w.Header().Set("Strict-Transport-Security", "max-age=631138519")
		// MIMEタイプの推測（スニッフィング）を禁止
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// IE向け：ダウンロードファイルをブラウザで直接開くのを禁止
		w.Header().Set("X-Download-Options", "noopen")
		// クリックジャッキング対策：iframe内での表示を一切禁止
		w.Header().Set("X-Frame-Options", "DENY")
		// Adobe関連リソースからのアクセスを禁止
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		// ブラウザ内蔵XSSフィルタの脆弱性回避のため明示的に無効化
		w.Header().Set("X-XSS-Protection", "0")
		// クロスオリジン遷移時のリファラ漏洩防止
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
