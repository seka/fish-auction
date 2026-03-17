package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateBuyerUseCase defines the interface for creating a buyer.
type CreateBuyerUseCase interface {
	// Execute creates a new buyer with authentication.
	Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error)
}

// createBuyerUseCase handles the creation of buyers
type createBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
	authRepo  repository.AuthenticationRepository
	txMgr     repository.TransactionManager
}

var _ CreateBuyerUseCase = (*createBuyerUseCase)(nil)

// NewCreateBuyerUseCase creates a new instance of CreateBuyerUseCase
func NewCreateBuyerUseCase(buyerRepo repository.BuyerRepository, authRepo repository.AuthenticationRepository, txMgr repository.TransactionManager) *createBuyerUseCase {
	return &createBuyerUseCase{buyerRepo: buyerRepo, authRepo: authRepo, txMgr: txMgr}
}

// Execute creates a new buyer with authentication
func (uc *createBuyerUseCase) Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var createdBuyer *model.Buyer
	err = uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create buyer
		buyer := &model.Buyer{
			Name:         name,
			Organization: organization,
			ContactInfo:  contactInfo,
		}

		buyerResult, err := uc.buyerRepo.Create(txCtx, buyer)
		if err != nil {
			return err
		}
		createdBuyer = buyerResult

		// Create authentication
		auth := &model.Authentication{
			BuyerID:      createdBuyer.ID,
			Email:        email,
			PasswordHash: string(hashedPassword),
			AuthType:     "password",
		}

		_, err = uc.authRepo.Create(txCtx, auth)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdBuyer, nil
}
