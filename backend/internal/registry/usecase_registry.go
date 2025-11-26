package registry

import (
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// UseCase defines the interface for creating use cases
type UseCase interface {
	NewCreateItemUseCase() item.CreateItemUseCase
	NewListItemsUseCase() item.ListItemsUseCase
	NewCreateBidUseCase() bid.CreateBidUseCase
	NewCreateBuyerUseCase() buyer.CreateBuyerUseCase
	NewListBuyersUseCase() buyer.ListBuyersUseCase
	NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase
	NewListFishermenUseCase() fisherman.ListFishermenUseCase
	NewListInvoicesUseCase() invoice.ListInvoicesUseCase
	NewLoginUseCase() auth.LoginUseCase
}

// useCaseRegistry implements the UseCase interface
type useCaseRegistry struct {
	repo Repository
}

// NewUseCaseRegistry creates a new UseCase registry
func NewUseCaseRegistry(repo Repository) UseCase {
	return &useCaseRegistry{repo: repo}
}

func (u *useCaseRegistry) NewCreateItemUseCase() item.CreateItemUseCase {
	return item.NewCreateItemUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewListItemsUseCase() item.ListItemsUseCase {
	return item.NewListItemsUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewCreateBidUseCase() bid.CreateBidUseCase {
	return bid.NewCreateBidUseCase(
		u.repo.NewItemRepository(),
		u.repo.NewBidRepository(),
		u.repo.NewTransactionManager(),
	)
}

func (u *useCaseRegistry) NewCreateBuyerUseCase() buyer.CreateBuyerUseCase {
	return buyer.NewCreateBuyerUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewListBuyersUseCase() buyer.ListBuyersUseCase {
	return buyer.NewListBuyersUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase {
	return fisherman.NewCreateFishermanUseCase(u.repo.NewFishermanRepository())
}

func (u *useCaseRegistry) NewListFishermenUseCase() fisherman.ListFishermenUseCase {
	return fisherman.NewListFishermenUseCase(u.repo.NewFishermanRepository())
}

func (u *useCaseRegistry) NewListInvoicesUseCase() invoice.ListInvoicesUseCase {
	return invoice.NewListInvoicesUseCase(u.repo.NewBidRepository())
}

func (u *useCaseRegistry) NewLoginUseCase() auth.LoginUseCase {
	return auth.NewLoginUseCase()
}
