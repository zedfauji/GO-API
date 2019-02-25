package generatetoken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"../../Utilities"
	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	url := "https://dev-z1-kqnnv.auth0.com/oauth/token"

	//TODO: Read credentials from Environment or any configuration files

	payload := strings.NewReader("{\"grant_type\":\"client_credentials\",\"client_id\": \"1lBr0bF30njM3qHTzHGsaYc5Z4RZaEL8\",\"client_secret\": \"5dUpVPFu6sof7u4aDjHHR59dzRadR1k1zh6q7x3dJuCQTIzhX9TDWGIlbpY76-tb\",\"audience\": \"zedlocal\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	s := new(Utility.ApiResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Whoops", err)
	}
	target := "/storeid/" + s.AccessToken
	http.Redirect(w, r, target, http.StatusSeeOther)

	//Improve Write to file Process
}

var mySigningKey = []byte("captainjacksparrowsayshi")

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
