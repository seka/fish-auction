package service

import "context"

type EmailService interface {
	SendBuyerPasswordReset(ctx context.Context, to, url string) error
	SendAdminPasswordReset(ctx context.Context, to, url string) error
}
