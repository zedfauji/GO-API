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

import _ "net/http/pprof"

var mySigningKey = []byte("captainjacksparrowsayshi")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Secret struct {
	ID     string `json:"id,omitempty"`
	Token  string `json:"hostname,omitempty"`
	Active bool
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
	secrets.Token = params["token"]
	deviceSecret = append(deviceSecret, secrets)
	json.NewEncoder(w).Encode(deviceSecret)
	fileName := fmt.Sprintf("%s_token", secrets.ID)
	fileLoc := filepath.Join("/tmp/data/secrets/", fileName)
	err := writeGob(fileLoc, deviceSecret)
	fmt.Println("I am over here")
	fmt.Fprintln(w, err)
	fmt.Println("Over here", err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		//http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Everything is Okay")
	}
	//m := map[string]string{
	//	"foo": "bar",
	//}
	//w.Header().Add("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	//_ = json.NewEncoder(w).Encode(m)

	//w.Write([]byte(secrets.ID))
}

//TODO: Improve Dumping into file process.
func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	//err = os.Chown(filePath, 999, 999)
	//os.Chmod(filePath, 0600)
	file.Close()
	return err
}

// our main function
// fun main()

func main() {
	router := mux.NewRouter()
	deviceSecret = append(deviceSecret, Secret{ID: "1", Token: "Device1", Active: true})
	deviceSecret = append(deviceSecret, Secret{ID: "2", Token: "Device2", Active: false})
	deviceSecret = append(deviceSecret, Secret{ID: "3", Token: "Device3-3", Active: true})
	//router.HandleFunc("/deviceSecret", GetIDs).Methods("GET")
	router.Handle("/deviceSecret", isAuthorized(GetIDs))
	router.HandleFunc("/deviceSecret/{id}", GetID).Methods("GET")
	//router.HandleFunc("/deviceSecret/{id}/{hostname}", CreateID).Methods("POST")
	router.Handle("/deviceSecret/{token}", isAuthorized(CreateID))
	//router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}
