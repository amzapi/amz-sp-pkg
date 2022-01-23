package rdt

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/amzapi/amz-sp-pkg/cache"
	"github.com/amzapi/amz-sp-pkg/signer"
	"github.com/amzapi/amz-sp-pkg/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Client struct {
	debug     bool           //
	cache     cache.Cache    //
	signer    *signer.Signer //
	userAgent string         //
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		debug: true,
		cache: cache.NewMemoryCache(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) GetRestrictedDataTokenSkipCache(ctx context.Context, endpoint string, awsRegion string, accessToken string, operations []*RestrictedOperation) (*RestrictedDataTokenResponse, error) {

	client := resty.New()
	client.SetDebug(true)
	client.SetError(&ErrorResponse{})
	client.SetHeader("User-Agent", c.userAgent)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
		err := c.signer.SignRequest(ctx, request, accessToken, awsRegion)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				return errors.New(http.StatusBadGateway, awsErr.Code(), awsErr.Message())
			} else {
				return errors.New(http.StatusBadGateway, "SignRequest", err.Error())
			}
		}
		return nil
	})

	url := fmt.Sprintf("%s%s", endpoint, "/tokens/2021-03-01/restrictedDataToken")
	var result RestrictedDataTokenResponse
	resp, err := client.R().SetContext(ctx).SetResult(&result).SetBody(&CreateRestrictedDataTokenRequest{RestrictedResources: operations}).Execute(http.MethodPost, url)
	if err != nil {
		return nil, errors.New(http.StatusBadGateway, "GetRestrictedDataToken", err.Error())
	}

	if resp.IsError() {
		if errorResponse, ok := resp.Error().(*ErrorResponse); ok {
			if len(errorResponse.Errors) > 0 {
				if errorResponse.Errors[0].Details == "" {
					return nil, errors.New(resp.StatusCode(), errorResponse.Errors[0].Code, errorResponse.Errors[0].Message)
				}
				return nil, errors.Newf(resp.StatusCode(), errorResponse.Errors[0].Code, "Message = %s , Details = %s", errorResponse.Errors[0].Message, errorResponse.Errors[0].Details)
			}
		}
		return nil, errors.New(http.StatusBadGateway, "GetRestrictedDataToken", resp.String())
	}

	return &result, nil
}

func (c *Client) GetRestrictedDataToken(ctx context.Context, sellerId string, endpoint string, awsRegion string, accessToken string, operations []*RestrictedOperation, isUniversal bool) (*string, error) {

	var cacheKey string

	if c.cache != nil {
		if isUniversal {
			cacheKey = fmt.Sprintf("amazon:rdt:%s", sellerId)
		} else if len(operations) > 0 {
			cacheKey = fmt.Sprintf("amazon:rdt:%s:%s", sellerId, operations[0].Path)
		}
		if cacheKey != "" {
			v, found, err := c.cache.GetToken(ctx, cacheKey)
			if found {
				return &v.AccessToken, nil
			} else if err != nil {
				return nil, errors.Errorf(http.StatusBadGateway, "GetRestrictedDataToken", "get token cache error: %v", err)
			}
		}
	}

	var body []*RestrictedOperation
	for _, operation := range operations {
		if operation.Enable && isUniversal == operation.IsUniversal {
			body = append(body, operation)
		}
	}

	result, err := c.GetRestrictedDataTokenSkipCache(ctx, endpoint, awsRegion, accessToken, body)
	if err != nil {
		return nil, err
	}

	if c.cache != nil && cacheKey != "" {
		token := types.Token{
			AccessToken: result.RestrictedDataToken,
			ExpiresIn:   result.ExpiresIn,
			Expiry:      result.Expiry,
		}
		err := c.cache.SetToken(ctx, cacheKey, &token, token.ExpiryDuration())
		if err != nil {
			return nil, errors.Errorf(http.StatusBadGateway, "GetRestrictedDataToken", "set token cache error: %v", err)
		}
	}

	return &result.RestrictedDataToken, nil
}
