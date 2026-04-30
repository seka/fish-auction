package mailhog

import "context"

type noopAdminEmailService struct{}

func (n *noopAdminEmailService) SendAdminPasswordReset(_ context.Context, _, _ string) error {
	return nil
}

type noopBuyerEmailService struct{}

func (n *noopBuyerEmailService) SendBuyerPasswordReset(_ context.Context, _, _ string) error {
	return nil
}
