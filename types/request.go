package types

type RequestMeta struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Scope        Scope  `json:"scope,omitempty"`
	Code         string `json:"code,omitempty"`
}
