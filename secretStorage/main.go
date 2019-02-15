package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

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
	secrets.Hostname = params["hostname"]
	deviceSecret = append(deviceSecret, secrets)
	json.NewEncoder(w).Encode(deviceSecret)
	err := writeGob("/data/secrets/sec12", deviceSecret)
	if err != nil {
		fmt.Println(err)
	}
}

// our main function
// fun main()

func main() {
	router := mux.NewRouter()
	deviceSecret = append(deviceSecret, Secret{ID: "1", Hostname: "Device1", Active: true})
	deviceSecret = append(deviceSecret, Secret{ID: "2", Hostname: "Device2", Active: false})
	deviceSecret = append(deviceSecret, Secret{ID: "3", Hostname: "Device3-3", Active: true})
	router.HandleFunc("/deviceSecret", GetIDs).Methods("GET")
	router.HandleFunc("/deviceSecret/{id}", GetID).Methods("GET")
	router.HandleFunc("/deviceSecret/{id}/{hostname}", CreateID).Methods("POST")
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
