package admin

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateAdminUseCase defines the interface for creating an admin.
type CreateAdminUseCase interface {
	// Execute creates a new admin with the given email and password.
	Execute(ctx context.Context, email, password string) (*model.Admin, error)
	// Count returns the total number of admins.
	Count(ctx context.Context) (int, error)
}

type createAdminUseCase struct {
	adminRepo repository.AdminRepository
}

var _ CreateAdminUseCase = (*createAdminUseCase)(nil)

// NewCreateAdminUseCase creates a new instance of CreateAdminUseCase
func NewCreateAdminUseCase(adminRepo repository.AdminRepository) CreateAdminUseCase {
	return &createAdminUseCase{adminRepo: adminRepo}
}

func (u *createAdminUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	pwd, err := model.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// 既に管理者が存在するか確認
	existing, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		// NotFound 以外はエラーとして扱う
		var nfErr *apperrors.NotFoundError
		if !errors.As(err, &nfErr) {
			return nil, fmt.Errorf("failed to check existing admin email: %w", err)
		}
	} else if existing != nil {
		return nil, &apperrors.ConflictError{Message: fmt.Sprintf("admin with email %s already exists", email)}
	}

	// 全体のカウントチェック
	count, err := u.adminRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count admins: %w", err)
	}
	if count > 0 {
		return nil, &apperrors.ConflictError{Message: "admin already exists"}
	}

	hashedPassword, err := pwd.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	admin := &model.Admin{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := u.adminRepo.Create(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	return admin, nil
}

func (u *createAdminUseCase) Count(ctx context.Context) (int, error) {
	count, err := u.adminRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count admins: %w", err)
	}
	return count, nil
}
