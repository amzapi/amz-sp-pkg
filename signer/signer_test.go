package signer

import (
	"context"
	"net/http"
	"testing"

	"gopkg.me/amz-sp-pkg/cache"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {

	ctx := context.Background()

	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := r.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	t.Run("NewSigner", func(t *testing.T) {

		s := NewSigner(
			WithAccessKeyID("<AWS AccessKeyID>"),
			WithSecretAccessKey("<AWS SecretAccessKey>"),
			WithRoleArn("<AWS RoleArn>"),
			WithDebug(true),
			WithCache(cache.NewRedisCache(r)),
		)

		assert.NotNil(t, s)

		t.Run("RefreshRoleCredentials", func(t *testing.T) {

			roleCredentials, err := s.RefreshRoleCredentials(ctx)

			assert.Nil(t, err)
			assert.NotNil(t, roleCredentials)

			t.Logf("%#v", roleCredentials)

			t.Run("SignRequest", func(t *testing.T) {

				req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
				assert.NotNil(t, req)

				for k, v := range req.Header {
					t.Logf("%s -> %s", k, v)
				}

				err = s.SignRequest(ctx, req, "test access token", "eu-west-1")
				assert.Nil(t, err)

				for k, v := range req.Header {
					t.Logf("%s -> %s", k, v)
				}
			})

		})

	})
}
