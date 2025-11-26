package testing

import (
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

type MockRegistry struct {
	CreateItemUC      item.CreateItemUseCase
	ListItemsUC       item.ListItemsUseCase
	CreateBidUC       bid.CreateBidUseCase
	CreateBuyerUC     buyer.CreateBuyerUseCase
	ListBuyersUC      buyer.ListBuyersUseCase
	CreateFishermanUC fisherman.CreateFishermanUseCase
	ListFishermenUC   fisherman.ListFishermenUseCase
	ListInvoicesUC    invoice.ListInvoicesUseCase
	LoginUC           auth.LoginUseCase
}

func (m *MockRegistry) NewItemRepository() repository.ItemRepository {
	return nil
}

func (m *MockRegistry) NewBidRepository() repository.BidRepository {
	return nil
}

func (m *MockRegistry) NewBuyerRepository() repository.BuyerRepository {
	return nil
}

func (m *MockRegistry) NewFishermanRepository() repository.FishermanRepository {
	return nil
}

func (m *MockRegistry) NewTransactionManager() repository.TransactionManager {
	return nil
}

func (m *MockRegistry) NewCreateItemUseCase() item.CreateItemUseCase {
	return m.CreateItemUC
}

func (m *MockRegistry) NewListItemsUseCase() item.ListItemsUseCase {
	return m.ListItemsUC
}

func (m *MockRegistry) NewCreateBidUseCase() bid.CreateBidUseCase {
	return m.CreateBidUC
}

func (m *MockRegistry) NewCreateBuyerUseCase() buyer.CreateBuyerUseCase {
	return m.CreateBuyerUC
}

func (m *MockRegistry) NewListBuyersUseCase() buyer.ListBuyersUseCase {
	return m.ListBuyersUC
}

func (m *MockRegistry) NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase {
	return m.CreateFishermanUC
}

func (m *MockRegistry) NewListFishermenUseCase() fisherman.ListFishermenUseCase {
	return m.ListFishermenUC
}

func (m *MockRegistry) NewListInvoicesUseCase() invoice.ListInvoicesUseCase {
	return m.ListInvoicesUC
}

func (m *MockRegistry) NewLoginUseCase() auth.LoginUseCase {
	return m.LoginUC
}

// Ensure MockRegistry implements registry.UseCase
var _ registry.UseCase = &MockRegistry{}
