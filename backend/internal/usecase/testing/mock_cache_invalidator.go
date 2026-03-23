package testing

import "context"

// MockCacheInvalidator is a mock implementation of CacheInvalidator for testing.
type MockCacheInvalidator struct {
	InvalidateCacheFunc func(ctx context.Context, id int) error
}

// InvalidateCache provides InvalidateCache related functionality.
func (m *MockCacheInvalidator) InvalidateCache(ctx context.Context, id int) error {
	if m.InvalidateCacheFunc != nil {
		return m.InvalidateCacheFunc(ctx, id)
	}
	return nil
}
