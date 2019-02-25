package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rs/xid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	_ "net/http/pprof"

	"./Routes/Read"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Secret struct {
	ID     string `json:"id,omitempty"`
	Token  string `json:"token,omitempty"`
	Active bool 
}

type ResponseSecret struct {
	ID     string
	Token  string
	Active string
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
	secrets.ID = xid.New().String()
	secrets.Token = params["token"]
	deviceSecret = append(deviceSecret, secrets)
	fileName := fmt.Sprintf("%s.json", secrets.ID)
	fileLoc := filepath.Join("/tmp/data/secrets/", fileName)
	err := writeGob(fileLoc, deviceSecret)
	m := map[string]string{}
	if err != nil {
		fmt.Println(err)
		m["status"] = "broken"
		//check(err)
	} else {
		m["status"] = "success"
		m["id"] = secrets.ID
	}
	_ = json.NewEncoder(w).Encode(m)
}

//TODO: Improve Dumping into file process.
func writeGob(filePath string, object interface{}) error {
	//file, err := os.Create(filePath)
	file, _ := json.MarshalIndent(object, "", "")
	_ = ioutil.WriteFile(filePath, file, 0600)

	//err = os.Chown(filePath, 999, 999)
	//os.Chmod(filePath, 0600)
	return nil
}

// our main function
// fun main()

func main() {
	router := mux.NewRouter()
	//deviceSecret = append(deviceSecret, Secret{ID: "1", Token: "Device1", Active: true})
	//deviceSecret = append(deviceSecret, Secret{ID: "2", Token: "Device2", Active: false})
	//deviceSecret = append(deviceSecret, Secret{ID: "3", Token: "Device3-3", Active: true})
	router.Handle("/deviceSecret", isAuthorized(GetIDs))
	router.HandleFunc("/deviceSecret/{id}", readoperation.ReadID).Methods("GET")
	router.Handle("/deviceSecret/{token}", isAuthorized(CreateID))
	log.Fatal(http.ListenAndServe(":8000", router))

}
