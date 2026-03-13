package admin

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdminUseCase defines the interface for creating an admin
type CreateAdminUseCase interface {
	Execute(ctx context.Context, email, password string) (*model.Admin, error)
	Count(ctx context.Context) (int, error)
}

type createAdminUseCase struct {
	adminRepo repository.AdminRepository
}

var _ CreateAdminUseCase = (*createAdminUseCase)(nil)

// NewCreateAdminUseCase creates a new instance of CreateAdminUseCase
func NewCreateAdminUseCase(adminRepo repository.AdminRepository) *createAdminUseCase {
	return &createAdminUseCase{adminRepo: adminRepo}
}

func (u *createAdminUseCase) Execute(ctx context.Context, email, password string) (*model.Admin, error) {
	// 既に管理者が存在するか確認 (個別のメールアドレス)
	existing, err := u.adminRepo.FindOneByEmail(ctx, email)
	if err != nil {
		// NotFound 以外はエラーとして扱う
		var nfErr *apperrors.NotFoundError
		if !errors.As(err, &nfErr) {
			return nil, err
		}
	} else if existing != nil {
		return nil, fmt.Errorf("admin with email %s already exists", email)
	}

	// 全体のカウントチェック (初期管理者のみ許可する場合などのため)
	count, err := u.adminRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("admin already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := u.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (uc *createAdminUseCase) Count(ctx context.Context) (int, error) {
	return uc.adminRepo.Count(ctx)
}
