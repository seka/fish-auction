package testing

import (
	"context"
	"time"
)

// MockRateLimitRepository is a no-op implementation of RateLimitRepository for testing.
// All methods return (0, nil) so requests always pass through (fail-open).
type MockRateLimitRepository struct{}

func (m *MockRateLimitRepository) IncrementAdminLogin(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, nil
}

func (m *MockRateLimitRepository) IncrementBuyerLogin(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, nil
}

func (m *MockRateLimitRepository) IncrementAdminReset(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, nil
}

func (m *MockRateLimitRepository) IncrementBuyerReset(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, nil
}
