package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateBuyerUseCase defines the interface for creating buyers
type CreateBuyerUseCase interface {
	Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error)
}

// createBuyerUseCase handles the creation of buyers
type createBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
	authRepo  repository.AuthenticationRepository
}

// NewCreateBuyerUseCase creates a new instance of CreateBuyerUseCase
func NewCreateBuyerUseCase(buyerRepo repository.BuyerRepository, authRepo repository.AuthenticationRepository) CreateBuyerUseCase {
	return &createBuyerUseCase{
		buyerRepo: buyerRepo,
		authRepo:  authRepo,
	}
}

// Execute creates a new buyer with authentication
func (uc *createBuyerUseCase) Execute(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create buyer
	buyer := &model.Buyer{
		Name:         name,
		Organization: organization,
		ContactInfo:  contactInfo,
	}

	createdBuyer, err := uc.buyerRepo.Create(ctx, buyer)
	if err != nil {
		return nil, err
	}

	// Create authentication
	auth := &model.Authentication{
		BuyerID:      createdBuyer.ID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		AuthType:     "password",
	}

	_, err = uc.authRepo.Create(ctx, auth)
	if err != nil {
		// TODO: Consider rollback or compensation logic here
		return nil, err
	}

	return createdBuyer, nil
}
