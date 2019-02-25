package Utility

type Secret struct {
	ID          string `json:"id,omitempty"`
	AccessToken string `json:"token,omitempty"`
	Active      bool
}

type ApiResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresTime int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type StoreResponse struct {
	WriteStatus string `json:"status"`
	GeneratedID string `json:"id"`
}
