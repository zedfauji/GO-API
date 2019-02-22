package Storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"../../provisioning/Utilities"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

func StoreToken(s Utility.ApiResponse) string {
	var response = "success"
	validToken, err := GenerateJWT()

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
	fmt.Println(err)
	fmt.Println(res.Header)
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		response = "broken"
		fmt.Println(string(b))
	}
	//if res.StatusCode > 200 {
	//		response = "broken"
	//}
	if err != nil {
		fmt.Println("Error:", err)
	}
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(body))
	return response
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Girish Dudhwal"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
