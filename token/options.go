package token

import (
	"gopkg.me/amz-sp-pkg/cache"
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

func WithTokenUrl(v string) Option {
	return func(s *Client) {
		s.tokenUrl = v
	}
}

func WithClientID(v string) Option {
	return func(s *Client) {
		s.clientId = v
	}
}

func WithClientSecret(v string) Option {
	return func(s *Client) {
		s.clientSecret = v
	}
}
