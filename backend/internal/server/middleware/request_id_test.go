package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRequestIDMiddleware_GeneratesIDWhenAbsent(t *testing.T) {
	m := NewRequestIDMiddleware()

	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured == "" {
		t.Fatalf("expected request id to be generated")
	}
	if got := rr.Header().Get(RequestIDHeader); got != captured {
		t.Fatalf("response header mismatch: got %q want %q", got, captured)
	}
}

func TestRequestIDMiddleware_PropagatesIncomingHeader(t *testing.T) {
	m := NewRequestIDMiddleware()

	incoming := uuid.NewString()
	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	ctx := context.WithValue(context.Background(), TrustedProxyKey, true)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, incoming)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured != incoming {
		t.Fatalf("context request id mismatch: got %q want %q", captured, incoming)
	}
	if got := rr.Header().Get(RequestIDHeader); got != incoming {
		t.Fatalf("response header mismatch: got %q want %q", got, incoming)
	}
}

func TestRequestIDMiddleware_RegeneratesWhenNotTrusted(t *testing.T) {
	m := NewRequestIDMiddleware()
	incoming := uuid.NewString()

	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, incoming)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured == incoming {
		t.Fatalf("expected request id to be regenerated for untrusted source")
	}
	if _, err := uuid.Parse(captured); err != nil {
		t.Fatalf("expected generated request id to be UUID, got %q", captured)
	}
}

func TestRequestIDMiddleware_RegeneratesWhenHeaderInvalid(t *testing.T) {
	m := NewRequestIDMiddleware()

	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	ctx := context.WithValue(context.Background(), TrustedProxyKey, true)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured == "not-a-uuid" {
		t.Fatalf("expected invalid request id to be rejected")
	}
	if _, err := uuid.Parse(captured); err != nil {
		t.Fatalf("expected generated request id to be UUID, got %q", captured)
	}
}

func TestRequestIDMiddleware_RegeneratesWhenHeaderTooLong(t *testing.T) {
	m := NewRequestIDMiddleware()

	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	ctx := context.WithValue(context.Background(), TrustedProxyKey, true)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, strings.Repeat("a", maxRequestIDLength+1))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if len(captured) > maxRequestIDLength {
		t.Fatalf("expected long request id to be rejected")
	}
	if _, err := uuid.Parse(captured); err != nil {
		t.Fatalf("expected generated request id to be UUID, got %q", captured)
	}
}

func TestRequestIDFromContext_NilContext(t *testing.T) {
	if got := RequestIDFromContext(nil); got != "" { //nolint:staticcheck // intentional nil check
		t.Fatalf("expected empty string for nil context, got %q", got)
	}
}

func TestRequestIDFromContext_NoValue(t *testing.T) {
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Fatalf("expected empty string when no request id, got %q", got)
	}
}
