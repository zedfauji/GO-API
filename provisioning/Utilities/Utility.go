package Utility

import (
	"fmt"
	"net/url"
)

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

func MakeURL(s string) string {
	apiUrl := "http://localhost:8000"
	urlPath := fmt.Sprintf("/deviceSecret/%s", s)
	resources := urlPath
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resources
	urlStr := u.String()
	return urlStr
}
