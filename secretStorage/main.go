package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/xid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Secret struct {
	ID       string `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Active   bool
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

var deviceSecret []Secret

func GetIDs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(deviceSecret)
}
func GetID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range deviceSecret {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Secret{})
}

//Create new ID

func CreateID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var secrets Secret
	_ = json.NewDecoder(r.Body).Decode(&deviceSecret)
	secrets.ID = params["id"]
	guid := xid.New()
	secrets.ID = guid.String()
	secrets.Hostname = params["hostname"]
	deviceSecret = append(deviceSecret, secrets)
	json.NewEncoder(w).Encode(deviceSecret)
	fileName := fmt.Sprintf("%s_token", secrets.ID)
	fileLoc := filepath.Join("/data/secrets/", fileName)
	err := writeGob(fileLoc, deviceSecret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	err1 := os.Chown(fileLoc, 999, 999)
	if err1 != nil {
		fmt.Println("Error in changing file permission", err1)
	}
}

// our main function
// fun main()

func main() {
	router := mux.NewRouter()
	deviceSecret = append(deviceSecret, Secret{ID: "1", Hostname: "Device1", Active: true})
	deviceSecret = append(deviceSecret, Secret{ID: "2", Hostname: "Device2", Active: false})
	deviceSecret = append(deviceSecret, Secret{ID: "3", Hostname: "Device3-3", Active: true})
	//router.HandleFunc("/deviceSecret", GetIDs).Methods("GET")
	router.Handle("/deviceSecret", isAuthorized(GetIDs))
	router.HandleFunc("/deviceSecret/{id}", GetID).Methods("GET")
	//router.HandleFunc("/deviceSecret/{id}/{hostname}", CreateID).Methods("POST")
	router.Handle("/deviceSecret/{hostname}", isAuthorized(CreateID))
	//router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}
