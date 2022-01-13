package cache

import (
	"context"
	"strconv"
	"testing"
	"time"

	"gopkg.me/amz-sp-pkg/types"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {

	ctx := context.Background()

	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := r.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	c := NewRedisCache(r)

	t.Parallel()

	t.Run("TestRedisCacheToken", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			err := c.SetToken(ctx, "Token"+strconv.Itoa(i), &types.Token{
				AccessToken: "test",
				Expiry:      time.Now().Add(time.Hour),
			}, 5*time.Second)
			assert.NoError(t, err)
		}
		for i := 0; i < 100; i++ {
			x, found, err := c.GetToken(ctx, "Token"+strconv.Itoa(i))
			assert.NoError(t, err)
			assert.Equal(t, found, true)
			assert.NotNil(t, x)
			assert.Equal(t, x.AccessToken, "test")
		}
	})

	t.Run("TestRedisCacheRoleCredentials", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			err := c.SetRoleCredentials(ctx, "RoleCredentials"+strconv.Itoa(i), &types.RoleCredentials{
				AccessKeyID: "test",
				Expiration:  time.Now().Add(time.Hour),
			}, 5*time.Second)
			assert.NoError(t, err)
		}
		for i := 0; i < 100; i++ {
			x, found, err := c.GetRoleCredentials(ctx, "RoleCredentials"+strconv.Itoa(i))
			assert.NoError(t, err)
			assert.Equal(t, found, true)
			assert.NotNil(t, x)
			assert.Equal(t, x.AccessKeyID, "test")
		}
	})

}
