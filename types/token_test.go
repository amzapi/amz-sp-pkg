package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenUnmarshalJSON(t *testing.T) {
	t.Run("TestTokenUnmarshalJSON", func(t *testing.T) {
		token := Token{
			AccessToken:  "test1",
			TokenType:    "test2",
			RefreshToken: "test3",
			ExpiresIn:    3600,
		}
		t.Run("Marshal", func(t *testing.T) {
			testdata, err := json.Marshal(token)
			assert.Nil(t, err)
			assert.NotNil(t, testdata)
			assert.Zero(t, token.Expiry)
			t.Run("Unmarshal", func(t *testing.T) {
				err = json.Unmarshal(testdata, &token)
				assert.Nil(t, err)
				assert.NotZero(t, token.Expiry)
				t.Logf("%+v", token)
			})
		})
	})
}
