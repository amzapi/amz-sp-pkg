package signer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/amzapi/amz-sp-pkg/cache"
	"github.com/amzapi/amz-sp-pkg/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

type Signer struct {
	cache           cache.Cache //
	cacheKey        string      //
	debug           bool        //
	accessKeyID     string      //AWS IAM User Access Key ID
	secretAccessKey string      //AWS IAM User Secret Key
	roleArn         string      //AWS IAM Role ARN
}

func NewSigner(opts ...Option) *Signer {
	s := &Signer{
		cache: cache.NewMemoryCache(),
		debug: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.cacheKey = fmt.Sprintf("amazon:sts:%s", s.accessKeyID)
	return s
}

func (s *Signer) RefreshRoleCredentialsSkipCache(ctx context.Context) (*types.RoleCredentials, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			s.accessKeyID,
			s.secretAccessKey,
			"",
		),
	}))

	svc := sts.New(sess)

	if s.debug {
		svc.Config.LogLevel = aws.LogLevel(aws.LogDebug)
	}

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(s.roleArn),
		RoleSessionName: aws.String("SPAPISession"),
	}

	result, err := svc.AssumeRoleWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	roleCredentials := &types.RoleCredentials{
		AccessKeyID:     aws.StringValue(result.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(result.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(result.Credentials.SessionToken),
		Expiration:      aws.TimeValue(result.Credentials.Expiration),
	}

	if s.cache != nil {
		err = s.cache.SetRoleCredentials(ctx, s.cacheKey, roleCredentials, roleCredentials.ExpiryDuration())
		if err != nil {
			return nil, errors.WithMessage(err, "[signer] set role credentials cache error")
		}
	}

	return roleCredentials, nil
}

func (s *Signer) RefreshRoleCredentials(ctx context.Context) (*types.RoleCredentials, error) {
	if s.cache != nil {
		v, found, err := s.cache.GetRoleCredentials(ctx, s.cacheKey)
		if found {
			return v, nil
		} else if err != nil {
			return nil, errors.WithMessage(err, "[signer] get role credentials cache error")
		}
	}
	return s.RefreshRoleCredentialsSkipCache(ctx)
}

func (s *Signer) SignRequest(ctx context.Context, r *http.Request, accessToken, awsRegion string) error {

	roleCredentials, err := s.RefreshRoleCredentials(ctx)
	if err != nil {
		return err
	}

	var body io.ReadSeeker
	if r.Body != nil {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
		r.Body = ioutil.NopCloser(body)
	}

	r.Header.Del("X-Amz-Access-Token")
	r.Header.Del("X-Amz-Date")
	r.Header.Del("X-Amz-Security-Token")
	r.Header.Add("X-Amz-Access-Token", accessToken)

	aws4Signer := v4.NewSigner(credentials.NewStaticCredentials(
		roleCredentials.AccessKeyID,
		roleCredentials.SecretAccessKey,
		roleCredentials.SessionToken),
		func(s *v4.Signer) {
			s.DisableURIPathEscaping = true
		},
	)

	_, err = aws4Signer.Sign(r, body, "execute-api", awsRegion, time.Now())

	return err
}
