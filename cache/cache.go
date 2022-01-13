package cache

import (
	"context"
	"time"

	"gopkg.me/amz-sp-pkg/types"
)

const (
	// expiryDelta determines how earlier a token should be considered
	// expired than its actual expiration time. It is used to avoid late
	// expirations due to client-server time mismatches.
	expiryDelta = 60 * time.Second

	// cleanupInterval cleanup interval
	cleanupInterval = 5 * time.Second

	// defaultExpiration default expiration duration
	defaultExpiration = 3600 * time.Second
)

type AuthLog struct {
	AwsRegion        string    // us-east-1 eu-west-1 us-west-2
	Endpoint         string    //
	State            string    // 唯一标识
	Region           string    // 区域代码
	CountryCode      []string  // 授权的国家代码
	CallBackUrl      string    // 回调URL
	SellingPartnerID string    // 卖家ID
	AccessToken      string    //
	TokenType        string    //
	RefreshToken     string    //
	ExpiresTime      time.Time //
}

type Cache interface {
	Delete(ctx context.Context, key string) error
	SetToken(ctx context.Context, key string, token *types.Token, expiration time.Duration) error
	GetToken(ctx context.Context, key string) (*types.Token, bool, error)
	SetRoleCredentials(ctx context.Context, key string, cred *types.RoleCredentials, expiration time.Duration) error
	GetRoleCredentials(ctx context.Context, key string) (*types.RoleCredentials, bool, error)
	SetAuthLog(ctx context.Context, key string, log *AuthLog, expiration time.Duration) error
	GetAuthLog(ctx context.Context, key string) (*AuthLog, bool, error)
}
