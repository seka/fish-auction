package service

import "context"

type BuyerEmailService interface {
	SendBuyerPasswordReset(ctx context.Context, to, url string) error
}

type AdminEmailService interface {
	SendAdminPasswordReset(ctx context.Context, to, url string) error
}
