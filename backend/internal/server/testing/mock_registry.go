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
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

type MockRegistry struct {
	CreateItemUC                item.CreateItemUseCase
	ListItemsUC                 item.ListItemsUseCase
	CreateBidUC                 bid.CreateBidUseCase
	CreateBuyerUC               buyer.CreateBuyerUseCase
	ListBuyersUC                buyer.ListBuyersUseCase
	CreateFishermanUC           fisherman.CreateFishermanUseCase
	ListFishermenUC             fisherman.ListFishermenUseCase
	ListInvoicesUC              invoice.ListInvoicesUseCase
	LoginUC                     auth.LoginUseCase
	CreateVenueUC               venue.CreateVenueUseCase
	ListVenuesUC                venue.ListVenuesUseCase
	GetVenueUC                  venue.GetVenueUseCase
	UpdateVenueUC               venue.UpdateVenueUseCase
	DeleteVenueUC               venue.DeleteVenueUseCase
	CreateAuctionUC             auction.CreateAuctionUseCase
	ListAuctionsUC              auction.ListAuctionsUseCase
	GetAuctionUC                auction.GetAuctionUseCase
	GetAuctionItemsUC           auction.GetAuctionItemsUseCase
	UpdateAuctionUC             auction.UpdateAuctionUseCase
	UpdateAuctionStatusUC       auction.UpdateAuctionStatusUseCase
	DeleteAuctionUC             auction.DeleteAuctionUseCase
	LoginBuyerUC                buyer.LoginBuyerUseCase
	GetBuyerPurchasesUC         buyer.GetBuyerPurchasesUseCase
	GetBuyerAuctionsUC          buyer.GetBuyerAuctionsUseCase
	UpdateBuyerPasswordUC       buyer.UpdatePasswordUseCase
	UpdateAdminPasswordUC       admin.UpdatePasswordUseCase
	RequestPasswordResetUC      auth.RequestPasswordResetUseCase
	ResetPasswordUC             auth.ResetPasswordUseCase
	RequestAdminPasswordResetUC admin.RequestPasswordResetUseCase
	ResetAdminPasswordUC        admin.ResetPasswordUseCase
	PushNotificationUC          notification.PushNotificationUseCase
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
	return m.LoginBuyerUC
}

func (m *MockRegistry) NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase {
	return m.GetBuyerPurchasesUC
}

func (m *MockRegistry) NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase {
	return m.GetBuyerAuctionsUC
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
	return m.UpdateAdminPasswordUC
}

func (m *MockRegistry) NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase {
	return m.UpdateBuyerPasswordUC
}

func (m *MockRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return m.RequestPasswordResetUC
}

func (m *MockRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return m.ResetPasswordUC
}

func (m *MockRegistry) NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase {
	return m.RequestAdminPasswordResetUC
}

func (m *MockRegistry) NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase {
	return m.ResetAdminPasswordUC
}

func (m *MockRegistry) NewPushNotificationUseCase() notification.PushNotificationUseCase {
	return m.PushNotificationUC
}

// Ensure MockRegistry implements registry.UseCase
var _ registry.UseCase = &MockRegistry{}
