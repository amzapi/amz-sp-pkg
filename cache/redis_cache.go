package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/amzapi/amz-sp-pkg/types"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisCache struct {
	cache *cache.Cache
}

func NewRedisCache(rdb redis.UniversalClient) Cache {
	return &RedisCache{
		cache: cache.New(&cache.Options{
			Redis:     rdb,
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		}),
	}
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}

func (r *RedisCache) SetToken(ctx context.Context, key string, token *types.Token, expiration time.Duration) error {
	if token.Expiry.Round(0).Add(-expiryDelta).Before(time.Now()) {
		return errors.New("token expired")
	}
	err := r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: token,
		TTL:   expiration,
	})
	if err != nil {
		return errors.Errorf("cache set error:%v", err)
	}
	return nil
}

func (r *RedisCache) GetToken(ctx context.Context, key string) (*types.Token, bool, error) {
	var t *types.Token
	err := r.cache.Get(ctx, key, &t)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return t, true, nil
}

func (r *RedisCache) SetRoleCredentials(ctx context.Context, key string, cred *types.RoleCredentials, expiration time.Duration) error {
	if cred.Expiration.Round(0).Add(-expiryDelta).Before(time.Now()) {
		return errors.New("token expired")
	}
	err := r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: cred,
		TTL:   expiration,
	})
	if err != nil {
		return errors.Errorf("cache set error:%v", err)
	}
	return nil
}

func (r *RedisCache) GetRoleCredentials(ctx context.Context, key string) (*types.RoleCredentials, bool, error) {
	var t *types.RoleCredentials
	err := r.cache.Get(ctx, key, &t)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return t, true, nil
}

func (r *RedisCache) SetAuthLog(ctx context.Context, key string, log *AuthLog, expiration time.Duration) error {
	err := r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: log,
		TTL:   expiration,
	})
	if err != nil {
		return errors.Errorf("cache set error:%v", err)
	}
	return nil
}

func (r *RedisCache) GetAuthLog(ctx context.Context, key string) (*AuthLog, bool, error) {
	var t *AuthLog
	err := r.cache.Get(ctx, key, &t)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return t, true, nil
}
