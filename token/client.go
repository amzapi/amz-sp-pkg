package token

import (
	"context"
	"crypto/md5"
	"fmt"
	"golang.org/x/sync/singleflight"

	"github.com/amzapi/amz-sp-pkg/cache"
	"github.com/amzapi/amz-sp-pkg/types"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Client struct {
	tokenUrl     string             //
	clientId     string             // SP-API LWA Client ID
	clientSecret string             // SP-API LWA Client Secret
	debug        bool               //
	cache        cache.Cache        //
	sf           singleflight.Group //
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		tokenUrl: "https://api.amazon.com/auth/o2/token",
		cache:    cache.NewMemoryCache(),
		debug:    true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) accessCacheKey(refreshToken string) string {
	h := md5.New()
	_, _ = h.Write([]byte(refreshToken))
	return fmt.Sprintf("amazon:token:%x", h.Sum(nil))
}

func (c *Client) grantLessCacheKey(scope types.Scope) string {
	return fmt.Sprintf("amazon:%s", scope)
}

// Exchange an authorization code received from a getAuthorizationCode operation for a refresh token
func (c *Client) Exchange(ctx context.Context, authCode string) (*types.Token, error) {

	if authCode == "" {
		return nil, errors.New("[token] please provide `authCode` argument")
	}

	requestMeta := types.RequestMeta{
		ClientID:     c.clientId,
		ClientSecret: c.clientSecret,
		GrantType:    "authorization_code",
		Code:         authCode,
	}

	token, err := c.retrieveToken(ctx, requestMeta)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		err := c.cache.SetToken(ctx, c.accessCacheKey(token.RefreshToken), token, token.ExpiryDuration())
		if err != nil {
			return nil, errors.Errorf("[token] set token cache error: %v", err)
		}
	}

	return token, nil
}

// GetAccessToken ...
func (c *Client) GetAccessToken(ctx context.Context, refreshToken string) (*types.Token, error) {

	if refreshToken == "" {
		return nil, errors.New("[token] please provide `refreshToken` argument")
	}

	if c.cache != nil {
		v, found, err := c.cache.GetToken(ctx, c.accessCacheKey(refreshToken))
		if found {
			return v, nil
		} else if err != nil {
			return nil, errors.Errorf("[token] get token cache error: %v", err)
		}
	}

	v, err, _ := c.sf.Do(refreshToken, func() (interface{}, error) {
		return c.GetAccessTokenSkipCache(ctx, refreshToken)
	})

	if err != nil {
		return nil, err
	}

	return v.(*types.Token), nil
}

// GetAccessTokenSkipCache ...
func (c *Client) GetAccessTokenSkipCache(ctx context.Context, refreshToken string) (*types.Token, error) {

	if refreshToken == "" {
		return nil, errors.New("[token] please provide `refreshToken` argument")
	}

	requestMeta := types.RequestMeta{
		ClientID:     c.clientId,
		ClientSecret: c.clientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	token, err := c.retrieveToken(ctx, requestMeta)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		err := c.cache.SetToken(ctx, c.accessCacheKey(token.RefreshToken), token, token.ExpiryDuration())
		if err != nil {
			return nil, errors.Errorf("[token] set token cache error: %v", err)
		}
	}

	return token, nil
}

// GetGrantLessAccessToken token for a grantless operation is requested
// scope should be one of: ['sellingpartnerapi::notifications', 'sellingpartnerapi::migration']
func (c *Client) GetGrantLessAccessToken(ctx context.Context, scope types.Scope) (*types.Token, error) {

	// get cache token
	if c.cache != nil {
		v, found, err := c.cache.GetToken(ctx, c.grantLessCacheKey(scope))
		if found {
			return v, nil
		} else if err != nil {
			return nil, fmt.Errorf("[token] get token cache error: %v", err)
		}
	}

	v, err, _ := c.sf.Do(string(scope), func() (interface{}, error) {
		return c.GetGrantLessAccessTokenSkipCache(ctx, scope)
	})

	if err != nil {
		return nil, err
	}

	return v.(*types.Token), nil
}

// GetGrantLessAccessTokenSkipCache ...
func (c *Client) GetGrantLessAccessTokenSkipCache(ctx context.Context, scope types.Scope) (*types.Token, error) {

	requestMeta := types.RequestMeta{
		ClientID:     c.clientId,
		ClientSecret: c.clientSecret,
		GrantType:    "client_credentials",
		Scope:        scope,
	}

	token, err := c.retrieveToken(ctx, requestMeta)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		err = c.cache.SetToken(ctx, c.grantLessCacheKey(scope), token, token.ExpiryDuration())
		if err != nil {
			return nil, errors.Errorf("[token] set token cache error: %v", err)
		}
	}

	return token, nil
}

// retrieveToken ...
func (c *Client) retrieveToken(ctx context.Context, requestMeta types.RequestMeta) (*types.Token, error) {

	token := &types.Token{}

	client := resty.New()
	client.SetDebug(c.debug)

	resp, err := client.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetHeader("Content-Type", "application/json").
		SetBody(requestMeta).
		SetResult(&token).
		Post(c.tokenUrl)

	if err != nil {
		return nil, errors.Errorf("[token] retrieve token error: %v", err)
	}

	if resp.IsError() {
		if err, ok := resp.Error().(*Error); ok {
			return nil, err
		}
		return nil, errors.Errorf("[token] retrieve token error: %v", resp.String())
	}

	return token, nil
}
