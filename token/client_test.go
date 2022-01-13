package token

import (
	"context"
	"testing"

	"github.com/amzapi/amz-sp-pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	const (
		clientID     = "amzn1.application-oa2-client.xxxxxxxx"
		clientSecret = "xxxxxxxx"
		authCode     = "xxxxxxxx"
	)

	ctx := context.Background()

	t.Run("TestNewClient", func(t *testing.T) {

		client := NewClient(
			WithClientID(clientID),
			WithClientSecret(clientSecret),
		)

		assert.NotNil(t, client)

		t.Run("TestExchange", func(t *testing.T) {
			token, err := client.Exchange(ctx, authCode)
			assert.Nil(t, err)
			assert.NotNil(t, token)
			printToken(t, token)
			t.Run("TestGetAccessToken", func(t *testing.T) {
				token2, err := client.GetAccessToken(ctx, token.RefreshToken)
				assert.Nil(t, err)
				assert.NotNil(t, token2)
				printToken(t, token2)
			})
		})

		t.Run("TestGetGrantLessAccessToken", func(t *testing.T) {
			token, err := client.GetGrantLessAccessToken(ctx, types.ScopeNotificationsApi)
			assert.Nil(t, err)
			assert.NotNil(t, token)
			printToken(t, token)
		})
	})
}

func printToken(t *testing.T, token *types.Token) {
	t.Logf("refresh_token = %s", token.RefreshToken)
	t.Logf("access_token = %s", token.AccessToken)
	t.Logf("token_type = %s", token.TokenType)
	t.Logf("expires_in = %d", token.ExpiresIn)
	t.Logf("expiry = %s", token.Expiry)
}
