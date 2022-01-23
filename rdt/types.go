package rdt

import (
	"encoding/json"
	"regexp"
	"time"
)

type RestrictedOperation struct {
	Name         string         `json:"-"`
	Method       string         `json:"method"`
	Path         string         `json:"path"`
	PathRegexp   *regexp.Regexp `json:"-"`
	DataElements []string       `json:"dataElements,omitempty"`
	IsUniversal  bool           `json:"-"`
	Enable       bool           `json:"-"`
}

type CreateRestrictedDataTokenRequest struct {
	RestrictedResources []*RestrictedOperation `json:"restrictedResources"`
	TargetApplication   *string                `json:"targetApplication,omitempty"`
}

type RestrictedDataTokenResponse struct {
	ExpiresIn           int32     `json:"expiresIn,omitempty"`
	RestrictedDataToken string    `json:"restrictedDataToken,omitempty"`
	Expiry              time.Time `json:"expiry,omitempty"`
}

func (t *RestrictedDataTokenResponse) UnmarshalJSON(data []byte) error {
	type xToken RestrictedDataTokenResponse
	x := &xToken{}
	if err := json.Unmarshal(data, x); err != nil {
		return err
	}
	x.Expiry = time.Now().Add(time.Duration(x.ExpiresIn) * time.Second)
	*t = RestrictedDataTokenResponse(*x)
	return nil
}

func (t *RestrictedDataTokenResponse) ExpiryDuration() time.Duration {
	return t.Expiry.Sub(time.Now())
}

func MatchOperation(operations []*RestrictedOperation, method, path string) *RestrictedOperation {
	for _, operation := range operations {
		if operation.Enable && operation.Method == method && operation.PathRegexp.MatchString(path) {
			return operation
		}
	}
	return nil
}

func MatchReportType(reportTypes []string, reportType string) bool {
	for _, s := range reportTypes {
		if s == reportType {
			return true
		}
	}
	return false
}
