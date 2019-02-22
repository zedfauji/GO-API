package Utility

type Secret struct {
	ID          string `json:"id,omitempty"`
	AccessToken string `json:"hostname,omitempty"`
	Active      bool
}

type ApiResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresTime int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
