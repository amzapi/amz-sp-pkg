package types

import "time"

type RoleCredentials struct {
	AccessKeyID     string    //
	SecretAccessKey string    //
	SessionToken    string    //
	Expiration      time.Time //
}

func (r *RoleCredentials) ExpiryDuration() time.Duration {
	return r.Expiration.Sub(time.Now())
}
