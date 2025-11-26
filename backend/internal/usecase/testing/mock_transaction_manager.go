package testing

import "context"

// MockTransactionManager is a mock implementation of TransactionManager
type MockTransactionManager struct {
	WithTransactionFunc func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.WithTransactionFunc != nil {
		return m.WithTransactionFunc(ctx, fn)
	}
	return fn(ctx)
}
