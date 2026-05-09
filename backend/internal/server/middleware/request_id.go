package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDHeader is the HTTP header used to propagate request IDs through ALB / proxies.
const RequestIDHeader = "X-Request-ID"

// RequestIDKey は request_id を context に保持するためのキー。
const RequestIDKey contextKey = "request_id"

// RequestIDMiddleware injects a request ID into the request context and reflects it back
// in the response so downstream services / clients can correlate logs.
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware creates a new RequestIDMiddleware.
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// Handle wraps next with request-id propagation.
// 既存の X-Request-ID があれば信頼してそのまま採用し、無ければ UUIDv4 を生成する。
// trustedProxy ミドルウェアの直後に配置することで ALB から付与された X-Request-ID
// が利用可能になる。
func (m *RequestIDMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		w.Header().Set(RequestIDHeader, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the request ID stored in ctx if present.
// 取得できない場合は空文字を返すため、呼び出し側で必須かどうかは設計次第で判断する。
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}
