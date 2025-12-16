package testing

import (
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type MockRegistry struct {
	CreateItemUC          item.CreateItemUseCase
	ListItemsUC           item.ListItemsUseCase
	CreateBidUC           bid.CreateBidUseCase
	CreateBuyerUC         buyer.CreateBuyerUseCase
	ListBuyersUC          buyer.ListBuyersUseCase
	CreateFishermanUC     fisherman.CreateFishermanUseCase
	ListFishermenUC       fisherman.ListFishermenUseCase
	ListInvoicesUC        invoice.ListInvoicesUseCase
	LoginUC               auth.LoginUseCase
	CreateVenueUC         venue.CreateVenueUseCase
	ListVenuesUC          venue.ListVenuesUseCase
	GetVenueUC            venue.GetVenueUseCase
	UpdateVenueUC         venue.UpdateVenueUseCase
	DeleteVenueUC         venue.DeleteVenueUseCase
	CreateAuctionUC       auction.CreateAuctionUseCase
	ListAuctionsUC        auction.ListAuctionsUseCase
	GetAuctionUC          auction.GetAuctionUseCase
	GetAuctionItemsUC     auction.GetAuctionItemsUseCase
	UpdateAuctionUC       auction.UpdateAuctionUseCase
	UpdateAuctionStatusUC auction.UpdateAuctionStatusUseCase
	DeleteAuctionUC       auction.DeleteAuctionUseCase
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

func (m *MockRegistry) NewLoginBuyerUseCase() buyer.LoginBuyerUseCase {
	return nil // Or add a field to MockRegistry if needed for tests
}

func (m *MockRegistry) NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase {
	return nil
}

func (m *MockRegistry) NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase {
	return nil
}

func (m *MockRegistry) NewCreateVenueUseCase() venue.CreateVenueUseCase {
	return m.CreateVenueUC
}

func (m *MockRegistry) NewListVenuesUseCase() venue.ListVenuesUseCase {
	return m.ListVenuesUC
}

func (m *MockRegistry) NewGetVenueUseCase() venue.GetVenueUseCase {
	return m.GetVenueUC
}

func (m *MockRegistry) NewUpdateVenueUseCase() venue.UpdateVenueUseCase {
	return m.UpdateVenueUC
}

func (m *MockRegistry) NewDeleteVenueUseCase() venue.DeleteVenueUseCase {
	return m.DeleteVenueUC
}

func (m *MockRegistry) NewCreateAuctionUseCase() auction.CreateAuctionUseCase {
	return m.CreateAuctionUC
}

func (m *MockRegistry) NewListAuctionsUseCase() auction.ListAuctionsUseCase {
	return m.ListAuctionsUC
}

func (m *MockRegistry) NewGetAuctionUseCase() auction.GetAuctionUseCase {
	return m.GetAuctionUC
}

func (m *MockRegistry) NewGetAuctionItemsUseCase() auction.GetAuctionItemsUseCase {
	return m.GetAuctionItemsUC
}

func (m *MockRegistry) NewUpdateAuctionUseCase() auction.UpdateAuctionUseCase {
	return m.UpdateAuctionUC
}

func (m *MockRegistry) NewUpdateAuctionStatusUseCase() auction.UpdateAuctionStatusUseCase {
	return m.UpdateAuctionStatusUC
}

func (m *MockRegistry) NewDeleteAuctionUseCase() auction.DeleteAuctionUseCase {
	return m.DeleteAuctionUC
}

func (m *MockRegistry) NewAdminUpdatePasswordUseCase() admin.UpdatePasswordUseCase {
	return nil
}

func (m *MockRegistry) NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase {
	return nil
}

func (m *MockRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return nil
}

func (m *MockRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return nil
}

func (m *MockRegistry) NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase {
	return nil
}

func (m *MockRegistry) NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase {
	return nil
}

// Ensure MockRegistry implements registry.UseCase
var _ registry.UseCase = &MockRegistry{}
