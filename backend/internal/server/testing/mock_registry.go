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

// MockRegistry is a mock implementation of Registry for testing.
type MockRegistry struct {
	CreateItemUC                item.CreateItemUseCase
	ListItemsUC                 item.ListItemsUseCase
	UpdateItemUC                item.UpdateItemUseCase
	DeleteItemUC                item.DeleteItemUseCase
	UpdateItemSortOrderUC       item.UpdateItemSortOrderUseCase
	ReorderItemsUC              item.ReorderItemsUseCase
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
	GetBuyerUC                  buyer.GetBuyerUseCase
	RequestPasswordResetUC      auth.RequestPasswordResetUseCase
	ResetPasswordUC             auth.ResetPasswordUseCase
	VerifyResetTokenUC          auth.VerifyResetTokenUseCase
	VerifyAdminResetTokenUC     admin.VerifyResetTokenUseCase
	RequestAdminPasswordResetUC admin.RequestPasswordResetUseCase
	ResetAdminPasswordUC        admin.ResetPasswordUseCase
	DeleteFishermanUC           fisherman.DeleteFishermanUseCase
	DeleteBuyerUC               buyer.DeleteBuyerUseCase
	SubscribeNotificationUC     notification.SubscribeNotificationUseCase
	PublishNotificationUC       notification.PublishNotificationUseCase
	CreateAdminUC               admin.CreateAdminUseCase
}

// NewItemRepository creates a new ItemRepository instance.
func (m *MockRegistry) NewItemRepository() repository.ItemRepository {
	return nil
}

// NewBidRepository creates a new BidRepository instance.
func (m *MockRegistry) NewBidRepository() repository.BidRepository {
	return nil
}

// NewBuyerRepository creates a new BuyerRepository instance.
func (m *MockRegistry) NewBuyerRepository() repository.BuyerRepository {
	return nil
}

// NewFishermanRepository creates a new FishermanRepository instance.
func (m *MockRegistry) NewFishermanRepository() repository.FishermanRepository {
	return nil
}

// NewTransactionManager creates a new TransactionManager instance.
func (m *MockRegistry) NewTransactionManager() repository.TransactionManager {
	return nil
}

// NewCreateItemUseCase creates a new CreateItemUseCase instance.
func (m *MockRegistry) NewCreateItemUseCase() item.CreateItemUseCase {
	return m.CreateItemUC
}

// NewListItemsUseCase creates a new ListItemsUseCase instance.
func (m *MockRegistry) NewListItemsUseCase() item.ListItemsUseCase {
	return m.ListItemsUC
}

// NewUpdateItemUseCase creates a new UpdateItemUseCase instance.
func (m *MockRegistry) NewUpdateItemUseCase() item.UpdateItemUseCase {
	return m.UpdateItemUC
}

// NewDeleteItemUseCase creates a new DeleteItemUseCase instance.
func (m *MockRegistry) NewDeleteItemUseCase() item.DeleteItemUseCase {
	return m.DeleteItemUC
}

// NewUpdateItemSortOrderUseCase creates a new UpdateItemSortOrderUseCase instance.
func (m *MockRegistry) NewUpdateItemSortOrderUseCase() item.UpdateItemSortOrderUseCase {
	return m.UpdateItemSortOrderUC
}

// NewReorderItemsUseCase creates a new ReorderItemsUseCase instance.
func (m *MockRegistry) NewReorderItemsUseCase() item.ReorderItemsUseCase {
	return m.ReorderItemsUC
}

// NewCreateBidUseCase creates a new CreateBidUseCase instance.
func (m *MockRegistry) NewCreateBidUseCase() bid.CreateBidUseCase {
	return m.CreateBidUC
}

// NewCreateBuyerUseCase creates a new CreateBuyerUseCase instance.
func (m *MockRegistry) NewCreateBuyerUseCase() buyer.CreateBuyerUseCase {
	return m.CreateBuyerUC
}

// NewListBuyersUseCase creates a new ListBuyersUseCase instance.
func (m *MockRegistry) NewListBuyersUseCase() buyer.ListBuyersUseCase {
	return m.ListBuyersUC
}

// NewCreateFishermanUseCase creates a new CreateFishermanUseCase instance.
func (m *MockRegistry) NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase {
	return m.CreateFishermanUC
}

// NewListFishermenUseCase creates a new ListFishermenUseCase instance.
func (m *MockRegistry) NewListFishermenUseCase() fisherman.ListFishermenUseCase {
	return m.ListFishermenUC
}

// NewListInvoicesUseCase creates a new ListInvoicesUseCase instance.
func (m *MockRegistry) NewListInvoicesUseCase() invoice.ListInvoicesUseCase {
	return m.ListInvoicesUC
}

// NewLoginUseCase creates a new LoginUseCase instance.
func (m *MockRegistry) NewLoginUseCase() auth.LoginUseCase {
	return m.LoginUC
}

// NewLoginBuyerUseCase creates a new LoginBuyerUseCase instance.
func (m *MockRegistry) NewLoginBuyerUseCase() buyer.LoginBuyerUseCase {
	return m.LoginBuyerUC
}

// NewGetBuyerPurchasesUseCase creates a new GetBuyerPurchasesUseCase instance.
func (m *MockRegistry) NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase {
	return m.GetBuyerPurchasesUC
}

// NewGetBuyerAuctionsUseCase creates a new GetBuyerAuctionsUseCase instance.
func (m *MockRegistry) NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase {
	return m.GetBuyerAuctionsUC
}

// NewGetBuyerUseCase creates a new GetBuyerUseCase instance.
func (m *MockRegistry) NewGetBuyerUseCase() buyer.GetBuyerUseCase {
	return m.GetBuyerUC
}

