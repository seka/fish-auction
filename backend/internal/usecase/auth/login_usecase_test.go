package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/service"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
)

type mockClock struct{}

func (c *mockClock) Now() time.Time                       { return time.Now() }
func (c *mockClock) NowIn(_ model.LocationName) time.Time { return time.Now() }

var _ service.Clock = (*mockClock)(nil)

type mockAdminRepository struct {
	admin              *model.Admin
	err                error
	failedAttempts     int64
	lockCalled         bool
	loginSuccessCalled bool
}

func (m *mockAdminRepository) FindByID(_ context.Context, id int) (*model.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin != nil && m.admin.ID == id {
		return m.admin, nil
	}
	return nil, &apperrors.NotFoundError{Resource: "Admin", ID: id}
}

func (m *mockAdminRepository) FindOneByEmail(_ context.Context, email string) (*model.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Return the admin if the email matches (simulating DB search)
	if m.admin != nil && m.admin.Email == email {
		return m.admin, nil
	}
	return nil, &apperrors.NotFoundError{Resource: "Admin", ID: 0}
}

func (m *mockAdminRepository) Create(_ context.Context, _ *model.Admin) error {
	return nil
}

func (m *mockAdminRepository) Count(_ context.Context) (int, error) {
	return 0, nil
}

func (m *mockAdminRepository) UpdatePassword(_ context.Context, _ int, _ string) error {
	return nil
}

func (m *mockAdminRepository) IncrementFailedAttempts(_ context.Context, _ int) (int64, error) {
	return m.failedAttempts, nil
}

func (m *mockAdminRepository) LockAccount(_ context.Context, _ int, _ time.Time) error {
	m.lockCalled = true
	return nil
}

func (m *mockAdminRepository) UpdateLoginSuccess(_ context.Context, _ int) error {
	m.loginSuccessCalled = true
	return nil
}

func TestLoginUseCase_AccountLocked(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	future := time.Now().Add(30 * time.Minute)
	repo := &mockAdminRepository{admin: &model.Admin{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		LockedUntil:  &future,
	}}
	uc := auth.NewLoginUseCase(repo, &mockClock{})

	_, err := uc.Execute(context.Background(), "admin@example.com", "password")
	if err == nil {
		t.Fatal("expected error for locked account, got nil")
	}
	var unauth *apperrors.UnauthorizedError
	if !errors.As(err, &unauth) {
		t.Errorf("expected UnauthorizedError, got %T: %v", err, err)
	}
}

func TestLoginUseCase_LockExpired(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	past := time.Now().Add(-1 * time.Minute)
	repo := &mockAdminRepository{admin: &model.Admin{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		LockedUntil:  &past,
	}}
	uc := auth.NewLoginUseCase(repo, &mockClock{})

	got, err := uc.Execute(context.Background(), "admin@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error after lock expiry, got %v", err)
	}
	if got == nil {
		t.Fatal("expected admin returned after lock expiry")
	}
}

func TestLoginUseCase_LockTriggeredAtThreshold(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	repo := &mockAdminRepository{
		admin:          &model.Admin{Email: "admin@example.com", PasswordHash: string(hash)},
		failedAttempts: auth.MaxAdminFailedLoginAttempts,
	}
	uc := auth.NewLoginUseCase(repo, &mockClock{})

	_, err := uc.Execute(context.Background(), "admin@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error on wrong password")
	}
	if !repo.lockCalled {
		t.Error("expected LockAccount to be called when attempts reach threshold")
	}
}

func TestLoginUseCase_SuccessCallsUpdateLoginSuccess(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	repo := &mockAdminRepository{admin: &model.Admin{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
	}}
	uc := auth.NewLoginUseCase(repo, &mockClock{})

	if _, err := uc.Execute(context.Background(), "admin@example.com", "password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.loginSuccessCalled {
		t.Error("expected UpdateLoginSuccess to be called on successful login")
	}
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Generate a valid has for "admin-password"
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin-password"), bcrypt.MinCost)

	validAdmin := &model.Admin{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
	}

	tests := []struct {
		name      string
		email     string
		password  string
		mockAdmin *model.Admin
		mockErr   error
		wantFound bool
		wantErr   bool
	}{
		{
			name:      "CorrectCredentials",
			email:     "admin@example.com",
			password:  "admin-password",
			mockAdmin: validAdmin,
			wantFound: true,
			wantErr:   false,
		},
		{
			name:      "WrongPassword",
			email:     "admin@example.com",
			password:  "wrong-password",
			mockAdmin: validAdmin,
			wantFound: false,
			wantErr:   true,
		},
		{
			name:      "UserNotFound",
			email:     "other@example.com",
			password:  "admin-password",
			mockAdmin: validAdmin,
			wantFound: false,
			wantErr:   true,
		},
		{
			name:      "RepoError",
			email:     "admin@example.com",
			password:  "admin-password",
			mockErr:   errors.New("db error"),
			wantFound: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepository{
				admin: tt.mockAdmin,
				err:   tt.mockErr,
			}
			uc := auth.NewLoginUseCase(repo, &mockClock{})

			gotAdmin, err := uc.Execute(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
			if (gotAdmin != nil) != tt.wantFound {
				t.Fatalf("expected found=%v, got admin=%v", tt.wantFound, gotAdmin)
			}
			if tt.wantFound && gotAdmin.ID != validAdmin.ID { // Assuming validAdmin has ID 0 if not set, but pointers should match if logic was mock based?
				// Ah, repo returns mockAdmin, so it should be same object or equivalent.
				if gotAdmin.Email != validAdmin.Email {
					t.Errorf("expected email %s, got %s", validAdmin.Email, gotAdmin.Email)
				}
			}
		})
	}
}
