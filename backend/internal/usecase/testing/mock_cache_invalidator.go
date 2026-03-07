package testing

import "context"

type MockCacheInvalidator struct {
	InvalidateCacheFunc func(ctx context.Context, id int) error
}

func (m *MockCacheInvalidator) InvalidateCache(ctx context.Context, id int) error {
	if m.InvalidateCacheFunc != nil {
		return m.InvalidateCacheFunc(ctx, id)
	}
	return nil
}
