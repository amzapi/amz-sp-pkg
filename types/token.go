package types

import (
	"encoding/json"
	"time"
)

// Token is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int32     `json:"expires_in"`
	Expiry       time.Time `json:"expiry,omitempty"`
	Scope        Scope     `json:"scope,omitempty"`
}

func (t *Token) UnmarshalJSON(data []byte) error {
	type xToken Token
	x := &xToken{}
	if err := json.Unmarshal(data, x); err != nil {
		return err
	}
	x.Expiry = time.Now().Add(time.Duration(x.ExpiresIn) * time.Second)
	*t = Token(*x)
	return nil
}

func (t *Token) ExpiryDuration() time.Duration {
	return t.Expiry.Sub(time.Now())
}
