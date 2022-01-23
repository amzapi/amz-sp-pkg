package rdt

import (
	"github.com/amzapi/amz-sp-pkg/cache"
	"github.com/amzapi/amz-sp-pkg/signer"
)

// Option is config option.
type Option func(client *Client)

func WithCache(v cache.Cache) Option {
	return func(s *Client) {
		s.cache = v
	}
}

func WithDebug(v bool) Option {
	return func(s *Client) {
		s.debug = v
	}
}

func WithUserAgent(v string) Option {
	return func(s *Client) {
		s.userAgent = v
	}
}

func WithSigner(v *signer.Signer) Option {
	return func(s *Client) {
		s.signer = v
	}
}
