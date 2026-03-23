package service

import "context"

// BuyerEmailService provides BuyerEmailService related functionality.
type BuyerEmailService interface {
	SendBuyerPasswordReset(ctx context.Context, to, url string) error
}

// AdminEmailService provides AdminEmailService related functionality.
type AdminEmailService interface {
	SendAdminPasswordReset(ctx context.Context, to, url string) error
}
