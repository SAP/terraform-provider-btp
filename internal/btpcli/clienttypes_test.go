package btpcli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRequest(t *testing.T) {
	t.Run("NewLoginRequest(...) doesn't set default idp", func(t *testing.T) {
		uut := NewLoginRequest("", "", "")
		assert.Empty(t, uut.IdentityProvider)
		assert.Empty(t, uut.GlobalAccountSubdomain)
		assert.Empty(t, uut.Username)
		assert.Empty(t, uut.Password)
	})
	t.Run("NewLoginRequest(...) uses all given values", func(t *testing.T) {
		uut := NewLoginRequest("my-subdomain", "my-user", "my-pass")
		assert.Empty(t, uut.IdentityProvider)
		assert.Equal(t, "my-subdomain", uut.GlobalAccountSubdomain)
		assert.Equal(t, "my-user", uut.Username)
		assert.Equal(t, "my-pass", uut.Password)
	})
	t.Run("NewLoginRequestWithCustomIDP(...) respects custom idp", func(t *testing.T) {
		uut := NewLoginRequestWithCustomIDP("my-idp", "", "", "")
		assert.Equal(t, "my-idp", uut.IdentityProvider)
	})
	t.Run("LoginRequest can be marshalled", func(t *testing.T) {
		uut := NewLoginRequestWithCustomIDP("my-idp", "my-subdomain", "my-user", "my-pass")

		b, err := json.Marshal(uut)

		if assert.NoError(t, err) {
			assert.Equal(t, `{"customIdp":"my-idp","subdomain":"my-subdomain","userName":"my-user","password":"my-pass"}`, string(b))
		}
	})
}
