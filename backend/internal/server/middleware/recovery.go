package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/seka/fish-auction/backend/internal/server/dto"
)

// RecoveryMiddleware gracefully handles panics, logging the error
// and returning a 500 Internal Server Error response.
type RecoveryMiddleware struct{}

// NewRecoveryMiddleware creates a new RecoveryMiddleware instance.
func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

// Handle provides Handle related functionality.
func (m *RecoveryMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// パニックエラーとスタックトレースをログ出力
				log.Printf("[PANIC RECOVERED] %v\n%s", err, debug.Stack())

				// クライアントには 500 Internal Server Error のJSONを返す
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
					Error:   "internal_error",
					Message: "Internal Server Error",
					Code:    http.StatusInternalServerError,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
