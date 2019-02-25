package Storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"../../provisioning/Routes/GenerateToken"
	"../../provisioning/Utilities"
)

func StoreToken(s Utility.ApiResponse) string {
	var response = "failed"
	validToken, err := generatetoken.GenerateJWT()

	var secrets Utility.Secret
	secrets.ID = "randonID"
	secrets.AccessToken = s.AccessToken
	apiUrl := "http://localhost:8000"
	urlPath := fmt.Sprintf("/deviceSecret/%s", secrets.AccessToken)
	resources := urlPath
	data := url.Values{}
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resources
	urlStr := u.String()
	client := &http.Client{}
	//req, _ := http.NewRequest("POST", "http://localhost:8000/deviceSecret/", nil)
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	//req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	var CurrentResponse Utility.StoreResponse
	_ = json.NewDecoder(res.Body).Decode(&CurrentResponse)
	if CurrentResponse.WriteStatus == "success" {
		response = CurrentResponse.GeneratedID
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
	return response
}
