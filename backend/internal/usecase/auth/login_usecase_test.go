package auth_test

import (
	"context"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
)

type mockAdminRepository struct {
	admin *entity.Admin
	err   error
}

func (m *mockAdminRepository) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.admin != nil && m.admin.ID == id {
		return m.admin, nil
	}
	return nil, nil
}

func (m *mockAdminRepository) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Return the admin if the email matches (simulating DB search)
	if m.admin != nil && m.admin.Email == email {
		return m.admin, nil
	}
	return nil, nil // Not found
}

func (m *mockAdminRepository) Create(ctx context.Context, admin *entity.Admin) error {
	return nil
}

func (m *mockAdminRepository) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *mockAdminRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	return nil
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Generate a valid has for "admin-password"
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin-password"), bcrypt.MinCost)

	validAdmin := &entity.Admin{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
	}

	tests := []struct {
		name      string
		email     string
		password  string
		mockAdmin *entity.Admin
		mockErr   error
		wantFound bool
	}{
		{
			name:      "CorrectCredentials",
			email:     "admin@example.com",
			password:  "admin-password",
			mockAdmin: validAdmin,
			wantFound: true,
		},
		{
			name:      "WrongPassword",
			email:     "admin@example.com",
			password:  "wrong-password",
			mockAdmin: validAdmin,
			wantFound: false,
		},
		{
			name:      "UserNotFound",
			email:     "other@example.com",
			password:  "admin-password",
			mockAdmin: validAdmin, // Repo has admin@example.com, but we search for other
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepository{
				admin: tt.mockAdmin,
				err:   tt.mockErr,
			}
			uc := auth.NewLoginUseCase(repo)

			gotAdmin, err := uc.Execute(context.Background(), tt.email, tt.password)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
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
