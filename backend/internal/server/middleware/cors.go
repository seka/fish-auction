package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware handles Cross-Origin Resource Sharing.
type CORSMiddleware struct {
	allowedOrigins []string
	c              *cors.Cors
}

// NewCORSMiddleware creates a new CORSMiddleware instance.
func NewCORSMiddleware(allowedOrigins []string) *CORSMiddleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
	})
	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
		c:              c,
	}
}

// Handle wraps an http.Handler with CORS logic.
func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return m.c.Handler(next)
}
