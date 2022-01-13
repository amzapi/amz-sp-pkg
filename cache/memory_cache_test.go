package cache

import (
	"context"
	"strconv"
	"testing"
	"time"

	"gopkg.me/amz-sp-pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {

	ctx := context.Background()
	c := NewMemoryCache()

	t.Parallel()

	t.Run("TestMemoryCacheToken", func(t *testing.T) {
		for i := 0; i < 100000; i++ {
			err := c.SetToken(ctx, "Token"+strconv.Itoa(i), &types.Token{
				AccessToken: "test",
				Expiry:      time.Now().Add(time.Hour),
			}, 5*time.Second)
			assert.NoError(t, err)
		}
		for i := 0; i < 100000; i++ {
			x, found, err := c.GetToken(ctx, "Token"+strconv.Itoa(i))
			assert.NoError(t, err)
			assert.Equal(t, found, true)
			assert.NotNil(t, x)
			assert.Equal(t, x.AccessToken, "test")
		}
	})

	t.Run("TestMemoryCacheRoleCredentials", func(t *testing.T) {
		for i := 0; i < 100000; i++ {
			err := c.SetRoleCredentials(ctx, "RoleCredentials"+strconv.Itoa(i), &types.RoleCredentials{
				AccessKeyID: "test",
				Expiration:  time.Now().Add(time.Hour),
			}, 5*time.Second)
			assert.NoError(t, err)
		}
		for i := 0; i < 100000; i++ {
			x, found, err := c.GetRoleCredentials(ctx, "RoleCredentials"+strconv.Itoa(i))
			assert.NoError(t, err)
			assert.Equal(t, found, true)
			assert.NotNil(t, x)
			assert.Equal(t, x.AccessKeyID, "test")
		}
	})
}
