package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

// Mock for CreateAdmin
type mockAdminRepositoryForCreate struct {
	existingAdmin *entity.Admin
	createErr     error
	repoErr       error
}

func (m *mockAdminRepositoryForCreate) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	return nil, nil
}
func (m *mockAdminRepositoryForCreate) FindOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	if m.existingAdmin != nil && m.existingAdmin.Email == email {
		return m.existingAdmin, nil
	}
	return nil, nil
}
func (m *mockAdminRepositoryForCreate) Create(ctx context.Context, admin *entity.Admin) error {
	return m.createErr
}
func (m *mockAdminRepositoryForCreate) Count(ctx context.Context) (int, error) {
	return 0, nil
}
func (m *mockAdminRepositoryForCreate) UpdatePassword(ctx context.Context, id int, hash string) error {
	return nil
}

func TestCreateAdminUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		existingAdmin *entity.Admin
		repoErr       error
		createErr     error
		wantErr       bool
	}{
		{
			name:     "Success",
			email:    "new@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:          "AlreadyExists",
			email:         "existing@example.com",
			password:      "password123",
			existingAdmin: &entity.Admin{ID: 1, Email: "existing@example.com"},
			wantErr:       true,
		},
		{
			name:     "RepoCheckError",
			email:    "error@example.com",
			password: "password123",
			repoErr:  errors.New("db error"),
			wantErr:  true,
		},
		{
			name:      "CreateError",
			email:     "new@example.com",
			password:  "password123",
			createErr: errors.New("create failed"),
			wantErr:   true,
		},
		{
			name:     "PasswordTooLong",
			email:    "long@example.com",
			password: "this_password_is_definitely_way_too_long_to_be_hashed_by_bcrypt_because_it_exceeds_seventy_two_bytes_limit",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepositoryForCreate{
				existingAdmin: tt.existingAdmin,
				repoErr:       tt.repoErr,
				createErr:     tt.createErr,
			}
			uc := admin.NewCreateAdminUseCase(repo)
			err := uc.Execute(context.Background(), tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateAdminUseCase_Count(t *testing.T) {
	uc := admin.NewCreateAdminUseCase(&mockAdminRepositoryForCreate{})
	count, err := uc.Count(context.Background())
	if err != nil {
		t.Errorf("Count() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Count() = %v, want 0", count)
	}
}
