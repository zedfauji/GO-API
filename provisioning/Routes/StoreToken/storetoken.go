package StoreToken

import (
	"encoding/json"
	"net/http"

	"../../Storage"
	"../../Utilities"
	"github.com/gorilla/mux"
)

var deviceSecret []Utility.Secret

func StoreId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//var secrets Utility.Secret
	secrets := new(Utility.ApiResponse)
	_ = json.NewDecoder(r.Body).Decode(&deviceSecret)
	secrets.AccessToken = params["token"]

	response := Storage.StoreToken(*secrets)
	m := map[string]string{}
	if response == "failed" {
		m["access_toke"] = "None"
		m["id"] = "None"
	} else {
		w.WriteHeader(http.StatusOK)
		m["access_token"] = secrets.AccessToken
		m["id"] = response

	}
	_ = json.NewEncoder(w).Encode(m)
}
