package mailhog

import (
	"net/smtp"
	"text/template"

	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
)

type mockTemplateLoader struct {
	realLoader *templates.TemplateLoader
	mockErr    bool
}

func (m *mockTemplateLoader) Get(name string) *template.Template {
	if m.mockErr {
		return nil
	}
	return m.realLoader.Get(name)
}

// setSendMailFunc replaces sendMailFunc for testing.
// Returns a function to restore the original value.
func setSendMailFunc(f func(addr string, a smtp.Auth, from string, to []string, msg []byte) error) func() {
	origAdmin := adminSendMailFunc
	origBuyer := buyerSendMailFunc
	adminSendMailFunc = f
	buyerSendMailFunc = f
	return func() {
		adminSendMailFunc = origAdmin
		buyerSendMailFunc = origBuyer
	}
}
