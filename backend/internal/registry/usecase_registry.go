package registry

import (
	"github.com/seka/fish-auction/backend/config"
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

// UseCase defines the interface for creating use cases
type UseCase interface {
	NewCreateItemUseCase() item.CreateItemUseCase
	NewListItemsUseCase() item.ListItemsUseCase
	NewUpdateItemUseCase() item.UpdateItemUseCase
	NewDeleteItemUseCase() item.DeleteItemUseCase
	NewUpdateItemSortOrderUseCase() item.UpdateItemSortOrderUseCase
	NewReorderItemsUseCase() item.ReorderItemsUseCase
	NewCreateBidUseCase() bid.CreateBidUseCase
	NewCreateBuyerUseCase() buyer.CreateBuyerUseCase
	NewListBuyersUseCase() buyer.ListBuyersUseCase
	NewLoginBuyerUseCase() buyer.LoginBuyerUseCase
	NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase
	NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase
	NewGetBuyerUseCase() buyer.GetBuyerUseCase
	NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase
	NewListFishermenUseCase() fisherman.ListFishermenUseCase
	NewDeleteFishermanUseCase() fisherman.DeleteFishermanUseCase
	NewDeleteBuyerUseCase() buyer.DeleteBuyerUseCase
	NewListInvoicesUseCase() invoice.ListInvoicesUseCase
	NewLoginUseCase() auth.LoginUseCase
	NewCreateVenueUseCase() venue.CreateVenueUseCase
	NewListVenuesUseCase() venue.ListVenuesUseCase
	NewGetVenueUseCase() venue.GetVenueUseCase
	NewUpdateVenueUseCase() venue.UpdateVenueUseCase
	NewDeleteVenueUseCase() venue.DeleteVenueUseCase
	NewCreateAuctionUseCase() auction.CreateAuctionUseCase
	NewListAuctionsUseCase() auction.ListAuctionsUseCase
	NewGetAuctionUseCase() auction.GetAuctionUseCase
	NewGetAuctionItemsUseCase() auction.GetAuctionItemsUseCase
	NewUpdateAuctionUseCase() auction.UpdateAuctionUseCase
	NewUpdateAuctionStatusUseCase() auction.UpdateAuctionStatusUseCase
	NewDeleteAuctionUseCase() auction.DeleteAuctionUseCase
	NewAdminUpdatePasswordUseCase() admin.UpdatePasswordUseCase
	NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase
	NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase
	NewResetPasswordUseCase() auth.ResetPasswordUseCase
	NewVerifyResetTokenUseCase() auth.VerifyResetTokenUseCase
	NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase
	NewVerifyAdminResetTokenUseCase() admin.VerifyResetTokenUseCase
	NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase
	NewSubscribeNotificationUseCase() notification.SubscribeNotificationUseCase
	NewPublishNotificationUseCase() notification.PublishNotificationUseCase
}

type useCaseRegistry struct {
	repo    Repository
	service Service
	cfg     *config.Config
}

// NewUseCaseRegistry creates a new UseCase registry
func NewUseCaseRegistry(repo Repository, service Service, cfg *config.Config) UseCase {
	return &useCaseRegistry{
		repo:    repo,
		service: service,
		cfg:     cfg,
	}
}

// ... (methods)

func (u *useCaseRegistry) NewCreateItemUseCase() item.CreateItemUseCase {
	return item.NewCreateItemUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewListItemsUseCase() item.ListItemsUseCase {
	return item.NewListItemsUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewUpdateItemUseCase() item.UpdateItemUseCase {
	return item.NewUpdateItemUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewDeleteItemUseCase() item.DeleteItemUseCase {
	return item.NewDeleteItemUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewUpdateItemSortOrderUseCase() item.UpdateItemSortOrderUseCase {
	return item.NewUpdateItemSortOrderUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewReorderItemsUseCase() item.ReorderItemsUseCase {
	return item.NewReorderItemsUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewCreateBidUseCase() bid.CreateBidUseCase {
	return bid.NewCreateBidUseCase(
		u.repo.NewItemRepository(),
		u.repo.NewBuyerRepository(),
		u.repo.NewBidRepository(),
		u.repo.NewAuctionRepository(),
		u.NewPublishNotificationUseCase(),
		u.repo.NewTransactionManager(),
		u.repo.NewItemCacheInvalidator(),
		u.service.NewClock(),
	)
}

func (u *useCaseRegistry) NewCreateBuyerUseCase() buyer.CreateBuyerUseCase {
	return buyer.NewCreateBuyerUseCase(u.repo.NewBuyerRepository(), u.repo.NewAuthenticationRepository(), u.repo.NewTransactionManager())
}

func (u *useCaseRegistry) NewListBuyersUseCase() buyer.ListBuyersUseCase {
	return buyer.NewListBuyersUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewLoginBuyerUseCase() buyer.LoginBuyerUseCase {
	return buyer.NewLoginBuyerUseCase(u.repo.NewBuyerRepository(), u.repo.NewAuthenticationRepository(), u.service.NewClock())
}

func (u *useCaseRegistry) NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase {
	return buyer.NewGetBuyerPurchasesUseCase(u.repo.NewBidRepository())
}

func (u *useCaseRegistry) NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase {
	return buyer.NewGetBuyerAuctionsUseCase(u.repo.NewBidRepository())
}

func (u *useCaseRegistry) NewGetBuyerUseCase() buyer.GetBuyerUseCase {
	return buyer.NewGetBuyerUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase {
	return fisherman.NewCreateFishermanUseCase(u.repo.NewFishermanRepository())
}

func (u *useCaseRegistry) NewListFishermenUseCase() fisherman.ListFishermenUseCase {
	return fisherman.NewListFishermenUseCase(u.repo.NewFishermanRepository())
}

func (u *useCaseRegistry) NewDeleteFishermanUseCase() fisherman.DeleteFishermanUseCase {
	return fisherman.NewDeleteFishermanUseCase(u.repo.NewFishermanRepository())
}

func (u *useCaseRegistry) NewDeleteBuyerUseCase() buyer.DeleteBuyerUseCase {
	return buyer.NewDeleteBuyerUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewListInvoicesUseCase() invoice.ListInvoicesUseCase {
	return invoice.NewListInvoicesUseCase(u.repo.NewBidRepository())
}

func (u *useCaseRegistry) NewLoginUseCase() auth.LoginUseCase {
	return auth.NewLoginUseCase(u.repo.NewAdminRepository())
}

func (u *useCaseRegistry) NewCreateVenueUseCase() venue.CreateVenueUseCase {
	return venue.NewCreateVenueUseCase(u.repo.NewVenueRepository())
}

func (u *useCaseRegistry) NewListVenuesUseCase() venue.ListVenuesUseCase {
	return venue.NewListVenuesUseCase(u.repo.NewVenueRepository())
}

func (u *useCaseRegistry) NewGetVenueUseCase() venue.GetVenueUseCase {
	return venue.NewGetVenueUseCase(u.repo.NewVenueRepository())
}

func (u *useCaseRegistry) NewUpdateVenueUseCase() venue.UpdateVenueUseCase {
	return venue.NewUpdateVenueUseCase(u.repo.NewVenueRepository())
}

func (u *useCaseRegistry) NewDeleteVenueUseCase() venue.DeleteVenueUseCase {
	return venue.NewDeleteVenueUseCase(u.repo.NewVenueRepository())
}

func (u *useCaseRegistry) NewCreateAuctionUseCase() auction.CreateAuctionUseCase {
	return auction.NewCreateAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewListAuctionsUseCase() auction.ListAuctionsUseCase {
	return auction.NewListAuctionsUseCase(u.repo.NewAuctionRepository(), u.service.NewClock())
}

func (u *useCaseRegistry) NewGetAuctionUseCase() auction.GetAuctionUseCase {
	return auction.NewGetAuctionUseCase(u.repo.NewAuctionRepository(), u.service.NewClock())
}

func (u *useCaseRegistry) NewGetAuctionItemsUseCase() auction.GetAuctionItemsUseCase {
	return auction.NewGetAuctionItemsUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewUpdateAuctionUseCase() auction.UpdateAuctionUseCase {
	return auction.NewUpdateAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewUpdateAuctionStatusUseCase() auction.UpdateAuctionStatusUseCase {
	return auction.NewUpdateAuctionStatusUseCase(
		u.repo.NewAuctionRepository(),
		u.repo.NewBuyerRepository(),
		u.NewPublishNotificationUseCase(),
	)
}

func (u *useCaseRegistry) NewDeleteAuctionUseCase() auction.DeleteAuctionUseCase {
	return auction.NewDeleteAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewAdminUpdatePasswordUseCase() admin.UpdatePasswordUseCase {
	return admin.NewUpdatePasswordUseCase(u.repo.NewAdminRepository(), u.repo.NewSessionRepository())
}

func (u *useCaseRegistry) NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase {
	return buyer.NewUpdatePasswordUseCase(u.repo.NewAuthenticationRepository(), u.repo.NewSessionRepository())
}

func (u *useCaseRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return auth.NewRequestPasswordResetUseCase(
		u.repo.NewBuyerRepository(),
		u.repo.PasswordReset(),
		u.service.NewBuyerEmailService(),
		u.cfg.FrontendURL,
		u.repo.NewTransactionManager(),
		u.service.NewClock(),
	)
}

func (u *useCaseRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return auth.NewResetPasswordUseCase(
		u.repo.PasswordReset(),
		u.repo.NewAuthenticationRepository(),
		u.repo.NewTransactionManager(),
		u.service.NewClock(),
	)
}

func (u *useCaseRegistry) NewVerifyResetTokenUseCase() auth.VerifyResetTokenUseCase {
	return auth.NewVerifyResetTokenUseCase(u.repo.PasswordReset(), u.service.NewClock())
}

func (u *useCaseRegistry) NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase {
	return admin.NewRequestPasswordResetUseCase(
		u.repo.NewAdminRepository(),
		u.repo.PasswordReset(),
		u.service.NewAdminEmailService(),
		u.cfg.FrontendURL,
		u.repo.NewTransactionManager(),
		u.service.NewClock(),
	)
}

func (u *useCaseRegistry) NewVerifyAdminResetTokenUseCase() admin.VerifyResetTokenUseCase {
	return admin.NewVerifyResetTokenUseCase(u.repo.NewAdminRepository(), u.repo.PasswordReset(), u.service.NewClock())
}

func (u *useCaseRegistry) NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase {
	return admin.NewResetPasswordUseCase(
		u.repo.PasswordReset(),
		u.repo.NewAdminRepository(),
		u.repo.NewTransactionManager(),
		u.service.NewClock(),
	)
}

func (u *useCaseRegistry) NewSubscribeNotificationUseCase() notification.SubscribeNotificationUseCase {
	return notification.NewSubscribeNotificationUseCase(u.repo.NewPushRepository())
}

func (u *useCaseRegistry) NewPublishNotificationUseCase() notification.PublishNotificationUseCase {
	return notification.NewPublishNotificationUseCase(
		u.repo.NewPushRepository(),
		u.service.NewPushNotificationService(),
	)
}
