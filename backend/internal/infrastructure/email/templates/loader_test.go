package templates_test

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/infrastructure/email/templates"
	"github.com/stretchr/testify/assert"
)

func TestTemplateLoader_Get(t *testing.T) {
	loader, err := templates.NewTemplateLoader()
	assert.NoError(t, err)
	assert.NotNil(t, loader)

	t.Run("GetBuyerPasswordReset", func(t *testing.T) {
		tmpl := loader.Get("buyer_password_reset.txt")
		assert.NotNil(t, tmpl)
		assert.Equal(t, "buyer_password_reset.txt", tmpl.Name())
	})

	t.Run("GetAdminPasswordReset", func(t *testing.T) {
		tmpl := loader.Get("admin_password_reset.txt")
		assert.NotNil(t, tmpl)
		assert.Equal(t, "admin_password_reset.txt", tmpl.Name())
	})

	t.Run("GetUnknown", func(t *testing.T) {
		tmpl := loader.Get("unknown.txt")
		assert.Nil(t, tmpl)
	})
}
