package signer

import (
	"gopkg.me/amz-sp-pkg/cache"
)

// Option is config option.
type Option func(*Signer)

func WithCache(v cache.Cache) Option {
	return func(s *Signer) {
		s.cache = v
	}
}

func WithDebug(v bool) Option {
	return func(s *Signer) {
		s.debug = v
	}
}

func WithAccessKeyID(v string) Option {
	return func(s *Signer) {
		s.accessKeyID = v
	}
}

func WithSecretAccessKey(v string) Option {
	return func(s *Signer) {
		s.secretAccessKey = v
	}
}

func WithRoleArn(v string) Option {
	return func(s *Signer) {
		s.roleArn = v
	}
}
