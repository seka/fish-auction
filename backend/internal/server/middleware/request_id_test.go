package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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

	const incoming = "test-fixed-id"
	var captured string
	handler := m.Handle(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
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
