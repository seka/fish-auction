package buyer

import (
	"context"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBuyerUseCase defines the interface for creating a buyer with authentication.
type CreateBuyerUseCase interface {
	Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error)
}

type createBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
	authRepo  repository.AuthenticationRepository
	txRepo    repository.TransactionManager
}

// NewCreateBuyerUseCase creates a new instance of CreateBuyerUseCase.
func NewCreateBuyerUseCase(
	buyerRepo repository.BuyerRepository,
	authRepo repository.AuthenticationRepository,
	txRepo repository.TransactionManager,
) CreateBuyerUseCase {
	return &createBuyerUseCase{
		buyerRepo: buyerRepo,
		authRepo:  authRepo,
		txRepo:    txRepo,
	}
}

// Execute creates a new buyer with authentication
func (uc *createBuyerUseCase) Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
	// 0. Validate password
	pwd, err := model.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := pwd.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var createdBuyer *model.Buyer

	// 1. Transactional creation
	err = uc.txRepo.WithTransaction(ctx, func(ctx context.Context) error {
		// 1-1. Create buyer profile
		buyer := &model.Buyer{
			Name:         name,
			Organization: organization,
			ContactInfo:  contactInfo,
		}
		res, err := uc.buyerRepo.Create(ctx, buyer)
		if err != nil {
			return fmt.Errorf("failed to create buyer profile: %w", err)
		}
		createdBuyer = res

		// 1-2. Create auth record
		auth := &model.Authentication{
			BuyerID:      createdBuyer.ID,
			Email:        email,
			PasswordHash: hashedPassword,
			AuthType:     "buyer",
		}
		_, err = uc.authRepo.Create(ctx, auth)
		if err != nil {
			return fmt.Errorf("failed to create auth record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create buyer with transaction: %w", err)
	}

	return createdBuyer, nil
}
