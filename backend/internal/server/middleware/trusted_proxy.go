package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// TrustedProxyMiddleware rewrites RemoteAddr / URL.Scheme based on
// X-Forwarded-For / X-Forwarded-Proto when the request originates from a
// trusted proxy network (e.g. ALB).
//
// 信頼ネットワークが空の場合は no-op として動作し、ローカル開発時の挙動を変えない。
type TrustedProxyMiddleware struct {
	trusted []*net.IPNet
}

// TrustedProxyKey indicates whether the request passed trusted proxy validation.
const TrustedProxyKey contextKey = "trusted_proxy"

// NewTrustedProxyMiddleware parses CIDR strings into a trusted proxy network set.
// 不正な CIDR は無視され、空配列なら本ミドルウェアは何もしない。
func NewTrustedProxyMiddleware(cidrs []string) *TrustedProxyMiddleware {
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		_, n, err := net.ParseCIDR(c)
		if err != nil {
			continue
		}
		nets = append(nets, n)
	}
	return &TrustedProxyMiddleware{trusted: nets}
}

// Handle wraps next with proxy header reflection.
func (m *TrustedProxyMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(m.trusted) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
			port = ""
		}
		if !m.isTrusted(host) {
			next.ServeHTTP(w, r)
			return
		}

		// X-Forwarded-For を右から走査し、信頼ネットワーク外の最初の IP を実クライアントとみなす。
		// ALB は client → ALB → app の経路を `client, alb-internal` のように追記する。
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			for i := len(parts) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(parts[i])
				if ip == "" {
					continue
				}
				if !m.isTrusted(ip) {
					// 後段で net.SplitHostPort を呼ぶ既存コードに合わせ host:port 形式を維持する。
					if port != "" {
						r.RemoteAddr = net.JoinHostPort(ip, port)
					} else {
						r.RemoteAddr = ip
					}
					break
				}
			}
		}

		// 終端が ALB の HTTPS リスナーであれば X-Forwarded-Proto=https が付く。
		// ハンドラが r.URL.Scheme を見るケース向けに反映する。
		if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			r.URL.Scheme = strings.ToLower(proto)
		}

		ctx := context.WithValue(r.Context(), TrustedProxyKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IsFromTrustedProxy returns true when trusted proxy middleware validated the request source.
func IsFromTrustedProxy(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	trusted, _ := ctx.Value(TrustedProxyKey).(bool)
	return trusted
}

func (m *TrustedProxyMiddleware) isTrusted(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, n := range m.trusted {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
