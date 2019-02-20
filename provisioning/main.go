package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type Secret struct {
	ID       string `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Active   bool
}

var deviceSecret []Secret

func getIds(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8000/deviceSecret", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(body))
}
func storeId(w http.ResponseWriter, r *http.Request) {
	//validToken, err := GenerateJWT()
	//if err != nil {
	//	fmt.Println("Failed to generate token")
	//}
	params := mux.Vars(r)
	var secrets Secret
	_ = json.NewDecoder(r.Body).Decode(&deviceSecret)
	secrets.ID = params["id"]
	fmt.Println(secrets.ID)

	secrets.Hostname = params["hostname"]
	fmt.Println(secrets.Hostname)
	apiUrl := "http://localhost:8000"
	urlPath := fmt.Sprintf("/deviceSecret/%s", secrets.Hostname)
	resources := urlPath
	data := url.Values{}
	data.Set("id", secrets.ID)
	data.Set("hostname", secrets.Hostname)
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resources
	urlStr := u.String()
	client := &http.Client{}
	//req, _ := http.NewRequest("POST", "http://localhost:8000/deviceSecret/", nil)
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	//req.Header.Set("Token", validToken)
	fmt.Println(urlStr)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(body))
}

type apiResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresTime int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	response := "broken"
	url := "https://dev-z1-kqnnv.auth0.com/oauth/token"

	payload := strings.NewReader("{\"grant_type\":\"client_credentials\",\"client_id\": \"1lBr0bF30njM3qHTzHGsaYc5Z4RZaEL8\",\"client_secret\": \"5dUpVPFu6sof7u4aDjHHR59dzRadR1k1zh6q7x3dJuCQTIzhX9TDWGIlbpY76-tb\",\"audience\": \"zedlocal\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	s := new(apiResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Whoops", err)
	}
	fmt.Println("Access Token", s.AccessToken)

	response = storeToken(*s)
	fmt.Println(response)

	if response == "broken" {
		fmt.Println("Its broken")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s.AccessToken))
	}

}

func storeToken(s apiResponse) string {
	fmt.Println("In StoreToken")
	var response = "success"
	validToken, err := GenerateJWT()

	var secrets Secret
	secrets.ID = "randonID"
	fmt.Println(secrets.ID)
	secrets.Hostname = s.AccessToken
	fmt.Println(secrets.Hostname)
	fmt.Println(secrets.Hostname)
	apiUrl := "http://localhost:8000"
	urlPath := fmt.Sprintf("/deviceSecret/%s", secrets.Hostname)
	resources := urlPath
	data := url.Values{}
	data.Set("id", secrets.ID)
	data.Set("hostname", secrets.Hostname)
	//fmt.Println(secrets.ID, secrets.Hostname)
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resources
	urlStr := u.String()
	client := &http.Client{}
	//req, _ := http.NewRequest("POST", "http://localhost:8000/deviceSecret/", nil)
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	fmt.Println(res.StatusCode)
	if res.StatusCode > 200 {
		response = "broken"
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
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

func handleRequests() {
	//http.HandleFunc("/", getIds)
	http.HandleFunc("/storeId/{id}/{hostname}", storeId)
	http.HandleFunc("/generateToken", generateToken)
	http.HandleFunc("/getids", getIds)

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	handleRequests()
}