// NewCreateVenueUseCase creates a new CreateVenueUseCase instance.
func (m *MockRegistry) NewCreateVenueUseCase() venue.CreateVenueUseCase {
	return m.CreateVenueUC
}

// NewListVenuesUseCase creates a new ListVenuesUseCase instance.
func (m *MockRegistry) NewListVenuesUseCase() venue.ListVenuesUseCase {
	return m.ListVenuesUC
}

// NewGetVenueUseCase creates a new GetVenueUseCase instance.
func (m *MockRegistry) NewGetVenueUseCase() venue.GetVenueUseCase {
	return m.GetVenueUC
}

// NewUpdateVenueUseCase creates a new UpdateVenueUseCase instance.
func (m *MockRegistry) NewUpdateVenueUseCase() venue.UpdateVenueUseCase {
	return m.UpdateVenueUC
}

// NewDeleteVenueUseCase creates a new DeleteVenueUseCase instance.
func (m *MockRegistry) NewDeleteVenueUseCase() venue.DeleteVenueUseCase {
	return m.DeleteVenueUC
}

// NewCreateAuctionUseCase creates a new CreateAuctionUseCase instance.
func (m *MockRegistry) NewCreateAuctionUseCase() auction.CreateAuctionUseCase {
	return m.CreateAuctionUC
}

// NewListAuctionsUseCase creates a new ListAuctionsUseCase instance.
func (m *MockRegistry) NewListAuctionsUseCase() auction.ListAuctionsUseCase {
	return m.ListAuctionsUC
}

// NewGetAuctionUseCase creates a new GetAuctionUseCase instance.
func (m *MockRegistry) NewGetAuctionUseCase() auction.GetAuctionUseCase {
	return m.GetAuctionUC
}

// NewGetAuctionItemsUseCase creates a new GetAuctionItemsUseCase instance.
func (m *MockRegistry) NewGetAuctionItemsUseCase() auction.GetAuctionItemsUseCase {
	return m.GetAuctionItemsUC
}

// NewUpdateAuctionUseCase creates a new UpdateAuctionUseCase instance.
func (m *MockRegistry) NewUpdateAuctionUseCase() auction.UpdateAuctionUseCase {
	return m.UpdateAuctionUC
}

// NewUpdateAuctionStatusUseCase creates a new UpdateAuctionStatusUseCase instance.
func (m *MockRegistry) NewUpdateAuctionStatusUseCase() auction.UpdateAuctionStatusUseCase {
	return m.UpdateAuctionStatusUC
}

// NewDeleteAuctionUseCase creates a new DeleteAuctionUseCase instance.
func (m *MockRegistry) NewDeleteAuctionUseCase() auction.DeleteAuctionUseCase {
	return m.DeleteAuctionUC
}

// NewAdminUpdatePasswordUseCase creates a new AdminUpdatePasswordUseCase instance.
func (m *MockRegistry) NewAdminUpdatePasswordUseCase() admin.UpdatePasswordUseCase {
	return m.UpdateAdminPasswordUC
}

// NewBuyerUpdatePasswordUseCase creates a new BuyerUpdatePasswordUseCase instance.
func (m *MockRegistry) NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase {
	return m.UpdateBuyerPasswordUC
}

// NewRequestPasswordResetUseCase creates a new RequestPasswordResetUseCase instance.
func (m *MockRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return m.RequestPasswordResetUC
}

// NewResetPasswordUseCase creates a new ResetPasswordUseCase instance.
func (m *MockRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return m.ResetPasswordUC
}

// NewVerifyResetTokenUseCase creates a new VerifyResetTokenUseCase instance.
func (m *MockRegistry) NewVerifyResetTokenUseCase() auth.VerifyResetTokenUseCase {
	return m.VerifyResetTokenUC
}

// NewRequestAdminPasswordResetUseCase creates a new RequestAdminPasswordResetUseCase instance.
func (m *MockRegistry) NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase {
	return m.RequestAdminPasswordResetUC
}

// NewVerifyAdminResetTokenUseCase creates a new VerifyAdminResetTokenUseCase instance.
func (m *MockRegistry) NewVerifyAdminResetTokenUseCase() admin.VerifyResetTokenUseCase {
	return m.VerifyAdminResetTokenUC
}

// NewResetAdminPasswordUseCase creates a new ResetAdminPasswordUseCase instance.
func (m *MockRegistry) NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase {
	return m.ResetAdminPasswordUC
}

// NewDeleteFishermanUseCase creates a new DeleteFishermanUseCase instance.
func (m *MockRegistry) NewDeleteFishermanUseCase() fisherman.DeleteFishermanUseCase {
	return m.DeleteFishermanUC
}

// NewDeleteBuyerUseCase creates a new DeleteBuyerUseCase instance.
func (m *MockRegistry) NewDeleteBuyerUseCase() buyer.DeleteBuyerUseCase {
	return m.DeleteBuyerUC
}

// NewSubscribeNotificationUseCase creates a new SubscribeNotificationUseCase instance.
func (m *MockRegistry) NewSubscribeNotificationUseCase() notification.SubscribeNotificationUseCase {
	return m.SubscribeNotificationUC
}

// NewPublishNotificationUseCase creates a new PublishNotificationUseCase instance.
func (m *MockRegistry) NewPublishNotificationUseCase() notification.PublishNotificationUseCase {
	return m.PublishNotificationUC
}

// NewCreateAdminUseCase creates a new CreateAdminUseCase instance.
func (m *MockRegistry) NewCreateAdminUseCase() admin.CreateAdminUseCase {
	return m.CreateAdminUC
}

// Ensure MockRegistry implements registry.UseCase
var _ registry.UseCase = &MockRegistry{}
