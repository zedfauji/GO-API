package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"../provisioning/Storage"
	"../provisioning/Utilities"

	"github.com/gorilla/mux"
)
import _ "net/http/pprof"

var mySigningKey = []byte("captainjacksparrowsayshi")

//type Secret struct {
//	ID       string `json:"id,omitempty"`
//	Hostname string `json:"hostname,omitempty"`
//	Active   bool
//}

//var deviceSecret []Secret
var deviceSecret []Utility.Secret

func getIds(w http.ResponseWriter, r *http.Request) {
	validToken, err := Storage.GenerateJWT()
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
	params := mux.Vars(r)
	//var secrets Utility.Secret
	secrets := new(Utility.ApiResponse)
	_ = json.NewDecoder(r.Body).Decode(&deviceSecret)
	secrets.AccessToken = params["token"]

	response := Storage.StoreToken(*secrets)
	fmt.Println(response)

	if response == "broken" {
		fmt.Println("Its broken")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(secrets.AccessToken))
	}
}

//func iSActive(s Utility.Secret) string,error {
//TODO: retrieve Accesstokena and its expiry time based on ID

//}
func generateToken(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	response := "broken"
	url := "https://dev-z1-kqnnv.auth0.com/oauth/token"

	//TODO: Read credentials from Environment or any configuration files

	payload := strings.NewReader("{\"grant_type\":\"client_credentials\",\"client_id\": \"1lBr0bF30njM3qHTzHGsaYc5Z4RZaEL8\",\"client_secret\": \"5dUpVPFu6sof7u4aDjHHR59dzRadR1k1zh6q7x3dJuCQTIzhX9TDWGIlbpY76-tb\",\"audience\": \"zedlocal\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	//var s []Utility.ApiResponse
	s := new(Utility.ApiResponse)
	//s := new(Utility.apiResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Whoops", err)
	}
	response = Storage.StoreToken(*s)
	fmt.Println(response)

	//Improve Write to file Process
	if response == "broken" {
		fmt.Println("Its broken")
	} else {
		w.WriteHeader(http.StatusOK)
		//TODO: Return ID with Accesstoken
		w.Write([]byte(s.AccessToken))
	}

}

func main() {
	router := mux.NewRouter()
	go func() {
		log.Fatal(http.ListenAndServe(":6060", http.DefaultServeMux))
	}()
	//http.HandleFunc("/", getIds)
	router.HandleFunc("/storeId/{token}}", storeId)
	router.HandleFunc("/generateToken", generateToken)
	router.HandleFunc("/getids", getIds)

	log.Fatal(http.ListenAndServe(":9001", router))
}
