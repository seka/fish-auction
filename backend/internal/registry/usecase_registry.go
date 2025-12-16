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
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

// UseCase defines the interface for creating use cases
type UseCase interface {
	NewCreateItemUseCase() item.CreateItemUseCase
	NewListItemsUseCase() item.ListItemsUseCase
	NewCreateBidUseCase() bid.CreateBidUseCase
	NewCreateBuyerUseCase() buyer.CreateBuyerUseCase
	NewListBuyersUseCase() buyer.ListBuyersUseCase
	NewLoginBuyerUseCase() buyer.LoginBuyerUseCase
	NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase
	NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase
	NewCreateFishermanUseCase() fisherman.CreateFishermanUseCase
	NewListFishermenUseCase() fisherman.ListFishermenUseCase
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
	NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase
	NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase
}

// useCaseRegistry implements the UseCase interface
type useCaseRegistry struct {
	repo Repository
	cfg  *config.Config
}

// NewUseCaseRegistry creates a new UseCase registry
func NewUseCaseRegistry(repo Repository, cfg *config.Config) UseCase {
	return &useCaseRegistry{repo: repo, cfg: cfg}
}

// ... (existing methods omitted for brevity, only changing what's needed or new)

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
		u.repo.NewAuctionRepository(),
		u.repo.NewTransactionManager(),
	)
}

func (u *useCaseRegistry) NewCreateBuyerUseCase() buyer.CreateBuyerUseCase {
	return buyer.NewCreateBuyerUseCase(u.repo.NewBuyerRepository(), u.repo.NewAuthenticationRepository())
}

func (u *useCaseRegistry) NewListBuyersUseCase() buyer.ListBuyersUseCase {
	return buyer.NewListBuyersUseCase(u.repo.NewBuyerRepository())
}

func (u *useCaseRegistry) NewLoginBuyerUseCase() buyer.LoginBuyerUseCase {
	return buyer.NewLoginBuyerUseCase(u.repo.NewBuyerRepository(), u.repo.NewAuthenticationRepository())
}

func (u *useCaseRegistry) NewGetBuyerPurchasesUseCase() buyer.GetBuyerPurchasesUseCase {
	return buyer.NewGetBuyerPurchasesUseCase(u.repo.NewBidRepository())
}

func (u *useCaseRegistry) NewGetBuyerAuctionsUseCase() buyer.GetBuyerAuctionsUseCase {
	return buyer.NewGetBuyerAuctionsUseCase(u.repo.NewBidRepository())
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
	return auction.NewListAuctionsUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewGetAuctionUseCase() auction.GetAuctionUseCase {
	return auction.NewGetAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewGetAuctionItemsUseCase() auction.GetAuctionItemsUseCase {
	return auction.NewGetAuctionItemsUseCase(u.repo.NewItemRepository())
}

func (u *useCaseRegistry) NewUpdateAuctionUseCase() auction.UpdateAuctionUseCase {
	return auction.NewUpdateAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewUpdateAuctionStatusUseCase() auction.UpdateAuctionStatusUseCase {
	return auction.NewUpdateAuctionStatusUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewDeleteAuctionUseCase() auction.DeleteAuctionUseCase {
	return auction.NewDeleteAuctionUseCase(u.repo.NewAuctionRepository())
}

func (u *useCaseRegistry) NewAdminUpdatePasswordUseCase() admin.UpdatePasswordUseCase {
	return admin.NewUpdatePasswordUseCase(u.repo.NewAdminRepository())
}

func (u *useCaseRegistry) NewBuyerUpdatePasswordUseCase() buyer.UpdatePasswordUseCase {
	return buyer.NewUpdatePasswordUseCase(u.repo.NewAuthenticationRepository())
}

func (u *useCaseRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return auth.NewRequestPasswordResetUseCase(
		u.repo.NewBuyerRepository(),
		u.repo.PasswordReset(),
		u.cfg,
	)
}

func (u *useCaseRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return auth.NewResetPasswordUseCase(
		u.repo.PasswordReset(),
		u.repo.NewAuthenticationRepository(),
	)
}

func (u *useCaseRegistry) NewRequestAdminPasswordResetUseCase() admin.RequestPasswordResetUseCase {
	return admin.NewRequestPasswordResetUseCase(
		u.repo.NewAdminRepository(),
		u.repo.AdminPasswordReset(),
		u.cfg,
	)
}

func (u *useCaseRegistry) NewResetAdminPasswordUseCase() admin.ResetPasswordUseCase {
	return admin.NewResetPasswordUseCase(
		u.repo.AdminPasswordReset(),
		u.repo.NewAdminRepository(),
	)
}
