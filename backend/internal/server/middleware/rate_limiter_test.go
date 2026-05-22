package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestRateLimiterMiddleware_Handle(t *testing.T) {
	const limit = 5
	window := 10 * time.Minute

	t.Run("under limit passes through", func(t *testing.T) {
		m := NewRateLimiterMiddleware(fixedCount(1), limit, window)
		rr := doRequest(m)
		if rr.Code != http.StatusOK {
			t.Errorf("got %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("at limit passes through", func(t *testing.T) {
		m := NewRateLimiterMiddleware(fixedCount(int64(limit)), limit, window)
		rr := doRequest(m)
		if rr.Code != http.StatusOK {
			t.Errorf("got %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("over limit returns 429", func(t *testing.T) {
		m := NewRateLimiterMiddleware(fixedCount(int64(limit)+1), limit, window)
		rr := doRequest(m)
		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("got %d, want %d", rr.Code, http.StatusTooManyRequests)
		}
	})

	t.Run("over limit sets Retry-After within window", func(t *testing.T) {
		m := NewRateLimiterMiddleware(fixedCount(int64(limit)+1), limit, window)
		rr := doRequest(m)

		raw := rr.Header().Get("Retry-After")
		if raw == "" {
			t.Fatal("Retry-After header is missing")
		}
		secs, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			t.Fatalf("Retry-After is not a number: %s", raw)
		}
		windowSecs := int64(window.Seconds())
		if secs < 1 || secs > windowSecs {
			t.Errorf("Retry-After %d is out of range [1, %d]", secs, windowSecs)
		}
	})

	t.Run("increment error fails open", func(t *testing.T) {
		m := NewRateLimiterMiddleware(errIncrement, limit, window)
		rr := doRequest(m)
		if rr.Code != http.StatusOK {
			t.Errorf("got %d, want %d (expected fail-open)", rr.Code, http.StatusOK)
		}
	})
}

func TestExtractIP(t *testing.T) {
	cases := []struct {
		remoteAddr string
		wantIP     string
	}{
		{"192.168.1.1:12345", "192.168.1.1"},
		{"[::1]:12345", "[::1]"},
		{"192.168.1.1", "192.168.1.1"},
	}
	for _, c := range cases {
		got := extractIP(c.remoteAddr)
		if got != c.wantIP {
			t.Errorf("extractIP(%q) = %q, want %q", c.remoteAddr, got, c.wantIP)
		}
	}
}

func fixedCount(n int64) RateLimitFunc {
	return func(_ context.Context, _ string, _ time.Duration) (int64, error) {
		return n, nil
	}
}

func errIncrement(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, errors.New("redis error")
}

func doRequest(m *RateLimiterMiddleware) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", http.NoBody)
	rr := httptest.NewRecorder()
	m.Handle(http.HandlerFunc(okHandler)).ServeHTTP(rr, req)
	return rr
}
