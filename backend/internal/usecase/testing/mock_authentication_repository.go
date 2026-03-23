package testing

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

// MockAuthenticationRepository is a mock implementation of AuthenticationRepository
type MockAuthenticationRepository struct {
	CreateFunc                  func(ctx context.Context, auth *model.Authentication) (*model.Authentication, error)
	FindByEmailFunc             func(ctx context.Context, email string) (*model.Authentication, error)
	FindByBuyerIDFunc           func(ctx context.Context, buyerID int) (*model.Authentication, error)
	UpdateLoginSuccessFunc      func(ctx context.Context, id int, loginAt time.Time) error
	IncrementFailedAttemptsFunc func(ctx context.Context, id int) error
	ResetFailedAttemptsFunc     func(ctx context.Context, id int) error
	LockAccountFunc             func(ctx context.Context, id int, until time.Time) error
	UpdatePasswordFunc          func(ctx context.Context, buyerID int, passwordHash string) error
}

// Create creates a new record.
func (m *MockAuthenticationRepository) Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
	return m.CreateFunc(ctx, auth)
}

// FindByEmail retrieves a record based on criteria.
func (m *MockAuthenticationRepository) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	return m.FindByEmailFunc(ctx, email)
}

// FindByBuyerID retrieves a record based on criteria.
func (m *MockAuthenticationRepository) FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error) {
	return m.FindByBuyerIDFunc(ctx, buyerID)
}

// UpdateLoginSuccess updates an existing record.
func (m *MockAuthenticationRepository) UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error {
	return m.UpdateLoginSuccessFunc(ctx, id, loginAt)
}

// IncrementFailedAttempts provides IncrementFailedAttempts related functionality.
func (m *MockAuthenticationRepository) IncrementFailedAttempts(ctx context.Context, id int) error {
	return m.IncrementFailedAttemptsFunc(ctx, id)
}

// ResetFailedAttempts provides ResetFailedAttempts related functionality.
func (m *MockAuthenticationRepository) ResetFailedAttempts(ctx context.Context, id int) error {
	return m.ResetFailedAttemptsFunc(ctx, id)
}

// LockAccount provides LockAccount related functionality.
func (m *MockAuthenticationRepository) LockAccount(ctx context.Context, id int, until time.Time) error {
	return m.LockAccountFunc(ctx, id, until)
}

// UpdatePassword updates an existing record.
func (m *MockAuthenticationRepository) UpdatePassword(ctx context.Context, buyerID int, passwordHash string) error {
	return m.UpdatePasswordFunc(ctx, buyerID, passwordHash)
}
