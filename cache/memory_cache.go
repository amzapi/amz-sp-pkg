package cache

import (
	"context"
	"time"

	"gopkg.me/amz-sp-pkg/types"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type MemoryCache struct {
	cache *cache.Cache
}

// NewMemoryCache cache object in go-cache,It's done in memory
func NewMemoryCache() Cache {
	return &MemoryCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (m *MemoryCache) SetToken(ctx context.Context, key string, token *types.Token, expiration time.Duration) error {
	if token.Expiry.Round(0).Add(-expiryDelta).Before(time.Now()) {
		return errors.New("token expired")
	}
	m.cache.Set(key, token, expiration-expiryDelta)
	return nil
}

func (m *MemoryCache) GetToken(ctx context.Context, key string) (*types.Token, bool, error) {
	if x, found := m.cache.Get(key); found {
		return x.(*types.Token), true, nil
	}
	return nil, false, nil
}

func (m *MemoryCache) SetRoleCredentials(ctx context.Context, key string, cred *types.RoleCredentials, expiration time.Duration) error {
	if cred.Expiration.Round(0).Add(-expiryDelta).Before(time.Now()) {
		return errors.New("token expired")
	}
	m.cache.Set(key, cred, expiration-expiryDelta)
	return nil
}

func (m *MemoryCache) GetRoleCredentials(ctx context.Context, key string) (*types.RoleCredentials, bool, error) {
	if x, found := m.cache.Get(key); found {
		return x.(*types.RoleCredentials), true, nil
	}
	return nil, false, nil
}

func (m *MemoryCache) SetAuthLog(ctx context.Context, key string, log *AuthLog, expiration time.Duration) error {
	m.cache.Set(key, log, expiration)
	return nil
}

func (m *MemoryCache) GetAuthLog(ctx context.Context, key string) (*AuthLog, bool, error) {
	if x, found := m.cache.Get(key); found {
		return x.(*AuthLog), true, nil
	}
	return nil, false, nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.cache.Delete(key)
	return nil
}
